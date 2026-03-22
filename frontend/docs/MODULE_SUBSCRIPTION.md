# 模块 PRD：订阅管理（MODULE_SUBSCRIPTION）

**版本：** V1.0
**所属Phase：** Phase 2（Sprint 15-16）
**优先级：** P1
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

订阅管理模块是MDM控制中台的商业化核心，负责管理用户的订阅计划、计费、支付、发票等商业化能力，支持多级订阅、自动续费、用量计费、Webhook事件通知，为产品商业化提供完整的基础设施。

### 1.2 核心价值

- **商业变现**：支持多种商业模式（订阅+按量）
- **自动化运营**：续费/升级/降级自动化
- **合规透明**：发票/账单清晰可查

### 1.3 范围边界

**包含：**
- 多级订阅计划管理
- 自动续费
- 订阅变更（升级/降级/取消）
- 用量计费
- API配额计费
- 发票管理
- 账单管理
- Webhook事件系统

**不包含：**
- 优惠券促销（会员管理模块）
- 支付网关集成（支付平台）
- 开放API市场（MODULE_PLATFORM_ECOSYSTEM）

---

## 二、功能详情

### 2.1 订阅计划管理

#### 2.1.1 订阅等级

| 等级 | 名称 | 定价 | 说明 |
|------|------|------|------|
| Free | 免费版 | ¥0 | 基础功能，设备数限制 |
| Basic | 基础版 | ¥29/月 | 1台设备，基础功能 |
| Pro | 专业版 | ¥99/月 | 3台设备，高级功能 |
| Enterprise | 企业版 | ¥299/月 | 10台设备，全功能 |
| Unlimited | 无限版 | ¥999/月 | 无限制设备数 |

#### 2.1.2 订阅权益

| 权益项 | Free | Basic | Pro | Enterprise | Unlimited |
|--------|------|-------|-----|------------|-----------|
| 设备数 | 1 | 1 | 3 | 10 | 无限制 |
| AI对话次数 | 50/天 | 200/天 | 1000/天 | 无限制 | 无限制 |
| 历史记录 | 7天 | 30天 | 1年 | 永久 | 永久 |
| 情感分析 | ❌ | ✅ | ✅ | ✅ | ✅ |
| 数字孪生 | ❌ | ❌ | ✅ | ✅ | ✅ |
| API访问 | ❌ | ❌ | ❌ | ✅ | ✅ |
| 优先客服 | ❌ | ❌ | ✅ | ✅ | ✅ |
| 数据导出 | ❌ | ❌ | 月度 | 实时 | 实时 |

#### 2.1.3 订阅配置

| 功能 | 说明 |
|------|------|
| 计划CRUD | 创建/编辑/下线订阅计划 |
| 功能开关 | 按订阅等级控制功能开关 |
| 配额配置 | 配置各等级配额限制 |
| 试用配置 | 支持免费试用配置 |
| 推荐配置 | 设置推荐展示的订阅计划 |

### 2.2 自动续费

#### 2.2.1 续费策略

| 功能 | 说明 |
|------|------|
| 自动续费 | 到期前自动扣款续费 |
| 续费提醒 | 到期前7/3/1天提醒 |
| 宽限期 | 到期后7天宽限期 |
| 续费优惠 | 续费可享受折扣 |
| 手动续费 | 支持用户手动续费 |

#### 2.2.2 续费流程

| 阶段 | 触发时间 | 操作 |
|------|----------|------|
| 续费提醒 | 到期前7天 | 发送续费提醒通知 |
| 续费提醒 | 到期前3天 | 发送续费提醒通知 |
| 续费提醒 | 到期前1天 | 发送续费提醒通知 |
| 自动续费 | 到期日 | 尝试自动扣款 |
| 功能降级 | 到期日+0 | 功能降级到Free |
| 宽限期 | 到期后7天 | 保持功能，宽限期计费 |
| 账号暂停 | 到期后7天 | 账号功能暂停 |
| 账号注销 | 到期后30天 | 注销账号（可配置） |

