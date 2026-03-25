-- ============================================
-- Sprint 16: 商业化/订阅
-- ============================================

-- Webhook表
CREATE TABLE IF NOT EXISTS webhooks (
    id BIGSERIAL PRIMARY KEY,
    webhook_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    webhook_name VARCHAR(100) NOT NULL,
    endpoint_url VARCHAR(500) NOT NULL,
    secret_key VARCHAR(255),
    events VARCHAR(50)[],
    is_active BOOLEAN DEFAULT true,
    retry_count INT DEFAULT 3,
    retry_interval INT DEFAULT 60,
    headers JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Webhook日志表
CREATE TABLE IF NOT EXISTS webhook_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    webhook_id BIGINT NOT NULL REFERENCES webhooks(id),
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    response_code INT,
    response_body TEXT,
    error_message TEXT,
    attempt INT DEFAULT 1,
    status VARCHAR(20), -- 'success', 'failed', 'pending'
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 账单表
CREATE TABLE IF NOT EXISTS billing_statements (
    id BIGSERIAL PRIMARY KEY,
    statement_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    billing_period VARCHAR(20) NOT NULL,
    statement_date DATE NOT NULL,
    due_date DATE,
    subtotal DECIMAL(10,2) DEFAULT 0,
    discount DECIMAL(10,2) DEFAULT 0,
    tax DECIMAL(10,2) DEFAULT 0,
    total DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(8) DEFAULT 'CNY',
    status VARCHAR(20) DEFAULT 'pending',
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 发票表
CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    invoice_id VARCHAR(36) NOT NULL UNIQUE,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    statement_id VARCHAR(36),
    billing_amount DECIMAL(10,2) NOT NULL,
    tax_rate DECIMAL(5,2) DEFAULT 0.06,
    tax_amount DECIMAL(10,2),
    total_amount DECIMAL(10,2) NOT NULL,
    invoice_type VARCHAR(20), -- 'personal', 'company'
    company_name VARCHAR(200),
    tax_id VARCHAR(50),
    billing_address TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    issued_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 使用量记录表
CREATE TABLE IF NOT EXISTS usage_records (
    id BIGSERIAL PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    subscription_id VARCHAR(36),
    usage_type VARCHAR(50) NOT NULL,
    usage_value DECIMAL(10,4) NOT NULL,
    unit VARCHAR(20),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    cost DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 17: 情感计算
-- ============================================

-- 宠物情绪行为表
CREATE TABLE IF NOT EXISTS pet_emotion_actions (
    id BIGSERIAL PRIMARY KEY,
    action_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    emotion_type VARCHAR(50) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    intensity FLOAT DEFAULT 0.5,
    trigger_source VARCHAR(50),
    context_data JSONB,
    action_result VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 家庭情绪表
CREATE TABLE IF NOT EXISTS family_emotions (
    id BIGSERIAL PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36),
    member_id BIGINT NOT NULL,
    emotion_type VARCHAR(50) NOT NULL,
    intensity FLOAT DEFAULT 0.5,
    context TEXT,
    source VARCHAR(50),
    detected_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 18: 数字孪生
-- ============================================

-- 数字孪生宠物表
CREATE TABLE IF NOT EXISTS digital_twin_pets (
    id BIGSERIAL PRIMARY KEY,
    twin_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    twin_state JSONB,
    sync_status VARCHAR(20) DEFAULT 'synced',
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 行为事件表
CREATE TABLE IF NOT EXISTS behavior_events (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    device_id VARCHAR(64),
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB,
    location JSONB,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 精彩瞬间表
CREATE TABLE IF NOT EXISTS highlight_moments (
    id BIGSERIAL PRIMARY KEY,
    moment_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    moment_type VARCHAR(50) NOT NULL,
    media_url VARCHAR(512),
    thumbnail_url VARCHAR(512),
    description TEXT,
    emotion_tags VARCHAR(50)[],
    captured_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 同步记录表
CREATE TABLE IF NOT EXISTS sync_records (
    id BIGSERIAL PRIMARY KEY,
    sync_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    sync_type VARCHAR(50) NOT NULL,
    sync_data JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
