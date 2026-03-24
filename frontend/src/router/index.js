import { createRouter, createWebHistory } from 'vue-router'

// 路由配置（登录页不需要 layout）
const routes = [
  // 公开路由
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

  // 受保护的路由（需要 layout）
  {
    path: '/',
    component: () => import('../layout/index.vue'),
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '设备大盘' }
      },
      {
        path: 'devices',
        name: 'Devices',
        component: () => import('../views/DeviceDashboard.vue'),
        meta: { title: '设备列表' }
      },
      {
        path: 'device/:id',
        name: 'DeviceDetail',
        component: () => import('../views/DeviceDetail.vue'),
        meta: { title: '设备详情' }
      },
      {
        path: 'ota',
        redirect: '/ota/packages'
      },
      {
        path: 'ota/packages',
        name: 'OtaPackages',
        component: () => import('../views/ota/OtaPackages.vue'),
        meta: { title: '固件包管理' }
      },
      {
        path: 'ota/deployments',
        name: 'OtaDeployments',
        component: () => import('../views/ota/OtaDeployments.vue'),
        meta: { title: '部署任务' }
      },
      {
        path: 'alert',
        name: 'Alert',
        component: () => import('../views/Alert.vue'),
        meta: { title: '告警管理' }
      },
      // 会员管理
      {
        path: 'members',
        name: 'Members',
        component: () => import('../views/members/MemberListView.vue'),
        meta: { title: '会员列表' }
      },
      {
        path: 'members/cards',
        name: 'MemberCards',
        component: () => import('../views/members/MemberCardView.vue'),
        meta: { title: '会员卡管理' }
      },
      {
        path: 'members/coupons',
        name: 'MemberCoupons',
        component: () => import('../views/members/CouponView.vue'),
        meta: { title: '优惠券' }
      },
      {
        path: 'members/stores',
        name: 'MemberStores',
        component: () => import('../views/members/StoreView.vue'),
        meta: { title: '门店管理' }
      },
      {
        path: 'members/levels',
        name: 'MemberLevels',
        component: () => import('../views/members/LevelView.vue'),
        meta: { title: '会员等级' }
      },
      {
        path: 'members/tags',
        name: 'MemberTags',
        component: () => import('../views/members/TagView.vue'),
        meta: { title: '标签管理' }
      },
      {
        path: 'members/points',
        name: 'MemberPoints',
        component: () => import('../views/members/PointsView.vue'),
        meta: { title: '积分管理' }
      },
      {
        path: 'members/promotions',
        name: 'MemberPromotions',
        component: () => import('../views/members/MemberPromotions.vue'),
        meta: { title: '促销活动' }
      },
      // 系统管理
      {
        path: 'system/monitor',
        name: 'SystemMonitor',
        component: () => import('../views/system/Monitor.vue'),
        meta: { title: '服务监控' }
      },
      {
        path: 'system/logs',
        name: 'SystemLogs',
        component: () => import('../views/system/Logs.vue'),
        meta: { title: '操作日志' }
      },
      {
        path: 'policies',
        name: 'Policies',
        component: () => import('../views/policies/PolicyList.vue'),
        meta: { title: '策略管理' }
      },
      // 租户管理
      {
        path: 'tenants/approval',
        name: 'TenantsApproval',
        component: () => import('../views/tenants/TenantApproval.vue'),
        meta: { title: '租户入驻审核' }
      },
      {
        path: 'tenants/management',
        name: 'TenantsManagement',
        component: () => import('../views/tenants/TenantManagement.vue'),
        meta: { title: '租户系统管理' }
      },
      {
        path: 'tenants/public-archives',
        name: 'TenantsPublicArchives',
        component: () => import('../views/tenants/PublicArchives.vue'),
        meta: { title: '公共档案' }
      },
      {
        path: 'tenants/system-info',
        name: 'TenantsSystemInfo',
        component: () => import('../views/tenants/SystemInfo.vue'),
        meta: { title: '系统信息' }
      },
      // 权限管理
      {
        path: 'permissions/groups',
        name: 'PermissionsGroups',
        component: () => import('../views/permissions/PermissionGroups.vue'),
        meta: { title: '权限组管理' }
      },
      {
        path: 'permissions/data-config',
        name: 'PermissionsDataConfig',
        component: () => import('../views/permissions/DataPermissionConfig.vue'),
        meta: { title: '数据权限配置' }
      },
      // 通知管理
      {
        path: 'notifications/list',
        name: 'NotificationsList',
        component: () => import('../views/notifications/NotificationList.vue'),
        meta: { title: '推送通知' }
      },
      {
        path: 'notifications/announcements',
        name: 'NotificationsAnnouncements',
        component: () => import('../views/notifications/Announcements.vue'),
        meta: { title: '公告管理' }
      },
      {
        path: 'notifications/templates',
        name: 'NotificationsTemplates',
        component: () => import('../views/notifications/NotificationTemplates.vue'),
        meta: { title: '通知模板' }
      },
      // AI 功能
      {
        path: 'ai/behavior',
        name: 'AiBehavior',
        component: () => import('../views/ai/AIBehaviorLogView.vue'),
        meta: { title: '行为分析' }
      },
      {
        path: 'ai/emotion',
        name: 'AiEmotion',
        component: () => import('../views/emotion/EmotionLogView.vue'),
        meta: { title: '情感识别' }
      },
      // 健康医疗
      {
        path: 'health/warnings',
        name: 'HealthWarning',
        component: () => import('../views/health/HealthWarningView.vue'),
        meta: { title: '健康预警' }
      },
      {
        path: 'health/sports',
        name: 'HealthSports',
        component: () => import('../views/health/ExerciseStatsView.vue'),
        meta: { title: '运动统计' }
      },
      {
        path: 'health/sleep',
        name: 'HealthSleep',
        component: () => import('../views/health/SleepAnalysisView.vue'),
        meta: { title: '睡眠分析' }
      },
    ]
  },

  // 404
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  // 公开路由
  if (to.path === '/login' || to.path === '/test-modals') {
    if (token && to.path === '/login') {
      next('/dashboard')
    } else {
      next()
    }
    return
  }
  
  // 需要登录
  if (!token) {
    next('/login')
    return
  }
  
  next()
})

export default router
