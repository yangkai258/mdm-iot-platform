# MDM 控制中台 — 产品路线图

**版本：** V2.2
**编制角色：** 产品经理 (agentcp)
**编制日期：** 2026-03-22
**文档状态：** 重大更新——76个功能点逐项列出，UI规范引用 PRD 文档，验收标准完善

---

## 一、项目概述

### 1.1 产品定位

MDM (Mobile Device Management) 控制中台是"端云协同"AI 电子宠物系统的核心管理平台，负责统一管理分布在全球的 M5Stack 硬件终端，维护设备全生命周期（注册、绑定、监控、OTA 升级、指令控制），并通过会员管理模块支持用户运营能力。

### 1.2 技术栈

| 层级 | 技术选型 |
|------|----------|
| 设备端 | M5Stack 硬件 + MimiClaw 固件 |
| 通信层 | MQTT (paho.mqtt.golang) + EMQX 5.0 |
| 后端 | Go + Gin + GORM + PostgreSQL + Redis |
| 前端 | Vue 3 + TypeScript + Vite + Arco Design Vue |
| 部署 | Docker + docker-compose + Nginx |

### 1.3 六大功能领域

| 领域 | 描述 | 涵盖模块 |
|------|------|----------|
| **D1 核心平台与AI能力** | AI训练、A/B测试、数字孪生、情感计算、具身智能、仿真测试 | AI工程、数字孪生、情感计算、具身智能、仿真 |
| **D2 安全与合规** | LDAP/AD、证书、远程锁定、音视频加密、数据合规 | 企业安全、合规策略 |
| **D3 设备管理与运维** | 地理围栏、跨境管控、DaaS、RTOS、BLE Mesh、端侧推理 | 设备管理、固件管理、运维工具 |
| **D4 用户与宠物管理** | 宠物登记、运动追踪、睡眠分析、寻回网络、家庭相册 | 宠物管理、用户管理、健康管理 |
| **D5 生态、商业与集成** | 订阅、用量计费、Webhook、开发者API、市场 | 订阅管理、开放平台、支付计费 |
| **D6 技术架构高阶** | RTOS优化、BLE Mesh、端侧推理、模型分片 | 技术架构 |

---

## 二、Phase 总览

| Phase | 周期 | 主题 | 核心目标 | Sprint |
|-------|------|------|----------|--------|
| **Phase 1** | Sprint 1-8 | 核心平台与AI | 完善核心平台，奠定AI系统工程基础 | S1-S8 |
| **Phase 2** | Sprint 9-16 | 企业级与安全合规 | 企业安全、全球化、商业化基础设施 | S9-S16 |
| **Phase 3** | Sprint 17-24 | 具身智能平台 | 数字孪生、情感计算、具身智能、仿真测试 | S17-S24 |
| **Phase 4** | Sprint 25-32 | 生态扩展 | 开放平台、市场、内容生态、第三方集成 | S25-S32 |

---

## 三、Phase 1：核心平台与AI（Sprint 1-8）

> **目标：** 完善核心平台能力，建立AI系统工程基础，实现设备与宠物基础管理

### Sprint 1-2：核心平台补齐

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| OTA Worker 实现 | P0 | MODULE_OTA_UPDATES.md |
| OTA 数据模型对齐 | P0 | MODULE_OTA_UPDATES.md |
| 设备影子（期望状态）完善 | P1 | MODULE_DEVICE_SHADOW.md |
| 设备模式管理（DND/安静/活跃） | P1 | MODULE_DEVICE_SHADOW.md |
| 认证授权 JWT 完善 | P1 | MODULE_PERMISSIONS.md |
| 设备指令下发完善 | P1 | MODULE_DEVICE_MANAGEMENT.md |
| CheckAlerts 集成 MQTT | P1 | MODULE_ALERT_SYSTEM.md |
| 告警通知渠道（邮件/SMS/Webhook） | P1 | MODULE_NOTIFICATION.md |
| 系统管理基础功能（用户/角色/权限/RBAC） | P1 | MODULE_SYSTEM_MANAGEMENT.md |
| 设备注册与配对流程 | P0 | MODULE_DEVICE_PAIRING.md |

**验收标准：**

- [ ] OTA Worker 后台部署任务自动执行，Worker 自动选中目标设备并下发 MQTT OTA 指令
- [ ] OTA 数据模型与 PRD 一致，DB migration + models 对齐
- [ ] `desired_config` MQTT 下发机制完整，设备影子状态同步正常
- [ ] 设备模式（DND/安静/活跃）CRUD 完整，MQTT 联动正常
- [ ] JWT 刷新 token + 权限矩阵落地，API 权限验证通过
- [ ] MQTT handler 完整实现，设备指令下发响应时间 < 2s
- [ ] CheckAlerts 设备心跳触发告警检查，告警触发准确率 > 95%
- [ ] 告警通知服务（邮件/SMS/Webhook）多渠道推送可用
- [ ] 系统管理：用户CRUD + 角色定义 + 权限矩阵配置 + RBAC 策略管理
- [ ] 设备配对：设备注册流程 + 配对码生成/验证 + 设备绑定

