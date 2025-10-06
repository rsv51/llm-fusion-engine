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
	c.Stream(func(w io.Writer) bool {
		// Copy a chunk of data from the downstream response to the client
		_, err := io.Copy(w, r.Response.Body)
		// If err is not nil, it means the stream has ended or an error occurred.
		// In either case, we should stop streaming.
		return err == nil
	})
}