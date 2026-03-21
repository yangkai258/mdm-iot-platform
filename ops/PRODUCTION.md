# MDM 平台生产部署指南

## 环境要求

- Docker Engine 24.0+
- Docker Compose v2.20+
- Git

## 快速启动

### 1. 克隆项目

```bash
cd mdm-project/ops
```

### 2. 生成必要的密钥

```bash
# 生成 JWT 密钥
export JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET=$JWT_SECRET"

# 生成 EMQX 管理员密码
export EMQX_ADMIN_PASSWORD=$(openssl rand -base64 24)
echo "EMQX_ADMIN_PASSWORD=$EMQX_ADMIN_PASSWORD"
```

### 3. 创建环境变量文件

```bash
cat > .env << EOF
# === 必须配置 ===

# JWT 密钥 (必填) - 生成方式: openssl rand -base64 32
JWT_SECRET=your-jwt-secret-here

# EMQX Dashboard 管理员密码 (必填)
EMQX_ADMIN_PASSWORD=your-emqx-password-here

# CORS 白名单 (必填) - 生产环境配置具体域名，禁用通配符
CORS_ALLOWED_ORIGINS=https://mdm.yourdomain.com

# PostgreSQL 密码 (必填)
POSTGRES_PASSWORD=your-postgres-password-here

# === 可选配置 ===

# PostgreSQL
POSTGRES_USER=mdm_user
POSTGRES_DB=mdm_db

# EMQX
EMQX_HOST=emqx
EMQX_ADMIN_USER=admin

# API Base URL (前端构建用)
VITE_API_BASE_URL=https://api.mdm.yourdomain.com

# 应用存储路径 (可选，默认 /data/apps)
# APP_STORAGE_PATH=/data/apps

# === 告警通知配置（可选）===
# SMTP 邮件通知
# SMTP_HOST=smtp.example.com
# SMTP_PORT=587
# SMTP_USER=alerts@example.com
# SMTP_PASSWORD=your-smtp-password
# SMTP_FROM=noreply@mdm.example.com
# SMTP_USE_TLS=true
# ALERT_ADMIN_EMAIL=admin@example.com

# Webhook 告警通知
# WEBHOOK_URL=https://hooks.example.com/alert
# WEBHOOK_TOKEN=your-webhook-secret
EOF
```

### 4. 创建数据目录

```bash
mkdir -p data/postgres data/redis data/emqx data/nginx/logs data/apps
```

### 5. 构建并启动所有服务

```bash
# 开发环境
docker-compose up -d --build

# 生产环境
docker-compose -f docker-compose.prod.yml up -d --build
```

### 6. 验证服务状态

```bash
docker-compose ps
docker-compose logs --tail=20
```

## 服务访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端 | http://localhost:80 | MDM 管理控制台 |
| 后端 API | http://localhost:8080 | REST API |
| EMQX Dashboard | http://localhost:18083 | MQTT 管理界面 |
| MQTT | tcp://localhost:1883 | 设备接入端口 |
| WebSocket MQTT | ws://localhost:8083 | WebSocket 设备接入 |

### 初始账号

- **EMQX Dashboard**: `admin` / (设置的 EMQX_ADMIN_PASSWORD)

## 服务启动顺序

Docker Compose `depends_on` + `condition: service_healthy` 确保以下顺序：

```
1. postgres (DB 就绪) → 
2. redis (缓存就绪) → 
3. emqx (MQTT Broker 就绪) → 
4. mdm-backend (API 服务) → 
5. mdm-frontend (前端) → 
6. nginx-proxy (反向代理)
```

## 环境变量说明

### JWT_SECRET

- **用途**: 签名 JWT Token 的密钥
- **要求**: 必须 ≥32 字符的随机字符串
- **生成**: `openssl rand -base64 32`
- **安全**: 禁止硬编码，禁止提交到 Git

### CORS_ALLOWED_ORIGINS

- **用途**: CORS 白名单域名
- **开发环境**: `http://localhost:5173,http://localhost:8080`
- **生产环境**: 必须配置具体域名，**禁止使用 `*`**
- **格式**: 多个域名用逗号分隔

### DATABASE_URL

- **格式**: `postgres://user:password@host:5432/dbname?sslmode=disable`
- **默认值**: `postgres://mdm_user:mdm_password@postgres:5432/mdm_db?sslmode=disable`
- **网络**: 使用 Docker 内部服务名 `postgres`，不走公网

### REDIS_URL

- **格式**: `redis://host:6379`
- **默认值**: `redis://redis:6379`
- **网络**: 使用 Docker 内部服务名 `redis`

### EMQX_BROKER_URL

- **格式**: `tcp://host:1883`
- **网络**: 使用 Docker 内部服务名 `emqx`

### APP_STORAGE_PATH

