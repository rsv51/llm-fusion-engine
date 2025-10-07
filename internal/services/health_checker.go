package services

import (
	"encoding/json"
	"fmt"
	"io"
	"llm-fusion-engine/internal/database"
	"log"
	"net/http"
	"strings"
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

	// Get API key from config
	apiKey, _ := config["apiKey"].(string)
	
	// For OpenAI-compatible APIs, test the /v1/models endpoint
	modelsURL := strings.TrimSuffix(baseURL, "/") + "/v1/models"
	
	log.Printf("[HealthCheck] Provider ID=%d, Name=%s, Testing URL: %s", providerID, provider.Name, modelsURL)
	
	// Create request
	req, err := http.NewRequest("GET", modelsURL, nil)
	if err != nil {
		log.Printf("[HealthCheck] Provider ID=%d: Failed to create request: %v", providerID, err)
		now := time.Now()
		provider.HealthStatus = "unhealthy"
		provider.Latency = nil
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}
	
	// Add authorization header if API key exists
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
		log.Printf("[HealthCheck] Provider ID=%d: Using API key authentication", providerID)
	} else {
		log.Printf("[HealthCheck] Provider ID=%d: No API key found in config", providerID)
	}
	
	// Send request and measure latency
	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	latency := time.Since(startTime).Milliseconds()
	now := time.Now()

	if err != nil {
		log.Printf("[HealthCheck] Provider ID=%d: Request error: %v", providerID, err)
		// Update as unhealthy due to request error
		provider.HealthStatus = "unhealthy"
		provider.Latency = nil
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}
	defer resp.Body.Close()

	// Read response body for debugging
	body, _ := io.ReadAll(resp.Body)
	log.Printf("[HealthCheck] Provider ID=%d: Status=%d, Latency=%dms, Body=%s",
		providerID, resp.StatusCode, latency, string(body))

	// Update health status based on response
	statusCode := resp.StatusCode
	
	// Determine health status with improved logic
	if statusCode >= 200 && statusCode < 300 {
		log.Printf("[HealthCheck] Provider ID=%d: Marked as HEALTHY (status %d)", providerID, statusCode)
		provider.HealthStatus = "healthy"
		provider.Latency = &latency
		provider.LastStatusCode = &statusCode
		provider.LastChecked = &now
		hc.db.Save(&provider)
	} else if statusCode == 401 || statusCode == 403 {
		// Authentication/authorization errors - provider is reachable but credentials may be invalid
		log.Printf("[HealthCheck] Provider ID=%d: Marked as DEGRADED (status %d - auth issue)", providerID, statusCode)
		provider.HealthStatus = "degraded"
		provider.Latency = &latency
		provider.LastStatusCode = &statusCode
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return fmt.Errorf("authentication/authorization failed with status code: %d", statusCode)
	} else {
		log.Printf("[HealthCheck] Provider ID=%d: Marked as UNHEALTHY (status %d)", providerID, statusCode)
		provider.HealthStatus = "unhealthy"
		provider.Latency = &latency
		provider.LastStatusCode = &statusCode
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return fmt.Errorf("health check failed with status code: %d", statusCode)
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