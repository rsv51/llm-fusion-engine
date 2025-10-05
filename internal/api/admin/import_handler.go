package admin

import (
	"encoding/csv"
	"fmt"
	"llm-fusion-engine/internal/database"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/xuri/excelize/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mime/multipart"
)

// ImportHandler handles data import operations
type ImportHandler struct {
	db *gorm.DB
}

// NewImportHandler creates a new ImportHandler
func NewImportHandler(db *gorm.DB) *ImportHandler {
	return &ImportHandler{db: db}
}

// ImportGroups imports groups from a CSV or XLSX file
func (h *ImportHandler) ImportGroups(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	ext := filepath.Ext(file.Filename)
	var groups []database.Group
	
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	if ext == ".csv" {
		groups, err = h.parseGroupsCSV(f)
	} else if ext == ".xlsx" {
		groups, err = h.parseGroupsXLSX(f, file.Size)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format. Please use CSV or XLSX."})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse file: %v", err)})
		return
	}

	if len(groups) > 0 {
		if err := h.db.CreateInBatches(&groups, 100).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import groups to database"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully imported %d groups.", len(groups)),
	})
}

// ImportKeys imports keys from a CSV or XLSX file
func (h *ImportHandler) ImportKeys(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	ext := filepath.Ext(file.Filename)
	var keys []database.ApiKey

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	if ext == ".csv" {
		keys, err = h.parseKeysCSV(f)
	} else if ext == ".xlsx" {
		keys, err = h.parseKeysXLSX(f)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format. Please use CSV or XLSX."})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse file: %v", err)})
		return
	}

	if len(keys) > 0 {
		if err := h.db.CreateInBatches(&keys, 100).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import keys to database"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully imported %d keys.", len(keys)),
	})
}

// ImportProviders imports providers from a CSV or XLSX file
func (h *ImportHandler) ImportProviders(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	ext := filepath.Ext(file.Filename)
	var providers []database.Provider

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	if ext == ".csv" {
		providers, err = h.parseProvidersCSV(f)
	} else if ext == ".xlsx" {
		providers, err = h.parseProvidersXLSX(f)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format. Please use CSV or XLSX."})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse file: %v", err)})
		return
	}

	if len(providers) > 0 {
		if err := h.db.CreateInBatches(&providers, 100).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import providers to database"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully imported %d providers.", len(providers)),
	})
}

func (h *ImportHandler) parseGroupsCSV(f multipart.File) ([]database.Group, error) {
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var groups []database.Group
	// Skip header row
	for _, record := range records[1:] {
		enabled, _ := strconv.ParseBool(record[2])
		priority, _ := strconv.Atoi(record[3])
		
		groups = append(groups, database.Group{
			Name:              record[1],
			Enabled:           enabled,
			Priority:          priority,
			LoadBalancePolicy: record[4],
		})
	}
	return groups, nil
}

func (h *ImportHandler) parseKeysCSV(f multipart.File) ([]database.ApiKey, error) {
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var keys []database.ApiKey
	for _, record := range records[1:] {
		providerID, _ := strconv.ParseUint(record[1], 10, 32)
		isHealthy, _ := strconv.ParseBool(record[3])
		rpmLimit, _ := strconv.Atoi(record[4])
		tpmLimit, _ := strconv.Atoi(record[5])

		keys = append(keys, database.ApiKey{
			ProviderID: uint(providerID),
			Key:        record[2],
			IsHealthy:  isHealthy,
			RpmLimit:   rpmLimit,
			TpmLimit:   tpmLimit,
		})
	}
	return keys, nil
}

func (h *ImportHandler) parseKeysXLSX(f multipart.File) ([]database.ApiKey, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}

	sheetName := "Keys"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("sheet 'Keys' not found or failed to read")
	}

	var keys []database.ApiKey
	for _, row := range rows[1:] {
		providerID, _ := strconv.ParseUint(row[1], 10, 32)
		isHealthy, _ := strconv.ParseBool(row[3])
		rpmLimit, _ := strconv.Atoi(row[4])
		tpmLimit, _ := strconv.Atoi(row[5])

		keys = append(keys, database.ApiKey{
			ProviderID: uint(providerID),
			Key:        row[2],
			IsHealthy:  isHealthy,
			RpmLimit:   rpmLimit,
			TpmLimit:   tpmLimit,
		})
	}
	return keys, nil
}

func (h *ImportHandler) parseProvidersCSV(f multipart.File) ([]database.Provider, error) {
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var providers []database.Provider
	for _, record := range records[1:] {
		groupID, _ := strconv.ParseUint(record[1], 10, 32)
		weight, _ := strconv.ParseUint(record[3], 10, 32)
		enabled, _ := strconv.ParseBool(record[4])

		providers = append(providers, database.Provider{
			GroupID:      uint(groupID),
			ProviderType: record[2],
			Weight:       uint(weight),
			Enabled:      enabled,
		})
	}
	return providers, nil
}

func (h *ImportHandler) parseProvidersXLSX(f multipart.File) ([]database.Provider, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}

	sheetName := "Providers"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("sheet 'Providers' not found or failed to read")
	}

	var providers []database.Provider
	for _, row := range rows[1:] {
		groupID, _ := strconv.ParseUint(row[1], 10, 32)
		weight, _ := strconv.ParseUint(row[3], 10, 32)
		enabled, _ := strconv.ParseBool(row[4])

		providers = append(providers, database.Provider{
			GroupID:      uint(groupID),
			ProviderType: row[2],
			Weight:       uint(weight),
			Enabled:      enabled,
		})
	}
	return providers, nil
}

func (h *ImportHandler) parseGroupsXLSX(f multipart.File, size int64) ([]database.Group, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}
	
	sheetName := "Groups"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("sheet 'Groups' not found or failed to read")
	}

	var groups []database.Group
	// Skip header row
	for _, row := range rows[1:] {
		enabled, _ := strconv.ParseBool(row[2])
		priority, _ := strconv.Atoi(row[3])

		groups = append(groups, database.Group{
			Name:              row[1],
			Enabled:           enabled,
			Priority:          priority,
			LoadBalancePolicy: row[4],
		})
	}
	return groups, nil
}