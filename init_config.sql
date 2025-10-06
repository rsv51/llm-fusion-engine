-- 初始化 llm-fusion-engine 配置
-- 此脚本将创建必要的模型、提供商和映射关系

-- 1. 创建模型
INSERT INTO models (name, remark, max_retry, timeout, enabled, created_at, updated_at)
VALUES ('qwen3-30b-a3b', 'Qwen 3 30B A3B 模型', 3, 60, 1, datetime('now'), datetime('now'));

-- 2. 创建提供商
INSERT INTO providers (name, type, config, enabled, priority, weight, health_status, created_at, updated_at)
VALUES (
    'ClawCloudRun',
    'openai',
    '{"baseUrl":"https://lwnpnfrzdigq.us-west-1.clawcloudrun.com","apiKey":"YOUR_API_KEY_HERE"}',
    1,
    100,
    1,
    'unknown',
    datetime('now'),
    datetime('now')
);

-- 3. 创建模型-提供商映射
INSERT INTO model_provider_mappings (model_id, provider_id, provider_model, enabled, weight, created_at, updated_at)
SELECT 
    m.id,
    p.id,
    'qwen3-30b-a3b',
    1,
    1,
    datetime('now'),
    datetime('now')
FROM models m
CROSS JOIN providers p
WHERE m.name = 'qwen3-30b-a3b' AND p.name = 'ClawCloudRun';

-- 4. 创建代理密钥(请替换为您自己的密钥)
INSERT INTO proxy_keys (user_id, key, enabled, created_at, updated_at)
VALUES (1, 'sk-your-proxy-key-here', 1, datetime('now'), datetime('now'));

-- 显示配置结果
SELECT 'Models:' as info;
SELECT * FROM models;

SELECT 'Providers:' as info;
SELECT id, name, type, enabled FROM providers;

SELECT 'Mappings:' as info;
SELECT mpm.id, m.name as model_name, p.name as provider_name, mpm.provider_model
FROM model_provider_mappings mpm
JOIN models m ON m.id = mpm.model_id
JOIN providers p ON p.id = mpm.provider_id;

SELECT 'Proxy Keys:' as info;
SELECT id, key, enabled FROM proxy_keys;