- **用途**: 应用安装包存储路径
- **默认值**: `/data/apps`
- **说明**: 容器内路径，对应 Docker Volume `app_storage_data`

### MQTT_TOPIC_PREFIX

- **用途**: MQTT Topic 前缀（用于通知下发）
- **默认值**: `/device`
- **说明**: 通知使用 `/device/{device_id}/down/notification`，OTA 使用 `/mdm/device/{device_id}/down/cmd`

### SMTP_HOST / SMTP_PASSWORD

- **用途**: 告警邮件发送（SMTP）
- **说明**: 配置 SMTP 服务器信息后，触发告警时自动发送邮件
- **参考值**: QQ 邮箱 `smtp.qq.com:587`，网易邮箱 `smtp.163.com:25`

### WEBHOOK_URL / WEBHOOK_TOKEN

- **用途**: 告警 Webhook 推送
- **说明**: 配置 Webhook 地址后，触发告警时 POST JSON 到该地址
- **安全**: 建议配置 `WEBHOOK_TOKEN` 用于签名验证

## 数据持久化

所有数据卷在 `data/` 目录下：

```
data/
├── postgres/   # PostgreSQL 数据文件
├── redis/      # Redis AOF + RDB 文件
├── emqx/       # EMQX 数据和日志
└── nginx/      # Nginx 日志
```

**重要**: 生产环境务必配置外部存储（如云存储挂载），不要只依赖本地目录。

### 应用存储

```
data/
└── apps/              # 应用安装包存储（IPA/APK/AAB/MSI）
    └── {app_code}/
        └── {version}/
            └── {package_file}
```

**生产环境建议**：使用云存储（S3/OSS）或 NFS 共享存储替代本地目录。

## 日志查看

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务
docker-compose logs -f mdm-backend
docker-compose logs -f emqx
docker-compose logs -f postgres

# 查看最近 100 行
docker-compose logs --tail=100 mdm-backend
```

## 健康检查

所有服务均配置了 healthcheck，可通过以下方式验证：

```bash
# 检查容器健康状态
docker-compose ps

# 手动检查
curl http://localhost:8080/health          # Backend
docker exec mdm-redis redis-cli ping         # Redis
docker exec mdm-postgres pg_isready -U mdm_user -d mdm_db  # PostgreSQL
docker exec mdm-emqx emqx _ctl status        # EMQX
```

## 备份策略

### PostgreSQL 备份

```bash
# 备份数据库
docker exec mdm-postgres pg_dump -U mdm_user mdm_db > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
docker exec -i mdm-postgres psql -U mdm_user mdm_db < backup_file.sql
```

### Redis 备份

```bash
# Redis 数据存储在 data/redis/ 目录
# AOF 文件: appendonly.aof
# 建议: 停止写入后复制备份
docker exec mdm-redis redis-cli BGSAVE
```

## 常见问题

### Q: 容器启动失败，提示 "JWT_SECRET must be set"

```bash
export JWT_SECRET=$(openssl rand -base64 32)
docker-compose up -d mdm-backend
```

### Q: EMQX 连接失败

检查 EMQX 是否完全启动（约需 60-90 秒）：
```bash
docker exec mdm-emqx emqx _ctl status
```

### Q: PostgreSQL 连接被拒绝

检查 healthcheck 是否通过：
```bash
docker-compose ps postgres
docker logs mdm-postgres | grep "database system is ready"
```

### Q: CORS 跨域错误

确认 `CORS_ALLOWED_ORIGINS` 环境变量已正确设置，**不要使用 `*`**，生产环境必须配置具体域名。

## 停止服务

```bash
# 停止（保留数据卷）
docker-compose down

# 完全停止并删除数据卷（危险！）
docker-compose down -v

# 生产环境优雅停止
docker-compose -f docker-compose.prod.yml down
```

## OTA 升级系统

### 概述

OTA（Over-The-Air）固件升级模块已集成到 `mdm-backend` 服务中，作为后台 goroutine 运行。每 5 分钟检查一次待下发的部署任务。

### MQTT Topic 列表

| Topic 模式 | 方向 | 说明 |
|------------|------|------|
| `/mdm/device/{device_id}/up/status` | 设备→服务器 | 设备心跳上报 |
| `/mdm/device/{device_id}/up/property` | 设备→服务器 | 设备属性上报 |
| `/mdm/device/{device_id}/down/cmd` | 服务器→设备 | 设备指令下发（含OTA指令）|
| `/mdm/device/{device_id}/down/desired` | 服务器→设备 | 期望状态下发 |

### OTA Worker 配置

OTA Worker 通过以下环境变量与 Redis 通信：

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `REDIS_URL` | Redis 连接地址 | `redis://redis:6379` |
| `MQTT_BROKER` | MQTT Broker 地址 | `tcp://localhost:1883` |
| `MQTT_USERNAME` | MQTT 认证用户名 | `admin` |
| `MQTT_PASSWORD` | MQTT 认证密码 | `public` |

