-- ============================================
-- AI Engineering 种子数据
-- ============================================

INSERT INTO ai_training_tasks (task_id, task_name, model_type, dataset_id, hyperparameters, status, progress, current_epoch, total_epochs, metrics, start_time) VALUES
('train_001', '情感识别模型训练', 'emotion_classification', 'ds_001', '{"learning_rate": 0.001, "batch_size": 32}'::jsonb, 'completed', 100, 100, 100, '{"accuracy": 0.95, "loss": 0.05}'::jsonb, NOW() - INTERVAL '1 day'),
('train_002', '动作决策模型训练', 'action_decision', 'ds_002', '{"learning_rate": 0.0005, "batch_size": 64}'::jsonb, 'running', 45, 45, 200, '{"accuracy": 0.82}'::jsonb, NOW() - INTERVAL '2 hours');

INSERT INTO ai_datasets (dataset_id, dataset_name, dataset_type, description, data_format, sample_count, features, labels, status) VALUES
('ds_001', '宠物情感数据集', 'emotion', '用于训练宠物情感识别模型', 'json', 10000, '["image", "audio"]'::jsonb, '["happy", "sad", "angry", "neutral"]'::jsonb, 'active'),
('ds_002', '宠物动作数据集', 'action', '用于训练动作决策模型', 'json', 50000, '["state", "context"]'::jsonb, '["sit", "stand", "walk", "run"]'::jsonb, 'active');

INSERT INTO ai_experiments (experiment_id, experiment_name, hypothesis, description, config, metrics, status, started_at) VALUES
('exp_001', '情感识别A/B测试', '使用更大数据集能提高准确率', '对比不同训练数据量的模型效果', '{"variants": ["control", "treatment"]}'::jsonb, '{"treatment_accuracy": 0.95}'::jsonb, 'running', NOW());

INSERT INTO ai_experiment_groups (group_id, group_name, experiment_id, variant_name, config, metrics, status) VALUES
('eg_001', '对照组', 'exp_001', 'control', '{"dataset_size": 5000}'::jsonb, '{"accuracy": 0.88}'::jsonb, 'completed'),
('eg_002', '实验组', 'exp_001', 'treatment', '{"dataset_size": 10000}'::jsonb, '{"accuracy": 0.95}'::jsonb, 'running');

INSERT INTO ai_monitoring_metrics (metric_id, model_id, model_version, metric_type, metric_name, metric_value, unit, dimensions, collected_at) VALUES
('amm_001', 'model_001', 'v2.1', 'performance', 'accuracy', 0.95, '%', '{"env": "production"}'::jsonb, NOW()),
('amm_002', 'model_001', 'v2.1', 'performance', 'latency', 45.2, 'ms', '{"env": "production"}'::jsonb, NOW());

INSERT INTO ai_alert_rules (rule_id, rule_name, alert_type, condition_expr, threshold_value, severity, notification_channels, is_active) VALUES
('aar_001', '准确率低于阈值', 'performance', 'accuracy < threshold', '0.9', 2, ARRAY['email', 'push'], true),
('aar_002', '延迟过高', 'performance', 'latency > threshold', '100', 1, ARRAY['sms', 'push'], true);

INSERT INTO ai_sandbox_environments (env_id, env_name, env_type, config, resources, status) VALUES
('sandbox_001', '测试环境1', 'emotion_test', '{"timeout": 300}'::jsonb, '{"cpu": "2core", "memory": "4GB"}'::jsonb, 'active'),
('sandbox_002', '测试环境2', 'action_test', '{"timeout": 600}'::jsonb, '{"cpu": "4core", "memory": "8GB"}'::jsonb, 'inactive');

INSERT INTO ai_sandbox_testcases (case_id, env_id, test_name, test_type, input_data, expected_output, actual_output, test_result, execution_time_ms) VALUES
('stc_001', 'sandbox_001', '情感识别边界测试', 'boundary', '{"image": "edge_case.jpg"}'::jsonb, '{"emotion": "neutral"}'::jsonb, '{"emotion": "neutral"}'::jsonb, 'pass', 52);

INSERT INTO ai_decision_logs (log_id, session_id, model_id, model_version, input_data, output_data, decision_reason, confidence_score, execution_time_ms) VALUES
('adl_001', 'sess_001', 'model_001', 'v2.1', '{"image": "test.jpg"}'::jsonb, '{"action": "pet", "emotion": "happy"}'::jsonb, '检测到正面情感', 0.92, 45);

