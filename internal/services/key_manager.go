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

// ValidateProxyKeyAsync checks if a proxy key is valid and enabled.
func (km *KeyManager) ValidateProxyKeyAsync(proxyKey string) (*database.ProxyKey, error) {
	var key database.ProxyKey
	err := km.db.Where("key = ? AND enabled = ?", proxyKey, true).First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}