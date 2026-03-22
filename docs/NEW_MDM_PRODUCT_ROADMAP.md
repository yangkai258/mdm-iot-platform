# MDM 控制中台 — 产品路线图

**版本：** V2.3
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
| **D4 用户与宠物管理** | 宠物登记、运动追踪、睡眠分析、寻回网络、家庭相册、**临时会员/访客会员** | 宠物管理、用户管理、健康管理、**临时会员** |
| **D5 生态、商业与集成** | 订阅、用量计费、Webhook、开发者API、市场 | 订阅管理、开放平台、支付计费 |
| **D6 技术架构高阶** | RTOS优化、BLE Mesh、端侧推理、模型分片、**OpenClaw AI 版本管理**、**系统参数/字典管理** | 技术架构、OpenClaw 控制台、系统管理 |

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
| **字典管理** | P1 | MODULE_SYSTEM_MANAGEMENT.md |
| 告警通知渠道（邮件/SMS/Webhook） | P1 | MODULE_NOTIFICATION.md |
| 系统管理基础功能（用户/角色/权限/RBAC） | P1 | MODULE_SYSTEM_MANAGEMENT.md |
| 设备注册与配对流程 | P0 | MODULE_DEVICE_PAIRING.md |
| **设备分组/标签管理** | P1 | MODULE_DEVICE_MANAGEMENT.md |
| **设备分享/授权** | P2 | MODULE_DEVICE_MANAGEMENT.md |
| **设备健康评分** | P1 | MODULE_DEVICE_MANAGEMENT.md |

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
| **告警升级机制**（持续未处理自动升级严重程度） | P1 | MODULE_ALERT_SYSTEM.md |
| **告警自愈建议**（根据告警类型推荐处理方案） | P1 | MODULE_ALERT_SYSTEM.md |
| 宠物记忆系统（交互历史/学习记忆） | P1 | MODULE_PET_MEMORY.md |
| 组织架构管理（部门/团队/成员） | P1 | MODULE_ORGANIZATION.md |
| **公司管理基础** | P1 | MODULE_ORGANIZATION.md |
| **OTA 灰度推送策略**（百分比/白名单） | P1 | MODULE_OTA_UPDATES.md |
| **OTA 升级预约** | P2 | MODULE_OTA_UPDATES.md |

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
| **OpenClaw AI 版本管理** | P2 | MODULE_OPENCLAW_VERSION.md |
| AI沙箱测试环境 | P1 | MODULE_AI_ENGINEERING.md |
| 边缘AI vs 云端AI路由 | P2 | MODULE_AI_ENGINEERING.md |
| **设备影子版本历史记录** | P1 | MODULE_DEVICE_SHADOW.md |
| **设备状态预测**（基于历史数据的趋势预测） | P2 | MODULE_DEVICE_SHADOW.md |
| **影子快照导出** | P2 | MODULE_DEVICE_SHADOW.md |
| **OTA 固件版本兼容性自动检测** | P1 | MODULE_OTA_UPDATES.md |

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
| **会员储值充值/扣款** | P1 | MODULE_MEMBER_MANAGEMENT.md |
| **会员成长值/等级自动升级触发器** | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 多宠物管理 | P1 | MODULE_PET_BEHAVIOR_ENGINE.md |
| 运动追踪基础 | P1 | MODULE_DIGITAL_TWIN.md |
| **临时会员/访客会员** | P2 | MODULE_MEMBER_MANAGEMENT.md |

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
| **监控告警规则配置**（基于监控指标的自定义告警） | P1 | MODULE_DEVICE_MONITORING.md |
| **多设备对比视图** | P2 | MODULE_DEVICE_MONITORING.md |
| 设备配对注册完善 | P0 | MODULE_DEVICE_PAIRING.md |
| **配对码批量生成** | P2 | MODULE_DEVICE_PAIRING.md |
| 知识库/天气/问答 | P1 | MODULE_KNOWLEDGE_BASE.md |
| **知识库版本管理** | P2 | MODULE_KNOWLEDGE_BASE.md |
| **知识库审核工作流** | P2 | MODULE_KNOWLEDGE_BASE.md |
| **应用权限管理**（应用所需权限列表和用户授权） | P1 | MODULE_APP_MANAGEMENT.md |
| **应用使用时长统计** | P2 | MODULE_APP_MANAGEMENT.md |

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
| **订阅赠送功能**（赠送好友体验） | P2 | MODULE_SUBSCRIPTION.md |
| **家庭计划**（多用户共享订阅） | P2 | MODULE_SUBSCRIPTION.md |
| **订阅试用期延长** | P3 | MODULE_SUBSCRIPTION.md |
| 用量计费 | P2 | MODULE_SUBSCRIPTION.md |
| API配额计费 | P2 | MODULE_SUBSCRIPTION.md |
| 优惠券/促销完善 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 发票/账单 | P2 | MODULE_SUBSCRIPTION.md |
| Webhook事件系统 | P1 | MODULE_SUBSCRIPTION.md |
| 开发者API | P1 | MODULE_PLATFORM_ECOSYSTEM.md |
| 会员等级权益配置 | P1 | MODULE_MEMBER_MANAGEMENT.md |
| 应用分发/企业应用商店 | P2 | MODULE_APP_MANAGEMENT.md |
| 文件库/内容分发管理 | P2 | MODULE_CONTENT_MANAGEMENT.md |
| **内容版本管理** | P2 | MODULE_CONTENT_MANAGEMENT.md |
| **内容评论/批注** | P3 | MODULE_CONTENT_MANAGEMENT.md |

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
| **情感预警机制**（情绪持续低落时通知家属） | P1 | MODULE_AFFECTIVE_COMPUTING.md |
| **情感数据脱敏和匿名化处理** | P2 | MODULE_AFFECTIVE_COMPUTING.md |
| **情感模型持续学习/重训练** | P2 | MODULE_AFFECTIVE_COMPUTING.md |
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
| **数字孪生备份/克隆** | P2 | MODULE_DIGITAL_TWIN.md |
| **第三方健康设备数据接入**（智能项圈等） | P2 | MODULE_DIGITAL_TWIN.md |
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
| **具身AI模型压缩和优化工具链** | P2 | MODULE_EMBODIED_AI.md |
| **动作模仿学习进度可视化** | P2 | MODULE_EMBODIED_AI.md |
| **具身AI安全审计日志** | P1 | MODULE_EMBODIED_AI.md |

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
| **仿真环境与CI/CD流水线集成** | P1 | MODULE_SIMULATION.md |
| **仿真结果自动生成测试报告** | P1 | MODULE_SIMULATION.md |
| **仿真数据集管理**（测试数据版本化） | P2 | MODULE_SIMULATION.md |
| **仿真资源配额管理** | P2 | MODULE_SIMULATION.md |
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

### P0 — MVP 必须（18个）

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
| 字典管理 | S2 | D1-核心平台 | MODULE_SYSTEM_MANAGEMENT.md |
| 宠物配置完善 | S3 | D4-宠物管理 | MODULE_PET_BEHAVIOR_ENGINE.md |
| 多用户交互 | S3 | D4-用户管理 | MODULE_OWNER_PROFILE.md |
| 宠物记忆系统 | S3 | D4-宠物管理 | MODULE_PET_MEMORY.md |
| 组织架构管理 | S3 | D4-用户管理 | MODULE_ORGANIZATION.md |
| 公司管理基础 | S3 | D4-宠物管理 | MODULE_ORGANIZATION.md |
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
| **OpenClaw AI 版本管理** | S5 | D6-生态 | MODULE_OPENCLAW_VERSION.md |
| **临时会员/访客会员** | S7 | D5-会员运营 | MODULE_MEMBER_MANAGEMENT.md |
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

