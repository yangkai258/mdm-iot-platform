# Sprint 16 规划

**时间**：2026-06-28
**状态**：待开始
**Sprint 周期**：2 周（2026-06-28 ～ 2026-07-11）

---

## 一、Sprint 目标

**目标：** 订阅和计费系统

在 Sprint 15（宠物生态）的基础上，实现完整的订阅和计费系统，包括订阅计划管理、自动续费、用量计费、Webhook 事件系统，为产品商业化提供完整的基础设施。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **订阅管理 API** | 完成 `/api/v1/subscriptions/*` 订阅 CRUD | subscription_controller.go | P0 |
| P0-2 | **订阅计划数据库** | 创建 subscription_plans + user_subscriptions 表 | models/subscription.go | P0 |
| P0-3 | **用量计费 API** | 完成 `/api/v1/usage/*` 用量查询和统计 | usage_controller.go | P0 |
| P0-4 | **Webhook 事件系统** | 完成 `/api/v1/webhooks/*` Webhook 配置和事件发送 | webhook_controller.go | P0 |
| P0-5 | **账单 API** | 完成 `/api/v1/billing/*` 账单查询接口 | billing_controller.go | P0 |
| P1-1 | **自动续费服务** | 实现定时检查+自动扣款续费 | subscription/renewal_service.go | P1 |
| P1-2 | **订阅变更服务** | 实现升级/降级/取消逻辑 | subscription/change_service.go | P1 |
| P1-3 | **发票 API** | 完成 `/api/v1/invoices/*` 发票申请和管理 | invoice_controller.go | P1 |
| P1-4 | **配额管理 API** | 完成 `/api/v1/usage/quotas/*` 配额管理 | quota_controller.go | P1 |
| P2-1 | **支付网关集成** | 集成支付宝/微信支付 | payment/gateway_service.go | P2 |
| P2-2 | **计费异常检测** | 用量异常自动告警 | usage/anomaly_detector.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **订阅管理前端** | 完成 SubscriptionManageView.vue 订阅计划管理 | SubscriptionManageView.vue | P0 |
| PF0-2 | **订阅计划展示** | 完成 SubscriptionPlansView.vue 计划列表/对比 | SubscriptionPlansView.vue | P0 |
| PF0-3 | **发票账单前端** | 完成 InvoiceBillingView.vue 发票申请/账单查看 | InvoiceBillingView.vue | P0 |
| PF0-4 | **用量查询页面** | 完成 UsageQueryView.vue 当前用量/配额查询 | UsageQueryView.vue | P0 |
| PF1-1 | **订阅升级/降级** | 完成 SubscriptionChangeView.vue 升级降级流程 | SubscriptionChangeView.vue | P1 |
| PF1-2 | **Webhook 配置页面** | 完成 WebhookConfigView.vue Webhook 创建/管理 | WebhookConfigView.vue | P1 |
| PF1-3 | **配额购买页面** | 完成 QuotaPurchaseView.vue 额外配额购买 | QuotaPurchaseView.vue | P1 |
| PF2-1 | **账单导出功能** | 完成 BillingExportView.vue 账单 PDF/Excel 导出 | BillingExportView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/subscriptions/plans` | GET | 订阅计划列表 |
| `GET /api/v1/subscriptions/plans/:id` | GET | 计划详情 |
| `GET /api/v1/subscriptions/current` | GET | 当前订阅 |
| `POST /api/v1/subscriptions/create` | POST | 创建订阅 |
| `POST /api/v1/subscriptions/upgrade` | POST | 升级订阅 |
| `POST /api/v1/subscriptions/downgrade` | POST | 降级订阅 |
| `POST /api/v1/subscriptions/cancel` | POST | 取消订阅 |
| `POST /api/v1/subscriptions/resume` | POST | 恢复订阅 |
| `POST /api/v1/subscriptions/renew` | POST | 手动续费 |
| `GET /api/v1/usage/current` | GET | 当前用量 |
| `GET /api/v1/usage/history` | GET | 用量历史 |
| `GET /api/v1/usage/prediction` | GET | 用量预测 |
| `GET /api/v1/usage/quotas` | GET | 配额查询 |
| `POST /api/v1/usage/quotas/purchase` | POST | 购买配额 |
| `GET /api/v1/billing/statements` | GET | 账单列表 |
| `GET /api/v1/billing/statements/:id` | GET | 账单详情 |
| `GET /api/v1/billing/summary` | GET | 账单汇总 |
| `GET /api/v1/invoices` | GET | 发票列表 |
| `GET /api/v1/invoices/:id` | GET | 发票详情 |
| `POST /api/v1/invoices` | POST | 申请开票 |
| `POST /api/v1/invoices/:id/download` | POST | 下载发票 |
| `POST /api/v1/invoices/titles` | POST | 添加发票抬头 |
| `GET /api/v1/webhooks` | GET | Webhook 列表 |
| `POST /api/v1/webhooks` | POST | 创建 Webhook |
| `GET /api/v1/webhooks/:id` | GET | Webhook 详情 |
| `PUT /api/v1/webhooks/:id` | PUT | 更新 Webhook |
| `DELETE /api/v1/webhooks/:id` | DELETE | 删除 Webhook |
| `GET /api/v1/webhooks/:id/logs` | GET | 事件日志 |
| `POST /api/v1/webhooks/:id/test` | POST | 测试 Webhook |

### 数据库设计

（继承 MODULE_SUBSCRIPTION 设计的完整表结构）

```sql
-- 订阅计划表
CREATE TABLE subscription_plans (
    id              BIGSERIAL PRIMARY KEY,
    plan_code       VARCHAR(50) NOT NULL UNIQUE,
    plan_name       VARCHAR(100) NOT NULL,
    plan_type       VARCHAR(20) NOT NULL,         -- 'free'/'paid'
    price_monthly   DECIMAL(10,2) DEFAULT 0,
    price_yearly    DECIMAL(10,2),
    features        JSONB,
    quotas          JSONB,
    is_active       BOOLEAN DEFAULT TRUE,
    is_recommended  BOOLEAN DEFAULT FALSE,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 用户订阅表
CREATE TABLE user_subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    plan_id         BIGINT NOT NULL REFERENCES subscription_plans(id),
    status          VARCHAR(20) NOT NULL,         -- 'active'/'cancelled'/'expired'/'trial'
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

-- 订阅变更记录表
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

-- 用量记录表
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

-- 账单表
CREATE TABLE billing_statements (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    statement_no    VARCHAR(50) NOT NULL UNIQUE,
    period_start    DATE NOT NULL,
    period_end      DATE NOT NULL,
    subtotal        DECIMAL(10,2) NOT NULL,
    discount        DECIMAL(10,2) DEFAULT 0,
    total           DECIMAL(10,2) NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',
    due_date        DATE,
    paid_at         TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 发票表
CREATE TABLE invoices (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    invoice_no      VARCHAR(50) NOT NULL UNIQUE,
    invoice_type    VARCHAR(20) NOT NULL,         -- 'normal'/'special'/'electronic'
    title           VARCHAR(255) NOT NULL,
    tax_number      VARCHAR(50),
    amount          DECIMAL(10,2) NOT NULL,
    tax_amount      DECIMAL(10,2),
    status          VARCHAR(20) DEFAULT 'pending',
    issued_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Webhook 配置表
CREATE TABLE webhooks (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    webhook_name    VARCHAR(255) NOT NULL,
    endpoint_url    VARCHAR(500) NOT NULL,
    secret_key      VARCHAR(255),
    events          VARCHAR(50)[],
    is_active       BOOLEAN DEFAULT TRUE,
    retry_count     INT DEFAULT 3,
    retry_interval  INT DEFAULT 60,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- Webhook 事件日志表
CREATE TABLE webhook_logs (
    id              BIGSERIAL PRIMARY KEY,
    webhook_id      BIGINT NOT NULL REFERENCES webhooks(id),
    event_type      VARCHAR(50) NOT NULL,
    payload         JSONB,
    response_code   INT,
    response_body   TEXT,
    error_message   TEXT,
    attempt         INT DEFAULT 1,
    status          VARCHAR(20),
    sent_at         TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 订阅创建 | 用户成功订阅，状态立即变为 active | 调用 API 验证 |
| 自动续费 | 到期日自动扣款续费 | 模拟到期测试 |
| 升级即时生效 | 升级操作完成后立即可用新功能 | 调用 API 验证 |
| 降级到周期末 | 降级操作记录，下次周期生效 | 调用 API 验证 |
| 取消成功 | 取消后不再自动扣款 | 调用 API 验证 |
| 用量实时性 | 用量数据延迟 < 5 分钟 | 调用 API 验证 |
| 预警准确性 | 用量超过阈值时准时预警 | 模拟超量测试 |
| Webhook 触发 | 事件发生后 < 1 分钟发送 | 触发事件验证 |
| 签名正确 | Webhook 签名验证 100% 准确 | 伪造签名测试 |
| 重试机制 | 失败后按配置的间隔重试 | 模拟失败测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 用量查询延迟 | <= 200ms |
| 订阅变更处理 | <= 1s |
| Webhook 发送延迟 | <= 1s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| 支付网关 | 支付宝/微信支付商户号 |
| 短信/邮件服务 | 续费提醒通知 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 支付失败 | 续费不成功 | 多次重试+人工催缴 |
| 用量数据延迟 | 计费不准确 | 对账机制 |
| Webhook 接收方不可用 | 事件丢失 | 重试+日志持久化 |
