# 会员管理模块 API 接口文档

**版本：** V1.0
**更新日期：** 2026-03-19
**Base URL：** `/api/v1/member`

---

## 通用说明

### 认证方式

所有接口需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

### 请求格式

- Content-Type: `application/json`
- 字符编码: `UTF-8`

### 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，0=成功，非0=失败 |
| message | string | 提示信息 |
| data | object | 返回数据（成功时） |

---

## 错误码定义

| 错误码 | 错误信息 | 说明 |
|--------|----------|------|
| 0 | success | 成功 |
| 1001 | PARAM_ERROR | 参数错误 |
| 1002 | PARAM_MISSING | 缺少必填参数 |
| 2001 | MEMBER_NOT_FOUND | 会员不存在 |
| 2002 | MEMBER_EXISTED | 会员已存在 |
| 2003 | MEMBER_DISABLED | 会员已被禁用 |
| 2004 | MEMBER_FROZEN | 会员已被冻结 |
| 2005 | LEVEL_NOT_FOUND | 会员等级不存在 |
| 2006 | LEVEL_HAS_MEMBER | 等级下存在会员，无法删除 |
| 2007 | CARD_TYPE_NOT_FOUND | 卡类型不存在 |
| 2008 | CARD_TYPE_HAS_MEMBER | 卡类型下存在会员，无法删除 |
| 3001 | COUPON_NOT_FOUND | 优惠券不存在 |
| 3002 | COUPON_STOCK_EMPTY | 优惠券库存不足 |
| 3003 | COUPON_EXPIRED | 优惠券已过期 |
| 3004 | COUPON_USED | 优惠券已使用 |
| 3005 | COUPON_NOT_OWNED | 会员未拥有该优惠券 |
| 4001 | TAG_NOT_FOUND | 标签不存在 |
| 4002 | TAG_HAS_MEMBER | 标签下存在会员，无法删除 |
| 5001 | STORE_NOT_FOUND | 门店不存在 |
| 5002 | STORE_HAS_MEMBER | 门店下存在会员，无法删除 |
| 6001 | POINTS_NOT_ENOUGH | 积分不足 |
| 6002 | POINTS_ADJUST_ERROR | 积分调整失败 |
| 7001 | ORDER_NOT_FOUND | 订单不存在 |
| 8001 | AUTH_FORBIDDEN | 无权限访问 |
| 9001 | SYSTEM_ERROR | 系统错误 |

---

## 1. 会员信息管理

### 1.1 会员列表查询

**接口描述：** 分页查询会员列表，支持多条件筛选

**请求方法：** GET  
**请求路径：** `/api/v1/member/list`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| pageSize | int | 否 | 每页条数，默认20，最大100 |
| keyword | string | 否 | 搜索关键字（手机号/姓名/会员编号） |
| levelId | int | 否 | 会员等级ID |
| status | int | 否 | 会员状态：1-正常 2-冻结 3-禁用 |
| storeId | int | 否 | 所属门店ID |
| cardTypeId | int | 否 | 卡类型ID |
| startDate | string | 否 | 注册开始日期，格式：YYYY-MM-DD |
| endDate | string | 否 | 注册结束日期，格式：YYYY-MM-DD |
| tags | string | 否 | 标签ID，多个用逗号分隔 |

**响应示例（成功）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "memberNo": "HY20260319000001",
        "name": "张三",
        "mobile": "13800138000",
        "gender": 1,
        "levelId": 2,
        "levelName": "银卡会员",
        "storeId": 1,
        "storeName": "旗舰店",
        "status": 1,
        "totalPoints": 1500,
        "totalConsume": 2500.00,
        "totalOrderCount": 5,
        "createdAt": "2026-01-15 08:00:00"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 1.2 会员详情查询

**请求方法：** GET  
**请求路径：** `/api/v1/member/detail/{id}`

**路径参数：** id - 会员ID（必填）

