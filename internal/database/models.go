package database

import (
	"time"
	"gorm.io/gorm"
)

// BaseModel defines the common fields for all models.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// User represents a user of the system.
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"` // Hashed password
	IsAdmin  bool   `gorm:"default:false"`
	ProxyKeys []ProxyKey
}

// ProxyKey is a key used by end-users to access the API.
type ProxyKey struct {
	BaseModel
	UserID             uint   `json:"userId"`
	Key                string `gorm:"uniqueIndex;not null" json:"key"`
	Enabled            bool   `gorm:"default:true" json:"enabled"`
	AllowedGroups      string `json:"allowedGroups"` // JSON array of group IDs
	GroupBalancePolicy string `gorm:"default:'failover'" json:"groupBalancePolicy"`
	GroupWeights       string `json:"groupWeights"` // JSON object for weighted balancing
	RpmLimit           int    `json:"rpmLimit"`
	TpmLimit           int    `json:"tpmLimit"`
}

// Group represents a collection of provider configurations for routing.
type Group struct {
	BaseModel
	Name              string     `gorm:"uniqueIndex;not null" json:"name"`
	Enabled           bool       `gorm:"default:true" json:"enabled"`
	Priority          int        `gorm:"default:0" json:"priority"`
	Models            string     `gorm:"type:text" json:"models"` // JSON array of supported models
	ModelAliases      string     `gorm:"type:text" json:"modelAliases"` // JSON object for model name mapping
	LoadBalancePolicy string     `gorm:"default:'failover'" json:"loadBalancePolicy"` // e.g., failover, round_robin, weighted
	Providers         []Provider `json:"providers"`                                   // Has-many relationship
}

// Provider holds the configuration for a specific LLM provider within a group.
type Provider struct {
	BaseModel
	GroupID      uint   `gorm:"index" json:"groupId"`
	ProviderType string `gorm:"not null" json:"providerType"` // e.g., openai, anthropic, gemini
	ApiKeys      []ApiKey `json:"apiKeys"` // Has-many relationship
	Weight       uint   `gorm:"default:1" json:"weight"`
	Enabled      bool   `gorm:"default:true" json:"enabled"`
	BaseURL      string `gorm:"type:varchar(255)" json:"baseUrl"`
	Timeout      int    `gorm:"default:30" json:"timeout"`    // 秒
	MaxRetries   int    `gorm:"default:3" json:"maxRetries"`
	HealthStatus string `gorm:"default:'unknown'" json:"healthStatus"` // healthy/unhealthy/unknown
	LastChecked  *time.Time `json:"lastChecked"`
	Latency      uint   `json:"latency"` // in milliseconds
}

// ApiKey stores an individual API key for a provider.
type ApiKey struct {
	BaseModel
	ProviderID uint      `gorm:"index" json:"providerId"`
	Key        string    `gorm:"uniqueIndex;not null" json:"key"`
	LastUsed   time.Time `json:"lastUsed"`
	IsHealthy  bool      `gorm:"default:true" json:"isHealthy"`
	RpmLimit   int       `json:"rpmLimit"`
	TpmLimit   int       `json:"tpmLimit"`
}

// RequestLog records API request details for monitoring and analytics.
type RequestLog struct {
	BaseModel
	GroupID          uint   `gorm:"index" json:"groupId"`
	ProviderID       uint   `gorm:"index" json:"providerId"`
	Model            string `gorm:"index" json:"model"`
	StatusCode       int    `gorm:"index" json:"statusCode"`
	LatencyMs        int64  `json:"latencyMs"`
	PromptTokens     int    `json:"promptTokens"`
	CompletionTokens int    `json:"completionTokens"`
	TotalTokens      int    `json:"totalTokens"`
	ErrorMessage     string `json:"errorMessage"`
	RequestIP        string `json:"requestIp"`
	UserAgent        string `json:"userAgent"`
}

// Model represents a specific LLM model available in the system.
type Model struct {
	BaseModel
	Name        string  `gorm:"uniqueIndex;not null" json:"name"`
	Provider    string  `gorm:"index" json:"provider"`
	Category    string  `json:"category"` // text/image/audio/video
	MaxTokens   int     `json:"maxTokens"`
	InputPrice  float64 `json:"inputPrice"`
	OutputPrice float64 `json:"outputPrice"`
	Description string  `json:"description"`
	Enabled     bool    `gorm:"default:true" json:"enabled"`
}

// ModelMapping allows aliasing model names to specific provider models.
type ModelMapping struct {
	BaseModel
	UserFriendlyName  string `gorm:"uniqueIndex;not null" json:"userFriendlyName"` // e.g., "fast-model"
	ProviderModelName string `gorm:"not null" json:"providerModelName"`           // e.g., "gpt-3.5-turbo"
	ProviderID        uint   `gorm:"not null" json:"providerId"`           // Foreign key to Provider
	Provider          Provider `gorm:"foreignKey:ProviderID" json:"provider"`
}