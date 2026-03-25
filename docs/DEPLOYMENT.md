# MDM IoT Platform 部署文档

## 目录

- [环境要求](#环境要求)
- [快速部署](#快速部署)
- [Docker Compose 部署](#docker-compose-部署)
- [生产环境部署](#生产环境部署)
- [GitHub Actions CI/CD](#github-actions-cicd)
- [环境变量配置](#环境变量配置)
- [服务健康检查](#服务健康检查)
- [故障排查](#故障排查)

---

## 环境要求

| 组件 | 版本 | 说明 |
|------|------|------|
| Docker | 20.10+ | 容器运行时 |
| Docker Compose | 2.0+ | 多容器编排 |
| Git | 2.30+ | 代码版本管理 |

---

## 快速部署

### 1. 克隆项目

```bash
git clone https://github.com/yangkai258/mdm-iot-platform.git
cd mdm-iot-platform
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入实际值
```

### 3. 一键启动

```bash
docker compose up -d
```

---

## Docker Compose 部署

### 服务架构

```
┌─────────────────────────────────────────────────────────────┐
│                      mdm_frontend                          │
│                   (Nginx + Vue3 SPA)                       │
│                        :80/:443                             │
└──────────────────────────┬────────────────────────────────┘
                           │
┌──────────────────────────▼────────────────────────────────┐
│                       mdm_backend                          │
│                  (Go + Gin + GORM)                        │
│                        :8080                               │
└──────┬─────────────────┬───────────────────┬──────────────┘
       │                 │                   │
┌──────▼──────┐  ┌───────▼───────┐  ┌──────▼────────┐
│  postgres   │  │     redis     │  │     emqx     │
│  :5432      │  │    :6379      │  │ :1883/:8083  │
│  (DB)       │  │  (缓存/影子)  │  │  (MQTT)      │
└─────────────┘  └───────────────┘  └──────────────┘
```

### 启动所有服务

```bash
docker compose up -d
```

### 查看服务状态

```bash
docker compose ps
```

### 查看日志

```bash
# 所有服务
docker compose logs -f

# 单个服务
docker compose logs -f mdm_backend
docker compose logs -f mdm_frontend
docker compose logs -f postgres
docker compose logs -f redis
docker compose logs -f emqx
```

### 健康检查端点

| 服务 | 端点 | 说明 |
|------|------|------|
| Backend | `http://localhost:8080/health` | Go 服务健康检查 |
| Frontend | `http://localhost:80/health` | Nginx 健康检查 |
| PostgreSQL | `pg_isready -U mdm_user` | 数据库就绪检查 |
| Redis | `redis-cli ping` | Redis ping |
| EMQX | `emqx ctl status` | MQTT Broker 状态 |

### 资源限制

| 服务 | CPU 限制 | 内存限制 |
|------|----------|----------|
| mdm_backend | 1 core | 512 MB |
| mdm_frontend | 0.5 core | 128 MB |
| postgres | 1 core | 512 MB |
| redis | 0.5 core | 256 MB |
| emqx | 1 core | 512 MB |

### 重启策略

所有服务配置为 `restart: unless-stopped`，确保服务异常退出后自动重启。

---

## 生产环境部署

### 1. 服务器准备

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh

# 安装 Docker Compose
apt install docker-compose-plugin

# 配置 Docker 开机自启
systemctl enable docker
```

### 2. 使用 systemd 管理服务

创建 `/etc/systemd/system/mdm.service`:

```ini
[Unit]
Description=MDM IoT Platform
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/mdm
ExecStart=/usr/local/bin/docker compose up -d
ExecStop=/usr/local/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

启用服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable mdm
sudo systemctl start mdm
```

### 3. Nginx 反向代理（可选）

```nginx
server {
    listen 80;
    server_name mdm.example.com;

    location / {
        proxy_pass http://127.0.0.1:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## GitHub Actions CI/CD

### 工作流说明

| 工作流文件 | 触发条件 | 说明 |
|------------|----------|------|
| `main.yml` | PR/推送至 main/master/develop | 构建 + 测试 |
| `deploy.yml` | main.yml 成功后自动触发 | 镜像构建 + 部署 |

### CI/CD 流程

```
┌─────────────┐    push/PR     ┌─────────────┐
│   main.yml  │ ────────────▶ │ Build & Test│
└─────────────┘                └──────┬──────┘
                                       │
                              (all jobs passed)
                                       │
                                       ▼
                               ┌─────────────┐
                               │ deploy.yml   │
                               └──────┬───────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                                     ▼
            ┌───────────────┐                   ┌───────────────┐
            │ Deploy Staging│                   │Deploy Prod    │
            │ (auto)        │                   │(manual approve│
            └───────────────┘                   └───────────────┘
```

### 部署配置

#### Staging 环境变量

在 GitHub仓库 Settings → Secrets 中配置：

| Secret 名称 | 说明 |
|-------------|------|
| `STAGING_HOST` | Staging 服务器 IP/域名 |
| `STAGING_USER` | SSH 用户名 |
| `STAGING_SSH_KEY` | SSH 私钥 |
| `STAGING_APP_DIR` | 应用部署目录 |

#### Production 环境变量

| Secret 名称 | 说明 |
|-------------|------|
| `PROD_HOST` | Production 服务器 IP/域名 |
| `PROD_USER` | SSH 用户名 |
| `PROD_SSH_KEY` | SSH 私钥 |
| `PROD_APP_DIR` | 应用部署目录 |

### 手动触发部署

```bash
# 通过 GitHub CLI
gh workflow run deploy.yml --ref main
```

---

## 环境变量配置

### .env 文件示例

```bash
# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://mdm.example.com

# EMQX
EMQX_DASHBOARD__DEFAULT_PASSWORD=your-secure-password
```

### 生产环境注意事项

⚠️ **重要**: 生产环境必须修改以下默认值：

- `JWT_SECRET` - 必须使用强随机密钥
- `EMQX_DASHBOARD__DEFAULT_PASSWORD` - Dashboard 密码
- `POSTGRES_PASSWORD` - 数据库密码

---

## 服务健康检查

### 手动检查

```bash
# Backend
curl http://localhost:8080/health

# Frontend
curl http://localhost:80/health

# PostgreSQL
docker exec mdm_postgres pg_isready -U mdm_user

# Redis
docker exec mdm_redis redis-cli ping

# EMQX
docker exec mdm_emqx emqx ctl status
```

### 健康检查脚本

```bash
#!/bin/bash
echo "=== MDM Service Health Check ==="

check_service() {
    local name=$1
    local cmd=$2
    echo -n "$name: "
    if eval "$cmd" > /dev/null 2>&1; then
        echo "✅ OK"
    else
        echo "❌ FAIL"
    fi
}

check_service "Backend" "curl -sf http://localhost:8080/health"
check_service "Frontend" "curl -sf http://localhost:80/health"
check_service "PostgreSQL" "docker exec mdm_postgres pg_isready -U mdm_user"
check_service "Redis" "docker exec mdm_redis redis-cli ping"
check_service "EMQX" "docker exec mdm_emqx emqx ctl status"
```

---

## 故障排查

### 查看实时日志

```bash
docker compose logs -f --tail=100
```

### 重启单个服务

```bash
docker compose restart mdm_backend
```

### 重建服务

```bash
docker compose up -d --force-recreate mdm_backend
```

### 清理并重新开始

```bash
docker compose down -v  # 删除数据卷
docker compose up -d     # 重新启动
```

### 常见问题

| 问题 | 解决方案 |
|------|----------|
| Backend 无法连接 PostgreSQL | 检查 `DATABASE_URL` 环境变量，确认 PostgreSQL 已启动 |
| 设备无法连接 MQTT | 检查 EMQX 是否正常运行，端口 1883 是否开放 |
| 前端 502 | 确认 Backend 服务已启动且健康检查通过 |
| 镜像拉取失败 | 检查网络连接，执行 `docker login ghcr.io` |

---

## 维护命令速查

```bash
# 启动
docker compose up -d

# 停止
docker compose down

# 重启
docker compose restart

# 重建
docker compose up -d --build

# 查看状态
docker compose ps

# 查看日志
docker compose logs -f

# 进入容器
docker exec -it mdm_backend /bin/sh
docker exec -it mdm_postgres psql -U mdm_user mdm_db
docker exec -it mdm_redis redis-cli
```
