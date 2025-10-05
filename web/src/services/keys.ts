// API 密钥管理服务

import api from './api'
import type {
  ApiKey,
  CreateKeyRequest,
  UpdateKeyRequest,
  KeyValidationResult,
  BatchValidationRequest,
  BatchValidationResponse,
  KeyUsageStats,
  PaginationResponse,
  PaginationParams,
} from '../types'

export const keysApi = {
  // 获取密钥列表
  async getKeys(params?: PaginationParams & { groupId?: number }): Promise<PaginationResponse<ApiKey>> {
    return api.get('/admin/keys', { params })
  },

  // 获取单个密钥详情
  async getKey(id: number): Promise<ApiKey> {
    return api.get(`/admin/keys/${id}`)
  },

  // 创建密钥
  async createKey(data: CreateKeyRequest): Promise<ApiKey> {
    return api.post('/admin/keys', data)
  },

  // 更新密钥
  async updateKey(id: number, data: UpdateKeyRequest): Promise<ApiKey> {
    return api.patch(`/admin/keys/${id}`, data)
  },

  // 删除密钥
  async deleteKey(id: number): Promise<void> {
    return api.delete(`/admin/keys/${id}`)
  },

  // 验证单个密钥
  async validateKey(id: number): Promise<KeyValidationResult> {
    return api.post(`/admin/keys/${id}/validate`)
  },

  // 批量验证密钥
  async batchValidateKeys(data: BatchValidationRequest): Promise<BatchValidationResponse> {
    return api.post('/admin/keys/batch-validate', data)
  },

  // 批量删除密钥
  async batchDeleteKeys(ids: number[]): Promise<void> {
    return api.post('/admin/keys/batch-delete', { ids })
  },

  // 获取密钥使用统计
  async getKeyStats(id: number): Promise<KeyUsageStats> {
    return api.get(`/admin/keys/${id}/stats`)
  },

  // 导出密钥列表
  async exportKeys(groupId?: number): Promise<Blob> {
    return api.get('/admin/keys/export', {
      params: { groupId },
      responseType: 'blob',
    })
  },

  // 导入密钥
  async importKeys(file: File): Promise<{ created: number; skipped: number; errors: any[] }> {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/admin/keys/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export default keysApi