# PRD：数据集开放平台

**版本：** V1.0
**所属Phase：** Phase 4（Sprint 31）
**优先级：** P3
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

数据集开放平台为学术研究机构提供宠物行为、情感、健康、语音等领域的高质量标注数据集，支持数据访问申请、审批、数据集版本管理，吸引学术合作，提升产品学术影响力，构建"学术研究+产品迭代"的正向循环。

### 1.2 核心价值

- **学术贡献**：开放高质量数据集，推动宠物AI研究
- **生态构建**：吸引学术机构，建立长期合作关系
- **持续迭代**：学术研究成果反哺产品能力

### 1.3 范围边界

**包含：**
- 数据集浏览和详情展示
- 数据访问申请和审批
- 数据集下载管理
- 数据集版本管理
- 数据集统计分析

**不包含：**
- 原始数据的直接提供（涉及隐私）
- 第三方数据集托管（仅托管自有数据集）
- 数据集标注平台（标注在本地完成）

---

## 二、系统架构

### 2.1 技术架构图

```
┌──────────────────────────────────────────────────────────────────────┐
│                        前端 (Vue 3 + Arco Design)                     │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌────────────┐  │
│  │ 数据集浏览   │ │ 申请管理    │ │  统计面板   │ │ 管理后台  │  │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └─────┬──────┘  │
└──────────┼───────────────┼───────────────┼────────────────┼─────────┘
           │               │               │                │
           └───────────────┴──────┐────────┴────────────────┘
                                  │ REST API
┌─────────────────────────────────┼────────────────────────────────────┐
│                          后端服务 (Go + Gin)                           │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │                     DatasetPlatform Service                     │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────┐  │   │
│  │  │DatasetAPI │  │AccessAPI   │  │VersionAPI  │  │StatsAPI│  │   │
│  │  └────────────┘  └────────────┘  └────────────┘  └────────┘  │   │
│  └───────────────────────────────────────────────────────────────┘   │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │               Data Access Control (ACL + 审批流)               │   │
│  └───────────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────┘
        │                    │                    │
        ▼                    ▼                    ▼
┌───────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  PostgreSQL   │  │      OSS        │  │    Redis        │
│ (元数据存储)   │  │ (数据集文件存储) │  │  (访问统计缓存)  │
└───────────────┘  └─────────────────┘  └─────────────────┘
```

---

## 三、功能详情

### 3.1 数据集浏览

| 功能 | 说明 |
|------|------|
| 数据集列表 | 按类型（行为/情感/健康/语音）分类展示 |
| 数据集搜索 | 按名称/类型/标签搜索 |
| 数据集详情 | 展示数据集描述、规模、格式、引用信息 |
| 数据集预览 | 预览部分数据样本（脱敏后） |
| 数据集引用 | 生成标准引用格式（APA/MLA/GB/T） |

### 3.2 数据集类型

| 类型 | 说明 | 记录数 | 字段 |
|------|------|--------|------|
| 宠物行为数据集 | 宠物日常行为标注数据 | >100,000条 | 时间戳/行为类型/环境/情绪 |
| 情感数据集 | 宠物情感反应标注数据 | >50,000条 | 情感类型/强度/触发因素 |
| 健康数据集 | 宠物健康指标标注数据 | >30,000条 | 指标类型/值/异常标记 |
| 语音数据集 | 宠物语音/叫声标注数据 | >20,000条 | 音频URL/语音类型/情绪标注 |

### 3.3 数据集详情

| 字段 | 说明 |
|------|------|
| 数据集名称 | 唯一标识名称 |
| 数据集类型 | 行为/情感/健康/语音 |
| 数据量 | 记录数/存储大小 |
| 数据格式 | JSON/CSV/Parquet |
| 授权类型 | CC-BY/CC0/研究专用 |
| 访问级别 | 公开/申请/受限 |
| 引用次数 | 学术引用次数 |
| 发表论文 | 关联论文链接 |

### 3.4 数据访问申请

| 功能 | 说明 |
|------|------|
| 申请入口 | 在数据集详情页发起申请 |
| 申请表单 | 机构信息/研究目的/数据安全保障承诺 |
| 审批流程 | 管理员初审+平台终审 |
| 审批通知 | 邮件+站内通知审批结果 |
| 访问期限 | 授权期限内可下载 |

### 3.5 数据集版本管理

| 功能 | 说明 |
|------|------|
| 版本列表 | 展示所有版本及变更说明 |
| 版本对比 | 对比两个版本的差异 |
| 下载指定版本 | 下载历史版本 |
| 版本标注 | 记录版本变更内容 |

