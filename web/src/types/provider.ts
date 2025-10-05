// 供应商相关类型定义

export interface Provider {
  id: number
  groupId: number
  providerType: string
  weight: number
  enabled: boolean
  baseUrl?: string
  timeout?: number
  maxRetries?: number
  healthStatus?: 'healthy' | 'unhealthy' | 'unknown'
  lastChecked?: string
  createdAt: string
  updatedAt: string
}

export interface CreateProviderRequest {
  groupId: number
  providerType: string
  weight?: number
  enabled?: boolean
  baseUrl?: string
  timeout?: number
  maxRetries?: number
}

export interface UpdateProviderRequest extends Partial<CreateProviderRequest> {}