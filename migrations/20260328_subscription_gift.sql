-- 订阅赠送表
CREATE TABLE IF NOT EXISTS subscription_gifts (
    id BIGSERIAL PRIMARY KEY,
    gift_code VARCHAR(64) UNIQUE NOT NULL,
    sender_id VARCHAR(64) NOT NULL,
    sender_name VARCHAR(100),
    recipient_id VARCHAR(64),
    recipient_name VARCHAR(100),
    recipient_email VARCHAR(200),
    plan_id BIGINT NOT NULL,
    plan_name VARCHAR(100),
    duration INTEGER DEFAULT 30,
    status VARCHAR(20) DEFAULT 'pending',
    claimed_at TIMESTAMP,
    expires_at TIMESTAMP,
    sent_at TIMESTAMP,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscription_gifts_sender ON subscription_gifts(sender_id);
CREATE INDEX IF NOT EXISTS idx_subscription_gifts_recipient ON subscription_gifts(recipient_id);
CREATE INDEX IF NOT EXISTS idx_subscription_gifts_code ON subscription_gifts(gift_code);

-- 订阅赠送使用记录表
CREATE TABLE IF NOT EXISTS subscription_gift_usage (
    id BIGSERIAL PRIMARY KEY,
    gift_id BIGINT NOT NULL,
    used_by_device VARCHAR(64),
    used_by_ip VARCHAR(50),
    claimed_feature VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sgu_gift ON subscription_gift_usage(gift_id);
