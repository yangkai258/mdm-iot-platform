# 设备配对与注册

**版本：** V1.0  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

设备配对与注册模块负责管理 MiniClaw 设备与 OpenClaw AI 系统的绑定关系。设备首次开机时需要完成配对流程，之后才能正常使用 AI 能力和远程控制。

**业务目标：**
- 实现设备安全配对流程
- 管理设备与用户/AI 的绑定关系
- 支持设备解绑和重新配对

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 |
|------|------|--------|----------|
| 配对码生成 | 生成设备配对码 | P0 | 人工 |
| 配对码验证 | 验证设备配对码有效性 | P0 | 自动 |
| 设备绑定 | 设备绑定到用户账户 | P0 | 自动 |
| AI 授权 | 设备绑定 AI 服务授权 | P0 | 自动 |
| 解绑设备 | 解绑设备与用户关系 | P0 | 人工 |
| 重新配对 | 设备恢复出厂设置后重新配对 | P1 | 人工 |
| 配对记录 | 查询设备配对历史 | P1 | 人工 |

---

## 3. 数据模型

### 3.1 配对记录表 (pairing_records)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| pairing_code | string | 配对码 | unique, not null, 6位数字 |
| device_id | string | 设备ID | FK to devices, nullable |
| user_id | uint | 用户ID | FK to sys_users, nullable |
| status | int | 状态 | 1=待激活 2=已配对 3=已解绑 |
| expires_at | datetime | 过期时间 | not null |
| paired_at | datetime | 配对时间 | nullable |
| unbound_at | datetime | 解绑时间 | nullable |
| unbound_reason | string | 解绑原因 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 设备 AI 授权表 (device_openclaw_binding)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, unique, not null |
| user_id | uint | 所属用户 | FK to sys_users, not null |
| openclaw_version_id | uint | AI 版本 | FK to openclaw_versions, not null |
| ai_model_type | string | AI 模型类型 | not null |
| auth_status | int | 授权状态 | 1=待授权 2=已授权 3=已取消 |
| auth_token | string | 授权Token | nullable |
| expires_at | datetime | Token过期时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

---

## 4. 接口定义

### 4.1 生成配对码

```
POST /api/v1/devices/pairing/code
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "pairing_code": "123456",
    "expires_at": "2026-03-22T12:00:00Z",
    "expires_in_seconds": 300
  }
}
```

### 4.2 设备配对

```
POST /api/v1/devices/pairing/verify
```

**请求示例：**
```json
{
  "pairing_code": "123456",
  "device_id": "pet-001",
  "device_info": {
    "firmware_version": "v1.3.0",
    "hardware_version": "v1.3",
    "mac_address": "AA:BB:CC:DD:EE:FF"
  }
}
```

### 4.3 解绑设备

```
POST /api/v1/devices/:device_id/unbind
```

### 4.4 查询配对记录

```
GET /api/v1/devices/pairing/history
```

---

## 5. 前端页面

| 页面 | 路由 | 功能 |
|------|------|------|
| 配对码管理 | /devices/pairing | 生成/管理配对码 |
| 设备绑定列表 | /devices/bindings | 查看绑定设备 |
| 配对历史 | /devices/pairing/history | 配对/解绑记录 |
