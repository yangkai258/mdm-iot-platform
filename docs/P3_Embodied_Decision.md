# 具身决策 PRD

## 1. 功能概述
具身决策模块是设备的大脑，综合环境感知、空间认知、行为记忆等信息，做出符合宠物当前状态和用户意图的智能决策。

## 2. 页面布局与交互

### 页面路径
`/embodied/decision` → `DecisionLogsView.vue`

### 决策日志
| 列 | 说明 |
|----|------|
| 时间 | created_at |
| 设备ID | device_id |
| 决策类型 | decision_type |
| 输入上下文 | context |
| 决策策略 | strategy |
| 结果 | result |
| 延迟 | latency_ms |

### 决策策略配置
- 策略列表（安全优先/体验优先/节能优先）
- 策略参数配置

## 3. API 契约

### 获取决策上下文
- 路径：`GET /api/v1/embodied/:device_id/decision/context`
- 响应：包含当前环境、宠物状态、用户偏好等

### 设置决策策略
- 路径：`POST /api/v1/embodied/:device_id/decision/strategy`
- 请求体：`{ "strategy": "safety_first", "params": {...} }`

### 获取决策日志
- 路径：`GET /api/v1/embodied/:device_id/decision/logs`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "device_id": "device-001",
        "decision_type": "navigation",
        "context": { "battery": 20, "goal": "home" },
        "strategy": "safety_first",
        "result": "proceed_to_home",
        "latency_ms": 50
      }
    ]
  }
}
```

## 4. 验收标准
- [ ] 决策日志分页加载正常
- [ ] 策略切换正常
- [ ] 决策延迟统计正确
- [ ] 上下文信息完整
