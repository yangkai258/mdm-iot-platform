<!-- components/footer/index.vue - 页脚 -->
<template>
  <div class="footer">
    <div class="footer-content">
      <span>MDM 物联网平台 © 2026</span>
      <span class="divider">|</span>
      <a-link>帮助文档</a-link>
      <a-link>关于我们</a-link>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped lang="less">
.footer {
  padding: 16px;
  text-align: center;
  background: var(--color-bg-2);
  border-top: 1px solid var(--color-border);
  
  &-content {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    font-size: 14px;
    color: var(--color-text-2);
    
    .divider {
      color: var(--color-border);
    }
  }
}
</style>
