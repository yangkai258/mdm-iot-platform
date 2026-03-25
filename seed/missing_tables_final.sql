-- ============================================
-- 剩余缺失的表
-- ============================================

-- 部门表 (06_权限规范.md)
CREATE TABLE IF NOT EXISTS departments (
    id BIGSERIAL PRIMARY KEY,
    dept_id VARCHAR(36) NOT NULL UNIQUE,
    parent_id BIGINT,
    dept_name VARCHAR(100) NOT NULL,
    dept_code VARCHAR(50),
    manager_id BIGINT,
    sort_order INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    tenant_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 租户配置表 (06_权限规范.md)
CREATE TABLE IF NOT EXISTS tenant_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    tenant_id UUID NOT NULL,
    config_key VARCHAR(100) NOT NULL,
    config_value TEXT,
    config_type VARCHAR(20),
    description VARCHAR(255),
    is_encrypted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Webhook市场配置表 (PRD_26_WEBHOOK_MARKET.md)
CREATE TABLE IF NOT EXISTS webhook_market_templates (
    id BIGSERIAL PRIMARY KEY,
    template_id VARCHAR(36) NOT NULL UNIQUE,
    template_name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    description TEXT,
    webhook_template JSONB,
    sample_payload JSONB,
    documentation_url VARCHAR(500),
    is_official BOOLEAN DEFAULT false,
    install_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Webhook市场配置记录表 (PRD_26_WEBHOOK_MARKET.md)
CREATE TABLE IF NOT EXISTS webhook_market_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    template_id VARCHAR(36),
    webhook_name VARCHAR(100),
    endpoint_url VARCHAR(500),
    custom_headers JSONB,
    settings JSONB,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Webhook市场日志表 (PRD_26_WEBHOOK_MARKET.md)
CREATE TABLE IF NOT EXISTS webhook_market_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    config_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    response_code INT,
    response_body TEXT,
    error_message TEXT,
    attempt INT DEFAULT 1,
    status VARCHAR(20),
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 仿真场景表 (MODULE_SIMULATION.md)
CREATE TABLE IF NOT EXISTS simulation_scenes (
    id BIGSERIAL PRIMARY KEY,
    scene_id VARCHAR(36) NOT NULL UNIQUE,
    scene_name VARCHAR(200) NOT NULL,
    scene_type VARCHAR(50) NOT NULL,
    description TEXT,
    scene_config JSONB,
    environment_id VARCHAR(36),
    duration_seconds INT,
    difficulty_level INT DEFAULT 1,
    success_criteria JSONB,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
