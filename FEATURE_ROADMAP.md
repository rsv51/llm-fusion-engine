# LLM Fusion Engine 功能开发路线图

本文档规划了所有待实现功能的详细设计方案和实现步骤。

## 第一阶段:核心数据管理功能

### 1. 数据导入导出功能 (CSV/XLSX)

#### 1.1 后端实现

**新增数据模型:**
- 无需新增,使用现有的 Group, Provider, ApiKey 模型

**新增 API 端点:**
```go
// /api/admin/export/groups?format=csv|xlsx
GET /api/admin/export/groups

// /api/admin/export/keys?format=csv|xlsx  
GET /api/admin/export/keys

// /api/admin/export/providers?format=csv|xlsx
GET /api/admin/export/providers

// /api/admin/import/groups
POST /api/admin/import/groups

// /api/admin/import/keys
POST /api/admin/import/keys

// /api/admin/import/providers
POST /api/admin/import/providers
```

**所需 Go 包:**
```go
import (
    "github.com/360EntSecGroup-Skylar/excelize/v2" // Excel 处理
    "encoding/csv"                                  // CSV 处理
)
```

**实现步骤:**
1. 创建 `internal/api/admin/export_handler.go`
2. 创建 `internal/api/admin/import_handler.go`
3. 实现 CSV/XLSX 序列化和反序列化逻辑
4. 添加文件验证和错误处理
5. 注册路由到 main.go

#### 1.2 前端实现

**新增页面:**
- 在 Settings 页面添加导入导出选项卡

**功能点:**
- 下载按钮(选择格式:CSV/XLSX)
- 文件上传组件
- 导入进度显示
- 导入结果反馈

**所需库:**
```json
{
  "xlsx": "^0.18.5"  // SheetJS for Excel handling
}
```

---

### 2. 模型映射功能 (模型别名配置)

#### 2.1 数据库模型扩展

**新增字段到 Group 表:**
```go
type Group struct {
    // ... 现有字段
    ModelAliases string `gorm:"type:text"` // JSON: {"gpt-4":"openai-gpt-4"}
}
```

#### 2.2 API 端点

```go
// 获取模型映射
GET /api/admin/groups/:id/model-aliases

// 更新模型映射
PUT /api/admin/groups/:id/model-aliases
Body: {
  "aliases": {
    "gpt-4": "openai-gpt-4",
    "claude-3": "anthropic-claude-3"
  }
}

// 批量设置映射
POST /api/admin/model-aliases/batch
Body: {
  "groupIds": [1, 2, 3],
  "aliases": {...}
}
```

#### 2.3 前端实现

**新增组件:**
- `ModelAliasEditor` - 键值对编辑器
- 在 Groups 详情页面添加"模型映射"选项卡

---

### 3. Provider 管理功能完善

#### 3.1 数据库模型扩展

**Provider 表新增字段:**
```go
type Provider struct {
    // ... 现有字段
    BaseURL      string `gorm:"type:varchar(255)"`
    Timeout      int    `gorm:"default:30"`    // 秒
    MaxRetries   int    `gorm:"default:3"`
    HealthStatus string `gorm:"default:'unknown'"` // healthy/unhealthy/unknown
    LastChecked  *time.Time
}
```

#### 3.2 API 端点

```go
// Provider CRUD
POST   /api/admin/providers
GET    /api/admin/providers?page=1&pageSize=20
GET    /api/admin/providers/:id
PUT    /api/admin/providers/:id
DELETE /api/admin/providers/:id

// Provider 批量操作
POST   /api/admin/providers/batch-create
POST   /api/admin/providers/batch-delete
POST   /api/admin/providers/batch-update

// Provider 健康检查
POST   /api/admin/providers/:id/health-check
POST   /api/admin/providers/health-check-all
```

#### 3.3 前端页面

**新增页面: Providers.tsx**
- Provider 列表(卡片/表格视图)
- 创建/编辑 Provider 模态框
- 健康状态指示器
- 批量操作工具栏

