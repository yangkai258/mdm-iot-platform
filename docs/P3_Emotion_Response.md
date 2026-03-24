# 情绪响应配置 PRD

## 1. 功能概述
情绪响应配置模块允许管理员为宠物配置不同情绪状态下的个性化响应动作，实现情绪关怀的自动化。

## 2. 页面布局与交互

### 页面路径
`/emotion/response` → `EmotionResponseConfigView.vue`

### 响应配置列表
| 列 | 说明 |
|----|------|
| 情绪类型 | emotion_type |
| 响应模式 | response_mode（安慰/鼓励/互动/安静）|
| 动作序列 | actions（JSON预览）|
| 启用状态 | enabled |
| 操作 | 编辑/删除 |

### 配置编辑弹窗
| 字段 | 类型 | 说明 |
|------|------|------|
| 情绪类型 | Select | - |
| 响应模式 | Select | - |
| 动作序列 | DragSort List | 拖拽排序的动作列表 |
| 启用 | Switch | - |

## 3. API 契约

### 获取情绪响应配置
- 路径：`GET /api/v1/emotions/pet/config`
- 响应：
```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "pet_id": 1,
      "emotion_type": "happy",
      "response_mode": "play",
      "actions": [
        { "type": "sound", "data": { "file": "happy.mp3" } },
        { "type": "movement", "data": { "action": "dance" } }
      ],
      "enabled": true
    }
  ]
}
```

### 更新配置
- 路径：`PUT /api/v1/emotions/pet/config`
- 请求体：
```json
{
  "pet_id": 1,
  "emotion_type": "happy",
  "response_mode": "play",
  "actions": [...],
  "enabled": true
}
```

## 4. 验收标准
- [ ] 配置列表加载正常
- [ ] 动作序列拖拽排序正常
- [ ] 保存配置成功
- [ ] 启用/禁用开关生效
- [ ] 响应动作正确执行
