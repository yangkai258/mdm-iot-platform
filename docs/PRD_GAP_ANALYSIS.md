# PRD 差距分析报告

**分析日期：** 2026-03-23
**分析范围：** 所有 MODULE*.md 文档 vs MDM_PRODUCT_ROADMAP.md
**分析方法：** 逐模块逐字段对比功能清单、数据模型、API接口、UI说明与Sprint规划

---

## 一、总体概述

### 1.1 Roadmap vs PRD 模块对照

| Roadmap阶段 | Roadmap中的模块 | PRD模块数量 | 覆盖状态 |
|------------|-----------------|------------|---------|
| Phase 1 (Sprint 1-4) | 设备注册绑定/设备影子/OTA升级/Web管理后台/消息推送/告警系统/会员基础 | 7 | ⚠️ 部分覆盖 |
| Phase 2 (Sprint 5-16) | 设备影子增强/OTA增强/高级告警/会员增强/应用管理/内容管理/数据分析/AI系统工程/订阅管理 | 9 | ✅ 较完整 |
| Phase 3 (Sprint 17-24) | 情感计算/数字孪生/具身智能/仿真测试 | 4 | ✅ 完整 |
| Phase 4 (Sprint 25-26) | 开放平台生态 | 1 | ✅ 完整 |
| **额外模块** | 组织架构/系统管理/权限管理/设备监控/设备配对/宠物行为引擎/宠物记忆/主人画像/知识库/OpenClaw控制台/OpenClaw版本管理 | 11 | ➕ Roadmap未规划但已产出 |

### 1.2 关键差距统计

| 差距类型 | 数量 | 说明 |
|----------|------|------|
| **Roadmap中缺失的PRD功能** | 约45项 | 部分功能在Roadmap中未被提及 |
| **Roadmap中缺失的UI说明** | 约30项 | 新兴AI模块缺少前端入口描述 |
| **Roadmap中缺失的数据模型** | 约15张表 | AI增强模块的新增表未被规划 |
| **未规划但已实现的模块** | 11个 | 组织架构/OpenClaw系列等 |
| **Sprint规划模糊** | 多个模块 | 大部分模块未标注具体Sprint |

---

## 二、每个模块的差距

---

### MODULE_DEVICE_MANAGEMENT（设备管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 设备分组/标签管理（PRD中有功能描述，但roadmap未提及）<br>② 设备分享/授权功能<br>③ 设备健康评分机制<br>④ 设备使用统计（使用时长/频率） |
| **缺失数据模型** | ① `device_tags` 设备标签关联表（PRD功能列表提及但未建模）<br>② `device_groups` 设备分组表 |
| **缺失API** | ① `GET /api/v1/devices/:id/health-score` 设备健康评分<br>② `POST /api/v1/devices/:id/share` 设备分享<br>③ `GET /api/v1/devices/groups` 设备分组 |
| **缺失UI说明** | ① 设备健康评分卡片的展示位置<br>② 设备标签的筛选UI入口<br>③ 设备详情页中的"分享"按钮 |
| **当前 Sprint 状态** | **Sprint 1-2**（Phase 1核心模块，但PRDOta.md未标注Sprint编号） |

---

### MODULE_OTA_UPDATES（OTA升级）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① OTA升级的国际化和多语言固件支持<br>② 固件版本兼容性自动检测（与OpenClaw版本管理联动）<br>③ 灰度推送策略（百分比/白名单）<br>④ OTA升级预约功能 |
| **缺失数据模型** | ✅ 数据模型较完整 |
| **缺失API** | ① `POST /api/v1/ota/deployments/:id/pause` 暂停升级<br>② `POST /api/v1/ota/deployments/:id/gray` 灰度发布<br>③ `GET /api/v1/ota/devices/:device_id/compatibility` 固件兼容性查询 |
| **缺失UI说明** | ① 灰度发布配置UI<br>② OTA升级统计大盘中的"升级趋势图"<br>③ 固件版本对比UI |
| **当前 Sprint 状态** | **Sprint 3**（Phase 1，但PRDOta.md未标注Sprint编号） |