**响应示例（成功）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "memberNo": "HY20260319000001",
    "name": "张三",
    "mobile": "13800138000",
    "gender": 1,
    "birthday": "1990-01-01",
    "levelId": 2,
    "levelName": "银卡会员",
    "cardTypeId": 1,
    "cardTypeName": "物理卡",
    "storeId": 1,
    "storeName": "旗舰店",
    "status": 1,
    "totalPoints": 1500,
    "availablePoints": 1200,
    "totalConsume": 2500.00,
    "totalOrderCount": 5,
    "tags": [{"id": 1, "name": "高价值"}],
    "createdAt": "2026-01-15 08:00:00"
  }
}
```

---

### 1.3 会员新增

**请求方法：** POST  
**请求路径：** `/api/v1/member/create`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| mobile | string | 是 | 手机号 |
| name | string | 是 | 姓名 |
| gender | int | 否 | 性别：0-未知 1-男 2-女 |
| birthday | string | 否 | 生日，格式：YYYY-MM-DD |
| email | string | 否 | 邮箱 |
| storeId | int | 是 | 所属门店ID |
| cardTypeId | int | 否 | 卡类型ID |

**请求示例：**

```json
{
  "mobile": "13800138000",
  "name": "张三",
  "gender": 1,
  "storeId": 1
}
```

**响应示例（成功）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "memberNo": "HY20260319000001"
  }
}
```

---

### 1.4 会员编辑

**请求方法：** PUT  
**请求路径：** `/api/v1/member/update/{id}`

**路径参数：** id - 会员ID（必填）

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 姓名 |
| gender | int | 否 | 性别 |
| birthday | string | 否 | 生日 |
| email | string | 否 | 邮箱 |
| cardTypeId | int | 否 | 卡类型ID |
| storeId | int | 否 | 所属门店ID |
| remark | string | 否 | 备注 |

---

### 1.5 会员删除

**请求方法：** DELETE  
**请求路径：** `/api/v1/member/delete/{id}`

**路径参数：** id - 会员ID（必填）

**响应示例：**

