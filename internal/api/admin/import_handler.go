package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// ImportHandler handles data import operations
type ImportHandler struct {
	db *gorm.DB
}

// NewImportHandler creates a new ImportHandler
func NewImportHandler(db *gorm.DB) *ImportHandler {
	return &ImportHandler{db: db}
}

// ImportAll imports all settings from a JSON, YAML or Excel file.
func (h *ImportHandler) ImportAll(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	// Check file extension
	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
	if ext != "xlsx" && ext != "json" && ext != "yaml" && ext != "yml" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type. Only .xlsx, .json, .yaml, .yml files are supported"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	switch ext {
	case "xlsx":
		h.importFromExcel(c, f, file.Filename)
	case "json":
		h.importFromJSON(c, f)
	case "yaml", "yml":
		h.importFromYAML(c, f)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
	}
}

// ImportFromExcel imports configuration from Excel file with three-sheet structure
func (h *ImportHandler) ImportFromExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	h.importFromExcel(c, f, file.Filename)
}

func (h *ImportHandler) importFromExcel(c *gin.Context, f interface{}, filename string) {
	var excelFile *excelize.File
	var err error

	switch v := f.(type) {
	case *excelize.File:
		excelFile = v
	default:
		// Read from file interface
		excelFile, err = excelize.OpenReader(f.(interface{ Read([]byte) (int, error) }))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read Excel file"})
			return
		}
		defer func() {
			if err := excelFile.Close(); err != nil {
				fmt.Printf("Error closing Excel file: %v\n", err)
			}
		}()
	}

	result := h.processExcelImport(excelFile)
	c.JSON(http.StatusOK, result)
}

