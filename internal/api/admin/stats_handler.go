package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StatsHandler handles statistics endpoints.
type StatsHandler struct {
	db        *gorm.DB
	startTime time.Time
}

// NewStatsHandler creates a new StatsHandler.
func NewStatsHandler(db *gorm.DB, startTime time.Time) *StatsHandler {
	return &StatsHandler{db: db, startTime: startTime}
}

// GetStats returns overall system statistics for the last 24 hours.
func (h *StatsHandler) GetStats(c *gin.Context) {
	// Define the time window for the stats (last 24 hours)
	since := time.Now().Add(-24 * time.Hour)

	// Get total requests
	var totalRequests int64
	h.db.Model(&database.RequestLog{}).Where("created_at > ?", since).Count(&totalRequests)

	// Get successful requests
	var successRequests int64
	h.db.Model(&database.RequestLog{}).Where("created_at > ? AND status_code >= 200 AND status_code < 300", since).Count(&successRequests)

	// Calculate success rate
	var successRate float64
	if totalRequests > 0 {
		successRate = (float64(successRequests) / float64(totalRequests)) * 100
	}

	// Get average response time
	var avgResponseTimeMs float64
	h.db.Model(&database.RequestLog{}).Where("created_at > ?", since).Select("avg(latency_ms)").Row().Scan(&avgResponseTimeMs)

	// Get active keys (keys used in the last 24 hours)
	var activeKeys int64
	h.db.Model(&database.RequestLog{}).Where("created_at > ?", since).Distinct("provider_id").Count(&activeKeys)

	// Get provider-specific stats
	type ProviderStatsResult struct {
		ProviderID      uint
		ProviderName    string
		RequestCount    int64
		SuccessCount    int64
		AvgResponseTime float64
	}
	var providerStats []ProviderStatsResult
	h.db.Model(&database.RequestLog{}).
		Select("provider_id, COUNT(*) as request_count, SUM(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 ELSE 0 END) as success_count, AVG(latency_ms) as avg_response_time").
		Where("created_at > ?", since).
		Group("provider_id").
		Scan(&providerStats)

	// Populate provider names
	for i, p := range providerStats {
		var provider database.Provider
		if err := h.db.First(&provider, p.ProviderID).Error; err == nil {
			providerStats[i].ProviderName = provider.Type
		}
	}

	stats := gin.H{
		"totalRequests":     totalRequests,
		"successRate":       successRate,
		"avgResponseTimeMs": avgResponseTimeMs,
		"activeKeys":        activeKeys,
		"providers":         providerStats,
		"startTime":         h.startTime,
	}

	c.JSON(http.StatusOK, stats)
}