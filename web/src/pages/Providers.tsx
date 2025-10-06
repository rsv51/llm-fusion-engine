import React, { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Search, Zap, Code, Copy, List, Download } from 'lucide-react' // Added icons for new features
import { Card, Button, Input, Modal, Badge } from '../components/ui'
import { api } from '../services'
import type { Provider, PaginationResponse, CreateProviderRequest } from '../types'
import { modelsApi } from '../services/models'

export const Providers: React.FC = () => {
  const [providers, setProviders] = useState<Provider[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProvider, setEditingProvider] = useState<Provider | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [providerModels, setProviderModels] = useState<{providerId: number, models: string[], providerName: string} | null>(null)
  const [isModelsModalOpen, setIsModelsModalOpen] = useState(false)
  const [selectedProvider, setSelectedProvider] = useState<Provider | null>(null)
  const [importing, setImporting] = useState(false)

  useEffect(() => {
    loadProviders()
  }, [])

  const loadProviders = async () => {
    setError(null); // 清除之前的错误
    try {
      setLoading(true)
      // api.get() 已经返回了 response.data，所以 response 就是后端返回的完整响应
      const response = await api.get<any>('/admin/providers')
      console.log('Providers API Response:', response) // 添加日志
      // 尝试安全地访问数据，后端返回格式为 { data: [...], pagination: {...} }
      if (response && Array.isArray(response.data)) {
        setProviders(response.data);
      } else {
        const errorMsg = 'API 响应数据格式不正确: ' + JSON.stringify(response);
        console.error(errorMsg);
        setError(errorMsg);
      }
    } catch (error: any) {
      console.error('加载供应商失败:', error)
      setError(error.message || JSON.stringify(error));
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingProvider(null)
    setIsModalOpen(true)
  }

  const handleEdit = (provider: Provider) => {
    setEditingProvider(provider)
    setIsModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除此供应商吗?')) return

    try {
      await api.delete(`/admin/providers/${id}`)
      await loadProviders()
    } catch (error) {
      console.error('删除供应商失败:', error)
      alert('删除失败,请重试')
    }
  }

  const handleSubmit = async (formData: Partial<Provider>) => {
    try {
      if (editingProvider) {
        await api.put(`/admin/providers/${editingProvider.id}`, formData)
      } else {
        await api.post('/admin/providers', formData)
      }
      setIsModalOpen(false)
      await loadProviders()
    } catch (error) {
      console.error('保存供应商失败:', error)
      alert('保存失败,请重试')
    }
  }

  const filteredProviders = React.useMemo(() => {
    if (!searchQuery.trim()) return providers;
    
    const query = searchQuery.toLowerCase();
    return providers.filter(provider =>
      provider?.name?.toLowerCase().includes(query) ||
      provider?.type?.toLowerCase().includes(query)
    );
  }, [providers, searchQuery]);

  const handleCheckHealth = async (id: number) => {
    try {
      await api.post(`/admin/health/providers/${id}`)
      await loadProviders()
    } catch (error) {
      console.error('健康检查失败:', error)
    }
  }

  const handleGetModels = async (provider: Provider) => {
    try {
      setSelectedProvider(provider)
      // 传入 provider.type，确保后端按实际供应商类型返回模型列表
      const response = await modelsApi.getProviderModels(provider.id, provider.type)
      setProviderModels({
        providerId: provider.id,
        models: response.models,
        providerName: response.providerName
      })
      setIsModelsModalOpen(true)
    } catch (error) {
      console.error('获取模型列表失败:', error)
      alert('获取模型列表失败，请重试')
    }
  }

  const handleImportAllModels = async () => {
    if (!providerModels || !selectedProvider) return
    
    try {
      setImporting(true)
      await modelsApi.importModels({
        providerId: providerModels.providerId,
        modelNames: providerModels.models
      })
      alert('所有模型导入成功')
      setIsModelsModalOpen(false)
    } catch (error) {
      console.error('导入模型失败:', error)
      alert('导入模型失败，请重试')
    } finally {
      setImporting(false)
    }
  }

  const handleImportSingleModel = async (modelName: string) => {
    if (!selectedProvider) return
    
    try {
      setImporting(true)
      await modelsApi.importModels({
        providerId: selectedProvider.id,
        modelNames: [modelName]
      })
      alert(`模型 ${modelName} 导入成功`)
    } catch (error) {
      console.error('导入模型失败:', error)
      alert('导入模型失败，请重试')
    } finally {
      setImporting(false)
    }
  }


  return (
    <div className="space-y-6">
      {/* 如果有错误，显示在页面上 */}
      {error && (
        <div style={{ color: 'red', border: '1px solid red', padding: '10px', backgroundColor: '#ffeeee' }}>
          <h2>前端错误:</h2>
          <pre>{error}</pre>
        </div>
      )}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">供应商管理</h1>
          <p className="text-sm text-gray-500 mt-1">管理所有可用的供应商</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          新建供应商
        </Button>
      </div>

      <Card className="p-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
          <Input
            placeholder="搜索供应商名称或类型..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
      </Card>

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <>
          {filteredProviders && filteredProviders.length > 0 ? (
            <Card>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">供应商名称</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">类型</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">状态</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">健康状态</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">权重</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">延迟</th>
                      <th className="px-4 py-3 text-right text-sm font-medium text-gray-700">操作</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {filteredProviders.map((provider) => (
                      <tr key={provider.id} className="hover:bg-gray-50">
                        <td className="px-4 py-3">
                          <div>
                            <div className="font-medium text-gray-900">{provider.name}</div>
                            <div className="text-sm text-gray-500">ID: {provider.id}</div>
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="text-sm text-gray-900">{provider.type}</div>
                        </td>
                        <td className="px-4 py-3">
                          <Badge variant={provider.enabled ? 'success' : 'default'}>
                            {provider.enabled ? '已启用' : '已禁用'}
                          </Badge>
                        </td>
                        <td className="px-4 py-3">
                          {provider.healthStatus && (
                            <Badge variant={provider.healthStatus === 'healthy' ? 'success' : provider.healthStatus === 'unhealthy' ? 'error' : 'warning'}>
                              {provider.healthStatus === 'healthy' ? '健康' : provider.healthStatus === 'unhealthy' ? '不健康' : '未知'}
                            </Badge>
                          )}
                        </td>
                        <td className="px-4 py-3">
                          <div className="text-sm text-gray-900">{provider.weight}</div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="text-sm text-gray-900">
                            {provider.latency ? `${provider.latency}ms` : '-'}
                          </div>
                        </td>
                        <td className="px-4 py-3 text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="secondary"
                              size="sm"
                              onClick={() => handleCheckHealth(provider.id)}
                              title="健康检查"
                            >
                              <Zap className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="secondary"
                              size="sm"
                              onClick={() => handleGetModels(provider)}
                              title="获取模型列表"
                            >
                              <List className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="secondary"
                              size="sm"
                              onClick={() => handleEdit(provider)}
                              title="编辑"
                            >
                              <Edit2 className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="danger"
                              size="sm"
                              onClick={() => handleDelete(provider.id)}
                              title="删除"
                            >
                              <Trash2 className="w-4 h-4" />
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </Card>
          ) : (
            <div className="text-center py-12">
              <p className="text-gray-500">暂无供应商，请点击右上角"新建供应商"按钮添加。</p>
            </div>
          )}
        </>
      )}

      <ProviderModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        provider={editingProvider}
      />

      <ProviderModelsModal
        isOpen={isModelsModalOpen}
        onClose={() => setIsModelsModalOpen(false)}
        providerModels={providerModels}
        onImportAll={handleImportAllModels}
        onImportSingle={handleImportSingleModel}
        importing={importing}
      />
    </div>
  )
}

interface ProviderModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: Partial<Provider>) => void
  provider: Provider | null
}

