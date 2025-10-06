package admin

import (
	"context"
	"fmt"
	"llm-fusion-engine/internal/database"
	"llm-fusion-engine/internal/providers"
	"net/http"
	"strconv"
	"time"

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

// DeleteProvider deletes a provider and its associated model mappings.
func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	
	// Start a transaction to ensure both operations succeed or fail together
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// First delete all associated ModelProviderMappings
	if err := tx.Where("provider_id = ?", id).Delete(&database.ModelProviderMapping{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated model mappings"})
		return
	}
	
	// Then delete the provider
	if err := tx.Delete(&database.Provider{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete provider"})
		return
	}
	
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Provider and associated model mappings deleted successfully"})
}

// GetProviderModels retrieves available models for a specific provider.
// It attempts to fetch models from the provider's actual API.
// If that fails, it falls back to a default list.
func (h *ProviderHandler) GetProviderModels(c *gin.Context) {
	id := c.Param("id")
	
	// Check if provider exists
	var provider database.Provider
	if err := h.db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}
	
	var models []string
	var fetchError error
	
	// Try to create a client and fetch real models from the provider's API
	client, err := providers.CreateClient(provider.Type, provider.Config)
	if err == nil {
		// Validate config before attempting to fetch
		if err := client.ValidateConfig(); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			// Fetch real models from the provider's API
			models, fetchError = client.GetModels(ctx)
		} else {
			fetchError = err
		}
	} else {
		fetchError = err
	}
	
	// If fetching from API failed, fall back to default list
	if fetchError != nil || len(models) == 0 {
		models = getDefaultModels(provider.Type)
		
		// Log the error for debugging but don't fail the request
		if fetchError != nil {
			c.JSON(http.StatusOK, gin.H{
				"models":       models,
				"providerName": provider.Name,
				"warning":      fmt.Sprintf("Failed to fetch from provider API, using default list: %v", fetchError),
				"source":       "default",
			})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"models":       models,
		"providerName": provider.Name,
		"source":       "api",
	})
}

// getDefaultModels returns a default list of models for a given provider type
func getDefaultModels(providerType string) []string {
	switch providerType {
	case "openai":
		return []string{
			"gpt-4", "gpt-4-turbo", "gpt-4o", "gpt-4o-mini",
			"gpt-3.5-turbo", "gpt-3.5-turbo-16k",
			"o1-preview", "o1-mini",
		}
	case "anthropic":
		return []string{
			"claude-3-opus-20240229",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
			"claude-3-5-sonnet-20240620",
			"claude-2.1", "claude-2.0",
		}
	case "gemini":
		return []string{
			"gemini-pro",
			"gemini-pro-vision",
			"gemini-1.5-pro-latest",
			"gemini-1.5-flash-latest",
			"gemini-ultra",
		}
	default:
		return []string{}
	}
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