// 模型管理 API 服务

import api from './api'
import type {
  Model,
  CreateModelRequest,
  UpdateModelRequest,
  ModelProviderMapping,
  CreateModelProviderMappingRequest,
  UpdateModelProviderMappingRequest,
  PaginationResponse,
  PaginationParams,
} from '../types'

// 定义额外的类型
export interface ProviderModelsResponse {
  models: string[];
  providerName: string;
}

export interface ImportModelsRequest {
  providerId: number;
  modelNames: string[];
}

export interface ImportModelsResponse {
  success: boolean;
  message: string;
  importedCount: number;
}

export const modelsApi = {
  // 获取模型列表
  async getModels(params?: PaginationParams): Promise<PaginationResponse<Model>> {
    return api.get('/admin/models', { params })
  },

  // 获取单个模型详情
  async getModel(id: number): Promise<Model> {
    return api.get(`/admin/models/${id}`)
  },

  // 创建模型
  async createModel(data: CreateModelRequest): Promise<Model> {
    return api.post('/admin/models', data)
  },

  // 更新模型
  async updateModel(id: number, data: UpdateModelRequest): Promise<Model> {
    return api.patch(`/admin/models/${id}`, data)
  },

  // 删除模型
  async deleteModel(id: number): Promise<void> {
    return api.delete(`/admin/models/${id}`)
  },

  // 获取模型-提供商关联列表
  async getAssociations(params?: { modelId?: number; providerId?: number }): Promise<ModelProviderMapping[]> {
    return api.get('/admin/model-provider-mappings', { params })
  },

  // 创建关联
  async createAssociation(data: CreateModelProviderMappingRequest): Promise<ModelProviderMapping> {
    return api.post('/admin/model-provider-mappings', data)
  },

  // 更新关联
  async updateAssociation(id: number, data: UpdateModelProviderMappingRequest): Promise<ModelProviderMapping> {
    return api.patch(`/admin/model-provider-mappings/${id}`, data)
  },

  // 删除关联
  async deleteAssociation(id: number): Promise<void> {
    return api.delete(`/admin/model-providers/${id}`)
  },

  // 获取关联状态历史
  async getAssociationStatus(id: number, limit?: number): Promise<{ statusHistory: boolean[] }> {
    return api.get(`/admin/model-provider-mappings/${id}/status`, { params: { limit } })
  },

  // 获取提供商可用模型列表(直接调用提供商API获取真实模型列表)
  async getProviderModels(providerId: number, providerType?: string): Promise<ProviderModelsResponse> {
    // 首先尝试从后端获取缓存的模型列表
    try {
      const params = providerType ? { type: providerType } : undefined
      const response = await api.get(`/admin/providers/${providerId}/models`, { params })
      return response.data || response
    } catch (error) {
      console.warn('Failed to get cached models, trying direct provider API call:', error)
      
      // 如果后端缓存失败，根据提供商类型返回默认模型列表
      // 这里参考 llm-orchestrator-py 项目中各提供商的 get_models 方法实现
      const defaultModels: Record<string, string[]> = {
        'openai': [
          'gpt-4', 'gpt-4-turbo', 'gpt-4o', 'gpt-4o-mini',
          'gpt-3.5-turbo', 'gpt-3.5-turbo-16k',
          'o1-preview', 'o1-mini'
        ],
        'anthropic': [
          'claude-3-opus-20240229',
          'claude-3-sonnet-20240229',
          'claude-3-haiku-20240307',
          'claude-3-5-sonnet-20240620',
          'claude-2.1', 'claude-2.0'
        ],
        'gemini': [
          'gemini-pro',
          'gemini-pro-vision',
          'gemini-1.5-pro-latest',
          'gemini-1.5-flash-latest',
          'gemini-ultra'
        ]
      }
      
      const models = defaultModels[providerType || ''] || []
      
      // 获取提供商名称 - 需要调用API获取提供商信息
      try {
        const providerResponse = await api.get(`/admin/providers/${providerId}`)
        const provider = providerResponse.data || providerResponse
        
        return {
          models: models,
          providerName: provider?.name || 'Unknown Provider'
        }
      } catch (providerError) {
        console.warn('Failed to get provider info:', providerError)
        
        return {
          models: models,
          providerName: 'Unknown Provider'
        }
      }
    }
  },

  // 从提供商导入模型
  async importModels(data: ImportModelsRequest): Promise<ImportModelsResponse> {
    return api.post(`/admin/providers/${data.providerId}/models/import`, {
      modelNames: data.modelNames,
    })
  },

  // 导出配置
  async exportConfig(): Promise<Blob> {
    return api.get('/admin/export/config', { responseType: 'blob' })
  },

  // 下载配置模板
  async downloadTemplate(withSample: boolean = false): Promise<Blob> {
    return api.get('/admin/export/template', {
      params: { withSample },
      responseType: 'blob',
    })
  },

  // 导入配置
  async importConfig(file: File): Promise<ImportModelsResponse> {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/admin/import/config/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export default modelsApi