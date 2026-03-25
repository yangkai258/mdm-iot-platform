# PRD：AI行为研究平台

**版本：** V1.0
**所属Phase：** Phase 4（Sprint 32）
**优先级：** P3
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

AI行为研究平台为学术研究机构提供宠物AI行为研究的在线实验环境，支持研究假设创建、实验环境配置、实验执行、结果分析、成果发表，构建"研究假设→实验验证→成果产出→产品迭代"的完整闭环。

### 1.2 核心价值

- **加速研究**：提供标准化实验环境，降低研究门槛
- **数据驱动**：基于真实用户数据进行实验
- **成果转化**：学术成果直接反哺产品能力

### 1.3 范围边界

**包含：**
- 研究机构注册和认证
- 研究假设管理
- 实验环境创建和配置
- 实验执行和监控
- 实验结果分析和对比
- 研究成果发表和引用

**不包含：**
- 原始数据直接提供（通过数据集平台提供）
- 实验代码托管（使用第三方Git服务）
- 论文撰写工具（仅提供数据支撑）

---

## 二、系统架构

### 2.1 技术架构图

```
┌──────────────────────────────────────────────────────────────────────┐
│                        前端 (Vue 3 + Arco Design)                     │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌────────────┐  │
│  │ 研究首页     │ │ 实验管理     │ │ 结果分析    │ │ 协作空间  │  │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └─────┬──────┘  │
└──────────┼───────────────┼───────────────┼────────────────┼─────────┘
           │               │               │                │
           └───────────────┴──────┐────────┴────────────────┘
                                  │ REST API / WebSocket
┌─────────────────────────────────┼────────────────────────────────────┐
│                          后端服务 (Go + Gin)                           │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │                  ResearchPlatform Service                      │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │   │
│  │  │PlatformAPI│ │Experiment│ │ResultAPI │ │CollabAPI │        │   │
│  │  │          │ │   API    │ │          │ │          │        │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │   │
│  └───────────────────────────────────────────────────────────────┘   │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │           Experiment Execution Engine (容器化)                  │   │
│  └───────────────────────────────────────────────────────────────┘   │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │              Data Isolation Layer (ACL + 审计)                  │   │
│  └───────────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────┘
        │                    │                    │
        ▼                    ▼                    ▼
┌───────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  PostgreSQL   │  │  Kubernetes    │  │      OSS        │
│ (元数据存储)   │  │ (实验执行环境)  │  │ (实验数据存储)  │
└───────────────┘  └─────────────────┘  └─────────────────┘
```

---

## 三、功能详情

### 3.1 研究机构管理

| 功能 | 说明 |
|------|------|
| 机构注册 | 研究机构注册账号，提交资质证明 |
| 机构认证 | 管理员审核机构资质，认证通过后开通研究功能 |
| 机构主页 | 展示机构信息、研究项目、发表成果 |
| 机构管理 | 机构管理员管理成员账户 |

### 3.2 研究假设管理

| 功能 | 说明 |
|------|------|
| 假设创建 | 创建研究假设，填写假设内容、预期结果 |
| 假设列表 | 展示所有研究假设及状态 |
| 假设详情 | 查看假设详情和关联实验 |
| 假设标签 | 按研究领域标签分类 |

### 3.3 实验环境管理

| 功能 | 说明 |
|------|------|
| 环境模板 | 提供标准实验环境模板（行为分析/情感计算/健康预测） |
| 环境创建 | 基于模板创建实验环境 |
| 环境配置 | 配置实验参数、数据集选择、评估指标 |
| 环境启动 | 启动实验环境，自动分配计算资源 |
| 环境监控 | 实时监控实验运行状态和资源使用 |
| 环境停止 | 停止实验环境，释放计算资源 |

### 3.4 实验执行

| 功能 | 说明 |
|------|------|
| 单一实验 | 运行单个实验，获取结果 |
| 参数扫描 | 对多个参数组合批量实验 |
| 实验对比 | 对比不同配置下的实验结果 |
| 中断恢复 | 实验中断后可从断点恢复 |

### 3.5 结果分析

| 功能 | 说明 |
|------|------|
| 原始结果 | 查看实验原始输出数据 |
| 统计摘要 | 自动生成统计摘要（均值/方差/p值） |
| 可视化 | 自动生成结果图表（折线/柱状/热力图） |
| 对比分析 | 对比多次实验的结果差异 |
| 结果导出 | 导出结果为CSV/JSON/PDF |

### 3.6 研究协作

