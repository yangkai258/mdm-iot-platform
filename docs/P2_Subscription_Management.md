# 订阅管理 PRD

## 1. 功能概述
订阅管理模块为用户提供宠物智能设备的订阅服务管理，包括订阅计划查看、续费管理、自动续费设置和订阅状态监控。

## 2. 页面布局与交互

### 页面路径
`/subscription` → `SubscriptionList.vue`

### 搜索表单
| 字段 | 类型 | 说明 |
|------|------|------|
| 用户ID | Input | 精确搜索 |
| 订阅状态 | Select | active / expired / cancelled / suspended |
| 计划名称 | Input | 模糊搜索 |
| 时间范围 | DateRange | 创建时间范围 |

### 数据表格
| 列 | 说明 |
|----|------|
| 订阅ID | id |
| 用户ID | user_id |
| 计划名称 | plan_name |
| 计划类型 | plan_type（monthly/yearly）|
| 价格 | price |
| 有效期 | start_date ~ end_date |
| 状态 | status 标签 |
| 自动续费 | auto_renew（开关）|
| 操作 | 详情/续费/取消 |

### 订阅详情面板
- 用户信息：用户ID、昵称、联系方式
- 订阅计划：计划名称、类型、价格
- 有效期：开始日期、结束日期、剩余天数
- 自动续费：开启/关闭
- 续费历史：最近5条续费记录
- 操作：启用自动续费/取消订阅/恢复订阅

### 按钮
- 「新建订阅」（管理员手动创建）
- 「续费」操作按钮
- 「取消自动续费」操作按钮
- 「恢复订阅」操作按钮（针对暂停的订阅）

## 3. API 契约

### 订阅列表
- 路径：`GET /api/v1/subscriptions`
- 实际路径（subscription_controller.go）：`GET /api/v1/subscriptions`
- 参数：`user_id`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "user_id": 100,
      "plan_name": "年度高级版",
      "plan_type": "yearly",
      "price": 299.00,
      "duration": 365,
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2025-01-01T00:00:00Z",
      "status": "active",
      "auto_renew": true,
      "renew_count": 1,
      "last_renew_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 创建订阅
- 路径：`POST /api/v1/subscriptions`
- 请求体：
```json
{
  "user_id": 100,
  "plan_name": "年度高级版",
  "plan_type": "yearly",
  "price": 299.00,
  "duration": 365
}
```

### 获取订阅详情
- 路径：`GET /api/v1/subscriptions/:id`
- 响应：同上单条数据

### 更新订阅
- 路径：`PUT /api/v1/subscriptions/:id`
- 请求体：可更新 plan_name, price, auto_renew 等

### 取消订阅
- 路径：`DELETE /api/v1/subscriptions/:id`
- 响应：`{ "code": 0, "message": "订阅已取消" }`

### 自动续费
- 路径：`POST /api/v1/subscriptions/:id/renew`
- 响应：
```json
{
  "code": 0,
  "message": "续费成功",
  "data": {
    "subscription_id": 1,
    "new_end_date": "2026-01-01T00:00:00Z"
  }
}
```

### 获取续费状态
- 路径：`GET /api/v1/subscriptions/:id/renewal-status`
- 响应：
```json
{
  "code": 0,
  "data": {
    "subscription_id": 1,
    "auto_renew": true,
    "end_date": "2025-01-01T00:00:00Z",
    "days_until_expiry": 30,
    "renewal_status": "active",
    "next_renewal_amount": 299.00
  }
}
```

### 取消自动续费
- 路径：`POST /api/v1/subscriptions/:id/cancel-renewal`
- 响应：`{ "code": 0, "message": "已取消自动续费" }`

### 恢复订阅
- 路径：`POST /api/v1/subscriptions/:id/resume`
- 前提：订阅状态为 suspended
- 响应：`{ "code": 0, "message": "订阅已恢复", "data": {...} }`

### 支付回调
- 路径：`POST /api/v1/subscriptions/webhook/payment`
- 请求体：
```json
{
  "order_no": "ORD-123456",
  "status": "success"
}
```

### 续费日志
- 路径：`GET /api/v1/subscriptions/:id/renewal-logs`
- 参数：`page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "subscription_id": 1,
        "renewed_at": "2024-01-01T00:00:00Z",
        "amount": 299.00,
        "status": "success"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

## 4. 数据模型

### Subscription（订阅表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID |
| plan_name | string | 计划名称 |
| plan_type | string | monthly/yearly |
| price | float64 | 价格 |
| duration | int | 天数 |
| start_date | datetime | 开始日期 |
| end_date | datetime | 结束日期 |
| status | string | active/expired/cancelled/suspended |
| auto_renew | bool | 自动续费 |
| renew_count | int | 续费次数 |
| last_renew_at | datetime | 最后续费时间 |
| retry_count | int | 重试次数 |
| renew_fail_reason | string | 续费失败原因 |
| created_at | datetime | 创建时间 |

### SubscriptionRenewalLog（续费日志）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| subscription_id | uint | 订阅ID |
| renewed_at | datetime | 续费时间 |
| amount | float64 | 续费金额 |
| status | string | success/failed |

## 5. 验收标准
- [ ] 订阅列表分页加载正常
- [ ] 按用户/状态/计划名筛选有效
- [ ] 创建订阅成功
- [ ] 续费操作使 end_date 延长
- [ ] 取消自动续费成功
- [ ] 取消订阅后 status 变为 cancelled
- [ ] 恢复订阅仅对 suspended 状态有效
- [ ] 支付回调正确处理 success/failed
- [ ] 续费日志完整记录
