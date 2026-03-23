# PRD：宠物用品电商接入

**版本：** V1.0
**所属Phase：** Phase 4（Sprint 27-28）
**优先级：** P3
**负责角色：** agentcp（产品）、agenthd（后端）

---

## 一、概述

### 1.1 模块定位

宠物用品电商接入模块为宠物主人提供一站式宠物食品、玩具、健康用品的购买体验，通过智能推荐算法根据宠物种类、年龄、健康状态推荐适合的商品，实现"健康追踪→推荐→购买"的完整闭环。

### 1.2 核心价值

- **便捷购买**：从健康数据直接推荐商品，减少决策成本
- **个性化推荐**：基于宠物画像精准推荐
- **健康联动**：根据体重、饮食数据推荐相应保健品

---

## 二、功能详情

### 2.1 商品浏览

| 功能 | 说明 |
|------|------|
| 分类浏览 | 按食品/玩具/保健品/日常用品分类 |
| 搜索 | 商品名称/品牌搜索 |
| 筛选排序 | 按价格/销量/评分筛选 |

### 2.2 智能推荐

| 功能 | 说明 |
|------|------|
| 健康推荐 | 根据体重推荐减肥粮/增重粮 |
| 年龄推荐 | 根据宠物年龄推荐阶段粮 |
| 场景推荐 | 夏季推荐降温垫，冬季推荐保暖衣 |

### 2.3 购物车与订单

| 功能 | 说明 |
|------|------|
| 加入购物车 | 从推荐页直接加入 |
| 订单创建 | 创建外部电商订单 |
| 订单状态同步 | 同步发货/配送状态 |

### 2.4 购买历史

| 功能 | 说明 |
|------|------|
| 历史订单 | 展示历史购买记录 |
| 复购 | 一键复购历史订单 |

---

## 三、API接口定义

### 3.1 商品接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/ecomm/products | 商品列表 |
| GET | /api/v1/ecomm/products/:id | 商品详情 |
| GET | /api/v1/ecomm/products/recommend | 智能推荐商品 |

### 3.2 购物车

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/ecomm/cart | 购物车列表 |
| POST | /api/v1/ecomm/cart | 加入购物车 |
| DELETE | /api/v1/ecomm/cart/:id | 删除购物车项 |

### 3.3 订单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/ecomm/orders | 创建订单 |
| GET | /api/v1/ecomm/orders | 订单列表 |
| GET | /api/v1/ecomm/orders/:id | 订单详情 |

---

## 四、数据库设计

### 4.1 推荐商品表 (eco_recommended_products)

```sql
CREATE TABLE eco_recommended_products (
    id              BIGSERIAL PRIMARY KEY,
    product_id      VARCHAR(100) NOT NULL,
    product_name    VARCHAR(255) NOT NULL,
    category        VARCHAR(50) NOT NULL,
    pet_type        VARCHAR(50),
    health_condition VARCHAR(100),
    recommendation_score DECIMAL(5,2),
    external_url    VARCHAR(500),
    image_url       VARCHAR(500),
    price           DECIMAL(10,2),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.2 购物车表 (eco_cart)

```sql
CREATE TABLE eco_cart (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    product_id      VARCHAR(100) NOT NULL,
    quantity        INT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.3 订单表 (eco_orders)

```sql
CREATE TABLE eco_orders (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    external_order_id VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'pending',
    total_amount    DECIMAL(10,2),
    items           JSONB,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 五、前端页面

| 页面 | 路由 | 说明 |
|------|------|------|
| 商品推荐 | /ecomm/recommend | 智能推荐商品 |
| 购物车 | /ecomm/cart | 购物车管理 |
| 订单列表 | /ecomm/orders | 订单历史 |

---

## 六、验收标准

- [ ] 智能推荐商品准确率>80%（基于宠物画像）
- [ ] 购物车CRUD功能正常
- [ ] 订单创建成功，状态同步正常
