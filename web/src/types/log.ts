// 请求日志相关类型定义

export interface RequestLog {
  id: string
  model: string
  providerName: string
  providerType: string
  statusCode: number
  latencyMs: number
  totalTokens: number
  promptTokens: number
  completionTokens: number
  errorMessage?: string
  createdAt: string
}

export interface LogQueryParams {
  page?: number
  pageSize?: number
  model?: string
  providerName?: string
  status?: 'success' | 'error'
  startDate?: string
  endDate?: string
}

export interface LogStats {
  totalRequests: number
  successRate: number
  avgLatencyMs: number
  totalTokens: number
  byModel: Record<string, number>
  byProvider: Record<string, number>
  byStatus: Record<string, number>
}