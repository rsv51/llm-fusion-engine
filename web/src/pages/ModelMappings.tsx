import React, { useState, useEffect } from 'react';
import { Plus, Edit2, Trash2, Search, Filter } from 'lucide-react';
import { Card, Button, Input, Modal, Badge, Select } from '../components/ui';
import { api, modelsApi } from '../services';
import type { ModelProviderMapping, Provider, Model } from '../types';

export const ModelMappings: React.FC = () => {
  const [mappings, setMappings] = useState<ModelProviderMapping[]>([]);
  const [models, setModels] = useState<Model[]>([]);
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingMapping, setEditingMapping] = useState<ModelProviderMapping | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedProviderId, setSelectedProviderId] = useState<number | ''>('');
  const [providerModels, setProviderModels] = useState<string[]>([]);
  const [loadingProviderModels, setLoadingProviderModels] = useState(false);
  const [healthStatusData, setHealthStatusData] = useState<Record<number, any[]>>({});

  useEffect(() => {
    loadData();
    loadHealthStatus();
  }, []);

  // 统一提取列表数据,兼容多种后端返回格式:
  // - 直接数组: [...]
  // - 顶层 items: { items: [...] }
  // - data: 数组: { data: [...] }
  // - data.items: { data: { items: [...] } }
  const pickList = (resp: any) => {
    if (Array.isArray(resp)) return resp;
    if (resp && Array.isArray(resp.data)) return resp.data;
    if (resp && Array.isArray(resp.items)) return resp.items;
    if (resp && resp.data && Array.isArray(resp.data.items)) return resp.data.items;
    return null;
  };

  const loadData = async () => {
    setError(null); // 清除之前的错误
    try {
      setLoading(true);
      // api.get() 已经返回了 response.data
      const [mappingsResponse, modelsResponse, providersResponse] = await Promise.all([
        api.get<any>('/admin/model-provider-mappings'),
        api.get<any>('/admin/models'),
        api.get<any>('/admin/providers'),
      ]) as [any, any, any];
      console.log('ModelProviderMappings API Response:', mappingsResponse); // 添加日志
      console.log('Models API Response:', modelsResponse); // 添加日志
      console.log('Providers API Response in ModelMappings:', providersResponse); // 添加日志

      // 安全地检查和设置数据
      const mappingsList = pickList(mappingsResponse);
      if (mappingsList) {
        setMappings(mappingsList);
      } else {
        const errorMsg = 'ModelProviderMappings API 响应数据格式不正确: ' + JSON.stringify(mappingsResponse);
        console.error(errorMsg);
        setError(errorMsg);
      }

      const modelsList = pickList(modelsResponse);
      if (modelsList) {
        setModels(modelsList);
      } else {
        const errorMsg = 'Models API 响应数据格式不正确: ' + JSON.stringify(modelsResponse);
        console.error(errorMsg);
        setError(errorMsg);
      }

      const providersList = pickList(providersResponse);
      if (providersList) {
        setProviders(providersList);
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

  const loadHealthStatus = async () => {
    try {
      const response = await api.get<Record<number, any[]>>('/admin/model-provider-mappings/health/all');
      setHealthStatusData(response.data || response);
    } catch (error) {
      console.error('加载健康状态失败:', error);
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
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">健康状态</th>
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
                          <HealthStatusIndicator healthStatus={healthStatusData[mapping.id] || []} />
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
          <ModelSelector
            models={models}
            providers={providers}
            selectedProviderId={formData.providerId || ''}
            onProviderChange={(providerId) => setFormData({ ...formData, providerId: providerId || 0 })}
            onModelSelect={(model) => setFormData({ ...formData, modelId: model.id })}
            selectedModel={models.find(m => m.id === formData.modelId)}
          />
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

// 模型选择器组件 - 用于按供应商筛选模型
const ModelSelector: React.FC<{
  models: Model[];
  providers: Provider[];
  selectedProviderId: number | '';
  onProviderChange: (providerId: number | '') => void;
  onModelSelect: (model: Model) => void;
  selectedModel?: Model | null;
}> = ({ models, providers, selectedProviderId, onProviderChange, onModelSelect, selectedModel }) => {
  const [filteredModels, setFilteredModels] = useState<Model[]>([]);
  const [modelSearch, setModelSearch] = useState('');
  
  // 当选择的供应商或搜索词变化时,筛选对应的模型
  useEffect(() => {
    let currentModels = models;
    
    if (modelSearch.trim()) {
      currentModels = currentModels.filter(m =>
        m.name.toLowerCase().includes(modelSearch.toLowerCase()) ||
        m.remark?.toLowerCase().includes(modelSearch.toLowerCase())
      );
    }
    
    setFilteredModels(currentModels);
  }, [selectedProviderId, models, providers, modelSearch]);

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">按供应商筛选</label>
        <Select
          value={selectedProviderId}
          onChange={(e) => onProviderChange(e.target.value === '' ? '' : parseInt(e.target.value))}
          options={[
            { value: '', label: '所有供应商' },
            ...providers.map(p => ({ value: p.id, label: `${p.name} (${p.type})` }))
          ]}
        />
      </div>
      
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">选择模型</label>
        <Input
          placeholder="搜索模型名称或备注..."
          value={modelSearch}
          onChange={(e) => setModelSearch(e.target.value)}
          className="mb-2"
        />
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2 max-h-60 overflow-y-auto p-1 border rounded">
          {filteredModels.length > 0 ? (
            filteredModels.map(model => (
              <div
                key={model.id}
                className={`p-2 border rounded cursor-pointer hover:bg-gray-50 transition-colors ${
                  selectedModel?.id === model.id ? 'bg-blue-50 border-blue-500' : 'border-gray-200'
                }`}
                onClick={() => onModelSelect(model)}
              >
                <div className="font-medium text-sm">{model.name}</div>
                <div className="text-xs text-gray-500">{model.remark || '无备注'}</div>
              </div>
            ))
          ) : (
            <div className="col-span-full text-center py-4 text-gray-500">
              没有找到相关模型
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// 健康状态指示器组件
interface HealthStatusIndicatorProps {
  healthStatus: Array<{
    timestamp: string;
    status: 'success' | 'error';
    statusCode: number;
    latencyMs?: number;
  }>;
}

const HealthStatusIndicator: React.FC<HealthStatusIndicatorProps> = ({ healthStatus }) => {
  if (!healthStatus || healthStatus.length === 0) {
    return (
      <div className="flex items-center gap-1">
        <div className="text-xs text-gray-400">暂无数据</div>
      </div>
    );
  }

  return (
    <div className="flex items-center gap-1" title="近期10次调用状态(最新→最旧)">
      {healthStatus.slice(0, 10).map((status, index) => (
        <div
          key={index}
          className={`w-3 h-3 rounded-full ${
            status.status === 'success'
              ? 'bg-green-500'
              : 'bg-red-500'
          }`}
          title={`${status.status === 'success' ? '成功' : '失败'} (${status.statusCode})${status.latencyMs ? ` - ${status.latencyMs}ms` : ''}\n${new Date(status.timestamp).toLocaleString()}`}
        />
      ))}
    </div>
  );
};

export default ModelMappings;