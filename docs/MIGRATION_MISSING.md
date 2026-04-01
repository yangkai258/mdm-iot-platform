# 旧前端 → 新前端 缺失文件清单

生成时间: 2026-04-02

## 🔴 核心页面（必须迁移）

| 文件 | 说明 | 优先级 |
|------|------|--------|
| `Dashboard.vue` | 旧版主仪表盘 | P0 |
| `DeviceDashboard.vue` | 设备仪表盘 | P0 |
| `DeviceStatus.vue` | 设备状态页 | P0 |
| `Member.vue` | 会员总览页 | P0 |
| `PetConfig.vue` | 宠物配置页 | P0 |
| `OtaFirmware.vue` | OTA固件管理 | P0 |
| `Alert.vue` | 告警总览页 | P0 |
| `Login.vue` | 旧登录页（参考） | P1 |
| `NotFound.vue` | 404页面 | P1 |
| `ModalTest.vue` | 测试页（可删） | P2 |

## 🟡 高级功能模块

### `advanced/` - 高级配置
| 文件 | 说明 |
|------|------|
| `ChildModeView.vue` | 儿童模式 |
| `ElderlyModeView.vue` | 老人模式 |
| `FamilyAlbumView.vue` | 家庭相册 |
| `FeatureConfigView.vue` | 功能配置 |
| `DeviceCertificates.vue` | 设备证书 |
| `DietRecordView.vue` | 饮食记录 |
| `PetFinderView.vue` | 宠物寻找 |
| `VaccinationView.vue` | 疫苗接种 |

### `ai/` - AI行为引擎
| 文件 | 说明 |
|------|------|
| `AIBehaviorView.vue` | AI行为总览 |
| `AIBehaviorDetailView.vue` | AI行为详情 |
| `AIBehaviorLogView.vue` | AI行为日志 |
| `AIDecisionLogsView.vue` | AI决策日志 |
| `AIEmotionView.vue` | AI情绪 |
| `AIQualityDashboardView.vue` | AI质量仪表盘 |
| `AISandboxView.vue` | AI沙盒 |
| `ModelVersionView.vue` | 模型版本 |
| `ModelPublishWorkflow.vue` | 模型发布流程 |
| `ModelRollbackView.vue` | 模型回滚 |

### `alert/` - 告警（子目录）
| 文件 | 说明 |
|------|------|
| `AlertHistoryView.vue` | 告警历史 |
| `AlertNotificationView.vue` | 告警通知 |
| `EmailChannelConfig.vue` | 邮件通道配置 |
| `SMSChannelConfig.vue` | 短信通道配置 |
| `WebhookChannelConfig.vue` | Webhook通道配置 |
| `NotificationLogsView.vue` | 通知日志 |
| `NotificationStatsView.vue` | 通知统计 |

### `analytics/` - 数据分析
| 文件 | 说明 |
|------|------|
| `DashboardView.vue` | 分析仪表盘 |
| `CohortAnalysisView.vue` | 队列分析 |
| `EventAnalyticsView.vue` | 事件分析 |
| `FunnelAnalysisView.vue` | 漏斗分析 |
| `RetentionAnalysisView.vue` | 留存分析 |

### `digital-twin/` - 数字孪生
| 文件 | 说明 |
|------|------|
| `VitalsDashboardView.vue` | 生命体征仪表盘 |
| `RealTimeVitalsChart.vue` | 实时生命体征 |
| `MomentsView.vue` | 精彩瞬间 |
| `BehaviorPredictionView.vue` | 行为预测 |
| `HistoricalReplayView.vue` | 历史回放 |

### `emotion/` - 情感计算
| 文件 | 说明 |
|------|------|
| `EmotionLogView.vue` | 情绪日志 |
| `EmotionRecognizeView.vue` | 情绪识别 |
| `EmotionReportView.vue` | 情绪报告 |
| `EmotionTrendView.vue` | 情绪趋势 |
| `EmotionResponseConfigView.vue` | 情绪响应配置 |
| `FamilyEmotionMapView.vue` | 家庭情绪地图 |
| `VoiceEmotionView.vue` | 语音情绪 |

### `health/` - 健康
| 文件 | 说明 |
|------|------|
| `HealthReportView.vue` | 健康报告 |
| `HealthWarningView.vue` | 健康预警 |
| `SleepAnalysisView.vue` | 睡眠分析 |
| `ExerciseStatsView.vue` | 运动统计 |

### `integration/` - 第三方集成
| 文件 | 说明 |
|------|------|
| `SmartHomeView.vue` | 智能家居 |
| `PetHospitalView.vue` | 宠物医院 |
| `PetShopView.vue` | 宠物商店 |

### `platform-evo/` - 平台增强
| 文件 | 说明 |
|------|------|
| `EdgeAIView.vue` | 边缘AI |
| `ModelShardView.vue` | 模型分片 |
| `RTOSView.vue` | RTOS设备 |
| `BLEMeshView.vue` | BLE Mesh |

