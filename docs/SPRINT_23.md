# Sprint 23 规划

**时间**：2026-09-13
**状态**：待开始
**Sprint 周期**：2 周（2026-09-13 ～ 2026-09-26）

---

## 一、Sprint 目标

**目标：** 仿真测试核心能力

实现完整的仿真与测试平台，包括虚拟宠物仿真运行环境、自动化测试框架、场景管理、回放系统，支持真机测试前的完整验证，大幅降低测试成本。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **虚拟宠物仿真 API** | 完成 `/api/v1/simulation/virtual-pet/*` 虚拟宠物管理 | virtual_pet_controller.go | P0 |
| P0-2 | **自动化测试框架 API** | 完成 `/api/v1/simulation/test-cases/*` 测试用例管理 | test_case_controller.go | P0 |
| P0-3 | **测试执行引擎** | 实现并发测试执行、WebSocket实时推送 | test_execution_engine.go | P0 |
| P0-4 | **回放系统 API** | 完成 `/api/v1/simulation/playback/*` 设备行为回放 | playback_controller.go | P0 |
| P1-1 | **仿真场景管理 API** | 完成 `/api/v1/simulation/scenarios/*` 场景管理 | scenario_controller.go | P1 |
| P1-2 | **测试报告生成** | 实现自动化测试报告生成 | test_report_generator.go | P1 |
| P2-1 | **用户行为模拟器** | 模拟真实用户交互行为 | user_behavior_simulator.go | P2 |
| P2-2 | **仿真数据集管理** | 完成 `/api/v1/simulation/datasets/*` 数据集管理 | dataset_controller.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **虚拟宠物运行面板** | 完成 VirtualPetRunView.vue 虚拟宠物运行控制 | VirtualPetRunView.vue | P0 |
| PF0-2 | **测试用例管理页面** | 完成 TestCaseManageView.vue 用例列表和执行 | TestCaseManageView.vue | P0 |
| PF0-3 | **测试执行监控页面** | 完成 TestExecutionView.vue 实时执行状态 | TestExecutionView.vue | P0 |
| PF0-4 | **回放系统页面** | 完成 PlaybackView.vue 设备行为回放 | PlaybackView.vue | P0 |
| PF1-1 | **场景编辑器页面** | 完成 ScenarioEditorView.vue 场景创建和配置 | ScenarioEditorView.vue | P1 |
| PF1-2 | **测试报告页面** | 完成 TestReportView.vue 测试报告查看 | TestReportView.vue | P1 |
| PF2-1 | **仿真数据集页面** | 完成 SimulationDatasetView.vue 数据集管理 | SimulationDatasetView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/simulation/virtual-pet` | GET | 获取虚拟宠物列表 |
| `POST /api/v1/simulation/virtual-pet` | POST | 创建虚拟宠物 |
| `GET /api/v1/simulation/virtual-pet/:id` | GET | 虚拟宠物详情 |
| `POST /api/v1/simulation/virtual-pet/:id/start` | POST | 启动虚拟宠物 |
| `POST /api/v1/simulation/virtual-pet/:id/stop` | POST | 停止虚拟宠物 |
| `GET /api/v1/simulation/test-cases` | GET | 测试用例列表 |
| `POST /api/v1/simulation/test-cases` | POST | 创建测试用例 |
| `GET /api/v1/simulation/test-cases/:id` | GET | 测试用例详情 |
| `PUT /api/v1/simulation/test-cases/:id` | PUT | 更新测试用例 |
| `DELETE /api/v1/simulation/test-cases/:id` | DELETE | 删除测试用例 |
| `POST /api/v1/simulation/test-cases/:id/run` | POST | 执行测试用例 |
| `POST /api/v1/simulation/test-cases/batch-run` | POST | 批量执行 |
| `GET /api/v1/simulation/executions/:id` | GET | 执行详情 |
| `GET /api/v1/simulation/playback` | GET | 回放列表 |
| `POST /api/v1/simulation/playback` | POST | 创建回放任务 |
| `GET /api/v1/simulation/playback/:id` | GET | 回放详情 |
| `GET /api/v1/simulation/scenarios` | GET | 场景列表 |
| `POST /api/v1/simulation/scenarios` | POST | 创建场景 |
| `GET /api/v1/simulation/scenarios/:id` | GET | 场景详情 |
| `PUT /api/v1/simulation/scenarios/:id` | PUT | 更新场景 |
| `GET /api/v1/simulation/reports/:execution_id` | GET | 测试报告 |
| `GET /api/v1/simulation/datasets` | GET | 数据集列表 |
| `POST /api/v1/simulation/datasets` | POST | 上传数据集 |

### 数据库设计

```sql
-- 虚拟宠物配置表
CREATE TABLE simulation_virtual_pets (
    id              BIGSERIAL PRIMARY KEY,
    pet_name        VARCHAR(100) NOT NULL,
    pet_type        VARCHAR(50) NOT NULL,
    config          JSONB NOT NULL,                   -- 虚拟宠物配置
    status          VARCHAR(20) DEFAULT 'stopped',   -- 'stopped'/'running'/'error'
    current_state   JSONB,
    runtime_hours   DECIMAL(10,2) DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 测试用例表
CREATE TABLE simulation_test_cases (
    id              BIGSERIAL PRIMARY KEY,
    case_name       VARCHAR(255) NOT NULL,
    case_type       VARCHAR(50) NOT NULL,            -- 'unit'/'integration'/'e2e'/'stress'
    description     TEXT,
    preconditions   JSONB,                           -- 前置条件
    test_steps      JSONB NOT NULL,                   -- 测试步骤
    expected_results JSONB,                          -- 预期结果
    test_data       JSONB,                           -- 测试数据
    tags            VARCHAR(50)[],
    priority        VARCHAR(10) DEFAULT 'P2',      -- 'P0'/'P1'/'P2'
    pass_count      INT DEFAULT 0,
    fail_count      INT DEFAULT 0,
    last_run_at     TIMESTAMP,
    last_run_result VARCHAR(20),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 测试执行记录表
CREATE TABLE simulation_test_executions (
    id              BIGSERIAL PRIMARY KEY,
    case_id         BIGINT NOT NULL REFERENCES simulation_test_cases(id),
    trigger_type    VARCHAR(30) NOT NULL,            -- 'manual'/'scheduled'/'cicd'
    triggered_by    VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'queued',    -- 'queued'/'running'/'passed'/'failed'/'cancelled'
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    duration_ms     INT,
    result_data     JSONB,                           -- 执行结果详情
    error_message   TEXT,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 回放记录表
CREATE TABLE simulation_playbacks (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100),
    replay_type     VARCHAR(50) NOT NULL,            -- 'device_behavior'/'user_interaction'
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP NOT NULL,
    playback_data   JSONB NOT NULL,                   -- 回放数据
    playback_speed  DECIMAL(3,2) DEFAULT 1.0,
    current_position TIMESTAMP,
    status          VARCHAR(20) DEFAULT 'ready',     -- 'ready'/'playing'/'paused'/'completed'
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 仿真场景表
CREATE TABLE simulation_scenarios (
    id              BIGSERIAL PRIMARY KEY,
    scenario_name   VARCHAR(255) NOT NULL,
    scenario_type   VARCHAR(50) NOT NULL,            -- 'functional'/'stress'/'regression'
    environment     JSONB,                           -- 环境配置
    pet_config      JSONB,                           -- 宠物配置
    device_config   JSONB,                           -- 设备配置
    description     TEXT,
    tags            VARCHAR(50)[],
    usage_count     INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 测试报告表
CREATE TABLE simulation_test_reports (
    id              BIGSERIAL PRIMARY KEY,
    execution_id    BIGINT NOT NULL REFERENCES simulation_test_executions(id),
    report_type     VARCHAR(50) NOT NULL,            -- 'summary'/'detailed'
    summary         JSONB,                           -- 摘要数据
    details         JSONB,                           -- 详细数据
    charts          JSONB,                           -- 图表数据
    generated_at    TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 虚拟宠物运行 | 虚拟宠物可在仿真环境中流畅运行>4小时 | 长时间运行测试 |
| 测试用例管理 | 用例CRUD正常，支持批量执行 | CRUD测试 |
| 测试执行 | 并发执行>10用例，WebSocket实时推送 | 性能测试 |
| 回放系统 | 设备行为完整回放，播放控制正常 | 回放功能测试 |
| 场景管理 | 场景创建/编辑/删除正常 | CRUD测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 单用例执行时间 | < 基准时间×1.2 |
| 并发执行10用例 | 总时间 < 单用例×3 |
| 回放启动时间 | < 2s |
| 测试报告生成 | < 10s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| MODULE_SIMULATION | 仿真测试完整 PRD |
| 仿真运行环境 | Docker容器化仿真环境 |
| CI/CD 集成 | GitHub Actions / Jenkins |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 仿真环境资源不足 | 测试性能差 | 资源监控+自动扩缩容 |
| 测试数据不足 | 覆盖率低 | 持续扩充数据集 |
| 仿真与真机差异 | 测试通过但真机失败 | 定期真机回归测试 |
