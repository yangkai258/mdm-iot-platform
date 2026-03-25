# Sprint 25 规划

**时间**：2026-09-27
**状态**：待开始
**Sprint 周期**：2 周（2026-09-27 ～ 2026-10-10）

---

## 一、Sprint 目标

**目标：** 开放平台核心 - 开发者生态与SDK

构建完整的开放平台开发者生态，提供REST API + GraphQL开发者API、完善SDK发布（iOS/Android/小程序/Web）、Webhook市场，实现第三方开发者的快速接入。

---

## 二、详细任务列表

### 后端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| P0-1 | **开发者门户 API** | 完成 `/api/v1/developer/*` 开发者注册/认证/应用管理 | developer_controller.go | P0 |
| P0-2 | **API Key管理** | 实现API Key生成/轮换/吊销 | apikey_service.go | P0 |
| P0-3 | **SDK服务端适配** | 确保现有API兼容多端SDK | api_sdk_compatibility.go | P0 |
| P0-4 | **Webhook市场 API** | 完成 `/api/v1/webhook-market/*` Webhook模板和配置 | webhook_market_controller.go | P0 |
| P1-1 | **SDK发布包构建** | 构建 iOS/Android/小程序/Web SDK发布包 | sdk_builder.go | P1 |
| P1-2 | **API限流与配额** | 实现开发者API限流和配额管理 | rate_limiter.go | P1 |
| P1-3 | **开发者认证** | 实现个人/企业开发者认证 | developer_verification.go | P1 |
| P2-1 | **API使用统计** | 实现API调用统计和报表 | api_usage_stats.go | P2 |
| P2-2 | **SDK文档站点** | 自动生成SDK文档 | sdk_docs_generator.go | P2 |

### 前端 P0/P1/P2 任务表

| # | 任务 | 说明 | 交付物 | 优先级 |
|---|------|------|--------|--------|
| PF0-1 | **开发者门户** | 完成 DeveloperPortalView.vue 开发者注册和首页 | DeveloperPortalView.vue | P0 |
| PF0-2 | **应用管理页面** | 完成 AppManageView.vue 创建和管理应用 | AppManageView.vue | P0 |
| PF0-3 | **API Key管理页面** | 完成 APIKeyManageView.vue Key的生成和配置 | APIKeyManageView.vue | P0 |
| PF0-4 | **Webhook市场页面** | 完成 WebhookMarketView.vue 模板浏览和配置 | WebhookMarketView.vue | P0 |
| PF1-1 | **SDK下载页面** | 完成 SDKDownloadView.vue 各平台SDK下载 | SDKDownloadView.vue | P1 |
| PF1-2 | **配额预警页面** | 完成 QuotaAlertView.vue 配额使用和预警 | QuotaAlertView.vue | P1 |
| PF2-1 | **API使用统计页面** | 完成 APIUsageStatsView.vue 调用统计 | APIUsageStatsView.vue | P2 |

---

## 三、技术方案

### API 路由设计

| 接口 | 方法 | 说明 |
|------|------|------|
| `POST /api/v1/developer/register` | POST | 开发者注册 |
| `GET /api/v1/developer/profile` | GET | 开发者信息 |
| `PUT /api/v1/developer/profile` | PUT | 更新开发者信息 |
| `POST /api/v1/developer/verify` | POST | 申请认证 |
| `GET /api/v1/developer/apps` | GET | 应用列表 |
| `POST /api/v1/developer/apps` | POST | 创建应用 |
| `GET /api/v1/developer/apps/:id` | GET | 应用详情 |
| `PUT /api/v1/developer/apps/:id` | PUT | 更新应用 |
| `DELETE /api/v1/developer/apps/:id` | DELETE | 删除应用 |
| `POST /api/v1/developer/apps/:id/keys` | POST | 生成API Key |
| `DELETE /api/v1/developer/keys/:key_id` | DELETE | 吊销API Key |
| `GET /api/v1/developer/usage` | GET | API使用统计 |
| `GET /api/v1/webhook-market/templates` | GET | Webhook模板列表 |
| `GET /api/v1/webhook-market/webhooks` | GET | 我的Webhook列表 |
| `POST /api/v1/webhook-market/webhooks` | POST | 创建Webhook |
| `POST /api/v1/webhook-market/webhooks/:id/test` | POST | 测试Webhook |

