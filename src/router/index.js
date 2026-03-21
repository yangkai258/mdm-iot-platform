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
  // 会员信息管理 (12.1)
  {
    path: '/member/card-types',
    name: 'MemberCardTypes',
    component: () => import('../views/member/MemberCardTypes.vue')
  },
  {
    path: '/member/card-groups',
    name: 'MemberCardGroups',
    component: () => import('../views/member/MemberCardGroups.vue')
  },
  {
    path: '/member/level-rules',
    name: 'MemberUpgradeRules',
    component: () => import('../views/member/MemberUpgradeRules.vue')
  },
  {
    path: '/member/settings',
    name: 'MemberSettings',
    component: () => import('../views/member/MemberSettings.vue')
  },
  // 会员订单 (12.2)
  {
    path: '/member/orders',
    name: 'MemberOrders',
    component: () => import('../views/member/MemberOrders.vue')
  },
  {
    path: '/member/occupation-types',
    name: 'OccupationTypes',
    component: () => import('../views/member/OccupationTypes.vue')
  },
  // 优惠券模块 (12.3)
  {
    path: '/member/coupon-inventory',
    name: 'CouponInventory',
    component: () => import('../views/member/CouponInventory.vue')
  },
  {
    path: '/member/coupon-messages',
    name: 'CouponMessages',
    component: () => import('../views/member/CouponMessages.vue')
  },
  // 会员标签 (12.4)
  {
    path: '/member/tags',
    name: 'MemberTags',
    component: () => import('../views/member/MemberTags.vue')
  },
  {
    path: '/member/tags/high-freq',
    name: 'HighFreqTags',
    component: () => import('../views/member/HighFreqTags.vue')
  },
  {
    path: '/member/tags/low-freq',
    name: 'LowFreqTags',
    component: () => import('../views/member/LowFreqTags.vue')
  },
  {
    path: '/member/tags/interest',
    name: 'InterestTags',
    component: () => import('../views/member/InterestTags.vue')
  },
  {
    path: '/member/tags/auto-clean',
    name: 'TagAutoClean',
    component: () => import('../views/member/TagAutoClean.vue')
  },
  {
    path: '/member/tags/report',
    name: 'TagReport',
    component: () => import('../views/member/TagReport.vue')
  },
  // 促销活动 (12.5)
  {
    path: '/member/redpackets',
    name: 'Redpackets',
    component: () => import('../views/member/Redpackets.vue')
  },
  {
    path: '/member/promotions/types',
    name: 'PromotionTypes',
    component: () => import('../views/member/PromotionTypes.vue')
  },
  {
    path: '/member/promotions/buy-gift',
    name: 'BuyGiftPromo',
    component: () => import('../views/member/BuyGiftPromo.vue')
  },
  {
    path: '/member/promotions/direct-reduce',
    name: 'DirectReducePromo',
    component: () => import('../views/member/DirectReducePromo.vue')
  },
  {
    path: '/member/promotions/amount-reduce',
    name: 'AmountReducePromo',
    component: () => import('../views/member/AmountReducePromo.vue')
  },
  {
    path: '/member/promotions/amount-discount',
    name: 'AmountDiscountPromo',
    component: () => import('../views/member/AmountDiscountPromo.vue')
  },
  {
    path: '/member/promotions/vip-exclusive',
    name: 'VipExclusivePromo',
    component: () => import('../views/member/VipExclusivePromo.vue')
  },
  // 会员礼包 (12.6)
  {
    path: '/member/gifts',
    name: 'MemberGifts',
    component: () => import('../views/member/MemberGifts.vue')
  },
  {
    path: '/member/gift-records',
    name: 'GiftRecords',
    component: () => import('../views/member/GiftRecords.vue')
  },
  // 会员服务 (12.7)
  {
    path: '/member/reception',
    name: 'MemberReception',
    component: () => import('../views/member/MemberReception.vue')
  },
  {
    path: '/member/articles',
    name: 'MemberArticles',
    component: () => import('../views/member/MemberArticles.vue')
  },
  {
    path: '/member/sms-templates',
    name: 'SmsTemplates',
    component: () => import('../views/member/SmsTemplates.vue')
  },
  {
    path: '/member/wechat-settings',
    name: 'WechatSettings',
    component: () => import('../views/member/WechatSettings.vue')
  },
  {
    path: '/member/benefits',
    name: 'MemberBenefits',
    component: () => import('../views/member/MemberBenefits.vue')
  },
  {
    path: '/member/sms-channels',
    name: 'SmsChannels',
    component: () => import('../views/member/SmsChannels.vue')
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