---

### Sprint 3-4：宠物基础管理

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 宠物配置完善 | P1 | MODULE_PET_BEHAVIOR_ENGINE.md |
| 宠物性格设定 | P1 | MODULE_OPENCLAW_VERSION.md |
| 免打扰规则设备端联动 | P1 | MODULE_DEVICE_SHADOW.md |
| 交互频率配置 | P1 | MODULE_OPENCLAW_VERSION.md |
| 宠物基础 CRUD | P1 | MODULE_PET_BEHAVIOR_ENGINE.md |
| 多用户交互识别 | P1 | MODULE_OWNER_PROFILE.md |
| 设备升级进度追踪 | P1 | MODULE_OTA_UPDATES.md |
| 告警确认/解决流程 | P1 | MODULE_ALERT_SYSTEM.md |
| 宠物记忆系统（交互历史/学习记忆） | P1 | MODULE_PET_MEMORY.md |
| 组织架构管理（部门/团队/成员） | P1 | MODULE_ORGANIZATION.md |

**验收标准：**

- [ ] 宠物配置 CRUD 完整，性格/交互频率正确下发到设备
- [ ] 性格类型枚举完整，AI 联动定义正确
- [ ] DND MQTT 下发 + 设备端联动正常
- [ ] 频率控制逻辑完整，前端配置 UI 可用
- [ ] 宠物档案管理页面完整，支持增删改查
- [ ] 家庭成员识别与差异化响应正确区分多用户
- [ ] OTAProgress 表 + 设备端回调，进度追踪实时准确
- [ ] 告警列表 + 状态流转（待确认→确认中→已解决）完整
- [ ] 宠物记忆：交互历史存储 + 长期记忆学习 + 记忆检索，记忆持久化正常
- [ ] 组织架构：部门/团队CRUD + 成员管理 + 层级关系，组织架构管理完整

---

### Sprint 5-6：AI系统工程基础

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| AI模型训练流水线框架 | P1 | MODULE_AI_ENGINEERING.md |
| A/B测试框架 | P1 | MODULE_AI_ENGINEERING.md |
| 模型监控仪表盘 | P1 | MODULE_AI_ENGINEERING.md |
| AI决策日志 | P1 | MODULE_AI_ENGINEERING.md |
| 模型版本管理 | P1 | MODULE_AI_ENGINEERING.md |
| 模型热回滚 | P1 | MODULE_AI_ENGINEERING.md |
| AI沙箱测试环境 | P1 | MODULE_AI_ENGINEERING.md |
| 边缘AI vs 云端AI路由 | P2 | MODULE_AI_ENGINEERING.md |

**验收标准：**

- [ ] 训练任务管理 + 状态追踪，训练任务可创建/暂停/取消
- [ ] A/B测试框架：实验配置支持多分组，分组流量均匀分布，效果数据自动统计
- [ ] 模型监控仪表盘：响应延迟/成功率实时展示，异常告警自动触发
- [ ] AI决策日志：决策记录表 + 查询接口，决策原因可追溯
- [ ] 模型版本管理：版本列表 + 上线状态，版本切换正常
- [ ] 模型热回滚：一键切换到历史版本，切换时间 < 30秒
- [ ] AI沙箱测试环境：隔离测试环境 + 测试报告，测试通过后才可上线
- [ ] 边缘AI vs 云端AI路由：智能路由策略配置，路由切换正常

---

### Sprint 7-8：会员管理与运营基础

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 会员信息管理 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 会员标签管理 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 优惠券管理 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 积分管理 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 促销管理 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 订单管理 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 多宠物管理 | P1 | MODULE_PET_BEHAVIOR_ENGINE.md |
| 运动追踪基础 | P1 | MODULE_DIGITAL_TWIN.md |

**验收标准：**

- [ ] 会员 CRUD + 会员等级完整，会员信息准确
- [ ] 标签 CRUD + 批量打标功能，标签应用正确
- [ ] 优惠券 CRUD + 发放/核销，优惠券使用核销正常
- [ ] 积分规则 + 获取/抵扣，积分计算准确
- [ ] 促销活动配置，促销规则正确应用
- [ ] 订单状态流转完整（待支付→已支付→已完成→已取消）
- [ ] 一账号多宠物 + 快速切换，宠物切换延迟 < 1s
- [ ] 运动数据采集 + 统计，运动数据展示准确

---

## 四、Phase 2：企业级与安全合规（Sprint 9-16）

> **目标：** 建立企业级安全能力，支持全球化部署，完善商业化基础设施

