# 模块 PRD：AI系统工程（MODULE_AI_ENGINEERING）

**版本：** V1.0
**所属Phase：** Phase 1（Sprint 5-6基础，Phase 2完善）
**优先级：** P1
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

AI系统工程模块是MDM控制中台的AI基础设施，负责管理宠物AI模型的训练、部署、监控、回滚全生命周期，支持A/B实验和边缘/云端智能路由，为具身智能和情感计算提供核心AI能力支撑。

### 1.2 核心价值

- **质量保障**：模型可观测、可回滚，线上问题秒级响应
- **效率提升**：训练流程自动化，A/B实验一键配置
- **成本优化**：边缘/云端智能路由，降低云端计算成本

### 1.3 范围边界

**包含：**
- AI模型训练流水线
- A/B测试框架
- 模型监控仪表盘
- 模型版本管理（含热回滚）
- AI沙箱测试环境
- AI决策日志与可解释性
- AI行为监控与异常告警
- 边缘AI vs 云端AI路由
- 模型分片加载

**不包含：**
- 具身智能核心（MODULE_EMBODIED_AI）
- 情感计算（MODULE_AFFECTIVE_COMPUTING）
- 端侧推理设备端实现（固件层）

---

## 二、功能详情

### 2.1 AI模型训练流水线

#### 2.1.1 功能描述

管理宠物AI模型的训练任务，支持数据准备、训练配置、状态追踪、结果评估全流程自动化。

#### 2.1.2 训练任务管理

| 功能 | 说明 |
|------|------|
| 创建训练任务 | 选择模型模板/数据集/超参数，提交训练任务 |
| 训练配置 | 学习率/批次大小/Epoch数/优化器配置 |
| 分布式训练 | 支持多GPU分布式训练 |
| 训练状态追踪 | 实时展示loss/accuracy/学习曲线 |
| 训练暂停/恢复 | 支持手动暂停和断点恢复 |
| 训练取消 | 支持取消正在运行的训练任务 |
| 训练结果评估 | 自动计算验证集指标，生成评估报告 |
| 模型导出 | 训练完成后自动导出为部署格式 |

#### 2.1.3 训练数据管理

| 功能 | 说明 |
|------|------|
| 数据集上传 | 支持图片/音频/视频/文本数据上传 |
| 数据集版本管理 | 每次数据变更生成新版本 |
| 数据集统计 | 数据量/类别分布/质量报告 |
| 数据标注工具 | 内置标注界面（图像分类/目标检测/语音标注） |
| 数据增强配置 | 翻转/旋转/颜色 jitter 等增强策略 |

#### 2.1.4 训练模板

| 模板名称 | 适用场景 | 预置配置 |
|----------|----------|----------|
| 视觉模型 | 宠物表情/行为识别 | ResNet50/EfficientNet |
| 语音模型 | 情绪识别/语音交互 | Transformer ASR |
| 多模态模型 | 跨模态理解 | CLIP-style |
| 强化学习 | 具身智能动作决策 | PPO/SAC |

### 2.2 A/B测试框架

#### 2.2.1 功能描述

支持AI模型的灰度发布和效果对比实验，支持多分组流量分配和统计分析。

#### 2.2.2 实验配置

| 功能 | 说明 |
|------|------|
| 创建实验 | 定义实验名称/描述/开始/结束时间 |
| 添加分组 | 控制组 + N个实验组 |
| 流量分配 | 支持百分比/用户ID哈希/设备ID哈希 |
| 流量分配规则 | 同用户始终路由到同一组（一致性） |
| 实验变量配置 | 每个分组可配置不同的模型版本/参数 |
| 实验目标设定 | 定义优化目标（点击率/停留时长/满意度） |
| 最小样本量 | 达到最小样本量才进行统计检验 |

#### 2.2.3 流量路由

| 功能 | 说明 |
|------|------|
| 自动路由 | 根据实验配置自动分流 |
| 强制分组 | 支持指定用户/设备进入特定分组 |
| 流量保护 | 单分组流量不超过配置上限 |
| 分组隔离 | 不同实验流量互不干扰 |

