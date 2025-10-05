import React, { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Search, Zap } from 'lucide-react'
import { Card, Button, Input, Modal, Badge } from '../components/ui'
import { api } from '../services'
import type { Provider, PaginationResponse } from '../types'

export const Providers: React.FC = () => {
  const [providers, setProviders] = useState<Provider[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProvider, setEditingProvider] = useState<Provider | null>(null)

  useEffect(() => {
    loadProviders()
  }, [])

  const loadProviders = async () => {
    setError(null); // 清除之前的错误
    try {
      setLoading(true)
      const response = await api.get<PaginationResponse<Provider>>('/admin/providers')
      console.log('Providers API Response:', response) // 添加日志
      // 尝试安全地访问数据
      if (response && response.data && Array.isArray(response.data.data)) {
        setProviders(response.data.data);
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

  let filteredProviders: Provider[] = []
  try {
    filteredProviders = providers.filter(provider =>
      provider?.providerType?.toLowerCase().includes(searchQuery.toLowerCase())
    )
  } catch (error) {
    console.error('Error filtering providers:', error, providers)
    filteredProviders = []
  }

  const handleCheckHealth = async (id: number) => {
    try {
      await api.post(`/admin/health/providers/${id}`)
      await loadProviders()
    } catch (error) {
      console.error('健康检查失败:', error)
    }
  }

  // 添加一个错误状态
  const [error, setError] = useState<string | null>(null);

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
            placeholder="搜索供应商类型..."
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
                        <h3 className="text-lg font-semibold text-gray-900">{provider.providerType}</h3>
                        <p className="text-sm text-gray-500 mt-1">ID: {provider.id}</p>
                      </div>
                      <Badge variant={provider.enabled ? 'success' : 'default'}>
                        {provider.enabled ? '已启用' : '已禁用'}
                      </Badge>
                    </div>

                    <div className="grid grid-cols-2 gap-4 pt-4 border-t border-gray-100">
                      <div>
                        <p className="text-xs text-gray-500">权重</p>
                        <p className="text-lg font-semibold text-gray-900">{provider.weight}</p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-500">分组ID</p>
                        <p className="text-lg font-semibold text-gray-900">{provider.groupId}</p>
                      </div>
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
              <p className="text-gray-500">暂无供应商，请点击右上角“新建供应商”按钮添加。</p>
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
  const [formData, setFormData] = useState<Partial<Provider>>({})

  useEffect(() => {
    if (provider) {
      setFormData(provider)
    } else {
      setFormData({
        providerType: '',
        enabled: true,
        weight: 1,
        groupId: 0,
        timeout: 30,
        maxRetries: 3,
      })
    }
  }, [provider])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={provider ? '编辑供应商' : '新建供应商'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">供应商类型</label>
          <Input
            value={formData.providerType}
            onChange={(e) => setFormData({ ...formData, providerType: e.target.value })}
            placeholder="例如: openai"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">分组ID</label>
          <Input
            type="number"
            value={formData.groupId}
            onChange={(e) => setFormData({ ...formData, groupId: parseInt(e.target.value) })}
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">权重</label>
          <Input
            type="number"
            value={formData.weight}
            onChange={(e) => setFormData({ ...formData, weight: parseInt(e.target.value) })}
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Base URL</label>
          <Input
            value={formData.baseUrl}
            onChange={(e) => setFormData({ ...formData, baseUrl: e.target.value })}
            placeholder="https://api.openai.com/v1"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">超时 (秒)</label>
          <Input
            type="number"
            value={formData.timeout}
            onChange={(e) => setFormData({ ...formData, timeout: parseInt(e.target.value) })}
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">最大重试次数</label>
          <Input
            type="number"
            value={formData.maxRetries}
            onChange={(e) => setFormData({ ...formData, maxRetries: parseInt(e.target.value) })}
          />
        </div>
        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={formData.enabled}
            onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
          />
          <label className="text-sm text-gray-700">启用此供应商</label>
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