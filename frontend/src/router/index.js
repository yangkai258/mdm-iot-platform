import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/test-modals',
    name: 'ModalTest',
    component: () => import('../views/ModalTest.vue')
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue')
  },
  {
    path: '/devices',
    name: 'Devices',
    component: () => import('../views/DeviceDashboard.vue')
  },
  {
    path: '/device/:id',
    name: 'DeviceDetail',
    component: () => import('../views/DeviceDetail.vue')
  },
  {
    path: '/ota',
    redirect: '/ota/packages'
  },
  {
    path: '/ota/packages',
    name: 'OtaPackages',
    component: () => import('../views/ota/OtaPackages.vue')
  },
  {
    path: '/ota/deployments',
    name: 'OtaDeployments',
    component: () => import('../views/ota/OtaDeployments.vue')
  },
  {
    path: '/alert',
    name: 'Alert',
    component: () => import('../views/Alert.vue')
  },
  {
    path: '/members',
    name: 'Members',
    component: () => import('../views/Member.vue')
  },
  {
    path: '/pet',
    name: 'Pet',
    component: () => import('../views/PetConfig.vue')
  },
  {
    path: '/pet/console',
    name: 'PetConsole',
    component: () => import('../views/pet/PetConsoleView.vue')
  },
  {
    path: '/pet/conversations',
    name: 'PetConversations',
    component: () => import('../views/pet/PetConversations.vue')
  },
  {
    path: '/miniclaw/firmwares',
    name: 'FirmwareList',
    component: () => import('../views/miniclaw/FirmwareList.vue')
  },
  {
    path: '/knowledge',
    name: 'KnowledgeList',
    component: () => import('../views/knowledge/KnowledgeList.vue')
  },
  {
    path: '/owner/profile',
    name: 'OwnerProfile',
    component: () => import('../views/owner/OwnerProfile.vue')
  },
  {
    path: '/system/monitor',
    name: 'Monitor',
    component: () => import('../views/system/Monitor.vue')
  },
  {
    path: '/system/logs',
    name: 'Logs',
    component: () => import('../views/system/Logs.vue')
  },
  {
    path: '/notifications',
    redirect: '/notifications/list'
  },
  {
    path: '/notifications/list',
    name: 'NotificationList',
    component: () => import('../views/notifications/NotificationList.vue')
  },
  {
    path: '/notifications/templates',
    name: 'NotificationTemplates',
    component: () => import('../views/notifications/NotificationTemplates.vue')
  },
  {
    path: '/notifications/announcements',
    name: 'Announcements',
    component: () => import('../views/notifications/Announcements.vue')
  },
  // 策略管理
  {
    path: '/policies/list',
    name: 'PolicyList',
    component: () => import('../views/policies/PolicyList.vue')
  },
  {
    path: '/policies/configs',
    name: 'PolicyConfigs',
    component: () => import('../views/policies/PolicyConfigs.vue')
  },
  {
    path: '/policies/compliance',
    name: 'ComplianceRules',
    component: () => import('../views/policies/ComplianceRules.vue')
  },
  // 会员管理
  {
    path: '/members/points',
    name: 'MemberPoints',
    component: () => import('../views/members/MemberPoints.vue')
  },
  {
    path: '/members/coupons',
    name: 'MemberCoupons',
    component: () => import('../views/members/MemberCoupons.vue')
  },
  {
    path: '/members/promotions',
    name: 'MemberPromotions',
    component: () => import('../views/members/MemberPromotions.vue')
  },
  // 告警中心
  {
    path: '/alerts/rules',
    name: 'AlertRules',
    component: () => import('../views/alerts/AlertRules.vue')
  },
  {
    path: '/alerts/list',
    name: 'AlertList',
    component: () => import('../views/alerts/AlertList.vue')
  },
  {
    path: '/alerts/settings',
    name: 'AlertSettings',
    component: () => import('../views/alerts/AlertSettings.vue')
  },
  // ============ Sprint 11: 告警通知 ============
  {
    path: '/alert/notification',
    name: 'AlertNotification',
    component: () => import('../views/alert/AlertNotificationView.vue'),
    meta: { title: '告警通知' }
  },
  {
    path: '/alert/history',
    name: 'AlertHistory',
    component: () => import('../views/alert/AlertHistoryView.vue'),
    meta: { title: '告警历史' }
  },
  {
    path: '/alert/notification-logs',
    name: 'NotificationLogs',
    component: () => import('../views/alert/NotificationLogsView.vue'),
    meta: { title: '通知日志' }
  },
  // 租户管理
  {
    path: '/tenants/approval',
    name: 'TenantApproval',
    component: () => import('../views/tenants/TenantApproval.vue')
  },
  {
    path: '/tenants/management',
    name: 'TenantManagement',
    component: () => import('../views/tenants/TenantManagement.vue')
  },
  {
    path: '/tenants/settings',
    name: 'TenantSettings',
    component: () => import('../views/tenants/TenantSettings.vue')
  },
  {
    path: '/tenants/public-archives',
    name: 'PublicArchives',
    component: () => import('../views/tenants/PublicArchives.vue')
  },
  {
    path: '/tenants/system-info',
    name: 'SystemInfo',
    component: () => import('../views/tenants/SystemInfo.vue')
  },
  // 多维权限
  {
    path: '/permissions/roles',
    name: 'PermissionRoles',
    component: () => import('../views/permissions/Roles.vue')
  },
  {
    path: '/permissions/menus',
    name: 'PermissionMenus',
    component: () => import('../views/permissions/Menus.vue')
  },
  {
    path: '/permissions/groups',
    name: 'PermissionGroups',
    component: () => import('../views/permissions/PermissionGroups.vue')
  },
  {
    path: '/permissions/data-config',
    name: 'DataPermissionConfig',
    component: () => import('../views/permissions/DataPermissionConfig.vue')
  },
  {
    path: '/permissions/api',
    name: 'ApiPermissions',
    component: () => import('../views/permissions/ApiPermissions.vue')
  },
  // 组织管理
  {
    path: '/org/companies',
    name: 'OrgCompanies',
    component: () => import('../views/org/Companies.vue')
  },
  {
    path: '/org/departments',
    name: 'OrgDepartments',
    component: () => import('../views/org/Departments.vue')
  },
  {
    path: '/org/posts',
    name: 'OrgPosts',
    component: () => import('../views/org/Posts.vue')
  },
  {
    path: '/org/employees',
    name: 'OrgEmployees',
    component: () => import('../views/org/Employees.vue')
  },

  // ============ Sprint 12: 安全功能 ============
  {
    path: '/security/permission',
    name: 'PermissionAssignment',
    component: () => import('../views/security/PermissionAssignmentView.vue'),
    meta: { title: '权限分配' }
  },
  {
    path: '/security/data-permission',
    name: 'DataPermission',
    component: () => import('../views/security/DataPermissionView.vue'),
    meta: { title: '数据权限' }
  },
  {
    path: '/security/certificate',
    name: 'CertificateManage',
    component: () => import('../views/security/CertificateManageView.vue'),
    meta: { title: '证书管理' }
  },
  {
    path: '/security/device',
    name: 'DeviceSecurity',
    component: () => import('../views/security/DeviceSecurityView.vue'),
    meta: { title: '设备安全' }
  },
  {
    path: '/security/ldap',
    name: 'LDAPConfig',
    component: () => import('../views/security/LDAPConfigView.vue'),
    meta: { title: 'LDAP配置' }
  },
  {
    path: '/security/user-sync',
    name: 'UserSync',
    component: () => import('../views/security/UserSyncView.vue'),
    meta: { title: '用户同步' }
  },

  // ============ Sprint 25: 安全与合规前端 ============
  {
    path: '/security/settings',
    name: 'SecuritySettings',
    component: () => import('../views/security/SecuritySettingsView.vue'),
    meta: { title: '安全设置' }
  },
  {
    path: '/security/privacy',
    name: 'DataPrivacy',
    component: () => import('../views/security/DataPrivacyView.vue'),
    meta: { title: '数据隐私' }
  },
  {
    path: '/security/audit-log',
    name: 'AuditLog',
    component: () => import('../views/security/AuditLogView.vue'),
    meta: { title: '审计日志' }
  },

  // ============ Sprint 13: 全球化设置 ============
  {
    path: '/globalization',
    name: 'GlobalizationSettings',
    component: () => import('../views/globalization/GlobalizationSettingsView.vue'),
    meta: { title: '全球化设置' }
  },
  {
    path: '/globalization/region',
    name: 'RegionManage',
    component: () => import('../views/globalization/RegionManageView.vue'),
    meta: { title: '区域管理' }
  },
  {
    path: '/globalization/timezone',
    name: 'TimezoneSettings',
    component: () => import('../views/globalization/TimezoneSettingsView.vue'),
    meta: { title: '时区设置' }
  },
  {
    path: '/globalization/data-residency',
    name: 'DataResidency',
    component: () => import('../views/globalization/DataResidencyView.vue'),
    meta: { title: '数据驻留规则' }
  },
  {
    path: '/globalization/ai-node',
    name: 'RegionalAINode',
    component: () => import('../views/globalization/RegionalAINodeView.vue'),
    meta: { title: '区域AI节点' }
  },
  {
    path: '/globalization/sync-status',
    name: 'RegionSyncStatus',
    component: () => import('../views/globalization/RegionSyncStatusView.vue'),
    meta: { title: '跨区域同步状态' }
  },

  // ============ Sprint 14: AI 行为监控 ============
  {
    path: '/ai/quality-dashboard',
    name: 'AIQualityDashboard',
    component: () => import('../views/ai/AIQualityDashboardView.vue'),
    meta: { title: 'AI 质量监控' }
  },
  {
    path: '/ai/behavior-log',
    name: 'AIBehaviorLog',
    component: () => import('../views/ai/AIBehaviorLogView.vue'),
    meta: { title: 'AI 行为日志' }
  },
  {
    path: '/ai/model-version',
    name: 'ModelVersion',
    component: () => import('../views/ai/ModelVersionView.vue'),
    meta: { title: '模型版本管理' }
  },
  {
    path: '/ai/sandbox',
    name: 'AISandbox',
    component: () => import('../views/ai/AISandboxView.vue'),
    meta: { title: 'AI 沙箱测试' }
  },
  {
    path: '/ai/behavior-detail/:id',
    name: 'AIBehaviorDetail',
    component: () => import('../views/ai/AIBehaviorDetailView.vue'),
    meta: { title: 'AI 行为详情' }
  },
  {
    path: '/ai/model-publish',
    name: 'ModelPublishWorkflow',
    component: () => import('../views/ai/ModelPublishWorkflow.vue'),
    meta: { title: '模型发布工作流' }
  },

  // ============ Sprint 17: 宠物情绪识别和响应 ============
  {
    path: '/emotion/recognize',
    name: 'EmotionRecognize',
    component: () => import('../views/emotion/EmotionRecognizeView.vue'),
    meta: { title: '情绪识别配置' }
  },
  {
    path: '/emotion/logs',
    name: 'EmotionLog',
    component: () => import('../views/emotion/EmotionLogView.vue'),
    meta: { title: '情绪日志' }
  },
  {
    path: '/emotion/response-config',
    name: 'EmotionResponseConfig',
    component: () => import('../views/emotion/EmotionResponseConfigView.vue'),
    meta: { title: '响应配置' }
  },
  {
    path: '/emotion/reports',
    name: 'EmotionReport',
    component: () => import('../views/emotion/EmotionReportView.vue'),
    meta: { title: '情绪报告' }
  },

  // ============ Sprint 18: 宠物数字孪生 ============
  {
    path: '/digital-twin/vitals',
    name: 'DigitalTwinVitals',
    component: () => import('../views/digital-twin/VitalsDashboardView.vue'),
    meta: { title: '生命体征仪表盘' }
  },
  {
    path: '/digital-twin/vitals-chart',
    name: 'DigitalTwinVitalsChart',
    component: () => import('../views/digital-twin/RealTimeVitalsChart.vue'),
    meta: { title: '实时体征曲线' }
  },
  {
    path: '/digital-twin/history',
    name: 'DigitalTwinHistory',
    component: () => import('../views/digital-twin/HistoricalReplayView.vue'),
    meta: { title: '历史回放' }
  },
  {
    path: '/digital-twin/behavior',
    name: 'DigitalTwinBehavior',
    component: () => import('../views/digital-twin/BehaviorPredictionView.vue'),
    meta: { title: '行为预测' }
  },

  // ============ Sprint 19: 健康追踪和预警 ============
  {
    path: '/health/warnings',
    name: 'HealthWarnings',
    component: () => import('../views/health/HealthWarningView.vue'),
    meta: { title: '健康预警' }
  },
  {
    path: '/health/exercise',
    name: 'ExerciseStats',
    component: () => import('../views/health/ExerciseStatsView.vue'),
    meta: { title: '运动统计' }
  },
  {
    path: '/health/sleep',
    name: 'SleepAnalysis',
    component: () => import('../views/health/SleepAnalysisView.vue'),
    meta: { title: '睡眠分析' }
  },
  {
    path: '/health/report',
    name: 'HealthReport',
    component: () => import('../views/health/HealthReportView.vue'),
    meta: { title: '健康报告' }
  },

  // ============ Sprint 20: 家庭和多用户场景 ============
  {
    path: '/family/members',
    name: 'FamilyMembers',
    component: () => import('../views/family/FamilyMembersView.vue'),
    meta: { title: '家庭成员' }
  },
  {
    path: '/family/child-mode',
    name: 'ChildMode',
    component: () => import('../views/family/ChildModeView.vue'),
    meta: { title: '儿童模式' }
  },
  {
    path: '/family/elder-mode',
    name: 'ElderMode',
    component: () => import('../views/family/ElderModeView.vue'),
    meta: { title: '老人陪伴模式' }
  },
  {
    path: '/family/album',
    name: 'FamilyAlbum',
    component: () => import('../views/family/FamilyAlbumView.vue'),
    meta: { title: '家庭相册' }
  },
  {
    path: '/family/settings',
    name: 'FamilySettings',
    component: () => import('../views/family/FamilySettingsView.vue'),
    meta: { title: '家庭设置' }
  },

  // ============ Sprint 21: 内容生态前端 ============
  {
    path: '/market/emoticons',
    name: 'EmoticonMarket',
    component: () => import('../views/market/EmoticonMarketView.vue'),
    meta: { title: '表情包市场' }
  },
  {
    path: '/market/actions',
    name: 'ActionMarket',
    component: () => import('../views/market/ActionMarketView.vue'),
    meta: { title: '动作资源库' }
  },
  {
    path: '/market/voices',
    name: 'VoiceConfig',
    component: () => import('../views/market/VoiceConfigView.vue'),
    meta: { title: '声音定制' }
  },

  // ============ Sprint 21: 具身智能前端 ============
  {
    path: '/embodied/:device_id/perception',
    name: 'EmbodiedPerception',
    component: () => import('../views/embodied/PerceptionView.vue'),
    meta: { title: '环境感知' }
  },
  {
    path: '/embodied/:device_id/map',
    name: 'EmbodiedMap',
    component: () => import('../views/embodied/MapView.vue'),
    meta: { title: '地图管理' }
  },
  {
    path: '/embodied/:device_id/navigate',
    name: 'EmbodiedNavigation',
    component: () => import('../views/embodied/NavigationView.vue'),
    meta: { title: '导航控制' }
  },
  {
    path: '/embodied/action-library',
    name: 'EmbodiedActionLibrary',
    component: () => import('../views/embodied/ActionLibraryView.vue'),
    meta: { title: '动作库' }
  },
  {
    path: '/embodied/:device_id/safety',
    name: 'EmbodiedSafety',
    component: () => import('../views/embodied/SafetyZonesView.vue'),
    meta: { title: '安全禁区' }
  },
  {
    path: '/embodied/:device_id/decision',
    name: 'EmbodiedDecision',
    component: () => import('../views/embodied/DecisionLogsView.vue'),
    meta: { title: '决策日志' }
  },

  // ============ Sprint 22: App 端页面 ============
  {
    path: '/app/devices',
    name: 'AppDeviceList',
    component: () => import('../views/app/AppDeviceListView.vue'),
    meta: { title: 'App设备列表', mobile: true }
  },
  {
    path: '/app/device/:id',
    name: 'AppDeviceControl',
    component: () => import('../views/app/AppDeviceControlView.vue'),
    meta: { title: 'App设备控制', mobile: true }
  },

  // ============ Sprint 22: 微信小程序 H5 页面 ============
  {
    path: '/miniapp/home',
    name: 'MiniAppHome',
    component: () => import('../views/app/MiniAppHomeView.vue'),
    meta: { title: '小程序首页', mobile: true }
  },
  {
    path: '/miniapp/devices',
    name: 'MiniAppDevices',
    component: () => import('../views/app/AppDeviceListView.vue'),
    meta: { title: '小程序设备列表', mobile: true }
  },
  {
    path: '/miniapp/device/:id',
    name: 'MiniAppDevice',
    component: () => import('../views/app/MiniAppDeviceView.vue'),
    meta: { title: '小程序设备控制', mobile: true }
  },

  // ============ Sprint 23: 第三方集成前端 ============
  {
    path: '/integration/smarthome',
    name: 'SmartHome',
    component: () => import('../views/integration/SmartHomeView.vue'),
    meta: { title: '智能家居' }
  },
  {
    path: '/integration/pet-hospital',
    name: 'PetHospital',
    component: () => import('../views/integration/PetHospitalView.vue'),
    meta: { title: '宠物医疗' }
  },
  {
    path: '/integration/pet-shop',
    name: 'PetShop',
    component: () => import('../views/integration/PetShopView.vue'),
    meta: { title: '宠物用品商城' }
  },

  // ============ Sprint 24: 研究平台前端 ============
  {
    path: '/research/data',
    name: 'ResearchData',
    component: () => import('../views/research/ResearchDataView.vue'),
    meta: { title: '数据集管理' }
  },
  {
    path: '/research/experiment',
    name: 'ResearchExperiment',
    component: () => import('../views/research/ExperimentView.vue'),
    meta: { title: 'AI 行为实验' }
  },

  // ============ Sprint 25: 开放平台前端 ============
  {
    path: '/platform/developer',
    name: 'PlatformDeveloper',
    component: () => import('../views/platform/DeveloperConsoleView.vue'),
    meta: { title: '开发者控制台' }
  },
  {
    path: '/platform/webhooks',
    name: 'PlatformWebhooks',
    component: () => import('../views/platform/WebhookMarketView.vue'),
    meta: { title: 'Webhook 市场' }
  },
  {
    path: '/platform/api-docs',
    name: 'PlatformApiDocs',
    component: () => import('../views/platform/ApiDocsView.vue'),
    meta: { title: 'API 文档' }
  }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫：检查登录状态
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  // 如果访问登录页面，直接放行
  if (to.path === '/login') {
    if (token) {
      // 已登录访问登录页，跳转到dashboard
      next('/dashboard')
    } else {
      next()
    }
    return
  }
  
  // 其他页面需要登录
  if (!token) {
    next('/login')
    return
  }
  
  // 已登录，允许访问
  next()
})

export default router
