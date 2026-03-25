# Sprint 22 规划

**时间**：2026-09-06
**状态**：待开始
**Sprint 周期**：2 周（2026-09-06 ～ 2026-09-19）

---

## 一、Sprint 目标

**目标：** 具身智能完善与多宠物协作

在 Sprint 21 基础上，完善具身智能安全边界、多宠物协作能力，并完成与数字孪生、情感计算模块的联动集成，构建完整的具身智能体系。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **多宠物协作 API** | 完成 `/api/v1/embodied/collaboration/*` 多宠物协作接口 | collaboration_controller.go | P0 |
| P0-2 | **具身AI安全审计** | 完成 `/api/v1/embodied/:device_id/safety/audit` 安全审计接口 | safety_audit_controller.go | P0 |
| P1-1 | **端云协同推理** | 实现边缘端+云端混合推理路由 | inference/edge_cloud_router.go | P1 |
| P1-2 | **具身AI模型压缩** | 实现模型压缩和优化工具链 | inference/compression.go | P1 |
| P1-3 | **具身与数字孪生联动** | 实现具身状态同步到数字孪生 | embodied/digital_twin_sync.go | P1 |
| P2-1 | **动作学习进度可视化 API** | `/api/v1/embodied/:device_id/learning/progress` | learning_controller.go | P2 |
| P2-2 | **具身行为预测** | 基于历史数据预测宠物下一步动作 | embodied/prediction.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **多宠物协作页面** | 完成 MultiPetCollaborationView.vue 协作配置和状态 | MultiPetCollaborationView.vue | P0 |
| PF0-2 | **具身安全审计页面** | 完成 EmbodiedSafetyAuditView.vue 安全日志查看 | EmbodiedSafetyAuditView.vue | P0 |
| PF1-1 | **端云推理配置页面** | 完成 EdgeCloudInferenceView.vue 推理路由配置 | EdgeCloudInferenceView.vue | P1 |
| PF1-2 | **具身-数字孪生联动** | 完成 EmbodiedTwinSyncView.vue 状态同步配置 | EmbodiedTwinSyncView.vue | P1 |
| PF2-1 | **动作学习进度页面** | 完成 LearningProgressView.vue 学习曲线和进度 | LearningProgressView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/embodied/collaboration/devices` | GET | 获取可协作的设备列表 |
| `POST /api/v1/embodied/collaboration/tasks` | POST | 创建协作任务 |
| `GET /api/v1/embodied/collaboration/tasks/:id` | GET | 协作任务详情 |
| `DELETE /api/v1/embodied/collaboration/tasks/:id` | DELETE | 取消协作任务 |
| `GET /api/v1/embodied/:device_id/safety/audit` | GET | 安全审计日志 |
| `GET /api/v1/embodied/:device_id/safety/report` | GET | 生成安全报告 |
| `GET /api/v1/embodied/:device_id/inference/config` | GET | 推理配置 |
| `PUT /api/v1/embodied/:device_id/inference/config` | PUT | 更新推理配置 |
| `POST /api/v1/embodied/:device_id/model/compress` | POST | 触发模型压缩 |
| `GET /api/v1/embodied/:device_id/learning/progress` | GET | 学习进度 |
| `GET /api/v1/embodied/:device_id/prediction` | GET | 行为预测结果 |

### 数据库设计

```sql
-- 多宠物协作任务表
CREATE TABLE embodied_collaboration_tasks (
    id              BIGSERIAL PRIMARY KEY,
    task_type       VARCHAR(50) NOT NULL,            -- 'exploration'/'play'/'patrol'
    initiator_device VARCHAR(100) NOT NULL,
    participant_devices VARCHAR(100)[] NOT NULL,
    target_position JSONB,                            -- 目标位置
    task_status     VARCHAR(20) DEFAULT 'pending',   -- 'pending'/'running'/'completed'/'failed'
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_task_status (task_status)
);

-- 协作交互记录表
CREATE TABLE embodied_collaboration_logs (
    id              BIGSERIAL PRIMARY KEY,
    task_id         BIGINT NOT NULL REFERENCES embodied_collaboration_tasks(id),
    device_id       VARCHAR(100) NOT NULL,
    action_type     VARCHAR(50),
    action_data     JSONB,
    synchronized    BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_task_logs (task_id)
);

-- 推理配置表
CREATE TABLE embodied_inference_configs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL UNIQUE,
    inference_mode  VARCHAR(20) DEFAULT 'auto',     -- 'edge'/'cloud'/'auto'
    cloud_endpoint  VARCHAR(500),
    fallback_enabled BOOLEAN DEFAULT TRUE,
    latency_threshold_ms INT DEFAULT 500,
    model_version   VARCHAR(50),
    compression_level VARCHAR(20) DEFAULT 'medium',  -- 'low'/'medium'/'high'
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 具身安全审计表
CREATE TABLE embodied_safety_audits (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL,
    audit_type      VARCHAR(50) NOT NULL,            -- 'daily'/'weekly'/'incident'
    report_data     JSONB NOT NULL,                   -- 审计报告数据
    overall_score   DECIMAL(5,2),                    -- 安全评分 0-100
    issues_found    INT DEFAULT 0,
    generated_at    TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_audits (device_id, generated_at DESC)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 多宠物协作 | 2台设备协作任务成功率>85% | 协作测试场景 |
| 安全审计 | 每日自动生成安全报告，报告完整率100% | 报告验证 |
| 端云协同 | 边缘优先，延迟>阈值自动切换云端，切换<1s | 延迟测试 |
| 模型压缩 | 压缩后模型体积减少>50%，准确率损失<5% | 精度对比测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 协作任务同步延迟 | < 500ms |
| 云端推理延迟 | < 2s（含网络） |
| 安全报告生成 | < 30s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 21 具身核心 | 环境感知、空间认知基础能力 |
| 云端推理服务 | 云端模型推理服务 |
| 设备端固件 | 具身协作协议支持 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 多设备通信不稳定 | 协作失败 | 本地优先+重试机制 |
| 云端推理服务故障 | 无法切换云端 | 边缘降级运行 |
| 安全报告泄露隐私 | 隐私合规风险 | 数据脱敏处理 |
