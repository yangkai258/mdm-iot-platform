# Sprint 11 规划

**时间**：2026-04-19
**状态**：✅ 后端已完成（2026-03-22）| 前端已完成
**Sprint 周期**：2 周（2026-04-19 ～ 2026-05-02）

---

## 一、Sprint 目标

**目标：** 实现完整的告警通知渠道

在 Sprint 10（设备监控面板）的基础上，实现完整的告警通知渠道，支持 SMTP 邮件通知、短信通知（SMS）、Webhook 通知，以及告警历史管理功能，确保设备异常能及时触达运营人员。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 | 状态 |
|---|------|------|--------|--------|------|
| P0-1 | **SMTP 邮件通知服务** | 实现邮件发送服务，支持模板渲染 | notification/email_service.go | P0 | ✅ 完成 |
| P0-2 | **短信通知服务** | 实现 SMS 发送服务（阿里云/腾讯云 SDK 集成） | notification/sms_service.go | P0 | ✅ 完成 |
| P0-3 | **Webhook 通知服务** | 实现 HTTP Webhook 推送，支持签名验证 | notification/webhook_service.go | P0 | ✅ 完成 |
| P0-4 | **通知渠道配置 API** | 完成 `/api/v1/notification/channels/*` CRUD | notification_controller.go | P0 | ✅ 完成 |
| P0-5 | **通知日志 API** | 完成 `/api/v1/notification/logs` 查询接口 | notification_controller.go | P0 | ✅ 完成 |
| P1-1 | **告警历史管理 API** | 完成 `/api/v1/alerts/history` 告警历史查询 | alert_history_controller.go | P1 | ✅ 完成 |
| P1-2 | **通知模板管理** | 完成 `/api/v1/notification/templates/*` CRUD | notification_controller.go | P1 | ✅ 完成 |
| P1-3 | **通知重试机制** | 实现失败通知自动重试（指数退避） | notification/retry_worker.go | P1 | ✅ 完成 |
| P1-4 | **告警升级策略** | 告警未处理时自动升级通知级别 | alert_escalation_service.go | P1 | ⏳ 待实现 |
| P2-1 | **通知渠道健康检查** | 实现通知渠道可用性检测 | notification/health_check.go | P2 | ✅ 完成 |
| P2-2 | **通知统计 API** | 完成 `/api/v1/notification/stats` 发送统计 | notification_controller.go | P2 | ✅ 完成 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 | 状态 |
|---|------|------|--------|--------|------|
| PF0-1 | **告警通知配置页面** | 完成 AlertNotificationView.vue 通知渠道配置 | AlertNotificationView.vue | P0 | ✅ 完成 |
| PF0-2 | **邮件渠道配置** | 完成 EmailChannelConfig.vue SMTP 配置表单 | EmailChannelConfig.vue | P0 | ✅ 完成 |
| PF0-3 | **短信渠道配置** | 完成 SMSChannelConfig.vue 短信服务商配置 | SMSChannelConfig.vue | P0 | ✅ 完成 |
| PF0-4 | **Webhook渠道配置** | 完成 WebhookChannelConfig.vue Webhook 配置 | WebhookChannelConfig.vue | P0 | ✅ 完成 |
| PF1-1 | **告警历史管理页面** | 完成 AlertHistoryView.vue 告警历史查询/导出 | AlertHistoryView.vue | P1 | ✅ 完成 |
| PF1-2 | **通知日志页面** | 完成 NotificationLogsView.vue 通知发送日志 | NotificationLogsView.vue | P1 | ✅ 完成 |
| PF1-3 | **通知模板编辑** | 完成 NotificationTemplateEditor.vue 模板编辑 | NotificationTemplateEditor.vue | P1 | ⏳ 待后端 |
| PF2-1 | **通知统计仪表盘** | 完成 NotificationStatsView.vue 发送统计看板 | NotificationStatsView.vue | P2 | ✅ 完成 |

---

## 六、前端交付清单

### ✅ 已完成文件

