package providers

import (
	"encoding/json"
	"fmt"
	"time"
)

// CreateClient creates a provider client based on the provider type
func CreateClient(providerType string, configJSON string) (ProviderClient, error) {
	// Parse the config JSON
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(configJSON), &configMap); err != nil {
		return nil, fmt.Errorf("failed to parse provider config: %w", err)
	}
	
	// Extract common config fields
	config := ProviderConfig{
		APIKey:     getStringFromConfig(configMap, "apiKey"),
		BaseURL:    getStringFromConfig(configMap, "baseUrl"),
		MaxRetries: getIntFromConfig(configMap, "maxRetries", 3),
	}
	
	// Parse timeout
	if timeoutSec := getIntFromConfig(configMap, "timeout", 30); timeoutSec > 0 {
		config.Timeout = time.Duration(timeoutSec) * time.Second
	} else {
		config.Timeout = 30 * time.Second
	}
	
	// Create client based on provider type
	switch providerType {
	case "openai":
		return NewOpenAIClient(config), nil
	case "anthropic":
		// Anthropic doesn't have a public models API endpoint
		// We'll need to return a static list or handle differently
		return nil, fmt.Errorf("anthropic provider does not support dynamic model listing")
	case "gemini":
		// Google Gemini also doesn't have a simple models endpoint
		return nil, fmt.Errorf("gemini provider does not support dynamic model listing")
	default:
		// For unknown types, try OpenAI-compatible API
		return NewOpenAIClient(config), nil
	}
}

// Helper functions to extract values from config map
func getStringFromConfig(config map[string]interface{}, key string) string {
	if val, ok := config[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getIntFromConfig(config map[string]interface{}, key string, defaultVal int) int {
	if val, ok := config[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return defaultVal
}