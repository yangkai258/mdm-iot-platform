# Agent YW - 运维工程师任务

**状态**: 🔄 进行中
**更新时间**: 2026-03-20 13:37 GMT+8

## 任务概述

为 AI 电子宠物 MDM 平台 Sprint 2 提供运维支持，完成应用管理和通知管理相关配置。

## Sprint 1.1 & 1.2 完成内容

### 1. Docker Compose MQTT 环境变量修复
- ✅ 将 `EMQX_BROKER_URL` 改为 `MQTT_BROKER`，与后端代码保持一致
- ✅ 添加 `MQTT_USERNAME` 和 `MQTT_PASSWORD` 环境变量
- ✅ 后端 `mqtt/handler.go` 从环境变量读取认证信息，不再硬编码 `admin/public`
- ✅ 修复 docker-compose.yml 和 docker-compose.prod.yml 中的 YAML 解析错误

### 2. Nginx 配置
- ✅ 创建 `ops/nginx/nginx.conf`
- ✅ 创建 `ops/nginx/conf.d/default.conf`
- ✅ 支持前端、后端、EMQX WebSocket 代理
- ✅ 创建 `ops/data/nginx/logs` 目录

### 3. OTA 升级文档
- ✅ PRODUCTION.md 添加 OTA 升级系统说明
- ✅ MQTT Topic 列表（`/mdm/device/{id}/up/status` 等）
- ✅ OTA Worker 环境变量说明
- ✅ OTA 数据库表说明
- ✅ OTA 指令格式说明
- ✅ EMQX 权限配置指引

### 4. 数据库迁移
- ✅ `ota_progress` 表通过 `db.AutoMigrate()` 自动创建
- ✅ `ota_packages` 和 `ota_deployments` 表同样自动创建

## Sprint 2 完成内容

### 1. 应用管理配置
- ✅ 添加 `APP_STORAGE_PATH` 环境变量（默认值 `/data/apps`）
- ✅ 添加 `app_storage_data` Docker Volume
- ✅ 挂载应用存储卷到 mdm-backend 容器
- ✅ 创建 `data/apps` 目录
- ✅ PRODUCTION.md 添加应用管理说明文档

### 2. 通知管理配置
- ✅ 确认 `/device/{id}/down/notification` MQTT Topic 可用
- ✅ EMQX 默认配置已允许所有 Topic（包括通知下发 Topic）
- ✅ PRODUCTION.md 添加通知管理系统说明
- ✅ 通知下发流程文档

### 3. 文档更新
- ✅ PRODUCTION.md 添加应用管理系统章节
- ✅ PRODUCTION.md 添加通知管理系统章节
- ✅ MQTT Topic 列表更新（包含 `/device/{device_id}/down/notification`）
- ✅ 数据目录创建命令更新（包含 `data/apps`）

## Git 提交记录

```
[master ccf1c2c] fix(ops): 重写 docker-compose 文件修复 YAML 解析错误
[master 1a807bb] fix(ops): 修复 MQTT 环境变量配置，添加 OTA 升级文档
[ops/sprint2 xxxxxxx] feat(ops): 添加应用管理和通知管理配置 (Sprint 2)
```

## 更新的文件

| 文件 | 说明 |
|------|------|
| `ops/docker-compose.prod.yml` | 添加 APP_STORAGE_PATH 环境变量和 app_storage_data 卷 |
| `ops/PRODUCTION.md` | 添加应用管理和通知管理系统文档 |
| `ops/data/apps/` | 创建应用存储目录 |
| `agent_tasks/agentyw_task.md` | 更新任务状态 |

## 部署验证

```bash
cd mdm-project/ops
export JWT_SECRET=$(openssl rand -base64 32)
export CORS_ALLOWED_ORIGINS=http://localhost:5173
export EMQX_ADMIN_PASSWORD=$(openssl rand -base64 24)
export POSTGRES_PASSWORD=$(openssl rand -base64 24)
mkdir -p data/postgres data/redis data/emqx data/nginx/logs data/apps
docker-compose -f docker-compose.prod.yml config  # 验证配置正确性
docker-compose -f docker-compose.prod.yml up -d   # 启动服务
```

## 待完成（需其他 Agent 配合）

- [ ] 编写 `backend/models/app_models.go`（应用管理数据模型）
- [ ] 编写 `backend/controllers/app_controller.go`（应用管理 API）
- [ ] 编写 `backend/models/notification_models.go`（通知管理数据模型）
- [ ] 编写 `backend/controllers/notification_controller.go`（通知管理 API）
- [ ] 添加通知下发 MQTT 发布功能（`/device/{id}/down/notification`）
- [ ] 编写 `frontend/Dockerfile.prod`（用于 docker-compose.prod.yml 生产构建）
- [ ] 编写 `.env.example` 模板文件
