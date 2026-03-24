# 虚拟宠物仿真 PRD

## 1. 功能概述
虚拟宠物仿真模块在虚拟环境中模拟宠物行为，用于新功能演示、算法验证和用户培训，无需真实硬件即可体验完整功能。

## 2. 页面布局与交互

### 页面路径
`/simulation/pet` → `VirtualPetSimulationView.vue`

### 仿真控制台
- 环境选择（室内/室外/公园）
- 宠物选择
- 「开始仿真」按钮
- 仿真速度控制

### 仿真视图
- 3D 虚拟宠物渲染
- 行为状态显示
- 传感器数据面板

## 3. API 契约

### 启动仿真
- 路径：`POST /api/v1/simulation/pet/start`
- 实际路径（simulation_controller.go）：`POST /api/v1/simulation/pet/start`
- 请求体：
```json
{
  "pet_id": 1,
  "environment": "indoor",
  "duration": 3600
}
```

### 停止仿真
- 路径：`POST /api/v1/simulation/pet/stop`
- 请求体：`{ "session_id": "xxx" }`

### 获取宠物状态
- 路径：`GET /api/v1/simulation/pet/:session_id`
- 响应：
```json
{
  "code": 0,
  "data": {
    "session_id": "xxx",
    "pet_id": 1,
    "status": "running",
    "elapsed_seconds": 120,
    "current_behavior": "exploring",
    "vitals": { "heart_rate": 85, "energy": 70 }
  }
}
```

## 4. 验收标准
- [ ] 仿真启动正常
- [ ] 仿真状态正确显示
- [ ] 仿真停止正常
- [ ] 仿真数据正确
