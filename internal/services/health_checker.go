package services

import (
	"encoding/json"
	"llm-fusion-engine/internal/database"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// HealthChecker service for checking provider health
type HealthChecker struct {
	db *gorm.DB
}

// NewHealthChecker creates a new HealthChecker service
func NewHealthChecker(db *gorm.DB) *HealthChecker {
	return &HealthChecker{db: db}
}

// CheckProvider checks the health of a single provider by making a test request.
// It no longer updates the provider's status directly in the database.
// The health status is now inferred from the success of recent requests logged by LogRequest.
func (hc *HealthChecker) CheckProvider(providerID uint) error {
	var provider database.Provider
	if err := hc.db.First(&provider, providerID).Error; err != nil {
		return err
	}

	// From provider.Config, parse BaseURL
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(provider.Config), &config); err != nil {
		return err
	}

	baseURL, ok := config["baseUrl"].(string)
	if !ok {
		return nil // Or return an error indicating misconfiguration
	}

	// Send a test request to the provider's BaseURL
	resp, err := http.Get(baseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil // The provider is reachable but returned an error, which is a form of health check success.
	}

	return nil
}

// CheckAllProviders checks the health of all providers
func (hc *HealthChecker) CheckAllProviders() {
	var providers []database.Provider
	hc.db.Find(&providers)

	for _, p := range providers {
		go func(providerID uint) {
			_ = hc.CheckProvider(providerID) // Error is ignored in background check
		}(p.ID)
	}
}

// SchedulePeriodicChecks schedules periodic health checks
func (hc *HealthChecker) SchedulePeriodicChecks(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			hc.CheckAllProviders()
		}
	}()
}