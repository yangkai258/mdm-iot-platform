<template>
  <div class="container">
    <a-card>
      <template #title><icon-skill-level /> 远程控制</template>
      <a-button type="primary" @click="handleRestart">重启设备</a-button>
      <a-button @click="handleUpgrade">固件升级</a-button>
    </a-card>
  </div>
</template>
<script setup>
const handleRestart = () => {}
const handleUpgrade = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
