package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"llm-fusion-engine/internal/database"
)

// ImportHandler handles data import operations
type ImportHandler struct {
	db *gorm.DB
}

// NewImportHandler creates a new ImportHandler
func NewImportHandler(db *gorm.DB) *ImportHandler {
	return &ImportHandler{db: db}
}

// ImportAll imports all settings from an Excel file.
func (h *ImportHandler) ImportAll(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	// Check file extension
	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
	if ext != "xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type. Only .xlsx files are supported"})
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


func (h *ImportHandler) processExcelImport(f *excelize.File) gin.H {
	result := gin.H{
		"filename": "",
		"result": gin.H{
			"providers":              gin.H{"total": 0, "imported": 0, "skipped": 0, "errors": []interface{}{}},
			"models":                 gin.H{"total": 0, "imported": 0, "skipped": 0, "errors": []interface{}{}},
			"modelProviderMappings":  gin.H{"total": 0, "imported": 0, "skipped": 0, "errors": []interface{}{}},
			"summary":                gin.H{"total_imported": 0, "total_skipped": 0, "total_errors": 0},
		},
	}

	// Process each sheet in order: Providers -> Models -> Associations
	sheetMap := f.GetSheetMap()
	sheetExists := func(name string) bool {
		for _, sheetName := range sheetMap {
			if sheetName == name {
				return true
			}
		}
		return false
	}

	if sheetExists("Providers") {
		h.processProvidersSheet(f, "Providers", result)
	}

	if sheetExists("Models") {
		h.processModelsSheet(f, "Models", result)
	}

	if sheetExists("ModelProviderMappings") {
		h.processModelProviderMappingsSheet(f, "ModelProviderMappings", result)
	}

	// Calculate summary
	providersResult := result["result"].(gin.H)["providers"].(gin.H)
	modelsResult := result["result"].(gin.H)["models"].(gin.H)
	modelProviderMappingsResult := result["result"].(gin.H)["modelProviderMappings"].(gin.H)
	
	summary := result["result"].(gin.H)["summary"].(gin.H)
	summary["total_imported"] = providersResult["imported"].(int) + modelsResult["imported"].(int) + modelProviderMappingsResult["imported"].(int)
	summary["total_skipped"] = providersResult["skipped"].(int) + modelsResult["skipped"].(int) + modelProviderMappingsResult["skipped"].(int)
	summary["total_errors"] = len(providersResult["errors"].([]interface{})) + len(modelsResult["errors"].([]interface{})) + len(modelProviderMappingsResult["errors"].([]interface{}))

	return result
}

func (h *ImportHandler) processProvidersSheet(f *excelize.File, sheetName string, result gin.H) {
	rows, err := f.GetRows(sheetName)
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
		if err := h.db.Where("type = ?", row[nameIdx]).First(&existingProvider).Error; err == nil {
			providersResult["skipped"] = providersResult["skipped"].(int) + 1
			continue
		}

		// Parse boolean values
		enabled := true
		if enabledIdx != -1 && len(row) > enabledIdx {
			enabled = parseBool(row[enabledIdx])
		}

		// Parse integer values
		if priorityIdx != -1 && len(row) > priorityIdx && row[priorityIdx] != "" {
			// priority is not used
		}

		weight := 100
		if weightIdx != -1 && len(row) > weightIdx && row[weightIdx] != "" {
			if val, err := strconv.Atoi(row[weightIdx]); err == nil {
				weight = val
			}
		}

	provider := database.Provider{
		Name:    row[nameIdx],
		Type:    row[typeIdx],
		Weight:  uint(weight),
		Enabled: enabled,
	}

		// Create config JSON with API key and baseURL
		config := make(map[string]interface{})
		
		if apiKeyIdx != -1 && len(row) > apiKeyIdx && row[apiKeyIdx] != "" {
			config["apiKey"] = row[apiKeyIdx]
		}
		
		if baseURLIdx != -1 && len(row) > baseURLIdx && row[baseURLIdx] != "" {
			config["baseUrl"] = row[baseURLIdx]
		}
		
		if len(config) > 0 {
			configJSON, _ := json.Marshal(config)
			provider.Config = string(configJSON)
		}

		if err := h.db.Create(&provider).Error; err != nil {
			providersResult["errors"] = append(providersResult["errors"].([]interface{}), gin.H{"row": i + 1, "field": "database", "error": err.Error()})
		} else {
			providersResult["imported"] = providersResult["imported"].(int) + 1
		}
	}
}

