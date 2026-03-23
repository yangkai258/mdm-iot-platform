# PRD：Webhook市场

**版本：** V1.0
**所属Phase：** Phase 4（Sprint 25-26）
**优先级：** P2
**负责角色：** agentcp（产品）、agenthd（后端）、agentqd（前端）

---

## 一、概述

### 1.1 模块定位

Webhook市场为开发者提供预置的Webhook模板库，支持第三方系统（Slack、钉钉、飞书、自定义）的事件订阅配置，实现MDM平台与外部系统的自动化事件驱动集成。

### 1.2 核心价值

- **开箱即用**：提供预置的Slack/钉钉/飞书模板
- **事件驱动**：设备告警/状态变更/订阅事件自动推送
- **可视化配置**：非技术人员也可配置Webhook

---

## 二、功能详情

### 2.1 预置模板

| 模板 | 触发事件 | 说明 |
|------|----------|------|
| Slack告警 | 设备告警触发 | 告警内容推送到Slack频道 |
| 钉钉群通知 | 设备告警/订阅变更 | 告警卡片推送到钉钉群 |
| 飞书群通知 | 设备告警/订阅变更 | 告警消息推送到飞书群 |
| 企业微信 | 设备告警/订阅变更 | 消息推送到企业微信 |
| 自定义Webhook | 任意事件 | 支持自定义URL和payload模板 |

### 2.2 事件类型

| 事件类型 | 触发条件 |
|----------|----------|
| device.alert | 设备告警触发 |
| device.online | 设备上线 |
| device.offline | 设备离线 |
| device.ota_progress | OTA升级进度更新 |
| subscription.created | 新订阅创建 |
| subscription.expired | 订阅到期 |
| pet.emotion_changed | 宠物情绪变化 |

### 2.3 Webhook配置

| 功能 | 说明 |
|------|------|
| 模板选择 | 从预置模板中选择 |
| URL配置 | 输入目标Webhook URL |
| 事件选择 | 选择要订阅的事件类型 |
| 认证配置 | 支持Secret签名认证 |
| 测试推送 | 发送测试事件验证配置 |

### 2.4 Webhook管理

| 功能 | 说明 |
|------|------|
| Webhook列表 | 展示所有已配置的Webhook |
| 启用/禁用 | 临时禁用Webhook |
| 编辑配置 | 修改Webhook配置 |
| 删除 | 删除Webhook |
| 推送日志 | 查看推送历史和状态 |

---

## 三、API接口定义

### 3.1 模板管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/webhook-market/templates | 预置模板列表 |
| GET | /api/v1/webhook-market/templates/:id | 模板详情 |

### 3.2 Webhook配置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/webhook-market/webhooks | Webhook列表 |
| POST | /api/v1/webhook-market/webhooks | 创建Webhook |
| GET | /api/v1/webhook-market/webhooks/:id | Webhook详情 |
| PUT | /api/v1/webhook-market/webhooks/:id | 更新Webhook |
| DELETE | /api/v1/webhook-market/webhooks/:id | 删除Webhook |
| POST | /api/v1/webhook-market/webhooks/:id/test | 发送测试推送 |

### 3.3 推送日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/webhook-market/webhooks/:id/logs | 推送日志 |
| POST | /api/v1/webhook-market/webhooks/:id/retry | 重试失败推送 |

---

## 四、数据库设计

### 4.1 Webhook模板表 (webhook_market_templates)

```sql
CREATE TABLE webhook_market_templates (
    id              BIGSERIAL PRIMARY KEY,
    template_name   VARCHAR(255) NOT NULL,
    platform        VARCHAR(50) NOT NULL,            -- 'slack'/'dingtalk'/'feishu'/'wecom'/'custom'
    description     TEXT,
    icon_url        VARCHAR(500),
    payload_template JSONB NOT NULL,                   -- 预置的payload模板
    event_types     VARCHAR(50)[],                   -- 支持的事件类型
    config_schema   JSONB,                           -- 配置字段定义
    is_official     BOOLEAN DEFAULT TRUE,
    usage_count     INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 4.2 Webhook配置表 (webhook_market_configs)

```sql
CREATE TABLE webhook_market_configs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    template_id     BIGINT REFERENCES webhook_market_templates(id),
    webhook_name    VARCHAR(255) NOT NULL,
    target_url      VARCHAR(500) NOT NULL,
    secret_key      VARCHAR(255),
    event_types     VARCHAR(50)[],                    -- 订阅的事件类型
    headers         JSONB,                            -- 自定义请求头
    is_active       BOOLEAN DEFAULT TRUE,
    retry_count     INT DEFAULT 3,
    timeout_seconds INT DEFAULT 30,
    last_triggered_at TIMESTAMP,
    last_status     VARCHAR(20),                       -- 'success'/'failed'
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_webhook_configs_user ON webhook_market_configs(user_id);
CREATE INDEX idx_webhook_configs_active ON webhook_market_configs(is_active);
```

### 4.3 推送日志表 (webhook_market_logs)

```sql
CREATE TABLE webhook_market_logs (
    id              BIGSERIAL PRIMARY KEY,
    webhook_id      BIGINT NOT NULL REFERENCES webhook_market_configs(id),
    event_type      VARCHAR(50) NOT NULL,
    payload         JSONB,
    response_status INT,
    response_body   TEXT,
    attempt_count   INT DEFAULT 1,
    status          VARCHAR(20) NOT NULL,              -- 'pending'/'success'/'failed'
    error_message   TEXT,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_webhook_logs_webhook ON webhook_market_logs(webhook_id, created_at DESC);
CREATE INDEX idx_webhook_logs_status ON webhook_market_logs(status);
```

---

## 五、前端页面

### 5.1 市场

| 页面 | 路由 | 说明 |
|------|------|------|
| 模板市场 | /webhook-market/templates | 预置模板浏览 |
| 模板详情 | /webhook-market/templates/:id | 模板详情和快速配置 |

### 5.2 Webhook管理

| 页面 | 路由 | 说明 |
|------|------|------|
| Webhook列表 | /webhook-market | 我的Webhook列表 |
| 创建Webhook | /webhook-market/create | 从模板创建Webhook |
| Webhook详情 | /webhook-market/:id | 查看配置和日志 |

---

## 六、验收标准

| 验收点 | 标准 |
|--------|------|
| 模板数量 | 至少5个预置模板（Slack/钉钉/飞书/企业微信/自定义） |
| 推送成功率 | 推送成功率>95% |
| 推送延迟 | 事件触发到推送<5秒 |
| 失败重试 | 失败自动重试3次 |
| 日志保留 | 推送日志保留30天 |
