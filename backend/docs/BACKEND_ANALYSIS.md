# MDM 后端架构分析报告

> 分析时间：2026-03-20  
> 技术栈：Go 1.26 + Gin 1.12 + GORM + PostgreSQL + Redis + MQTT (paho.mqtt.golang)  
> 分析范围：models/、controllers/、middleware/、mqtt/、utils/、main.go

---

## 一、架构问题清单

### 🔴 严重问题 (P0)

| # | 问题 | 位置 | 说明 |
|---|------|------|------|
| 1 | **硬编码 JWT 密钥** | `middleware/jwt.go:12` | `jwtSecret = []byte("mdm-secret-key-change-in-production")` 必须改为从环境变量或配置中心读取 |
| 2 | **MQTT CommandController 未注入** | `main.go:63` | `MQTT: nil` 导致指令下发功能完全失效 |
| 3 | **CORS 完全开放** | `main.go:46` | `Access-Control-Allow-Origin: *` 生产环境必须限制来源域名 |
| 4 | **JWT 无刷新机制** | `middleware/jwt.go` | Token 24h 过期后用户需重新登录，无 refresh_token 机制 |
| 5 | **密码未加密存储** | `auth_controller.go` | SysUser.Password 字段声明了 `gorm:"-"` 不返回，但创建用户时明文存储风险 |

### 🟠 高优先级 (P1)

| # | 问题 | 位置 | 说明 |
|---|------|------|------|
| 6 | **无 Service 层** | 全部 controllers/ | 所有业务逻辑直接写在 Controller 中，违反单一职责原则 |
| 7 | **无 DTO/Request 分离** | 全部 controllers/ | Request struct 内嵌在 controller 文件中，重复定义无法复用 |
| 8 | **无统一错误响应** | 全部 controllers/ | 有的返回 `{"code": 400, "message": ""}`，有的返回 `{"code": 4005, ...}`，无统一标准 |
| 9 | **设备列表内存过滤在线状态** | `device_controller.go:155` | 先查全量设备到内存，再过滤 online/offline，大数据量必崩 |
| 10 | **操作日志 goroutine 无错误回调** | `middleware/operation_log.go:68` | `go func()` 内 db.Create 失败仅打印 log，日志可能丢失 |
| 11 | **Redis 连接无连接池配置** | `utils/redis.go` | 未设置 `PoolSize`、`MinIdleConns` 等参数 |
| 12 | **缺少请求 ID 追踪** | 全部 | 无法在日志中追踪完整请求链路 |

### 🟡 中优先级 (P2)

| # | 问题 | 位置 | 说明 |
|---|------|------|------|
| 13 | **路由分散在两处** | `main.go` + `device_controller.go` | 系统管理路由在 main.go，业务路由在 RegisterRoutes，难以维护 |
| 14 | **无数据库事务管理** | 全部写操作 | 涉及多表操作（如会员注册送积分）无事务保护 |
| 15 | **DeviceShadow 未持久化** | `mqtt/handler.go` | Redis 宕机后设备影子数据全丢失，无降级方案 |
| 16 | **宠物模型与会员模块混入设备项目** | `models/` | PetProfile、MemberOrder 等属于不同业务域，混淆了 MDM 核心职责 |
| 17 | **缺少业务层配置中心** | `utils/redis.go` | 数据库/Redis 连接参数从环境变量读取，但无结构化配置管理 |
| 18 | **未使用 GORM Callbacks 审计字段** | `models/` | CreatedAt/UpdatedAt 依赖 GORM 默认行为，无自定义审计逻辑 |
| 19 | **部门树递归无缓存** | `org_controller.go` | 每次请求都从数据库重新构建树形结构 |
| 20 | **心跳检测遍历所有 shadow key** | `mqtt/handler.go:134` | `h.Redis.GetAllShadowKeys()` 无分片，设备量过万时阻塞 |

---

## 二、代码问题清单

### 2.1 代码规范问题

