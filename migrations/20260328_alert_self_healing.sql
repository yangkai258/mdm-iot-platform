-- 告警自愈建议表
CREATE TABLE IF NOT EXISTS alert_self_healing (
    id BIGSERIAL PRIMARY KEY,
    alert_type VARCHAR(50) NOT NULL,
    alert_sub_type VARCHAR(50),
    severity INTEGER DEFAULT 2,
    title VARCHAR(200) NOT NULL,
    root_cause TEXT,
    recommendation TEXT NOT NULL,
    steps_json TEXT,
    success_rate DECIMAL(5,2) DEFAULT 0,
    used_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    tags VARCHAR(500),
    created_by VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_alert_self_healing_type ON alert_self_healing(alert_type);
CREATE INDEX IF NOT EXISTS idx_alert_self_healing_active ON alert_self_healing(is_active);

-- 告警自愈执行记录表
CREATE TABLE IF NOT EXISTS alert_self_healing_records (
    id BIGSERIAL PRIMARY KEY,
    alert_id VARCHAR(64) NOT NULL,
    self_healing_id BIGINT NOT NULL,
    alert_type VARCHAR(50),
    trigger_condition TEXT,
    steps_executed INTEGER DEFAULT 0,
    steps_total INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    executed_by VARCHAR(64),
    result TEXT,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    duration INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_self_healing_records_alert ON alert_self_healing_records(alert_id);
CREATE INDEX IF NOT EXISTS idx_self_healing_records_self_healing ON alert_self_healing_records(self_healing_id);
