# 模块 PRD：仿真与测试（MODULE_SIMULATION）

**版本：** V1.0
**所属Phase：** Phase 3（Sprint 23-24）
**优先级：** P1
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

仿真与测试模块为AI电子宠物提供完整的虚拟仿真环境和自动化测试能力，支持虚拟宠物运行、自动化测试用例管理、设备行为回放、场景仿真，以及A/B实验预验证，大幅降低真机测试成本，加速产品迭代。

### 1.2 核心价值

- **成本优化**：减少真机测试，缩短迭代周期
- **质量保障**：自动化回归测试，覆盖率持续提升
- **问题定位**：回放系统快速复现问题

### 1.3 范围边界

**包含：**
- 虚拟宠物仿真
- 自动化测试框架
- 回放系统
- 仿真场景管理
- 压力测试
- A/B实验仿真

**不包含：**
- 具身智能（MODULE_EMBODIED_AI）
- 情感计算（MODULE_AFFECTIVE_COMPUTING）
- CI/CD流水线（DevOps）

---

## 二、功能详情

### 2.1 虚拟宠物仿真

#### 2.1.1 虚拟宠物环境

| 功能 | 说明 |
|------|------|
| 虚拟宠物渲染 | 3D虚拟宠物形象展示 |
| 环境模拟 | 模拟不同家居环境 |
| 物理模拟 | 虚拟宠物的物理运动 |
| 多宠物支持 | 同时运行多个虚拟宠物 |
| 设备模拟 | 模拟设备传感器数据 |

#### 2.1.2 虚拟交互

| 功能 | 说明 |
|------|------|
| 触摸交互 | 模拟触摸虚拟宠物 |
| 语音交互 | 模拟语音指令和对话 |
| 手势交互 | 模拟手势指令 |
| 环境交互 | 虚拟宠物与虚拟环境交互 |

#### 2.1.3 仿真配置

| 功能 | 说明 |
|------|------|
| 宠物类型 | 支持多种宠物类型（猫/狗/兔子等） |
| 性格配置 | 配置虚拟宠物的性格 |
| 能力配置 | 配置虚拟宠物的AI能力 |
| 环境配置 | 配置虚拟环境参数 |

### 2.2 自动化测试框架

#### 2.2.1 测试用例管理

| 功能 | 说明 |
|------|------|
| 用例创建 | 创建测试用例 |
| 用例分类 | 按模块/功能/优先级分类 |
| 用例版本 | 用例版本管理 |
| 用例依赖 | 定义用例间依赖关系 |
| 用例标签 | 用例标签管理 |

#### 2.2.2 测试执行

| 功能 | 说明 |
|------|------|
| 单用例执行 | 执行单个测试用例 |
| 批量执行 | 批量执行多个用例 |
| 定时执行 | 配置定时执行任务 |
| 并发执行 | 支持多并发执行 |
| 失败重试 | 失败用例自动重试 |

#### 2.2.3 测试报告

| 功能 | 说明 |
|------|------|
| 实时报告 | 执行过程中实时展示结果 |
| 详细报告 | 包含步骤截图/日志 |
| 趋势分析 | 通过率趋势图 |
| 失败分析 | 失败原因分析 |
| 覆盖率统计 | 代码/功能覆盖率统计 |

#### 2.2.4 测试用例类型

| 类型 | 说明 | 示例 |
|------|------|------|
| 功能测试 | 验证功能正确性 | 情绪识别准确率测试 |
| 性能测试 | 验证性能指标 | 响应延迟测试 |
| 压力测试 | 验证极限能力 | 高并发测试 |
| 回归测试 | 验证修改未破坏功能 | 全量回归 |
| 冒烟测试 | 验证核心功能可用 | 快速冒烟 |

### 2.3 回放系统

#### 2.3.1 行为录制

| 功能 | 说明 |
|------|------|
| 自动录制 | 设备行为自动录制 |
| 手动标记 | 手动标记关键事件 |
| 传感器数据 | 录制完整传感器数据 |
| 用户操作 | 录制用户操作序列 |

#### 2.3.2 行为回放

| 功能 | 说明 |
|------|------|
| 精准回放 | 完整复现历史行为 |
| 变量替换 | 回放时替换变量 |
| 断点设置 | 回放中设置断点 |
| 速度控制 | 调整回放速度 |
| 部分回放 | 从指定时间点回放 |

