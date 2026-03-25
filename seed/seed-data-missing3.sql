-- ============================================
-- 更多空表补充数据 - 正确column名称
-- ============================================

-- DaaS设备租赁 (daas_device_rentals)
INSERT INTO daas_device_rentals (rental_no, contract_id, tenant_id, user_id, device_id, device_sn, device_name, action, deposit_amount, status) VALUES
('RT2024010001', 1, 1, 1, 1, 'SN001', '测试设备1', 'rent', 100.00, 'rented'),
('RT2024010002', 2, 2, 2, 2, 'SN002', '测试设备2', 'rent', 200.00, 'rented');

-- 租户申请 (tenant_applications)
INSERT INTO tenant_applications (application_code, company_name, contact_name, contact_phone, contact_email, status) VALUES
('APP2024010001', '测试公司A', '张三', '13800138001', 'zhangsan@test.com', 'approved'),
('APP2024010002', '测试公司B', '李四', '13800138002', 'lisi@test.com', 'pending');

-- 租户配额 (tenant_quotas)
INSERT INTO tenant_quotas (tenant_id, user_count, device_count, dept_count, store_count) VALUES
('550e8400-e29b-41d4-a716-446655440000', 10, 45, 5, 3),
('550e8400-e29b-41d4-a716-446655440001', 5, 12, 2, 1);

-- 开发者应用 (developer_apps)
INSERT INTO developer_apps (user_id, app_name, app_key, description, platform, status) VALUES
(1, '我的设备管理App', 'ak_app_001', '管理我的智能设备', 'iOS,Android', 1),
(2, '数据分析平台', 'ak_app_002', '设备数据分析', 'Web', 1);

-- 优惠券发放记录 (coupon_grants)
INSERT INTO coupon_grants (coupon_id, member_id, grant_time, status) VALUES
(1, 1, NOW() - INTERVAL '10 days', 1),
(2, 2, NOW() - INTERVAL '5 days', 1);

-- 情绪报告 (emotion_reports)
INSERT INTO emotion_reports (pet_id, report_type, start_date, end_date, summary, emotion_stats, recommendations, generated_at) VALUES
(1, 'daily', NOW() - INTERVAL '1 day', NOW(), '{"happy": 0.75, "sad": 0.05}'::jsonb, '{"neutral": 0.2}'::jsonb, '{"recommendation": "继续保持"}'::jsonb, NOW()),
(2, 'weekly', NOW() - INTERVAL '7 days', NOW(), '{"happy": 0.82, "excited": 0.70}'::jsonb, '{"neutral": 0.1}'::jsonb, '{"recommendation": "增加互动"}'::jsonb, NOW());