## 十一、API 接口规范（补充）

> **说明：** 以下 API 为根据 PRD_GAP_ANALYSIS.md 差距分析补充的缺失接口，与现有 API 共同构成完整接口规范。

### 11.1 设备管理（MODULE_DEVICE_MANAGEMENT）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/devices/:id/health-score` | GET | 设备健康评分 | S1-2 |
| `POST /api/v1/devices/:id/share` | POST | 设备分享/授权 | S1-2 |
| `DELETE /api/v1/devices/:id/share/:target_user_id` | DELETE | 取消设备分享 | S1-2 |
| `GET /api/v1/devices/groups` | GET | 设备分组列表 | S1-2 |
| `POST /api/v1/devices/groups` | POST | 创建设备分组 | S1-2 |
| `PUT /api/v1/devices/groups/:id` | PUT | 更新设备分组 | S1-2 |
| `DELETE /api/v1/devices/groups/:id` | DELETE | 删除设备分组 | S1-2 |
| `POST /api/v1/devices/:id/tags` | POST | 为设备添加标签 | S1-2 |
| `DELETE /api/v1/devices/:id/tags/:tag` | DELETE | 移除设备标签 | S1-2 |
| `GET /api/v1/devices/:id/usage-stats` | GET | 设备使用时长/频率统计 | S1-2 |

### 11.2 OTA 升级（MODULE_OTA_UPDATES）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/ota/deployments/:id/pause` | POST | 暂停升级任务 | S3 |
| `POST /api/v1/ota/deployments/:id/resume` | POST | 恢复升级任务 | S3 |
| `POST /api/v1/ota/deployments/:id/gray` | POST | 灰度发布（百分比/白名单） | S3 |
| `GET /api/v1/ota/devices/:device_id/compatibility` | GET | 固件兼容性查询 | S5 |
| `POST /api/v1/ota/deployments/:id/schedule` | POST | OTA 升级预约 | S3 |
| `DELETE /api/v1/ota/deployments/:id/schedule` | DELETE | 取消预约升级 | S3 |

### 11.3 设备影子（MODULE_DEVICE_SHADOW）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/shadows/:device_id/history` | GET | 影子版本历史 | S5 |
| `GET /api/v1/shadows/:device_id/prediction` | GET | 设备状态预测 | S5 |
| `POST /api/v1/shadows/:device_id/export` | POST | 导出影子快照 | S5 |
| `POST /api/v1/shadows/:device_id/anomaly/mark` | POST | 标记状态异常 | S5 |

### 11.4 告警系统（MODULE_ALERT_SYSTEM）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/alerts/:id/escalate` | POST | 告警升级（提升严重程度） | S4 |
| `GET /api/v1/alerts/suggestions/:id` | GET | 自愈建议（推荐处理方案） | S4 |
| `GET /api/v1/alerts/rules/:id/trigger-history` | GET | 规则触发历史 | S4 |
| `POST /api/v1/alerts/rules/:id/suppress` | POST | 配置告警抑制/去重规则 | S4 |

### 11.5 会员管理（MODULE_MEMBER_MANAGEMENT）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/members/:id/balance/recharge` | POST | 储值充值 | S7 |
| `POST /api/v1/members/:id/balance/deduct` | POST | 储值扣款 | S7 |
| `GET /api/v1/members/:id/balance/transactions` | GET | 储值交易流水 | S7 |
| `POST /api/v1/members/:id/upgrade/check` | POST | 触发等级检查 | S7 |
| `DELETE /api/v1/members/:id` | DELETE | 会员注销 | S7 |
| `GET /api/v1/members/:id/profile-tags` | GET | 会员画像标签 | S7 |

### 11.6 通知与消息（MODULE_NOTIFICATION）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/notifications/:id/recall` | POST | 撤回通知 | S4 |
| `PUT /api/v1/notifications/:id/read-receipt` | PUT | 已读回执确认 | S4 |
| `GET /api/v1/notifications/templates/:id/preview` | GET | 模板预览 | S4 |
| `POST /api/v1/notifications/schedules` | POST | 创建定时发送任务 | S4 |
| `DELETE /api/v1/notifications/schedules/:id` | DELETE | 删除定时发送任务 | S4 |

### 11.7 应用管理（MODULE_APP_MANAGEMENT）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/apps/:id/permissions` | GET | 应用所需权限列表 | S10 |
| `POST /api/v1/apps/:id/permissions/approve` | POST | 授权应用权限 | S10 |
| `DELETE /api/v1/apps/:id/permissions/:permission_id` | DELETE | 撤销应用权限 | S10 |
| `GET /api/v1/apps/:id/usage-stats` | GET | 应用使用时长统计 | S10 |

### 11.8 内容管理（MODULE_CONTENT_MANAGEMENT）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/contents/:id/versions` | GET | 内容版本历史 | S15 |
| `POST /api/v1/contents/:id/comment` | POST | 添加内容评论 | S15 |
| `GET /api/v1/contents/:id/comments` | GET | 获取内容评论列表 | S15 |
| `POST /api/v1/contents/:id/subscribe` | POST | 订阅内容更新 | S15 |
| `DELETE /api/v1/contents/:id/subscribe` | DELETE | 取消内容订阅 | S15 |

### 11.9 数据分析（MODULE_DATA_ANALYSIS）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/stats/reports/custom` | POST | 生成自定义报表 | S13 |
| `GET /api/v1/stats/reports/export` | GET | 导出报表（Excel/PDF） | S13 |
| `POST /api/v1/stats/reports/subscribe` | POST | 报表订阅/定时推送 | S13 |
| `GET /api/v1/stats/reports/:id` | GET | 获取报表详情 | S13 |

### 11.10 系统管理（MODULE_SYSTEM_MANAGEMENT）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/system/users/:id/force-logout` | POST | 强制下线用户 | S4 |
| `GET /api/v1/system/health` | GET | 系统健康检查 | S4 |
| `GET /api/v1/system/config/history` | GET | 配置变更历史 | S4 |

### 11.11 订阅管理（MODULE_SUBSCRIPTION）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/subscriptions/gift` | POST | 赠送订阅 | S15 |
| `POST /api/v1/subscriptions/family/join` | POST | 加入家庭计划 | S15 |
| `GET /api/v1/subscriptions/family/:plan_id` | GET | 家庭计划详情 | S15 |
| `GET /api/v1/subscriptions/analytics` | GET | 订阅数据分析（付费转化率/流失率） | S15 |

### 11.12 数字孪生（MODULE_DIGITAL_TWIN）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/digital-twin/:pet_id/backup` | POST | 创建数字孪生备份 | S19 |
| `GET /api/v1/digital-twin/:pet_id/backups` | GET | 备份列表 | S19 |
| `POST /api/v1/digital-twin/:pet_id/restore/:backup_id` | POST | 恢复备份 | S19 |
| `POST /api/v1/digital-twin/:pet_id/health-devices/sync` | POST | 第三方健康设备数据同步 | S19 |
| `GET /api/v1/digital-twin/:pet_id/health-devices` | GET | 已绑定健康设备列表 | S19 |

