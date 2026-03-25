-- ============================================
-- Sprint 28: 宠物社交
-- ============================================

CREATE TABLE IF NOT EXISTS social_share_records (
    id BIGSERIAL PRIMARY KEY,
    share_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    share_type VARCHAR(50) NOT NULL,
    content_id VARCHAR(36),
    content_type VARCHAR(50),
    platform VARCHAR(50),
    share_url VARCHAR(500),
    click_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    shared_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pet_social_interactions (
    id BIGSERIAL PRIMARY KEY,
    interaction_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    target_pet_id VARCHAR(36),
    interaction_type VARCHAR(50) NOT NULL,
    interaction_data JSONB,
    is_positive BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 30: 健康分析
-- ============================================

CREATE TABLE IF NOT EXISTS sleep_analysis (
    id BIGSERIAL PRIMARY KEY,
    analysis_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    analysis_date DATE NOT NULL,
    total_sleep_minutes INT,
    deep_sleep_minutes INT,
    light_sleep_minutes INT,
    rem_sleep_minutes INT,
    sleep_quality_score DECIMAL(5,2),
    sleep_patterns JSONB,
    recommendations JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS diet_meals (
    id BIGSERIAL PRIMARY KEY,
    meal_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    meal_type VARCHAR(20) NOT NULL,
    food_name VARCHAR(200) NOT NULL,
    food_amount DECIMAL(10,2),
    unit VARCHAR(20),
    calories DECIMAL(10,4),
    nutrients JSONB,
    feeding_time TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS nutrition_reports (
    id BIGSERIAL PRIMARY KEY,
    report_id VARCHAR(36) NOT NULL UNIQUE,
    pet_id VARCHAR(36) NOT NULL,
    report_type VARCHAR(50) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    daily_calories_avg DECIMAL(10,2),
    nutrients_summary JSONB,
    diet_evaluation JSONB,
    recommendations JSONB,
    generated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 31: 数据集开放
-- ============================================

CREATE TABLE IF NOT EXISTS research_datasets (
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
    access_level VARCHAR(20) DEFAULT 'restricted',
    license_type VARCHAR(50),
    owner_id BIGINT,
    citation TEXT,
    published_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dataset_versions (
    id BIGSERIAL PRIMARY KEY,
    version_id VARCHAR(36) NOT NULL UNIQUE,
    dataset_id VARCHAR(36) NOT NULL,
    version_number VARCHAR(20) NOT NULL,
    changes TEXT,
    diff_summary JSONB,
    new_records INT,
    updated_records INT,
    deleted_records INT,
    storage_path VARCHAR(500),
    size_bytes BIGINT,
    is_current BOOLEAN DEFAULT false,
    released_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dataset_access_requests (
    id BIGSERIAL PRIMARY KEY,
    request_id VARCHAR(36) NOT NULL UNIQUE,
    dataset_id VARCHAR(36) NOT NULL,
    applicant_id BIGINT NOT NULL,
    access_type VARCHAR(50) NOT NULL,
    purpose TEXT,
    research_plan TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    reviewed_by BIGINT,
    reviewed_at TIMESTAMP,
    review_comment TEXT,
    granted_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dataset_download_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    dataset_id VARCHAR(36) NOT NULL,
    user_id BIGINT NOT NULL,
    download_size_bytes BIGINT,
    download_method VARCHAR(20),
    ip_address VARCHAR(45),
    user_agent TEXT,
    downloaded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Sprint 32: AI研究平台
-- ============================================

CREATE TABLE IF NOT EXISTS research_platforms (
    id BIGSERIAL PRIMARY KEY,
    platform_id VARCHAR(36) NOT NULL UNIQUE,
    platform_name VARCHAR(200) NOT NULL,
    platform_type VARCHAR(50) NOT NULL,
    description TEXT,
    endpoint_url VARCHAR(500),
    auth_config JSONB,
    is_enabled BOOLEAN DEFAULT false,
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS research_experiments (
    id BIGSERIAL PRIMARY KEY,
    experiment_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_name VARCHAR(200) NOT NULL,
    platform_id VARCHAR(36),
    hypothesis TEXT,
    config JSONB,
    status VARCHAR(20) DEFAULT 'draft',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    results JSONB,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS research_collaborators (
    id BIGSERIAL PRIMARY KEY,
    collaborator_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_id VARCHAR(36) NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(50) DEFAULT 'member',
    permissions JSONB,
    invited_at TIMESTAMP,
    joined_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS research_experiment_results (
    id BIGSERIAL PRIMARY KEY,
    result_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_id VARCHAR(36) NOT NULL,
    result_type VARCHAR(50) NOT NULL,
    metrics JSONB,
    visualizations JSONB,
    conclusion TEXT,
    is_significant BOOLEAN,
    peer_review_status VARCHAR(20),
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS research_hypotheses (
    id BIGSERIAL PRIMARY KEY,
    hypothesis_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_id VARCHAR(36),
    hypothesis_text TEXT NOT NULL,
    variables JSONB,
    prediction TEXT,
    status VARCHAR(20) DEFAULT 'proposed',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS research_institutions (
    id BIGSERIAL PRIMARY KEY,
    institution_id VARCHAR(36) NOT NULL UNIQUE,
    institution_name VARCHAR(200) NOT NULL,
    institution_type VARCHAR(50),
    country VARCHAR(50),
    website VARCHAR(255),
    contact_email VARCHAR(100),
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS research_publications (
    id BIGSERIAL PRIMARY KEY,
    publication_id VARCHAR(36) NOT NULL UNIQUE,
    title VARCHAR(500) NOT NULL,
    authors JSONB,
    institution_id VARCHAR(36),
    publication_type VARCHAR(50),
    journal VARCHAR(200),
    doi VARCHAR(100),
    abstract TEXT,
    keywords VARCHAR(50)[],
    citation_count INT DEFAULT 0,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS research_discussions (
    id BIGSERIAL PRIMARY KEY,
    discussion_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_id VARCHAR(36),
    user_id BIGINT NOT NULL,
    discussion_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    parent_id BIGINT,
    upvotes INT DEFAULT 0,
    is_answered BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Module AI Engineering
-- ============================================

CREATE TABLE IF NOT EXISTS ai_model_shards (
    id BIGSERIAL PRIMARY KEY,
    shard_id VARCHAR(36) NOT NULL UNIQUE,
    model_id VARCHAR(64) NOT NULL,
    shard_index INT NOT NULL,
    total_shards INT NOT NULL,
    shard_size_bytes BIGINT,
    shard_path VARCHAR(500),
    checksum VARCHAR(64),
    version VARCHAR(20),
    is_loaded BOOLEAN DEFAULT false,
    loaded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_model_versions (
    id BIGSERIAL PRIMARY KEY,
    version_id VARCHAR(36) NOT NULL UNIQUE,
    model_id VARCHAR(64) NOT NULL,
    version_number VARCHAR(20) NOT NULL,
    changelog TEXT,
    performance_metrics JSONB,
    is_stable BOOLEAN DEFAULT false,
    is_current BOOLEAN DEFAULT false,
    deprecated_at TIMESTAMP,
    released_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Module Simulation
-- ============================================

CREATE TABLE IF NOT EXISTS simulation_test_cases (
    id BIGSERIAL PRIMARY KEY,
    case_id VARCHAR(36) NOT NULL UNIQUE,
    case_name VARCHAR(200) NOT NULL,
    case_type VARCHAR(50) NOT NULL,
    scenario_id VARCHAR(36),
    preconditions JSONB,
    test_steps JSONB,
    expected_results JSONB,
    priority INT DEFAULT 3,
    automated BOOLEAN DEFAULT false,
    script_path VARCHAR(500),
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_test_executions (
    id BIGSERIAL PRIMARY KEY,
    execution_id VARCHAR(36) NOT NULL UNIQUE,
    testcase_id VARCHAR(36) NOT NULL,
    execution_env VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration_ms INT,
    actual_results JSONB,
    pass_count INT DEFAULT 0,
    fail_count INT DEFAULT 0,
    error_message TEXT,
    logs_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_test_reports (
    id BIGSERIAL PRIMARY KEY,
    report_id VARCHAR(36) NOT NULL UNIQUE,
    execution_id VARCHAR(36),
    report_type VARCHAR(50) NOT NULL,
    summary JSONB,
    details JSONB,
    metrics JSONB,
    failure_analysis TEXT,
    recommendations TEXT,
    generated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_playbacks (
    id BIGSERIAL PRIMARY KEY,
    playback_id VARCHAR(36) NOT NULL UNIQUE,
    simulation_id VARCHAR(36),
    session_id VARCHAR(64),
    recording_data JSONB NOT NULL,
    duration_seconds INT,
    playback_speed DECIMAL(3,2) DEFAULT 1.0,
    status VARCHAR(20) DEFAULT 'ready',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_stress_tests (
    id BIGSERIAL PRIMARY KEY,
    test_id VARCHAR(36) NOT NULL UNIQUE,
    test_name VARCHAR(200) NOT NULL,
    test_type VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    target_metrics JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    current_concurrent_users INT DEFAULT 0,
    target_concurrent_users INT,
    results JSONB,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS simulation_virtual_pets (
    id BIGSERIAL PRIMARY KEY,
    pet_id VARCHAR(36) NOT NULL UNIQUE,
    simulation_id VARCHAR(36),
    virtual_pet_data JSONB,
    initial_state JSONB,
    personality_params JSONB,
    behavior_model VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Legacy/User/Permission Tables
-- ============================================

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20),
    password_hash VARCHAR(255),
    real_name VARCHAR(100),
    avatar_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active',
    last_login_at TIMESTAMP,
    tenant_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    permission_id VARCHAR(36) NOT NULL UNIQUE,
    permission_name VARCHAR(100) NOT NULL,
    permission_code VARCHAR(100) NOT NULL UNIQUE,
    resource_type VARCHAR(50),
    description TEXT,
    parent_id BIGINT,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    tenant_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    operation_type VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id VARCHAR(100),
    request_data JSONB,
    response_data JSONB,
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    status VARCHAR(20),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ota_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id VARCHAR(36) NOT NULL UNIQUE,
    task_name VARCHAR(200) NOT NULL,
    device_id VARCHAR(64) NOT NULL,
    target_version VARCHAR(50),
    task_type VARCHAR(20) DEFAULT 'ota',
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
