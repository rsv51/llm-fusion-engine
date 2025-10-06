import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Badge } from '../components/ui';
import { logsApi } from '../services';
import type { Log } from '../types';

export const LogDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [log, setLog] = useState<Log | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (id) {
      const fetchLog = async () => {
        try {
          setLoading(true);
          const response = await logsApi.getLog(id);
          setLog(response);
        } catch (err: any) {
          setError(err.message || 'Failed to fetch log details');
        } finally {
          setLoading(false);
        }
      };
      fetchLog();
    }
  }, [id]);

  const renderJson = (jsonString: string) => {
    try {
      const obj = JSON.parse(jsonString);
      return <pre className="whitespace-pre-wrap break-all">{JSON.stringify(obj, null, 2)}</pre>;
    } catch {
      return <pre className="whitespace-pre-wrap break-all">{jsonString}</pre>;
    }
  };

  if (loading) {
    return <div className="text-center py-12">加载中...</div>;
  }

  if (error) {
    return <div className="text-center py-12 text-red-500">错误: {error}</div>;
  }

  if (!log) {
    return <div className="text-center py-12">未找到日志记录。</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">请求详情</h1>
      <Card>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-4 p-6">
          <DetailItem label="请求时间" value={new Date(log.timestamp).toLocaleString('zh-CN')} />
          <DetailItem label="代理密钥" value={log.proxy_key} />
          <DetailItem label="提供商" value={log.provider} />
          <DetailItem label="模型" value={log.model} />
          <DetailItem label="状态码" value={<Badge variant={log.is_success ? 'success' : 'error'}>{log.response_status}</Badge>} />
          <DetailItem label="请求类型" value={(log.request_body || '').includes('"stream":true') ? '流式' : '非流式'} />
          <DetailItem label="响应时间" value={`${log.latency}ms`} />
          <DetailItem label="Token使用量" value={log.total_tokens} />
          <DetailItem label="客户端IP" value={log.request_url} /> 
        </div>
      </Card>

      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold mb-4">请求内容</h2>
          <div className="bg-gray-100 p-4 rounded-lg text-sm font-mono">
            {renderJson(log.request_body || '')}
          </div>
        </div>
      </Card>

      <Card>
        <div className="p-6">
          <h2 className="text-lg font-semibold mb-4">响应内容</h2>
          <div className="bg-gray-100 p-4 rounded-lg text-sm font-mono">
            {renderJson(log.response_body || '')}
          </div>
        </div>
      </Card>
    </div>
  );
};

const DetailItem: React.FC<{ label: string; value: React.ReactNode }> = ({ label, value }) => (
  <div className="border-b border-gray-200 py-2">
    <div className="text-sm font-medium text-gray-500">{label}</div>
    <div className="text-base text-gray-900 mt-1">{value}</div>
  </div>
);

export default LogDetail;