### 数据库设计

```sql
-- 开发者表
CREATE TABLE developers (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    developer_name  VARCHAR(100) NOT NULL,
    developer_type  VARCHAR(20) DEFAULT 'personal', -- 'personal'/'enterprise'
    company_name    VARCHAR(200),
    verified_status VARCHAR(20) DEFAULT 'pending', -- 'pending'/'verified'/'rejected'
    verified_at     TIMESTAMP,
    contact_email  VARCHAR(255),
    website_url    VARCHAR(500),
    description     TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 开发者应用表
CREATE TABLE developer_apps (
    id              BIGSERIAL PRIMARY KEY,
    developer_id    BIGINT NOT NULL REFERENCES developers(id),
    app_name        VARCHAR(100) NOT NULL,
    app_type        VARCHAR(30) NOT NULL,            -- 'ios'/'android'/'miniapp'/'web'/'other'
    app_description TEXT,
    app_icon_url    VARCHAR(500),
    redirect_urls   VARCHAR(500)[],
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- API Key表
CREATE TABLE developer_api_keys (
    id              BIGSERIAL PRIMARY KEY,
    app_id          BIGINT NOT NULL REFERENCES developer_apps(id),
    key_prefix      VARCHAR(20) NOT NULL,            -- 'mdm_live_' / 'mdm_test_'
    key_hash        VARCHAR(255) NOT NULL,            -- SHA256 hash
    key_secret_enc  VARCHAR(255),                      -- encrypted secret
    scopes          VARCHAR(50)[],                    -- 权限范围
    rate_limit      INT DEFAULT 1000,                -- 每小时限流
    is_active       BOOLEAN DEFAULT TRUE,
    last_used_at    TIMESTAMP,
    expires_at      TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- API使用记录表
CREATE TABLE developer_api_usage (
    id              BIGSERIAL PRIMARY KEY,
    key_id          BIGINT NOT NULL REFERENCES developer_api_keys(id),
    endpoint        VARCHAR(200) NOT NULL,
    method          VARCHAR(10) NOT NULL,
    status_code     INT,
    response_time_ms INT,
    request_size    BIGINT,
    response_size   BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    INDEX idx_key_usage (key_id, created_at DESC)
);

-- Webhook市场表（已在 PRD_26_WEBHOOK_MARKET.md 定义）
-- 复用 webhook_market_templates / webhook_market_configs / webhook_market_logs
```

---

## 四、验收标准

### 4.1 功能验收

| 用例 | 验收条件 | 测试方法 |
|------|----------|----------|
| 开发者注册 | 注册流程完整，认证提交成功 | 注册测试 |
| 应用管理 | 应用CRUD正常，API Key生成成功 | CRUD测试 |
| Webhook市场 | 预置模板>5个，Webhook配置正常，推送成功率>95% | Webhook测试 |
| SDK构建 | iOS/Android/小程序SDK可下载，Demo可运行 | 下载测试 |

### 4.2 性能验收

| 验收点 | 标准 |
|--------|------|
| API Key验证延迟 | < 10ms |
| Webhook推送延迟 | < 5s |
| 开发者门户首屏 | < 2s |

---

## 五、依赖与风险

### 依赖

| 依赖 | 说明 |
|------|------|
| PRD_26_WEBHOOK_MARKET | Webhook市场 PRD |
| PRD_26_SDK_RELEASE | SDK发布 PRD |
| 现有 REST API v1 | API基础能力 |

### 风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| API安全漏洞 | 开发者数据泄露 | 安全审计+渗透测试 |
| Webhook推送失败率高 | 集成体验差 | 重试机制+监控告警 |
| SDK兼容性问题 | 开发者投诉 | 多版本测试+降级策略 |
