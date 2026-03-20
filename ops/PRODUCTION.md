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
EOF
```

### 4. 创建数据目录

```bash
mkdir -p data/postgres data/redis data/emqx data/nginx/logs
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
```