#### 2.2.4 效果评估

| 功能 | 说明 |
|------|------|
| 实时统计 | 各分组核心指标实时展示 |
| 统计显著性 | t检验/贝叶斯分析判断显著性 |
| 效果趋势图 | 指标随时间变化曲线 |
| 实验报告 | 自动生成实验结论和建议 |
| 胜出分组 | 自动标记胜出分组 |
| 一键全量 | 胜出分组一键全量发布 |

### 2.3 模型监控

#### 2.3.1 功能描述

实时监控线上AI模型的性能指标，及时发现异常并告警。

#### 2.3.2 监控指标

| 指标类型 | 具体指标 | 告警阈值 |
|----------|----------|----------|
| 性能指标 | 响应延迟（P50/P95/P99） | P99>500ms |
| 性能指标 | QPS/并发数 | 超过容量80% |
| 质量指标 | 准确率/召回率 | 低于基线5% |
| 质量指标 | 用户满意度 | <4.0分 |
| 业务指标 | 交互次数/会话数 | 突降>30% |
| 资源指标 | GPU/CPU利用率 | >90% |
| 资源指标 | 内存使用率 | >85% |

#### 2.3.3 监控面板

| 功能 | 说明 |
|------|------|
| 实时大盘 | 核心指标实时展示 |
| 历史趋势 | 指标历史数据查询 |
| 模型对比 | 多版本模型指标对比 |
| 告警记录 | 历史告警查询 |
| 仪表盘自定义 | 用户自定义指标卡片布局 |

### 2.4 模型版本管理

#### 2.4.1 功能描述

管理AI模型版本，支持版本列表、上线状态切换、热回滚。

#### 2.4.2 版本管理

| 功能 | 说明 |
|------|------|
| 版本列表 | 展示所有模型版本（训练任务关联） |
| 版本详情 | 版本号/训练数据/评估指标/上线时间 |
| 版本标签 | 支持给版本打标签（stable/beta/candidate） |
| 上线状态 | 线上/灰度/下线/草稿 |
| 上线审批 | 重要版本上线需审批 |
| 版本对比 | 两个版本指标对比 |

#### 2.4.3 模型热回滚

| 功能 | 说明 |
|------|------|
| 一键回滚 | 选择历史版本，一键切换 |
| 回滚时间 | 切换时间<30秒 |
| 回滚记录 | 记录回滚原因/操作人/时间 |
| 回滚确认 | 回滚后需确认观测一段时间 |
| 自动回滚 | 指标异常自动触发回滚 |

### 2.5 AI沙箱测试

#### 2.5.1 功能描述

提供隔离的测试环境，在正式发布前验证新模型/策略的效果。

#### 2.5.2 沙箱环境

| 功能 | 说明 |
|------|------|
| 测试环境隔离 | 独立资源，与生产完全隔离 |
| 测试数据 | 支持模拟数据/脱敏真实数据 |
| 环境配置 | 可配置不同的模型/参数 |
| 并发限制 | 沙箱环境有并发上限 |

#### 2.5.3 测试用例

| 功能 | 说明 |
|------|------|
| 用例管理 | 创建/编辑/删除测试用例 |
| 用例分类 | 正常case/边界case/异常case |
| 用例执行 | 批量执行测试用例 |
| 执行报告 | 测试结果/耗时/错误详情 |
| 自动化测试 | 支持定时自动执行 |

#### 2.5.4 灰度发布

| 功能 | 说明 |
|------|------|
| 灰度策略 | 百分比/白名单/设备型号 |
| 灰度监控 | 灰度期间实时监控指标 |
| 灰度升级 | 指标正常自动提升灰度比例 |
| 灰度回滚 | 指标异常自动/手动回滚 |

### 2.6 AI决策日志与可解释性

#### 2.6.1 决策日志

