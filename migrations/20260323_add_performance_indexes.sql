-- Migration: 20260323_add_performance_indexes.sql
-- Description: 添加高频查询字段索引以优化数据库性能
-- Date: 2026-03-23
-- Author: db-optimization agent

BEGIN;

-- =====================================================
-- 1. HealthWarning 健康预警 - 缺失索引
-- =====================================================

-- 宠物+租户+状态 组合查询 (预警列表筛选)
CREATE INDEX IF NOT EXISTS idx_health_warnings_pet_tenant_status 
ON health_warnings(pet_uuid, tenant_id, status);

-- 宠物+类别 组合查询 (按类别筛选)
CREATE INDEX IF NOT EXISTS idx_health_warnings_pet_category 
ON health_warnings(pet_uuid, category);

-- 宠物+级别 组合查询 (按级别筛选)
CREATE INDEX IF NOT EXISTS idx_health_warnings_pet_level 
ON health_warnings(pet_uuid, level);

-- 租户+状态+时间 组合查询 (仪表盘)
CREATE INDEX IF NOT EXISTS idx_health_warnings_tenant_status_time 
ON health_warnings(tenant_id, status, created_at);

-- =====================================================
-- 2. HealthAlert 健康预警 - 缺失索引
-- =====================================================

-- 宠物+状态 组合查询 (活跃预警)
CREATE INDEX IF NOT EXISTS idx_health_alerts_pet_status 
ON health_alerts(pet_uuid, status);

-- 租户+状态+时间 组合查询 (时间范围查询)
CREATE INDEX IF NOT EXISTS idx_health_alerts_tenant_status_time 
ON health_alerts(tenant_id, status, occurred_at);

-- =====================================================
-- 3. ExerciseRecord 运动记录 - 缺失索引
-- =====================================================

-- 宠物+运动类型+时间 组合查询 (运动历史)
CREATE INDEX IF NOT EXISTS idx_exercise_records_pet_type_time 
ON exercise_records(pet_uuid, exercise_type, start_time);

-- 租户+时间 组合查询 (租户运动统计)
CREATE INDEX IF NOT EXISTS idx_exercise_records_tenant_time 
ON exercise_records(tenant_id, start_time);

-- 宠物+日期 组合查询 (每日汇总)
CREATE INDEX IF NOT EXISTS idx_exercise_records_pet_date 
ON exercise_records(pet_uuid, summary_date);

-- =====================================================
-- 4. SleepRecord 睡眠记录 - 缺失索引
-- =====================================================

-- 宠物+睡眠日期 组合查询 (睡眠历史)
CREATE INDEX IF NOT EXISTS idx_sleep_records_pet_date 
ON sleep_records(pet_uuid, sleep_date);

-- 租户+睡眠日期 组合查询 (租户睡眠统计)
CREATE INDEX IF NOT EXISTS idx_sleep_records_tenant_date 
ON sleep_records(tenant_id, sleep_date);

-- =====================================================
-- 5. VitalRecord 生命体征记录 - 缺失索引
-- =====================================================

-- 宠物+体征类型+时间 组合查询 (体征历史时间线)
CREATE INDEX IF NOT EXISTS idx_vital_records_pet_type_time 
ON vital_records(pet_uuid, vital_type, recorded_at);

-- 租户+时间 组合查询 (租户体征统计)
CREATE INDEX IF NOT EXISTS idx_vital_records_tenant_time 
ON vital_records(tenant_id, recorded_at);

-- 宠物+是否异常+时间 组合查询 (异常体征查询)
CREATE INDEX IF NOT EXISTS idx_vital_records_pet_abnormal_time 
ON vital_records(pet_uuid, is_abnormal, recorded_at);

-- =====================================================
-- 6. BehaviorEvent 行为事件 - 缺失索引
-- =====================================================

-- 宠物+行为类型+时间 组合查询 (行为历史)
CREATE INDEX IF NOT EXISTS idx_behavior_events_pet_type_time 
ON behavior_events(pet_uuid, behavior_type, start_time);

-- 租户+时间 组合查询 (租户行为统计)
CREATE INDEX IF NOT EXISTS idx_behavior_events_tenant_time 
ON behavior_events(tenant_id, start_time);

-- 宠物+异常标记 组合查询 (异常行为查询)
CREATE INDEX IF NOT EXISTS idx_behavior_events_pet_anomaly 
ON behavior_events(pet_uuid, is_anomaly);

-- =====================================================
-- 7. AIBehaviorLog AI行为日志 - 缺失索引
-- =====================================================

-- 设备+创建时间 组合查询 (设备行为日志)
CREATE INDEX IF NOT EXISTS idx_ai_behavior_logs_device_time 
ON ai_behavior_logs(device_id, created_at);

