// 请求日志相关类型定义

export interface Log {
  id: string
  proxy_key: string
  model: string
  provider: string
  request_url: string
  response_status: number
  is_success: boolean
  latency: number // in milliseconds
  timestamp: string
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
  request_body?: string;
  response_body?: string;
}

export interface LogQueryParams {
  page?: number
  pageSize?: number
  model?: string
  provider?: string
  status?: number
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