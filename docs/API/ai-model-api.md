# AI 模型 API

**控制器：** `controllers/ai_model_controller.go`  
**路由前缀：** `/api/v1/ai`

---

## 概述

AI 模型管理 API 用于管理平台的 AI 模型配置，包括模型的注册、部署、版本管理和推理调用。

### 模型状态

| 状态值 | 说明 |
|--------|------|
| pending | 待处理 |
| ready | 就绪 |
| deploying | 部署中 |
| online | 已上线 |
| offline | 已下线 |
| error | 异常 |

### 模型类型

| 类型 | 说明 |
|------|------|
| llm | 大语言模型 |
| tts | 语音合成 |
| stt | 语音识别 |
| vision | 视觉模型 |
| embedding | 向量嵌入 |
| agent | Agent模型 |

### AI 提供商

| 提供商 | 说明 |
|--------|------|
| openai | OpenAI |
| anthropic | Anthropic |
| gemini | Google Gemini |
| local | 本地部署 |
| custom | 自定义 |

---

## 1. 模型列表

获取 AI 模型配置列表，支持分页和筛选。

### 请求

```
GET /api/v1/ai/models
```

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20 |
| provider | string | 否 | 提供商筛选 |
| model_type | string | 否 | 模型类型筛选 |
| status | string | 否 | 状态筛选 |
| keyword | string | 否 | 关键词搜索（名称/模型标识） |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "GPT-4",
        "model_key": "gpt-4",
        "description": "OpenAI GPT-4 大语言模型",
        "provider": "openai",
        "model_type": "llm",
        "model_size": null,
        "status": "online",
        "capabilities": "[\"chat\",\"completion\",\"function_call\"]",
        "quota_daily": 10000,
        "quota_monthly": 300000,
        "price_per_1k": 0.03,
        "created_at": "2026-03-01T10:00:00Z",
        "updated_at": "2026-03-22T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 15,
      "total_pages": 1
    }
  }
}
```

---

## 2. 创建模型

注册/上传新的 AI 模型配置。

### 请求

```
POST /api/v1/ai/models
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 模型名称 |
| description | string | 否 | 模型描述 |
| provider | string | 是 | 提供商：`openai` / `anthropic` / `gemini` / `local` / `custom` |
| model_type | string | 是 | 模型类型：`llm` / `tts` / `stt` / `vision` / `embedding` / `agent` |
| model_size | string | 否 | 模型规模参数 |
| file_path | string | 否 | 模型文件路径（本地部署时） |
| file_size | int | 否 | 文件大小（字节） |
| checksum | string | 否 | 文件校验和 |
| config | string | 否 | 额外配置（JSON格式） |
| capabilities | string | 否 | 模型能力列表（JSON数组格式） |

### 请求示例

```json
{
  "name": "GPT-4",
  "model_key": "gpt-4",
  "description": "OpenAI GPT-4 大语言模型",
  "provider": "openai",
  "model_type": "llm",
  "capabilities": "[\"chat\",\"completion\",\"function_call\"]",
  "config": "{\"max_tokens\": 4096, \"temperature\": 0.7}",
  "price_per_1k": 0.03
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "GPT-4",
    "model_key": "gpt-4",
    "description": "OpenAI GPT-4 大语言模型",
    "provider": "openai",
    "model_type": "llm",
    "status": "pending",
    "capabilities": "[\"chat\",\"completion\",\"function_call\"]",
    "created_at": "2026-03-22T10:00:00Z"
  }
}
```

### 错误码

| code | 说明 |
|------|------|
| 4009 | 模型已存在（相同 model_key） |

---

## 3. 获取模型详情

根据 ID 获取单个模型的详细信息。

### 请求

```
GET /api/v1/ai/models/:id
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 模型ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "GPT-4",
    "model_key": "gpt-4",
    "description": "OpenAI GPT-4 大语言模型",
    "provider": "openai",
    "model_type": "llm",
    "model_size": null,
    "file_path": null,
    "status": "online",
    "config": "{\"max_tokens\": 4096, \"temperature\": 0.7}",
    "capabilities": "[\"chat\",\"completion\",\"function_call\"]",
    "quota_daily": 10000,
    "quota_monthly": 300000,
    "price_per_1k": 0.03,
    "usage_today": 1500,
    "usage_month": 45000,
    "create_user_id": 1,
    "org_id": 1,
    "created_at": "2026-03-01T10:00:00Z",
    "updated_at": "2026-03-22T10:00:00Z"
  }
}
```

