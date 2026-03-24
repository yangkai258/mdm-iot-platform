<!-- App.vue - 应用入口（Arco Design Pro 风格） -->
<template>
  <a-config-provider :locale="locale">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import zhCN from '@arco-design/web-vue/es/locale/lang/zh-cn';
import enUS from '@arco-design/web-vue/es/locale/lang/en-us';
import useLocale from '@/hooks/locale';

const { currentLocale } = useLocale();

// 根据当前语言返回对应的 locale
const locale = computed(() => {
  return currentLocale.value === 'zh-CN' ? zhCN : enUS;
});
</script>

<style>
/* 全局样式变量 */
:root {
  --color-bg-1: #f5f7fa;
  --color-bg-2: #ffffff;
  --color-bg-3: #f0f2f5;
  --color-text-1: #303133;
  --color-text-2: #606266;
  --color-text-3: #909399;
  --color-border: #dcdfe6;
  --color-primary: #165dff;
  --color-primary-light-1: #e8f0ff;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  width: 100%;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100%;
  height: 100%;
}

/* 路由切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
