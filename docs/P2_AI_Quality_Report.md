# AI 质量报告 PRD

## 1. 功能概述
AI 质量报告模块对 AI 模型的质量指标进行周期性汇总分析，包括准确率、召回率、F1分数、响应延迟等，生成可视化的质量评估报告。

## 2. 页面布局与交互

### 页面路径
`/ai/quality` → `AIQualityDashboardView.vue`

### 质量概览
- 模型健康度评分
- 准确率趋势图
- 延迟分布图
- Top 问题类型

### 详细报告
- 按模型分页查看
- 导出报告（PDF/Excel）

## 3. API 契约

### 质量报告列表
- 路径：`GET /api/v1/ai/quality/reports`
- 参数：`model_id`, `period`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "model_id": "emotion-v1",
        "model_name": "情感识别模型",
        "period": "weekly",
        "start_date": "2024-01-01",
        "end_date": "2024-01-07",
        "accuracy": 0.92,
        "recall": 0.89,
        "f1_score": 0.905,
        "avg_latency_ms": 150,
        "total_inferences": 50000,
        "error_count": 500,
        "created_at": "2024-01-07T00:00:00Z"
      }
    ],
    "total": 10
  }
}
```

### 报告详情
- 路径：`GET /api/v1/ai/quality/reports/:id`
- 响应：包含完整质量指标和图表数据

### 生成报告
- 路径：`POST /api/v1/ai/quality/reports/generate`
- 请求体：`{ "model_id": "emotion-v1", "period": "weekly" }`

## 4. 验收标准
- [ ] 质量报告列表加载正常
- [ ] 准确率/召回率/F1计算正确
- [ ] 延迟分布统计正确
- [ ] 报告导出正常
- [ ] 趋势图正确显示
