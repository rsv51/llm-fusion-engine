package core

import "llm-fusion-engine/internal/database"

// ProviderRouteResult defines the result of a routing decision.
type ProviderRouteResult struct {
	Group         *database.Group
	ApiKey        string
	ResolvedModel string
}

// IProviderRouter is responsible for routing a request to the appropriate provider group.
type IProviderRouter interface {
	// RouteRequestAsync selects a provider group based on the model, proxy key, and other strategies.
	RouteRequestAsync(model, proxyKey string, excludedGroups []string) (*ProviderRouteResult, error)
}

// IKeyManager manages the API keys for different provider groups.
type IKeyManager interface {
	// GetNextKeyAsync retrieves the next available and healthy API key for a given group.
	GetNextKeyAsync(groupID uint) (string, error)
	// ValidateProxyKeyAsync checks if a proxy key is valid and returns it.
	ValidateProxyKeyAsync(proxyKey string) (*database.ProxyKey, error)
}

// IProvider represents a specific LLM provider (e.g., OpenAI, Anthropic).
type IProvider interface {
	// TODO: Define methods for making chat completion requests.
}

import "net/http"

// IProviderFactory creates instances of IProvider.
type IProviderFactory interface {
	// GetProvider gets a provider instance by its type name (e.g., "openai").
	GetProvider(providerType string) (IProvider, error)
}

// IMultiProviderService coordinates the routing and execution of requests.
type IMultiProviderService interface {
	ProcessChatCompletionHttpAsync(
		requestBody map[string]interface{},
		proxyKey string,
	) (*http.Response, error)
}