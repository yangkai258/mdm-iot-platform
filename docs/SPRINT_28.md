# Sprint 28 规划

**时间**：2026-11-08
**状态**：待开始
**Sprint 周期**：2 周（2026-11-08 ～ 2026-11-21）

---

## 一、Sprint 目标

**目标：** 宠物生活服务闭环与社交

完善宠物生活服务闭环（用品电商、社交分享），并开始宠物用品电商接入和社交平台分享功能，实现宠物主从健康管理到商品购买、内容分享的完整闭环。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **宠物用品电商 API** | 完成 `/api/v1/ecomm/*` 电商接口（基于 PRD_28_PET_SUPPLIES_ECOMMERCE） | ecomm_controller.go | P0 |
| P0-2 | **智能推荐服务** | 实现基于宠物画像的个性化推荐 | recommendation_engine.go | P0 |
| P0-3 | **社交分享 API** | 完成 `/api/v1/social/share/*` 社交分享 | social_share_controller.go | P0 |
| P1-1 | **订单状态同步** | 实现电商订单与发货状态同步 | order_status_sync.go | P1 |
| P1-2 | **购物车联动** | 实现购物车与推荐系统联动 | cart_recommendation.go | P1 |
| P2-1 | **宠物朋友圈** | 实现宠物社交功能 | pet_social_feed.go | P2 |
| P2-2 | **宠物视频互动** | 实现宠物视频分享和互动 | pet_video互动.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **商品推荐页面** | 完成 EcommRecommendView.vue 智能推荐商品 | EcommRecommendView.vue | P0 |
| PF0-2 | **购物车页面** | 完成 CartView.vue 购物车管理 | CartView.vue | P0 |
| PF0-3 | **订单页面** | 完成 OrderView.vue 订单列表和详情 | OrderView.vue | P0 |
| PF0-4 | **社交分享页面** | 完成 SocialShareView.vue 微信/抖音分享 | SocialShareView.vue | P0 |
| PF1-1 | **订单跟踪页面** | 完成 OrderTrackingView.vue 物流跟踪 | OrderTrackingView.vue | P1 |
| PF2-1 | **宠物朋友圈页面** | 完成 PetSocialView.vue 宠物动态feed | PetSocialView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/ecomm/products` | GET | 商品列表 |
| `GET /api/v1/ecomm/products/recommend` | GET | 个性化推荐 |
| `GET /api/v1/ecomm/cart` | GET | 购物车列表 |
| `POST /api/v1/ecomm/cart` | POST | 加入购物车 |
| `PUT /api/v1/ecomm/cart/:id` | PUT | 更新购物车 |
| `DELETE /api/v1/ecomm/cart/:id` | DELETE | 删除购物车项 |
| `POST /api/v1/ecomm/orders` | POST | 创建订单 |
| `GET /api/v1/ecomm/orders` | GET | 订单列表 |
| `GET /api/v1/ecomm/orders/:id` | GET | 订单详情 |
| `GET /api/v1/ecomm/orders/:id/tracking` | GET | 物流跟踪 |
| `POST /api/v1/social/share` | POST | 分享到社交平台 |
| `GET /api/v1/social/share/config` | GET | 分享平台配置 |
| `GET /api/v1/social/feed` | GET | 宠物朋友圈 |
| `POST /api/v1/social/feed` | POST | 发布动态 |
| `GET /api/v1/social/feed/:id` | GET | 动态详情 |

### 数据库设计

```sql
-- 社交分享记录表
CREATE TABLE social_share_records (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    pet_id          BIGINT REFERENCES pets(id),
    platform        VARCHAR(30) NOT NULL,            -- 'wechat'/'douyin'/'weibo'
    content_type    VARCHAR(50) NOT NULL,            -- 'photo'/'video'/'moment'
    media_urls      VARCHAR(500)[],
    share_url       VARCHAR(500),
    description     VARCHAR(255),
    shared_at       TIMESTAMP DEFAULT NOW(),
    INDEX idx_user_share (user_id, shared_at DESC)
);

-- 宠物朋友圈表
CREATE TABLE pet_social_posts (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    pet_id          BIGINT REFERENCES pets(id),
    content         TEXT,
    media_urls      VARCHAR(500)[],
    location_name   VARCHAR(100),
    like_count      INT DEFAULT 0,
    comment_count   INT DEFAULT 0,
    visibility      VARCHAR(20) DEFAULT 'public',   -- 'public'/'friends'/'private'
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_public_feed (visibility, created_at DESC)
);

-- 宠物朋友圈互动表
CREATE TABLE pet_social_interactions (
    id              BIGSERIAL PRIMARY KEY,
    post_id         BIGINT NOT NULL REFERENCES pet_social_posts(id),
    user_id         BIGINT NOT NULL REFERENCES users(id),
    interaction_type VARCHAR(30) NOT NULL,           -- 'like'/'comment'/'share'
    content         TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_post_interactions (post_id)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 智能推荐 | 推荐准确率>80%，响应<1s | A/B测试 |
| 购物车 | CRUD正常，推荐联动正常 | 流程测试 |
| 订单创建 | 订单创建成功，状态同步正常 | 订单测试 |
| 社交分享 | 微信/抖音分享正常，落地页可访问 | 分享测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 推荐响应时间 | < 500ms |
| 朋友圈feed加载 | < 2s |
| 分享回调 | < 10s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| PRD_28_PET_SUPPLIES_ECOMMERCE | 宠物用品电商 PRD |
| Sprint 27 第三方集成 | 宠物医疗/保险基础 |
| 宠物档案 | 宠物画像数据 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 电商平台API不稳定 | 订单失败 | 异步队列+重试机制 |
| 社交平台政策变更 | 分享功能失效 | 及时跟进平台政策 |
| 推荐效果不佳 | 转化率低 | 持续优化推荐算法 |
