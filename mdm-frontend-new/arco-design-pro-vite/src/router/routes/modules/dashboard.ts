import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const DASHBOARD: AppRouteRecordRaw = {
  path: '/dashboard',
  name: 'dashboard',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.dashboard',
    requiresAuth: true,
    icon: 'icon-dashboard',
    order: 0,
  },
  children: [
    {
      path: 'workplace',
      name: 'Workplace',
      component: () => import('@/views/dashboard/workplace/index.vue'),
      meta: {
        locale: 'menu.dashboard.workplace',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'monitor',
      name: 'Monitor',
      component: () => import('@/views/dashboard/monitor/index.vue'),
      meta: {
        locale: 'menu.dashboard.monitor',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
  ],
};

// MDM 业务路由
const MDM_ROUTES: AppRouteRecordRaw[] = [
  {
    path: '/members',
    name: 'Members',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.members',
      requiresAuth: true,
      icon: 'icon-user',
      order: 1,
    },
    children: [
      { path: '', redirect: '/members/list' },
      {
        path: 'list',
        name: 'MemberList',
        component: () => import('@/views/members/MemberListView.vue'),
        meta: { locale: 'menu.members.list', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'points',
        name: 'MemberPoints',
        component: () => import('@/views/PointsRulesView.vue'),
        meta: { locale: 'menu.members.points', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'coupons',
        name: 'MemberCoupons',
        component: () => import('@/views/CouponsView.vue'),
        meta: { locale: 'menu.members.coupons', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/devices',
    name: 'Devices',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.devices',
      requiresAuth: true,
      icon: 'icon-device',
      order: 1,
    },
    children: [
      { path: '', redirect: '/devices/list' },
      {
        path: 'list',
        name: 'DeviceList',
        component: () => import('@/views/DeviceList.vue'),
        meta: { locale: 'menu.devices.list', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'detail/:id',
        name: 'DeviceDetail',
        component: () => import('@/views/DeviceDetail.vue'),
        meta: { locale: 'menu.devices.detail', requiresAuth: true, roles: ['*'], hideInMenu: true },
      },
      {
        path: 'pairing',
        name: 'DevicePairing',
        component: () => import('@/views/DevicePairing.vue'),
        meta: { locale: 'menu.devices.pairing', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'groups',
        name: 'DeviceGroups',
        component: () => import('@/views/DeviceGroupsView.vue'),
        meta: { locale: 'menu.devices.groups', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/pets',
    name: 'Pets',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.pets',
      requiresAuth: true,
      icon: 'icon-heart',
      order: 1,
    },
    children: [
      { path: '', redirect: '/pets/list' },
      {
        path: 'list',
        name: 'PetList',
        component: () => import('@/views/PetsView.vue'),
        meta: { locale: 'menu.pets.list', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'health',
        name: 'PetHealth',
        component: () => import('@/views/HealthTrackingView.vue'),
        meta: { locale: 'menu.pets.health', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/orders',
    name: 'Orders',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.orders',
      requiresAuth: true,
      icon: 'icon-shopping-cart',
      order: 1,
    },
    children: [
      {
        path: 'list',
        name: 'OrderList',
        component: () => import('@/views/OrdersView.vue'),
        meta: { locale: 'menu.orders.list', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/roles',
    name: 'Roles',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.roles',
      requiresAuth: true,
      icon: 'icon-team',
      order: 1,
    },
    children: [
      {
        path: '',
        redirect: '/roles/list',
      },
      {
        path: 'list',
        name: 'RoleList',
        component: () => import('@/views/RolesView.vue'),
        meta: { locale: 'menu.roles.list', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/permissions',
    name: 'Permissions',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.permissions',
      requiresAuth: true,
      icon: 'icon-safe',
      order: 1,
    },
    children: [
      {
        path: '',
        redirect: '/permissions/list',
      },
      {
        path: 'list',
        name: 'PermissionList',
        component: () => import('@/views/PermissionsView.vue'),
        meta: { locale: 'menu.permissions.list', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/alerts',
    name: 'Alerts',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.alerts',
      requiresAuth: true,
      icon: 'icon-warning',
      order: 1,
    },
    children: [
      {
        path: '',
        redirect: '/alerts/list',
      },
      {
        path: 'list',
        name: 'AlertList',
        component: () => import('@/views/AlertsView.vue'),
        meta: { locale: 'menu.alerts.list', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'rules',
        name: 'AlertRules',
        component: () => import('@/views/AlertRulesView.vue'),
        meta: { locale: 'menu.alerts.rules', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/ota',
    name: 'OTA',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.ota',
      requiresAuth: true,
      icon: 'icon-upload',
      order: 1,
    },
    children: [
      {
        path: '',
        redirect: '/ota/list',
      },
      {
        path: 'list',
        name: 'OTAList',
        component: () => import('@/views/OTAView.vue'),
        meta: { locale: 'menu.ota.list', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/ai',
    name: 'AI',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.ai',
      requiresAuth: true,
      icon: 'icon-robot',
      order: 1,
    },
    children: [
      {
        path: 'models',
        name: 'AIModels',
        component: () => import('@/views/AIModelsView.vue'),
        meta: { locale: 'menu.ai.models', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/subscriptions',
    name: 'Subscriptions',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.subscriptions',
      requiresAuth: true,
      icon: 'icon-subscription',
      order: 1,
    },
    children: [
      {
        path: '',
        redirect: '/subscriptions/list',
      },
      {
        path: 'list',
        name: 'SubscriptionList',
        component: () => import('@/views/SubscriptionsView.vue'),
        meta: { locale: 'menu.subscriptions.list', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/tenants',
    name: 'Tenants',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.tenants',
      requiresAuth: true,
      icon: 'icon-building',
      order: 1,
    },
    children: [
      {
        path: '',
        redirect: '/tenants/list',
      },
      {
        path: 'list',
        name: 'TenantList',
        component: () => import('@/views/TenantsView.vue'),
        meta: { locale: 'menu.tenants.list', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/notifications',
    name: 'Notifications',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.notifications',
      requiresAuth: true,
      icon: 'icon-message',
      order: 1,
    },
    children: [
      {
        path: '',
        redirect: '/notifications/list',
      },
      {
        path: 'list',
        name: 'NotificationList',
        component: () => import('@/views/NotificationsView.vue'),
        meta: { locale: 'menu.notifications.list', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
  {
    path: '/system',
    name: 'System',
    component: DEFAULT_LAYOUT,
    meta: {
      locale: 'menu.system',
      requiresAuth: true,
      icon: 'icon-settings',
      order: 1,
    },
    children: [
      {
        path: 'users',
        name: 'SystemUsers',
        component: () => import('@/views/UsersView.vue'),
        meta: { locale: 'menu.system.users', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'settings',
        name: 'SystemSettings',
        component: () => import('@/views/SystemSettingsView.vue'),
        meta: { locale: 'menu.system.settings', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'dictionaries',
        name: 'Dictionaries',
        component: () => import('@/views/DictionariesView.vue'),
        meta: { locale: 'menu.system.dictionaries', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'departments',
        name: 'Departments',
        component: () => import('@/views/DepartmentsView.vue'),
        meta: { locale: 'menu.system.departments', requiresAuth: true, roles: ['*'] },
      },
      {
        path: 'logs',
        name: 'OperationLogs',
        component: () => import('@/views/OperationLogsView.vue'),
        meta: { locale: 'menu.system.logs', requiresAuth: true, roles: ['*'] },
      },
    ],
  },
];

export default [DASHBOARD, ...MDM_ROUTES];
