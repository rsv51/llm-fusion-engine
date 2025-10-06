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
	h.db.Model(&database.Log{}).Where("timestamp > ?", since).Count(&totalRequests)

	// Get successful requests
	var successRequests int64
	h.db.Model(&database.Log{}).Where("timestamp > ? AND response_status >= 200 AND response_status < 300", since).Count(&successRequests)

	// Calculate success rate
	var successRate float64
	if totalRequests > 0 {
		successRate = (float64(successRequests) / float64(totalRequests)) * 100
	}

	// Get average response time
	var avgResponseTimeMs float64
	h.db.Model(&database.Log{}).Where("timestamp > ?", since).Select("avg(latency)").Row().Scan(&avgResponseTimeMs)

	// Get active keys (keys used in the last 24 hours)
	var activeKeys int64
	h.db.Model(&database.Log{}).Where("timestamp > ?", since).Distinct("proxy_key").Count(&activeKeys)

	// Get provider-specific stats
	type ProviderStatsResult struct {
		Provider        string  `json:"provider"`
		RequestCount    int64   `json:"requestCount"`
		SuccessCount    int64   `json:"successCount"`
		ErrorCount      int64   `json:"errorCount"`
		AvgResponseTime float64 `json:"avgResponseTimeMs"`
	}
	var providerStats []ProviderStatsResult
	h.db.Model(&database.Log{}).
		Select("provider, COUNT(*) as request_count, "+
			"SUM(CASE WHEN response_status >= 200 AND response_status < 300 THEN 1 ELSE 0 END) as success_count, "+
			"SUM(CASE WHEN response_status >= 400 THEN 1 ELSE 0 END) as error_count, "+
			"AVG(latency) as avg_response_time").
		Where("timestamp > ?", since).
		Group("provider").
		Scan(&providerStats)

	// Populate provider names
	// The provider name is now directly in the stats, so no need for a second query

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