### Sprint 9-10：企业安全与设备管理

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| LDAP/AD 集成 | P1 | MODULE_PERMISSIONS.md |
| 证书管理 | P1 | MODULE_DEVICE_MANAGEMENT.md |
| 远程锁定/擦除 | P1 | MODULE_DEVICE_MANAGEMENT.md |
| 合规策略强制 | P1 | MODULE_POLICY_MANAGEMENT.md |
| 数据脱敏 | P2 | MODULE_DATA_ANALYSIS.md |
| 数据最小化 | P2 | MODULE_DATA_ANALYSIS.md |
| AI伦理/公平性测试 | P2 | MODULE_AI_ENGINEERING.md |
| 设备监控面板（实时状态/指标/告警） | P1 | MODULE_DEVICE_MONITORING.md |
| 设备配对注册完善 | P0 | MODULE_DEVICE_PAIRING.md |
| 知识库/天气/问答 | P1 | MODULE_KNOWLEDGE_BASE.md |

**验收标准：**

- [ ] 企业用户SSO登录成功，LDAP/AD部门/用户自动同步
- [ ] 设备证书自动颁发/续期/吊销，证书管理完整
- [ ] 锁定指令5秒内到达设备，擦除后数据不可恢复
- [ ] 策略规则引擎 + 强制执行，不合规设备隔离正常
- [ ] 敏感字段自动脱敏，脱敏规则验证通过
- [ ] 采集最小化配置，数据采集符合最小必要原则
- [ ] AI公平性检测工具，公平性报告生成正常
- [ ] 设备监控面板：实时状态展示 + 指标监控 + 告警列表，监控数据延迟 < 5s
- [ ] 设备配对：配对流程优化 + 批量配对支持，配对成功率 > 99%
- [ ] 知识库：问答检索 + 天气查询 + 知识图谱，知识库召回率 > 90%

---

### Sprint 11-12：全球化与数据合规

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 地理围栏 | P1 | MODULE_DEVICE_MANAGEMENT.md |
| 跨境设备统一管控 | P1 | MODULE_DEVICE_MANAGEMENT.md |
| 设备即服务（DaaS） | P2 | MODULE_DEVICE_MANAGEMENT.md |
| 数据驻留合规 | P1 | MODULE_DATA_ANALYSIS.md |
| 区域AI节点 | P1 | MODULE_AI_ENGINEERING.md |
| 多时区支持 | P1 | MODULE_AI_ENGINEERING.md |
| RTOS优化 | P2 | MODULE_MINICLAW_FIRMWARE.md |
| BLE Mesh组网 | P3 | MODULE_MINICLAW_FIRMWARE.md |

**验收标准：**

- [ ] 围栏配置 + 进入/离开告警，围栏边界触发准确率 > 95%
- [ ] 多区域设备统一管理界面，跨境设备管控正常
- [ ] 设备租赁管理 + 使用计费，DaaS 流程完整
- [ ] 多区域数据库 + 数据访问控制，数据驻留合规
- [ ] 边缘AI部署 + 智能路由，区域AI节点响应正常
- [ ] 宠物作息按时区 + UTC存储，多时区切换正常
- [ ] 设备端性能优化，RTOS响应延迟降低 > 20%
- [ ] 多设备Mesh组网协议，Mesh组网连接正常

---

### Sprint 13-14：AI系统工程完善

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| AI行为监控/异常告警 | P1 | MODULE_AI_ENGINEERING.md |
| AI决策可解释性 | P1 | MODULE_AI_ENGINEERING.md |
| A/B测试完善 | P1 | MODULE_AI_ENGINEERING.md |
| 模型分片加载 | P2 | MODULE_AI_ENGINEERING.md |
| 端侧推理 | P2 | MODULE_AI_ENGINEERING.md |
| AI质量报告 | P1 | MODULE_AI_ENGINEERING.md |

**验收标准：**

- [ ] 异常行为检测 + 告警触发，AI行为监控准确率 > 90%
- [ ] 决策解释界面 + 原因追溯，用户可查看为什么
- [ ] 灰度发布 + 统计分析，灰度流量分配正确
- [ ] 按需加载 + 增量更新，模型分片加载正常
- [ ] 设备端小模型部署，端侧推理延迟 < 100ms
- [ ] 自动生成质量报告，AI质量报告完整

---

### Sprint 15-16：商业化基础设施

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 订阅管理 | P1 | MODULE_SUBSCRIPTION.md |
| 用量计费 | P2 | MODULE_SUBSCRIPTION.md |
| API配额计费 | P2 | MODULE_SUBSCRIPTION.md |
| 优惠券/促销完善 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 发票/账单 | P2 | MODULE_SUBSCRIPTION.md |
| Webhook事件系统 | P1 | MODULE_SUBSCRIPTION.md |
| 开发者API | P1 | MODULE_PLATFORM_ECOSYSTEM.md |
| 会员等级权益配置 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 应用分发/企业应用商店 | P2 | MODULE_APP_MANAGEMENT.md |
| 文件库/内容分发管理 | P2 | MODULE_CONTENT_MANAGEMENT.md |