### `simulation/` - 仿真测试
| 文件 | 说明 |
|------|------|
| `VirtualPetSimulationView.vue` | 虚拟宠物仿真 |
| `SimulationScenesView.vue` | 仿真场景 |
| `SimulationRecordingsView.vue` | 仿真记录 |
| `ABExperimentView.vue` | AB实验 |
| `ReplayView.vue` | 回放 |
| `StressTestView.vue` | 压力测试 |
| `TestFrameworkView.vue` | 测试框架 |

### `policies/` - 合规策略
| 文件 | 说明 |
|------|------|
| `PolicyList.vue` | 策略列表 |
| `PolicyTemplatesView.vue` | 策略模板 |
| `PolicyDistributionView.vue` | 策略分发 |
| `PolicyAuditView.vue` | 策略审计 |
| `PolicyConfigs.vue` | 策略配置 |
| `ComplianceRules.vue` | 合规模型 |
| `DeviceCompliance.vue` | 设备合规 |

### `security/` - 安全
| 文件 | 说明 |
|------|------|
| `AuditLogView.vue` | 审计日志 |
| `CertificateManageView.vue` | 证书管理 |
| `DataPermissionView.vue` | 数据权限 |
| `DataPrivacyView.vue` | 数据隐私 |
| `DeviceSecurityView.vue` | 设备安全 |
| `LDAPConfigView.vue` | LDAP配置 |
| `PermissionAssignmentView.vue` | 权限分配 |
| `SecuritySettingsView.vue` | 安全设置 |
| `UserSyncView.vue` | 用户同步 |

### `globalization/` - 全球化
| 文件 | 说明 |
|------|------|
| `LanguagePackView.vue` | 语言包 |
| `CurrencySettingsView.vue` | 货币设置 |
| `TimezoneSettingsView.vue` | 时区设置 |
| `DataResidencyView.vue` | 数据驻留 |
| `ContentDeliveryView.vue` | 内容分发 |
| `RegionalAINodeView.vue` | 区域AI节点 |
| `RegionManageView.vue` | 区域管理 |
| `RegionNodeView.vue` | 区域节点 |
| `RegionSyncStatusView.vue` | 区域同步状态 |
| `GlobalizationSettingsView.vue` | 全球化设置 |

### `research/` - 研究平台
| 文件 | 说明 |
|------|------|
| `DatasetLibraryView.vue` | 数据集库 |
| `DatasetDetailView.vue` | 数据集详情 |
| `DatasetView.vue` | 数据集 |
| `ExperimentView.vue` | 实验 |
| `ExperimentRunView.vue` | 实验运行 |
| `ResearchProjectView.vue` | 研究项目 |

## 🟢 业务功能

### `app/` - 应用
| 文件 | 说明 |
|------|------|
| `AppDeviceListView.vue` | App设备列表 |
| `AppDeviceControlView.vue` | App设备控制 |
| `MiniAppHomeView.vue` | 小程序首页 |
| `MiniAppDeviceView.vue` | 小程序设备 |

### `apps/` - 应用分发
| 文件 | 说明 |
|------|------|
| `AppList.vue` | 应用列表 |
| `AppVersions.vue` | 应用版本 |
| `AppDistributions.vue` | 应用分发 |

### `content/` - 内容管理
| 文件 | 说明 |
|------|------|
| `AppStoreView.vue` | 应用商店 |
| `ContentLibraryView.vue` | 内容库 |

### `knowledge/` - 知识库
| 文件 | 说明 |
|------|------|
| `KnowledgeList.vue` | 知识列表 |

### `member/` (子目录) - 会员详细功能
| 文件 | 说明 |
|------|------|
| `member/MemberSettings.vue` | 会员设置 |
| `member/MemberBenefits.vue` | 会员权益 |
| `member/MemberOrders.vue` | 会员订单 |
| `member/MemberCardTypes.vue` | 卡类型 |
| `member/MemberCardGroups.vue` | 卡组 |
| `member/MemberReception.vue` | 会员接待 |
| `member/MemberUpgradeRules.vue` | 升级规则 |
| `member/HighFreqTags.vue` | 高频标签 |
| `member/LowFreqTags.vue` | 低频标签 |
| `member/InterestTags.vue` | 兴趣标签 |
| `member/MemberArticles.vue` | 会员文章 |
| `member/MemberGifts.vue` | 会员礼品 |
| `member/SmsChannels.vue` | 短信通道 |
| `member/SmsTemplates.vue` | 短信模板 |
| `member/WechatSettings.vue` | 微信设置 |
| `member/OccupationTypes.vue` | 职业类型 |
| `member/PromotionTypes.vue` | 促销类型 |
| `member/Redpackets.vue` | 红包 |
| `member/AmountDiscountPromo.vue` | 满减促销 |
| `member/AmountReducePromo.vue` | 满折促销 |
| `member/BuyGiftPromo.vue` | 买赠促销 |
| `member/DirectReducePromo.vue` | 直减促销 |
| `member/VipExclusivePromo.vue` | VIP专属促销 |
| `member/CouponInventory.vue` | 优惠券库存 |
| `member/CouponMessages.vue` | 优惠券消息 |
| `member/GiftRecords.vue` | 礼品记录 |
| `member/TagAutoClean.vue` | 标签自动清理 |
| `member/TagReport.vue` | 标签报告 |