| 文件路径 | 说明 |
|----------|------|
| `src/api/notification.ts` | 通知相关 API 定义 |
| `src/composables/useNotification.ts` | 通知相关 Composable |
| `src/views/alert/AlertNotificationView.vue` | 告警通知主页面 |
| `src/views/alert/EmailChannelConfig.vue` | 邮件渠道配置 |
| `src/views/alert/SMSChannelConfig.vue` | 短信渠道配置 |
| `src/views/alert/WebhookChannelConfig.vue` | Webhook 配置 |
| `src/views/alert/AlertHistoryView.vue` | 告警历史页面 |
| `src/views/alert/NotificationLogsView.vue` | 通知日志页面 |
| `src/views/alert/NotificationStatsView.vue` | 统计报表页面 |
| `src/router/index.js` | 已添加新路由 |

### 路由配置

```typescript
{
  path: '/alert/notification',
  name: 'AlertNotification',
  component: () => import('@/views/alert/AlertNotificationView.vue')
},
{
  path: '/alert/history',
  name: 'AlertHistory',
  component: () => import('@/views/alert/AlertHistoryView.vue')
},
{
  path: '/alert/notification-logs',
  name: 'NotificationLogs',
  component: () => import('@/views/alert/NotificationLogsView.vue')
}
```

---

## 七、后端交付清单

### ✅ 已完成文件

| 文件路径 | 说明 |
|----------|------|
| `backend/notification/email_service.go` | SMTP 邮件服务（支持 TLS/STARTTLS） |
| `backend/notification/sms_service.go` | SMS 短信服务（支持阿里云/腾讯云） |
| `backend/notification/webhook_service.go` | Webhook 服务（支持 HMAC-SHA256 签名） |
| `backend/notification/retry_worker.go` | 重试机制（指数退避：1s,2s,4s,8s,16s，最大5次） |
| `backend/notification/health_check.go` | 渠道健康检查（支持 SMTP/Webhook/SMS） |
| `backend/models/notification_channel.go` | 通知渠道扩展模型 |
| `backend/models/notification_log.go` | 通知日志模型 + 统计结构体 |
| `backend/models/alert_history.go` | 告警历史模型 |
| `backend/controllers/notification_controller.go` | 新增 `/api/v1/notification/*` 路由 |
| `backend/controllers/alert_history_controller.go` | 告警历史 API |
| `backend/migrations/005_sprint11_notification.sql` | 数据库迁移 SQL |
| `backend/main.go` | 已添加 AlertHistory AutoMigrate + 路由注册 |

