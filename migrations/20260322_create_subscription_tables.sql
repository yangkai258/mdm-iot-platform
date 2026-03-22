-- Sprint 16: 订阅和计费系统数据库迁移
-- 执行时间: 2026-03-22

-- =============================================
-- 1. 订阅计划表 (subscription_plans)
-- =============================================
CREATE TABLE IF NOT EXISTS subscription_plans (
    id              BIGSERIAL PRIMARY KEY,
    plan_id         VARCHAR(64) NOT NULL UNIQUE,
    plan_name       VARCHAR(64) NOT NULL,
    plan_type       VARCHAR(32) NOT NULL,          -- free/basic/pro/enterprise
    price           DECIMAL(10,2) DEFAULT 0,         -- 月费
    currency        VARCHAR(8) DEFAULT 'CNY',
    duration_days   INT DEFAULT 30,                  -- 订阅周期（天）
    features        JSONB DEFAULT '{}',              -- 功能列表
    quotas          JSONB DEFAULT '{}',              -- 配额限制
    status          VARCHAR(20) DEFAULT 'active',   -- active/inactive
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscription_plans_status ON subscription_plans(status);
CREATE INDEX IF NOT EXISTS idx_subscription_plans_plan_type ON subscription_plans(plan_type);

-- =============================================
-- 2. 用户订阅表 (user_subscriptions)
-- =============================================
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    sub_id          VARCHAR(64) NOT NULL UNIQUE,
    user_id         BIGINT NOT NULL,
    plan_id         VARCHAR(64) NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',   -- active/expired/cancelled/pending
    start_time      TIMESTAMP NOT NULL,
    expire_time     TIMESTAMP NOT NULL,
    auto_renew      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_plan_id ON user_subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_status ON user_subscriptions(status);

-- =============================================
-- 3. 订阅变更记录表 (subscription_changes)
-- =============================================
CREATE TABLE IF NOT EXISTS subscription_changes (
    id              BIGSERIAL PRIMARY KEY,
    change_id       VARCHAR(64) NOT NULL UNIQUE,
    user_id         BIGINT NOT NULL,
    sub_id          VARCHAR(64),
    change_type     VARCHAR(20) NOT NULL,          -- upgrade/downgrade/renew/cancel/create
    from_plan_id    VARCHAR(64),
    to_plan_id      VARCHAR(64) NOT NULL,
    amount          DECIMAL(10,2) DEFAULT 0,
    change_reason   TEXT,
    effective_at    TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscription_changes_user_id ON subscription_changes(user_id);
CREATE INDEX IF NOT EXISTS idx_subscription_changes_sub_id ON subscription_changes(sub_id);
CREATE INDEX IF NOT EXISTS idx_subscription_changes_type ON subscription_changes(change_type);

-- =============================================
-- 4. 用量记录表 (usage_records)
-- =============================================
CREATE TABLE IF NOT EXISTS usage_records (
    id              BIGSERIAL PRIMARY KEY,
    record_id       VARCHAR(64) NOT NULL UNIQUE,
    user_id         BIGINT NOT NULL,
    usage_type      VARCHAR(32) NOT NULL,          -- api_call/device/storage/bandwidth
    usage_value     DECIMAL(12,2) DEFAULT 0,
    unit            VARCHAR(16),                    -- 次/GB/MB/个
    quota_limit     DECIMAL(12,2) DEFAULT 0,
    quota_used      DECIMAL(12,2) DEFAULT 0,
    record_date     TIMESTAMP NOT NULL,
    period_start    TIMESTAMP NOT NULL,
    period_end      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_usage_type ON usage_records(usage_type);
CREATE INDEX IF NOT EXISTS idx_usage_records_record_date ON usage_records(record_date);
CREATE INDEX IF NOT EXISTS idx_usage_records_period ON usage_records(period_start, period_end);

-- =============================================
-- 5. 用户配额表 (user_quotas)
-- =============================================
CREATE TABLE IF NOT EXISTS user_quotas (
    id              BIGSERIAL PRIMARY KEY,
    quota_id        VARCHAR(64) NOT NULL UNIQUE,
    user_id         BIGINT NOT NULL,
    quota_type      VARCHAR(32) NOT NULL,          -- api_call/device/storage/bandwidth
    quota_limit     DECIMAL(12,2) DEFAULT 0,
    quota_used      DECIMAL(12,2) DEFAULT 0,
    unit            VARCHAR(16),
    period_type     VARCHAR(20) DEFAULT 'monthly', -- monthly/yearly/daily
    period_start    TIMESTAMP NOT NULL,
    period_end      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_quotas_user_id ON user_quotas(user_id);
CREATE INDEX IF NOT EXISTS idx_user_quotas_type ON user_quotas(quota_type);
CREATE INDEX IF NOT EXISTS idx_user_quotas_period ON user_quotas(period_start, period_end);

-- =============================================
-- 6. Webhook 配置表 (webhooks)
-- =============================================
CREATE TABLE IF NOT EXISTS webhooks (
    id              BIGSERIAL PRIMARY KEY,
    webhook_id      VARCHAR(64) NOT NULL UNIQUE,
    name            VARCHAR(128) NOT NULL,
    url             VARCHAR(512) NOT NULL,
    secret          VARCHAR(256),
    event_types     JSONB DEFAULT '[]',             -- 订阅的事件类型
    status          VARCHAR(20) DEFAULT 'active',   -- active/inactive
    tenant_id       BIGINT NOT NULL,
    headers         JSONB DEFAULT '{}',             -- 自定义请求头
    retry_count     INT DEFAULT 3,                  -- 重试次数
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhooks_webhook_id ON webhooks(webhook_id);
CREATE INDEX IF NOT EXISTS idx_webhooks_tenant_id ON webhooks(tenant_id);
CREATE INDEX IF NOT EXISTS idx_webhooks_status ON webhooks(status);

-- =============================================
-- 7. Webhook 事件表 (webhook_events)
-- =============================================
CREATE TABLE IF NOT EXISTS webhook_events (
    id              BIGSERIAL PRIMARY KEY,
    event_id        VARCHAR(64) NOT NULL UNIQUE,
    webhook_id      VARCHAR(64) NOT NULL,
    event_type      VARCHAR(64) NOT NULL,
    payload         JSONB DEFAULT '{}',
    status          VARCHAR(20) DEFAULT 'pending', -- pending/success/failed
    attempts        INT DEFAULT 0,
    max_attempts    INT DEFAULT 3,
    last_error      VARCHAR(512),
    response_code   SMALLINT,
    response_body   TEXT,
    delivered_at    TIMESTAMP,
    next_retry_at   TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_event_id ON webhook_events(event_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_webhook_id ON webhook_events(webhook_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_status ON webhook_events(status);
CREATE INDEX IF NOT EXISTS idx_webhook_events_created ON webhook_events(created_at);

-- =============================================
-- 8. 账单记录表 (billing_records)
-- =============================================
CREATE TABLE IF NOT EXISTS billing_records (
    id              BIGSERIAL PRIMARY KEY,
    bill_id         VARCHAR(64) NOT NULL UNIQUE,
    user_id         BIGINT NOT NULL,
    type            VARCHAR(20) NOT NULL,          -- subscription/upgrade/quota
    amount          DECIMAL(10,2) NOT NULL,
    currency        VARCHAR(8) DEFAULT 'CNY',
    status          VARCHAR(20) DEFAULT 'pending',  -- pending/paid/refunded
    pay_method      VARCHAR(32),                    -- alipay/wechat/card
    pay_time        TIMESTAMP,
    description     VARCHAR(256),
    order_no        VARCHAR(64),                    -- 外部订单号
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_billing_records_bill_id ON billing_records(bill_id);
CREATE INDEX IF NOT EXISTS idx_billing_records_user_id ON billing_records(user_id);
CREATE INDEX IF NOT EXISTS idx_billing_records_status ON billing_records(status);
CREATE INDEX IF NOT EXISTS idx_billing_records_type ON billing_records(type);
CREATE INDEX IF NOT EXISTS idx_billing_records_created ON billing_records(created_at);

-- =============================================
-- 9. 发票表 (invoices)
-- =============================================
CREATE TABLE IF NOT EXISTS invoices (
    id              BIGSERIAL PRIMARY KEY,
    invoice_id      VARCHAR(64) NOT NULL UNIQUE,
    bill_id         VARCHAR(64),
    user_id         BIGINT NOT NULL,
    title           VARCHAR(128) NOT NULL,          -- 发票抬头
    tax_no          VARCHAR(32),                    -- 税号
    amount          DECIMAL(10,2) NOT NULL,
    tax_amount      DECIMAL(10,2) DEFAULT 0,
    tax_rate        DECIMAL(5,2) DEFAULT 6,          -- 税率
    status          VARCHAR(20) DEFAULT 'pending',  -- pending/issued/rejected
    invoice_type    VARCHAR(20) DEFAULT 'electronic', -- normal/special/electronic
    receiver_email  VARCHAR(128),
    receiver_phone  VARCHAR(32),
    issued_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_invoices_invoice_id ON invoices(invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoices_bill_id ON invoices(bill_id);
CREATE INDEX IF NOT EXISTS idx_invoices_user_id ON invoices(user_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);

-- =============================================
-- 10. 插入默认订阅计划
-- =============================================
INSERT INTO subscription_plans (plan_id, plan_name, plan_type, price, currency, duration_days, features, quotas, status, sort_order)
VALUES 
    ('plan_free', '免费版', 'free', 0, 'CNY', 30, 
     '{"ai_chat": false, "storage": false, "priority_support": false}',
     '{"api_call": 100, "device": 1, "storage": 0.1, "bandwidth": 1}',
     'active', 1),
    ('plan_basic', '基础版', 'basic', 29.9, 'CNY', 30,
     '{"ai_chat": true, "storage": true, "priority_support": false}',
     '{"api_call": 1000, "device": 5, "storage": 1, "bandwidth": 10}',
     'active', 2),
    ('plan_pro', '专业版', 'pro', 99.9, 'CNY', 30,
     '{"ai_chat": true, "storage": true, "priority_support": true}',
     '{"api_call": 10000, "device": 20, "storage": 10, "bandwidth": 100}',
     'active', 3),
    ('plan_enterprise', '企业版', 'enterprise', 299.9, 'CNY', 30,
     '{"ai_chat": true, "storage": true, "priority_support": true, "custom": true}',
     '{"api_call": -1, "device": -1, "storage": 100, "bandwidth": 1000}',
     'active', 4)
ON CONFLICT (plan_id) DO NOTHING;

-- =============================================
-- 11. 创建更新时间戳触发器函数
-- =============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有新表创建更新触发器
DO $$
DECLARE
    t text;
BEGIN
    FOREACH t IN ARRAY ARRAY[
        'subscription_plans',
        'user_subscriptions',
        'subscription_changes',
        'usage_records',
        'user_quotas',
        'webhooks',
        'webhook_events',
        'billing_records',
        'invoices'
    ]
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS update_%s_updated_at ON %s', t, t);
        EXECUTE format('CREATE TRIGGER update_%s_updated_at BEFORE UPDATE ON %s FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()', t, t);
    END LOOP;
END;
$$ LANGUAGE plpgsql;