| 文件 | 行 | 问题 |
|------|-----|------|
| `models/device.go` | 19 | `BeforeCreate` 使用全局 UUID 生成器，高并发下无 Snowflake 优化 |
| `controllers/auth_controller.go` | 50 | 登录日志记录登录成功，但退出登录时从 context 断言 userID 可能 panic |
| `controllers/auth_controller.go` | 55 | `ctx.Get("user_id")` 返回 `interface{}`，未做类型校验 |
| `controllers/device_controller.go` | 56 | `gorm.ErrRecordNotFound` 判断后继续执行 `Save()` 无 else |
| `controllers/device_crud.go` | 37 | `ShouldBindJSON` 后直接覆盖 model 原字段，忽略 `BindUserID` 等已有值 |
| `controllers/command_controller.go` | 42 | MQTT Publish QoS=0，指令可靠性无保障 |
| `controllers/org_controller.go` | 87 | 递归构建部门树时使用 GORM 查询结果遍历，未使用指针效率低 |
| `middleware/operation_log.go` | 34 | `io.ReadAll` 读取 body 后 `Request.Body` 被消耗，后续 handler 读不到 body |
| `utils/redis.go` | 31 | `fmt.Sscanf(redisURL, ...)` 解析 Redis URL 未处理错误 |
| `models/org_models.go` | 35 | `Children []Department` 使用 `gorm:"-"` 忽略标签，混入了 Model 定义 |

### 2.2 安全隐患

| 风险 | 位置 | 描述 |
|------|------|------|
| SQL 注入 | `device_controller.go:125` | `query.Where("device_id LIKE ? OR sn_code LIKE ?", search, search)` 虽使用参数化，但 search 本身可能含 `%`/`_` 需预处理 |
| 水平越权 | `command_controller.go` | 未校验请求发起人是否有权操作目标 device_id |
| 垂直越权 | `middleware/jwt.go` | JWT 只存 role_id，但无接口级别的权限校验 |
| 敏感信息泄露 | `operation_log.go` | 请求参数直接记录日志，可能含密码等敏感字段 |
| JWT 密钥泄露 | `middleware/jwt.go:12` | 密钥硬编码在代码中，泄露后无法轮换 |

### 2.3 性能瓶颈

| 场景 | 位置 | 问题 |
|------|------|------|
| 设备列表 + 在线状态 | `device_controller.go:143-163` | N+1 查询：循环内逐个调 `GetDeviceShadow` |
| 在线设备筛选 | `device_controller.go:155` | 全量查 DB 再内存过滤，PageSize=100 但 total 可能上万 |
| 离线检测 | `mqtt/handler.go:134` | 每 30s 全量扫描所有 shadow key，O(n) 扫描 |
| 会员列表 | `member_controller.go` | `Preload("Card")` 对每条记录都 JOIN，未做预加载优化 |

---

## 三、API 质量审查

### 3.1 RESTful 规范问题

| 接口 | 问题 |
|------|------|
| `POST /devices/bind/:sn_code` | 绑定操作应使用 `PUT /devices/:sn_code/bind` |
| `POST /devices/unbind/:sn_code` | 解绑操作应使用 `DELETE /devices/:sn_code/bind` |
| `GET /dashboard/stats` | stats 应为名词，`GET /dashboard/statistics` |
| `GET /devices/:device_id/profile` | profile 是设备子资源，应为 `GET /devices/:device_id/pet-profile` |
| `PUT /devices/:device_id/profile` | 同上 |
| `POST /devices/:device_id/commands` | 指令是设备子资源，符合 RESTful |
| `GET /ota/devices/:device_id/check` | 检查 OTA 应为 `GET /ota/check/:device_id` 或查询参数 |

### 3.2 错误响应不一致

**统一响应格式缺失**，当前混用以下格式：

```json
// 格式 A
{"code": 0, "message": "success", "data": {...}}

// 格式 B
{"code": 0, "data": {...}}

// 格式 C
{"code": 400, "message": "参数错误"}

// 格式 D
{"code": 4005, "message": "...", "error_code": "ERR_VALIDATION"}
```

**HTTP 状态码使用混乱**：
- 参数错误：有的返回 200（业务错误），有的返回 400（语义错误）
- 建议：业务逻辑错误返回 200 + 业务 code，HTTP 状态码专用于传输层

### 3.3 数据验证缺失

| 字段 | 验证缺失 |
|------|----------|
| `Member.phone` | 未校验手机号格式 |
| `Member.email` | 未校验邮箱格式 |
| `MemberOrder.order_no` | 未校验唯一性冲突处理 |
| `Device.sn_code` | 格式校验仅靠 binding tag |
| `SysUser.username` | 未限制用户名长度/字符集 |
| `OTAPackage.md5_hash` | 未校验 MD5 固定 32 位格式 |

