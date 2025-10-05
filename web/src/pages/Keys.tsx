import React, { useState, useEffect } from 'react'
import { Plus, Eye, EyeOff, Copy, Trash2, Search } from 'lucide-react'
import { Card, Button, Input, Modal, Badge } from '../components/ui'
import { keysApi, groupsApi } from '../services'
import type { ApiKey, Group, PaginationParams } from '../types'

export const Keys: React.FC = () => {
  const [keys, setKeys] = useState<ApiKey[]>([])
  const [groups, setGroups] = useState<Group[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedGroup, setSelectedGroup] = useState<number | null>(null)
  const [pagination, setPagination] = useState<PaginationParams>({ page: 1, pageSize: 20 })
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [visibleKeys, setVisibleKeys] = useState<Set<number>>(new Set())

  useEffect(() => {
    loadGroups()
  }, [])

  useEffect(() => {
    loadKeys()
  }, [pagination, selectedGroup])

  const loadGroups = async () => {
    try {
      const response = await groupsApi.getGroups()
      setGroups(response.data)
     } catch (error) {
      console.error('加载分组失败:', error)
    }
  }

  const loadKeys = async () => {
    try {
      setLoading(true)
      const response = await keysApi.getKeys({
        ...pagination,
        groupId: selectedGroup || undefined
       })
       setKeys(response.data)
       setTotal(response.pagination.total)
      } catch (error) {
       console.error('加载密钥失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setIsModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除此密钥吗?')) return

    try {
      await keysApi.deleteKey(id)
      await loadKeys()
    } catch (error) {
      console.error('删除密钥失败:', error)
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

  const filteredKeys = keys.filter(key =>
    (key.description?.toLowerCase() || '').includes(searchQuery.toLowerCase()) ||
    key.key.toLowerCase().includes(searchQuery.toLowerCase()) ||
    key.groupName.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <div className="space-y-6">
      {/* 页面标题和操作栏 */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">密钥管理</h1>
          <p className="text-sm text-gray-500 mt-1">管理 API 密钥,共 {total} 个</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          新建密钥
        </Button>
      </div>

      {/* 搜索和筛选栏 */}
      <Card className="p-4">
        <div className="flex flex-col sm:flex-row gap-4">
          {/* 搜索框 */}
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <Input
              placeholder="搜索密钥名称或密钥值..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>

          {/* 分组筛选 */}
          <div className="sm:w-64">
            <select
              value={selectedGroup || ''}
              onChange={(e) => setSelectedGroup(e.target.value ? Number(e.target.value) : null)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">所有分组</option>
              {groups.map(group => (
                <option key={group.id} value={group.id}>{group.name}</option>
              ))}
            </select>
          </div>
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
                    名称
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    密钥
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    分组
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    状态
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    使用次数
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredKeys.map((key) => {
                  const isVisible = visibleKeys.has(key.id)
                  const groupName = groups.find(g => g.id === key.groupId)?.name || '未分组'

                  return (
                    <tr key={key.id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900">{key.description || `密钥 ${key.id}`}</div>
                      </td>
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
                        <span className="text-sm text-gray-900">{groupName}</span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Badge variant={key.enabled ? 'success' : 'default'}>
                          {key.enabled ? '已启用' : '已禁用'}
                        </Badge>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="text-sm text-gray-900">{key.usageCount || 0}</span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
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

      {/* 创建模态框 */}
      <KeyModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSuccess={loadKeys}
        groups={groups}
      />
    </div>
  )
}

interface KeyModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
  groups: Group[]
}

const KeyModal: React.FC<KeyModalProps> = ({ isOpen, onClose, onSuccess, groups }) => {
  const [formData, setFormData] = useState({
    description: '',
    key: '',
    groupId: 1,
    enabled: true
  })

  const generateRandomKey = () => {
    const prefix = 'fk-';
    const randomPart = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
    setFormData({ ...formData, key: prefix + randomPart });
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.key.trim() || !formData.groupId) {
      alert('请填写完整信息')
      return
    }

    try {
      await keysApi.createKey({
        key: formData.key,
        groupId: formData.groupId,
        description: formData.description || undefined,
        enabled: formData.enabled
      })
      onClose()
      onSuccess()
      setFormData({ description: '', key: '', groupId: 1, enabled: true })
    } catch (error) {
      console.error('创建密钥失败:', error)
      alert('创建失败,请重试')
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="新建密钥">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            描述信息
          </label>
          <Input
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="例如: OpenAI 主密钥"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            密钥值 <span className="text-red-500">*</span>
          </label>
          <div className="flex gap-2">
            <Input
              value={formData.key}
              onChange={(e) => setFormData({ ...formData, key: e.target.value })}
              placeholder="sk-..."
              required
              className="flex-1"
            />
            <Button type="button" variant="secondary" onClick={generateRandomKey}>
              随机生成
            </Button>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            所属分组 <span className="text-red-500">*</span>
          </label>
          <select
            value={formData.groupId}
            onChange={(e) => setFormData({ ...formData, groupId: Number(e.target.value) })}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          >
            {groups.map(group => (
              <option key={group.id} value={group.id}>{group.name}</option>
            ))}
          </select>
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
            创建
          </Button>
        </div>
      </form>
    </Modal>
  )
}

export default Keys