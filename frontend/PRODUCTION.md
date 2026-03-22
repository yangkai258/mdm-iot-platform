# MDM 控制中台 - 生产环境部署方案

## 架构概览

```
                    ┌─────────────────┐
                    │   负载均衡器    │
                    │   (Nginx/Caddy) │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              ▼
        ┌──────────┐  ┌──────────┐  ┌──────────┐
        │ Frontend │  │ Frontend │  │ Frontend │
        │  (Vue3)  │  │  (Vue3)  │  │  (Vue3)  │
        └────┬─────┘  └────┬─────┘  └────┬─────┘
             │             │             │
             └─────────────┼─────────────┘
                           │
                    ┌──────┴──────┐
                    │  Nginx API  │
                    │   Gateway   │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
              ▼            ▼            ▼
        ┌──────────┐  ┌──────────┐  ┌──────────┐
        │ Backend  │  │ Backend  │  │ Backend  │
        │  (Go)    │  │  (Go)    │  │  (Go)    │
        └────┬─────┘  └────┬─────┘  └────┬─────┘
             │            │            │
       ┌─────┴─────┐      │      ┌──────┴──────┐
       │           │      │      │             │
       ▼           ▼      ▼      ▼             ▼
   ┌──────┐  ┌────────┐ ┌────┐ ┌────────┐  ┌────────┐
   │PostgreSQL│  │ Redis │ │MQTT│ │PostgreSQL│  │ Redis │
   │ (主从)  │  │Cluster│ │    │ │ (从)   │  │(副本)  │
   └────────┘  └───────┘ └────┘ └────────┘  └────────┘
```

## 1. Docker Compose 生产配置

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  # PostgreSQL 主库
  postgres:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - mdm_network
    command: >
      postgres
      -c max_connections=200
      -c shared_buffers=256MB
      -c effective_cache_size=1GB
      -c maintenance_work_mem=64MB
      -c checkpoint_completion_target=0.9
      -c wal_buffers=16MB
      -c default_statistics_target=100
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 缓存
  redis:
    image: redis:7-alpine
    restart: always
    command: >
      redis-server
      --appendonly yes
      --maxmemory 512mb
      --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - mdm_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # EMQX MQTT Broker
  emqx:
    image: emqx/emqx:latest
    restart: always
    environment:
      EMQX_NAME: mdm_emqx
      EMQX_HOST: 0.0.0.0
      EMQX_DASHBOARD__DEFAULT_USERNAME: ${MQTT_ADMIN_USER}
      EMQX_DASHBOARD__DEFAULT_PASSWORD: ${MQTT_ADMIN_PASSWORD}
      EMQX_ALLOW_ANONYMOUS: "false"
    volumes:
      - emqx_data:/opt/emqx/data
      - emqx_log:/opt/emqx/log
    ports:
      - "1883:1883"
      - "8083:8083"
      - "18083:18083"
    networks:
      - mdm_network
    healthcheck:
      test: ["CMD", "emqx", "ctl", "status"]
      interval: 10s
      timeout: 5s
      retries: 5

  # 后端服务 (多实例)
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.prod
    restart: always
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
      REDIS_URL: redis://redis:6379
      MQTT_BROKER: mqtt://emqx:1883
      LOG_LEVEL: info
      GIN_MODE: release
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      emqx:
        condition: service_healthy
    networks:
      - mdm_network
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M

  # Nginx 网关
  gateway:
    image: nginx:alpine
    restart: always
    volumes:
      - ./gateway/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./gateway/conf.d:/etc/nginx/conf.d:ro
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - backend
    networks:
      - mdm_network

  # 前端 (可选，也可以用 CDN)
  frontend:
    image: nginx:alpine
    restart: always
    volumes:
      - ./frontend/dist:/usr/share/nginx/html:ro
      - ./frontend/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - gateway
    networks:
      - mdm_network

volumes:
  postgres_data:
  redis_data:
  emqx_data:
  emqx_log:

networks:
  mdm_network:
    driver: bridge
```

## 2. 后端 Dockerfile (生产)

```dockerfile
# backend/Dockerfile.prod
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mdm-server .

# 运行镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/mdm-server .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

USER nonroot:nonroot

CMD ["./mdm-server"]
```

## 3. 环境变量配置

```bash
# .env.production
# 数据库
DB_USER=mdm_prod_user
DB_PASSWORD=<secure-password>
DB_NAME=mdm_production

# MQTT
MQTT_ADMIN_USER=admin
MQTT_ADMIN_PASSWORD=<secure-password>

# JWT
JWT_SECRET=<very-secure-random-string>

# Redis
REDIS_PASSWORD=<secure-password>

# 域名
DOMAIN=mdm.yourdomain.com
```

## 4. Nginx 网关配置

```nginx
# gateway/nginx.conf
events {
    worker_connections 2048;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript 
               application/javascript application/xml+rss 
               application/json application/x-javascript;

    # 上游后端
    upstream backend_servers {
        least_conn;
        server backend:8080 max_fails=3 fail_timeout=30s;
        keepalive 32;
    }

    server {
        listen 80;
        server_name mdm.yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name mdm.yourdomain.com;

        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
        ssl_prefer_server_ciphers off;

        # 前端静态文件
        location / {
            root /usr/share/nginx/html;
            index index.html;
            try_files $uri $uri/ /index.html;
            
            # 缓存
            location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
                expires 1y;
                add_header Cache-Control "public, immutable";
            }
        }

        # API 代理
        location /api/ {
            proxy_pass http://backend_servers;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_connect_timeout 30s;
            proxy_send_timeout 30s;
            proxy_read_timeout 30s;
        }

        # WebSocket
        location /ws/ {
            proxy_pass http://backend_servers;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_read_timeout 86400;
        }
    }
}
```

## 5. 部署脚本

```bash
#!/bin/bash
# deploy.sh

set -e

echo "🚀 开始部署 MDM 控制中台..."

# 1. 拉取最新代码
git pull origin main

# 2. 构建前端
cd frontend
npm install
npm run build
cd ..

# 3. 构建后端
cd backend
go build -o mdm-server .
cd ..

# 4. 启动服务
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build

# 5. 检查状态
sleep 10
docker-compose -f docker-compose.prod.yml ps

echo "✅ 部署完成！"
echo "访问 https://mdm.yourdomain.com"
```

## 6. 监控配置

```yaml
# docker-compose.monitoring.yml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3001:3000"
    volumes:
      - ./monitoring/grafana/dashboards:/etc/grafana/provisioning/dashboards
```

## 7. 关键优化点

| 项目 | 开发环境 | 生产环境 |
|------|----------|----------|
| 后端实例数 | 1 | 2-3 |
| 日志 | stdout | 文件 + 结构化 |
| 缓存 | 256MB | 512MB-1GB |
| 连接池 | 默认 | 50-100 |
| SSL | 无 | 必须 |
| 监控 | 无 | Prometheus+Grafana |
| 备份 | 手动 | 自动化 |

---

**当前开发环境状态**：正常运行，无需修改  
**生产部署建议**：使用上述配置替换 docker-compose.yml
