import React, { useState } from 'react';
import { Button, Card, Input } from '../components/ui';
import { useNavigate } from 'react-router-dom';
import { authApi } from '../services/auth';

export const LoginPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await authApi.login({ username, password });
      navigate('/');
    } catch (err) {
      setError('登录失败，请检查您的凭据。');
      console.error('Login failed:', err);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <Card className="w-full max-w-sm p-8 space-y-6">
        <h1 className="text-2xl font-bold text-center">登录</h1>
        <form onSubmit={handleLogin} className="space-y-4">
          <Input
            label="用户名"
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <Input
            label="密码"
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {error && <p className="text-sm text-red-600">{error}</p>}
          <Button type="submit" className="w-full">
            登录
          </Button>
        </form>
      </Card>
    </div>
  );
};

export default LoginPage;