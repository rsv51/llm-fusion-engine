import React, { useState } from 'react'
import { Download, Upload, FileText, AlertCircle, CheckCircle, X } from 'lucide-react'
import { Card, Button, Modal, Badge } from '../components/ui'
import { api } from '../services'

interface ImportResult {
  filename: string
  result: {
    providers: {
      total: number
      imported: number
      skipped: number
      errors: Array<{
        row: number
        field: string
        error: string
      }>
    }
    models: {
      total: number
      imported: number
      skipped: number
      errors: Array<{
        row: number
        field: string
        error: string
      }>
    }
    associations: {
      total: number
      imported: number
      skipped: number
      errors: Array<{
        row: number
        field: string
        error: string
      }>
    }
    summary: {
      total_imported: number
      total_skipped: number
      total_errors: number
    }
  }
}

export const ImportExport: React.FC = () => {
  const [isUploading, setIsUploading] = useState(false)
  const [importResult, setImportResult] = useState<ImportResult | null>(null)
  const [showResultModal, setShowResultModal] = useState(false)

  const downloadTemplate = async (withSample: boolean = false) => {
    try {
      const response = await api.get(`/api/admin/export/template?with_sample=${withSample}`, {
        responseType: 'blob'
      })
      
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `llm_fusion_engine_template${withSample ? '_with_sample' : ''}.xlsx`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('下载模板失败:', error)
      alert('下载模板失败，请重试')
    }
  }

  const exportConfig = async (format: string = 'excel') => {
    try {
      const response = await api.get(`/api/admin/export/all?format=${format}`, {
        responseType: 'blob'
      })
      
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)
      const extension = format === 'excel' ? 'xlsx' : format
      const filename = `llm_fusion_engine_config_${timestamp}.${extension}`
      
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('导出配置失败:', error)
      alert('导出配置失败，请重试')
    }
  }

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    // Check file type
    const allowedTypes = [
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', // .xlsx
      'application/json',
      'application/x-yaml',
      'text/yaml'
    ]
    
    if (!allowedTypes.includes(file.type) && !file.name.endsWith('.xlsx') && !file.name.endsWith('.json') && !file.name.endsWith('.yaml') && !file.name.endsWith('.yml')) {
      alert('不支持的文件类型，请上传 .xlsx、.json、.yaml 或 .yml 文件')
      return
    }

    setIsUploading(true)
    const formData = new FormData()
    formData.append('file', file)

    try {
      const response = await api.post('/api/admin/import/all', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      
      setImportResult(response.data)
      setShowResultModal(true)
    } catch (error: any) {
      console.error('导入失败:', error)
      const errorMessage = error.response?.data?.error || error.message || '导入失败，请重试'
      alert(errorMessage)
    } finally {
      setIsUploading(false)
      // Reset file input
      if (event.target) {
        (event.target as HTMLInputElement).value = ''
      }
    }
  }

  const renderResultSummary = (result: ImportResult['result']) => (
    <div className="space-y-4">
      {/* Summary */}
      <div className="grid grid-cols-3 gap-4">
        <Card className="p-4 text-center">
          <div className="text-2xl font-bold text-green-600">{result.summary.total_imported}</div>
          <div className="text-sm text-gray-600">成功导入</div>
        </Card>
        <Card className="p-4 text-center">
          <div className="text-2xl font-bold text-yellow-600">{result.summary.total_skipped}</div>
          <div className="text-sm text-gray-600">跳过重复</div>
        </Card>
        <Card className="p-4 text-center">
          <div className="text-2xl font-bold text-red-600">{result.summary.total_errors}</div>
          <div className="text-sm text-gray-600">错误数量</div>
        </Card>
      </div>

      {/* Providers */}
      <Card className="p-4">
        <h3 className="font-semibold mb-2 flex items-center">
          <FileText className="w-4 h-4 mr-2" />
          提供商 ({result.providers.imported}/{result.providers.total})
        </h3>
        <div className="text-sm text-gray-600 mb-2">
          导入: {result.providers.imported}, 跳过: {result.providers.skipped}, 错误: {result.providers.errors.length}
        </div>
        {result.providers.errors.length > 0 && (
          <div className="space-y-1">
            {result.providers.errors.slice(0, 5).map((error, index) => (
              <div key={index} className="text-xs text-red-600 bg-red-50 p-2 rounded">
                第 {error.row} 行 ({error.field}): {error.error}
              </div>
            ))}
            {result.providers.errors.length > 5 && (
              <div className="text-xs text-gray-500">
                还有 {result.providers.errors.length - 5} 个错误...
              </div>
            )}
          </div>
        )}
      </Card>

      {/* Models */}
      <Card className="p-4">
        <h3 className="font-semibold mb-2 flex items-center">
          <FileText className="w-4 h-4 mr-2" />
          模型 ({result.models.imported}/{result.models.total})
        </h3>
        <div className="text-sm text-gray-600 mb-2">
          导入: {result.models.imported}, 跳过: {result.models.skipped}, 错误: {result.models.errors.length}
        </div>
        {result.models.errors.length > 0 && (
          <div className="space-y-1">
            {result.models.errors.slice(0, 5).map((error, index) => (
              <div key={index} className="text-xs text-red-600 bg-red-50 p-2 rounded">
                第 {error.row} 行 ({error.field}): {error.error}
              </div>
            ))}
            {result.models.errors.length > 5 && (
              <div className="text-xs text-gray-500">
                还有 {result.models.errors.length - 5} 个错误...
              </div>
            )}
          </div>
        )}
      </Card>

      {/* Associations */}
      <Card className="p-4">
        <h3 className="font-semibold mb-2 flex items-center">
          <FileText className="w-4 h-4 mr-2" />
          关联 ({result.associations.imported}/{result.associations.total})
        </h3>
        <div className="text-sm text-gray-600 mb-2">
          导入: {result.associations.imported}, 跳过: {result.associations.skipped}, 错误: {result.associations.errors.length}
        </div>
        {result.associations.errors.length > 0 && (
          <div className="space-y-1">
            {result.associations.errors.slice(0, 5).map((error, index) => (
              <div key={index} className="text-xs text-red-600 bg-red-50 p-2 rounded">
                第 {error.row} 行 ({error.field}): {error.error}
              </div>
            ))}
            {result.associations.errors.length > 5 && (
              <div className="text-xs text-gray-500">
                还有 {result.associations.errors.length - 5} 个错误...
              </div>
            )}
          </div>
        )}
      </Card>
    </div>
  )

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">导入导出管理</h1>
          <p className="text-sm text-gray-500 mt-1">批量导入导出系统配置，支持Excel、JSON、YAML格式</p>
        </div>
      </div>

      {/* Export Section */}
      <Card className="p-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center">
          <Download className="w-5 h-5 mr-2" />
          导出配置
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Button
            onClick={() => downloadTemplate(false)}
            className="flex flex-col items-center h-auto py-4"
            variant="secondary"
          >
            <FileText className="w-8 h-8 mb-2" />
            <span className="text-sm">空白模板</span>
            <span className="text-xs text-gray-500">下载Excel模板</span>
          </Button>
          
          <Button
            onClick={() => downloadTemplate(true)}
            className="flex flex-col items-center h-auto py-4"
            variant="secondary"
          >
            <FileText className="w-8 h-8 mb-2" />
            <span className="text-sm">示例模板</span>
            <span className="text-xs text-gray-500">含示例数据</span>
          </Button>
          
          <Button
            onClick={() => exportConfig('excel')}
            className="flex flex-col items-center h-auto py-4"
          >
            <Download className="w-8 h-8 mb-2" />
            <span className="text-sm">导出Excel</span>
            <span className="text-xs text-gray-500">当前配置</span>
          </Button>
          
          <Button
            onClick={() => exportConfig('json')}
            className="flex flex-col items-center h-auto py-4"
            variant="secondary"
          >
            <Download className="w-8 h-8 mb-2" />
            <span className="text-sm">导出JSON</span>
            <span className="text-xs text-gray-500">备份格式</span>
          </Button>
        </div>
      </Card>

      {/* Import Section */}
      <Card className="p-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center">
          <Upload className="w-5 h-5 mr-2" />
          导入配置
        </h2>
        
        <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
          <Upload className="w-12 h-12 mx-auto text-gray-400 mb-4" />
          <p className="text-lg font-medium text-gray-700 mb-2">上传配置文件</p>
          <p className="text-sm text-gray-500 mb-4">
            支持 .xlsx、.json、.yaml、.yml 格式
          </p>
          
          <input
            type="file"
            id="file-upload"
            className="hidden"
            accept=".xlsx,.json,.yaml,.yml"
            onChange={handleFileUpload}
            disabled={isUploading}
          />
          
          <label
            htmlFor="file-upload"
            className={`inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 cursor-pointer ${
              isUploading ? 'opacity-50 cursor-not-allowed' : ''
            }`}
          >
            {isUploading ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                导入中...
              </>
            ) : (
              <>
                <Upload className="w-4 h-4 mr-2" />
                选择文件
              </>
            )}
          </label>
        </div>

        <div className="mt-4 p-4 bg-blue-50 rounded-lg">
          <h3 className="font-medium text-blue-900 mb-2 flex items-center">
            <AlertCircle className="w-4 h-4 mr-2" />
            导入说明
          </h3>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>• Excel文件必须包含三个工作表：Providers、Models、Associations</li>
            <li>• 系统会自动跳过重复的数据（按名称去重）</li>
            <li>• 建议先下载模板，按格式填写数据</li>
            <li>• 导入前建议先导出当前配置作为备份</li>
          </ul>
        </div>
      </Card>

      {/* Result Modal */}
      <Modal
        isOpen={showResultModal}
        onClose={() => setShowResultModal(false)}
        title="导入结果"
        size="lg"
      >
        {importResult && renderResultSummary(importResult.result)}
        
        <div className="flex justify-end gap-3 pt-4">
          <Button
            variant="secondary"
            onClick={() => setShowResultModal(false)}
          >
            关闭
          </Button>
          {importResult?.result.summary.total_errors === 0 && (
            <Badge variant="success" className="flex items-center">
              <CheckCircle className="w-4 h-4 mr-1" />
              导入成功
            </Badge>
          )}
        </div>
      </Modal>
    </div>
  )
}

export default ImportExport