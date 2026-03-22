# Sprint 14 规划

**时间**：2026-05-31
**状态**：进行中
**Sprint 周期**：2 周（2026-05-31 ～ 2026-06-13）

---

## 一、Sprint 目标

**目标：** AI 行为监控和模型管理

在 Sprint 13（全球化）的基础上，实现 AI 行为监控、模型热回滚、AI 沙箱测试功能，提供完整的 AI 工程化能力，确保 AI 系统的稳定性和可观测性。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **AI 行为监控 API** | 完成 `/api/v1/ai/monitor/*` 行为数据采集接口 | ai_monitor_controller.go | P0 |
| P0-2 | **AI 行为日志存储** | 创建 ai_behavior_logs 表存储 AI 行为 | models/ai_behavior_log.go | P0 |
| P0-3 | **模型热回滚 API** | 完成 `POST /api/v1/ai/models/{id}/rollback` | model_rollback_controller.go | P0 |
| P0-4 | **AI 沙箱测试 API** | 完成 `/api/v1/ai/sandbox/*` 沙箱测试接口 | ai_sandbox_controller.go | P0 |
| P0-5 | **AI 质量指标 API** | 完成 `GET /api/v1/ai/quality/metrics` 质量指标 | ai_quality_controller.go | P0 |
| P1-1 | **AI 行为实时分析** | 实现 AI 行为异常检测服务 | ai/behavior_analyzer.go | P1 |
| P1-2 | **模型版本管理** | 完成 `/api/v1/ai/models/*` 版本 CRUD | model_version_controller.go | P1 |
| P1-3 | **沙箱测试报告** | 沙箱测试结果生成详细报告 | sandbox/report_generator.go | P1 |
| P2-1 | **AI 行为告警** | AI 行为异常时自动触发告警 | ai/behavior_alert_service.go | P2 |
| P2-2 | **AI 质量趋势分析** | 完成质量趋势图表数据接口 | ai_quality_controller.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **AI 质量仪表盘** | 完成 AIQualityDashboardView.vue AI 质量总览 | AIQualityDashboardView.vue | P0 |
| PF0-2 | **AI 行为日志页面** | 完成 AIBehaviorLogView.vue 行为日志查看 | AIBehaviorLogView.vue | P0 |
| PF0-3 | **模型版本管理页面** | 完成 ModelVersionView.vue 版本列表/发布/回滚 | ModelVersionView.vue | P0 |
| PF0-4 | **AI 沙箱测试页面** | 完成 AISandboxView.vue 沙箱测试配置和结果 | AISandboxView.vue | P0 |
| PF1-1 | **AI 行为详情页** | 完成 AIBehaviorDetailView.vue 单次行为详情 | AIBehaviorDetailView.vue | P1 |
| PF1-2 | **模型发布工作流** | 完成 ModelPublishWorkflow.vue 版本发布审批流 | ModelPublishWorkflow.vue | P1 |
| PF2-1 | **AI 质量趋势图表** | 完成 AIQualityTrendChart.vue 质量趋势可视化 | AIQualityTrendChart.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/v1/ai/monitor/events` | POST | 上报 AI 行为事件 |
| `GET /api/v1/ai/monitor/events` | GET | AI 行为事件列表 |
| `GET /api/v1/ai/monitor/events/:id` | GET | 行为事件详情 |
| `GET /api/v1/ai/quality/metrics` | GET | AI 质量指标（延迟/准确率/错误率） |
| `GET /api/v1/ai/models` | GET | 模型版本列表 |
| `POST /api/v1/ai/models` | POST | 注册新模型版本 |
| `GET /api/v1/ai/models/:id` | GET | 模型版本详情 |
| `PUT /api/v1/ai/models/:id` | PUT | 更新模型版本信息 |
| `POST /api/v1/ai/models/:id/rollback` | POST | 回滚到指定版本 |
| `GET /api/v1/ai/models/:id/rollback/history` | GET | 回滚历史 |
| `GET /api/v1/ai/sandbox/test` | GET | 沙箱测试列表 |
| `POST /api/v1/ai/sandbox/test` | POST | 创建沙箱测试任务 |
| `GET /api/v1/ai/sandbox/test/:id` | GET | 获取测试结果 |
| `POST /api/v1/ai/sandbox/test/:id/run` | POST | 执行沙箱测试 |

### 数据库设计

```sql
-- AI 行为日志表
CREATE TABLE ai_behavior_logs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(64),
    pet_id          BIGINT,
    user_id         BIGINT,
    model_version   VARCHAR(50) NOT NULL,
    behavior_type   VARCHAR(50) NOT NULL,         -- 'intent_recognition'/'response_generation'/'action_selection'
    input_summary   JSONB,                          -- 输入摘要
    output_summary  JSONB,                          -- 输出摘要
    latency_ms      INT,
    confidence      DECIMAL(5,4),
    error_code      VARCHAR(50),
    error_message   TEXT,
    is_anomaly      BOOLEAN DEFAULT FALSE,
    anomaly_score   DECIMAL(5,4),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_id (device_id),
    INDEX idx_model_version (model_version),
    INDEX idx_behavior_type (behavior_type),
    INDEX idx_is_anomaly (is_anomaly),
    INDEX idx_created_at (created_at DESC)
);

