// 请求日志 API 服务

import api from './api'
import type {
  Log,
  LogQueryParams,
  LogStats,
  PaginationResponse,
} from '../types'

export const logsApi = {
  // 获取日志列表
  async getLogs(params?: LogQueryParams): Promise<PaginationResponse<Log>> {
    return api.get('/admin/logs', { params })
  },

  // 获取单条日志
  async getLog(id: string): Promise<Log> {
    return api.get(`/admin/logs/${id}`);
  },

  // 获取日志统计
  async getLogStats(hours: number = 24): Promise<LogStats> {
    return api.get('/admin/logs/stats', { params: { hours } })
  },

  // 导出日志
  async exportLogs(params?: LogQueryParams): Promise<Blob> {
    return api.get('/admin/logs/export', {
      params,
      responseType: 'blob',
    })
  },

  // 清空日志
  async clearLogs(beforeDate?: string): Promise<void> {
    return api.delete('/admin/logs', { params: { beforeDate } })
  },
}

export default logsApi