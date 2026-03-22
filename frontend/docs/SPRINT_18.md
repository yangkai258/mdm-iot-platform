# Sprint 18 规划

**时间**：2026-07-26
**状态**：待开始
**Sprint 周期**：2 周（2026-07-26 ～ 2026-08-08）

---

## 一、Sprint 目标

**目标：** 实时体征和历史回放

在 Sprint 17（情感计算）的基础上，实现宠物数字孪生的核心功能，包括实时生命体征同步、行为预测、历史回放，为宠物健康管理提供数据基础。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **实时生命体征 API** | 完成 `/api/v1/digital-twin/:pet_id/vitals` 实时体征 | vitals_controller.go | P0 |
| P0-2 | **生命体征数据库** | 创建 vital_records + behavior_events 表 | models/vital_record.go | P0 |
| P0-3 | **历史回放 API** | 完成 `/api/v1/digital-twin/:pet_id/timeline` 时间轴 | timeline_controller.go | P0 |
| P0-4 | **行为预测 API** | 完成 `/api/v1/digital-twin/:pet_id/predictions` | prediction_controller.go | P0 |
| P0-5 | **体征上报处理** | 实现 `/api/v1/digital-twin/:pet_id/vitals/report` 接收设备数据 | vitals_report_handler.go | P0 |
| P1-1 | **精彩瞬间生成** | 完成 `/api/v1/digital-twin/:pet_id/memories` 回忆集 | highlight_moments_controller.go | P1 |
| P1-2 | **跨设备状态同步** | 完成 `/api/v1/digital-twin/:pet_id/sync/*` 同步 | sync_controller.go | P1 |
| P1-3 | **健康预警接口** | 完成 `/api/v1/digital-twin/:pet_id/alerts` 健康预警 | health_alert_controller.go | P1 |
| P2-1 | **历史数据导出** | 完成 `/api/v1/digital-twin/:pet_id/export` 数据导出 | export_controller.go | P2 |
| P2-2 | **体征趋势分析** | 完成体征趋势图表数据接口 | vitals_trend_controller.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **生命体征仪表盘** | 完成 VitalsDashboardView.vue 体征总览 | VitalsDashboardView.vue | P0 |
| PF0-2 | **实时体征曲线** | 完成 RealTimeVitalsChart.vue 实时心跳/呼吸曲线 | RealTimeVitalsChart.vue | P0 |
| PF0-3 | **历史回放前端** | 完成 HistoricalReplayView.vue 时间轴回放 | HistoricalReplayView.vue | P0 |
| PF0-4 | **行为预测展示** | 完成 BehaviorPredictionView.vue 预测结果展示 | BehaviorPredictionView.vue | P0 |
| PF1-1 | **健康预警前端** | 完成 HealthAlertView.vue 健康预警列表 | HealthAlertView.vue | P1 |
| PF1-2 | **回忆集展示** | 完成 MemoryCollectionView.vue AI 回忆集 | MemoryCollectionView.vue | P1 |
| PF1-3 | **精彩瞬间相册** | 完成 HighlightMomentsView.vue 精彩瞬间 | HighlightMomentsView.vue | P1 |
| PF2-1 | **数据导出配置** | 完成 DataExportView.vue 导出格式配置 | DataExportView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/digital-twin/:pet_id/vitals` | GET | 获取实时生命体征 |
| `GET /api/v1/digital-twin/:pet_id/vitals/history` | GET | 生命体征历史 |
| `GET /api/v1/digital-twin/:pet_id/vitals/heartbeat` | GET | 心跳曲线 |
| `GET /api/v1/digital-twin/:pet_id/vitals/respiration` | GET | 呼吸曲线 |
| `GET /api/v1/digital-twin/:pet_id/vitals/temperature` | GET | 体温曲线 |
| `POST /api/v1/digital-twin/:pet_id/vitals/report` | POST | 设备上报体征数据 |
| `GET /api/v1/digital-twin/:pet_id/predictions` | GET | 获取行为预测 |
| `GET /api/v1/digital-twin/:pet_id/predictions/short-term` | GET | 短期动作预测 |
| `GET /api/v1/digital-twin/:pet_id/predictions/intent` | GET | 意图识别结果 |
| `GET /api/v1/digital-twin/:pet_id/alerts` | GET | 健康预警列表 |
| `POST /api/v1/digital-twin/:pet_id/alerts/:id/ack` | POST | 确认预警 |
| `POST /api/v1/digital-twin/:pet_id/alerts/:id/ignore` | POST | 忽略预警 |
| `GET /api/v1/digital-twin/:pet_id/timeline` | GET | 时间轴事件 |
| `GET /api/v1/digital-twin/:pet_id/events/:event_id` | GET | 事件详情 |
| `GET /api/v1/digital-twin/:pet_id/memories` | GET | 回忆集 |
| `GET /api/v1/digital-twin/:pet_id/memories/:id/highlights` | GET | 精彩瞬间 |
| `POST /api/v1/digital-twin/:pet_id/export` | POST | 导出数据 |
| `GET /api/v1/digital-twin/:pet_id/sync/status` | GET | 同步状态 |
| `POST /api/v1/digital-twin/:pet_id/sync/pull` | POST | 拉取最新状态 |
| `POST /api/v1/digital-twin/:pet_id/sync/push` | POST | 推送本地变更 |

### 数据库设计

```sql
-- 宠物数字孪生表
CREATE TABLE digital_twin_pets (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    current_state   JSONB,
    last_sync_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(pet_id)
);