-- 状态+创建时间 组合查询 (状态筛选+时间排序)
CREATE INDEX IF NOT EXISTS idx_ai_behavior_logs_status_time 
ON ai_behavior_logs(status, created_at);

-- =====================================================
-- 8. PetLongTermMemory 长期记忆 - 缺失索引
-- =====================================================

-- 用户+设备 组合查询 (用户记忆查询)
CREATE INDEX IF NOT EXISTS idx_long_term_memory_user_device 
ON long_term_memory(user_id, device_id);

-- 记忆类别+设备 组合查询 (类别记忆查询)
CREATE INDEX IF NOT EXISTS idx_long_term_memory_category_device 
ON long_term_memory(memory_category, device_id);

-- =====================================================
-- 9. DeviceMetric 设备指标 - 缺失索引
-- =====================================================

-- 设备+指标类型+时间戳 组合查询 (指标历史)
CREATE INDEX IF NOT EXISTS idx_device_metrics_device_type_time 
ON device_metrics(device_id, metric_type, timestamp);

-- 租户+时间戳 组合查询 (租户指标统计)
CREATE INDEX IF NOT EXISTS idx_device_metrics_tenant_time 
ON device_metrics(tenant_id, timestamp);

-- =====================================================
-- 10. DeviceAlertRule 设备告警规则 - 缺失索引
-- =====================================================

-- 启用状态+告警类型 组合查询 (规则匹配)
CREATE INDEX IF NOT EXISTS idx_device_alert_rules_enabled_type 
ON device_alert_rules(enabled, alert_type);

-- =====================================================
-- 11. CommandHistory 指令历史 - 缺失索引
-- =====================================================

-- 设备+状态+时间 组合查询 (待处理指令)
-- 注意: CommandHistory 无 TableName，GORM 默认复数化
CREATE INDEX IF NOT EXISTS idx_command_history_device_status_time 
ON command_histories(device_id, status, created_at);

-- 指令类型+时间 组合查询 (类型统计)
CREATE INDEX IF NOT EXISTS idx_command_history_type_time 
ON command_histories(cmd_type, created_at);

-- =====================================================
-- 12. HealthAlertRule 健康预警规则 - 缺失索引
-- =====================================================

-- 宠物+启用状态 组合查询 (宠物规则)
CREATE INDEX IF NOT EXISTS idx_health_alert_rules_pet_enabled 
ON health_alert_rules(pet_uuid, is_enabled);

-- 租户+启用状态 组合查询 (租户规则)
CREATE INDEX IF NOT EXISTS idx_health_alert_rules_tenant_enabled 
ON health_alert_rules(tenant_id, is_enabled);

-- =====================================================
-- 13. AIRollbackTask AI回滚任务 - 缺失索引
-- =====================================================

-- 模型+状态 组合查询 (模型回滚查询)
CREATE INDEX IF NOT EXISTS idx_ai_rollback_tasks_model_status 
ON ai_rollback_tasks(model_id, status);

-- 状态+创建时间 组合查询 (待处理回滚)
CREATE INDEX IF NOT EXISTS idx_ai_rollback_tasks_status_time 
ON ai_rollback_tasks(status, created_at);

-- =====================================================
-- 14. OTAProgress OTA升级进度 - 缺失索引
-- =====================================================

-- 部署ID+设备ID 组合查询 (特定设备进度)
CREATE INDEX IF NOT EXISTS idx_ota_progress_deployment_device 
ON ota_progress(deployment_id, device_id);

-- 状态+更新时间 组合查询 (状态+排序)
CREATE INDEX IF NOT EXISTS idx_ota_progress_status_time 
ON ota_progress(ota_status, updated_at);

-- =====================================================
-- 15. MemberPointsRecord 会员积分流水 - 缺失索引
-- =====================================================

-- 会员+创建时间 组合查询 (积分历史)
CREATE INDEX IF NOT EXISTS idx_member_points_records_member_time 
ON member_points_records(member_id, created_at);

-- 积分类型+创建时间 组合查询 (类型统计)
CREATE INDEX IF NOT EXISTS idx_member_points_records_type_time 
ON member_points_records(points_type, created_at);

-- =====================================================
-- 16. CouponGrant 优惠券发放记录 - 缺失索引
-- =====================================================

-- 优惠券+状态 组合查询 (优惠券使用情况)
CREATE INDEX IF NOT EXISTS idx_coupon_grants_coupon_status 
ON coupon_grants(coupon_id, status);

-- 会员+状态 组合查询 (会员优惠券)
CREATE INDEX IF NOT EXISTS idx_coupon_grants_member_status 
ON coupon_grants(member_id, status);

