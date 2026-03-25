-- ============================================
-- 新增缺失表的种子数据
-- ============================================

-- Sprint 13: 全球化区域数据
INSERT INTO regions (region_code, region_name, region_name_en, continent, country, timezone, ai_node_url, is_active, sort_order) VALUES
('CN-EAST', '华东区域', 'East China', 'Asia', 'China', 'Asia/Shanghai', 'https://ai-cn-east.miniclaw.com', true, 1),
('CN-NORTH', '华北区域', 'North China', 'Asia', 'China', 'Asia/Shanghai', 'https://ai-cn-north.miniclaw.com', true, 2),
('US-EAST', '美东区域', 'US East', 'North America', 'United States', 'America/New_York', 'https://ai-us-east.miniclaw.com', true, 3),
('EU-WEST', '欧洲区域', 'EU West', 'Europe', 'Germany', 'Europe/Berlin', 'https://ai-eu-west.miniclaw.com', true, 4);

INSERT INTO region_nodes (node_id, region_id, node_name, node_type, endpoint_url, health_status, is_master, weight) VALUES
('node-cn-east-1', 1, '华东节点1', 'primary', 'https://cn-east-1.miniclaw.com', 'healthy', true, 100),
('node-cn-east-2', 1, '华东节点2', 'replica', 'https://cn-east-2.miniclaw.com', 'healthy', false, 80),
('node-us-east-1', 3, '美东节点1', 'primary', 'https://us-east-1.miniclaw.com', 'healthy', true, 100);

INSERT INTO data_residency_rules (rule_id, data_type, source_region, allowed_regions, retention_days, description, is_active) VALUES
('data-rule-001', 'user_data', 'CN', ARRAY['CN'], 365, '中国用户数据必须留在国内', true),
('data-rule-002', 'device_data', 'CN', ARRAY['CN', 'HK'], 180, '设备数据可在大陆和香港', true);

INSERT INTO timezone_configs (config_key, timezone_id, utc_offset, description, is_active) VALUES
('tz-shanghai', 'Asia/Shanghai', '+08:00', '北京时间', true),
('tz-tokyo', 'Asia/Tokyo', '+09:00', '东京时间', true),
('tz-newyork', 'America/New_York', '-05:00', '纽约时间', true);

-- Sprint 14: AI系统工程
INSERT INTO ai_behavior_logs (log_id, device_id, behavior_type, input_data, output_data, confidence_score, model_version, latency_ms) VALUES
('log-001', 'dev_001', 'emotion_recognition', '{"image": "base64..."}'::jsonb, '{"emotion": "happy", "confidence": 0.95}'::jsonb, 0.95, 'model_v2.1', 45),
('log-002', 'dev_002', 'action_decision', '{"context": "playing"}'::jsonb, '{"action": "jump", "confidence": 0.88}'::jsonb, 0.88, 'model_v2.1', 32);

INSERT INTO ai_sandbox_tests (test_id, test_name, test_type, test_scenario, expected_output, actual_output, test_result, executed_at) VALUES
('test-001', '情感识别测试', 'emotion', '{"input": "happy_face.jpg"}'::jsonb, '{"emotion": "happy"}'::jsonb, '{"emotion": "happy"}'::jsonb, 'pass', NOW()),
('test-002', '动作决策测试', 'action', '{"context": "idle"}'::jsonb, '{"action": "sit"}'::jsonb, '{"action": "sit"}'::jsonb, 'pass', NOW());

INSERT INTO model_rollback_records (record_id, model_id, from_version, to_version, rollback_reason, rollback_status, performed_by, rolled_back_at) VALUES
('rollback-001', 'model-001', 'v2.2', 'v2.1', '发现v2.2存在决策偏差', 'completed', 'admin', NOW() - INTERVAL '2 days');

-- Sprint 15: 宠物生态
INSERT INTO pets (pet_id, device_id, owner_id, pet_name, species, breed, birth_date, gender, weight, personality, health_status) VALUES
('pet_001', 'dev_001', 1, '小白', 'cat', '中华田园猫', '2022-03-15', 'female', 4.5, '{"playful": 0.8, "friendly": 0.9}'::jsonb, 'healthy'),
('pet_002', 'dev_002', 1, '旺财', 'dog', '金毛', '2021-06-20', 'male', 25.0, '{"loyal": 1.0, "energetic": 0.9}'::jsonb, 'healthy'),
('pet_003', 'dev_003', 2, '咪咕', 'bird', '虎皮鹦鹉', '2023-01-10', 'male', 0.1, '{"curious": 0.9, "talkative": 0.7}'::jsonb, 'healthy');

