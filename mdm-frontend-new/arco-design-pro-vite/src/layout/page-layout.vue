<template>
  <div class="page-layout">
    <Breadcrumb :items="breadcrumbItems" />
    <router-view v-slot="{ Component, route }">
      <transition name="fade" mode="out-in" appear>
        <component
          :is="Component"
          v-if="Component"
          :key="route.fullPath"
        />
        <div v-else class="page-loading">
          <a-spin size="large" tip="页面加载中..." />
        </div>
      </transition>
    </router-view>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue';
  import { useRoute } from 'vue-router';
  import { useTabBarStore } from '@/store';
  import Breadcrumb from '@/components/breadcrumb/index.vue';

  const route = useRoute();
  const tabBarStore = useTabBarStore();

  const breadcrumbItems = computed(() => {
    const items: string[] = [];
    const matched = route.matched.filter(
      (r) => r.meta && (r.meta.locale || r.meta.title)
    );
    matched.forEach((r) => {
      if (r.meta.locale) items.push(r.meta.locale as string);
    });
    if (route.meta.title) {
      items.push(route.meta.title as string);
    } else if (items.length === 0) {
      items.push('menu.dashboard');
    }
    return items;
  });
</script>

<style scoped lang="less">
  .page-layout {
    padding: 0 16px 16px;
  }
  .page-loading {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 400px;
  }
</style>
