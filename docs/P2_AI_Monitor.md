# AI 模型监控 PRD

## 1. 功能概述
AI 模型监控模块对平台 AI 模型（行为分析、情感识别、具身决策等）的运行状态、性能指标和调用情况进行监控告警，确保 AI 服务的稳定性。

## 2. 页面布局与交互

### 页面路径
`/ai/monitor` → `AIMonitorView.vue`

### 监控概览（顶部卡片）
- 在线模型数
- 总调用次数（今日）
- 平均响应时间
- 错误率

### 监控面板
- 模型列表（表格）
- 调用趋势图
- 错误分布图

### 模型列表
| 列 | 说明 |
|----|------|
| 模型名称 | model_name |
| 版本 | version |
| 状态 | status（running/stopped/error）|
| 调用量（今日）| calls_today |
| 平均延迟 | avg_latency_ms |
| 错误率 | error_rate |
| 操作 | 停止/重启/查看详情 |

## 3. API 契约

### 模型监控列表
- 路径：`GET /api/v1/ai/monitor/models`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "model_id": "behavior-v1",
        "model_name": "行为分析模型",
        "version": "1.2.3",
        "status": "running",
        "calls_today": 5000,
        "avg_latency_ms": 120,
        "error_rate": 0.5,
        "last_inference_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 模型指标详情
- 路径：`GET /api/v1/ai/monitor/models/:id`
- 响应：包含详细指标

### 指标历史
- 路径：`GET /api/v1/ai/monitor/metrics`
- 参数：`model_id`, `hours`
- 响应：
```json
{
  "code": 0,
  "data": {
    "latency_trend": [
      { "time": "2024-01-01T00:00:00Z", "avg_ms": 120 }
    ],
    "error_trend": [
      { "time": "2024-01-01T00:00:00Z", "count": 5 }
    ],
    "throughput_trend": [
      { "time": "2024-01-01T00:00:00Z", "calls": 500 }
    ]
  }
}
```

## 4. 验收标准
- [ ] 模型列表加载正常
- [ ] 监控指标实时更新
- [ ] 错误率计算正确
- [ ] 告警触发正常
- [ ] 停止/重启模型成功
