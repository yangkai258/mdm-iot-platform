# MEMORY.md - agenthd 知识库

## 技术栈

### 核心技能
- Go 1.21+
- Gin Web Framework
- GORM v2
- PostgreSQL
- Redis
- JWT Authentication

### 熟悉工具
- GoLand / VSCode
- Postman / Insomnia
- pgAdmin
- Redis Desktop Manager

## 项目经验

### MDM 控制中台
- 路径: backend/
- 框架: Gin + GORM
- 数据库: PostgreSQL
- 端口: 8080
- 认证: JWT

## API 实现记录

### 会员管理模块
| 路由 | 方法 | 说明 |
|------|------|------|
| /api/v1/members | GET/POST | 会员列表/创建 |
| /api/v1/members/:id | GET/PUT/DELETE | 会员详情/更新/删除 |
| /api/v1/member/cards | CRUD | 会员卡管理 |
| /api/v1/member/coupons | CRUD | 优惠券管理 |
| /api/v1/member/levels | CRUD | 会员等级 |
| /api/v1/member/stores | CRUD | 店铺管理 |
| /api/v1/member/tags | CRUD | 会员标签 |
| /api/v1/member/promotions | CRUD | 促销活动 |
| /api/v1/member/points/rules | CRUD | 积分规则 |
| /api/v1/member/orders | CRUD | 订单管理 |

## 踩坑记录

### 1. GORM 软删除
```go
// 正确用法
DB.Find(&users) // 自动过滤已删除
DB.Unscoped().Find(&users) // 包含已删除
```

### 2. JWT 过期时间
```go
// 设置合理的过期时间
expiredAt := time.Now().Add(24 * time.Hour)
```

### 3. CORS 配置
```go
// 生产环境要限制来源
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3000"},
}))
```

## 学习计划

- [ ] 深入理解 Go 协程和通道
- [ ] 掌握微服务架构
- [ ] 学习消息队列（Kafka/RabbitMQ）
- [ ] 掌握 Docker 和 K8s

---

_技术为业务服务，不要为了技术而技术。_
