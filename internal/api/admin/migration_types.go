package admin

import "llm-fusion-engine/internal/database"

// MigrationData is a unified structure for importing and exporting all settings.
type MigrationData struct {
	Groups        []database.Group        `json:"groups" yaml:"groups"`
	Providers     []database.Provider     `json:"providers" yaml:"providers"`
	ApiKeys       []database.ApiKey       `json:"apiKeys" yaml:"apiKeys"`
	ModelMappings []database.ModelMapping `json:"modelMappings" yaml:"modelMappings"`
}