func (h *ImportHandler) importFromJSON(c *gin.Context, f interface{}) {
	var migrationData MigrationData
	
	var err error
	switch v := f.(type) {
	case interface{ Read([]byte) (int, error) }:
		err = json.NewDecoder(v.(interface{ Read([]byte) (int, error) })).Decode(&migrationData)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file reader"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		return h.importMigrationData(tx, migrationData)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}

func (h *ImportHandler) importFromYAML(c *gin.Context, f interface{}) {
	var migrationData MigrationData
	
	var err error
	switch v := f.(type) {
	case interface{ Read([]byte) (int, error) }:
		err = yaml.NewDecoder(v.(interface{ Read([]byte) (int, error) })).Decode(&migrationData)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file reader"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YAML format"})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		return h.importMigrationData(tx, migrationData)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}

func (h *ImportHandler) importMigrationData(tx *gorm.DB, data MigrationData) error {
	// Clear existing data
	if err := tx.Exec("DELETE FROM model_mappings").Error; err != nil { return err }
	if err := tx.Exec("DELETE FROM api_keys").Error; err != nil { return err }
	if err := tx.Exec("DELETE FROM providers").Error; err != nil { return err }
	if err := tx.Exec("DELETE FROM groups").Error; err != nil { return err }

	// Import new data
	if len(data.Groups) > 0 {
		if err := tx.Create(&data.Groups).Error; err != nil { return err }
	}
	if len(data.Providers) > 0 {
		if err := tx.Create(&data.Providers).Error; err != nil { return err }
	}
	if len(data.ApiKeys) > 0 {
		if err := tx.Create(&data.ApiKeys).Error; err != nil { return err }
	}
	if len(data.ModelMappings) > 0 {
		if err := tx.Create(&data.ModelMappings).Error; err != nil { return err }
	}

	return nil
}

func (h *ImportHandler) processExcelImport(f *excelize.File) gin.H {
	result := gin.H{
		"filename": "",
		"result": gin.H{
			"providers":    gin.H{"total": 0, "imported": 0, "skipped": 0, "errors": []interface{}{}},
			"models":      gin.H{"total": 0, "imported": 0, "skipped": 0, "errors": []interface{}{}},
			"associations": gin.H{"total": 0, "imported": 0, "skipped": 0, "errors": []interface{}{}},
			"summary":     gin.H{"total_imported": 0, "total_skipped": 0, "total_errors": 0},
		},
	}

	// Process each sheet in order: Providers -> Models -> Associations
	if providers, ok := f.GetSheetMap()["Providers"]; ok {
		h.processProvidersSheet(f, providers, result)
	}
	
	if models, ok := f.GetSheetMap()["Models"]; ok {
		h.processModelsSheet(f, models, result)
	}
	
	if associations, ok := f.GetSheetMap()["Associations"]; ok {
		h.processAssociationsSheet(f, associations, result)
	}

	// Calculate summary
	providersResult := result["result"].(gin.H)["providers"].(gin.H)
	modelsResult := result["result"].(gin.H)["models"].(gin.H)
	associationsResult := result["result"].(gin.H)["associations"].(gin.H)
	
	summary := result["result"].(gin.H)["summary"].(gin.H)
	summary["total_imported"] = providersResult["imported"].(int) + modelsResult["imported"].(int) + associationsResult["imported"].(int)
	summary["total_skipped"] = providersResult["skipped"].(int) + modelsResult["skipped"].(int) + associationsResult["skipped"].(int)
	summary["total_errors"] = len(providersResult["errors"].([]interface{})) + len(modelsResult["errors"].([]interface{})) + len(associationsResult["errors"].([]interface{}))

	return result
}

func (h *ImportHandler) processProvidersSheet(f *excelize.File, sheet string, result gin.H) {
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		return
	}

	providersResult := result["result"].(gin.H)["providers"].(gin.H)
	headers := rows[0]
	
	// Find column indices
	nameIdx := findColumnIndex(headers, "name")
	typeIdx := findColumnIndex(headers, "type")
	apiKeyIdx := findColumnIndex(headers, "api_key")
	baseURLIdx := findColumnIndex(headers, "base_url")
	priorityIdx := findColumnIndex(headers, "priority")
	weightIdx := findColumnIndex(headers, "weight")
	enabledIdx := findColumnIndex(headers, "enabled")

	if nameIdx == -1 || typeIdx == -1 || apiKeyIdx == -1 {
		providersResult["errors"] = append(providersResult["errors"].([]interface{}), gin.H{"row": 1, "field": "headers", "error": "Required columns (name, type, api_key) not found"})
		return
	}

	// Process data rows
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) <= nameIdx || row[nameIdx] == "" {
			continue // Skip empty rows
		}

		providersResult["total"] = providersResult["total"].(int) + 1

		// Check if provider already exists
		var existingProvider database.Provider
		if err := h.db.Where("provider_type = ?", row[nameIdx]).First(&existingProvider).Error; err == nil {
			providersResult["skipped"] = providersResult["skipped"].(int) + 1
			continue
		}

		// Parse boolean values
		enabled := true
		if enabledIdx != -1 && len(row) > enabledIdx {
			enabled = parseBool(row[enabledIdx])
		}

		// Parse integer values
		priority := 100
		if priorityIdx != -1 && len(row) > priorityIdx && row[priorityIdx] != "" {
			if val, err := strconv.Atoi(row[priorityIdx]); err == nil {
				priority = val
			}
		}

		weight := 100
		if weightIdx != -1 && len(row) > weightIdx && row[weightIdx] != "" {
			if val, err := strconv.Atoi(row[weightIdx]); err == nil {
				weight = val
			}
		}

		provider := database.Provider{
			ProviderType: row[nameIdx],
			Weight:      uint(weight),
			Enabled:     enabled,
		}

		if baseURLIdx != -1 && len(row) > baseURLIdx && row[baseURLIdx] != "" {
			provider.BaseURL = row[baseURLIdx]
		}

		if err := h.db.Create(&provider).Error; err != nil {
			providersResult["errors"] = append(providersResult["errors"].([]interface{}), gin.H{"row": i + 1, "field": "database", "error": err.Error()})
		} else {
			providersResult["imported"] = providersResult["imported"].(int) + 1
		}
	}
}

func (h *ImportHandler) processModelsSheet(f *excelize.File, sheet string, result gin.H) {
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		return
	}

	modelsResult := result["result"].(gin.H)["models"].(gin.H)
	headers := rows[0]
	
	// Find column indices
	nameIdx := findColumnIndex(headers, "name")
	remarkIdx := findColumnIndex(headers, "remark")
	maxRetryIdx := findColumnIndex(headers, "max_retry")
	timeoutIdx := findColumnIndex(headers, "timeout")

	if nameIdx == -1 {
		modelsResult["errors"] = append(modelsResult["errors"].([]interface{}), gin.H{"row": 1, "field": "headers", "error": "Required column 'name' not found"})
		return
	}

	// Process data rows
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) <= nameIdx || row[nameIdx] == "" {
			continue // Skip empty rows
		}

		modelsResult["total"] = modelsResult["total"].(int) + 1

		// Check if model mapping already exists
		var existingMapping database.ModelMapping
		if err := h.db.Where("user_friendly_name = ?", row[nameIdx]).First(&existingMapping).Error; err == nil {
			modelsResult["skipped"] = modelsResult["skipped"].(int) + 1
			continue
		}

		// Parse integer values
		maxRetry := 3
		if maxRetryIdx != -1 && len(row) > maxRetryIdx && row[maxRetryIdx] != "" {
			if val, err := strconv.Atoi(row[maxRetryIdx]); err == nil {
				maxRetry = val
			}
		}

		timeout := 30
		if timeoutIdx != -1 && len(row) > timeoutIdx && row[timeoutIdx] != "" {
			if val, err := strconv.Atoi(row[timeoutIdx]); err == nil {
				timeout = val
			}
		}

		// Create model mapping (since we don't have separate models table in current schema)
		mapping := database.ModelMapping{
			UserFriendlyName:  row[nameIdx],
			ProviderModelName: row[nameIdx], // Use same name for provider model
			ProviderID:        1, // Default provider ID, should be configurable
		}

		if remarkIdx != -1 && len(row) > remarkIdx && row[remarkIdx] != "" {
			// Store remark in a different way since we don't have remark field
			// We could use userFriendlyName to include remark
			mapping.UserFriendlyName = row[nameIdx]
			if row[remarkIdx] != "" {
				mapping.UserFriendlyName = row[remarkIdx]
			}
		}

		if err := h.db.Create(&mapping).Error; err != nil {
			modelsResult["errors"] = append(modelsResult["errors"].([]interface{}), gin.H{"row": i + 1, "field": "database", "error": err.Error()})
		} else {
			modelsResult["imported"] = modelsResult["imported"].(int) + 1
		}
	}
}

