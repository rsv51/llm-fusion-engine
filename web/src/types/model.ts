// 模型相关类型定义

export interface Model {
  id: number
  name: string
  description?: string
  maxRetry: number
  timeout: number
  enabled: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateModelRequest {
  name: string
  description?: string
  maxRetry?: number
  timeout?: number
  enabled?: boolean
}

export interface UpdateModelRequest extends Partial<CreateModelRequest> {}

export interface ModelProviderAssociation {
  id: number
  modelId: number
  modelName: string
  providerId: number
  providerName: string
  providerType: string
  providerModel: string
  weight: number
  toolCall: boolean
  structuredOutput: boolean
  image: boolean
  enabled: boolean
  statusHistory?: boolean[]
  createdAt: string
  updatedAt: string
}

export interface CreateAssociationRequest {
  modelId: number
  providerId: number
  providerModel: string
  weight?: number
  toolCall?: boolean
  structuredOutput?: boolean
  image?: boolean
  enabled?: boolean
}

export interface UpdateAssociationRequest extends Partial<CreateAssociationRequest> {}

export interface ProviderModelsResponse {
  providerId: number
  providerName: string
  models: string[]
}

export interface ImportModelsRequest {
  providerId: number
  modelNames: string[] | null // null = import all
}

export interface ImportModelsResponse {
  created: number
  skipped: number
  total: number
  errors: ImportError[]
}

export interface ImportError {
  row: number
  field: string
  error: string
}