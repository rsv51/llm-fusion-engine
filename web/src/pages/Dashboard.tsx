import React, { useState, useEffect } from 'react'
import { TrendingUp, TrendingDown, Activity, Key, Zap, AlertCircle } from 'lucide-react'
import { Card, Badge } from '../components/ui'
import { systemApi, logsApi } from '../services'
import type { SystemStats } from '../types'

interface DisplayLog {
  id: string
  model: string
  providerName: string
  statusCode: number
  createdAt: string
}

export const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<SystemStats | null>(null)
  const [recentLogs, setRecentLogs] = useState<DisplayLog[]>([])
  const [loading, setLoading] = useState(true)
  const [uptime, setUptime] = useState('')

  useEffect(() => {
    loadData()
    const interval = setInterval(loadData, 30000) // 每30秒刷新
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (stats?.startTime) {
      const updateUptime = () => {
        const startTime = new Date(stats.startTime)
        const now = new Date()
        const diff = now.getTime() - startTime.getTime()

        const days = Math.floor(diff / (1000 * 60 * 60 * 24))
        const hours = Math.floor((diff / (1000 * 60 * 60)) % 24)
        const minutes = Math.floor((diff / 1000 / 60) % 60)

        setUptime(`${days}d ${hours}h ${minutes}m`)
      }

      updateUptime()
      const uptimeInterval = setInterval(updateUptime, 60000) // 每分钟更新
      return () => clearInterval(uptimeInterval)
    }
  }, [stats])

  const loadData = async () => {
    try {
      const [statsData, logsData] = await Promise.all([
        systemApi.getStats(),
        logsApi.getLogs({ page: 1, pageSize: 5 })
      ])
      setStats(statsData)
      setRecentLogs(logsData.data.map((log: any) => ({
      	id: log.id,
      	model: log.model,
      	providerName: log.providerName,
      	statusCode: log.statusCode,
      	createdAt: log.createdAt
      })))
    } catch (error) {
      console.error('加载数据失败:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  const statCards = [
    {
      title: '总请求数',
      value: stats?.totalRequests || 0,
      icon: Activity,
      color: 'blue'
    },
    {
      title: '活跃密钥',
      value: stats?.activeKeys || 0,
      icon: Key,
      color: 'green'
    },
    {
      title: '平均响应时间',
      value: `${stats?.avgResponseTimeMs?.toFixed(0) || 0}ms`,
      icon: Zap,
      color: 'purple'
    },
    {
      title: '成功率',
      value: `${stats?.successRate?.toFixed(1) || 0}%`,
      icon: AlertCircle,
      color: 'green'
    }
  ]

  return (
    <div className="space-y-6">
      {/* 页面标题 */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">仪表盘</h1>
        <p className="text-sm text-gray-500 mt-1">实时监控系统运行状态</p>
      </div>

      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {statCards.map((card) => {
          const Icon = card.icon
          return (
            <Card key={card.title} className="p-6">
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <p className="text-sm text-gray-500">{card.title}</p>
                  <p className="text-2xl font-bold text-gray-900 mt-2">{card.value}</p>
                </div>
                <div className={`w-12 h-12 rounded-full bg-${card.color}-100 flex items-center justify-center`}>
                  <Icon className={`w-6 h-6 text-${card.color}-600`} />
                </div>
              </div>
            </Card>
          )
        })}
      </div>

      {/* 提供商状态 */}
      <Card className="p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">提供商状态</h2>
        <div className="space-y-3">
          {stats?.providers && stats.providers.length > 0 ? (
            stats.providers.map((provider) => {
              const successRate = provider.requestCount > 0
                ? ((provider.successCount / provider.requestCount) * 100).toFixed(1)
                : '0.0'
              
              return (
                <div key={provider.providerId} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className={`w-2 h-2 rounded-full ${provider.errorCount === 0 ? 'bg-green-500' : 'bg-yellow-500'}`}></div>
                    <div>
                      <p className="font-medium text-gray-900">{provider.providerName}</p>
                      <p className="text-xs text-gray-500">响应时间: {provider.avgResponseTimeMs}ms</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium text-gray-900">{provider.requestCount} 次请求</p>
                    <Badge variant={provider.errorCount === 0 ? 'success' : 'warning'}>
                      成功率 {successRate}%
                    </Badge>
                  </div>
                </div>
              )
            })
          ) : (
            <p className="text-sm text-gray-500 text-center py-4">暂无数据</p>
          )}
        </div>
      </Card>

      {/* 最近日志 */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">最近日志</h2>
          <a href="/logs" className="text-sm text-blue-600 hover:text-blue-700">查看全部</a>
        </div>
        <div className="space-y-2">
          {recentLogs.length > 0 ? (
            recentLogs.map((log) => (
              <div key={log.id} className="flex items-center justify-between p-3 border-l-2 border-gray-200 hover:border-blue-500 hover:bg-gray-50 transition-colors">
                <div>
                  <p className="text-sm text-gray-900">{log.model} · {log.providerName}</p>
                  <p className="text-xs text-gray-500">{new Date(log.createdAt).toLocaleString()}</p>
                </div>
                <Badge variant={log.statusCode < 400 ? 'success' : 'error'}>
                  {log.statusCode}
                </Badge>
              </div>
            ))
          ) : (
            <p className="text-sm text-gray-500 text-center py-4">暂无日志</p>
          )}
        </div>
      </Card>
    </div>
  )
}

export default Dashboard