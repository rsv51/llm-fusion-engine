import type { Provider } from './provider';

export interface ModelMapping {
  id: number;
  userFriendlyName: string;
  providerModelName: string;
  providerId: number;
  provider?: Provider;
}