-- ============================================================
-- MDM 系统多租户数据库迁移脚本
-- 版本: 001
-- 描述: 创建 tenants 表，为所有业务表添加 tenant_id 字段
-- 执行方式: psql -U postgres -d mdm_db -f 001_multi_tenant.sql
-- ============================================================

BEGIN;

-- ============================================================
-- 1. 创建 tenants 租户表
-- ============================================================
CREATE TABLE IF NOT EXISTS tenants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_code     VARCHAR(50) UNIQUE NOT NULL,          -- 租户标识，如 "company_a"
    name            VARCHAR(200) NOT NULL,                 -- 公司名称
    contact_name    VARCHAR(100),                          -- 联系人
    contact_phone   VARCHAR(20),                           -- 手机
    contact_email   VARCHAR(100),                          -- 邮箱
    plan            VARCHAR(50) DEFAULT 'free',            -- free/basic/professional/enterprise
    status          VARCHAR(20) DEFAULT 'pending',        -- pending/active/suspended/expired
    logo_url        VARCHAR(500),
    domain          VARCHAR(200),                          -- 自定义域名
    expires_at      TIMESTAMP,
    settings        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE tenants IS '租户主表';
COMMENT ON COLUMN tenants.tenant_code IS '租户唯一标识代码';
COMMENT ON COLUMN tenants.plan IS '套餐: free(免费)/basic(基础)/professional(专业)/enterprise(企业)';
COMMENT ON COLUMN tenants.status IS '状态: pending(待激活)/active(正常)/suspended(停用)/expired(过期)';
COMMENT ON COLUMN tenants.settings IS '租户自定义配置(JSON)';

-- ============================================================
-- 2. 为 sys_users 表添加 tenant_id
-- ============================================================
ALTER TABLE sys_users ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN sys_users.tenant_id IS '所属租户ID';

-- 为 sys_users 添加外键约束（软参照，不阻断业务）
ALTER TABLE sys_users
    ADD CONSTRAINT fk_sys_users_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_sys_users_tenant_id ON sys_users(tenant_id);

-- ============================================================
-- 3. 为 devices 表添加 tenant_id
-- ============================================================
ALTER TABLE devices ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN devices.tenant_id IS '所属租户ID';

ALTER TABLE devices
    ADD CONSTRAINT fk_devices_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_devices_tenant_id ON devices(tenant_id);

-- ============================================================
-- 4. 为 members 表添加 tenant_id
-- ============================================================
ALTER TABLE members ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN members.tenant_id IS '所属租户ID';

ALTER TABLE members
    ADD CONSTRAINT fk_members_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_members_tenant_id ON members(tenant_id);

-- ============================================================
-- 5. 为 ota_packages 表添加 tenant_id
-- ============================================================
ALTER TABLE ota_packages ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN ota_packages.tenant_id IS '所属租户ID';

ALTER TABLE ota_packages
    ADD CONSTRAINT fk_ota_packages_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_ota_packages_tenant_id ON ota_packages(tenant_id);

-- ============================================================
-- 6. 为 ota_deployments 表添加 tenant_id
-- ============================================================
ALTER TABLE ota_deployments ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN ota_deployments.tenant_id IS '所属租户ID';

ALTER TABLE ota_deployments
    ADD CONSTRAINT fk_ota_deployments_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_ota_deployments_tenant_id ON ota_deployments(tenant_id);

-- ============================================================
-- 7. 为 ota_progress 表添加 tenant_id
-- ============================================================
ALTER TABLE ota_progress ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN ota_progress.tenant_id IS '所属租户ID';

ALTER TABLE ota_progress
    ADD CONSTRAINT fk_ota_progress_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_ota_progress_tenant_id ON ota_progress(tenant_id);

-- ============================================================
-- 8. 为 device_alerts 表添加 tenant_id
-- ============================================================
ALTER TABLE device_alerts ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN device_alerts.tenant_id IS '所属租户ID';

ALTER TABLE device_alerts
    ADD CONSTRAINT fk_device_alerts_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_device_alerts_tenant_id ON device_alerts(tenant_id);

-- ============================================================
-- 9. 为 device_alert_rules 表添加 tenant_id
-- ============================================================
ALTER TABLE device_alert_rules ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN device_alert_rules.tenant_id IS '所属租户ID';