INSERT INTO ai_routing_policies (policy_id, policy_name, policy_type, conditions, actions, priority, is_active) VALUES
('arp_001', '高优先级路由', 'load_balance', '{"load": {"gt": 80}}'::jsonb, '{"route_to": "backup_node"}'::jsonb, 100, true);

-- ============================================
-- Embodied AI 种子数据
-- ============================================

INSERT INTO embodied_maps (map_id, device_id, map_name, map_data, resolution, coverage_area, is_active) VALUES
('map_001', 'dev_001', '家庭平面图', '{"walls": [], "furniture": []}'::jsonb, 0.05, '{"width": 10, "height": 8}'::jsonb, true),
('map_002', 'dev_002', '办公室地图', '{"walls": [], "obstacles": []}'::jsonb, 0.05, '{"width": 15, "height": 12}'::jsonb, true);

INSERT INTO spatial_positions (position_id, device_id, map_id, x, y, z, orientation, timestamp) VALUES
('pos_001', 'dev_001', 'map_001', 2.5, 3.5, 0, 90.0, NOW()),
('pos_002', 'dev_002', 'map_002', 5.0, 6.0, 0, 180.0, NOW());

INSERT INTO action_executions (execution_id, device_id, action_id, action_name, parameters, start_time, end_time, duration_ms, status) VALUES
('ae_001', 'dev_001', 'act_001', 'walk_forward', '{"distance": 1.0}'::jsonb, NOW() - INTERVAL '1 minute', NOW(), 2500, 'success'),
('ae_002', 'dev_001', 'act_002', 'turn_left', '{"angle": 90}'::jsonb, NOW() - INTERVAL '30 seconds', NOW(), 1200, 'success');

INSERT INTO safety_zones (zone_id, device_id, zone_name, zone_type, boundary, center_lat, center_lng, radius, is_active, alert_enabled) VALUES
('sz_001', 'dev_001', '厨房危险区', 'restricted', '{"type": "rect", "x1": 0, "y1": 0, "x2": 2, "y2": 2}'::jsonb, 30.5728, 114.2772, 1.0, true, true),
('sz_002', 'dev_001', '安全休息区', 'safe', '{"type": "circle", "cx": 5, "cy": 5, "r": 2}'::jsonb, 30.5729, 114.2773, 2.0, true, false);

INSERT INTO embodied_decision_logs (log_id, device_id, decision_type, context_data, decision_result, confidence_score, execution_path, latency_ms) VALUES
('edl_001', 'dev_001', 'navigation', '{"current_pos": {"x": 1, "y": 1}}'::jsonb, '{"path": [1,2,3], "action": "move"}'::jsonb, 0.89, '["sense", "plan", "act"]'::jsonb, 32);

INSERT INTO safety_logs (log_id, device_id, event_type, severity, description, location, triggered_at, resolution) VALUES
('sl_001', 'dev_001', 'zone_violation', 'warning', '进入厨房危险区', '{"x": 1.5, "y": 1.5}'::jsonb, NOW() - INTERVAL '10 minutes', '已发出告警');

-- ============================================
-- Platform Ecosystem 种子数据
-- ============================================

INSERT INTO plugins (plugin_id, plugin_name, plugin_type, version, description, author, is_enabled, is_official, install_count, rating) VALUES
('plugin_001', '智能家居集成', 'integration', '1.0.0', '支持主流智能家居设备', 'MiniClaw Team', true, true, 1500, 4.8),
('plugin_002', '健康监测增强', 'health', '1.2.0', '高级健康数据分析', 'MiniClaw Team', true, true, 980, 4.6);

INSERT INTO emotion_packs (pack_id, pack_name, category, preview_urls, emotions, price, is_free, download_count, rating, status) VALUES
('ep_001', '开心表情包', 'emotion', ARRAY['https://cdn.miniclaw.com/emotion/happy1.png'], '["happy", "excited"]'::jsonb, 0, true, 5000, 4.9, 'active'),
('ep_002', '可爱动物包', 'animal', ARRAY['https://cdn.miniclaw.com/emotion/cute1.png'], '["cute", "adorable"]'::jsonb, 6.00, false, 2000, 4.7, 'active');