---

### MODULE_DEVICE_SHADOW（设备影子）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 设备影子版本历史记录<br>② 设备影子快照导出功能<br>③ 设备状态预测（基于历史数据的趋势预测）<br>④ 设备状态异常自动标记 |
| **缺失数据模型** | ① `device_shadow_versions` 影子版本历史表（未在PRD中建模）<br>② `device_state_predictions` 状态预测表 |
| **缺失API** | ① `GET /api/v1/shadows/:device_id/history` 影子版本历史<br>② `GET /api/v1/shadows/:device_id/prediction` 状态预测<br>③ `POST /api/v1/shadows/:device_id/export` 导出快照 |
| **缺失UI说明** | ① 设备影子详情页中的"历史版本"Tab<br>② 状态预测曲线的展示位置 |
| **当前 Sprint 状态** | **Sprint 2**（Phase 1，但PRD未标注Sprint编号） |

---

### MODULE_ALERT_SYSTEM（告警系统）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 告警抑制/去重规则（PRD流程图有提及但功能细节不足）<br>② 告警升级机制（持续未处理自动升级严重程度）<br>③ 告警自愈建议（根据告警类型推荐处理方案）<br>④ 告警与运维工单系统联动 |
| **缺失数据模型** | ✅ 数据模型完整（含 rules 和 alerts 两张核心表） |
| **缺失API** | ① `POST /api/v1/alerts/:id/escalate` 告警升级<br>② `GET /api/v1/alerts/suggestions/:id` 自愈建议<br>③ `GET /api/v1/alerts/rules/:id/trigger-history` 规则触发历史 |
| **缺失UI说明** | ① 告警升级规则的配置UI入口<br>② 自愈建议卡片的展示位置<br>③ 告警统计页中的"趋势图"和"分布图"的前端组件 |
| **当前 Sprint 状态** | **Sprint 4**（Phase 1，但PRD未标注Sprint编号） |

---

### MODULE_APP_MANAGEMENT（应用管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 应用市场/应用发现功能（仅内部托管，缺少公开市场）<br>② 应用权限管理（应用需要的权限列表和用户授权）<br>③ 应用使用时长统计<br>④ 应用自动更新策略 |
| **缺失数据模型** | ① `app_permissions` 应用权限表（PRD功能列表提及但未建模） |
| **缺失API** | ① `GET /api/v1/apps/:id/permissions` 应用所需权限<br>② `POST /api/v1/apps/:id/permissions/approve` 权限授权<br>③ `GET /api/v1/apps/:id/usage-stats` 使用时长统计 |
| **缺失UI说明** | ① 应用详情页中的"权限管理"Tab<br>② 应用统计页中的"使用时长"图表 |
| **当前 Sprint 状态** | **Sprint 10-11**（Phase 2，但PRD未标注Sprint编号） |

---

### MODULE_MEMBER_MANAGEMENT（会员管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 会员成长值/等级自动升级机制（PRD有升级规则但升级触发器未描述）<br>② 会员储值扣款流程（储值余额的扣减逻辑）<br>③ 会员注销/迁移流程<br>④ 会员画像标签（与主人画像库MODULE_OWNER_PROFILE的联动） |
| **缺失数据模型** | ① `member_upgrade_triggers` 升级触发记录表<br>② `member_balance_transactions` 储值交易流水表（充值/消费） |
| **缺失API** | ① `POST /api/v1/members/:id/balance/recharge` 储值充值<br>② `POST /api/v1/members/:id/balance/deduct` 储值扣款<br>③ `POST /api/v1/members/:id/upgrade/check` 触发等级检查<br>④ `DELETE /api/v1/members/:id` 会员注销 |
| **缺失UI说明** | ① 会员详情页中的"储值充值"按钮<br>② 等级自动升级的"升级动画"展示 |
| **当前 Sprint 状态** | **Sprint 7**（Phase 1会员基础，但PRD未标注Sprint编号） |