### 错误码

| code | 说明 |
|------|------|
| 4004 | 模型不存在 |
| 4005 | 无效的模型ID |

---

## 4. 更新模型信息

更新模型的配置和参数。

### 请求

```
PUT /api/v1/ai/models/:id
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 模型ID |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 模型名称 |
| description | string | 否 | 模型描述 |
| model_type | string | 否 | 模型类型 |
| config | string | 否 | 额外配置（JSON格式） |
| capabilities | string | 否 | 模型能力列表 |
| quota_daily | int | 否 | 每日配额 |
| quota_monthly | int | 否 | 每月配额 |
| price_per_1k | float | 否 | 每千token价格 |

### 请求示例

```json
{
  "name": "GPT-4-32K",
  "description": "GPT-4 32K 上下文版本",
  "quota_daily": 20000,
  "price_per_1k": 0.06
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "GPT-4-32K",
    "description": "GPT-4 32K 上下文版本",
    "quota_daily": 20000,
    "price_per_1k": 0.06,
    "updated_at": "2026-03-22T11:00:00Z"
  }
}
```

---

## 5. 删除模型

删除模型配置。

### 请求

```
DELETE /api/v1/ai/models/:id
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 模型ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success"
}
```

### 错误码

| code | 说明 |
|------|------|
| 4003 | 不允许删除在线模型，请先下线 |

---

## 6. 部署模型

将模型部署上线。

### 请求

```
POST /api/v1/ai/models/:id/deploy
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 模型ID |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| deployment_env | string | 否 | 部署环境：`production` / `staging`，默认 `production` |
| replica_count | int | 否 | 副本数量，默认 1 |
| resources | object | 否 | 资源限制配置 |

### 请求示例

```json
{
  "deployment_env": "production",
  "replica_count": 2,
  "resources": {
    "cpu": "2",
    "memory": "4Gi"
  }
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "model_id": 1,
    "deployment_id": "deploy_xyz789",
    "status": "deploying",
    "deployed_at": "2026-03-22T10:00:00Z"
  }
}
```

---

## 7. 获取部署历史

获取模型的部署记录。

### 请求

```
GET /api/v1/ai/models/:id/deploy-history
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 模型ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "deployment_id": "deploy_001",
        "model_id": 1,
        "deployment_env": "production",
        "status": "success",
        "replica_count": 2,
        "started_at": "2026-03-20T10:00:00Z",
        "completed_at": "2026-03-20T10:05:00Z"
      },
      {
        "deployment_id": "deploy_002",
        "model_id": 1,
        "deployment_env": "production",
        "status": "failed",
        "error_message": "资源不足",
        "started_at": "2026-03-22T10:00:00Z",
        "completed_at": "2026-03-22T10:02:00Z"
      }
    ]
  }
}
```

---

## 8. 模型推理调用

调用 AI 模型进行推理（通过 AI 推理控制器）。

### 请求

```
POST /api/v1/ai/inference
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| model_id | int | 否 | 模型ID（与 model_key 二选一） |
| model_key | string | 否 | 模型标识（与 model_id 二选一） |
| prompt | string | 是 | 输入提示 |
| messages | array | 否 | 对话消息列表（ChatML格式） |
| max_tokens | int | 否 | 最大生成token数 |
| temperature | float | 否 | 采样温度 0-2 |
| top_p | float | 否 | 核采样概率 |
| stream | bool | 否 | 是否流式返回，默认 false |

### 请求示例（非流式）

```json
{
  "model_key": "gpt-4",
  "messages": [
    {"role": "system", "content": "你是一个有帮助的助手"},
    {"role": "user", "content": "你好，请介绍一下你自己"}
  ],
  "max_tokens": 1000,
  "temperature": 0.7
}
```

### 响应示例（非流式）

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "model": "gpt-4",
    "choices": [
      {
        "message": {
          "role": "assistant",
          "content": "你好！我是..."
        },
        "finish_reason": "stop"
      }
    ],
    "usage": {
      "prompt_tokens": 50,
      "completion_tokens": 150,
      "total_tokens": 200
    },
    "id": "chatcmpl_abc123",
    "created": 1677654321
  }
}
```

