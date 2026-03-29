-- 功能配置表迁移
-- 用于管理功能模块的分组和排序

-- 功能分组表
CREATE TABLE IF NOT EXISTS feature_groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(100) NOT NULL,
    group_code VARCHAR(50) UNIQUE,
    icon VARCHAR(100),
    color VARCHAR(20),
    sort INT DEFAULT 0,
    description TEXT,
    status INT DEFAULT 1,  -- 1=启用 0=禁用
    tenant_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feature_groups_tenant ON feature_groups(tenant_id);
CREATE INDEX IF NOT EXISTS idx_feature_groups_deleted ON feature_groups(deleted_at);

-- 功能项表
CREATE TABLE IF NOT EXISTS feature_items (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES feature_groups(id) ON DELETE CASCADE,
    feature_name VARCHAR(100) NOT NULL,
    feature_code VARCHAR(50) UNIQUE,
    icon VARCHAR(100),
    route_path VARCHAR(255),
    component VARCHAR(255),
    api_paths TEXT,  -- JSON数组，存储关联的API路径
    permission VARCHAR(100),
    sort INT DEFAULT 0,
    status INT DEFAULT 1,  -- 1=启用 0=禁用
    is_default INT DEFAULT 0,  -- 1=默认选中
    badge VARCHAR(50),  -- 徽章，如"新"、"Beta"
    description TEXT,
    tenant_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feature_items_group ON feature_items(group_id);
CREATE INDEX IF NOT EXISTS idx_feature_items_tenant ON feature_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_feature_items_deleted ON feature_items(deleted_at);

-- 插入默认分组和功能
INSERT INTO feature_groups (group_name, group_code, icon, color, sort, description, tenant_id) VALUES
('设备管理', 'device', 'Cpu', '#165dff', 1, '设备注册、绑定、OTA升级', '00000000-0000-0000-0000-000000000001'),
('会员管理', 'member', 'User', '#00b42a', 2, '会员信息、积分、优惠券', '00000000-0000-0000-0000-000000000001'),
('告警管理', 'alert', 'Bell', '#ff7d00', 3, '设备告警、通知设置', '00000000-0000-0000-0000-000000000001'),
('AI功能', 'ai', 'Robot', '#722ed1', 4, 'AI聊天、情感计算', '00000000-0000-0000-0000-000000000001'),
('数据分析', 'analytics', 'Chart', '#eb0aa4', 5, '数据统计、报表', '00000000-0000-0000-0000-000000000001')
ON CONFLICT (group_code) DO NOTHING;

-- 插入默认功能项
INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default) 
SELECT g.id, '设备列表', 'device_list', 'Device', '/devices', 'DeviceListView', 'device:view', 1, g.tenant_id, 1
FROM feature_groups g WHERE g.group_code = 'device'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '设备注册', 'device_register', 'Plus', '/devices/register', 'DeviceRegisterView', 'device:manage', 2, g.tenant_id, 0
FROM feature_groups g WHERE g.group_code = 'device'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, 'OTA管理', 'ota', 'CloudUpload', '/ota', 'OTAManageView', 'ota:view', 3, g.tenant_id, 0
FROM feature_groups g WHERE g.group_code = 'device'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '会员列表', 'member_list', 'User', '/members', 'MemberListView', 'member:view', 1, g.tenant_id, 1
FROM feature_groups g WHERE g.group_code = 'member'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '会员等级', 'member_levels', 'Star', '/members/levels', 'MemberLevelsView', 'member:manage', 2, g.tenant_id, 0
FROM feature_groups g WHERE g.group_code = 'member'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '优惠券', 'coupons', 'Ticket', '/coupons', 'CouponListView', 'member:manage', 3, g.tenant_id, 0
FROM feature_groups g WHERE g.group_code = 'member'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '告警列表', 'alert_list', 'Bell', '/alerts', 'AlertListView', 'alert:view', 1, g.tenant_id, 1
FROM feature_groups g WHERE g.group_code = 'alert'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '告警规则', 'alert_rules', 'Setting', '/alerts/rules', 'AlertRulesView', 'alert:manage', 2, g.tenant_id, 0
FROM feature_groups g WHERE g.group_code = 'alert'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, 'AI聊天', 'ai_chat', 'Chat', '/ai/chat', 'AIChatView', 'ai:use', 1, g.tenant_id, 1
FROM feature_groups g WHERE g.group_code = 'ai'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '情感识别', 'emotion', 'Heart', '/emotion', 'EmotionView', 'ai:use', 2, g.tenant_id, 0
FROM feature_groups g WHERE g.group_code = 'ai'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '数据看板', 'dashboard', 'Chart', '/dashboard', 'DashboardView', 'data:view', 1, g.tenant_id, 1
FROM feature_groups g WHERE g.group_code = 'analytics'
ON CONFLICT (feature_code) DO NOTHING;

INSERT INTO feature_items (group_id, feature_name, feature_code, icon, route_path, component, permission, sort, tenant_id, is_default)
SELECT g.id, '统计报表', 'reports', 'DataAnalysis', '/reports', 'ReportsView', 'data:view', 2, g.tenant_id, 0
FROM feature_groups g WHERE g.group_code = 'analytics'
ON CONFLICT (feature_code) DO NOTHING;
