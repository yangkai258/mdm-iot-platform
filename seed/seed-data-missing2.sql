-- ============================================
-- 空表补充数据 - 修正column名称
-- ============================================

-- 积分规则 (points_rules)
INSERT INTO points_rules (rule_code, rule_name, rule_type, points, amount, remark, status) VALUES
('points_001', '消费得积分', 1, 1, 1.00, '每消费1元得1积分', 1),
('points_002', '充值得积分', 1, 2, 1.00, '每充值1元得2积分', 1),
('points_003', '推荐得积分', 2, 500, NULL, '推荐新用户得500积分', 1);

-- 策略表 (policies)
INSERT INTO policies (name, policy_type, description, priority, config_ids, rule_ids, enabled, status, platform, scope) VALUES
('数据保留策略', 'data_retention', '用户数据保留期限策略', 50, '[]'::jsonb, '[]'::jsonb, true, 'active', 'all', 'all'),
('设备安全策略', 'device_security', '设备安全配置策略', 80, '[]'::jsonb, '[]'::jsonb, true, 'active', 'all', 'all');

-- 策略绑定 (policy_bindings) - 需要先有policies的id
INSERT INTO policy_bindings (policy_id, target_type, target_id, target_name, bound_by, bound_at, status) VALUES
(1, 'tenant', '550e8400-e29b-41d4-a716-446655440000', '测试租户', 'admin', NOW(), 1),
(2, 'device_group', 'default', '默认设备组', 'admin', NOW(), 1);

-- 通知日志 (notification_logs)
INSERT INTO notification_logs (channel_type, recipient, subject, content, status, sent_at) VALUES
('email', 'user1@example.com', '设备告警', '您的设备电量低于20%', 'sent', NOW()),
('sms', '13800138001', '健康提醒', '请记得测量血压', 'sent', NOW()),
('push', 'device_token_xxx', '固件更新', '新版本v1.2.5已发布', 'sent', NOW());

-- 指令历史 (command_histories)
INSERT INTO command_histories (device_id, cmd_type, action, status, sent_at) VALUES
('dev_001', 'control', 'restart', 'success', NOW()),
('dev_001', 'ota', 'update', 'success', NOW()),
('dev_002', 'config', 'update', 'success', NOW());

-- DaaS合同 (daas_contracts)
INSERT INTO daas_contracts (contract_no, tenant_id, user_id, device_id, plan_name, daily_rate, monthly_rate, deposit_amount, contract_period, start_date, end_date, status) VALUES
('CT2024010001', 1, 1, 1, '企业版', 0.50, 15.00, 100.00, 12, NOW() - INTERVAL '30 days', NOW() + INTERVAL '335 days', 'active'),
('CT2024010002', 2, 2, 2, '旗舰版', 0.40, 12.00, 200.00, 24, NOW() - INTERVAL '60 days', NOW() + INTERVAL '665 days', 'active');

-- DaaS计费记录 (daas_billings)
INSERT INTO daas_billings (bill_no, tenant_id, user_id, contract_id, bill_type, period_start, period_end, days, daily_rate, amount, status) VALUES
('Bill2024010001', 1, 1, 1, 'rental', NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 day', 29, 0.50, 14.50, 'paid'),
('Bill2024010002', 2, 2, 2, 'rental', NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 day', 29, 0.40, 11.60, 'paid');

-- DaaS设备租赁 (daas_device_rentals)
INSERT INTO daas_device_rentals (device_id, tenant_id, contract_id, start_date, monthly_rate, deposit_amount, status) VALUES
('dev_001', 1, 1, NOW() - INTERVAL '30 days', 15.00, 100.00, 'active'),
('dev_002', 2, 2, NOW() - INTERVAL '15 days', 12.00, 200.00, 'active');

-- 租户申请 (tenant_applications)
INSERT INTO tenant_applications (company_name, contact_name, contact_email, contact_phone, plan_type, status, submitted_at) VALUES
('测试公司A', '张三', 'zhangsan@test.com', '13800138001', 'enterprise', 'approved', NOW() - INTERVAL '30 days'),
('测试公司B', '李四', 'lisi@test.com', '13800138002', 'standard', 'pending', NOW() - INTERVAL '2 days');

-- 租户配额 (tenant_quotas)
INSERT INTO tenant_quotas (tenant_id, quota_key, quota_value, used_value, unit) VALUES
(1, 'devices', 100, 45, 'count'),
(1, 'storage', 1000, 350, 'GB'),
(2, 'devices', 50, 12, 'count'),
(2, 'storage', 500, 120, 'GB');

-- 开发者应用 (developer_apps)
INSERT INTO developer_apps (developer_id, app_name, app_type, description, app_key, app_secret, status) VALUES
(1, '我的设备管理App', 'mobile', '管理我的智能设备', 'ak_xxxxxxxxxxxx', 'sk_xxxxxxxxxxxx', 'active'),
(2, '数据分析平台', 'web', '设备数据分析', 'ak_yyyyyyyyyyyy', 'sk_yyyyyyyyyyyy', 'active');

-- 优惠券发放记录 (coupon_grants)
INSERT INTO coupon_grants (coupon_id, user_id, granted_by, source_type, source_no) VALUES
(1, 1, 1, 'register', 'user_001'),
(2, 2, 1, 'promotion', 'promo_001');

-- 情绪报告 (emotion_reports)
INSERT INTO emotion_reports (device_id, report_type, period_start, period_end, avg_happiness, avg_excitement, dominant_emotion, insights) VALUES
('dev_001', 'daily', NOW() - INTERVAL '1 day', NOW(), 0.75, 0.60, 'happy', '宠物今日情绪稳定，积极性较高'),
('dev_002', 'weekly', NOW() - INTERVAL '7 days', NOW(), 0.82, 0.70, 'excited', '本周宠物活跃度显著提升');
