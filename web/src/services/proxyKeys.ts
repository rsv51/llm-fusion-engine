// ProxyKey 管理服务

import api from './api'
import type {
  ProxyKey,
  CreateProxyKeyRequest,
  UpdateProxyKeyRequest,
  PaginationResponse,
  PaginationParams,
} from '../types'

export const proxyKeysApi = {
  // 获取代理密钥列表
  async getProxyKeys(params?: PaginationParams): Promise<PaginationResponse<ProxyKey>> {
    return api.get('/admin/proxy-keys', { params })
  },

  // 获取单个代理密钥详情
  async getProxyKey(id: number): Promise<ProxyKey> {
    return api.get(`/admin/proxy-keys/${id}`)
  },

  // 创建代理密钥
  async createProxyKey(data: CreateProxyKeyRequest): Promise<ProxyKey> {
    return api.post('/admin/proxy-keys', data)
  },

  // 更新代理密钥
  async updateProxyKey(id: number, data: UpdateProxyKeyRequest): Promise<ProxyKey> {
    return api.put(`/admin/proxy-keys/${id}`, data)
  },

  // 删除代理密钥
  async deleteProxyKey(id: number): Promise<void> {
    return api.delete(`/admin/proxy-keys/${id}`)
  },
}

export default proxyKeysApi