### OTA 数据库表

以下表在 `db.AutoMigrate()` 时自动创建：

- `ota_packages` - 固件包记录
- `ota_deployments` - 部署任务
- `ota_progress` - 设备升级进度

### OTA 指令格式

通过 MQTT 下发的 OTA 指令：

```json
{
  "cmd_id": "ota-{deployment_id}-{device_id}",
  "cmd_type": "ota",
  "ota": {
    "version": "v1.3.0",
    "url": "https://cdn.example.com/firmware/v1.3.0.bin",
    "md5": "d41d8cd98f00b204e9800998ecf8427e"
  },
  "timestamp": "2026-03-20T12:00:00Z"
}
```

### EMQX 权限配置

EMQX 默认配置已允许上述 Topic。生产环境如需精细化配置 ACL，可通过 Dashboard 或 REST API 配置：

```bash
# 通过 EMQX Dashboard API 创建 ACL 规则
curl -X POST http://localhost:18083/api/v5/authentication \
  -u admin:your-password \
  -H "Content-Type: application/json" \
  -d '{
    "type": "acl_file",
    "enabled": true
  }'
```

### 设备影子与 OTA

OTA 期望版本可在设备影子的 `desired_config.desired_firmware` 字段中设置。OTA Worker 会检查该字段并在适当时机触发升级。

---

## 应用管理系统

### 概述

应用管理模块为 MDM 中台提供企业应用分发能力，支持上传和管理企业内部应用（IPA/APK/AAB/MSI），并通过分发策略向设备推送安装、强制更新或卸载。

### 环境变量

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `APP_STORAGE_PATH` | 应用安装包存储路径 | `/data/apps` |

### 应用存储配置

应用安装包（IPA/APK/AAB/MSI）存储在 Docker Volume `app_storage_data` 中，容器内路径为 `/data/apps`。

**目录结构：**
```
data/
└── apps/
    └── {app_code}/
        └── {version}/
            └── {package_file}
```

**生产环境建议：**
- 使用云存储（如 S3、OSS）替代本地存储
- 配置 CDN 加速应用包下载
- 定期备份应用仓库

```bash
# 生产环境使用对象存储（示例：S3）
export APP_STORAGE_PATH=s3://mdm-bucket/apps
# 或使用 NFS 共享存储
export APP_STORAGE_PATH=/mnt/nfs/apps
```

### 数据表

应用管理模块使用以下数据表（通过 `db.AutoMigrate()` 自动创建）：

| 表名 | 说明 |
|------|------|
| `apps` | 应用主表 |
| `app_versions` | 应用版本表 |
| `app_distributions` | 应用分发任务表 |
| `app_licenses` | 应用许可证表（VPP）|
| `app_configurations` | 应用托管配置表 |
| `app_installations` | 应用安装记录表 |

### EMQX 权限配置

应用分发通知通过 MQTT Topic `/device/{device_id}/down/notification` 下发。该 Topic 需要在 EMQX 中配置发布权限。

**EMQX 默认配置已允许所有 Topic。** 生产环境如需精细化配置 ACL，可通过 Dashboard 或 REST API 配置：

```bash
# 通过 EMQX Dashboard API 配置精细化 ACL
# 1. 登录 EMQX Dashboard (http://localhost:18083)
# 2. 访问 访问控制 → ACL > 内置数据库
# 3. 添加 ACL 规则：
#    - 允许客户端发布到 /device/%u/down/notification（%u 为用户名）
#    - 允许客户端订阅 /device/%u/up/#

# 或通过 CLI 创建 ACL 文件
docker exec mdm-emqx emqx ctl acl reload
```

---

## 通知管理系统

### 概述

通知与消息模块为 MDM 中台提供多渠道推送能力，支持向设备发送文本通知、企业公告，并提供命令反馈查看能力。

### MQTT Topic 列表

| Topic 模式 | 方向 | 说明 |
|------------|------|------|
| `/mdm/device/{device_id}/up/status` | 设备→服务器 | 设备心跳上报 |
| `/mdm/device/{device_id}/up/property` | 设备→服务器 | 设备属性上报 |
| `/mdm/device/{device_id}/down/cmd` | 服务器→设备 | 设备指令下发（含OTA指令）|
| `/mdm/device/{device_id}/down/desired` | 服务器→设备 | 期望状态下发 |
| `/device/{device_id}/down/notification` | 服务器→设备 | **通知/消息下发** |

### 通知下发 Topic

通知通过 MQTT Topic `/device/{device_id}/down/notification` 下发到设备：

**通知消息格式：**
```json
{
  "notification_id": "notif-uuid-001",
  "title": "固件升级通知",
  "content": "有新版本固件可用，请及时更新",
  "notification_type": "push",
  "priority": "normal",
  "timestamp": "2026-03-20T12:00:00Z"
}
```

