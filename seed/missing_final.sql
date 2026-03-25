-- 补充遗漏的具身智能安全日志表

CREATE TABLE IF NOT EXISTS embodied_safety_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20),
    description TEXT,
    location JSONB,
    triggered_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP,
    resolution TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_embodied_safety_logs_device_id ON embodied_safety_logs (device_id);
CREATE INDEX IF NOT EXISTS idx_embodied_safety_logs_triggered_at ON embodied_safety_logs (triggered_at DESC);
