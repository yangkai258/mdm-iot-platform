# MDM 设备管理平台 API 接口文档

**版本：** V1.0
**更新日期：** 2026-03-20
**Base URL：** `/api/v1`
**状态：** 与后端代码一致 ✅

---

## 1. 通用说明

### 1.1 认证方式

除登录接口外，所有接口需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

- Token 有效期：**24小时**
- Token 在 `middleware/jwt.go` 中生成和校验
- JWT 载荷包含字段：`user_id`、`username`、`role_id`

### 1.2 CORS 跨域

后端已配置全局 CORS 中间件（`main.go`）：

| 响应头 | 值 |
|--------|-----|
| `Access-Control-Allow-Origin` | `*` |
| `Access-Control-Allow-Methods` | `GET, POST, PUT, DELETE, OPTIONS` |
| `Access-Control-Allow-Headers` | `Content-Type, Authorization` |
| 预检请求（OPTIONS） | 返回 `204 No Content`，不进入业务路由 |

> ⚠️ 注意：`Authorization` 头已在 CORS 白名单中，前端无需额外配置。

### 1.3 请求格式

- Content-Type: `application/json`
- 字符编码: `UTF-8`

### 1.4 响应格式

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

---

## 2. 认证接口

### 2.1 用户登录

**接口描述：** 用户登录，验证用户名密码后返回 JWT Token

**请求方法：** POST
**请求路径：** `/api/v1/auth/login`
**是否需要认证：** 否

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码（明文） |

**请求示例：**

```json
{
  "username": "admin",
  "password": "password123"
}
```

**响应示例（成功）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "email": "admin@example.com",
      "role_id": 1
    }
  }
}
```

**响应示例（失败 - 密码错误）：**

```json
{
  "code": 401,
  "message": "用户名或密码错误"
}
```

**响应示例（失败 - 账号禁用）：**

```json
{
  "code": 403,
  "message": "账号已被禁用"
}
```

---

### 2.2 获取当前用户信息

**请求方法：** GET
**请求路径：** `/api/v1/auth/userinfo`
**是否需要认证：** 是

**响应示例：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "管理员",
    "email": "admin@example.com",
    "phone": "13800138000",
    "role_id": 1,
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

---

### 2.3 用户登出

**请求方法：** POST
**请求路径：** `/api/v1/auth/logout`
**是否需要认证：** 是

**响应示例：**

```json
{
  "code": 0,
  "message": "success"
}
```

---

## 3. 设备管理接口

### 3.1 设备注册

**接口描述：** 设备首次联网时注册 MAC 地址和硬件信息，返回设备 ID

**请求方法：** POST
**请求路径：** `/api/v1/devices/register`
**是否需要认证：** 否（设备端直连）

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| mac_address | string | 是 | MAC 地址，格式：`AA:BB:CC:DD:EE:FF` |
| sn_code | string | 是 | 设备序列号 |
| hardware_model | string | 是 | 硬件型号 |
| firmware_version | string | 是 | 固件版本 |
| device_id | string | 否 | 已有设备 ID（更新时传入） |

**请求示例：**

```json
{
  "mac_address": "AA:BB:CC:DD:EE:FF",
  "sn_code": "SN2026032000001",
  "hardware_model": "M5Stack-Core2",
  "firmware_version": "v1.2.3"
}
```

**响应示例（成功）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "DEV001",
    "lifecycle_status": 1,
    "created_at": "2026-03-20T08:00:00Z"
  }
}
```

**lifecycle_status 取值：**

| 值 | 说明 |
|----|------|
| 1 | 待激活 |
| 2 | 服役中 |
| 3 | 离线 |
| 4 | 报废 |

**错误响应（无效 MAC 格式）：**

```json
{
  "code": 4005,
  "message": "无效的MAC地址格式",
  "error_code": "ERR_VALIDATION"
}
```

---

### 3.2 设备列表查询

**接口描述：** 分页查询设备列表，支持按状态、型号、关键字筛选，在线状态从 Redis 设备影子获取