### 11.13 情感计算（MODULE_AFFECTIVE_COMPUTING）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/emotion/model/retrain` | POST | 触发情感模型重训练 | S17 |
| `POST /api/v1/emotion/:pet_id/alert/notify` | POST | 情感预警通知（通知家属） | S17 |
| `GET /api/v1/emotion/:pet_id/reports/daily` | GET | 每日情绪报告 | S17 |
| `GET /api/v1/emotion/:pet_id/reports/weekly` | GET | 每周情绪趋势报告 | S17 |

### 11.14 具身智能（MODULE_EMBODIED_AI）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/embodied/:device_id/model/compress` | POST | 触发模型压缩 | S21 |
| `GET /api/v1/embodied/:device_id/learning/progress` | GET | 动作模仿学习进度 | S21 |
| `GET /api/v1/embodied/:device_id/safety/logs` | GET | 安全审计日志 | S21 |

### 11.15 仿真测试（MODULE_SIMULATION）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/simulation/cicd/integrate` | POST | CI/CD 集成配置 | S23 |
| `GET /api/v1/simulation/cicd/status` | GET | CI/CD 集成状态 | S23 |
| `GET /api/v1/simulation/datasets` | GET | 仿真数据集列表 | S23 |
| `POST /api/v1/simulation/datasets` | POST | 创建仿真数据集 | S23 |
| `GET /api/v1/simulation/quotas` | GET | 仿真资源配额查询 | S23 |

### 11.16 知识库（MODULE_KNOWLEDGE_BASE）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/knowledge/entries/:id/history` | GET | 知识版本历史 | S9 |
| `POST /api/v1/knowledge/entries/batch-import` | POST | 批量导入知识 | S9 |
| `GET /api/v1/knowledge/top-queries` | GET | 高频查询排行 | S9 |
| `POST /api/v1/knowledge/entries/:id/audit` | POST | 提交知识审核 | S9 |

### 11.17 设备监控（MODULE_DEVICE_MONITORING）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/monitoring/alerts/rules` | POST | 配置监控告警规则 | S9 |
| `GET /api/v1/monitoring/alerts/rules` | GET | 监控告警规则列表 | S9 |
| `GET /api/v1/devices/:device_id/sensor-history/export` | GET | 导出传感器历史数据 | S9 |

### 11.18 设备配对（MODULE_DEVICE_PAIRING）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/devices/pairing/codes/batch` | POST | 批量生成配对码 | S9 |
| `DELETE /api/v1/devices/pairing/codes/cleanup` | DELETE | 清理过期配对码 | S9 |
| `GET /api/v1/devices/pairing/stats` | GET | 配对统计看板 | S9 |

### 11.19 AI 工程（MODULE_AI_ENGINEERING）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/ai/gpu/status` | GET | GPU 资源状态 | S5 |
| `POST /api/v1/ai/scaling/trigger` | POST | 触发自动扩缩容 | S5 |
| `GET /api/v1/ai/scaling/history` | GET | 扩缩容历史记录 | S5 |

### 11.20 宠物行为引擎（MODULE_PET_BEHAVIOR_ENGINE）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/behavior/sequences/:id/versions` | POST | 创建动作序列新版本 | S8 |
| `GET /api/v1/behavior/actions/compatibility/:device_model` | GET | 查询设备兼容动作 | S8 |

### 11.21 宠物记忆（MODULE_PET_MEMORY）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `POST /api/v1/memory/:memory_id/restore` | POST | 恢复误删记忆 | S8 |
| `GET /api/v1/memory/:device_id/quality-report` | GET | 记忆质量报告 | S8 |
| `GET /api/v1/memory/:device_id/recycle-bin` | GET | 记忆回收站列表 | S8 |

### 11.22 主人画像（MODULE_OWNER_PROFILE）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/owner-profile/:user_id/export` | GET | 画像数据导出 | S8 |
| `PUT /api/v1/owner-profile/:user_id/privacy` | PUT | 更新隐私设置 | S8 |

### 11.23 订阅赠送（MODULE_SUBSCRIPTION）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/subscriptions/gifts/:gift_id/status` | GET | 赠送状态查询 | S15 |

### 11.24 开放平台生态（MODULE_PLATFORM_ECOSYSTEM）

| 接口 | 方法 | 说明 | Sprint |
|------|------|------|--------|
| `GET /api/v1/developer/level` | GET | 查询开发者等级 | S25 |
| `POST /api/v1/developer/refund` | POST | API 调用退款 | S25 |

---

## 十二、UI 页面规范（补充）

> **说明：** 以下 UI 规范为根据 PRD_GAP_ANALYSIS.md 差距分析补充的缺失页面规范，覆盖所有 AI 增强模块。

### 12.1 设备管理（Sprint 1-2）

#### 设备分组管理页面
- **页面路径**: `/devices/groups`
- **页面布局**: 列表式（左树右表）
- **核心组件**: 树形分组列表、设备卡片表格、分组设备数统计
- **交互规范**: 新建/编辑分组弹出 Drawer；删除前二次确认；拖拽调整分组顺序
- **按钮**: [新建分组] 靠左, [导出] 靠右

#### 设备标签管理页面
- **页面路径**: `/devices/tags`
- **页面布局**: 卡片式（标签云 + 筛选面板）
- **核心组件**: 标签云、已标设备列表、设备标签筛选器
- **交互规范**: 点击标签筛选设备；标签颜色自定义；支持批量打标
- **按钮**: [新建标签] 靠左, [批量打标] 靠右

#### 设备健康评分卡片
- **页面路径**: `/devices/:id/detail`（内嵌）
- **页面布局**: 卡片式（评分圆环 + 趋势图）
- **核心组件**: 健康评分圆环（0-100）、7日趋势折线图、评分因素列表
- **交互规范**: 评分低于阈值高亮警示；点击因素查看详细指标；刷新间隔 5 分钟
- **按钮**: [刷新评分] 靠右

#### 设备分享页面
- **页面路径**: `/devices/:id/share`
- **页面布局**: 抽屉式（右侧滑出）
- **核心组件**: 分享用户列表、权限级别选择（只读/控制/管理）、有效期设置
- **交互规范**: 添加分享对象自动发送通知；修改权限即时生效；移除需二次确认
- **按钮**: [添加分享对象] 靠左, [关闭] 靠右

### 12.2 OTA 升级（Sprint 3）

#### 灰度发布配置页面
- **页面路径**: `/ota/deployments/:id/gray`
- **页面布局**: 步骤式（三步引导）
- **核心组件**: 策略选择（百分比/白名单）、设备选择器、灰度进度条
- **交互规范**: 百分比拖拽实时预览影响设备数；白名单支持批量导入；灰度暂停/恢复
- **按钮**: [上一步] [下一步] 靠右, [取消] 靠左

#### OTA 固件兼容性检测页面
- **页面路径**: `/ota/compatibility-check`
- **页面布局**: 表单 + 结果展示
- **核心组件**: 固件版本选择器、设备模型下拉、兼容性结果表格
- **交互规范**: 批量选择多个固件版本对比；结果支持导出 Excel
- **按钮**: [开始检测] 靠左, [导出报告] 靠右

### 12.3 设备影子（Sprint 2）

#### 影子版本历史页面
- **页面路径**: `/devices/:id/shadow/history`
- **页面布局**: 时间线式（左侧时间轴 + 右侧详情）
- **核心组件**: 版本时间轴、版本详情卡片（JSON 展示）、diff 对比视图
- **交互规范**: 点击版本查看完整内容；两版本对比高亮差异；支持快照导出
- **按钮**: [导出快照] 靠左, [对比版本] 靠右

