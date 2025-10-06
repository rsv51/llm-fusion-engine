import React, { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Search, Zap, Code } from 'lucide-react' // Added Code icon for config
import { Card, Button, Input, Modal, Badge } from '../components/ui'
import { api } from '../services'
import type { Provider, PaginationResponse, CreateProviderRequest } from '../types'

export const Providers: React.FC = () => {
  const [providers, setProviders] = useState<Provider[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProvider, setEditingProvider] = useState<Provider | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadProviders()
  }, [])

  const loadProviders = async () => {
    setError(null); // 清除之前的错误
    try {
      setLoading(true)
      // api.get() 已经返回了 response.data，所以 response 就是 PaginationResponse<Provider>
      const paginationResponse = await api.get<PaginationResponse<Provider>>('/admin/providers')
      console.log('Providers API Response:', paginationResponse) // 添加日志
      // 尝试安全地访问数据
      if (paginationResponse && Array.isArray(paginationResponse.data)) {
        setProviders(paginationResponse.data);
      } else {
        const errorMsg = 'API 响应数据格式不正确: ' + JSON.stringify(paginationResponse);
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
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredProviders.map((provider) => (
                <Card key={provider.id} className="p-6 hover:shadow-lg transition-shadow">
                  <div className="space-y-4">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h3 className="text-lg font-semibold text-gray-900">{provider.name}</h3>
                        <p className="text-sm text-gray-500 mt-1">ID: {provider.id}</p>
                      </div>
                      <Badge variant={provider.enabled ? 'success' : 'default'}>
                        {provider.enabled ? '已启用' : '已禁用'}
                      </Badge>
                    </div>

                    <div className="space-y-2">
                      <div>
                        <p className="text-xs text-gray-500">类型</p>
                        <p className="text-sm font-medium text-gray-900">{provider.type}</p>
                      </div>
                      {provider.healthStatus && (
                        <div>
                          <p className="text-xs text-gray-500">健康状态</p>
                          <Badge variant={provider.healthStatus === 'healthy' ? 'success' : provider.healthStatus === 'unhealthy' ? 'error' : 'warning'}>
                            {provider.healthStatus === 'healthy' ? '健康' : provider.healthStatus === 'unhealthy' ? '不健康' : '未知'}
                          </Badge>
                        </div>
                      )}
                    </div>

                    <div className="grid grid-cols-2 gap-4 pt-4 border-t border-gray-100">
                      <div>
                        <p className="text-xs text-gray-500">权重</p>
                        <p className="text-lg font-semibold text-gray-900">{provider.weight}</p>
                      </div>
                      {provider.latency && (
                        <div>
                          <p className="text-xs text-gray-500">延迟</p>
                          <p className="text-lg font-semibold text-gray-900">{provider.latency}ms</p>
                        </div>
                      )}
                    </div>

                    <div className="flex gap-2 pt-2">
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => handleCheckHealth(provider.id)}
                        className="flex-1"
                      >
                        <Zap className="w-4 h-4 mr-1" />
                        检查
                      </Button>
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => handleEdit(provider)}
                      >
                        <Edit2 className="w-4 h-4" />
                      </Button>
                      <Button
                        variant="danger"
                        size="sm"
                        onClick={() => handleDelete(provider.id)}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
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

export default Providers