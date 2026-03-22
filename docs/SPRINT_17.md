# Sprint 17 规划

**时间**：2026-07-12
**状态**：待开始
**Sprint 周期**：2 周（2026-07-12 ～ 2026-07-25）

---

## 一、Sprint 目标

**目标：** 宠物情绪识别和响应

在 Sprint 16（商业化）的基础上，实现宠物情绪识别和响应功能，包括多模态情绪识别（语音/表情/行为）、情绪响应策略、情绪日志管理，构建有温度的情感陪伴体验。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **情绪识别 API** | 完成 `/api/v1/emotion/recognize/*` 多模态识别 | emotion_recognize_controller.go | P0 |
| P0-2 | **情绪识别数据库** | 创建 emotion_records 表存储情绪数据 | models/emotion_record.go | P0 |
| P0-3 | **情绪响应 API** | 完成 `/api/v1/emotion/:pet_id/response` 响应接口 | emotion_response_controller.go | P0 |
| P0-4 | **情绪响应配置** | 完成 `/api/v1/emotion/:pet_id/response-config` 配置 | emotion_config_controller.go | P0 |
| P0-5 | **情绪日志 API** | 完成 `/api/v1/emotion/logs` 日志查询接口 | emotion_log_controller.go | P0 |
| P1-1 | **情绪动作库** | 完成 `/api/v1/emotion/action-library` CRUD | emotion_action_controller.go | P1 |
| P1-2 | **情绪报告生成** | 完成 `/api/v1/emotion/:pet_id/reports/*` 日报/周报/月报 | emotion_report_controller.go | P1 |
| P1-3 | **家庭情绪地图 API** | 完成 `/api/v1/emotion/family/:family_id/*` | family_emotion_controller.go | P1 |
| P2-1 | **情绪趋势分析** | 完成情绪趋势图表数据接口 | emotion_trend_controller.go | P2 |
| P2-2 | **情绪预警通知** | 情绪异常时自动通知用户 | emotion/alert_service.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **情绪识别配置前端** | 完成 EmotionRecognizeView.vue 情绪识别面板 | EmotionRecognizeView.vue | P0 |
| PF0-2 | **情绪日志查看页面** | 完成 EmotionLogView.vue 情绪历史记录 | EmotionLogView.vue | P0 |
| PF0-3 | **情绪响应配置页面** | 完成 EmotionResponseConfigView.vue 响应策略配置 | EmotionResponseConfigView.vue | P0 |
| PF0-4 | **情绪报告查看页面** | 完成 EmotionReportView.vue 日报/周报/月报 | EmotionReportView.vue | P0 |
| PF1-1 | **情绪实时面板** | 完成 EmotionRealTimePanel.vue 实时情绪展示 | EmotionRealTimePanel.vue | P1 |
| PF1-2 | **情绪动作库页面** | 完成 EmotionActionLibraryView.vue 动作库管理 | EmotionActionLibraryView.vue | P1 |
| PF2-1 | **家庭情绪地图前端** | 完成 FamilyEmotionMapView.vue 家庭情绪可视化 | FamilyEmotionMapView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/v1/emotion/recognize/voice` | POST | 语音情绪识别 |
| `POST /api/v1/emotion/recognize/text` | POST | 文字情绪识别 |
| `POST /api/v1/emotion/recognize/pet` | POST | 宠物情绪识别 |
| `POST /api/v1/emotion/recognize/batch` | POST | 批量情绪识别 |
| `GET /api/v1/emotion/:pet_id/response-strategies` | GET | 获取响应策略 |
| `POST /api/v1/emotion/:pet_id/response` | POST | 触发情绪响应 |
| `PUT /api/v1/emotion/:pet_id/response-config` | PUT | 更新响应配置 |
| `GET /api/v1/emotion/action-library` | GET | 动作库列表 |
| `POST /api/v1/emotion/action-library` | POST | 添加自定义动作 |
| `PUT /api/v1/emotion/action-library/:id` | PUT | 更新动作 |
| `DELETE /api/v1/emotion/action-library/:id` | DELETE | 删除动作 |
| `GET /api/v1/emotion/logs` | GET | 情绪日志列表 |
| `GET /api/v1/emotion/logs/:id` | GET | 日志详情 |
| `PUT /api/v1/emotion/logs/:id` | PUT | 更新日志（添加备注） |
| `DELETE /api/v1/emotion/logs/:id` | DELETE | 删除日志 |
| `GET /api/v1/emotion/:pet_id/reports/daily` | GET | 每日情绪报告 |
| `GET /api/v1/emotion/:pet_id/reports/weekly` | GET | 每周情绪报告 |
| `GET /api/v1/emotion/:pet_id/reports/monthly` | GET | 每月情绪报告 |
| `GET /api/v1/emotion/family/:family_id/map` | GET | 家庭情绪地图 |
| `GET /api/v1/emotion/family/:family_id/trends` | GET | 家庭情绪趋势 |
| `GET /api/v1/emotion/family/:family_id/suggestions` | GET | 家庭情绪建议 |

### 数据库设计

```sql
-- 情绪记录表
CREATE TABLE emotion_records (
    id              BIGSERIAL PRIMARY KEY,
    subject_type    VARCHAR(20) NOT NULL,         -- 'user'/'pet'
    subject_id      BIGINT NOT NULL,
    emotion_type    VARCHAR(30) NOT NULL,         -- 'happy'/'sad'/'angry'/'surprised'/'calm'/'anxious'/'lonely'/'tired'
    intensity       DECIMAL(5,2) NOT NULL,         -- 情绪强度 0-100
    source          VARCHAR(20) NOT NULL,         -- 'voice'/'text'/'face'/'behavior'
    confidence      DECIMAL(5,4),
    context         JSONB,
    trigger_event   TEXT,
    tags            VARCHAR(100)[],
    note            TEXT,
    recorded_at     TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_subject (subject_type, subject_id, recorded_at DESC),
    INDEX idx_emotion_type (emotion_type, recorded_at DESC)
);

