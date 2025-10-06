package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProviderHandler handles CRUD operations for providers.
type ProviderHandler struct {
	db *gorm.DB
}

// NewProviderHandler creates a new ProviderHandler.
func NewProviderHandler(db *gorm.DB) *ProviderHandler {
	return &ProviderHandler{db: db}
}

// CreateProvider creates a new provider.
func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var provider database.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create provider"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// GetProviders retrieves all providers with pagination.
func (h *ProviderHandler) GetProviders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var providers []database.Provider
	var total int64

	// Count total records
	if err := h.db.Model(&database.Provider{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count providers"})
		return
	}

	// Get paginated records
	// Note: Preload("ApiKeys") is removed as ApiKeys are now part of the JSON config.
	// If direct access to ApiKeys is still needed, a separate endpoint or logic would be required.
	if err := h.db.Offset(offset).Limit(pageSize).Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve providers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": providers,
		"pagination": gin.H{
			"page":      page,
			"pageSize":  pageSize,
			"total":     total,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetProvider retrieves a single provider by ID.
func (h *ProviderHandler) GetProvider(c *gin.Context) {
	var provider database.Provider
	id := c.Param("id")
	if err := h.db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}
	c.JSON(http.StatusOK, provider)
}

// UpdateProvider updates an existing provider.
func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	var provider database.Provider
	id := c.Param("id")
	if err := h.db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&provider)
	c.JSON(http.StatusOK, provider)
}

// DeleteProvider deletes a provider.
func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.Provider{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete provider"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Provider deleted successfully"})
}

// GetProviderModels retrieves available models for a specific provider.
func (h *ProviderHandler) GetProviderModels(c *gin.Context) {
	id := c.Param("id")
	
	// Check if provider exists
	var provider database.Provider
	if err := h.db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}
	
	// For now, return a mock list of models based on provider type
	// In a real implementation, this would query the provider's API
	var models []string
	switch provider.Type {
	case "openai":
		models = []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo", "gpt-3.5-turbo-16k"}
	case "anthropic":
		models = []string{"claude-3-opus", "claude-3-sonnet", "claude-3-haiku"}
	case "gemini":
		models = []string{"gemini-pro", "gemini-pro-vision"}
	default:
		models = []string{"default-model"}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"models":       models,
		"providerName": provider.Name,
	})
}

// ImportProviderModels imports models for a specific provider.
func (h *ProviderHandler) ImportProviderModels(c *gin.Context) {
	id := c.Param("id")
	
	// Check if provider exists
	var provider database.Provider
	if err := h.db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}
	
	var req struct {
		ModelNames []string `json:"modelNames"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Create models if they don't exist
	var createdCount int
	for _, modelName := range req.ModelNames {
		var existingModel database.Model
		if err := h.db.Where("name = ?", modelName).First(&existingModel).Error; err != nil {
			// Model doesn't exist, create it
			newModel := database.Model{
				Name:     modelName,
				Remark:   fmt.Sprintf("Imported from %s", provider.Name),
				MaxRetry: 3,
				Timeout:  30,
				Enabled:  true,
			}
			if err := h.db.Create(&newModel).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create model %s", modelName)})
				return
			}
			createdCount++
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Models imported successfully",
		"importedCount": createdCount,
	})
}