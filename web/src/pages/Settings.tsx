import React, { useState } from 'react'
import { Card, Button } from '../components/ui'
import { Save, RefreshCw, Database, Shield, Bell, Globe, Upload, Download } from 'lucide-react'

export const Settings: React.FC = () => {
  const [saving, setSaving] = useState(false)
  const [activeTab, setActiveTab] = useState('general')
  const [importing, setImporting] = useState(false)
  const [importFile, setImportFile] = useState<File | null>(null)
  const [importType, setImportType] = useState('groups')

  const handleSave = async () => {
    setSaving(true)
    // 模拟保存操作
    await new Promise(resolve => setTimeout(resolve, 1000))
    setSaving(false)
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setImportFile(e.target.files[0])
    }
  }

  const handleImport = async () => {
    if (!importFile) {
      alert('请选择要导入的文件')
      return
    }
    setImporting(true)
    const formData = new FormData()
    formData.append('file', importFile)
    
    try {
      const response = await fetch(`/api/admin/import/${importType}`, {
        method: 'POST',
        body: formData,
      })
      const result = await response.json()
      if (response.ok) {
        alert(result.message)
      } else {
        throw new Error(result.error || '导入失败')
      }
    } catch (error) {
      console.error('导入失败:', error)
      if (error instanceof Error) {
        alert(`导入失败: ${error.message}`)
      } else {
        alert('导入失败: 发生未知错误')
      }
    } finally {
      setImporting(false)
      setImportFile(null)
    }
  }

  const handleExport = (dataType: string, format: 'csv' | 'xlsx') => {
    window.location.href = `/api/admin/export/${dataType}?format=${format}`
  }

  const renderGeneralSettings = () => (
    <>
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
    </>
  )

  const renderImportExport = () => (
    <div className="space-y-6">
      {/* 数据导出 */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <Download className="w-5 h-5 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">数据导出</h2>
        </div>
        <div className="space-y-4">
          <div className="flex items-center justify-between p-4 border rounded-lg">
            <div>
              <p className="font-medium text-gray-900">导出分组数据</p>
              <p className="text-sm text-gray-500">将所有分组配置导出为 CSV 或 XLSX 文件</p>
            </div>
            <div className="flex gap-2">
              <Button variant="secondary" onClick={() => handleExport('groups', 'csv')}>CSV</Button>
              <Button variant="secondary" onClick={() => handleExport('groups', 'xlsx')}>XLSX</Button>
            </div>
          </div>
          <div className="flex items-center justify-between p-4 border rounded-lg">
            <div>
              <p className="font-medium text-gray-900">导出密钥数据</p>
              <p className="text-sm text-gray-500">将所有密钥信息导出为 CSV 或 XLSX 文件</p>
            </div>
            <div className="flex gap-2">
              <Button variant="secondary" onClick={() => handleExport('keys', 'csv')}>CSV</Button>
              <Button variant="secondary" onClick={() => handleExport('keys', 'xlsx')}>XLSX</Button>
            </div>
          </div>
        </div>
      </Card>

      {/* 数据导入 */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <Upload className="w-5 h-5 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">数据导入</h2>
        </div>
        <div className="p-4 border rounded-lg space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              选择导入类型
            </label>
            <select
              value={importType}
              onChange={(e) => setImportType(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="groups">分组</option>
              <option value="keys">密钥</option>
              <option value="providers">供应商</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              选择文件 (.csv, .xlsx)
            </label>
            <div className="flex items-center gap-4">
              <input
                type="file"
                onChange={handleFileChange}
                accept=".csv, .xlsx"
                className="flex-1 text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
              />
              <Button onClick={handleImport} disabled={!importFile || importing}>
                {importing ? '导入中...' : '开始导入'}
              </Button>
            </div>
            {importFile && <p className="text-sm text-gray-500 mt-2">已选择文件: {importFile.name}</p>}
          </div>
        </div>
      </Card>
    </div>
  )

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

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8" aria-label="Tabs">
          <button
            onClick={() => setActiveTab('general')}
            className={`whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'general'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            通用设置
          </button>
          <button
            onClick={() => setActiveTab('import-export')}
            className={`whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'import-export'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            导入/导出
          </button>
        </nav>
      </div>

      {/* Tab Content */}
      <div className="pt-6">
        {activeTab === 'general' && renderGeneralSettings()}
        {activeTab === 'import-export' && renderImportExport()}
      </div>
    </div>
  )
}

export default Settings