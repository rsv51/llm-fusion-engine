package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"llm-fusion-engine/internal/core"
	"net/http"
	"net/url"
	"strings"

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

	// 2. Get the provider from route result
	provider := routeResult.Provider
	if provider == nil {
		return nil, errors.New("no provider found in route result")
	}

	// 3. Parse provider config to get baseUrl and apiKey
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(provider.Config), &config); err != nil {
		return nil, errors.New("failed to parse provider config")
	}

	baseUrl, ok := config["baseUrl"].(string)
	if !ok || baseUrl == "" {
		return nil, errors.New("baseUrl not found in provider config")
	}

	// Get API key from config or use routed key
	apiKey := routeResult.ApiKey
	if configKey, ok := config["apiKey"].(string); ok && configKey != "" {
		apiKey = configKey
	}

	// 4. Modify the request with the resolved model
	requestBody["model"] = routeResult.ResolvedModel

	// Clean up "[undefined]" values sent by some clients
	for key, value := range requestBody {
		if strVal, ok := value.(string); ok && strVal == "[undefined]" {
			delete(requestBody, key)
		}
	}

	// 5. Construct the full API endpoint URL
	// Construct the full API endpoint URL
	var apiEndpoint string
	if chatEndpoint, ok := config["chatEndpoint"].(string); ok && chatEndpoint != "" {
		// If a specific chatEndpoint is defined, use it.
		apiEndpoint = baseUrl
		if strings.HasSuffix(apiEndpoint, "/") {
			apiEndpoint = apiEndpoint[:len(apiEndpoint)-1]
		}
		if !strings.HasPrefix(chatEndpoint, "/") {
			chatEndpoint = "/" + chatEndpoint
		}
		apiEndpoint += chatEndpoint
	} else {
		// Otherwise, use the baseUrl, but only append the default path if the baseUrl doesn't already have a path.
		apiEndpoint = baseUrl
		parsedUrl, err := url.Parse(apiEndpoint)
		if err == nil && (parsedUrl.Path == "" || parsedUrl.Path == "/") {
			if !strings.HasSuffix(apiEndpoint, "/") {
				apiEndpoint += "/"
			}
			apiEndpoint += "v1/chat/completions"
		}
	}

	// 6. Create request body
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// 7. Create HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// 8. Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 9. Execute the request
	client := &http.Client{}
	return client.Do(req)
}