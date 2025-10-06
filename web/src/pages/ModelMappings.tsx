import React, { useState, useEffect } from 'react';
import { Plus, Edit2, Trash2, Search } from 'lucide-react';
import { Card, Button, Input, Modal, Badge } from '../components/ui';
import { api } from '../services';
import type { ModelProviderMapping, Provider, Model, PaginationResponse } from '../types';

export const ModelMappings: React.FC = () => {
  const [mappings, setMappings] = useState<ModelProviderMapping[]>([]);
  const [models, setModels] = useState<Model[]>([]);
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingMapping, setEditingMapping] = useState<ModelProviderMapping | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setError(null); // 清除之前的错误
    try {
      setLoading(true);
      // api.get() 已经返回了 response.data
      const [mappingsResponse, modelsResponse, providersResponse] = await Promise.all([
        api.get<any>('/admin/model-provider-mappings'),
        api.get<any>('/admin/models'),
        api.get<any>('/admin/providers'),
      ]);
      console.log('ModelProviderMappings API Response:', mappingsResponse); // 添加日志
      console.log('Models API Response:', modelsResponse); // 添加日志
      console.log('Providers API Response in ModelMappings:', providersResponse); // 添加日志

      // 安全地检查和设置数据
      if (mappingsResponse && Array.isArray(mappingsResponse.data)) {
        setMappings(mappingsResponse.data);
      } else {
        const errorMsg = 'ModelProviderMappings API 响应数据格式不正确: ' + JSON.stringify(mappingsResponse);
        console.error(errorMsg);
        setError(errorMsg);
      }

      if (modelsResponse && Array.isArray(modelsResponse.data)) {
        setModels(modelsResponse.data);
      } else if (modelsResponse && modelsResponse.data && Array.isArray(modelsResponse.data.items)) {
        setModels(modelsResponse.data.items);
      } else {
        const errorMsg = 'Models API 响应数据格式不正确: ' + JSON.stringify(modelsResponse);
        console.error(errorMsg);
        setError(errorMsg);
      }

      if (providersResponse && Array.isArray(providersResponse.data)) {
        setProviders(providersResponse.data);
      } else if (providersResponse && providersResponse.data && Array.isArray(providersResponse.data.items)) {
        setProviders(providersResponse.data.items);
      } else {
        const errorMsg = 'Providers API (in ModelMappings) 响应数据格式不正确: ' + JSON.stringify(providersResponse);
        console.error(errorMsg);
        setError(errorMsg);
      }
    } catch (error: any) {
      console.error('加载数据失败:', error);
      setError(error.message || JSON.stringify(error));
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingMapping(null);
    setIsModalOpen(true);
  };

  const handleEdit = (mapping: ModelProviderMapping) => {
    setEditingMapping(mapping);
    setIsModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除此模型映射吗?')) return;
    try {
      await api.delete(`/admin/model-provider-mappings/${id}`);
      await loadData();
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  const handleSubmit = async (formData: Partial<ModelProviderMapping>) => {
    try {
      if (editingMapping) {
        await api.put(`/admin/model-provider-mappings/${editingMapping.id}`, formData);
      } else {
        await api.post('/admin/model-provider-mappings', formData);
      }
      setIsModalOpen(false);
      await loadData();
    } catch (error) {
      console.error('保存失败:', error);
    }
  };

  // 过滤映射
  const filteredMappings = React.useMemo(() => {
    if (!searchQuery.trim()) return mappings;
    
    const query = searchQuery.toLowerCase();
    return mappings.filter(mapping => (
      mapping.model?.name?.toLowerCase().includes(query) ||
      mapping.provider?.name?.toLowerCase().includes(query) ||
      mapping.providerModel?.toLowerCase().includes(query)
    ));
  }, [mappings, searchQuery]);

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
          <h1 className="text-2xl font-bold">模型映射</h1>
          <p className="text-sm text-gray-500 mt-1">管理模型与提供商之间的映射关系</p>
        </div>
        <Button onClick={handleCreate}><Plus className="w-4 h-4 mr-2" />新建映射</Button>
      </div>

      <Card className="p-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
          <Input
            placeholder="搜索模型、提供商或模型名称..."
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
          {filteredMappings && filteredMappings.length > 0 ? (
            <Card>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">模型名称</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">提供商模型名称</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">提供商</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">功能</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">权重</th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">状态</th>
                      <th className="px-4 py-3 text-right text-sm font-medium text-gray-700">操作</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {filteredMappings.map((mapping) => (
                      <tr key={mapping.id} className="hover:bg-gray-50">
                        <td className="px-4 py-3">
                          <div className="font-medium text-gray-900">{mapping.model?.name}</div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="text-sm text-gray-900">{mapping.providerModel}</div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="text-sm text-gray-900">{mapping.provider?.name} ({mapping.provider?.type})</div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex gap-1 flex-wrap">
                            {mapping.toolCall && <Badge variant="info">工具调用</Badge>}
                            {mapping.structuredOutput && <Badge variant="info">结构化输出</Badge>}
                            {mapping.image && <Badge variant="info">图像输入</Badge>}
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="text-sm text-gray-900">{mapping.weight}</div>
                        </td>
                        <td className="px-4 py-3">
                          <Badge variant={mapping.enabled ? 'success' : 'default'}>
                            {mapping.enabled ? '启用' : '禁用'}
                          </Badge>
                        </td>
                        <td className="px-4 py-3 text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="secondary"
                              size="sm"
                              onClick={() => handleEdit(mapping)}
                              title="编辑"
                            >
                              <Edit2 className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="danger"
                              size="sm"
                              onClick={() => handleDelete(mapping.id)}
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
              <p className="text-gray-500">暂无模型映射，请点击右上角"新建映射"按钮添加。</p>
            </div>
          )}
        </>
      )}

      <ModelMappingModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        mapping={editingMapping}
        models={models}
        providers={providers}
      />
    </div>
  );
};

interface ModelMappingModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: Partial<ModelProviderMapping>) => void;
  mapping: ModelProviderMapping | null;
  models: Model[];
  providers: Provider[];
}