---

## 9. 模型监控

获取模型的运行监控数据。

### 请求

```
GET /api/v1/ai/models/:id/monitor
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 模型ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| period | string | 否 | 时间周期：`hour` / `day` / `week`，默认 `day` |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "model_id": 1,
    "period": "day",
    "metrics": {
      "total_requests": 1500,
      "successful_requests": 1480,
      "failed_requests": 20,
      "avg_latency_ms": 250,
      "p95_latency_ms": 500,
      "p99_latency_ms": 800,
      "quota_usage_percent": 15.0
    },
    "time_series": [
      {
        "timestamp": "2026-03-22T00:00:00Z",
        "requests": 100,
        "latency_ms": 240
      }
    ]
  }
}
```

---

## 10. 模型质量评估

获取模型的质量评估数据。

### 请求

```
GET /api/v1/ai/models/:id/quality
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "model_id": 1,
    "overall_score": 4.5,
    "dimensions": {
      "accuracy": 4.6,
      "fluency": 4.8,
      "coherence": 4.5,
      "safety": 4.9
    },
    "sample_count": 1000,
    "updated_at": "2026-03-22T10:00:00Z"
  }
}
```

---

## 11. 模型分片管理

管理大型模型的分片配置。

### 模型分片列表

```
GET /api/v1/ai/model-shards/:model_id
```

### 创建模型分片

```
POST /api/v1/ai/model-shards
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| model_id | int | 是 | 模型ID |
| shard_index | int | 是 | 分片索引 |
| file_path | string | 是 | 分片文件路径 |
| file_size | int | 是 | 文件大小 |
| checksum | string | 是 | 校验和 |

### 删除模型分片

```
DELETE /api/v1/ai/model-shards/:id
```

---

## 12. 模型版本管理

管理模型的版本历史。

### 版本列表

```
GET /api/v1/ai/model-versions/:model_id
```

### 创建版本

```
POST /api/v1/ai/model-versions
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| model_id | int | 是 | 模型ID |
| version | string | 是 | 版本号，如 `1.0.0` |
| changelog | string | 否 | 变更日志 |
| file_path | string | 否 | 版本文件路径 |

### 回滚版本

```
POST /api/v1/ai/model-versions/:id/rollback
```

---

## 13. AI 训练管理

模型训练相关接口。

### 训练任务列表

```
GET /api/v1/ai/trainings
```

### 创建训练任务

```
POST /api/v1/ai/trainings
```

### 获取训练任务详情

```
GET /api/v1/ai/trainings/:id
```

### 取消训练任务

```
POST /api/v1/ai/trainings/:id/cancel
```

---

## 14. AI Sandbox 管理

沙箱环境管理接口。

### Sandbox 列表

```
GET /api/v1/ai/sandboxes
```

### 创建 Sandbox

```
POST /api/v1/ai/sandboxes
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 沙箱名称 |
| model_id | int | 否 | 关联模型ID |
| config | object | 否 | 沙箱配置 |

### Sandbox 详情

```
GET /api/v1/ai/sandboxes/:id
```

### 删除 Sandbox

```
DELETE /api/v1/ai/sandboxes/:id
```

---

## 相关数据模型

### AIModelConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 模型名称 |
| model_key | string | 模型唯一标识 |
| description | string | 模型描述 |
| provider | string | 提供商 |
| model_type | string | 模型类型 |
| model_size | string | 模型规模 |
| file_path | string | 模型文件路径 |
| file_size | int | 文件大小 |
| checksum | string | 校验和 |
| status | string | 状态 |
| config | text | 额外配置JSON |
| capabilities | text | 能力列表JSON |
| quota_daily | int | 每日配额 |
| quota_monthly | int | 每月配额 |
| price_per_1k | float | 价格 |
| create_user_id | uint | 创建用户ID |
| org_id | uint | 组织ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

---

## 错误码

| code | error_code | 说明 |
|------|-------------|------|
| 4004 | ERR_NOT_FOUND | 模型不存在 |
| 4005 | ERR_VALIDATION | 参数校验失败 |
| 4009 | ERR_CONFLICT | 模型已存在 |
| 4003 | ERR_INVALID_STATUS | 状态不允许（如删除在线模型） |
| 5001 | ERR_INTERNAL | 服务器内部错误 |
