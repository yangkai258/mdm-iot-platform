-- 知识库版本管理表
CREATE TABLE IF NOT EXISTS knowledge_versions (
    id BIGSERIAL PRIMARY KEY,
    knowledge_id BIGINT NOT NULL,
    version VARCHAR(32) NOT NULL,
    content TEXT,
    change_log TEXT,
    content_hash VARCHAR(64),
    file_url VARCHAR(512),
    file_size BIGINT DEFAULT 0,
    change_type VARCHAR(20),
    status VARCHAR(20) DEFAULT 'draft',
    published_at TIMESTAMP,
    published_by VARCHAR(64),
    created_by VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_knowledge_versions_kid ON knowledge_versions(knowledge_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_versions_status ON knowledge_versions(status);

-- 知识库版本审核记录表
CREATE TABLE IF NOT EXISTS knowledge_version_reviews (
    id BIGSERIAL PRIMARY KEY,
    version_id BIGINT NOT NULL,
    review_status VARCHAR(20) DEFAULT 'pending',
    review_comment TEXT,
    reviewer_id VARCHAR(64),
    reviewer_name VARCHAR(100),
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_kv_reviews_vid ON knowledge_version_reviews(version_id);
