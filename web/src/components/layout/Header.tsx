import React from 'react'
import { Menu, Bell, User, LogOut } from 'lucide-react'
import { Badge } from '../ui'
import { useAuth } from '../../hooks/useAuth'
import { useNavigate } from 'react-router-dom'

interface HeaderProps {
  onMenuClick: () => void
}

export const Header: React.FC<HeaderProps> = ({ onMenuClick }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="sticky top-0 z-30 h-16 bg-white border-b border-gray-200">
      <div className="flex items-center justify-between h-full px-4 lg:px-6">
        {/* 左侧 - 菜单按钮 */}
        <button
          onClick={onMenuClick}
          className="lg:hidden text-gray-500 hover:text-gray-700 transition-colors"
        >
          <Menu className="w-6 h-6" />
        </button>

        {/* 中间 - 系统状态 */}
        <div className="hidden lg:flex items-center gap-4">
          <Badge variant="success">系统正常</Badge>
        </div>

        {/* 右侧 - 操作按钮 */}
        <div className="flex items-center gap-3">
          {/* 通知 */}
          <button className="relative p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors">
            <Bell className="w-5 h-5" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>

          {/* 用户菜单 */}
          {user ? (
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors">
              <div className="w-8 h-8 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center">
                <User className="w-4 h-4" />
              </div>
              <span className="hidden sm:block text-sm font-medium text-gray-700">{user.username}</span>
            </div>
          ) : (
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors" onClick={() => navigate('/login')}>
              <div className="w-8 h-8 bg-gray-200 text-gray-600 rounded-full flex items-center justify-center">
                <User className="w-4 h-4" />
              </div>
              <span className="hidden sm:block text-sm font-medium text-gray-700">登录</span>
            </div>
          )}


          {/* 登出 */}
          {user && (
            <button
              onClick={handleLogout}
              className="p-2 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
              title="退出登录"
            >
              <LogOut className="w-5 h-5" />
            </button>
          )}
        </div>
      </div>
    </header>
  )
}

export default Header