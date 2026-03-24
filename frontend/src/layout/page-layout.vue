<!-- layout/page-layout.vue - 页面内容布局 -->
<template>
  <div class="page-layout">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped lang="less">
.page-layout {
  width: 100%;
  min-height: 100%;
}
</style>
