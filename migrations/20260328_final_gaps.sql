-- 设备健康评分表
CREATE TABLE IF NOT EXISTS device_health_scores (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL,
    total_score DECIMAL(5,2) DEFAULT 0,
    grade VARCHAR(1),
    uptime_score DECIMAL(5,2) DEFAULT 100,
    perf_score DECIMAL(5,2) DEFAULT 100,
    security_score DECIMAL(5,2) DEFAULT 100,
    behavior_score DECIMAL(5,2) DEFAULT 100,
    issues_json TEXT,
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dhs_device ON device_health_scores(device_id);

-- 告警去重规则表
CREATE TABLE IF NOT EXISTS alert_deduplication_rules (
    id BIGSERIAL PRIMARY KEY,
    rule_id VARCHAR(64) UNIQUE NOT NULL,
    alert_type VARCHAR(50) NOT NULL,
    device_pattern VARCHAR(200),
    severity_min INTEGER DEFAULT 1,
    severity_max INTEGER DEFAULT 5,
    dedup_window_seconds INTEGER DEFAULT 300,
    dedup_strategy VARCHAR(20) DEFAULT 'first',
    max_count_per_window INTEGER DEFAULT 1,
    suppression_type VARCHAR(20) DEFAULT 'none',
    is_active BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_by VARCHAR(64),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_adr_alert_type ON alert_deduplication_rules(alert_type);

-- 告警去重记录表
CREATE TABLE IF NOT EXISTS alert_deduplication_records (
    id BIGSERIAL PRIMARY KEY,
    rule_id VARCHAR(64) NOT NULL,
    alert_type VARCHAR(50),
    device_id VARCHAR(64) NOT NULL,
    window_start TIMESTAMP NOT NULL,
    window_end TIMESTAMP NOT NULL,
    alert_count INTEGER DEFAULT 1,
    alert_summary TEXT,
    final_alert_id VARCHAR(64),
    final_alert_snapshot TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_adr_device ON alert_deduplication_records(device_id);
CREATE INDEX IF NOT EXISTS idx_adr_rule ON alert_deduplication_records(rule_id);

-- 地图服务集成配置表
CREATE TABLE IF NOT EXISTS map_integration_configs (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(50) UNIQUE NOT NULL,
    api_key VARCHAR(500),
    api_secret VARCHAR(500),
    is_active BOOLEAN DEFAULT FALSE,
    services VARCHAR(200),
    quota_limit INTEGER DEFAULT 0,
    quota_used INTEGER DEFAULT 0,
    quota_reset_at TIMESTAMP,
    status VARCHAR(20),
    error_message TEXT,
    created_by VARCHAR(64),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 地图服务调用日志表
CREATE TABLE IF NOT EXISTS map_service_logs (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(50),
    service_type VARCHAR(50),
    endpoint VARCHAR(200),
    request_data TEXT,
    response_data TEXT,
    latency INTEGER DEFAULT 0,
    status_code INTEGER DEFAULT 0,
    cost DECIMAL(10,4) DEFAULT 0,
    device_id VARCHAR(64),
    user_id VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_msl_device ON map_service_logs(device_id);
CREATE INDEX IF NOT EXISTS idx_msl_user ON map_service_logs(user_id);

-- 宠物位置记录表
CREATE TABLE IF NOT EXISTS pet_locations (
    id BIGSERIAL PRIMARY KEY,
    pet_id VARCHAR(64) NOT NULL,
    device_id VARCHAR(64),
    latitude DECIMAL(10,7) NOT NULL,
    longitude DECIMAL(10,7) NOT NULL,
    altitude DECIMAL(10,2),
    address VARCHAR(500),
    poi_name VARCHAR(200),
    accuracy DECIMAL(8,2),
    battery_level INTEGER DEFAULT 100,
    location_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pl_pet ON pet_locations(pet_id);
CREATE INDEX IF NOT EXISTS idx_pl_device ON pet_locations(device_id);
CREATE INDEX IF NOT EXISTS idx_pl_created ON pet_locations(created_at);