func (h *ImportHandler) processModelsSheet(f *excelize.File, sheetName string, result gin.H) {
	rows, err := f.GetRows(sheetName)
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

		// Check if model already exists
		var existingModel database.Model
		if err := h.db.Where("name = ?", row[nameIdx]).First(&existingModel).Error; err == nil {
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

		// Create model
		model := database.Model{
			Name:     row[nameIdx],
			MaxRetry: maxRetry,
			Timeout:  timeout,
			Enabled:  true,
		}

		if remarkIdx != -1 && len(row) > remarkIdx && row[remarkIdx] != "" {
			model.Remark = row[remarkIdx]
		}

		if err := h.db.Create(&model).Error; err != nil {
			modelsResult["errors"] = append(modelsResult["errors"].([]interface{}), gin.H{"row": i + 1, "field": "database", "error": err.Error()})
		} else {
			modelsResult["imported"] = modelsResult["imported"].(int) + 1
		}
	}
}

func (h *ImportHandler) processModelProviderMappingsSheet(f *excelize.File, sheetName string, result gin.H) {
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return
	}

	modelProviderMappingsResult := result["result"].(gin.H)["modelProviderMappings"].(gin.H)
	headers := rows[0]

	// Find column indices
	modelNameIdx := findColumnIndex(headers, "modelname")
	providerNameIdx := findColumnIndex(headers, "providername")
	providerModelIdx := findColumnIndex(headers, "providermodel")
	toolCallIdx := findColumnIndex(headers, "toolcall")
	structuredOutputIdx := findColumnIndex(headers, "structuredoutput")
	imageIdx := findColumnIndex(headers, "image")
	weightIdx := findColumnIndex(headers, "weight")
	enabledIdx := findColumnIndex(headers, "enabled")

	if modelNameIdx == -1 || providerNameIdx == -1 || providerModelIdx == -1 {
		modelProviderMappingsResult["errors"] = append(modelProviderMappingsResult["errors"].([]interface{}), gin.H{"row": 1, "field": "headers", "error": "Required columns (ModelName, ProviderName, ProviderModel) not found"})
		return
	}

	// Process data rows
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) <= modelNameIdx || row[modelNameIdx] == "" {
			continue // Skip empty rows
		}

		modelProviderMappingsResult["total"] = modelProviderMappingsResult["total"].(int) + 1

		// Find provider by name
		var provider database.Provider
		if err := h.db.Where("name = ?", row[providerNameIdx]).First(&provider).Error; err != nil {
			modelProviderMappingsResult["errors"] = append(modelProviderMappingsResult["errors"].([]interface{}), gin.H{"row": i+1, "field": "provider_name", "error": "Provider not found: " + row[providerNameIdx]})
			continue
		}

		// Find model by name
		var model database.Model
		if err := h.db.Where("name = ?", row[modelNameIdx]).First(&model).Error; err != nil {
			modelProviderMappingsResult["errors"] = append(modelProviderMappingsResult["errors"].([]interface{}), gin.H{"row": i+1, "field": "model_name", "error": "Model not found: " + row[modelNameIdx]})
			continue
		}

		// Check if association already exists
		var existingMapping database.ModelProviderMapping
		if err := h.db.Where("model_id = ? AND provider_id = ?", model.ID, provider.ID).First(&existingMapping).Error; err == nil {
			modelProviderMappingsResult["skipped"] = modelProviderMappingsResult["skipped"].(int) + 1
			continue
		}

		// Parse boolean values
		toolCall := false
		if toolCallIdx != -1 && len(row) > toolCallIdx {
			toolCall = parseBool(row[toolCallIdx])
		}

		structuredOutput := false
		if structuredOutputIdx != -1 && len(row) > structuredOutputIdx {
			structuredOutput = parseBool(row[structuredOutputIdx])
		}

		image := false
		if imageIdx != -1 && len(row) > imageIdx {
			image = parseBool(row[imageIdx])
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

		mapping := database.ModelProviderMapping{
			ModelID:          model.ID,
			ProviderID:       provider.ID,
			ProviderModel:    row[providerModelIdx],
			ToolCall:         &toolCall,
			StructuredOutput: &structuredOutput,
			Image:            &image,
			Weight:           weight,
			Enabled:          enabled,
		}

		if err := h.db.Create(&mapping).Error; err != nil {
			modelProviderMappingsResult["errors"] = append(modelProviderMappingsResult["errors"].([]interface{}), gin.H{"row": i+1, "field": "database", "error": err.Error()})
		} else {
			modelProviderMappingsResult["imported"] = modelProviderMappingsResult["imported"].(int) + 1
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