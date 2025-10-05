// 系统健康和统计 API 服务

import api from './api'
import type {
  SystemHealth,
  SystemStats,
  DashboardData,
} from '../types'

export const systemApi = {
  // 获取系统健康状态
  async getHealth(): Promise<SystemHealth> {
    return api.get('/admin/health')
  },

  // 获取系统统计信息
  async getStats(hours: number = 24): Promise<SystemStats> {
    return api.get('/admin/stats', { params: { hours } })
  },

  // 获取仪表盘数据
  async getDashboardData(): Promise<DashboardData> {
    const [health, stats] = await Promise.all([
      this.getHealth(),
      this.getStats(24),
    ])
    
    return {
      health,
      stats,
      recentLogs: [],
      chartData: {
        requestTrend: [],
        tokenUsage: [],
        providerDistribution: [],
        successRate: [],
      },
    }
  },

  // 用户认证
  async login(password: string): Promise<{ token: string }> {
    return api.post('/admin/login', { password })
  },

  // 用户登出
  async logout(): Promise<void> {
    return api.post('/admin/logout')
  },
}

export default systemApi