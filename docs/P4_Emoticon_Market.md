# 表情包市场 PRD

## 1. 功能概述
表情包市场提供宠物数字表情包的浏览、购买、下载和管理功能，丰富宠物与主人之间的互动体验。

## 2. 页面布局与交互

### 页面路径
`/market/emojis` → `EmoticonMarketView.vue`

### 表情包商店
- 分类标签（可爱/搞笑/节日/自定义）
- 表情包卡片列表
- 预览动画效果

### 我的表情包
- 已购表情包列表
- 表情包管理

## 3. API 契约

### 表情包列表
- 路径：`GET /api/v1/market/emoticons`
- 实际路径（market_controller.go）：`GET /api/v1/market/emoticons`
- 参数：`pack_type`, `keyword`, `status`, `sort_by`, `order`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "items": [
    {
      "id": 1,
      "pack_name": "可爱猫咪表情",
      "pack_type": "cute",
      "thumbnail_urls": ["https://..."],
      "emoticon_count": 16,
      "price": 6.00,
      "downloads": 5000,
      "rating_avg": 4.8,
      "status": "published"
    }
  ],
  "total": 50,
  "page": 1,
  "page_size": 20
}
```

### 创建表情包
- 路径：`POST /api/v1/market/emoticons`
- 请求体：`{ "pack_name": "...", "pack_type": "...", "emoticons": [...], "price": 6.00 }`

### 获取表情包详情
- 路径：`GET /api/v1/market/emoticons/:id`

### 发布表情包
- 路径：`POST /api/v1/market/emoticons/:id/publish`

### 删除表情包
- 路径：`DELETE /api/v1/market/emoticons/:id`

## 4. 数据模型

### EmoticonPack（表情包）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pack_name | string | 表情包名称 |
| pack_type | string | cute/funny/holiday/custom |
| emoticon_urls | JSONB | 表情URL数组 |
| price | float | 价格 |
| downloads | int | 下载次数 |
| rating_avg | float | 平均评分 |
| status | string | draft/published/archived |

## 5. 验收标准
- [ ] 表情包列表分类加载正常
- [ ] 预览动画正常播放
- [ ] 购买/下载正常
- [ ] 表情包管理正常
- [ ] 评分功能正常
