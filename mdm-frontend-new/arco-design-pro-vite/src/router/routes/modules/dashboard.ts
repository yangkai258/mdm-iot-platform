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
    {
      path: 'members',
      name: 'DashboardMembers',
      component: () => import('@/views/members/MemberListView.vue'),
      meta: {
        locale: 'menu.members',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'devices',
      name: 'DashboardDevices',
      component: () => import('@/views/DeviceList.vue'),
      meta: {
        locale: 'menu.devices',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'roles',
      name: 'DashboardRoles',
      component: () => import('@/views/RolesView.vue'),
      meta: {
        locale: 'menu.roles',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'permissions',
      name: 'DashboardPermissions',
      component: () => import('@/views/PermissionsView.vue'),
      meta: {
        locale: 'menu.permissions',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'alerts',
      name: 'DashboardAlerts',
      component: () => import('@/views/AlertsView.vue'),
      meta: {
        locale: 'menu.alerts',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'ota',
      name: 'DashboardOTA',
      component: () => import('@/views/OTAView.vue'),
      meta: {
        locale: 'menu.ota',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'tenants',
      name: 'DashboardTenants',
      component: () => import('@/views/TenantsView.vue'),
      meta: {
        locale: 'menu.tenants',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'subscriptions',
      name: 'DashboardSubscriptions',
      component: () => import('@/views/SubscriptionsView.vue'),
      meta: {
        locale: 'menu.subscriptions',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default DASHBOARD;