| 功能 | 说明 |
|------|------|
| 邀请协作者 | 通过邮件邀请研究伙伴加入 |
| 权限管理 | 设置协作者权限（查看/编辑/管理） |
| 讨论区 | 实验讨论和问题解答 |
| 版本历史 | 实验配置版本管理 |

### 3.7 成果发表

| 功能 | 说明 |
|------|------|
| 成果提交 | 提交研究成果（论文/专利/报告） |
| 成果审核 | 管理员审核成果内容 |
| 成果展示 | 审核通过后在平台展示 |
| 引用追踪 | 追踪成果被引用情况 |

---

## 四、API接口定义

### 4.1 机构管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/research/institutions | 注册机构 |
| GET | /api/v1/research/institutions/:id | 机构详情 |
| PUT | /api/v1/research/institutions/:id | 更新机构信息 |
| POST | /api/v1/research/institutions/:id/members | 添加成员 |
| DELETE | /api/v1/research/institutions/:id/members/:userId | 移除成员 |

### 4.2 实验管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/research/experiments | 实验列表 |
| POST | /api/v1/research/experiments | 创建实验 |
| GET | /api/v1/research/experiments/:id | 实验详情 |
| PUT | /api/v1/research/experiments/:id | 更新实验 |
| DELETE | /api/v1/research/experiments/:id | 删除实验 |
| POST | /api/v1/research/experiments/:id/start | 启动实验 |
| POST | /api/v1/research/experiments/:id/stop | 停止实验 |
| GET | /api/v1/research/experiments/:id/status | 实验状态 |
| GET | /api/v1/research/experiments/:id/logs | 实验日志 |

### 4.3 环境管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/research/env-templates | 环境模板列表 |
| GET | /api/v1/research/env-templates/:id | 模板详情 |
| POST | /api/v1/research/experiments/:id/env | 创建实验环境 |
| GET | /api/v1/research/experiments/:id/env | 实验环境详情 |
| POST | /api/v1/research/experiments/:id/env/start | 启动环境 |
| POST | /api/v1/research/experiments/:id/env/stop | 停止环境 |
| GET | /api/v1/research/experiments/:id/env/metrics | 环境监控指标 |

### 4.4 结果分析

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/research/experiments/:id/results | 实验结果 |
| GET | /api/v1/research/experiments/:id/results/summary | 统计摘要 |
| GET | /api/v1/research/experiments/:id/results/charts | 图表数据 |
| GET | /api/v1/research/experiments/:id/results/export | 导出结果 |
| POST | /api/v1/research/experiments/compare | 对比多个实验 |

### 4.5 协作管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/research/experiments/:id/collaborators | 协作者列表 |
| POST | /api/v1/research/experiments/:id/collaborators | 添加协作者 |
| PUT | /api/v1/research/experiments/:id/collaborators/:userId | 更新权限 |
| DELETE | /api/v1/research/experiments/:id/collaborators/:userId | 移除协作者 |
| GET | /api/v1/research/experiments/:id/discussions | 讨论列表 |
| POST | /api/v1/research/experiments/:id/discussions | 发帖 |

### 4.6 成果管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/research/publications | 成果列表 |
| POST | /api/v1/research/publications | 提交成果 |
| GET | /api/v1/research/publications/:id | 成果详情 |
| GET | /api/v1/research/publications/:id/citations | 引用列表 |

---

## 五、数据库设计

### 5.1 研究机构表 (research_institutions)

```sql
CREATE TABLE research_institutions (
    id                  BIGSERIAL PRIMARY KEY,
    institution_name    VARCHAR(200) NOT NULL,
    institution_type    VARCHAR(50) NOT NULL,            -- 'university'/'research_institute'/'company'
    country             VARCHAR(50),
    website_url         VARCHAR(500),
    description         TEXT,
    logo_url            VARCHAR(500),
    verified_status     VARCHAR(20) DEFAULT 'pending',    -- 'pending'/'verified'/'rejected'
    verified_at         TIMESTAMP,
    verified_by         BIGINT,
    member_count        INT DEFAULT 0,
    publication_count   INT DEFAULT 0,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);
```

### 5.2 研究假设表 (research_hypotheses)