```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

---

### 1.6 会员状态变更

**请求方法：** PUT  
**请求路径：** `/api/v1/member/status/{id}`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | int | 是 | 状态：1-正常 2-冻结 3-禁用 |
| reason | string | 否 | 变更原因 |

---

### 1.7 会员等级调整

**请求方法：** PUT  
**请求路径：** `/api/v1/member/level/{id}`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| levelId | int | 是 | 目标等级ID |
| reason | string | 是 | 调整原因：upgrade/downgrade/vip/other |
| remark | string | 否 | 备注说明 |

---

## 2. 会员卡类型管理

### 2.1 卡类型列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/card-type/list`

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "物理卡",
      "discountRate": 0.95,
      "pointsRate": 1.5,
      "description": "实体会员卡",
      "status": 1,
      "memberCount": 100
    }
  ]
}
```

---

### 2.2 卡类型新增

**请求方法：** POST  
**请求路径：** `/api/v1/member/card-type/create`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 卡类型名称 |
| discountRate | decimal | 是 | 折扣率（0-1） |
| pointsRate | decimal | 是 | 积分倍率 |
| description | string | 否 | 描述 |

---

### 2.3 卡类型编辑

**请求方法：** PUT  
**请求路径：** `/api/v1/member/card-type/update/{id}`

---

### 2.4 卡类型删除

**请求方法：** DELETE  
**请求路径：** `/api/v1/member/card-type/delete/{id}`

---

## 3. 会员等级管理

### 3.1 等级列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/level/list`

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "普通会员",
      "code": "NORMAL",
      "minAmount": 0,
      "maxAmount": 999,
      "discountRate": 1.0,
      "pointsRate": 1.0,
      "benefits": ["新人礼包"],
      "memberCount": 1000
    },
    {
      "id": 2,
      "name": "银卡会员",
      "code": "SILVER",
      "minAmount": 1000,
      "maxAmount": 4999,
      "discountRate": 0.95,
      "pointsRate": 1.5,
      "benefits": ["95折优惠", "专属客服"],
      "memberCount": 500
    },
    {
      "id": 3,
      "name": "金卡会员",
      "code": "GOLD",
      "minAmount": 5000,
      "maxAmount": 9999,
      "discountRate": 0.9,
      "pointsRate": 2.0,
      "benefits": ["9折优惠", "专属客服", "优先发货"],
      "memberCount": 200
    },
    {
      "id": 4,
      "name": "钻石会员",
      "code": "DIAMOND",
      "minAmount": 10000,
      "maxAmount": null,
      "discountRate": 0.85,
      "pointsRate": 3.0,
      "benefits": ["85折优惠", "专属客服", "优先发货", "专属活动", "免费配送"],
      "memberCount": 50
    }
  ]
}
```

---

### 3.2 等级新增

**请求方法：** POST  
**请求路径：** `/api/v1/member/level/create`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 等级名称 |
| code | string | 是 | 等级编码 |
| minAmount | decimal | 是 | 升级门槛金额 |
| maxAmount | decimal | 否 | 最高门槛，为空表示无上限 |
| discountRate | decimal | 是 | 折扣率 |
| pointsRate | decimal | 是 | 积分倍率 |
| description | string | 否 | 描述 |
| benefits | array | 否 | 权益列表 |
| sort | int | 否 | 排序值 |

---

### 3.3 等级编辑

**请求方法：** PUT  
**请求路径：** `/api/v1/member/level/update/{id}`

---

### 3.4 等级删除

**请求方法：** DELETE  
**请求路径：** `/api/v1/member/level/delete/{id}`

---

### 3.5 升级规则设置

**请求方法：** PUT  
**请求路径：** `/api/v1/member/level/upgrade-rules`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| autoUpgrade | boolean | 是 | 是否开启自动升级 |
| autoDowngrade | boolean | 是 | 是否开启自动降级 |
| downgradeCycle | int | 否 | 降级周期（月） |
| downgradeNoticeDays | int | 否 | 降级预警提前天数 |

---

## 4. 优惠券管理

### 4.1 优惠券列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/coupon/list`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| pageSize | int | 否 | 每页条数 |
| type | int | 否 | 类型：1-满减券 2-折扣券 3-兑换券 |
| status | int | 否 | 状态：1-未开始 2-进行中 3-已结束 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "新人100元券",
        "type": 1,
        "discountValue": 100,
        "minAmount": 500,
        "totalStock": 10000,
        "remainStock": 8500,
        "status": 2,
        "usedCount": 1500
      }
    ],
    "total": 10
  }
}
```

---

### 4.2 优惠券新增

**请求方法：** POST  
**请求路径：** `/api/v1/member/coupon/create`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 优惠券名称 |
| type | int | 是 | 类型：1-满减券 2-折扣券 3-兑换券 |
| discountType | int | 是 | 优惠方式：1-满减 2-折扣 |
| discountValue | decimal | 是 | 优惠值 |
| minAmount | decimal | 否 | 使用门槛 |
| totalStock | int | 是 | 总库存 |
| perLimit | int | 否 | 每人限领数量 |
| validType | int | 是 | 有效期类型：1-固定日期 2-领取后N天 |
| validDays | int | 否 | 有效天数 |
| startTime | string | 否 | 开始时间 |
| endTime | string | 否 | 结束时间 |
| description | string | 否 | 使用说明 |

---

### 4.3 优惠券编辑

**请求方法：** PUT  
**请求路径：** `/api/v1/member/coupon/update/{id}`

---

### 4.4 优惠券删除

**请求方法：** DELETE  
**请求路径：** `/api/v1/member/coupon/delete/{id}`

---

### 4.5 优惠券发放

**请求方法：** POST  
**请求路径：** `/api/v1/member/coupon/grant`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| couponId | int | 是 | 优惠券ID |
| memberIds | array | 是 | 会员ID列表（最多100个） |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "successCount": 5,
    "failCount": 0
  }
}
```

---

### 4.6 会员优惠券列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/coupon/member-list`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| memberId | int | 是 | 会员ID |
| status | int | 否 | 状态：0-未使用 1-已使用 2-已过期 |

---

### 4.7 优惠券使用记录

**请求方法：** GET  
**请求路径：** `/api/v1/member/coupon/use-record`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| couponId | int | 否 | 优惠券ID |
| memberId | int | 否 | 会员ID |
| startDate | string | 否 | 开始日期 |
| endDate | string | 否 | 结束日期 |
| page | int | 否 | 页码 |
| pageSize | int | 否 | 每页条数 |

---

## 5. 会员标签管理

