# 会员管理模块差距分析
**日期：** 2026-03-29
**分析依据：** MODULE_MEMBER_MANAGEMENT.md (PRD V1.4)

---

## 一、PRD要求 vs 实际实现对比

### 1.1 数据模型

| PRD要求 | 状态 | 说明 |
|---------|------|------|
| members | ✅ | 已有Member模型 |
| member_cards | ✅ | 已有MemberCard模型 |
| member_levels | ✅ | 已有MemberLevel模型 |
| member_upgrade_rules | ❌ | **缺失** - 无升级规则表 |
| coupons | ✅ | 已有Coupon模型 |
| points_rules | ✅ | 已有PointsRule模型 |
| member_points_records | ✅ | 已有MemberPointsRecord模型 |
| member_orders | ✅ | 已有MemberOrder模型 |

### 1.2 API接口

| PRD要求 | 路径 | 状态 | 说明 |
|---------|------|------|------|
| 会员列表 | GET /api/v1/members | ✅ | |
| 会员详情 | GET /api/v1/members/:id | ✅ | |
| 创建会员 | POST /api/v1/members | ✅ | |
| 更新会员 | PUT /api/v1/members/:id | ✅ | |
| 删除会员 | DELETE /api/v1/members/:id | ✅ | |
| 等级列表 | GET /api/v1/member/levels | ✅ | |
| 创建等级 | POST /api/v1/member/levels | ✅ | |
| 会员卡列表 | GET /api/v1/member/cards | ✅ | |
| 创建会员卡 | POST /api/v1/member/cards | ✅ | |
| 优惠券列表 | GET /api/v1/member/coupons | ✅ | |
| 创建优惠券 | POST /api/v1/member/coupons | ✅ | |
| **发放优惠券** | POST /api/v1/member/coupons/:id/grant | ❌ | **缺失** |
| **核销优惠券** | POST /api/v1/member/coupon-grants/:id/use | ❌ | **缺失** |
| 积分规则列表 | GET /api/v1/member/points/rules | ✅ | |
| 积分流水 | GET /api/v1/member/points/records | ✅ | |
| **积分调整** | POST /api/v1/members/:id/points/adjust | ❌ | **缺失** |
| 促销列表 | GET /api/v1/member/promotions | ✅ | |
| 创建促销 | POST /api/v1/member/promotions | ✅ | |
| 订单列表 | GET /api/v1/member/orders | ✅ | |
| 创建订单 | POST /api/v1/member/orders | ✅ | |
| 店铺列表 | GET /api/v1/member/stores | ✅ | |
| 创建店铺 | POST /api/v1/member/stores | ✅ | |
| 标签列表 | GET /api/v1/member/tags | ✅ | |
| 创建标签 | POST /api/v1/member/tags | ✅ | |
| **批量打标** | POST /api/v1/members/:id/tags | ❌ | **缺失** |

### 1.3 前端视图

| PRD要求 | 文件 | 状态 | 说明 |
|---------|------|------|------|
| 会员列表 | MemberListView.vue | ✅ | |
| 会员卡管理 | MemberCards.vue | ❌ | **缺失** |
| 优惠券管理 | CouponCenter.vue, CouponsView.vue | ✅ | |
| 积分管理 | PointsRulesView.vue, MemberPoints.vue | ✅ | |
| 会员等级 | MemberLevels.vue | ✅ | |
| **升级规则** | MemberUpgradeRules.vue | ❌ | **缺失** |
| 会员标签 | MemberTags.vue | ✅ | |
| 促销活动 | PromotionsView.vue | ✅ | |
| 订单管理 | OrdersView.vue, MemberOrders.vue | ✅ | |
| **店铺管理** | MemberStores.vue | ❌ | **缺失** |

---

## 二、缺失功能详细说明

### 2.1 后端缺失

#### 2.1.1 member_upgrade_rules 升级规则表
```sql
CREATE TABLE member_upgrade_rules (
    id UINT PRIMARY KEY,
    from_level UINT,
    to_level UINT,
    points_threshold INT64,
    amount_threshold FLOAT64,
    points_reward INT64,
    status INT
);
```

#### 2.1.2 API缺失

1. **发放优惠券** - 需要实现批量发放给会员
   ```go
   // POST /api/v1/member/coupons/:id/grant
   func GrantCoupon(ctx *gin.Context) // 发放优惠券给指定会员
   ```

2. **核销优惠券** - 需要实现核销接口
   ```go
   // POST /api/v1/member/coupon-grants/:id/use
   func UseCoupon(ctx *gin.Context) // 核销已发放的优惠券
   ```

3. **积分调整** - 需要实现积分手动调整
   ```go
   // POST /api/v1/members/:id/points/adjust
   // Body: { "points": 100, "points_type": 1, "remark": "手动调整" }
   ```

4. **批量打标** - 需要实现批量打标签
   ```go
   // POST /api/v1/members/:id/tags
   // Body: { "tag_ids": [1,2,3] }
   ```

### 2.2 前端缺失

#### 2.2.1 MemberCards.vue - 会员卡管理页面
**功能要求：**
- 会员卡列表（table）
- 新建会员卡（drawer/form）
- 编辑会员卡
- 删除会员卡
- 卡类型：储值卡/积分卡/打折卡

**UI参考：** 见PRD 8.2会员卡管理页面

#### 2.2.2 MemberUpgradeRules.vue - 升级规则页面
**功能要求：**
- 升级规则列表
- 新建升级规则（from_level → to_level）
- 设置积分阈值/消费金额阈值
- 设置升级赠送积分

#### 2.2.3 MemberStores.vue - 会员店铺管理页面
**功能要求：**
- 店铺列表
- 新建店铺
- 编辑店铺
- 区域管理

---

## 三、补全计划

### Sprint X: 会员管理补全

#### Phase 1: 后端补全 (agenthd)
| 任务 | 优先级 | 工作量 |
|------|--------|--------|
| 创建 member_upgrade_rules 模型 | P0 | 30min |
| 实现发放优惠券 API | P0 | 1h |
| 实现核销优惠券 API | P0 | 1h |
| 实现积分调整 API | P1 | 1h |
| 实现批量打标 API | P1 | 1h |

#### Phase 2: 前端补全 (agentqd)
| 任务 | 优先级 | 工作量 |
|------|--------|--------|
| MemberCards.vue | P0 | 2h |
| MemberUpgradeRules.vue | P1 | 2h |
| MemberStores.vue | P2 | 2h |

---

## 四、验收标准

### P0 验收
- [ ] 会员CRUD完整
- [ ] 会员等级CRUD完整
- [ ] 会员卡CRUD完整
- [ ] 优惠券CRUD + 发放 + 核销

### P1 验收
- [ ] 积分规则 + 流水完整
- [ ] 促销CRUD完整
- [ ] 订单CRUD完整
- [ ] 批量打标功能

### P2 验收
- [ ] 升级规则完整
- [ ] 店铺管理完整