ALTER TABLE device_alert_rules
    ADD CONSTRAINT fk_device_alert_rules_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_device_alert_rules_tenant_id ON device_alert_rules(tenant_id);

-- ============================================================
-- 10. 为 geofence_alerts 表添加 tenant_id
-- ============================================================
ALTER TABLE geofence_alerts ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN geofence_alerts.tenant_id IS '所属租户ID';

ALTER TABLE geofence_alerts
    ADD CONSTRAINT fk_geofence_alerts_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_geofence_alerts_tenant_id ON geofence_alerts(tenant_id);

-- ============================================================
-- 11. 为 geofence_rules 表添加 tenant_id
-- ============================================================
ALTER TABLE geofence_rules ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN geofence_rules.tenant_id IS '所属租户ID';

ALTER TABLE geofence_rules
    ADD CONSTRAINT fk_geofence_rules_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_geofence_rules_tenant_id ON geofence_rules(tenant_id);

-- ============================================================
-- 12. 为 alert_notifications 表添加 tenant_id
-- ============================================================
ALTER TABLE alert_notifications ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN alert_notifications.tenant_id IS '所属租户ID';

ALTER TABLE alert_notifications
    ADD CONSTRAINT fk_alert_notifications_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_alert_notifications_tenant_id ON alert_notifications(tenant_id);

-- ============================================================
-- 13. 为 notifications 表添加 tenant_id
-- ============================================================
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN notifications.tenant_id IS '所属租户ID';

ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id ON notifications(tenant_id);

-- ============================================================
-- 14. 为 notification_templates 表添加 tenant_id
-- ============================================================
ALTER TABLE notification_templates ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN notification_templates.tenant_id IS '所属租户ID';

ALTER TABLE notification_templates
    ADD CONSTRAINT fk_notification_templates_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_notification_templates_tenant_id ON notification_templates(tenant_id);

-- ============================================================
-- 15. 为 announcements 表添加 tenant_id
-- ============================================================
ALTER TABLE announcements ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN announcements.tenant_id IS '所属租户ID';

ALTER TABLE announcements
    ADD CONSTRAINT fk_announcements_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_announcements_tenant_id ON announcements(tenant_id);

-- ============================================================
-- 16. 为 policies 表添加 tenant_id
-- ============================================================
ALTER TABLE policies ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN policies.tenant_id IS '所属租户ID';

ALTER TABLE policies
    ADD CONSTRAINT fk_policies_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_policies_tenant_id ON policies(tenant_id);

-- ============================================================
-- 17. 为 policy_configs 表添加 tenant_id
-- ============================================================
ALTER TABLE policy_configs ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN policy_configs.tenant_id IS '所属租户ID';

ALTER TABLE policy_configs
    ADD CONSTRAINT fk_policy_configs_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_policy_configs_tenant_id ON policy_configs(tenant_id);

-- ============================================================
-- 18. 为 policy_bindings 表添加 tenant_id
-- ============================================================
ALTER TABLE policy_bindings ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN policy_bindings.tenant_id IS '所属租户ID';

ALTER TABLE policy_bindings
    ADD CONSTRAINT fk_policy_bindings_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_policy_bindings_tenant_id ON policy_bindings(tenant_id);

-- ============================================================
-- 19. 为 compliance_policies 表添加 tenant_id
-- ============================================================
ALTER TABLE compliance_policies ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN compliance_policies.tenant_id IS '所属租户ID';

ALTER TABLE compliance_policies
    ADD CONSTRAINT fk_compliance_policies_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_compliance_policies_tenant_id ON compliance_policies(tenant_id);

-- ============================================================
-- 20. 为 compliance_violations 表添加 tenant_id
-- ============================================================
ALTER TABLE compliance_violations ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN compliance_violations.tenant_id IS '所属租户ID';

ALTER TABLE compliance_violations
    ADD CONSTRAINT fk_compliance_violations_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_compliance_violations_tenant_id ON compliance_violations(tenant_id);

-- ============================================================
-- 21. 为 apps 表添加 tenant_id
-- ============================================================
ALTER TABLE apps ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN apps.tenant_id IS '所属租户ID';

ALTER TABLE apps
    ADD CONSTRAINT fk_apps_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_apps_tenant_id ON apps(tenant_id);

