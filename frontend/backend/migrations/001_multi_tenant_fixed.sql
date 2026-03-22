-- ============================================================
-- MDM 系统多租户数据库迁移脚本 (修复版)
-- 版本: 001
-- 描述: 创建 tenants 表，为所有存在的业务表添加 tenant_id 字段
-- 执行方式: psql -U mdm_user -d mdm_db -f 001_multi_tenant_fixed.sql
-- ============================================================

-- ============================================================
-- 1. 创建 tenants 租户表
-- ============================================================
BEGIN;
CREATE TABLE IF NOT EXISTS tenants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_code     VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(200) NOT NULL,
    contact_name    VARCHAR(100),
    contact_phone   VARCHAR(20),
    contact_email   VARCHAR(100),
    plan            VARCHAR(50) DEFAULT 'free',
    status          VARCHAR(20) DEFAULT 'pending',
    logo_url        VARCHAR(500),
    domain          VARCHAR(200),
    expires_at      TIMESTAMP,
    settings        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
COMMENT ON TABLE tenants IS '租户主表';
COMMENT ON COLUMN tenants.tenant_code IS '租户唯一标识代码';
COMMENT ON COLUMN tenants.plan IS '套餐: free/basic/professional/enterprise';
COMMENT ON COLUMN tenants.status IS '状态: pending/active/suspended/expired';
COMMENT ON COLUMN tenants.settings IS '租户自定义配置(JSON)';
COMMIT;

-- ============================================================
-- 2. sys_users 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE sys_users ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE sys_users ADD CONSTRAINT fk_sys_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_sys_users_tenant_id ON sys_users(tenant_id);
COMMENT ON COLUMN sys_users.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 3. devices 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE devices ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE devices ADD CONSTRAINT fk_devices_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_devices_tenant_id ON devices(tenant_id);
COMMENT ON COLUMN devices.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 4. members 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE members ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE members ADD CONSTRAINT fk_members_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_members_tenant_id ON members(tenant_id);
COMMENT ON COLUMN members.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 5. ota_packages 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE ota_packages ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE ota_packages ADD CONSTRAINT fk_ota_packages_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_ota_packages_tenant_id ON ota_packages(tenant_id);
COMMENT ON COLUMN ota_packages.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 6. ota_deployments 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE ota_deployments ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE ota_deployments ADD CONSTRAINT fk_ota_deployments_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_ota_deployments_tenant_id ON ota_deployments(tenant_id);
COMMENT ON COLUMN ota_deployments.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 7. ota_progresses 表添加 tenant_id (注意: 表名是 ota_progresses)
-- ============================================================
BEGIN;
ALTER TABLE ota_progresses ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE ota_progresses ADD CONSTRAINT fk_ota_progresses_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_ota_progresses_tenant_id ON ota_progresses(tenant_id);
COMMENT ON COLUMN ota_progresses.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 8. device_alerts 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE device_alerts ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE device_alerts ADD CONSTRAINT fk_device_alerts_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_device_alerts_tenant_id ON device_alerts(tenant_id);
COMMENT ON COLUMN device_alerts.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 9. device_alert_rules 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE device_alert_rules ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE device_alert_rules ADD CONSTRAINT fk_device_alert_rules_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_device_alert_rules_tenant_id ON device_alert_rules(tenant_id);
COMMENT ON COLUMN device_alert_rules.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 10. geofence_alerts 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE geofence_alerts ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE geofence_alerts ADD CONSTRAINT fk_geofence_alerts_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_geofence_alerts_tenant_id ON geofence_alerts(tenant_id);
COMMENT ON COLUMN geofence_alerts.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 11. geofence_rules 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE geofence_rules ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE geofence_rules ADD CONSTRAINT fk_geofence_rules_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_geofence_rules_tenant_id ON geofence_rules(tenant_id);
COMMENT ON COLUMN geofence_rules.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 12. alert_notifications 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE alert_notifications ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE alert_notifications ADD CONSTRAINT fk_alert_notifications_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_alert_notifications_tenant_id ON alert_notifications(tenant_id);
COMMENT ON COLUMN alert_notifications.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 13. notifications 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE notifications ADD CONSTRAINT fk_notifications_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id ON notifications(tenant_id);
COMMENT ON COLUMN notifications.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 14. notification_templates 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE notification_templates ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE notification_templates ADD CONSTRAINT fk_notification_templates_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_notification_templates_tenant_id ON notification_templates(tenant_id);
COMMENT ON COLUMN notification_templates.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 15. announcements 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE announcements ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE announcements ADD CONSTRAINT fk_announcements_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_announcements_tenant_id ON announcements(tenant_id);
COMMENT ON COLUMN announcements.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 16. policies 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE policies ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE policies ADD CONSTRAINT fk_policies_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_policies_tenant_id ON policies(tenant_id);
COMMENT ON COLUMN policies.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 17. policy_configs 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE policy_configs ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE policy_configs ADD CONSTRAINT fk_policy_configs_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_policy_configs_tenant_id ON policy_configs(tenant_id);
COMMENT ON COLUMN policy_configs.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 18. policy_bindings 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE policy_bindings ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE policy_bindings ADD CONSTRAINT fk_policy_bindings_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_policy_bindings_tenant_id ON policy_bindings(tenant_id);
COMMENT ON COLUMN policy_bindings.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 19. compliance_policies 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE compliance_policies ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE compliance_policies ADD CONSTRAINT fk_compliance_policies_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_compliance_policies_tenant_id ON compliance_policies(tenant_id);
COMMENT ON COLUMN compliance_policies.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 20. compliance_violations 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE compliance_violations ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE compliance_violations ADD CONSTRAINT fk_compliance_violations_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_compliance_violations_tenant_id ON compliance_violations(tenant_id);
COMMENT ON COLUMN compliance_violations.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 21. apps 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE apps ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE apps ADD CONSTRAINT fk_apps_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_apps_tenant_id ON apps(tenant_id);
COMMENT ON COLUMN apps.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 22. app_versions 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE app_versions ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE app_versions ADD CONSTRAINT fk_app_versions_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_app_versions_tenant_id ON app_versions(tenant_id);
COMMENT ON COLUMN app_versions.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 23. app_distributions 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE app_distributions ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE app_distributions ADD CONSTRAINT fk_app_distributions_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_app_distributions_tenant_id ON app_distributions(tenant_id);
COMMENT ON COLUMN app_distributions.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 24. app_install_records 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE app_install_records ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE app_install_records ADD CONSTRAINT fk_app_install_records_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_app_install_records_tenant_id ON app_install_records(tenant_id);
COMMENT ON COLUMN app_install_records.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 25. app_licenses 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE app_licenses ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE app_licenses ADD CONSTRAINT fk_app_licenses_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_app_licenses_tenant_id ON app_licenses(tenant_id);
COMMENT ON COLUMN app_licenses.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 26. coupons 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE coupons ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE coupons ADD CONSTRAINT fk_coupons_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_coupons_tenant_id ON coupons(tenant_id);
COMMENT ON COLUMN coupons.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 27. coupon_grants 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE coupon_grants ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE coupon_grants ADD CONSTRAINT fk_coupon_grants_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_coupon_grants_tenant_id ON coupon_grants(tenant_id);
COMMENT ON COLUMN coupon_grants.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 28. promotions 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE promotions ADD CONSTRAINT fk_promotions_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_promotions_tenant_id ON promotions(tenant_id);
COMMENT ON COLUMN promotions.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 29. stores 表添加 tenant_id
-- ============================================================
BEGIN;
ALTER TABLE stores ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE stores ADD CONSTRAINT fk_stores_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_stores_tenant_id ON stores(tenant_id);
COMMENT ON COLUMN stores.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 30. command_histories 表添加 tenant_id (注意: 表名是 command_histories)
-- ============================================================
BEGIN;
ALTER TABLE command_histories ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE command_histories ADD CONSTRAINT fk_command_histories_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX IF NOT EXISTS idx_command_histories_tenant_id ON command_histories(tenant_id);
COMMENT ON COLUMN command_histories.tenant_id IS '所属租户ID';
COMMIT;

