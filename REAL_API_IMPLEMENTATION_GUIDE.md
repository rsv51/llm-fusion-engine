# 真实API获取模型列表实现指南

## 修复概述

现在系统已经实现了从供应商真实API获取模型列表的功能,而不再依赖硬编码的模型列表。

## 实现的改进

### 1. 创建了供应商客户端接口

**文件**: [`internal/providers/interface.go`](llm-fusion-engine/internal/providers/interface.go)

定义了通用的供应商客户端接口:
```go
type ProviderClient interface {
    GetModels(ctx context.Context) ([]string, error)
    ValidateConfig() error
}
```

### 2. 实现了OpenAI兼容的客户端

**文件**: [`internal/providers/openai.go`](llm-fusion-engine/internal/providers/openai.go)

实现了调用 OpenAI API `/v1/models` 端点获取真实模型列表的功能。这个客户端也兼容任何使用 OpenAI API 格式的供应商(如您的 qwen 服务)。

### 3. 创建了供应商工厂

**文件**: [`internal/providers/factory.go`](llm-fusion-engine/internal/providers/factory.go)

根据供应商类型和配置创建相应的客户端实例。

### 4. 更新了后端API处理器

**文件**: [`internal/api/admin/provider_handler.go`](llm-fusion-engine/internal/api/admin/provider_handler.go:120-200)

`GetProviderModels` 函数现在:
1. 尝试从供应商的实际API获取模型列表
2. 如果API调用失败,回退到默认模型列表
3. 返回模型来源信息(`source: "api"` 或 `source: "default"`)

### 5. 修复了前端UI文本

**文件**: [`web/src/pages/Providers.tsx`](llm-fusion-engine/web/src/pages/Providers.tsx:555-585)

- 将"一键复制所有模型"改为"导入所有模型到配置"
- 将单个模型的"复制"按钮改为"导入"
- 更准确地描述了按钮的实际功能

## 工作流程

1. **用户点击"获取模型列表"**
   - 前端调用 `/admin/providers/:id/models`

2. **后端处理**
   - 从数据库获取供应商配置
   - 解析供应商的 `config` JSON(包含 apiKey, baseUrl 等)
   - 创建对应的客户端(如 OpenAIClient)
   - 调用供应商真实API: `GET {baseUrl}/models`
   - 解析响应并返回模型列表

3. **对于您的 qwen 供应商**
   - 如果配置正确,会调用您提供的 API 地址
   - 返回该 API 实际提供的模型(如 `qwen3-coder-plus`)
   - 不再返回硬编码的 OpenAI 模型列表

## 供应商配置要求

为了正确获取模型列表,供应商的 `config` JSON 必须包含:

```json
{
  "apiKey": "your-api-key",
  "baseUrl": "https://your-provider-api.com/v1",
  "timeout": 30,
  "maxRetries": 3
}
```

## 故障处理

如果从 API 获取模型失败:
- 系统会返回默认模型列表
- 在响应中包含警告信息和错误原因
- 前端会显示获取到的模型,用户仍然可以继续操作

## 测试建议

1. **测试真实API调用**
   ```bash
   # 确保供应商配置正确
   # 点击"获取模型列表"
   # 应该看到实际API返回的模型
   ```

2. **测试回退机制**
   ```bash
   # 将供应商的 apiKey 设置为无效值
   # 点击"获取模型列表"
   # 应该看到默认模型列表和警告信息
   ```

3. **测试模型导入**
   ```bash
   # 获取模型列表后
   # 点击"导入所有模型到配置"
   # 检查模型管理页面,应该看到新导入的模型
   ```

## 下一步优化

1. **添加缓存机制**
   - 缓存模型列表24小时
   - 减少对供应商API的调用频率

2. **支持更多供应商类型**
   - 为 Anthropic 和 Gemini 实现专门的客户端
   - 支持其他OpenAI兼容的供应商

3. **添加手动刷新功能**
   - 在UI上添加"刷新模型列表"按钮
   - 强制重新从API获取最新模型

4. **改进错误提示**
   - 在UI上显示更详细的错误信息
   - 提供故障排查建议

## 常见问题

**Q: 为什么我的供应商仍然显示默认模型列表?**

A: 检查以下几点:
- 供应商的 `config` 中是否包含正确的 `apiKey` 和 `baseUrl`
- API地址是否可访问
- API是否返回标准的 OpenAI 格式响应
- 查看后端日志中的错误信息

**Q: 如何知道模型是从API获取还是默认列表?**

A: 查看后端响应中的 `source` 字段:
- `"source": "api"` = 从真实API获取
- `"source": "default"` = 使用默认列表

**Q: OpenAI格式的响应应该是什么样的?**

A: 标准格式:
```json
{
  "data": [
    {
      "id": "model-name",
      "object": "model",
      "created": 1234567890,
      "owned_by": "provider"
    }
  ],
  "object": "list"
}
```

## 相关文件

- [`internal/providers/`](llm-fusion-engine/internal/providers/) - 供应商客户端实现
- [`internal/api/admin/provider_handler.go`](llm-fusion-engine/internal/api/admin/provider_handler.go) - API处理器
- [`web/src/pages/Providers.tsx`](llm-fusion-engine/web/src/pages/Providers.tsx) - 前端页面
- [`web/src/services/models.ts`](llm-fusion-engine/web/src/services/models.ts) - 前端服务