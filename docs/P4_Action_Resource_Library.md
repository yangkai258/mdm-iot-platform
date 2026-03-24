# 动作资源库 PRD

## 1. 功能概述
动作资源库提供宠物动作动画的共享平台，用户可以上传、下载、评分和管理动作资源。

## 2. 页面布局与交互

### 页面路径
`/market/actions` → `ActionMarketView.vue`

### 动作库
- 分类（行走/跑步/跳舞/特技/休息）
- 动作卡片列表
- 预览动画

## 3. API 契约

### 动作资源列表
- 路径：`GET /api/v1/market/actions`
- 响应：
```json
{
  "code": 0,
  "items": [
    {
      "id": 1,
      "name": "欢快跑步",
      "category": "running",
      "thumbnail_url": "https://...",
      "duration": 3.0,
      "difficulty": "medium",
      "downloads": 2000,
      "rating": 4.6
    }
  ]
}
```

### 创建动作资源
- 路径：`POST /api/v1/market/actions`
- 请求体：`{ "name": "...", "category": "...", "animation_data": "...", "price": 0 }`

### 发布动作
- 路径：`POST /api/v1/market/actions/:id/publish`

## 4. 验收标准
- [ ] 动作列表加载正常
- [ ] 预览播放正常
- [ ] 下载功能正常
- [ ] 评分正常
