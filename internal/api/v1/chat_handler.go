package v1

import (
	"io"
	"llm-fusion-engine/internal/core"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles chat completion requests.
type ChatHandler struct {
	service core.IMultiProviderService
	keyManager core.IKeyManager
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(service core.IMultiProviderService, keyManager core.IKeyManager) *ChatHandler {
	return &ChatHandler{service: service, keyManager: keyManager}
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

	// 3. Validate proxy key
	if _, err := h.keyManager.ValidateProxyKeyAsync(proxyKey); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	// 4. Process the request
	resp, err := h.service.ProcessChatCompletionHttpAsync(requestBody, proxyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 5. Proxy the response
	isStreaming := false
	if stream, ok := requestBody["stream"].(bool); ok && stream {
		isStreaming = true
	}

	// Copy headers from the downstream response to the client response
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Status(resp.StatusCode)

	if isStreaming {
		// For streaming responses, we need to continuously read from the downstream
		// response body and write to the client's response writer.
		c.Stream(func(w io.Writer) bool {
			// Copy a chunk of data from the downstream response to the client
			_, err := io.CopyN(w, resp.Body, 2048) // Read in 2KB chunks
			if err != nil {
				// If we've reached the end of the stream, io.EOF will be reported.
				// In that case, we return false to stop streaming.
				return false
			}
			return true
		})
	} else {
		// For non-streaming responses, just copy the entire body at once.
		io.Copy(c.Writer, resp.Body)
	}
}