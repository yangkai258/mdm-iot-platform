<template>
  <div class="container">
    <a-card>
      <template #title><icon-scan /> 数字孪生概览</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="建模设备" :value="stats.devices" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="同步状态" value="正常" /></a-card></a-col>
      </a-row>
      <a-button type="primary" style="margin-top: 16px">进入数字孪生</a-button>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ devices: 128 })
</script>
<style scoped>.container { padding: 16px; }</style>
