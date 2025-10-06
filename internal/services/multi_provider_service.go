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
	c *gin.Context,
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
			// Create a unique request ID for logging
			requestID := uuid.New().String()
			c.Set("requestID", requestID) // Store it in context for later use
			s.LogRequest(requestID, requestBody, proxyKey, provider.Name, apiEndpoint, nil, false, latency, 0, 0, 0)
			continue // Retry with the next provider
		}

		// For logging, we need to read the body and then replace it.
		// This logic is now centralized within the LogRequest function.
		// We pass a placeholder for token usage for now, which will be updated.
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// We will parse the body for tokens later
			requestID := uuid.New().String()
			c.Set("requestID", requestID)
			s.LogRequest(requestID, requestBody, proxyKey, provider.Name, apiEndpoint, resp, true, latency, 0, 0, 0)
			return resp, nil // Success
		}

		// Handle non-2xx responses
		requestID := uuid.New().String()
		c.Set("requestID", requestID)
		s.LogRequest(requestID, requestBody, proxyKey, provider.Name, apiEndpoint, resp, false, latency, 0, 0, 0)
		// The original response body is closed within LogRequest, so we don't do it here.

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

	// Normalize base URL
	baseUrl = strings.TrimSuffix(baseUrl, "/")

	// Determine the standard suffix for the given provider type
	var fullPath string
	var pathAfterV1 string

	switch strings.ToLower(providerType) {
	case "anthropic":
		fullPath = "/v1/messages"
		pathAfterV1 = "/messages"
	case "gemini":
		fullPath = "/v1beta/models/gemini-pro:generateContent"
		// Gemini is special, doesn't follow /v1 pattern, so just append
		return baseUrl + fullPath, nil
	default:
		fullPath = "/v1/chat/completions"
		pathAfterV1 = "/chat/completions"
	}

	// Case 1: URL is already complete
	if strings.HasSuffix(baseUrl, fullPath) {
		return baseUrl, nil
	}

	// Case 2: URL ends with /v1
	if strings.HasSuffix(baseUrl, "/v1") {
		return baseUrl + pathAfterV1, nil
	}

	// Case 3: Bare domain or other path
	return baseUrl + fullPath, nil
}

// LogRequest logs the details of an API request and its response.
func (s *MultiProviderService) LogRequest(
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
		ID:               requestID,
		ProxyKey:         proxyKey,
		Model:            requestBody["model"].(string),
		Provider:         providerName,
		RequestURL:       requestUrl,
		RequestBody:      string(reqBodyBytes),
		ResponseBody:     string(respBodyBytes),
		ResponseStatus:   status,
		IsSuccess:        isSuccess,
		Latency:          latency.Milliseconds(),
		Timestamp:        time.Now(),
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokens,
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