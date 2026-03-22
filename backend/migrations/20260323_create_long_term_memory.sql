-- ============================================================
-- Sprint 9: 长期记忆表
-- 创建时间: 2026-03-23
-- 描述: 存储宠物的长期记忆和知识
-- ============================================================

CREATE TABLE IF NOT EXISTS long_term_memory (
    id BIGSERIAL PRIMARY KEY,
    memory_id VARCHAR(64) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    memory_category VARCHAR(32) NOT NULL,
    content JSONB NOT NULL,
    keywords JSONB DEFAULT '[]',
    embedding JSONB,
    confidence FLOAT DEFAULT 0.8,
    reinforcement_count INT DEFAULT 1,
    last_reinforced_at TIMESTAMP,
    decay_score FLOAT DEFAULT 1.0,
    is_locked BOOLEAN DEFAULT FALSE,
    source_memory_id VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_long_term_memory_device_id ON long_term_memory(device_id);
CREATE INDEX IF NOT EXISTS idx_long_term_memory_category ON long_term_memory(memory_category);

COMMENT ON TABLE long_term_memory IS '长期记忆表';
COMMENT ON COLUMN long_term_memory.memory_id IS '记忆唯一ID';
COMMENT ON COLUMN long_term_memory.device_id IS '设备ID';
COMMENT ON COLUMN long_term_memory.user_id IS '用户ID';
COMMENT ON COLUMN long_term_memory.memory_category IS '记忆分类: preference, habit, knowledge, relationship';
COMMENT ON COLUMN long_term_memory.content IS '记忆内容JSON';
COMMENT ON COLUMN long_term_memory.keywords IS '关键词列表';
COMMENT ON COLUMN long_term_memory.embedding IS '向量嵌入数据';
COMMENT ON COLUMN long_term_memory.confidence IS '置信度 0.0-1.0';
COMMENT ON COLUMN long_term_memory.reinforcement_count IS '强化次数';
COMMENT ON COLUMN long_term_memory.last_reinforced_at IS '最后强化时间';
COMMENT ON COLUMN long_term_memory.decay_score IS '衰减分数';
COMMENT ON COLUMN long_term_memory.is_locked IS '是否锁定(不参与遗忘)';
COMMENT ON COLUMN long_term_memory.source_memory_id IS '来源记忆ID';
