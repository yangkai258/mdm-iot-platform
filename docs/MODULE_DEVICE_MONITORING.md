# 设备监控面板

**版本：** V1.0  
**模块负责人：** agentcp  
**编制日期：** 2026-03-22  

---

## 1. 概述

设备监控面板提供实时传感器数据展示、在线设备地图、设备健康状态监控等能力，帮助运营人员实时掌握设备运行状态。

**业务目标：**
- 实时展示设备传感器数据
- 展示在线设备分布地图
- 设备健康状态监控和告警
- 支持历史数据查询

---

## 2. 功能列表

| 功能 | 描述 | 优先级 | 触发方式 |
|------|------|--------|----------|
| 实时传感器数据 | 实时显示电池/位置/温度等 | P0 | 自动 |
| 在线设备地图 | GIS 地图展示设备位置 | P1 | 自动 |
| 设备健康状态 | 健康/亚健康/异常状态展示 | P0 | 自动 |
| 历史数据查询 | 查询历史传感器数据 | P1 | 人工 |
| 设备日志 | 设备操作日志查看 | P1 | 人工 |
| 告警通知 | 设备异常时发送通知 | P0 | 自动 |
| 远程诊断 | 远程获取设备诊断信息 | P1 | 人工 |

---

## 3. 数据模型

### 3.1 传感器数据记录表 (sensor_data_log)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, not null |
| sensor_type | string | 传感器类型 | battery/location/temperature/ultrasonic |
| sensor_value | json | 传感器值 | not null |
| recorded_at | datetime | 记录时间 | not null |
| created_at | datetime | 创建时间 | auto |

### 3.2 设备健康状态表 (device_health)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, unique, not null |
| health_status | int | 健康状态 | 1=健康 2=亚健康 3=异常 |
| battery_level | int | 电池电量 | 0-100 |
| is_online | bool | 在线状态 | default false |
| last_seen_at | datetime | 最后在线时间 | auto |
| error_count_24h | int | 24小时错误次数 | default 0 |
| alert_count_24h | int | 24小时告警次数 | default 0 |
| updated_at | datetime | 更新时间 | auto |

### 3.3 设备日志表 (device_logs)

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| id | uint | 主键 | PK, auto_increment |
| device_id | string | 设备ID | FK to devices, not null |
| log_type | string | 日志类型 | system/action/error/alert |
| log_level | int | 日志级别 | 1=debug 2=info 3=warning 4=error |
| message | text | 日志内容 | not null |
| metadata | json | 扩展数据 | nullable |
| created_at | datetime | 创建时间 | auto |

---

## 4. 接口定义

### 4.1 获取设备实时状态

```
GET /api/v1/devices/:device_id/realtime-status
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "pet-001",
    "is_online": true,
    "battery_level": 85,
    "position": { "x": 120.5, "y": 80.3 },
    "temperature": 25.5,
    "last_seen_at": "2026-03-22T10:30:00Z",
    "health_status": 1
  }
}
```

### 4.2 获取设备日志

```
GET /api/v1/devices/:device_id/logs
```

### 4.3 获取传感器历史数据

```
GET /api/v1/devices/:device_id/sensor-history
```

### 4.4 获取健康状态列表

```
GET /api/v1/devices/health
```

### 4.5 远程诊断

```
POST /api/v1/devices/:device_id/diagnostics
```

---

## 5. 前端页面

| 页面 | 路由 | 功能 |
|------|------|------|
| 设备监控面板 | /monitoring/dashboard | 实时传感器、健康状态 |
| 在线设备地图 | /monitoring/map | GIS 地图设备分布 |
| 设备日志 | /monitoring/logs | 设备日志查询 |
| 健康状态 | /monitoring/health | 设备健康列表 |
