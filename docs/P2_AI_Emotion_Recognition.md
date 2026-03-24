# AI 情感识别 PRD

## 1. 功能概述
AI 情感识别模块通过分析宠物的表情、动作和叫声等数据，识别宠物当前的情绪状态（开心/悲伤/愤怒/平静/兴奋/焦虑），为宠物情绪健康管理提供数据支持。

## 2. 页面布局与交互

### 页面路径
`/ai/emotion` → `AIEmotionView.vue`

### 实时情绪监测（顶部）
- 设备选择器
- 实时情绪状态（表情图标+文字）
- 情绪强度（1-10刻度条）

### 情绪趋势图（中部）
- 24小时情绪曲线
- 情绪分布饼图

### 情绪日志列表
| 列 | 说明 |
|----|------|
| 时间 | recorded_at |
| 宠物 | pet_id |
| 情绪类型 | emotion_type |
| 强度 | intensity |
| 触发原因 | trigger |
| AI响应 | ai_response |
| 操作 | 查看详情 |

## 3. API 契约

### 情绪记录列表
- 路径：`GET /api/v1/emotions/records`
- 实际路径（emotion_controller.go）：`GET /api/v1/emotions/records`
- 参数：`pet_id`, `user_id`, `emotion_type`, `start_date`, `end_date`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "pet_id": 1,
        "user_id": 100,
        "source": "voice",
        "emotion_type": "happy",
        "intensity": 8,
        "trigger": "用户回家",
        "context": "检测到尾巴摇摆",
        "ai_response": "播放欢迎音乐",
        "recorded_at": "2024-01-01T08:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 创建情绪记录
- 路径：`POST /api/v1/emotions/records`
- 请求体：
```json
{
  "pet_id": 1,
  "user_id": 100,
  "source": "voice",
  "emotion_type": "happy",
  "intensity": 8,
  "trigger": "用户回家",
  "context": "检测到尾巴摇摆"
}
```

### 情绪统计
- 路径：`GET /api/v1/emotions/records/stats`
- 参数：`pet_id`, `period`
- 响应：
```json
{
  "code": 0,
  "data": {
    "total_records": 500,
    "dominant_emotion": "happy",
    "avg_intensity": 7.2,
    "by_type": {
      "happy": 200,
      "calm": 150,
      "excited": 100,
      "anxious": 50
    }
  }
}
```

### 情绪响应动作
- 路径：`GET /api/v1/emotions/pet/actions`
- 路径：`POST /api/v1/emotions/pet/actions`
- 路径：`PUT /api/v1/emotions/pet/actions/:id`
- 路径：`DELETE /api/v1/emotions/pet/actions/:id`

### 情绪响应配置
- 路径：`GET /api/v1/emotions/pet/config`
- 路径：`PUT /api/v1/emotions/pet/config`

### 情绪报告
- 路径：`GET /api/v1/emotions/report`
- 参数：`pet_id`, `period` (daily/weekly/monthly)
- 响应：
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "pet_id": 1,
    "period": "weekly",
    "start_date": "2024-01-01",
    "end_date": "2024-01-07",
    "summary": { "happy": 60, "calm": 30 },
    "avg_intensity": 7.5,
    "dominant_emotion": "happy",
    "trend": "improving"
  }
}
```

## 4. 数据模型

### EmotionRecord（情绪记录）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pet_id | uint | 宠物ID |
| user_id | uint | 用户ID |
| source | string | voice/text/pet_behavior |
| emotion_type | string | happy/sad/angry/calm/excited/anxious |
| intensity | int | 1-10 |
| trigger | string | 触发原因 |
| context | string | 上下文 |
| ai_response | string | AI响应动作 |
| recorded_at | datetime | 记录时间 |

### PetEmotionAction（宠物情绪响应动作）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pet_id | uint | 宠物ID |
| emotion_type | string | 情绪类型 |
| action_type | string | gesture/sound/movement |
| action_data | JSONB | 动作数据 |
| priority | int | 优先级 |
| enabled | bool | 是否启用 |

## 5. 验收标准
- [ ] 情绪日志列表分页加载正常
- [ ] 实时情绪状态正确显示
- [ ] 情绪趋势图正确
- [ ] 情绪类型分布饼图正确
- [ ] 情绪配置保存成功
- [ ] 情绪报告生成正确