**验收标准：**

- [ ] 多级订阅 + 自动续费 + 变更，订阅变更即时生效
- [ ] 按对话量/存储量计费，用量统计准确
- [ ] 开发者API用量统计与计费，配额扣费正常
- [ ] 促销规则配置 + 发放策略，促销应用正确
- [ ] 自动开票 + 月度账单，发票生成正常
- [ ] 事件订阅 + 事件类型配置，Webhook推送正常
- [ ] 开放平台 API 文档 + 沙箱，API文档完整可用
- [ ] 等级权益差异化配置，权益应用正确
- [ ] 应用商店：应用上架/审核/下架 + 批量分发 + 安装追踪，企业应用商店可用
- [ ] 内容管理：文件上传/存储/分发 + 媒体库管理，内容分发正常

---

## 五、Phase 3：具身智能平台（Sprint 17-24）

> **目标：** 构建数字孪生、情感计算、具身智能核心能力，建立仿真测试体系

### Sprint 17-18：情感计算

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 用户语音情绪识别 | P1 | MODULE_AFFECTIVE_COMPUTING.md |
| 用户文字情绪识别 | P1 | MODULE_AFFECTIVE_COMPUTING.md |
| 宠物表情情绪识别 | P1 | MODULE_AFFECTIVE_COMPUTING.md |
| 情绪响应策略 | P1 | MODULE_AFFECTIVE_COMPUTING.md |
| 情绪低落安慰 | P1 | MODULE_AFFECTIVE_COMPUTING.md |
| 情绪日志 | P2 | MODULE_AFFECTIVE_COMPUTING.md |
| 家庭情绪地图 | P2 | MODULE_AFFECTIVE_COMPUTING.md |
| 情绪趋势分析 | P3 | MODULE_AFFECTIVE_COMPUTING.md |

**验收标准：**

- [ ] 语音情绪分析模型，语音情绪识别准确率 > 85%
- [ ] NLP情绪分析，文字情绪识别准确率 > 85%
- [ ] 宠物表情/行为情绪判断，表情识别准确率 > 80%
- [ ] 基于情绪的响应动作正确触发，响应延迟 < 500ms
- [ ] 宠物主动安慰动作触发，用户满意度评分 > 4.0
- [ ] 情绪记录 + 历史查询，情绪日志完整可查
- [ ] 多宠物家庭情绪感知，家庭情绪地图展示正常
- [ ] 周/月情绪曲线，情绪趋势分析正常

---

### Sprint 19-20：数字孪生

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 实时生命体征数字孪生 | P1 | MODULE_DIGITAL_TWIN.md |
| 行为预测 | P2 | MODULE_DIGITAL_TWIN.md |
| 健康预警 | P2 | MODULE_DIGITAL_TWIN.md |
| 历史回放 | P2 | MODULE_DIGITAL_TWIN.md |
| 精彩瞬间AI筛选 | P2 | MODULE_DIGITAL_TWIN.md |
| 跨设备状态同步 | P1 | MODULE_DEVICE_SHADOW.md |
| 离线支持 | P2 | MODULE_DEVICE_SHADOW.md |

**验收标准：**

- [ ] 心跳/呼吸实时曲线，生命体征延迟 < 2秒
- [ ] 宠物下一步动作预测，预测准确率 > 75%
- [ ] 行为异常提前预警，预警准确率 > 80%
- [ ] 宠物完整生命记录，历史数据完整可查
- [ ] 高光时刻自动剪辑，精彩瞬间筛选正常
- [ ] 云端状态同步 + 切换设备续连，切换延迟 < 3s
- [ ] 本地缓存 + 断网续传，离线数据同步正常

---

### Sprint 21-22：具身智能核心

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 环境感知 | P1 | MODULE_EMBODIED_AI.md |
| 空间认知 | P1 | MODULE_EMBODIED_AI.md |
| 自主探索 | P1 | MODULE_EMBODIED_AI.md |
| 动作模仿 | P1 | MODULE_EMBODIED_AI.md |
| 具身AI决策引擎 | P1 | MODULE_EMBODIED_AI.md |
| 具身AI安全边界 | P1 | MODULE_EMBODIED_AI.md |
| 多宠物协作 | P2 | MODULE_EMBODIED_AI.md |

**验收标准：**

