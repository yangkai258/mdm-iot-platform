-- ============================================
-- 空表补充数据 - 修正版
-- ============================================

-- 老人档案
INSERT INTO elderly_profiles (profile_id, device_id, name, nickname, age, emergency_contact, emergency_phone, caregiver_user_id, status) VALUES
('elderly_001', 'dev_001', '张爷爷', '老张', 78, '张叔叔', '13800138001', '1', 1),
('elderly_002', 'dev_002', '李奶奶', '老李', 82, '李阿姨', '13800138002', '2', 1);

-- 儿童档案
INSERT INTO children_profiles (profile_id, device_id, name, nickname, age, parent_user_id, status) VALUES
('child_001', 'dev_001', '小明', '小名', 8, '3', 1),
('child_002', 'dev_002', '小红', '红红', 6, '4', 1);

-- 老人提醒
INSERT INTO elderly_reminders (reminder_id, profile_id, title, content, reminder_type, medicine_name, schedule_time, schedule_type, is_enabled) VALUES
('reminder_001', 'elderly_001', '用药提醒', '该吃降压药了', 'medication', '降压药', '08:00', 'daily', true),
('reminder_002', 'elderly_001', '运动提醒', '该散步了', 'exercise', NULL, '18:00', 'daily', true),
('reminder_003', 'elderly_002', '用药提醒', '该吃降糖药了', 'medication', '降糖药', '07:30', 'daily', true);

-- 老人关怀配置
INSERT INTO elderly_care_configs (user_id, device_id, is_enabled, health_monitor_enabled, heart_rate_alert_high, heart_rate_alert_low, activity_goal, fall_detection_enabled, emergency_contact_name, emergency_contact_phone) VALUES
(1, 'dev_001', true, true, 100, 50, 6000, true, '张叔叔', '13800138001'),
(2, 'dev_002', true, true, 95, 55, 5000, true, '李阿姨', '13800138002');

-- 儿童模式配置
INSERT INTO child_mode_configs (user_id, device_id, is_enabled, content_filter_level, daily_time_limit, allowed_start_time, allowed_end_time) VALUES
(3, 'dev_001', true, 'strict', 120, '08:00', '20:00'),
(4, 'dev_002', true, 'moderate', 90, '09:00', '19:00');

-- 虚拟宠物
INSERT INTO virtual_pets (pet_id, name, species, personality, mood, health, hunger, energy, happiness, age, create_user_id) VALUES
('vp_001', '小v', 'cat', '好奇', 'happy', 100, 30, 80, 90, 2, 1),
('vp_002', '旺财', 'dog', '忠诚', 'excited', 95, 45, 85, 95, 3, 2),
('vp_003', '咪咕', 'bird', '活泼', 'happy', 90, 20, 90, 90, 1, 1);

-- 性能指标
INSERT INTO performance_metrics (device_id, metric_name, metric_type, metric_value, unit, source, tenant_id) VALUES
('dev_001', 'cpu_usage', 'cpu', 45.5, '%', 'system', '550e8400-e29b-41d4-a716-446655440000'),
('dev_001', 'memory_usage', 'memory', 62.3, '%', 'system', '550e8400-e29b-41d4-a716-446655440000'),
('dev_001', 'temperature', 'thermal', 38.2, '°C', 'sensor', '550e8400-e29b-41d4-a716-446655440000'),
('dev_002', 'cpu_usage', 'cpu', 38.2, '%', 'system', '550e8400-e29b-41d4-a716-446655440000'),
('dev_002', 'network_latency', 'network', 25, 'ms', 'ping', '550e8400-e29b-41d4-a716-446655440000');

-- 审批历史
INSERT INTO approval_histories (application_id, action, action_text, operator, comment) VALUES
(1, 'approved', '审批通过', 'admin', '设备符合安全标准'),
(2, 'approved', '审批通过', 'admin', '固件通过测试'),
(3, 'rejected', '审批拒绝', 'admin', '资质不符合');

-- 地理围栏告警
INSERT INTO geofence_alerts (rule_id, device_id, alert_type, latitude, longitude, severity, message, status) VALUES
(1, 'dev_001', 'enter', 30.5728, 114.2772, 2, '设备进入安全区域', 1),
(2, 'dev_002', 'exit', 31.2304, 121.4737, 3, '设备离开安全区域', 1);

-- 家庭相册评论 (需要先有photo_uuid，这里用占位符)
INSERT INTO family_album_comments (photo_uuid, user_id, user_name, content) VALUES
('photo_001', 2, '用户2', '这张照片好可爱！'),
('photo_001', 3, '用户3', '确实，萌化了'),
('photo_002', 1, '用户1', '值得纪念的一天');

-- 家庭相册点赞
INSERT INTO family_album_likes (photo_uuid, user_id) VALUES
('photo_001', 2),
('photo_001', 3),
('photo_002', 1);

-- 宠物社交动态评论
INSERT INTO pet_social_comments (post_id, pet_id, content) VALUES
(1, 'dev_001', '太可爱了！'),
(1, 'dev_002', '想rua一下'),
(2, 'dev_001', '看起来很开心');

-- 宠物社交关注
INSERT INTO pet_social_follows (follower_id, following_id) VALUES
('dev_001', 'dev_002'),
('dev_002', 'dev_001'),
('dev_003', 'dev_001'),
('dev_001', 'dev_003');

-- 宠物社交点赞
INSERT INTO pet_social_likes (post_id, pet_id) VALUES
(1, 'dev_002'),
(1, 'dev_003'),
(2, 'dev_001');
