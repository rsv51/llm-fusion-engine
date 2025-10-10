-- 初始化 llm-fusion-engine 数据库结构
-- 此脚本基于 Go 模型定义创建所有必要的表

-- 启用外键约束
PRAGMA foreign_keys = ON;

-- 1. 创建 users 表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 2. 创建 groups 表
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,
    models TEXT,
    model_aliases TEXT,
    load_balance_policy TEXT DEFAULT 'failover',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 3. 创建 providers 表
CREATE TABLE IF NOT EXISTS providers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,
    config TEXT,
    console VARCHAR(255),
    enabled BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,
    weight INTEGER DEFAULT 100,
    timeout INTEGER DEFAULT 300,
    health_status VARCHAR(50),
    latency INTEGER,
    last_status_code INTEGER,
    last_checked DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 4. 创建 api_keys 表
CREATE TABLE IF NOT EXISTS api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    provider_id INTEGER,
    key TEXT UNIQUE NOT NULL,
    last_used DATETIME,
    is_healthy BOOLEAN DEFAULT TRUE,
    rpm_limit INTEGER,
    tpm_limit INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (provider_id) REFERENCES providers(id)
);

-- 5. 创建 models 表
CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    remark TEXT,
    max_retry INTEGER DEFAULT 3,
    timeout INTEGER DEFAULT 30,
    enabled BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 6. 创建 model_provider_mappings 表
CREATE TABLE IF NOT EXISTS model_provider_mappings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    provider_id INTEGER NOT NULL,
    provider_model TEXT NOT NULL,
    tool_call BOOLEAN,
    structured_output BOOLEAN,
    image BOOLEAN,
    weight INTEGER DEFAULT 1,
    enabled BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (model_id) REFERENCES models(id),
    FOREIGN KEY (provider_id) REFERENCES providers(id)
);

-- 7. 创建 proxy_keys 表
CREATE TABLE IF NOT EXISTS proxy_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    key TEXT UNIQUE NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    allowed_groups TEXT,
    group_balance_policy TEXT DEFAULT 'failover',
    group_weights TEXT,
    rpm_limit INTEGER,
    tpm_limit INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 8. 创建 logs 表
CREATE TABLE IF NOT EXISTS logs (
    id TEXT PRIMARY KEY,
    proxy_key TEXT,
    model TEXT,
    provider TEXT,
    request_url TEXT,
    request_body TEXT,
    response_body TEXT,
    response_status INTEGER,
    is_success BOOLEAN,
    latency INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    prompt_tokens INTEGER,
    completion_tokens INTEGER,
    total_tokens INTEGER
);

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);
CREATE INDEX IF NOT EXISTS idx_providers_name ON providers(name);
CREATE INDEX IF NOT EXISTS idx_providers_type ON providers(type);
CREATE INDEX IF NOT EXISTS idx_providers_enabled ON providers(enabled);
CREATE INDEX IF NOT EXISTS idx_api_keys_provider_id ON api_keys(provider_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);
CREATE INDEX IF NOT EXISTS idx_models_name ON models(name);
CREATE INDEX IF NOT EXISTS idx_model_provider_mappings_model_provider ON model_provider_mappings(model_id, provider_id);
CREATE INDEX IF NOT EXISTS idx_proxy_keys_key ON proxy_keys(key);
CREATE INDEX IF NOT EXISTS idx_proxy_keys_user_id ON proxy_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_proxy_key ON logs(proxy_key);
CREATE INDEX IF NOT EXISTS idx_logs_model ON logs(model);
CREATE INDEX IF NOT EXISTS idx_logs_provider ON logs(provider);
CREATE INDEX IF NOT EXISTS idx_logs_response_status ON logs(response_status);
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);

-- 插入默认管理员用户（密码: admin）
INSERT OR IGNORE INTO users (username, password, is_admin) 
VALUES ('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', TRUE);

-- 显示创建的表结构
SELECT 'Users table structure:' as info;
PRAGMA table_info(users);

SELECT 'Groups table structure:' as info;
PRAGMA table_info(groups);

SELECT 'Providers table structure:' as info;
PRAGMA table_info(providers);

SELECT 'Models table structure:' as info;
PRAGMA table_info(models);

SELECT 'Model Provider Mappings table structure:' as info;
PRAGMA table_info(model_provider_mappings);

SELECT 'Proxy Keys table structure:' as info;
PRAGMA table_info(proxy_keys);

SELECT 'API Keys table structure:' as info;
PRAGMA table_info(api_keys);

SELECT 'Logs table structure:' as info;
PRAGMA table_info(logs);