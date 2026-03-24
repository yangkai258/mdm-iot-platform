# AI 沙箱测试 PRD

## 1. 功能概述
AI 沙箱测试模块提供隔离的测试环境，支持在不影响生产的情况下对 AI 模型进行功能测试、性能压测和效果评估。

## 2. 页面布局与交互

### 页面路径
`/ai/sandbox` → `AISandboxView.vue`

### 测试场景
- 场景选择（新建/选择已有）
- 输入配置（测试数据）
- 执行测试
- 结果展示

### 测试配置
| 字段 | 类型 | 说明 |
|------|------|------|
| 场景名称 | Input | - |
| 测试模型 | Select | 选择模型版本 |
| 输入类型 | Select | 文本/图片/语音 |
| 测试数据 | Textarea/File | - |
| 预期输出 | Textarea | - |

## 3. API 契约

### 创建测试
- 路径：`POST /api/v1/ai/sandbox/tests`
- 请求体：
```json
{
  "name": "情感识别测试v1",
  "model_id": "emotion-v1",
  "model_version": "1.2.3",
  "input_type": "text",
  "input_data": "宠物表现出开心的样子",
  "expected_output": "happy"
}
```

### 执行测试
- 路径：`POST /api/v1/ai/sandbox/tests/:id/run`
- 响应：
```json
{
  "code": 0,
  "data": {
    "test_id": "test-uuid",
    "status": "running"
  }
}
```

### 获取测试结果
- 路径：`GET /api/v1/ai/sandbox/tests/:id`
- 响应：
```json
{
  "code": 0,
  "data": {
    "test_id": "test-uuid",
    "status": "completed",
    "result": {
      "output": "happy",
      "confidence": 0.95,
      "latency_ms": 120
    },
    "passed": true,
    "actual_vs_expected": "matched"
  }
}
```

### 批量测试
- 路径：`POST /api/v1/ai/sandbox/batch`
- 请求体：`{ "tests": [...] }`

## 4. 验收标准
- [ ] 测试场景创建成功
- [ ] 测试执行正常
- [ ] 测试结果正确展示
- [ ] 批量测试正常
- [ ] 性能指标（延迟、吞吐量）正确
