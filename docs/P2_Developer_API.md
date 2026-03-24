# 开发者 API PRD

## 1. 功能概述
开发者 API 模块为第三方开发者提供 API 访问密钥管理、API 调用配额监控和使用量统计功能，支持开发者自助接入 MDM 平台能力。

## 2. 页面布局与交互

### 页面路径
`/developer` → `DeveloperConsoleView.vue`

### 开发者应用管理（Tab1）
#### 数据表格
| 列 | 说明 |
|----|------|
| 应用名称 | app_name |
| AppKey | app_key（显示前4位+***）|
| AppSecret | app_secret（脱敏）|
| 权限范围 | scopes |
| 状态 | status（active/disabled）|
| 创建时间 | created_at |
| 操作 | 查看/编辑/禁用/删除 |

#### 新建应用弹窗
| 字段 | 类型 | 说明 |
|------|------|------|
| 应用名称 | Input | 必填 |
| 应用描述 | Textarea | - |
| 权限范围 | Checkbox Group | 设备查询/设备控制/OTA升级等 |

### API 密钥管理（Tab2）
- AppKey 和 AppSecret 显示
- 「重新生成 Secret」按钮（需确认）
- 授权回调 URL 配置

### API 调用统计（Tab3）
- 统计卡片：今日调用量/本月调用量/配额使用率
- 调用趋势图（折线图）
- Top API 排行（表格）

## 3. API 契约

### 应用列表
- 路径：`GET /api/v1/developer/apps`
- 参数：`page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "app_id": "uuid",
        "app_name": "第三方宠物APP",
        "app_key": "mdm_a1b2***",
        "scopes": ["device:read", "device:control"],
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 5
  }
}
```

### 创建应用
- 路径：`POST /api/v1/developer/apps`
- 请求体：
```json
{
  "app_name": "第三方宠物APP",
  "description": "用于XXX功能",
  "scopes": ["device:read", "device:control"],
  "redirect_uri": "https://example.com/callback"
}
```

### 获取应用详情
- 路径：`GET /api/v1/developer/apps/:id`
- 响应：包含完整的 app_key 和 app_secret

### 更新应用
- 路径：`PUT /api/v1/developer/apps/:id`
- 请求体：可更新 name, scopes, redirect_uri, status

### 删除应用
- 路径：`DELETE /api/v1/developer/apps/:id`

### 重新生成 Secret
- 路径：`POST /api/v1/developer/apps/:id/rotate-secret`
- 响应：
```json
{
  "code": 0,
  "data": {
    "app_key": "mdm_a1b2***",
    "app_secret": "新Secret（仅此次显示）"
  }
}
```

### API 调用统计
- 路径：`GET /api/v1/developer/stats`
- 参数：`app_id`, `period` (today/month/custom)
- 响应：
```json
{
  "code": 0,
  "data": {
    "total_calls": 100000,
    "success_calls": 99000,
    "error_calls": 1000,
    "quota_limit": 1000000,
    "quota_used_percent": 10.0,
    "top_apis": [
      { "api": "/api/v1/devices", "calls": 50000 },
      { "api": "/api/v1/devices/:id/status", "calls": 30000 }
    ],
    "daily_trend": [
      { "date": "2024-01-01", "calls": 5000 }
    ]
  }
}
```

### API 配额管理
- 路径：`GET /api/v1/developer/quota`
- 响应：
```json
{
  "code": 0,
  "data": {
    "daily_limit": 100000,
    "daily_used": 50000,
    "monthly_limit": 1000000,
    "monthly_used": 500000,
    "rate_limit": 1000,
    "rate_window": 60
  }
}
```

### 禁用/启用应用
- 路径：`POST /api/v1/developer/apps/:id/disable`
- 路径：`POST /api/v1/developer/apps/:id/enable`

## 4. 数据模型

### DeveloperApp（开发者应用）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| app_id | string | UUID |
| app_name | string | 应用名称 |
| app_key | string | API Key |
| app_secret | string | API Secret（加密）|
| description | string | 描述 |
| scopes | JSONB | 权限范围 |
| redirect_uri | string | 回调URL |
| status | string | active/disabled |
| tenant_id | uint | 租户ID |
| created_at | datetime | 创建时间 |

### APICallLog（API调用日志）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| app_id | string | 应用ID |
| api_path | string | API路径 |
| method | string | HTTP方法 |
| status_code | int | 响应码 |
| latency_ms | int | 耗时ms |
| called_at | datetime | 调用时间 |

## 5. 验收标准
- [ ] 应用列表加载正常
- [ ] 创建应用生成正确的 AppKey/AppSecret
- [ ] 权限范围多选正确
- [ ] Secret重新生成后旧Secret立即失效
- [ ] API调用统计按时段正确统计
- [ ] 配额超限提示正确
- [ ] 禁用应用后API调用返回403
- [ ] Top API排行正确