| 功能 | 说明 |
|------|------|
| 决策记录 | 每次AI决策记录输入/输出/上下文 |
| 日志存储 | 支持按时间/用户/设备查询 |
| 日志保留 | 默认保留90天，可配置 |
| 日志导出 | 支持导出为JSON/CSV |

#### 2.6.2 AI可解释性

| 功能 | 说明 |
|------|------|
| 决策原因 | 记录决策的主要影响因素 |
| 解释界面 | 用户可在App查看AI为什么这么做 |
| 解释级别 | 简单模式/详细模式 |
| 解释可信度 | 同时展示解释本身的置信度 |

### 2.7 AI行为监控与异常告警

#### 2.7.1 行为监控

| 功能 | 说明 |
|------|------|
| 行为基线 | 基于历史数据建立正常行为基线 |
| 异常检测 | 实时检测偏离基线的异常行为 |
| 行为报告 | 定期生成AI行为分析报告 |

#### 2.7.2 异常告警

| 功能 | 说明 |
|------|------|
| 告警规则 | 配置异常检测规则（阈值/趋势） |
| 告警级别 | 提示/警告/严重 |
| 告警通知 | 邮件/SMS/Webhook通知 |
| 告警抑制 | 避免告警风暴（去重/聚合） |
| 告警升级 | 持续未处理自动升级 |

### 2.8 边缘AI vs 云端AI路由

#### 2.8.1 路由策略

| 策略 | 说明 | 适用场景 |
|------|------|----------|
| 边缘优先 | 优先端侧推理，失败时走云端 | 低延迟敏感 |
| 云端优先 | 优先云端推理，边缘做备份 | 精度优先 |
| 智能路由 | 根据网络/任务类型/模型大小自动选择 | 通用场景 |
| 强制边缘 | 仅端侧推理 | 隐私敏感 |
| 强制云端 | 仅云端推理 | 端侧模型不支持 |

#### 2.8.2 路由配置

| 功能 | 说明 |
|------|------|
| 任务类型路由 | 不同任务类型走不同推理路径 |
| 网络质量感知 | 根据网络延迟自动切换 |
| 模型大小阈值 | 小于X MB走边缘，大于走云端 |
| 降级策略 | 云端不可用时自动降级到边缘 |

### 2.9 模型分片加载

#### 2.9.1 分片策略

| 功能 | 说明 |
|------|------|
| 按层分片 | 模型按层拆分，按需加载 |
| 按功能分片 | 基础能力/高级能力分片 |
| 动态分片 | 根据用户使用情况动态调整 |
| 分片预加载 | 预测用户行为提前加载 |

#### 2.9.2 分片管理

| 功能 | 说明 |
|------|------|
| 分片列表 | 展示所有模型分片 |
| 分片依赖 | 定义分片间依赖关系 |
| 分片大小监控 | 监控分片加载时间 |
| 分片缓存 | 热门分片本地缓存 |

---

## 三、API接口定义

### 3.1 训练流水线

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/ai/training-tasks | 创建训练任务 |
| GET | /api/v1/ai/training-tasks | 训练任务列表 |
| GET | /api/v1/ai/training-tasks/:id | 训练任务详情 |
| PUT | /api/v1/ai/training-tasks/:id | 更新训练任务 |
| POST | /api/v1/ai/training-tasks/:id/pause | 暂停训练 |
| POST | /api/v1/ai/training-tasks/:id/resume | 恢复训练 |
| POST | /api/v1/ai/training-tasks/:id/cancel | 取消训练 |
| GET | /api/v1/ai/training-tasks/:id/metrics | 训练指标 |

### 3.2 数据集管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/ai/datasets | 创建数据集 |
| GET | /api/v1/ai/datasets | 数据集列表 |
| GET | /api/v1/ai/datasets/:id | 数据集详情 |
| PUT | /api/v1/ai/datasets/:id | 更新数据集 |
| DELETE | /api/v1/ai/datasets/:id | 删除数据集 |
| POST | /api/v1/ai/datasets/:id/upload | 上传数据 |
| GET | /api/v1/ai/datasets/:id/stats | 数据集统计 |