**请求方法：** GET
**请求路径：** `/api/v1/devices`
**是否需要认证：** 是

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20，最大 100 |
| lifecycle_status | int | 否 | 生命周期状态：1-待激活 2-服役中 3-离线 4-报废 |
| hardware_model | string | 否 | 硬件型号（精确匹配） |
| status | string | 否 | 在线状态筛选：`online` / `offline`（内存过滤） |
| search | string | 否 | 搜索关键字（模糊匹配 device_id 或 sn_code） |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "device_id": "DEV001",
        "mac_address": "AA:BB:CC:DD:EE:FF",
        "sn_code": "SN2026032000001",
        "hardware_model": "M5Stack-Core2",
        "firmware_version": "v1.2.3",
        "bind_user_id": "user_001",
        "lifecycle_status": 2,
        "created_at": "2026-03-20T08:00:00Z",
        "updated_at": "2026-03-20T10:00:00Z",
        "is_online": true,
        "battery_level": 85
      },
      {
        "device_id": "DEV002",
        "mac_address": "11:22:33:44:55:66",
        "sn_code": "SN2026032000002",
        "hardware_model": "M5Stack-Core2",
        "firmware_version": "v1.2.2",
        "bind_user_id": null,
        "lifecycle_status": 1,
        "created_at": "2026-03-19T08:00:00Z",
        "updated_at": "2026-03-19T08:00:00Z",
        "is_online": false,
        "battery_level": 0
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

> **字段说明：**
> - `is_online` 和 `battery_level` 从 Redis 设备影子实时读取
> - `battery_level` 未获取到时默认为 0
> - `status=online/offline` 为内存过滤，在分页前生效

---

### 3.3 扫码绑定设备

**接口描述：** 用户扫码后将设备绑定到自己账号

**请求方法：** POST
**请求路径：** `/api/v1/devices/bind/:sn_code`
**是否需要认证：** 是

**路径参数：** `sn_code` — 设备序列号

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| bind_user_id | string | 是 | 绑定用户 ID |

**请求示例：**

```json
{
  "bind_user_id": "user_001"
}
```

**响应示例（成功）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "DEV001",
    "bind_user_id": "user_001",
    "lifecycle_status": 2,
    "message": "绑定成功"
  }
}
```

**错误响应（设备不存在）：**

```json
{
  "code": 4002,
  "message": "设备不存在",
  "error_code": "ERR_DEVICE_002"
}
```

**错误响应（设备状态非法）：**

```json
{
  "code": 4003,
  "message": "非法设备状态，无法绑定",
  "error_code": "ERR_DEVICE_003"
}
```

---

### 3.4 解绑设备

**请求方法：** POST
**请求路径：** `/api/v1/devices/unbind/:sn_code`
**是否需要认证：** 是

**响应示例：**

```json
{
  "code": 0,
  "message": "success"
}
```

---

### 3.5 获取设备详情

**请求方法：** GET
**请求路径：** `/api/v1/devices/:device_id`
**是否需要认证：** 是

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "DEV001",
    "mac_address": "AA:BB:CC:DD:EE:FF",
    "sn_code": "SN2026032000001",
    "hardware_model": "M5Stack-Core2",
    "firmware_version": "v1.2.3",
    "bind_user_id": "user_001",
    "lifecycle_status": 2,
    "shadow": {
      "is_online": true,
      "battery_level": 85,
      "last_heartbeat": "2026-03-20T10:00:00Z"
    },
    "pet_profile": {
      "name": "小白",
      "species": "cat",
      "age": 2
    },
    "created_at": "2026-03-20T08:00:00Z",
    "updated_at": "2026-03-20T10:00:00Z"
  }
}
```

---

### 3.6 更新设备信息

**请求方法：** PUT
**请求路径：** `/api/v1/devices/:device_id`
**是否需要认证：** 是

---

### 3.7 删除设备

**请求方法：** DELETE
**请求路径：** `/api/v1/devices/:device_id`
**是否需要认证：** 是

---

### 3.8 更新设备状态

**请求方法：** PUT
**请求路径：** `/api/v1/devices/:device_id/status`
**是否需要认证：** 是

---

### 3.9 获取宠物配置

**请求方法：** GET
**请求路径：** `/api/v1/devices/:device_id/profile`
**是否需要认证：** 是

---

### 3.10 更新宠物配置

**请求方法：** PUT
**请求路径：** `/api/v1/devices/:device_id/profile`
**是否需要认证：** 是

---

### 3.11 发送设备指令

**请求方法：** POST
**请求路径：** `/api/v1/devices/:device_id/commands`
**是否需要认证：** 是

---

### 3.12 指令历史

**请求方法：** GET
**请求路径：** `/api/v1/devices/:device_id/commands`
**是否需要认证：** 是

---

## 4. OTA 固件升级接口

### 4.1 创建升级包

