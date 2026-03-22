-- Sprint 14: AI 行为监控和模型管理
-- Migration: 20260322_create_ai_tables.sql

-- AI 行为日志表
CREATE TABLE IF NOT EXISTS ai_behavior_logs (
    id              BIGSERIAL PRIMARY KEY,
    log_id          VARCHAR(64) NOT NULL UNIQUE,
    device_id       VARCHAR(64),
    user_id         BIGINT,
    model_name      VARCHAR(64),
    model_version   VARCHAR(32),
    event_type      VARCHAR(32), -- inference/abnormal/rollback
    input_summary   TEXT,
    output_summary  TEXT,
    latency_ms      INT,
    error_rate      DECIMAL(5,4) DEFAULT 0,
    confidence      DECIMAL(5,4) DEFAULT 0,
    status          VARCHAR(20), -- success/failed/anomaly
    error_msg       VARCHAR(512),
    metadata        JSONB,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_behavior_device_id ON ai_behavior_logs(device_id);
CREATE INDEX IF NOT EXISTS idx_ai_behavior_user_id ON ai_behavior_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_behavior_event_type ON ai_behavior_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_ai_behavior_status ON ai_behavior_logs(status);
CREATE INDEX IF NOT EXISTS idx_ai_behavior_created_at ON ai_behavior_logs(created_at DESC);

-- AI 模型版本表
CREATE TABLE IF NOT EXISTS ai_model_versions (
    id              BIGSERIAL PRIMARY KEY,
    model_id        VARCHAR(64) NOT NULL,
    model_name      VARCHAR(64),
    version         VARCHAR(32) NOT NULL,
    status          VARCHAR(20), -- testing/staging/production/deprecated
    model_path      VARCHAR(512),
    config          JSONB,
    metrics         JSONB,
    rollback_from   VARCHAR(32),
    published_by    BIGINT,
    published_at    TIMESTAMP,
    deprecated_at   TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(model_id, version)
);

CREATE INDEX IF NOT EXISTS idx_ai_model_model_id ON ai_model_versions(model_id);
CREATE INDEX IF NOT EXISTS idx_ai_model_status ON ai_model_versions(status);
CREATE INDEX IF NOT EXISTS idx_ai_model_created_at ON ai_model_versions(created_at DESC);

-- AI 沙箱测试任务表
CREATE TABLE IF NOT EXISTS ai_sandbox_tests (
    id              BIGSERIAL PRIMARY KEY,
    task_id         VARCHAR(64) NOT NULL UNIQUE,
    model_id        VARCHAR(64),
    test_data_id    VARCHAR(64),
    test_type       VARCHAR(30), -- unit/integration/stress
    test_name       VARCHAR(100),
    test_cases      JSONB NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending', -- pending/running/completed/failed
    result          JSONB,
    report_path     VARCHAR(512),
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_sandbox_model_id ON ai_sandbox_tests(model_id);
CREATE INDEX IF NOT EXISTS idx_ai_sandbox_status ON ai_sandbox_tests(status);
CREATE INDEX IF NOT EXISTS idx_ai_sandbox_created_at ON ai_sandbox_tests(created_at DESC);

-- AI 回滚任务表
CREATE TABLE IF NOT EXISTS ai_rollback_tasks (
    id              BIGSERIAL PRIMARY KEY,
    task_id         VARCHAR(64) NOT NULL UNIQUE,
    model_id        VARCHAR(64),
    from_version    VARCHAR(32),
    to_version      VARCHAR(32),
    reason          TEXT,
    status          VARCHAR(20), -- pending/in_progress/completed/failed
    triggered_by    BIGINT,
    affected_count  INT DEFAULT 0,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_rollback_task_id ON ai_rollback_tasks(task_id);
CREATE INDEX IF NOT EXISTS idx_ai_rollback_model_id ON ai_rollback_tasks(model_id);
CREATE INDEX IF NOT EXISTS idx_ai_rollback_status ON ai_rollback_tasks(status);
CREATE INDEX IF NOT EXISTS idx_ai_rollback_created_at ON ai_rollback_tasks(created_at DESC);