### API 路由清单

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/v1/notification/channels` | GET | 通知渠道列表 |
| `/api/v1/notification/channels` | POST | 创建通知渠道 |
| `/api/v1/notification/channels/:id` | GET | 渠道详情 |
| `/api/v1/notification/channels/:id` | PUT | 更新渠道配置 |
| `/api/v1/notification/channels/:id` | DELETE | 删除渠道 |
| `/api/v1/notification/channels/:id/test` | POST | 测试渠道连通性 |
| `/api/v1/notification/logs` | GET | 通知日志（支持多条件筛选） |
| `/api/v1/notification/stats` | GET | 通知发送统计 |
| `/api/v1/notification/templates` | GET/POST | 通知模板 CRUD |
| `/api/v1/notification/templates/:id` | PUT/DELETE | 通知模板更新/删除 |
| `/api/v1/alerts/history` | GET | 告警历史列表 |
| `/api/v1/alerts/history/:id` | GET | 告警历史详情 |
| `/api/v1/alerts/history/archive` | POST | 归档告警到历史表 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `GET /api/v1/notification/channels` | GET | 通知渠道列表 |
| `POST /api/v1/notification/channels` | POST | 创建通知渠道 |
| `GET /api/v1/notification/channels/:id` | GET | 渠道详情 |
| `PUT /api/v1/notification/channels/:id` | PUT | 更新渠道配置 |
| `DELETE /api/v1/notification/channels/:id` | DELETE | 删除渠道 |
| `POST /api/v1/notification/channels/:id/test` | POST | 测试渠道连通性 |
| `GET /api/v1/notification/logs` | GET | 通知日志（支持按渠道/状态/时间筛选） |
| `GET /api/v1/notification/templates` | GET | 通知模板列表 |
| `POST /api/v1/notification/templates` | POST | 创建模板 |
| `PUT /api/v1/notification/templates/:id` | PUT | 更新模板 |
| `DELETE /api/v1/notification/templates/:id` | DELETE | 删除模板 |
| `GET /api/v1/alerts/history` | GET | 告警历史（支持按设备/类型/时间/状态筛选） |
| `GET /api/v1/notification/stats` | GET | 通知发送统计 |

### 数据库设计

```sql
-- 通知渠道表
CREATE TABLE notification_channels (
    id              BIGSERIAL PRIMARY KEY,
    channel_type    VARCHAR(20) NOT NULL,         -- 'email'/'sms'/'webhook'
    channel_name    VARCHAR(100) NOT NULL,
    config          JSONB NOT NULL,               -- 渠道配置（加密存储敏感信息）
    enabled         BOOLEAN DEFAULT TRUE,
    is_default      BOOLEAN DEFAULT FALSE,
    priority        INT DEFAULT 0,
    health_status   VARCHAR(20) DEFAULT 'unknown',
    last_checked_at TIMESTAMP,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 通知模板表
CREATE TABLE notification_templates (
    id              BIGSERIAL PRIMARY KEY,
    template_name   VARCHAR(100) NOT NULL,
    channel_type    VARCHAR(20) NOT NULL,
    template_type   VARCHAR(50) NOT NULL,         -- 'alert'/'reminder'/'system'
    subject_template VARCHAR(255),
    body_template   TEXT NOT NULL,
    variables       JSONB,                          -- 支持的变量列表
    is_active       BOOLEAN DEFAULT TRUE,
    created_by      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 通知日志表
CREATE TABLE notification_logs (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT REFERENCES notification_channels(id),
    channel_type    VARCHAR(20) NOT NULL,
    alert_id        BIGINT,
    recipient       VARCHAR(255),                    -- 邮箱/手机号/URL
    subject         VARCHAR(255),
    body            TEXT,
    status          VARCHAR(20) NOT NULL,         -- 'pending'/'sent'/'failed'
    error_code      VARCHAR(50),
    error_message   TEXT,
    attempt_count   INT DEFAULT 0,
    sent_at         TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_channel_type_status (channel_type, status),
    INDEX idx_alert_id (alert_id),
    INDEX idx_created_at (created_at DESC)
);

-- 告警历史表（归档已解决的告警）
CREATE TABLE alert_history (
    id              BIGSERIAL PRIMARY KEY,
    original_id     BIGINT NOT NULL,
    rule_id         BIGINT,
    device_id       VARCHAR(64) NOT NULL,
    alert_type      VARCHAR(50) NOT NULL,
    severity        INT NOT NULL,
    message         TEXT,
    trigger_value   VARCHAR(100),
    threshold       VARCHAR(100),
    status          INT NOT NULL,
    notified_channels JSONB,                        -- 通知渠道列表
    confirmed_at    TIMESTAMP,
    confirmed_by    BIGINT,
    resolved_at     TIMESTAMP,
    resolved_by     BIGINT,
    resolve_remark  TEXT,
    created_at      TIMESTAMP NOT NULL,
    resolved_at_h   TIMESTAMP,                     -- 解决时间（归档后）
    archived_at     TIMESTAMP DEFAULT NOW(),
    INDEX idx_device_id (device_id),
    INDEX idx_alert_type (alert_type),
    INDEX idx_created_at (created_at DESC)
);
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| SMTP 邮件发送 | 告警触发后 30s 内收到邮件 | 配置邮箱触发测试告警 |
| SMS 短信发送 | 告警触发后 60s 内收到短信 | 配置短信触发测试告警 |
| Webhook 推送 | 告警触发后 HTTP 请求正确送达 | 配置 Webhook 验证请求 |
| 渠道配置 CRUD | 完整增删改查 | 调用各接口验证 |
| 通知重试 | 失败后自动重试 3 次 | 模拟发送失败 |
| 告警历史查询 | 支持多条件筛选和导出 | 分页+筛选测试 |
| 通知模板 | 支持变量替换 | 发送测试模板 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| 邮件发送延迟 | <= 5s |
| 短信发送延迟 | <= 10s |
| Webhook 推送延迟 | <= 3s |
| 通知日志查询 | <= 500ms（10000 条内） |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| Sprint 10 告警规则引擎 | 告警触发通知 |
| 阿里云/腾讯云 SMS SDK | 短信发送 |
| SMTP 服务 | 邮件发送 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 第三方 SMS/邮件服务不可用 | 通知无法送达 | 降级到其他渠道+告警 |
| Webhook 超时阻塞 | 告警处理延迟 | 异步发送+超时配置 |
| 通知日志数据量过大 | DB 存储压力 | 分表+定期归档 |