---

### MODULE_NOTIFICATION（通知与消息）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 通知草稿/定时发送管理<br>② 通知已读回执确认<br>③ 通知撤回功能（发送后取消）<br>④ 通知免打扰时间段配置 |
| **缺失数据模型** | ① `notification_schedules` 定时发送表<br>② `notification_blacklist` 免打扰黑名单表 |
| **缺失API** | ① `POST /api/v1/notifications/:id/recall` 撤回通知<br>② `PUT /api/v1/notifications/:id/read-receipt` 已读回执<br>③ `GET /api/v1/notifications/templates/:id/preview` 模板预览 |
| **缺失UI说明** | ① 通知发送页中的"定时发送"配置UI<br>② 已撤回通知的显示状态 |
| **当前 Sprint 状态** | **Sprint 4**（Phase 1推送通知，但PRD未标注Sprint编号） |

---

### MODULE_CONTENT_MANAGEMENT（内容与文档管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 内容版本管理（同一文件多次上传的版本控制）<br>② 内容水印添加功能<br>③ 内容评论/批注功能<br>④ 内容订阅/收藏功能 |
| **缺失数据模型** | ① `content_versions` 内容版本表<br>② `content_comments` 内容评论表<br>③ `content_subscriptions` 内容订阅表 |
| **缺失API** | ① `GET /api/v1/contents/:id/versions` 版本历史<br>② `POST /api/v1/contents/:id/comment` 添加评论<br>③ `POST /api/v1/contents/:id/subscribe` 订阅内容 |
| **缺失UI说明** | ① 内容详情页中的"版本历史"Tab<br>② 内容评论区的展示位置 |
| **当前 Sprint 状态** | **Sprint 12**（Phase 2，但PRD未标注Sprint编号） |

---

### MODULE_SYSTEM_MANAGEMENT（系统管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 系统配置变更审计（配置修改历史）<br>② 系统健康检查接口（/health）的可视化展示<br>③ 操作日志实时推送（WebSocket）<br>④ 用户在线状态管理（强制下线） |
| **缺失数据模型** | ✅ 数据模型完整（含 sys_operation_logs / sys_login_logs / device_event_logs / apns_config） |
| **缺失API** | ① `POST /api/v1/system/users/:id/force-logout` 强制下线<br>② `GET /api/v1/system/health` 健康检查<br>③ `GET /api/v1/system/config/history` 配置变更历史 |
| **缺失UI说明** | ① 系统健康状态仪表盘<br>② 强制下线按钮的位置 |
| **当前 Sprint 状态** | **Sprint 4**（Phase 1，但PRD未标注Sprint编号） |

---

### MODULE_DATA_ANALYSIS（数据分析）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 自定义报表生成器（用户可自定义维度和指标）<br>② 数据导出功能（Excel/PDF）<br>③ 报表订阅/定时推送<br>④ 异常检测自动告警（数据异常时自动推送） |
| **缺失数据模型** | ① `report_templates` 报表模板表<br>② `report_subscriptions` 报表订阅表 |
| **缺失API** | ① `POST /api/v1/stats/reports/custom` 生成自定义报表<br>② `GET /api/v1/stats/reports/export` 导出报表<br>③ `POST /api/v1/stats/reports/subscribe` 报表订阅 |
| **缺失UI说明** | ① 报表生成器的拖拽配置UI<br>② 报表导出按钮的位置 |
| **当前 Sprint 状态** | **Sprint 13**（Phase 2，但PRD未标注Sprint编号） |

---

### MODULE_ORGANIZATION（组织架构）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 组织架构变更记录（部门合并/拆分历史）<br>② 批量员工调岗/调部门<br>③ 任职资格管理<br>④ 员工通讯录导出 |
| **缺失数据模型** | ① `org_change_logs` 组织变更日志表 |
| **缺失API** | ① `POST /api/v1/org/employees/batch-transfer` 批量调岗<br>② `GET /api/v1/org/employees/export` 导出通讯录 |
| **缺失UI说明** | ① 批量调岗的操作入口<br>② 组织变更时间线视图 |
| **当前 Sprint 状态** | **未规划（额外模块）** |

