-- ============================================
-- AI Engineering 模块
-- ============================================

-- AI训练任务表
CREATE TABLE IF NOT EXISTS ai_training_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id VARCHAR(36) NOT NULL UNIQUE,
    task_name VARCHAR(200) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    dataset_id VARCHAR(36),
    hyperparameters JSONB,
    training_config JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    current_epoch INT DEFAULT 0,
    total_epochs INT,
    metrics JSONB,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    error_message TEXT,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AI数据集表
CREATE TABLE IF NOT EXISTS ai_datasets (
    id BIGSERIAL PRIMARY KEY,
    dataset_id VARCHAR(36) NOT NULL UNIQUE,
    dataset_name VARCHAR(200) NOT NULL,
    dataset_type VARCHAR(50) NOT NULL,
    description TEXT,
    data_format VARCHAR(20),
    data_size BIGINT,
    sample_count INT,
    features JSONB,
    labels JSONB,
    storage_path VARCHAR(500),
    version VARCHAR(20) DEFAULT '1.0',
    status VARCHAR(20) DEFAULT 'active',
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AI实验表
CREATE TABLE IF NOT EXISTS ai_experiments (
    id BIGSERIAL PRIMARY KEY,
    experiment_id VARCHAR(36) NOT NULL UNIQUE,
    experiment_name VARCHAR(200) NOT NULL,
    description TEXT,
    hypothesis TEXT,
    config JSONB,
    metrics JSONB,
    status VARCHAR(20) DEFAULT 'draft',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AI实验组表
CREATE TABLE IF NOT EXISTS ai_experiment_groups (
    id BIGSERIAL PRIMARY KEY,
    group_id VARCHAR(36) NOT NULL UNIQUE,
    group_name VARCHAR(200) NOT NULL,
    experiment_id VARCHAR(36),
    variant_name VARCHAR(50) NOT NULL,
    config JSONB,
    metrics JSONB,
    status VARCHAR(20) DEFAULT 'running',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AI监控指标表
CREATE TABLE IF NOT EXISTS ai_monitoring_metrics (
    id BIGSERIAL PRIMARY KEY,
    metric_id VARCHAR(36) NOT NULL UNIQUE,
    model_id VARCHAR(64),
    model_version VARCHAR(50),
    metric_type VARCHAR(50) NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(10,4) NOT NULL,
    unit VARCHAR(20),
    threshold DECIMAL(10,4),
    dimensions JSONB,
    collected_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI告警规则表
CREATE TABLE IF NOT EXISTS ai_alert_rules (
    id BIGSERIAL PRIMARY KEY,
    rule_id VARCHAR(36) NOT NULL UNIQUE,
    rule_name VARCHAR(200) NOT NULL,
    alert_type VARCHAR(50) NOT NULL,
    condition_expr TEXT NOT NULL,
    threshold_value VARCHAR(100),
    severity INT DEFAULT 2,
    notification_channels VARCHAR(50)[],
    is_active BOOLEAN DEFAULT true,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AI沙箱环境表
CREATE TABLE IF NOT EXISTS ai_sandbox_environments (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(36) NOT NULL UNIQUE,
    env_name VARCHAR(100) NOT NULL,
    env_type VARCHAR(50) NOT NULL,
    config JSONB,
    resources JSONB,
    status VARCHAR(20) DEFAULT 'inactive',
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AI沙箱测试用例表
CREATE TABLE IF NOT EXISTS ai_sandbox_testcases (
    id BIGSERIAL PRIMARY KEY,
    case_id VARCHAR(36) NOT NULL UNIQUE,
    env_id VARCHAR(36),
    test_name VARCHAR(200) NOT NULL,
    test_type VARCHAR(50),
    input_data JSONB,
    expected_output JSONB,
    actual_output JSONB,
    test_result VARCHAR(20),
    execution_time_ms INT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI决策日志表
CREATE TABLE IF NOT EXISTS ai_decision_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    session_id VARCHAR(64),
    model_id VARCHAR(64),
    model_version VARCHAR(50),
    input_data JSONB,
    output_data JSONB,
    decision_reason TEXT,
    confidence_score FLOAT,
    execution_time_ms INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI路由策略表
CREATE TABLE IF NOT EXISTS ai_routing_policies (
    id BIGSERIAL PRIMARY KEY,
    policy_id VARCHAR(36) NOT NULL UNIQUE,
    policy_name VARCHAR(200) NOT NULL,
    policy_type VARCHAR(50) NOT NULL,
    conditions JSONB NOT NULL,
    actions JSONB NOT NULL,
    priority INT DEFAULT 50,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Embodied AI 模块
-- ============================================

-- 环境地图表
CREATE TABLE IF NOT EXISTS embodied_maps (
    id BIGSERIAL PRIMARY KEY,
    map_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64),
    map_name VARCHAR(100) NOT NULL,
    map_data JSONB NOT NULL,
    resolution DECIMAL(6,2),
    coverage_area JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 空间位置表
CREATE TABLE IF NOT EXISTS spatial_positions (
    id BIGSERIAL PRIMARY KEY,
    position_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    map_id VARCHAR(36),
    x DECIMAL(10,4) NOT NULL,
    y DECIMAL(10,4) NOT NULL,
    z DECIMAL(10,4) DEFAULT 0,
    orientation DECIMAL(6,2),
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 动作执行记录表
CREATE TABLE IF NOT EXISTS action_executions (
    id BIGSERIAL PRIMARY KEY,
    execution_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    action_id VARCHAR(64) NOT NULL,
    action_name VARCHAR(100),
    parameters JSONB,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_ms INT,
    status VARCHAR(20),
    result_data JSONB,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 安全区域表
CREATE TABLE IF NOT EXISTS safety_zones (
    id BIGSERIAL PRIMARY KEY,
    zone_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    zone_name VARCHAR(100) NOT NULL,
    zone_type VARCHAR(50),
    boundary JSONB NOT NULL,
    center_lat DECIMAL(10,6),
    center_lng DECIMAL(10,6),
    radius DECIMAL(10,4),
    is_active BOOLEAN DEFAULT true,
    alert_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 具身决策日志表
CREATE TABLE IF NOT EXISTS embodied_decision_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    decision_type VARCHAR(50) NOT NULL,
    context_data JSONB,
    decision_result JSONB,
    confidence_score FLOAT,
    execution_path JSONB,
    latency_ms INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 安全日志表
CREATE TABLE IF NOT EXISTS safety_logs (
    id BIGSERIAL PRIMARY KEY,
    log_id VARCHAR(36) NOT NULL UNIQUE,
    device_id VARCHAR(64) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20),
    description TEXT,
    location JSONB,
    triggered_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP,
    resolution TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
