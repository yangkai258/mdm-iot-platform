# 动作模仿 PRD

## 1. 功能概述
动作模仿模块记录和学习宠物的动作行为，并能够复现或改编这些动作，提供丰富的动作资源库。

## 2. 页面布局与交互

### 页面路径
`/embodied/motion` → `ActionLibraryView.vue`

### 动作库
- 动作分类（行走/跑步/跳舞/特技）
- 动作列表（缩略图+名称）
- 「录制新动作」按钮

### 动作详情
- 动作预览动画
- 动作参数配置
- 「执行」「分享」按钮

## 3. API 契约

### 动作库列表
- 路径：`GET /api/v1/embodied/action-library`
- 实际路径（embodied_controller.go）：`GET /api/v1/embodied/action-library`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "欢快跳舞",
        "category": "dance",
        "thumbnail_url": "https://...",
        "duration": 5,
        "difficulty": "easy"
      }
    ]
  }
}
```

### 录制动作
- 路径：`POST /api/v1/embodied/action-library/record`
- 请求体：`{ "device_id": "xxx", "name": "...", "motion_data": [...] }`

### 学习动作
- 路径：`POST /api/v1/embodied/action-library/:id/learn`
- 请求体：`{ "pet_id": 1 }`

### 执行动作
- 路径：`POST /api/v1/embodied/:device_id/action/execute`
- 请求体：`{ "action_id": 1, "params": { "speed": 1.0 } }`

### 停止动作
- 路径：`POST /api/v1/embodied/:device_id/action/stop`

## 4. 验收标准
- [ ] 动作库列表加载正常
- [ ] 动作录制成功
- [ ] 动作执行正常
- [ ] 动作停止正常
- [ ] 动作分享正常
