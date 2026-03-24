<!-- components/tab-bar/index.vue - 标签页 -->
<template>
  <div class="tab-bar" v-if="tabs.length > 0">
    <div class="tabs-wrapper">
      <div 
        v-for="tab in tabs" 
        :key="tab.path"
        :class="['tab', { active: currentPath === tab.path }]"
        @click="handleTabClick(tab)"
      >
        <span class="tab-title">{{ tab.title }}</span>
        <a-button 
          v-if="tabs.length > 1" 
          type="text" 
          size="mini"
          class="tab-close"
          @click.stop="handleTabClose(tab)"
        >
          <template #icon><icon-close /></template>
        </a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';

interface Tab {
  path: string;
  title: string;
  name?: string;
}

const router = useRouter();
const route = useRoute();

const tabs = ref<Tab[]>([
  { path: '/dashboard', title: '设备大盘', name: 'Dashboard' }
]);

const currentPath = computed(() => route.path);

// 监听路由变化，自动添加标签
watch(
  () => route.path,
  (path) => {
    if (route.meta?.title && !tabs.value.find(t => t.path === path)) {
      tabs.value.push({
        path,
        title: route.meta.title as string,
        name: route.name as string,
      });
    }
  },
  { immediate: true }
);

const handleTabClick = (tab: Tab) => {
  router.push(tab.path);
};

const handleTabClose = (tab: Tab) => {
  const index = tabs.value.findIndex(t => t.path === tab.path);
  tabs.value.splice(index, 1);
  
  // 如果关闭的是当前页，切换到最后一个
  if (currentPath.value === tab.path && tabs.value.length > 0) {
    router.push(tabs.value[tabs.value.length - 1].path);
  }
};

import { computed } from 'vue';
</script>

<style scoped lang="less">
.tab-bar {
  padding: 8px 16px;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
  
  .tabs-wrapper {
    display: flex;
    gap: 8px;
    overflow-x: auto;
    
    &::-webkit-scrollbar {
      height: 4px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: rgba(255, 255, 255, 0.2);
      border-radius: 2px;
    }
  }
  
  .tab {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    background: var(--color-bg-3);
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
    
    &:hover {
      background: var(--color-primary-light-1);
    }
    
    &.active {
      background: var(--color-primary);
      color: #fff;
    }
    
    .tab-title {
      font-size: 13px;
    }
    
    .tab-close {
      padding: 2px;
      opacity: 0.6;
      
      &:hover {
        opacity: 1;
      }
    }
  }
}
</style>
