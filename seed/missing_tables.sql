-- Sprint 9 & Sprint 10 缺失表

-- 宠物状态表
CREATE TABLE IF NOT EXISTS pet_status (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL UNIQUE,
    pet_name VARCHAR(32) NOT NULL DEFAULT '小爪',
    pet_type VARCHAR(16) DEFAULT 'cat',
    personality JSONB,
    appearance JSONB,
    mood INT DEFAULT 50,
    energy INT DEFAULT 100,
    hunger INT DEFAULT 0,
    position_x FLOAT DEFAULT 0,
    position_y FLOAT DEFAULT 0,
    current_expression VARCHAR(32) DEFAULT 'happy',
    current_action VARCHAR(32),
    is_online BOOLEAN DEFAULT FALSE,
    last_seen_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_pet_status_device_id ON pet_status (device_id);

-- 对话记录表
CREATE TABLE IF NOT EXISTS conversations (
    id BIGSERIAL PRIMARY KEY,
    conversation_id VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    device_id VARCHAR(64),
    title VARCHAR(256) NOT NULL,
    last_message TEXT,
    last_message_at TIMESTAMP,
    message_count INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations (user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_device_id ON conversations (device_id);

-- 对话消息表
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    message_id VARCHAR(64) NOT NULL UNIQUE,
    conversation_id VARCHAR(64) NOT NULL,
    sender_type SMALLINT NOT NULL,
    sender_id BIGINT,
    content TEXT NOT NULL,
    content_type SMALLINT DEFAULT 1,
    media_url VARCHAR(512),
    intent VARCHAR(64),
    confidence FLOAT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages (conversation_id);

-- 动作库表
CREATE TABLE IF NOT EXISTS action_library (
    id BIGSERIAL PRIMARY KEY,
    action_id VARCHAR(64) NOT NULL UNIQUE,
    action_name VARCHAR(64) NOT NULL,
    action_name_en VARCHAR(64),
    category VARCHAR(32) NOT NULL,
    description TEXT,
    duration_ms INT NOT NULL,
    priority INT DEFAULT 5,
    is_emergency BOOLEAN DEFAULT FALSE,
    compatible_models JSONB NOT NULL,
    parameters JSONB,
    animation_data JSONB,
    motor_commands JSONB,
    audio_file VARCHAR(256),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

-- 决策规则表
CREATE TABLE IF NOT EXISTS decision_rules (
    id BIGSERIAL PRIMARY KEY,
    rule_id VARCHAR(64) NOT NULL UNIQUE,
    rule_name VARCHAR(128) NOT NULL,
    rule_type VARCHAR(32) NOT NULL,
    conditions JSONB NOT NULL,
    condition_logic VARCHAR(8) DEFAULT 'AND',
    actions JSONB NOT NULL,
    priority INT DEFAULT 50,
    is_active BOOLEAN DEFAULT TRUE,
    device_id VARCHAR(64),
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_decision_rules_device_id ON decision_rules (device_id);
CREATE INDEX IF NOT EXISTS idx_decision_rules_is_active ON decision_rules (is_active);

-- 短期记忆表
CREATE TABLE IF NOT EXISTS short_term_memory (
    id BIGSERIAL PRIMARY KEY,
    memory_id VARCHAR(64) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    session_id VARCHAR(64) NOT NULL,
    message_id VARCHAR(64),
    memory_type VARCHAR(32) NOT NULL,
    content JSONB NOT NULL,
    importance FLOAT DEFAULT 0.5,
    access_count INT DEFAULT 0,
    last_accessed_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_short_term_memory_device_id ON short_term_memory (device_id);
CREATE INDEX IF NOT EXISTS idx_short_term_memory_session_id ON short_term_memory (session_id);
CREATE INDEX IF NOT EXISTS idx_short_term_memory_expires_at ON short_term_memory (expires_at);

-- 长期记忆表
CREATE TABLE IF NOT EXISTS long_term_memory (
    id BIGSERIAL PRIMARY KEY,
    memory_id VARCHAR(64) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    memory_category VARCHAR(32) NOT NULL,
    content JSONB NOT NULL,
    keywords JSONB,
    embedding JSONB,
    confidence FLOAT DEFAULT 0.8,
    reinforcement_count INT DEFAULT 1,
    last_reinforced_at TIMESTAMP,
    decay_score FLOAT DEFAULT 1.0,
    is_locked BOOLEAN DEFAULT FALSE,
    source_memory_id VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_long_term_memory_device_id ON long_term_memory (device_id);
CREATE INDEX IF NOT EXISTS idx_long_term_memory_category ON long_term_memory (memory_category);

-- 传感器事件表
CREATE TABLE IF NOT EXISTS sensor_events (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL,
    sensor_type VARCHAR(50) NOT NULL,
    value JSONB NOT NULL,
    unit VARCHAR(20),
    quality VARCHAR(20),
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_sensor_events_device_sensor_time ON sensor_events (device_id, sensor_type, timestamp DESC);

-- 设备操作日志表
CREATE TABLE IF NOT EXISTS device_operation_logs (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL,
    operator_id BIGINT NOT NULL,
    operator_type VARCHAR(20),
    operation_type VARCHAR(50) NOT NULL,
    operation_data JSONB,
    result VARCHAR(20),
    error_message TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);

CREATE INDEX IF NOT EXISTS idx_device_operation_logs_device_id ON device_operation_logs (device_id);
CREATE INDEX IF NOT EXISTS idx_device_operation_logs_operator_id ON device_operation_logs (operator_id);
CREATE INDEX IF NOT EXISTS idx_device_operation_logs_created_at ON device_operation_logs (created_at DESC);

-- 告警规则模板表
CREATE TABLE IF NOT EXISTS alert_templates (
    id BIGSERIAL PRIMARY KEY,
    template_name VARCHAR(100) NOT NULL,
    alert_type VARCHAR(50) NOT NULL,
    condition_expr TEXT NOT NULL,
    threshold_value VARCHAR(100),
    severity INT DEFAULT 2,
    notify_ways VARCHAR(50)[],
    remark TEXT,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    tenant_id UUID
);
