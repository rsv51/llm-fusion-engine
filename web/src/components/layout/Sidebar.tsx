import React from 'react'
import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Users,
  Key,
  Box,
  FileText,
  Settings,
  X
} from 'lucide-react'

interface SidebarProps {
  isOpen: boolean
  onToggle: () => void
}

const navigation = [
  { name: '仪表盘', href: '/', icon: LayoutDashboard },
  { name: '分组管理', href: '/groups', icon: Users },
  { name: '密钥管理', href: '/keys', icon: Key },
  { name: '模型配置', href: '/models', icon: Box },
  { name: '请求日志', href: '/logs', icon: FileText },
  { name: '系统设置', href: '/settings', icon: Settings },
]

export const Sidebar: React.FC<SidebarProps> = ({ isOpen, onToggle }) => {
  return (
    <>
      {/* 移动端遮罩 */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={onToggle}
        />
      )}
      
      {/* 侧边栏 */}
      <div 
        className={`fixed lg:static inset-y-0 left-0 z-50 w-64 bg-white border-r border-gray-200 transform transition-transform duration-200 ease-in-out ${
          isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        }`}
      >
        {/* Logo */}
        <div className="flex items-center justify-between h-16 px-6 border-b border-gray-200">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold">
              LF
            </div>
            <span className="text-lg font-semibold text-gray-900">LLM Fusion</span>
          </div>
          <button 
            onClick={onToggle}
            className="lg:hidden text-gray-500 hover:text-gray-700"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* 导航菜单 */}
        <nav className="flex-1 px-4 py-4 space-y-1 overflow-y-auto">
          {navigation.map((item) => (
            <NavLink
              key={item.name}
              to={item.href}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                }`
              }
            >
              <item.icon className="w-5 h-5" />
              <span>{item.name}</span>
            </NavLink>
          ))}
        </nav>

        {/* 底部信息 */}
        <div className="px-6 py-4 border-t border-gray-200">
          <div className="text-xs text-gray-500">
            <p>Version 1.0.0</p>
            <p className="mt-1">© 2025 LLM Fusion Engine</p>
          </div>
        </div>
      </div>
    </>
  )
}

export default Sidebar