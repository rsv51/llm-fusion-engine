# LLM Fusion Engine

## 简介

LLM Fusion Engine 是一个高性能、可扩展的下一代大型语言模型 (LLM) 网关。它巧妙地融合了 `llmio-master` 的现代化前端与全面管理 API，以及 `OrchestrationApi-main` 强大而灵活的路由与编排核心逻辑。

本项目旨在为企业和开发者提供一个统一的、可靠的、可观测的 LLM 服务访问层，简化多模型、多提供商环境下的复杂性。

## 核心功能

- **OpenAI 兼容 API:** 提供与 OpenAI API 完全兼容的 `/v1/chat/completions` 端点，允许无缝集成现有生态。
- **高级路由与编排:**
    - **分组策略:** 以“分组”为核心，灵活配置提供商、模型和路由规则。
    - **多种负载均衡:** 支持故障转移 (failover)、轮询 (round-robin)、加权 (weighted) 和随机 (random) 等多种负载均衡策略。
    - **模型别名:** 允许将通用模型名称（如 `gpt-4-best`）映射到特定提供商的具体模型。
- **全面的管理 UI:**
    - 基于 React 和 Vite 构建的现代化单页应用 (SPA)。
    - 提供对分组、提供商、密钥等所有核心资源的可视化 CRUD 操作。
    - 包含仪表盘、实时日志等监控功能。
- **Docker 原生:**
    - 提供多阶段 `Dockerfile`，可构建轻量、自包含的生产镜像。
    - 简化部署和运维流程。

## 项目结构

```
llm-fusion-engine/
├── cmd/server/main.go      # Go 应用程序主入口
├── internal/               # 后端核心业务逻辑
│   ├── api/                # API 处理器 (v1 和 admin)
│   ├── core/               # 核心业务接口定义
│   ├── database/           # 数据库模型和初始化
│   └── services/           # 业务服务实现
├── web/                    # 前端 React 项目
│   ├── src/                # 前端源代码
│   ├── package.json
│   └── vite.config.ts
├── Dockerfile              # 多阶段 Docker 构建文件
└── README.md               # 项目文档
```

## 如何构建和运行

### 先决条件

- Docker

### 构建

在项目根目录下运行以下命令来构建 Docker 镜像：

```bash
docker build -t llm-fusion-engine .
```

### 运行

使用以下命令来运行应用程序：

```bash
docker run -p 8080:8080 -d llm-fusion-engine
```

应用程序将在 `http://localhost:8080` 上可用。

- **API 端点:** `http://localhost:8080/v1`
- **管理后台:** `http://localhost:8080/`
