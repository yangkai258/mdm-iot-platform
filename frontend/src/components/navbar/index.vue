<!-- components/navbar/index.vue - 顶部导航栏 -->
<template>
  <div class="navbar">
    <div class="navbar-left">
      <!-- 移动端菜单按钮 -->
      <a-button 
        v-if="device === 'mobile'" 
        type="text" 
        @click="toggleDrawerMenu"
      >
        <template #icon><icon-menu /></template>
      </a-button>
      
      <!-- 面包屑 -->
      <a-breadcrumb v-if="device !== 'mobile'" class="breadcrumb">
        <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
          {{ item.title }}
        </a-breadcrumb-item>
      </a-breadcrumb>
    </div>
    
    <div class="navbar-right">
      <!-- 全屏切换 -->
      <a-tooltip :content="isFullscreen ? '退出全屏' : '全屏'">
        <a-button type="text" @click="toggleFullscreen">
          <template #icon>
            <icon-fullscreen-exit v-if="isFullscreen" />
            <icon-fullscreen v-else />
          </template>
        </a-button>
      </a-tooltip>
      
      <!-- 语言切换 -->
      <a-dropdown @select="handleLanguageChange">
        <a-button type="text">
          <template #icon><icon-language /></template>
        </a-button>
        <template #menu>
          <a-doption value="zh-CN">中文</a-doption>
          <a-doption value="en-US">English</a-doption>
        </template>
      </a-dropdown>
      
      <!-- 用户信息 -->
      <a-dropdown trigger="click">
        <div class="user-info">
          <a-avatar :size="32">
            <icon-user />
          </a-avatar>
          <span class="username">{{ username }}</span>
        </div>
        <template #menu>
          <a-doption value="profile">个人设置</a-doption>
          <a-divider />
          <a-doption value="logout" @click="handleLogout">退出登录</a-doption>
        </template>
      </a-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '@/store';
import {
  IconMenu,
  IconFullscreen,
  IconFullscreenExit,
  IconLanguage,
  IconUser,
} from '@arco-design/web-vue/es/icon';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

// 响应式
const device = ref('desktop');
const isFullscreen = ref(false);

// 注入切换抽屉菜单的方法
const toggleDrawerMenu = inject<() => void>('toggleDrawerMenu', () => {});

// 用户名
const username = computed(() => {
  return userStore.userInfo?.nickname || userStore.userInfo?.username || 'Admin';
});

// 面包屑
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title);
  return matched.map(item => ({
    path: item.path,
    title: item.meta?.title as string,
  }));
});

// 切换全屏
const toggleFullscreen = async () => {
  if (!document.fullscreenElement) {
    await document.documentElement.requestFullscreen();
    isFullscreen.value = true;
  } else {
    await document.exitFullscreen();
    isFullscreen.value = false;
  }
};

// 切换语言
const handleLanguageChange = (lang: string) => {
  localStorage.setItem('locale', lang);
  window.location.reload();
};

// 退出登录
const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

// 初始化
onMounted(() => {
  device.value = window.innerWidth < 992 ? 'mobile' : 'desktop';
  window.addEventListener('resize', () => {
    device.value = window.innerWidth < 992 ? 'mobile' : 'desktop';
  });
});
</script>

<style scoped lang="less">
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 16px;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
  
  &-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  
  &-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  .breadcrumb {
    margin-left: 16px;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 8px;
    border-radius: 4px;
    cursor: pointer;
    transition: background 0.2s;
    
    &:hover {
      background: var(--color-bg-3);
    }
    
    .username {
      font-size: 14px;
      color: var(--color-text-1);
    }
  }
}
</style>
