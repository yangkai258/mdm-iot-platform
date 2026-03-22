# 设备管理 API

**控制器：** `controllers/device_controller.go`  
**路由前缀：** `/api/v1/devices`

---

## 1. 设备注册

设备首次上电时调用，注册设备信息。

### 请求

```
POST /api/v1/devices/register
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| mac_address | string | 是 | MAC 地址，格式 `XX:XX:XX:XX:XX:XX` |
| sn_code | string | 是 | 设备序列号 |
| hardware_model | string | 是 | 硬件型号，如 `M5Stack-Core2` |
| firmware_version | string | 是 | 固件版本号，如 `1.0.0` |
| device_id | string | 否 | 设备ID（可选，服务端自动生成） |

### 请求示例

```json
{
  "mac_address": "AA:BB:CC:DD:EE:FF",
  "sn_code": "SN20260322001",
  "hardware_model": "M5Stack-Core2",
  "firmware_version": "1.0.0"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "dev_abc123",
    "lifecycle_status": 1,
    "created_at": "2026-03-22T10:00:00Z"
  }
}
```

### 响应字段

| 字段 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备唯一标识 |
| lifecycle_status | int | 生命周期状态：1=待激活 2=服役中 3=维修 4=报废 |
| created_at | datetime | 注册时间 |

---

## 2. 设备列表

获取设备列表，支持分页和筛选。

### 请求

```
GET /api/v1/devices
```

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20，最大 100 |
| status | string | 否 | 在线状态筛选：`online` / `offline` |
| lifecycle_status | int | 否 | 生命周期状态：1=待激活 2=服役中 3=维修 4=报废 |
| hardware_model | string | 否 | 硬件型号筛选 |
| device_type | string | 否 | 设备类型筛选（等同于 hardware_model） |
| tenant_id | string | 否 | 租户ID筛选 |
| search | string | 否 | 关键词搜索（device_id / sn_code / mac_address） |
| start_time | string | 否 | 创建时间范围开始，格式 `2006-01-02 15:04:05` |
| end_time | string | 否 | 创建时间范围结束，格式 `2006-01-02 15:04:05` |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "device_id": "dev_abc123",
        "sn_code": "SN20260322001",
        "mac_address": "AA:BB:CC:DD:EE:FF",
        "hardware_model": "M5Stack-Core2",
        "firmware_version": "1.0.0",
        "lifecycle_status": 2,
        "is_online": true,
        "battery_level": 85,
        "created_at": "2026-03-22T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

---

## 3. 扫码绑定设备

用户扫描设备二维码，将设备绑定到当前用户。

### 请求

```
POST /api/v1/devices/bind/:sn_code
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| sn_code | string | 是 | 设备序列号 |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| bind_user_id | string | 是 | 绑定用户ID |

### 请求示例

```json
{
  "bind_user_id": "user_123"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "dev_abc123",
    "bind_user_id": "user_123",
    "lifecycle_status": 2,
    "message": "绑定成功"
  }
}
```

### 错误码

| code | error_code | 说明 |
|------|-------------|------|
| 4002 | ERR_DEVICE_002 | 设备不存在 |
| 4003 | ERR_DEVICE_003 | 设备状态不允许绑定 |

---

## 4. 解绑设备

将设备从用户账户解除绑定。

### 请求

```
POST /api/v1/devices/unbind/:sn_code
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| sn_code | string | 是 | 设备序列号 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "dev_abc123",
    "lifecycle_status": 1,
    "message": "解绑成功"
  }
}
```

---

## 5. 获取设备详情

根据设备ID获取设备详细信息。

### 请求

```
GET /api/v1/devices/:device_id
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "device_id": "dev_abc123",
    "sn_code": "SN20260322001",
    "mac_address": "AA:BB:CC:DD:EE:FF",
    "hardware_model": "M5Stack-Core2",
    "firmware_version": "1.0.0",
    "lifecycle_status": 2,
    "bind_user_id": "user_123",
    "is_online": true,
    "battery_level": 85,
    "created_at": "2026-03-22T10:00:00Z",
    "updated_at": "2026-03-22T12:00:00Z"
  }
}
```

---

## 6. 更新设备信息

更新设备基本信息。

### 请求

```
PUT /api/v1/devices/:device_id
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| hardware_model | string | 否 | 硬件型号 |
| firmware_version | string | 否 | 固件版本 |
| device_name | string | 否 | 设备名称 |
| device_type | string | 否 | 设备类型 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": { /* 更新后的设备对象 */ }
}
```

---

## 7. 删除设备

