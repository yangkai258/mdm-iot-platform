# MDM 平台 API 文档

**版本：** V1.0  
**更新日期：** 2026-03-22  
**Base URL：** `/api/v1`  
**协议：** REST + MQTT

---

## 概述

本文档涵盖 MDM 智能设备管理平台的全部 REST API，按功能模块划分为：

| 模块 | 文档 | 说明 |
|------|------|------|
| 设备管理 | [device-api.md](./device-api.md) | 设备注册、绑定、心跳、OTA |
| 会员管理 | [member-api.md](./member-api.md) | 会员、卡券、积分、店铺 |
| 告警通知 | [alert-notification-api.md](./alert-notification-api.md) | 告警规则、告警记录、推送通知 |
| 宠物管理 | [pet-api.md](./pet-api.md) | 宠物状态、消息、动作 |
| AI 模型 | [ai-model-api.md](./ai-model-api.md) | 模型注册、部署、推理 |

---

## 通用规范

### 认证方式

除登录接口外，所有接口需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

- Token 有效期：**24小时**
- JWT 载荷包含：`user_id`、`username`、`role_id`

### 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0=成功，非0=失败 |
| message | string | 提示信息 |
| data | object | 返回数据（成功时） |

### 分页响应

列表接口统一返回分页结构：

```json
{
  "code": 0,
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

### 通用错误码

| code | 说明 |
|------|------|
| 0 | 成功 |
| 4005 | 参数校验失败 |
| 4002 | 资源不存在 |
| 4003 | 业务状态不允许 |
| 4009 | 资源冲突 |
| 5001 | 服务器内部错误 |

---

## 目录结构

```
docs/API/
├── README.md                      # 本文档
├── device-api.md                  # 设备管理 API
├── member-api.md                  # 会员管理 API
├── alert-notification-api.md      # 告警通知 API
├── pet-api.md                     # 宠物管理 API
└── ai-model-api.md                # AI 模型 API
```
