-- ============================================
-- Sprint 13: 全球化/多区域数据库
-- ============================================

-- 区域表
CREATE TABLE IF NOT EXISTS regions (
    id BIGSERIAL PRIMARY KEY,
    region_code VARCHAR(20) NOT NULL UNIQUE,
    region_name VARCHAR(100) NOT NULL,
    region_name_en VARCHAR(100),
    continent VARCHAR(50),
    country VARCHAR(50),
    timezone VARCHAR(50),
    ai_node_url VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 区域节点表
CREATE TABLE IF NOT EXISTS region_nodes (
    id BIGSERIAL PRIMARY KEY,
    node_id VARCHAR(36) NOT NULL UNIQUE,
    region_id BIGINT REFERENCES regions(id),
    node_name VARCHAR(100) NOT NULL,
    node_type VARCHAR(20), -- 'primary', 'replica', 'cache'
    endpoint_url VARCHAR(255),
    health_status VARCHAR(20) DEFAULT 'healthy',
    is_master BOOLEAN DEFAULT false,
    weight INT DEFAULT 100,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 数据驻留规则表
CREATE TABLE IF NOT EXISTS data_residency_rules (
    id BIGSERIAL PRIMARY KEY,
    rule_id VARCHAR(36) NOT NULL UNIQUE,
    data_type VARCHAR(50) NOT NULL, -- 'user_data', 'device_data', 'behavior_data'
    source_region VARCHAR(20) NOT NULL,
    allowed_regions VARCHAR(20)[] NOT NULL,
    retention_days INT DEFAULT 365,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 时区配置表
CREATE TABLE IF NOT EXISTS timezone_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(50) NOT NULL UNIQUE,
    timezone_id VARCHAR(50) NOT NULL,
    utc_offset VARCHAR(10),
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 区域同步记录表
CREATE TABLE IF NOT EXISTS region_sync_records (
    id BIGSERIAL PRIMARY KEY,
    sync_id VARCHAR(36) NOT NULL UNIQUE,
    source_region VARCHAR(20) NOT NULL,
    target_region VARCHAR(20) NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    record_count INT DEFAULT 0,
    sync_status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 14: AI系统工程
-- ============================================

-- AI行为日志表
CREATE TABLE IF NOT EXISTS ai_behavior_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64),
    behavior_type VARCHAR(50) NOT NULL,
    input_data JSONB,
    output_data JSONB,
    decision_path JSONB,
    confidence_score FLOAT,
    model_version VARCHAR(50),
    latency_ms INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI沙箱测试表
CREATE TABLE IF NOT EXISTS ai_sandbox_tests (
    id BIGSERIAL PRIMARY KEY,
    test_id VARCHAR(36) NOT NULL UNIQUE,
    test_name VARCHAR(100) NOT NULL,
    test_type VARCHAR(50) NOT NULL,
    test_scenario JSONB,
    input_data JSONB,
    expected_output JSONB,
    actual_output JSONB,
    test_result VARCHAR(20), -- 'pass', 'fail', 'pending'
    error_message TEXT,
    executed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 模型回滚记录表
CREATE TABLE IF NOT EXISTS model_rollback_records (
    id BIGSERIAL PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL UNIQUE,
    model_id VARCHAR(64) NOT NULL,
    from_version VARCHAR(50) NOT NULL,
    to_version VARCHAR(50) NOT NULL,
    rollback_reason TEXT,
    rollback_status VARCHAR(20) DEFAULT 'pending',
    performed_by VARCHAR(64),
    rolled_back_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 15: 宠物生态
-- ============================================

-- 宠物基础信息表
CREATE TABLE IF NOT EXISTS pets (
    id BIGSERIAL PRIMARY KEY,
    pet_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64),
    owner_id BIGINT NOT NULL,
    pet_name VARCHAR(32) NOT NULL,
    species VARCHAR(32) NOT NULL, -- 'cat', 'dog', 'bird', etc.
    breed VARCHAR(50),
    birth_date DATE,
    gender VARCHAR(10),
    weight DECIMAL(6,2),
    height DECIMAL(6,2),
    color VARCHAR(50),
    personality JSONB,
    health_status VARCHAR(20) DEFAULT 'healthy',
    avatar_url VARCHAR(512),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 宠物设备绑定表
CREATE TABLE IF NOT EXISTS pet_device_bindings (
    id BIGSERIAL PRIMARY KEY,
    binding_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL REFERENCES pets(pet_id),
    device_id VARCHAR(64) NOT NULL,
    binding_type VARCHAR(20) DEFAULT 'primary', -- 'primary', 'secondary'
    is_active BOOLEAN DEFAULT true,
    bound_at TIMESTAMP DEFAULT NOW(),
    unbound_at TIMESTAMP
);

-- 寻宠报告表 (丢失报告)
CREATE TABLE IF NOT EXISTS lost_found_reports (
    id BIGSERIAL PRIMARY KEY,
    report_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    report_type VARCHAR(20) NOT NULL, -- 'lost', 'found'
    last_seen_location TEXT,
    last_seen_lat DECIMAL(10,6),
    last_seen_lng DECIMAL(10,6),
    description TEXT,
    reward DECIMAL(10,2),
    contact_phone VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active',
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 家庭宠物邀请表
CREATE TABLE IF NOT EXISTS household_pet_invites (
    id BIGSERIAL PRIMARY KEY,
    invite_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    household_id VARCHAR(36),
    inviter_id BIGINT NOT NULL,
    invitee_id BIGINT,
    invite_token VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',
    expires_at TIMESTAMP,
    responded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
