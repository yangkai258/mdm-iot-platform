<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb v-if="breadcrumbs.length > 0" class="breadcrumb">
      <a-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index">
        <a v-if="item.path" :href="'#' + item.path">{{ item.label }}</a>
        <span v-else>{{ item.label }}</span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div v-if="$slots.search" class="search-bar">
      <slot name="search" />
    </div>

    <!-- 操作按钮栏 -->
    <div v-if="$slots.actions" class="action-bar">
      <slot name="actions" />
    </div>

    <!-- 筛选面板插槽 -->
    <div v-if="$slots.filter" class="filter-panel">
      <slot name="filter" />
    </div>

    <!-- 主内容区 -->
    <div class="content-area">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  breadcrumbs?: Array<{ label: string; path?: string }>
}>()
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-bar {
  margin-bottom: 16px;
}

.action-bar {
  margin-bottom: 16px;
}

.filter-panel {
  margin-bottom: 16px;
}

.content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
