package admin

import (
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

// GetGroups retrieves all groups.
func (h *GroupHandler) GetGroups(c *gin.Context) {
	var groups []database.Group
	if err := h.db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}
	c.JSON(http.StatusOK, groups)
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