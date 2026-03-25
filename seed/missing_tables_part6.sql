-- ============================================
-- Simulation 仿真测试模块
-- ============================================

-- 仿真宠物表
CREATE TABLE IF NOT EXISTS simulation_pets (
    id BIGSERIAL PRIMARY KEY,
    pet_id VARCHAR(36) NOT NULL UNIQUE,
    simulation_id VARCHAR(36),
    virtual_pet_data JSONB,
    initial_state JSONB,
    personality_params JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 仿真场景表
CREATE TABLE IF NOT EXISTS simulation_scenarios (
    id BIGSERIAL PRIMARY KEY,
    scenario_id VARCHAR(36) NOT NULL UNIQUE,
    scenario_name VARCHAR(200) NOT NULL,
    scenario_type VARCHAR(50) NOT NULL,
    description TEXT,
    environment_config JSONB,
    pet_configs JSONB,
    duration_minutes INT,
    expected_outcomes JSONB,
    difficulty_level INT DEFAULT 1,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 测试执行记录表
CREATE TABLE IF NOT EXISTS test_executions (
    id BIGSERIAL PRIMARY KEY,
    execution_id VARCHAR(36) NOT NULL UNIQUE,
    testcase_id VARCHAR(36),
    scenario_id VARCHAR(36),
    execution_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration_ms INT,
    result_data JSONB,
    error_message TEXT,
    logs_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 测试报告表
CREATE TABLE IF NOT EXISTS test_reports (
    id BIGSERIAL PRIMARY KEY,
    report_id VARCHAR(36) NOT NULL UNIQUE,
    execution_id VARCHAR(36),
    report_type VARCHAR(50) NOT NULL,
    summary JSONB,
    details JSONB,
    metrics JSONB,
    charts JSONB,
    recommendations TEXT,
    generated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 回放记录表
CREATE TABLE IF NOT EXISTS playback_records (
    id BIGSERIAL PRIMARY KEY,
    playback_id VARCHAR(36) NOT NULL UNIQUE,
    simulation_id VARCHAR(36),
    session_data JSONB NOT NULL,
    duration_seconds INT,
    recording_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'recording',
    created_at TIMESTAMP DEFAULT NOW()
);

-- 压力测试表
CREATE TABLE IF NOT EXISTS stress_tests (
    id BIGSERIAL PRIMARY KEY,
    test_id VARCHAR(36) NOT NULL UNIQUE,
    test_name VARCHAR(200) NOT NULL,
    test_config JSONB NOT NULL,
    target_metrics JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    current_load INT DEFAULT 0,
    target_load INT,
    results JSONB,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 仿真数据集表
CREATE TABLE IF NOT EXISTS simulation_datasets (
    id BIGSERIAL PRIMARY KEY,
    dataset_id VARCHAR(36) NOT NULL UNIQUE,
    dataset_name VARCHAR(200) NOT NULL,
    dataset_type VARCHAR(50) NOT NULL,
    description TEXT,
    data_schema JSONB,
    storage_path VARCHAR(500),
    size_bytes BIGINT,
    record_count INT,
    version VARCHAR(20) DEFAULT '1.0',
    status VARCHAR(20) DEFAULT 'active',
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 仿真数据集版本表
CREATE TABLE IF NOT EXISTS simulation_dataset_versions (
    id BIGSERIAL PRIMARY KEY,
    version_id VARCHAR(36) NOT NULL UNIQUE,
    dataset_id VARCHAR(36) NOT NULL,
    version_number VARCHAR(20) NOT NULL,
    changes TEXT,
    storage_path VARCHAR(500),
    record_count INT,
    is_current BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 仿真集成表
CREATE TABLE IF NOT EXISTS simulation_integrations (
    id BIGSERIAL PRIMARY KEY,
    integration_id VARCHAR(36) NOT NULL UNIQUE,
    integration_name VARCHAR(200) NOT NULL,
    integration_type VARCHAR(50) NOT NULL,
    endpoint_url VARCHAR(500),
    auth_config JSONB,
    is_enabled BOOLEAN DEFAULT false,
    last_test_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- CI/CD任务表
CREATE TABLE IF NOT EXISTS simulation_cicd_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(36) NOT NULL UNIQUE,
    job_name VARCHAR(200) NOT NULL,
    job_type VARCHAR(50) NOT NULL,
    pipeline_config JSONB,
    trigger_type VARCHAR(20),
    trigger_config JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    build_number INT,
    commit_sha VARCHAR(40),
    branch VARCHAR(100),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    duration_ms INT,
    artifacts JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 数据集任务表
CREATE TABLE IF NOT EXISTS dataset_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(36) NOT NULL UNIQUE,
    job_name VARCHAR(200) NOT NULL,
    dataset_id VARCHAR(36),
    operation VARCHAR(50) NOT NULL,
    config JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    records_processed INT DEFAULT 0,
    error_count INT DEFAULT 0,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AB实验表
CREATE TABLE IF NOT EXISTS ab_experiments (
    id BIGSERIAL PRIMARY KEY,
    experiment_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_name VARCHAR(200) NOT NULL,
    hypothesis TEXT,
    description TEXT,
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

-- 数据集样本表
CREATE TABLE IF NOT EXISTS dataset_samples (
    id BIGSERIAL PRIMARY KEY,
    sample_id VARCHAR(36) NOT NULL UNIQUE,
    dataset_id VARCHAR(36) NOT NULL,
    sample_data JSONB NOT NULL,
    label VARCHAR(100),
    features JSONB,
    split_type VARCHAR(20) DEFAULT 'train', -- 'train', 'validation', 'test'
    quality_score DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- PRD_28: 宠物电商
-- ============================================

-- 推荐产品表
CREATE TABLE IF NOT EXISTS eco_recommended_products (
    id BIGSERIAL PRIMARY KEY,
    recommendation_id VARCHAR(36) NOT NULL UNIQUE,
    product_id VARCHAR(36) NOT NULL,
    pet_id VARCHAR(36),
    user_id BIGINT,
    recommendation_type VARCHAR(50),
    score DECIMAL(5,4),
    reason TEXT,
    is_clicked BOOLEAN DEFAULT false,
    is_purchased BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 购物车表
CREATE TABLE IF NOT EXISTS eco_cart (
    id BIGSERIAL PRIMARY KEY,
    cart_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    quantity INT DEFAULT 1,
    price DECIMAL(10,2),
    added_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 电商订单表
CREATE TABLE IF NOT EXISTS eco_orders (
    id BIGSERIAL PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL UNIQUE,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    shipping_amount DECIMAL(10,2) DEFAULT 0,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    final_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    shipping_address JSONB,
    payment_method VARCHAR(20),
    paid_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- PRD_30: 家庭相册增强
-- ============================================

-- 相册照片表
CREATE TABLE IF NOT EXISTS album_photos (
    id BIGSERIAL PRIMARY KEY,
    photo_id VARCHAR(36) NOT NULL UNIQUE,
    album_id VARCHAR(36) NOT NULL,
    uploader_id BIGINT NOT NULL,
    photo_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    caption TEXT,
    location JSONB,
    taken_at TIMESTAMP,
    ai_tags VARCHAR(50)[],
    face_count INT DEFAULT 0,
    is_public BOOLEAN DEFAULT false,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AI相册表
CREATE TABLE IF NOT EXISTS album_ai_albums (
    id BIGSERIAL PRIMARY KEY,
    ai_album_id VARCHAR(36) NOT NULL UNIQUE,
    ai_album_name VARCHAR(200) NOT NULL,
    album_type VARCHAR(50) NOT NULL, -- 'face', 'location', 'date', 'event'
    criteria JSONB,
    cover_photo_url VARCHAR(500),
    photo_count INT DEFAULT 0,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 相册分类表
CREATE TABLE IF NOT EXISTS album_categories (
    id BIGSERIAL PRIMARY KEY,
    category_id VARCHAR(36) NOT NULL UNIQUE,
    category_name VARCHAR(100) NOT NULL,
    parent_id VARCHAR(36),
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 相册家庭成员表
CREATE TABLE IF NOT EXISTS album_family_members (
    id BIGSERIAL PRIMARY KEY,
    member_id VARCHAR(36) NOT NULL UNIQUE,
    album_id VARCHAR(36) NOT NULL,
    user_id BIGINT NOT NULL,
    face_encoding JSONB,
    face_photo_url VARCHAR(500),
    name VARCHAR(100),
    relation VARCHAR(50),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);