#### 2.2.3 支付方式

| 支付方式 | 说明 |
|----------|------|
| 支付宝 | 支付宝自动扣款 |
| 微信支付 | 微信自动扣款 |
| 信用卡 | Visa/Mastercard自动扣款 |
| Apple Pay | iOS内购 |
| Google Pay | Android内购 |

### 2.3 订阅变更

#### 2.3.1 升级

| 功能 | 说明 |
|------|------|
| 即时升级 | 升级立即生效 |
| 差价计算 | 按剩余时间计算差价 |
| 升级优惠 | 首次升级享受优惠 |
| 升级历史 | 记录升级历史 |

#### 2.3.2 降级

| 功能 | 说明 |
|------|------|
| 到期降级 | 降级到下个周期生效 |
| 立即降级 | 可选择立即降级（不退款） |
| 降级限制 | 部分功能不支持降级 |
| 数据保留 | 降级后超出配额数据保留30天 |

#### 2.3.3 取消

| 功能 | 说明 |
|------|------|
| 取消订阅 | 取消自动续费 |
| 取消原因 | 收集取消原因 |
| 取消优惠 | 取消时可提供挽留优惠 |
| 取消确认 | 确认取消操作 |

### 2.4 用量计费

#### 2.4.1 用量类型

| 用量类型 | 计费单位 | 定价 |
|----------|----------|------|
| AI对话次数 | 次 | ¥0.01/次（超出配额后） |
| 存储空间 | GB/月 | ¥0.5/GB/月 |
| API调用次数 | 次 | ¥0.001/次（开发者API） |
| 高级表情包 | 个 | ¥1-5/个 |
| 高级动作 | 个 | ¥2-10/个 |

#### 2.4.2 用量监控

| 功能 | 说明 |
|------|------|
| 用量实时 | 实时展示各类型用量 |
| 用量预测 | 预测本月用量 |
| 用量预警 | 用量超过80%预警 |
| 用量排名 | 用量最大的用户/设备 |

#### 2.4.3 用量账单

| 功能 | 说明 |
|------|------|
| 实时账单 | 实时展示当月费用 |
| 账单明细 | 各用量类型明细 |
| 费用分摊 | 多设备费用分摊 |

### 2.5 API配额计费

#### 2.5.1 配额配置

| 功能 | 说明 |
|------|------|
| 配额等级 | 按订阅等级配置API配额 |
| 配额刷新 | 按分钟/小时/天刷新 |
| 配额预留 | 支持预留额外配额 |
| 配额预警 | 配额使用超过阈值预警 |

#### 2.5.2 配额管理

| 功能 | 说明 |
|------|------|
| 配额查询 | 查询当前配额使用情况 |
| 配额调整 | 管理员调整用户配额 |
| 配额购买 | 用户购买额外配额 |
| 配额转让 | 企业内配额转让 |

### 2.6 发票管理

#### 2.6.1 发票类型

| 类型 | 说明 |
|------|------|
| 增值税普通发票 | 个人用户 |
| 增值税专用发票 | 企业用户 |
| 电子发票 | 所有用户 |

#### 2.6.2 发票申请

| 功能 | 说明 |
|------|------|
| 发票抬头 | 管理发票抬头信息 |
| 发票资质 | 审核企业发票资质 |
| 发票申请 | 申请开票 |
| 发票历史 | 历史发票记录 |
| 发票下载 | 下载/发送电子发票 |

#### 2.6.3 自动开票

| 功能 | 说明 |
|------|------|
| 自动开票 | 每月自动生成发票 |
| 开票阈值 | 满一定金额自动开票 |
| 合并开票 | 多笔订单合并开票 |

### 2.7 账单管理

#### 2.7.1 账单类型

