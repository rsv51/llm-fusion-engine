# LLM Fusion Engine - 现代化 Web 界面设计文档

## 1. 设计概述

### 1.1 设计目标
- 综合 llm-orchestrator-py 和 OrchestrationApi-main 两个项目的优点
- 打造现代化、高效、易用的管理界面
- 实现模块化、可维护的代码架构
- 提供出色的用户体验和响应式设计

### 1.2 技术栈选择

**核心技术**:
- **React 18** - 现代化 UI 框架
- **TypeScript** - 类型安全
- **Vite** - 快速构建工具
- **Tailwind CSS** - 实用优先的 CSS 框架
- **React Router v6** - 客户端路由

**关键依赖**:
- `react-router-dom` - 页面路由管理
- `axios` - HTTP 请求客户端
- `chart.js` + `react-chartjs-2` - 数据可视化
- `lucide-react` - 现代化图标库
- `react-hook-form` - 表单管理
- `zustand` - 轻量级状态管理
- `xlsx` - Excel 导入导出

## 2. 架构设计

### 2.1 项目目录结构

```
llm-fusion-engine/web/src/
├── components/          # 可复用组件
│   ├── ui/             # 基础 UI 组件
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   ├── Modal.tsx
│   │   ├── Input.tsx
│   │   ├── Select.tsx
│   │   ├── Table.tsx
│   │   ├── Badge.tsx
│   │   └── Sidebar.tsx
│   ├── charts/         # 图表组件
│   │   ├── LineChart.tsx
│   │   ├── BarChart.tsx
│   │   └── PieChart.tsx
│   └── layout/         # 布局组件
│       ├── Layout.tsx
│       ├── Header.tsx
│       └── Sidebar.tsx
├── pages/              # 页面组件
│   ├── Dashboard.tsx
│   ├── Groups.tsx
│   ├── Keys.tsx
│   ├── Models.tsx
│   ├── ProxyKeys.tsx
│   ├── Logs.tsx
│   └── Settings.tsx
├── services/           # API 服务
│   ├── api.ts
│   ├── groups.ts
│   ├── keys.ts
│   ├── models.ts
│   ├── proxyKeys.ts
│   └── logs.ts
├── hooks/              # 自定义 Hooks
│   ├── useGroups.ts
│   ├── useKeys.ts
│   ├── useModels.ts
│   └── useHealth.ts
├── types/              # TypeScript 类型
│   ├── group.ts
│   ├── key.ts
│   ├── model.ts
│   └── common.ts
├── store/              # 状态管理
│   ├── authStore.ts
│   └── uiStore.ts
├── utils/              # 工具函数
│   ├── format.ts
│   ├── validation.ts
│   └── export.ts
└── App.tsx             # 根组件
```

### 2.2 导航结构

采用**侧边栏导航 + 顶部状态栏**的布局:

```
+------------------+--------------------------------+
| Logo             | System Status    User    Logout|
+------------------+--------------------------------+
| 📊 Dashboard     |                                |
| 👥 Groups        |                                |
| 🔑 Keys          |        Main Content            |
| 🤖 Models        |                                |
| 🔐 Proxy Keys    |                                |
| 📝 Logs          |                                |
| ⚙️  Settings     |                                |
+------------------+--------------------------------+
```

## 3. 功能模块设计

### 3.1 Dashboard (仪表盘)

**功能概览**:
- 系统运行状态 (运行时间、版本、数据库状态)
- 关键指标统计卡片 (总请求数、成功率、平均延迟、Token 使用)
- 请求趋势图表 (Chart.js 折线图)
- 提供商健康状态可视化 (参考 llm-orchestrator-py 的状态条设计)
- Token 使用分布 (饼图)

**借鉴来源**:
- OrchestrationApi-main: 系统健康监控卡片设计
- llm-orchestrator-py: 状态条可视化、统计数据展示

### 3.2 Groups (服务商分组管理)