INSERT INTO pet_device_bindings (binding_id, pet_id, device_id, binding_type, is_active) VALUES
('pdb_001', 'pet_001', 'dev_001', 'primary', true),
('pdb_002', 'pet_002', 'dev_002', 'primary', true);

INSERT INTO lost_found_reports (report_id, pet_id, report_type, last_seen_location, last_seen_lat, last_seen_lng, description, status) VALUES
('lfr_001', 'pet_001', 'found', '小区花园', 30.5728, 114.2772, '发现一只橘猫，很亲人', 'resolved');

INSERT INTO household_pet_invites (invite_id, pet_id, inviter_id, invitee_id, status, expires_at) VALUES
('hpi_001', 'pet_001', 1, 2, 'accepted', NOW() + INTERVAL '7 days');

-- Sprint 16: 商业化/订阅
INSERT INTO webhooks (webhook_id, user_id, webhook_name, endpoint_url, secret_key, events, is_active, retry_count) VALUES
('wh_001', 1, '设备状态通知', 'https://example.com/webhook/device', 'secret_xxx', ARRAY['device.status', 'device.alert'], true, 3),
('wh_002', 2, '订单通知', 'https://example.com/webhook/order', 'secret_yyy', ARRAY['order.created', 'order.paid'], true, 2);

INSERT INTO webhook_logs (log_id, webhook_id, event_type, payload, response_code, status, sent_at) VALUES
('wlog_001', 1, 'device.status', '{"device_id": "dev_001", "status": "online"}'::jsonb, 200, 'success', NOW());

INSERT INTO billing_statements (statement_id, user_id, billing_period, statement_date, due_date, subtotal, discount, tax, total, status) VALUES
('stmt_001', 1, '2024-01', '2024-02-01', '2024-02-15', 99.00, 10.00, 8.90, 97.90, 'paid'),
('stmt_002', 2, '2024-01', '2024-02-01', '2024-02-15', 199.00, 20.00, 17.90, 196.90, 'pending');

INSERT INTO invoices (invoice_id, invoice_number, user_id, billing_amount, tax_rate, tax_amount, total_amount, invoice_type, company_name, tax_id, status) VALUES
('inv_001', 'INV2024010001', 1, 97.90, 0.06, 5.87, 103.77, 'company', '测试公司A', '91330000MA2XXXXXX', 'issued'),
('inv_002', 'INV2024010002', 2, 196.90, 0.06, 11.81, 208.71, 'personal', NULL, NULL, 'issued');

INSERT INTO usage_records (record_id, user_id, subscription_id, usage_type, usage_value, unit, period_start, period_end, cost) VALUES
('ur_001', 1, 'sub_001', 'api_calls', 15000, 'count', '2024-01-01', '2024-01-31', 15.00),
('ur_002', 2, 'sub_002', 'api_calls', 50000, 'count', '2024-01-01', '2024-01-31', 40.00);

-- Sprint 17: 情感计算
INSERT INTO pet_emotion_actions (action_id, pet_id, emotion_type, action_type, intensity, trigger_source, action_result) VALUES
('pea_001', 'pet_001', 'happy', 'tail_wag', 0.8, 'interaction', 'success'),
('pea_002', 'pet_002', 'excited', 'jump', 0.9, 'command', 'success');

INSERT INTO family_emotions (record_id, household_id, member_id, emotion_type, intensity, context, detected_at) VALUES
('fe_001', 'hh_001', 1, 'happy', 0.85, '家庭聚会', NOW());

-- Sprint 18: 数字孪生
INSERT INTO digital_twin_pets (twin_id, pet_id, twin_state, sync_status, last_sync_at) VALUES
('twin_001', 'pet_001', '{"mood": "happy", "energy": 85}'::jsonb, 'synced', NOW()),
('twin_002', 'pet_002', '{"mood": "excited", "energy": 95}'::jsonb, 'synced', NOW());

INSERT INTO behavior_events (event_id, pet_id, device_id, event_type, event_data, timestamp) VALUES
('be_001', 'pet_001', 'dev_001', 'movement', '{"x": 100, "y": 200}'::jsonb, NOW()),
('be_002', 'pet_002', 'dev_002', 'feeding', '{"food_type": "dry", "amount": 50}'::jsonb, NOW());

INSERT INTO highlight_moments (moment_id, pet_id, moment_type, media_url, description, emotion_tags, captured_at) VALUES
('hm_001', 'pet_001', 'playful', 'https://cdn.miniclaw.com/moments/001.mp4', '玩耍时刻', ARRAY['happy', 'playful'], NOW()),
('hm_002', 'pet_002', 'cute', 'https://cdn.miniclaw.com/moments/002.jpg', '睡觉萌照', ARRAY['cute', 'peaceful'], NOW());