---

### MODULE_DEVICE_MONITORING（设备监控面板）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 监控告警规则配置（基于监控指标的自定义告警）<br>② 监控数据导出功能<br>③ 多设备对比视图<br>④ 监控大屏全屏模式 |
| **缺失数据模型** | ① `monitoring_alert_rules` 监控告警规则表 |
| **缺失API** | ① `POST /api/v1/monitoring/alerts/rules` 配置监控告警规则<br>② `GET /api/v1/devices/:device_id/sensor-history/export` 导出传感器历史 |
| **缺失UI说明** | ① 监控大屏的全屏入口<br>② 告警规则配置UI |
| **当前 Sprint 状态** | **未规划（额外模块）** |

---

### MODULE_DEVICE_PAIRING（设备配对与注册）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 配对码批量生成<br>② 配对码过期自动清理<br>③ 配对记录统计分析<br>④ 配对失败原因分类统计 |
| **缺失数据模型** | ✅ 数据模型完整（pairing_records / device_openclaw_binding） |
| **缺失API** | ① `POST /api/v1/devices/pairing/codes/batch` 批量生成配对码<br>② `DELETE /api/v1/devices/pairing/codes/cleanup` 清理过期配对码<br>③ `GET /api/v1/devices/pairing/stats` 配对统计 |
| **缺失UI说明** | ① 批量生成配对码UI入口<br>② 配对统计看板 |
| **当前 Sprint 状态** | **未规划（额外模块）** |

---

### MODULE_PET_BEHAVIOR_ENGINE（宠物行为引擎）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 动作序列版本管理<br>② 动作执行优先级动态调整（基于宠物状态）<br>③ 动作异常自动上报和统计<br>④ 动作序列导入/导出 |
| **缺失数据模型** | ✅ 数据模型完整（action_library / action_sequences / decision_rules / sensor_events / action_executions） |
| **缺失API** | ① `POST /api/v1/behavior/sequences/:id/versions` 序列版本管理<br>② `GET /api/v1/behavior/actions/compatibility/:device_model` 设备兼容动作查询 |
| **缺失UI说明** | ① 动作序列版本历史对比UI<br>② 设备兼容动作的筛选展示 |
| **当前 Sprint 状态** | **未规划（额外模块，Phase 1-2 AI层核心）** |

---

### MODULE_PET_MEMORY（宠物记忆库）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 记忆加密存储（主人隐私保护）<br>② 记忆跨宠物迁移（换宠物时迁移记忆偏好）<br>③ 记忆质量评分机制<br>④ 记忆误删恢复（回收站） |
| **缺失数据模型** | ① `memory_recycle_bin` 记忆回收站表 |
| **缺失API** | ① `POST /api/v1/memory/:memory_id/restore` 恢复误删记忆<br>② `GET /api/v1/memory/:device_id/quality-report` 记忆质量报告 |
| **缺失UI说明** | ① 记忆回收站的展示入口<br>② 记忆质量评分的前端展示 |
| **当前 Sprint 状态** | **未规划（额外模块，Phase 2 AI层核心）** |

---

### MODULE_OWNER_PROFILE（主人画像库）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 画像数据导入/导出<br>② 多宠物主人管理（一个主人多个宠物时的偏好隔离）<br>③ 画像相似度匹配（用于社区功能）<br>④ 主人隐私设置（画像数据的可见性） |
| **缺失数据模型** | ① `owner_privacy_settings` 主人隐私设置表 |
| **缺失API** | ① `GET /api/v1/owner-profile/:user_id/export` 画像导出<br>② `PUT /api/v1/owner-profile/:user_id/privacy` 隐私设置 |
| **缺失UI说明** | ① 隐私设置页面的入口 |
| **当前 Sprint 状态** | **未规划（额外模块，Phase 2 AI层核心）** |

