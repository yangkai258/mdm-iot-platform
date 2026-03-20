import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