---

## 四、优化建议

### 4.1 架构重构建议

```
当前: main.go → Controllers → GORM (全在 Controller 层)
建议: main.go → Handlers → Services → Repositories → GORM
```

**分层职责**：
- **Handler**：参数解析、响应组装、HTTP 状态码
- **Service**：业务逻辑、事务管理、领域规则
- **Repository**：数据访问、查询构造、缓存读写

**推荐新增目录结构**：
```
backend/
├── config/          # 配置加载（Viper/ENV）
├── internal/         # 内部业务
│   ├── handler/     # HTTP handlers
│   ├── service/     # 业务逻辑
│   ├── repository/  # 数据访问
│   └── dto/         # 请求/响应 DTO
├── pkg/             # 可独立发布的包（mqtt, redis等）
└── main.go
```

### 4.2 关键优化项

| 优化项 | 当前 | 建议 |
|--------|------|------|
| JWT 安全 | 硬编码密钥 | 使用 RSA 非对称密钥或 Vault |
| 在线状态查询 | Redis N+1 | 使用 Redis MGET 批量获取 |
| 离线检测 | 定时全表扫描 | 使用 Redis Key过期事件 或 sorted set |
| 日志记录 | goroutine 异步 | 写入本地 ring buffer 或 Kafka |
| 设备影子 | 仅 Redis | Redis + MySQL 双写，定期持久化 |
| 配置管理 | 环境变量 | 使用 Viper 支持多环境 YAML |

### 4.3 安全加固

1. **JWT 密钥**：改为 RSA-256 或从 Vault/KMS 获取
2. **CORS**：配置化白名单，移除 `*`
3. **Rate Limiting**：登录接口添加 IP/用户维度限流
4. **参数校验**：引入 `go-playground/validator` 替代 Gin 内置
5. **敏感日志**：过滤 password、token、phone 等字段
6. **权限校验**：在 middleware 层实现接口级权限码校验

---

## 五、重构优先级

### 🔥 P0 - 必须修复 (安全/致命缺陷)

- [ ] 硬编码 JWT 密钥 → 改为环境变量或配置中心
- [ ] CORS 开放 `*` → 配置化白名单
- [ ] MQTT CommandController nil 注入 → 正确注入 MQTT 客户端
- [ ] 登录接口限流防护 → 防暴力破解

### 🎯 P1 - 尽快修复 (功能/性能问题)

- [ ] 引入 Service 层，解耦 Controller 和业务逻辑
- [ ] 统一错误响应格式（建议 APIGateway 统一封装）
- [ ] 设备列表在线状态：N+1 → Redis MGET 批量
- [ ] 添加请求 ID（request_id）追踪
- [ ] 操作日志 goroutine → ring buffer 写入
- [ ] Redis 连接池配置

### 📋 P2 - 规划修复 (工程化)

- [ ] DTO/Request 结构分离到 dto/ 包
- [ ] 数据库事务封装（Service 层）
- [ ] 会员/宠物模型拆分到独立模块或单独服务
- [ ] 引入依赖注入框架（wire/fx/dig）
- [ ] 部门树结果缓存（Redis TTL）
- [ ] 设备影子持久化（Redis → MySQL 定期同步）

### 🚀 P3 - 长期优化 (架构演进)

- [ ] MQTT 指令 QoS 升级（0 → 1）
- [ ] 引入配置中心（Apollo/Nacos）
- [ ] 链路追踪（Jaeger/OpenTelemetry）
- [ ] 数据库读写分离
- [ ] 微服务拆分（会员服务 / 设备服务）

---

## 六、总结

当前 MDM 后端 **核心功能可用**，但在 **架构规范、安全性、性能扩展** 方面存在较大改进空间：

1. **最大风险**：JWT 硬编码 + CORS 全开放，属于安全红线
2. **最大性能瓶颈**：设备列表的在线状态过滤（N+1 查询）
3. **最大工程债务**：Service 层缺失，导致业务逻辑散落在 Controller 中

建议优先完成 P0 安全修复，再逐步推进 P1 性能优化，最后完成 P2 工程化改造。
