<!-- layout/index.vue - 主布局（Arco Design） -->
<template>
  <a-layout class="layout" :class="{ mobile: isMobile }">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-if="renderMenu"
      class="layout-sider"
      :collapsed="collapsed"
      :collapsible="true"
      :width="220"
      :collapsed-width="64"
      :hide-trigger="true"
      breakpoint="xl"
      @collapse="handleCollapse"
    >
      <Menu />
    </a-layout-sider>
    
    <!-- 移动端抽屉菜单 -->
    <a-drawer
      v-if="isMobile"
      :visible="drawerVisible"
      placement="left"
      :footer="false"
      mask-closable
      :closable="false"
      @cancel="drawerVisible = false"
    >
      <Menu />
    </a-drawer>
    
    <!-- 主内容区 -->
    <a-layout class="layout-content-wrapper">
      <!-- 顶部导航 -->
      <a-layout-header class="layout-header">
        <NavBar />
      </a-layout-header>
      
      <!-- 标签栏 -->
      <TabBar v-if="appStore.tabBar" />
      
      <!-- 内容 -->
      <a-layout-content class="layout-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>
      
      <!-- 页脚 -->
      <a-layout-footer v-if="appStore.footer">
        <Footer />
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, provide, onMounted } from 'vue';
import { useAppStore } from '@/store';
import NavBar from '@/components/navbar/index.vue';
import Menu from '@/components/menu/index.vue';
import Footer from '@/components/footer/index.vue';
import TabBar from '@/components/tab-bar/index.vue';

const appStore = useAppStore();

// 响应式
const isMobile = ref(false);
const drawerVisible = ref(false);

// 计算属性
const collapsed = computed(() => appStore.menuCollapse);
const renderMenu = computed(() => appStore.menu);

// 折叠/展开
const handleCollapse = (val: boolean) => {
  appStore.updateSettings({ menuCollapse: val });
};

// 切换抽屉菜单
const toggleDrawer = () => {
  drawerVisible.value = !drawerVisible.value;
};

// 提供给子组件
provide('toggleDrawerMenu', toggleDrawer);

// 初始化
onMounted(() => {
  const checkMobile = () => {
    isMobile.value = window.innerWidth < 992;
  };
  
  checkMobile();
  window.addEventListener('resize', checkMobile);
  
  // 初始化用户信息
  const userStore = useUserStore();
  userStore.initFromStorage();
});

import { useUserStore } from '@/store';
</script>

<style scoped lang="less">
@nav-size-height: 60px;

.layout {
  width: 100%;
  height: 100vh;
  
  &.mobile {
    .layout-sider {
      display: none;
    }
  }
}

.layout-sider {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 100;
  height: 100%;
  background: var(--color-bg-2);
  border-right: 1px solid var(--color-border);
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
  
  :deep(.arco-layout-sider-children) {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
}

.layout-content-wrapper {
  min-height: 100vh;
  transition: padding-left 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
  
  &:not(.mobile) {
    margin-left: 220px;
    
    .layout-sider:not(.arco-layout-sider-collapsed) ~ & {
      margin-left: 220px;
    }
    
    .layout-sider.arco-layout-sider-collapsed ~ & {
      margin-left: 64px;
    }
  }
}

.layout-header {
  position: sticky;
  top: 0;
  z-index: 99;
  height: @nav-size-height;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
}

.layout-content {
  min-height: calc(100vh - @nav-size-height);
  padding: 16px;
  background: var(--color-bg-1);
  overflow-y: auto;
}

// 路由切换动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