**EMQX 权限要求：**
- 后端服务（mdm-backend）需要对 `/device/+/down/notification` 有 **发布权限**
- EMQX 默认配置已允许该 Topic，无需额外配置

```bash
# 验证 EMQX Topic 权限
docker exec mdm-emqx emqx ctl brokers
docker exec mdm-emqx emqx ctl vm
```

### 数据表

通知管理模块使用以下数据表（通过 `db.AutoMigrate()` 自动创建）：

| 表名 | 说明 |
|------|------|
| `notifications` | 通知主表 |
| `notification_templates` | 通知模板表 |
| `announcements` | 公告表 |
| `device_notifications` | 设备通知记录表 |

### 通知发送流程

```
管理员发送推送通知
    │
    ▼
POST /api/v1/notifications/push
    │
    ├─→ 创建 notifications 记录 (status=pending)
    ├─→ 解析 target_ids 批量创建 device_notifications
    │
    ▼
通知发送 Worker (后台处理)
    │
    ▼
遍历 device_notifications (status=pending)
    │
    ├─→ 通过 MQTT 发送到设备 /device/{id}/down/notification
    │
    ├─→ 成功 ──► status=delivered, delivered_at=now()
    │
    └─→ 失败 ──► status=failed, error_message=...
    │
    ▼
更新 notifications.sent_count / failed_count
```

### 常见问题

**Q: 设备收不到通知**

1. 检查设备是否在线：
```bash
docker exec mdm-redis redis-cli GET "device_shadow:{device_id}"
```

2. 检查 EMQX 连接状态：
```bash
docker exec mdm-emqx emqx ctl client list
```

3. 检查后端日志：
```bash
docker-compose logs mdm-backend | grep notification
```

---

## 策略管理系统

### 概述

策略管理模块通过合规策略（CompliancePolicy）对设备进行实时监控和自动处置。当设备数据上报触发策略条件时，系统自动记录违规并执行预配置的补救措施（隔离、阻止、通知、擦除）。

### 数据表

合规策略模块使用以下数据表（通过 `db.AutoMigrate()` 自动创建）：

| 表名 | 说明 |
|------|------|
| `compliance_policies` | 合规策略主表 |
| `compliance_violations` | 违规记录表 |

### CompliancePolicy 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | varchar(100) | 策略名称 |
| `description` | varchar(255) | 策略描述 |
| `policy_type` | varchar(50) | 策略类型：`firmware_version`, `battery_level`, `region_lock`, `encryption_required` |
| `target_value` | varchar(100) | 目标值，如版本号、最低电量等 |
| `condition` | varchar(20) | 条件：`=`, `!=`, `>=`, `<=`, `<`, `>` |
| `severity` | int | 严重程度：1-低 2-中 3-高 4-严重 |
| `remediation_action` | varchar(50) | 补救措施：`isolate`（隔离）、`wipe`（擦除）、`notify`（通知）、`block`（阻止） |
| `enabled` | bool | 是否启用 |
| `enforce_scope` | varchar(50) | 生效范围：`all`（全部）、`group`（分组）、`individual`（单个） |

### 策略检查流程

```
设备上报数据 (MQTT /device/{id}/up/status 或 /up/property)
    │
    ▼
CheckCompliance() 回调处理
    │
    ▼
遍历所有启用的 compliance_policies
    │
    ├─→ battery_level: 比较电池电量与 target_value
    ├─→ offline_duration: 比较离线时长与阈值
    └─→ is_online: 检测设备在线状态
    │
    ▼
条件命中 → 创建 compliance_violations 记录
    │
    ▼
executeRemediation() 执行补救措施
    │
    ├─→ notify   ──► 创建告警记录（alert_type=compliance_violation）
    ├─→ isolate  ──► 创建严重告警 + 更新设备影子 current_mode=isolated
    ├─→ block    ──► 创建严重告警（设备阻止）
    └─→ wipe     ──► 创建紧急告警 + 状态置为处理中
```

### API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/compliance/policies` | 获取策略列表 |
| POST | `/api/v1/compliance/policies` | 创建策略 |
| PUT | `/api/v1/compliance/policies/:id` | 更新策略 |
| DELETE | `/api/v1/compliance/policies/:id` | 删除策略 |
| GET | `/api/v1/compliance/violations` | 获取违规记录 |
| PUT | `/api/v1/compliance/violations/:id/resolve` | 标记违规已处理 |

### 常见问题

**Q: 策略没有触发**

1. 检查策略是否启用：
```bash
docker-compose exec mdm-backend curl http://localhost:8080/api/v1/compliance/policies
```

2. 检查设备数据是否上报：
```bash
docker exec mdm-redis redis-cli GET "device_shadow:{device_id}"
```