```sql
CREATE TABLE research_hypotheses (
    id                  BIGSERIAL PRIMARY KEY,
    researcher_id       BIGINT NOT NULL REFERENCES users(id),
    institution_id      BIGINT REFERENCES research_institutions(id),
    hypothesis_title    VARCHAR(300) NOT NULL,
    hypothesis_content  TEXT NOT NULL,
    research_domain     VARCHAR(50) NOT NULL,            -- 'behavior'/'emotion'/'health'/'voice'
    expected_outcome    TEXT,
    related_work        TEXT,
    tags                VARCHAR(50)[],
    status              VARCHAR(20) DEFAULT 'draft',    -- 'draft'/'approved'/'running'/'completed'
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_hypotheses_domain ON research_hypotheses(research_domain);
CREATE INDEX idx_hypotheses_researcher ON research_hypotheses(researcher_id);
```

### 5.3 实验环境模板表 (experiment_env_templates)

```sql
CREATE TABLE experiment_env_templates (
    id                  BIGSERIAL PRIMARY KEY,
    template_name       VARCHAR(100) NOT NULL,
    template_type       VARCHAR(50) NOT NULL,            -- 'behavior'/'emotion'/'health'
    description         TEXT,
    config_schema       JSONB NOT NULL,
    compute_resources   JSONB NOT NULL,
    software_stack      JSONB,
    dataset_binding     JSONB,
    is_public           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMP DEFAULT NOW()
);
```

### 5.4 研究实验表 (research_experiments)

```sql
CREATE TABLE research_experiments (
    id                  BIGSERIAL PRIMARY KEY,
    experiment_name     VARCHAR(200) NOT NULL,
    researcher_id       BIGINT NOT NULL REFERENCES users(id),
    institution_id      BIGINT REFERENCES research_institutions(id),
    hypothesis_id      BIGINT REFERENCES research_hypotheses(id),
    experiment_type     VARCHAR(50) NOT NULL,
    methodology         TEXT,
    hypothesis          TEXT,
    config              JSONB NOT NULL,
    dataset_ids         BIGINT[],
    status              VARCHAR(20) DEFAULT 'draft',  -- 'draft'/'pending'/'approved'/'running'/'paused'/'completed'/'failed'
    env_id              VARCHAR(100),
    started_at          TIMESTAMP,
    completed_at        TIMESTAMP,
    estimated_duration  INT,
    actual_duration     INT,
    results_summary     JSONB,
    is_public           BOOLEAN DEFAULT FALSE,
    published_at        TIMESTAMP,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_experiments_researcher ON research_experiments(researcher_id);
CREATE INDEX idx_experiments_status ON research_experiments(status);
CREATE INDEX idx_experiments_institution ON research_experiments(institution_id);
```

### 5.5 实验协作者表 (research_collaborators)

```sql
CREATE TABLE research_collaborators (
    id                  BIGSERIAL PRIMARY KEY,
    experiment_id       BIGINT NOT NULL REFERENCES research_experiments(id) ON DELETE CASCADE,
    user_id             BIGINT NOT NULL REFERENCES users(id),
    role                VARCHAR(30) DEFAULT 'viewer',  -- 'viewer'/'contributor'/'admin'
    invited_by          BIGINT NOT NULL REFERENCES users(id),
    invited_at          TIMESTAMP DEFAULT NOW(),
    accepted_at         TIMESTAMP,
    status              VARCHAR(20) DEFAULT 'pending',  -- 'pending'/'accepted'/'declined'
    UNIQUE(experiment_id, user_id)
);

CREATE INDEX idx_collaborators_experiment ON research_collaborators(experiment_id);
```

### 5.6 实验结果表 (research_experiment_results)

```sql
CREATE TABLE research_experiment_results (
    id                  BIGSERIAL PRIMARY KEY,
    experiment_id       BIGINT NOT NULL REFERENCES research_experiments(id) ON DELETE CASCADE,
    result_type         VARCHAR(50) NOT NULL,            -- 'metric'/'model'/'analysis'/'visualization'
    metric_name         VARCHAR(100),
    metric_value        JSONB,
    raw_data_url        VARCHAR(500),
    chart_config        JSONB,
    notes               TEXT,
    created_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_results_experiment ON research_experiment_results(experiment_id);
```

### 5.7 研究成果表 (research_publications)

