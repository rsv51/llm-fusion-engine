package v1

import (
	"llm-fusion-engine/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModelHandler struct {
	db *gorm.DB
}

func NewModelHandler(db *gorm.DB) *ModelHandler {
	return &ModelHandler{db: db}
}

// GetModels returns a list of available models.
// This is a public endpoint and does not require authentication.
func (h *ModelHandler) GetModels(c *gin.Context) {
	var models []database.Model
	if err := h.db.Find(&models).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve models"})
		return
	}

	// Format the response to be compatible with OpenAI's API
	type ModelInfo struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	}

	var responseData []ModelInfo
	for _, model := range models {
		responseData = append(responseData, ModelInfo{
			ID:      model.Name,
			Object:  "model",
			Created: model.CreatedAt.Unix(),
			OwnedBy: "llm-fusion-engine",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   responseData,
	})
}
