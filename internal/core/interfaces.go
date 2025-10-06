package core

import (
	"llm-fusion-engine/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ProviderRouteResult defines the result of a routing decision.
type ProviderRouteResult struct {
	Group         *database.Group
	Provider      *database.Provider
	ApiKey        string
	ResolvedModel string
	RetryCount    int
	RetryAfter    time.Duration
}

// IProviderRouter is responsible for routing a request to the appropriate provider group.
type IProviderRouter interface {
	// RouteRequestAsync selects a provider based on the model, proxy key, and failover/load-balancing strategies.
	RouteRequestAsync(model, proxyKey string, excludedProviders []uint) (*ProviderRouteResult, error)
}

// IKeyManager manages the API keys for different provider groups.
type IKeyManager interface {
	// ValidateProxyKeyAsync checks if a proxy key is valid and returns it.
	ValidateProxyKeyAsync(proxyKey string) (*database.ProxyKey, error)
	UpdateLogTokens(requestID string, promptTokens, completionTokens, totalTokens int)
}

// IProvider represents a specific LLM provider (e.g., OpenAI, Anthropic).
type IProvider interface {
	// TODO: Define methods for making chat completion requests.
}

// IProviderFactory creates instances of IProvider.
type IProviderFactory interface {
	// GetProvider gets a provider instance by its type name (e.g., "openai").
	GetProvider(providerType string) (IProvider, error)
}

// IMultiProviderService coordinates the routing and execution of requests.
type IMultiProviderService interface {
	ProcessChatCompletionHttpAsync(
		c *gin.Context,
		requestBody map[string]interface{},
		proxyKey string,
	) (*http.Response, error)
	LogRequest(
		requestID string,
		requestBody map[string]interface{},
		proxyKey string,
		providerName string,
		requestUrl string,
		response *http.Response,
		isSuccess bool,
		latency time.Duration,
		promptTokens int,
		completionTokens int,
		totalTokens int,
	)
}