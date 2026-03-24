# 历史回放 PRD

## 1. 功能概述
历史回放模块允许用户回放宠物设备在某段时间内的完整运行状态和行为数据，还原真实的宠物生活场景。

## 2. 页面布局与交互

### 页面路径
`/digital-twin/playback` → `HistoricalReplayView.vue`

### 回放控制
- 时间选择器（开始时间~结束时间）
- 回放速度（1x/2x/4x/8x）
- 播放/暂停/停止按钮
- 进度条（可拖拽）

### 回放视图
- 左侧：时间轴事件列表
- 右侧：地图/行为动画
- 底部：实时数据面板

## 3. API 契约

### 获取回放数据
- 路径：`GET /api/v1/digital-twin/replay/:pet_id`
- 实际路径（digital_twin_controller.go）：`GET /api/v1/digital-twin/replay/:pet_id`
- 参数：`start_time`, `end_time`
- 响应：
```json
{
  "code": 0,
  "data": {
    "pet_id": 1,
    "start_time": "2024-01-01T08:00:00Z",
    "end_time": "2024-01-01T12:00:00Z",
    "events": [
      {
        "time": "2024-01-01T08:00:00Z",
        "event_type": "woke_up",
        "location": { "lat": 39.9, "lng": 116.4 }
      },
      {
        "time": "2024-01-01T08:30:00Z",
        "event_type": "eating",
        "location": { "lat": 39.9, "lng": 116.4 }
      }
    ],
    "vitals": [
      { "time": "2024-01-01T08:00:00Z", "heart_rate": 80 }
    ]
  }
}
```

## 4. 验收标准
- [ ] 回放数据正确加载
- [ ] 播放/暂停/停止正常
- [ ] 速度切换正常
- [ ] 进度拖拽正常
- [ ] 时间轴事件点击跳转正常
