# AI 决策日志 PRD

## 1. 功能概述
AI 决策日志模块记录 AI 决策引擎的所有推理过程和决策结果，便于问题追溯、模型优化和合规审计。

## 2. 页面布局与交互

### 页面路径
`/ai/logs` → `AIDecisionLogsView.vue`

### 搜索表单
| 字段 | 类型 | 说明 |
|------|------|------|
| 设备ID | Input | - |
| 决策类型 | Select | behavior/emotion/action |
| 时间范围 | DateRange | - |

### 日志列表
| 列 | 说明 |
|----|------|
| 时间 | created_at |
| 设备ID | device_id |
| 决策类型 | decision_type |
| 输入 | input_data |
| 决策结果 | output |
| 置信度 | confidence |
| 延迟 | latency_ms |

## 3. API 契约

### 决策日志列表
- 路径：`GET /api/v1/ai/decision/logs`
- 参数：`device_id`, `decision_type`, `start_time`, `end_time`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "device_id": "device-001",
        "decision_type": "emotion_response",
        "input_data": { "emotion": "happy" },
        "output": { "action": "play_music", "params": { "song": "happy" } },
        "confidence": 0.95,
        "latency_ms": 50,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1000,
    "page": 1,
    "page_size": 20
  }
}
```

### 决策详情
- 路径：`GET /api/v1/ai/decision/logs/:id`

## 4. 验收标准
- [ ] 日志列表分页加载正常
- [ ] 按设备/类型/时间筛选有效
- [ ] 日志详情显示完整输入输出
- [ ] 日志导出功能正常
