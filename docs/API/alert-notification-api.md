# 告警通知 API

**控制器：** `controllers/alert_controller.go`, `controllers/notification_controller.go`  
**路由前缀：** `/api/v1`

---

## 告警管理 (alert_controller)

### 1. 告警规则

#### 1.1 获取告警规则列表

```
GET /api/v1/alerts/rules
```

#### 响应示例

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "低电量告警",
        "device_id": "dev_abc123",
        "alert_type": "low_battery",
        "condition": "lt",
        "threshold": 20,
        "severity": "warning",
        "enabled": true,
        "notify_ways": ["push", "sms"],
        "remark": "电量低于20%时触发",
        "created_at": "2026-03-22T10:00:00Z"
      }
    ]
  }
}
```

#### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| name | string | 规则名称 |
| device_id | string | 关联设备ID |
| alert_type | string | 告警类型：`low_battery` / `offline` / `temperature` / `geofence_violation` 等 |
| condition | string | 条件：`lt`（小于）/ `gt`（大于）/ `eq`（等于） |
| threshold | float | 阈值 |
| severity | string | 严重程度：`info` / `warning` / `critical` |
| enabled | bool | 是否启用 |
| notify_ways | array | 通知方式：`push` / `sms` / `email` |

---

#### 1.2 创建告警规则

```
POST /api/v1/alerts/rules
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 规则名称 |
| device_id | string | 否 | 关联设备ID（可选，关联所有设备时为空） |
| alert_type | string | 是 | 告警类型 |
| condition | string | 是 | 比较条件：`lt` / `gt` / `eq` |
| threshold | float | 是 | 阈值 |
| severity | string | 是 | 严重程度 |
| enabled | bool | 否 | 是否启用，默认 true |
| notify_ways | array | 否 | 通知方式列表 |
| remark | string | 否 | 备注 |

#### 请求示例

```json
{
  "name": "低电量告警",
  "device_id": "dev_abc123",
  "alert_type": "low_battery",
  "condition": "lt",
  "threshold": 20,
  "severity": "warning",
  "enabled": true,
  "notify_ways": ["push", "sms"],
  "remark": "电量低于20%时触发"
}
```

#### 响应示例

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "低电量告警",
    "device_id": "dev_abc123",
    "alert_type": "low_battery",
    "condition": "lt",
    "threshold": 20,
    "severity": "warning",
    "enabled": true,
    "notify_ways": ["push", "sms"],
    "remark": "电量低于20%时触发",
    "created_at": "2026-03-22T10:00:00Z"
  }
}
```

---

#### 1.3 更新告警规则

```
PUT /api/v1/alerts/rules/:id
```

#### 请求参数 (JSON)

与创建规则相同，所有字段可选。

#### 响应示例

```json
{
  "code": 0,
  "data": { /* 更新后的规则对象 */ }
}
```

---

#### 1.4 删除告警规则

```
DELETE /api/v1/alerts/rules/:id
```

#### 响应示例

```json
{
  "code": 0,
  "message": "删除成功"
}
```

---

#### 1.5 批量删除告警规则

```
POST /api/v1/alerts/rules/batch-delete
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| ids | array | 是 | 规则ID列表 |

#### 请求示例

```json
{
  "ids": [1, 2, 3]
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "删除成功"
}
```

---

### 2. 告警记录

#### 2.1 获取告警记录

```
GET /api/v1/alerts
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 否 | 设备ID |
| status | string | 否 | 状态：`pending`（待处理）/ `confirmed`（已确认）/ `resolved`（已解决）/ `ignored`（已忽略） |
| alert_type | string | 否 | 告警类型 |
| severity | string | 否 | 严重程度 |

#### 响应示例

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "device_id": "dev_abc123",
        "alert_type": "low_battery",
        "severity": "warning",
        "status": 1,
        "message": "设备电量低于20%",
        "trigger_value": 15,
        "threshold": 20,
        "created_at": "2026-03-22T10:00:00Z",
        "confirmed_at": null,
        "resolved_at": null
      }
    ]
  }
}
```

#### 告警状态值

| 值 | 说明 |
|----|------|
| 1 | 待处理 |
| 2 | 已确认 |
| 3 | 已解决 |
| 4 | 已忽略 |

---

#### 2.2 确认告警

```
POST /api/v1/alerts/:id/confirm
```

#### 响应示例

```json
{
  "code": 0,
  "message": "告警已确认",
  "data": {
    "id": 1,
    "status": 2,
    "confirmed_at": "2026-03-22T11:00:00Z",
    "confirmed_by": "user_123"
  }
}
```

---

#### 2.3 解决告警

```
POST /api/v1/alerts/:id/resolve
```

#### 响应示例

```json
{
  "code": 0,
  "message": "告警已解决",
  "data": {
    "id": 1,
    "status": 3,
    "resolved_at": "2026-03-22T12:00:00Z",
    "resolved_by": "user_123"
  }
}
```

---

#### 2.4 忽略告警

```
POST /api/v1/alerts/:id/ignore
```

#### 响应示例

```json
{
  "code": 0,
  "message": "告警已忽略",
  "data": {
    "id": 1,
    "status": 4,
    "ignored_at": "2026-03-22T12:00:00Z",
    "ignored_by": "user_123"
  }
}
```

---

#### 2.5 批量确认告警

