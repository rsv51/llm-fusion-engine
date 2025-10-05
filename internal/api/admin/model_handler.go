package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ModelHandler handles CRUD operations for models.
type ModelHandler struct {
	db *gorm.DB
}

// NewModelHandler creates a new ModelHandler.
func NewModelHandler(db *gorm.DB) *ModelHandler {
	return &ModelHandler{db: db}
}

// CreateModel creates a new model.
func (h *ModelHandler) CreateModel(c *gin.Context) {
	var model database.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model"})
		return
	}

	c.JSON(http.StatusOK, model)
}

// GetModels retrieves all models with pagination.
func (h *ModelHandler) GetModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	search := c.Query("search")

	var models []database.Model
	var total int64

	query := h.db.Model(&database.Model{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&models)

	c.JSON(http.StatusOK, gin.H{
		"items":      models,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetModel retrieves a single model by ID.
func (h *ModelHandler) GetModel(c *gin.Context) {
	var model database.Model
	id := c.Param("id")
	if err := h.db.First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}
	c.JSON(http.StatusOK, model)
}

// UpdateModel updates an existing model.
func (h *ModelHandler) UpdateModel(c *gin.Context) {
	var model database.Model
	id := c.Param("id")
	if err := h.db.First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&model)
	c.JSON(http.StatusOK, model)
}

// DeleteModel deletes a model.
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.Model{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete model"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Model deleted successfully"})
}

// CloneModel clones an existing model.
func (h *ModelHandler) CloneModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	var originalModel database.Model
	if err := h.db.First(&originalModel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	var req struct {
		NewName string `json:"newName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clonedModel := originalModel
	clonedModel.ID = 0 // GORM will treat this as a new record
	clonedModel.Name = req.NewName

	if err := h.db.Create(&clonedModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clone model"})
		return
	}

	c.JSON(http.StatusOK, clonedModel)
}