package database

import (
	"time"
	"gorm.io/gorm"
)

// BaseModel defines the common fields for all models.
type BaseModel struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
	UserID            uint
	Key               string `gorm:"uniqueIndex;not null"`
	Enabled           bool   `gorm:"default:true"`
	AllowedGroups     string // JSON array of group IDs
	GroupBalancePolicy string `gorm:"default:'failover'"`
	GroupWeights      string // JSON object for weighted balancing
	RpmLimit          int
	TpmLimit          int
}

// Group represents a collection of provider configurations for routing.
type Group struct {
	BaseModel
	Name              string `gorm:"uniqueIndex;not null"`
	Enabled           bool   `gorm:"default:true"`
	Priority          int    `gorm:"default:0"`
	Models            string `gorm:"type:text"` // JSON array of supported models
	ModelAliases      string `gorm:"type:text"` // JSON object for model name mapping
	LoadBalancePolicy string `gorm:"default:'failover'"` // e.g., failover, round_robin, weighted
	Providers         []Provider // Has-many relationship
}

// Provider holds the configuration for a specific LLM provider within a group.
type Provider struct {
	BaseModel
	GroupID      uint   `gorm:"index"`
	ProviderType string `gorm:"not null"` // e.g., openai, anthropic, gemini
	ApiKeys      []ApiKey // Has-many relationship
	Weight       uint   `gorm:"default:1"`
	Enabled      bool   `gorm:"default:true"`
	BaseURL      string `gorm:"type:varchar(255)"`
	Timeout      int    `gorm:"default:30"`    // 秒
	MaxRetries   int    `gorm:"default:3"`
	HealthStatus string `gorm:"default:'unknown'"` // healthy/unhealthy/unknown
	LastChecked  *time.Time
	Latency      uint   // in milliseconds
}

// ApiKey stores an individual API key for a provider.
type ApiKey struct {
	BaseModel
	ProviderID  uint   `gorm:"index"`
	Key         string `gorm:"uniqueIndex;not null"`
	LastUsed    time.Time
	IsHealthy   bool   `gorm:"default:true"`
	RpmLimit    int    // Requests per minute limit
	TpmLimit    int    // Tokens per minute limit
}

// RequestLog records API request details for monitoring and analytics.
type RequestLog struct {
	BaseModel
	GroupID          uint   `gorm:"index"`
	ProviderID       uint   `gorm:"index"`
	Model            string `gorm:"index"`
	StatusCode       int    `gorm:"index"`
	LatencyMs        int64
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	ErrorMessage     string
	RequestIP        string
	UserAgent        string
}

// Model represents a specific LLM model available in the system.
type Model struct {
	BaseModel
	Name         string `gorm:"uniqueIndex;not null"`
	Provider     string `gorm:"index"`
	Category     string // text/image/audio/video
	MaxTokens    int
	InputPrice   float64
	OutputPrice  float64
	Description  string
	Enabled      bool `gorm:"default:true"`
}

// ModelMapping allows aliasing model names to specific provider models.
type ModelMapping struct {
	BaseModel
	UserFriendlyName string `gorm:"uniqueIndex;not null"` // e.g., "fast-model"
	ProviderModelName string `gorm:"not null"`           // e.g., "gpt-3.5-turbo"
	ProviderID        uint   `gorm:"not null"`           // Foreign key to Provider
	Provider          Provider `gorm:"foreignKey:ProviderID"`
}