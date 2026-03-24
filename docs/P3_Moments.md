# 精彩瞬间 PRD

## 1. 功能概述
精彩瞬间模块自动捕捉并记录宠物生活中的有趣、可爱或重要的时刻，支持相册展示和社交分享。

## 2. 页面布局与交互

### 页面路径
`/digital-twin/moments` → `MomentsView.vue`

### 相册视图
- 筛选：类型（可爱/里程碑/搞笑）
- 时间轴/网格切换
- 照片墙展示

### 瞬间详情
- 大图预览
- 拍摄时间/地点
- AI描述
- 「分享」「删除」按钮

## 3. API 契约

### 获取精彩瞬间
- 路径：`GET /api/v1/digital-twin/highlights`
- 实际路径（digital_twin_controller.go）：`GET /api/v1/digital-twin/highlights`
- 参数：`pet_id`, `type`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "pet_id": 1,
        "type": "cute",
        "title": "第一次站立",
        "media_url": "https://...",
        "captured_at": "2024-01-01T08:00:00Z",
        "description": "AI自动识别：宠物第一次站立成功"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

### 标注瞬间
- 路径：`POST /api/v1/digital-twin/highlights`
- 请求体：`{ "pet_id": 1, "type": "milestone", "title": "...", "media_url": "...", "description": "..." }`

## 4. 数据模型

### HighlightMoment（精彩瞬间）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pet_id | uint | 宠物ID |
| type | string | cute/milestone/funny |
| title | string | 标题 |
| media_url | string | 媒体URL |
| captured_at | datetime | 拍摄时间 |
| description | string | AI描述 |
| tenant_id | string | 租户ID |

## 5. 验收标准
- [ ] 瞬间列表加载正常
- [ ] 按类型筛选有效
- [ ] 大图预览正常
- [ ] 分享功能正常
- [ ] 删除功能正常
