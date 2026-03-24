# 实时生命体征 PRD

## 1. 功能概述
实时生命体征模块从设备获取并展示宠物的心率、呼吸率、体温等实时生理数据，提供历史趋势和健康告警。

## 2. 页面布局与交互

### 页面路径
`/digital-twin/vitals` → `RealTimeVitalsChart.vue` + `VitalsDashboardView.vue`

### 实时仪表盘
- 设备选择器
- 三大指标卡片：心率(bpm)/呼吸率(次/分)/体温(℃)
- 实时波形图（动态刷新）
- 正常范围参考线

### 历史趋势
- 时间范围选择（24h/7d/30d）
- 指标趋势折线图
- 异常标记（红色高亮点）

### 告警面板
- 当前告警列表
- 确认/忽略按钮

## 3. API 契约

### 生命体征仪表盘
- 路径：`GET /api/v1/digital-twin/vitals/dashboard`
- 实际路径（digital_twin_controller.go）：`GET /api/v1/digital-twin/vitals/dashboard`
- 参数：`pet_id`
- 响应：
```json
{
  "code": 0,
  "data": {
    "heart_rate": { "latest": 85, "avg": 82, "min": 70, "max": 95, "unit": "bpm" },
    "breathing": { "latest": 25, "avg": 24, "unit": "次/分" },
    "temperature": { "latest": 38.5, "avg": 38.3, "unit": "℃" }
  }
}
```

### 实时生命体征
- 路径：`GET /api/v1/digital-twin/vitals/realtime`
- 参数：`pet_id`
- 响应：
```json
{
  "code": 0,
  "data": {
    "pet_id": 1,
    "heart_rate": 85,
    "breathing_rate": 25,
    "temperature": 38.5,
    "recorded_at": "2024-01-01T00:00:00Z"
  }
}
```

### 生命体征历史
- 路径：`GET /api/v1/digital-twin/vitals/history`
- 参数：`pet_id`, `type`, `start_time`, `end_time`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "pet_id": 1,
        "type": "heart_rate",
        "value": 85.0,
        "unit": "bpm",
        "recorded_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1000,
    "page": 1,
    "page_size": 100
  }
}
```

### 健康告警
- 路径：`GET /api/v1/digital-twin/vitals/alerts`
- 参数：`pet_id`, `status`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "pet_id": 1,
        "type": "abnormal_heart_rate",
        "level": "critical",
        "message": "心率异常：连续3次超过100bpm",
        "status": "pending",
        "detected_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10
  }
}
```

### 确认告警
- 路径：`POST /api/v1/digital-twin/vitals/alerts/:id/confirm`

### 忽略告警
- 路径：`POST /api/v1/digital-twin/vitals/alerts/:id/ignore`

## 4. 数据模型

### VitalRecord（生命体征记录）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pet_id | uint | 宠物ID |
| type | string | heart_rate/breathing/temperature |
| value | float64 | 数值 |
| unit | string | 单位 |
| recorded_at | datetime | 记录时间 |
| source | string | device/manual |

### HealthAlert（健康告警）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pet_id | uint | 宠物ID |
| type | string | 告警类型 |
| level | string | warning/critical |
| message | string | 告警信息 |
| status | string | pending/confirmed/ignored/resolved |
| detected_at | datetime | 检测时间 |

## 5. 验收标准
- [ ] 实时指标卡片数据正确
- [ ] 实时波形图动态刷新
- [ ] 历史趋势图加载正常
- [ ] 健康告警触发正常
- [ ] 确认/忽略告警操作正常
- [ ] 异常数据正确标红
