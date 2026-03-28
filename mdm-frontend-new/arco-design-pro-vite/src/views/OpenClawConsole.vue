<template>
  <div class="container">
    <a-card>
      <template #title><icon-desktop /> OpenClaw控制台</template>
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="版本">v2.1.0</a-descriptions-item>
        <a-descriptions-item label="状态"><a-badge status="success" text="运行中" /></a-descriptions-item>
      </a-descriptions>
      <a-button type="primary" style="margin-top: 16px">进入控制台</a-button>
    </a-card>
  </div>
</template>
<script setup>
</script>
<style scoped>.container { padding: 16px; }</style>