### 5.1 标签列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/tag/list`

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "高价值",
      "type": 1,
      "typeName": "属性标签",
      "color": "#FF6B6B",
      "memberCount": 100
    },
    {
      "id": 2,
      "name": "活跃用户",
      "type": 2,
      "typeName": "行为标签",
      "color": "#4ECDC4",
      "memberCount": 500
    }
  ]
}
```

**标签类型：** 1-属性标签 2-行为标签 3-自定义标签

---

### 5.2 标签新增

**请求方法：** POST  
**请求路径：** `/api/v1/member/tag/create`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 标签名称 |
| type | int | 是 | 标签类型：1-属性 2-行为 3-自定义 |
| color | string | 否 | 标签颜色 |
| description | string | 否 | 描述 |

---

### 5.3 标签编辑

**请求方法：** PUT  
**请求路径：** `/api/v1/member/tag/update/{id}`

---

### 5.4 标签删除

**请求方法：** DELETE  
**请求路径：** `/api/v1/member/tag/delete/{id}`

---

### 5.5 批量打标

**请求方法：** POST  
**请求路径：** `/api/v1/member/tag/batch-tag`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| tagId | int | 是 | 标签ID |
| memberFilters | object | 是 | 会员筛选条件 |
| memberFilters.levelIds | array | 否 | 等级ID列表 |
| memberFilters.minConsume | decimal | 否 | 最低消费 |
| memberFilters.maxConsume | decimal | 否 | 最高消费 |
| memberFilters.storeIds | array | 否 | 门店ID列表 |

---

### 5.6 批量移除标签

**请求方法：** POST  
**请求路径：** `/api/v1/member/tag/batch-untag`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| tagId | int | 是 | 标签ID |
| memberIds | array | 是 | 会员ID列表 |

---

## 6. 促销活动管理

### 6.1 活动列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/promotion/list`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| pageSize | int | 否 | 每页条数 |
| type | int | 否 | 活动类型：1-买赠 2-直减 3-满额减 4-满额折 |
| status | int | 否 | 状态：1-未开始 2-进行中 3-已结束 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "春季大促",
        "type": 3,
        "typeName": "满额减",
        "rules": [{"threshold": 200, "discount": 20}],
        "startTime": "2026-03-01 00:00:00",
        "endTime": "2026-03-31 23:59:59",
        "status": 2
      }
    ]
  }
}
```

---

### 6.2 活动新增

**请求方法：** POST  
**请求路径：** `/api/v1/member/promotion/create`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 活动名称 |
| type | int | 是 | 活动类型：1-买赠 2-直减 3-满额减 4-满额折 |
| rules | array | 是 | 活动规则 |
| rules[].threshold | decimal | 条件 | 门槛金额（满减/满折时） |
| rules[].discount | decimal | 是 | 优惠值 |
| rules[].buyProductId | int | 条件 | 购买商品ID（买赠时） |
| rules[].giftProductId | int | 条件 | 赠品ID（买赠时） |
| startTime | string | 是 | 开始时间 |
| endTime | string | 是 | 结束时间 |
| applicableStores | array | 否 | 适用门店 |
| memberLevels | array | 否 | 适用会员等级 |
| description | string | 否 | 说明 |

---

### 6.3 活动编辑

**请求方法：** PUT  
**请求路径：** `/api/v1/member/promotion/update/{id}`

---

### 6.4 活动删除

**请求方法：** DELETE  
**请求路径：** `/api/v1/member/promotion/delete/{id}`

---

## 7. 会员积分管理

### 7.1 积分规则设置

**请求方法：** PUT  
**请求路径：** `/api/v1/member/points/rules`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| baseRate | decimal | 是 | 基础积分倍率（1元=N积分） |
| pointsToMoney | int | 是 | 抵扣比例（100积分=1元） |
| maxDeductRate | decimal | 否 | 单笔最高抵扣比例，默认0.5 |
| minPointsToUse | int | 否 | 最低抵扣积分，默认100 |
| birthdayDouble | boolean | 否 | 生日是否双倍 |
| signInEnabled | boolean | 否 | 是否开启签到积分 |
| signInPoints | int | 否 | 签到积分 |

---

### 7.2 积分规则查询

**请求方法：** GET  
**请求路径：** `/api/v1/member/points/rules`

---

### 7.3 积分流水查询

**请求方法：** GET  
**请求路径：** `/api/v1/member/points/flow`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| memberId | int | 是 | 会员ID |
| type | int | 否 | 类型：1-获取 2-抵扣 3-调整 |
| startDate | string | 否 | 开始日期 |
| endDate | string | 否 | 结束日期 |
| page | int | 否 | 页码 |
| pageSize | int | 否 | 每页条数 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "type": 1,
        "typeName": "获取",
        "points": 500,
        "balance": 1500,
        "source": "order",
        "description": "订单消费获得积分",
        "createdAt": "2026-03-15 10:30:00"
      }
    ]
  }
}
```

---

### 7.4 积分调整

**请求方法：** POST  
**请求路径：** `/api/v1/member/points/adjust`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| memberId | int | 是 | 会员ID |
| points | int | 是 | 调整积分（正数增加，负数扣减） |
| reason | string | 是 | 调整原因 |
| remark | string | 否 | 备注 |