-- 生命体征记录表
CREATE TABLE vital_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    vital_type      VARCHAR(30) NOT NULL,         -- 'heartbeat'/'respiration'/'temperature'/'activity'
    value           JSONB NOT NULL,
    recorded_at     TIMESTAMP NOT NULL,
    device_id       VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_pet_vital_time (pet_id, vital_type, recorded_at DESC)
);

-- 行为事件表
CREATE TABLE behavior_events (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    event_type      VARCHAR(50) NOT NULL,         -- 'eating'/'sleeping'/'playing'/'walking'
    event_data      JSONB,
    start_time      TIMESTAMP NOT NULL,
    end_time        TIMESTAMP,
    duration_seconds INT,
    confidence      DECIMAL(5,4),
    device_id       VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_pet_time (pet_id, start_time DESC)
);

-- 健康预警表
CREATE TABLE health_alerts (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    alert_type      VARCHAR(50) NOT NULL,         -- 'fever'/'heart_rate'/'breathing'/'behavior'
    severity        VARCHAR(20) NOT NULL,           -- 'info'/'warning'/'critical'
    trigger_value   VARCHAR(100),
    threshold       VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'pending',
    suggestion      TEXT,
    acknowledged_at TIMESTAMP,
    acknowledged_by BIGINT REFERENCES users(id),
    ignore_reason   TEXT,
    resolved_at     TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 精彩瞬间表
CREATE TABLE highlight_moments (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    moment_type     VARCHAR(30) NOT NULL,         -- 'emotion_peak'/'interaction'/'rare_action'
    media_url       VARCHAR(500),
    thumbnail_url   VARCHAR(500),
    description     TEXT,
    emotion_score   DECIMAL(5,4),
    captured_at     TIMESTAMP NOT NULL,
    is_auto_generated BOOLEAN DEFAULT TRUE,
    is_exported     BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 同步记录表
CREATE TABLE sync_records (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    device_id       VARCHAR(100) NOT NULL,
    sync_type       VARCHAR(20) NOT NULL,         -- 'push'/'pull'/'batch'
    sync_data       JSONB,
    sync_status     VARCHAR(20) NOT NULL,
    records_synced  INT DEFAULT 0,
    error_message   TEXT,
    synced_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 体征延迟 | 设备上报到展示延迟 < 2 秒 | 计时测试 |
| 心跳曲线 | 最近 24 小时曲线正确展示 | 调用 API 验证 |
| 体征准确性 | 与设备传感器数据偏差在允许范围内 | 对比测试 |
| 短期动作准确率 | > 75% | 标注数据对比 |
| 意图识别准确率 | > 70% | 标注数据对比 |
| 预警触发 | 符合条件时 5 分钟内发出预警 | 模拟测试 |
| 数据完整性 | 历史事件记录不丢失 | 调用 API 验证 |
| 查询性能 | 1000 条记录内查询 < 1 秒 | 性能测试 |
| 同步延迟 | 状态变更到其他设备 < 5 秒 | 计时测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 体征写入吞吐量 | >= 500 records/s |
| 历史查询延迟 | <= 500ms |
| 体征曲线渲染数据生成 | <= 200ms |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 10 传感器数据 | 体征数据来源 |
| Sprint 14 AI 预测 | 行为预测模型 |
| Sprint 17 情感计算 | 情绪指数数据 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 体征数据量大 | DB 存储压力 | 分表+TTL+压缩 |
| 设备离线 | 体征中断 | 本地缓存+断网续传 |
| 预测模型延迟 | 预测结果不及时 | 模型优化+边缘推理 |
