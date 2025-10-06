import React, { useState, useEffect } from 'react'
import { Plus, Eye, EyeOff, Copy, Trash2, Search, Edit } from 'lucide-react'
import { Card, Button, Input, Modal, Badge } from '../components/ui'
import { proxyKeysApi } from '../services'
import type { ProxyKey, PaginationParams } from '../types'

export const ProxyKeys: React.FC = () => {
  const [proxyKeys, setProxyKeys] = useState<ProxyKey[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [pagination, setPagination] = useState<PaginationParams>({ page: 1, pageSize: 20 })
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingKey, setEditingKey] = useState<ProxyKey | null>(null)
  const [visibleKeys, setVisibleKeys] = useState<Set<number>>(new Set())

  useEffect(() => {
    loadProxyKeys()
  }, [pagination])

  const loadProxyKeys = async () => {
    try {
      setLoading(true)
      const response = await proxyKeysApi.getProxyKeys(pagination)
      setProxyKeys(response.data)
      setTotal(response.pagination.total)
    } catch (error) {
      console.error('加载代理密钥失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingKey(null)
    setIsModalOpen(true)
  }

  const handleEdit = (key: ProxyKey) => {
    setEditingKey(key)
    setIsModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除此代理密钥吗?')) return

    try {
      await proxyKeysApi.deleteProxyKey(id)
      await loadProxyKeys()
    } catch (error) {
      console.error('删除代理密钥失败:', error)
      alert('删除失败,请重试')
    }
  }

  const toggleKeyVisibility = (keyId: number) => {
    setVisibleKeys(prev => {
      const newSet = new Set(prev)
      if (newSet.has(keyId)) {
        newSet.delete(keyId)
      } else {
        newSet.add(keyId)
      }
      return newSet
    })
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    alert('已复制到剪贴板')
  }

  const maskKey = (key: string) => {
    if (key.length <= 8) return key
    return key.substring(0, 4) + '•'.repeat(key.length - 8) + key.substring(key.length - 4)
  }

  const filteredKeys = proxyKeys.filter(key =>
    key.key.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <div className="space-y-6">
      {/* 页面标题和操作栏 */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">代理密钥管理</h1>
          <p className="text-sm text-gray-500 mt-1">管理用户访问密钥,共 {total} 个</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          新建密钥
        </Button>
      </div>

      {/* 搜索栏 */}
      <Card className="p-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
          <Input
            placeholder="搜索密钥值..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
      </Card>

      {/* 密钥列表 */}
      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <Card>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    密钥
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    RPM限制
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    TPM限制
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    状态
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredKeys.map((key) => {
                  const isVisible = visibleKeys.has(key.id)

                  return (
                    <tr key={key.id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <code className="text-sm text-gray-600 font-mono">
                            {isVisible ? key.key : maskKey(key.key)}
                          </code>
                          <button
                            onClick={() => toggleKeyVisibility(key.id)}
                            className="text-gray-400 hover:text-gray-600"
                          >
                            {isVisible ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                          </button>
                          <button
                            onClick={() => copyToClipboard(key.key)}
                            className="text-gray-400 hover:text-gray-600"
                          >
                            <Copy className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="text-sm text-gray-900">
                          {key.rpmLimit || '无限制'}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="text-sm text-gray-900">
                          {key.tpmLimit || '无限制'}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Badge variant={key.enabled ? 'success' : 'default'}>
                          {key.enabled ? '已启用' : '已禁用'}
                        </Badge>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm space-x-2">
                        <Button
                          variant="secondary"
                          size="sm"
                          onClick={() => handleEdit(key)}
                        >
                          <Edit className="w-4 h-4" />
                        </Button>
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={() => handleDelete(key.id)}
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>

            {filteredKeys.length === 0 && (
              <div className="text-center py-12">
                <p className="text-gray-500">
                  {searchQuery ? '未找到匹配的密钥' : '暂无密钥,点击右上角新建密钥'}
                </p>
              </div>
            )}
          </div>

          {/* 分页 */}
          {total > pagination.pageSize && (
            <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
              <div className="text-sm text-gray-500">
                共 {total} 条记录,第 {pagination.page} / {Math.ceil(total / pagination.pageSize)} 页
              </div>
              <div className="flex gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  disabled={pagination.page === 1}
                  onClick={() => setPagination({ ...pagination, page: pagination.page - 1 })}
                >
                  上一页
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  disabled={pagination.page >= Math.ceil(total / pagination.pageSize)}
                  onClick={() => setPagination({ ...pagination, page: pagination.page + 1 })}
                >
                  下一页
                </Button>
              </div>
            </div>
          )}
        </Card>
      )}

      {/* 创建/编辑模态框 */}
      <ProxyKeyModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSuccess={loadProxyKeys}
        editingKey={editingKey}
      />
    </div>
  )
}

interface ProxyKeyModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
  editingKey: ProxyKey | null
}

const ProxyKeyModal: React.FC<ProxyKeyModalProps> = ({ isOpen, onClose, onSuccess, editingKey }) => {
  const [formData, setFormData] = useState({
    key: '',
    enabled: true,
    rpmLimit: 0,
    tpmLimit: 0,
  })

  useEffect(() => {
    if (editingKey) {
      setFormData({
        key: editingKey.key,
        enabled: editingKey.enabled,
        rpmLimit: editingKey.rpmLimit || 0,
        tpmLimit: editingKey.tpmLimit || 0,
      })
    } else {
      setFormData({ key: '', enabled: true, rpmLimit: 0, tpmLimit: 0 })
    }
  }, [editingKey, isOpen])

  const generateRandomKey = () => {
    const prefix = 'pk-'
    const randomPart = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
    setFormData({ ...formData, key: prefix + randomPart })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.key.trim()) {
      alert('请填写密钥值')
      return
    }

    try {
      if (editingKey) {
        await proxyKeysApi.updateProxyKey(editingKey.id, {
          enabled: formData.enabled,
          rpmLimit: formData.rpmLimit || 0,
          tpmLimit: formData.tpmLimit || 0,
        })
      } else {
        await proxyKeysApi.createProxyKey({
          key: formData.key,
          enabled: formData.enabled,
          rpmLimit: formData.rpmLimit || 0,
          tpmLimit: formData.tpmLimit || 0,
        })
      }
      onClose()
      onSuccess()
    } catch (error) {
      console.error('保存密钥失败:', error)
      alert('保存失败,请重试')
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={editingKey ? '编辑代理密钥' : '新建代理密钥'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            密钥值 <span className="text-red-500">*</span>
          </label>
          <div className="flex gap-2">
            <Input
              value={formData.key}
              onChange={(e) => setFormData({ ...formData, key: e.target.value })}
              placeholder="pk-..."
              required
              disabled={!!editingKey}
              className="flex-1"
            />
            {!editingKey && (
              <Button type="button" variant="secondary" onClick={generateRandomKey}>
                随机生成
              </Button>
            )}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            RPM限制 (每分钟请求数)
          </label>
          <Input
            type="number"
            value={formData.rpmLimit}
            onChange={(e) => setFormData({ ...formData, rpmLimit: parseInt(e.target.value) || 0 })}
            placeholder="0 = 无限制"
            min="0"
          />
          <p className="text-xs text-gray-500 mt-1">设置为0表示无限制</p>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            TPM限制 (每分钟Token数)
          </label>
          <Input
            type="number"
            value={formData.tpmLimit}
            onChange={(e) => setFormData({ ...formData, tpmLimit: parseInt(e.target.value) || 0 })}
            placeholder="0 = 无限制"
            min="0"
          />
          <p className="text-xs text-gray-500 mt-1">设置为0表示无限制</p>
        </div>

        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            id="enabled"
            checked={formData.enabled}
            onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
            className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
          />
          <label htmlFor="enabled" className="text-sm text-gray-700">
            启用此密钥
          </label>
        </div>

        <div className="flex gap-3 pt-4">
          <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
            取消
          </Button>
          <Button type="submit" className="flex-1">
            {editingKey ? '保存' : '创建'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}

export default ProxyKeys