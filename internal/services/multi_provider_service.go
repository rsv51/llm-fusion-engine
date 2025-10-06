package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"llm-fusion-engine/internal/core"
	"llm-fusion-engine/internal/database"
	"net/http"
	"gorm.io/gorm"
)

// MultiProviderService coordinates routing and API requests.
type MultiProviderService struct {
	router         core.IProviderRouter
	providerFactory core.IProviderFactory
	db             *gorm.DB
}

// NewMultiProviderService creates a new MultiProviderService.
func NewMultiProviderService(router core.IProviderRouter, factory core.IProviderFactory, db *gorm.DB) *MultiProviderService {
	return &MultiProviderService{
		router:         router,
		providerFactory: factory,
		db:             db,
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
	// Since we removed the direct relationship, we need to query the database.
	var providers []database.Provider
	if err := s.db.Find(&providers).Error; err != nil {
		return nil, errors.New("failed to retrieve providers")
	}
	
	if len(providers) == 0 {
		return nil, errors.New("no providers available")
	}
	
	// For now, use the first provider's type
	// TODO: Implement proper provider selection logic based on the group
	providerType := providers[0].Type
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