#### 2.3.3 回放分析

| 功能 | 说明 |
|------|------|
| 差异对比 | 回放结果与录制对比 |
| 性能分析 | 回放过程性能分析 |
| 问题定位 | 快速定位问题原因 |
| 分享功能 | 分享回放给开发人员 |

### 2.4 仿真场景管理

#### 2.4.1 场景库

| 功能 | 说明 |
|------|------|
| 预置场景 | 预置常见测试场景 |
| 自定义场景 | 用户创建自定义场景 |
| 场景导入导出 | 场景的导入导出 |
| 场景评分 | 用户对场景评分 |

#### 2.4.2 场景编辑器

| 功能 | 说明 |
|------|------|
| 环境编辑 | 编辑虚拟环境 |
| 物体摆放 | 摆放物体和障碍物 |
| 灯光配置 | 配置光照条件 |
| 声音配置 | 配置背景声音 |
| 事件配置 | 配置触发事件 |

#### 2.4.3 预置场景

| 场景名称 | 说明 |
|----------|------|
| 客厅日常 | 模拟客厅日常活动 |
| 餐厅用餐 | 模拟用餐场景 |
| 夜间睡眠 | 模拟夜间睡眠 |
| 户外散步 | 模拟户外活动 |
| 多人互动 | 模拟多人家庭成员 |
| 紧急情况 | 模拟紧急事件 |

### 2.5 压力测试

#### 2.5.1 并发测试

| 功能 | 说明 |
|------|------|
| 设备并发 | 模拟多设备并发 |
| 用户并发 | 模拟多用户并发 |
| 请求并发 | 模拟API高并发 |
| 逐步加压 | 逐步增加并发量 |

#### 2.5.2 性能测试

| 功能 | 说明 |
|------|------|
| 响应时间 | 测试各接口响应时间 |
| 吞吐量 | 测试系统吞吐量 |
| 资源使用 | 监控CPU/内存/网络 |
| 性能瓶颈 | 定位性能瓶颈 |

#### 2.5.3 稳定性测试

| 功能 | 说明 |
|------|------|
| 长时间运行 | 7x24小时连续运行 |
| 内存泄漏 | 检测内存泄漏 |
| 连接泄漏 | 检测连接泄漏 |
| 故障恢复 | 测试故障自动恢复 |

### 2.6 A/B实验仿真

#### 2.6.1 实验预验证

| 功能 | 说明 |
|------|------|
| 参数仿真 | 仿真不同参数效果 |
| 流量模拟 | 模拟不同流量场景 |
| 效果预估 | 预估实验效果 |
| 风险评估 | 评估实验风险 |

#### 2.6.2 实验对比

| 功能 | 说明 |
|------|------|
| 多方案对比 | 同时对比多个方案 |
| 指标对比 | 对比各项指标 |
| 成本对比 | 对比资源消耗 |
| 推荐建议 | 推荐最优方案 |

---

## 三、API接口定义

### 3.1 虚拟宠物

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/simulation/pets | 创建虚拟宠物 |
| GET | /api/v1/simulation/pets/:id | 获取虚拟宠物 |
| PUT | /api/v1/simulation/pets/:id | 更新虚拟宠物配置 |
| DELETE | /api/v1/simulation/pets/:id | 删除虚拟宠物 |
| POST | /api/v1/simulation/pets/:id/interact | 虚拟交互 |

### 3.2 自动化测试

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/testcases | 测试用例列表 |
| POST | /api/v1/simulation/testcases | 创建测试用例 |
| GET | /api/v1/simulation/testcases/:id | 用例详情 |
| PUT | /api/v1/simulation/testcases/:id | 更新用例 |
| DELETE | /api/v1/simulation/testcases/:id | 删除用例 |
| POST | /api/v1/simulation/testcases/:id/execute | 执行用例 |
| POST | /api/v1/simulation/testcases/batch-execute | 批量执行 |
| GET | /api/v1/simulation/reports/:id | 测试报告 |
| GET | /api/v1/simulation/reports | 报告列表 |