### 3.3 A/B实验

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/ai/experiments | 创建实验 |
| GET | /api/v1/ai/experiments | 实验列表 |
| GET | /api/v1/ai/experiments/:id | 实验详情 |
| PUT | /api/v1/ai/experiments/:id | 更新实验 |
| POST | /api/v1/ai/experiments/:id/start | 启动实验 |
| POST | /api/v1/ai/experiments/:id/stop | 停止实验 |
| GET | /api/v1/ai/experiments/:id/results | 实验结果 |
| POST | /api/v1/ai/experiments/:id/rollback | 全量发布胜出组 |

### 3.4 模型管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/ai/models | 模型列表 |
| GET | /api/v1/ai/models/:id | 模型详情 |
| POST | /api/v1/ai/models/:id/deploy | 部署模型 |
| POST | /api/v1/ai/models/:id/rollback | 回滚模型 |
| POST | /api/v1/ai/models/:id/offline | 下线模型 |
| GET | /api/v1/ai/models/:id/versions | 模型版本列表 |
| POST | /api/v1/ai/models/:id/versions | 创建新版本 |

### 3.5 监控告警

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/ai/monitoring/dashboard | 监控仪表盘 |
| GET | /api/v1/ai/monitoring/metrics | 指标查询 |
| GET | /api/v1/ai/monitoring/alerts | 告警列表 |
| POST | /api/v1/ai/monitoring/alerts/:id/ack | 确认告警 |
| GET | /api/v1/ai/monitoring/alerts/rules | 告警规则列表 |
| POST | /api/v1/ai/monitoring/alerts/rules | 创建告警规则 |

### 3.6 沙箱测试

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/ai/sandbox/environments | 创建测试环境 |
| GET | /api/v1/ai/sandbox/environments/:id | 环境详情 |
| DELETE | /api/v1/ai/sandbox/environments/:id | 删除环境 |
| POST | /api/v1/ai/sandbox/testcases | 创建测试用例 |
| POST | /api/v1/ai/sandbox/testcases/execute | 执行测试用例 |
| GET | /api/v1/ai/sandbox/reports/:id | 测试报告 |

### 3.7 AI决策日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/ai/decisions | 决策日志列表 |
| GET | /api/v1/ai/decisions/:id | 决策详情 |
| GET | /api/v1/ai/decisions/explain/:id | 决策解释 |

### 3.8 路由策略

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/ai/routing/policies | 路由策略列表 |
| POST | /api/v1/ai/routing/policies | 创建路由策略 |
| PUT | /api/v1/ai/routing/policies/:id | 更新路由策略 |
| DELETE | /api/v1/ai/routing/policies/:id | 删除路由策略 |

---

## 四、数据库设计

### 4.1 训练任务表 (ai_training_tasks)

```sql
CREATE TABLE ai_training_tasks (
    id              BIGSERIAL PRIMARY KEY,
    task_name       VARCHAR(255) NOT NULL,
    model_type      VARCHAR(50) NOT NULL,           -- 'vision'/'audio'/'multimodal'/'rl'
    model_template  VARCHAR(100),
    dataset_id      BIGINT REFERENCES ai_datasets(id),
    config          JSONB,                          -- 超参数配置
    status          VARCHAR(20) NOT NULL,          -- 'pending'/'running'/'paused'/'completed'/'failed'/'cancelled'
    started_at      TIMESTAMP,
    finished_at     TIMESTAMP,
    metrics         JSONB,                          -- 最终指标结果
    model_version_id BIGINT REFERENCES ai_model_versions(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_training_tasks_status ON ai_training_tasks(status);
CREATE INDEX idx_training_tasks_dataset ON ai_training_tasks(dataset_id);
```

### 4.2 数据集表 (ai_datasets)

