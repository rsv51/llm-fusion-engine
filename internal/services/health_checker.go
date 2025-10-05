package services

import (
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

// CheckProvider checks the health of a single provider
func (hc *HealthChecker) CheckProvider(providerID uint) (*database.Provider, error) {
	var provider database.Provider
	if err := hc.db.First(&provider, providerID).Error; err != nil {
		return nil, err
	}

	// 模拟健康检查逻辑
	// 在实际应用中,这里应该向 provider.BaseURL 发送一个测试请求
	startTime := time.Now()
	resp, err := http.Get(provider.BaseURL)
	latency := time.Since(startTime)

	now := time.Now()
	provider.LastChecked = &now

	if err != nil || resp.StatusCode >= 400 {
		provider.HealthStatus = "unhealthy"
	} else {
		provider.HealthStatus = "healthy"
	}

	// 更新数据库
	hc.db.Save(&provider)

	return &provider, nil
}

// CheckAllProviders checks the health of all providers
func (hc *HealthChecker) CheckAllProviders() {
	var providers []database.Provider
	hc.db.Find(&providers)

	for _, p := range providers {
		go hc.CheckProvider(p.ID)
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