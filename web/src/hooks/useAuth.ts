import { useState, useEffect } from 'react';
import { authApi } from '../services/auth';
import { User } from '../types';

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(authApi.isAuthenticated());

  useEffect(() => {
    if (isAuthenticated) {
      // 在实际应用中，您应该在此处获取用户个人资料。
      // 目前，我们只设置一个模拟用户。
      // 注意：此模拟用户数据应在实现真实的用户获取逻辑后移除。
      const mockUser: User = {
        id: 1,
        username: 'admin',
        name: 'Admin User',
        email: 'admin@example.com',
        role: 1,
        status: 1,
      };
      setUser(mockUser);
    } else {
      setUser(null);
    }
  }, [isAuthenticated]);

  const login = async (credentials: {username: string, password?: string, github_token?: string}) => {
    const data = await authApi.login(credentials);
    setUser(data.user);
    setIsAuthenticated(true);
  };

  const logout = () => {
    authApi.logout();
    setUser(null);
    setIsAuthenticated(false);
  };

  return { user, isAuthenticated, login, logout };
};