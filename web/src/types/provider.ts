// 供应商相关类型定义

export interface ProviderConfig {
  baseUrl?: string;
  timeout?: number;
  maxRetries?: number;
  enabled?: boolean;
  apiKey?: string; // Example, actual config depends on provider type
  // Add other provider-specific fields as needed
}

export interface Provider {
  id: number;
  name: string; // e.g., "MyOpenAIInstance"
  type: string; // e.g., "openai", "anthropic"
  config: string; // JSON string of ProviderConfig
  console?: string;
  enabled: boolean;
  weight: number;
  healthStatus?: 'healthy' | 'unhealthy' | 'unknown';
  lastChecked?: string;
  latency?: number; // in milliseconds
  createdAt: string;
  updatedAt: string;
  // Note: ApiKeys are now part of the 'config' JSON string.
}

export interface CreateProviderRequest {
  name: string;
  type: string;
  config: string; // JSON string
  console?: string;
  enabled?: boolean;
  weight?: number;
}

export interface UpdateProviderRequest extends Partial<CreateProviderRequest> {}