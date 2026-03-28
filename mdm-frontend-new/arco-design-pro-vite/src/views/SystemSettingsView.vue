<template>
  <div class="container">
    <a-card>
      <template #title><icon-settings /> 系统设置</template>
      <a-list>
        <a-list-item><a-list-item-meta title="基础设置" description="系统基础配置" /><template #actions><a-link>配置</a-link></template></a-list-item>
        <a-list-item><a-list-item-meta title="安全设置" description="密码、安全策略" /><template #actions><a-link>配置</a-link></template></a-list-item>
        <a-list-item><a-list-item-meta title="通知设置" description="邮件、短信通知" /><template #actions><a-link>配置</a-link></template></a-list-item>
      </a-list>
    </a-card>
  </div>
</template>
<script setup>
</script>
<style scoped>.container { padding: 16px; }</style>
