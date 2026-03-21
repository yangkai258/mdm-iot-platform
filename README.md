# MDM 物联网设备管理中台

> 🔧 **AICD + 设备管理双引擎驱动** — 面向 M5Stack 物联网设备的智能管控平台

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3-4FC08D?style=flat-square&logo=vue.js)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

---

## 🎯 项目简介

MDM (Mobile Device Management) 中台是 AICD（AI Companion Dashboard）架构的核心组件，负责统一管理分布在全球的 M5Stack 物联网设备。

### 核心能力

| 模块 | 功能 | 状态 |
|------|------|------|
| **设备管理** | 注册、绑定、解绑、状态追踪、OTA升级 | ✅ 基础完成 |
| **设备影子** | Redis 实时状态缓存 (TTL 90秒) | ✅ 基础完成 |
| **OTA 升级** | 固件下发、灰度发布、后台 Worker | ✅ 完成 |
| **告警系统** | 设备异常告警、CheckAlerts 监控 | ✅ 完成 |
| **会员管理** | 会员体系、积分、优惠券、等级 | ✅ 基础完成 |
| **数据分析** | 设备/会员数据报表 | 🔜 规划中 |
| **多租户系统** | 租户管理、部门岗位、权限组 | ✅ 完成 |
| **数据权限** | 行级数据隔离 (DataScope 中间件) | ✅ 完成 |

---

## 🏗️ 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                    前端 (Vue3 + Arco Design Pro)            │
│                    http://localhost:5174 (新)                │
│                    http://localhost:3000 (旧)                │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTTP / JWT
┌─────────────────────────▼───────────────────────────────────┐
│                    后端 (Go + Gin + GORM)                   │
│                    http://localhost:16666                    │
├───────────────┬─────────────────┬───────────────────────────┤
│   PostgreSQL  │      Redis     │         EMQX              │
│   设备台账    │    设备影子     │        MQTT Broker         │
│   会员数据    │    会话缓存     │     设备心跳/指令          │
│   业务数据    │    实时状态     │     ws://localhost:8083    │
└───────────────┴─────────────────┴───────────────────────────┘
```

---

## 🚀 Sprint 开发进度

### Sprint 8 ✅ 完成 (最新)
- **DataScope 数据权限中间件** - 实现行级数据隔离
  - 支持 5 种权限范围：全部/本部门/本部门及下属/仅本人/自定义
  - 自动过滤 GORM 查询结果
  - 租户隔离支持

### Sprint 7 ✅ 完成
- **OTA 后台 Worker** - 固件升级任务队列
- **CheckAlerts 告警监控** - 设备异常自动检测
- **设备影子增强** - Redis TTL 90秒心跳
- **NRD 前端迁移** - Vue3 + Arco Design Pro

### Sprint 6 ✅ 完成
- **多租户系统** - 租户/公司/部门/岗位/员工
- **权限系统增强** - 权限组/角色/菜单/API权限
- **Position Templates** - 岗位模板克隆/启用/禁用

### Sprint 1-5 ✅ 完成
- 设备注册/绑定/心跳
- 会员 CRUD 9个子模块
- 基础 API 接口

---

## 👥 Agent 团队

本项目采用 **多 Agent 协作开发模式**，由 OpenClaw 驱动的 AI Agent 团队共同维护：

| Agent | 角色 | 核心职责 |
|-------|------|----------|
| **zg** | 架构师 | 技术选型、任务调度、质量管控 |
| **agentcp** | 产品经理 | 需求分析、API 契约设计、产品规划 |
| **agenthd** | 后端开发 | Go 服务端开发、数据库设计 |
| **agentqd** | 前端开发 | Vue3 组件、UI/UX 实现 |
| **agentcs** | 测试工程师 | 压测、功能测试、质量把关 |
| **agentyw** | 运维工程师 | Docker 部署、监控告警 |

### 团队协作规范

```
用户需求 → agentcp 分析 → zg 审核 → agenthd/agentqd 并行开发 → agentcs 测试 → agentyw 部署
```

---

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- EMQX 5.0

### 1. 克隆项目

```bash
git clone https://github.com/yangkai258/mdm-iot-platform.git
cd mdm-iot-platform
```

### 2. 一键启动

```bash
# 启动所有依赖服务
docker-compose up -d

# 启动后端
cd backend
go run main.go

