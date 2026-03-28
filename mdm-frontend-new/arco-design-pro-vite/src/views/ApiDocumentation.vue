<template>
  <div class="container">
    <a-card>
      <template #title><icon-book /> API文档</template>
      <a-tabs>
        <a-tab-pane key="overview" title="概览">
          <a-typography-title :level="4">MDM 控制中台 API</a-typography-title>
          <a-typography-paragraph>版本: v2.0</a-typography-paragraph>
        </a-tab-pane>
        <a-tab-pane key="auth" title="认证">认证相关接口</a-tab-pane>
        <a-tab-pane key="devices" title="设备">设备管理接口</a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>
<script setup>
</script>
<style scoped>.container { padding: 16px; }</style>
