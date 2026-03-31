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
    component: () => import('../views/pet/PetConsole.vue')
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
  // 会员管理 - 会员积分
  {
    path: '/members/points',
    name: 'MemberPoints',
    component: () => import('../views/members/MemberPoints.vue')
  },
  {
    path: '/members/points/rules',
    name: 'PointsRules',
    component: () => import('../views/members/PointsRules.vue')
  },
  {
    path: '/members/points/inventory',
    name: 'PointsInventory',
    component: () => import('../views/members/PointsInventory.vue')
  },
  {
    path: '/members/points/records',
    name: 'PointsRecords',
    component: () => import('../views/members/PointsRecords.vue')
  },
  {
    path: '/members/points/settings',
    name: 'PointsSettings',
    component: () => import('../views/members/PointsSettings.vue')
  },
  {
    path: '/members/points/exclude',
    name: 'PointsExclude',
    component: () => import('../views/members/PointsExclude.vue')
  },
  // 会员管理 - 优惠券
  {
    path: '/members/coupons',
    name: 'MemberCoupons',
    component: () => import('../views/members/MemberCoupons.vue')
  },
  // 会员管理 - 促销
  {
    path: '/members/promotions',
    name: 'MemberPromotions',
    component: () => import('../views/members/MemberPromotions.vue')
  },
  // 会员管理 - 店铺管理
  {
    path: '/members/stores',
    name: 'MemberStores',
    component: () => import('../views/members/MemberStores.vue')
  },
  {
    path: '/members/store-sources',
    name: 'StoreSources',
    component: () => import('../views/members/StoreSources.vue')
  },
  {
    path: '/members/store-locations',
    name: 'StoreLocations',
    component: () => import('../views/members/StoreLocations.vue')
  },
  {
    path: '/members/channels',
    name: 'MemberChannels',
    component: () => import('../views/members/MemberChannels.vue')
  },
  {
    path: '/members/miniprogram',
    name: 'MiniProgram',
    component: () => import('../views/members/MiniProgram.vue')
  },
  {
    path: '/members/printers',
    name: 'Printers',
    component: () => import('../views/members/Printers.vue')
  },
  // 会员管理 - 临时会员
  {
    path: '/members/temp-members',
    name: 'TempMembers',
    component: () => import('../views/members/TempMembers.vue')
  },
  {
    path: '/members/temp-coupons',
    name: 'TempCoupons',
    component: () => import('../views/members/TempCoupons.vue')
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
  // 健康管理
  {
    path: '/health/exercise-stats',
    name: 'ExerciseStats',
    component: () => import('../views/health/ExerciseStatsView.vue')
  },
  {
    path: '/health/reports',
    name: 'HealthReports',
    component: () => import('../views/health/HealthReportView.vue')
  },
  {
    path: '/health/warnings',
    name: 'HealthWarnings',
    component: () => import('../views/health/HealthWarningView.vue')
  },
  {
    path: '/health/sleep',
    name: 'SleepAnalysis',
    component: () => import('../views/health/SleepAnalysisView.vue')
  },
  // 应用管理
  {
    path: '/apps/list',
    name: 'AppList',
    component: () => import('../views/apps/AppList.vue')
  },
  {
    path: '/apps/versions',
    name: 'AppVersions',
    component: () => import('../views/apps/AppVersions.vue')
  },
  {
    path: '/apps/distributions',
    name: 'AppDistributions',
    component: () => import('../views/apps/AppDistributions.vue')
  },
  // 第三方接入
  {
    path: '/integration/pet-hospitals',
    name: 'PetHospitals',
    component: () => import('../views/integration/PetHospitalView.vue')
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
  {
    path: '/org/standard-positions',
    name: 'StandardPositions',
    component: () => import('../views/org/StandardPositions.vue')
  },
  // 租户申请
  {
    path: '/tenants/application',
    name: 'TenantApplication',
    component: () => import('../views/tenants/TenantApplication.vue')
  },
  // 设备状态
  {
    path: '/devices/status',
    name: 'DeviceStatus',
    component: () => import('../views/DeviceStatus.vue')
  },
  // OTA固件
  {
    path: '/ota/firmware',
    name: 'OtaFirmware',
    component: () => import('../views/OtaFirmware.vue')
  },
  // 门户
  {
    path: '/portal',
    name: 'Portal',
    component: () => import('../views/Portal.vue')
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
