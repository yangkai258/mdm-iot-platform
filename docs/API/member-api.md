# 会员管理 API

**控制器：** `controllers/member_controller.go`, `controllers/member_enhanced_controller.go`  
**路由前缀：** `/api/v1`

---

## 1. 会员管理

### 1.1 会员列表

获取会员列表，支持分页和筛选。

#### 请求

```
GET /api/v1/members
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 10 |
| keyword | string | 否 | 关键词搜索（姓名/会员号/手机号） |
| level | string | 否 | 会员等级 |
| member_level | string | 否 | 会员等级（别名） |
| points_min | int | 否 | 最小积分 |
| points_max | int | 否 | 最大积分 |
| status | string | 否 | 状态筛选 |
| start_time | string | 否 | 创建时间开始，格式 `YYYY-MM-DD HH:mm:ss` |
| end_time | string | 否 | 创建时间结束，格式 `YYYY-MM-DD HH:mm:ss` |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "member_code": "M20260322001",
        "member_name": "张三",
        "phone": "13800138000",
        "member_level": "gold",
        "points": 5000,
        "status": "active",
        "card": {
          "id": 1,
          "card_code": "CARD001"
        },
        "created_at": "2026-03-22T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 1.2 创建会员

#### 请求

```
POST /api/v1/members
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| member_name | string | 是 | 会员姓名 |
| phone | string | 否 | 手机号 |
| member_code | string | 否 | 会员编号 |
| member_level | string | 否 | 会员等级 |
| email | string | 否 | 邮箱 |
| gender | string | 否 | 性别：`male` / `female` |
| birthday | string | 否 | 生日 |

#### 请求示例

```json
{
  "member_name": "张三",
  "phone": "13800138000",
  "member_code": "M20260322001",
  "member_level": "silver"
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "member_code": "M20260322001",
    "member_name": "张三",
    "phone": "13800138000",
    "member_level": "silver",
    "points": 0,
    "status": "active",
    "created_at": "2026-03-22T10:00:00Z"
  }
}
```

---

### 1.3 更新会员

#### 请求

```
PUT /api/v1/members/:id
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| member_name | string | 否 | 会员姓名 |
| phone | string | 否 | 手机号 |
| member_level | string | 否 | 会员等级 |
| status | string | 否 | 状态 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": { /* 更新后的会员对象 */ }
}
```

---

### 1.4 删除会员

#### 请求

```
DELETE /api/v1/members/:id
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success"
}
```

---

### 1.5 会员详情

#### 请求

```
GET /api/v1/members/:id
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "member_code": "M20260322001",
    "member_name": "张三",
    "phone": "13800138000",
    "member_level": "gold",
    "points": 5000,
    "status": "active",
    "card": {
      "id": 1,
      "card_code": "CARD001",
      "card_name": "金卡"
    },
    "created_at": "2026-03-22T10:00:00Z",
    "updated_at": "2026-03-22T12:00:00Z"
  }
}
```

---

## 2. 会员卡管理

### 2.1 会员卡列表

```
GET /api/v1/member/cards
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 10 |
| keyword | string | 否 | 关键词搜索（卡名/卡号） |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "card_code": "CARD001",
        "card_name": "金卡",
        "discount_rate": 0.9,
        "points_multiplier": 2.0,
        "created_at": "2026-03-22T10:00:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 2.2 创建会员卡

```
POST /api/v1/member/cards
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card_name | string | 是 | 卡名 |
| card_code | string | 否 | 卡号 |
| discount_rate | float | 否 | 折扣率，如 0.9 表示 9 折 |
| points_multiplier | float | 否 | 积分倍数 |
| description | string | 否 | 描述 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "card_code": "CARD001",
    "card_name": "金卡",
    "discount_rate": 0.9,
    "points_multiplier": 2.0
  }
}
```

---

### 2.3 更新会员卡

```
PUT /api/v1/member/cards/:id
```

### 2.4 删除会员卡

```
DELETE /api/v1/member/cards/:id
```

---

## 3. 优惠券管理

### 3.1 优惠券列表

```
GET /api/v1/member/coupons
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 10 |
| keyword | string | 否 | 关键词搜索 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "coupon_code": "COUPON001",
        "coupon_name": "新人优惠券",
        "coupon_type": "discount",
        "discount_amount": 10.00,
        "min_order_amount": 100.00,
        "total_stock": 1000,
        "remain_stock": 500,
        "valid_from": "2026-03-01",
        "valid_until": "2026-03-31",
        "created_at": "2026-03-22T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 3.2 创建优惠券

```
POST /api/v1/member/coupons
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| coupon_name | string | 是 | 优惠券名称 |
| coupon_code | string | 否 | 优惠券代码 |
| coupon_type | string | 否 | 类型：`discount` / `cash` / `gift` |
| discount_amount | float | 否 | 优惠金额 |
| min_order_amount | float | 否 | 最低订单金额 |
| total_stock | int | 否 | 发放总数量 |
| valid_from | string | 否 | 有效期开始 |
| valid_until | string | 否 | 有效期结束 |

---

### 3.3 更新优惠券

```
PUT /api/v1/member/coupons/:id
```

### 3.4 删除优惠券

```
DELETE /api/v1/member/coupons/:id
```

---

## 4. 店铺管理

### 4.1 店铺列表

```
GET /api/v1/member/stores
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 10 |
| keyword | string | 否 | 关键词搜索 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "store_code": "STORE001",
        "store_name": "旗舰店",
        "address": "北京市朝阳区xxx",
        "phone": "010-12345678",
        "status": "active"
      }
    ],
    "total": 30,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 4.2 创建店铺

