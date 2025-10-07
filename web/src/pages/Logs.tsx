import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom';
import { Card, Badge } from '../components/ui'
import { logsApi } from '../services'
import type { Log } from '../types'
import { Search, Filter, RefreshCw } from 'lucide-react'

export const Logs: React.FC = () => {
  const [logs, setLogs] = useState<Log[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)

  const pickData = (resp: any) => {
    if (!resp) return { data: [], pagination: {} };
    if (Array.isArray(resp.data) && resp.pagination) return resp;
    if (Array.isArray(resp.items)) return { data: resp.items, pagination: { page: resp.page, totalPage: resp.totalPages } };
    if (Array.isArray(resp)) return { data: resp, pagination: {} };
    return { data: [], pagination: {} };
  };

  const fetchLogs = async (pageNum: number) => {
    setLoading(true)
    try {
      const response = await logsApi.getLogs({
        page: pageNum,
        pageSize: 20
      })
      const { data, pagination } = pickData(response);
      setLogs(data || [])
      setPage(pagination?.page || 1)
      setTotalPages(pagination?.totalPage || 1)
    } catch (error) {
      console.error('Failed to fetch logs:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchLogs(1)
  }, [])

  const getStatusBadge = (isSuccess: boolean) => {
    if (isSuccess) {
      return <Badge variant="success">成功</Badge>
    }
    return <Badge variant="error">失败</Badge>
  }

  const formatDate = (date: string) => {
    return new Date(date).toLocaleString('zh-CN')
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">请求日志</h1>
        <button
          onClick={() => fetchLogs(page)}
          disabled={loading}
          className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
          刷新
        </button>
      </div>

      {/* 筛选栏 */}
      <Card>
        <div className="flex gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="搜索模型、分组..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
          <button className="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">
            <Filter className="w-4 h-4" />
            筛选
          </button>
        </div>
      </Card>

      {/* 日志列表 */}
      <Card>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">时间</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">模型</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">供应商</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">状态</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">延迟</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">Token</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                    加载中...
                  </td>
                </tr>
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                    暂无日志记录
                  </td>
                </tr>
              ) : (
                logs.map((log) => (
                  <tr key={log.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="px-4 py-3 text-sm text-gray-900">
                      <Link to={`/logs/${log.id}`} className="text-blue-600 hover:underline">
                        {formatDate(log.timestamp)}
                      </Link>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-900">{log.model}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">{log.provider}</td>
                    <td className="px-4 py-3 text-sm">{getStatusBadge(log.is_success)}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">{log.latency}ms</td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {log.total_tokens}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* 分页 */}
        <div className="flex justify-center items-center gap-4 px-4 py-3 border-t border-gray-200">
          <button
            onClick={() => fetchLogs(page - 1)}
            disabled={page <= 1}
            className="px-3 py-1 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
          >
            上一页
          </button>
          <span className="px-3 py-1 text-sm text-gray-700">
            第 {page} / {totalPages} 页
          </span>
          <button
            onClick={() => fetchLogs(page + 1)}
            disabled={page >= totalPages}
            className="px-3 py-1 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50"
          >
            下一页
          </button>
        </div>
      </Card>
    </div>
  )
}

export default Logs