### `org/` - 组织架构
| 文件 | 说明 |
|------|------|
| `Companies.vue` | 公司列表 |
| `CompanyList.vue` | 公司管理 |
| `Departments.vue` | 部门列表 |
| `DepartmentList.vue` | 部门管理 |
| `Employees.vue` | 员工列表 |
| `EmployeeList.vue` | 员工管理 |
| `Posts.vue` | 岗位列表 |
| `PostList.vue` | 岗位管理 |
| `PositionList.vue` | 职位列表 |
| `StandardPositions.vue` | 标准职位 |

### `ota/` - OTA
| 文件 | 说明 |
|------|------|
| `OtaPackages.vue` | OTA包管理 |
| `OtaDeployments.vue` | OTA部署 |

### `pet-social/` - 宠物社交
| 文件 | 说明 |
|------|------|
| `PetFeedView.vue` | 宠物动态 |
| `PetFollowView.vue` | 关注 |
| `PetPlaydateView.vue` | 相约 |
| `MomentsView.vue` | 瞬间 |
| `LikeRecordView.vue` | 点赞记录 |
| `FollowView.vue` | 关注记录 |

### `platform/` - 开发者平台
| 文件 | 说明 |
|------|------|
| `DeveloperApi.vue` | API文档 |
| `DeveloperConsoleView.vue` | 开发者控制台 |
| `AgentView.vue` | Agent管理 |
| `AiModelVersionView.vue` | AI模型版本 |
| `WebhookMarketView.vue` | Webhook市场 |

### `security-evo/` - 安全增强
| 文件 | 说明 |
|------|------|
| `AuditLogView.vue` | 审计日志 |
| `ComplianceView.vue` | 合规 |
| `DataExportView.vue` | 数据导出 |
| `GDPRView.vue` | GDPR |

### `store/` - 店铺
| 文件 | 说明 |
|------|------|
| `StoreLocationsView.vue` | 门店位置 |
| `StoreSourcesView.vue` | 门店来源 |
| `PromotionTypesView.vue` | 促销类型 |

### `system/` - 系统设置
| 文件 | 说明 |
|------|------|
| `SystemSettings.vue` | 系统设置 |
| `Logs.vue` | 日志 |
| `Monitor.vue` | 监控 |
| `EmailTemplates.vue` | 邮件模板 |

### `tenants/` - 租户
| 文件 | 说明 |
|------|------|
| `TenantManagement.vue` | 租户管理 |
| `TenantApproval.vue` | 租户审批 |
| `TenantSettings.vue` | 租户设置 |
| `PublicArchives.vue` | 公共档案 |
| `SystemInfo.vue` | 系统信息 |

### `webhooks/` - Webhook
| 文件 | 说明 |
|------|------|
| `WebhookList.vue` | Webhook列表 |
| `WebhookLogs.vue` | Webhook日志 |

### `notifications/` - 通知
| 文件 | 说明 |
|------|------|
| `Announcements.vue` | 公告 |
| `NotificationList.vue` | 通知列表 |
| `NotificationTemplates.vue` | 通知模板 |

### `billing/` - 账单
| 文件 | 说明 |
|------|------|
| `BillingList.vue` | 账单列表 |

### `owner/` - 主人
| 文件 | 说明 |
|------|------|
| `OwnerProfile.vue` | 主人档案 |

### `pet/` - 宠物（旧版）
| 文件 | 说明 |
|------|------|
| `PetConsole.vue` | 宠物控制台 |
| `PetConsoleView.vue` | 宠物控制台视图 |
| `PetConversations.vue` | 宠物对话 |
| `PetFeedView.vue` | 宠物动态 |
| `PetFollowView.vue` | 宠物关注 |
| `PetPlaydateView.vue` | 宠物相约 |

### `miniclaw/` - MiniClaw
| 文件 | 说明 |
|------|------|
| `FirmwareList.vue` | 固件列表 |

### `portal/` - 门户
| 文件 | 说明 |
|------|------|
| `Dashboard.vue` | 门户仪表盘 |
| `Workbench.vue` | 工作台 |
| `WorkbenchPortal.vue` | 工作台门户 |
| `PersonalDesktop.vue` | 个人桌面 |

## 迁移建议

1. **优先核心页面** - Dashboard, Device, Member, Pet, Alert, OTA
2. **按模块批量迁移** - 每个模块一起迁移，保持一致性
3. **API联调** - 每个页面迁移后立即对接后端API
4. **旧前端保留** - 迁移完成前不删除旧前端

## 统计
- 旧前端总数: ~270 个 .vue 文件
- 新前端已有: ~260 个 .vue 文件
- 缺失: ~270 个文件
