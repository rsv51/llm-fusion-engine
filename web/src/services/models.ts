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

  // 获取提供商可用模型列表(可带类型,提高后端判定准确性)
  async getProviderModels(providerId: number, providerType?: string): Promise<ProviderModelsResponse> {
    const params = providerType ? { type: providerType } : undefined
    return api.get(`/admin/providers/${providerId}/models`, { params })
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