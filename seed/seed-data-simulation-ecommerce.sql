-- ============================================
-- Simulation 仿真测试模块 种子数据
-- ============================================

INSERT INTO simulation_pets (pet_id, simulation_id, virtual_pet_data, initial_state, personality_params) VALUES
('sp_001', 'sim_001', '{"name": "虚拟小白", "type": "cat"}'::jsonb, '{"mood": 80, "energy": 90}'::jsonb, '{"playful": 0.8}'::jsonb),
('sp_002', 'sim_001', '{"name": "虚拟旺财", "type": "dog"}'::jsonb, '{"mood": 75, "energy": 85}'::jsonb, '{"loyal": 1.0}'::jsonb);

INSERT INTO simulation_scenarios (scenario_id, scenario_name, scenario_type, description, environment_config, pet_configs, duration_minutes, difficulty_level, created_by) VALUES
('ss_001', '室内自由活动', 'free_play', '宠物在室内自由探索', '{"area": "indoor", "obstacles": []}'::jsonb, '["cat", "dog"]'::jsonb, 60, 1, 1),
('ss_002', '障碍穿越挑战', 'challenge', '穿越障碍物到达目标', '{"area": "indoor", "obstacles": ["box", "tunnel"]}'::jsonb, '["dog"]'::jsonb, 30, 3, 1);

INSERT INTO test_executions (execution_id, testcase_id, scenario_id, execution_type, status, start_time, end_time, duration_ms, result_data) VALUES
('te_001', 'tc_001', 'ss_001', 'simulation', 'completed', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '50 minutes', 600000, '{"success": true}'::jsonb),
('te_002', 'tc_002', 'ss_002', 'stress_test', 'running', NOW() - INTERVAL '10 minutes', NULL, NULL, NULL);

INSERT INTO test_reports (report_id, execution_id, report_type, summary, details, metrics, recommendations, generated_at) VALUES
('tr_001', 'te_001', 'simulation_report', '{"total_time": 3600, "actions": 50}'::jsonb, '{"events": []}'::jsonb, '{"avg_response_time": 45}'::jsonb, '系统运行正常', NOW());

INSERT INTO playback_records (playback_id, simulation_id, session_data, duration_seconds, status) VALUES
('pr_001', 'sim_001', '{"frames": []}'::jsonb, 3600, 'completed');

INSERT INTO stress_tests (test_id, test_name, test_config, target_metrics, status, progress, current_load, target_load, started_at) VALUES
('st_001', '并发连接测试', '{"concurrent_users": 1000, "duration": 300}'::jsonb, '{"response_time": 200, "error_rate": 0.01}'::jsonb, 'running', 65, 650, 1000, NOW() - INTERVAL '10 minutes');

INSERT INTO simulation_datasets (dataset_id, dataset_name, dataset_type, description, data_schema, record_count, version, status) VALUES
('sd_001', '宠物行为数据集V1', 'behavior', '收集的宠物行为数据', '{"features": ["action", "emotion"]}'::jsonb, 50000, '1.0', 'active'),
('sd_002', '环境感知数据集', 'perception', '环境传感器数据', '{"features": ["image", "depth"]}'::jsonb, 20000, '1.0', 'active');

INSERT INTO simulation_dataset_versions (version_id, dataset_id, version_number, changes, record_count, is_current) VALUES
('sdv_001', 'sd_001', '1.0', '初始版本', 50000, true),
('sdv_002', 'sd_001', '1.1', '新增1000条数据', 51000, false);

INSERT INTO simulation_integrations (integration_id, integration_name, integration_type, endpoint_url, is_enabled, last_test_at) VALUES
('si_001', '外部AI服务', 'ai_service', 'https://ai.external.com/api', true, NOW());

INSERT INTO simulation_cicd_jobs (job_id, job_name, job_type, pipeline_config, trigger_type, status, build_number, branch, started_at) VALUES
('cicd_001', '模型训练CI', 'training', '{"stages": ["build", "test", "deploy"]}'::jsonb, 'push', 'success', 45, 'main', NOW() - INTERVAL '1 hour'),
('cicd_002', '仿真测试CD', 'testing', '{"stages": ["test", "deploy"]}'::jsonb, 'manual', 'running', 12, 'release/v2', NOW() - INTERVAL '5 minutes');

INSERT INTO dataset_jobs (job_id, job_name, dataset_id, operation, config, status, progress, records_processed, started_at) VALUES
('dj_001', '数据清洗任务', 'sd_001', 'clean', '{"remove_duplicates": true}'::jsonb, 'completed', 100, 50000, NOW() - INTERVAL '2 hours'),
('dj_002', '数据增强任务', 'sd_002', 'augment', '{"methods": ["flip", "rotate"]}'::jsonb, 'running', 45, 9000, NOW() - INTERVAL '30 minutes');