---

## 第二阶段:增强功能

### 4. 模型按名称分页加载

#### 4.1 后端实现

**新增 Model 表:**
```go
type Model struct {
    BaseModel
    Name         string `gorm:"uniqueIndex;not null"`
    Provider     string `gorm:"index"`
    Category     string // text/image/audio/video
    MaxTokens    int
    InputPrice   float64
    OutputPrice  float64
    Description  string
    Enabled      bool `gorm:"default:true"`
}
```

**API 端点:**
```go
GET /api/admin/models?page=1&pageSize=20&search=gpt&provider=openai&category=text
GET /api/admin/models/:id
POST /api/admin/models
PUT /api/admin/models/:id
DELETE /api/admin/models/:id
```

#### 4.2 前端实现

- 模型搜索框(实时搜索)
- 虚拟滚动列表(react-window)
- 分类筛选器
- 模型详情抽屉

---

### 5. 健康检测功能

#### 5.1 后端实现

**健康检查服务:**
```go
// internal/services/health_checker.go
type HealthChecker struct {
    db *gorm.DB
}

func (hc *HealthChecker) CheckProvider(providerID uint) (*HealthStatus, error)
func (hc *HealthChecker) CheckAllProviders() ([]HealthStatus, error)
func (hc *HealthChecker) SchedulePeriodicChecks(interval time.Duration)
```

**Health Status 结构:**
```go
type HealthStatus struct {
    ProviderID    uint
    Status        string // healthy/unhealthy/unknown
    ResponseTime  int64  // ms
    LastChecked   time.Time
    ErrorMessage  string
}
```

#### 5.2 API 端点

```go
GET  /api/admin/health/providers
GET  /api/admin/health/providers/:id
POST /api/admin/health/check
POST /api/admin/health/check/:id
```

#### 5.3 前端实现

- 实时健康状态仪表盘
- 健康历史图表
- 自动刷新(WebSocket 或轮询)
- 手动触发检查按钮

---

### 6. 模型一键复制功能

#### 6.1 API 端点

```go
// 复制单个模型配置
POST /api/admin/models/:id/clone
Body: {
  "newName": "gpt-4-clone",
  "targetGroupId": 2
}

// 批量复制
POST /api/admin/models/batch-clone
Body: {
  "modelIds": [1, 2, 3],
  "targetGroupId": 2,
  "namePrefix": "clone-"
}
```

#### 6.2 前端实现

- 模型卡片上的"复制"按钮
- 复制配置对话框
- 批量选择和复制
- 复制预览

---

## 第三阶段:安全功能

### 7. JWT 认证和登录功能

#### 7.1 数据库模型

**用户认证表:**
```go
type User struct {
    BaseModel
    Username     string `gorm:"uniqueIndex;not null"`
    Email        string `gorm:"uniqueIndex;not null"`
    PasswordHash string `gorm:"not null"`
    Role         string `gorm:"default:'user'"` // admin/user
    LastLoginAt  *time.Time
    IsActive     bool `gorm:"default:true"`
}

type Session struct {
    BaseModel
    UserID       uint
    Token        string `gorm:"uniqueIndex;not null"`
    ExpiresAt    time.Time
    IPAddress    string
    UserAgent    string
}
```

#### 7.2 认证流程

**所需包:**
```go
import (
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)
```

**API 端点:**
```go
POST /api/auth/register
POST /api/auth/login
POST /api/auth/logout
POST /api/auth/refresh
GET  /api/auth/me
PUT  /api/auth/password
```

**JWT Payload:**
```json
{
  "sub": "user_id",
  "username": "admin",
  "role": "admin",
  "exp": 1234567890,
  "iat": 1234567890
}
```

#### 7.3 中间件

```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // 验证 JWT
        // 设置用户信息到 context
    }
}

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查用户角色
    }
}
```

