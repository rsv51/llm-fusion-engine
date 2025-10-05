package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

// ExportAll exports all settings to a JSON, YAML or Excel file.
func (h *ExportHandler) ExportAll(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	
	switch format {
	case "excel":
		h.exportToExcel(c)
	case "yaml":
		h.exportToYAML(c)
	default:
		h.exportToJSON(c)
	}
}

// ExportTemplate exports an Excel template with optional sample data
func (h *ExportHandler) ExportTemplate(c *gin.Context) {
	withSample := c.DefaultQuery("with_sample", "false") == "true"
	
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing Excel file: %v\n", err)
		}
	}()

	// Create Providers sheet
	h.createProvidersSheet(f, withSample)
	
	// Create Models sheet
	h.createModelsSheet(f, withSample)
	
	// Create Associations sheet
	h.createAssociationsSheet(f, withSample)

	// Set response headers
	timestamp := time.Now().Format("20060102_150405")
	filename := "llm_fusion_engine_template"
	if withSample {
		filename += "_with_sample"
	}
	filename += "_" + timestamp + ".xlsx"
	
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel template"})
	}
}

func (h *ExportHandler) exportToJSON(c *gin.Context) {
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

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=llm-fusion-engine-backup.json")
	if err := json.NewEncoder(c.Writer).Encode(migrationData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JSON"})
	}
}

func (h *ExportHandler) exportToYAML(c *gin.Context) {
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

	c.Header("Content-Type", "application/x-yaml")
	c.Header("Content-Disposition", "attachment; filename=llm-fusion-engine-backup.yaml")
	if err := yaml.NewEncoder(c.Writer).Encode(migrationData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate YAML"})
	}
}

func (h *ExportHandler) exportToExcel(c *gin.Context) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing Excel file: %v\n", err)
		}
	}()

	// Get actual data from database
	var groups []database.Group
	var providers []database.Provider
	var apiKeys []database.ApiKey
	var modelMappings []database.ModelMapping

	h.db.Find(&groups)
	h.db.Find(&providers)
	h.db.Find(&apiKeys)
	h.db.Find(&modelMappings)

	// Create sheets with actual data
	h.createProvidersSheetWithData(f, providers)
	h.createModelsSheetWithData(f, modelMappings)
	h.createAssociationsSheetWithData(f, modelMappings)

	// Set response headers
	timestamp := time.Now().Format("20060102_150405")
	filename := "llm_fusion_engine_config_" + timestamp + ".xlsx"
	
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file"})
	}
}

