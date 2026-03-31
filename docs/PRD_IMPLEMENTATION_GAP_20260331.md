# PRD vs 实现差距分析
**日期:** 2026-03-31
**评审人:** 架构师 (zg)

---

## 一、P0 核心功能缺口（阻断性）

| PRD 模块 | PRD 要求 | 实现状态 | 问题 |
|----------|----------|----------|------|
| **设备影子** | desired/reported/state-diff | ⚠️ partial | desired-state ✅，reported-state ❌ 404，state-diff ❌ 404 |
| **宠物行为引擎** | /behavior/generate, /behavior/actions, /behavior/rules, /behavior/sequences | ❌ 全部404 | 控制器未注册 |
| **宠物记忆** | /memory/* 对话持久化 | ❌ 未测试 | API 路径未知 |
| **OTA Worker** | 后台自动下发升级 | ⚠️ 路由存在 | Worker 是否运行未知 |
| **设备配对** | /device/pairing/* 首次配对 | ❌ 404 | 路由未注册 |
| **AI 版本管理** | /ai/models 模型版本控制 | ❌ 404 | 路由未注册 |
| **设备监控面板** | 实时传感器/地图 | ❌ 404 | 路由未注册 |

---

## 二、P1 高优先级缺口

| PRD 模块 | PRD 要求 | 实现状态 | 问题 |
|----------|----------|----------|------|
| **告警通知** | SMTP/SMS/Webhook 渠道 | ⚠️ 部分 | 路由存在，但需验证实际发送 |
| **宠物交互频率** | DND/免打扰设置 | ❌ 未找到路由 | 前端可能有 UI，后端缺失 |
| **会员卡管理** | /cards CRUD | ❌ 404 | 路由未注册 |
| **会员标签** | /members/:id/tags | ❌ 404 | 路由未注册 |
| **店铺管理** | /stores CRUD | ⚠️ 部分 | GET ✅，POST ❌ 500 |
| **权限分配 UI** | 可视化权限配置 | ⚠️ 部分 | 路由存在，需验证功能 |
| **数据权限** | 行级/列级权限 | ⚠️ 路由存在 | 功能未完整验证 |

---

## 三、P2 中优先级缺口

| PRD 模块 | 实现状态 |
|----------|----------|
| 会员订单 `/orders` | ❌ 404 |
| 会员礼包 `/gift` | ⚠️ 路由存在 (subscription_gift)，功能未知 |
| 临时会员 `/temp-members` | ❌ 404 |
| 会员服务 `/services` | ❌ 404 |
| 设备日志 `/device/logs` | ❌ 404 |
| 远程诊断 `/device/diagnosis` | ❌ 404 |
| 批量操作 `/devices/batch-*` | ⚠️ 路由存在 (import/export/batch-delete/batch-status) |
| BPMN 流程 `/flow/*` | ❌ 404 |

---

## 四、实测正常的核心 API ✅

| 类别 | API |
|------|-----|
| 认证 | `/api/v1/auth/login` |
| 设备 | `/api/v1/devices` GET/POST, `/api/v1/devices/:id/desired-state` |
| 告警 | `/api/v1/alerts`, `/api/v1/alerts/rules` |
| 会员 | `/api/v1/members` GET |
| OTA | `/api/v1/ota/packages`, `/api/v1/ota/deployments` |
| 组织 | `/api/v1/org/companies`, `/org/departments`, `/org/employees` |
| 角色权限 | `/api/v1/roles`, `/api/v1/menus`, `/api/v1/permissions` |
| AI 对话 | `/api/v1/ai/chat` (POST) |
| 数字孪生 | `/api/v1/digital-twin/*` 健康预警 ✅ |
| 健康管理 | `/api/v1/health/warnings` ✅ |
| 离线缓存 | `/api/v1/offline/cache` ✅ |
| 订阅 | `/api/v1/subscriptions/auto-renewal/*` ✅ |
| 数据脱敏 | `/api/v1/data-masking/rules` ✅ |
| 租户 | `/api/v1/tenant-approvals` ✅ |

---

## 五、根本原因

1. **路由未注册** — 控制器存在，但 `RegisterRoutes()` 未在 main.go 中调用
2. **路径不匹配** — 前端调用路径与后端注册路径不一致
3. **500 错误** — 部分 POST API 数据库模型与代码不匹配
4. **功能缺失** — PRD 定义的功能完全未开发（行为引擎/配对/会员卡等）

---

## 六、建议修复顺序

### 第一批（修复路由注册）
1. 修复 `/devices/:id/reported-state` — 路由存在但返回404
2. 注册 `/behavior/*` — 宠物行为引擎
3. 注册 `/device/pairing/*` — 设备配对
4. 注册 `/ai/models` — AI版本管理

### 第二批（修复数据库错误）
5. 修复 `POST /stores` 500 错误
6. 修复 `POST /members` 500 错误

### 第三批（新增缺失功能）
7. 会员卡管理 `/cards`
8. 会员标签 `/members/:id/tags`
9. 临时会员 `/temp-members`
10. 设备日志 `/device/logs`
