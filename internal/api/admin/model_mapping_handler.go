package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ModelMappingHandler handles CRUD operations for model mappings.
type ModelMappingHandler struct {
	db *gorm.DB
}

// NewModelMappingHandler creates a new ModelMappingHandler.
func NewModelMappingHandler(db *gorm.DB) *ModelMappingHandler {
	return &ModelMappingHandler{db: db}
}

// CreateModelMapping creates a new model mapping.
func (h *ModelMappingHandler) CreateModelMapping(c *gin.Context) {
	var mapping database.ModelMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&mapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model mapping"})
		return
	}

	c.JSON(http.StatusOK, mapping)
}

// GetModelMappings retrieves all model mappings with pagination.
func (h *ModelMappingHandler) GetModelMappings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var mappings []database.ModelMapping
	var total int64

	if err := h.db.Model(&database.ModelMapping{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count model mappings"})
		return
	}

	if err := h.db.Offset(offset).Limit(pageSize).Find(&mappings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve model mappings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mappings,
		"pagination": gin.H{
			"page":      page,
			"pageSize":  pageSize,
			"total":     total,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetModelMapping retrieves a single model mapping by ID.
func (h *ModelMappingHandler) GetModelMapping(c *gin.Context) {
	var mapping database.ModelMapping
	id := c.Param("id")
	if err := h.db.Preload("Provider").First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model mapping not found"})
		return
	}
	c.JSON(http.StatusOK, mapping)
}

// UpdateModelMapping updates an existing model mapping.
func (h *ModelMappingHandler) UpdateModelMapping(c *gin.Context) {
	var mapping database.ModelMapping
	id := c.Param("id")
	if err := h.db.First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model mapping not found"})
		return
	}

	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&mapping)
	c.JSON(http.StatusOK, mapping)
}

// DeleteModelMapping deletes a model mapping.
func (h *ModelMappingHandler) DeleteModelMapping(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.ModelMapping{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete model mapping"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Model mapping deleted successfully"})
}