const ModelMappingModal: React.FC<ModelMappingModalProps> = ({ isOpen, onClose, onSubmit, mapping, models, providers }) => {
  const getInitialFormData = (): Partial<ModelProviderMapping> => {
    if (mapping) {
      return mapping;
    }
    
    return {
      modelId: models[0]?.id || 0,
      providerId: providers[0]?.id || 0,
      providerModel: '',
      toolCall: false,
      structuredOutput: false,
      image: false,
      weight: 1,
      enabled: true,
    };
  };

  const [formData, setFormData] = useState<Partial<ModelProviderMapping>>(getInitialFormData());

  useEffect(() => {
    setFormData(getInitialFormData());
  }, [mapping, models, providers]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={mapping ? '编辑模型映射' : '新建模型映射'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">模型</label>
          <select
            value={formData.modelId || ''}
            onChange={(e) => setFormData({ ...formData, modelId: parseInt(e.target.value) })}
            className="w-full p-2 border rounded"
            required
          >
            {models.map(m => <option key={m.id} value={m.id}>{m.name}</option>)}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">提供商</label>
          <select
            value={formData.providerId || ''}
            onChange={(e) => setFormData({ ...formData, providerId: parseInt(e.target.value) })}
            className="w-full p-2 border rounded"
            required
          >
            {providers.map(p => <option key={p.id} value={p.id}>{p.name} ({p.type})</option>)}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">提供商模型名称</label>
          <Input
            value={formData.providerModel || ''}
            onChange={(e) => setFormData({ ...formData, providerModel: e.target.value })}
            placeholder="例如: gpt-4-turbo"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">权重</label>
          <Input
            type="number"
            value={formData.weight || 0}
            onChange={(e) => setFormData({ ...formData, weight: parseInt(e.target.value) || 0 })}
            min="1"
            required
          />
        </div>
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700">功能支持</label>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={formData.toolCall || false}
                onChange={(e) => setFormData({ ...formData, toolCall: e.target.checked })}
                id="toolCall-checkbox"
              />
              <label htmlFor="toolCall-checkbox" className="text-sm text-gray-700">工具调用</label>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={formData.structuredOutput || false}
                onChange={(e) => setFormData({ ...formData, structuredOutput: e.target.checked })}
                id="structuredOutput-checkbox"
              />
              <label htmlFor="structuredOutput-checkbox" className="text-sm text-gray-700">结构化输出</label>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={formData.image || false}
                onChange={(e) => setFormData({ ...formData, image: e.target.checked })}
                id="image-checkbox"
              />
              <label htmlFor="image-checkbox" className="text-sm text-gray-700">图像输入</label>
            </div>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            checked={formData.enabled || false}
            onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
            id="enabled-checkbox"
          />
          <label htmlFor="enabled-checkbox" className="text-sm text-gray-700">启用此映射</label>
        </div>
        <div className="flex justify-end gap-2 pt-4">
          <Button type="button" variant="secondary" onClick={onClose}>取消</Button>
          <Button type="submit">保存</Button>
        </div>
      </form>
    </Modal>
  );
};

export default ModelMappings;