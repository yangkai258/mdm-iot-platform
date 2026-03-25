-- 修复缺失的表 (因之前语法错误未创建)

-- 疾病模式表
CREATE TABLE IF NOT EXISTS disease_patterns (
    id BIGSERIAL PRIMARY KEY,
    pattern_id VARCHAR(36) NOT NULL UNIQUE,
    pattern_name VARCHAR(100) NOT NULL,
    disease_name VARCHAR(100) NOT NULL,
    species VARCHAR(32),
    symptoms JSONB,
    risk_factors JSONB,
    severity VARCHAR(20),
    description TEXT,
    doc_references TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
