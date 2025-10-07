-- 为 providers 表添加健康状态字段
-- 如果数据库中已有这些字段,此脚本将会失败,这是正常的

-- 添加健康状态字段
ALTER TABLE providers ADD COLUMN health_status VARCHAR(50) DEFAULT 'unknown';

-- 添加延迟字段(毫秒)
ALTER TABLE providers ADD COLUMN latency INTEGER DEFAULT NULL;

-- 添加HTTP状态码字段
ALTER TABLE providers ADD COLUMN last_status_code INTEGER DEFAULT NULL;

-- 添加最后检查时间字段
ALTER TABLE providers ADD COLUMN last_checked DATETIME DEFAULT NULL;

-- 更新现有记录为未知状态
UPDATE providers SET health_status = 'unknown' WHERE health_status IS NULL OR health_status = '';