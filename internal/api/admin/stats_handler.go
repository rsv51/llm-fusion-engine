package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StatsHandler handles statistics endpoints.
type StatsHandler struct {
	db *gorm.DB
}

// NewStatsHandler creates a new StatsHandler.
func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

// GetStats returns overall system statistics.
func (h *StatsHandler) GetStats(c *gin.Context) {
	var groupCount int64
	var providerCount int64
	var keyCount int64
	var enabledGroups int64

	h.db.Model(&database.Group{}).Count(&groupCount)
	h.db.Model(&database.Provider{}).Count(&providerCount)
	h.db.Model(&database.ApiKey{}).Count(&keyCount)
	h.db.Model(&database.Group{}).Where("enabled = ?", true).Count(&enabledGroups)

	stats := gin.H{
		"total_groups":    groupCount,
		"total_providers": providerCount,
		"total_keys":      keyCount,
		"enabled_groups":  enabledGroups,
	}

	c.JSON(http.StatusOK, stats)
}