package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"llm-fusion-engine/internal/database"
)

// ExportHandler handles data export operations
type ExportHandler struct {
	db *gorm.DB
}

// NewExportHandler creates a new ExportHandler
func NewExportHandler(db *gorm.DB) *ExportHandler {
	return &ExportHandler{db: db}
}

// ExportAll exports all settings to an Excel file.
func (h *ExportHandler) ExportAll(c *gin.Context) {
	h.exportToExcel(c)
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

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")
	
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
	var modelProviderMappings []database.ModelProviderMapping
	var models []database.Model

	h.db.Find(&groups)
	h.db.Find(&providers)
	h.db.Find(&apiKeys)
	h.db.Find(&models)
	h.db.Preload("Provider").Preload("Model").Find(&modelProviderMappings)

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")
	
	// Create sheets with actual data
	h.createProvidersSheetWithData(f, providers)
	h.createModelsSheetWithData(f, models)
	h.createModelProviderMappingsSheetWithData(f, modelProviderMappings)

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
	// Create ModelProviderMappings sheet
	idx, _ := f.NewSheet("ModelProviderMappings")
	f.SetActiveSheet(idx)
	
	// Set headers
	headers := []string{"ModelName", "ProviderName", "ProviderModel", "ToolCall", "StructuredOutput", "Image", "Weight", "Enabled"}
	for i, header := range headers {
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	if withSample {
		// Add sample data
		sampleData := [][]interface{}{
			{"gpt-4o", "OpenAI-Main", "gpt-4o-2024-05-13", true, false, true, 100, true},
			{"claude-3.5-sonnet", "Anthropic-Main", "claude-3-5-sonnet-20241022", true, true, false, 100, true},
		}
		for i, row := range sampleData {
			for j, value := range row {
				f.SetCellValue("ModelProviderMappings", fmt.Sprintf("%c%d", 'A'+j, i+2), value)
			}
		}
	}
}

func (h *ExportHandler) createProvidersSheetWithData(f *excelize.File, providers []database.Provider) {
	// Create Providers sheet with actual data
	idx, _ := f.NewSheet("Providers")
	f.SetActiveSheet(idx)
	
	// Set headers
	headers := []string{"ID", "Name", "Type", "Config", "Console", "Enabled", "Weight", "HealthStatus", "LastChecked", "Latency"}
	for i, header := range headers {
		f.SetCellValue("Providers", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	// Add actual data
	for i, provider := range providers {
		row := i + 2
		f.SetCellValue("Providers", fmt.Sprintf("A%d", row), provider.ID)
		f.SetCellValue("Providers", fmt.Sprintf("B%d", row), provider.Name)
		f.SetCellValue("Providers", fmt.Sprintf("C%d", row), provider.Type)
		f.SetCellValue("Providers", fmt.Sprintf("D%d", row), provider.Config)
		f.SetCellValue("Providers", fmt.Sprintf("E%d", row), provider.Console)
		f.SetCellValue("Providers", fmt.Sprintf("F%d", row), provider.Enabled)
		f.SetCellValue("Providers", fmt.Sprintf("G%d", row), provider.Weight)
		f.SetCellValue("Providers", fmt.Sprintf("H%d", row), provider.HealthStatus)
		f.SetCellValue("Providers", fmt.Sprintf("I%d", row), provider.LastChecked)
		f.SetCellValue("Providers", fmt.Sprintf("J%d", row), provider.Latency)
	}
}

func (h *ExportHandler) createModelsSheetWithData(f *excelize.File, models []database.Model) {
	// Create Models sheet with actual data
	idx, _ := f.NewSheet("Models")
	f.SetActiveSheet(idx)
	
	// Set headers for models
	headers := []string{"ID", "Name", "Remark", "MaxRetry", "Timeout", "Enabled"}
	for i, header := range headers {
		f.SetCellValue("Models", fmt.Sprintf("%c1", 'A'+i), header)
	}
	
	// Add actual data
	for i, model := range models {
		row := i + 2
		f.SetCellValue("Models", fmt.Sprintf("A%d", row), model.ID)
		f.SetCellValue("Models", fmt.Sprintf("B%d", row), model.Name)
		f.SetCellValue("Models", fmt.Sprintf("C%d", row), model.Remark)
		f.SetCellValue("Models", fmt.Sprintf("D%d", row), model.MaxRetry)
		f.SetCellValue("Models", fmt.Sprintf("E%d", row), model.Timeout)
		f.SetCellValue("Models", fmt.Sprintf("F%d", row), model.Enabled)
	}
}

func (h *ExportHandler) createModelProviderMappingsSheetWithData(f *excelize.File, modelProviderMappings []database.ModelProviderMapping) {
	// Create ModelProviderMappings sheet with actual data
	idx, _ := f.NewSheet("ModelProviderMappings")
	f.SetActiveSheet(idx)

	// Set headers
	headers := []string{"ModelName", "ProviderName", "ProviderModel", "ToolCall", "StructuredOutput", "Image", "Weight", "Enabled"}
	for i, header := range headers {
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("%c1", 'A'+i), header)
	}

	// Add actual data with model and provider info
	for i, mapping := range modelProviderMappings {
		row := i + 2
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("A%d", row), mapping.Model.Name)
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("B%d", row), mapping.Provider.Name)
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("C%d", row), mapping.ProviderModel)
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("D%d", row), mapping.ToolCall)
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("E%d", row), mapping.StructuredOutput)
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("F%d", row), mapping.Image)
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("G%d", row), mapping.Weight)
		f.SetCellValue("ModelProviderMappings", fmt.Sprintf("H%d", row), mapping.Enabled)
	}
}