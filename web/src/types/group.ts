// 分组相关类型定义

export interface Group {
  id: number
  name: string
  description?: string
  priority: number
  enabled: boolean
  loadBalanceStrategy: 'round_robin' | 'weighted' | 'random' | 'failover'
  healthCheckInterval: number
  maxRetries: number
  timeout: number
  proxyUrl?: string
  createdAt: string
  updatedAt: string
}

export interface CreateGroupRequest {
  name: string
  description?: string
  priority?: number
  enabled?: boolean
  loadBalanceStrategy?: Group['loadBalanceStrategy']
  healthCheckInterval?: number
  maxRetries?: number
  timeout?: number
  proxyUrl?: string
}

export interface UpdateGroupRequest extends Partial<CreateGroupRequest> {}

export interface GroupWithProviders extends Group {
  providers: GroupProvider[]
}

export interface GroupProvider {
  id: number
  groupId: number
  providerId: number
  providerName: string
  providerType: string
  apiKey: string
  baseUrl?: string
  weight: number
  enabled: boolean
  priority: number
  createdAt: string
}

export interface GroupHealthStatus {
  groupId: number
  groupName: string
  isHealthy: boolean
  responseTimeMs: number
  lastCheckAt: string
  providers: ProviderHealthStatus[]
}

export interface ProviderHealthStatus {
  providerId: number
  providerName: string
  isHealthy: boolean
  responseTimeMs: number
  errorMessage?: string
}