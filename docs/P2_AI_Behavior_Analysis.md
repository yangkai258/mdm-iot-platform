# AI 行为分析 PRD

## 1. 功能概述
AI 行为分析模块对宠物设备的日常行为数据进行采集、存储和分析，通过机器学习识别宠物行为模式，为宠物健康监测和个性化服务提供数据支撑。

## 2. 页面布局与交互

### 页面路径
`/ai/behavior` → `AIBehaviorView.vue`

### 数据概览（顶部卡片）
- 今日活跃设备数
- 行为事件总数（今日）
- 异常行为数
- 平均响应延迟

### 行为分析（中部）
- 左侧：行为类型分布（饼图）
- 右侧：行为趋势（折线图）

### 行为日志列表
| 列 | 说明 |
|----|------|
| 设备ID | device_id |
| 宠物昵称 | pet_name |
| 行为类型 | behavior_type |
| 开始时间 | start_time |
| 持续时长 | duration |
| 详情 | metadata（JSON预览）|

## 3. API 契约

### 行为列表
- 路径：`GET /api/v1/ai/behaviors`
- 实际路径（behavior_controller.go）：`GET /api/v1/ai/pets/:device_id/behaviors`
- 参数：`device_id`, `behavior_type`, `start_time`, `end_time`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "pet_id": 1,
        "type": "eating",
        "start_time": "2024-01-01T08:00:00Z",
        "end_time": "2024-01-01T08:15:00Z",
        "duration": 900,
        "metadata": { "food_amount": "normal" }
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 创建行为记录
- 路径：`POST /api/v1/ai/pets/:device_id/behaviors/trigger`
- 请求体：
```json
{
  "behavior_type": "eating",
  "start_time": "2024-01-01T08:00:00Z",
  "duration": 900,
  "metadata": { "food_amount": "normal" }
}
```

### 行为统计
- 路径：`GET /api/v1/ai/behaviors/stats`
- 参数：`device_id`, `period` (day/week/month)
- 响应：
```json
{
  "code": 0,
  "data": {
    "total_events": 500,
    "by_type": {
      "eating": 100,
      "sleeping": 200,
      "playing": 150,
      "other": 50
    },
    "avg_duration": { "eating": 900, "sleeping": 28800 }
  }
}
```

### 异常行为检测
- 路径：`GET /api/v1/ai/behaviors/anomalies`
- 参数：`device_id`, `severity`, `page`, `page_size`
- 响应：异常行为列表

## 4. 数据模型

### BehaviorEvent（行为事件）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pet_id | uint | 宠物ID |
| type | string | eating/sleeping/playing 等 |
| start_time | datetime | 开始时间 |
| end_time | datetime | 结束时间 |
| duration | int | 持续秒数 |
| metadata | string | JSON额外数据 |
| tenant_id | string | 租户ID |

## 5. 验收标准
- [ ] 行为日志列表分页加载正常
- [ ] 按设备/类型/时间筛选有效
- [ ] 行为类型分布饼图正确
- [ ] 行为趋势折线图正确
- [ ] 异常行为标记正常
- [ ] 触发行为接口正常