### 3.3 回放系统

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/playbacks | 回放列表 |
| POST | /api/v1/simulation/playbacks | 创建回放（录制） |
| GET | /api/v1/simulation/playbacks/:id | 回放详情 |
| POST | /api/v1/simulation/playbacks/:id/play | 开始回放 |
| POST | /api/v1/simulation/playbacks/:id/stop | 停止回放 |
| POST | /api/v1/simulation/playbacks/:id/compare | 差异对比 |
| DELETE | /api/v1/simulation/playbacks/:id | 删除回放 |

### 3.4 仿真场景

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/simulation/scenarios | 场景列表 |
| POST | /api/v1/simulation/scenarios | 创建场景 |
| GET | /api/v1/simulation/scenarios/:id | 场景详情 |
| PUT | /api/v1/simulation/scenarios/:id | 更新场景 |
| DELETE | /api/v1/simulation/scenarios/:id | 删除场景 |
| POST | /api/v1/simulation/scenarios/:id/run | 运行场景 |
| POST | /api/v1/simulation/scenarios/import | 导入场景 |
| POST | /api/v1/simulation/scenarios/export/:id | 导出场景 |

### 3.5 压力测试

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/simulation/stress-tests | 创建压力测试 |
| GET | /api/v1/simulation/stress-tests/:id | 测试详情 |
| POST | /api/v1/simulation/stress-tests/:id/start | 开始测试 |
| POST | /api/v1/simulation/stress-tests/:id/stop | 停止测试 |
| GET | /api/v1/simulation/stress-tests/:id/report | 测试报告 |

---

## 四、数据库设计

### 4.1 虚拟宠物表 (simulation_pets)

