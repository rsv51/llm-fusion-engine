// 模型相关类型定义

export interface Model {
  id: number;
  name: string; // e.g., "GPT-4-Turbo"
  remark?: string;
  maxRetry: number;
  timeout: number; // in seconds
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateModelRequest {
  name: string;
  remark?: string;
  maxRetry?: number;
  timeout?: number;
  enabled?: boolean;
}

export interface UpdateModelRequest extends Partial<CreateModelRequest> {}