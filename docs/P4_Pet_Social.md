# 宠物社交 PRD

## 1. 功能概述
宠物社交模块为宠物主人提供社区交流平台，支持发布宠物动态、关注其他宠物、约伴玩耍（Playdate）等功能。

## 2. 页面布局与交互

### 页面路径
`/pet/social` → `PetFeedView.vue` + `PetFollowView.vue` + `PetPlaydateView.vue`

### Tab1：宠物动态
- 动态信息流
- 发布动态入口
- 点赞/评论/分享

### Tab2：关注
- 关注列表
- 粉丝列表
- 关注/取消关注按钮

### Tab3：约伴玩耍
- 活动列表
- 创建活动
- 报名/取消报名

## 3. API 契约

### 动态列表
- 路径：`GET /api/v1/pet-social/posts`
- 实际路径（pet_social_controller.go）：`GET /api/v1/pet-social/posts`
- 参数：`pet_id`, `user_id`, `page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "posts": [
    {
      "id": 1,
      "pet_id": 1,
      "author_id": 100,
      "content": "今天天气真好！",
      "media_urls": ["https://..."],
      "like_count": 25,
      "comment_count": 5,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

### 发布动态
- 路径：`POST /api/v1/pet-social/posts`
- 请求体：`{ "pet_id": 1, "author_id": 100, "content": "...", "media_urls": [...], "is_public": true }`

### 点赞
- 路径：`POST /api/v1/pet-social/posts/:id/like`
- 路径：`DELETE /api/v1/pet-social/posts/:id/like`

### 评论
- 路径：`GET /api/v1/pet-social/posts/:id/comments`
- 路径：`POST /api/v1/pet-social/posts/:id/comments`
- 请求体：`{ "user_id": 100, "content": "好可爱！" }`

### 关注
- 路径：`GET /api/v1/pet-social/following`
- 路径：`POST /api/v1/pet-social/follow`
- 请求体：`{ "follower_id": 100, "followee_id": 200, "follow_type": "pet" }`
- 路径：`DELETE /api/v1/pet-social/follow/:id`

### 粉丝列表
- 路径：`GET /api/v1/pet-social/followers`

### 约伴列表
- 路径：`GET /api/v1/pet-social/playdates`
- 响应：
```json
{
  "code": 0,
  "list": [
    {
      "id": 1,
      "organizer_id": 100,
      "title": "周末宠物派对",
      "location": "朝阳公园",
      "start_time": "2024-01-07T10:00:00Z",
      "end_time": "2024-01-07T16:00:00Z",
      "status": "pending",
      "max_pets": 10
    }
  ]
}
```

### 创建约伴
- 路径：`POST /api/v1/pet-social/playdates`
- 请求体：`{ "organizer_id": 100, "title": "...", "location": "...", "start_time": "...", "end_time": "...", "max_pets": 10 }`

### 参加约伴
- 路径：`POST /api/v1/pet-social/playdates/:id/join`

## 4. 数据模型

### Post（动态）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| pet_id | uint | 宠物ID |
| author_id | uint | 作者ID |
| content | string | 内容 |
| media_urls | string | JSON媒体URL数组 |
| like_count | int | 点赞数 |
| comment_count | int | 评论数 |
| is_public | bool | 是否公开 |

### Follow（关注关系）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| follower_id | uint | 关注者 |
| followee_id | uint | 被关注者 |
| follow_type | string | pet/user |

### PetPlaydate（约伴）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| organizer_id | uint | 组织者 |
| title | string | 活动名称 |
| location | string | 地点 |
| start_time | datetime | 开始时间 |
| end_time | datetime | 结束时间 |
| status | string | pending/confirmed/cancelled |

## 5. 验收标准
- [ ] 动态信息流加载正常
- [ ] 发布动态成功
- [ ] 点赞/评论正常
- [ ] 关注/取消关注正常
- [ ] 约伴创建/参加正常
- [ ] 动态分享正常
