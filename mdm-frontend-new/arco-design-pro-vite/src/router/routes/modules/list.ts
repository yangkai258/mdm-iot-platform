import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const LIST: AppRouteRecordRaw = {
  path: '/list',
  name: 'list',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.list',
    requiresAuth: true,
    icon: 'icon-list',
    order: 2,
  },
  children: [
    {
      path: 'search-table', // The midline path complies with SEO specifications
      name: 'SearchTable',
      component: () => import('@/views/list/search-table/index.vue'),
      meta: {
        locale: 'menu.list.searchTable',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'card',
      name: 'Card',
      component: () => import('@/views/list/card/index.vue'),
      meta: {
        locale: 'menu.list.cardList',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

// MDM Business Routes
export const MDM_ROUTES: AppRouteRecordRaw[] = [
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
      {
        path: '',
        redirect: '/members/list',
      },
      {
        path: 'list',
        name: 'MemberList',
        component: () => import('@/views/members/MemberListView.vue'),
        meta: {
          locale: 'menu.members.list',
          requiresAuth: true,
          roles: ['*'],
        },
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
      {
        path: '',
        redirect: '/devices/list',
      },
      {
        path: 'list',
        name: 'DeviceList',
        component: () => import('@/views/DeviceList.vue'),
        meta: {
          locale: 'menu.devices.list',
          requiresAuth: true,
          roles: ['*'],
        },
      },
    ],
  },
];
