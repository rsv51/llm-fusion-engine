// 系统健康和统计相关类型定义

export interface SystemHealth {
  status: 'healthy' | 'unhealthy'
  databaseStatus: 'connected' | 'disconnected'
  cacheStatus: 'connected' | 'disconnected' | 'not_configured'
  uptime: number
  version: string
  providers: ProviderHealth[]
}

export interface ProviderHealth {
  providerId: number
  providerName: string
  providerType: string
  isHealthy: boolean
  responseTimeMs?: number
  errorMessage?: string
  lastCheckAt: string
}

export interface SystemStats {
  totalRequests: number
  successRate: number
  avgResponseTimeMs: number
  totalTokens: number
  activeProviders: number
  activeKeys: number
  requestsChange?: number
  providers?: ProviderStats[]
  startTime: string
}

export interface ProviderStats {
  providerId: number
  providerName: string
  requestCount: number
  successCount: number
  errorCount: number
  totalTokens: number
  avgResponseTimeMs: number
}

export interface DashboardData {
  health: SystemHealth
  stats: SystemStats
  recentLogs: any[]
  chartData: ChartData
}

export interface ChartData {
  requestTrend: TimeSeriesData[]
  tokenUsage: CategoryData[]
  providerDistribution: CategoryData[]
  successRate: TimeSeriesData[]
}

export interface TimeSeriesData {
  timestamp: string
  value: number
}

export interface CategoryData {
  label: string
  value: number
}