#### 设备状态预测页面
- **页面路径**: `/devices/:id/shadow/prediction`
- **页面布局**: 图表式（大屏看板）
- **核心组件**: 预测趋势折线图、异常标记点、置信区间区域
- **交互规范**: 悬停显示具体时间点预测值；异常点点击查看详情；刷新间隔 1 分钟
- **按钮**: [刷新预测] 靠左, [导出数据] 靠右

### 12.4 告警系统（Sprint 4）

#### 告警升级规则配置页面
- **页面路径**: `/alerts/rules/:id/escalation`
- **页面布局**: 表单式
- **核心组件**: 升级条件配置（时间/次数）、升级目标级别选择、通知渠道配置
- **交互规范**: 条件预览；规则启用/禁用开关；测试升级按钮
- **按钮**: [保存规则] 靠左, [测试] 靠右

#### 自愈建议卡片（内嵌告警详情页）
- **页面路径**: `/alerts/:id`（内嵌）
- **页面布局**: 卡片式（建议列表 + 执行按钮）
- **核心组件**: 建议类型图标、建议描述、执行成功率和历史执行结果
- **交互规范**: 一键执行建议；执行结果 Toast 反馈；建议根据告警类型动态推荐
- **按钮**: [执行建议] 靠左, [忽略] 靠右

#### 告警统计页面
- **页面路径**: `/alerts/stats`
- **页面布局**: 大盘式（顶部 KPI 卡片 + 下方图表）
- **核心组件**: 告警总数/今日新增/未处理趋势折线图、告警类型分布饼图、Top5 告警源柱状图
- **交互规范**: 时间范围选择器（今日/本周/本月）；点击图表钻取详情
- **按钮**: [时间范围] 靠左, [导出报表] 靠右

### 12.5 会员管理（Sprint 7）

#### 储值充值/扣款页面
- **页面路径**: `/members/:id/balance`
- **页面布局**: 卡片式（余额展示 + 操作面板）
- **核心组件**: 当前余额卡片、充值金额快捷选择、充值方式选择、交易流水列表
- **交互规范**: 充值金额自定义输入；扣款需二次确认并填写原因；流水实时更新
- **按钮**: [充值] [扣款] 靠左, [交易明细] 靠右

#### 会员等级自动升级配置页面
- **页面路径**: `/members/levels/upgrade-rules`
- **页面布局**: 表单 + 规则列表
- **核心组件**: 升级条件配置（成长值阈值）、升级动画预览、等级权益预览
- **交互规范**: 保存前预览升级效果；支持模拟升级计算；升级时播放动画
- **按钮**: [新增规则] 靠左, [模拟计算] 靠右

### 12.6 应用管理（Sprint 10）

#### 应用权限管理页面
- **页面路径**: `/apps/:id/permissions`
- **页面布局**: 表格式（权限列表 + 审批流程）
- **核心组件**: 权限名称、权限类型（敏感/普通）、用户授权状态、审批按钮
- **交互规范**: 批量授权；权限申请审批工作流；临时权限显示到期时间
- **按钮**: [申请权限] 靠左, [授权管理] 靠右

#### 应用使用统计页面
- **页面路径**: `/apps/:id/usage-stats`
- **页面布局**: 图表式
- **核心组件**: 使用时长柱状图（日/周/月）、DAU 趋势图、使用峰值时段热力图
- **交互规范**: 时间范围切换；点击柱状图钻取单日详情；支持数据导出
- **按钮**: [时间范围] 靠左, [导出] 靠右

### 12.7 通知与消息（Sprint 4）

#### 定时发送配置页面
- **页面路径**: `/notifications/schedules`
- **页面布局**: 列表 + 表单抽屉
- **核心组件**: 定时任务列表（发送时间/接收人/状态）、草稿编辑表单
- **交互规范**: 定时任务可编辑/暂停/删除；发送前预览；草稿自动保存
- **按钮**: [新建定时发送] 靠左, [批量操作] 靠右

### 12.8 内容管理（Sprint 15）

#### 内容版本历史页面
- **页面路径**: `/contents/:id/versions`
- **页面布局**: 时间线式
- **核心组件**: 版本时间轴、版本对比视图（diff）、版本备注
- **交互规范**: 点击版本恢复；两版本对比高亮变更；版本标签管理
- **按钮**: [版本对比] 靠左, [下载指定版本] 靠右

#### 内容评论区
- **页面路径**: `/contents/:id/comments`
- **页面布局**: 嵌套评论式
- **核心组件**: 评论列表、评论输入框、评论点赞/踩
- **交互规范**: 支持@提及；评论回复嵌套显示；敏感词过滤
- **按钮**: [发表评论] 靠左, [收起评论] 靠右

### 12.9 数据分析（Sprint 13）

#### 自定义报表生成器页面
- **页面路径**: `/stats/reports/custom`
- **页面布局**: 拖拽式（左侧指标面板 + 中间画布 + 右侧配置）
- **核心组件**: 维度选择面板、指标拖拽画布、图表类型选择器、预览区域
- **交互规范**: 拖拽配置；实时预览；配置可保存为模板；支持定时推送
- **按钮**: [保存报表] 靠左, [预览] [生成报表] 靠右

### 12.10 系统管理（Sprint 4）

#### 系统健康检查大盘
- **页面路径**: `/system/health`
- **页面布局**: 大盘式（KPI 卡片 + 服务状态列表）
- **核心组件**: 系统状态指示灯（正常/警告/异常）、服务组件列表、响应时间趋势图
- **交互规范**: 异常项高亮；点击服务查看详情；历史健康数据查询
- **按钮**: [刷新状态] 靠左, [导出报告] 靠右

#### 操作日志实时推送页面
- **页面路径**: `/system/logs`
- **页面布局**: 实时滚动列表式
- **核心组件**: 日志流（WebSocket 实时）、日志级别过滤、操作类型筛选
- **交互规范**: 日志流自动滚动（可暂停）；点击日志查看详情；支持关键词搜索
- **按钮**: [暂停/继续] 靠左, [导出日志] 靠右

### 12.11 订阅管理（Sprint 15）

#### 订阅赠送页面
- **页面路径**: `/subscriptions/gift`
- **页面布局**: 表单 + 分享面板
- **核心组件**: 套餐选择、接收人填写、有效期设置、分享链接/二维码
- **交互规范**: 生成一次性领取链接；领取状态实时追踪；赠送记录列表
- **按钮**: [生成赠送链接] 靠左, [我的赠送记录] 靠右

#### 家庭计划管理页面
- **页面路径**: `/subscriptions/family`
- **页面布局**: 卡片式（计划概览 + 成员管理）
- **核心组件**: 家庭计划卡片（共享时长/成员数/配额）、成员列表、配额分配配置
- **交互规范**: 成员添加/移除；配额灵活分配；超配额提示；成员使用排行
- **按钮**: [邀请成员] 靠左, [管理配额] 靠右

### 12.12 情感计算（Sprint 17）

#### 情绪识别配置页面
- **页面路径**: `/emotion/config`
- **页面布局**: 表单 + 预览
- **核心组件**: 情绪识别开关（语音/文字/表情）、敏感度调节、响应策略配置、模型版本选择
- **交互规范**: 修改后实时预览效果；支持分宠物配置；配置变更记录
- **按钮**: [保存配置] 靠左, [恢复默认] 靠右

