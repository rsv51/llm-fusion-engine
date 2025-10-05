package admin

import (
	"encoding/csv"
	"fmt"
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
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

// ExportGroups exports all groups to CSV or XLSX format
func (h *ExportHandler) ExportGroups(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	
	var groups []database.Group
	if err := h.db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}

	if format == "xlsx" {
		h.exportGroupsXLSX(c, groups)
	} else {
		h.exportGroupsCSV(c, groups)
	}
}

// exportGroupsCSV exports groups to CSV format
func (h *ExportHandler) exportGroupsCSV(c *gin.Context, groups []database.Group) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=groups.csv")
	
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "Name", "Enabled", "Priority", "LoadBalancePolicy", "CreatedAt"})

	// Write data
	for _, group := range groups {
		writer.Write([]string{
			strconv.FormatUint(uint64(group.ID), 10),
			group.Name,
			strconv.FormatBool(group.Enabled),
			strconv.Itoa(group.Priority),
			group.LoadBalancePolicy,
			group.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

// exportGroupsXLSX exports groups to XLSX format
func (h *ExportHandler) exportGroupsXLSX(c *gin.Context, groups []database.Group) {
	f := excelize.NewFile()
	defer f.Close()
	
	sheetName := "Groups"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	
	// Set header
	headers := []string{"ID", "Name", "Enabled", "Priority", "LoadBalancePolicy", "CreatedAt"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	
	// Set data
	for i, group := range groups {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), group.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), group.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), group.Enabled)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), group.Priority)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), group.LoadBalancePolicy)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), group.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=groups.xlsx")
	
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file"})
	}
}

// ExportKeys exports all API keys to CSV or XLSX format
func (h *ExportHandler) ExportKeys(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	
	var keys []database.ApiKey
	if err := h.db.Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve keys"})
		return
	}

	if format == "xlsx" {
		h.exportKeysXLSX(c, keys)
	} else {
		h.exportKeysCSV(c, keys)
	}
}

// exportKeysCSV exports keys to CSV format
func (h *ExportHandler) exportKeysCSV(c *gin.Context, keys []database.ApiKey) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=keys.csv")
	
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "ProviderID", "Key", "IsHealthy", "RpmLimit", "TpmLimit", "LastUsed"})

	// Write data
	for _, key := range keys {
		lastUsed := ""
		if !key.LastUsed.IsZero() {
			lastUsed = key.LastUsed.Format("2006-01-02 15:04:05")
		}
		writer.Write([]string{
			strconv.FormatUint(uint64(key.ID), 10),
			strconv.FormatUint(uint64(key.ProviderID), 10),
			key.Key,
			strconv.FormatBool(key.IsHealthy),
			strconv.Itoa(key.RpmLimit),
			strconv.Itoa(key.TpmLimit),
			lastUsed,
		})
	}
}

// exportKeysXLSX exports keys to XLSX format
func (h *ExportHandler) exportKeysXLSX(c *gin.Context, keys []database.ApiKey) {
	f := excelize.NewFile()
	defer f.Close()
	
	sheetName := "Keys"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	
	// Set header
	headers := []string{"ID", "ProviderID", "Key", "IsHealthy", "RpmLimit", "TpmLimit", "LastUsed"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	
	// Set data
	for i, key := range keys {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), key.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), key.ProviderID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), key.Key)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), key.IsHealthy)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), key.RpmLimit)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), key.TpmLimit)
		if !key.LastUsed.IsZero() {
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), key.LastUsed.Format("2006-01-02 15:04:05"))
		}
	}
	
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=keys.xlsx")
	
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file"})
	}
}

// ExportProviders exports all providers to CSV or XLSX format
func (h *ExportHandler) ExportProviders(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	
	var providers []database.Provider
	if err := h.db.Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve providers"})
		return
	}

	if format == "xlsx" {
		h.exportProvidersXLSX(c, providers)
	} else {
		h.exportProvidersCSV(c, providers)
	}
}

// exportProvidersCSV exports providers to CSV format
func (h *ExportHandler) exportProvidersCSV(c *gin.Context, providers []database.Provider) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=providers.csv")
	
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "GroupID", "ProviderType", "Weight", "Enabled", "CreatedAt"})

	// Write data
	for _, provider := range providers {
		writer.Write([]string{
			strconv.FormatUint(uint64(provider.ID), 10),
			strconv.FormatUint(uint64(provider.GroupID), 10),
			provider.ProviderType,
			strconv.FormatUint(uint64(provider.Weight), 10),
			strconv.FormatBool(provider.Enabled),
			provider.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

// exportProvidersXLSX exports providers to XLSX format
func (h *ExportHandler) exportProvidersXLSX(c *gin.Context, providers []database.Provider) {
	f := excelize.NewFile()
	defer f.Close()
	
	sheetName := "Providers"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	
	// Set header
	headers := []string{"ID", "GroupID", "ProviderType", "Weight", "Enabled", "CreatedAt"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}
	
	// Set data
	for i, provider := range providers {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), provider.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), provider.GroupID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), provider.ProviderType)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), provider.Weight)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), provider.Enabled)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), provider.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=providers.xlsx")
	
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file"})
	}
}