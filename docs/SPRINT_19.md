# Sprint 19 规划

**时间**：2026-08-09
**状态**：待开始
**Sprint 周期**：2 周（2026-08-09 ～ 2026-08-22）

---

## 一、Sprint 目标

**目标：** 健康追踪和预警

在 Sprint 18（数字孪生）的基础上，实现宠物健康追踪和早期疾病预警，包括早期疾病预警、运动追踪、睡眠分析，为宠物主人提供全面的健康管理服务。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **早期疾病预警 API** | 完成 `/api/v1/health/:pet_id/early-warning` 预警接口 | health_warning_controller.go | P0 |
| P0-2 | **疾病预警数据库** | 创建 health_warnings + disease_patterns 表 | models/health_warning.go | P0 |
| P0-3 | **运动追踪 API** | 完成 `/api/v1/health/:pet_id/exercise/*` 运动数据 | exercise_controller.go | P0 |
| P0-4 | **睡眠分析 API** | 完成 `/api/v1/health/:pet_id/sleep/*` 睡眠数据 | sleep_controller.go | P0 |
| P0-5 | **健康报告 API** | 完成 `/api/v1/health/:pet_id/report` 综合健康报告 | health_report_controller.go | P0 |
| P1-1 | **疾病模式库** | 完成 `/api/v1/health/disease-patterns` 疾病模式管理 | disease_pattern_controller.go | P1 |
| P1-2 | **健康基线计算** | 实现宠物个体健康基线自动计算 | health/baseline_service.go | P1 |
| P1-3 | **运动目标管理** | 完成 `/api/v1/health/:pet_id/exercise/goals` 目标设置 | exercise_goal_controller.go | P1 |
| P2-1 | **健康专家建议** | AI 生成个性化健康建议 | health/expert_suggestion_service.go | P2 |
| P2-2 | **健康数据分享** | 支持导出健康报告分享给兽医 | health/share_service.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **健康预警前端** | 完成 HealthWarningView.vue 早期疾病预警展示 | HealthWarningView.vue | P0 |
| PF0-2 | **运动统计前端** | 完成 ExerciseStatsView.vue 运动数据统计 | ExerciseStatsView.vue | P0 |
| PF0-3 | **睡眠分析前端** | 完成 SleepAnalysisView.vue 睡眠分析展示 | SleepAnalysisView.vue | P0 |
| PF0-4 | **综合健康报告** | 完成 HealthReportView.vue 健康报告查看 | HealthReportView.vue | P0 |
| PF1-1 | **健康趋势图表** | 完成 HealthTrendChart.vue 健康指标趋势 | HealthTrendChart.vue | P1 |
| PF1-2 | **运动目标设置** | 完成 ExerciseGoalSettingView.vue 运动目标配置 | ExerciseGoalSettingView.vue | P1 |
| PF2-1 | **健康数据分享** | 完成 HealthShareView.vue 分享给兽医 | HealthShareView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/health/:pet_id/early-warning` | GET | 早期疾病预警列表 |
| `GET /api/v1/health/:pet_id/early-warning/:id` | GET | 预警详情 |
| `POST /api/v1/health/:pet_id/early-warning/:id/ack` | POST | 确认预警 |
| `POST /api/v1/health/:pet_id/early-warning/:id/dismiss` | POST | 忽略预警 |
| `GET /api/v1/health/:pet_id/exercise` | GET | 运动数据列表 |
| `POST /api/v1/health/:pet_id/exercise` | POST | 上报运动数据 |
| `GET /api/v1/health/:pet_id/exercise/summary` | GET | 运动汇总统计 |
| `GET /api/v1/health/:pet_id/exercise/goals` | GET | 运动目标 |
| `POST /api/v1/health/:pet_id/exercise/goals` | POST | 设置运动目标 |
| `GET /api/v1/health/:pet_id/sleep` | GET | 睡眠数据列表 |
| `POST /api/v1/health/:pet_id/sleep` | POST | 上报睡眠数据 |
| `GET /api/v1/health/:pet_id/sleep/analysis` | GET | 睡眠分析报告 |
| `GET /api/v1/health/:pet_id/report` | GET | 综合健康报告 |
| `GET /api/v1/health/:pet_id/report/weekly` | GET | 周健康报告 |
| `GET /api/v1/health/:pet_id/report/monthly` | GET | 月健康报告 |
| `GET /api/v1/health/disease-patterns` | GET | 疾病模式列表 |
| `POST /api/v1/health/disease-patterns` | POST | 添加疾病模式 |
| `GET /api/v1/health/:pet_id/share` | POST | 分享健康数据 |

### 数据库设计

```sql
-- 疾病预警规则表
CREATE TABLE disease_patterns (
    id              BIGSERIAL PRIMARY KEY,
    pattern_name    VARCHAR(100) NOT NULL,
    disease_name    VARCHAR(100) NOT NULL,
    symptom_combinations JSONB NOT NULL,            -- 症状组合
    risk_factors    JSONB,                          -- 风险因子
    severity_level  INT DEFAULT 1,                   -- 1=低 2=中 3=高
    early_indicator JSONB,                          -- 早期指标
    recommendation  TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 健康预警记录表
CREATE TABLE health_warnings (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    pattern_id      BIGINT REFERENCES disease_patterns(id),
    warning_type    VARCHAR(50) NOT NULL,
    severity        INT NOT NULL,
    confidence      DECIMAL(5,4),
    trigger_data    JSONB,                          -- 触发数据
    description     TEXT,
    recommendation  TEXT,
    status          VARCHAR(20) DEFAULT 'active',  -- 'active'/'acknowledged'/'dismissed'/'resolved'
    detected_at     TIMESTAMP NOT NULL,
    acknowledged_at TIMESTAMP,
    acknowledged_by BIGINT REFERENCES users(id),
    resolved_at     TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_pet_status (pet_id, status),
    INDEX idx_detected_at (detected_at DESC)
);

-- 运动记录表
CREATE TABLE exercise_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    exercise_type   VARCHAR(50) NOT NULL,         -- 'walk'/'run'/'play'/'swim'
    duration_minutes INT NOT NULL,
    distance_km     DECIMAL(5,2),
    calories_burned INT,
    intensity       VARCHAR(20),                   -- 'low'/'medium'/'high'
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP,
    route_data      JSONB,                          -- 轨迹数据
    device_id       VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_pet_time (pet_id, start_time DESC)
);

