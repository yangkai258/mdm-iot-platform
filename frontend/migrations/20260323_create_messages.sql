-- ============================================================
-- Sprint 9: 消息表
-- 创建时间: 2026-03-23
-- 描述: 存储会话中的消息记录
-- ============================================================

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
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);

COMMENT ON TABLE messages IS '消息表';
COMMENT ON COLUMN messages.message_id IS '消息唯一ID';
COMMENT ON COLUMN messages.conversation_id IS '所属会话ID';
COMMENT ON COLUMN messages.sender_type IS '发送者类型: 1=用户, 2=宠物';
COMMENT ON COLUMN messages.sender_id IS '发送者ID';
COMMENT ON COLUMN messages.content IS '消息内容';
COMMENT ON COLUMN messages.content_type IS '内容类型: 1=text, 2=voice, 3=image, 4=action';
COMMENT ON COLUMN messages.media_url IS '媒体文件URL';
COMMENT ON COLUMN messages.intent IS '意图识别结果';
COMMENT ON COLUMN messages.confidence IS '置信度';