- [ ] 摄像头视觉感知 + 障碍物检测，环境感知准确率 > 90%
- [ ] 室内地图构建 + 定位，空间认知准确率 > 85%
- [ ] 自主导航 + 路径规划，自主导航成功率 > 85%
- [ ] 动作学习 + 动作库，动作模仿准确率 > 80%
- [ ] 感知-认知-决策-执行闭环，决策延迟 < 1s
- [ ] 物理交互安全限制，安全边界触发准确率 100%
- [ ] 宠物间协作行为，多宠物协作流畅

---

### Sprint 23-24：仿真与测试

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 虚拟宠物仿真 | P1 | MODULE_SIMULATION.md |
| 自动化测试框架 | P1 | MODULE_SIMULATION.md |
| 回放系统 | P1 | MODULE_SIMULATION.md |
| 仿真场景管理 | P1 | MODULE_SIMULATION.md |
| 压力测试 | P2 | MODULE_SIMULATION.md |
| A/B实验仿真 | P2 | MODULE_SIMULATION.md |
| 用户行为模拟 | P2 | MODULE_SIMULATION.md |

**验收标准：**

- [ ] 虚拟宠物运行环境稳定，虚拟宠物运行流畅
- [ ] 自动化测试用例管理，测试用例可执行
- [ ] 设备行为回放 + 问题复现，回放功能正常
- [ ] 场景库 + 场景编辑器，场景管理完整
- [ ] 并发/性能压力测试，压力测试报告正常
- [ ] 实验效果预验证，仿真结果与真机误差 < 10%
- [ ] 模拟真实用户交互，用户行为模拟正常

---

## 六、Phase 4：生态扩展（Sprint 25-32）

> **目标：** 构建开放平台生态，内容市场，第三方集成

### Sprint 25-26：开放平台

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 开发者API完善 | P1 | MODULE_PLATFORM_ECOSYSTEM.md |
| App/插件市场 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| 表情包市场 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| 动作资源库 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| 声音定制 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| Webhook市场 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| SDK发布 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |

**验收标准：**

- [ ] REST API + GraphQL，API文档完整，API测试通过
- [ ] 插件上传/审核/发布，插件市场流程完整
- [ ] 表情上传/商店/推荐，表情包管理正常
- [ ] 动作编辑器 + 动作市场，动作管理正常
- [ ] TTS个性化 + 语音克隆，声音定制正常
- [ ] 预置Webhook模板，Webhook市场可用
- [ ] iOS/Android/小程序SDK，SDK文档完整

---

### Sprint 27-28：第三方集成

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 智能家居对接 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| 宠物医疗对接 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| 宠物保险对接 | P3 | MODULE_PLATFORM_ECOSYSTEM.md |
| 宠物用品电商 | P3 | MODULE_PLATFORM_ECOSYSTEM.md |
| 社交平台分享 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |
| 地图服务对接 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |

**验收标准：**

- [ ] 米家/天猫精灵/HomeKit集成，智能家居控制正常
- [ ] 兽医预约 + 病历同步，宠物医疗对接正常
- [ ] 保险理赔流程集成，宠物保险对接正常
- [ ] 食品/玩具推荐 + 一键购买，电商推荐转化率 > 5%
- [ ] 微信/抖音分享能力，社交分享正常
- [ ] 高德/Google地图集成，地图服务正常

---

### Sprint 29-30：高级功能

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 儿童模式完善 | P1 | MODULE_OPENCLAW_CONSOLE.md |
| 老人陪伴模式 | P2 | MODULE_OPENCLAW_CONSOLE.md |
| 家庭相册 | P2 | MODULE_OPENCLAW_CONSOLE.md |
| 寻回网络 | P1 | MODULE_DEVICE_MANAGEMENT.md |
| 睡眠分析完善 | P1 | MODULE_DIGITAL_TWIN.md |
| 体重追踪完善 | P1 | MODULE_DIGITAL_TWIN.md |
| 饮食记录 | P2 | MODULE_DIGITAL_TWIN.md |
| 宠物社交 | P2 | MODULE_PLATFORM_ECOSYSTEM.md |

**验收标准：**

- [ ] 使用时长 + 内容过滤 + 报告，儿童模式功能完整
- [ ] 简化UI + 主动问候 + 紧急求助，老人陪伴模式正常
- [ ] 宠物视角照片 + 家庭共享，家庭相册功能完整
- [ ] 丢失宠物协查 + 定位追踪，寻回网络正常
- [ ] 深睡/浅睡/REM分析，睡眠分析报告完整
- [ ] 体重曲线 + 营养建议，体重追踪正常
- [ ] 每餐记录 + 热量统计，饮食记录正常
- [ ] 宠物视频互动 + 朋友圈，宠物社交正常

---

### Sprint 31-32：平台演进

**功能清单：**

