package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProxyKeyHandler handles CRUD operations for proxy keys.
type ProxyKeyHandler struct {
	db *gorm.DB
}

// NewProxyKeyHandler creates a new ProxyKeyHandler.
func NewProxyKeyHandler(db *gorm.DB) *ProxyKeyHandler {
	return &ProxyKeyHandler{db: db}
}

// CreateProxyKey creates a new proxy key.
func (h *ProxyKeyHandler) CreateProxyKey(c *gin.Context) {
	var proxyKey database.ProxyKey
	if err := c.ShouldBindJSON(&proxyKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&proxyKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy key"})
		return
	}

	c.JSON(http.StatusOK, proxyKey)
}

// GetProxyKeys retrieves all proxy keys with pagination.
func (h *ProxyKeyHandler) GetProxyKeys(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var proxyKeys []database.ProxyKey
	var total int64

	// Count total records
	if err := h.db.Model(&database.ProxyKey{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count proxy keys"})
		return
	}

	// Get paginated records
	if err := h.db.Offset(offset).Limit(pageSize).Find(&proxyKeys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve proxy keys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": proxyKeys,
		"pagination": gin.H{
			"page":      page,
			"pageSize":  pageSize,
			"total":     total,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetProxyKey retrieves a single proxy key by ID.
func (h *ProxyKeyHandler) GetProxyKey(c *gin.Context) {
	var proxyKey database.ProxyKey
	id := c.Param("id")
	if err := h.db.First(&proxyKey, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proxy key not found"})
		return
	}
	c.JSON(http.StatusOK, proxyKey)
}

// UpdateProxyKey updates an existing proxy key.
func (h *ProxyKeyHandler) UpdateProxyKey(c *gin.Context) {
	var proxyKey database.ProxyKey
	id := c.Param("id")
	if err := h.db.First(&proxyKey, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proxy key not found"})
		return
	}

	if err := c.ShouldBindJSON(&proxyKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&proxyKey)
	c.JSON(http.StatusOK, proxyKey)
}

// DeleteProxyKey deletes a proxy key.
func (h *ProxyKeyHandler) DeleteProxyKey(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.ProxyKey{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete proxy key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Proxy key deleted successfully"})
}