package services

import (
	"encoding/json"
	"fmt"
	"io"
	"llm-fusion-engine/internal/constants"
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
		provider.HealthStatus = string(constants.HealthStatusUnhealthy)
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}

	baseURL, ok := config["baseUrl"].(string)
	if !ok || baseURL == "" {
		// Update as unknown due to missing baseUrl
		now := time.Now()
		provider.HealthStatus = string(constants.HealthStatusUnknown)
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return nil
	}

	// Get API key from config
	apiKey, _ := config["apiKey"].(string)
	
	// Use a default model for testing based on provider type
	// Many proxy services don't support /v1/models endpoint reliably
	var testModel string
	switch strings.ToLower(provider.Type) {
	case "anthropic":
		testModel = "claude-3-haiku-20240307"
	case "gemini":
		testModel = "gemini-1.5-flash"
	default:
		testModel = "gpt-3.5-turbo" // Default for OpenAI-compatible APIs
	}
	
	log.Printf("[HealthCheck] Provider ID=%d, Name=%s, Type=%s, Using default test model: %s",
		providerID, provider.Name, provider.Type, testModel)
	
	// Send chat completion request with "hi" directly
	chatURL := strings.TrimSuffix(baseURL, "/") + "/chat/completions"
	chatPayload := map[string]interface{}{
		"model": testModel,
		"messages": []map[string]string{
			{"role": "user", "content": "hi"},
		},
		"max_tokens": 10,
	}
	
	chatBody, err := json.Marshal(chatPayload)
	if err != nil {
		log.Printf("[HealthCheck] Provider ID=%d: Failed to marshal chat request: %v", providerID, err)
		now := time.Now()
		provider.HealthStatus = string(constants.HealthStatusUnhealthy)
		provider.Latency = nil
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}
	
	log.Printf("[HealthCheck] Provider ID=%d: Sending chat request to %s with model %s", providerID, chatURL, testModel)
	
	chatReq, err := http.NewRequest("POST", chatURL, strings.NewReader(string(chatBody)))
	if err != nil {
		log.Printf("[HealthCheck] Provider ID=%d: Failed to create chat request: %v", providerID, err)
		now := time.Now()
		provider.HealthStatus = string(constants.HealthStatusUnhealthy)
		provider.Latency = nil
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}
	
	chatReq.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		chatReq.Header.Set("Authorization", "Bearer "+apiKey)
	}
	
	// Measure latency for chat request
	client := &http.Client{Timeout: 10 * time.Second}
	startTime := time.Now()
	resp, err := client.Do(chatReq)
	latency := time.Since(startTime).Milliseconds()
	now := time.Now()

	if err != nil {
		log.Printf("[HealthCheck] Provider ID=%d: Chat request error: %v", providerID, err)
		provider.HealthStatus = string(constants.HealthStatusUnhealthy)
		provider.Latency = nil
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return err
	}
	defer resp.Body.Close()

	// Read chat response for debugging
	body, _ := io.ReadAll(resp.Body)
	log.Printf("[HealthCheck] Provider ID=%d: Chat response status=%d, latency=%dms, body=%s",
		providerID, resp.StatusCode, latency, string(body))

	// Validate chat response
	statusCode := resp.StatusCode
	
	if statusCode >= 200 && statusCode < 300 {
		// Parse response to verify we got actual content
		var chatResponse map[string]interface{}
		if err := json.Unmarshal(body, &chatResponse); err != nil {
			log.Printf("[HealthCheck] Provider ID=%d: Failed to parse chat response: %v", providerID, err)
			provider.HealthStatus = string(constants.HealthStatusDegraded)
			provider.Latency = &latency
			provider.LastStatusCode = &statusCode
			provider.LastChecked = &now
			hc.db.Save(&provider)
			return fmt.Errorf("invalid chat response format")
		}
		
		// Check if response has content (indicates successful generation)
		choices, ok := chatResponse["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			log.Printf("[HealthCheck] Provider ID=%d: No choices in chat response", providerID)
			provider.HealthStatus = string(constants.HealthStatusDegraded)
			provider.Latency = &latency
			provider.LastStatusCode = &statusCode
			provider.LastChecked = &now
			hc.db.Save(&provider)
			return fmt.Errorf("no response content")
		}
		
		firstChoice, ok := choices[0].(map[string]interface{})
		if !ok {
			log.Printf("[HealthCheck] Provider ID=%d: Invalid choice format", providerID)
			provider.HealthStatus = string(constants.HealthStatusDegraded)
			provider.Latency = &latency
			provider.LastStatusCode = &statusCode
			provider.LastChecked = &now
			hc.db.Save(&provider)
			return fmt.Errorf("invalid choice format")
		}
		
		message, ok := firstChoice["message"].(map[string]interface{})
		if !ok {
			log.Printf("[HealthCheck] Provider ID=%d: No message in choice", providerID)
			provider.HealthStatus = string(constants.HealthStatusDegraded)
			provider.Latency = &latency
			provider.LastStatusCode = &statusCode
			provider.LastChecked = &now
			hc.db.Save(&provider)
			return fmt.Errorf("no message content")
		}
		
		content, ok := message["content"].(string)
		if !ok || content == "" {
			log.Printf("[HealthCheck] Provider ID=%d: Empty content in response", providerID)
			provider.HealthStatus = string(constants.HealthStatusDegraded)
			provider.Latency = &latency
			provider.LastStatusCode = &statusCode
			provider.LastChecked = &now
			hc.db.Save(&provider)
			return fmt.Errorf("empty response content")
		}
		
		// Success - provider responded with actual content
		log.Printf("[HealthCheck] Provider ID=%d: HEALTHY - Got valid response: '%s' (latency: %dms)",
			providerID, content, latency)
		provider.HealthStatus = string(constants.HealthStatusHealthy)
		provider.Latency = &latency
		provider.LastStatusCode = &statusCode
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return nil
		
	} else if statusCode == 401 || statusCode == 403 {
		// Authentication/authorization errors
		log.Printf("[HealthCheck] Provider ID=%d: DEGRADED (status %d - auth issue)", providerID, statusCode)
		provider.HealthStatus = string(constants.HealthStatusDegraded)
		provider.Latency = &latency
		provider.LastStatusCode = &statusCode
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return fmt.Errorf("authentication/authorization failed with status code: %d", statusCode)
	} else {
		log.Printf("[HealthCheck] Provider ID=%d: UNHEALTHY (status %d)", providerID, statusCode)
		provider.HealthStatus = string(constants.HealthStatusUnhealthy)
		provider.Latency = &latency
		provider.LastStatusCode = &statusCode
		provider.LastChecked = &now
		hc.db.Save(&provider)
		return fmt.Errorf("health check failed with status code: %d", statusCode)
	}
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