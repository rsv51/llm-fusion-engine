package admin

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LogHandler handles request log queries.
type LogHandler struct {
	db *gorm.DB
}

// NewLogHandler creates a new LogHandler.
func NewLogHandler(db *gorm.DB) *LogHandler {
	return &LogHandler{db: db}
}

// GetLogs retrieves request logs with pagination and filtering.
func (h *LogHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	offset := (page - 1) * pageSize
	
	var logs []database.Log
	var total int64

	query := h.db.Model(&database.Log{})

	// Apply filters if provided
	if model := c.Query("model"); model != "" {
		query = query.Where("model = ?", model)
	}
	if status := c.Query("status"); status != "" {
		statusCode, _ := strconv.Atoi(status)
		query = query.Where("status_code = ?", statusCode)
	}
	if groupID := c.Query("group_id"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	
	// Count total records
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count logs"})
		return
	}
	
	// Get paginated records, ordered by creation time descending
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"page":      page,
			"pageSize":  pageSize,
			"total":     total,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetLog retrieves a single request log by ID.
func (h *LogHandler) GetLog(c *gin.Context) {
	var log database.Log
	id := c.Param("id")
	if err := h.db.First(&log, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}
	c.JSON(http.StatusOK, log)
}

// DeleteLogs deletes logs based on filters (e.g., older than X days).
func (h *LogHandler) DeleteLogs(c *gin.Context) {
	// Optional: implement bulk deletion based on date range
	daysOld, _ := strconv.Atoi(c.DefaultQuery("daysOld", "30"))
	
	// Delete logs older than specified days
	// Note: The date function might vary depending on the SQL dialect.
	// This uses a generic approach that might need adjustment for specific databases.
	// For SQLite, it would be something like: "timestamp < date('now', '-' || ? || ' day')"
	// For simplicity, we'll stick to a more direct time comparison.
	deleteTime := time.Now().AddDate(0, 0, -daysOld)
	result := h.db.Where("timestamp < ?", deleteTime).Delete(&database.Log{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete logs"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Logs deleted successfully",
		"deleted": result.RowsAffected,
	})
}