-- 运动目标表
CREATE TABLE exercise_goals (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    goal_type       VARCHAR(50) NOT NULL,           -- 'daily_steps'/'daily_distance'/'weekly_exercise'
    target_value    INT NOT NULL,
    current_value   INT DEFAULT 0,
    period          VARCHAR(20) DEFAULT 'daily',  -- 'daily'/'weekly'
    start_date      DATE,
    is_active       BOOLEAN DEFAULT TRUE,
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(pet_id, goal_type)
);

-- 睡眠记录表
CREATE TABLE sleep_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    sleep_type      VARCHAR(20) NOT NULL,         -- 'deep'/'light'/'rem'
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP,
    duration_minutes INT,
    quality_score   INT,                          -- 1-100
    interruptions   INT DEFAULT 0,
    device_id       VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_pet_time (pet_id, start_time DESC)
);

-- 健康基线表
CREATE TABLE health_baselines (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    baseline_type   VARCHAR(50) NOT NULL,         -- 'heart_rate'/'sleep'/'exercise'/'weight'
    baseline_value  DECIMAL(10,2) NOT NULL,
    variance        DECIMAL(10,2),
    sample_count    INT DEFAULT 0,
    calculated_at   TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(pet_id, baseline_type)
);

-- 健康报告表
CREATE TABLE health_reports (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    report_type     VARCHAR(20) NOT NULL,         -- 'daily'/'weekly'/'monthly'
    period_start    TIMESTAMP NOT NULL,
    period_end      TIMESTAMP NOT NULL,
    summary         JSONB,
    exercise_summary JSONB,
    sleep_summary   JSONB,
    warning_summary JSONB,
    overall_score   INT,                           -- 1-100
    generated_at    TIMESTAMP DEFAULT NOW(),
    INDEX idx_pet_period (pet_id, period_start DESC)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 早期疾病预警 | 符合模式时准确触发预警 | 注入模拟数据 |
| 预警准确率 | 误报率 < 20% | 人工评估 |
| 运动追踪 | 运动数据正确记录 | 上报数据验证 |
| 运动目标达成 | 目标达成率计算正确 | 模拟数据验证 |
| 睡眠分析 | 深睡/浅睡时间正确 | 传感器数据对比 |
| 健康报告 | 汇总数据准确 | 人工核对 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 预警生成延迟 | <= 5 分钟（从数据到预警） |
| 报告生成 | <= 10s |
| 查询性能 | <= 500ms（1000 条内） |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 18 体征数据 | 运动/睡眠数据来源 |
| Sprint 14 AI 模型 | 疾病模式识别 |
| 宠物基础数据 | 品种/年龄/体重 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 预警误报率高 | 用户不信任 | 持续优化模式库 |
| 宠物数据不足 | 基线不准确 | 积累足够数据再启用 |
| 隐私问题 | 用户对健康数据敏感 | 数据脱敏+授权 |
