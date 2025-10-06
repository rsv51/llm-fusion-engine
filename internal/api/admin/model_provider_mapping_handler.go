package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ModelProviderMappingHandler handles CRUD operations for Models and their Provider Mappings.
type ModelProviderMappingHandler struct {
	db *gorm.DB
}

// NewModelProviderMappingHandler creates a new ModelProviderMappingHandler.
func NewModelProviderMappingHandler(db *gorm.DB) *ModelProviderMappingHandler {
	return &ModelProviderMappingHandler{db: db}
}

// --- Model CRUD ---

// CreateModel creates a new model definition.
func (h *ModelProviderMappingHandler) CreateModel(c *gin.Context) {
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

// GetModels retrieves all model definitions with pagination.
func (h *ModelProviderMappingHandler) GetModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var models []database.Model
	var total int64

	if err := h.db.Model(&database.Model{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count models"})
		return
	}

	if err := h.db.Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve models"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": models,
		"pagination": gin.H{
			"page":      page,
			"pageSize":  pageSize,
			"total":     total,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetModel retrieves a single model definition by ID.
func (h *ModelProviderMappingHandler) GetModel(c *gin.Context) {
	var model database.Model
	id := c.Param("id")
	if err := h.db.First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}
	c.JSON(http.StatusOK, model)
}

// UpdateModel updates an existing model definition.
func (h *ModelProviderMappingHandler) UpdateModel(c *gin.Context) {
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

// DeleteModel deletes a model definition.
// Note: This should also handle cascading deletes for ModelProviderMappings if not handled by DB constraint.
func (h *ModelProviderMappingHandler) DeleteModel(c *gin.Context) {
	id := c.Param("id")
	// Manually delete related mappings first if DB doesn't handle ON DELETE CASCADE well or for more control
	if err := h.db.Where("model_id = ?", id).Delete(&database.ModelProviderMapping{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related provider mappings for model"})
		return
	}
	if err := h.db.Delete(&database.Model{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete model"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Model and its mappings deleted successfully"})
}

// --- ModelProviderMapping CRUD ---

// CreateModelProviderMapping creates a new mapping between a model and a provider.
func (h *ModelProviderMappingHandler) CreateModelProviderMapping(c *gin.Context) {
	var mapping database.ModelProviderMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that ModelID and ProviderID exist
	var modelCount int64
	if h.db.Model(&database.Model{}).Where("id = ?", mapping.ModelID).Count(&modelCount); modelCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Referenced ModelID does not exist"})
		return
	}
	var providerCount int64
	if h.db.Model(&database.Provider{}).Where("id = ?", mapping.ProviderID).Count(&providerCount); providerCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Referenced ProviderID does not exist"})
		return
	}

	if err := h.db.Create(&mapping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model provider mapping"})
		return
	}

	// Preload related data for the response
	h.db.Preload("Model").Preload("Provider").First(&mapping, mapping.ID)
	c.JSON(http.StatusOK, mapping)
}

// GetModelProviderMappings retrieves all model-provider mappings with pagination.
func (h *ModelProviderMappingHandler) GetModelProviderMappings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var mappings []database.ModelProviderMapping
	var total int64

	if err := h.db.Model(&database.ModelProviderMapping{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count model provider mappings"})
		return
	}

	if err := h.db.Offset(offset).Limit(pageSize).Preload("Model").Preload("Provider").Find(&mappings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve model provider mappings"})
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

// GetModelProviderMapping retrieves a single model-provider mapping by ID.
func (h *ModelProviderMappingHandler) GetModelProviderMapping(c *gin.Context) {
	var mapping database.ModelProviderMapping
	id := c.Param("id")
	if err := h.db.Preload("Model").Preload("Provider").First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model provider mapping not found"})
		return
	}
	c.JSON(http.StatusOK, mapping)
}

// UpdateModelProviderMapping updates an existing model-provider mapping.
func (h *ModelProviderMappingHandler) UpdateModelProviderMapping(c *gin.Context) {
	var mapping database.ModelProviderMapping
	id := c.Param("id")
	if err := h.db.First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model provider mapping not found"})
		return
	}

	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&mapping)
	// Preload related data for the response
	h.db.Preload("Model").Preload("Provider").First(&mapping, mapping.ID)
	c.JSON(http.StatusOK, mapping)
}

// DeleteModelProviderMapping deletes a model-provider mapping.
func (h *ModelProviderMappingHandler) DeleteModelProviderMapping(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.ModelProviderMapping{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete model provider mapping"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Model provider mapping deleted successfully"})
}

// GetMappingHealthStatus retrieves the recent health status for a model-provider mapping.
func (h *ModelProviderMappingHandler) GetMappingHealthStatus(c *gin.Context) {
	id := c.Param("id")
	
	var mapping database.ModelProviderMapping
	if err := h.db.First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model provider mapping not found"})
		return
	}

	// Get the last 10 request logs for this mapping
	var logs []database.RequestLog
	err := h.db.Where("provider_id = ? AND model = ?", mapping.ProviderID, mapping.ProviderModel).
		Order("created_at DESC").
		Limit(10).
		Find(&logs).Error
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve health status"})
		return
	}

	// Process logs to create health status array
	healthStatus := make([]gin.H, 0, len(logs))
	for _, log := range logs {
		status := "success"
		if log.StatusCode >= 400 {
			status = "error"
		}
		healthStatus = append(healthStatus, gin.H{
			"timestamp":  log.CreatedAt,
			"status":     status,
			"statusCode": log.StatusCode,
			"latencyMs":  log.LatencyMs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"mappingId":    mapping.ID,
		"healthStatus": healthStatus,
	})
}

// GetAllMappingsHealthStatus retrieves health status for all mappings.
func (h *ModelProviderMappingHandler) GetAllMappingsHealthStatus(c *gin.Context) {
	var mappings []database.ModelProviderMapping
	if err := h.db.Find(&mappings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mappings"})
		return
	}

	result := make(map[uint][]gin.H)
	
	for _, mapping := range mappings {
		var logs []database.RequestLog
		err := h.db.Where("provider_id = ? AND model = ?", mapping.ProviderID, mapping.ProviderModel).
			Order("created_at DESC").
			Limit(10).
			Find(&logs).Error
		
		if err != nil {
			continue
		}

		healthStatus := make([]gin.H, 0, len(logs))
		for _, log := range logs {
			status := "success"
			if log.StatusCode >= 400 {
				status = "error"
			}
			healthStatus = append(healthStatus, gin.H{
				"timestamp":  log.CreatedAt,
				"status":     status,
				"statusCode": log.StatusCode,
			})
		}
		
		result[mapping.ID] = healthStatus
	}

	c.JSON(http.StatusOK, result)
}