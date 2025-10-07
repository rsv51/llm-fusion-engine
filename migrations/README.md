# 数据库迁移脚本使用指南

## 概述

本目录包含llm-fusion-engine项目的数据库迁移脚本。这些脚本用于更新现有数据库结构以支持新功能。

## 迁移脚本列表

### 1. add_health_status_fields.sql

**目的**: 为providers表添加健康状态监控相关字段

**添加字段**:
- `health_status`: 健康状态 (VARCHAR(50), 默认'unknown')
- `latency`: 延迟时间，单位毫秒 (INTEGER, 可为NULL)
- `last_checked`: 最后检查时间 (DATETIME, 可为NULL)

**使用方法**:

#### SQLite数据库:
```bash
# 进入llm-fusion-engine目录
cd llm-fusion-engine

# 应用迁移脚本
sqlite3 fusion.db < migrations/add_health_status_fields.sql
```

#### 其他数据库:
根据你使用的数据库类型,使用相应的命令行工具执行SQL脚本。

## 验证迁移

执行迁移后,可以验证字段是否已添加:

```sql
-- 查看providers表结构
PRAGMA table_info(providers);

-- 或使用
.schema providers
```

你应该能看到新增的三个字段: health_status, latency, last_checked

## 注意事项

1. **备份数据库**: 在执行任何迁移之前,建议先备份数据库
   ```bash
   cp fusion.db fusion.db.backup
   ```

2. **幂等性**: 如果表中已存在这些字段,迁移脚本将会失败,这是正常的保护机制

3. **数据完整性**: 迁移脚本会将所有现有记录的health_status设置为'unknown'

4. **应用顺序**: 按文件名顺序应用迁移脚本

## 故障排除

### 字段已存在错误
如果看到"duplicate column name"错误,说明字段已经存在,无需再次执行迁移。

### 权限错误
确保你有写入数据库文件的权限。

### 数据库锁定
如果数据库正在被其他进程使用,请先停止应用程序再执行迁移。

## 回滚

如果需要回滚此迁移:

```sql
-- 删除添加的字段
ALTER TABLE providers DROP COLUMN health_status;
ALTER TABLE providers DROP COLUMN latency;
ALTER TABLE providers DROP COLUMN last_checked;
```

**注意**: SQLite对DROP COLUMN的支持有限,可能需要使用表重建方式。