#### 7.4 前端实现

**新增页面:**
- `Login.tsx` - 登录页面
- `Register.tsx` - 注册页面(可选)

**认证状态管理:**
```typescript
// src/contexts/AuthContext.tsx
interface AuthContextType {
  user: User | null
  token: string | null
  login: (username: string, password: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
}
```

**路由保护:**
```typescript
<Route element={<ProtectedRoute />}>
  <Route path="/dashboard" element={<Dashboard />} />
  {/* 其他受保护路由 */}
</Route>
```

---

## 实现优先级建议

### Phase 1 (立即开始)
1. ✅ 修复当前显示问题
2. 🔄 数据导入导出(CSV/XLSX)
3. Provider 管理完善

### Phase 2 (核心功能)
4. 模型映射功能
5. 健康检测功能
6. 模型按名称分页

### Phase 3 (增强功能)
7. 模型一键复制
8. JWT 认证登录

---

## 技术栈补充

### 后端新增依赖
```go
go get github.com/360EntSecGroup-Skylar/excelize/v2
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

### 前端新增依赖
```json
{
  "xlsx": "^0.18.5",
  "react-window": "^1.8.10",
  "recharts": "^2.10.0",
  "@types/react-window": "^1.8.8"
}
```

---

## 数据库迁移计划

### Migration 1: 添加 Model 表
```sql
CREATE TABLE models (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    provider VARCHAR(100),
    category VARCHAR(50),
    max_tokens INTEGER,
    input_price DECIMAL(10,6),
    output_price DECIMAL(10,6),
    description TEXT,
    enabled BOOLEAN DEFAULT TRUE,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE INDEX idx_models_provider ON models(provider);
CREATE INDEX idx_models_category ON models(category);
```

### Migration 2: 扩展 Provider 表
```sql
ALTER TABLE providers ADD COLUMN base_url VARCHAR(255);
ALTER TABLE providers ADD COLUMN timeout INTEGER DEFAULT 30;
ALTER TABLE providers ADD COLUMN max_retries INTEGER DEFAULT 3;
ALTER TABLE providers ADD COLUMN health_status VARCHAR(20) DEFAULT 'unknown';
ALTER TABLE providers ADD COLUMN last_checked DATETIME;
```

### Migration 3: 添加认证表
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    last_login_at DATETIME,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token VARCHAR(500) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
```

---

## 测试计划

### 单元测试
- [ ] Export/Import 功能测试
- [ ] 模型映射逻辑测试
- [ ] JWT 生成和验证测试
- [ ] 健康检查逻辑测试

### 集成测试
- [ ] 完整的导入导出流程
- [ ] 认证和授权流程
- [ ] API 端到端测试

### 性能测试
- [ ] 大量模型加载性能
- [ ] 健康检查并发测试
- [ ] 数据导出性能测试

---

## 文档更新

需要更新的文档:
- [ ] API 文档 - 添加所有新端点
- [ ] 用户手册 - 功能使用说明
- [ ] 开发者文档 - 架构更新
- [ ] 部署文档 - 新依赖和配置

---

## 预估工作量

| 功能模块 | 后端开发 | 前端开发 | 测试 | 总计 |
|---------|---------|---------|------|------|
| 数据导入导出 | 2天 | 1天 | 1天 | 4天 |
| 模型映射 | 1天 | 1天 | 0.5天 | 2.5天 |
| Provider管理 | 2天 | 2天 | 1天 | 5天 |
| 健康检测 | 2天 | 1天 | 1天 | 4天 |
| 模型分页加载 | 1天 | 1.5天 | 0.5天 | 3天 |
| 一键复制 | 0.5天 | 0.5天 | 0.5天 | 1.5天 |
| JWT认证 | 3天 | 2天 | 1天 | 6天 |
| **总计** | **11.5天** | **9天** | **5.5天** | **26天** |

注:以上为单人全职开发的预估时间。