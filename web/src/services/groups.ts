// 分组管理 API 服务

import api from './api'
import type {
  Group,
  CreateGroupRequest,
  UpdateGroupRequest,
  GroupWithProviders,
  GroupHealthStatus,
  PaginationResponse,
  PaginationParams,
} from '../types'

export const groupsApi = {
  // 获取分组列表
  async getGroups(params?: PaginationParams): Promise<PaginationResponse<Group>> {
    return api.get('/admin/groups', { params })
  },

  // 获取单个分组详情
  async getGroup(id: number): Promise<GroupWithProviders> {
    return api.get(`/admin/groups/${id}`)
  },

  // 创建分组
  async createGroup(data: CreateGroupRequest): Promise<Group> {
    return api.post('/admin/groups', data)
  },

  // 更新分组
  async updateGroup(id: number, data: UpdateGroupRequest): Promise<Group> {
  	return api.put(`/admin/groups/${id}`, data)
  },

  // 删除分组
  async deleteGroup(id: number): Promise<void> {
    return api.delete(`/admin/groups/${id}`)
  },

  // 启用/禁用分组
  async toggleGroup(id: number, enabled: boolean): Promise<Group> {
    return api.patch(`/admin/groups/${id}`, { enabled })
  },

  // 批量删除分组
  async batchDeleteGroups(ids: number[]): Promise<void> {
    return api.post('/admin/groups/batch-delete', { ids })
  },

  // 获取分组健康状态
  async getGroupHealth(id: number): Promise<GroupHealthStatus> {
    return api.get(`/admin/groups/${id}/health`)
  },

  // 检查所有分组健康状态
  async checkAllGroupsHealth(): Promise<GroupHealthStatus[]> {
    return api.post('/admin/groups/health-check')
  },
}

export default groupsApi