# MDM 控制中台 — 周开发计划
**日期：** 2026-03-23 至 2026-03-29
**执行：** agentcp + agenthd + agentqd + agentcs + agentyw

---

## 📅 Day 1 · 2026-03-23（周一）

### Sprint 21：具身智能核心
| Agent | 任务 | 交付物 |
|-------|------|--------|
| agenthd | 环境感知/空间认知 API | embodied_controller.go + 模型 |
| agentqd | 感知面板/地图管理/导航控制页面 | 6个 Vue 组件 |
| agentcp | 仿真测试 PRD 补充 | MODULE_SIMULATION.md 完善 |

---

## 📅 Day 2 · 2026-03-24（周二）

### Sprint 22：具身智能前端 + 仿真测试
| Agent | 任务 | 交付物 |
|-------|------|--------|
| agenthd | 动作模仿/决策引擎/安全边界 API | embodied API 完善 |
| agentqd | 动作库/安全禁区/决策日志页面 | 3个 Vue 页面 |
| agentcs | 具身智能 API 单元测试 | 仿真测试用例 |
| agentcp | 第三方集成 PRD | MODULE_PLATFORM_ECOSYSTEM.md 补充 |

---

## 📅 Day 3 · 2026-03-25（周三）

### Sprint 23：仿真测试框架
| Agent | 任务 | 交付物 |
|-------|------|--------|
| agenthd | 虚拟宠物仿真/自动化测试 API | simulation_controller.go |
| agentqd | 虚拟宠物仿真前端 | VirtualPetSimulationView |
| agentcp | 开放平台 PRD | MODULE_PLATFORM_ECOSYSTEM.md 完善 |
| agentyw | 仿真环境 CI/CD 配置 | Dockerfile + pipeline |

---

## 📅 Day 4 · 2026-03-26（周四）

### Sprint 23-24：仿真测试完善
| Agent | 任务 | 交付物 |
|-------|------|--------|
| agenthd | 场景管理/回放系统/压力测试 API | simulation API 完善 |
| agentqd | 场景管理/压力测试页面 | StressTestView + SceneManageView |
| agentcs | 仿真测试集成测试 | 端到端测试 |
| agentyw | 仿真资源配额管理 | quota_service.go |

---

## 📅 Day 5 · 2026-03-27（周五）

### Sprint 25：开放平台基础
| Agent | 任务 | 交付物 |
|-------|------|--------|
| agenthd | 开发者 API 完善/Webhook 市场 API | developer + webhook API |
| agentqd | 开发者控制台/应用商店前端 | DeveloperConsoleView |
| agentcp | 高级功能 PRD 补充 | Sprint 29-30 PRD |
| agentcs | API 契约测试 | API 测试用例 |

---

## 📅 Day 6 · 2026-03-28（周六）

### Sprint 25-26：内容市场
| Agent | 任务 | 交付物 |
|-------|------|--------|
| agenthd | 表情包/动作/声音资源 API | market_controller.go |
| agentqd | 表情包市场/动作库/声音定制页面 | EmoticonMarket + ActionResource |
| agentcp | 第三方集成接口定义 | Integration API Spec |
| agentyw | 开放平台安全配置 | OAuth2 + API Key 管理 |

---

## 📅 Day 7 · 2026-03-29（周日）

### Sprint 27：第三方集成
| Agent | 任务 | 交付物 |
|-------|------|--------|
| agenthd | 智能家居/宠物医疗/地图服务 API | integration_controller.go |
| agentqd | 第三方集成配置页面 | IntegrationView |
| agentcs | 第三方 API 集成测试 | Integration tests |
| agentyw | 第三方服务监控配置 | Monitoring setup |

---

## 📊 周总结

| 完成 Sprint | 后端 API | 前端页面 |
|-------------|---------|---------|
| Sprint 21 | 27个 | 6个 |
| Sprint 22 | - | 3个 |
| Sprint 23 | 12个 | 4个 |
| Sprint 24 | 8个 | 3个 |
| Sprint 25 | 10个 | 3个 |
| Sprint 26 | 8个 | 4个 |
| Sprint 27 | 12个 | 4个 |

**总计：** 77个 API · 27个页面 · 1个 PRD 完善 · 1个 CI/CD 配置

---

## 🔑 明天唤醒指令

> **"继续推进 MD"**

发送这句话给我，我就会自动安排 5 个 agent 继续 Sprint 21-27 的开发工作。
