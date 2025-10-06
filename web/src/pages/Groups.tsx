import React, { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Search } from 'lucide-react'
import { Card, Button, Input, Modal, Badge } from '../components/ui'
import { groupsApi } from '../services'
import type { Group } from '../types'

export const Groups: React.FC = () => {
  const [groups, setGroups] = useState<Group[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingGroup, setEditingGroup] = useState<Group | null>(null)

  useEffect(() => {
    loadGroups()
  }, [])

  const loadGroups = async () => {
    try {
      setLoading(true)
      const response = await groupsApi.getGroups()
      setGroups(response.data)
     } catch (error) {
      console.error('加载分组失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = () => {
    setEditingGroup(null)
    setIsModalOpen(true)
  }

  const handleEdit = (group: Group) => {
    setEditingGroup(group)
    setIsModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除此分组吗?')) return

    try {
      await groupsApi.deleteGroup(id)
      await loadGroups()
    } catch (error) {
      console.error('删除分组失败:', error)
      alert('删除失败,请重试')
    }
  }

  const handleSubmit = async (formData: Partial<Group>) => {
    try {
      if (editingGroup) {
        await groupsApi.updateGroup(editingGroup.id, formData)
      } else {
        // 确保必填字段存在
        if (!formData.name) {
          alert('分组名称不能为空')
          return
        }
        await groupsApi.createGroup(formData as any)
      }
      setIsModalOpen(false)
      await loadGroups()
    } catch (error) {
      console.error('保存分组失败:', error)
      alert('保存失败,请重试')
    }
  }

  const filteredGroups = groups.filter(group =>
    group.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    group.description?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <div className="space-y-6">
      {/* 页面标题和操作栏 */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">分组管理</h1>
          <p className="text-sm text-gray-500 mt-1">管理 API 密钥分组</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          新建分组
        </Button>
      </div>

      {/* 搜索栏 */}
      <Card className="p-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
          <Input
            placeholder="搜索分组名称或描述..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
      </Card>

      {/* 分组列表 */}
      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredGroups.map((group) => (
            <Card key={group.id} className="p-6 hover:shadow-lg transition-shadow">
              <div className="space-y-4">
                {/* 分组头部 */}
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-gray-900">{group.name}</h3>
                    {group.description && (
                      <p className="text-sm text-gray-500 mt-1">{group.description}</p>
                    )}
                  </div>
                  <Badge variant={group.enabled ? 'success' : 'default'}>
                    {group.enabled ? '已启用' : '已禁用'}
                  </Badge>
                </div>

                {/* 统计信息 */}
                <div className="grid grid-cols-2 gap-4 pt-4 border-t border-gray-100">
                  <div>
                    <p className="text-xs text-gray-500">优先级</p>
                    <p className="text-lg font-semibold text-gray-900">{group.priority}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">负载策略</p>
                    <p className="text-sm font-medium text-gray-900">{group.loadBalanceStrategy}</p>
                  </div>
                </div>

                {/* 操作按钮 */}
                <div className="flex gap-2 pt-2">
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => handleEdit(group)}
                    className="flex-1"
                  >
                    <Edit2 className="w-4 h-4 mr-1" />
                    编辑
                  </Button>
                  <Button
                    variant="danger"
                    size="sm"
                    onClick={() => handleDelete(group.id)}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </Card>
          ))}

          {filteredGroups.length === 0 && (
            <div className="col-span-full text-center py-12">
              <p className="text-gray-500">
                {searchQuery ? '未找到匹配的分组' : '暂无分组,点击右上角新建分组'}
              </p>
            </div>
          )}
        </div>
      )}

      {/* 编辑/创建模态框 */}
      <GroupModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        group={editingGroup}
      />
    </div>
  )
}

interface GroupModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: Partial<Group>) => void
  group: Group | null
}

const GroupModal: React.FC<GroupModalProps> = ({ isOpen, onClose, onSubmit, group }) => {
  const [activeTab, setActiveTab] = useState('general')
  const [formData, setFormData] = useState<Partial<Group>>({
    name: '',
    description: '',
    enabled: true,
    modelAliases: {}
  })

  useEffect(() => {
    if (group) {
      setFormData({
        name: group.name,
        description: group.description,
        enabled: group.enabled,
        modelAliases: group.modelAliases || {}
      })
    } else {
      setFormData({
        name: '',
        description: '',
        enabled: true,
        modelAliases: {}
      })
    }
    setActiveTab('general')
  }, [group])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.name?.trim()) {
      alert('请输入分组名称')
      return
    }
    onSubmit(formData)
  }

  const renderGeneralSettings = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          分组名称 <span className="text-red-500">*</span>
        </label>
        <Input
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="例如: 主要分组"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          描述
        </label>
        <textarea
          value={formData.description || ''}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          placeholder="分组描述信息"
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
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
          启用此分组
        </label>
      </div>
    </div>
  )

  const renderModelAliases = () => {
    const aliases = formData.modelAliases || {}
    const aliasEntries = Object.entries(aliases) as [string, string][]

    const handleAliasChange = (index: number, key: string, value: string) => {
      const newAliasEntries = [...aliasEntries];
      newAliasEntries[index] = [key, value];
      const newAliases = Object.fromEntries(newAliasEntries);
      setFormData({ ...formData, modelAliases: newAliases });
    }

    const addAlias = () => {
      const newAliases = { ...aliases, [`new_model_${aliasEntries.length}`]: '' }
      setFormData({ ...formData, modelAliases: newAliases })
    }

    const removeAlias = (key: string) => {
      const newAliases = { ...aliases }
      delete newAliases[key]
      setFormData({ ...formData, modelAliases: newAliases })
    }

    return (
      <div className="space-y-4">
        {aliasEntries.map(([key, value], index) => (
          <div key={index} className="flex items-center gap-2">
            <Input
              placeholder="原始模型名称"
              value={key}
              onChange={(e) => handleAliasChange(index, e.target.value, value)}
              className="flex-1"
            />
            <Input
              placeholder="映射后名称"
              value={value}
              onChange={(e) => handleAliasChange(index, key, e.target.value)}
              className="flex-1"
            />
            <Button variant="danger" onClick={() => removeAlias(key)}>
              <Trash2 className="w-4 h-4" />
            </Button>
          </div>
        ))}
        <Button onClick={addAlias} variant="secondary" className="w-full">
          <Plus className="w-4 h-4 mr-2" />
          添加映射
        </Button>
      </div>
    )
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={group ? '编辑分组' : '新建分组'}>
      <div className="border-b border-gray-200 mb-4">
        <nav className="-mb-px flex space-x-8" aria-label="Tabs">
          <button
            onClick={() => setActiveTab('general')}
            className={`whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'general'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            通用设置
          </button>
          {group && (
            <button
              onClick={() => setActiveTab('aliases')}
              className={`whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'aliases'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              模型映射
            </button>
          )}
        </nav>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        {activeTab === 'general' && renderGeneralSettings()}
        {activeTab === 'aliases' && renderModelAliases()}

        <div className="flex gap-3 pt-4">
          <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
            取消
          </Button>
          <Button type="submit" className="flex-1">
            {group ? '保存' : '创建'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}

export default Groups