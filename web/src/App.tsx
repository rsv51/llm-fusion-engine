import React, { useState } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Sidebar, Header } from './components/layout'
import { Dashboard, Groups, Keys, Logs, Models, Settings } from './pages'

export const App: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false)

  return (
    <BrowserRouter>
      <div className="flex h-screen bg-gray-50">
        {/* 侧边栏 */}
        <Sidebar isOpen={sidebarOpen} onToggle={() => setSidebarOpen(!sidebarOpen)} />

        {/* 主内容区 */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* 顶部栏 */}
          <Header onMenuClick={() => setSidebarOpen(!sidebarOpen)} />

          {/* 页面内容 */}
          <main className="flex-1 overflow-y-auto p-6">
            <Routes>
              <Route path="/" element={<Navigate to="/dashboard" replace />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/groups" element={<Groups />} />
              <Route path="/keys" element={<Keys />} />
              <Route path="/models" element={<Models />} />
              <Route path="/logs" element={<Logs />} />
              <Route path="/settings" element={<Settings />} />
              <Route path="*" element={<Navigate to="/dashboard" replace />} />
            </Routes>
          </main>
        </div>
      </div>
    </BrowserRouter>
  )
}

export default App