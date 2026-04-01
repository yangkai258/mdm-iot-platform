<template>
  <div class="page-layout">
    <Breadcrumb :items="breadcrumbItems" />
    <router-view v-slot="{ Component, route }">
      <transition name="fade" mode="out-in" appear>
        <component
          :is="Component"
          v-if="route.meta.ignoreCache"
          :key="route.fullPath"
        />
        <keep-alive v-else :include="cacheList">
          <component :is="Component" :key="route.fullPath" />
        </keep-alive>
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

  const cacheList = computed(() => tabBarStore.getCacheList);

  const breadcrumbItems = computed(() => {
    const items: string[] = [];
    const matched = route.matched.filter(
      (r) => r.meta && (r.meta.locale || r.meta.title)
    );
    matched.forEach((r) => {
      if (r.meta.locale) items.push(r.meta.locale as string);
    });
    // 最后一项用页面标题（locale 不存在时用 title）
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
</style>