-- ============================================================
-- 22. 为 app_versions 表添加 tenant_id
-- ============================================================
ALTER TABLE app_versions ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN app_versions.tenant_id IS '所属租户ID';

ALTER TABLE app_versions
    ADD CONSTRAINT fk_app_versions_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_app_versions_tenant_id ON app_versions(tenant_id);

-- ============================================================
-- 23. 为 app_distributions 表添加 tenant_id
-- ============================================================
ALTER TABLE app_distributions ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN app_distributions.tenant_id IS '所属租户ID';

ALTER TABLE app_distributions
    ADD CONSTRAINT fk_app_distributions_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_app_distributions_tenant_id ON app_distributions(tenant_id);

-- ============================================================
-- 24. 为 app_install_records 表添加 tenant_id
-- ============================================================
ALTER TABLE app_install_records ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN app_install_records.tenant_id IS '所属租户ID';

ALTER TABLE app_install_records
    ADD CONSTRAINT fk_app_install_records_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_app_install_records_tenant_id ON app_install_records(tenant_id);

-- ============================================================
-- 25. 为 app_licenses 表添加 tenant_id
-- ============================================================
ALTER TABLE app_licenses ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN app_licenses.tenant_id IS '所属租户ID';

ALTER TABLE app_licenses
    ADD CONSTRAINT fk_app_licenses_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_app_licenses_tenant_id ON app_licenses(tenant_id);

-- ============================================================
-- 26. 为 coupons 表添加 tenant_id
-- ============================================================
ALTER TABLE coupons ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN coupons.tenant_id IS '所属租户ID';

ALTER TABLE coupons
    ADD CONSTRAINT fk_coupons_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_coupons_tenant_id ON coupons(tenant_id);

-- ============================================================
-- 27. 为 coupon_grants 表添加 tenant_id
-- ============================================================
ALTER TABLE coupon_grants ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN coupon_grants.tenant_id IS '所属租户ID';

ALTER TABLE coupon_grants
    ADD CONSTRAINT fk_coupon_grants_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_coupon_grants_tenant_id ON coupon_grants(tenant_id);

-- ============================================================
-- 28. 为 promotions 表添加 tenant_id
-- ============================================================
ALTER TABLE promotions ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN promotions.tenant_id IS '所属租户ID';

ALTER TABLE promotions
    ADD CONSTRAINT fk_promotions_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_promotions_tenant_id ON promotions(tenant_id);

-- ============================================================
-- 29. 为 stores 表添加 tenant_id
-- ============================================================
ALTER TABLE stores ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN stores.tenant_id IS '所属租户ID';

ALTER TABLE stores
    ADD CONSTRAINT fk_stores_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_stores_tenant_id ON stores(tenant_id);

-- ============================================================
-- 30. 为 pet_conversations 表添加 tenant_id
-- ============================================================
ALTER TABLE pet_conversations ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN pet_conversations.tenant_id IS '所属租户ID';

ALTER TABLE pet_conversations
    ADD CONSTRAINT fk_pet_conversations_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_pet_conversations_tenant_id ON pet_conversations(tenant_id);

-- ============================================================
-- 31. 为 pet_messages 表添加 tenant_id
-- ============================================================
ALTER TABLE pet_messages ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN pet_messages.tenant_id IS '所属租户ID';

ALTER TABLE pet_messages
    ADD CONSTRAINT fk_pet_messages_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_pet_messages_tenant_id ON pet_messages(tenant_id);

-- ============================================================
-- 32. 为 long_term_memories 表添加 tenant_id
-- ============================================================
ALTER TABLE long_term_memories ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN long_term_memories.tenant_id IS '所属租户ID';

ALTER TABLE long_term_memories
    ADD CONSTRAINT fk_long_term_memories_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_long_term_memories_tenant_id ON long_term_memories(tenant_id);

-- ============================================================
-- 33. 为 pet_behavior_actions 表添加 tenant_id
-- ============================================================
ALTER TABLE pet_behavior_actions ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN pet_behavior_actions.tenant_id IS '所属租户ID';

ALTER TABLE pet_behavior_actions
    ADD CONSTRAINT fk_pet_behavior_actions_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_pet_behavior_actions_tenant_id ON pet_behavior_actions(tenant_id);

