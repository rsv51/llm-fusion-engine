package admin

import (
	"llm-fusion-engine/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check operations.
type HealthHandler struct {
	db *gorm.DB
	healthChecker *services.HealthChecker
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *gorm.DB, healthChecker *services.HealthChecker) *HealthHandler {
	return &HealthHandler{db: db, healthChecker: healthChecker}
}

// CheckProviderHealth handles a health check for a single provider.
func (h *HealthHandler) CheckProviderHealth(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider ID"})
		return
	}

	provider, err := h.healthChecker.CheckProvider(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check provider health"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// CheckAllProvidersHealth triggers a health check for all providers.
func (h *HealthHandler) CheckAllProvidersHealth(c *gin.Context) {
	go h.healthChecker.CheckAllProviders()
	c.JSON(http.StatusAccepted, gin.H{"message": "Health check for all providers has been initiated."})
}