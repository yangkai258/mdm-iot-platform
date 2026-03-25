-- ============================================
-- Sprint 19: 健康医疗
-- ============================================

-- 疾病模式表
CREATE TABLE IF NOT EXISTS disease_patterns (
    id BIGSERIAL PRIMARY KEY,
    pattern_id VARCHAR(36) NOT NULL UNIQUE,
    pattern_name VARCHAR(100) NOT NULL,
    disease_name VARCHAR(100) NOT NULL,
    species VARCHAR(32),
    symptoms JSONB,
    risk_factors JSONB,
    severity VARCHAR(20),
    description TEXT,
    doc_references TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 睡眠记录表
CREATE TABLE IF NOT EXISTS sleep_records (
    id BIGSERIAL PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    sleep_start TIMESTAMP NOT NULL,
    sleep_end TIMESTAMP,
    duration_minutes INT,
    sleep_quality VARCHAR(20),
    sleep_stages JSONB,
    interruptions INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 健康基线表
CREATE TABLE IF NOT EXISTS health_baselines (
    id BIGSERIAL PRIMARY KEY,
    baseline_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    metric_type VARCHAR(50) NOT NULL,
    baseline_value DECIMAL(10,4),
    normal_range_min DECIMAL(10,4),
    normal_range_max DECIMAL(10,4),
    measurement_unit VARCHAR(20),
    last_measured_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 健康报告表
CREATE TABLE IF NOT EXISTS health_reports (
    id BIGSERIAL PRIMARY KEY,
    report_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    report_type VARCHAR(50) NOT NULL, -- 'daily', 'weekly', 'monthly'
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    health_summary JSONB,
    activity_summary JSONB,
    recommendations JSONB,
    vet_recommendation TEXT,
    generated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 20: 家庭场景
-- ============================================

-- 家庭表
CREATE TABLE IF NOT EXISTS households (
    id BIGSERIAL PRIMARY KEY,
    household_id VARCHAR(36) NOT NULL UNIQUE,
    household_name VARCHAR(100) NOT NULL,
    owner_id BIGINT NOT NULL,
    address TEXT,
    location JSONB,
    family_photo_url VARCHAR(512),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 家庭成员表
CREATE TABLE IF NOT EXISTS household_members (
    id BIGSERIAL PRIMARY KEY,
    member_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36) NOT NULL,
    user_id BIGINT,
    member_name VARCHAR(64) NOT NULL,
    relation VARCHAR(50), -- 'owner', 'spouse', 'child', 'parent', 'guest'
    role VARCHAR(20) DEFAULT 'member', -- 'admin', 'member', 'guest'
    avatar_url VARCHAR(512),
    is_active BOOLEAN DEFAULT true,
    joined_at TIMESTAMP DEFAULT NOW()
);

-- 家庭邀请表
CREATE TABLE IF NOT EXISTS household_invites (
    id BIGSERIAL PRIMARY KEY,
    invite_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36) NOT NULL,
    inviter_id BIGINT NOT NULL,
    invitee_email VARCHAR(100),
    invitee_phone VARCHAR(20),
    invite_token VARCHAR(100),
    role VARCHAR(20) DEFAULT 'member',
    status VARCHAR(20) DEFAULT 'pending',
    expires_at TIMESTAMP,
    responded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 家庭模式配置表
CREATE TABLE IF NOT EXISTS family_mode_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36),
    mode_type VARCHAR(50) NOT NULL, -- 'child_mode', 'elderly_mode', 'general'
    is_enabled BOOLEAN DEFAULT false,
    settings JSONB,
    schedule JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 互动记录表
CREATE TABLE IF NOT EXISTS interaction_records (
    id BIGSERIAL PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    device_id VARCHAR(64),
    user_id BIGINT NOT NULL,
    interaction_type VARCHAR(50) NOT NULL,
    interaction_data JSONB,
    duration_seconds INT,
    feedback VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 家庭活动表
CREATE TABLE IF NOT EXISTS family_activities (
    id BIGSERIAL PRIMARY KEY,
    activity_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36) NOT NULL,
    activity_type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    location TEXT,
    scheduled_at TIMESTAMP,
    duration_minutes INT,
    participants BIGINT[],
    cover_image_url VARCHAR(512),
    status VARCHAR(20) DEFAULT 'scheduled',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Subscription: API配额
-- ============================================

-- API配额表
CREATE TABLE IF NOT EXISTS api_quotas (
    id BIGSERIAL PRIMARY KEY,
    quota_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    api_name VARCHAR(100) NOT NULL,
    quota_limit INT NOT NULL,
    quota_used INT DEFAULT 0,
    period_type VARCHAR(20) DEFAULT 'monthly', -- 'daily', 'monthly', 'yearly'
    period_start DATE,
    period_end DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
