package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// KeyHandler handles CRUD operations for API keys.
type KeyHandler struct {
	db *gorm.DB
}

// NewKeyHandler creates a new KeyHandler.
func NewKeyHandler(db *gorm.DB) *KeyHandler {
	return &KeyHandler{db: db}
}

// CreateKey creates a new API key.
func (h *KeyHandler) CreateKey(c *gin.Context) {
	var key database.ApiKey
	if err := c.ShouldBindJSON(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create key"})
		return
	}

	c.JSON(http.StatusOK, key)
}

// GetKeys retrieves all API keys with pagination.
func (h *KeyHandler) GetKeys(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	offset := (page - 1) * pageSize
	
	var keys []database.ApiKey
	var total int64
	
	// Count total records
	if err := h.db.Model(&database.ApiKey{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count keys"})
		return
	}
	
	// Get paginated records
	if err := h.db.Offset(offset).Limit(pageSize).Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve keys"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": keys,
		"pagination": gin.H{
			"page":      page,
			"pageSize":  pageSize,
			"total":     total,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetKey retrieves a single API key by ID.
func (h *KeyHandler) GetKey(c *gin.Context) {
	var key database.ApiKey
	id := c.Param("id")
	if err := h.db.First(&key, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}
	c.JSON(http.StatusOK, key)
}

// UpdateKey updates an existing API key.
func (h *KeyHandler) UpdateKey(c *gin.Context) {
	var key database.ApiKey
	id := c.Param("id")
	if err := h.db.First(&key, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	if err := c.ShouldBindJSON(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&key)
	c.JSON(http.StatusOK, key)
}

// DeleteKey deletes an API key.
func (h *KeyHandler) DeleteKey(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.ApiKey{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Key deleted successfully"})
}