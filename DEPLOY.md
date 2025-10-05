# LLM Fusion Engine 部署指南

本文档提供详细的部署步骤,帮助您快速部署 LLM Fusion Engine。

## 目录
- [环境要求](#环境要求)
- [本地开发部署](#本地开发部署)
- [Docker 部署](#docker-部署)
- [生产环境部署](#生产环境部署)
- [环境变量配置](#环境变量配置)
- [故障排除](#故障排除)

## 环境要求

### 本地开发
- Go 1.21 或更高版本
- Node.js 18 或更高版本
- npm 或 pnpm

### Docker 部署
- Docker 20.10 或更高版本
- Docker Compose 2.0 或更高版本

## 本地开发部署

### 1. 克隆项目
```bash
git clone <repository-url>
cd llm-fusion-engine
```

### 2. 安装依赖

**后端依赖:**
```bash
go mod download
```

**前端依赖:**
```bash
cd web
npm install
# 或使用 pnpm
pnpm install
```

### 3. 开发模式运行

**启动后端 (终端 1):**
```bash
go run cmd/server/main.go
```

**启动前端开发服务器 (终端 2):**
```bash
cd web
npm run dev
```

访问 http://localhost:5173 查看应用。

### 4. 构建生产版本

**构建前端:**
```bash
cd web
npm run build
```

**构建后端:**
```bash
go build -o llm-fusion-engine cmd/server/main.go
```

## Docker 部署

### 方式一: 使用 docker-compose (推荐)

1. **确保 Docker 和 Docker Compose 已安装**
   ```bash
   docker --version
   docker-compose --version
   ```

2. **构建并启动服务**
   ```bash
   docker-compose up -d
   ```

3. **查看日志**
   ```bash
   docker-compose logs -f
   ```

4. **停止服务**
   ```bash
   docker-compose down
   ```

### 方式二: 直接使用 Docker

1. **构建镜像**
   ```bash
   docker build -t llm-fusion-engine:latest .
   ```

2. **运行容器**
   ```bash
   docker run -d \
     --name llm-fusion-engine \
     -p 8080:8080 \
     -v $(pwd)/data:/app/data \
     llm-fusion-engine:latest
   ```

3. **查看日志**
   ```bash
   docker logs -f llm-fusion-engine
   ```

4. **停止容器**
   ```bash
   docker stop llm-fusion-engine
   docker rm llm-fusion-engine
   ```

## 生产环境部署

### 使用 Docker Hub

1. **登录 Docker Hub**
   ```bash
   docker login
   ```

2. **标记镜像**
   ```bash
   docker tag llm-fusion-engine:latest <your-username>/llm-fusion-engine:latest
   ```

3. **推送镜像**
   ```bash
   docker push <your-username>/llm-fusion-engine:latest
   ```

4. **在生产服务器上拉取并运行**
   ```bash
   docker pull <your-username>/llm-fusion-engine:latest
   docker run -d \
     --name llm-fusion-engine \
     -p 8080:8080 \
     -v /path/to/data:/app/data \
     --restart unless-stopped \
     <your-username>/llm-fusion-engine:latest
   ```

### 使用 GitHub Container Registry

1. **登录 GHCR**
   ```bash
   echo $GITHUB_TOKEN | docker login ghcr.io -u <username> --password-stdin
   ```

2. **标记并推送镜像**
   ```bash
   docker tag llm-fusion-engine:latest ghcr.io/<username>/llm-fusion-engine:latest
   docker push ghcr.io/<username>/llm-fusion-engine:latest
   ```

3. **在生产环境拉取**
   ```bash
   docker pull ghcr.io/<username>/llm-fusion-engine:latest
   docker run -d \
     --name llm-fusion-engine \
     -p 8080:8080 \
     -v /path/to/data:/app/data \
     --restart unless-stopped \
     ghcr.io/<username>/llm-fusion-engine:latest
   ```

## 环境变量配置

创建 `.env` 文件配置环境变量:

```env
# 服务器配置
PORT=8080
GIN_MODE=release

# 数据库配置
DATABASE_TYPE=sqlite
DATABASE_PATH=./data/fusion.db

# PostgreSQL 配置 (如果使用)
# DATABASE_TYPE=postgres
# DATABASE_HOST=localhost
# DATABASE_PORT=5432
# DATABASE_USER=postgres
# DATABASE_PASSWORD=your_password
# DATABASE_NAME=llm_fusion

# 日志配置
LOG_LEVEL=info
LOG_FILE=./logs/app.log

# 安全配置
ENABLE_AUTH=true
API_KEY_HEADER=X-API-Key
```

## 故障排除

### 问题 1: 前端无法连接后端 API

**症状:** 浏览器控制台显示 CORS 错误或 404

**解决方案:**
1. 检查后端是否正常运行: `curl http://localhost:8080/api/admin/stats`
2. 检查 Vite 代理配置 (开发模式): `web/vite.config.ts`
3. 确保后端路由正确注册

### 问题 2: Docker 构建失败

**症状:** `npm install` 或 `go build` 失败

**解决方案:**
1. 清理 Docker 缓存: `docker system prune -a -f`
2. 检查网络连接
3. 使用国内镜像源 (如果在中国):
   ```dockerfile
   # 在 Dockerfile 中添加
   RUN npm config set registry https://registry.npmmirror.com
   ```

### 问题 3: 页面路由404错误

**症状:** 刷新页面或直接访问 `/keys`、`/groups` 等路径返回 404

**解决方案:**
- 确保后端的 NoRoute 处理器正确配置为返回 `index.html`
- 检查 `main.go` 中的路由配置

### 问题 4: 数据库连接失败

**症状:** 启动时报错 "Failed to connect to database"

**解决方案:**
1. 检查数据库配置
2. 确保数据目录有写权限
3. 如果使用 PostgreSQL,确保数据库服务正在运行

### 问题 5: 端口已被占用

**症状:** "address already in use"

**解决方案:**
```bash
# 查找占用端口的进程
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# 修改端口
PORT=8081 go run cmd/server/main.go
```

## 性能优化建议

### 1. 使用反向代理

推荐使用 Nginx 作为反向代理:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. 启用 HTTPS

使用 Let's Encrypt 获取免费 SSL 证书:

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 3. 配置数据库连接池

在生产环境使用 PostgreSQL 并配置连接池:

```go
// 在 database.go 中
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 4. 启用日志轮转

使用 logrotate 管理日志文件:

```bash
/var/log/llm-fusion/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
}
```

## 监控和维护

### 健康检查

访问健康检查端点:
```bash
curl http://localhost:8080/api/health
```

### 查看系统统计

```bash
curl http://localhost:8080/api/admin/stats
```

### 数据备份

**SQLite 备份:**
```bash
cp data/fusion.db data/fusion.db.backup.$(date +%Y%m%d)
```

**PostgreSQL 备份:**
```bash
pg_dump -U postgres -d llm_fusion > backup.sql
```

## 更新部署

### 更新 Docker 镜像

```bash
# 拉取最新代码
git pull

# 重新构建镜像
docker-compose build

# 重启服务
docker-compose down
docker-compose up -d
```

### 滚动更新

```bash
# 构建新镜像
docker build -t llm-fusion-engine:v2 .

# 启动新容器
docker run -d --name llm-fusion-engine-v2 -p 8081:8080 llm-fusion-engine:v2

# 测试新版本
curl http://localhost:8081/api/health

# 切换流量 (更新 Nginx 配置)
# 停止旧容器
docker stop llm-fusion-engine
docker rm llm-fusion-engine
```

## 联系支持

如果遇到无法解决的问题,请:
1. 查看 [GitHub Issues](https://github.com/your-repo/issues)
2. 提交新的 Issue 并附上详细的错误信息和日志
3. 联系技术支持团队

## 相关文档

- [README.md](README.md) - 项目概览
- [API 文档](docs/API.md) - API 接口说明
- [架构设计](docs/ARCHITECTURE.md) - 系统架构