-- Migration: 20260323_create_device_metrics
-- Description: 创建设备监控指标表

CREATE TABLE IF NOT EXISTS device_metrics (
    id              BIGSERIAL PRIMARY KEY,
    metric_id       VARCHAR(64) NOT NULL UNIQUE,
    device_id       VARCHAR(64) NOT NULL,
    metric_type     VARCHAR(32) NOT NULL,
    metric_name     VARCHAR(64) NOT NULL,
    metric_value    DOUBLE PRECISION NOT NULL,
    unit            VARCHAR(16),
    timestamp       TIMESTAMP NOT NULL,
    tags            JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW(),
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_device_metrics_device_id ON device_metrics(device_id);
CREATE INDEX IF NOT EXISTS idx_device_metrics_metric_id ON device_metrics(metric_id);
CREATE INDEX IF NOT EXISTS idx_device_metrics_metric_type ON device_metrics(metric_type);
CREATE INDEX IF NOT EXISTS idx_device_metrics_timestamp ON device_metrics(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_device_metrics_device_type_time ON device_metrics(device_id, metric_type, timestamp DESC);
