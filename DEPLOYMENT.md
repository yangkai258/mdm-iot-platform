# MDM 控制中台 - 部署说明文档

> 本文档涵盖开发环境、GitHub Actions CI/CD、以及生产环境部署流程。

---

## 目录

- [1. 环境要求](#1-环境要求)
- [2. 项目结构](#2-项目结构)
- [3. 本地开发部署](#3-本地开发部署)
- [4. Docker Compose 一键启动](#4-docker-compose-一键启动)
- [5. GitHub Actions CI/CD](#5-github-actions-cicd)
- [6. 生产环境部署](#6-生产环境部署)
- [7. 健康检查与监控](#7-健康检查与监控)
- [8. 常见问题](#8-常见问题)

---

## 1. 环境要求

| 组件 | 版本 | 说明 |
|------|------|------|
| Go | ≥ 1.23 | 后端服务 |
| Node.js | ≥ 20 | 前端构建 |
| Docker | ≥ 24.0 | 容器运行时 |
| Docker Compose | ≥ 2.20 | 多服务编排 |
| PostgreSQL | 15 | 主数据库 |
| Redis | 7 | 缓存/设备影子 |
| EMQX | 5.x | MQTT Broker |

---

## 2. 项目结构

```
mdm-project/
├── .github/
│   └── workflows/
│       └── ci-cd.yml          # GitHub Actions CI/CD 流水线
├── backend/                   # (已废弃，后端源码在根目录)
├── frontend/
│   ├── src/                   # Vue3 源码
│   ├── dist/                  # 构建产物（npm run build 生成）
│   ├── Dockerfile             # 前端镜像构建
│   └── nginx.conf             # Nginx 配置
├── backend/                   # (已废弃)
├── main.go                    # 后端入口 (Go)
├── go.mod                     # Go 模块定义
├── Dockerfile                 # 后端多阶段构建镜像
├── docker-compose.yml         # 开发环境编排
├── docker-compose.prod.yml    # 生产环境编排
├── vite.config.js             # Vite 构建配置
└── package.json               # 前端依赖定义
```

---

## 3. 本地开发部署

### 3.1 前置条件

```bash
# 安装 Go
go version  # 确认 >= 1.23

# 安装 Node.js
node --version  # 确认 >= 20

# 安装 Docker Desktop (Windows) 或 Docker Engine (Linux)
docker --version
docker-compose --version
```

### 3.2 启动基础设施服务

```bash
# 仅启动数据库、Redis、EMQX（不启动业务服务）
docker-compose up -d postgres redis emqx

# 验证服务健康
docker-compose ps
```

### 3.3 启动后端

```bash
# 下载 Go 依赖
go mod download

# 运行后端（开发模式）
go run main.go

# 或者构建并运行
go build -o mdm-server main.go
./mdm-server
```

后端默认监听 `http://localhost:8080`。

### 3.4 启动前端

```bash
# 安装前端依赖
npm install

# 开发模式热重载
npm run dev

# 生产构建
npm run build
```

前端默认监听 `http://localhost:3000`，API 代理到 `http://localhost:16666`（旧端口）。

---

## 4. Docker Compose 一键启动

### 4.1 开发环境

```bash
# 启动全部服务（包含健康检查）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f mdm_backend
docker-compose logs -f mdm_frontend

# 停止全部服务
docker-compose down

# 停止并清除数据卷（重置数据库）
docker-compose down -v
```

### 4.2 健康检查说明

| 服务 | 健康检查方式 | 端口 |
|------|-------------|------|
| `postgres` | `pg_isready -U mdm_user` | 5432 |
| `redis` | `redis-cli ping` | 6379 |
| `emqx` | `emqx ctl status` | 18083 |
| `mdm_backend` | `curl http://localhost:8080/health` | 8080 |
| `mdm_frontend` | `wget http://localhost:80/health` | 80/443 |

所有服务均配置了 `restart: unless-stopped`，异常退出后会自动重启。

### 4.3 资源限制

- **mdm_backend**: 限制 512MB 内存，1 CPU；预留 128MB，0.25 CPU
- **mdm_frontend**: 无强制限制（静态文件服务，资源占用低）

### 4.4 生产环境

```bash
# 使用生产配置
docker-compose -f docker-compose.prod.yml up -d --build

# 或分步构建
docker build -t mdm-backend:latest -f Dockerfile .
docker build -t mdm-frontend:latest -f frontend/Dockerfile ./frontend
docker-compose -f docker-compose.prod.yml up -d
```

---

## 5. GitHub Actions CI/CD

### 5.1 流水线概览

```
触发条件: push/PR 到 main/master/develop 分支

阶段一: 并行构建
├── build-backend   → Go 编译 + go vet + 单元测试
└── build-frontend  → npm ci + ESLint + Vite 构建

阶段二: Docker 镜像构建
└── build-docker    → 构建并推送 ghcr.io 镜像

阶段三: 集成测试
└── integration-test → 启动 Postgres/Redis/EMQX 服务，跑 API Smoke Test

阶段四: 部署（仅 push 到 main/master）
└── deploy-staging   → SSH 到预发布服务器执行 docker compose pull && up
```

### 5.2 流水线文件

路径: `.github/workflows/ci-cd.yml`

### 5.3 所需 Secrets（仓库 Settings → Secrets）

| Secret | 说明 | 示例 |
|--------|------|------|
| `STAGING_SSH_KEY` | 部署服务器 SSH 私钥 | `-----BEGIN OPENSSH...` |
| `STAGING_HOST` | 部署服务器 IP/域名 | `staging.mdm.example.com` |
| `STAGING_USER` | 部署服务器用户名 | `deploy` |

### 5.4 查看流水线状态

```
GitHub 仓库 → Actions → 选择 workflow → 查看各 Job 状态
```

### 5.5 发布流程

1. 代码合并到 `main` 分支
2. GitHub Actions 自动触发构建
3. 构建成功后自动部署到 Staging
4. 手动审核后，在 GitHub Releases 发布正式版本

---

## 6. 生产环境部署

### 6.1 服务器准备

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh
sudo systemctl enable docker
sudo systemctl start docker

# 安装 Docker Compose v2
sudo apt-get install docker-compose-v2
```

### 6.2 部署步骤

```bash
# 1. SSH 登录服务器
ssh deploy@staging.mdm.example.com

# 2. 创建部署目录
mkdir -p /opt/mdm
cd /opt/mdm

# 3. 克隆（或更新）仓库
git clone https://github.com/yangkai258/mdm-iot-platform.git .
git checkout main

# 4. 创建 .env.production 环境变量文件
cat > .env << EOF
JWT_SECRET=<随机字符串，建议 64 位>
DB_USER=mdm_prod_user
DB_PASSWORD=<强密码>
DB_NAME=mdm_production
MQTT_ADMIN_USER=admin
MQTT_ADMIN_PASSWORD=<强密码>
EOF

# 5. 启动服务
docker-compose -f docker-compose.prod.yml up -d --build

# 6. 验证健康状态
docker-compose ps
curl -sf http://localhost:8080/health
curl -sf http://localhost:80/health
```

### 6.3 生产环境关键配置

```yaml
# docker-compose.prod.yml 关键优化
postgres:
  command: >
    postgres
    -c max_connections=200
    -c shared_buffers=256MB
    -c effective_cache_size=1GB
    -c checkpoint_completion_target=0.9

redis:
  command: >
    redis-server
    --appendonly yes
    --maxmemory 512mb
    --maxmemory-policy allkeys-lru

mdm_backend:
  deploy:
    replicas: 2  # 后端多实例
    resources:
      limits:
        cpus: '1'
        memory: 1G
```

### 6.4 Nginx HTTPS 配置

```nginx
server {
    listen 443 ssl http2;
    server_name mdm.example.com;

    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;

    # 强制 HTTPS
    add_header Strict-Transport-Security "max-age=31536000" always;
}
```

### 6.5 备份策略

```bash
# PostgreSQL 每日备份（crontab）
0 2 * * * docker exec mdm_postgres pg_dump -U mdm_user mdm_db > /backup/mdm_$(date +\%Y\%m\%d).sql

# Redis AOF 持久化（已启用）
# EMQX 数据卷持久化（已配置）
```

---

## 7. 健康检查与监控

### 7.1 Docker 健康检查端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `GET /health` | Backend | 返回 `{"status":"healthy"}` |
| `GET /health` | Frontend | Nginx 返回纯文本 `healthy` |
| `GET /api/v1/health` | Backend | 备用健康检查 |

### 7.2 Docker 健康状态

```bash
# 查看所有容器健康状态
docker inspect --format='{{.Name}}: {{.State.Health.Status}}' $(docker ps -q)

# 查看健康检查失败日志
docker inspect --format='{{.Name}}: {{.State.Health.Log}}' mdm_backend
```

### 7.3 EMQX Dashboard

- 地址: `http://localhost:18083`
- 用户名: `admin`
- 密码: `public`（开发环境，生产环境请修改）

---

## 8. 常见问题

### Q1: 后端启动失败，提示 `connection refused`?

确保 PostgreSQL/Redis/EMQX 已启动并通过健康检查：
```bash
docker-compose up -d postgres redis emqx
sleep 10  # 等待健康检查通过
docker-compose up -d mdm_backend
```

### Q2: 前端构建失败 `npm ERR! peer dep`?

```bash
# 清理 node_modules 后重装
rm -rf node_modules package-lock.json
npm install
```

### Q3: Docker 镜像构建OOM？

```bash
# 增加 Docker 内存限制（Docker Desktop → Settings → Resources）
# 推荐至少 4GB
```

### Q4: GitHub Actions 构建失败？

检查 `ci-cd.yml` 中的 Go 版本和 Node 版本是否与本地一致：
```yaml
GO_VERSION: '1.23'   # 与 go.mod 中 go 1.23 匹配
NODE_VERSION: '20'   # 与 package.json 引擎要求匹配
```

### Q5: `curl` 健康检查在 alpine 镜像中不存在？

Dockerfile 已安装 `ca-certificates` 和 `curl`。如遇问题，确认 `apk add curl` 执行成功。备用方案：
```yaml
healthcheck:
  test: ["CMD-SHELL", "wget -qO- http://localhost:8080/health || exit 1"]
```

---

*本文档由 CI/CD Agent 自动生成，最后更新: 2026-03-24*
