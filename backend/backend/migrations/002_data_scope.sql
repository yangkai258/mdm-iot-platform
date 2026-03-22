-- 数据权限范围表
-- 用于配置角色的数据访问范围

CREATE TABLE IF NOT EXISTS data_scopes (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL DEFAULT 0,
    scope_type VARCHAR(50) NOT NULL DEFAULT 'all', -- all/org/org_and_children/self/custom
    scope_value VARCHAR(255), -- custom时存储具体条件，如 "dept_id:1,2,3"
    tenant_id VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_data_scope_role FOREIGN KEY (role_id) REFERENCES sys_roles(id) ON DELETE CASCADE,
    INDEX idx_data_scope_role (role_id),
    INDEX idx_data_scope_tenant (tenant_id)
);

-- COMMENT
COMMENT ON TABLE data_scopes IS '数据权限配置表';
COMMENT ON COLUMN data_scopes.role_id IS '角色ID';
COMMENT ON COLUMN data_scopes.scope_type IS '权限范围类型: all(全部), org(本部门), org_and_children(本部门及下级), self(仅本人), custom(自定义)';
COMMENT ON COLUMN data_scopes.scope_value IS '自定义条件，如 dept_id:1,2,3';
COMMENT ON COLUMN data_scopes.tenant_id IS '租户ID';

-- 为现有业务表添加 org_id 和 create_user_id 字段（如尚未存在）
-- 注意：这些字段会在 GORM AutoMigrate 时自动创建

-- 设备表添加数据权限字段
ALTER TABLE devices ADD COLUMN IF NOT EXISTS org_id BIGINT DEFAULT 0;
ALTER TABLE devices ADD COLUMN IF NOT EXISTS create_user_id BIGINT DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_devices_org ON devices(org_id);
CREATE INDEX IF NOT EXISTS idx_devices_create_user ON devices(create_user_id);

-- OTA部署表添加数据权限字段
ALTER TABLE ota_deployments ADD COLUMN IF NOT EXISTS create_user_id BIGINT DEFAULT 0;
ALTER TABLE ota_deployments ADD COLUMN IF NOT EXISTS org_id BIGINT DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_ota_deployments_org ON ota_deployments(org_id);
CREATE INDEX IF NOT EXISTS idx_ota_deployments_create_user ON ota_deployments(create_user_id);

-- 告警表添加数据权限字段
ALTER TABLE device_alerts ADD COLUMN IF NOT EXISTS org_id BIGINT DEFAULT 0;
ALTER TABLE device_alerts ADD COLUMN IF NOT EXISTS create_user_id BIGINT DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_device_alerts_org ON device_alerts(org_id);
CREATE INDEX IF NOT EXISTS idx_device_alerts_create_user ON device_alerts(create_user_id);

-- 系统用户表添加租户ID字段
ALTER TABLE sys_users ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(64);
CREATE INDEX IF NOT EXISTS idx_sys_users_tenant ON sys_users(tenant_id);