#### 情绪日志查看页面
- **页面路径**: `/emotion/logs`
- **页面布局**: 时间线 + 筛选面板
- **核心组件**: 情绪时间轴（颜色编码情绪类型）、情绪详情卡片、情绪诱因标注
- **交互规范**: 按日期/情绪类型筛选；点击情绪点查看详情；支持导出 Excel
- **按钮**: [筛选] 靠左, [导出日志] 靠右

#### 情绪报告页面
- **页面路径**: `/emotion/reports`
- **页面布局**: 报告式（图文混排）
- **核心组件**: 周/月情绪曲线图、情绪分布饼图、高频情绪事件、情绪触发因素分析
- **交互规范**: 自动生成报告；支持分享；报告模板可定制；数据脱敏展示
- **按钮**: [生成报告] 靠左, [分享报告] [导出 PDF] 靠右

#### 情感预警配置页面
- **页面路径**: `/emotion/alerts`
- **页面布局**: 表单 + 通知预览
- **核心组件**: 预警规则配置（连续低落次数/时长阈值）、通知家属配置、预警等级设置
- **交互规范**: 规则测试；预警记录查询；通知模板自定义
- **按钮**: [保存规则] 靠左, [测试通知] 靠右

### 12.13 数字孪生（Sprint 19）

#### 生命体征仪表盘
- **页面路径**: `/digital-twin/:pet_id/dashboard`
- **页面布局**: 大盘式（实时数据 + 历史图表）
- **核心组件**: 心率/呼吸实时曲线（2 秒刷新）、体温/活动量卡片、睡眠质量雷达图
- **交互规范**: 实时数据高亮异常；点击曲线钻取时间点详情；支持全屏大屏模式
- **按钮**: [全屏] 靠左, [导出数据] 靠右

#### 历史回放页面
- **页面路径**: `/digital-twin/:pet_id/playback`
- **页面布局**: 时间轴 + 回放控制
- **核心组件**: 24 小时时间轴、行为事件标记（进食/睡眠/活动）、关键生命体征曲线
- **交互规范**: 拖拽时间轴回放；事件点快速跳转；支持倍速播放（0.5x/1x/2x）
- **按钮**: [播放] [暂停] 靠左, [导出回放] 靠右

#### 健康预警页面
- **页面路径**: `/digital-twin/:pet_id/health-alerts`
- **页面布局**: 列表 + 详情
- **核心组件**: 预警列表（预警类型/级别/时间）、预警详情（趋势图 + 建议）、处理状态
- **交互规范**: 预警分级颜色标识；一键处理/忽略；关联宠物医疗记录
- **按钮**: [处理] 靠左, [忽略] 靠右

#### 数字孪生备份管理页面
- **页面路径**: `/digital-twin/:pet_id/backups`
- **页面布局**: 卡片列表
- **核心组件**: 备份时间轴、备份大小、备份类型（全量/增量）、恢复预览
- **交互规范**: 一键备份；选择备份点恢复；备份加密存储；备份自动清理策略
- **按钮**: [立即备份] 靠左, [设置自动备份] 靠右

#### 第三方健康设备绑定页面
- **页面路径**: `/digital-twin/:pet_id/health-devices`
- **页面布局**: 设备列表 + 绑定向导
- **核心组件**: 已绑定设备卡片、数据同步状态、同步历史
- **交互规范**: 扫码绑定；同步频率配置；数据冲突处理；设备解绑二次确认
- **按钮**: [绑定新设备] 靠左, [同步设置] 靠右

### 12.14 具身智能（Sprint 21）

#### 环境感知配置页面
- **页面路径**: `/embodied/:device_id/sensors`
- **页面布局**: 表单 + 可视化预览
- **核心组件**: 摄像头标定、障碍物检测阈值、感知灵敏度配置、感知范围可视化图
- **交互规范**: 实时预览感知效果；参数调整即时生效；感知异常告警配置
- **按钮**: [保存配置] 靠左, [标定向导] 靠右

#### 空间认知地图页面
- **页面路径**: `/embodied/:device_id/maps`
- **页面布局**: 2D/3D 地图编辑器
- **核心组件**: 室内地图展示、兴趣点标记、安全区/危险区绘制、实时定位
- **交互规范**: 地图缩放/平移；点击添加兴趣点；安全区颜色编码；定位实时更新
- **按钮**: [编辑地图] 靠左, [保存地图] [全屏] 靠右

#### 探索任务页面
- **页面路径**: `/embodied/:device_id/exploration`
- **页面布局**: 任务卡片 + 地图
- **核心组件**: 探索任务列表（待执行/进行中/已完成）、任务路径预览、执行状态
- **交互规范**: 创建探索任务；任务优先级调整；执行中断可续传；任务历史查询
- **按钮**: [新建任务] 靠左, [任务历史] 靠右

#### 具身 AI 学习进度页面
- **页面路径**: `/embodied/:device_id/learning`
- **页面布局**: 进度仪表盘
- **核心组件**: 学习进度环形图、动作模仿成功率趋势、已学习动作库、设备学习历史
- **交互规范**: 进度实时更新；点击动作查看学习详情；支持重新学习
- **按钮**: [刷新进度] 靠左, [重新学习] 靠右

#### 具身 AI 安全审计日志页面
- **页面路径**: `/embodied/:device_id/safety/logs`
- **页面布局**: 日志列表 + 详情
- **核心组件**: 安全事件时间轴、事件详情（触发条件/响应动作）、安全边界配置
- **交互规范**: 按类型/级别筛选；事件详情关联地图位置；安全规则配置入口
- **按钮**: [导出日志] 靠左, [配置安全边界] 靠右

### 12.15 仿真测试（Sprint 23）

#### CI/CD 集成配置页面
- **页面路径**: `/simulation/cicd`
- **页面布局**: 表单 + 连接状态
- **核心组件**: CI/CD 平台选择（GitHub/GitLab/Jenkins）、仓库配置、Webhook 设置、触发规则
- **交互规范**: 连接测试；触发规则配置；构建历史查看；失败自动告警
- **按钮**: [测试连接] 靠左, [保存配置] 靠右

#### 仿真数据集管理页面
- **页面路径**: `/simulation/datasets`
- **页面布局**: 表格 + 版本管理
- **核心组件**: 数据集列表（名称/版本/大小）、版本历史、导入进度条
- **交互规范**: 批量导入；版本对比；数据集标签管理；存储配额显示
- **按钮**: [导入数据集] 靠左, [批量操作] 靠右

### 12.16 宠物行为引擎（Sprint 8）

#### 动作序列版本历史页面
- **页面路径**: `/behavior/sequences/:id/versions`
- **页面布局**: 版本列表 + diff 对比
- **核心组件**: 版本时间轴、版本差异高亮、版本标签（稳定版/测试版）
- **交互规范**: 版本切换；两版本对比；版本回滚；版本发布审批
- **按钮**: [对比] 靠左, [回滚] [发布] 靠右

#### 设备兼容动作查询页面
- **页面路径**: `/behavior/compatibility`
- **页面布局**: 表单 + 结果表格
- **核心组件**: 设备模型选择、固件版本选择、兼容动作列表（支持/不支持）
- **交互规范**: 选择设备自动筛选兼容动作；点击动作查看详情；支持批量导出
- **按钮**: [查询] 靠左, [导出] 靠右

### 12.17 宠物记忆库（Sprint 8）

