import React, { useState, useEffect } from 'react'
import { Card, Button, Input, Modal, Badge } from '../components/ui'
import { Box, Settings, CheckCircle, XCircle, Plus, Edit2, Trash2, Search, Copy } from 'lucide-react'
import { api } from '../services'
import type { Model, PaginationResponse } from '../types'

export const Models: React.FC = () => {
  const [models, setModels] = useState<Model[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingModel, setEditingModel] = useState<Model | null>(null)

  useEffect(() => {
    loadModels(1)
  }, [])

  const pickData = (resp: any) => {
    if (!resp) return { data: [], pagination: {} };
    if (Array.isArray(resp.data) && resp.pagination) return resp;
    if (Array.isArray(resp.items)) return { data: resp.items, pagination: { page: resp.page, totalPage: resp.totalPages } };
    if (Array.isArray(resp)) return { data: resp, pagination: {} };
    return { data: [], pagination: {} };
  };

  const loadModels = async (pageNum: number) => {
    try {
      setLoading(true);
      const response = await api.get('/admin/models', {
        params: { page: pageNum, pageSize: 12, search: searchQuery },
      });
      const { data, pagination } = pickData(response);
      setModels(data);
      setPage(pagination?.page || 1);
      setTotalPages(pagination?.totalPage || 1);
    } catch (error) {
      console.error('加载模型失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    loadModels(1)
  }

  const handleCreate = () => {
    setEditingModel(null)
    setIsModalOpen(true)
  }

  const handleEdit = (model: Model) => {
    setEditingModel(model)
    setIsModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除此模型吗?')) return
    try {
      await api.delete(`/admin/models/${id}`)
      await loadModels(page)
    } catch (error) {
      console.error('删除模型失败:', error)
    }
  }

  const handleClone = async (model: Model) => {
    try {
      const newName = `${model.name}-copy`;
      await api.post(`/admin/models/${model.id}/clone`, { newName })
      await loadModels(page)
    } catch (error) {
      console.error('克隆模型失败:', error)
    }
  }

  const handleSubmit = async (formData: Partial<Model>) => {
    try {
      if (editingModel) {
        await api.put(`/admin/models/${editingModel.id}`, formData)
      } else {
        await api.post('/admin/models', formData)
      }
      setIsModalOpen(false)
      await loadModels(page)
    } catch (error) {
      console.error('保存模型失败:', error)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">模型配置</h1>
        <Button onClick={handleCreate}>
          <Plus className="w-4 h-4 mr-2" />
          新建模型
        </Button>
      </div>

      <Card className="p-4">
        <div className="flex gap-4">
          <Input
            placeholder="搜索模型名称..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="flex-1"
          />
          <Button onClick={handleSearch}>
            <Search className="w-4 h-4 mr-2" />
            搜索
          </Button>
        </div>
      </Card>

      {loading ? (
        <div className="text-center py-12">加载中...</div>
      ) : (
        <>
          <Card>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">模型名称</th>
                    <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">备注</th>
                    <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">最大重试次数</th>
                    <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">超时时间</th>
                    <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">状态</th>
                    <th className="px-4 py-3 text-right text-sm font-medium text-gray-700">操作</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {models.map((model) => (
                    <tr key={model.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                            <Box className="w-4 h-4 text-blue-600" />
                          </div>
                          <div className="font-medium text-gray-900">{model.name}</div>
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <div className="text-sm text-gray-900">{model.remark || '无备注'}</div>
                      </td>
                      <td className="px-4 py-3">
                        <div className="text-sm text-gray-900">{model.maxRetry}</div>
                      </td>
                      <td className="px-4 py-3">
                        <div className="text-sm text-gray-900">{model.timeout}秒</div>
                      </td>
                      <td className="px-4 py-3">
                        <Badge variant={model.enabled ? 'success' : 'default'}>
                          {model.enabled ? '启用' : '禁用'}
                        </Badge>
                      </td>
                      <td className="px-4 py-3 text-right">
                        <div className="flex justify-end gap-2">
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() => handleClone(model)}
                            title="克隆"
                          >
                            <Copy className="w-4 h-4" />
                          </Button>
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() => handleEdit(model)}
                            title="编辑"
                          >
                            <Edit2 className="w-4 h-4" />
                          </Button>
                          <Button
                            variant="danger"
                            size="sm"
                            onClick={() => handleDelete(model.id)}
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
          <div className="flex justify-center items-center gap-4">
            <Button onClick={() => loadModels(page - 1)} disabled={page <= 1}>
              上一页
            </Button>
            <span>第 {page} / {totalPages} 页</span>
            <Button onClick={() => loadModels(page + 1)} disabled={page >= totalPages}>
              下一页
            </Button>
          </div>
        </>
      )}

      <ModelModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        model={editingModel}
      />
    </div>
  )
}

interface ModelModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: Partial<Model>) => void
  model: Model | null
}

const ModelModal: React.FC<ModelModalProps> = ({ isOpen, onClose, onSubmit, model }) => {
  const [formData, setFormData] = useState<Partial<Model>>({})

  useEffect(() => {
    if (model) {
      setFormData({
        ...model,
        enabled: model.enabled ?? true,
      })
    } else {
      setFormData({ name: '', remark: '', maxRetry: 3, timeout: 30, enabled: true })
    }
  }, [model])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.name?.trim()) {
      alert('模型名称不能为空')
      return
    }
    onSubmit(formData)
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={model ? '编辑模型' : '新建模型'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">模型名称</label>
          <Input
            value={formData.name || ''}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">备注</label>
          <Input
            value={formData.remark || ''}
            onChange={(e) => setFormData({ ...formData, remark: e.target.value })}
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">最大重试次数</label>
          <Input
            type="number"
            value={formData.maxRetry || 0}
            onChange={(e) => setFormData({ ...formData, maxRetry: parseInt(e.target.value) || 0 })}
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">超时时间（秒）</label>
          <Input
            type="number"
            value={formData.timeout || 0}
            onChange={(e) => setFormData({ ...formData, timeout: parseInt(e.target.value) || 0 })}
            min="1"
          />
        </div>
        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={formData.enabled ?? false}
            onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
          />
          <label>启用</label>
        </div>
        <div className="flex justify-end gap-2">
          <Button type="button" variant="secondary" onClick={onClose}>取消</Button>
          <Button type="submit">保存</Button>
        </div>
      </form>
    </Modal>
  )
}

export default Models