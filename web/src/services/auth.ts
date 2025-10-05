import api from './api';
import { LoginRequest, LoginResponse } from '../types';

export const authApi = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = (await api.post('/auth/login', credentials)) as LoginResponse;
    if (response.token) {
      localStorage.setItem('token', response.token);
    }
    return response;
  },

  logout() {
    localStorage.removeItem('token');
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },
};

export default authApi;