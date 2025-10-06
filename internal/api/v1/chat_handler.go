package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"llm-fusion-engine/internal/core"
	"llm-fusion-engine/internal/util"
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
	resp, err := h.service.ProcessChatCompletionHttpAsync(c, requestBody, proxyKey)
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

	// Create a TeeReader to capture the response body while streaming
	teeReader := util.NewTeeReader(resp.Body)
	resp.Body = teeReader

	if isStreaming {
		// Use TransparentStreamingActionResult for streaming
		actionResult := NewTransparentStreamingActionResult(resp)
		actionResult.ExecuteResultAsync(c)
	} else {
		// For non-streaming responses, copy headers and body
		defer resp.Body.Close()
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}

	// After the response has been sent, parse the captured body for token usage
	go func() {
		bodyBytes := teeReader.GetContent()
		var promptTokens, completionTokens, totalTokens int

		// This is a simplified parsing logic. A more robust solution would handle
		// different SSE formats and potential JSON errors.
		if isStreaming {
			// For streaming, find the last `data:` block that contains "usage"
			lines := strings.Split(string(bodyBytes), "\n")
			for i := len(lines) - 1; i >= 0; i-- {
				if strings.HasPrefix(lines[i], "data:") {
					jsonData := strings.TrimPrefix(lines[i], "data: ")
					if strings.Contains(jsonData, `"usage"`) {
						var usageEvent struct {
							Usage struct {
								PromptTokens     int `json:"prompt_tokens"`
								CompletionTokens int `json:"completion_tokens"`
								TotalTokens      int `json:"total_tokens"`
							} `json:"usage"`
						}
						if json.Unmarshal([]byte(jsonData), &usageEvent) == nil {
							promptTokens = usageEvent.Usage.PromptTokens
							completionTokens = usageEvent.Usage.CompletionTokens
							totalTokens = usageEvent.Usage.TotalTokens
							break // Found the usage, stop searching
						}
					}
				}
			}
		} else {
			// For non-streaming, unmarshal the whole body
			var responseData struct {
				Usage struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				} `json:"usage"`
			}
			if json.Unmarshal(bodyBytes, &responseData) == nil {
				promptTokens = responseData.Usage.PromptTokens
				completionTokens = responseData.Usage.CompletionTokens
				totalTokens = responseData.Usage.TotalTokens
			}
		}

		// Get the unique request ID from the context
		if requestID, exists := c.Get("requestID"); exists {
			// Update the log entry with token usage
			h.keyManager.UpdateLogTokens(requestID.(string), promptTokens, completionTokens, totalTokens)
		}
	}()
}