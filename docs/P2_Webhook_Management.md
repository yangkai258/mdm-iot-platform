# Webhook 管理 PRD

## 1. 功能概述
Webhook 管理模块允许管理员配置系统事件通知的 HTTP 回调，当指定事件（如订阅创建、支付成功、设备告警）发生时，自动向预设 URL 发送 HTTP POST 请求。

## 2. 页面布局与交互

### 页面路径
`/webhooks` → `WebhookList.vue` + `WebhookLogs.vue`

### Tab1：Webhook 列表
#### 搜索表单
| 字段 | 类型 | 说明 |
|------|------|------|
| Webhook名称 | Input | 模糊搜索 |
| 事件类型 | Multi-Select | subscription.created / payment.success 等 |
| 状态 | Select | active / inactive |

#### 数据表格
| 列 | 说明 |
|----|------|
| Webhook名称 | name |
| URL地址 | url（脱敏显示）|
| 事件类型 | event_types（Tag列表）|
| 状态 | status（active/inactive）|
| 创建时间 | created_at |
| 操作 | 编辑/删除/测试 |

### 新建/编辑弹窗
| 字段 | 类型 | 说明 |
|------|------|------|
| Webhook名称 | Input | 必填 |
| 回调URL | Input | 必填（URL格式校验）|
| 事件类型 | Multi-Select | 必选 |
| 签名密钥 | Input | 可选（用于HMAC签名）|
| 请求头 | Key-Value | 可选自定义请求头 |
| 重试次数 | Input | 默认3次 |
| 启用状态 | Switch | 默认开启 |

### Tab2：Webhook 日志
#### 筛选条件
| 字段 | 类型 | 说明 |
|------|------|------|
| Webhook | Select | 选择webhook |
| 事件类型 | Select | - |
| 状态 | Select | pending/success/failed |
| 时间范围 | DateRange | - |

#### 数据表格
| 列 | 说明 |
|----|------|
| 事件ID | event_id |
| Webhook名称 | webhook_name |
| 事件类型 | event_type |
| 状态 | status（成功绿/失败红/进行中黄）|
| 尝试次数 | attempts |
| 响应码 | response_code |
| 创建时间 | created_at |
| 操作 | 查看详情/重试 |

### 详情抽屉
- 请求信息：URL、请求头、请求体
- 响应信息：状态码、响应体、耗时
- 错误信息（如有）

## 3. API 契约

### Webhook列表
- 路径：`GET /api/v1/webhooks`
- 参数：`page`, `page_size`, `status`, `event_type`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "webhook_id": "uuid",
        "name": "支付回调",
        "url": "https://example.com/webhook",
        "event_types": ["subscription.created", "payment.success"],
        "status": "active",
        "retry_count": 3,
        "headers": { "X-Custom": "value" },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### 创建Webhook
- 路径：`POST /api/v1/webhooks`
- 请求体：
```json
{
  "name": "支付回调",
  "url": "https://example.com/webhook",
  "event_types": ["subscription.created", "payment.success"],
  "secret": "可选签名密钥",
  "headers": { "X-Custom": "value" },
  "retry_count": 3,
  "status": "active"
}
```

### 更新Webhook
- 路径：`PUT /api/v1/webhooks/:id`
- 请求体：同上

### 删除Webhook
- 路径：`DELETE /api/v1/webhooks/:id`

### 测试Webhook
- 路径：`POST /api/v1/webhooks/:id/test`
- 请求体：可选指定测试 payload
- 响应：
```json
{
  "code": 0,
  "data": {
    "success": true,
    "response_code": 200,
    "response_body": "OK",
    "duration_ms": 150
  }
}
```

### Webhook日志列表
- 路径：`GET /api/v1/webhooks/:id/logs`
- 参数：`page`, `page_size`, `status`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "event_id": "uuid",
        "event_type": "subscription.created",
        "status": "success",
        "attempts": 1,
        "response_code": 200,
        "last_error": null,
        "delivered_at": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

### 重试Webhook事件
- 路径：`POST /api/v1/webhooks/events/:event_id/retry`
- 响应：`{ "code": 0, "message": "重试已加入队列" }`

### 获取事件详情
- 路径：`GET /api/v1/webhooks/events/:event_id`
- 响应：
```json
{
  "code": 0,
  "data": {
    "event_id": "uuid",
    "event_type": "subscription.created",
    "payload": { "subscription_id": 1, "user_id": 100 },
    "status": "success",
    "attempts": 1,
    "response_code": 200,
    "response_body": "OK",
    "delivered_at": "2024-01-01T00:00:00Z"
  }
}
```

## 4. 数据模型

### Webhook（Webhook配置）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| webhook_id | string | UUID |
| name | string | 名称 |
| url | string | 回调URL |
| secret | string | 签名密钥（加密）|
| event_types | JSONB | 事件类型数组 |
| status | string | active/inactive |
| tenant_id | uint | 租户ID |
| headers | JSONB | 自定义请求头 |
| retry_count | int | 重试次数 |
| created_at | datetime | 创建时间 |

### WebhookEvent（Webhook事件）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| event_id | string | UUID |
| webhook_id | string | 关联webhook ID |
| event_type | string | 事件类型 |
| payload | JSONB | 请求内容 |
| status | string | pending/success/failed |
| attempts | int | 尝试次数 |
| max_attempts | int | 最大尝试 |
| last_error | string | 最后错误 |
| response_code | int | HTTP响应码 |
| response_body | string | 响应内容 |
| delivered_at | datetime | 成功时间 |
| next_retry_at | datetime | 下次重试时间 |

## 5. 验收标准
- [ ] Webhook列表分页加载正常
- [ ] 新建Webhook成功（URL格式校验）
- [ ] 事件类型多选正常
- [ ] 编辑/删除Webhook成功
- [ ] 测试Webhook发送成功并显示响应
- [ ] 日志列表按状态/类型筛选有效
- [ ] 查看事件详情显示完整请求/响应
- [ ] 重试功能正常
- [ ] HMAC签名正确生成
