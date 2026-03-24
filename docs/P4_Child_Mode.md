# 儿童模式 PRD

## 1. 功能概述
儿童模式专为有儿童的家庭设计，提供适龄的内容过滤、使用时间管理和家长监护功能。

## 2. 页面布局与交互

### 页面路径
`/family/child-mode` → `ChildModeView.vue`

### 儿童模式配置
- 开关控制
- 可访问功能配置
- 内容过滤级别
- 使用时间限制

### 时间管理
- 每日可用时段
- 禁用时段
- 单次使用时长

### 家长监护
- 使用报告
- 远程锁定

## 3. API 契约

### 获取儿童模式配置
- 路径：`GET /api/v1/family/child-mode/config`
- 响应：
```json
{
  "code": 0,
  "data": {
    "enabled": true,
    "content_filter_level": "strict",
    "daily_limit_minutes": 60,
    "allowed_features": ["games", "education"],
    "blocked_features": ["social", "payment"],
    "time_slots": [
      { "start": "08:00", "end": "20:00" }
    ]
  }
}
```

### 更新配置
- 路径：`PUT /api/v1/family/child-mode/config`
- 请求体：同上

### 使用报告
- 路径：`GET /api/v1/family/child-mode/reports`
- 响应：
```json
{
  "code": 0,
  "data": {
    "usage_today_minutes": 45,
    "usage_history": [
      { "date": "2024-01-01", "minutes": 60 }
    ],
    "top_features": [
      { "feature": "games", "minutes": 30 }
    ]
  }
}
```

## 4. 验收标准
- [ ] 儿童模式开关控制正常
- [ ] 内容过滤生效
- [ ] 时间限制正常
- [ ] 使用报告准确
- [ ] 远程锁定正常
