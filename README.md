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
| **设备管理** | 注册、绑定、解绑、状态追踪 | ✅ 基础完成 |
| **设备影子** | Redis 实时状态缓存 | ✅ 基础完成 |
| **OTA 升级** | 固件下发、灰度发布 | ⚠️ 执行器待完善 |
| **会员管理** | 会员体系、积分、优惠券 | ✅ 基础完成 |
| **告警系统** | 设备异常告警 | ⚠️ 触发器待接入 |
| **数据分析** | 设备/会员数据报表 | 🔜 规划中 |

---

## 🏗️ 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                      前端 (Vue3 + Arco Design)              │
│                    http://localhost:3000                      │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTTP / JWT
┌─────────────────────────▼───────────────────────────────────┐
│                    后端 (Go + Gin)                           │
│                    http://localhost:8080                     │
├───────────────┬─────────────────┬───────────────────────────┤
│   PostgreSQL  │      Redis     │         EMQX              │
│   设备台账    │    设备影子     │        MQTT Broker         │
│   会员数据    │    会话缓存     │     设备心跳/指令          │
│   业务数据    │    实时状态     │     ws://localhost:8083    │
└───────────────┴─────────────────┴───────────────────────────┘
```

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
| 控制台 | http://localhost:3000 | admin / admin123 |
| EMQX Dashboard | http://localhost:18083 | admin / public |

---

## 📁 项目结构

```
mdm-project/
├── backend/                    # Go 后端服务
│   ├── controllers/           # API 控制器
│   ├── models/                # 数据模型
│   ├── middleware/            # 中间件 (JWT/CORS)
│   ├── mqtt/                  # MQTT 消息处理
│   └── main.go                # 入口
│
├── frontend/                   # Vue3 前端
│   ├── dist/                  # 静态页面
│   └── src/                   # 源码
│
├── docs/                      # 文档中心
│   ├── PRODUCT_ANALYSIS.md    # 产品分析
│   ├── BACKEND_ANALYSIS.md   # 后端分析
│   ├── FRONTEND_ANALYSIS.md  # 前端分析
│   ├── ARCHITECTURE_ANALYSIS.md # 架构分析
│   └── MEMBER_REQUIREMENTS.md # 会员模块需求
│
├── testing/                   # 测试工程师档案
├── ops/                       # 运维工程师档案
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

### 会员管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/members` | GET/POST | 会员列表/创建 |
| `/api/v1/members/:id` | GET/PUT/DELETE | 会员 CRUD |
| `/api/v1/member/cards` | GET/POST | 会员卡管理 |
| `/api/v1/member/coupons` | GET/POST | 优惠券管理 |
| `/api/v1/member/levels` | GET | 会员等级 |
| `/api/v1/member/points/rules` | GET/POST | 积分规则 |

### MQTT Topic

```
# 设备 → 云端
/device/{device_id}/up/status     # 心跳
/device/{device_id}/up/property    # 属性

# 云端 → 设备
/device/{device_id}/down/cmd      # 指令
/device/{device_id}/down/config    # 配置
```

---

## ⚠️ 当前已知问题

团队深度分析后发现以下问题，**正在修复中**：

### P0 阻断性 (🔴 优先修复)

| 问题 | 影响 | 状态 |
|------|------|------|
| JWT 密钥硬编码 | 安全风险 | 🔜 修复中 |
| MQTT 注入为 nil | 指令下发失效 | 🔜 修复中 |
| OTA 无后台 Worker | 升级指令发不出去 | 🔜 修复中 |
| 前端删除功能失效 | 数据不删除 | 🔜 修复中 |

### P1 高优先级 (🟡 规划修复)

| 问题 | 影响 |
|------|------|
| N+1 查询 | 性能瓶颈 |
| Redis 配置失效 | 生产环境连接失败 |
| 前端文件损坏 | member.html 需重构 |

---

## 📊 分析报告

团队已完成项目深度分析，报告位于 `docs/` 目录：

| 报告 | 内容 |
|------|------|
| `PRODUCT_ANALYSIS.md` | 产品需求分析、模块优先级 |
| `BACKEND_ANALYSIS.md` | 后端架构问题、安全隐患 |
| `FRONTEND_ANALYSIS.md` | UI/UX 问题、组件规范 |
| `ARCHITECTURE_ANALYSIS.md` | 系统瓶颈、演进路线 |

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
