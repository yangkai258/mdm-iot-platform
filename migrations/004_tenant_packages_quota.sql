-- ============================================================
-- MDM 系统多租户套餐配额迁移脚本
-- 版本: 004
-- 描述: 创建套餐表 (packages) 和套餐配额表 (package_quotas)
-- 执行方式: psql -U mdm_user -d mdm_db -f 004_tenant_packages_quota.sql
-- ============================================================

BEGIN;

-- ============================================================
-- 1. packages 套餐表
-- ============================================================
CREATE TABLE IF NOT EXISTS packages (
    id              BIGSERIAL PRIMARY KEY,
    package_code    VARCHAR(50) UNIQUE NOT NULL,
    package_name    VARCHAR(100) NOT NULL,
    plan_type       VARCHAR(20) DEFAULT 'free'           NOT NULL, -- free, basic, professional, enterprise
    description     VARCHAR(500),
    price_monthly   DECIMAL(10,2) DEFAULT 0,
    price_yearly    DECIMAL(10,2) DEFAULT 0,
    is_active       BOOLEAN DEFAULT true,
    is_default      BOOLEAN DEFAULT false,                -- 是否为默认套餐
    sort_order      INTEGER DEFAULT 0,
    features        JSONB DEFAULT '{}',                   -- 功能特性列表
    quota_config    JSONB DEFAULT '{}',                    -- 配额配置
    settings        JSONB DEFAULT '{}',                    -- 附加设置
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    deleted_at      TIMESTAMP
);
COMMENT ON TABLE packages IS '租户套餐表';
COMMENT ON COLUMN packages.package_code IS '套餐代码';
COMMENT ON COLUMN packages.plan_type IS '套餐类型: free/basic/professional/enterprise';
COMMENT ON COLUMN packages.features IS '功能特性(JSON): {"ota": true, "analytics": false}';
COMMENT ON COLUMN packages.quota_config IS '配额配置(JSON): {"devices": 10, "users": 5}';

-- ============================================================
-- 2. package_quotas 套餐配额表
-- ============================================================
CREATE TABLE IF NOT EXISTS package_quotas (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       UUID NOT NULL,
    package_id      BIGINT NOT NULL,
    quota_type      VARCHAR(50) NOT NULL,                 -- user, device, store, dept, ota_deployment, app, notification, alert
    quota_limit     INTEGER NOT NULL DEFAULT 0,             -- 配额上限，0表示无限制
    quota_used      INTEGER NOT NULL DEFAULT 0,             -- 当前使用量
    quota_warn_at   INTEGER,                               -- 警告阈值（百分比）
    updated_at      TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_package_quotas_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_package_quotas_package FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE RESTRICT,
    CONSTRAINT uq_package_quotas_tenant_type UNIQUE (tenant_id, quota_type)
);
COMMENT ON TABLE package_quotas IS '租户套餐配额表';
COMMENT ON COLUMN package_quotas.quota_type IS '配额类型: user/device/store/dept/ota_deployment/app/notification/alert';
COMMENT ON COLUMN package_quotas.quota_limit IS '配额上限，0表示无限制(-1)';

-- 索引
CREATE INDEX IF NOT EXISTS idx_package_quotas_tenant_id ON package_quotas(tenant_id);
CREATE INDEX IF NOT EXISTS idx_package_quotas_package_id ON package_quotas(package_id);

COMMIT;

BEGIN;

-- ============================================================
-- 3. 插入默认套餐数据
-- ============================================================

-- Free 套餐
INSERT INTO packages (package_code, package_name, plan_type, description,
    price_monthly, price_yearly, is_active, is_default, sort_order,
    features, quota_config, settings)
VALUES (
    'free',
    '免费版',
    'free',
    '适合个人或小规模团队试用，包含基础设备管理功能',
    0, 0,
    true, true, 0,
    '{"device_management": true, "ota_upgrade": true, "basic_alerts": true, "messaging": false, "analytics": false, "custom_branding": false, "api_access": false, "sla": null}'::jsonb,
    '{"devices": 5, "users": 2, "stores": 1, "departments": 1, "ota_deployments": 1, "apps": 1, "notifications": 100, "alerts": 50}'::jsonb,
    '{"max_device_per_user": 3, "data_retention_days": 7}'::jsonb
) ON CONFLICT (package_code) DO NOTHING;

-- Basic 套餐
INSERT INTO packages (package_code, package_name, plan_type, description,
    price_monthly, price_yearly, is_active, is_default, sort_order,
    features, quota_config, settings)
VALUES (
    'basic',
    '基础版',
    'basic',
    '适合中小企业，支持更多设备和用户',
    99, 990,
    true, false, 1,
    '{"device_management": true, "ota_upgrade": true, "basic_alerts": true, "messaging": true, "analytics": true, "custom_branding": false, "api_access": false, "sla": "99.5"}'::jsonb,
    '{"devices": 50, "users": 10, "stores": 5, "departments": 5, "ota_deployments": 5, "apps": 5, "notifications": 1000, "alerts": 500}'::jsonb,
    '{"max_device_per_user": 10, "data_retention_days": 30}'::jsonb
) ON CONFLICT (package_code) DO NOTHING;