3. 检查后端日志中的合规检查：
```bash
docker-compose logs mdm-backend | grep Compliance
```

---

## 告警通知系统

### 概述

告警通知模块提供设备告警规则配置和真实告警通知（邮件、Webhook）能力。当设备触发告警规则或合规策略时，系统根据规则配置的 `notify_ways` 字段，通过邮件或 Webhook 发送告警通知。

### 告警规则

`DeviceAlertRule` 字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | varchar(100) | 规则名称 |
| `device_id` | varchar(36) | 关联设备ID（空表示所有设备） |
| `alert_type` | varchar(50) | 告警类型：`battery_low`, `offline`, `temperature_high` 等 |
| `condition` | varchar(100) | 条件：`<`, `>`, `=`, `>=`, `<=` |
| `threshold` | float64 | 阈值 |
| `severity` | int | 严重程度：1-低 2-中 3-高 4-严重 |
| `enabled` | bool | 是否启用 |
| `notify_ways` | varchar(100) | 通知方式：`email`, `sms`, `webhook`（逗号分隔） |
| `remark` | varchar(255) | 备注 |

### 数据表

| 表名 | 说明 |
|------|------|
| `device_alert_rules` | 告警规则表 |
| `device_alerts` | 告警记录表（包含触发值、告警状态） |

### 告警触发流程

```
MQTT 设备数据上报
    │
    ▼
CheckAlerts(db, deviceID, data)
    │
    ▼
遍历所有启用的 device_alert_rules
    │
    ├─→ battery_low: data["battery"] vs threshold
    └─→ offline:     data["is_online"] == false
    │
    ▼
触发条件命中
    │
    ▼
创建 device_alerts 记录 (status=1:未处理)
    │
    ▼
后台 Worker 查询 notify_ways
    │
    ├─→ email   ──► SMTP 发送邮件到 ALERT_ADMIN_EMAIL
    └─→ webhook ──► POST 到 WEBHOOK_URL
```

### 环境变量配置

在 `.env` 文件中添加以下可选配置（启用邮件或 Webhook 通知时必须）：

```bash
# === 告警通知配置 ===

# SMTP 邮件通知（可选，不配置则不发送邮件）
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=alerts@example.com
SMTP_PASSWORD=your-smtp-password
SMTP_FROM=noreply@mdm.example.com
SMTP_USE_TLS=true

# Webhook 告警通知（可选，不配置则不发送 Webhook）
WEBHOOK_URL=https://hooks.example.com/alert
WEBHOOK_TOKEN=your-webhook-secret

# 告警管理员邮箱（接收严重告警邮件）
ALERT_ADMIN_EMAIL=admin@example.com
```

### SMTP 配置说明

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `SMTP_HOST` | SMTP 服务器地址 | （空，不启用）|
| `SMTP_PORT` | SMTP 端口 | `587` |
| `SMTP_USER` | SMTP 用户名 | （空）|
| `SMTP_PASSWORD` | SMTP 密码 | （空）|
| `SMTP_FROM` | 发件人地址 | `noreply@mdm.example.com` |
| `SMTP_USE_TLS` | 是否使用 TLS | `true` |
| `ALERT_ADMIN_EMAIL` | 告警接收邮箱 | （空）|

**SMTP 配置示例（QQ 邮箱）：**
```bash
SMTP_HOST=smtp.qq.com
SMTP_PORT=587
SMTP_USER=your-email@qq.com
SMTP_PASSWORD=your-authorization-code
SMTP_FROM=your-email@qq.com
SMTP_USE_TLS=true
ALERT_ADMIN_EMAIL=admin@example.com
```

**获取 QQ 邮箱授权码：**
1. 登录 QQ 邮箱 → 设置 → 账户
2. 开启 POP3/SMTP 服务
3. 生成授权码（填入 `SMTP_PASSWORD`）

### Webhook 配置说明

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `WEBHOOK_URL` | Webhook 接收地址 | （空，不启用）|
| `WEBHOOK_TOKEN` | 签名密钥（用于 HMAC 签名） | （空）|

**Webhook POST 请求格式：**
```json
{
  "alert_id": 123,
  "device_id": "device-001",
  "alert_type": "battery_low",
  "severity": 3,
  "message": "电池电量过低",
  "trigger_val": 15.0,
  "threshold": 20.0,
  "timestamp": "2026-03-20T12:00:00Z"
}
```

**请求头：**
- `Content-Type: application/json`
- `X-Webhook-Token: {WEBHOOK_TOKEN}`（如果配置了 `WEBHOOK_TOKEN`）

**Webhook 签名验证（推荐）：**
```go
// 服务器端 HMAC-SHA256 签名验证
signature := hmac.New(sha256.New, []byte(WEBHOOK_TOKEN))
signature.Write(body)
expected := hex.EncodeToString(signature.Sum(nil))
```

