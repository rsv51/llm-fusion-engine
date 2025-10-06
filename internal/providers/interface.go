package providers

import (
	"context"
	"time"
)

// ProviderClient defines the interface for interacting with LLM providers
type ProviderClient interface {
	// GetModels retrieves the list of available models from the provider
	GetModels(ctx context.Context) ([]string, error)
	
	// ValidateConfig validates the provider configuration
	ValidateConfig() error
}

// ProviderConfig represents common configuration for all providers
type ProviderConfig struct {
	APIKey    string
	BaseURL   string
	Timeout   time.Duration
	MaxRetries int
}