-- ============================================================
-- 34. 为 companies 表添加 tenant_id（与 tenants 合一，但保留以兼容现有结构）
-- ============================================================
ALTER TABLE companies ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN companies.tenant_id IS '关联租户ID';

ALTER TABLE companies
    ADD CONSTRAINT fk_companies_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_companies_tenant_id ON companies(tenant_id);

-- ============================================================
-- 35. 为 departments 表添加 tenant_id
-- ============================================================
ALTER TABLE departments ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN departments.tenant_id IS '所属租户ID';

ALTER TABLE departments
    ADD CONSTRAINT fk_departments_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_departments_tenant_id ON departments(tenant_id);

-- ============================================================
-- 36. 为 employees 表添加 tenant_id
-- ============================================================
ALTER TABLE employees ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN employees.tenant_id IS '所属租户ID';

ALTER TABLE employees
    ADD CONSTRAINT fk_employees_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_employees_tenant_id ON employees(tenant_id);

-- ============================================================
-- 37. 为 command_history 表添加 tenant_id
-- ============================================================
ALTER TABLE command_history ADD COLUMN IF NOT EXISTS tenant_id UUID;
COMMENT ON COLUMN command_history.tenant_id IS '所属租户ID';

ALTER TABLE command_history
    ADD CONSTRAINT fk_command_history_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

CREATE INDEX IF NOT EXISTS idx_command_history_tenant_id ON command_history(tenant_id);

-- ============================================================
-- ============================================================
-- 3. 数据迁移：将现有数据归属到默认租户
-- ============================================================

-- 3.1 创建默认租户记录
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

-- 3.2 将所有现有业务数据的 tenant_id 更新为默认租户
-- 注意：仅更新 tenant_id 为 NULL 的记录，避免重复迁移
UPDATE sys_users            SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE devices              SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE members              SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE ota_packages         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE ota_deployments      SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE ota_progress         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE device_alerts        SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE device_alert_rules   SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE geofence_alerts      SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE geofence_rules       SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE alert_notifications  SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE notifications         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE notification_templates SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE announcements         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE policies              SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE policy_configs        SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE policy_bindings       SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE compliance_policies  SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE compliance_violations SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE apps                  SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_versions          SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_distributions     SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_install_records   SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE app_licenses          SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE coupons               SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE coupon_grants         SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE promotions             SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE stores                SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE pet_conversations     SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE pet_messages          SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE long_term_memories    SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE pet_behavior_actions SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE companies             SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE departments            SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE employees             SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
UPDATE command_history       SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;

-- ============================================================
-- 4. 记录迁移日志
-- ============================================================
INSERT INTO sys_operation_logs (user_id, username, module, operation, method, path, params, result, status, created_at)
VALUES (
    0,
    'system',
    'database_migration',
    '001_multi_tenant',
    'SQL',
    'migrations/001_multi_tenant.sql',
    '{"tenants_created": true, "default_tenant_id": "00000000-0000-0000-0000-000000000001", "default_tenant_code": "default"}',
    'success',
    1,
    NOW()
);

-- ============================================================
-- 5. 验证迁移结果
-- ============================================================
DO $$
DECLARE
    v_tenant_count      INTEGER;
    v_tables_with_tenant INTEGER;
BEGIN
    -- 检查 tenants 表
    SELECT COUNT(*) INTO v_tenant_count FROM tenants;
    RAISE NOTICE '[迁移验证] tenants 表记录数: %', v_tenant_count;

    -- 检查已添加 tenant_id 字段的表数量
    SELECT COUNT(*) INTO v_tables_with_tenant
    FROM information_schema.columns
    WHERE table_schema = 'public'
      AND column_name = 'tenant_id';

    RAISE NOTICE '[迁移验证] 已添加 tenant_id 字段的表数量: %', v_tables_with_tenant;

    IF v_tenant_count = 0 THEN
        RAISE EXCEPTION '迁移失败: tenants 表没有记录';
    END IF;

    IF v_tables_with_tenant < 30 THEN
        RAISE EXCEPTION '迁移警告: 仅 % 个表添加了 tenant_id 字段，可能存在遗漏', v_tables_with_tenant;
    END IF;

    RAISE NOTICE '[迁移成功] 多租户字段迁移完成，共 % 个表', v_tables_with_tenant;
END;
$$;

COMMIT;