**调整原因：** complaint_compensation/activity_reward/birthday_gift/error_correction/violation_punish/other

---

### 7.5 会员积分查询

**请求方法：** GET  
**请求路径：** `/api/v1/member/points/balance/{memberId}`

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "memberId": 1,
    "totalPoints": 1500,
    "availablePoints": 1200,
    "frozenPoints": 300
  }
}
```

---

## 8. 店铺管理

### 8.1 店铺列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/store/list`

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "旗舰店",
      "code": "STORE_001",
      "province": "广东省",
      "city": "深圳市",
      "district": "南山区",
      "address": "科技园路100号",
      "longitude": 113.953345,
      "latitude": 22.538567,
      "contact": "0755-12345678",
      "managerName": "李经理",
      "managerMobile": "13800138000",
      "defaultCardTypeId": 1,
      "status": 1,
      "memberCount": 500
    }
  ]
}
```

---

### 8.2 店铺新增

**请求方法：** POST  
**请求路径：** `/api/v1/member/store/create`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 门店名称 |
| code | string | 是 | 门店编码 |
| province | string | 是 | 省份 |
| city | string | 是 | 城市 |
| district | string | 否 | 区县 |
| address | string | 是 | 详细地址 |
| longitude | decimal | 否 | 经度 |
| latitude | decimal | 否 | 纬度 |
| contact | string | 否 | 联系电话 |
| managerName | string | 否 | 店长姓名 |
| managerMobile | string | 否 | 店长手机号 |
| defaultCardTypeId | int | 否 | 默认卡类型ID |

---

### 8.3 店铺编辑

**请求方法：** PUT  
**请求路径：** `/api/v1/member/store/update/{id}`

---

### 8.4 店铺删除

**请求方法：** DELETE  
**请求路径：** `/api/v1/member/store/delete/{id}`

---

## 9. 会员订单管理

### 9.1 订单列表

**请求方法：** GET  
**请求路径：** `/api/v1/member/order/list`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| memberId | int | 否 | 会员ID |
| storeId | int | 否 | 门店ID |
| status | int | 否 | 状态：1-待支付 2-已支付 3-已完成 4-已取消 |
| startDate | string | 否 | 开始日期 |
| endDate | string | 否 | 结束日期 |
| page | int | 否 | 页码 |
| pageSize | int | 否 | 每页条数 |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "orderNo": "ORD20260319000001",
        "memberId": 1,
        "memberName": "张三",
        "storeId": 1,
        "storeName": "旗舰店",
        "totalAmount": 500.00,
        "discountAmount": 50.00,
        "payAmount": 440.00,
        "pointsEarned": 440,
        "status": 3,
        "statusName": "已完成",
        "createdAt": "2026-03-19 14:00:00"
      }
    ]
  }
}
```

---

### 9.2 订单详情

**请求方法：** GET  
**请求路径：** `/api/v1/member/order/detail/{id}`

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "orderNo": "ORD20260319000001",
    "memberId": 1,
    "memberName": "张三",
    "storeId": 1,
    "storeName": "旗舰店",
    "items": [
      {
        "productId": 1001,
        "productName": "商品A",
        "quantity": 2,
        "price": 200.00,
        "subtotal": 400.00
      }
    ],
    "totalAmount": 500.00,
    "discountAmount": 50.00,
    "payAmount": 440.00,
    "pointsEarned": 440,
    "status": 3,
    "payTime": "2026-03-19 14:30:00",
    "createdAt": "2026-03-19 14:00:00"
  }
}
```

---

### 9.3 订单状态变更

**请求方法：** PUT  
**请求路径：** `/api/v1/member/order/status/{id}`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | int | 是 | 目标状态：1-待支付 2-已支付 3-已完成 4-已取消 |
| reason | string | 否 | 变更原因（取消时必填） |

---

## 附录：错误码速查表

| 模块 | 错误码范围 | 说明 |
|------|------------|------|
| 通用 | 1001-1099 | 参数错误 |
| 会员 | 2001-2099 | 会员相关错误 |
| 优惠券 | 3001-3099 | 优惠券相关错误 |
| 标签 | 4001-4099 | 标签相关错误 |
| 门店 | 5001-5099 | 门店相关错误 |
| 积分 | 6001-6099 | 积分相关错误 |
| 订单 | 7001-7099 | 订单相关错误 |
| 权限 | 8001-8099 | 权限相关错误 |
| 系统 | 9001-9099 | 系统相关错误 |
