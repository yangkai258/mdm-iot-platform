-- Migration: 20260323_create_batch_tasks
-- Description: 创建批量任务表

CREATE TABLE IF NOT EXISTS batch_tasks (
    id              BIGSERIAL PRIMARY KEY,
    task_id         VARCHAR(64) NOT NULL UNIQUE,
    task_type       VARCHAR(32) NOT NULL,
    total           INT NOT NULL DEFAULT 0,
    success         INT DEFAULT 0,
    failed          INT DEFAULT 0,
    pending         INT DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    results         JSONB DEFAULT '[]',
    creator_id      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    completed_at    TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_batch_tasks_task_id ON batch_tasks(task_id);
CREATE INDEX IF NOT EXISTS idx_batch_tasks_task_type ON batch_tasks(task_type);
CREATE INDEX IF NOT EXISTS idx_batch_tasks_status ON batch_tasks(status);
CREATE INDEX IF NOT EXISTS idx_batch_tasks_creator_id ON batch_tasks(creator_id);
CREATE INDEX IF NOT EXISTS idx_batch_tasks_created_at ON batch_tasks(created_at DESC);
