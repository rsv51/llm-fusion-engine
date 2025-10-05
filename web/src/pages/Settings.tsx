import React, { useState } from 'react'
import { Card } from '../components/ui'
import { Save, RefreshCw, Database, Shield, Bell, Globe } from 'lucide-react'

export const Settings: React.FC = () => {
  const [saving, setSaving] = useState(false)

  const handleSave = async () => {
    setSaving(true)
    // 模拟保存操作
    await new Promise(resolve => setTimeout(resolve, 1000))
    setSaving(false)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">系统设置</h1>
        <button
          onClick={handleSave}
          disabled={saving}
          className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          {saving ? (
            <RefreshCw className="w-4 h-4 animate-spin" />
          ) : (
            <Save className="w-4 h-4" />
          )}
          {saving ? '保存中...' : '保存设置'}
        </button>
      </div>

      {/* 数据库设置 */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <Database className="w-5 h-5 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">数据库设置</h2>
        </div>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              数据库类型
            </label>
            <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
              <option value="sqlite">SQLite</option>
              <option value="postgres">PostgreSQL</option>
              <option value="mysql">MySQL</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              连接字符串
            </label>
            <input
              type="text"
              defaultValue="fusion.db"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </Card>

      {/* 安全设置 */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <Shield className="w-5 h-5 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">安全设置</h2>
        </div>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium text-gray-900">启用 API 密钥认证</p>
              <p className="text-sm text-gray-500">要求所有请求提供有效的 API 密钥</p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" className="sr-only peer" defaultChecked />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium text-gray-900">启用请求日志</p>
              <p className="text-sm text-gray-500">记录所有 API 请求以供审计</p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" className="sr-only peer" defaultChecked />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              日志保留天数
            </label>
            <input
              type="number"
              defaultValue="30"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </Card>

      {/* 通知设置 */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <Bell className="w-5 h-5 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">通知设置</h2>
        </div>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium text-gray-900">服务异常通知</p>
              <p className="text-sm text-gray-500">当服务不可用时发送通知</p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" className="sr-only peer" defaultChecked />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium text-gray-900">配额预警</p>
              <p className="text-sm text-gray-500">当使用量接近配额时发送通知</p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" className="sr-only peer" />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>
        </div>
      </Card>

      {/* 系统信息 */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <Globe className="w-5 h-5 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">系统信息</h2>
        </div>
        <div className="space-y-3">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">版本号</span>
            <span className="font-medium text-gray-900">v1.0.0</span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">运行时间</span>
            <span className="font-medium text-gray-900">2 天 14 小时</span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">Go 版本</span>
            <span className="font-medium text-gray-900">1.21.0</span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">数据库大小</span>
            <span className="font-medium text-gray-900">24.5 MB</span>
          </div>
        </div>
      </Card>
    </div>
  )
}

export default Settings