func (h *ImportHandler) processAssociationsSheet(f *excelize.File, sheet string, result gin.H) {
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		return
	}

	associationsResult := result["result"].(gin.H)["associations"].(gin.H)
	headers := rows[0]
	
	// Find column indices
	modelNameIdx := findColumnIndex(headers, "model_name")
	providerNameIdx := findColumnIndex(headers, "provider_name")
	providerModelIdx := findColumnIndex(headers, "provider_model")
	supportsToolsIdx := findColumnIndex(headers, "supports_tools")
	supportsVisionIdx := findColumnIndex(headers, "supports_vision")
	weightIdx := findColumnIndex(headers, "weight")
	enabledIdx := findColumnIndex(headers, "enabled")

	if modelNameIdx == -1 || providerNameIdx == -1 || providerModelIdx == -1 {
		associationsResult["errors"] = append(associationsResult["errors"].([]interface{}), gin.H{"row": 1, "field": "headers", "error": "Required columns (model_name, provider_name, provider_model) not found"})
		return
	}

	// Process data rows
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) <= modelNameIdx || row[modelNameIdx] == "" {
			continue // Skip empty rows
		}

		associationsResult["total"] = associationsResult["total"].(int) + 1

		// Check if association already exists by user friendly name
		var existingMapping database.ModelMapping
		if err := h.db.Where("user_friendly_name = ?", row[modelNameIdx]).First(&existingMapping).Error; err == nil {
			associationsResult["skipped"] = associationsResult["skipped"].(int) + 1
			continue
		}

		// Parse boolean values
		supportsTools := false
		if supportsToolsIdx != -1 && len(row) > supportsToolsIdx {
			supportsTools = parseBool(row[supportsToolsIdx])
		}

		supportsVision := false
		if supportsVisionIdx != -1 && len(row) > supportsVisionIdx {
			supportsVision = parseBool(row[supportsVisionIdx])
		}

		enabled := true
		if enabledIdx != -1 && len(row) > enabledIdx {
			enabled = parseBool(row[enabledIdx])
		}

		// Parse weight
		weight := 100
		if weightIdx != -1 && len(row) > weightIdx && row[weightIdx] != "" {
			if val, err := strconv.Atoi(row[weightIdx]); err == nil {
				weight = val
			}
		}

		// Find provider by type
		var provider database.Provider
		if err := h.db.Where("provider_type = ?", row[providerNameIdx]).First(&provider).Error; err != nil {
			associationsResult["errors"] = append(associationsResult["errors"].([]interface{}), gin.H{"row": i + 1, "field": "provider_name", "error": "Provider not found: " + row[providerNameIdx]})
			continue
		}

		mapping := database.ModelMapping{
			UserFriendlyName:  row[modelNameIdx],
			ProviderModelName: row[providerModelIdx],
			ProviderID:        provider.ID,
		}

		if err := h.db.Create(&mapping).Error; err != nil {
			associationsResult["errors"] = append(associationsResult["errors"].([]interface{}), gin.H{"row": i + 1, "field": "database", "error": err.Error()})
		} else {
			associationsResult["imported"] = associationsResult["imported"].(int) + 1
		}
	}
}

// Helper functions
func findColumnIndex(headers []string, columnName string) int {
	for i, header := range headers {
		if strings.ToLower(strings.TrimSpace(header)) == strings.ToLower(columnName) {
			return i
		}
	}
	return -1
}

func parseBool(value string) bool {
	val := strings.ToLower(strings.TrimSpace(value))
	return val == "true" || val == "1" || val == "yes" || val == "t"
}