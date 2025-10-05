// API 密钥相关类型定义

export interface ApiKey {
  id: number
  groupId: number
  groupName: string
  key: string
  description?: string
  status: 'valid' | 'invalid' | 'untested'
  enabled: boolean
  lastValidatedAt?: string
  usageCount: number
  errorCount: number
  createdAt: string
  updatedAt: string
}

export interface CreateKeyRequest {
  groupId: number
  key: string
  description?: string
  enabled?: boolean
}

export interface UpdateKeyRequest {
  description?: string
  enabled?: boolean
}

export interface KeyValidationResult {
  keyId: number
  isValid: boolean
  responseTimeMs?: number
  errorMessage?: string
}

export interface BatchValidationRequest {
  keyIds: number[]
}

export interface BatchValidationResponse {
  results: KeyValidationResult[]
  totalValid: number
  totalInvalid: number
}

export interface KeyUsageStats {
  keyId: number
  usageCount: number
  errorCount: number
  lastUsedAt: string
  avgResponseTimeMs: number
}