func (h *ExportHandler) createProvidersSheet(f *excelize.File, withSample bool) {
	// Create Providers sheet
	idx, _ := f.NewSheet("Providers")
	f.SetActiveSheet(idx)
	
	// Set headers
	headers := []string{"name", "type", "api_key", "base_url", "priority", "weight", "enabled"}
	for i, header := range headers {
		f.SetCellValue("Providers", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	if withSample {
		// Add sample data
		sampleData := [][]interface{}{
			{"OpenAI-Main", "openai", "sk-xxx", "https://api.openai.com/v1", 100, 100, true},
			{"Anthropic-Main", "anthropic", "sk-ant-xxx", "https://api.anthropic.com/v1", 100, 100, true},
		}
		for i, row := range sampleData {
			for j, value := range row {
				f.SetCellValue("Providers", fmt.Sprintf("%c%d", 'A'+j, i+2), value)
			}
		}
	}
}

func (h *ExportHandler) createModelsSheet(f *excelize.File, withSample bool) {
	// Create Models sheet
	idx, _ := f.NewSheet("Models")
	f.SetActiveSheet(idx)
	
	// Set headers
	headers := []string{"name", "remark", "max_retry", "timeout"}
	for i, header := range headers {
		f.SetCellValue("Models", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	if withSample {
		// Add sample data
		sampleData := [][]interface{}{
			{"gpt-4o", "GPT-4 Optimized", 3, 60},
			{"claude-3.5-sonnet", "Claude 3.5 Sonnet", 3, 60},
		}
		for i, row := range sampleData {
			for j, value := range row {
				f.SetCellValue("Models", fmt.Sprintf("%c%d", 'A'+j, i+2), value)
			}
		}
	}
}

func (h *ExportHandler) createAssociationsSheet(f *excelize.File, withSample bool) {
	// Create Associations sheet
	idx, _ := f.NewSheet("Associations")
	f.SetActiveSheet(idx)
	
	// Set headers
	headers := []string{"model_name", "provider_name", "provider_model", "supports_tools", "supports_vision", "weight", "enabled"}
	for i, header := range headers {
		f.SetCellValue("Associations", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	if withSample {
		// Add sample data
		sampleData := [][]interface{}{
			{"gpt-4o", "OpenAI-Main", "gpt-4o-2024-05-13", true, true, 100, true},
			{"claude-3.5-sonnet", "Anthropic-Main", "claude-3-5-sonnet-20241022", true, true, 100, true},
		}
		for i, row := range sampleData {
			for j, value := range row {
				f.SetCellValue("Associations", fmt.Sprintf("%c%d", 'A'+j, i+2), value)
			}
		}
	}
}

func (h *ExportHandler) createProvidersSheetWithData(f *excelize.File, providers []database.Provider) {
	// Create Providers sheet with actual data
	idx, _ := f.NewSheet("Providers")
	f.SetActiveSheet(idx)
	
	// Set headers
	headers := []string{"ID", "GroupID", "ProviderType", "Weight", "Enabled", "BaseURL", "Timeout", "MaxRetries", "HealthStatus", "LastChecked"}
	for i, header := range headers {
		f.SetCellValue("Providers", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	// Add actual data
	for i, provider := range providers {
		row := i + 2
		f.SetCellValue("Providers", fmt.Sprintf("A%d", row), provider.ID)
		f.SetCellValue("Providers", fmt.Sprintf("B%d", row), provider.GroupID)
		f.SetCellValue("Providers", fmt.Sprintf("C%d", row), provider.ProviderType)
		f.SetCellValue("Providers", fmt.Sprintf("D%d", row), provider.Weight)
		f.SetCellValue("Providers", fmt.Sprintf("E%d", row), provider.Enabled)
		f.SetCellValue("Providers", fmt.Sprintf("F%d", row), provider.BaseURL)
		f.SetCellValue("Providers", fmt.Sprintf("G%d", row), provider.Timeout)
		f.SetCellValue("Providers", fmt.Sprintf("H%d", row), provider.MaxRetries)
		f.SetCellValue("Providers", fmt.Sprintf("I%d", row), provider.HealthStatus)
		f.SetCellValue("Providers", fmt.Sprintf("J%d", row), provider.LastChecked)
	}
}

func (h *ExportHandler) createModelsSheetWithData(f *excelize.File, modelMappings []database.ModelMapping) {
	// Create Models sheet with actual data
	idx, _ := f.NewSheet("Models")
	f.SetActiveSheet(idx)
	
	// Set headers for model mappings
	headers := []string{"ID", "UserFriendlyName", "ProviderModelName", "ProviderID"}
	for i, header := range headers {
		f.SetCellValue("Models", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	// Add actual data
	for i, mapping := range modelMappings {
		row := i + 2
		f.SetCellValue("Models", fmt.Sprintf("A%d", row), mapping.ID)
		f.SetCellValue("Models", fmt.Sprintf("B%d", row), mapping.UserFriendlyName)
		f.SetCellValue("Models", fmt.Sprintf("C%d", row), mapping.ProviderModelName)
		f.SetCellValue("Models", fmt.Sprintf("D%d", row), mapping.ProviderID)
	}
}

func (h *ExportHandler) createAssociationsSheetWithData(f *excelize.File, modelMappings []database.ModelMapping) {
	// Create Associations sheet with actual data
	idx, _ := f.NewSheet("Associations")
	f.SetActiveSheet(idx)
	
	// Set headers
	headers := []string{"ID", "UserFriendlyName", "ProviderModelName", "ProviderID", "ProviderName", "ProviderType"}
	for i, header := range headers {
		f.SetCellValue("Associations", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	// Add actual data with provider info
	for i, mapping := range modelMappings {
		row := i + 2
		f.SetCellValue("Associations", fmt.Sprintf("A%d", row), mapping.ID)
		f.SetCellValue("Associations", fmt.Sprintf("B%d", row), mapping.UserFriendlyName)
		f.SetCellValue("Associations", fmt.Sprintf("C%d", row), mapping.ProviderModelName)
		f.SetCellValue("Associations", fmt.Sprintf("D%d", row), mapping.ProviderID)
		f.SetCellValue("Associations", fmt.Sprintf("E%d", row), mapping.Provider.ProviderType)
		f.SetCellValue("Associations", fmt.Sprintf("F%d", row), mapping.ProviderType)
	}
}