```sql
CREATE TABLE ai_datasets (
    id              BIGSERIAL PRIMARY KEY,
    dataset_name    VARCHAR(255) NOT NULL,
    dataset_type    VARCHAR(50) NOT NULL,           -- 'image'/'audio'/'video'/'text'
    version         VARCHAR(20) DEFAULT '1.0',
    parent_version  BIGINT REFERENCES ai_datasets(id),
    stats           JSONB,                          -- 数据量/类别分布
    tags            VARCHAR(100)[],
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.3 A/B实验表 (ai_experiments)

```sql
CREATE TABLE ai_experiments (
    id              BIGSERIAL PRIMARY KEY,
    experiment_name VARCHAR(255) NOT NULL,
    description     TEXT,
    model_id        BIGINT REFERENCES ai_models(id),
    status          VARCHAR(20) NOT NULL,          -- 'draft'/'running'/'paused'/'stopped'/'completed'
    start_time      TIMESTAMP,
    end_time        TIMESTAMP,
    target_metric   VARCHAR(50),                    -- 优化目标
    min_sample_size INT DEFAULT 1000,
    winner_group_id BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.4 实验分组表 (ai_experiment_groups)

```sql
CREATE TABLE ai_experiment_groups (
    id              BIGSERIAL PRIMARY KEY,
    experiment_id   BIGINT REFERENCES ai_experiments(id) ON DELETE CASCADE,
    group_name      VARCHAR(100) NOT NULL,
    group_type      VARCHAR(20) NOT NULL,           -- 'control'/'treatment'
    traffic_percent DECIMAL(5,2) NOT NULL,         -- 流量占比 0-100
    model_version_id BIGINT REFERENCES ai_model_versions(id),
    config          JSONB,                          -- 实验配置参数
    is_control      BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.5 模型表 (ai_models)

```sql
CREATE TABLE ai_models (
    id              BIGSERIAL PRIMARY KEY,
    model_name      VARCHAR(255) NOT NULL,
    model_type      VARCHAR(50) NOT NULL,           -- 'vision'/'audio'/'multimodal'/'rl'
    description     TEXT,
    current_version_id BIGINT,
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.6 模型版本表 (ai_model_versions)

```sql
CREATE TABLE ai_model_versions (
    id              BIGSERIAL PRIMARY KEY,
    model_id        BIGINT REFERENCES ai_models(id),
    version         VARCHAR(20) NOT NULL,
    version_tag     VARCHAR(50),                   -- 'stable'/'beta'/'candidate'
    status          VARCHAR(20) NOT NULL,          -- 'online'/'gray'/'offline'/'draft'
    model_file_url  VARCHAR(500),                  -- 模型文件存储地址
    model_size_mb   DECIMAL(10,2),
    config          JSONB,                          -- 模型配置
    metrics         JSONB,                          -- 评估指标
    training_task_id BIGINT REFERENCES ai_training_tasks(id),
    deployed_at     TIMESTAMP,
    offline_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(model_id, version)
);
```

### 4.7 监控指标表 (ai_monitoring_metrics)

```sql
CREATE TABLE ai_monitoring_metrics (
    id              BIGSERIAL PRIMARY KEY,
    model_version_id BIGINT REFERENCES ai_model_versions(id),
    metric_name     VARCHAR(50) NOT NULL,
    metric_value    DECIMAL(15,5) NOT NULL,
    metric_unit     VARCHAR(20),
    percentile      VARCHAR(10),                   -- 'p50'/'p95'/'p99'
    recorded_at     TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_monitoring_metrics_model ON ai_monitoring_metrics(model_version_id, recorded_at);
```

### 4.8 告警规则表 (ai_alert_rules)

```sql
CREATE TABLE ai_alert_rules (
    id              BIGSERIAL PRIMARY KEY,
    rule_name       VARCHAR(255) NOT NULL,
    metric_name     VARCHAR(50) NOT NULL,
    condition       VARCHAR(20) NOT NULL,         -- '>'/'<'/'>='/'<='/=='
    threshold       DECIMAL(15,5) NOT NULL,
    severity        VARCHAR(20) NOT NULL,          -- 'info'/'warning'/'critical'
    notify_ways     VARCHAR(50)[],                -- 'email'/'sms'/'webhook'
    enabled         BOOLEAN DEFAULT TRUE,
    cooldown_minutes INT DEFAULT 10,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.9 沙箱环境表 (ai_sandbox_environments)

```sql
CREATE TABLE ai_sandbox_environments (
    id              BIGSERIAL PRIMARY KEY,
    env_name        VARCHAR(255) NOT NULL,
    env_config      JSONB,                          -- 环境配置
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMP DEFAULT NOW(),
    expired_at      TIMESTAMP
);
```

### 4.10 测试用例表 (ai_sandbox_testcases)

```sql
CREATE TABLE ai_sandbox_testcases (
    id              BIGSERIAL PRIMARY KEY,
    env_id          BIGINT REFERENCES ai_sandbox_environments(id),
    case_name       VARCHAR(255) NOT NULL,
    case_type       VARCHAR(20) NOT NULL,         -- 'normal'/'boundary'/'error'
    input_data      JSONB NOT NULL,
    expected_output JSONB,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.11 AI决策日志表 (ai_decision_logs)

```sql
CREATE TABLE ai_decision_logs (
    id              BIGSERIAL PRIMARY KEY,
    model_version_id BIGINT REFERENCES ai_model_versions(id),
    device_id       VARCHAR(100),
    user_id         BIGINT,
    task_type       VARCHAR(50) NOT NULL,
    input_data      JSONB NOT NULL,
    output_data     JSONB NOT NULL,
    decision_reason JSONB,                          -- 决策原因
    latency_ms      INTEGER,
    inference_mode  VARCHAR(20),                   -- 'edge'/'cloud'
    created_at      TIMESTAMP NOT NULL
);

CREATE INDEX idx_decision_logs_device ON ai_decision_logs(device_id, created_at);
CREATE INDEX idx_decision_logs_user ON ai_decision_logs(user_id, created_at);
CREATE INDEX idx_decision_logs_time ON ai_decision_logs(created_at);
```

### 4.12 路由策略表 (ai_routing_policies)

```sql
CREATE TABLE ai_routing_policies (
    id              BIGSERIAL PRIMARY KEY,
    policy_name     VARCHAR(255) NOT NULL,
    task_type       VARCHAR(50) NOT NULL,
    strategy        VARCHAR(20) NOT NULL,           -- 'edge_first'/'cloud_first'/'smart'/'force_edge'/'force_cloud'
    config          JSONB,                          -- 策略配置参数
    priority        INTEGER DEFAULT 0,
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.13 模型分片表 (ai_model_shards)

```sql
CREATE TABLE ai_model_shards (
    id              BIGSERIAL PRIMARY KEY,
    model_version_id BIGINT REFERENCES ai_model_versions(id),
    shard_name      VARCHAR(100) NOT NULL,
    shard_index     INTEGER NOT NULL,
    file_url        VARCHAR(500),
    size_mb         DECIMAL(10,2),
    dependencies    BIGINT[],                       -- 依赖的分片ID
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 五、前端页面清单

### 5.1 AI工作台

| 页面 | 路由 | 说明 |
|------|------|------|
| AI工作台首页 | /ai/dashboard | 核心指标大盘 + 快捷入口 |
| 模型列表 | /ai/models | 所有模型 + 版本状态 |

### 5.2 训练中心

| 页面 | 路由 | 说明 |
|------|------|------|
| 训练任务列表 | /ai/training | 训练任务列表 |
| 创建训练任务 | /ai/training/create | 训练配置表单 |
| 训练详情/监控 | /ai/training/:id | 训练详情 + 实时曲线 |
| 数据集管理 | /ai/datasets | 数据集列表 |
| 数据集详情 | /ai/datasets/:id | 数据集统计 + 标注 |

### 5.3 A/B实验

| 页面 | 路由 | 说明 |
|------|------|------|
| 实验列表 | /ai/experiments | 所有实验 |
| 创建实验 | /ai/experiments/create | 实验配置向导 |
| 实验详情 | /ai/experiments/:id | 实验数据 + 效果对比 |
| 实验报告 | /ai/experiments/:id/report | 自动生成报告 |

### 5.4 模型管理

| 页面 | 路由 | 说明 |
|------|------|------|
| 模型列表 | /ai/models | 模型列表 + 状态 |
| 模型详情 | /ai/models/:id | 版本列表 + 指标对比 |
| 部署配置 | /ai/models/:id/deploy | 部署策略配置 |
| 回滚确认 | /ai/models/:id/rollback | 回滚操作 |

### 5.5 监控中心

| 页面 | 路由 | 说明 |
|------|------|------|
| 监控仪表盘 | /ai/monitoring | 实时监控大盘 |
| 指标详情 | /ai/monitoring/metrics | 指标历史查询 |
| 告警规则 | /ai/monitoring/rules | 规则配置 |
| 告警记录 | /ai/monitoring/alerts | 历史告警 |

### 5.6 沙箱测试

| 页面 | 路由 | 说明 |
|------|------|------|
| 沙箱环境 | /ai/sandbox/environments | 环境管理 |
| 测试用例 | /ai/sandbox/testcases | 用例管理 |
| 测试报告 | /ai/sandbox/reports | 测试报告列表 |

### 5.7 AI决策

| 页面 | 路由 | 说明 |
|------|------|------|
| 决策日志 | /ai/decisions | 决策日志查询 |
| 决策详情 | /ai/decisions/:id | 单次决策 + 解释 |
| 路由策略 | /ai/routing | 路由配置 |

---

## 六、验收标准

### 6.1 AI训练流水线

| 验收点 | 标准 |
|--------|------|
| 任务创建 | 填写配置后成功创建训练任务，状态变为pending |
| 状态追踪 | 训练进行中实时更新loss曲线，延迟<5秒 |
| 任务暂停/恢复 | 暂停后训练停止，恢复后从断点继续 |
| 任务取消 | 取消后训练立即停止，资源释放 |
| 模型导出 | 训练完成后自动生成可部署模型文件 |

### 6.2 A/B测试

| 验收点 | 标准 |
|--------|------|
| 分流准确性 | 各分组流量比例与配置偏差<2% |
| 分组一致性 | 同一用户多次请求始终路由到同一分组 |
| 效果统计 | 实验结束后自动计算p值和置信区间 |
| 一键全量 | 点击后胜出分组在<30秒内全量生效 |

### 6.3 模型监控

| 验收点 | 标准 |
|--------|------|
| 指标实时性 | 指标延迟<10秒 |
| 告警触发 | 指标超过阈值后5分钟内发出告警 |
| 告警去重 | 同一规则告警5分钟内不重复发送 |

### 6.4 模型回滚

| 验收点 | 标准 |
|--------|------|
| 回滚时间 | 点击回滚到生效<30秒 |
| 回滚记录 | 回滚操作记录包含操作人/时间/原因 |
| 回滚验证 | 回滚后模型质量指标与目标版本一致 |

### 6.5 沙箱测试

| 验收点 | 标准 |
|--------|------|
| 环境隔离 | 沙箱环境与生产环境完全隔离 |
| 测试执行 | 测试用例执行结果准确，错误详情清晰 |
| 测试报告 | 报告包含通过率/耗时/失败原因 |

### 6.6 AI决策日志

| 验收点 | 标准 |
|--------|------|
| 日志记录 | 每笔决策完整记录输入/输出/原因 |
| 查询性能 | 100万条数据内查询响应<2秒 |
| 解释展示 | 用户可在App看到决策原因（简化版） |

### 6.7 路由策略

| 验收点 | 标准 |
|--------|------|
| 策略生效 | 路由策略变更后<1分钟生效 |
| 降级切换 | 云端不可用时自动切换到边缘<5秒 |
| 流量统计 | 边缘/云端流量统计准确 |


---

## 七、页面布局规范

### 7.1 AI工作台首页（/ai/dashboard）

**布局结构：**
1. 面包屑 → 页面标题
2. 统计卡片（在线模型数/活跃实验数/告警数/推理请求量）—— 白色
3. 核心指标图表（延迟趋势/成功率/资源利用率）
4. 快捷入口卡片

### 7.2 训练中心页面（/ai/training）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片 #F2F3F5）：任务状态 / 模型类型
3. 操作栏（创建训练任务靠左）
4. 训练任务列表表格

**按钮规范：**
- [创建训练任务] — 左对齐
- [详情] [暂停/恢复] [取消] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 任务名称 | 200px | - |
| 模型类型 | 100px | - |
| 数据集 | 150px | - |
| 状态 | 100px | pending/running/paused/completed/failed |
| 创建时间 | 150px | - |
| 操作 | 120px | 详情/暂停/恢复/取消 |

**分页：** 右下角，10/20/50/100 条

### 7.3 数据集管理页面（/ai/datasets）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：数据类型 / 标签
3. 操作栏（创建数据集/上传数据靠左）
4. 数据集列表表格

**按钮规范：**
- [创建数据集] [上传数据] — 左对齐
- [详情] [编辑] [删除] — 行内右对齐

**分页：** 右下角，10/20/50/100 条

### 7.4 A/B实验页面（/ai/experiments）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：实验状态 / 开始时间
3. 操作栏（创建实验靠左）
4. 实验列表表格

**按钮规范：**
- [创建实验] — 左对齐
- [详情] [启动/停止] [全量发布] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 实验名称 | 200px | - |
| 关联模型 | 150px | - |
| 状态 | 100px | draft/running/paused/completed |
| 流量分配 | 100px | - |
| 胜出分组 | 100px | - |
| 操作 | 120px | 详情/启动/停止/全量发布 |

**分页：** 右下角，10/20/50/100 条

### 7.5 模型管理页面（/ai/models）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：模型类型 / 状态
3. 操作栏（部署模型靠右）
4. 模型列表表格

**按钮规范：**
- [部署] [回滚] — 行内右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 模型名称 | 200px | - |
| 模型类型 | 100px | - |
| 当前版本 | 120px | - |
| 状态 | 100px | online/gray/offline/draft |
| 部署时间 | 150px | - |
| 操作 | 120px | 详情/部署/回滚/下线 |

**分页：** 右下角，10/20/50/100 条

### 7.6 监控仪表盘页面（/ai/monitoring）

**布局结构：**
1. 面包屑 → 页面标题
2. 指标卡片（响应延迟P99/成功率/资源利用率/并发数）—— 白色
3. 实时趋势图表
4. 告警列表表格

**分页：** 右下角，10/20/50/100 条

### 7.7 沙箱测试页面（/ai/sandbox/environments）

**布局结构：**
1. 面包屑 → 页面标题
2. 操作栏（创建环境靠左）
3. 沙箱环境列表表格
4. 测试用例Tab

**按钮规范：**
- [创建环境] [创建测试用例] — 左对齐
- [执行] [删除] — 行内右对齐

**分页：** 右下角，10/20/50/100 条

### 7.8 AI决策日志页面（/ai/decisions）

**布局结构：**
1. 面包屑 → 页面标题
2. 筛选区（浅灰卡片）：设备ID / 用户ID / 任务类型 / 时间范围
3. 操作栏（导出靠右）
4. 决策日志列表表格

**按钮规范：**
- [导出] — 右对齐

**表格列：**
| 列名 | 宽度 | 说明 |
|------|------|------|
| 序号 | 60px | - |
| 决策ID | 150px | - |
| 设备ID | 150px | - |
| 任务类型 | 120px | - |
| 推理模式 | 100px | edge/cloud |
| 延迟ms | 100px | - |
| 决策时间 | 150px | - |
| 操作 | 120px | 详情/解释 |

**分页：** 右下角，10/20/50/100 条

### 7.9 弹窗规范

| 类型 | 使用场景 |
|------|----------|
| Drawer 抽屉 | 创建/编辑训练任务、创建/编辑实验、创建数据集 |
| Dialog 对话框 | 确认暂停/取消/删除 |
| 全屏模态 | 暂无复杂表单场景 |
