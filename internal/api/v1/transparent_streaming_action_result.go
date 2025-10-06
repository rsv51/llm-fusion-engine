package v1

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TransparentStreamingActionResult handles streaming responses transparently.
type TransparentStreamingActionResult struct {
	Response *http.Response
}

// NewTransparentStreamingActionResult creates a new TransparentStreamingActionResult.
func NewTransparentStreamingActionResult(resp *http.Response) *TransparentStreamingActionResult {
	return &TransparentStreamingActionResult{Response: resp}
}

// ExecuteResultAsync streams the response body to the client.
func (r *TransparentStreamingActionResult) ExecuteResultAsync(c *gin.Context) {
	defer r.Response.Body.Close()

	// Copy headers from the downstream response to the client response
	for key, values := range r.Response.Header {
		// Skip headers that are automatically handled or can cause issues
		if strings.EqualFold(key, "Transfer-Encoding") {
			continue
		}
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Ensure necessary headers for streaming are set
	if c.GetHeader("Cache-Control") == "" {
		c.Header("Cache-Control", "no-cache")
	}
	if c.GetHeader("Connection") == "" {
		c.Header("Connection", "keep-alive")
	}

	c.Status(r.Response.StatusCode)

	// Stream the body
	// Use a buffer to read from the response body chunk by chunk
	buffer := make([]byte, 4096) // 4KB buffer
	c.Stream(func(w io.Writer) bool {
		n, err := r.Response.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := w.Write(buffer[:n]); writeErr != nil {
				// Error writing to client, stop streaming
				return false
			}
		}

		if err != nil {
			// If EOF or any other error, stop streaming
			return false
		}

		// Continue streaming
		return true
	})
}