INSERT INTO sync_records (sync_id, device_id, sync_type, sync_data, status, completed_at) VALUES
('sync_001', 'dev_001', 'state_sync', '{"state": "updated"}'::jsonb, 'success', NOW());

-- Sprint 19: 健康医疗
INSERT INTO disease_patterns (pattern_id, pattern_name, disease_name, species, symptoms, severity, description) VALUES
('dp_001', '猫感冒', 'feline_upper_respiratory', 'cat', '["sneezing", "runny_nose", "lethargy"]'::jsonb, 'mild', '猫咪上呼吸道感染'),
('dp_002', '狗皮肤病', 'canine_dermatitis', 'dog', '["itching", "hair_loss", "redness"]'::jsonb, 'moderate', '狗狗皮肤炎症');

INSERT INTO sleep_records (record_id, pet_id, sleep_start, sleep_end, duration_minutes, sleep_quality, sleep_stages, interruptions) VALUES
('sleep_001', 'pet_001', NOW() - INTERVAL '8 hours', NOW() - INTERVAL '2 hours', 360, 'good', '{"deep": 180, "light": 120, "rem": 60}'::jsonb, 1),
('sleep_002', 'pet_002', NOW() - INTERVAL '10 hours', NOW() - INTERVAL '1 hour', 540, 'excellent', '{"deep": 270, "light": 180, "rem": 90}'::jsonb, 0);

INSERT INTO health_baselines (baseline_id, pet_id, metric_type, baseline_value, normal_range_min, normal_range_max, measurement_unit, last_measured_at) VALUES
('hb_001', 'pet_001', 'weight', 4.5, 3.5, 5.5, 'kg', NOW()),
('hb_002', 'pet_002', 'weight', 25.0, 22.0, 28.0, 'kg', NOW());

INSERT INTO health_reports (report_id, pet_id, report_type, period_start, period_end, health_summary, activity_summary, recommendations, generated_at) VALUES
('hr_001', 'pet_001', 'weekly', NOW() - INTERVAL '7 days', NOW(), '{"overall": "good", "concerns": []}'::jsonb, '{"steps": 5000, "active_hours": 6}'::jsonb, '{"recommendation": "继续保持当前饮食"}'::jsonb, NOW());

-- Sprint 20: 家庭场景
INSERT INTO households (household_id, household_name, owner_id, address, family_photo_url, is_active) VALUES
('hh_001', '小明家', 1, '武汉市洪山区某小区', 'https://cdn.miniclaw.com/family/001.jpg', true),
('hh_002', '小红家', 2, '北京市朝阳区某小区', 'https://cdn.miniclaw.com/family/002.jpg', true);

INSERT INTO household_members (member_id, household_id, user_id, member_name, relation, role, is_active) VALUES
('hm_001', 'hh_001', 1, '小明', 'owner', 'admin', true),
('hm_002', 'hh_001', 2, '小红', 'spouse', 'member', true);

INSERT INTO household_invites (invite_id, household_id, inviter_id, invitee_email, invite_token, role, status, expires_at) VALUES
('hi_001', 'hh_001', 1, 'zhangsan@example.com', 'token_xxx', 'member', 'pending', NOW() + INTERVAL '7 days');

INSERT INTO family_mode_configs (config_id, household_id, mode_type, is_enabled, settings) VALUES
('fmc_001', 'hh_001', 'child_mode', true, '{"time_limit": 120, "content_filter": "strict"}'::jsonb);

INSERT INTO interaction_records (record_id, pet_id, device_id, user_id, interaction_type, interaction_data, duration_seconds, feedback) VALUES
('ir_001', 'pet_001', 'dev_001', 1, 'play', '{"toy": "feather"}'::jsonb, 300, 'positive'),
('ir_002', 'pet_002', 'dev_002', 1, 'walk', '{"duration": 30}'::jsonb, 1800, 'positive');

INSERT INTO family_activities (activity_id, household_id, activity_type, title, description, scheduled_at, duration_minutes, participants, status) VALUES
('fa_001', 'hh_001', 'outing', '周末郊游', '带宠物去公园', NOW() + INTERVAL '2 days', 180, ARRAY[1, 2], 'scheduled');

-- API配额
INSERT INTO api_quotas (quota_id, user_id, api_name, quota_limit, quota_used, period_type, period_start, period_end, is_active) VALUES
('aq_001', 1, '/api/v1/devices', 10000, 1500, 'monthly', '2024-01-01', '2024-01-31', true),
('aq_002', 2, '/api/v1/devices', 50000, 8000, 'monthly', '2024-01-01', '2024-01-31', true);