删除设备记录。

### 请求

```
DELETE /api/v1/devices/:device_id
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success"
}
```

---

## 8. 更新设备状态

更新设备生命周期状态。

### 请求

```
PUT /api/v1/devices/:device_id/status
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| lifecycle_status | int | 是 | 状态值：1=待激活 2=服役中 3=维修 4=报废 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "dev_abc123",
    "lifecycle_status": 3
  }
}
```

---

## 9. 获取设备影子

获取设备的期望状态（Desired State）。

### 请求

```
GET /api/v1/devices/:device_id/desired-state
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "dev_abc123",
    "desired_state": {
      "mood": "happy",
      "energy": 80,
      "language": "zh-CN"
    },
    "timestamp": "2026-03-22T10:00:00Z"
  }
}
```

---

## 10. 设置设备影子

设置设备的期望状态。

### 请求

```
PUT /api/v1/devices/:device_id/desired-state
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| mood | string | 否 | 心情状态 |
| energy | int | 否 | 能量值 0-100 |
| language | string | 否 | 语言设置 |
| custom_params | object | 否 | 自定义参数 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "dev_abc123",
    "desired_state": {
      "mood": "happy",
      "energy": 80,
      "language": "zh-CN"
    }
  }
}
```

---

## 11. 发送设备指令

向设备下发控制指令。

### 请求

```
POST /api/v1/devices/:device_id/commands
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| command_type | string | 是 | 指令类型，如 `action`、`query`、`control` |
| action_id | string | 否 | 动作库ID（当 command_type=action 时） |
| parameters | object | 否 | 指令参数 |

### 请求示例

```json
{
  "command_type": "action",
  "action_id": "action_wiggle",
  "parameters": {
    "duration_ms": 500,
    "intensity": 3
  }
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "command_id": "cmd_xyz789",
    "device_id": "dev_abc123",
    "status": "sent",
    "sent_at": "2026-03-22T10:00:00Z"
  }
}
```

---

## 12. 获取指令历史

查询设备指令下发历史。

### 请求

```
GET /api/v1/devices/:device_id/commands
```

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20 |
| command_type | string | 否 | 指令类型筛选 |
| status | string | 否 | 状态筛选：`pending` / `sent` / `delivered` / `executed` |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "command_id": "cmd_xyz789",
        "command_type": "action",
        "action_id": "action_wiggle",
        "status": "executed",
        "sent_at": "2026-03-22T10:00:00Z",
        "executed_at": "2026-03-22T10:00:01Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

---

## 13. 设备宠物档案-获取

获取设备关联的宠物档案。

### 请求

```
GET /api/v1/devices/:device_id/profile
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "dev_abc123",
    "pet_name": "小爪",
    "pet_type": "cat",
    "birthday": "2024-06-01",
    "gender": "male",
    "breed": "橘猫",
    "avatar_url": "https://example.com/avatar.png",
    "personality": {
      "trait": "curious",
      "favorite_toy": "feather"
    }
  }
}
```

---

## 14. 设备宠物档案-更新

更新设备关联的宠物档案。

### 请求

```
PUT /api/v1/devices/:device_id/profile
```

### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pet_name | string | 否 | 宠物名称 |
| pet_type | string | 否 | 宠物类型 |
| birthday | string | 否 | 生日，格式 `YYYY-MM-DD` |
| gender | string | 否 | 性别：`male` / `female` |
| breed | string | 否 | 品种 |
| avatar_url | string | 否 | 头像URL |
| personality | object | 否 | 性格属性 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": { /* 更新后的宠物档案 */ }
}
```

---

## OTA 固件升级 API

### 15. 创建固件包

```
POST /api/v1/ota/packages
```

### 16. 固件包列表

```
GET /api/v1/ota/packages
```

### 17. 创建部署任务

```
POST /api/v1/ota/deployments
```

### 18. 部署任务列表

```
GET /api/v1/ota/deployments
```

### 19. 获取部署详情

```
GET /api/v1/ota/deployments/:id
```

### 20. 暂停部署

```
POST /api/v1/ota/deployments/:id/pause
```

### 21. 恢复部署

```
POST /api/v1/ota/deployments/:id/resume
```

### 22. 取消部署

```
POST /api/v1/ota/deployments/:id/cancel
```

### 23. 获取部署进度

```
GET /api/v1/ota/deployments/:id/progress
```

### 24. 设备 OTA 上报

```
POST /api/v1/ota/devices/:device_id/report
```

### 25. 设备检查 OTA

```
GET /api/v1/ota/devices/:device_id/check
```
