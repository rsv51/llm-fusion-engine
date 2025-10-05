# LLM Fusion Engine - 部署指南

本文档提供 LLM Fusion Engine 的详细部署说明,包括 Docker 部署和本地开发部署两种方式。

## 目录

- [系统要求](#系统要求)
- [Docker 部署(推荐)](#docker-部署推荐)
- [本地开发部署](#本地开发部署)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

## 系统要求

### 最低配置
- CPU: 2 核心
- 内存: 2GB RAM
- 磁盘: 10GB 可用空间
- 操作系统: Linux/Windows/macOS

### 推荐配置
- CPU: 4 核心
- 内存: 4GB RAM
- 磁盘: 20GB 可用空间
- 操作系统: Linux (Ubuntu 20.04+)

### 软件依赖
- Docker 20.10+ (Docker 部署)
- Docker Compose 2.0+ (Docker 部署)
- Go 1.21+ (本地开发)
- Node.js 18+ (本地开发)
- PostgreSQL 14+ (本地开发)

## Docker 部署(推荐)

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/llm-fusion-engine.git
cd llm-fusion-engine
```

### 2. 配置环境变量

创建 `.env` 文件:

```bash
cp .env.example .env
```

编辑 `.env` 文件,设置必要的环境变量:

```env
# 数据库配置
DB_HOST=postgres
DB_PORT=5432
DB_USER=llmfusion
DB_PASSWORD=your_secure_password_here
DB_NAME=llmfusion

# 应用配置
APP_PORT=8080
APP_ENV=production
LOG_LEVEL=info

# 管理员凭证
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_secure_admin_password_here
```

### 3. 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 检查服务状态
docker-compose ps
```

### 4. 访问应用

- Web 管理界面: http://localhost:8080
- API 端点: http://localhost:8080/v1
- 管理 API: http://localhost:8080/api

默认管理员账户:
- 用户名: admin
- 密码: (在 `.env` 文件中设置的 ADMIN_PASSWORD)

### 5. 停止服务

```bash
# 停止所有服务
docker-compose stop

# 停止并删除所有容器
docker-compose down

# 停止并删除所有容器和数据卷
docker-compose down -v
```

## 本地开发部署

### 1. 准备数据库

安装并启动 PostgreSQL:

```bash
# Ubuntu/Debian
sudo apt-get install postgresql-14

# macOS
brew install postgresql@14

# 启动 PostgreSQL
sudo service postgresql start  # Linux
brew services start postgresql@14  # macOS
```

创建数据库和用户:

```sql
CREATE DATABASE llmfusion;
CREATE USER llmfusion WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE llmfusion TO llmfusion;
```

### 2. 配置后端

```bash
cd backend

# 复制配置文件
cp config.example.yaml config.yaml

# 编辑配置文件
nano config.yaml
```

配置数据库连接:

```yaml
database:
  host: localhost
  port: 5432
  user: llmfusion
  password: your_password
  database: llmfusion
```

### 3. 启动后端服务

```bash
# 安装依赖
go mod download

# 运行数据库迁移
go run cmd/migrate/main.go

# 启动服务
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 4. 配置前端

```bash
cd web

# 安装依赖
npm install
```

### 5. 启动前端开发服务器

```bash
# 开发模式
npm run dev

# 生产构建
npm run build
npm run preview
```

前端开发服务器将在 `http://localhost:5173` 启动。

## 配置说明

### 环境变量说明

| 变量名 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `DB_HOST` | 数据库主机地址 | `localhost` | 是 |
| `DB_PORT` | 数据库端口 | `5432` | 是 |
| `DB_USER` | 数据库用户名 | `llmfusion` | 是 |
| `DB_PASSWORD` | 数据库密码 | - | 是 |
| `DB_NAME` | 数据库名称 | `llmfusion` | 是 |
| `APP_PORT` | 应用端口 | `8080` | 否 |
| `APP_ENV` | 运行环境 | `production` | 否 |
| `LOG_LEVEL` | 日志级别 | `info` | 否 |
| `ADMIN_USERNAME` | 管理员用户名 | `admin` | 是 |
| `ADMIN_PASSWORD` | 管理员密码 | - | 是 |

### 配置文件说明

后端配置文件 [`config.yaml`](backend/config.yaml:1) 包含以下主要配置项:

```yaml
# 服务器配置
server:
  port: 8080
  host: 0.0.0.0
  read_timeout: 30s
  write_timeout: 30s

# 数据库配置
database:
  host: localhost
  port: 5432
  user: llmfusion
  password: your_password
  database: llmfusion
  max_connections: 100
  max_idle_connections: 10

# 日志配置
logging:
  level: info
  format: json
  output: stdout

# 缓存配置 (可选)
cache:
  enabled: false
  redis_url: redis://localhost:6379
```

## 常见问题

### 1. 数据库连接失败

**问题**: 应用启动时提示数据库连接失败

**解决方案**:
- 检查数据库是否正在运行: `docker-compose ps` (Docker) 或 `systemctl status postgresql` (本地)
- 验证数据库连接信息是否正确
- 确保数据库用户有足够的权限
- 检查防火墙设置

### 2. 端口冲突

**问题**: 启动时提示端口已被占用

**解决方案**:
- 修改 `docker-compose.yml` 中的端口映射
- 或停止占用该端口的其他服务
- 使用 `netstat -tlnp | grep 8080` 查看端口占用情况

### 3. 前端无法连接后端

**问题**: 前端页面无法加载数据

**解决方案**:
- 检查后端服务是否正常运行
- 验证 API 代理配置是否正确
- 检查浏览器控制台的网络请求
- 确保 CORS 配置正确

### 4. Docker 构建失败

**问题**: Docker 镜像构建时出错

**解决方案**:
- 清理 Docker 缓存: `docker system prune -a`
- 检查 Dockerfile 语法
- 确保网络连接正常
- 增加 Docker 内存限制

### 5. 性能问题

**问题**: 系统响应缓慢

**解决方案**:
- 检查数据库连接池配置
- 增加服务器资源
- 启用 Redis 缓存
- 优化数据库查询
- 检查日志级别设置

## 更新和维护

### 更新应用

```bash
# Docker 部署
docker-compose pull
docker-compose up -d

# 本地部署
git pull
cd backend && go build
cd web && npm run build
```

### 备份数据

```bash
# 备份 PostgreSQL 数据库
docker-compose exec postgres pg_dump -U llmfusion llmfusion > backup.sql

# 恢复数据库
docker-compose exec -T postgres psql -U llmfusion llmfusion < backup.sql
```

### 查看日志

```bash
# Docker 部署
docker-compose logs -f app

# 本地部署
tail -f logs/app.log
```

## 安全建议

1. **修改默认密码**: 首次部署后立即修改管理员密码
2. **使用 HTTPS**: 生产环境建议配置 SSL/TLS 证书
3. **限制访问**: 使用防火墙限制不必要的端口访问
4. **定期更新**: 保持系统和依赖项的最新版本
5. **备份数据**: 定期备份数据库和配置文件

## 技术支持

如有问题或需要帮助,请:
- 查看项目 [README.md](README.md:1)
- 提交 Issue: https://github.com/yourusername/llm-fusion-engine/issues
- 查看文档: https://docs.llmfusion.example.com

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE:1) 文件。