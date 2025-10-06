import React, { useState } from 'react'
import { Download, Upload, FileSpreadsheet, AlertCircle, CheckCircle, Info } from 'lucide-react'
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
    modelProviderMappings: {
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
  const [isExporting, setIsExporting] = useState(false)
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
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)
      link.setAttribute('download', `配置模板${withSample ? '_示例数据' : ''}_${timestamp}.xlsx`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('下载模板失败:', error)
      alert('下载模板失败，请重试')
    }
  }

  const exportConfig = async () => {
    setIsExporting(true)
    try {
      const response = await api.get('/api/admin/export/all', {
        responseType: 'blob'
      })
      
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)
      const filename = `系统配置_${timestamp}.xlsx`
      
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
    } finally {
      setIsExporting(false)
    }
  }

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    // Check file type - only accept Excel files
    if (!file.name.endsWith('.xlsx')) {
      alert('不支持的文件类型，请上传 .xlsx 文件')
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
      {/* Summary Cards */}
      <div className="grid grid-cols-3 gap-4">
        <Card className="p-4 text-center bg-green-50 border-green-200">
          <div className="text-3xl font-bold text-green-600">{result.summary.total_imported}</div>
          <div className="text-sm text-green-700 font-medium mt-1">成功导入</div>
        </Card>
        <Card className="p-4 text-center bg-yellow-50 border-yellow-200">
          <div className="text-3xl font-bold text-yellow-600">{result.summary.total_skipped}</div>
          <div className="text-sm text-yellow-700 font-medium mt-1">跳过重复</div>
        </Card>
        <Card className="p-4 text-center bg-red-50 border-red-200">
          <div className="text-3xl font-bold text-red-600">{result.summary.total_errors}</div>
          <div className="text-sm text-red-700 font-medium mt-1">导入错误</div>
        </Card>
      </div>

      {/* Detailed Results */}
      <div className="space-y-3">
        {/* Providers */}
        {result.providers.total > 0 && (
          <Card className="p-4">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold text-gray-900">
                提供商 ({result.providers.imported}/{result.providers.total})
              </h3>
              <div className="text-xs text-gray-500">
                跳过: {result.providers.skipped} | 错误: {result.providers.errors.length}
              </div>
            </div>
            {result.providers.errors.length > 0 && (
              <div className="space-y-1 mt-2">
                {result.providers.errors.slice(0, 3).map((error, index) => (
                  <div key={index} className="text-xs text-red-600 bg-red-50 p-2 rounded">
                    第 {error.row} 行 - {error.field}: {error.error}
                  </div>
                ))}
                {result.providers.errors.length > 3 && (
                  <div className="text-xs text-gray-500 text-center">
                    还有 {result.providers.errors.length - 3} 个错误...
                  </div>
                )}
              </div>
            )}
          </Card>
        )}

        {/* Models */}
        {result.models.total > 0 && (
          <Card className="p-4">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold text-gray-900">
                模型 ({result.models.imported}/{result.models.total})
              </h3>
              <div className="text-xs text-gray-500">
                跳过: {result.models.skipped} | 错误: {result.models.errors.length}
              </div>
            </div>
            {result.models.errors.length > 0 && (
              <div className="space-y-1 mt-2">
                {result.models.errors.slice(0, 3).map((error, index) => (
                  <div key={index} className="text-xs text-red-600 bg-red-50 p-2 rounded">
                    第 {error.row} 行 - {error.field}: {error.error}
                  </div>
                ))}
                {result.models.errors.length > 3 && (
                  <div className="text-xs text-gray-500 text-center">
                    还有 {result.models.errors.length - 3} 个错误...
                  </div>
                )}
              </div>
            )}
          </Card>
        )}

        {/* Model Provider Mappings */}
        {result.modelProviderMappings.total > 0 && (
          <Card className="p-4">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold text-gray-900">
                模型映射 ({result.modelProviderMappings.imported}/{result.modelProviderMappings.total})
              </h3>
              <div className="text-xs text-gray-500">
                跳过: {result.modelProviderMappings.skipped} | 错误: {result.modelProviderMappings.errors.length}
              </div>
            </div>
            {result.modelProviderMappings.errors.length > 0 && (
              <div className="space-y-1 mt-2">
                {result.modelProviderMappings.errors.slice(0, 3).map((error, index) => (
                  <div key={index} className="text-xs text-red-600 bg-red-50 p-2 rounded">
                    第 {error.row} 行 - {error.field}: {error.error}
                  </div>
                ))}
                {result.modelProviderMappings.errors.length > 3 && (
                  <div className="text-xs text-gray-500 text-center">
                    还有 {result.modelProviderMappings.errors.length - 3} 个错误...
                  </div>
                )}
              </div>
            )}
          </Card>
        )}
      </div>
    </div>
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">系统设置迁移</h1>
        <p className="text-sm text-gray-500 mt-1">通过Excel文件导出或导入系统配置</p>
      </div>

      {/* Export Section */}
      <Card className="p-6">
        <div className="flex items-center mb-4">
          <Download className="w-5 h-5 mr-2 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">导出配置</h2>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Button
            onClick={() => downloadTemplate(false)}
            variant="secondary"
            className="h-24 flex flex-col items-center justify-center"
          >
            <FileSpreadsheet className="w-8 h-8 mb-2 text-gray-600" />
            <div className="text-sm font-medium">空白模板</div>
            <div className="text-xs text-gray-500">下载空白Excel模板</div>
          </Button>
          
          <Button
            onClick={() => downloadTemplate(true)}
            variant="secondary"
            className="h-24 flex flex-col items-center justify-center"
          >
            <FileSpreadsheet className="w-8 h-8 mb-2 text-gray-600" />
            <div className="text-sm font-medium">示例模板</div>
            <div className="text-xs text-gray-500">包含示例数据</div>
          </Button>
          
          <Button
            onClick={exportConfig}
            disabled={isExporting}
            className="h-24 flex flex-col items-center justify-center"
          >
            {isExporting ? (
              <>
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white mb-2"></div>
                <div className="text-sm font-medium">导出中...</div>
              </>
            ) : (
              <>
                <Download className="w-8 h-8 mb-2" />
                <div className="text-sm font-medium">导出当前配置</div>
                <div className="text-xs opacity-80">保存为Excel</div>
              </>
            )}
          </Button>
        </div>
      </Card>

      {/* Import Section */}
      <Card className="p-6">
        <div className="flex items-center mb-4">
          <Upload className="w-5 h-5 mr-2 text-blue-600" />
          <h2 className="text-lg font-semibold text-gray-900">导入配置</h2>
        </div>
        
        <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-blue-400 transition-colors">
          <FileSpreadsheet className="w-16 h-16 mx-auto text-gray-400 mb-4" />
          <p className="text-lg font-medium text-gray-700 mb-2">上传Excel配置文件</p>
          <p className="text-sm text-gray-500 mb-4">仅支持 .xlsx 格式的文件</p>
          
          <input
            type="file"
            id="file-upload"
            className="hidden"
            accept=".xlsx"
            onChange={handleFileUpload}
            disabled={isUploading}
          />
          
          <label
            htmlFor="file-upload"
            className={`inline-flex items-center px-6 py-3 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 cursor-pointer transition-colors ${
              isUploading ? 'opacity-50 cursor-not-allowed' : ''
            }`}
          >
            {isUploading ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                正在导入...
              </>
            ) : (
              <>
                <Upload className="w-4 h-4 mr-2" />
                选择文件
              </>
            )}
          </label>
        </div>

        {/* Import Instructions */}
        <div className="mt-6 p-4 bg-blue-50 rounded-lg border border-blue-200">
          <div className="flex items-start">
            <Info className="w-5 h-5 text-blue-600 mr-3 mt-0.5 flex-shrink-0" />
            <div>
              <h3 className="font-medium text-blue-900 mb-2">导入说明</h3>
              <ul className="text-sm text-blue-800 space-y-1">
                <li>• Excel文件必须包含三个工作表：Providers、Models、ModelProviderMappings</li>
                <li>• 系统会自动跳过已存在的重复数据</li>
                <li>• 建议先下载模板，按照模板格式填写数据</li>
                <li>• 导入前建议先导出当前配置进行备份</li>
              </ul>
            </div>
          </div>
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
        
        <div className="flex justify-between items-center pt-4 mt-4 border-t">
          <div>
            {importResult?.result.summary.total_errors === 0 ? (
              <Badge variant="success" className="flex items-center">
                <CheckCircle className="w-4 h-4 mr-1" />
                导入成功
              </Badge>
            ) : (
              <Badge variant="error" className="flex items-center">
                <AlertCircle className="w-4 h-4 mr-1" />
                部分导入失败
              </Badge>
            )}
          </div>
          <Button
            onClick={() => setShowResultModal(false)}
          >
            关闭
          </Button>
        </div>
      </Modal>
    </div>
  )
}

export default ImportExport