# 设备远程锁定/擦除 PRD

## 1. 功能概述
设备远程锁定/擦除功能允许管理员通过管理后台对丢失或被盗的设备执行远程锁定、数据擦除操作，保护用户隐私和数据安全。

## 2. 页面布局与交互

### 页面路径
`/devices/remote` → `DeviceRemoteControl.vue`

### 页面布局
- 顶部：设备选择器（选择目标设备）
- 左侧：设备信息卡片（设备ID、型号、在线状态、绑定用户）
- 右侧：操作面板
  - 锁定设备
  - 解锁设备
  - 远程擦除（数据清除）
  - 恢复出厂设置

### 操作说明
- 「锁定设备」：下发 MQTT 锁定命令，设备收到后进入锁定状态
- 「解锁设备」：下发解锁命令，设备恢复正常
- 「远程擦除」：擦除设备上所有用户数据（需二次确认+输入验证码）
- 「恢复出厂」：擦除+恢复出厂设置（需最高权限确认）

### 设备列表（快速操作）
| 列 | 说明 |
|----|------|
| 设备ID | device_id |
| 设备型号 | hardware_model |
| 在线状态 | is_online（绿点/灰点）|
| 绑定用户 | bind_user_id |
| 操作 | 锁定/解锁/擦除按钮 |

## 3. API 契约

### 锁定设备
- 路径：`POST /api/v1/devices/:device_id/security/lock`
- 路径（device_controller.go中实际定义）：`POST /api/v1/devices/:device_id/lock`
- 请求体：无（或可选 reason）
```json
{
  "reason": "丢失设备"
}
```
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "status": "lock_sent",
    "message": "锁定命令已下发"
  },
  "message": "锁定命令已下发"
}
```

### 解锁设备
- 路径：`POST /api/v1/devices/:device_id/security/unlock`
- 路径（实际）：`POST /api/v1/devices/:device_id/unlock`
- 请求体：
```json
{
  "unlock_code": "123456"
}
```
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "status": "unlock_sent"
  }
}
```

### 远程擦除
- 路径：`POST /api/v1/devices/:device_id/security/wipe`
- 请求体：
```json
{
  "wipe_type": "data_only",
  "confirm_code": "123456"
}
```
- wpe_type: `data_only`（仅数据）/ `factory_reset`（恢复出厂）
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "task_id": "wipe-task-uuid",
    "status": "wipe_initiated"
  }
}
```

### 查询设备安全状态
- 路径：`GET /api/v1/devices/:device_id/security/status`
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "is_locked": false,
    "is_wiped": false,
    "lock_history": [],
    "wipe_history": []
  }
}
```

### 获取擦除历史
- 路径：`GET /api/v1/devices/:device_id/wipe-history`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "task_id": "xxx",
        "wipe_type": "data_only",
        "status": "completed",
        "initiated_by": "admin",
        "initiated_at": "2024-01-01T00:00:00Z",
        "completed_at": "2024-01-01T00:05:00Z"
      }
    ]
  }
}
```

## 4. 数据模型

### WipeHistory（擦除历史）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| device_id | string | 设备ID |
| task_id | string | 任务UUID |
| wipe_type | string | data_only / factory_reset |
| status | string | initiated / completed / failed |
| initiated_by | string | 操作人 |
| initiated_at | datetime | 发起时间 |
| completed_at | datetime | 完成时间 |
| error_message | string | 错误信息 |

### DeviceSecurityStatus（设备安全状态）
| 字段 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |
| is_locked | bool | 是否被锁定 |
| is_jailbroken | bool | 是否越狱 |
| last_lock_at | datetime | 最近锁定时间 |
| last_unlock_at | datetime | 最近解锁时间 |

## 5. 验收标准
- [ ] 设备列表显示在线/离线状态
- [ ] 成功下发锁定命令（MQTT QoS=1）
- [ ] 成功下发解锁命令
- [ ] 远程擦除需二次确认+验证码
- [ ] 操作日志完整记录
- [ ] 设备不在线时提示用户
- [ ] 擦除任务可追踪状态