---

## 四、API接口定义

### 4.1 数据集管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/datasets | 数据集列表 |
| GET | /api/v1/datasets/:id | 数据集详情 |
| GET | /api/v1/datasets/:id/samples | 预览样本（脱敏） |
| GET | /api/v1/datasets/:id/citation | 生成引用格式 |
| GET | /api/v1/datasets/:id/versions | 版本列表 |
| GET | /api/v1/datasets/:id/versions/:version | 指定版本详情 |
| GET | /api/v1/datasets/:id/stats | 数据集统计 |

### 4.2 访问申请

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/datasets/:id/access-request | 申请数据访问 |
| GET | /api/v1/datasets/access-requests | 我的申请列表 |
| GET | /api/v1/datasets/access-requests/:id | 申请详情 |
| POST | /api/v1/admin/datasets/access-requests/:id/approve | 审批通过 |
| POST | /api/v1/admin/datasets/access-requests/:id/reject | 审批拒绝 |
| GET | /api/v1/admin/datasets/access-requests | 所有申请列表（管理端） |

### 4.3 数据集下载

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/datasets/:id/download | 获取下载链接 |
| POST | /api/v1/datasets/:id/download/log | 记录下载日志 |

### 4.4 数据集管理（管理员）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/admin/datasets | 创建数据集 |
| PUT | /api/v1/admin/datasets/:id | 更新数据集 |
| POST | /api/v1/admin/datasets/:id/versions | 上传新版本 |
| DELETE | /api/v1/admin/datasets/:id | 删除数据集 |

---

## 五、数据库设计

### 5.1 数据集表 (research_datasets)

```sql
CREATE TABLE research_datasets (
    id                  BIGSERIAL PRIMARY KEY,
    dataset_name        VARCHAR(200) NOT NULL,
    dataset_type        VARCHAR(50) NOT NULL,            -- 'behavior'/'emotion'/'health'/'voice'
    description         TEXT,
    size_gb             DECIMAL(10,2),
    record_count        INT,
    data_format         VARCHAR(30),                     -- 'json'/'csv'/'parquet'
    download_url        VARCHAR(500),
    storage_path        VARCHAR(500),                    -- OSS存储路径
    license_type        VARCHAR(50),                     -- 'cc-by'/'cc0'/'research-only'
    access_level        VARCHAR(20) DEFAULT 'application', -- 'open'/'application'/'restricted'
    paper_url           VARCHAR(500),
    contact_email       VARCHAR(255),
    download_count      INT DEFAULT 0,
    citation_count      INT DEFAULT 0,
    tags                VARCHAR(50)[],
    preview_sample_url  VARCHAR(500),                   -- 预览样本URL
    is_published        BOOLEAN DEFAULT FALSE,
    published_at        TIMESTAMP,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_datasets_type ON research_datasets(dataset_type);
CREATE INDEX idx_datasets_access ON research_datasets(access_level);
CREATE INDEX idx_datasets_published ON research_datasets(is_published);
```

### 5.2 数据集版本表 (dataset_versions)

```sql
CREATE TABLE dataset_versions (
    id                  BIGSERIAL PRIMARY KEY,
    dataset_id          BIGINT NOT NULL REFERENCES research_datasets(id) ON DELETE CASCADE,
    version             VARCHAR(20) NOT NULL,
    changes             TEXT,                             -- 版本变更说明
    size_gb             DECIMAL(10,2),
    record_count        INT,
    download_url        VARCHAR(500),
    storage_path        VARCHAR(500),
    is_latest           BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMP DEFAULT NOW(),
    UNIQUE(dataset_id, version)
);

CREATE INDEX idx_version_dataset ON dataset_versions(dataset_id);
```

### 5.3 数据访问申请表 (dataset_access_requests)

```sql
CREATE TABLE dataset_access_requests (
    id                  BIGSERIAL PRIMARY KEY,
    dataset_id          BIGINT NOT NULL REFERENCES research_datasets(id),
    applicant_id        BIGINT NOT NULL REFERENCES users(id),
    institution         VARCHAR(200) NOT NULL,
    institution_type    VARCHAR(50),                     -- 'university'/'research_institute'/'company'
    department          VARCHAR(100),
    researcher_name    VARCHAR(100) NOT NULL,
    research_purpose   TEXT NOT NULL,
    intended_use        TEXT,                            -- 具体使用计划
    data_security      TEXT NOT NULL,                    -- 数据安全保障承诺
    research_duration  VARCHAR(50),                      -- 研究周期
    expected_outcome   TEXT,                            -- 预期成果
    status              VARCHAR(20) DEFAULT 'pending',    -- 'pending'/'approved'/'rejected'/'expired'
    reviewer_id        BIGINT,
    reviewed_at        TIMESTAMP,
    approval_note       TEXT,
    granted_until      TIMESTAMP,                        -- 授权截止日期
    rejection_reason   TEXT,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_requests_dataset ON dataset_access_requests(dataset_id, status);
CREATE INDEX idx_requests_applicant ON dataset_access_requests(applicant_id);
CREATE INDEX idx_requests_status ON dataset_access_requests(status);
```