---

### MODULE_KNOWLEDGE_BASE（知识库）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 知识库版本管理（知识条目的变更历史）<br>② 知识库审核工作流（用户提交知识需审核）<br>③ 知识库批量导入/导出（Excel格式）<br>④ 知识库使用排行（高频知识TOP） |
| **缺失数据模型** | ① `knowledge_audit_logs` 知识审核日志表<br>② `knowledge_versions` 知识版本表 |
| **缺失API** | ① `GET /api/v1/knowledge/entries/:id/history` 知识版本历史<br>② `POST /api/v1/knowledge/entries/batch-import` 批量导入<br>③ `GET /api/v1/knowledge/top-queries` 高频查询排行 |
| **缺失UI说明** | ① 知识审核工作流UI<br>② 批量导入的Excel模板下载入口 |
| **当前 Sprint 状态** | **未规划（额外模块，Phase 2 AI层核心）** |

---

### MODULE_AI_ENGINEERING（AI系统工程）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 模型训练任务的分布式调度（Kubernetes集群集成）<br>② GPU资源池管理<br>③ 模型推理服务自动扩缩容<br>④ A/B实验与订阅系统联动（不同订阅等级使用不同模型） |
| **缺失数据模型** | ✅ 数据模型完整（含训练任务/数据集/A-B实验/模型版本/监控/沙箱/决策日志/路由策略/分片共9张表） |
| **缺失API** | ① `GET /api/v1/ai/gpu/status` GPU资源状态<br>② `POST /api/v1/ai/scaling/trigger` 触发自动扩缩容 |
| **缺失UI说明** | ① GPU资源监控大盘<br>② 模型推理服务健康状态 |
| **当前 Sprint 状态** | **Sprint 5-6（Phase 2，PRD标注Phase 1基础/Phase 2完善）** |

---

### MODULE_EMBODIED_AI（具身智能）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 具身AI模型的端云协同推理（边缘端+云端混合）<br>② 具身AI模型压缩和优化工具链<br>③ 动作模仿的学习进度可视化<br>④ 具身AI安全审计日志 |
| **缺失数据模型** | ✅ 数据模型完整（含 embodied_maps / spatial_positions / action_library / action_executions / safety_zones / embodied_decision_logs / safety_logs） |
| **缺失API** | ① `POST /api/v1/embodied/:device_id/model/compress` 模型压缩<br>② `GET /api/v1/embodied/:device_id/learning/progress` 学习进度 |
| **缺失UI说明** | ① 模型压缩配置UI<br>② 学习进度的进度条展示 |
| **当前 Sprint 状态** | **Sprint 21-22（Phase 3，PRD明确标注）** |

---

### MODULE_AFFECTIVE_COMPUTING（情感计算）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 情感识别模型的训练和更新（支持持续学习）<br>② 情感预警机制（用户情绪持续低落时通知家属）<br>③ 情感数据的脱敏和匿名化处理（隐私合规）<br>④ 宠物情感和主人情感的互动影响模型 |
| **缺失数据模型** | ✅ 数据模型完整（含 emotion_records / pet_emotion_actions / emotion_response_configs / family_emotions / emotion_reports） |
| **缺失API** | ① `POST /api/v1/emotion/model/retrain` 触发情感模型重训练<br>② `POST /api/v1/emotion/:pet_id/alert/notify` 情感预警通知 |
| **缺失UI说明** | ① 情感预警规则的配置UI<br>② 家庭情绪预警的推送展示 |
| **当前 Sprint 状态** | **Sprint 17-18（Phase 3，PRD明确标注）** |

---