```
POST /api/v1/member/stores
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| store_name | string | 是 | 店铺名称 |
| store_code | string | 否 | 店铺编码 |
| address | string | 否 | 地址 |
| phone | string | 否 | 联系电话 |

---

### 4.3 更新店铺

```
PUT /api/v1/member/stores/:id
```

### 4.4 删除店铺

```
DELETE /api/v1/member/stores/:id
```

---

## 5. 会员标签

### 5.1 标签列表

```
GET /api/v1/member/tags
```

### 5.2 创建标签

```
POST /api/v1/member/tags
```

### 5.3 更新标签

```
PUT /api/v1/member/tags/:id
```

### 5.4 删除标签

```
DELETE /api/v1/member/tags/:id
```

---

## 6. 促销活动

### 6.1 促销活动列表

```
GET /api/v1/member/promotions
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "promotion_name": "春季大促",
        "promotion_type": "discount",
        "start_time": "2026-03-01",
        "end_time": "2026-03-31",
        "status": "active"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 6.2 创建促销活动

```
POST /api/v1/member/promotions
```

### 6.3 更新促销活动

```
PUT /api/v1/member/promotions/:id
```

### 6.4 删除促销活动

```
DELETE /api/v1/member/promotions/:id
```

---

## 7. 会员等级

### 7.1 等级列表

```
GET /api/v1/member/levels
```

### 7.2 创建等级

```
POST /api/v1/member/levels
```

### 7.3 更新等级

```
PUT /api/v1/member/levels/:id
```

### 7.4 删除等级

```
DELETE /api/v1/member/levels/:id
```

---

## 8. 积分管理

### 8.1 积分流水

```
GET /api/v1/member/points/records
```

#### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| member_id | string | 否 | 会员ID |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页条数 |

---

### 8.2 积分规则列表

```
GET /api/v1/member/points/rules
```

### 8.3 创建积分规则

```
POST /api/v1/member/points/rules
```

### 8.4 更新积分规则

```
PUT /api/v1/member/points/rules/:id
```

### 8.5 删除积分规则

```
DELETE /api/v1/member/points/rules/:id
```

---

## 9. 会员订单

### 9.1 订单列表

```
GET /api/v1/member/orders
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "order_no": "ORDER20260322001",
        "member_id": 1,
        "total_amount": 200.00,
        "points_earned": 200,
        "status": "completed",
        "created_at": "2026-03-22T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 9.2 创建订单

```
POST /api/v1/member/orders
```

### 9.3 订单详情

```
GET /api/v1/member/orders/:id
```

---

## 10. 等级调整记录

### 10.1 升级记录列表

```
GET /api/v1/member/upgrade-records
```

---

## 11. 增强功能 (member_enhanced_controller)

### 11.1 添加积分

```
POST /api/v1/members/:id/points/add
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| points | int | 是 | 积分数量（正数） |
| reason | string | 否 | 原因 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "member_id": 1,
    "points_added": 100,
    "current_balance": 5100,
    "reason": "消费返积分"
  }
}
```

---

### 11.2 扣减积分

```
POST /api/v1/members/:id/points/deduct
```

#### 请求参数 (JSON)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| points | int | 是 | 积分数量（正数） |
| reason | string | 否 | 原因 |

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "member_id": 1,
    "points_deducted": 50,
    "current_balance": 5050,
    "reason": "兑换优惠券"
  }
}
```

---

### 11.3 查询积分余额

```
GET /api/v1/members/:id/points/balance
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "member_id": 1,
    "current_balance": 5050,
    "frozen_points": 0
  }
}
```

---

### 11.4 查询积分日志

```
GET /api/v1/members/:id/points/logs
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "points": 100,
        "type": "earn",
        "reason": "消费返积分",
        "created_at": "2026-03-22T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 11.5 会员积分列表（新路径）

```
GET /api/v1/members/points
```

### 11.6 批量调整积分

```
POST /api/v1/members/points/adjust
```

### 11.7 会员优惠券列表

```
GET /api/v1/members/:id/coupons
```

---

## 12. 优惠券新路径

### 12.1 优惠券列表（新）

```
GET /api/v1/coupons
```

### 12.2 创建优惠券（新）

```
POST /api/v1/coupons
```

### 12.3 发放优惠券

```
POST /api/v1/coupons/:id/issue
```

### 12.4 使用优惠券

```
POST /api/v1/coupons/:id/use
```

---

## 13. 促销活动新路径

### 13.1 促销活动列表（新）

```
GET /api/v1/promotions
```

### 13.2 创建促销活动（新）

```
POST /api/v1/promotions
```

### 13.3 更新促销活动（新）

```
PUT /api/v1/promotions/:id
```

### 13.4 删除促销活动（新）

```
DELETE /api/v1/promotions/:id
```

### 13.5 促销活动详情（新）

```
GET /api/v1/promotions/:id
```