const ProviderModal: React.FC<ProviderModalProps> = ({ isOpen, onClose, onSubmit, provider }) => {
  // 解析配置对象
  const parseConfig = (configStr: string) => {
    try {
      return JSON.parse(configStr);
    } catch {
      return { baseUrl: '', timeout: 30, maxRetries: 3, enabled: true, apiKey: '' };
    }
  };

  const getDefaultFormData = (): Partial<CreateProviderRequest> => ({
    name: '',
    type: 'openai',
    config: JSON.stringify({ baseUrl: '', timeout: 30, maxRetries: 3, enabled: true, apiKey: '' }, null, 2),
    enabled: true,
    weight: 1,
  });

  const getInitialFormData = (): Partial<CreateProviderRequest> => {
    if (provider) {
      return {
        name: provider.name,
        type: provider.type,
        config: provider.config,
        console: provider.console,
        enabled: provider.enabled,
        weight: provider.weight,
      };
    }
    return getDefaultFormData();
  };

  const [formData, setFormData] = useState<Partial<CreateProviderRequest>>(getInitialFormData());
  const [configError, setConfigError] = useState<string | null>(null);
  const [config, setConfig] = useState(parseConfig(formData.config || ''));

  useEffect(() => {
    setFormData(getInitialFormData());
    setConfigError(null);
    setConfig(parseConfig(getInitialFormData().config || ''));
  }, [provider]);

  // 当配置对象改变时，更新 JSON 字符串
  useEffect(() => {
    try {
      const configStr = JSON.stringify(config, null, 2);
      setFormData({ ...formData, config: configStr });
      setConfigError(null);
    } catch (err) {
      setConfigError("无效的配置格式");
    }
  }, [config]);

  const handleConfigChange = (field: string, value: string | number | boolean) => {
    setConfig({ ...config, [field]: value });
  };

  const handleDirectConfigChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const configStr = e.target.value;
    setFormData({ ...formData, config: configStr });
    try {
      const parsedConfig = JSON.parse(configStr);
      setConfig(parsedConfig);
      setConfigError(null);
    } catch (err) {
      setConfigError("无效的 JSON 格式");
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (configError) {
      alert("请修正配置中的 JSON 错误")
      return
    }
    onSubmit(formData)
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={provider ? '编辑供应商' : '新建供应商'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">供应商名称</label>
          <Input
            value={formData.name || ''}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="例如: MyOpenAIInstance"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">供应商类型</label>
          <select
            value={formData.type || 'openai'}
            onChange={(e) => setFormData({ ...formData, type: e.target.value })}
            className="w-full p-2 border rounded"
            required
          >
            <option value="openai">OpenAI</option>
            <option value="anthropic">Anthropic</option>
            <option value="gemini">Google Gemini</option>
            {/* Add more provider types as needed */}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">API 密钥</label>
          <Input
            type="password"
            value={config.apiKey || ''}
            onChange={(e) => handleConfigChange('apiKey', e.target.value)}
            placeholder="输入 API 密钥"
            className="w-full"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">基础 URL</label>
          <Input
            value={config.baseUrl || ''}
            onChange={(e) => handleConfigChange('baseUrl', e.target.value)}
            placeholder="例如: https://api.openai.com/v1"
            className="w-full"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">超时时间 (秒)</label>
          <Input
            type="number"
            value={config.timeout || 30}
            onChange={(e) => handleConfigChange('timeout', parseInt(e.target.value) || 30)}
            className="w-full"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">最大重试次数</label>
          <Input
            type="number"
            value={config.maxRetries || 3}
            onChange={(e) => handleConfigChange('maxRetries', parseInt(e.target.value) || 3)}
            className="w-full"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            高级配置 (JSON 格式)
            <Code className="w-4 h-4 inline ml-2" />
          </label>
          <textarea
            value={formData.config || ''}
            onChange={handleDirectConfigChange}
            placeholder={`例如:\n{\n  "baseUrl": "https://api.openai.com/v1",\n  "apiKey": "sk-...",\n  "timeout": 30\n}`}
            rows={6}
            className="w-full p-2 border rounded font-mono text-sm"
          />
          {configError && <p className="text-red-500 text-sm">{configError}</p>}
          <p className="text-xs text-gray-500 mt-1">
            提示：您可以直接编辑上面的 JSON 配置，或者使用上方的表单字段。修改表单字段会自动更新此处的 JSON 配置。
          </p>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">控制台地址 (可选)</label>
          <Input
            value={formData.console || ''}
            onChange={(e) => setFormData({ ...formData, console: e.target.value })}
            placeholder="https://console.example.com"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">权重</label>
          <Input
            type="number"
            value={formData.weight || 0}
            onChange={(e) => setFormData({ ...formData, weight: parseInt(e.target.value) || 0 })}
            required
          />
        </div>
        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={formData.enabled || false}
            onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
            id="enabled-checkbox"
          />
          <label htmlFor="enabled-checkbox" className="text-sm text-gray-700">启用此供应商</label>
        </div>
        <div className="flex gap-3 pt-4">
          <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
            取消
          </Button>
          <Button type="submit" className="flex-1">
            {provider ? '保存' : '创建'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}

interface ProviderModelsModalProps {
  isOpen: boolean
  onClose: () => void
  providerModels: {providerId: number, models: string[], providerName: string} | null
  onImportAll: () => void
  onImportSingle: (modelName: string) => void
  importing: boolean
}

const ProviderModelsModal: React.FC<ProviderModelsModalProps> = ({
  isOpen,
  onClose,
  providerModels,
  onImportAll,
  onImportSingle,
  importing
}) => {
  if (!providerModels) return null

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`${providerModels.providerName} - 模型列表`}>
      <div className="space-y-4">
        <div>
          <p className="text-sm text-gray-500 mb-2">
            找到 {providerModels.models.length} 个可用模型
          </p>
          <div className="max-h-60 overflow-y-auto border rounded-lg">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left text-sm font-medium text-gray-700">模型名称</th>
                  <th className="px-4 py-2 text-right text-sm font-medium text-gray-700">操作</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {providerModels.models.map((model, index) => (
                  <tr key={index}>
                    <td className="px-4 py-2 text-sm text-gray-900">{model}</td>
                    <td className="px-4 py-2 text-right">
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => onImportSingle(model)}
                        disabled={importing}
                      >
                        <Download className="w-4 h-4 mr-1" />
                        导入
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
        
        <div className="flex justify-between pt-4 border-t border-gray-100">
          <Button
            variant="secondary"
            onClick={onClose}
          >
            取消
          </Button>
          <Button
            onClick={onImportAll}
            disabled={importing}
          >
            <Download className="w-4 h-4 mr-2" />
            {importing ? '导入中...' : '导入所有模型到配置'}
          </Button>
        </div>
      </div>
    </Modal>
  )
}

export default Providers