#### 记忆回收站页面
- **页面路径**: `/memory/:device_id/recycle-bin`
- **页面布局**: 卡片列表
- **核心组件**: 已删除记忆卡片（预览内容）、删除时间、剩余保留天数
- **交互规范**: 恢复误删记忆（30 天内）；永久删除需二次确认；自动清理过期记忆
- **按钮**: [批量恢复] 靠左, [永久删除] 靠右

#### 记忆质量报告页面
- **页面路径**: `/memory/:device_id/quality`
- **页面布局**: 报告式
- **核心组件**: 质量评分雷达图（完整性/准确性/时效性）、问题记忆列表、优化建议
- **交互规范**: 自动生成报告；点击问题项跳转处理；支持定期自动生成
- **按钮**: [生成报告] 靠左, [导出报告] 靠右

### 12.18 主人画像库（Sprint 8）

#### 画像数据导出页面
- **页面路径**: `/owner-profile/:user_id/export`
- **页面布局**: 表单 + 预览
- **核心组件**: 导出字段选择、格式选择（JSON/Excel/CSV）、预览区域、导出历史
- **交互规范**: 字段预览；导出进度条；历史导出下载；脱敏规则说明
- **按钮**: [开始导出] 靠左, [导出历史] 靠右

#### 主人隐私设置页面
- **页面路径**: `/owner-profile/:user_id/privacy`
- **页面布局**: 分组表单式
- **核心组件**: 数据可见性配置（公开/好友/私有）、画像数据授权管理、第三方共享记录
- **交互规范**: 配置变更即时保存；数据删除请求；隐私政策展示
- **按钮**: [保存设置] 靠左, [删除我的数据] 靠右

### 12.19 知识库（Sprint 9）

#### 知识版本历史页面
- **页面路径**: `/knowledge/entries/:id/history`
- **页面布局**: 时间线式
- **核心组件**: 版本时间轴、版本对比视图（diff）、版本审批状态
- **交互规范**: 点击版本查看详情；版本恢复需审批；版本标签管理
- **按钮**: [版本对比] 靠左, [恢复版本] 靠右

#### 知识审核工作流页面
- **页面路径**: `/knowledge/audit`
- **页面布局**: 待审核列表 + 审核详情
- **核心组件**: 待审核知识列表（提交人/时间/类型）、知识内容预览、审核操作按钮
- **交互规范**: 通过/拒绝需填写理由；支持批量审核；审核历史可追溯
- **按钮**: [批量通过] [批量拒绝] 靠右

#### 知识批量导入页面
- **页面路径**: `/knowledge/import`
- **页面布局**: 上传向导
- **核心组件**: Excel 模板下载、上传进度条、导入结果预览（成功/失败列表）
- **交互规范**: 下载标准模板；失败行高亮并显示原因；支持断点续传
- **按钮**: [下载模板] 靠左, [开始导入] 靠右

### 12.20 设备监控（Sprint 9）

#### 监控大屏全屏页面
- **页面路径**: `/monitoring/large-screen`
- **页面布局**: 大盘全屏（无侧边栏）
- **核心组件**: 设备状态地图、实时指标图表、告警滚动列表
- **交互规范**: 全屏展示；定时刷新；支持 TV 模式（无鼠标操作）
- **按钮**: [退出全屏] 靠右上角

#### 监控告警规则配置页面
- **页面路径**: `/monitoring/alerts/rules`
- **页面布局**: 规则列表 + 表单抽屉
- **核心组件**: 规则列表（指标/阈值/告警级别）、规则启用状态、触发统计
- **交互规范**: 新建/编辑规则弹出 Drawer；删除前二次确认；支持规则复制
- **按钮**: [新建规则] 靠左, [批量启用/禁用] 靠右

### 12.21 设备配对（Sprint 9）

#### 配对码批量生成页面
- **页面路径**: `/devices/pairing/batch`
- **页面布局**: 表单 + 生成结果
- **核心组件**: 数量输入、设备模型选择、有效期设置、配对码列表
- **交互规范**: 批量生成支持 100+；结果支持导出 Excel/CSV；配对码可复制
- **按钮**: [批量生成] 靠左, [导出全部] 靠右

#### 配对统计看板
- **页面路径**: `/devices/pairing/stats`
- **页面布局**: 大盘式
- **核心组件**: 配对成功率趋势图、配对失败原因分布饼图、Top 设备型号排行
- **交互规范**: 时间范围筛选；点击图表钻取详情；失败原因分类统计
- **按钮**: [时间范围] 靠左, [导出报表] 靠右

### 12.22 OpenClaw 控制台（Sprint 6-7）

#### 控制台主题定制页面
- **页面路径**: `/console/themes`
- **页面布局**: 主题预览 + 配置面板
- **核心组件**: 主题预览卡片（深色/浅色/自定义）、配色方案选择、皮肤切换
- **交互规范**: 实时预览主题效果；支持自定义配色；主题切换即时生效
- **按钮**: [保存主题] 靠左, [恢复默认] 靠右

#### 多语言切换页面
- **页面路径**: `/console/language`
- **页面布局**: 语言列表 + 切换确认
- **核心组件**: 语言列表（中文/English/日本語等）、当前语言标识、翻译进度
- **交互规范**: 切换语言即时生效；未翻译内容显示原语言；支持自定义翻译
- **按钮**: [保存设置] 靠左

### 12.23 OpenClaw AI 版本管理（Sprint 7）

#### 自动回滚规则配置页面
- **页面路径**: `/openclaw/versions/:id/auto-rollback`
- **页面布局**: 表单 + 规则列表
- **核心组件**: 回滚触发条件配置（错误率/延迟阈值）、回滚版本选择、通知配置
- **交互规范**: 条件预览；支持规则测试；回滚执行需二次确认
- **按钮**: [保存规则] 靠左, [测试回滚] 靠右

### 12.24 权限管理（Sprint 4）

#### 权限申请表单页面
- **页面路径**: `/permissions/request`
- **页面布局**: 表单式
- **核心组件**: 权限类型选择、申请理由输入、有效期设置、审批流程预览
- **交互规范**: 权限类型级联选择；申请提交后进入审批流；进度可查询
- **按钮**: [提交申请] 靠左, [查看审批进度] 靠右

#### 权限审批工作流页面（管理员）
- **页面路径**: `/permissions/approve`
- **页面布局**: 待审批列表 + 审批详情
- **核心组件**: 待审批申请列表（申请人/权限类型/申请时间）、详情展开、审批操作
- **交互规范**: 通过/拒绝需填写意见；支持批量审批；审批历史完整记录
- **按钮**: [批量通过] [批量拒绝] 靠右

### 12.25 组织架构（Sprint 4）

#### 批量调岗操作页面
- **页面路径**: `/org/employees/batch-transfer`
- **页面布局**: 表单 + 预览
- **核心组件**: 员工选择（批量）、目标部门/岗位选择、调岗生效日期、影响预览
- **交互规范**: 批量选择员工；调岗前预览影响范围；调岗需二次确认
- **按钮**: [选择员工] 靠左, [执行调岗] 靠右

#### 组织变更时间线页面
- **页面路径**: `/org/change-logs`
- **页面布局**: 时间线式
- **核心组件**: 组织变更时间轴（部门合并/拆分/迁移）、变更详情、影响分析
- **交互规范**: 点击变更查看详情；支持按部门/时间筛选；变更可导出
- **按钮**: [筛选] 靠左, [导出记录] 靠右

---

## 十三、数据模型规划（补充）

> **说明：** 以下数据表为根据 PRD_GAP_ANALYSIS.md 差距分析补充的缺失数据模型，与现有数据模型共同构成完整设计。

