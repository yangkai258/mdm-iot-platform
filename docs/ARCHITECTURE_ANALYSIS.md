# MDM 项目架构分析报告

**分析日期:** 2026-03-20  
**分析角色:** 系统架构师 (agentyw)  
**项目版本:** V1.0

---

## 一、系统架构审查

### 1.1 技术栈总览

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend                                │
│              Vue3 + TypeScript + ArcoDesign                 │
│                   (Nginx 静态部署)                          │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST
┌────────────────────────▼────────────────────────────────────┐
│                      Backend                                │
│                 Go 1.26 + Gin 1.12                          │
│  ┌──────────┬──────────┬──────────┬──────────┐             │
│  │Controllers│ Middleware│  Models  │   MQTT   │             │
│  └─────┬────┘──────────┴────┬─────┴────┬─────┘             │
│        │                     │          │                   │
│  ┌─────▼─────────────────────▼──────────▼─────┐             │
│  │              GORM (PostgreSQL)               │             │
│  └─────────────────────────────────────────────┘             │
└──────┬──────────────┬──────────────┬─────────────────────────┘
       │              │              │
  ┌────▼────┐   ┌─────▼─────┐ ┌─────▼──────┐
  │PostgreSQL│   │   Redis   │ │  EMQX 5.x  │
  │ (数据)   │   │ (设备影子) │ │   (MQTT)   │
  └─────────┘   └───────────┘ └─────┬──────┘
                                     │
                               ┌─────▼──────┐
                               │ M5Stack    │
                               │ 设备端     │
                               └────────────┘
```

**技术选型评估:**

| 组件 | 选型 | 评分 | 说明 |
|------|------|------|------|
| 后端框架 | Go + Gin | ⭐⭐⭐⭐ | 高性能，协程友好，生态成熟 |
| ORM | GORM | ⭐⭐⭐ | 易用但缺编译期安全，适合快速开发 |
| 数据库 | PostgreSQL 15 | ⭐⭐⭐⭐ | 适合 IoT 时序+关系混合场景 |
| 缓存 | Redis 7 | ⭐⭐⭐⭐ | 设备影子模式正确 |
| MQTT Broker | EMQX 5.x | ⭐⭐⭐⭐⭐ | 国产顶级，水平扩展能力强 |
| 前端 | Vue3 + ArcoDesign | ⭐⭐⭐⭐ | 企业级 UI 组件，适合中台 |
| 容器化 | Docker | ⭐⭐⭐⭐ | 标准部署方式 |

### 1.2 分层架构分析

**当前分层结构:**
```
路由层 (r.GET / r.POST)
    ↓
中间件层 (JWTAuth, OperationLog, CORS)
    ↓
控制器层 (DeviceController, MemberController...)
    ↓
模型层 (GORM struct + Database)
```

**问题: 缺少 Service 层 (业务逻辑层)**

当前所有业务逻辑直接写在 Controller 中，违反了分层架构原则:

```go
// ❌ 当前: Controller 直接操作 DB
func (c *DeviceController) Register(ctx *gin.Context) {
    // 校验、创建、更新...全部在 Controller
    result := c.DB.Where("mac_address = ?", req.MacAddress).First(&device)
    ...
}

// ✅ 推荐: 引入 Service 层
func (c *DeviceController) Register(ctx *gin.Context) {
    device, err := c.deviceService.Register(ctx, req)
    ...
}
```

**缺少 Service 层的后果:**
- 业务逻辑无法复用
- 单元测试困难
- 代码膨胀后难以维护
- 跨模块调用混乱 (Controller A 调用 Controller B)

### 1.3 模块边界分析

**问题: 设备管理模块代码分散在多个文件**

```
controllers/
  device_controller.go      ← Register, Bind, List, Get (DeviceController)
  device_crud.go           ← Delete, Update (独立函数，非方法)
  device_management_controller.go  ← 未知，可能有重叠
```

- `DeleteDevice` 和 `UpdateDevice` 是独立函数，不属于任何结构体
- `DeviceController` 是值类型方法，但 `DeleteDevice` / `UpdateDevice` 是函数式中间件
- 路由注册风格不统一

**会员模块与设备模块完全隔离:**
- 会员模块 (`member_controller.go`) 没有任何设备关联逻辑
- 设备模块没有会员关联逻辑
- 两者通过 `bind_user_id` 外键关联，但没有级联操作

---

## 二、数据流分析

### 2.1 设备数据流

```
设备上电
   │
   ▼
