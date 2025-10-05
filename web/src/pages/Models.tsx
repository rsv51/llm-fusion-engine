import React from 'react'
import { Card } from '../components/ui'
import { Box, Settings, CheckCircle, XCircle } from 'lucide-react'

export const Models: React.FC = () => {
  // 模拟数据 - 后续可以从 API 获取
  const models = [
    { id: '1', name: 'gpt-4', provider: 'OpenAI', status: 'active', requests: 1250 },
    { id: '2', name: 'gpt-3.5-turbo', provider: 'OpenAI', status: 'active', requests: 5420 },
    { id: '3', name: 'claude-3-opus', provider: 'Anthropic', status: 'active', requests: 890 },
    { id: '4', name: 'claude-3-sonnet', provider: 'Anthropic', status: 'inactive', requests: 0 },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">模型配置</h1>
        <button className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          <Settings className="w-4 h-4" />
          配置模型
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {models.map((model) => (
          <Card key={model.id}>
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                  <Box className="w-6 h-6 text-blue-600" />
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900">{model.name}</h3>
                  <p className="text-sm text-gray-500">{model.provider}</p>
                </div>
              </div>
              {model.status === 'active' ? (
                <CheckCircle className="w-5 h-5 text-green-500" />
              ) : (
                <XCircle className="w-5 h-5 text-gray-400" />
              )}
            </div>
            <div className="mt-4 pt-4 border-t border-gray-200">
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-600">总请求数</span>
                <span className="font-medium text-gray-900">{model.requests.toLocaleString()}</span>
              </div>
            </div>
          </Card>
        ))}
      </div>

      <Card>
        <h2 className="text-lg font-semibold text-gray-900 mb-4">模型使用统计</h2>
        <div className="space-y-3">
          {models.map((model) => (
            <div key={model.id} className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="text-sm font-medium text-gray-900">{model.name}</span>
                <span className="text-sm text-gray-500">({model.provider})</span>
              </div>
              <div className="flex items-center gap-4">
                <div className="w-48 bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-blue-600 h-2 rounded-full" 
                    style={{ width: `${(model.requests / 5420) * 100}%` }}
                  />
                </div>
                <span className="text-sm font-medium text-gray-900 w-16 text-right">
                  {model.requests}
                </span>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  )
}

export default Models