### MODULE_DIGITAL_TWIN（数字孪生）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 数字孪生3D渲染引擎的前端集成（PRD标注"不包含前端3D渲染引擎"但未规划接口）<br>② 宠物克隆/备份功能（数字孪生数据备份）<br>③ 第三方健康设备数据接入（智能项圈等）<br>④ 数字孪生与宠物医疗记录联动 |
| **缺失数据模型** | ① `health_device_sync` 第三方健康设备同步表<br>② `digital_twin_backups` 数字孪生备份表 |
| **缺失API** | ① `POST /api/v1/digital-twin/:pet_id/backup` 创建备份<br>② `POST /api/v1/digital-twin/:pet_id/restore/:backup_id` 恢复备份<br>③ `POST /api/v1/digital-twin/:pet_id/health-devices/sync` 第三方设备数据同步 |
| **缺失UI说明** | ① 数字孪生备份/恢复入口<br>② 第三方健康设备绑定UI |
| **当前 Sprint 状态** | **Sprint 19-20（Phase 3，PRD明确标注）** |

---

### MODULE_SUBSCRIPTION（订阅管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 订阅赠送功能（赠送好友体验）<br>② 家庭计划（多用户共享订阅）<br>③ 订阅试用期延长<br>④ 订阅数据分析（付费转化率/流失率） |
| **缺失数据模型** | ① `subscription_gifts` 订阅赠送表<br>② `family_plans` 家庭计划表 |
| **缺失API** | ① `POST /api/v1/subscriptions/gift` 赠送订阅<br>② `POST /api/v1/subscriptions/family/join` 加入家庭计划<br>③ `GET /api/v1/subscriptions/analytics` 订阅数据分析 |
| **缺失UI说明** | ① 订阅赠送的分享入口<br>② 家庭计划的管理UI |
| **当前 Sprint 状态** | **Sprint 15-16（Phase 2，PRD明确标注）** |

---

### MODULE_SIMULATION（仿真与测试）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 仿真环境与CI/CD流水线集成<br>② 仿真结果自动生成测试报告并推送至Slack/钉钉<br>③ 仿真数据集管理（测试数据的版本化）<br>④ 仿真资源的配额管理和计费 |
| **缺失数据模型** | ① `simulation_datasets` 仿真数据集表<br>② `simulation_quotas` 仿真配额表 |
| **缺失API** | ① `POST /api/v1/simulation/cicd/integrate` CI/CD集成配置<br>② `GET /api/v1/simulation/datasets` 仿真数据集管理 |
| **缺失UI说明** | ① CI/CD集成配置UI<br>② 仿真数据集版本管理界面 |
| **当前 Sprint 状态** | **Sprint 23-24（Phase 3，PRD明确标注）** |

---

### MODULE_OPENCLAW_CONSOLE（OpenClaw宠物控制台）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 控制台皮肤/主题定制<br>② 控制台多语言支持（国际化）<br>③ 控制台辅助功能支持（无障碍访问）<br>④ 控制台意见反馈和评分 |
| **缺失数据模型** | ① `console_themes` 控制台主题表<br>② `console_feedback` 用户反馈表 |
| **缺失API** | ① `GET /api/v1/console/themes` 主题列表<br>② `POST /api/v1/console/feedback` 提交反馈 |
| **缺失UI说明** | ① 控制台主题切换入口<br>② 多语言切换按钮的位置 |
| **当前 Sprint 状态** | **未规划（额外模块，Phase 1-2前端核心入口）** |

---

### MODULE_OPENCLAW_VERSION（OpenClaw AI版本管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① AI版本灰度发布自动化（与MODULE_AI_ENGINEERING联动）<br>② AI版本回滚的自动化触发规则<br>③ 设备端AI能力发现协议（设备上报支持的AI能力） |
| **缺失数据模型** | ✅ 数据模型完整（openclaw_versions / version_compatibility / device_openclaw_version） |
| **缺失API** | ① `POST /api/v1/openclaw/versions/:id/auto-rollback-rule` 配置自动回滚规则 |
| **缺失UI说明** | ① 自动回滚规则的配置UI |
| **当前 Sprint 状态** | **未规划（额外模块，与Sprint 5-6 AI系统工程强相关）** |

---