| 类型 | 说明 |
|------|------|
| 订阅账单 | 订阅费用账单 |
| 用量账单 | 超额用量账单 |
| 综合账单 | 综合月度账单 |

#### 2.7.2 账单功能

| 功能 | 说明 |
|------|------|
| 账单列表 | 月度账单列表 |
| 账单详情 | 账单明细 |
| 账单支付 | 账单支付 |
| 账单导出 | 导出账单PDF/Excel |

### 2.8 Webhook事件系统

#### 2.8.1 事件类型

| 事件类型 | 说明 | 触发时机 |
|----------|------|----------|
| subscription.created | 订阅创建 | 新订阅创建时 |
| subscription.renewed | 订阅续费 | 续费成功时 |
| subscription.upgraded | 订阅升级 | 升级成功时 |
| subscription.downgraded | 订阅降级 | 降级成功时 |
| subscription.cancelled | 订阅取消 | 取消续费时 |
| subscription.expired | 订阅过期 | 到期时 |
| usage.threshold_reached | 用量达标 | 用量超过阈值时 |
| invoice.created | 账单创建 | 新账单生成时 |
| invoice.paid | 账单支付 | 账单支付成功时 |
| payment.failed | 支付失败 | 支付失败时 |

#### 2.8.2 Webhook配置

| 功能 | 说明 |
|------|------|
| Webhook创建 | 创建Webhook订阅 |
| 事件选择 | 选择订阅的事件类型 |
| 签名验证 | 签名密钥验证 |
| 重试机制 | 失败自动重试 |
| 事件日志 | 事件发送历史 |

---

## 三、API接口定义

### 3.1 订阅管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/subscriptions/plans | 订阅计划列表 |
| GET | /api/v1/subscriptions/plans/:id | 计划详情 |
| GET | /api/v1/subscriptions/current | 当前订阅 |
| POST | /api/v1/subscriptions/create | 创建订阅 |
| POST | /api/v1/subscriptions/upgrade | 升级订阅 |
| POST | /api/v1/subscriptions/downgrade | 降级订阅 |
| POST | /api/v1/subscriptions/cancel | 取消订阅 |
| POST | /api/v1/subscriptions/resume | 恢复订阅 |
| POST | /api/v1/subscriptions/renew | 手动续费 |

### 3.2 用量计费

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/usage/current | 当前用量 |
| GET | /api/v1/usage/history | 用量历史 |
| GET | /api/v1/usage/prediction | 用量预测 |
| GET | /api/v1/usage/quotas | 配额查询 |
| POST | /api/v1/usage/quotas/purchase | 购买配额 |

### 3.3 支付

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/payments/pay | 发起支付 |
| GET | /api/v1/payments/:id | 支付详情 |
| POST | /api/v1/payments/:id/refund | 申请退款 |

### 3.4 发票

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/invoices | 发票列表 |
| GET | /api/v1/invoices/:id | 发票详情 |
| POST | /api/v1/invoices | 申请开票 |
| POST | /api/v1/invoices/:id/download | 下载发票 |
| POST | /api/v1/invoices/titles | 添加发票抬头 |

### 3.5 账单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/billing/statements | 账单列表 |
| GET | /api/v1/billing/statements/:id | 账单详情 |
| GET | /api/v1/billing/summary | 账单汇总 |
| GET | /api/v1/billing/export | 导出账单 |

### 3.6 Webhook

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/webhooks | Webhook列表 |
| POST | /api/v1/webhooks | 创建Webhook |
| GET | /api/v1/webhooks/:id | Webhook详情 |
| PUT | /api/v1/webhooks/:id | 更新Webhook |
| DELETE | /api/v1/webhooks/:id | 删除Webhook |
| GET | /api/v1/webhooks/:id/logs | 事件日志 |
| POST | /api/v1/webhooks/:id/test | 测试Webhook |

---

## 四、数据库设计

### 4.1 订阅计划表 (subscription_plans)