### API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/alerts/rules` | 获取告警规则列表 |
| POST | `/api/v1/alerts/rules` | 创建告警规则 |
| GET | `/api/v1/alerts` | 获取告警记录 |
| PUT | `/api/v1/alerts/:id` | 更新告警状态（确认/解决）|
| GET | `/api/v1/dashboard/stats` | 大盘统计数据（含待处理告警数）|

### 常见问题

**Q: 告警触发了但没有收到邮件**

1. 确认 SMTP 环境变量已配置：
```bash
docker-compose exec mdm-backend env | grep SMTP
```

2. 检查后端日志中的邮件发送错误：
```bash
docker-compose logs mdm-backend | grep -i "smtp\|email\|mail"
```

3. 检查垃圾邮件文件夹

**Q: Webhook 没有收到请求**

1. 确认 `WEBHOOK_URL` 已配置
2. 检查 Webhook 服务端是否正常可达
3. 检查后端日志：
```bash
docker-compose logs mdm-backend | grep -i webhook
```

**Q: 如何禁用邮件通知**

将 `SMTP_HOST` 留空，或在告警规则的 `notify_ways` 中不包含 `email`。

---

## 数据库迁移说明

### AutoMigrate 自动迁移

所有数据表在服务首次启动时通过 `db.AutoMigrate()` 自动创建，包括：

| 模块 | 表 |
|------|-----|
| 设备 | `devices`, `device_shadows` |
| OTA | `ota_packages`, `ota_deployments`, `ota_progress` |
| 应用管理 | `apps`, `app_versions`, `app_distributions`, `app_install_records`, `app_licenses` |
| 通知 | `notifications`, `notification_templates`, `announcements`, `device_notifications` |
| 告警 | `device_alert_rules`, `device_alerts` |
| 合规策略 | `compliance_policies`, `compliance_violations` |
| 系统 | `sys_users`, `sys_roles`, `sys_menus`, `sys_dictionaries`, `sys_operation_logs`, `sys_login_logs` |
| 会员 | `member_orders`, `member_upgrade_records` |

**注意**：`AutoMigrate` 只会创建不存在的表，不会修改已存在表的结构。如需执行结构变更，请手动编写 ALTER TABLE 迁移脚本。

---

## 多租户部署说明

### 概述

MDM 平台支持多租户部署模式，通过 `tenant_id` 字段实现数据隔离。每个租户拥有独立的设备、会员、告警等数据，API 层面通过中间件自动注入租户上下文。

### 租户隔离机制

| 隔离维度 | 实现方式 |
|----------|----------|
| 数据隔离 | 所有数据表含 `tenant_id` 字段，GORM Scope 自动过滤 |
| MQTT Topic | `/mdm/tenant/{tenant_id}/device/{device_id}/...` |
| API 隔离 | 中间件从 JWT claims 提取 `tenant_id` |

### 租户 API 请求示例

```bash
# 登录时获取租户 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"tenant_admin","password":"xxx","tenant_id":"tenant_001"}'

# 请求时携带租户上下文
curl http://localhost:8080/api/v1/devices \
  -H "Authorization: Bearer <token>"
# JWT Payload 中包含: { "tenant_id": "tenant_001", "user_id": "...", "role": "tenant_admin" }
```

### 部署多租户模式

#### 方式一：共享 PostgreSQL + Schema 隔离（推荐）

```bash
# docker-compose.prod.yml 中配置多数据库
services:
  postgres:
    environment:
      POSTGRES_USER: mdm_user
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: mdm_db
    command:
      - postgres
      - "-c"
      - "search_path=public,tenant_a,tenant_b"
```

#### 方式二：独立 PostgreSQL 实例（高隔离需求）

```bash
# 为每个租户创建独立数据库
docker exec mdm-postgres psql -U mdm_user -d mdm_db -c \
  "CREATE DATABASE tenant_a_db;"
docker exec mdm-postgres psql -U mdm_user -d mdm_db -c \
  "CREATE DATABASE tenant_b_db;"
```

#### 方式三：多租户共用 Schema（成本优先）

```bash
# 通过 tenant_id 字段隔离，共用数据库实例
# 环境变量配置
TENANT_ISOLATION_MODE=shared_schema
DEFAULT_TENANT_ID=default
```

### 环境变量说明

| 环境变量 | 说明 | 可选值 |
|----------|------|--------|
| `TENANT_ISOLATION_MODE` | 租户隔离模式 | `shared_schema`(默认), `schema_per_tenant`, `db_per_tenant` |
| `DEFAULT_TENANT_ID` | 默认租户 ID | `default` |
| `MAX_TENANTS` | 最大租户数限制 | `100` |