### MODULE_PERMISSIONS（权限管理）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 权限的有效期管理（临时权限，到期自动回收）<br>② 权限申请审批工作流（用户申请更高权限）<br>③ 权限数据导出（Excel格式）<br>④ 跨租户权限隔离（未来多租户场景） |
| **缺失数据模型** | ① `permission_requests` 权限申请记录表<br>② `temporary_permissions` 临时权限表 |
| **缺失API** | ① `POST /api/v1/permissions/request` 申请权限<br>② `POST /api/v1/permissions/approve/:request_id` 审批权限 |
| **缺失UI说明** | ① 权限申请表单UI<br>② 权限审批工作流UI（管理员侧） |
| **当前 Sprint 状态** | **未规划（额外模块，属于Phase 1基础安全能力）** |

---

### MODULE_PLATFORM_ECOSYSTEM（开放平台生态）

| 差距维度 | 详情 |
|----------|------|
| **缺失功能** | ① 开发者等级体系（初级/中级/高级开发者，不同API配额）<br>② API调用退款机制<br>③ 开放平台社区论坛<br>④ 开发者活动/黑客马拉松管理 |
| **缺失数据模型** | ① `developer_levels` 开发者等级表<br>② `platform_events` 平台活动表 |
| **缺失API** | ① `GET /api/v1/developer/level` 查询开发者等级<br>② `POST /api/v1/developer/refund` API调用退款 |
| **缺失UI说明** | ① 开发者等级展示和升级入口<br>② 平台活动页面 |
| **当前 Sprint 状态** | **Sprint 25-26（Phase 4，PRD明确标注）** |

---

## 三、共性问题汇总

### 3.1 Sprint 规划未标注

**几乎所有模块的PRD都未标注具体的Sprint编号**，仅 `MODULE_AI_ENGINEERING`、`MODULE_EMBODIED_AI`、`MODULE_AFFECTIVE_COMPUTING`、`MODULE_DIGITAL_TWIN`、`MODULE_SUBSCRIPTION`、`MODULE_SIMULATION` 在文档中标注了Phase和Sprint。

**影响：** 无法判断开发优先级和排期合理性。

### 3.2 AI增强模块缺少前端UI入口

MODULE_PET_BEHAVIOR_ENGINE、MODULE_PET_MEMORY、MODULE_OWNER_PROFILE、MODULE_KNOWLEDGE_BASE 等AI层模块大量缺少：
- 前端入口路由
- 前端页面组件设计
- 前端按钮规范
- 响应式适配说明

### 3.3 额外模块未进入Roadmap

以下已产出PRD的模块未在 Roadmap 中体现：
- MODULE_ORGANIZATION（组织架构）
- MODULE_SYSTEM_MANAGEMENT（系统管理）
- MODULE_PERMISSIONS（权限管理）
- MODULE_DEVICE_MONITORING（设备监控面板）
- MODULE_DEVICE_PAIRING（设备配对与注册）
- MODULE_PET_BEHAVIOR_ENGINE（宠物行为引擎）
- MODULE_PET_MEMORY（宠物记忆库）
- MODULE_OWNER_PROFILE（主人画像库）
- MODULE_KNOWLEDGE_BASE（知识库）
- MODULE_OPENCLAW_CONSOLE（OpenClaw控制台）
- MODULE_OPENCLAW_VERSION（OpenClaw AI版本管理）

### 3.4 跨模块联动设计缺失

多个模块PRD中描述了跨模块联动（如行为引擎→知识库→主人画像），但未制定：
- 联动接口的版本管理
- 联动失败降级策略
- 联动数据的一致性保证机制

### 3.5 测试模块严重缺失

除了仿真测试(MODULE_SIMULATION)外，没有独立的：
- 单元测试规范
- 集成测试策略
- 性能测试指标
- 安全测试要求

---

## 四、优先级建议

### P0（必须修复，Sprint 1-4内）

