package services

import (
	"encoding/json"
	"errors"
	"llm-fusion-engine/internal/core"
	"llm-fusion-engine/internal/database"
	"math/rand"
	"sort"
	"strings"
	"time"

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

// RouteRequestAsync selects a provider group based on the model, proxy key, and other strategies.
func (r *ProviderRouter) RouteRequestAsync(model, proxyKey string, excludedGroups []string) (*core.ProviderRouteResult, error) {
	// 1. Validate proxy key
	validatedProxyKey, err := r.keyManager.ValidateProxyKeyAsync(proxyKey)
	if err != nil {
		return nil, errors.New("invalid proxy key")
	}

	// 2. Find candidate groups that support the model
	candidateGroups, err := r.findCandidateGroups(model, validatedProxyKey, excludedGroups)
	if err != nil || len(candidateGroups) == 0 {
		return nil, errors.New("no available provider group found for the given model")
	}

	// 3. Select a group based on the load balancing policy
	selectedGroup, err := r.selectGroup(candidateGroups, validatedProxyKey)
	if err != nil {
		return nil, err
	}

	// 4. Get the next available API key for the selected group
	apiKey, err := r.keyManager.GetNextKeyAsync(selectedGroup.ID)
	if err != nil {
		// This could be a point for failover logic in the future
		return nil, errors.New("no available API key in the selected group")
	}

	// 5. Resolve model alias
	resolvedModel := r.resolveModelAlias(model, selectedGroup.ModelAliases)

	return &core.ProviderRouteResult{
		Group:         selectedGroup,
		ApiKey:        apiKey,
		ResolvedModel: resolvedModel,
	}, nil
}

// findCandidateGroups finds all groups that can handle the request.
func (r *ProviderRouter) findCandidateGroups(model string, proxyKey *database.ProxyKey, excludedGroups []string) ([]*database.Group, error) {
	var groups []*database.Group
	query := r.db.Where("enabled = ?", true)

	// Exclude already attempted groups
	if len(excludedGroups) > 0 {
		query = query.Where("name NOT IN ?", excludedGroups)
	}

	// Filter by groups allowed by the proxy key
	if proxyKey != nil && proxyKey.AllowedGroups != "" {
		var allowedGroupIDs []uint
		if err := json.Unmarshal([]byte(proxyKey.AllowedGroups), &allowedGroupIDs); err == nil {
			query = query.Where("id IN ?", allowedGroupIDs)
		}
	}

	if err := query.Find(&groups).Error; err != nil {
		return nil, err
	}

	// Filter groups that support the model
	var candidateGroups []*database.Group
	for _, group := range groups {
		var supportedModels []string
		if err := json.Unmarshal([]byte(group.Models), &supportedModels); err != nil {
			continue // Skip group if models JSON is invalid
		}

		isSupported := false
		for _, m := range supportedModels {
			if m == model {
				isSupported = true
				break
			}
		}

		if isSupported || r.resolveModelAlias(model, group.ModelAliases) != model {
			candidateGroups = append(candidateGroups, group)
		}
	}

	return candidateGroups, nil
}

// selectGroup applies the load balancing policy to choose a group.
func (r *ProviderRouter) selectGroup(groups []*database.Group, proxyKey *database.ProxyKey) (*database.Group, error) {
	if len(groups) == 0 {
		return nil, errors.New("no candidate groups to select from")
	}

	policy := "failover" // default policy
	if proxyKey != nil && proxyKey.GroupBalancePolicy != "" {
		policy = proxyKey.GroupBalancePolicy
	}

	switch strings.ToLower(policy) {
	case "round_robin":
		// TODO: Implement round-robin logic (requires state/cache)
		return groups[0], nil // Placeholder
	case "weighted":
		return r.selectByWeightedRandom(groups, proxyKey)
	case "random":
		return groups[rand.Intn(len(groups))], nil
	case "failover":
		fallthrough
	default:
		return r.selectByFailover(groups)
	}
}

// selectByFailover selects the group with the highest priority.
func (r *ProviderRouter) selectByFailover(groups []*database.Group) (*database.Group, error) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Priority > groups[j].Priority
	})
	return groups[0], nil
}

// selectByWeightedRandom selects a group based on weights.
func (r *ProviderRouter) selectByWeightedRandom(groups []*database.Group, proxyKey *database.ProxyKey) (*database.Group, error) {
	// TODO: Implement weighted random selection based on provider weights within the group
	// and group weights in the proxy key.
	return r.selectByFailover(groups) // Fallback to failover for now
}

// resolveModelAlias checks if the model has an alias in the group's config.
func (r *ProviderRouter) resolveModelAlias(model, modelAliasesJSON string) string {
	if modelAliasesJSON == "" {
		return model
	}

	var aliases map[string]string
	if err := json.Unmarshal([]byte(modelAliasesJSON), &aliases); err != nil {
		return model // Return original model if JSON is invalid
	}

	if alias, ok := aliases[model]; ok {
		return alias
	}

	return model
}

func init() {
	rand.Seed(time.Now().UnixNano())
}