-- AI 模型版本表（扩展 MODULE_OPENCLAW_VERSION）
CREATE TABLE ai_model_versions (
    id              BIGSERIAL PRIMARY KEY,
    model_id        VARCHAR(50) NOT NULL,           -- 'openclaw'/'pet-llm'/'emotion'
    version_id      VARCHAR(50) NOT NULL,           -- semver
    version_name    VARCHAR(100),
    model_type      VARCHAR(30) NOT NULL,         -- 'openai'/'anthropic'/'local'
    model_name      VARCHAR(100) NOT NULL,
    endpoint_url    VARCHAR(500),
    api_key_ref     VARCHAR(100),                   -- 引用 secrets 表
    capabilities    JSONB,
    config          JSONB,                          -- 模型参数配置
    status          VARCHAR(20) NOT NULL,         -- 'dev'/'testing'/'stable'/'deprecated'/'rolling_back'
    rollout_percent INT DEFAULT 0,                  -- 灰度发布比例
    metrics         JSONB,                          -- 运行时指标
    test_report_url VARCHAR(500),
    published_at    TIMESTAMP,
    deprecated_at   TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(model_id, version_id)
);

-- AI 沙箱测试任务表
CREATE TABLE ai_sandbox_tests (
    id              BIGSERIAL PRIMARY KEY,
    test_name       VARCHAR(100) NOT NULL,
    model_version_id VARCHAR(50),
    test_type       VARCHAR(30) NOT NULL,         -- 'unit'/'integration'/'performance'/'safety'
    test_cases      JSONB NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending', -- 'pending'/'running'/'completed'/'failed'
    results         JSONB,
    report_url      VARCHAR(500),
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 模型回滚记录表
CREATE TABLE model_rollback_records (
    id              BIGSERIAL PRIMARY KEY,
    model_id        VARCHAR(50) NOT NULL,
    from_version    VARCHAR(50) NOT NULL,
    to_version      VARCHAR(50) NOT NULL,
    rollback_reason TEXT,
    triggered_by    BIGINT NOT NULL,
    rollback_status VARCHAR(20) NOT NULL,         -- 'initiated'/'in_progress'/'completed'/'failed'
    affected_devices INT DEFAULT 0,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| AI 行为日志 | 每次 AI 推理行为正确记录 | 调用 AI 接口验证日志 |
| 行为异常检测 | 异常行为 5 分钟内标记 | 注入异常测试 |
| 模型热回滚 | 回滚操作 30s 内完成 | 执行回滚测试 |
| 沙箱测试 | 完整测试报告生成 | 运行沙箱测试 |
| 质量指标 | 指标计算正确 | 对比人工统计 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 行为日志写入延迟 | <= 50ms |
| 质量指标查询 | <= 200ms |
| 模型回滚时间 | <= 30s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 9 AI 层 | AI 行为触发来源 |
| AI 模型服务 | OpenAI/Claude 等 API |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| AI 服务不可用 | AI 功能不可用 | 降级到本地模型 |
| 回滚过程中请求失败 | 用户体验差 | 灰度回滚+备用节点 |
| 行为日志数据量大 | DB 压力 | 分表+定期归档 |

---

## 六、完成清单

### 前端完成情况

| # | 任务 | 交付物 | 状态 |
|---|------|--------|------|
| PF0-1 | AI 质量仪表盘 | `views/ai/AIQualityDashboardView.vue` | ✅ 完成 |
| PF0-2 | AI 行为日志页面 | `views/ai/AIBehaviorLogView.vue` | ✅ 完成 |
| PF0-3 | 模型版本管理页面 | `views/ai/ModelVersionView.vue` | ✅ 完成 |
| PF0-4 | AI 沙箱测试页面 | `views/ai/AISandboxView.vue` | ✅ 完成 |
| PF1-1 | AI 行为详情页 | `views/ai/AIBehaviorDetailView.vue` | ✅ 完成 |
| PF1-2 | 模型发布工作流 | `views/ai/ModelPublishWorkflow.vue` | ✅ 完成 |
| — | API 层 | `api/ai.ts` | ✅ 完成 |
| — | Composable | `composables/useAIQuality.ts` | ✅ 完成 |
| — | 路由配置 | `router/index.js` | ✅ 完成 |
| — | 图表依赖 | `echarts` | ✅ 已安装 |

### 后端完成情况

| # | 任务 | 状态 |
|---|------|------|
| P0-1 | AI 行为监控 API | ⏳ 待开发 |
| P0-2 | AI 行为日志存储 | ⏳ 待开发 |
| P0-3 | 模型热回滚 API | ⏳ 待开发 |
| P0-4 | AI 沙箱测试 API | ⏳ 待开发 |
| P0-5 | AI 质量指标 API | ⏳ 待开发 |
