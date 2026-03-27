<!-- ArcoDesign Pro 标准列表页布局组件 -->
<!-- 使用方式: <ProPageLayout module="设备管理" page="设备列表"> -->
<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item v-if="module">{{ module }}</a-breadcrumb-item>
      <a-breadcrumb-item v-if="page">{{ page }}</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索框区域 -->
    <div v-if="$slots.search" class="pro-search-bar">
      <slot name="search" />
    </div>

    <!-- 操作按钮组 -->
    <div v-if="$slots.actions" class="pro-action-bar">
      <slot name="actions" />
    </div>

    <!-- 主内容区域 -->
    <div class="pro-content-area">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  module?: string
  page?: string
}>()
</script>

<style scoped>
.pro-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.pro-breadcrumb {
  margin-bottom: 16px;
}

.pro-search-bar {
  margin-bottom: 12px;
}

.pro-action-bar {
  margin-bottom: 16px;
}

.pro-content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
