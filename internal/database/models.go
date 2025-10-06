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

// Provider holds the configuration for a specific LLM provider.
// Note: The original GroupID is removed to align with a more direct provider management model.
// Group-based routing can be implemented at a higher level if needed.
type Provider struct {
	BaseModel
	Name         string `gorm:"uniqueIndex;not null" json:"name"` // e.g., "MyOpenAIInstance"
	Type         string `gorm:"index;not null" json:"type"`     // e.g., openai, anthropic, gemini
	Config       string `gorm:"type:text" json:"config"`       // JSON string for provider-specific settings (e.g., API key, base URL)
	Console      string `gorm:"type:varchar(255)" json:"console"` // Optional console URL for the provider
	Enabled      bool   `gorm:"default:true" json:"enabled"`
	Weight       uint   `gorm:"default:1" json:"weight"` // For load balancing if used directly
	HealthStatus string `gorm:"default:'unknown'" json:"healthStatus"` // healthy/unhealthy/unknown
	LastChecked  *time.Time `json:"lastChecked"`
	Latency      uint   `json:"latency"` // in milliseconds
	ApiKeys      []ApiKey `json:"apiKeys"` // Has-many relationship
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

// Model represents a user-friendly definition of a model with common configurations.
type Model struct {
	BaseModel
	Name     string `gorm:"uniqueIndex;not null" json:"name"` // e.g., "GPT-4-Turbo", "Claude-3-Sonnet"
	Remark   string `gorm:"type:text" json:"remark"`           // Description or notes for the model
	MaxRetry int    `gorm:"default:3" json:"maxRetry"`         // Global retry limit for this model
	Timeout  int    `gorm:"default:30" json:"timeout"`         // Global timeout in seconds for this model
	Enabled  bool   `gorm:"default:true" json:"enabled"`       // Whether this model definition is active
}

// ModelProviderMapping links a Model definition to a specific Provider instance,
// defining how the model is served by that provider and its specific capabilities.
type ModelProviderMapping struct {
	BaseModel
	ModelID          uint   `gorm:"index:idx_model_provider;not null" json:"modelId"`
	ProviderID       uint   `gorm:"index:idx_model_provider;not null" json:"providerId"`
	ProviderModel    string `gorm:"not null" json:"providerModel"` // The actual model ID on the provider's platform (e.g., "gpt-4-0125-preview")
	ToolCall         *bool  `json:"toolCall"`         // Can this model instance accept tool calls?
	StructuredOutput *bool  `json:"structuredOutput"` // Can this model instance accept structured output requests?
	Image            *bool  `json:"image"`            // Can this model instance accept image inputs (vision)?
	Weight           int    `gorm:"default:1" json:"weight"` // Weight for load balancing among multiple provider instances for the same model
	Enabled          bool   `gorm:"default:true" json:"enabled"` // Is this specific mapping enabled?
	Model            Model   `gorm:"foreignKey:ModelID" json:"model"`
	Provider         Provider `gorm:"foreignKey:ProviderID" json:"provider"`
}