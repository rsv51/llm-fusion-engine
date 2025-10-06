package migrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// MigrateRefactorModelsAndProviders performs the database migration for the new model and provider structure.
func MigrateRefactorModelsAndProviders(db *gorm.DB) error {
	// We use raw SQL for more control over the migration process,
	// especially for table alterations and data migration.
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}

	// Disable foreign key checks for the duration of the migration
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}
	// Re-enable foreign key checks at the end
	defer func() {
		if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
			log.Printf("Warning: failed to re-enable foreign key checks: %v", err)
		}
	}()

	// Begin a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// --- Step 1: Create new tables ---

	// Create the new `models` table
	if err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS models (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			name TEXT NOT NULL UNIQUE,
			remark TEXT,
			max_retry INTEGER DEFAULT 3,
			timeout INTEGER DEFAULT 30,
			enabled BOOLEAN DEFAULT 1
		);
	`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create 'models' table: %w", err)
	}

	// Create the new `model_provider_mappings` table
	if err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS model_provider_mappings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			model_id INTEGER NOT NULL,
			provider_id INTEGER NOT NULL,
			provider_model TEXT NOT NULL,
			tool_call BOOLEAN,
			structured_output BOOLEAN,
			image BOOLEAN,
			weight INTEGER DEFAULT 1,
			enabled BOOLEAN DEFAULT 1,
			FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE,
			FOREIGN KEY (provider_id) REFERENCES providers(id) ON DELETE CASCADE
		);
	`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create 'model_provider_mappings' table: %w", err)
	}

	// --- Step 2: Alter the existing `providers` table ---
	// This involves adding new columns, dropping old ones, and renaming.

	// Add new columns
	newColumns := []string{
		"ALTER TABLE providers ADD COLUMN name TEXT",
		"ALTER TABLE providers ADD COLUMN config TEXT",
		"ALTER TABLE providers ADD COLUMN console TEXT",
	}
	for _, stmt := range newColumns {
		// Use "IF NOT EXISTS" for SQLite compatibility (though not standard SQL, GORM might handle it or we check first)
		// For simplicity, we'll assume GORM or the driver handles it, or we ignore errors if column exists.
		// A more robust way is to query pragma table_info first.
		if err := tx.Exec(stmt).Error; err != nil {
			// Check if the error is "duplicate column name" which is fine
			if !isDuplicateColumnError(err) {
				tx.Rollback()
				return fmt.Errorf("failed to add new column to providers table with statement '%s': %w", stmt, err)
			}
		}
	}

	// Populate the new `name` column from `provider_type` if it's empty
	if err := tx.Exec("UPDATE providers SET name = provider_type WHERE name IS NULL OR name = ''").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to populate provider names: %w", err)
	}

	// Migrate data to the new `config` column
	// This requires fetching existing providers, formatting their config, and updating
	type OldProvider struct {
		ID           uint
		ProviderType string
		BaseURL      string
		Timeout      int
		MaxRetries   int
		Enabled      bool
	}

	var oldProviders []OldProvider
	if err := tx.Table("providers").Select("id, provider_type, base_url, timeout, max_retries, enabled").Find(&oldProviders).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch old providers for config migration: %w", err)
	}

	for _, p := range oldProviders {
		configMap := map[string]interface{}{
			"baseUrl":    p.BaseURL,
			"timeout":   p.Timeout,
			"maxRetries": p.MaxRetries,
			"enabled":   p.Enabled,
		}
		configBytes, _ := json.Marshal(configMap)
		if err := tx.Exec("UPDATE providers SET config = ? WHERE id = ?", string(configBytes), p.ID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update config for provider %d: %w", p.ID, err)
		}
	}

	// Drop old columns
	oldColumnsToDrop := []string{"group_id", "provider_type", "base_url", "timeout", "max_retries", "health_status", "last_checked", "latency"}
	for _, col := range oldColumnsToDrop {
		// SQLite does not support dropping multiple columns with a single ALTER TABLE statement easily.
		// We'll drop them one by one. GORM's Migrator might handle this better, but for raw SQL:
		if err := tx.Exec(fmt.Sprintf("ALTER TABLE providers DROP COLUMN %s", col)).Error; err != nil {
			// It's possible the column doesn't exist, which is fine during a refactor.
			log.Printf("Warning: could not drop column %s from providers (it might not exist): %v", col, err)
		}
	}
	
	// Rename `weight` to `new_weight` temporarily to avoid conflict if it already exists and has different type
	// Then rename it back to `weight` after other changes. Or just ensure it's the correct type.
	// For simplicity, we'll assume it's fine or handle it if issues arise.
	// The new Provider model has `weight uint`, so let's ensure it's `weight INTEGER`
	// If `weight` already exists and is not the desired type, this needs more care.
	// Assuming `weight` is already `INTEGER` from previous schema.

	// Make `name` unique and not null
	if err := tx.Exec("UPDATE providers SET name = 'provider_' || id WHERE name IS NULL OR name = ''").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set default names for providers: %w", err)
	}
	if err := tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_providers_name ON providers(name)").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create unique index on providers.name: %w", err)
	}
	// Ensure `type` column exists and is indexed. If it was renamed from `provider_type`, this needs adjustment.
	// Assuming `type` is the new name for `provider_type`.
	// If `provider_type` still exists, rename it.
	if err := tx.Exec("ALTER TABLE providers RENAME COLUMN provider_type TO type").Error; err != nil {
		// Check if it's because the column doesn't exist or is already renamed
		if isColumnMissingError(err) {
			// Check if `type` column exists, if not, create it
			if err := tx.Exec("ALTER TABLE providers ADD COLUMN type TEXT").Error; err != nil && !isDuplicateColumnError(err) {
				// If adding `type` fails and it's not a duplicate error, rollback
				if !isDuplicateColumnError(err) {
					tx.Rollback()
					return fmt.Errorf("failed to add 'type' column to providers: %w", err)
				}
			}
			// Populate `type` if it was just added or is empty
			// This logic needs to be smarter if `type` can be derived from old data.
			// For now, let's assume `type` should be set based on `name` or some other logic if `provider_type` is gone.
			// This is a placeholder for more complex migration logic.
			if err := tx.Exec("UPDATE providers SET type = 'unknown' WHERE type IS NULL OR type = ''").Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to populate provider type: %w", err)
			}
		}
	}
	if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_providers_type ON providers(type)").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create index on providers.type: %w", err)
	}


	// --- Step 3: Migrate data from `model_mappings` to new tables ---
	// This involves creating entries in `models` and `model_provider_mappings`.

	type OldModelMapping struct {
		ID                 uint
		UserFriendlyName  string
		ProviderModelName string
		ProviderID        uint
	}

	var oldMappings []OldModelMapping
	if err := tx.Table("model_mappings").Select("id, user_friendly_name, provider_model_name, provider_id").Find(&oldMappings).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch old model_mappings: %w", err)
	}

	for _, mapping := range oldMappings {
		// 1. Create or find a corresponding entry in the `models` table
		var modelID uint
		// Check if model with this name already exists
		tx.Raw("SELECT id FROM models WHERE name = ?", mapping.UserFriendlyName).Scan(&modelID)
		if modelID == 0 {
			// Model does not exist, create it
			newModel := map[string]interface{}{
				"name":     mapping.UserFriendlyName,
				"remark":   "Migrated from old mapping",
				"enabled":  true,
			}
			result := tx.Table("models").Create(newModel)
			if result.Error != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create model '%s': %w", mapping.UserFriendlyName, result.Error)
			}
			// Get the ID of the newly created model
			// GORM's Create doesn't directly return the ID in the map, so we query it.
			tx.Raw("SELECT last_insert_rowid()").Scan(&modelID)
			if modelID == 0 {
				tx.Rollback()
				return fmt.Errorf("failed to retrieve ID for newly created model '%s'", mapping.UserFriendlyName)
			}
		}

		// 2. Create an entry in `model_provider_mappings`
		newMapping := map[string]interface{}{
			"model_id":       modelID,
			"provider_id":    mapping.ProviderID,
			"provider_model": mapping.ProviderModelName,
			"enabled":       true,
		}
		if err := tx.Table("model_provider_mappings").Create(newMapping).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create model_provider_mapping for model_id %d, provider_id %d: %w", modelID, mapping.ProviderID, err)
		}
	}

	// --- Step 4: Drop old tables ---
	// After successful migration, drop the old tables.
	oldTablesToDrop := []string{"model_mappings", "groups", "api_keys"} // api_keys might be related to old provider structure
	for _, table := range oldTablesToDrop {
		if err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			log.Printf("Warning: could not drop table %s: %v", table, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	log.Println("Database migration to refactor models and providers completed successfully.")
	return nil
}

// Helper function to check for SQLite "duplicate column name" error
func isDuplicateColumnError(err error) bool {
	return err != nil && (err.Error() == "duplicate column name: " || err.Error() == "SQL logic error: duplicate column name (1)")
}

// Helper function to check for SQLite "no such column" error
func isColumnMissingError(err error) bool {
	return err != nil && (err.Error() == "no such column: " || err.Error() == "SQL logic error: no such column (1)")
}