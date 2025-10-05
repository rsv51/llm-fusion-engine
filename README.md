# LLM Fusion Engine

一个现代化、模块化的 LLM 网关，结合了 llmio 和 OrchestrationApi 的优点。

## ✨ 特性

- **多提供商支持**: 跨多个 LLM 提供商路由请求
- **智能负载均衡**: 基于优先级和轮询的路由策略
- **密钥管理**: 安全的 API 密钥轮换和管理
- **分组配置**: 将提供商组织成逻辑组
- **OpenAI 兼容 API**: OpenAI 端点的直接替代品
- **现代化 Web UI**: 基于 React 的管理界面
- **Docker 就绪**: 单命令部署

## 🚀 快速开始

### 使用 Docker

```bash
# 拉取镜像
docker pull ghcr.io/YOUR_USERNAME/llm-fusion-engine:latest

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name llm-fusion-engine \
  ghcr.io/YOUR_USERNAME/llm-fusion-engine:latest
```

访问 `http://localhost:8080` 查看 Web 管理界面。

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/YOUR_USERNAME/llm-fusion-engine.git
cd llm-fusion-engine

# 构建 Docker 镜像
docker build -t llm-fusion-engine:latest .

# 运行容器
docker run -d -p 8080:8080 llm-fusion-engine:latest
```

## 📖 使用指南

### 1. 访问管理界面

打开浏览器访问 `http://localhost:8080`，首次访问会显示登录页面。

**默认管理员账号：**
- 用户名：`admin`
- 密码：`admin`

> ⚠️ **安全提示**：首次登录后，建议立即修改默认密码以确保系统安全。

登录后，你会看到：

- **Dashboard**: 系统统计概览
  - 总组数、已启用组数
  - 提供商总数
  - API 密钥总数
  
- **Groups**: 管理提供商组
  - 创建新组
  - 配置优先级
  - 启用/禁用组

### 2. 创建第一个组

1. 点击 "Groups" 导航到组管理页面
2. 点击 "New Group" 按钮
3. 填写表单：
   - **Name**: 组名称（例如：`openai-group`）
   - **Priority**: 优先级（数字越大优先级越高）
   - **Enabled**: 是否启用该组
4. 点击 "Create Group"

### 3. 使用 API

#### OpenAI 兼容端点

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

#### 管理 API

获取系统统计：
```bash
curl http://localhost:8080/api/admin/stats
```

获取所有组：
```bash
curl http://localhost:8080/api/admin/groups
```

创建新组：
```bash
curl -X POST http://localhost:8080/api/admin/groups \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "new-group",
    "Priority": 10,
    "Enabled": true
  }'
```

## 🏗️ 架构设计

### 核心组件

```
┌─────────────────────────────────────────┐
│           Web UI (React)                │
│  Dashboard | Groups | Providers | Keys  │
└─────────────────┬───────────────────────┘
                  │ HTTP
┌─────────────────▼───────────────────────┐
│         API Layer (Gin)                 │
│  /v1/* (OpenAI)  |  /api/* (Admin)     │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Service Layer                      │
│  ┌────────────┐  ┌──────────────┐      │
│  │  Provider  │  │     Key      │      │
│  │   Router   │  │   Manager    │      │
│  └────────────┘  └──────────────┘      │
│  ┌────────────────────────────┐        │
│  │  MultiProviderService      │        │
│  └────────────────────────────┘        │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Database Layer (GORM)              │
│  Groups | Providers | Keys | Stats      │
│            SQLite                        │
└─────────────────────────────────────────┘
```

### 主要服务

- **ProviderRouter**: 智能路由逻辑
  - 优先级路由
  - 轮询负载均衡
  - 故障转移
  
- **KeyManager**: 密钥管理
  - 安全存储
  - 自动轮换
  - 使用跟踪
  
- **MultiProviderService**: 请求编排
  - 提供商选择
  - 请求转发
  - 响应聚合

## 📡 API 参考

### OpenAI 兼容端点

#### POST /v1/chat/completions

与 OpenAI Chat Completions API 完全兼容。

**请求示例：**
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  "temperature": 0.7,
  "stream": false
}
```

### 管理 API 端点

#### GET /api/admin/stats

获取系统统计信息。

**响应示例：**
```json
{
  "total_groups": 5,
  "enabled_groups": 3,
  "total_providers": 12,
  "total_keys": 25
}
```

#### GET /api/admin/groups

获取所有组列表。

#### POST /api/admin/groups

创建新组。

**请求体：**
```json
{
  "Name": "group-name",
  "Priority": 10,
  "Enabled": true
}
```

#### GET /api/admin/groups/:id

获取指定组的详情。

#### PUT /api/admin/groups/:id

更新指定组。

#### DELETE /api/admin/groups/:id

删除指定组。

## ⚙️ 配置

### 环境变量

- `PORT` - 服务器端口（默认：8080）
- `DB_PATH` - 数据库文件路径（默认：fusion.db）
- `GIN_MODE` - Gin 运行模式（release/debug，默认：debug）

### 用户管理

系统在首次启动时会自动创建默认管理员账号：
- 用户名：`admin`
- 密码：`admin`

如果数据库中已存在用户，则不会创建默认账号。

### 数据持久化

系统使用 SQLite 存储配置。建议使用 Docker volume 持久化数据：

```bash
docker run -d \
  -p 8080:8080 \
  -v ./data:/app/data \
  -e DB_PATH=/app/data/fusion.db \
  llm-fusion-engine:latest
```

## 🛠️ 开发指南

### 前置要求

- Go 1.21+
- Node.js 18+
- pnpm

### 后端开发

```bash
cd llm-fusion-engine

# 安装依赖
go mod download

# 运行开发服务器
go run cmd/server/main.go
```

### 前端开发

```bash
cd web

# 安装依赖
pnpm install

# 运行开发服务器
pnpm dev
```

### 构建

```bash
# 构建后端
go build -o server ./cmd/server

# 构建前端
cd web
pnpm build
```

## 🐳 Docker 部署

### 单容器部署

```bash
docker build -t llm-fusion-engine:latest .
docker run -d -p 8080:8080 llm-fusion-engine:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  llm-fusion-engine:
    image: ghcr.io/YOUR_USERNAME/llm-fusion-engine:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - DB_PATH=/app/data/fusion.db
      - GIN_MODE=release
    restart: unless-stopped
```

## 📝 技术栈

### 后端
- **Go 1.21**: 高性能后端语言
- **Gin**: Web 框架
- **GORM**: ORM 库
- **SQLite**: 嵌入式数据库（glebarez/sqlite - 纯 Go 实现）

### 前端
- **React 18**: UI 框架
- **TypeScript**: 类型安全
- **Vite**: 构建工具
- **Tailwind CSS**: 样式框架

### DevOps
- **Docker**: 容器化
- **GitHub Actions**: CI/CD

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

## 📄 许可证

MIT License
