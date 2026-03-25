-- ============================================
-- Platform Ecosystem 模块
-- ============================================

-- 插件表
CREATE TABLE IF NOT EXISTS plugins (
    id BIGSERIAL PRIMARY KEY,
    plugin_id VARCHAR(36) NOT NULL UNIQUE,
    plugin_name VARCHAR(100) NOT NULL,
    plugin_type VARCHAR(50) NOT NULL,
    version VARCHAR(20),
    description TEXT,
    author VARCHAR(100),
    homepage_url VARCHAR(500),
    icon_url VARCHAR(500),
    manifest JSONB,
    permissions JSONB,
    is_enabled BOOLEAN DEFAULT false,
    is_official BOOLEAN DEFAULT false,
    install_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 表情包表
CREATE TABLE IF NOT EXISTS emotion_packs (
    id BIGSERIAL PRIMARY KEY,
    pack_id VARCHAR(36) NOT NULL UNIQUE,
    pack_name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    preview_urls TEXT[],
    emotions JSONB,
    price DECIMAL(10,2) DEFAULT 0,
    is_free BOOLEAN DEFAULT true,
    download_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 动作资源表
CREATE TABLE IF NOT EXISTS actions (
    id BIGSERIAL PRIMARY KEY,
    action_id VARCHAR(36) NOT NULL UNIQUE,
    action_name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    animation_data JSONB,
    audio_url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    compatible_models JSONB,
    price DECIMAL(10,2) DEFAULT 0,
    is_free BOOLEAN DEFAULT true,
    download_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 声音定制表
CREATE TABLE IF NOT EXISTS voices (
    id BIGSERIAL PRIMARY KEY,
    voice_id VARCHAR(36) NOT NULL UNIQUE,
    voice_name VARCHAR(100) NOT NULL,
    voice_type VARCHAR(50),
    preview_url VARCHAR(500),
    audio_samples JSONB,
    characteristics JSONB,
    price DECIMAL(10,2) DEFAULT 0,
    is_free BOOLEAN DEFAULT true,
    download_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 用户购买记录表
CREATE TABLE IF NOT EXISTS user_purchases (
    id BIGSERIAL PRIMARY KEY,
    purchase_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    item_type VARCHAR(50) NOT NULL, -- 'emotion_pack', 'action', 'voice'
    item_id VARCHAR(36) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(8) DEFAULT 'CNY',
    payment_method VARCHAR(20),
    transaction_id VARCHAR(100),
    status VARCHAR(20) DEFAULT 'completed',
    purchased_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 评分表
CREATE TABLE IF NOT EXISTS ratings (
    id BIGSERIAL PRIMARY KEY,
    rating_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    item_type VARCHAR(50) NOT NULL,
    item_id VARCHAR(36) NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    is_anonymous BOOLEAN DEFAULT false,
    helpful_count INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'visible',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 第三方集成表
CREATE TABLE IF NOT EXISTS integrations (
    id BIGSERIAL PRIMARY KEY,
    integration_id VARCHAR(36) NOT NULL UNIQUE,
    integration_name VARCHAR(100) NOT NULL,
    integration_type VARCHAR(50) NOT NULL, -- 'smart_home', 'health', 'ecommerce'
    provider VARCHAR(50),
    description TEXT,
    api_endpoint VARCHAR(500),
    auth_config JSONB,
    is_enabled BOOLEAN DEFAULT false,
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 智能家居设备表
CREATE TABLE IF NOT EXISTS smart_home_devices (
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

-- 智能家居触发器表
CREATE TABLE IF NOT EXISTS smart_home_triggers (
    id BIGSERIAL PRIMARY KEY,
    trigger_id VARCHAR(36) NOT NULL UNIQUE,
    household_id VARCHAR(36),
    trigger_name VARCHAR(100) NOT NULL,
    trigger_type VARCHAR(50) NOT NULL,
    conditions JSONB,
    actions JSONB,
    is_active BOOLEAN DEFAULT true,
    last_triggered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 兽医预约表
CREATE TABLE IF NOT EXISTS vet_appointments (
    id BIGSERIAL PRIMARY KEY,
    appointment_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    clinic_name VARCHAR(200),
    vet_name VARCHAR(100),
    appointment_date TIMESTAMP NOT NULL,
    appointment_type VARCHAR(50),
    reason TEXT,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'scheduled',
    reminder_sent BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 宠物医疗记录表
CREATE TABLE IF NOT EXISTS pet_medical_records (
    id BIGSERIAL PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    record_type VARCHAR(50) NOT NULL, -- 'checkup', 'vaccination', 'surgery', 'medication'
    record_date DATE NOT NULL,
    clinic_name VARCHAR(200),
    vet_name VARCHAR(100),
    diagnosis TEXT,
    treatment TEXT,
    medications JSONB,
    cost DECIMAL(10,2),
    next_appointment_date DATE,
    attachments JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 保险单表
CREATE TABLE IF NOT EXISTS insurance_policies (
    id BIGSERIAL PRIMARY KEY,
    policy_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    insurance_company VARCHAR(100) NOT NULL,
    policy_number VARCHAR(100) NOT NULL,
    plan_name VARCHAR(100),
    coverage_start DATE NOT NULL,
    coverage_end DATE,
    coverage_amount DECIMAL(10,2),
    premium DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 地图配置表
CREATE TABLE IF NOT EXISTS map_configs (
    id BIGSERIAL PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL UNIQUE,
    map_type VARCHAR(50) NOT NULL, -- 'floor_plan', 'outdoor', 'custom'
    config_data JSONB NOT NULL,
    resolution DECIMAL(6,2),
    scale DECIMAL(10,2),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
