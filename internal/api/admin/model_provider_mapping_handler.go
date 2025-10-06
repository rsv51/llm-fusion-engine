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

	query := h.db.Model(&database.ModelProviderMapping{})

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count model provider mappings"})
		return
	}

	// Get paginated records
	if err := query.Offset(offset).Limit(pageSize).Preload("Model").Preload("Provider").Find(&mappings).Error; err != nil {
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
	if err := h.db.Preload("Provider").Preload("Model").First(&mapping, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model provider mapping not found"})
		return
	}

	// Define a struct to hold the aggregation results
	type HealthStats struct {
		TotalRequests   int64   `json:"totalRequests"`
		Successful      int64   `json:"successful"`
		AverageLatency  float64 `json:"averageLatency"`
	}

	var stats HealthStats
	
	// Query the last 20 logs for this specific mapping to calculate health
	h.db.Model(&database.Log{}).
		Select("count(*) as total_requests, sum(case when is_success = 1 then 1 else 0 end) as successful, avg(latency) as average_latency").
		Where("provider = ? AND model = ?", mapping.Provider.Name, mapping.ProviderModel).
		Order("timestamp DESC").
		Limit(20).
		Scan(&stats)

	var successRate float64
	if stats.TotalRequests > 0 {
		successRate = (float64(stats.Successful) / float64(stats.TotalRequests)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"mappingId":      mapping.ID,
		"totalRequests":  stats.TotalRequests,
		"successRate":    successRate,
		"averageLatency": stats.AverageLatency,
	})
}

// GetAllMappingsHealthStatus retrieves aggregated health status for all mappings.
func (h *ModelProviderMappingHandler) GetAllMappingsHealthStatus(c *gin.Context) {
	type MappingHealthStats struct {
		MappingID        uint    `json:"mappingId"`
		TotalRequests    int64   `json:"totalRequests"`
		SuccessRate      float64 `json:"successRate"`
		AverageLatency   float64 `json:"averageLatency"`
	}

	// This query is more complex. It calculates stats for each provider/model pair from the logs.
	var results []struct {
		Provider         string
		Model            string
		TotalRequests    int64
		Successful       int64
		AverageLatency   float64
	}
	
	// We get the stats for all provider/model pairs found in the logs
	h.db.Model(&database.Log{}).
		Select("provider, model, count(*) as total_requests, sum(case when is_success = 1 then 1 else 0 end) as successful, avg(latency) as average_latency").
		Group("provider, model").
		Scan(&results)

	// Now, we need to map these stats back to our ModelProviderMapping IDs
	var mappings []database.ModelProviderMapping
	h.db.Preload("Provider").Find(&mappings)

	statsMap := make(map[string]MappingHealthStats)
	for _, r := range results {
		key := r.Provider + "::" + r.Model
		var successRate float64
		if r.TotalRequests > 0 {
			successRate = (float64(r.Successful) / float64(r.TotalRequests)) * 100
		}
		statsMap[key] = MappingHealthStats{
			TotalRequests:  r.TotalRequests,
			SuccessRate:    successRate,
			AverageLatency: r.AverageLatency,
		}
	}

	finalResult := make(map[uint]MappingHealthStats)
	for _, m := range mappings {
		key := m.Provider.Name + "::" + m.ProviderModel
		if stats, ok := statsMap[key]; ok {
			stats.MappingID = m.ID
			finalResult[m.ID] = stats
		} else {
			// If no logs found for this mapping, return zero stats
			finalResult[m.ID] = MappingHealthStats{MappingID: m.ID}
		}
	}

	c.JSON(http.StatusOK, finalResult)
}