| 功能 | 优先级 | 对应 PRD 文档 |
|------|--------|--------------|
| 端侧推理完善 | P1 | MODULE_MINICLAW_FIRMWARE.md |
| 模型分片加载完善 | P1 | MODULE_AI_ENGINEERING.md |
| BLE Mesh完善 | P2 | MODULE_MINICLAW_FIRMWARE.md |
| RTOS深度优化 | P2 | MODULE_MINICLAW_FIRMWARE.md |
| 数据集开放 | P3 | MODULE_PLATFORM_ECOSYSTEM.md |
| AI行为研究平台 | P3 | MODULE_PLATFORM_ECOSYSTEM.md |

**验收标准：**

- [ ] 设备端模型优化 + 隐私保护，端侧推理延迟 < 100ms
- [ ] 动态分片 + 增量更新，模型分片加载正常
- [ ] 宠物间直接通信，BLE Mesh通信正常
- [ ] 实时性能优化，RTOS性能提升 > 30%
- [ ] 学术研究数据共享，数据集开放正常
- [ ] 学术合作实验平台，AI行为研究平台正常

---

## 七、功能优先级矩阵

### P0 — MVP 必须（16个）

| 功能 | Sprint | 所属领域 | 对应 PRD 文档 |
|------|--------|----------|--------------|
| OTA Worker | S1 | D3-设备管理 | MODULE_OTA_UPDATES.md |
| AI 模型监控 | S5 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| 模型热回滚 | S5 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| AI 沙箱测试 | S6 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| AI 行为监控/异常告警 | S13 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| LDAP/AD 集成 | S9 | D2-安全 | MODULE_PERMISSIONS.md |
| 证书管理 | S9 | D2-安全 | MODULE_DEVICE_MANAGEMENT.md |
| 远程锁定/擦除 | S9 | D2-安全 | MODULE_DEVICE_MANAGEMENT.md |
| 订阅管理 | S15 | D5-商业化 | MODULE_SUBSCRIPTION.md |
| 设备注册与配对 | S1 | D3-设备管理 | MODULE_DEVICE_PAIRING.md |
| 设备配对注册完善 | S9 | D3-设备管理 | MODULE_DEVICE_PAIRING.md |
| 情绪识别 | S17 | D1-情感计算 | MODULE_AFFECTIVE_COMPUTING.md |
| 情绪响应 | S17 | D1-情感计算 | MODULE_AFFECTIVE_COMPUTING.md |
| 实时生命体征 | S19 | D1-数字孪生 | MODULE_DIGITAL_TWIN.md |
| 环境感知 | S21 | D1-具身智能 | MODULE_EMBODIED_AI.md |
| 空间认知 | S21 | D1-具身智能 | MODULE_EMBODIED_AI.md |
| 儿童模式 | S29 | D4-用户管理 | MODULE_OPENCLAW_CONSOLE.md |
| 寻回网络 | S29 | D4-用户管理 | MODULE_DEVICE_MANAGEMENT.md |

### P1 — 完整产品（20个）

| 功能 | Sprint | 所属领域 | 对应 PRD 文档 |
|------|--------|----------|--------------|
| 设备影子完善 | S1 | D3-设备管理 | MODULE_DEVICE_SHADOW.md |
| 告警通知渠道 | S1 | D3-运维 | MODULE_NOTIFICATION.md |
| 系统管理基础功能 | S1 | D3-设备管理 | MODULE_SYSTEM_MANAGEMENT.md |
| 宠物配置完善 | S3 | D4-宠物管理 | MODULE_PET_BEHAVIOR_ENGINE.md |
| 多用户交互 | S3 | D4-用户管理 | MODULE_OWNER_PROFILE.md |
| 宠物记忆系统 | S3 | D4-宠物管理 | MODULE_PET_MEMORY.md |
| 组织架构管理 | S3 | D4-用户管理 | MODULE_ORGANIZATION.md |
| AI模型训练流水线 | S5 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| A/B测试框架 | S5 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| 地理围栏 | S11 | D3-设备管理 | MODULE_DEVICE_MANAGEMENT.md |
| 跨境设备统一管控 | S11 | D3-设备管理 | MODULE_DEVICE_MANAGEMENT.md |
| 数据驻留合规 | S11 | D2-安全 | MODULE_DATA_ANALYSIS.md |
| 区域AI节点 | S11 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| 优惠券/促销 | S8 | D5-商业化 | MODULE_MEMBER_MANAGEMENT.md |
| Webhook事件系统 | S16 | D5-商业化 | MODULE_SUBSCRIPTION.md |
| AI决策可解释性 | S13 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| AI决策日志 | S5 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| 运动追踪 | S8 | D4-健康管理 | MODULE_DIGITAL_TWIN.md |
| 多宠物管理 | S8 | D4-宠物管理 | MODULE_PET_BEHAVIOR_ENGINE.md |
| 数据脱敏 | S10 | D2-安全 | MODULE_DATA_ANALYSIS.md |
| 合规策略强制 | S9 | D2-安全 | MODULE_POLICY_MANAGEMENT.md |
| 多时区支持 | S11 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| 开发者API | S16 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 设备监控面板 | S9 | D3-设备管理 | MODULE_DEVICE_MONITORING.md |
| 知识库/天气/问答 | S9 | D1-AI工程 | MODULE_KNOWLEDGE_BASE.md |

