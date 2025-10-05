package admin

import (
	"encoding/json"
	"llm-fusion-engine/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GroupHandler handles CRUD operations for groups.
type GroupHandler struct {
	db *gorm.DB
}

// NewGroupHandler creates a new GroupHandler.
func NewGroupHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

// CreateGroup creates a new group.
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group database.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroups retrieves all groups with pagination response format.
func (h *GroupHandler) GetGroups(c *gin.Context) {
	var groups []database.Group
	if err := h.db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}
	
	// Return in pagination response format for frontend compatibility
	c.JSON(http.StatusOK, gin.H{
		"items":      groups,
		"total":      len(groups),
		"page":       1,
		"pageSize":   len(groups),
		"totalPages": 1,
	})
}

// GetGroup retrieves a single group by ID.
func (h *GroupHandler) GetGroup(c *gin.Context) {
	var group database.Group
	id := c.Param("id")
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
}

// UpdateGroup updates an existing group.
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	var group database.Group
	id := c.Param("id")
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&group)
	c.JSON(http.StatusOK, group)
}

// DeleteGroup deletes a group.
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&database.Group{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// GetModelAliases retrieves the model aliases for a group.
func (h *GroupHandler) GetModelAliases(c *gin.Context) {
	var group database.Group
	id := c.Param("id")
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var aliases map[string]string
	if group.ModelAliases != "" {
		if err := json.Unmarshal([]byte(group.ModelAliases), &aliases); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse model aliases"})
			return
		}
	}

	c.JSON(http.StatusOK, aliases)
}

// UpdateModelAliases updates the model aliases for a group.
func (h *GroupHandler) UpdateModelAliases(c *gin.Context) {
	var group database.Group
	id := c.Param("id")
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var req struct {
		Aliases map[string]string `json:"aliases"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aliasesJSON, err := json.Marshal(req.Aliases)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize model aliases"})
		return
	}

	group.ModelAliases = string(aliasesJSON)
	if err := h.db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model aliases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model aliases updated successfully"})
}