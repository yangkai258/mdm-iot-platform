-- ============================================================
-- Sprint 9: 会话表
-- 创建时间: 2026-03-23
-- 描述: 存储用户与宠物之间的对话会话
-- ============================================================

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
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_device_id ON conversations(device_id);

COMMENT ON TABLE conversations IS '对话会话表';
COMMENT ON COLUMN conversations.conversation_id IS '会话唯一ID';
COMMENT ON COLUMN conversations.user_id IS '用户ID';
COMMENT ON COLUMN conversations.device_id IS '设备ID';
COMMENT ON COLUMN conversations.title IS '会话标题';
COMMENT ON COLUMN conversations.last_message IS '最后一条消息';
COMMENT ON COLUMN conversations.last_message_at IS '最后消息时间';
COMMENT ON COLUMN conversations.message_count IS '消息数量';
COMMENT ON COLUMN conversations.status IS '状态: 1=正常, 2=已归档';