**请求方法：** POST
**请求路径：** `/api/v1/ota/packages`
**是否需要认证：** 是

---

### 4.2 升级包列表

**请求方法：** GET
**请求路径：** `/api/v1/ota/packages`
**是否需要认证：** 是

---

### 4.3 创建升级任务

**请求方法：** POST
**请求路径：** `/api/v1/ota/deployments`
**是否需要认证：** 是

---

### 4.4 设备检查 OTA

**请求方法：** GET
**请求路径：** `/api/v1/ota/devices/:device_id/check`
**是否需要认证：** 否（设备端直连）

---

## 5. 组织管理接口

### 5.1 公司管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/org/companies` | 公司列表 |
| POST | `/api/v1/org/companies` | 创建公司 |
| PUT | `/api/v1/org/companies/:id` | 更新公司 |
| DELETE | `/api/v1/org/companies/:id` | 删除公司 |

### 5.2 部门管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/org/departments` | 部门列表 |
| GET | `/api/v1/org/departments/tree` | 部门树形结构 |
| POST | `/api/v1/org/departments` | 创建部门 |
| PUT | `/api/v1/org/departments/:id` | 更新部门 |
| DELETE | `/api/v1/org/departments/:id` | 删除部门 |

### 5.3 岗位管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/org/positions` | 岗位列表 |
| POST | `/api/v1/org/positions` | 创建岗位 |
| PUT | `/api/v1/org/positions/:id` | 更新岗位 |
| DELETE | `/api/v1/org/positions/:id` | 删除岗位 |

### 5.4 员工管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/org/employees` | 员工列表 |
| POST | `/api/v1/org/employees` | 创建员工 |
| PUT | `/api/v1/org/employees/:id` | 更新员工 |
| DELETE | `/api/v1/org/employees/:id` | 删除员工 |

### 5.5 基准岗位管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/org/standard-positions` | 基准岗位列表 |
| POST | `/api/v1/org/standard-positions` | 创建基准岗位 |
| PUT | `/api/v1/org/standard-positions/:id` | 更新基准岗位 |
| DELETE | `/api/v1/org/standard-positions/:id` | 删除基准岗位 |

---

## 6. 权限与角色管理接口

### 6.1 权限管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/permissions` | 权限列表 |
| POST | `/api/v1/permissions` | 创建权限 |
| PUT | `/api/v1/permissions/:id` | 更新权限 |
| DELETE | `/api/v1/permissions/:id` | 删除权限 |

### 6.2 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/roles` | 角色列表 |
| POST | `/api/v1/roles` | 创建角色 |
| PUT | `/api/v1/roles/:id` | 更新角色 |
| DELETE | `/api/v1/roles/:id` | 删除角色 |
| GET | `/api/v1/roles/:id/perms` | 获取角色权限 |
| PUT | `/api/v1/roles/:id/perms` | 设置角色权限 |

---

## 7. 系统管理接口

### 7.1 菜单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/menus/tree` | 菜单树形结构 |

### 7.2 字典管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/dicts/:type` | 按类型获取字典 |

### 7.3 日志管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/logs/operations` | 操作日志 |
| GET | `/api/v1/logs/login` | 登录日志 |

---

## 8. 告警管理接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/alerts/rules` | 告警规则列表 |
| POST | `/api/v1/alerts/rules` | 创建告警规则 |
| GET | `/api/v1/alerts` | 告警记录列表 |
| GET | `/api/v1/dashboard/stats` | 仪表盘统计 |

---

## 9. 健康检查

### 9.1 服务健康检查

**请求方法：** GET
**请求路径：** `/health`
**是否需要认证：** 否

**响应示例：**

```json
{
  "status": "ok"
}
```

---

## 附录：错误码定义

| 错误码 | error_code | 说明 |
|--------|------------|------|
| 0 | — | 成功 |
| 400 | — | 参数错误（通用） |
| 401 | — | 未授权 / Token 过期 |
| 403 | — | 账号禁用 |
| 404 | — | 资源不存在 |
| 4002 | ERR_DEVICE_002 | 设备不存在 |
| 4003 | ERR_DEVICE_003 | 非法设备状态 |
| 4005 | ERR_VALIDATION | 参数校验失败 |
| 5001 | ERR_INTERNAL | 服务器内部错误 |

---

## 修订记录

| 日期 | 版本 | 修改内容 |
|------|------|----------|
| 2026-03-20 | V1.0 | 初始版本，与后端代码对齐 |