POST /api/v1/devices/register (首次，无需 Token)
   │
   ▼
DeviceController.Register()
   │
   ├─→ 写 PostgreSQL (devices 表, lifecycle_status=1)
   │
   ▼
用户扫码 → POST /api/v1/devices/bind/:sn_code
   │
   ▼
DeviceController.Bind()
   │
   ├─→ 写 PostgreSQL (bind_user_id, lifecycle_status=2)
   │     (❌ 未同步更新 Redis 设备影子)
   │
   ▼
设备 MQTT 连接 → /mdm/device/{id}/up/status (心跳 30s)
   │
   ▼
MQTT Handler.StatusMessageHandler()
   │
   ├─→ 解析 JSON payload
   ├─→ 写 Redis shadow:{device_id} (TTL 90s)
   │     device_id, is_online, battery_level, current_mode, last_heartbeat
   │
   ▼
前端 GET /api/v1/devices/list
   │
   ▼
DeviceController.List()
   │
   ├─→ 查 PostgreSQL (分页列表)
   └─→ 逐个查 Redis GetDeviceShadow (❌ N+1 问题)
```

### 2.2 会员数据流

```
会员注册/创建
   │
   ▼
MemberController.MemberCreate()
   │
   ├─→ 写 PostgreSQL (members 表)
   ├─→ 可选: 关联会员卡 (card_id FK)
   │
   ▼
会员消费 → MemberOrder
   │
   ├─→ 写 MemberOrder 表
   ├─→ 增减积分 (MemberPointsRecord)
   ├─→ 检查升级规则 (❌ 无自动升级逻辑)
```

**会员模块问题:**
- 无缓存层 (会员数据每次 HTTP 请求都查 DB)
- 积分规则计算依赖人工触发，无自动引擎
- 无消息队列，消费记录异步处理

### 2.3 MQTT 消息流

```
                    ┌──────────────┐
                    │   EMQX 5.x   │
                    └──────┬───────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   /up/status    /up/property    /down/cmd (发布)
   (订阅)         (订阅)          (发布)
        │                  │                  │
        ▼                  ▼                  ▼
 StatusMessageHandler  PropertyHandler  PublishCommand
        │                                    
        ▼                                    
  Redis shadow                           
  (TTL 90s)                              

  ┌─────────────────────────────────────────┐
  │  心跳检查器 (30s 间隔轮询所有 shadow:*)  │
  │  → 超时 90s → shadow.IsOnline=false     │
  └─────────────────────────────────────────┘
```

---

## 三、扩展性分析

### 3.1 水平扩展能力

| 组件 | 当前单点? | 扩展方案 |
|------|-----------|----------|
| Backend | ✅ 单实例 | K8s HPA / 多副本 + LB |
| PostgreSQL | ✅ 单节点 | 主从流复制 → PgBouncer 连接池 |
| Redis | ✅ 单节点 | Redis Cluster / 主从 |
| EMQX | ✅ 单节点 | EMQX Cluster (水平扩展) |
| Frontend | ✅ 静态 | CDN 加速 |

**EMQX Cluster 已就绪:** `docker-compose.prod.yml` 使用 `emqx/emqx:5.3`，支持集群模式。

### 3.2 垂直扩展能力

**PostgreSQL 连接池配置 (main.go → utils/redis.go):**
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)   // ← 合理，但可按需调大
sqlDB.SetConnMaxLifetime(time.Hour)
```

**Redis 连接:** 单客户端，无连接池配置。

### 3.3 瓶颈识别

#### 🔴 严重瓶颈

**1. Redis KEYS 命令在心跳检查中被滥用**

```go
// utils/redis.go
func (r *RedisClient) GetAllShadowKeys() ([]string, error) {
    ctx := context.Background()
    keys, err := r.client.Keys(ctx, "shadow:*").Result()  // ❌ O(N) 阻塞命令
    ...
}
```