-- =====================================================
-- 17. MemberTagRecord 会员标签流水 - 缺失索引
-- =====================================================

-- 标签+创建时间 组合查询 (标签使用历史)
CREATE INDEX IF NOT EXISTS idx_member_tag_records_tag_time 
ON member_tag_records(tag_id, created_at);

-- 会员+标签 组合查询 (会员标签查询)
CREATE INDEX IF NOT EXISTS idx_member_tag_records_member_tag 
ON member_tag_records(member_id, tag_id);

-- =====================================================
-- 18. PetMessage 宠物对话消息 - 缺失索引
-- =====================================================

-- 会话ID+创建时间 组合查询 (消息历史)
CREATE INDEX IF NOT EXISTS idx_pet_messages_conversation_time 
ON pet_messages(conversation_id, created_at);

-- =====================================================
-- 19. BehaviorPrediction 行为预测 - 缺失索引
-- =====================================================

-- 宠物+预测类型+时间窗口 组合查询 (预测查询)
CREATE INDEX IF NOT EXISTS idx_behavior_predictions_pet_type_window 
ON behavior_predictions(pet_uuid, prediction_type, time_window_start);

-- =====================================================
-- 20. Notification 通知记录 - 缺失索引
-- =====================================================

-- 状态+创建时间 组合查询 (通知队列)
CREATE INDEX IF NOT EXISTS idx_notifications_status_time 
ON notifications(status, created_at);

-- =====================================================
-- 21. DeviceAlert 设备告警记录 - 缺失索引
-- =====================================================

-- 状态+创建时间 组合查询 (告警处理)
CREATE INDEX IF NOT EXISTS idx_device_alerts_status_time 
ON device_alerts(status, created_at);

-- =====================================================
-- 22. SubscriptionChange 订阅变更记录 - 缺失索引
-- =====================================================

-- 用户+创建时间 组合查询 (订阅历史)
CREATE INDEX IF NOT EXISTS idx_subscription_changes_user_time 
ON subscription_changes(user_id, created_at);

-- =====================================================
-- 23. UserSubscription 用户订阅 - 缺失索引
-- =====================================================

-- 计划+状态 组合查询 (计划订阅查询)
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_plan_status 
ON user_subscriptions(plan_id, status);

-- 过期时间+状态 组合查询 (过期订阅查询)
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_expire_status 
ON user_subscriptions(expire_time, status);

-- =====================================================
-- 24. PetEmotionAction 宠物情绪动作 - 缺失索引
-- =====================================================

-- 情绪类型+启用状态 组合查询 (情绪动作查询)
CREATE INDEX IF NOT EXISTS idx_pet_emotion_actions_emotion_enabled 
ON pet_emotion_actions(emotion_type, enabled);

-- =====================================================
-- 25. EmotionResponseConfig 情绪响应配置 - 缺失索引
-- =====================================================

-- 宠物+启用状态 组合查询 (宠物情绪配置)
CREATE INDEX IF NOT EXISTS idx_emotion_response_configs_pet_enabled 
ON emotion_response_configs(pet_id, enabled);

-- =====================================================
-- 26. AlertNotification 告警通知记录 - 缺失索引
-- =====================================================

-- 告警ID+状态 组合查询 (通知状态)
CREATE INDEX IF NOT EXISTS idx_alert_notifications_alert_status 
ON alert_notifications(alert_id, status);

-- =====================================================
-- 27. GeofenceAlert 地理围栏告警 - 缺失索引
-- =====================================================

-- 设备+状态 组合查询 (设备围栏告警)
CREATE INDEX IF NOT EXISTS idx_geofence_alerts_device_status 
ON geofence_alerts(device_id, status);

-- =====================================================
-- 28. HealthWarning 预警记录 - 补充索引
-- =====================================================

-- 宠物+状态+创建时间 组合查询 (宠物预警排序)
CREATE INDEX IF NOT EXISTS idx_health_warnings_pet_status_time 
ON health_warnings(pet_uuid, status, created_at);

-- =====================================================
-- 29. SensorEvent 传感器事件 - 缺失索引
-- =====================================================

-- 设备+传感器类型+时间 组合查询 (传感器历史)
CREATE INDEX IF NOT EXISTS idx_sensor_events_device_type_time 
ON sensor_events(device_id, sensor_type, created_at);

-- =====================================================
-- 30. LoginLog 登录日志 - 缺失索引
-- =====================================================

-- 租户+登录时间 组合查询 (登录日志查询)
CREATE INDEX IF NOT EXISTS idx_login_logs_tenant_time 
ON login_logs(tenant_id, login_time);

COMMIT;