### 租户管理 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/tenants` | 获取租户列表 |
| POST | `/api/v1/tenants` | 创建租户 |
| GET | `/api/v1/tenants/:id` | 获取租户详情 |
| PUT | `/api/v1/tenants/:id` | 更新租户配置 |
| DELETE | `/api/v1/tenants/:id` | 删除租户 |
| GET | `/api/v1/tenants/:id/stats` | 租户统计数据 |

### 租户设备 Topic 规范

```
# 设备上报
/mdm/tenant/{tenant_id}/device/{device_id}/up/status
/mdm/tenant/{tenant_id}/device/{device_id}/up/property

# 服务下发
/mdm/tenant/{tenant_id}/device/{device_id}/down/cmd
/mdm/tenant/{tenant_id}/device/{device_id}/down/notification
/mdm/tenant/{tenant_id}/device/{device_id}/down/desired
```

### 常见问题

**Q: 租户间数据泄露**

确认所有 API 查询均通过 GORM Scope 过滤：
```go
// 正确：在查询时强制注入 tenant_id
func ListDevices(db *gorm.DB, tenantID string) []Device {
    return db.Scopes(TenantScope(tenantID)).Find(&devices)
}

// 错误：直接查询未过滤
func ListDevices(db *gorm.DB) []Device {
    return db.Find(&devices) // 危险！
}
```

---

## 单元测试运行说明

### 概述

MDM 平台使用 pytest 作为 Python 测试框架，P0 测试覆盖核心功能回归。测试代码位于 `testing/p0_tests/` 目录。

### 前置条件

```bash
# 确保已在宿主机安装 Python 3.8+ 和 pip
python --version   # >= 3.8

# 安装测试依赖
cd C:\Users\YKing\.openclaw\workspace\mdm-project\backend\testing\p0_tests
pip install -r requirements.txt
```

### 运行测试

#### 方式一：直接运行（需要 MDM 服务已启动）

```bash
# 方式一：直接运行（需要 MDM 服务已启动）

# 设置环境变量（可选，有默认值）
set MDM_API_BASE_URL=http://localhost:8080
set TEST_USERNAME=admin
set TEST_PASSWORD=admin123

# 运行所有 P0 测试
pytest p0_tests/ -v

# 运行单个测试文件
pytest p0_tests/test_jwt_config.py -v
pytest p0_tests/test_cors.py -v
pytest p0_tests/test_mqtt_command.py -v
pytest p0_tests/test_login_enter_key.py -v

# 带详细输出
pytest p0_tests/ -v --tb=short

# 生成 HTML 报告
pytest p0_tests/ -v --html=report.html --self-contained-html
```

#### 方式二：Docker 容器内运行

```bash
# 启动测试容器（需要服务已启动）
docker run --rm \
  -v C:\Users\YKing\.openclaw\workspace\mdm-project\backend:/app \
  -w /app/testing/p0_tests \
  -e MDM_API_BASE_URL=http://host.docker.internal:8080 \
  python:3.11-slim \
  sh -c "pip install -r requirements.txt && pytest . -v"
```

#### 方式三：通过 docker-compose 自动运行

```bash
# 在项目根目录执行
docker-compose -f docker-compose.prod.yml up -d
docker-compose -f docker-compose.prod.yml exec backend \
  sh -c "cd /app/testing/p0_tests && pytest . -v"
```

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MDM_API_BASE_URL` | `http://localhost:8080` | 后端 API 地址 |
| `MDM_BACKEND_PATH` | `项目目录/backend` | 后端源码路径 |
| `TEST_USERNAME` | `admin` | 测试用户名 |
| `TEST_PASSWORD` | `admin123` | 测试密码 |
| `MQTT_BROKER` | `tcp://localhost:1883` | MQTT Broker 地址 |
| `MQTT_USERNAME` | `admin` | MQTT 用户名 |
| `MQTT_PASSWORD` | `public` | MQTT 密码 |

### 测试用例说明

| 测试文件 | 测试目标 | 覆盖范围 |
|----------|----------|----------|
| `test_jwt_config.py` | JWT 密钥配置 | 硬编码密钥检测、环境变量生效验证 |
| `test_cors.py` | CORS 安全配置 | 跨域响应头检查、预检请求验证 |
| `test_mqtt_command.py` | MQTT 指令下发 | 离线设备处理、指令参数校验 |
| `test_login_enter_key.py` | 登录页 Enter 键 | Vue 事件绑定分析 |
| `test_device_shadow.py` | 设备影子 | Redis 设备状态读写 |
| `test_alert_trigger.py` | 告警触发 | 告警规则条件匹配 |
| `test_notification.py` | 通知下发 | MQTT 通知发送 |
| `test_policy_management.py` | 合规策略 | 策略 CRUD 和执行 |

### CI/CD 集成

在 CI pipeline 中集成测试：

```yaml
# .github/workflows/test.yml
- name: Run P0 Tests
  run: |
    docker-compose -f docker-compose.prod.yml up -d
    sleep 30  # 等待服务就绪
    pytest testing/p0_tests/ -v --junitxml=report.xml
  env:
    MDM_API_BASE_URL: http://localhost:8080
```

