# MEMORY.md - 全局长期记忆

## 项目概述

### MDM 控制中台
- **类型**: 移动设备管理 + 会员管理系统
- **技术栈**: Go + Gin + GORM / Vue 3 + Arco Design
- **路径**: C:\Users\YKing\.openclaw\workspace\mdm-project

## 团队架构

| Agent | 角色 | 目录 |
|-------|------|------|
| zg (主) | 架构师/主管 | workspace |
| agentcp | 产品经理 | mdm-project/docs |
| agenthd | 后端开发 | mdm-project/backend |
| agentqd | 前端开发 | mdm-project/frontend |
| agentcs | 测试工程师 | mdm-project/testing |
| agentyw | 运维工程师 | mdm-project/ops |

## 技能库

已安装以下技能提升团队效率：
- agent-task-tracker
- agent-orchestrate
- team-status-tracker
- skill-vetter
- self-monitor
- task-development-workflow
- product-manager
- requirements-analysis
- memory-setup
- proactive-agent

## 完成的功能

### 设备管理
- 设备注册/绑定/解绑
- 设备状态管理
- OTA 固件管理
- 指令下发

### 会员管理
- 会员 CRUD (9个子模块)
- 会员卡管理
- 优惠券管理
- 会员等级
- 积分规则
- 促销活动
- 店铺管理
- 订单管理

## 访问信息
- 前端: http://localhost:3000
- 后端: http://localhost:8080
- 账号: admin / admin123

## 重要教训

1. **架构师不自己干活** - Agent 超时应该分析原因、重新派任务，不是自己写代码
2. **任务要小而明确** - 避免大任务导致超时
3. **检查点机制** - 监控 Agent 状态，及时干预
4. **会话快照** - 已设置每小时自动快照，防止系统重启后丢失工作进度

## ⚠️ 已知问题

### 会话丢失问题 (2026-03-21)
- **现象**: 用户昨晚工作到2点，今早发现我没有那段时间的任何记忆
- **原因**: OpenClaw重启或会话被清理时会丢失历史
- **解决**: 创建了每小时cron任务自动保存会话快照到 `memory/SESSION_SNAPSHOT.md`

## 前端迁移策略 (2026-03-21 更新)

**目标目录**：`mdm-frontend-new/arco-design-pro-vite/`（Arco Design Pro 模板）
**旧目录**：`frontend/`（不再维护）

**迁移进度**：昨晚2点创建了全新ArcoPro工程，已有多模块view文件

**UI规范**：标准三段式布局
1. 面包屑 → 页面标题
2. 搜索筛选区（浅灰卡片，三列栅格）
3. 操作栏（新建靠左，下载靠右）
4. 数据表格

**参考设计**：`screenshot_reference.png`

---

## 🚨 核心开发流程（强制执行）

**需求定义 → 产品+架构师评审 → 开发**

```
用户需求
    ↓
agentcp（产品经理）输出 PRD/接口契约
    ↓
zg（架构师）+ agentcp 联合评审
    ↓
确认需求完整、接口清晰、无歧义
    ↓
才能派给 agenthd/agentqd 开始开发
```

**关键约束：所有需求必须 agentcp + zg 联合定义完整后才交给开发Agent。**
- 禁止需求不完整就交给开发
- 禁止开发Agent自己补充需求细节
- 如果开发中发现需求模糊，必须打回 agentcp 重新定义

---

## 2026-03-20 重大更新

### 团队深度分析
- 4个Agent并行分析MDM项目
- 产出4份分析报告 (PRODUCT/BACKEND/FRONTEND/ARCHITECTURE)

### 识别的核心问题

**P0阻断性 (9个):** ✅ 已全部修复
- JWT密钥硬编码
- CORS全开放
- MQTT注入为nil (指令下发失效)
- OTA无后台Worker
- CheckAlerts从未调用
- 回车登录失效
- deleteRecord不删数据
- member.html文件损坏
- Dashboard统计查询错误字段

**P1高优先级 (5个):** ✅ 已全部修复
- N+1查询问题
- Redis URL解析失效
- 主色拼写错误(165qff)
- API硬编码localhost
- 无登录限流

### GitHub
- https://github.com/yangkai258/mdm-iot-platform

## 2026-03-21 Sprint 3 状态

### Sprint 3 完成 ✅
- **agenthd**: 告警SMTP/Webhook通知服务 + 合规策略API
- **agentqd**: 前端源码commit (18文件) + 新模块commit (22文件)
- **新backend commit**: `e60944b`
- **新frontend commit**: `f349003`, `21c213b`

### 今日完成汇总
| 类型 | 数量 | 状态 |
|------|------|------|
| P0 问题 | 9个 | ✅ 全部修复+推送 |
| P1 问题 | 5个 | ✅ 全部修复+推送 |
| 新增前端模块 | 22文件 | ✅ 已推送 |
| Sprint 3 前端 | 18文件 | ✅ 已推送 |
| Sprint 3 后端 | 告警通知+合规API | ✅ 已推送 |

> ✅ Subagents当前无活跃会话，P0+P1已全部修复推送

### 进行中的开发
- **多租户系统**: 新增 tenant_controller, company_controller, department_controller, employee_controller, position_template_controller 等
- **权限系统增强**: permission_controller, permission_group_controller, role_controller, menu_controller
- **新中间件**: tenant.go, permission.go, quota_check.go
- **新模型**: permission_models.go, tenant.go
- **新文档**: MULTI_TENANT_PRD.md + 8个模块PRD文档

### Sprint 6 UI规范
- 风格: ArcoDesign Pro
- 面包屑: 左上角
- 搜索框: 面包屑下方，左对齐
- 按钮组: 搜索框下方，靠左排列

### 最近提交
- `e60944b` feat(backend): add alert SMTP/Webhook notification service and compliance policies API routes
- `f349003` feat: commit Sprint 3 frontend source files
- `21c213b` feat: 添加新模块 - org/tenants/permissions/pet/owner/knowledge/miniclaw 等
- `d1fa70d` fix: 实现登录限流保护
- `080a059` fix: 移除 API base URL 硬编码 localhost，改为相对路径 /api

### 新增前端模块 (今日)
22个新文件: org/(4), tenants/(4), permissions/(5), pet/(2), owner/(1), knowledge/(1), miniclaw/(1), components/(2), assets/(1), views/(1)

---

_持续更新，记录项目的成长轨迹。_
