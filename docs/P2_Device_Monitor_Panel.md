# 设备监控面板 PRD

## 1. 功能概述
设备监控面板提供全平台设备的实时运行状态监控、告警统计、在线率分析和设备分布概览，帮助运维人员快速掌握设备群整体健康状况。

## 2. 页面布局与交互

### 页面路径
`/devices/monitor` → `DeviceMonitorPanel.vue`

### 页面布局
- 顶部：统计卡片区（4个KPI卡片）
- 中部左侧：在线设备列表
- 中部右侧：告警趋势图表（7天折线图）
- 底部：设备分布/固件版本分布/告警类型分布

### 统计卡片（顶部）
| 卡片 | 说明 |
|------|------|
| 设备总数 | total_devices |
| 在线设备 | online_devices |
| 在线率 | online_rate (百分比) |
| 活跃告警 | active_alerts |

### 筛选条件
| 字段 | 类型 | 说明 |
|------|------|------|
| 设备型号 | Select（M5Stack/ESP32等）|
| 时间范围 | DateRange（近24h/7d/30d）|
| 租户 | Select（多租户场景）|

### 在线设备列表
| 列 | 说明 |
|----|------|
| 设备ID | device_id |
| 在线状态 | is_online（绿点/离线）|
| 最后心跳 | last_heartbeat |
| 电池电量 | battery_level（%）|
| 当前模式 | current_mode |

## 3. API 契约

### 监控仪表板汇总
- 路径：`GET /api/v1/devices/monitor/dashboard`
- 实际路径（device_monitor_controller.go）：`GET /api/v1/devices/monitor/dashboard`
- 响应：
```json
{
  "code": 0,
  "data": {
    "summary": {
      "total_devices": 1000,
      "online_devices": 850,
      "offline_devices": 150,
      "online_rate": 85.0,
      "active_alerts": 12,
      "today_alerts": 45,
      "today_resolved": 38
    },
    "status_distribution": [
      { "status": 1, "count": 100 },
      { "status": 2, "count": 800 }
    ],
    "model_distribution": [
      { "model": "M5Stack", "count": 600 },
      { "model": "ESP32", "count": 400 }
    ],
    "alert_trend": [
      { "date": "2024-01-01", "count": 10 },
      { "date": "2024-01-02", "count": 15 }
    ],
    "online_devices": [
      {
        "device_id": "device-001",
        "is_online": true,
        "last_seen": "2024-01-01T00:00:00Z",
        "battery_level": 85,
        "current_mode": "normal"
      }
    ]
  }
}
```

### 设备指标统计
- 路径：`GET /api/v1/devices/monitor/metrics`
- 参数：`hours` (int, default 24, max 720)
- 响应：
```json
{
  "code": 0,
  "data": {
    "period_hours": 24,
    "total_devices": 1000,
    "active_devices": 950,
    "total_alerts": 120,
    "high_severity_alerts": 5,
    "online_rate": 85.0,
    "firmware_distribution": [
      { "version": "v1.2.3", "count": 500 }
    ],
    "uptime_distribution": [
      { "bucket": "<1小时", "count": 50 },
      { "bucket": "1-6小时", "count": 200 },
      { "bucket": "6-24小时", "count": 300 }
    ],
    "alert_type_distribution": [
      { "alert_type": "battery_low", "count": 50 }
    ],
    "mode_distribution": [
      { "mode": "normal", "count": 600 }
    ]
  }
}
```

### 单设备实时监控数据
- 路径：`GET /api/v1/devices/:device_id/monitor/realtime`
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "is_online": true,
    "battery_level": 85,
    "current_mode": "normal",
    "last_heartbeat": "2024-01-01T00:00:00Z",
    "last_ip": "192.168.1.100",
    "is_jailbroken": false,
    "root_status": "normal",
    "location": {
      "latitude": 39.9,
      "longitude": 116.4
    },
    "online_duration": "3天5小时",
    "runtime": "30天5小时12分",
    "lifecycle_status": 2,
    "hardware_model": "M5Stack",
    "firmware_version": "v1.2.3",
    "recent_alerts": []
  }
}
```

### 创建设备告警规则
- 路径：`POST /api/v1/devices/:device_id/monitor/alert-rules`
- 请求体：
```json
{
  "name": "低电量告警",
  "alert_type": "battery_low",
  "condition": "<",
  "threshold": 20,
  "severity": 2,
  "enabled": true,
  "notify_ways": "email,webhook"
}
```

## 4. 数据模型

### DeviceAlert（设备告警）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| device_id | string | 设备ID |
| alert_type | string | 告警类型 |
| severity | int | 严重程度1-4 |
| message | string | 告警消息 |
| status | int | 状态1-未确认/2-处理中/3-已解决 |
| created_at | datetime | 发生时间 |
| resolved_at | datetime | 解决时间 |

### DeviceAlertRule（设备告警规则）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| device_id | string | 设备ID |
| name | string | 规则名称 |
| alert_type | string | battery_low/offline/high_temp |
| condition | string | </>/=/ >= / <= |
| threshold | float | 阈值 |
| severity | int | 1-4 |
| enabled | bool | 是否启用 |
| notify_ways | string | email,sms,webhook |

## 5. 验收标准
- [ ] 仪表板统计卡片数据正确加载
- [ ] 在线/离线设备数量计算正确
- [ ] 7天告警趋势折线图显示正常
- [ ] 设备型号分布饼图显示正常
- [ ] 设备在线时长分布柱状图显示正常
- [ ] 单设备实时监控数据刷新正常
- [ ] 告警规则创建成功
- [ ] 时间范围筛选生效
