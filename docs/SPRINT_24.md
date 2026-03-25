# Sprint 24 规划

**时间**：2026-09-20
**状态**：待开始
**Sprint 周期**：2 周（2026-09-20 ～ 2026-10-03）

---

## 一、Sprint 目标

**目标：** 仿真测试完善与CI/CD集成

在 Sprint 23 基础上，完善压力测试、A/B实验仿真、CI/CD集成，以及仿真资源配额管理，构建完整的仿真测试体系，实现与研发流水线的无缝对接。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **压力测试引擎** | 实现自研k6-style压力测试引擎 | stress_test_engine.go | P0 |
| P0-2 | **CI/CD集成 API** | 完成 `/api/v1/simulation/cicd/*` CI/CD集成配置 | cicd_controller.go | P0 |
| P0-3 | **A/B实验仿真 API** | 完成 `/api/v1/simulation/ab-experiments/*` A/B实验仿真 | ab_experiment_controller.go | P0 |
| P1-1 | **资源配额管理** | 实现仿真资源配额管理和计费 | quota_manager.go | P1 |
| P1-2 | **仿真结果分析** | 实现仿真结果与真机误差分析 | result_analyzer.go | P1 |
| P2-1 | **仿真告警通知** | 仿真异常自动通知到Slack/钉钉 | simulation_notifier.go | P2 |
| P2-2 | **仿真报告导出** | 支持PDF/HTML格式报告导出 | report_exporter.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **压力测试配置页面** | 完成 StressTestView.vue 压力测试配置和执行 | StressTestView.vue | P0 |
| PF0-2 | **CI/CD集成页面** | 完成 CICDIntegrationView.vue CI/CD流水线配置 | CICDIntegrationView.vue | P0 |
| PF0-3 | **A/B实验仿真页面** | 完成 ABExperimentSimulationView.vue A/B实验仿真配置 | ABExperimentSimulationView.vue | P0 |
| PF1-1 | **资源配额页面** | 完成 QuotaManagementView.vue 配额查看和管理 | QuotaManagementView.vue | P1 |
| PF1-2 | **仿真结果分析页面** | 完成 SimulationAnalysisView.vue 仿真与真机对比分析 | SimulationAnalysisView.vue | P1 |
| PF2-1 | **报告导出功能** | 完成 ReportExportView.vue 报告导出 | ReportExportView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/v1/simulation/stress-test` | POST | 创建压力测试任务 |
| `GET /api/v1/simulation/stress-test/:id` | GET | 压力测试详情 |
| `POST /api/v1/simulation/stress-test/:id/start` | POST | 启动压力测试 |
| `GET /api/v1/simulation/stress-test/:id/stop` | GET | 停止压力测试 |
| `GET /api/v1/simulation/stress-test/:id/results` | GET | 压力测试结果 |
| `GET /api/v1/simulation/cicd/pipelines` | GET | CI/CD流水线列表 |
| `POST /api/v1/simulation/cicd/pipelines` | POST | 创建流水线 |
| `PUT /api/v1/simulation/cicd/pipelines/:id` | PUT | 更新流水线 |
| `DELETE /api/v1/simulation/cicd/pipelines/:id` | DELETE | 删除流水线 |
| `POST /api/v1/simulation/cicd/pipelines/:id/trigger` | POST | 触发CI/CD |
| `GET /api/v1/simulation/cicd/webhook/notify` | GET | CI/CD Webhook回调 |
| `GET /api/v1/simulation/ab-experiments` | GET | A/B实验列表 |
| `POST /api/v1/simulation/ab-experiments` | POST | 创建A/B实验 |
| `GET /api/v1/simulation/ab-experiments/:id` | GET | 实验详情 |
| `POST /api/v1/simulation/ab-experiments/:id/run` | POST | 执行仿真 |
| `GET /api/v1/simulation/ab-experiments/:id/results` | GET | 实验对比结果 |
| `GET /api/v1/simulation/quotas` | GET | 配额使用情况 |
| `PUT /api/v1/simulation/quotas` | PUT | 更新配额配置 |
| `GET /api/v1/simulation/analysis/:execution_id` | GET | 仿真结果分析 |
| `GET /api/v1/simulation/reports/:id/export` | GET | 导出报告 |

### 数据库设计

```sql
-- 压力测试任务表
CREATE TABLE simulation_stress_tests (
    id              BIGSERIAL PRIMARY KEY,
    test_name       VARCHAR(255) NOT NULL,
    target_endpoint VARCHAR(500) NOT NULL,
    load_profile    JSONB NOT NULL,                   -- 负载配置：并发数/RPS/持续时间
    thresholds      JSONB,                           -- 性能阈值
    status          VARCHAR(20) DEFAULT 'draft',     -- 'draft'/'running'/'completed'/'failed'
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    results_summary JSONB,                           -- 摘要结果
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- CI/CD集成配置表
CREATE TABLE simulation_cicd_configs (
    id              BIGSERIAL PRIMARY KEY,
    pipeline_name   VARCHAR(255) NOT NULL,
    cicd_type       VARCHAR(30) NOT NULL,            -- 'github_actions'/'jenkins'/'gitlab_ci'
    webhook_url     VARCHAR(500),
    webhook_secret  VARCHAR(255),
    trigger_events  VARCHAR(50)[],                   -- 'push'/'pull_request'/'schedule'
    simulation_config JSONB,                          -- 仿真配置
    is_active       BOOLEAN DEFAULT TRUE,
    last_triggered_at TIMESTAMP,
    last_status     VARCHAR(20),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- A/B实验仿真表
CREATE TABLE simulation_ab_experiments (
    id              BIGSERIAL PRIMARY KEY,
    experiment_name VARCHAR(255) NOT NULL,
    hypothesis      TEXT,
    variant_a       JSONB NOT NULL,                   -- A组配置
    variant_b       JSONB NOT NULL,                   -- B组配置
    traffic_split   DECIMAL(5,2) DEFAULT 0.5,        -- 流量分配比例
    metrics         JSONB NOT NULL,                   -- 评估指标
    run_config      JSONB,                           -- 仿真运行配置
    status          VARCHAR(20) DEFAULT 'draft',     -- 'draft'/'running'/'completed'
    results         JSONB,                           -- 实验结果
    p_value         DECIMAL(6,4),                    -- 统计显著性
    conclusion      TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 仿真配额表
CREATE TABLE simulation_quotas (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    quota_type      VARCHAR(50) NOT NULL,            -- 'compute_hours'/'storage_gb'/'api_calls'
    total_quota     DECIMAL(10,2) NOT NULL,
    used_quota      DECIMAL(10,2) DEFAULT 0,
    reset_day       INT DEFAULT 1,                   -- 每月重置日
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, quota_type)
);

-- 仿真通知配置表
CREATE TABLE simulation_notifications (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    channel         VARCHAR(30) NOT NULL,            -- 'slack'/'dingtalk'/'email'
    webhook_url     VARCHAR(500),
    notify_on       VARCHAR(50)[],                   -- 'test_failed'/'quota_warning'/'report_ready'
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 压力测试 | 支持1000并发，压测报告准确 | 标准化压测场景 |
| CI/CD集成 | GitHub Actions触发自动仿真，报告回传 | 集成测试 |
| A/B实验仿真 | 仿真结果与真机误差<10% | 真机对比测试 |
| 配额管理 | 配额超限自动拦截，预警及时 | 配额测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 压力测试并发能力 | 支持1000并发 |
| CI/CD触发延迟 | < 30s |
| A/B结果分析时间 | < 60s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 23 仿真测试核心 | 基础测试框架 |
| CI/CD平台 | GitHub Actions / Jenkins API |
| 仿真计算资源 | GPU/CPU资源池 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 仿真资源被滥用 | 资源浪费 | 配额+计费机制 |
| CI/CD平台API变更 | 集成失效 | API版本管理 |
| 仿真与真机差异大 | 误导决策 | 定期真机对比验证 |
