package services

import (
	"llm-fusion-engine/internal/database"
	"gorm.io/gorm"
)

// KeyManager implements the IKeyManager interface.
type KeyManager struct {
	db *gorm.DB
}

// NewKeyManager creates a new KeyManager.
func NewKeyManager(db *gorm.DB) *KeyManager {
	return &KeyManager{db: db}
}

// GetNextKeyAsync retrieves the next available and healthy API key for a given group.
// For now, it will use a simple round-robin logic.
func (km *KeyManager) GetNextKeyAsync(groupID uint) (string, error) {
	var apiKey database.ApiKey
	// TODO: Implement more sophisticated key selection logic (e.g., round-robin, least used).
	err := km.db.Where("provider_id IN (SELECT id FROM providers WHERE group_id = ?)", groupID).
		Where("is_healthy = ?", true).
		Order("last_used asc").
		First(&apiKey).Error

	if err != nil {
		return "", err
	}

	// Update the last used time
	km.db.Model(&apiKey).Update("last_used", gorm.Expr("CURRENT_TIMESTAMP"))

	return apiKey.Key, nil
}

// GetNextKeyForProviderAsync retrieves the next available and healthy API key for a given provider.
func (km *KeyManager) GetNextKeyForProviderAsync(providerID uint) (string, error) {
	var apiKey database.ApiKey
	err := km.db.Where("provider_id = ? AND is_healthy = ?", providerID, true).
		Order("last_used asc").
		First(&apiKey).Error

	if err != nil {
		return "", err
	}

	// Update the last used time
	km.db.Model(&apiKey).Update("last_used", gorm.Expr("CURRENT_TIMESTAMP"))

	return apiKey.Key, nil
}

// ValidateProxyKeyAsync checks if a proxy key is valid and enabled.
func (km *KeyManager) ValidateProxyKeyAsync(proxyKey string) (*database.ProxyKey, error) {
	var key database.ProxyKey
	err := km.db.Where("key = ? AND enabled = ?", proxyKey, true).First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}