```sql
CREATE TABLE simulation_pets (
    id              BIGSERIAL PRIMARY KEY,
    pet_name        VARCHAR(100) NOT NULL,
    pet_type        VARCHAR(50) NOT NULL,
    personality     JSONB,
    capabilities    JSONB,
    environment_id  BIGINT,
    status          VARCHAR(20) DEFAULT 'idle',
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.2 仿真场景表 (simulation_scenarios)

```sql
CREATE TABLE simulation_scenarios (
    id              BIGSERIAL PRIMARY KEY,
    scenario_name   VARCHAR(255) NOT NULL,
    scenario_type   VARCHAR(50) NOT NULL,         -- 'preset'/'custom'
    environment     JSONB,                          -- 环境配置
    objects         JSONB[],                        -- 物体列表
    events          JSONB[],                        -- 触发事件
    config          JSONB,                          -- 场景配置
    is_public       BOOLEAN DEFAULT FALSE,
    score           DECIMAL(5,2),
    downloads       INT DEFAULT 0,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.3 测试用例表 (simulation_testcases)

```sql
CREATE TABLE simulation_testcases (
    id              BIGSERIAL PRIMARY KEY,
    case_name       VARCHAR(255) NOT NULL,
    case_type       VARCHAR(30) NOT NULL,         -- 'functional'/'performance'/'stress'/'regression'/'smoke'
    module          VARCHAR(50),
    priority        VARCHAR(20) DEFAULT 'medium',  -- 'high'/'medium'/'low'
    description     TEXT,
    preconditions   TEXT,
    test_steps      JSONB,                          -- 测试步骤
    expected_result TEXT,
    tags            VARCHAR(100)[],
    dependencies    BIGINT[],
    version         INT DEFAULT 1,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.4 测试执行记录表 (test_executions)

```sql
CREATE TABLE test_executions (
    id              BIGSERIAL PRIMARY KEY,
    testcase_id     BIGINT REFERENCES simulation_testcases(id),
    execution_type  VARCHAR(20),                  -- 'manual'/'scheduled'/'api'
    trigger_params  JSONB,
    status          VARCHAR(20),                  -- 'pending'/'running'/'passed'/'failed'/'skipped'
    start_time      TIMESTAMP,
    end_time        TIMESTAMP,
    duration_ms     INT,
    environment     JSONB,
    result_details  JSONB,                        -- 详细结果
    screenshots     VARCHAR(500)[],
    logs            TEXT,
    error_message   TEXT,
    retry_count     INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_test_executions_case ON test_executions(testcase_id, created_at DESC);
```

### 4.5 测试报告表 (test_reports)

```sql
CREATE TABLE test_reports (
    id              BIGSERIAL PRIMARY KEY,
    report_name     VARCHAR(255),
    execution_ids   BIGINT[],
    summary         JSONB,                          -- 通过率/覆盖率等摘要
    pass_count      INT,
    fail_count      INT,
    skip_count      INT,
    total_count     INT,
    coverage        JSONB,                          -- 覆盖率数据
    trend_data      JSONB,                          -- 趋势数据
    generated_at    TIMESTAMP DEFAULT NOW()
);
```

### 4.6 回放记录表 (playback_records)

```sql
CREATE TABLE playback_records (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100),
    pet_id          BIGINT,
    record_type     VARCHAR(20) NOT NULL,         -- 'auto'/'manual'
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP,
    duration_ms     INT,
    sensor_data     JSONB,                          -- 传感器数据
    user_actions    JSONB[],                        -- 用户操作
    events          JSONB[],                        -- 事件
    playback_url    VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'recording',
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.7 压力测试表 (stress_tests)

```sql
CREATE TABLE stress_tests (
    id              BIGSERIAL PRIMARY KEY,
    test_name       VARCHAR(255) NOT NULL,
    test_type       VARCHAR(20) NOT NULL,         -- 'concurrent'/'performance'/'stability'
    config          JSONB NOT NULL,               -- 测试配置
    status          VARCHAR(20) DEFAULT 'draft',
    start_time      TIMESTAMP,
    end_time        TIMESTAMP,
    metrics         JSONB,                          -- 性能指标
    report_url      VARCHAR(500),
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 五、前端页面清单

### 5.1 虚拟宠物

| 页面 | 路由 | 说明 |
|------|------|------|
| 虚拟宠物列表 | /simulation/pets | 虚拟宠物列表 |
| 创建虚拟宠物 | /simulation/pets/create | 创建虚拟宠物 |
| 虚拟宠物详情 | /simulation/pets/:id | 虚拟宠物详情 |
| 虚拟交互 | /simulation/pets/:id/interact | 虚拟交互界面 |

### 5.2 测试用例

| 页面 | 路由 | 说明 |
|------|------|------|
| 用例列表 | /simulation/testcases | 测试用例列表 |
| 创建用例 | /simulation/testcases/create | 创建测试用例 |
| 用例详情 | /simulation/testcases/:id | 用例详情/编辑 |
| 执行结果 | /simulation/testcases/:id/results | 用例执行结果 |

### 5.3 测试报告

| 页面 | 路由 | 说明 |
|------|------|------|
| 报告列表 | /simulation/reports | 报告列表 |
| 报告详情 | /simulation/reports/:id | 报告详情 |
| 趋势分析 | /simulation/reports/trends | 通过率趋势 |

### 5.4 回放系统

| 页面 | 路由 | 说明 |
|------|------|------|
| 回放列表 | /simulation/playbacks | 回放列表 |
| 回放详情 | /simulation/playbacks/:id | 回放详情 |
| 回放播放器 | /simulation/playbacks/:id/player | 回放播放界面 |

### 5.5 仿真场景

| 页面 | 路由 | 说明 |
|------|------|------|
| 场景列表 | /simulation/scenarios | 场景库列表 |
| 创建场景 | /simulation/scenarios/create | 场景编辑器 |
| 场景详情 | /simulation/scenarios/:id | 场景详情 |
| 运行场景 | /simulation/scenarios/:id/run | 场景运行 |
| 场景导入 | /simulation/scenarios/import | 导入场景 |
| 场景导出 | /simulation/scenarios/export/:id | 导出场景 |

### 5.6 压力测试

| 页面 | 路由 | 说明 |
|------|------|------|
| 压力测试列表 | /simulation/stress-tests | 测试列表 |
| 创建测试 | /simulation/stress-tests/create | 创建压力测试 |
| 测试详情 | /simulation/stress-tests/:id | 测试详情/监控 |
| 测试报告 | /simulation/stress-tests/:id/report | 性能报告 |

---

## 六、验收标准

### 6.1 虚拟宠物仿真

| 验收点 | 标准 |
|--------|------|
| 虚拟宠物运行 | 虚拟宠物稳定运行不掉帧 |
| 交互响应 | 交互响应<200ms |
| 多宠物支持 | 同时运行5个虚拟宠物流畅 |

### 6.2 自动化测试

| 验收点 | 标准 |
|--------|------|
| 用例执行 | 用例按预期执行，结果准确 |
| 批量执行 | 100个用例并发执行正常 |
| 报告生成 | 执行完成5分钟内生成报告 |
| 覆盖率 | 核心功能覆盖率>80% |

### 6.3 回放系统

| 验收点 | 标准 |
|--------|------|
| 回放精度 | 回放与录制偏差<1% |
| 差异对比 | 差异对比准确 |
| 分享可用 | 回放可被其他用户播放 |

### 6.4 仿真场景

| 验收点 | 标准 |
|--------|------|
| 场景运行 | 预置场景100%可运行 |
| 场景编辑 | 自定义场景编辑保存成功 |
| 场景导入导出 | 导入导出成功率>99% |

### 6.5 压力测试

| 验收点 | 标准 |
|--------|------|
| 性能数据准确性 | 性能数据与实际一致 |
| 并发支持 | 支持100+并发连接 |
| 稳定性测试 | 7x24小时运行无内存泄漏 |


---

## 六、页面布局规范

### 6.1 虚拟宠物列表页面（/simulation/pets）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（创建虚拟宠物靠左）
3. 虚拟宠物卡片网格（展示宠物形象/类型/状态）

**按钮规范：**
- [创建虚拟宠物] — 左对齐
- [交互] [详情] [删除] — 卡片内右对齐

### 6.2 测试用例列表页面（/simulation/testcases）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：用例类型 / 模块 / 优先级
3. 操作栏（创建用例靠左，批量执行靠右）
4. 测试用例列表表格

**按钮规范：**
- [创建用例] — 左对齐
- [批量执行] — 操作栏中间
- [详情] [编辑] [执行] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 用例名称 | 200px | - |
| 类型 | 100px | 功能/性能/压力/回归/冒烟 |
| 优先级 | 100px | high/medium/low |
| 模块 | 120px | - |
| 状态 | 80px | - |
| 操作 | 120px | 详情/编辑/执行/删除 |

**分页：** 右下角，10/20/50/100 条

### 6.3 测试报告页面（/simulation/reports）

**布局结构：**
1. 面包屑 → 页面标题
2. 统计概览卡片（通过率/覆盖率/总用例数/失败数）—— 白色
3. 趋势图表
4. 测试报告列表表格

**按钮规范：**
- [查看详情] — 行内右对齐

**分页：** 右下角，10/20/50/100 条

### 6.4 回放列表页面（/simulation/playbacks）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：设备ID / 日期范围
3. 操作栏（开始录制靠左）
4. 回放列表表格

**按钮规范：**
- [开始录制] — 左对齐
- [播放] [对比] [删除] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 回放ID | 150px | - |
| 设备ID | 150px | - |
| 录制类型 | 100px | auto/manual |
| 时长 | 100px | - |
| 状态 | 100px | recording/completed |
| 创建时间 | 150px | - |
| 操作 | 120px | 播放/对比/删除 |

**分页：** 右下角，10/20/50/100 条

### 6.5 仿真场景页面（/simulation/scenarios）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：场景类型 / 是否公开
3. 操作栏（创建场景/导入靠左，导出靠右）
4. 场景列表卡片网格

**按钮规范：**
- [创建场景] [导入] — 左对齐
- [导出] — 操作栏中间
- [编辑] [运行] [删除] — 卡片内右对齐

### 6.6 压力测试页面（/simulation/stress-tests）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（创建测试靠左）
3. 压力测试列表表格

**按钮规范：**
- [创建测试] — 左对齐
- [开始] [停止] [查看报告] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 测试名称 | 200px | - |
| 测试类型 | 100px | 并发/性能/稳定性 |
| 状态 | 100px | draft/running/completed |
| 创建时间 | 150px | - |
| 操作 | 120px | 开始/停止/查看报告 |

**分页：** 右下角，10/20/50/100 条

### 6.7 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 创建/编辑测试用例、创建/编辑场景、测试报告详情 |
| Dialog 对话框 | 确认删除、确认开始/停止测试 |
| 全屏模态 | 暂无复杂表单场景 |
