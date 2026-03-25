-- ============================================
-- Sprint 21: 具身智能感知
-- ============================================

CREATE TABLE IF NOT EXISTS embodied_perceptions (
    id BIGSERIAL PRIMARY KEY,
    perception_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    perception_type VARCHAR(50) NOT NULL,
    raw_data JSONB NOT NULL,
    processed_data JSONB,
    confidence_score FLOAT,
    environment_context JSONB,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_embodied_perceptions_device_id ON embodied_perceptions (device_id);
CREATE INDEX IF NOT EXISTS idx_embodied_perceptions_timestamp ON embodied_perceptions (timestamp DESC);

-- ============================================
-- Sprint 22: 多宠物协作
-- ============================================

CREATE TABLE IF NOT EXISTS embodied_collaboration_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id VARCHAR(36) NOT NULL UNIQUE,
    task_name VARCHAR(200) NOT NULL,
    task_type VARCHAR(50) NOT NULL,
    initiator_device_id VARCHAR(64) NOT NULL,
    participant_device_ids VARCHAR(64)[],
    task_config JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    result_data JSONB,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS embodied_collaboration_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    task_id VARCHAR(36),
    device_id VARCHAR(64) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    action_data JSONB,
    result VARCHAR(20),
    latency_ms INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS embodied_inference_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64),
    inference_mode VARCHAR(20) DEFAULT 'cloud',
    fallback_enabled BOOLEAN DEFAULT true,
    local_model_path VARCHAR(255),
    cloud_endpoint VARCHAR(255),
    retry_config JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS embodied_safety_audits (
    id BIGSERIAL PRIMARY KEY,
    audit_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    audit_type VARCHAR(50) NOT NULL,
    audit_result JSONB,
    risk_level VARCHAR(20),
    recommendations JSONB,
    auditor VARCHAR(64),
    audited_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 24: 仿真高级功能
-- ============================================

CREATE TABLE IF NOT EXISTS simulation_cicd_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    pipeline_name VARCHAR(200) NOT NULL,
    pipeline_config JSONB NOT NULL,
    trigger_rules JSONB,
    is_enabled BOOLEAN DEFAULT true,
    last_run_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_quotas (
    id BIGSERIAL PRIMARY KEY,
    quota_id VARCHAR(36) NOT NULL UNIQUE,
    tenant_id UUID,
    quota_type VARCHAR(50) NOT NULL,
    quota_limit INT NOT NULL,
    quota_used INT DEFAULT 0,
    period_type VARCHAR(20) DEFAULT 'monthly',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_notifications (
    id BIGSERIAL PRIMARY KEY,
    notification_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    notification_type VARCHAR(50) NOT NULL,
    title VARCHAR(200),
    content TEXT,
    data JSONB,
    is_read BOOLEAN DEFAULT false,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_ab_experiments (
    id BIGSERIAL PRIMARY KEY,
    experiment_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_name VARCHAR(200) NOT NULL,
    hypothesis TEXT,
    variants JSONB NOT NULL,
    traffic_allocation JSONB,
    metrics JSONB,
    status VARCHAR(20) DEFAULT 'draft',
    start_date DATE,
    end_date DATE,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 25: 开发者平台
-- ============================================

CREATE TABLE IF NOT EXISTS developers (
    id BIGSERIAL PRIMARY KEY,
    developer_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    company_name VARCHAR(200),
    website VARCHAR(255),
    description TEXT,
    logo_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active',
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS developer_api_keys (
    id BIGSERIAL PRIMARY KEY,
    key_id VARCHAR(36) NOT NULL UNIQUE,
    developer_id VARCHAR(36) NOT NULL,
    api_key VARCHAR(64) NOT NULL UNIQUE,
    api_secret VARCHAR(128),
    permissions JSONB,
    rate_limit INT DEFAULT 1000,
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS developer_api_usage (
    id BIGSERIAL PRIMARY KEY,
    usage_id VARCHAR(36) NOT NULL UNIQUE,
    developer_id VARCHAR(36) NOT NULL,
    api_key_id VARCHAR(36),
    endpoint VARCHAR(200) NOT NULL,
    method VARCHAR(10),
    status_code INT,
    response_time_ms INT,
    request_size_bytes BIGINT,
    response_size_bytes BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 26: 应用市场
-- ============================================

CREATE TABLE IF NOT EXISTS market_apps (
    id BIGSERIAL PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL UNIQUE,
    app_name VARCHAR(200) NOT NULL,
    developer_id VARCHAR(36),
    category VARCHAR(50) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    screenshots TEXT[],
    version VARCHAR(20),
    price DECIMAL(10,2) DEFAULT 0,
    is_free BOOLEAN DEFAULT true,
    download_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS market_emojis (
    id BIGSERIAL PRIMARY KEY,
    emoji_id VARCHAR(36) NOT NULL UNIQUE,
    pack_id VARCHAR(36),
    emoji_name VARCHAR(100) NOT NULL,
    emoji_url VARCHAR(500) NOT NULL,
    category VARCHAR(50),
    tags VARCHAR(50)[],
    price DECIMAL(10,2) DEFAULT 0,
    download_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS market_actions (
    id BIGSERIAL PRIMARY KEY,
    action_id VARCHAR(36) NOT NULL UNIQUE,
    action_name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    animation_data JSONB,
    thumbnail_url VARCHAR(500),
    compatible_devices JSONB,
    price DECIMAL(10,2) DEFAULT 0,
    download_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS market_voices (
    id BIGSERIAL PRIMARY KEY,
    voice_id VARCHAR(36) NOT NULL UNIQUE,
    voice_name VARCHAR(100) NOT NULL,
    voice_type VARCHAR(50),
    preview_url VARCHAR(500),
    audio_samples JSONB,
    price DECIMAL(10,2) DEFAULT 0,
    download_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS smarthome_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64),
    integration_id VARCHAR(36),
    config_data JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pet_medical_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    config_name VARCHAR(100) NOT NULL,
    config_type VARCHAR(50) NOT NULL,
    config_data JSONB NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 27: 第三方集成
-- ============================================

CREATE TABLE IF NOT EXISTS smarthome_devices (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36),
    integration_id VARCHAR(36),
    device_name VARCHAR(100) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    serial_number VARCHAR(100),
    control_data JSONB,
    status VARCHAR(20) DEFAULT 'online',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS smarthome_linkages (
    id BIGSERIAL PRIMARY KEY,
    linkage_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36),
    linkage_name VARCHAR(100) NOT NULL,
    trigger_type VARCHAR(50) NOT NULL,
    trigger_config JSONB,
    actions JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_triggered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
