// store/modules/user.ts - 用户状态管理
import { defineStore } from 'pinia';

interface UserState {
  userInfo: {
    id?: string | number;
    username?: string;
    nickname?: string;
    email?: string;
    avatar?: string;
    role?: string;
  } | null;
  token: string | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    userInfo: null,
    token: localStorage.getItem('token') || null,
  }),

  getters: {
    isLoggedIn(): boolean {
      return !!this.token;
    },
    role(): string {
      return this.userInfo?.role || 'admin';
    },
  },

  actions: {
    setUserInfo(info: UserState['userInfo']) {
      this.userInfo = info;
    },
    
    setToken(token: string) {
      this.token = token;
      localStorage.setItem('token', token);
    },
    
    logout() {
      this.userInfo = null;
      this.token = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    },
    
    initFromStorage() {
      const token = localStorage.getItem('token');
      const user = localStorage.getItem('user');
      if (token) this.token = token;
      if (user) {
        try {
          this.userInfo = JSON.parse(user);
        } catch (e) {
          console.error('Failed to parse user info');
        }
      }
    },
  },
});