**核心功能**:
- 分组列表展示 (表格 + 卡片视图)
- 完整的 CRUD 操作
- 分区表单设计 (基本信息、服务商配置、高级参数、代理设置)
- 分组启用/禁用
- 批量操作 (批量删除、批量导入导出)
- 分组健康检查

**表单设计** (参考 OrchestrationApi-main):
```
┌─ 基本信息 ─────────────────────┐
│ • 分组名称                      │
│ • 描述                          │
│ • 优先级                        │
└────────────────────────────────┘

┌─ 服务商配置 ───────────────────┐
│ • 服务商类型选择                │
│ • API Key 列表                  │
│ • Base URL                      │
└────────────────────────────────┘

┌─ 高级参数 ─────────────────────┐
│ • 负载均衡策略                  │
│ • 重试次数                      │
│ • 超时时间                      │
│ • 健康检查间隔                  │
└────────────────────────────────┘

┌─ 代理设置 ─────────────────────┐
│ • 代理服务器地址                │
│ • 认证信息                      │
└────────────────────────────────┘
```

**借鉴来源**:
- OrchestrationApi-main: 完整的分区表单设计
- llm-orchestrator-py: 简洁的提供商管理逻辑

### 3.3 Keys (API 密钥管理)

**核心功能**:
- 密钥列表展示
- 添加/编辑/删除密钥
- 密钥验证状态 (有效/无效/未检测)
- 批量验证
- 批量导入导出 (Excel)
- 密钥使用统计

**状态指示**:
```
✅ 有效 (绿色)
❌ 无效 (红色)
⏳ 验证中 (黄色)
❓ 未检测 (灰色)
```

**借鉴来源**:
- OrchestrationApi-main: 密钥验证和状态管理
- llm-orchestrator-py: Excel 批量操作

### 3.4 Models (模型配置)

**核心功能**:
- 模型列表管理
- 添加/编辑/删除模型
- 模型-服务商关联配置
- 模型能力标记 (工具调用、结构化输出、视觉)
- 批量导入导出
- 从服务商获取可用模型 (参考 llm-orchestrator-py)

**模型导入流程** (参考 llm-orchestrator-py):
```
1. 选择服务商
2. 获取可用模型列表
3. 预览模型
4. 一键导入全部 / 选择性导入
5. 自动创建模型配置和关联
```

**借鉴来源**:
- llm-orchestrator-py: 提供商模型一键导入功能
- OrchestrationApi-main: 模型重命名和别名智能推荐

### 3.5 Proxy Keys (代理密钥管理)

**核心功能**:
- 代理密钥生成
- 权限控制 (允许的模型、速率限制)
- 分组选择策略 (轮询、权重、随机、故障转移)
- 使用统计查看
- 密钥启用/禁用

**选择策略配置** (参考 OrchestrationApi-main):
```
┌─ 选择策略 ─────────────────────┐
│ ○ 轮询 (Round Robin)           │
│ ○ 权重 (Weight)                │
│ ○ 随机 (Random)                │
│ ○ 故障转移 (Failover)          │
└────────────────────────────────┘

┌─ 分组权重配置 ─────────────────┐
│ Group A: ████░░░░░░ 40%        │
│ Group B: ██████░░░░ 60%        │
└────────────────────────────────┘
```

**借鉴来源**:
- OrchestrationApi-main: 完整的代理密钥系统和选择策略

### 3.6 Logs (请求日志)

**核心功能**:
- 日志列表展示
- 分页、搜索、筛选
- 按时间、模型、服务商筛选
- 日志详情查看
- 导出日志

**借鉴来源**:
- llm-orchestrator-py: 简洁的日志展示
- OrchestrationApi-main: 分页和筛选功能

### 3.7 Settings (系统设置)

**核心功能**:
- API 密钥配置
- 管理员密钥设置
- 系统参数配置
- Excel 批量管理 (导入/导出配置模板)
- 用户管理

**Excel 批量管理** (参考 llm-orchestrator-py):
```
📥 导出所有配置 (3工作表: Providers, Models, Associations)
📄 下载空白模板
📝 下载带示例模板
📤 导入配置
```