-- 宠物情绪动作表
CREATE TABLE pet_emotion_actions (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    emotion_type    VARCHAR(30) NOT NULL,
    action_type     VARCHAR(30) NOT NULL,         -- 'body'/'sound'/'expression'/'behavior'
    action_name     VARCHAR(100) NOT NULL,
    action_data     JSONB,
    priority        INT DEFAULT 0,
    min_intensity   DECIMAL(5,2) DEFAULT 0,
    max_intensity   DECIMAL(5,2) DEFAULT 100,
    is_custom       BOOLEAN DEFAULT FALSE,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 情绪响应配置表
CREATE TABLE emotion_response_configs (
    id              BIGSERIAL PRIMARY KEY,
    pet_id          BIGINT NOT NULL REFERENCES pets(id),
    sensitivity     DECIMAL(3,2) DEFAULT 0.5,
    response_delay_ms INT DEFAULT 0,
    response_freq_limit INT DEFAULT 10,
    disabled_actions VARCHAR(50)[],
    enabled         BOOLEAN DEFAULT TRUE,
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 家庭情绪表
CREATE TABLE family_emotions (
    id              BIGSERIAL PRIMARY KEY,
    family_id       BIGINT NOT NULL REFERENCES families(id),
    recorded_at     TIMESTAMP NOT NULL,
    overall_index   DECIMAL(5,2),
    dominant_emotion VARCHAR(30),
    member_states   JSONB,
    influence_map   JSONB,
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(family_id, recorded_at)
);

-- 情绪报告表
CREATE TABLE emotion_reports (
    id              BIGSERIAL PRIMARY KEY,
    subject_type    VARCHAR(20) NOT NULL,
    subject_id      BIGINT NOT NULL,
    report_type     VARCHAR(20) NOT NULL,         -- 'daily'/'weekly'/'monthly'
    period_start    TIMESTAMP NOT NULL,
    period_end      TIMESTAMP NOT NULL,
    summary         JSONB,
    emotion_trends  JSONB,
    anomaly_events  JSONB,
    suggestions     TEXT,
    generated_at    TIMESTAMP DEFAULT NOW()
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 语音情绪识别 | 准确率 > 85% | 标注数据对比测试 |
| 文字情绪识别 | 准确率 > 88% | 标注数据对比测试 |
| 宠物表情识别 | 准确率 > 80% | 标注数据对比测试 |
| 识别延迟 | 语音 < 500ms，文字 < 200ms | 计时测试 |
| 响应触发 | 情绪识别后 5 秒内触发响应 | 计时测试 |
| 响应准确率 | 响应动作与情绪匹配度 > 80% | 人工评估 |
| 响应频率限制 | 不超过配置的频率限制 | 调用 API 验证 |
| 禁用生效 | 禁用的动作类型不触发 | 配置禁用验证 |
| 情绪日志 | 所有情绪变化自动记录 | 调用 API 验证 |
| 报告生成 | 日报/周报/月报按时自动生成 | 定时任务验证 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 情绪识别吞吐量 | >= 100 requests/s |
| 日志查询 | <= 500ms（1000 条内） |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 9 宠物基础 | 宠物模型 |
| AI 模型服务 | 情绪识别 AI 模型 |
| Sprint 14 AI 监控 | AI 行为埋点 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 情绪识别准确率不足 | 响应不恰当 | 持续迭代模型 |
| 隐私问题 | 用户对情绪监控有顾虑 | 明确告知+数据保护 |
| 情绪日志数据量大 | DB 存储压力 | 分表+定期归档 |