### 5.4 数据集下载日志表 (dataset_download_logs)

```sql
CREATE TABLE dataset_download_logs (
    id                  BIGSERIAL PRIMARY KEY,
    dataset_id          BIGINT NOT NULL REFERENCES research_datasets(id),
    version_id          BIGINT REFERENCES dataset_versions(id),
    user_id             BIGINT NOT NULL REFERENCES users(id),
    download_size       BIGINT,
    ip_address          VARCHAR(50),
    user_agent          VARCHAR(500),
    downloaded_at       TIMESTAMP DEFAULT NOW(),
    INDEX idx_dataset_downloads (dataset_id, downloaded_at DESC)
);
```

---

## 六、前端页面

### 6.1 数据集浏览

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 数据集列表 | /research/datasets | 页面顶部右侧 | 「筛选」「搜索」 |
| 数据集详情 | /research/datasets/:id | 内容区左下 | 「申请访问」「收藏」「分享」 |
| 引用生成 | /research/datasets/:id/citation | 详情页左侧 | 「复制引用」 |

### 6.2 申请管理

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 我的申请 | /research/my-applications | 列表顶部 | 「新建申请」「筛选」 |
| 申请详情 | /research/my-applications/:id | 操作区左 | 「撤回」「跟进」 |

### 6.3 管理后台

| 页面 | 路由 | 按钮位置（靠左） | 按钮命名 |
|------|------|-----------------|---------|
| 申请审批列表 | /admin/datasets/applications | 列表顶部 | 「批量审批」「导出」 |
| 申请详情 | /admin/datasets/applications/:id | 操作区左 | 「批准」「拒绝」「联系申请人」 |
| 数据集管理 | /admin/datasets | 列表顶部 | 「创建数据集」「导入」 |
| 数据集编辑 | /admin/datasets/:id/edit | 表单底部 | 「保存」「发布」「添加版本」 |

### 6.4 前端UI规范

- 列表页使用卡片式布局，每个数据集一张卡片
- 详情页左侧为数据集信息，右侧为操作面板
- 申请表单使用步骤式，第一步机构信息，第二步研究目的，第三步承诺确认
- 数据集列表支持按类型、访问级别、标签筛选
- 统计面板使用图表展示下载量、引用量趋势

---

## 七、验收标准

### 7.1 功能验收

| 验收点 | 标准 |
|--------|------|
| 数据集列表 | 按类型筛选正常，搜索响应<500ms |
| 数据集详情 | 完整展示数据集信息，引用格式正确 |
| 访问申请 | 申请提交后3个工作日内完成审批 |
| 下载记录 | 每次下载准确记录，无遗漏 |
| 版本管理 | 版本对比功能正常，历史版本可下载 |

### 7.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 数据集列表加载 | < 1s（100个数据集） |
| 申请审批 | 审批操作响应<200ms |
| 下载链接生成 | < 2s |
| 统计面板 | < 3s（含图表渲染） |

### 7.3 安全验收

| 验收点 | 标准 |
|--------|------|
| 数据隔离 | 未授权用户无法访问数据集文件 |
| 访问控制 | 过期授权自动回收 |
| 操作审计 | 所有操作有完整日志 |

---

## 八、附录

### 8.1 引用格式示例

**APA格式：**
```
MDM Research Team. (2026). [数据集名称] [数据集类型] (Version [版本]). MDM Platform. https://mdm.example.com/datasets/[id]
```

**GB/T 7714格式：**
```
MDM研究团队. [数据集名称][数据集类型][DB/OL]. MDM平台, 2026-03-24. https://mdm.example.com/datasets/[id]
```

### 8.2 许可证类型说明

| 许可证 | 说明 | 商业使用 |
|--------|------|----------|
| CC0 | 公有领域，无任何限制 | ✅ 可用 |
| CC-BY | 署名即可，可商业使用 | ✅ 可用 |
| 研究专用 | 仅限学术研究，不可商用 | ❌ 不可用 |

### 8.3 修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|---------|
| V1.0 | 2026-03-24 | agentcp | 初稿创建 |
