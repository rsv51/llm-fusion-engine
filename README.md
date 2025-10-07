# LLM Fusion Engine

一个现代化、企业级的 LLM 网关系统,提供统一的 API 接口管理多个 LLM 提供商,支持智能路由、负载均衡、健康监控等功能。

## ✨ 核心特性

### 🔄 多提供商管理
- **灵活的提供商配置**: 支持 OpenAI、Anthropic、Gemini 等多种 LLM 提供商
- **动态添加/修改**: 通过 Web UI 实时管理提供商配置
- **健康状态监控**: 自动检测提供商健康状态(健康/降级/不可用)
- **性能指标追踪**: 实时显示延迟、响应时间等关键指标

### 📊 模型管理
- **统一模型目录**: 集中管理所有提供商的可用模型
- **模型映射**: 将客户端请求的模型名映射到实际提供商模型
- **批量导入/导出**: 支持 Excel 批量导入模型配置
- **模型分组**: 按提供商或用途组织模型

### 🎯 智能路由
- **优先级路由**: 基于提供商优先级的智能路由策略
- **负载均衡**: 轮询算法实现请求分发
- **自动故障转移**: 检测到故障时自动切换到备用提供商
- **请求日志**: 详细记录每次 API 调用的完整信息

### 🛡️ 企业级功能
- **认证与授权**: 基于 Token 的安全认证机制
- **请求限流**: 可配置的请求速率限制
- **数据持久化**: SQLite 数据库存储配置和日志
- **健康检查**: 定期自动检查提供商可用性
- **详细日志**: 完整的请求/响应日志追踪

### 🎨 现代化 Web UI
- **仪表盘**: 实时统计数据和系统健康状态
- **提供商管理**: 可视化管理提供商配置
- **模型配置**: 图形化界面管理模型映射
- **日志查看**: 实时查看和搜索请求日志
- **响应式设计**: 支持桌面和移动设备访问

## 🚀 快速开始

### 使用 Docker

```bash
# 拉取镜像
docker pull ghcr.io/rsv51/llm-fusion-engine:latest

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name llm-fusion-engine \
  ghcr.io/rsv51/llm-fusion-engine:latest
```

访问 `http://localhost:8080` 查看 Web 管理界面。

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/rsv51/llm-fusion-engine.git
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

登录后，你会看到五个主要功能模块：

#### 📊 Dashboard (仪表盘)
- **系统统计**: 提供商总数、模型总数、请求总数
- **提供商状态**: 实时显示各提供商健康状态
- **延迟监控**: 显示每个提供商的响应延迟
- **快速操作**: 快速访问常用功能

#### 🔧 Providers (提供商管理)
- **添加提供商**: 配置新的 LLM 提供商
  - 提供商名称和类型(OpenAI/Anthropic/Gemini等)
  - Base URL 和 API Key
  - 优先级设置
- **健康监控**: 实时显示提供商状态
  - 🟢 健康: 正常可用
  - 🟡 降级: 部分功能受限(如认证问题)
  - 🔴 不可用: 无法连接或错误
- **性能指标**: 查看延迟和最后检查时间
- **批量管理**: 启用/禁用/删除提供商

#### 🎯 Models (模型配置)
- **模型列表**: 查看所有可用模型
- **模型详情**: 查看模型关联的提供商
- **批量导入**: 通过 Excel 批量导入模型配置
- **批量导出**: 导出当前配置用于备份

#### 🔀 Model Mappings (模型映射)
- **映射规则**: 定义客户端模型名到提供商模型的映射
- **优先级控制**: 设置不同映射的优先级
- **灵活路由**: 支持一对多映射实现负载均衡

#### 📝 Logs (请求日志)
- **实时日志**: 查看所有 API 请求记录
- **详细信息**: 
  - 请求时间和响应时间
  - 使用的模型和提供商
  - Token 使用统计
  - HTTP 状态码
- **日志筛选**: 按提供商、模型、状态码筛选
- **详情查看**: 查看完整的请求和响应内容

### 2. 配置第一个提供商

1. 点击左侧菜单的 **"Providers"**
2. 点击右上角的 **"New Provider"** 按钮
3. 填写提供商信息：
   ```
   Name: my-openai-provider
   Type: openai
   Base URL: https://api.openai.com/v1
   API Key: sk-your-api-key-here
   Priority: 100
   Enabled: ✓
   ```