INSERT INTO ab_experiments (experiment_id, experiment_name, hypothesis, variants, traffic_allocation, metrics, status, start_date) VALUES
('ab_001', '新推荐算法测试', '新算法能提高转化率', '{"control": 0.5, "treatment": 0.5}'::jsonb, '{"control": 50, "treatment": 50}'::jsonb, '{"ctr": 0.05}'::jsonb, 'running', '2024-02-01');

INSERT INTO dataset_samples (sample_id, dataset_id, sample_data, label, split_type, quality_score) VALUES
('dsamp_001', 'sd_001', '{"action": "sit", "emotion": "happy"}'::jsonb, 'sit', 'train', 0.95),
('dsamp_002', 'sd_001', '{"action": "run", "emotion": "excited"}'::jsonb, 'run', 'train', 0.92),
('dsamp_003', 'sd_001', '{"action": "sleep", "emotion": "peaceful"}'::jsonb, 'sleep', 'validation', 0.88);

-- ============================================
-- PRD_28: 宠物电商 种子数据
-- ============================================

INSERT INTO eco_recommended_products (recommendation_id, product_id, pet_id, user_id, recommendation_type, score, reason, is_clicked, is_purchased) VALUES
('rec_001', 'prod_001', 'pet_001', 1, 'similar_pet', 0.95, '和您的宠物很配', true, false),
('rec_002', 'prod_002', 'pet_002', 1, 'bestseller', 0.88, '热销产品', false, false);

INSERT INTO eco_cart (cart_id, user_id, product_id, quantity, price, added_at) VALUES
('cart_001', 1, 'prod_001', 2, 29.90, NOW() - INTERVAL '1 day'),
('cart_002', 1, 'prod_003', 1, 59.00, NOW() - INTERVAL '2 hours');

INSERT INTO eco_orders (order_id, order_number, user_id, total_amount, discount_amount, shipping_amount, tax_amount, final_amount, status, shipping_address, paid_at) VALUES
('order_001', 'ECO2024020001', 1, 89.90, 10.00, 0, 4.79, 84.69, 'paid', '{"name": "小明", "phone": "13800138000", "address": "武汉市..."}'::jsonb, NOW() - INTERVAL '1 day'),
('order_002', 'ECO2024020002', 2, 159.00, 0, 10, 10.14, 179.14, 'shipped', '{"name": "小红", "phone": "13900139000", "address": "北京市..."}'::jsonb, NOW() - INTERVAL '2 days');

-- ============================================
-- PRD_30: 家庭相册增强 种子数据
-- ============================================

INSERT INTO album_photos (photo_id, album_id, uploader_id, photo_url, thumbnail_url, caption, location, taken_at, ai_tags, face_count, like_count) VALUES
('photo_001', 'album_001', 1, 'https://cdn.miniclaw.com/photos/001.jpg', 'https://cdn.miniclaw.com/photos/001_thumb.jpg', '小白的萌照', '{"lat": 30.5728, "lng": 114.2772}'::jsonb, NOW() - INTERVAL '3 days', ARRAY['cat', 'cute'], 1, 25),
('photo_002', 'album_001', 1, 'https://cdn.miniclaw.com/photos/002.jpg', 'https://cdn.miniclaw.com/photos/002_thumb.jpg', '旺财在公园', '{"lat": 30.5729, "lng": 114.2773}'::jsonb, NOW() - INTERVAL '5 days', ARRAY['dog', 'outdoor'], 1, 18);

INSERT INTO album_ai_albums (ai_album_id, ai_album_name, album_type, criteria, cover_photo_url, photo_count) VALUES
('ai_album_001', '猫咪时光', 'face', '{"pet_type": "cat"}'::jsonb, 'https://cdn.miniclaw.com/photos/001.jpg', 15),
('ai_album_002', '户外活动', 'location', '{"location_type": "park"}'::jsonb, 'https://cdn.miniclaw.com/photos/002.jpg', 8);

INSERT INTO album_categories (category_id, category_name, parent_id, sort_order, is_active) VALUES
('cat_001', '宠物照片', NULL, 1, true),
('cat_002', '家庭活动', NULL, 2, true),
('cat_003', '猫咪', 'cat_001', 1, true),
('cat_004', '狗狗', 'cat_001', 2, true);

INSERT INTO album_family_members (member_id, album_id, user_id, face_encoding, face_photo_url, name, relation, is_verified) VALUES
('afm_001', 'album_001', 1, '{"encoding": "..."}'::jsonb, 'https://cdn.miniclaw.com/faces/001.jpg', '小明', 'owner', true),
('afm_002', 'album_001', 2, '{"encoding": "..."}'::jsonb, 'https://cdn.miniclaw.com/faces/002.jpg', '小白', 'pet', true);