`KEYS shadow:*` 会遍历整个 Redis keyspace，在 10 万设备时可能导致 Redis 阻塞数秒。

**解决方案:** 改用 `SCAN` 或维护设备列表的 Redis Set。

---

**2. Redis URL 解析逻辑完全失效**

```go
// utils/redis.go InitRedis()
addr := "localhost:6379"   // ← 被硬编码，忽略 REDIS_URL
password := ""
db := 0
fmt.Sscanf(redisURL, "redis://%*s@%s", &addr)  // ← fmt.Sscanf 解析不了 URL
```

生产环境 `REDIS_URL=redis://:redis_password@redis:6379` 会被完全忽略，始终连接 `localhost:6379`。

---

**3. 设备影子与数据库状态不同步**

当设备离线超过 90s，`Redis shadow.IsOnline=false`，但 PostgreSQL `devices.lifecycle_status` 不会更新。
设备列表查询时，Redis 和 DB 数据不一致，导致管理后台看到设备"离线"但 DB 状态仍是"服役中"。

---

#### 🟡 中等瓶颈

**4. N+1 查询问题**

```go
// device_controller.go List()
for i, d := range devices {
    shadow, err := c.Redis.GetDeviceShadow(d.DeviceID)  // N 次 Redis 调用
    ...
}
```

100 台设备 = 1 次 DB + 100 次 Redis GET。

**解决方案:** 使用 Redis `MGET` 批量获取，或 Lua 脚本一次获取。

---

**5. 操作日志 goroutine 泄漏风险**

```go
go func() {
    if err := db.Create(&logEntry).Error; err != nil {
        log.Printf("Failed to create operation log: %v", err)
    }
}()
```

无缓冲 channel，无限创建 goroutine，高并发时可能导致 OOM。

---

**6. 无 MQTT 重连机制**

```go
// mqtt/handler.go InitMQTT()
client := mqtt.NewClient(opts)
token := client.Connect()
```

EMQX 连接断开后无自动重连，paho.mqtt.golang 默认有重连但没有指数退避。

---

#### 🟢 轻微瓶颈

**7. 前端 dist 目录无法通过 docker-compose 构建**

`docker-compose.prod.yml` 中:
```yaml
frontend:
  image: nginx:alpine
  volumes:
    - ./frontend/dist:/usr/share/nginx/html:ro  # ← dist 需提前 build
```

没有 `docker build` 步骤，需要手动 `npm run build`。

---

## 四、安全分析

### 4.1 认证/授权

#### 🔴 JWT 密钥硬编码

```go
// middleware/jwt.go
var jwtSecret = []byte("mdm-secret-key-change-in-production")  // ❌ 危险
```

密钥写在源代码中，若代码泄露则所有 Token 可被伪造。

**修复:** 使用环境变量 `os.Getenv("JWT_SECRET")` 或 K8s Secret。

---

#### 🔴 CORS 全通配

```go
// main.go
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  // ❌ 生产环境危险
```

允许任何来源访问 API，无法防止 CSRF。

---

#### 🟡 无 Token 刷新机制

JWT 有效期 24 小时，无 refresh_token 端点。用户需每 24 小时重新登录。

---

#### 🟡 RBAC 表存在但未执行

- `SysRole`, `SysMenu`, `SysPermission` 表已定义
- 但 Controller 层无权限校验，所有登录用户可访问所有 API
- `jwt.go` 只提取 `role_id` 存入 context，从未使用

---

#### 🟢 无登录限流

登录接口 `/api/v1/auth/login` 无 rate limiting，可被暴力破解。

---

### 4.2 数据安全

#### 🟡 密码加密合规

使用 `golang.org/x/crypto/bcrypt`，符合安全规范。

---

#### 🟡 数据库凭证硬编码在 docker-compose

```yaml
# docker-compose.prod.yml
postgres:
  POSTGRES_PASSWORD: mdm_password   # ← 硬编码
redis:
  command: redis-server --requirepass redis_password  # ← 硬编码
```

应使用 `.env` 文件或 K8s Secret。

---

#### 🟢 敏感数据无脱敏

操作日志记录完整请求 body：
```go
if len(bodyBytes) > 2000 {
    logEntry.Params = string(bodyBytes[:2000]) + "...(truncated)"
}
```

