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

// CheckProvider checks the health of a single provider by making a test request
// and updates the provider's health status, latency, and last checked time in the database.
func (hc *HealthChecker) CheckProvider(providerID uint) error {
	var provider database.Provider
	if err := hc.db.First(&provider, providerID).Error; err != nil {
		return err
	}

	// From provider.Config, parse BaseURL
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(provider.Config), &config); err != nil {
		// Update as unhealthy due to config error
		now := time.Now()
		provider.HealthStatus = "unhealthy"
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}

	baseURL, ok := config["baseUrl"].(string)
	if !ok || baseURL == "" {
		// Update as unknown due to missing baseUrl
		now := time.Now()
		provider.HealthStatus = "unknown"
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return nil
	}

	// Send a test request to the provider's BaseURL and measure latency
	startTime := time.Now()
	resp, err := http.Get(baseURL)
	latency := time.Since(startTime).Milliseconds()
	now := time.Now()

	if err != nil {
		// Update as unhealthy due to request error
		provider.HealthStatus = "unhealthy"
		provider.Latency = nil
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}
	defer resp.Body.Close()

	// Update health status based on response
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		provider.HealthStatus = "healthy"
		provider.Latency = &latency
		provider.LastChecked = &now
		hc.db.Save(&provider)
	} else {
		provider.HealthStatus = "unhealthy"
		provider.Latency = &latency
		provider.LastChecked = &now
		hc.db.Save(&provider)
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