INSERT INTO actions (action_id, action_name, category, description, animation_data, compatible_models, price, is_free, download_count, status) VALUES
('action_001', '打招呼', 'social', '宠物主动打招呼', '{"frames": 30}'::jsonb, '["miniclaw_v1", "miniclaw_v2"]'::jsonb, 0, true, 8000, 'active'),
('action_002', '后空翻', 'entertainment', '高难度后空翻动作', '{"frames": 60}'::jsonb, '["miniclaw_v2"]'::jsonb, 12.00, false, 1500, 'active');

INSERT INTO voices (voice_id, voice_name, voice_type, preview_url, characteristics, price, is_free, download_count, status) VALUES
('voice_001', '甜美音', 'cute', 'https://cdn.miniclaw.com/voice/sweet_preview.mp3', '{"pitch": "high", "speed": "fast"}'::jsonb, 0, true, 6000, 'active'),
('voice_002', '磁性音', 'magnetic', 'https://cdn.miniclaw.com/voice/magnetic_preview.mp3', '{"pitch": "low", "speed": "slow"}'::jsonb, 8.00, false, 3000, 'active');

INSERT INTO user_purchases (purchase_id, user_id, item_type, item_id, price, currency, payment_method, status, purchased_at) VALUES
('purchase_001', 1, 'emotion_pack', 'ep_002', 6.00, 'CNY', 'wechat', 'completed', NOW() - INTERVAL '5 days'),
('purchase_002', 2, 'action', 'action_002', 12.00, 'CNY', 'alipay', 'completed', NOW() - INTERVAL '3 days');

INSERT INTO ratings (rating_id, user_id, item_type, item_id, rating, review, helpful_count) VALUES
('rating_001', 1, 'emotion_pack', 'ep_002', 5, '非常可爱！', 10),
('rating_002', 2, 'action', 'action_002', 4, '动作流畅', 5);

INSERT INTO integrations (integration_id, integration_name, integration_type, provider, description, is_enabled, last_sync_at) VALUES
('int_001', '小米智能家居', 'smart_home', 'Xiaomi', '米家设备集成', true, NOW() - INTERVAL '1 hour'),
('int_002', '华为健康', 'health', 'Huawei', '华为健康数据同步', true, NOW() - INTERVAL '2 hours');

INSERT INTO smart_home_devices (device_id, household_id, integration_id, device_name, device_type, manufacturer, status, is_active) VALUES
('shd_001', 'hh_001', 'int_001', '米家台灯', 'light', '小米', 'online', true),
('shd_002', 'hh_001', 'int_001', '米家空调', 'ac', '小米', 'online', true);

INSERT INTO smart_home_triggers (trigger_id, household_id, trigger_name, trigger_type, conditions, actions, is_active) VALUES
('sht_001', 'hh_001', '回家开灯', 'location', '{"event": "enter_home"}'::jsonb, '{"action": "turn_on", "device": "shd_001"}'::jsonb, true);

INSERT INTO vet_appointments (appointment_id, pet_id, clinic_name, vet_name, appointment_date, appointment_type, reason, status) VALUES
('vet_001', 'pet_001', '宠物医院A', '张医生', NOW() + INTERVAL '3 days', 'checkup', '年度体检', 'scheduled');

INSERT INTO pet_medical_records (record_id, pet_id, record_type, record_date, clinic_name, vet_name, diagnosis, treatment, medications, cost) VALUES
('pmr_001', 'pet_001', 'vaccination', '2024-01-15', '宠物医院A', '张医生', '健康', '狂犬疫苗接种', '["rabies_vaccine"]'::jsonb, 80.00);

INSERT INTO insurance_policies (policy_id, pet_id, insurance_company, policy_number, plan_name, coverage_start, coverage_end, coverage_amount, premium, status) VALUES
('ins_001', 'pet_001', '宠物保险公司A', 'POL20240001', '基础版', '2024-01-01', '2025-01-01', 50000.00, 1200.00, 'active');

INSERT INTO map_configs (config_id, map_type, config_data, resolution, scale, is_default) VALUES
('mc_001', 'floor_plan', '{"layout": "2d", "units": "meters"}'::jsonb, 0.05, 1.0, true);