### P2 — 差异化竞争（25个）

| 功能 | Sprint | 所属领域 | 对应 PRD 文档 |
|------|--------|----------|--------------|
| 边缘AI vs 云端AI | S6 | D1-AI工程 | MODULE_AI_ENGINEERING.md |
| 宠物表情情绪识别 | S17 | D1-情感计算 | MODULE_AFFECTIVE_COMPUTING.md |
| 情绪日志 | S17 | D1-情感计算 | MODULE_AFFECTIVE_COMPUTING.md |
| 行为预测 | S19 | D1-数字孪生 | MODULE_DIGITAL_TWIN.md |
| 健康预警 | S19 | D1-数字孪生 | MODULE_DIGITAL_TWIN.md |
| 历史回放 | S19 | D1-数字孪生 | MODULE_DIGITAL_TWIN.md |
| 跨设备状态同步 | S19 | D1-数字孪生 | MODULE_DEVICE_SHADOW.md |
| 自主探索 | S21 | D1-具身智能 | MODULE_EMBODIED_AI.md |
| 动作模仿 | S21 | D1-具身智能 | MODULE_EMBODIED_AI.md |
| 虚拟宠物仿真 | S23 | D1-仿真测试 | MODULE_SIMULATION.md |
| 自动化测试框架 | S23 | D1-仿真测试 | MODULE_SIMULATION.md |
| App/插件市场 | S25 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 表情包市场 | S25 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 动作资源库 | S25 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 声音定制 | S25 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 订阅自动续费 | S15 | D5-商业化 | MODULE_SUBSCRIPTION.md |
| 用量计费 | S15 | D5-商业化 | MODULE_SUBSCRIPTION.md |
| 睡眠分析 | S29 | D4-健康管理 | MODULE_DIGITAL_TWIN.md |
| 体重追踪 | S29 | D4-健康管理 | MODULE_DIGITAL_TWIN.md |
| 老人陪伴模式 | S29 | D4-用户管理 | MODULE_OPENCLAW_CONSOLE.md |
| 家庭相册 | S29 | D4-用户管理 | MODULE_OPENCLAW_CONSOLE.md |
| 智能家居对接 | S27 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 宠物医疗对接 | S27 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 端侧推理 | S14 | D6-技术架构 | MODULE_AI_ENGINEERING.md |
| 模型分片加载 | S14 | D6-技术架构 | MODULE_AI_ENGINEERING.md |
| 应用分发/企业应用商店 | S15 | D5-商业化 | MODULE_APP_MANAGEMENT.md |
| 文件库/内容分发管理 | S15 | D5-商业化 | MODULE_CONTENT_MANAGEMENT.md |

### P3 — 生态护城河（15个）

| 功能 | Sprint | 所属领域 | 对应 PRD 文档 |
|------|--------|----------|--------------|
| 家庭情绪地图 | S17 | D1-情感计算 | MODULE_AFFECTIVE_COMPUTING.md |
| 情绪趋势分析 | S17 | D1-情感计算 | MODULE_AFFECTIVE_COMPUTING.md |
| 精彩瞬间AI筛选 | S19 | D1-数字孪生 | MODULE_DIGITAL_TWIN.md |
| 离线支持 | S19 | D1-数字孪生 | MODULE_DEVICE_SHADOW.md |
| 多宠物协作 | S21 | D1-具身智能 | MODULE_EMBODIED_AI.md |
| 回放系统 | S23 | D1-仿真测试 | MODULE_SIMULATION.md |
| 压力测试 | S23 | D1-仿真测试 | MODULE_SIMULATION.md |
| BLE Mesh组网 | S12 | D6-技术架构 | MODULE_MINICLAW_FIRMWARE.md |
| RTOS深度优化 | S31 | D6-技术架构 | MODULE_MINICLAW_FIRMWARE.md |
| 宠物社交 | S29 | D4-宠物管理 | MODULE_PLATFORM_ECOSYSTEM.md |
| 发票/账单 | S16 | D5-商业化 | MODULE_SUBSCRIPTION.md |
| API配额计费 | S16 | D5-商业化 | MODULE_SUBSCRIPTION.md |
| 宠物保险对接 | S27 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| 数据集开放 | S31 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |
| AI行为研究平台 | S31 | D5-生态 | MODULE_PLATFORM_ECOSYSTEM.md |

---

## 八、模块依赖关系

