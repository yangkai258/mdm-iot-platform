import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

// MDM路由作为顶级菜单
const MDMMenu: AppRouteRecordRaw = {
  path: '/mdm',
  name: 'Mdm',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.mdm',
    requiresAuth: true,
    icon: 'icon-setup',
    order: 1,
  },
  children: [
    {
      path: 'members',
      name: 'MdmMembers',
      component: () => import('@/views/members/MemberListView.vue'),
      meta: {
        locale: 'menu.members',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'devices',
      name: 'MdmDevices',
      component: () => import('@/views/DeviceList.vue'),
      meta: {
        locale: 'menu.devices',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default MDMMenu;
