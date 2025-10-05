package v1

import (
	"encoding/json"
	"io"
	"llm-fusion-engine/internal/core"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles chat completion requests.
type ChatHandler struct {
	service core.IMultiProviderService
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(service core.IMultiProviderService) *ChatHandler {
	return &ChatHandler{service: service}
}

// ChatCompletions is the handler for the /v1/chat/completions endpoint.
func (h *ChatHandler) ChatCompletions(c *gin.Context) {
	// 1. Parse request body
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. Extract proxy key
	authHeader := c.GetHeader("Authorization")
	proxyKey := strings.TrimPrefix(authHeader, "Bearer ")
	if proxyKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	// 3. Process the request
	resp, err := h.service.ProcessChatCompletionHttpAsync(requestBody, proxyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 4. Proxy the response
	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Copy status code
	c.Status(resp.StatusCode)

	// Copy body
	io.Copy(c.Writer, resp.Body)
}