1. **为所有模块PRD补充具体Sprint编号**（参考 Roadmap 中标注了Sprint的模块）
2. **为 MODULE_DEVICE_MANAGEMENT 补充**设备分组/标签功能的API和数据模型
3. **为 MODULE_MEMBER_MANAGEMENT 补充**储值充值/扣款API和储值交易流水表
4. **Roadmap 补充**以下未规划但已产出的模块：
   - MODULE_ORGANIZATION → Sprint 4
   - MODULE_PERMISSIONS → Sprint 4
   - MODULE_DEVICE_MONITORING → Sprint 6
   - MODULE_OPENCLAW_CONSOLE → Sprint 6-7
   - MODULE_OPENCLAW_VERSION → Sprint 7
   - 行为引擎/记忆/画像/知识库 → Sprint 8-9

### P1（Phase 2内完善）

5. 为 AI 层模块（行为引擎/记忆/画像/知识库）补充完整的前端UI页面设计
6. 补充跨模块联动接口的版本管理和降级策略
7. 补充 MODULE_SUBSCRIPTION 的订阅赠送和家庭计划功能

### P2（Phase 3-4完善）

8. 补充所有模块的数据导出功能
9. 补充权限管理模块的审批工作流
10. 补充开放平台生态的开发者等级和活动管理

---

## 五、附录：模块PRD清单

| 文件名 | 版本 | Sprint归属（推测） |
|--------|------|-------------------|
| MODULE_DEVICE_MANAGEMENT.md | V1.4 | Sprint 1-2 |
| MODULE_OTA_UPDATES.md | V1.4 | Sprint 3 |
| MODULE_DEVICE_SHADOW.md | V1.2 | Sprint 2 |
| MODULE_POLICY_MANAGEMENT.md | V1.3 | Sprint 4 |
| MODULE_ALERT_SYSTEM.md | V1.5 | Sprint 4 |
| MODULE_APP_MANAGEMENT.md | V1.3 | Sprint 10-11 |
| MODULE_MEMBER_MANAGEMENT.md | V1.4 | Sprint 7 |
| MODULE_NOTIFICATION.md | V1.3 | Sprint 4 |
| MODULE_CONTENT_MANAGEMENT.md | V1.3 | Sprint 12 |
| MODULE_SYSTEM_MANAGEMENT.md | V1.4 | Sprint 4 |
| MODULE_DATA_ANALYSIS.md | V1.4 | Sprint 13 |
| MODULE_ORGANIZATION.md | V1.4 | ⚠️ 未规划 |
| MODULE_DEVICE_MONITORING.md | V1.0 | ⚠️ 未规划 |
| MODULE_DEVICE_PAIRING.md | V1.0 | ⚠️ 未规划 |
| MODULE_PET_BEHAVIOR_ENGINE.md | V1.1 | ⚠️ 未规划 |
| MODULE_PET_MEMORY.md | V1.1 | ⚠️ 未规划 |
| MODULE_OWNER_PROFILE.md | V1.0 | ⚠️ 未规划 |
| MODULE_KNOWLEDGE_BASE.md | V1.0 | ⚠️ 未规划 |
| MODULE_AI_ENGINEERING.md | V1.0 | Sprint 5-6 |
| MODULE_EMBODIED_AI.md | V1.0 | Sprint 21-22 |
| MODULE_AFFECTIVE_COMPUTING.md | V1.0 | Sprint 17-18 |
| MODULE_DIGITAL_TWIN.md | V1.0 | Sprint 19-20 |
| MODULE_SUBSCRIPTION.md | V1.0 | Sprint 15-16 |
| MODULE_SIMULATION.md | V1.0 | Sprint 23-24 |
| MODULE_OPENCLAW_CONSOLE.md | V1.1 | ⚠️ 未规划 |
| MODULE_OPENCLAW_VERSION.md | V1.0 | ⚠️ 未规划 |
| MODULE_PERMISSIONS.md | V1.0 | ⚠️ 未规划 |
| MODULE_PLATFORM_ECOSYSTEM.md | V