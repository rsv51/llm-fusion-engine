package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"llm-fusion-engine/internal/core"
	"net/http"
)

// MultiProviderService coordinates routing and API requests.
type MultiProviderService struct {
	router         core.IProviderRouter
	providerFactory core.IProviderFactory
}

// NewMultiProviderService creates a new MultiProviderService.
func NewMultiProviderService(router core.IProviderRouter, factory core.IProviderFactory) *MultiProviderService {
	return &MultiProviderService{
		router:         router,
		providerFactory: factory,
	}
}

// ProcessChatCompletionHttpAsync handles the chat completion request.
func (s *MultiProviderService) ProcessChatCompletionHttpAsync(
	requestBody map[string]interface{},
	proxyKey string,
) (*http.Response, error) {
	model, ok := requestBody["model"].(string)
	if !ok {
		return nil, errors.New("model not specified in request")
	}

	// 1. Route the request
	routeResult, err := s.router.RouteRequestAsync(model, proxyKey, nil)
	if err != nil {
		return nil, err
	}

	// 2. Get the provider instance
	// We need to get the provider type from the selected group.
	// This assumes a group has a field like `ProviderType`.
	// Let's assume the first provider in the group determines the type.
	if len(routeResult.Group.Providers) == 0 {
		return nil, errors.New("selected group has no providers")
	}
	// providerType := routeResult.Group.Providers[0].ProviderType
	// TODO: Use provider factory when implementing actual provider logic
	// provider, err := s.providerFactory.GetProvider(providerType)
	// if err != nil {
	// 	return nil, err
	// }

	// 3. Prepare and forward the request
	// The actual implementation of this will be in the provider instance.
	// For now, we'll just simulate it.
	// TODO: Replace this with actual provider.ProcessRequest(...)
	
	// Modify the request with the resolved model
	requestBody["model"] = routeResult.ResolvedModel
	
	// Create a new HTTP request to the downstream service
	downstreamURL := "http://localhost:8080/mock" // This should be the actual provider endpoint
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", downstreamURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set headers, including the real API key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+routeResult.ApiKey)

	// Execute the request
	client := &http.Client{}
	return client.Do(req)
}