```
POST /api/v1/alerts/batch-confirm
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| alert_ids | array | 是 | 告警ID列表 |

#### 请求示例

```json
{
  "alert_ids": [1, 2, 3]
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "已确认 3 条告警"
}
```

---

#### 2.6 批量解决告警

```
POST /api/v1/alerts/batch-resolve
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| alert_ids | array | 是 | 告警ID列表 |

#### 响应示例

```json
{
  "code": 0,
  "message": "已解决 3 条告警"
}
```

---

## 通知管理 (notification_controller)

### 3. 通知列表

```
GET /api/v1/notifications
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 否 | 设备ID |
| status | string | 否 | 状态：`pending` / `sent` / `read` |
| priority | string | 否 | 优先级：`low` / `normal` / `high` |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20，最大 100 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "device_id": "dev_abc123",
        "title": "系统升级通知",
        "content": "系统将于今晚10点进行升级",
        "priority": "normal",
        "channel": "push",
        "status": "sent",
        "sent_at": "2026-03-22T10:00:00Z",
        "read_at": null,
        "created_at": "2026-03-22T09:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 4. 发送设备通知

向指定设备下发通知消息。

```
POST /api/v1/notifications/devices/:device_id
```

#### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 是 | 设备ID |

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 通知标题 |
| content | string | 是 | 通知内容 |
| priority | int | 否 | 优先级，1=低 2=中 3=高，默认 2 |
| channel | string | 否 | 渠道：`push` / `mqtt`，默认 `push` |
| template_id | int | 否 | 模板ID（使用时 title/content 将被模板内容替换） |
| variables | object | 否 | 模板变量（当 template_id 存在时） |
| created_by | string | 否 | 创建人 |

#### 请求示例

```json
{
  "title": "系统升级通知",
  "content": "系统将于今晚10点进行升级，请及时保存数据",
  "priority": 2,
  "channel": "push"
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "device_id": "dev_abc123",
    "status": "sent",
    "sent_at": "2026-03-22T10:00:00Z"
  }
}
```

---

### 5. 获取通知详情

```
GET /api/v1/notifications/:id
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "device_id": "dev_abc123",
    "title": "系统升级通知",
    "content": "系统将于今晚10点进行升级",
    "priority": 2,
    "channel": "push",
    "status": "sent",
    "sent_at": "2026-03-22T10:00:00Z",
    "created_by": "admin",
    "created_at": "2026-03-22T09:00:00Z"
  }
}
```

---

### 6. 通过模板发送通知

使用预设模板向多个设备发送通知。

```
POST /api/v1/notifications/template
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| template_id | int | 是 | 模板ID |
| target_type | string | 是 | 目标类型：`all`（全部设备）/ `device`（指定设备）/ `user`（指定用户） |
| target_ids | array | 否 | 目标ID列表（当 target_type 为 `device` 或 `user` 时） |
| variables | object | 否 | 模板变量 |
| created_by | string | 否 | 创建人 |

#### 请求示例

```json
{
  "template_id": 1,
  "target_type": "device",
  "target_ids": ["dev_abc123", "dev_def456"],
  "variables": {
    "username": "张三",
    "updatetime": "今晚10点"
  },
  "created_by": "admin"
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "sent_count": 2,
    "notifications": [
      {
        "id": 1,
        "device_id": "dev_abc123",
        "title": "升级提醒",
        "content": "张三您好，系统将于今晚10点进行升级",
        "status": "sent"
      },
      {
        "id": 2,
        "device_id": "dev_def456",
        "title": "升级提醒",
        "content": "李四您好，系统将于今晚10点进行升级",
        "status": "sent"
      }
    ]
  }
}
```

---

### 7. 删除通知

```
DELETE /api/v1/notifications/:id
```

#### 响应示例

```json
{
  "code": 0,
  "message": "删除成功"
}
```

---

## 告警设置

### 8. 告警设置列表

```
GET /api/v1/alert-settings
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "setting_type": "email",
        "enabled": true,
        "config": {
          "smtp_host": "smtp.example.com",
          "smtp_port": 587,
          "recipients": ["admin@example.com"]
        }
      }
    ]
  }
}
```

---

### 9. 创建告警设置

```
POST /api/v1/alert-settings
```

### 10. 更新告警设置

```
PUT /api/v1/alert-settings/:id
```

### 11. 删除告警设置

```
DELETE /api/v1/alert-settings/:id
```

---

## 告警历史

### 12. 告警历史列表

```
GET /api/v1/alert-history
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| device_id | string | 否 | 设备ID |
| alert_type | string | 否 | 告警类型 |
| start_time | string | 否 | 开始时间 |
| end_time | string | 否 | 结束时间 |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页条数 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "device_id": "dev_abc123",
        "alert_type": "low_battery",
        "message": "设备电量低于20%",
        "severity": "warning",
        "status": "resolved",
        "created_at": "2026-03-22T10:00:00Z",
        "resolved_at": "2026-03-22T12:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

## MQTT 通知下发

通知通过 MQTT 协议下发到设备。

### Topic 格式

```
/device/{device_id}/down/notification
```

### Payload 格式

```json
{
  "type": "notification",
  "title": "通知标题",
  "content": "通知内容",
  "priority": 2,
  "timestamp": "2026-03-22T10:00:00Z"
}
```

### 设备回应 Topic

```
/device/{device_id}/up/notification
```

### 设备回应 Payload

```json
{
  "notification_id": 1,
  "status": "delivered",
  "delivered_at": "2026-03-22T10:00:01Z"
}
```