```sql
CREATE TABLE research_publications (
    id                  BIGSERIAL PRIMARY KEY,
    experiment_id       BIGINT REFERENCES research_experiments(id),
    researcher_id       BIGINT NOT NULL REFERENCES users(id),
    institution_id      BIGINT REFERENCES research_institutions(id),
    publication_type    VARCHAR(30) NOT NULL,            -- 'paper'/'patent'/'report'
    title               VARCHAR(300) NOT NULL,
    abstract            TEXT,
    paper_url           VARCHAR(500),
    publication_date    DATE,
    venue               VARCHAR(200),
    citation_count      INT DEFAULT 0,
    status              VARCHAR(20) DEFAULT 'pending',   -- 'pending'/'approved'/'published'
    reviewed_at         TIMESTAMP,
    created_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_publications_researcher ON research_publications(researcher_id);
CREATE INDEX idx_publications_status ON research_publications(status);
```

### 5.8 实验讨论表 (research_discussions)

```sql
CREATE TABLE research_discussions (
    id                  BIGSERIAL PRIMARY KEY,
    experiment_id       BIGINT NOT NULL REFERENCES research_experiments(id) ON DELETE CASCADE,
    user_id             BIGINT NOT NULL REFERENCES users(id),
    parent_id           BIGINT REFERENCES research_discussions(id),
    content             TEXT NOT NULL,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_discussions_experiment ON research_discussions(experiment_id);
```

---

## 六、前端页面

### 6.1 研究首页

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 研究首页 | /research | 顶部右侧 | 「创建实验」「我的实验」 |
| 机构注册 | /research/institution/register | 表单底部 | 「提交注册」 |

### 6.2 实验管理

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 实验列表 | /research/experiments | 列表顶部 | 「新建实验」「筛选」「导出」 |
| 实验详情 | /research/experiments/:id | 操作区左 | 「启动」「停止」「编辑」「对比」 |
| 实验配置 | /research/experiments/:id/config | 表单底部 | 「保存配置」「启动实验」 |

### 6.3 环境管理

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 环境模板 | /research/env-templates | 模板卡片左 | 「使用模板」 |
| 环境详情 | /research/experiments/:id/env | 操作区左 | 「启动环境」「停止环境」「查看监控」 |

### 6.4 结果分析

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 结果概览 | /research/experiments/:id/results | 页面顶部 | 「导出」「对比」「分享」 |
| 结果对比 | /research/experiments/compare | 表单底部 | 「开始对比」 |
| 图表详情 | /research/experiments/:id/results/charts/:chartId | 图表工具栏 | 「下载图表」「全屏」 |

### 6.5 协作与成果

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 协作者管理 | /research/experiments/:id/collaborators | 页面顶部 | 「邀请协作者」「筛选」 |
| 讨论区 | /research/experiments/:id/discussions | 发帖框下方 | 「发帖」「回复」 |
| 成果列表 | /research/publications | 页面顶部 | 「提交成果」「筛选」 |

### 6.6 前端UI规范

- 研究首页使用仪表盘布局，展示关键指标
- 实验列表使用表格布局，支持排序和筛选
- 实验详情页使用标签页，左侧为导航，右侧为内容
- 结果分析使用图表为主，配合数据表格
- 环境监控使用实时图表，WebSocket推送更新

---

## 七、验收标准

### 7.1 功能验收

| 验收点 | 标准 |
|--------|------|
| 机构注册 | 机构注册流程完整，审核<5个工作日 |
| 实验创建 | 实验创建<2min，环境就绪<10min |
| 实验执行 | 实时状态推送，状态更新<1s |
| 结果分析 | 图表生成<5s，支持CSV/JSON导出 |
| 协作功能 | 协作者权限正确，结果隔离100% |

### 7.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 实验列表加载 | < 1s（100个实验） |
| 实验启动 | 环境就绪<10min |
| 结果查询 | < 2s |
| 图表渲染 | < 3s |

### 7.3 安全验收

| 验收点 | 标准 |
|--------|------|
| 数据隔离 | 实验间数据完全隔离 |
| 权限控制 | 协作者权限100%生效 |
| 操作审计 | 所有操作有完整日志 |

---

## 八、附录

### 8.1 环境模板类型

| 模板类型 | 适用场景 | 计算资源 |
|----------|----------|----------|
| 行为分析 | 宠物行为模式研究 | CPU 8核/16GB |
| 情感计算 | 宠物情感识别研究 | GPU/16GB |
| 健康预测 | 宠物健康趋势预测 | CPU 8核/16GB |

### 8.2 研究领域标签

| 标签 | 说明 |
|------|------|
| behavior | 宠物行为分析 |
| emotion | 情感计算 |
| health | 健康预测 |
| voice | 语音分析 |
| embodied | 具身智能 |
| digital-twin | 数字孪生 |

### 8.3 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|---------|
| V1.0 | 2026-03-24 | agentcp | 初稿创建 |
