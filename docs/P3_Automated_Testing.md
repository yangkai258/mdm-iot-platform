# 自动化测试 PRD

## 1. 功能概述
自动化测试模块提供设备功能自动化测试框架，支持创建测试用例、批量执行测试、生成测试报告。

## 2. 页面布局与交互

### 页面路径
`/simulation/test` → `TestFrameworkView.vue` + `ABExperimentView.vue`

### 测试用例管理
- 用例列表
- 「新建用例」按钮
- 用例编辑（DSL编写）
- 执行历史

### A/B 实验
- 实验列表
- 实验配置
- 结果分析

## 3. API 契约

### 测试场景列表
- 路径：`GET /api/v1/simulation/scenes`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "OTA升级测试",
        "scene_type": "ota",
        "description": "测试OTA升级流程"
      }
    ]
  }
}
```

### 创建测试
- 路径：`POST /api/v1/simulation/scenes`
- 请求体：`{ "name": "...", "scene_type": "...", "config": {...} }`

### 执行测试
- 路径：`POST /api/v1/simulation/scenes/:id/run`

### 获取测试结果
- 路径：`GET /api/v1/simulation/sessions/:id/results`
- 响应：
```json
{
  "code": 0,
  "data": {
    "session_id": "xxx",
    "status": "completed",
    "passed": 95,
    "failed": 5,
    "total": 100,
    "duration_seconds": 300
  }
}
```

### A/B 实验
- 路径：`POST /api/v1/simulation/ab-experiments`
- 路径：`GET /api/v1/simulation/ab-experiments/:id`
- 路径：`POST /api/v1/simulation/ab-experiments/:id/run`
- 路径：`GET /api/v1/simulation/ab-experiments/:id/results`

## 4. 验收标准
- [ ] 测试场景创建成功
- [ ] 测试执行正常
- [ ] 测试结果正确
- [ ] A/B实验创建和执行正常
- [ ] 测试报告生成正确