### 常见问题

**Q: 测试提示 "connection refused"**

确认 MDM 后端服务已启动并监听 8080 端口：
```bash
docker-compose ps
curl http://localhost:8080/health
```

**Q: JWT 测试失败**

确认环境变量 `JWT_SECRET` 已正确设置（非默认值）：
```bash
docker-compose exec backend env | grep JWT_SECRET
```

---

## 日志配置说明

### 日志架构

MDM 平台采用分层日志架构：

```
应用日志 (Gin/Zap)
    ↓
标准输出 (stdout)
    ↓
Docker 日志驱动 (json-file)
    ↓
宿主机日志文件 /var/lib/docker/containers/
```

### 日志级别配置

| 环境变量 | 可选值 | 默认值 | 说明 |
|----------|--------|--------|------|
| `GIN_MODE` | `debug`, `release`, `test` | `release` | Gin 框架模式 |
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error` | `info` | 应用日志级别 |
| `LOG_FORMAT` | `json`, `console` | `json` | 日志输出格式 |

### 配置示例

```bash
# .env 文件中添加
GIN_MODE=release
LOG_LEVEL=info
LOG_FORMAT=json
```

### 应用日志配置

```go
// Gin 中间件日志配置 (main.go)
import "github.com/gin-contrib/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

router.Use(zap.Ginzap(logger, time.RFC3339, true))
router.Use(zap.RecoveryWithZap(logger, true))
```

### 日志字段规范

所有日志使用 JSON 格式，包含以下标准字段：

```json
{
  "level": "info",
  "ts": "2026-03-22T00:40:00.000+08:00",
  "caller": "handler.go:42",
  "msg": "device command sent",
  "tenant_id": "tenant_001",
  "device_id": "device-abc-123",
  "trace_id": "trace-xyz-789",
  "duration_ms": 15
}
```

### 查看日志

```bash
# 宿主机层面查看容器日志
docker logs mdm-backend --tail=100 -f

# 使用 docker-compose 查看
docker-compose logs -f mdm-backend
docker-compose logs --tail=200 mdm-backend | grep "error"

# 过滤特定租户/设备日志
docker-compose logs -f mdm-backend | grep "tenant_001"

# 导出日志到文件
docker-compose logs mdm-backend > backend.log 2>&1
```

### Nginx 日志配置

```nginx
# nginx/conf.d/default.conf

http {
    # 访问日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for" '
                    'rt=$request_time uct="$upstream_connect_time" '
                    'uht="$upstream_header_time" urt="$upstream_response_time"';

    access_log /var/log/nginx/access.log main;

    # 错误日志级别
    error_log /var/log/nginx/error.log warn;

    # 日志轮转（通过 logrotate 配置）
    # /etc/logrotate.d/nginx
}
```

### 日志轮转配置

```bash
# /etc/logrotate.d/mdm
/data/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 root root
    sharedscripts
    postrotate
        docker-compose -f /path/to/docker-compose.yml kill -s USR1 mdm-nginx-proxy
    endscript
}
```

### 结构化日志分析

```bash
# 使用 jq 分析 JSON 日志
docker-compose logs -f mdm-backend --tail=500 | jq '. | select(.level=="error")'

# 统计错误类型
docker-compose logs mdm-backend --since=1h | jq -r '.msg' | sort | uniq -c | sort -rn

# 慢请求分析（>1s）
docker-compose logs mdm-backend --since=1h | jq -r 'select(.duration_ms > 1000) | .msg'
```

### ELK/Grafana Loki 集成

```yaml
# docker-compose.prod.yml 添加日志收集
services:
  # ... 其他服务 ...

  # Loki 日志收集（替代默认 json-file 驱动）
  promtail:
    image: grafana/promtail:latest
    volumes:
      - ./logs:/var/log/mdm
      - ./promtail.yaml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
```

```yaml
# promtail.yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: mdm-backend
    static_configs:
      - targets:
          - localhost
        labels:
          job: mdm-backend
          __path__: /var/log/mdm/backend/*.log
```

### 日志保留策略

| 日志类型 | 保留时间 | 存储位置 |
|----------|----------|----------|
| 应用日志 | 30 天 | `/data/logs/app/` |
| Nginx 访问日志 | 90 天 | `/data/logs/nginx/` |
| Nginx 错误日志 | 90 天 | `/data/logs/nginx/` |
| EMQX 日志 | 7 天 | `/data/emqx/log/` |
| Docker 容器日志 | 7 天 | `/var/lib/docker/containers/` |

**注意**：生产环境建议将日志推送至对象存储（S3/OSS）或日志服务（ELK、Grafana Loki），不要仅依赖本地磁盘存储。
```
