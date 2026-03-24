// store/modules/app.ts - 应用状态管理
import { defineStore } from 'pinia';
import { Notification } from '@arco-design/web-vue';

interface AppState {
  // 导航栏
  navbar: boolean;
  // 菜单相关
  menu: boolean;
  menuWidth: number;
  menuCollapse: boolean;
  // 顶部菜单模式
  topMenu: boolean;
  // 标签栏
  tabBar: boolean;
  // 设备类型
  device: 'desktop' | 'mobile';
  // 页脚
  footer: boolean;
  // 主题
  theme: 'light' | 'dark';
}

export const useAppStore = defineStore('app', {
  state: (): AppState => ({
    navbar: true,
    menu: true,
    menuWidth: 220,
    menuCollapse: false,
    topMenu: false,
    tabBar: true,
    device: 'desktop',
    footer: false,
    theme: 'light',
  }),

  getters: {
    hideMenu(): boolean {
      return this.device === 'mobile';
    },
  },

  actions: {
    updateSettings(payload: Partial<AppState>) {
      Object.assign(this, payload);
    },
    
    toggleMenuCollapse() {
      this.menuCollapse = !this.menuCollapse;
    },
    
    notifySuccess(message: string) {
      Notification.success({ content: message, duration: 3000 });
    },
    
    notifyError(message: string) {
      Notification.error({ content: message, duration: 5000 });
    },
  },
});