4. 点击 **"Create Provider"**
5. 系统会自动进行健康检查，显示提供商状态

### 3. 添加模型配置

有两种方式添加模型：

#### 方式一：手动添加
1. 进入 **"Models"** 页面
2. 点击 **"New Model"**
3. 填写：
   - Model Name: `gpt-3.5-turbo`
   - Provider: 选择刚创建的提供商
   - Enabled: ✓

#### 方式二：Excel 批量导入
1. 下载模板：在 **"Models"** 页面点击 **"Export Excel"** 获取模板
2. 编辑 Excel 文件，添加多个模型
3. 点击 **"Import Excel"** 上传文件
4. 系统自动批量导入

### 4. 配置模型映射(可选)

如果需要将客户端的模型名映射到不同的实际模型：

1. 进入 **"Model Mappings"** 页面
2. 点击 **"New Mapping"**
3. 配置映射规则：
   ```
   Client Model: gpt-4
   Provider Model: gpt-3.5-turbo
   Provider: my-openai-provider
   Priority: 10
   ```
4. 保存后，客户端请求 `gpt-4` 时会自动路由到 `gpt-3.5-turbo`

### 5. 使用 API 服务

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

**流式响应：**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ],
    "stream": true
  }'
```

#### 管理 API

**获取系统统计：**
```bash
curl http://localhost:8080/api/admin/stats
```

**获取所有提供商：**
```bash
curl http://localhost:8080/api/admin/providers
```

**创建新提供商：**
```bash
curl -X POST http://localhost:8080/api/admin/providers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-provider",
    "type": "openai",
    "config": "{\"baseUrl\":\"https://api.openai.com/v1\",\"apiKey\":\"sk-xxx\"}",
    "priority": 100,
    "enabled": true
  }'
```

**触发健康检查：**
```bash
curl -X POST http://localhost:8080/api/admin/providers/1/health-check
```

**获取请求日志：**
```bash
curl http://localhost:8080/api/admin/logs?page=1&page_size=20
```

### 6. 监控和维护

#### 健康检查机制
系统每5分钟自动检查所有提供商的健康状态：
- ✅ **Healthy**: 提供商正常，延迟在可接受范围内
- ⚠️ **Degraded**: 提供商可用但存在问题(如认证错误)
- ❌ **Unhealthy**: 提供商不可用或响应超时

健康检查通过发送测试聊天请求实现：
- OpenAI 类型: 使用 `gpt-3.5-turbo` 模型
- Anthropic 类型: 使用 `claude-3-haiku-20240307` 模型
- Gemini 类型: 使用 `gemini-1.5-flash` 模型

#### 日志管理
- 所有请求都会记录在数据库中
- 包含完整的请求/响应内容
- 可通过 Web UI 查看和搜索
- 支持按时间、提供商、模型筛选

#### 数据库维护
定期备份数据库文件 `fusion.db`：
```bash
# 停止服务
docker stop llm-fusion-engine

# 备份数据库
cp data/fusion.db data/fusion.db.backup

# 重启服务
docker start llm-fusion-engine
```

## 🏗️ 架构设计

### 核心组件

```
┌─────────────────────────────────────────┐
│           Web UI (React)                │
│  Dashboard | Providers | Models | Logs  │
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
│  │  Provider  │  │ Health Check │      │
│  │   Router   │  │  Service     │      │
│  └────────────┘  └──────────────┘      │
│  ┌────────────────────────────┐        │
│  │  MultiProviderService      │        │
│  └────────────────────────────┘        │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Database Layer (GORM)              │
│  Providers | Models | Logs | Health     │
│            SQLite                        │
└─────────────────────────────────────────┘
```

### 主要服务

- **ProviderRouter**: 智能路由逻辑
  - 优先级路由
  - 轮询负载均衡
  - 故障转移

- **HealthChecker**: 健康检查服务
  - 定期自动检查
  - 实时状态更新
  - 性能指标收集

- **MultiProviderService**: 请求编排
  - 提供商选择
  - 请求转发
  - 响应聚合

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
    image: ghcr.io/rsv51/llm-fusion-engine:latest
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