密码等敏感字段未被过滤。

---

### 4.3 通信安全

#### 🟡 生产环境无 TLS

- PostgreSQL: `sslmode=disable`
- MQTT: `tcp://` (无 TLS)
- Frontend → Backend: HTTP (无 HTTPS)

---

#### 🟡 MQTT 认证仅靠用户名密码

EMQX 配置了 `EMQX_AUTH__SIMPLE__DEFAULT_PASSWORD`，但无 ACL 规则。设备间可能互相订阅主题。

---

## 五、架构问题清单

### 严重程度分级

| 级别 | 数量 | 说明 |
|------|------|------|
| 🔴 严重 (必须修复) | 4 | 影响安全或可用性 |
| 🟡 中等 (近期修复) | 9 | 影响性能或可维护性 |
| 🟢 轻微 (规划修复) | 6 | 最佳实践优化 |

### 🔴 严重问题

| # | 问题 | 位置 | 影响 |
|---|------|------|------|
| S1 | JWT 密钥硬编码在源码 | `middleware/jwt.go:14` | 安全: Token 可被伪造 |
| S2 | Redis URL 解析完全失效 | `utils/redis.go:InitRedis()` | 功能: 生产环境 Redis 永远连 localhost |
| S3 | CORS Allow-Origin 全通配 | `main.go` | 安全: 任意站点的 CSRF 攻击 |
| S4 | 设备影子与 DB 状态不同步 | MQTT Handler / DeviceController | 功能: 离线设备 DB 仍显示在线 |

### 🟡 中等问题

| # | 问题 | 位置 | 影响 |
|---|------|------|------|
| M1 | Redis KEYS 命令阻塞心跳检查 | `utils/redis.go:GetAllShadowKeys()` | 性能: 10万设备时 Redis 可阻塞数秒 |
| M2 | N+1 查询设备影子 | `device_controller.go:List()` | 性能: 100设备=100次Redis调用 |
| M3 | 无 Service 层 | 整体架构 | 可维护性: 业务逻辑散落各处 |
| M4 | 操作日志 goroutine 泄漏风险 | `middleware/operation_log.go` | 稳定性: 高并发时 OOM |
| M5 | MQTT 无重连策略 | `mqtt/handler.go:InitMQTT()` | 可用性: 网络抖动后无自动恢复 |
| M6 | RBAC 表存在但未执行 | Controller 层 | 安全: 所有登录用户等权限 |
| M7 | 无 Token 刷新机制 | `middleware/jwt.go` | UX: 每24小时强制重新登录 |
| M8 | 无登录限流 | `controllers/auth_controller.go` | 安全: 暴力破解风险 |
| M9 | 数据库凭证硬编码 | `docker-compose.prod.yml` | 安全: 泄露风险 |

### 🟢 轻微问题

| # | 问题 | 位置 | 影响 |
|---|------|------|------|
| L1 | 前端 dist 需手动构建 | `docker-compose.prod.yml` | DevOps: 无自动化构建 |
| L2 | 设备管理代码分散 3 个文件 | `controllers/device_*.go` | 可维护性 |
| L3 | 会员模块与设备模块完全隔离 | 业务设计 | 功能: 无法按设备查会员 |
| L4 | PostgreSQL sslmode=disable | `utils/redis.go` | 安全: 网络层明文传输 |
| L5 | MQTT 无 TLS | `mqtt/handler.go` | 安全: 网络层明文 |
| L6 | 无数据库迁移工具 | 整体 | DevOps: schema 管理混乱 |

---

## 六、优化建议

### 6.1 第一阶段: 紧急修复 (1-2 周)

**S2: 修复 Redis URL 解析**
```go
// utils/redis.go
func InitRedis() (*RedisClient, error) {
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        redisURL = "redis://localhost:6379"
    }
    
    // 使用 net/url 正确解析 redis:// scheme
    u, err := url.Parse(redisURL)
    if err != nil {
        return nil, err
    }
    
    addr := u.Host
    if u.User != nil {
        password, _ = u.User.Password()
    }
    
    return NewRedisClient(addr, password, 0)
}
```

**S4: 设备影子与 DB 同步**

