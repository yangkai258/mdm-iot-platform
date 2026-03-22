-- Migration: 20260323_create_alert_rules
-- Description: 创建告警规则表

CREATE TABLE IF NOT EXISTS alert_rules (
    id              BIGSERIAL PRIMARY KEY,
    rule_id         VARCHAR(64) NOT NULL UNIQUE,
    rule_name       VARCHAR(128) NOT NULL,
    rule_type       VARCHAR(32) NOT NULL DEFAULT 'device',
    conditions      JSONB DEFAULT '[]',
    actions         JSONB DEFAULT '[]',
    enabled         BOOLEAN DEFAULT TRUE,
    priority        INT DEFAULT 50,
    device_id       VARCHAR(64),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alert_rules_rule_id ON alert_rules(rule_id);
CREATE INDEX IF NOT EXISTS idx_alert_rules_rule_type ON alert_rules(rule_type);
CREATE INDEX IF NOT EXISTS idx_alert_rules_device_id ON alert_rules(device_id);
CREATE INDEX IF NOT EXISTS idx_alert_rules_enabled ON alert_rules(enabled);
CREATE INDEX IF NOT EXISTS idx_alert_rules_priority ON alert_rules(priority DESC);
