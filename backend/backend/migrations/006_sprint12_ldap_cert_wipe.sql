-- Sprint 12: LDAP/AD集成、证书管理、远程设备锁定/擦除、数据权限
-- 执行时间: 2026-03-22

-- ============================================
-- 1. LDAP 配置相关表
-- ============================================

-- LDAP 配置表
CREATE TABLE IF NOT EXISTS ldap_configs (
    id              BIGSERIAL PRIMARY KEY,
    config_name     VARCHAR(100),
    host            VARCHAR(255) NOT NULL,
    port            INT DEFAULT 389,
    base_dn         VARCHAR(255) NOT NULL,
    bind_dn         VARCHAR(255),
    bind_password   VARCHAR(255),                  -- AES 加密存储
    use_ssl         BOOLEAN DEFAULT FALSE,
    use_tls         BOOLEAN DEFAULT FALSE,
    user_filter     VARCHAR(500) DEFAULT '(objectClass=user)',
    group_filter    VARCHAR(500) DEFAULT '(objectClass=group)',
    sync_interval   INT DEFAULT 3600,              -- 秒
    is_enabled      BOOLEAN DEFAULT FALSE,
    last_sync_at    TIMESTAMP,
    status          VARCHAR(20) DEFAULT 'inactive', -- active/inactive/error
    tenant_id       VARCHAR(50),
    description     TEXT,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- LDAP 用户映射表
CREATE TABLE IF NOT EXISTS ldap_user_mappings (
    id              BIGSERIAL PRIMARY KEY,
    ldap_dn         VARCHAR(255) NOT NULL UNIQUE,
    local_user_id   BIGINT,
    username        VARCHAR(100),
    email           VARCHAR(255),
    display_name    VARCHAR(100),
    ldap_groups     TEXT,                           -- JSON 数组
    sync_status     VARCHAR(20) DEFAULT 'synced',  -- synced/pending/removed
    last_synced_at  TIMESTAMP,
    tenant_id       VARCHAR(50),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ldap_user_mappings_tenant ON ldap_user_mappings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ldap_user_mappings_username ON ldap_user_mappings(username);

-- LDAP 分组-角色映射表
CREATE TABLE IF NOT EXISTS ldap_group_role_mappings (
    id              BIGSERIAL PRIMARY KEY,
    ldap_group_dn   VARCHAR(255) NOT NULL,
    ldap_group_name VARCHAR(100),
    role_id         BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    tenant_id       VARCHAR(50),
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(ldap_group_dn, tenant_id)
);

CREATE INDEX IF NOT EXISTS idx_ldap_group_role_mappings_tenant ON ldap_group_role_mappings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ldap_group_role_mappings_role ON ldap_group_role_mappings(role_id);

-- ============================================
-- 2. 证书管理相关表
-- ============================================

-- 证书表
CREATE TABLE IF NOT EXISTS certificates (
    id              BIGSERIAL PRIMARY KEY,
    cert_id         VARCHAR(64) NOT NULL UNIQUE,
    cert_name       VARCHAR(128),
    cert_type       VARCHAR(32),                    -- device/client/server/ca
    serial_number    VARCHAR(128),
    subject         VARCHAR(256),
    issuer          VARCHAR(256),
    thumbprint      VARCHAR(64) UNIQUE,
    not_before      TIMESTAMP NOT NULL,
    not_after       TIMESTAMP NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',   -- active/expired/revoked/pending
    cert_file       VARCHAR(512),                  -- 证书文件路径 (公钥)
    key_file        VARCHAR(512),                  -- 私钥文件路径 (敏感，仅文件权限保护)
    tenant_id       VARCHAR(50),
    description     TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_certificates_tenant ON certificates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_certificates_status ON certificates(status);
CREATE INDEX IF NOT EXISTS idx_certificates_not_after ON certificates(not_after);
CREATE INDEX IF NOT EXISTS idx_certificates_serial ON certificates(serial_number);
CREATE INDEX IF NOT EXISTS idx_certificates_thumbprint ON certificates(thumbprint);

-- ============================================
-- 3. 设备安全相关表
-- ============================================

-- 设备擦除历史表
CREATE TABLE IF NOT EXISTS wipe_history (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    operator_id     BIGINT NOT NULL,
    operator_name   VARCHAR(100),
    wipe_type       VARCHAR(32),                    -- full/selective
    status          VARCHAR(20) DEFAULT 'pending', -- pending/executing/completed/failed
    confirm_token   VARCHAR(64),                    -- 二次确认token
    confirmed_at    TIMESTAMP,
    executed_at     TIMESTAMP,
    completed_at    TIMESTAMP,
    result          TEXT,                          -- 操作结果/错误信息
    reason          TEXT,                          -- 擦除原因
    tenant_id       VARCHAR(50),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wipe_history_device ON wipe_history(device_id);
CREATE INDEX IF NOT EXISTS idx_wipe_history_tenant ON wipe_history(tenant_id);
CREATE INDEX IF NOT EXISTS idx_wipe_history_operator ON wipe_history(operator_id);
CREATE INDEX IF NOT EXISTS idx_wipe_history_status ON wipe_history(status);
CREATE INDEX IF NOT EXISTS idx_wipe_history_created ON wipe_history(created_at);

-- ============================================
-- 4. 数据权限相关表
-- ============================================

-- 数据权限规则表
CREATE TABLE IF NOT EXISTS data_permission_rules (
    id              BIGSERIAL PRIMARY KEY,
    rule_name       VARCHAR(100) NOT NULL,
    resource_type   VARCHAR(50) NOT NULL,          -- device/pet/member/org
    rule_type       VARCHAR(20) NOT NULL,          -- row/column
    resource_ids    TEXT[],                        -- 资源ID列表，NULL表示全部
    permission_expr JSONB NOT NULL,                -- 权限表达式
    priority        INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    description     TEXT,
    tenant_id       VARCHAR(50),
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_data_permission_rules_tenant ON data_permission_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_data_permission_rules_resource ON data_permission_rules(resource_type);
CREATE INDEX IF NOT EXISTS idx_data_permission_rules_active ON data_permission_rules(is_active);

-- 用户数据权限表（行级/列级权限）
CREATE TABLE IF NOT EXISTS user_data_permissions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    role_id         BIGINT DEFAULT 0,               -- 可选，角色级别配置
    resource_type   VARCHAR(50) NOT NULL,
    rule_type       VARCHAR(20) NOT NULL,           -- row/column
    column_fields   TEXT[],                         -- 列级权限：可访问的字段列表
    data_scope      JSONB,                          -- 行级权限：数据范围表达式
    filter_expr     TEXT,                           -- 自定义过滤表达式
    priority        INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    tenant_id       VARCHAR(50),
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_data_permissions_user ON user_data_permissions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_data_permissions_role ON user_data_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_user_data_permissions_tenant ON user_data_permissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_data_permissions_resource ON user_data_permissions(resource_type);

-- ============================================
-- 5. 更新现有表（可选字段）
-- ============================================

-- 为设备表添加锁定状态字段（可选）
-- ALTER TABLE devices ADD COLUMN IF NOT EXISTS is_locked BOOLEAN DEFAULT FALSE;
-- ALTER TABLE devices ADD COLUMN locked_at TIMESTAMP;
-- ALTER TABLE devices ADD COLUMN locked_by VARCHAR(100);

-- 添加注释
COMMENT ON TABLE certificates IS '设备证书表 - 存储设备/客户端/服务端证书';
COMMENT ON TABLE wipe_history IS '设备擦除历史表 - 记录所有远程擦除操作（审计用）';
COMMENT ON TABLE ldap_configs IS 'LDAP/AD 配置表';
COMMENT ON TABLE data_permission_rules IS '数据权限规则表 - 行级/列级权限配置';
COMMENT ON COLUMN certificates.key_file IS '私钥文件路径，权限应为600，不通过API返回';
COMMENT ON COLUMN ldap_configs.bind_password IS '管理员密码，AES加密存储';
