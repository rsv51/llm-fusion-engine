package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIClient implements ProviderClient for OpenAI-compatible APIs
type OpenAIClient struct {
	config ProviderConfig
	client *http.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config ProviderConfig) *OpenAIClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}
	
	return &OpenAIClient{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// OpenAIModelsResponse represents the response from OpenAI /v1/models endpoint
type OpenAIModelsResponse struct {
	Data []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
	Object string `json:"object"`
}

// GetModels retrieves the list of available models from OpenAI-compatible API
func (c *OpenAIClient) GetModels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/models", c.config.BaseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var modelsResp OpenAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	models := make([]string, 0, len(modelsResp.Data))
	for _, model := range modelsResp.Data {
		models = append(models, model.ID)
	}
	
	return models, nil
}

// ValidateConfig validates the provider configuration
func (c *OpenAIClient) ValidateConfig() error {
	if c.config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	return nil
}