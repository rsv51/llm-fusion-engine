// ProxyKey 相关类型定义

export interface ProxyKey {
  id: number
  userId: number
  key: string
  enabled: boolean
  allowedGroups: string // JSON array of group IDs
  groupBalancePolicy: string
  groupWeights: string // JSON object for weighted balancing
  rpmLimit: number
  tpmLimit: number
  createdAt: string
  updatedAt: string
}

export interface CreateProxyKeyRequest {
  key: string
  userId?: number
  enabled?: boolean
  allowedGroups?: string
  groupBalancePolicy?: string
  groupWeights?: string
  rpmLimit?: number
  tpmLimit?: number
}

export interface UpdateProxyKeyRequest {
  enabled?: boolean
  allowedGroups?: string
  groupBalancePolicy?: string
  groupWeights?: string
  rpmLimit?: number
  tpmLimit?: number
}