-- ============================================================
-- Sprint 9: 短期记忆表
-- 创建时间: 2026-03-23
-- 描述: 存储宠物的短期会话记忆
-- ============================================================

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
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_short_term_memory_device_id ON short_term_memory(device_id);
CREATE INDEX IF NOT EXISTS idx_short_term_memory_session_id ON short_term_memory(session_id);
CREATE INDEX IF NOT EXISTS idx_short_term_memory_expires_at ON short_term_memory(expires_at);

COMMENT ON TABLE short_term_memory IS '短期记忆表';
COMMENT ON COLUMN short_term_memory.memory_id IS '记忆唯一ID';
COMMENT ON COLUMN short_term_memory.device_id IS '设备ID';
COMMENT ON COLUMN short_term_memory.user_id IS '用户ID';
COMMENT ON COLUMN short_term_memory.session_id IS '会话ID';
COMMENT ON COLUMN short_term_memory.message_id IS '关联消息ID';
COMMENT ON COLUMN short_term_memory.memory_type IS '记忆类型: interaction, conversation, preference, context';
COMMENT ON COLUMN short_term_memory.content IS '记忆内容JSON';
COMMENT ON COLUMN short_term_memory.importance IS '重要性 0.0-1.0';
COMMENT ON COLUMN short_term_memory.access_count IS '访问次数';
COMMENT ON COLUMN short_term_memory.last_accessed_at IS '最后访问时间';
COMMENT ON COLUMN short_term_memory.expires_at IS '过期时间';
