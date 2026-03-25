# 产品数据库设计文档 (PRODUCT_DATABASE_SCHEMA)

**版本**: V1.0  
**日期**: 2026-03-24  
**状态**: ✅ 完整

---

## 一、文档说明

本文档记录了MDM控制中台完整数据库设计，包括：
- **PRD定义表**: 192张 - Sprint 1-32及模块文档中定义的核心表
- **实现扩展表**: 127张 - 实现过程中必需的扩展表

---

## 二、核心系统表 (System Core)

### 2.1 用户与认证

| 表名 | 列数 | 说明 |
|------|------|------|
| users | 11 | 用户主表 |
| user_roles | 3 | 用户角色关联 |
| user_subscriptions | 11 | 用户订阅 |
| login_logs | 12 | 登录日志 |
| sys_login_logs | 11 | 系统登录日志 |

```sql
-- users 用户表
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    user_id         VARCHAR(36) UNIQUE,
    username        VARCHAR(64) NOT NULL UNIQUE,
    email           VARCHAR(100) NOT NULL UNIQUE,
    phone           VARCHAR(20),
    password_hash   VARCHAR(255),
    real_name       VARCHAR(100),
    avatar_url      VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'active',
    last_login_at   TIMESTAMP,
    tenant_id       UUID,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- user_roles 用户角色关联
CREATE TABLE user_roles (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL,
    role_id     BIGINT NOT NULL,
    tenant_id   UUID,
    created_at  TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, role_id)
);
```

### 2.2 权限与角色

| 表名 | 列数 | 说明 |
|------|------|------|
| roles | 8 | 角色表 |
| permissions | 15 | 权限表 |
| role_permissions | 4 | 角色权限关联 |
| menus | 14 | 菜单表 |
| sys_menus | 12 | 系统菜单 |

