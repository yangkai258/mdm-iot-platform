-- 会员360度画像表
CREATE TABLE IF NOT EXISTS member_360_profiles (
    id BIGSERIAL PRIMARY KEY,
    member_id VARCHAR(64) UNIQUE NOT NULL,
    member_name VARCHAR(100),
    age INTEGER DEFAULT 0,
    gender VARCHAR(10),
    location VARCHAR(200),
    occupation VARCHAR(100),
    interests TEXT,
    total_spend DECIMAL(12,2) DEFAULT 0,
    avg_order_value DECIMAL(12,2) DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    last_order_date TIMESTAMP,
    preferred_payment VARCHAR(50),
    price_sensitivity DECIMAL(5,2) DEFAULT 0,
    login_frequency INTEGER DEFAULT 0,
    avg_session_duration INTEGER DEFAULT 0,
    active_days_per_week INTEGER DEFAULT 0,
    feature_usage_map TEXT,
    pet_engagement DECIMAL(5,2) DEFAULT 0,
    preferred_pet_type VARCHAR(50),
    preferred_ai_style VARCHAR(50),
    notification_pref VARCHAR(20),
    member_level VARCHAR(20),
    member_since TIMESTAMP,
    churn_risk DECIMAL(5,2) DEFAULT 0,
    ltv DECIMAL(12,2) DEFAULT 0,
    tags TEXT,
    custom_labels TEXT,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mp_member ON member_360_profiles(member_id);

-- 会员画像洞察表
CREATE TABLE IF NOT EXISTS member_portrait_insights (
    id BIGSERIAL PRIMARY KEY,
    member_id VARCHAR(64) NOT NULL,
    insight_type VARCHAR(50),
    title VARCHAR(200),
    description TEXT,
    confidence DECIMAL(5,2),
    action TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mpi_member ON member_portrait_insights(member_id);

-- 动作模仿学习进度表
CREATE TABLE IF NOT EXISTS action_learning_progress (
    id BIGSERIAL PRIMARY KEY,
    pet_id VARCHAR(64) NOT NULL,
    member_id VARCHAR(64) NOT NULL,
    action_id BIGINT NOT NULL,
    action_name VARCHAR(100),
    level INTEGER DEFAULT 1,
    exp INTEGER DEFAULT 0,
    exp_to_next_level INTEGER DEFAULT 100,
    mastery_rate DECIMAL(5,2) DEFAULT 0,
    practice_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    total_duration INTEGER DEFAULT 0,
    avg_accuracy DECIMAL(5,2) DEFAULT 0,
    last_practice_at TIMESTAMP,
    next_practice_at TIMESTAMP,
    learned_at TIMESTAMP,
    mastered_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'learning',
    star_rating INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_alp_pet ON action_learning_progress(pet_id);
CREATE INDEX IF NOT EXISTS idx_alp_member ON action_learning_progress(member_id);

-- 动作模仿练习记录表
CREATE TABLE IF NOT EXISTS action_learning_sessions (
    id BIGSERIAL PRIMARY KEY,
    progress_id BIGINT NOT NULL,
    pet_id VARCHAR(64) NOT NULL,
    action_id BIGINT NOT NULL,
    action_name VARCHAR(100),
    duration INTEGER DEFAULT 0,
    accuracy DECIMAL(5,2),
    is_success BOOLEAN DEFAULT FALSE,
    score INTEGER DEFAULT 0,
    exp_gained INTEGER DEFAULT 0,
    feedback TEXT,
    video_url VARCHAR(512),
    screenshot_url VARCHAR(512),
    session_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_als_progress ON action_learning_sessions(progress_id);

-- 内容版本表
CREATE TABLE IF NOT EXISTS content_versions (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL,
    version VARCHAR(32) NOT NULL,
    version_num INTEGER DEFAULT 1,
    title VARCHAR(200),
    description TEXT,
    file_url VARCHAR(512),
    file_size BIGINT DEFAULT 0,
    content_hash VARCHAR(64),
    thumbnail_url VARCHAR(512),
    change_log TEXT,
    change_type VARCHAR(20),
    change_size BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft',
    is_latest BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP,
    published_by VARCHAR(64),
    created_by VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cv_content ON content_versions(content_id);
CREATE INDEX IF NOT EXISTS idx_cv_status ON content_versions(status);

-- 内容版本审核表
CREATE TABLE IF NOT EXISTS content_version_reviews (
    id BIGSERIAL PRIMARY KEY,
    version_id BIGINT NOT NULL,
    review_type VARCHAR(20),
    review_status VARCHAR(20) DEFAULT 'pending',
    review_score DECIMAL(5,2),
    review_report TEXT,
    issues TEXT,
    reviewer_id VARCHAR(64),
    reviewer_name VARCHAR(100),
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cvr_version ON content_version_reviews(version_id);
