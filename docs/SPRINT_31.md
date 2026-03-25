# Sprint 31 规划

**时间**：2026-12-06
**状态**：待开始
**Sprint 周期**：2 周（2026-12-06 ～ 2026-12-19）

---

## 一、Sprint 目标

**目标：** 端侧推理完善与数据集开放

完善端侧推理能力（模型分片加载、动态更新）、BLE Mesh组网完善、RTOS深度优化，并开放学术研究数据集，构建MDM平台的持续演进能力。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **端侧推理完善** | 完善设备端模型推理，延迟<100ms | edge_inference_engine.go | P0 |
| P0-2 | **模型分片加载完善** | 实现动态分片和增量更新 | model_sharding.go | P0 |
| P0-3 | **BLE Mesh完善** | 实现宠物间直接通信，Mesh组网稳定 | ble_mesh_enhanced.go | P0 |
| P1-1 | **RTOS深度优化** | 实现实时性能优化，性能提升>30% | rtos_optimizer.go | P1 |
| P1-2 | **数据集开放平台** | 实现 `/api/v1/datasets/*` 学术数据集开放 | dataset_platform_controller.go | P1 |
| P1-3 | **数据访问控制** | 实现数据集访问权限和审批 | data_access_control.go | P1 |
| P2-1 | **数据集管理后台** | 实现数据集上传、标注、版本管理 | dataset_admin.go | P2 |
| P2-2 | **数据集统计分析** | 数据集使用统计和影响力追踪 | dataset_analytics.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **端侧推理监控页面** | 完成 EdgeInferenceView.vue 推理状态监控 | EdgeInferenceView.vue | P0 |
| PF0-2 | **BLE Mesh拓扑页面** | 完成 BLEMeshTopologyView.vue Mesh网络可视化 | BLEMeshTopologyView.vue | P0 |
| PF1-1 | **数据集开放平台页面** | 完成 DatasetPlatformView.vue 数据集浏览和申请 | DatasetPlatformView.vue | P1 |
| PF1-2 | **数据集管理页面** | 完成 DatasetAdminView.vue 数据集上传和管理 | DatasetAdminView.vue | P1 |
| PF2-1 | **RTOS性能监控页面** | 完成 RTOSPerfView.vue 性能优化监控 | RTOSPerfView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/edge-inference/status` | GET | 端侧推理状态 |
| `GET /api/v1/edge-inference/models` | GET | 已加载模型列表 |
| `POST /api/v1/edge-inference/models/:id/shard` | POST | 加载模型分片 |
| `GET /api/v1/ble-mesh/topology` | GET | Mesh拓扑图 |
| `POST /api/v1/ble-mesh/devices/:id/connect` | POST | 连接Mesh设备 |
| `GET /api/v1/datasets` | GET | 数据集列表 |
| `GET /api/v1/datasets/:id` | GET | 数据集详情 |
| `POST /api/v1/datasets/:id/access-request` | POST | 申请数据访问 |
| `GET /api/v1/datasets/:id/download` | GET | 下载数据集 |
| `POST /api/v1/datasets` | POST | 上传数据集（管理员） |
| `GET /api/v1/datasets/analytics` | GET | 数据集统计 |

### 数据库设计

```sql
-- 数据集表
CREATE TABLE research_datasets (
    id              BIGSERIAL PRIMARY KEY,
    dataset_name    VARCHAR(200) NOT NULL,
    dataset_type    VARCHAR(50) NOT NULL,            -- 'behavior'/'emotion'/'health'/'voice'
    description     TEXT,
    size_gb         DECIMAL(10,2),
    record_count    INT,
    data_format     VARCHAR(30),                     -- 'json'/'csv'/'parquet'
    download_url    VARCHAR(500),
    license_type    VARCHAR(50),                     -- 'cc-by'/'cc0'/'research-only'
    access_level    VARCHAR(20) DEFAULT 'application', -- 'open'/'application'/'restricted'
    paper_url       VARCHAR(500),
    contact_email   VARCHAR(255),
    download_count  INT DEFAULT 0,
    citation_count  INT DEFAULT 0,
    is_published    BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 数据访问申请表
CREATE TABLE dataset_access_requests (
    id              BIGSERIAL PRIMARY KEY,
    dataset_id      BIGINT NOT NULL REFERENCES research_datasets(id),
    applicant_id    BIGINT NOT NULL REFERENCES users(id),
    institution     VARCHAR(200),
    research_purpose TEXT NOT NULL,
    intended_use    TEXT,
    data_security   VARCHAR(100),                    -- 数据安全保障承诺
    status          VARCHAR(20) DEFAULT 'pending',    -- 'pending'/'approved'/'rejected'
    reviewed_by     BIGINT,
    reviewed_at     TIMESTAMP,
    approval_note   TEXT,
    granted_until   TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_dataset_requests (dataset_id, status)
);

-- 数据集版本表
CREATE TABLE dataset_versions (
    id              BIGSERIAL PRIMARY KEY,
    dataset_id      BIGINT NOT NULL REFERENCES research_datasets(id),
    version         VARCHAR(20) NOT NULL,
    changes         TEXT,
    size_gb         DECIMAL(10,2),
    record_count    INT,
    download_url    VARCHAR(500),
    is_latest       BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 端侧推理 | 推理延迟<100ms，隐私数据不出设备 | 延迟测试 |
| 模型分片加载 | 分片加载时间<5s，增量更新正常 | 加载测试 |
| BLE Mesh | 宠物间直接通信成功率>95% | Mesh测试 |
| RTOS优化 | 性能提升>30%，响应延迟降低 | 性能基准测试 |
| 数据集开放 | 数据集发布/申请/审批流程完整 | 全流程测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 端侧推理延迟 | < 100ms |
| BLE Mesh通信 | < 200ms |
| 数据集下载 | < 60s（1GB） |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| MODULE_MINICLAW_FIRMWARE | 固件优化能力 |
| MODULE_AI_ENGINEERING | 模型分片加载基础 |
| 学术合作机构 | 数据集来源 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 端侧推理性能瓶颈 | 用户体验差 | 持续优化模型压缩 |
| 数据集隐私泄露 | 法律风险 | 严格数据脱敏+审批 |
| BLE Mesh不稳定 | 设备通信失败 | 多路径冗余+重试 |