-- ============================================================
-- 31. 创建默认租户
-- ============================================================
BEGIN;
INSERT INTO tenants (id, tenant_code, name, plan, status, settings, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'default',
    '默认租户',
    'free',
    'active',
    '{"created_by": "migration", "note": "系统初始化默认租户，所有存量数据归属此租户"}'::jsonb,
    NOW(),
    NOW()
)
ON CONFLICT (tenant_code) DO NOTHING;
COMMIT;

-- ============================================================
-- 32. 将所有现有业务数据的 tenant_id 更新为默认租户
-- ============================================================
BEGIN;
UPDATE sys_users             SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE devices               SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE members               SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE ota_packages          SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE ota_deployments       SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE ota_progresses        SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE device_alerts         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE device_alert_rules    SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE geofence_alerts       SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE geofence_rules        SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE alert_notifications   SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE notifications          SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE notification_templates SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE announcements         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE policies               SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE policy_configs         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE policy_bindings        SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE compliance_policies  SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE compliance_violations SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE apps                   SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_versions           SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_distributions      SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_install_records   SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_licenses           SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE coupons                SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE coupon_grants          SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE promotions              SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE stores                 SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE command_histories      SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
COMMIT;

-- ============================================================
-- 33. 记录迁移日志
-- ============================================================
BEGIN;
INSERT INTO sys_operation_logs (user_id, username, module, operation, method, path, params, result, status, created_at)
VALUES (
    0,
    'system',
    'database_migration',
    '001_multi_tenant_fixed',
    'SQL',
    'migrations/001_multi_tenant_fixed.sql',
    '{"tenants_created": true, "default_tenant_id": "00000000-0000-0000-0000-000000000001", "default_tenant_code": "default"}',
    'success',
    1,
    NOW()
);
COMMIT;

-- ============================================================
-- 34. 验证迁移结果
-- ============================================================
DO $$
DECLARE
    v_tenant_count      INTEGER;
    v_tables_with_tenant INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_tenant_count FROM tenants;
    RAISE NOTICE '[迁移验证] tenants 表记录数: %', v_tenant_count;

    SELECT COUNT(*) INTO v_tables_with_tenant
    FROM information_schema.columns
    WHERE table_schema = 'public'
      AND column_name = 'tenant_id';

    RAISE NOTICE '[迁移验证] 已添加 tenant_id 字段的表数量: %', v_tables_with_tenant;

    IF v_tenant_count = 0 THEN
        RAISE EXCEPTION '迁移失败: tenants 表没有记录';
    END IF;

    RAISE NOTICE '[迁移成功] 多租户字段迁移完成，共 % 个表', v_tables_with_tenant;
END;
$$;