```sql
-- roles 角色表
CREATE TABLE roles (
    id          BIGSERIAL PRIMARY KEY,
    role_id     VARCHAR(36) UNIQUE,
    role_name   VARCHAR(100) NOT NULL,
    role_code   VARCHAR(50) UNIQUE,
    description TEXT,
    status      VARCHAR(20) DEFAULT 'active',
    tenant_id   UUID,
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- permissions 权限表
CREATE TABLE permissions (
    id              BIGSERIAL PRIMARY KEY,
    permission_id    VARCHAR(36) UNIQUE,
    permission_name  VARCHAR(100) NOT NULL,
    permission_code  VARCHAR(100) NOT NULL UNIQUE,
    resource_type   VARCHAR(50),
    description     TEXT,
    parent_id       BIGINT,
    sort_order      INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 2.3 租户管理

| 表名 | 列数 | 说明 |
|------|------|------|
| tenants | 14 | 租户表 |
| tenant_configs | 10 | 租户配置 |
| tenant_quotas | 7 | 租户配额 |
| tenant_applications | 18 | 租户申请 |
| departments | 11 | 部门表 |

```sql
-- tenants 租户表
CREATE TABLE tenants (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       UUID UNIQUE,
    tenant_name     VARCHAR(200) NOT NULL,
    tenant_type     VARCHAR(50),
    contact_name    VARCHAR(100),
    contact_email   VARCHAR(100),
    contact_phone   VARCHAR(20),
    status          VARCHAR(20) DEFAULT 'active',
    plan_type       VARCHAR(50),
    expires_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 三、设备管理 (Device Management)

### 3.1 设备主表

| 表名 | 列数 | 说明 |
|------|------|------|
| devices | 14 | 设备主表 |
| device_shadows | 21 | 设备影子 |
| device_alerts | 19 | 设备告警 |
| device_alert_rules | 17 | 设备告警规则 |

```sql
-- devices 设备表
CREATE TABLE devices (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL UNIQUE,
    device_name     VARCHAR(100),
    device_type     VARCHAR(50),
    model           VARCHAR(100),
    firmware_version VARCHAR(50),
    status          VARCHAR(20) DEFAULT 'offline',
    last_seen_at    TIMESTAMP,
    owner_id        BIGINT,
    tenant_id       UUID,
    location        JSONB,
    metadata        JSONB,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- device_shadows 设备影子
CREATE TABLE device_shadows (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64) NOT NULL,
    desired_state   JSONB,
    reported_state  JSONB,
    delta           JSONB,
    version         BIGINT,
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 3.2 设备运维

| 表名 | 列数 | 说明 |
|------|------|------|
| device_certificates | 17 | 设备证书 |
| device_lock_records | 15 | 设备锁定记录 |
| device_operation_logs | 11 | 设备操作日志 |
| device_performance_history | 11 | 设备性能历史 |

### 3.3 OTA升级

| 表名 | 列数 | 说明 |
|------|------|------|
| ota_tasks | 16 | OTA任务 |
| ota_packages | 29 | OTA包 |
| ota_deployments | 25 | OTA部署 |
| ota_progress | 15 | OTA进度 |
| firmware_optimization_configs | 12 | 固件优化配置 |

---

## 四、会员管理 (Member Management)

### 4.1 会员主表

| 表名 | 列数 | 说明 |
|------|------|------|
| members | 19 | 会员主表 |
| member_cards | 13 | 会员卡 |
| member_levels | 10 | 会员等级 |
| member_tags | 11 | 会员标签 |
| member_card_groups | 6 | 会员卡组 |

### 4.2 会员业务

| 表名 | 列数 | 说明 |
|------|------|------|
| member_orders | 11 | 会员订单 |
| member_points_records | 12 | 积分记录 |
| coupons | 15 | 优惠券 |
| coupon_grants | 8 | 优惠券发放 |
| promotions | 10 | 促销活动 |
| points_rules | 10 | 积分规则 |
| member_upgrade_rules | 9 | 会员升级规则 |
| member_upgrade_records | 7 | 会员升级记录 |
| member_operation_records | 7 | 会员操作记录 |
| member_tag_records | 6 | 会员标签记录 |

---

## 五、AI与智能 (AI Intelligence)

### 5.1 AI模型管理

| 表名 | 列数 | 说明 |
|------|------|------|
| ai_models | 23 | AI模型 |
| ai_model_versions | 11 | AI模型版本 |
| ai_model_versions_v2 | 22 | AI模型版本V2 |
| ai_model_shards | 13 | AI模型分片 |
| ai_model_deploy_history | 12 | AI模型部署历史 |

```sql
-- ai_models AI模型表
CREATE TABLE ai_models (
    id              BIGSERIAL PRIMARY KEY,
    model_id        VARCHAR(64) NOT NULL UNIQUE,
    model_name      VARCHAR(100) NOT NULL,
    model_type      VARCHAR(50),
    version         VARCHAR(20),
    framework       VARCHAR(50),
    status          VARCHAR(20) DEFAULT 'training',
    accuracy        DECIMAL(5,4),
    metrics         JSONB,
    config          JSONB,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 5.2 AI训练与推理

| 表名 | 列数 | 说明 |
|------|------|------|
| ai_training | 33 | AI训练 |
| ai_training_tasks | 18 | AI训练任务 |
| ai_datasets | 16 | AI数据集 |
| ai_experiments | 13 | AI实验 |
| ai_inference | 27 | AI推理 |

### 5.3 AI监控与质量

| 表名 | 列数 | 说明 |
|------|------|------|
| ai_monitoring_metrics | 12 | AI监控指标 |
| ai_alert_rules | 12 | AI告警规则 |
| ai_audit_logs | 24 | AI审计日志 |
| ai_audit_reports | 17 | AI审计报告 |
| ai_fairness_tests | 22 | AI公平性测试 |
| ai_fairness_metrics | 15 | AI公平性指标 |
| ai_bias_detections | 19 | AI偏见检测 |

### 5.4 AI决策与路由

| 表名 | 列数 | 说明 |
|------|------|------|
| ai_decision_logs | 11 | AI决策日志 |
| ai_routing_policies | 10 | AI路由策略 |
| ai_behavior_logs | 11 | AI行为日志 |

### 5.5 AI沙箱

| 表名 | 列数 | 说明 |
|------|------|------|
| ai_sandbox_environments | 10 | AI沙箱环境 |
| ai_sandbox_testcases | 12 | AI沙箱测试用例 |
| ai_sandbox_tests | 12 | AI沙箱测试 |
| ai_deploy_sharded | 15 | AI分片部署 |

---

## 六、具身智能 (Embodied AI)

| 表名 | 列数 | 说明 |
|------|------|------|
| embodied_perceptions | 10 | 具身感知 |
| embodied_maps | 10 | 环境地图 |
| embodied_action_library | 17 | 动作库 |
| embodied_collaboration_tasks | 13 | 多宠物协作任务 |
| embodied_collaboration_logs | 9 | 协作日志 |
| embodied_inference_configs | 11 | 推理配置 |
| embodied_ai_states | 30 | 具身AI状态 |
| embodied_decision_logs | 10 | 具身决策日志 |
| embodied_safety_logs | 11 | 具身安全日志 |
| embodied_safety_audits | 10 | 具身安全审计 |

```sql
-- embodied_perceptions 具身感知
CREATE TABLE embodied_perceptions (
    id              BIGSERIAL PRIMARY KEY,
    perception_id   VARCHAR(36) UNIQUE,
    device_id       VARCHAR(64) NOT NULL,
    perception_type VARCHAR(50) NOT NULL,
    raw_data        JSONB NOT NULL,
    processed_data  JSONB,
    confidence_score FLOAT,
    environment_context JSONB,
    timestamp       TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 七、宠物管理 (Pet Management)

### 7.1 宠物主表

| 表名 | 列数 | 说明 |
|------|------|------|
| pets | 18 | 宠物主表 |
| pet_profiles | 9 | 宠物档案 |
| pet_status | 18 | 宠物状态 |
| pet_device_bindings | 8 | 宠物设备绑定 |

### 7.2 宠物健康

| 表名 | 列数 | 说明 |
|------|------|------|
| pet_health_records | 25 | 宠物健康记录 |
| pet_vaccinations | 20 | 宠物疫苗 |
| pet_diet_records | 23 | 宠物饮食记录 |
| pet_lost_reports | 24 | 宠物丢失报告 |
| pet_finder_alerts | 8 | 寻宠告警 |
| pet_finder_reports | 21 | 寻宠报告 |
| pet_finder_sightings | 14 | 宠物目击 |

### 7.3 宠物社交

| 表名 | 列数 | 说明 |
|------|------|------|
| pet_social_posts | 10 | 宠物动态 |
| pet_social_comments | 6 | 动态评论 |
| pet_social_likes | 4 | 动态点赞 |
| pet_social_follows | 4 | 关注 |
| pet_social_interactions | 8 | 宠物互动 |

### 7.4 宠物医疗与保险

| 表名 | 列数 | 说明 |
|------|------|------|
| pet_medical_records | 15 | 宠物医疗记录 |
| pet_medical_configs | 9 | 医疗配置 |
| insurance_policies | 13 | 保险单 |
| insurance_claims | 23 | 保险理赔 |
| insurance_products | 27 | 保险产品 |

---

## 八、数字孪生 (Digital Twin)

| 表名 | 列数 | 说明 |
|------|------|------|
| digital_twin_pets | 8 | 数字孪生宠物 |
| vital_records | 22 | 生命体征记录 |
| vital_trends | 11 | 体征趋势 |
| behavior_events | 9 | 行为事件 |
| health_alerts | 33 | 健康告警 |
| health_warnings | 33 | 健康预警 |
| health_baselines | 11 | 健康基线 |
| health_reports | 12 | 健康报告 |
| highlight_moments | 10 | 精彩瞬间 |
| sync_records | 10 | 同步记录 |

---

## 九、健康与运动 (Health & Exercise)

| 表名 | 列数 | 说明 |
|------|------|------|
| exercise_goals | 17 | 运动目标 |
| exercise_records | 34 | 运动记录 |
| exercise_summaries | 23 | 运动摘要 |
| exercise_trends | 5 | 运动趋势 |
| sleep_records | 10 | 睡眠记录 |
| sleep_analysis | 12 | 睡眠分析 |
| diet_meals | 12 | 饮食记录 |
| nutrition_reports | 12 | 营养报告 |
| health_monitor_settings | 24 | 健康监控设置 |
| health_alert_rules | 17 | 健康告警规则 |

---

## 十、情感计算 (Affective Computing)

| 表名 | 列数 | 说明 |
|------|------|------|
| emotion_records | 13 | 情感记录 |
| emotion_reports | 12 | 情感报告 |
| emotion_response_configs | 13 | 情感响应配置 |
| pet_emotion_actions | 10 | 宠物情感行为 |
| family_emotions | 10 | 家庭情感 |

---

## 十一、通知与告警 (Notification & Alert)

### 11.1 通知渠道

| 表名 | 列数 | 说明 |
|------|------|------|
| notification_channels | 25 | 通知渠道 |
| notification_templates | 12 | 通知模板 |
| notification_logs | 14 | 通知日志 |
| notifications | 13 | 通知记录 |

### 11.2 告警管理

| 表名 | 列数 | 说明 |
|------|------|------|
| alert_templates | 12 | 告警模板 |
| alert_history | 20 | 告警历史 |
| alert_notifications | 10 | 告警通知 |
| alert_settings | 19 | 告警设置 |

---

## 十二、订阅与计费 (Subscription & Billing)

| 表名 | 列数 | 说明 |
|------|------|------|
| subscription_plans | 14 | 订阅计划 |
| subscription_changes | 13 | 订阅变更 |
| usage_limits | 18 | 使用限制 |
| usage_records | 11 | 使用记录 |
| billing_statements | 14 | 账单 |
| invoices | 16 | 发票 |
| daas_billings | 21 | DaaS计费 |
| daas_contracts | 22 | DaaS合同 |
| daas_device_rentals | 18 | DaaS设备租赁 |

---

## 十三、全球化 (Globalization)

| 表名 | 列数 | 说明 |
|------|------|------|
| regions | 12 | 区域 |
| region_nodes | 11 | 区域节点 |
| data_residency_rules | 10 | 数据驻留规则 |
| timezone_configs | 8 | 时区配置 |
| region_sync_records | 11 | 区域同步记录 |
| translations | 10 | 翻译 |

---

## 十四、仿真测试 (Simulation)

| 表名 | 列数 | 说明 |
|------|------|------|
| simulation_pets | 8 | 仿真宠物 |
| simulation_virtual_pets | 9 | 虚拟宠物 |
| simulation_environments | 14 | 仿真环境 |
| simulation_scenarios | 13 | 仿真场景 |
| simulation_scenes | 13 | 仿真场景 |
| simulation_test_cases | 14 | 测试用例 |
| simulation_testcases | 14 | 测试用例 |
| simulation_test_executions | 14 | 测试执行 |
| simulation_test_reports | 11 | 测试报告 |
| simulation_runs | 19 | 仿真运行 |
| simulation_metrics | 15 | 仿真指标 |
| simulation_playbacks | 9 | 回放记录 |
| simulation_stress_tests | 16 | 压力测试 |
| simulation_ab_experiments | 13 | AB实验 |
| simulation_cicd_configs | 9 | CI/CD配置 |
| simulation_cicd_jobs | 17 | CI/CD任务 |
| simulation_datasets | 14 | 仿真数据集 |
| simulation_dataset_versions | 9 | 数据集版本 |
| simulation_integrations | 10 | 仿真集成 |
| simulation_quotas | 10 | 仿真配额 |
| simulation_notifications | 10 | 仿真通知 |
| stress_tests | 15 | 压力测试 |
| test_executions | 13 | 测试执行 |
| test_reports | 11 | 测试报告 |
| playback_records | 8 | 回放记录 |
| dataset_jobs | 13 | 数据集任务 |
| dataset_samples | 9 | 数据集样本 |
| experiment_env_templates | 10 | 实验环境模板 |

---

## 十五、开放平台 (Open Platform)

### 15.1 开发者

| 表名 | 列数 | 说明 |
|------|------|------|
| developers | 11 | 开发者 |
| developer_apps | 16 | 开发者应用 |
| developer_api_keys | 11 | API密钥 |
| developer_api_usage | 11 | API使用记录 |

### 15.2 应用市场

| 表名 | 列数 | 说明 |
|------|------|------|
| apps | 12 | 应用 |
| store_apps | 26 | 商店应用 |
| store_app_versions | 17 | 应用版本 |
| store_installations | 15 | 应用安装 |
| store_reviews | 11 | 应用评价 |
| app_distributions | 18 | 应用分发 |
| app_install_records | 10 | 安装记录 |
| app_licenses | 11 | 应用许可 |

### 15.3 内容市场

| 表名 | 列数 | 说明 |
|------|------|------|
| market_apps | 16 | 市场应用 |
| market_emojis | 10 | 表情包 |
| market_actions | 11 | 动作包 |
| market_voices | 10 | 声音包 |
| emotion_packs | 14 | 表情包 |
| actions | 16 | 动作资源 |
| voices | 14 | 声音资源 |
| plugins | 17 | 插件 |
| ratings | 12 | 评分 |
| user_purchases | 12 | 用户购买 |

---

## 十六、Webhook管理

| 表名 | 列数 | 说明 |
|------|------|------|
| webhooks | 13 | Webhook |
| webhook_logs | 12 | Webhook日志 |
| webhook_market_templates | 14 | Webhook市场模板 |
| webhook_market_configs | 11 | Webhook市场配置 |
| webhook_market_logs | 12 | Webhook市场日志 |

---

## 十七、家庭场景 (Family)

| 表名 | 列数 | 说明 |
|------|------|------|
| households | 10 | 家庭 |
| household_members | 10 | 家庭成员 |
| household_invites | 12 | 家庭邀请 |
| family_mode_configs | 9 | 家庭模式配置 |
| family_albums | 19 | 家庭相册 |
| album_photos | 16 | 相册照片 |
| album_ai_albums | 10 | AI相册 |
| album_categories | 8 | 相册分类 |
| album_family_members | 10 | 相册成员 |
| family_activities | 14 | 家庭活动 |
| family_emotions | 10 | 家庭情感 |
| interaction_records | 10 | 互动记录 |
| household_pet_invites | 11 | 家庭宠物邀请 |

### 老人与儿童

| 表名 | 列数 | 说明 |
|------|------|------|
| elderly_profiles | 22 | 老人档案 |
| elderly_care_configs | 23 | 老人关怀配置 |
| elderly_reminders | 22 | 老人提醒 |
| children_profiles | 15 | 儿童档案 |
| child_mode_configs | 21 | 儿童模式配置 |
| virtual_pets | 19 | 虚拟宠物 |

---

## 十八、智能家居 (Smart Home)

| 表名 | 列数 | 说明 |
|------|------|------|
| integrations | 12 | 第三方集成 |
| smart_home_devices | 14 | 智能家居设备 |
| smart_home_triggers | 11 | 智能家居触发器 |
| smarthome_devices | 14 | 智能家居设备 |
| smarthome_configs | 8 | 智能家居配置 |
| smarthome_linkages | 11 | 智能家居联动 |

---

## 十九、宠物电商 (Pet E-commerce)

| 表名 | 列数 | 说明 |
|------|------|------|
| pet_products | 23 | 宠物商品 |
| pet_shop_orders | 23 | 宠物订单 |
| eco_recommended_products | 11 | 推荐产品 |
| eco_cart | 8 | 购物车 |
| eco_orders | 17 | 电商订单 |
| social_share_records | 13 | 社交分享记录 |

---

## 二十、API管理

| 表名 | 列数 | 说明 |
|------|------|------|
| api_keys | 15 | API密钥 |
| api_permissions | 10 | API权限 |
| api_quotas | 12 | API配额 |

---

## 二十一、安全与合规 (Security & Compliance)

| 表名 | 列数 | 说明 |
|------|------|------|
| encryption_keys | 12 | 加密密钥 |
| certificates | 17 | 证书 |
| ldap_configs | 20 | LDAP配置 |
| ldap_user_mappings | 12 | LDAP用户映射 |
| ldap_group_role_mappings | 7 | LDAP组角色映射 |
| compliance_policies | 12 | 合规策略 |
| compliance_violations | 12 | 合规违规 |
| gdpr_requests | 17 | GDPR请求 |
| data_permission_rules | 14 | 数据权限规则 |
| data_scopes | 7 | 数据范围 |
| data_anonymization_records | 13 | 数据脱敏记录 |
| content_filter_rules | 19 | 内容过滤规则 |
| device_lock_records | 15 | 设备锁定记录 |

---

## 二十二、运营与分析 (Operations & Analytics)

| 表名 | 列数 | 说明 |
|------|------|------|
| announcements | 16 | 公告 |
| activity_logs | 12 | 活动日志 |
| operation_logs | 15 | 操作日志 |
| audit_logs | 19 | 审计日志 |
| sys_operation_logs | 19 | 系统操作日志 |
| sys_configs | 8 | 系统配置 |
| sys_dictionaries | 10 | 系统字典 |
| analytics_records | 11 | 分析记录 |
| custom_reports | 16 | 自定义报表 |
| api_keys | 15 | API密钥 |

---

## 二十三、地理位置 (Geolocation)

| 表名 | 列数 | 说明 |
|------|------|------|
| geofence_rules | 13 | 地理围栏规则 |
| geofence_alerts | 11 | 地理围栏告警 |
| spatial_positions | 10 | 空间位置 |
| map_configs | 9 | 地图配置 |
| environment_maps | 27 | 环境地图 |
| finder_alerts | 15 | 查找器告警 |
| sighting_reports | 13 | 目击报告 |

---

## 二十四、会员扩展

| 表名 | 列数 | 说明 |
|------|------|------|
| plans | 14 | 计划 |
| temp_members | 8 | 临时会员 |
| lost_found_reports | 14 | 失物招领 |

---

## 二十五、RTOS系统

| 表名 | 列数 | 说明 |
|------|------|------|
| rtos_memory | 15 | RTOS内存 |
| rtos_stats | 15 | RTOS统计 |
| rtos_tasks | 15 | RTOS任务 |
| performance_metrics | 10 | 性能指标 |

---

## 二十六、研究平台

| 表名 | 列数 | 说明 |
|------|------|------|
| research_datasets | 18 | 研究数据集 |
| dataset_access_requests | 14 | 数据访问请求 |
| dataset_versions | 14 | 数据集版本 |
| dataset_download_logs | 10 | 数据集下载日志 |
| research_platforms | 11 | 研究平台 |
| research_experiments | 13 | 研究实验 |
| research_collaborators | 9 | 研究协作者 |
| research_experiment_results | 11 | 实验结果 |
| research_hypotheses | 8 | 研究假设 |
| research_institutions | 9 | 研究机构 |
| research_publications | 13 | 研究出版物 |
| research_discussions | 11 | 研究讨论 |
| ab_experiments | 14 | AB实验 |

---

## 二十七、其他表

| 表名 | 列数 | 说明 |
|------|------|------|
| knowledge | 7 | 知识库 |
| vaccination_reminders | 10 | 疫苗提醒 |
| mesh_devices | 18 | Mesh设备 |
| mesh_networks | 15 | Mesh网络 |
| mesh_network_members | 11 | Mesh网络成员 |
| order_items | 7 | 订单项 |
| product_categories | 10 | 产品分类 |
| policies | 16 | 策略 |
| policy_configs | 13 | 策略配置 |
| policy_bindings | 12 | 策略绑定 |
| position_templates | 12 | 岗位模板 |
| role_api_permissions | 4 | API权限 |
| role_menus | 4 | 角色菜单 |
| role_permission_groups | 4 | 角色权限组 |
| user_data_permissions | 14 | 用户数据权限 |
| tenant_quotas | 7 | 租户配额 |
| data_residency_rules | 10 | 数据驻留规则 |
| export_jobs | 17 | 导出任务 |
| exploration_sessions | 29 | 探索会话 |
| insurance_claim_documents | 16 | 保险理赔文档 |

---

## 二十八、数据库统计

| 类别 | 数量 |
|------|------|
| PRD定义表 | 192 张 |
| 实现扩展表 | 127 张 |
| **数据库总表数** | **319 张** |

### 表空间分布

| 模块 | 表数量 |
|------|--------|
| AI与智能 | ~35 |
| 仿真测试 | ~25 |
| 宠物管理 | ~20 |
| 会员管理 | ~15 |
| 设备管理 | ~12 |
| 具身智能 | ~10 |
| 数字孪生 | ~10 |
| 订阅计费 | ~9 |
| 开放平台 | ~20 |
| 安全合规 | ~13 |
| 其他 | ~50 |

---

_文档生成时间: 2026-03-24_
