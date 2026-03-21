# 设备影子与期望状态

**版本：** V1.0  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

设备影子模块实现 MQTT 的 desired/reported 模式，用于管理设备的期望状态和实际状态。运营人员可以预设设备行为，设备定期同步并上报执行结果。

**业务目标：**
- 实现设备期望状态下发
- 跟踪设备实际状态
- 支持设备配置批量同步

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 |
|------|------|--------|----------|
| 期望状态设置 | 设置设备的期望行为/配置 | P0 | 人工 |
| 期望状态查询 | 设备查询待执行的期望状态 | P0 | 自动 |
| 实际状态上报 | 设备上报当前实际状态 | P0 | 自动 |
| 状态同步 | 定期同步 desired 和 reported | P0 | 自动 |
| 配置批量下发 | 批量设置多设备期望状态 | P1 | 人工 |
| 状态差异告警 | desired 与 reported 不一致告警 | P1 | 自动 |

---

## 3. 数据模型

### 3.1 设备期望状态表 (device_desired_state)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, unique, not null |
| desired_state | json | 期望状态 JSON | not null |
| version | int | 状态版本号 | not null, auto increment |
| set_by | uint | 设置人 | FK to sys_users, not null |
| set_at | datetime | 设置时间 | auto |
| applied_at | datetime | 设备应用时间 | nullable |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

### 3.2 设备实际状态表 (device_reported_state)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, unique, not null |
| reported_state | json | 实际上报状态 | not null |
| version | int | 状态版本号 | not null |
| reported_at | datetime | 上报时间 | auto |
| sync_status | int | 同步状态 | 1=同步 2=差异 3=超时 |
| created_at | datetime | 创建时间 | auto |
| updated_at | datetime | 更新时间 | auto |

---

## 4. 接口定义

### 4.1 设置期望状态

```
PUT /api/v1/devices/:device_id/desired-state
```

**请求示例：**
```json
{
  "state": {
    "pet_name": "小爪",
    "personality": "活泼",
    "dnd_enabled": true,
    "dnd_start": "22:00",
    "dnd_end": "08:00",
    "interaction_frequency": "normal"
  }
}
```

### 4.2 获取期望状态

```
GET /api/v1/devices/:device_id/desired-state
```

### 4.3 获取实际状态

```
GET /api/v1/devices/:device_id/reported-state
```

### 4.4 获取状态差异

```
GET /api/v1/devices/:device_id/state-diff
```

### 4.5 批量设置期望状态

```
POST /api/v1/devices/batch-desired-state
```

---

## 5. MQTT 协议

### 5.1 期望状态下发

```
Topic: /miniclaw/{device_id}/down/desired_state
```

**Payload:**
```json
{
  "version": 5,
  "state": {
    "pet_name": "小爪",
    "dnd_enabled": true
  },
  "timestamp": 1710902400000
}
```

### 5.2 实际状态上报

```
Topic: /miniclaw/{device_id}/up/reported_state
```

**Payload:**
```json
{
  "version": 4,
  "state": {
    "pet_name": "小爪",
    "current_action": "idle",
    "mood": 75
  },
  "timestamp": 1710902400000
}
```
