-- ============================================================
-- Sprint 9: 动作库表
-- 创建时间: 2026-03-23
-- 描述: 存储宠物的动作定义和参数配置
-- ============================================================

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
    parameters JSONB DEFAULT '{}',
    animation_data JSONB DEFAULT '{}',
    motor_commands JSONB DEFAULT '{}',
    audio_file VARCHAR(256),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_action_library_category ON action_library(category);
CREATE INDEX IF NOT EXISTS idx_action_library_action_id ON action_library(action_id);

COMMENT ON TABLE action_library IS '动作库表';
COMMENT ON COLUMN action_library.action_id IS '动作唯一ID';
COMMENT ON COLUMN action_library.action_name IS '动作名称';
COMMENT ON COLUMN action_library.action_name_en IS '动作英文名称';
COMMENT ON COLUMN action_library.category IS '动作分类: emotion, greeting, play, utility';
COMMENT ON COLUMN action_library.description IS '动作描述';
COMMENT ON COLUMN action_library.duration_ms IS '动作持续时间(毫秒)';
COMMENT ON COLUMN action_library.priority IS '优先级 1-10';
COMMENT ON COLUMN action_library.is_emergency IS '是否紧急动作';
COMMENT ON COLUMN action_library.compatible_models IS '兼容的设备型号列表';
COMMENT ON COLUMN action_library.parameters IS '动作参数配置';
COMMENT ON COLUMN action_library.animation_data IS '动画数据';
COMMENT ON COLUMN action_library.motor_commands IS '电机指令';
COMMENT ON COLUMN action_library.audio_file IS '关联的音频文件';