### 13.1 设备管理

#### device_tags（设备标签关联表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| device_id | BIGINT FK | 设备 ID |
| tag_key | VARCHAR(64) | 标签键 |
| tag_value | VARCHAR(256) | 标签值 |
| created_at | TIMESTAMP | 创建时间 |
| created_by | BIGINT FK | 创建人 |

#### device_groups（设备分组表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| group_name | VARCHAR(128) | 分组名称 |
| parent_id | BIGINT FK | 父分组 ID（支持树形） |
| description | TEXT | 分组描述 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

#### device_usage_stats（设备使用统计表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| device_id | BIGINT FK | 设备 ID |
| usage_date | DATE | 统计日期 |
| total_minutes | INT | 当日使用时长（分钟） |
| session_count | INT | 当日会话次数 |
| avg_response_time_ms | INT | 平均响应时长（毫秒） |

### 13.2 设备影子

#### device_shadow_versions（影子版本历史表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| device_id | BIGINT FK | 设备 ID |
| version | INT | 版本号 |
| shadow_state | JSONB | 影子状态快照 |
| diff_from_previous | JSONB | 与上一版本差异 |
| created_at | TIMESTAMP | 创建时间 |
| created_by | VARCHAR(64) | 创建方式（auto/manual） |

#### device_state_predictions（设备状态预测表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| device_id | BIGINT FK | 设备 ID |
| predicted_at | TIMESTAMP | 预测时间点 |
| predicted_state | JSONB | 预测状态 |
| confidence | DECIMAL(5,2) | 置信度（0-100） |
| actual_state | JSONB | 实际状态（回填） |
| prediction_error | DECIMAL(5,2) | 预测误差 |

### 13.3 会员管理

#### member_balance_transactions（储值交易流水表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| member_id | BIGINT FK | 会员 ID |
| transaction_type | VARCHAR(16) | 交易类型（recharge/deduct/refund） |
| amount | DECIMAL(12,2) | 交易金额 |
| balance_before | DECIMAL(12,2) | 交易前余额 |
| balance_after | DECIMAL(12,2) | 交易后余额 |
| payment_method | VARCHAR(32) | 支付方式 |
| order_id | VARCHAR(64) | 关联订单号 |
| reason | VARCHAR(256) | 扣款原因（仅扣款） |
| created_at | TIMESTAMP | 创建时间 |

#### member_upgrade_triggers（升级触发记录表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| member_id | BIGINT FK | 会员 ID |
| trigger_type | VARCHAR(32) | 触发类型（manual/auto/time_threshold） |
| old_level | VARCHAR(32) | 原等级 |
| new_level | VARCHAR(32) | 新等级 |
| growth_value | INT | 成长值 |
| triggered_at | TIMESTAMP | 触发时间 |

### 13.4 通知与消息

#### notification_schedules（定时发送表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| template_id | BIGINT FK | 消息模板 ID |
| scheduled_at | TIMESTAMP | 计划发送时间 |
| recipient_type | VARCHAR(32) | 接收者类型（user/device/group） |
| recipient_ids | JSONB | 接收者 ID 列表 |
| status | VARCHAR(16) | 状态（pending/sent/failed/cancelled） |
| created_at | TIMESTAMP | 创建时间 |
| sent_at | TIMESTAMP | 实际发送时间 |

#### notification_blacklist（免打扰黑名单表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| user_id | BIGINT FK | 用户 ID |
| channel | VARCHAR(16) | 通知渠道（push/email/sms） |
| start_time | TIME | 免打扰开始时间 |
| end_time | TIME | 免打扰结束时间 |
| created_at | TIMESTAMP | 创建时间 |

### 13.5 应用管理

#### app_permissions（应用权限表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| app_id | BIGINT FK | 应用 ID |
| permission_key | VARCHAR(64) | 权限标识 |
| permission_name | VARCHAR(128) | 权限名称 |
| permission_type | VARCHAR(16) | 权限类型（sensitive/normal） |
| description | TEXT | 权限描述 |

#### app_user_permissions（用户应用授权表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| app_id | BIGINT FK | 应用 ID |
| user_id | BIGINT FK | 用户 ID |
| permission_id | BIGINT FK | 权限 ID |
| status | VARCHAR(16) | 状态（pending/approved/rejected） |
| expires_at | TIMESTAMP | 有效期（临时权限） |
| approved_at | TIMESTAMP | 审批时间 |
| approved_by | BIGINT FK | 审批人 |

### 13.6 内容管理

#### content_versions（内容版本表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| content_id | BIGINT FK | 内容 ID |
| version | INT | 版本号 |
| file_url | VARCHAR(512) | 文件 URL |
| file_hash | VARCHAR(64) | 文件哈希 |
| file_size | BIGINT | 文件大小 |
| change_summary | VARCHAR(256) | 变更说明 |
| created_by | BIGINT FK | 创建人 |
| created_at | TIMESTAMP | 创建时间 |

#### content_comments（内容评论表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| content_id | BIGINT FK | 内容 ID |
| user_id | BIGINT FK | 评论用户 |
| parent_id | BIGINT FK | 父评论 ID（嵌套） |
| content | TEXT | 评论内容 |
| like_count | INT | 点赞数 |
| status | VARCHAR(16) | 状态（active/deleted） |
| created_at | TIMESTAMP | 创建时间 |

### 13.7 数据分析

#### report_templates（报表模板表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| template_name | VARCHAR(128) | 模板名称 |
| dimensions | JSONB | 维度配置 |
| metrics | JSONB | 指标配置 |
| chart_types | JSONB | 图表类型 |
| created_by | BIGINT FK | 创建人 |
| is_public | BOOLEAN | 是否公开 |
| created_at | TIMESTAMP | 创建时间 |

#### report_subscriptions（报表订阅表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| user_id | BIGINT FK | 订阅用户 |
| report_template_id | BIGINT FK | 报表模板 ID |
| schedule_type | VARCHAR(16) | 调度类型（daily/weekly/monthly） |
| schedule_cron | VARCHAR(64) | Cron 表达式 |
| delivery_channels | JSONB | 投递渠道（email/push） |
| last_sent_at | TIMESTAMP | 上次发送时间 |
| created_at | TIMESTAMP | 创建时间 |

### 13.8 数字孪生

#### digital_twin_backups（数字孪生备份表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| pet_id | BIGINT FK | 宠物 ID |
| backup_type | VARCHAR(16) | 备份类型（full/incremental） |
| backup_size | BIGINT | 备份大小（字节） |
| encrypted_data | BYTEA | 加密备份数据 |
| status | VARCHAR(16) | 状态（creating/ready/failed） |
| created_at | TIMESTAMP | 创建时间 |

#### health_device_sync（第三方健康设备同步表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| pet_id | BIGINT FK | 宠物 ID |
| device_type | VARCHAR(64) | 设备类型（smart_collar/weighing_scale） |
| device_id | VARCHAR(128) | 设备序列号 |
| last_sync_at | TIMESTAMP | 最后同步时间 |
| sync_status | VARCHAR(16) | 同步状态 |
| sync_error | TEXT | 同步错误信息 |

### 13.9 仿真测试

#### simulation_datasets（仿真数据集表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| dataset_name | VARCHAR(128) | 数据集名称 |
| version | VARCHAR(32) | 版本号 |
| description | TEXT | 数据集描述 |
| file_count | INT | 文件数量 |
| total_size | BIGINT | 总大小（字节） |
| tags | JSONB | 标签 |
| created_by | BIGINT FK | 创建人 |
| created_at | TIMESTAMP | 创建时间 |

