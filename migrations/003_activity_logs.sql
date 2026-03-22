-- ============================================================
-- Sprint 9: Activity Logs 审计日志迁移
-- ============================================================

-- 创建 activity_logs 表
CREATE TABLE IF NOT EXISTS activity_logs (
    id              SERIAL PRIMARY KEY,
    user_id         INT,
    username        VARCHAR(64),
    action          VARCHAR(64) NOT NULL,  -- create/update/delete/login/logout
    resource_type   VARCHAR(64) NOT NULL,  -- device/member/role/config/menu
    resource_id     INT,
    resource_name   VARCHAR(255),
    details         JSONB DEFAULT '{}',
    ip              VARCHAR(32),
    user_agent      VARCHAR(255),
    tenant_id       UUID,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_username ON activity_logs(username);
CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
CREATE INDEX IF NOT EXISTS idx_activity_logs_resource_type ON activity_logs(resource_type);
CREATE INDEX IF NOT EXISTS idx_activity_logs_resource_id ON activity_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_tenant_id ON activity_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs(created_at);

-- 创建 login_logs 表（独立登录日志）
CREATE TABLE IF NOT EXISTS login_logs (
    id          SERIAL PRIMARY KEY,
    user_id     INT,
    username    VARCHAR(64),
    ip          VARCHAR(32),
    location    VARCHAR(255),
    browser     VARCHAR(50),
    os          VARCHAR(50),
    status      INT DEFAULT 1,  -- 1:成功 0:失败
    msg         VARCHAR(255),
    tenant_id   UUID,
    login_time  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_username ON login_logs(username);
CREATE INDEX IF NOT EXISTS idx_login_logs_tenant_id ON login_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_login_time ON login_logs(login_time);

-- 迁移说明：
-- 如果表已存在但缺少 tenant_id 列，执行以下 ALTER 语句：
-- ALTER TABLE activity_logs ADD COLUMN IF NOT EXISTS tenant_id UUID;
-- ALTER TABLE login_logs ADD COLUMN IF NOT EXISTS tenant_id UUID;
