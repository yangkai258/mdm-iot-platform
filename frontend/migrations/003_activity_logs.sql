-- Activity Logs and Login Logs Migration
-- Created: 2026-03-21

-- Activity Logs table
CREATE TABLE IF NOT EXISTS activity_logs (
    id SERIAL PRIMARY KEY,
    user_id INT DEFAULT 0,
    username VARCHAR(64) DEFAULT '',
    action VARCHAR(64) DEFAULT '',
    resource_type VARCHAR(64) DEFAULT '',
    resource_id INT DEFAULT 0,
    details JSONB DEFAULT '{}',
    ip VARCHAR(32) DEFAULT '',
    user_agent TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Index for activity_logs
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
CREATE INDEX IF NOT EXISTS idx_activity_logs_resource_type ON activity_logs(resource_type);
CREATE INDEX IF NOT EXISTS idx_activity_logs_resource_id ON activity_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs(created_at);

-- Login Logs table
CREATE TABLE IF NOT EXISTS login_logs (
    id SERIAL PRIMARY KEY,
    user_id INT DEFAULT 0,
    username VARCHAR(64) DEFAULT '',
    status VARCHAR(32) DEFAULT '',
    ip VARCHAR(32) DEFAULT '',
    location VARCHAR(128) DEFAULT '',
    device VARCHAR(128) DEFAULT '',
    browser VARCHAR(128) DEFAULT '',
    login_at TIMESTAMP DEFAULT NOW(),
    logout_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Index for login_logs
CREATE INDEX IF NOT EXISTS idx_login_logs_user_id ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_username ON login_logs(username);
CREATE INDEX IF NOT EXISTS idx_login_logs_status ON login_logs(status);
CREATE INDEX IF NOT EXISTS idx_login_logs_login_at ON login_logs(login_at);

-- Comments
COMMENT ON TABLE activity_logs IS 'Activity/Audit logs for tracking user actions';
COMMENT ON TABLE login_logs IS 'Login and logout tracking';
