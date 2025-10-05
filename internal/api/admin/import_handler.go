package admin

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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

// ImportAll imports all settings from a JSON or YAML file.
func (h *ImportHandler) ImportAll(c *gin.Context) {
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

	var migrationData MigrationData
	contentType := file.Header.Get("Content-Type")

	if contentType == "application/json" {
		if err := json.NewDecoder(f).Decode(&migrationData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
	} else if contentType == "application/x-yaml" || contentType == "text/yaml" {
		if err := yaml.NewDecoder(f).Decode(&migrationData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YAML format"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}

	// Use a transaction to ensure all or nothing
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Clear existing data
		if err := tx.Exec("DELETE FROM model_mappings").Error; err != nil { return err }
		if err := tx.Exec("DELETE FROM api_keys").Error; err != nil { return err }
		if err := tx.Exec("DELETE FROM providers").Error; err != nil { return err }
		if err := tx.Exec("DELETE FROM groups").Error; err != nil { return err }

		// Import new data
		if len(migrationData.Groups) > 0 {
			if err := tx.Create(&migrationData.Groups).Error; err != nil { return err }
		}
		if len(migrationData.Providers) > 0 {
			if err := tx.Create(&migrationData.Providers).Error; err != nil { return err }
		}
		if len(migrationData.ApiKeys) > 0 {
			if err := tx.Create(&migrationData.ApiKeys).Error; err != nil { return err }
		}
		if len(migrationData.ModelMappings) > 0 {
			if err := tx.Create(&migrationData.ModelMappings).Error; err != nil { return err }
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}