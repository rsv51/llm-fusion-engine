import React, { useState, useEffect } from 'react';
import { Plus, Edit2, Trash2, Search } from 'lucide-react';
import { Card, Button, Input, Modal, Badge } from '../components/ui';
import { api } from '../services';
import type { ModelMapping, Provider, PaginationResponse } from '../types';

export const ModelMappings: React.FC = () => {
  const [mappings, setMappings] = useState<ModelMapping[]>([]);
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingMapping, setEditingMapping] = useState<ModelMapping | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setError(null); // 清除之前的错误
    try {
      setLoading(true);
      // api.get() 已经返回了 response.data
      const [mappingsPaginationResponse, providersPaginationResponse] = await Promise.all([
        api.get<PaginationResponse<ModelMapping>>('/admin/model-mappings'),
        api.get<PaginationResponse<Provider>>('/admin/providers'),
      ]);
      console.log('ModelMappings API Response:', mappingsPaginationResponse); // 添加日志
      console.log('Providers API Response in ModelMappings:', providersPaginationResponse); // 添加日志

      // 安全地检查和设置数据
      if (mappingsPaginationResponse && Array.isArray(mappingsPaginationResponse.data)) {
        setMappings(mappingsPaginationResponse.data);
      } else {
        const errorMsg = 'ModelMappings API 响应数据格式不正确: ' + JSON.stringify(mappingsPaginationResponse);
        console.error(errorMsg);
        setError(errorMsg);
      }

      if (providersPaginationResponse && Array.isArray(providersPaginationResponse.data)) {
        setProviders(providersPaginationResponse.data);
      } else {
        const errorMsg = 'Providers API (in ModelMappings) 响应数据格式不正确: ' + JSON.stringify(providersPaginationResponse);
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

  const handleEdit = (mapping: ModelMapping) => {
    setEditingMapping(mapping);
    setIsModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除此模型映射吗?')) return;
    try {
      await api.delete(`/admin/model-mappings/${id}`);
      await loadData();
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  const handleSubmit = async (formData: Partial<ModelMapping>) => {
    try {
      if (editingMapping) {
        await api.put(`/admin/model-mappings/${editingMapping.id}`, formData);
      } else {
        await api.post('/admin/model-mappings', formData);
      }
      setIsModalOpen(false);
      await loadData();
    } catch (error) {
      console.error('保存失败:', error);
    }
  };

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
        <h1 className="text-2xl font-bold">模型映射</h1>
        <Button onClick={handleCreate}><Plus className="w-4 h-4 mr-2" />新建映射</Button>
      </div>

      <Card>
        <table className="w-full">
          <thead>
            <tr>
              <th className="px-4 py-2 text-left">用户友好名称</th>
              <th className="px-4 py-2 text-left">提供商模型名称</th>
              <th className="px-4 py-2 text-left">提供商</th>
              <th className="px-4 py-2 text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            {mappings.map((mapping) => (
              <tr key={mapping.id}>
                <td className="px-4 py-2">{mapping.userFriendlyName}</td>
                <td className="px-4 py-2">{mapping.providerModelName}</td>
                <td className="px-4 py-2">{mapping.provider?.providerType}</td>
                <td className="px-4 py-2 text-right">
                  <Button variant="secondary" size="sm" onClick={() => handleEdit(mapping)}><Edit2 className="w-4 h-4" /></Button>
                  <Button variant="danger" size="sm" onClick={() => handleDelete(mapping.id)} className="ml-2"><Trash2 className="w-4 h-4" /></Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </Card>

      <ModelMappingModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
        mapping={editingMapping}
        providers={providers}
      />
    </div>
  );
};

interface ModelMappingModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: Partial<ModelMapping>) => void;
  mapping: ModelMapping | null;
  providers: Provider[];
}

const ModelMappingModal: React.FC<ModelMappingModalProps> = ({ isOpen, onClose, onSubmit, mapping, providers }) => {
  const [formData, setFormData] = useState<Partial<ModelMapping>>({});

  useEffect(() => {
    if (mapping) {
      setFormData(mapping);
    } else {
      setFormData({
        userFriendlyName: '',
        providerModelName: '',
        providerId: providers[0]?.id || 0,
      });
    }
  }, [mapping, providers]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={mapping ? '编辑模型映射' : '新建模型映射'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label>用户友好名称</label>
          <Input value={formData.userFriendlyName} onChange={(e) => setFormData({ ...formData, userFriendlyName: e.target.value })} required />
        </div>
        <div>
          <label>提供商模型名称</label>
          <Input value={formData.providerModelName} onChange={(e) => setFormData({ ...formData, providerModelName: e.target.value })} required />
        </div>
        <div>
          <label>提供商</label>
          <select value={formData.providerId} onChange={(e) => setFormData({ ...formData, providerId: parseInt(e.target.value) })} className="w-full p-2 border rounded">
            {providers.map(p => <option key={p.id} value={p.id}>{p.providerType}</option>)}
          </select>
        </div>
        <div className="flex justify-end gap-2">
          <Button type="button" variant="secondary" onClick={onClose}>取消</Button>
          <Button type="submit">保存</Button>
        </div>
      </form>
    </Modal>
  );
};

export default ModelMappings;