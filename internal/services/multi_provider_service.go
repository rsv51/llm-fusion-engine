package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"llm-fusion-engine/internal/core"
	"llm-fusion-engine/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MultiProviderService coordinates routing and API requests.
type MultiProviderService struct {
	router          core.IProviderRouter
	providerFactory core.IProviderFactory
	db              *gorm.DB
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

	var excludedProviders []uint
	var lastErr error

	for i := 0; i < 5; i++ { // Allow up to 5 retries (initial + 4 retries)
		// 1. Route the request
		routeResult, err := s.router.RouteRequestAsync(model, proxyKey, excludedProviders)
		if err != nil {
			return nil, err
		}

		// 2. Get the provider and prepare the request
		provider := routeResult.Provider
		if provider == nil {
			return nil, errors.New("no provider found in route result")
		}

		// Exclude this provider from future retries in this request
		excludedProviders = append(excludedProviders, provider.ID)

		var config map[string]interface{}
		if err := json.Unmarshal([]byte(provider.Config), &config); err != nil {
			lastErr = fmt.Errorf("failed to parse config for provider %s: %w", provider.Name, err)
			continue
		}

		baseUrl, _ := config["baseUrl"].(string)
		apiKey := routeResult.ApiKey
		if configKey, ok := config["apiKey"].(string); ok && configKey != "" {
			apiKey = configKey
		}

		requestBody["model"] = routeResult.ResolvedModel
		requestBody = cleanupUndefined(requestBody).(map[string]interface{})

		apiEndpoint, err := getRequestURL(provider.Type, baseUrl)
		if err != nil {
			lastErr = err
			continue
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		// 3. Execute the request and handle retries
		startTime := time.Now()
		req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{Timeout: time.Duration(provider.Timeout) * time.Second}
		resp, err := client.Do(req)
		latency := time.Since(startTime)

		if err != nil {
			lastErr = err
			s.LogRequest(requestBody, proxyKey, provider.Name, apiEndpoint, nil, false, latency)
			continue // Retry with the next provider
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			s.LogRequest(requestBody, proxyKey, provider.Name, apiEndpoint, resp, true, latency)
			return resp, nil // Success
		}

		// Handle non-2xx responses
		s.LogRequest(requestBody, proxyKey, provider.Name, apiEndpoint, resp, false, latency)
		resp.Body.Close() // Close body to allow for reuse of connection

		// Decide if we should retry
		shouldRetry := false
		if policy, ok := config["retryPolicy"].(map[string]interface{}); ok {
			if codes, ok := policy["statusCodes"].([]interface{}); ok {
				for _, code := range codes {
					if int(code.(float64)) == resp.StatusCode {
						shouldRetry = true
						break
					}
				}
			}
		} else {
			// Default retry policy for common transient errors
			defaultRetryCodes := []int{429, 500, 502, 503, 504}
			for _, code := range defaultRetryCodes {
				if code == resp.StatusCode {
					shouldRetry = true
					break
				}
			}
		}

		if !shouldRetry {
			return nil, fmt.Errorf("provider %s returned non-retriable status code %d", provider.Name, resp.StatusCode)
		}

		lastErr = fmt.Errorf("provider %s failed with status %d", provider.Name, resp.StatusCode)
		// Loop will continue to the next provider
	}

	return nil, fmt.Errorf("all retries failed. last error: %w", lastErr)
}

// getRequestURL acts as an adapter to get the correct API endpoint for different provider types.
func getRequestURL(providerType, baseUrl string) (string, error) {
	if baseUrl == "" {
		return "", errors.New("baseUrl is not configured for the provider")
	}

	// Normalize base URL by removing any trailing slashes
	baseUrl = strings.TrimSuffix(baseUrl, "/")

	// Determine the API path based on provider type
	var apiPath string
	switch strings.ToLower(providerType) {
	case "anthropic":
		apiPath = "/v1/messages"
	case "gemini":
		// This is a simplification. A real Gemini adapter would be more complex.
		apiPath = "/v1beta/models/gemini-pro:generateContent"
	default:
		// Default to OpenAI-compatible path for "openai", "azure", "openrouter", etc., and unknown types.
		apiPath = "/v1/chat/completions"
	}

	// Smartly combine baseUrl and apiPath
	// If baseUrl already ends with the path, don't append it again.
	if strings.HasSuffix(baseUrl, apiPath) {
		return baseUrl, nil
	}

	// If baseUrl contains a different /v1/ path, it's likely a custom endpoint. Trust it.
	if strings.Contains(baseUrl, "/v1/") && !strings.HasSuffix(baseUrl, apiPath) {
		return baseUrl, nil
	}
	
	return baseUrl + apiPath, nil
}

// LogRequest logs the details of an API request and its response.
func (s *MultiProviderService) LogRequest(
	requestBody map[string]interface{},
	proxyKey string,
	providerName string,
	requestUrl string,
	response *http.Response,
	isSuccess bool,
	latency time.Duration,
) {
	reqBodyBytes, _ := json.Marshal(requestBody)
	var respBodyBytes []byte
	var status int

	if response != nil {
		status = response.StatusCode
		// Read and then replace the body to allow it to be read again
		respBodyBytes, _ = ioutil.ReadAll(response.Body)
		response.Body.Close() // Close the original body
		response.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))
	}

	logEntry := database.Log{
		ID:             uuid.New().String(),
		ProxyKey:       proxyKey,
		Model:          requestBody["model"].(string),
		Provider:       providerName,
		RequestURL:     requestUrl,
		RequestBody:    string(reqBodyBytes),
		ResponseBody:   string(respBodyBytes),
		ResponseStatus: status,
		IsSuccess:      isSuccess,
		Latency:        latency.Milliseconds(),
		Timestamp:      time.Now(),
	}

	s.db.Create(&logEntry)
}

// cleanupUndefined recursively removes keys with "[undefined]" string values from a map.
func cleanupUndefined(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if strVal, ok := value.(string); ok && strVal == "[undefined]" {
				delete(v, key)
			} else {
				v[key] = cleanupUndefined(value)
			}
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = cleanupUndefined(value)
		}
		return v
	default:
		return data
	}
}