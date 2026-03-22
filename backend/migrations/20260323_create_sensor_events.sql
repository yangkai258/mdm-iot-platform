-- Migration: 20260323_create_sensor_events
-- Description: 创建传感器事件表和阈值配置表

-- 传感器事件表
CREATE TABLE IF NOT EXISTS sensor_events (
    id              BIGSERIAL PRIMARY KEY,
    event_id        VARCHAR(64) NOT NULL UNIQUE,
    device_id       VARCHAR(64) NOT NULL,
    sensor_type     VARCHAR(32) NOT NULL,
    sensor_value    DOUBLE PRECISION NOT NULL,
    unit            VARCHAR(16),
    threshold       DOUBLE PRECISION DEFAULT 0,
    is_abnormal     BOOLEAN DEFAULT FALSE,
    event_type      VARCHAR(32) NOT NULL DEFAULT 'normal',
    description     TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    deleted_at      TIMESTAMP,
    CONSTRAINT idx_sensor_events_device_time UNIQUE (device_id, sensor_type, created_at)
);

CREATE INDEX IF NOT EXISTS idx_sensor_events_device_id ON sensor_events(device_id);
CREATE INDEX IF NOT EXISTS idx_sensor_events_sensor_type ON sensor_events(sensor_type);
CREATE INDEX IF NOT EXISTS idx_sensor_events_created_at ON sensor_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sensor_events_event_type ON sensor_events(event_type);
CREATE INDEX IF NOT EXISTS idx_sensor_events_is_abnormal ON sensor_events(is_abnormal);

-- 传感器阈值配置表
CREATE TABLE IF NOT EXISTS sensor_thresholds (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    sensor_type     VARCHAR(32) NOT NULL,
    min_value       DOUBLE PRECISION DEFAULT 0,
    max_value       DOUBLE PRECISION DEFAULT 100,
    unit            VARCHAR(16),
    enabled         BOOLEAN DEFAULT TRUE,
    updated_at      TIMESTAMP DEFAULT NOW(),
    CONSTRAINT idx_sensor_thresholds_device_sensor UNIQUE (device_id, sensor_type)
);
