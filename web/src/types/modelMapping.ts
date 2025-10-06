import type { Provider } from './provider';
import type { Model } from './model';

export interface ModelProviderMapping {
  id: number;
  modelId: number;
  providerId: number;
  providerModel: string; // The actual model ID on the provider's platform
  toolCall?: boolean;
  structuredOutput?: boolean;
  image?: boolean;
  weight: number;
  enabled: boolean;
  model?: Model; // Preloaded model details
  provider?: Provider; // Preloaded provider details
}

export interface CreateModelProviderMappingRequest {
  modelId: number;
  providerId: number;
  providerModel: string;
  toolCall?: boolean;
  structuredOutput?: boolean;
  image?: boolean;
  weight?: number;
  enabled?: boolean;
}

export interface UpdateModelProviderMappingRequest extends Partial<CreateModelProviderMappingRequest> {}