**借鉴来源**:
- llm-orchestrator-py: Excel 批量导入导出功能
- OrchestrationApi-main: 用户管理和系统设置

## 4. UI 组件设计

### 4.1 基础组件

#### Button 组件
```typescript
interface ButtonProps {
  variant: 'primary' | 'secondary' | 'danger' | 'success'
  size: 'sm' | 'md' | 'lg'
  loading?: boolean
  icon?: ReactNode
  onClick?: () => void
}
```

#### Card 组件
```typescript
interface CardProps {
  title?: string
  actions?: ReactNode
  children: ReactNode
  className?: string
}
```

#### Modal 组件
```typescript
interface ModalProps {
  isOpen: boolean
  title: string
  onClose: () => void
  children: ReactNode
  footer?: ReactNode
}
```

#### Table 组件
```typescript
interface TableProps<T> {
  columns: Column<T>[]
  data: T[]
  pagination?: PaginationConfig
  onRowClick?: (row: T) => void
  selection?: SelectionConfig
}
```

### 4.2 特色组件

#### StatusBadge (状态徽章)
```typescript
interface StatusBadgeProps {
  status: 'success' | 'error' | 'warning' | 'info'
  text: string
}
```

颜色系统:
- `success`: 绿色 (#10b981)
- `error`: 红色 (#ef4444)
- `warning`: 黄色 (#f59e0b)
- `info`: 蓝色 (#3b82f6)

#### HealthIndicator (健康指示器)
参考 llm-orchestrator-py 的状态条设计:
```typescript
interface HealthIndicatorProps {
  history: boolean[]  // 最近10次状态
}
```

渲染为状态条:
```
████████░░  (8/10 成功)
```

#### Chart (图表组件)
基于 Chart.js 封装:
```typescript
// LineChart - 请求趋势
// BarChart - Token 使用对比
// PieChart - 服务商分布
```

### 4.3 布局组件

#### Layout
```
┌─────────────────────────────────────┐
│ Header (Logo + Status + User)       │
├─────────┬───────────────────────────┤
│         │                           │
│ Sidebar │     Main Content          │
│         │                           │
│         │                           │
└─────────┴───────────────────────────┘
```

#### Sidebar
- 响应式设计 (移动端可折叠)
- 图标 + 文字导航
- 活动状态高亮

## 5. 交互设计

### 5.1 表单交互

**智能验证** (react-hook-form):
- 实时验证
- 错误提示
- 必填字段标记

**分步表单** (复杂配置):
```
Step 1: 基本信息 → Step 2: 服务商配置 → Step 3: 高级设置 → Step 4: 确认
```

### 5.2 批量操作

**选择模式**:
- 单选框选择
- 全选/反选
- 批量操作按钮 (启用/禁用/删除/导出)

**操作确认**:
- 危险操作需二次确认
- 显示受影响的项目数量

### 5.3 实时更新

**状态轮询**:
- 健康状态每 30 秒更新
- 统计数据每 60 秒更新

**加载状态**:
- Skeleton 占位符
- 加载动画
- 进度条

### 5.4 响应式设计

**断点系统** (Tailwind):
```
sm: 640px   (手机横屏)
md: 768px   (平板)
lg: 1024px  (笔记本)
xl: 1280px  (桌面)
2xl: 1536px (大屏)
```

**移动端优化**:
- 侧边栏折叠为汉堡菜单
- 表格切换为卡片视图
- 触摸友好的按钮尺寸 (min-height: 44px)

## 6. 设计系统

### 6.1 颜色系统

```css
/* 主色调 */
--primary: #2563eb;        /* 蓝色 */
--primary-hover: #1d4ed8;

/* 功能色 */
--success: #10b981;        /* 绿色 */
--danger: #ef4444;         /* 红色 */
--warning: #f59e0b;        /* 黄色 */
--info: #3b82f6;           /* 浅蓝 */

/* 中性色 */
--gray-50: #f9fafb;
--gray-100: #f3f4f6;
--gray-200: #e5e7eb;
--gray-500: #6b7280;
--gray-900: #111827;
```

### 6.2 间距系统

使用 Tailwind 的间距比例:
```
4px, 8px, 12px, 16px, 20px, 24px, 32px, 40px, 48px, 64px
```

### 6.3 圆角系统

```
rounded-sm: 2px
rounded: 4px
rounded-md: 6px
rounded-lg: 8px
rounded-xl: 12px
```

### 6.4 阴影系统

```
shadow-sm: 0 1px 2px rgba(0,0,0,0.05)
shadow: 0 1px 3px rgba(0,0,0,0.1)
shadow-md: 0 4px 6px rgba(0,0,0,0.1)
shadow-lg: 0 10px 15px rgba(0,0,0,0.1)
```

## 7. 性能优化

### 7.1 代码分割

```typescript
// 路由懒加载
const Dashboard = lazy(() => import('./pages/Dashboard'))
const Groups = lazy(() => import('./pages/Groups'))
// ...
```

### 7.2 数据缓存

使用自定义 Hooks 实现数据缓存:
```typescript
const { data, loading, refetch } = useGroups({
  cache: true,
  cacheTime: 5 * 60 * 1000 // 5分钟
})
```

### 7.3 虚拟滚动

对于长列表使用虚拟滚动:
```typescript
import { useVirtualizer } from '@tanstack/react-virtual'
```

## 8. 实施计划

### Phase 1: 基础架构 (第1-2天)
- [x] 创建项目目录结构
- [ ] 安装必要依赖
- [ ] 配置路由
- [ ] 创建基础布局组件

### Phase 2: UI 组件库 (第3-4天)
- [ ] 实现基础组件 (Button, Card, Modal, Input, Select)
- [ ] 实现表格组件 (支持分页、搜索、筛选)
- [ ] 实现图表组件 (Chart.js 封装)
- [ ] 实现状态组件 (Badge, HealthIndicator)

### Phase 3: API 服务层 (第5天)
- [ ] 创建 API 客户端
- [ ] 实现各模块 API 服务
- [ ] 创建自定义 Hooks

### Phase 4: 页面开发 (第6-10天)
- [ ] Dashboard 页面
- [ ] Groups 页面
- [ ] Keys 页面
- [ ] Models 页面
- [ ] Proxy Keys 页面
- [ ] Logs 页面
- [ ] Settings 页面

### Phase 5: 测试和优化 (第11-12天)
- [ ] 功能测试
- [ ] 性能优化
- [ ] 响应式测试
- [ ] 文档完善

## 9. 技术亮点

### 9.1 从 llm-orchestrator-py 借鉴
- ✅ Excel 批量导入导出功能
- ✅ 状态历史可视化 (状态条)
- ✅ 提供商模型一键导入
- ✅ 简洁的 API 请求封装

### 9.2 从 OrchestrationApi-main 借鉴
- ✅ 完整的分区表单设计
- ✅ API 密钥验证和状态管理
- ✅ 代理密钥选择策略配置
- ✅ Chart.js 图表可视化
- ✅ 模型重命名和别名推荐
- ✅ 响应式移动端优化
- ✅ 自定义模态框动画

### 9.3 创新改进
- ✅ React + TypeScript 现代化架构
- ✅ 组件化和模块化设计
- ✅ 侧边栏导航 (替代 Tab 切换)
- ✅ React Router 客户端路由
- ✅ 统一的设计系统和颜色方案
- ✅ 性能优化 (懒加载、缓存、虚拟滚动)

## 10. 用户体验优化

### 10.1 加载体验
- Skeleton 占位符
- 加载动画
- 进度指示器

### 10.2 错误处理
- 友好的错误提示
- 错误边界组件
- 重试机制

### 10.3 快捷操作
- 键盘快捷键支持
- 批量操作
- 快速搜索

### 10.4 反馈机制
- Toast 通知
- 操作确认
- 成功/失败提示

---

**设计版本**: 1.0  
**最后更新**: 2025-10-05  
**设计师**: Kilo Code (AI Assistant)