# 启动前端 (新终端)
cd frontend
npm install
npm run dev
```

### 3. 访问

| 服务 | 地址 | 账号 |
|------|------|------|
| 控制台 | http://localhost:5174 | admin / admin123 |
| EMQX Dashboard | http://localhost:18083 | admin / public |

---

## 📁 项目结构

```
mdm-project/
├── backend/                    # Go 后端服务
│   ├── controllers/           # API 控制器
│   ├── models/                # 数据模型
│   ├── middleware/           # 中间件 (JWT/CORS/DataScope)
│   ├── plugins/              # GORM 插件 (DataScope)
│   ├── mqtt/                 # MQTT 消息处理
│   ├── ota/                  # OTA 升级 Worker
│   ├── cmd/genhash/          # 密码生成工具
│   └── main.go               # 入口
│
├── frontend/                  # Vue3 前端 (旧)
│   └── dist/                  # 静态页面
│
├── mdm-frontend-new/           # Vue3 前端 (新 - ArcoPro)
│   └── arco-design-pro-vite/
│
├── docs/                      # 文档中心
│   ├── SPRINT_1-7_REVIEW.md  # Sprint 1-7 评审报告
│   ├── SPRINT_8.md           # Sprint 8 开发文档
│   └── *.md                  # 其他PRD文档
│
├── testing/                   # 测试工程师档案
├── ops/                      # 运维工程师档案
└── docker-compose.yml        # 一键部署
```

---

## 🔌 API 接口

### 设备管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/devices/register` | POST | 设备注册 |
| `/api/v1/devices` | GET | 设备列表 |
| `/api/v1/devices/:id` | GET/PUT/DELETE | 设备 CRUD |
| `/api/v1/devices/:id/status` | PUT | 更新状态 |
| `/api/v1/devices/:id/ota` | POST | 下发OTA指令 |

### 会员管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/members` | GET/POST | 会员列表/创建 |
| `/api/v1/members/:id` | GET/PUT/DELETE | 会员 CRUD |
| `/api/v1/member/cards` | GET/POST | 会员卡管理 |
| `/api/v1/member/coupons` | GET/POST | 优惠券管理 |
| `/api/v1/member/levels` | GET | 会员等级 |
| `/api/v1/member/points/rules` | GET/POST | 积分规则 |

### 组织架构

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/tenants` | GET/POST | 租户管理 |
| `/api/v1/companies` | GET/POST | 公司管理 |
| `/api/v1/departments` | GET/POST | 部门管理 |
| `/api/v1/employees` | GET/POST | 员工管理 |
| `/api/v1/positions` | GET/POST | 岗位管理 |

### 权限管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/roles` | GET/POST | 角色管理 |
| `/api/v1/permissions` | GET | 权限列表 |
| `/api/v1/permission-groups` | GET/POST | 权限组 |
| `/api/v1/data-scopes` | GET/POST | 数据权限配置 |

### MQTT Topic

```
# 设备 → 云端
/device/{device_id}/up/status     # 心跳
/device/{device_id}/up/property    # 属性

# 云端 → 设备
/device/{device_id}/down/cmd      # 指令
/device/{device_id}/down/config    # 配置
/device/{device_id}/down/ota      # OTA升级
```

---

## 🔐 数据权限 (DataScope)

DataScope 中间件实现行级数据隔离，支持 5 种权限范围：

| 范围类型 | 说明 | 适用场景 |
|----------|------|----------|
| `DataScopeAll` | 全部数据 | 超管 |
| `DataScopeOrg` | 本部门数据 | 部门主管 |
| `DataScopeOrgAndChildren` | 本部门及下属 | 高层管理 |
| `DataScopeSelf` | 仅本人数据 | 普通员工 |
| `DataScopeCustom` | 自定义数据范围 | 特殊权限 |

使用方式：
```go
// 在 GORM 查询中自动注入数据权限过滤
db = db.Session(&gorm.Session{}).Set("data_scope", &DataScopeInfo{
    UserID:    userID,
    RoleID:    roleID,
    ScopeType: scopeType,
})
```

---

## 📊 Sprint 评审

| Sprint | 状态 | 核心产出 |
|--------|------|----------|
| Sprint 1-5 | ✅ | 设备管理、会员管理基础 |
| Sprint 6 | ✅ | 多租户系统、权限系统增强 |
| Sprint 7 | ✅ | OTA Worker、CheckAlerts、NRD前端 |
| Sprint 8 | ✅ | DataScope 数据权限中间件 |

详细评审报告见 `docs/SPRINT_*-REVIEW.md`

---

## 🤝 贡献指南

欢迎提交 Issue 和 PR！

1. Fork 本仓库
2. 创建特性分支 `git checkout -b feature/AmazingFeature`
3. 提交更改 `git commit -m 'Add some AmazingFeature'`
4. 推送分支 `git push origin feature/AmazingFeature`
5. 创建 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - 高性能 Go Web 框架
- [Vue3](https://github.com/vuejs/core) - 渐进式 JavaScript 框架
- [Arco Design](https://arco.design/) - 企业级设计系统
- [EMQX](https://www.emqx.io/) - 全球最具扩展性的 MQTT broker

---

> **💡 提示**: 本项目由 AI Agent 团队协作开发维护，持续迭代中。
