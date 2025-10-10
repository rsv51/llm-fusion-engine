package constants

// HealthStatus 定义健康状态类型
type HealthStatus string

const (
	// HealthStatusHealthy 健康状态
	HealthStatusHealthy HealthStatus = "healthy"
	
	// HealthStatusDegraded 降级状态
	HealthStatusDegraded HealthStatus = "degraded"
	
	// HealthStatusUnhealthy 不健康状态
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	
	// HealthStatusUnknown 未知状态
	HealthStatusUnknown HealthStatus = "unknown"
)

// IsValid 检查健康状态是否有效
func (s HealthStatus) IsValid() bool {
	switch s {
	case HealthStatusHealthy, HealthStatusDegraded, HealthStatusUnhealthy, HealthStatusUnknown:
		return true
	default:
		return false
	}
}

// String 返回健康状态的字符串表示
func (s HealthStatus) String() string {
	return string(s)
}