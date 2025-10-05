package admin

import (
	"encoding/json"
	"fmt"
	"llm-fusion-engine/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// ExportHandler handles data export operations
type ExportHandler struct {
	db *gorm.DB
}

// NewExportHandler creates a new ExportHandler
func NewExportHandler(db *gorm.DB) *ExportHandler {
	return &ExportHandler{db: db}
}

// ExportAll exports all settings to a JSON or YAML file.
func (h *ExportHandler) ExportAll(c *gin.Context) {
	var migrationData MigrationData
	var err error

	// Fetch all data
	if err = h.db.Find(&migrationData.Groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}
	if err = h.db.Find(&migrationData.Providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve providers"})
		return
	}
	if err = h.db.Find(&migrationData.ApiKeys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve api keys"})
		return
	}
	if err = h.db.Find(&migrationData.ModelMappings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve model mappings"})
		return
	}

	format := c.DefaultQuery("format", "json")

	if format == "yaml" {
		c.Header("Content-Type", "application/x-yaml")
		c.Header("Content-Disposition", "attachment; filename=llm-fusion-engine-backup.yaml")
		if err := yaml.NewEncoder(c.Writer).Encode(migrationData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate YAML"})
		}
	} else {
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=llm-fusion-engine-backup.json")
		if err := json.NewEncoder(c.Writer).Encode(migrationData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JSON"})
		}
	}
}