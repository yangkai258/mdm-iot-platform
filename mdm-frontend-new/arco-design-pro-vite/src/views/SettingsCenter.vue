<template>
  <div class="container">
    <a-card>
      <template #title><icon-settings /> 设置中心</template>
      <a-list>
        <a-list-item><a-list-item-meta title="主题设置" description="自定义界面主题" /><template #actions><a-link>配置</a-link></template></a-list-item>
        <a-list-item><a-list-item-meta title="通知设置" description="配置通知方式" /><template #actions><a-link>配置</a-link></template></a-list-item>
        <a-list-item><a-list-item-meta title="安全设置" description="密码、安全" /><template #actions><a-link>配置</a-link></template></a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
</script>
<style scoped>.container { padding: 16px; }</style>