```sql
CREATE TABLE subscription_plans (
    id              BIGSERIAL PRIMARY KEY,
    plan_code       VARCHAR(50) NOT NULL UNIQUE,
    plan_name       VARCHAR(100) NOT NULL,
    plan_type       VARCHAR(20) NOT NULL,         -- 'free'/'paid'
    price_monthly   DECIMAL(10,2) DEFAULT 0,
    price_yearly    DECIMAL(10,2),
    features        JSONB,                          -- 权益配置
    quotas          JSONB,                          -- 配额配置
    is_active       BOOLEAN DEFAULT TRUE,
    is_recommended  BOOLEAN DEFAULT FALSE,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.2 用户订阅表 (user_subscriptions)

```sql
CREATE TABLE user_subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    plan_id         BIGINT NOT NULL REFERENCES subscription_plans(id),
    status          VARCHAR(20) NOT NULL,          -- 'active'/'cancelled'/'expired'/'trial'
    started_at      TIMESTAMP NOT NULL,
    expires_at      TIMESTAMP NOT NULL,
    next_billing_at TIMESTAMP,
    cancelled_at    TIMESTAMP,
    cancel_reason   TEXT,
    trial_ends_at   TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id)
);
```

### 4.3 订阅变更记录表 (subscription_changes)

```sql
CREATE TABLE subscription_changes (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    change_type     VARCHAR(20) NOT NULL,         -- 'upgrade'/'downgrade'/'renew'/'cancel'/'create'
    from_plan_id    BIGINT REFERENCES subscription_plans(id),
    to_plan_id      BIGINT NOT NULL REFERENCES subscription_plans(id),
    amount          DECIMAL(10,2),
    change_reason   TEXT,
    effective_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.4 用量记录表 (usage_records)

```sql
CREATE TABLE usage_records (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    usage_type      VARCHAR(30) NOT NULL,         -- 'ai_chat'/'storage'/'api_call'/'emotion_pack'/'action_pack'
    usage_count     BIGINT DEFAULT 0,
    quota_limit     BIGINT,
    period_start    TIMESTAMP NOT NULL,
    period_end      TIMESTAMP NOT NULL,
    cost            DECIMAL(10,2) DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_usage_records_user_period ON usage_records(user_id, usage_type, period_start DESC);
```

### 4.5 配额表 (api_quotas)

```sql
CREATE TABLE api_quotas (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    quota_type      VARCHAR(30) NOT NULL,
    quota_limit     BIGINT NOT NULL,
    quota_used      BIGINT DEFAULT 0,
    quota_reserved  BIGINT DEFAULT 0,
    reset_interval  VARCHAR(20),                  -- 'minute'/'hour'/'day'/'month'
    last_reset_at  TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, quota_type)
);
```

### 4.6 账单表 (billing_statements)

```sql
CREATE TABLE billing_statements (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    statement_no    VARCHAR(50) NOT NULL UNIQUE,
    period_start    DATE NOT NULL,
    period_end      DATE NOT NULL,
    subtotal        DECIMAL(10,2) NOT NULL,
    discount        DECIMAL(10,2) DEFAULT 0,
    total           DECIMAL(10,2) NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',  -- 'pending'/'paid'/'overdue'
    due_date        DATE,
    paid_at         TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_billing_statements_user ON billing_statements(user_id, period_start DESC);
```

### 4.7 发票表 (invoices)

```sql
CREATE TABLE invoices (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    invoice_no      VARCHAR(50) NOT NULL UNIQUE,
    invoice_type    VARCHAR(20) NOT NULL,          -- 'normal'/'special'/'electronic'
    title           VARCHAR(255) NOT NULL,
    tax_number      VARCHAR(50),
    amount          DECIMAL(10,2) NOT NULL,
    tax_amount      DECIMAL(10,2),
    status          VARCHAR(20) DEFAULT 'pending', -- 'pending'/'issued'/'cancelled'
    issued_at       TIMESTAMP,
    statement_ids   BIGINT[],
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_invoices_user ON invoices(user_id, created_at DESC);
```

### 4.8 Webhook配置表 (webhooks)

```sql
CREATE TABLE webhooks (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    webhook_name    VARCHAR(255) NOT NULL,
    endpoint_url    VARCHAR(500) NOT NULL,
    secret_key      VARCHAR(255),
    events          VARCHAR(50)[],                -- 订阅的事件类型
    is_active       BOOLEAN DEFAULT TRUE,
    retry_count     INT DEFAULT 3,
    retry_interval  INT DEFAULT 60,               -- 秒
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.9 Webhook事件日志表 (webhook_logs)

```sql
CREATE TABLE webhook_logs (
    id              BIGSERIAL PRIMARY KEY,
    webhook_id      BIGINT NOT NULL REFERENCES webhooks(id),
    event_type      VARCHAR(50) NOT NULL,
    payload         JSONB,
    response_code   INT,
    response_body   TEXT,
    error_message   TEXT,
    attempt         INT DEFAULT 1,
    status          VARCHAR(20),                   -- 'success'/'failed'/'pending'
    sent_at         TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_webhook_logs_webhook ON webhook_logs(webhook_id, created_at DESC);
```

---

## 五、前端页面清单

### 5.1 订阅管理

| 页面 | 路由 | 说明 |
|------|------|------|
| 订阅计划 | /subscription/plans | 订阅计划列表/对比 |
| 当前订阅 | /subscription/current | 当前订阅详情 |
| 订阅升级 | /subscription/upgrade | 升级确认/支付 |
| 订阅降级 | /subscription/downgrade | 降级确认 |
| 取消订阅 | /subscription/cancel | 取消流程/挽留 |

### 5.2 用量计费

| 页面 | 路由 | 说明 |
|------|------|------|
| 用量中心 | /usage | 当前用量/配额 |
| 用量历史 | /usage/history | 历史用量查询 |
| 用量预测 | /usage/prediction | 本月用量预测 |
| 购买配额 | /usage/quotas/purchase | 购买额外配额 |

### 5.3 支付

| 页面 | 路由 | 说明 |
|------|------|------|
| 支付页面 | /payment/:order_id | 支付页面 |
| 支付方式 | /payment/methods | 支付方式管理 |

### 5.4 发票

| 页面 | 路由 | 说明 |
|------|------|------|
| 发票列表 | /invoices | 发票历史 |
| 申请开票 | /invoices/apply | 申请发票 |
| 发票抬头 | /invoices/titles | 发票抬头管理 |
| 发票详情 | /invoices/:id | 发票详情 |

### 5.5 账单

| 页面 | 路由 | 说明 |
|------|------|------|
| 账单列表 | /billing/statements | 账单列表 |
| 账单详情 | /billing/statements/:id | 账单明细 |
| 账单汇总 | /billing/summary | 费用汇总 |

### 5.6 Webhook

| 页面 | 路由 | 说明 |
|------|------|------|
| Webhook列表 | /webhooks | Webhook配置列表 |
| 创建Webhook | /webhooks/create | 创建Webhook |
| Webhook详情 | /webhooks/:id | Webhook详情 |
| Webhook日志 | /webhooks/:id/logs | 事件日志 |

---

## 六、验收标准

### 6.1 订阅管理

| 验收点 | 标准 |
|--------|------|
| 订阅创建 | 用户成功订阅，状态立即变为active |
| 自动续费 | 到期日自动扣款续费 |
| 升级即时生效 | 升级操作完成后立即可用新功能 |
| 降级到周期末 | 降级操作记录，下次周期生效 |
| 取消成功 | 取消后不再自动扣款 |

### 6.2 用量计费

| 验收点 | 标准 |
|--------|------|
| 用量实时性 | 用量数据延迟<5分钟 |
| 预警准确性 | 用量超过阈值时准时预警 |
| 账单准确性 | 账单金额与实际用量一致 |

### 6.3 发票

| 验收点 | 标准 |
|--------|------|
| 发票开具 | 申请后3个工作日内开具 |
| 发票内容 | 发票内容与实际交易一致 |
| 电子发票 | 电子发票可下载/推送 |

### 6.4 Webhook

| 验收点 | 标准 |
|--------|------|
| 事件触发 | 事件发生后<1分钟发送 |
| 签名正确 | 签名验证100%准确 |
| 重试机制 | 失败后按配置的间隔重试 |
| 事件不丢 | 事件发送成功前不丢失 |


---

## 六、页面布局规范

### 6.1 订阅计划页面（/subscription/plans）

**布局结构：**
1. 面包屑 → 页面标题
2. 订阅计划卡片网格（Free/Basic/Pro/Enterprise/Unlimited 并排展示）
3. 计划对比表格

**按钮规范：**
- [立即订阅] — 各卡片底部居中
- [联系我们] — 企业版卡片

### 6.2 当前订阅页面（/subscription/current）

**布局结构：**
1. 面包屑 → 页面标题
2. 当前订阅信息卡片（白色）：计划名称/到期时间/自动续费状态
3. 操作按钮（升级/降级/取消靠右）
4. 订阅历史表格

**按钮规范：**
- [升级] [降级] [取消订阅] — 右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 变更类型 | 100px | upgrade/downgrade/renew/cancel |
| 原计划 | 100px | - |
| 新计划 | 100px | - |
| 金额 | 100px | - |
| 生效时间 | 150px | - |
| 操作时间 | 150px | - |

**分页：** 右下角，10/20/50/100 条

### 6.3 用量中心页面（/usage）

**布局结构：**
1. 面包屑 → 页面标题
2. 当前配额卡片（白色）：AI对话次数/存储空间/API调用
3. 用量趋势图表
4. 用量明细表格

**按钮规范：**
- [购买配额] — 右对齐

**分页：** 右下角，10/20/50/100 条

### 6.4 发票列表页面（/invoices）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（申请开票靠左）
3. 发票列表表格

**按钮规范：**
- [申请开票] — 左对齐
- [详情] [下载] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 发票号 | 150px | - |
| 发票类型 | 100px | 普通/专用/电子 |
| 抬头 | 200px | - |
| 金额 | 100px | - |
| 状态 | 100px | pending/issued/cancelled |
| 开票时间 | 150px | - |
| 操作 | 120px | 详情/下载 |

**分页：** 右下角，10/20/50/100 条

### 6.5 账单页面（/billing/statements）

**布局结构：**
1. 面包屑 → 页面标题
2. 账单汇总卡片（白色）：本期应付/已付/逾期）—— 白色
3. 账单列表表格

**按钮规范：**
- [支付] [导出] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 账单号 | 150px | - |
| 账期 | 150px | - |
| 金额 | 100px | - |
| 状态 | 100px | pending/paid/overdue |
| 到期日 | 120px | - |
| 操作 | 120px | 详情/支付/导出 |

**分页：** 右下角，10/20/50/100 条

### 6.6 Webhook配置页面（/webhooks）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（创建Webhook靠左）
3. Webhook列表表格

**按钮规范：**
- [创建Webhook] — 左对齐
- [详情] [编辑] [测试] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| Webhook名称 | 200px | - |
| 端点URL | 250px | - |
| 订阅事件 | 200px | - |
| 状态 | 80px | active/inactive |
| 操作 | 120px | 详情/编辑/测试/删除 |

**分页：** 右下角，10/20/50/100 条

### 6.7 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 创建/编辑Webhook、订阅升级/降级确认 |
| Dialog 对话框 | 确认取消订阅、确认删除Webhook |
| 全屏模态 | 暂无复杂表单场景 |