#### simulation_quotas（仿真配额表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| user_id | BIGINT FK | 用户 ID |
| quota_type | VARCHAR(32) | 配额类型（cpu/memory/gpu/time） |
| total_quota | BIGINT | 总额度 |
| used_quota | BIGINT | 已使用额度 |
| reset_at | TIMESTAMP | 重置时间 |
| updated_at | TIMESTAMP | 更新时间 |

### 13.10 宠物记忆

#### memory_recycle_bin（记忆回收站表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| memory_id | VARCHAR(64) | 原记忆 ID |
| device_id | BIGINT FK | 设备 ID |
| memory_content | JSONB | 记忆内容快照 |
| deleted_by | BIGINT FK | 删除人 |
| deleted_at | TIMESTAMP | 删除时间 |
| expires_at | TIMESTAMP | 过期时间（30 天后） |
| status | VARCHAR(16) | 状态（deleted/permanently_deleted/restored） |

### 13.11 主人画像

#### owner_privacy_settings（主人隐私设置表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| user_id | BIGINT FK | 用户 ID（唯一） |
| profile_visibility | VARCHAR(16) | 画像可见性（public/friends/private） |
| allow_data_export | BOOLEAN | 允许数据导出 |
| allow_third_party_share | BOOLEAN | 允许第三方共享 |
| updated_at | TIMESTAMP | 更新时间 |

### 13.12 知识库

#### knowledge_audit_logs（知识审核日志表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| knowledge_entry_id | BIGINT FK | 知识条目 ID |
| auditor_id | BIGINT FK | 审核人 |
| action | VARCHAR(16) | 审核动作（approve/reject） |
| comment | TEXT | 审核意见 |
| created_at | TIMESTAMP | 审核时间 |

#### knowledge_versions（知识版本表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| knowledge_entry_id | BIGINT FK | 知识条目 ID |
| version | INT | 版本号 |
| content_hash | VARCHAR(64) | 内容哈希 |
| change_summary | VARCHAR(256) | 变更说明 |
| created_by | BIGINT FK | 创建人 |
| created_at | TIMESTAMP | 创建时间 |

### 13.13 订阅管理

#### subscription_gifts（订阅赠送表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| gift_code | VARCHAR(64) | 赠送码（唯一） |
| sender_id | BIGINT FK | 发送者 ID |
| recipient_id | BIGINT FK | 接收者 ID |
| subscription_plan_id | BIGINT FK | 订阅计划 ID |
| gift_type | VARCHAR(16) | 赠送类型（trial/paid） |
| status | VARCHAR(16) | 状态（pending/claimed/expired） |
| expires_at | TIMESTAMP | 过期时间 |
| claimed_at | TIMESTAMP | 领取时间 |

#### family_plans（家庭计划表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| plan_name | VARCHAR(128) | 计划名称 |
| owner_id | BIGINT FK | 计划所有者 |
| max_members | INT | 最大成员数 |
| shared_quota | JSONB | 共享配额 |
| created_at | TIMESTAMP | 创建时间 |

### 13.14 权限管理

#### permission_requests（权限申请记录表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| requester_id | BIGINT FK | 申请人 |
| permission_id | BIGINT FK | 权限 ID |
| reason | TEXT | 申请理由 |
| status | VARCHAR(16) | 状态（pending/approved/rejected） |
| approver_id | BIGINT FK | 审批人 |
| approved_at | TIMESTAMP | 审批时间 |
| comment | TEXT | 审批意见 |
| created_at | TIMESTAMP | 创建时间 |

#### temporary_permissions（临时权限表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| user_id | BIGINT FK | 用户 ID |
| permission_id | BIGINT FK | 权限 ID |
| granted_by | BIGINT FK | 授权人 |
| expires_at | TIMESTAMP | 到期时间 |
| reason | VARCHAR(256) | 授权原因 |
| created_at | TIMESTAMP | 创建时间 |

### 13.15 开放平台生态

#### developer_levels（开发者等级表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| level_name | VARCHAR(32) | 等级名称（primary/intermediate/advanced） |
| api_quota_per_day | INT | 每日 API 配额 |
| rate_limit_per_minute | INT | 每分钟限流 |
| features | JSONB | 可用功能列表 |

#### platform_events（平台活动表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| event_name | VARCHAR(128) | 活动名称 |
| event_type | VARCHAR(32) | 活动类型（hackathon/webinar） |
| start_at | TIMESTAMP | 开始时间 |
| end_at | TIMESTAMP | 结束时间 |
| description | TEXT | 活动描述 |
| status | VARCHAR(16) | 状态 |
| created_at | TIMESTAMP | 创建时间 |

### 13.16 设备监控

#### monitoring_alert_rules（监控告警规则表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| rule_name | VARCHAR(128) | 规则名称 |
| metric_type | VARCHAR(64) | 指标类型 |
| threshold | DECIMAL(12,2) | 阈值 |
| comparison | VARCHAR(8) | 比较符（gt/lt/eq） |
| severity | VARCHAR(16) | 严重程度 |
| alert_channels | JSONB | 告警渠道 |
| enabled | BOOLEAN | 是否启用 |
| created_at | TIMESTAMP | 创建时间 |

### 13.17 组织架构

#### org_change_logs（组织变更日志表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| change_type | VARCHAR(32) | 变更类型（merge/split/transfer） |
| source_entity | VARCHAR(64) | 源实体（部门/岗位） |
| target_entity | VARCHAR(64) | 目标实体 |
| affected_employees | JSONB | 受影响员工列表 |
| change_summary | TEXT | 变更摘要 |
| executed_by | BIGINT FK | 执行人 |
| executed_at | TIMESTAMP | 执行时间 |

### 13.18 OpenClaw 控制台

#### console_themes（控制台主题表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| theme_name | VARCHAR(64) | 主题名称 |
| theme_config | JSONB | 主题配置（颜色/字体/布局） |
| is_system | BOOLEAN | 是否系统主题 |
| created_by | BIGINT FK | 创建人 |
| created_at | TIMESTAMP | 创建时间 |

#### console_feedback（用户反馈表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| user_id | BIGINT FK | 用户 ID |
| feedback_type | VARCHAR(32) | 反馈类型（bug/suggestion/compliment） |
| content | TEXT | 反馈内容 |
| rating | INT | 评分（1-5） |
| status | VARCHAR(16) | 处理状态 |
| created_at | TIMESTAMP | 创建时间 |

---

## 十四、修订记录

| 版本 | 日期 | 修订人 | 修订内容 |
|------|------|--------|----------|
| V1.0 | 2026-03-20 | agentcp | 初稿创建，模块现状分析和路线图规划 |
| V1.1 | 2026-03-20 | agentcp | 新增多租户Sprint规划 |
| V2.0 | 2026-03-22 | agentcp | 按6大领域重组，新增Phase 1-4完整规划（32个Sprint） |
| V2.1 | 2026-03-22 | agentcp | 76个功能点逐项列出，UI规范引用PRD文档，验收标准完善 |
| V2.2 | 2026-03-22 | agentcp | 补充缺失的9个MODULE PRD引用，完善功能清单与验收标准 |
| V2.3 | 2026-03-23 | agentcp | 补充15张缺失数据表、30个缺失API、30+UI页面规范；修正P0数量（18个）；补充OTA灰度/告警升级/储值等遗漏功能 |