心跳超时后，更新 DB 中设备的最后在线时间:
```go
// mqtt/handler.go checkOfflineDevices()
if elapsed > 90*time.Second && shadow.IsOnline {
    shadow.IsOnline = false
    h.Redis.SetDeviceShadow(shadow.DeviceID, *shadow, 0)
    
    // 异步更新 DB
    go func() {
        h.DB.Model(&models.Device{}).
            Where("device_id = ?", shadow.DeviceID).
            Update("last_heartbeat", shadow.LastHeartbeat)
    }()
}
```

**S1: JWT 密钥环境变量化**
```go
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
if string(jwtSecret) == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
```

**S3: CORS 配置化**
```go
allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGINS")
if allowedOrigin == "" {
    allowedOrigin = "*"  // Dev 默认
}
r.Use(cors.New(cors.Config{
    AllowOrigins:     strings.Split(allowedOrigin, ","),
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,  // 配合具体域名
}))
```

### 6.2 第二阶段: 性能优化 (2-4 周)

**M1: 替换 KEYS 为 SCAN + Redis Set**

维护设备列表:
```go
// 设备注册时加入 Set
r.client.SAdd(ctx, "device_ids", deviceID)

// 心跳检查时使用 SSCAN
iter := r.client.SScan(ctx, "device_ids", 0, "shadow:*", 100).Iterator()
```

**M2: 批量获取设备影子**
```go
// List() 中使用 MGET
keys := make([]string, len(devices))
for i, d := range devices {
    keys[i] = fmt.Sprintf("shadow:%s", d.DeviceID)
}
results, _ := r.client.MGet(ctx, keys...).Result()
```

**M5: MQTT 重连 + 指数退避**
```go
opts.SetAutoReconnect(true)
opts.SetMaxReconnectInterval(60 * time.Second)
```

### 6.3 第三阶段: 架构演进 (1-2 月)

**引入 Service 层**

```
controllers/    →  接收请求，参数校验，调用 Service
services/       →  业务逻辑，事务管理
repositories/   →  数据访问，查询构建
```

```
services/
  device_service.go
  member_service.go
  ota_service.go
```

**引入消息队列 (可选)**

当设备量 > 1 万时，MQTT 消息处理可能成为瓶颈:
```
EMQX → Kafka/RabbitMQ → Backend Workers → DB
         ↑
    异步解耦，削峰填谷
```

---

## 七、演进路线图

```
当前状态                    3个月后                    6个月后                    12个月后
   │                          │                          │                          │
   ▼                          ▼                          ▼                          ▼
┌──────────┐            ┌──────────┐              ┌──────────┐              ┌──────────┐
│ 单体架构 │     →      │  Service │       →      │ 消息队列 │       →      │ 分布式微 │
│ 粗糙分层 │            │   层引入  │              │  解耦    │              │   服务   │
└──────────┘            └──────────┘              └──────────┘              └──────────┘

修复:                    优化:                      扩展:
- JWT 环境变量           - Redis SCAN              - EMQX 集群
- Redis URL 解析         - 批量 MGET               - PostgreSQL 主从
- CORS 配置化            - 操作日志 channel         - Redis Cluster
- 影子/DB 同步           - MQTT 重连策略            - CDN 加速前端
                        - RBAC 权限执行             - K8s 部署

安全:                    监控:                      高级:
- 登录限流               - Prometheus + Grafana    - OTA 灰度发布
- Token 刷新             - EMQX 监控                - 设备分组策略
- TLS 加密               - 日志聚合                 - 多租户隔离
```

---

## 八、总结

MDM 项目整体架构**方向正确**:
- MQTT + 设备影子的架构设计合理
- PostgreSQL + Redis 分层存储符合 IoT 场景
- EMQX 作为 MQTT Broker 具备企业级扩展能力
- Vue3 + ArcoDesign 前端技术栈选型适合管理中台

但存在**分层不清晰**和**关键路径有安全隐患**的问题，需要尽快修复 S1-S4 和 M1-M2，尤其 Redis URL 解析失效和设备影子不同步是生产环境运行的**阻断性问题**。

Service 层的引入是中期最重要的架构改进，它将为未来的功能扩展和团队协作奠定基础。