-- Professional 套餐
INSERT INTO packages (package_code, package_name, plan_type, description,
    price_monthly, price_yearly, is_active, is_default, sort_order,
    features, quota_config, settings)
VALUES (
    'professional',
    '专业版',
    'professional',
    '适合成长型企业，全功能支持',
    299, 2990,
    true, false, 2,
    '{"device_management": true, "ota_upgrade": true, "basic_alerts": true, "messaging": true, "analytics": true, "custom_branding": true, "api_access": true, "sla": "99.9"}'::jsonb,
    '{"devices": 200, "users": 50, "stores": 20, "departments": 20, "ota_deployments": 20, "apps": 20, "notifications": 5000, "alerts": 2000}'::jsonb,
    '{"max_device_per_user": 20, "data_retention_days": 90}'::jsonb
) ON CONFLICT (package_code) DO NOTHING;

-- Enterprise 套餐
INSERT INTO packages (package_code, package_name, plan_type, description,
    price_monthly, price_yearly, is_active, is_default, sort_order,
    features, quota_config, settings)
VALUES (
    'enterprise',
    '企业版',
    'enterprise',
    '适合大型企业，无限制配额，专属支持',
    999, 9990,
    true, false, 3,
    '{"device_management": true, "ota_upgrade": true, "basic_alerts": true, "messaging": true, "analytics": true, "custom_branding": true, "api_access": true, "sla": "99.99", "dedicated_support": true, "custom_integrations": true}'::jsonb,
    '{"devices": -1, "users": -1, "stores": -1, "departments": -1, "ota_deployments": -1, "apps": -1, "notifications": -1, "alerts": -1}'::jsonb,
    '{"max_device_per_user": -1, "data_retention_days": 365}'::jsonb
) ON CONFLICT (package_code) DO NOTHING;

COMMIT;

BEGIN;

-- ============================================================
-- 4. 为现有默认租户初始化套餐配额记录
-- ============================================================
DO $$
DECLARE
    v_default_tenant_id UUID := '00000000-0000-0000-0000-000000000001';
    v_free_package_id   BIGINT;
BEGIN
    -- 获取 free 套餐 ID
    SELECT id INTO v_free_package_id FROM packages WHERE package_code = 'free' LIMIT 1;

    -- 插入各配额类型记录
    INSERT INTO package_quotas (tenant_id, package_id, quota_type, quota_limit, quota_used, quota_warn_at)
    VALUES
        (v_default_tenant_id, v_free_package_id, 'user',              2,   0, 80),
        (v_default_tenant_id, v_free_package_id, 'device',           5,   0, 80),
        (v_default_tenant_id, v_free_package_id, 'store',            1,   0, 80),
        (v_default_tenant_id, v_free_package_id, 'dept',             1,   0, 80),
        (v_default_tenant_id, v_free_package_id, 'ota_deployment',   1,   0, 80),
        (v_default_tenant_id, v_free_package_id, 'app',              1,   0, 80),
        (v_default_tenant_id, v_free_package_id, 'notification',   100,   0, 80),
        (v_default_tenant_id, v_free_package_id, 'alert',            50,   0, 80)
    ON CONFLICT (tenant_id, quota_type) DO NOTHING;

    RAISE NOTICE '[迁移004] packages 和 package_quotas 表创建完成，默认租户配额初始化完成';
END;
$$;

COMMIT;

-- ============================================================
-- 5. 记录迁移日志
-- ============================================================
BEGIN;
INSERT INTO sys_operation_logs (user_id, username, module, operation, method, path, params, result, status, created_at)
VALUES (
    0,
    'system',
    'database_migration',
    '004_tenant_packages_quota',
    'SQL',
    'migrations/004_tenant_packages_quota.sql',
    '{"tables_created": ["packages", "package_quotas"], "default_packages": ["free", "basic", "professional", "enterprise"]}',
    'success',
    1,
    NOW()
);
COMMIT;

-- ============================================================
-- 6. 验证迁移结果
-- ============================================================
DO $$
DECLARE
    v_package_count  INTEGER;
    v_quota_count    INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_package_count FROM packages WHERE deleted_at IS NULL;
    RAISE NOTICE '[迁移验证] packages 表记录数: %', v_package_count;

    SELECT COUNT(*) INTO v_quota_count FROM package_quotas
    WHERE tenant_id = '00000000-0000-0000-0000-000000000001';
    RAISE NOTICE '[迁移验证] 默认租户 package_quotas 记录数: %', v_quota_count;

    IF v_package_count < 1 THEN
        RAISE EXCEPTION '迁移失败: packages 表没有记录';
    END IF;

    RAISE NOTICE '[迁移成功] 套餐配额表创建完成，共 % 个套餐，% 条配额记录',
        v_package_count, v_quota_count;
END;
$$;
