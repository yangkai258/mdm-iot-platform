# 情绪识别 PRD

## 1. 功能概述
情绪识别模块通过多模态数据（语音语调、面部表情、肢体动作、文字）对用户和宠物的情绪状态进行实时识别和记录。

## 2. 页面布局与交互

### 页面路径
`/emotion/recognition` → `EmotionRecognizeView.vue`

### 实时情绪面板
- 设备选择器
- 当前情绪状态（图标+文字）
- 情绪强度仪表盘
- 情绪持续时间

### 情绪历史
- 日志列表（同 AI 情感识别）

## 3. API 契约

### 获取实时情绪
- 路径：`GET /api/v1/emotions/realtime/:device_id`
- 响应：
```json
{
  "code": 0,
  "data": {
    "device_id": "device-001",
    "emotion_type": "happy",
    "intensity": 8,
    "confidence": 0.95,
    "source": "multimodal",
    "detected_at": "2024-01-01T00:00:00Z"
  }
}
```

### 情绪记录列表
- 路径：`GET /api/v1/emotions/records`
- 响应：同 AI 情感识别

## 4. 验收标准
- [ ] 实时情绪正确显示
- [ ] 情绪历史加载正常
- [ ] 多模态融合识别工作
