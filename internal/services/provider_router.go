package services

import (
	"errors"
	"llm-fusion-engine/internal/core"
	"llm-fusion-engine/internal/database"
	"sort"

	"gorm.io/gorm"
)

// ProviderRouter implements the IProviderRouter interface.
type ProviderRouter struct {
	db         *gorm.DB
	keyManager core.IKeyManager
}

// NewProviderRouter creates a new ProviderRouter.
func NewProviderRouter(db *gorm.DB, keyManager core.IKeyManager) *ProviderRouter {
	return &ProviderRouter{
		db:         db,
		keyManager: keyManager,
	}
}

// RouteRequestAsync selects a provider based on model mappings and performs failover.
func (r *ProviderRouter) RouteRequestAsync(model, proxyKey string, excludedProviders []uint) (*core.ProviderRouteResult, error) {
	// 1. Validate proxy key
	_, err := r.keyManager.ValidateProxyKeyAsync(proxyKey)
	if err != nil {
		return nil, errors.New("invalid proxy key")
	}

	// 2. Find all candidate providers via ModelProviderMapping
	var mappings []database.ModelProviderMapping
	query := r.db.Joins("JOIN models ON models.id = model_provider_mappings.model_id").
		Where("models.name = ?", model).
		Preload("Provider")

	if len(excludedProviders) > 0 {
		query = query.Where("provider_id NOT IN ?", excludedProviders)
	}

	if err := query.Find(&mappings).Error; err != nil || len(mappings) == 0 {
		return nil, errors.New("no provider mapping found for the given model")
	}

	// 3. Sort providers by priority (for failover)
	sort.Slice(mappings, func(i, j int) bool {
		// Higher priority value means it comes first
		return mappings[i].Provider.Priority > mappings[j].Provider.Priority
	})

	// 4. Iterate through sorted providers and try to get a key (this is the failover mechanism)
	for _, mapping := range mappings {
		provider := &mapping.Provider

		// Try to get a key for this provider
		var apiKey string
		var keyErr error
		if km, ok := r.keyManager.(*KeyManager); ok {
			apiKey, keyErr = km.GetNextKeyForProviderAsync(provider.ID)
		} else {
			keyErr = errors.New("key manager does not support direct provider key retrieval")
		}

		if keyErr == nil {
			// Success! We found a working provider and key.
			return &core.ProviderRouteResult{
				Group:         nil, // No group context when using direct mapping
				Provider:      provider,
				ApiKey:        apiKey,
				ResolvedModel: mapping.ProviderModel,
			}, nil
		}
		// If getting a key fails, the loop will continue to the next provider with lower priority.
	}

	// 5. If the loop completes, it means no provider in the mapping had a working key
	return nil, errors.New("no available API key for any of the mapped providers")
}

func init() {
	// rand.Seed is no longer needed as rand is not used
}