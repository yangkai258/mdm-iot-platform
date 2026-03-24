# 行为预测 PRD

## 1. 功能概述
行为预测模块基于历史行为数据，运用机器学习模型预测宠物的短期行为趋势，为用户提供个性化的护理建议。

## 2. 页面布局与交互

### 页面路径
`/digital-twin/prediction` → `BehaviorPredictionView.vue`

### 预测结果展示
- 设备/宠物选择
- 预测时间范围（未来6小时/12小时/24小时）
- 预测行为分布图
- 护理建议列表

## 3. API 契约

### 获取行为预测
- 路径：`GET /api/v1/digital-twin/behavior/prediction`
- 实际路径（digital_twin_controller.go）：`GET /api/v1/digital-twin/behavior/prediction`
- 参数：`device_id`, `hours` (6/12/24)
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "prediction_time": "2024-01-01T12:00:00Z",
    "predictions": [
      {
        "time_slot": "08:00-09:00",
        "predicted_behavior": "eating",
        "confidence": 0.85,
        "recommendation": "建议此时进行喂食"
      }
    ],
    "confidence_avg": 0.82
  }
}
```

### 行为历史
- 路径：`GET /api/v1/digital-twin/behavior/history`
- 参数：`device_id`, `start_time`, `end_time`

## 4. 验收标准
- [ ] 预测结果正确展示
- [ ] 预测图表正常显示
- [ ] 护理建议准确
- [ ] 时间范围选择正常