```
Phase 1 (Sprint 1-8)
    │
    ├─ D3设备管理 ─────────→ D1-AI工程基础
    ├─ D4宠物管理 ─────────→ D1-情感计算
    └─ D5商业化 ───────────→ Phase 2
                                │
Phase 2 (Sprint 9-16) ─────────┤
    │                           │
    ├─ D2安全 ─────────────────→┘
    ├─ D1-AI工程 ──────────────→ Phase 3
    └─ D5商业化 ──────────────→ Phase 4
                                    │
Phase 3 (Sprint 17-24) ────────────┤
    │                               │
    ├─ D1-情感计算 ────────────────┤
    ├─ D1-数字孪生 ────────────────┤
    ├─ D1-具身智能 ────────────────┤
    └─ D1-仿真测试 ────────────────┘
                                        │
Phase 4 (Sprint 25-32) ←───────────────┘
    │
    ├─ D5-生态（市场/第三方集成）
    ├─ D6-技术架构（端侧推理/BLE Mesh）
    └─ D4-高级功能（儿童/老人/社交）
```

---

## 九、估算工作量

| Phase | Sprint | 功能点数 | 估算人天 | 说明 |
|-------|--------|----------|---------|------|
| Phase 1 | S1-S8 | 32 | ~128 | 核心平台补齐+AI工程基础 |
| Phase 2 | S9-S16 | 29 | ~116 | 企业安全+全球化+商业化 |
| Phase 3 | S17-S24 | 27 | ~108 | 具身智能+情感+数字孪生+仿真 |
| Phase 4 | S25-S32 | 26 | ~104 | 生态+第三方集成+平台演进 |
| **总计** | **S1-S32** | **76** | **~456** | |

---

## 十、验收标准总览

### Phase 1 验收标准

| 验收点 | 标准 |
|--------|------|
| OTA Worker | 1. 创建部署任务后，Worker自动选中目标设备并下发MQTT OTA指令<br>2. 设备上报进度后，OTAProgress表正确更新 |
| AI模型训练 | 1. 训练任务可创建/暂停/取消<br>2. 训练状态实时追踪<br>3. 训练完成自动触发评估 |
| A/B测试 | 1. 实验配置支持多分组<br>2. 分组流量均匀分布<br>3. 实验效果数据自动统计 |
| 模型监控 | 1. 响应延迟/成功率实时展示<br>2. 异常告警自动触发 |
| 模型热回滚 | 1. 一键切换到历史版本<br>2. 切换时间<30秒 |
| 宠物配置 | 1. 宠物配置CRUD完整<br>2. 性格/交互频率正确下发到设备 |
| 告警通知 | 1. 设备异常自动触发邮件/SMS通知<br>2. 告警确认/解决流程完整 |

### Phase 2 验收标准

| 验收点 | 标准 |
|--------|------|
| LDAP/AD集成 | 1. 企业用户SSO登录成功<br>2. 部门/用户自动同步 |
| 证书管理 | 1. 设备证书自动颁发/续期<br>2. 丢失设备可远程吊销 |
| 远程锁定/擦除 | 1. 锁定指令5秒内到达设备<br>2. 擦除后数据不可恢复 |
| 地理围栏 | 1. 围栏边界触发准确率>95%<br>2. 进入/离开事件实时通知 |
| 订阅管理 | 1. 多级订阅自动续费<br>2. 升级/降级即时生效 |

### Phase 3 验收标准

| 验收点 | 标准 |
|--------|------|
| 情绪识别 | 1. 语音/文字/表情情绪识别准确率>85%<br>2. 响应延迟<500ms |
| 情绪响应 | 1. 基于情绪的响应动作正确触发<br>2. 用户满意度评分>4.0 |
| 数字孪生 | 1. 生命体征延迟<2秒<br>2. 历史数据完整可查 |
| 具身智能 | 1. 环境感知准确率>90%<br>2. 自主导航成功率>85% |
| 仿真测试 | 1. 虚拟宠物运行稳定<br>2. 自动化测试覆盖率>80% |

### Phase 4 验收标准

| 验收点 | 标准 |
|--------|------|
| 开放平台 | 1. 开发者API文档完整<br>2. 插件审核流程可用 |
| 第三方集成 | 1. 至少3个智能家居平台对接完成<br>2. 电商推荐转化率>5% |
| 端侧推理 | 1. 设备端推理延迟<100ms<br>2. 隐私数据不出设备 |

---

## 十一、修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿创建，模块现状分析和路线图规划 |
| V1.1 | 2026-03-20 | agentcp | 新增多租户Sprint规划 |
| V2.0 | 2026-03-22 | agentcp | 按6大领域重组，新增Phase 1-4完整规划（32个Sprint） |
| V2.2 | 2026-03-22 | agentcp | 补充缺失的9个MODULE PRD引用，完善功能清单与验收标准 |
| V2.1 | 2026-03-22 | agentcp | 76个功能点逐项列出，UI规范引用PRD文档，验收标准完善 |
