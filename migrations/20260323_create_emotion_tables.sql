-- Sprint 17: 情感计算模块数据库迁移
-- 情绪记录表
CREATE TABLE IF NOT EXISTS emotion_records (
    id SERIAL PRIMARY KEY,
    subject_type VARCHAR(20) NOT NULL,
    subject_id INTEGER NOT NULL,
    emotion_type VARCHAR(30) NOT NULL,
    intensity DOUBLE PRECISION NOT NULL,
    source VARCHAR(20),
    confidence DOUBLE PRECISION DEFAULT 0,
    context JSONB,
    trigger_event TEXT,
    tags TEXT[],
    note TEXT,
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_emotion_records_subject ON emotion_records(subject_type, subject_id);
CREATE INDEX IF NOT EXISTS idx_emotion_records_emotion ON emotion_records(emotion_type);
CREATE INDEX IF NOT EXISTS idx_emotion_records_recorded ON emotion_records(recorded_at);

-- 宠物情绪动作表
CREATE TABLE IF NOT EXISTS pet_emotion_actions (
    id SERIAL PRIMARY KEY,
    emotion_type VARCHAR(30) NOT NULL,
    action_name VARCHAR(128) NOT NULL,
    action_code VARCHAR(64) NOT NULL UNIQUE,
    description TEXT,
    parameters JSONB,
    priority INTEGER DEFAULT 0,
    min_intensity DOUBLE PRECISION DEFAULT 0,
    max_intensity DOUBLE PRECISION DEFAULT 100,
    duration INTEGER DEFAULT 0,
    enabled BOOLEAN DEFAULT TRUE,
    icon VARCHAR(64),
    animation_url VARCHAR(255),
    sound_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pet_emotion_actions_emotion ON pet_emotion_actions(emotion_type);

-- 情绪响应配置表
CREATE TABLE IF NOT EXISTS emotion_response_configs (
    id SERIAL PRIMARY KEY,
    pet_id INTEGER NOT NULL,
    emotion_type VARCHAR(30) NOT NULL,
    strategy VARCHAR(20) NOT NULL,
    action_code VARCHAR(64),
    action_param JSONB,
    response_delay INTEGER DEFAULT 0,
    enabled BOOLEAN DEFAULT TRUE,
    threshold DOUBLE PRECISION DEFAULT 30,
    cooldown INTEGER DEFAULT 60000,
    last_triggered TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_emotion_response_configs_pet ON emotion_response_configs(pet_id);

-- 情绪报告表
CREATE TABLE IF NOT EXISTS emotion_reports (
    id SERIAL PRIMARY KEY,
    pet_id INTEGER NOT NULL,
    report_type VARCHAR(20) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    summary JSONB,
    emotion_stats JSONB,
    interaction_stats JSONB,
    trend_analysis JSONB,
    recommendations JSONB,
    generated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_emotion_reports_pet ON emotion_reports(pet_id);
