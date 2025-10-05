package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"

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

// GetProviders retrieves all providers.
func (h *ProviderHandler) GetProviders(c *gin.Context) {
	var providers []database.Provider
	if err := h.db.Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve providers"})
		return
	}
	c.JSON(http.StatusOK, providers)
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