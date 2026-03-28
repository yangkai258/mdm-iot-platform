<template>
  <div class="container">
    <a-card>
      <template #title><icon-cloud /> 云端备份</template>
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6"><a-card><a-statistic title="备份数量" :value="stats.count" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="总大小" :value="stats.size" suffix="GB" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="最后备份" :value="stats.lastBackup" /></a-card></a-col>
      </a-row>
      <a-button type="primary" @click="handleBackup"><template #icon><icon-plus /></template>立即备份</a-button>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ count: 10, size: 50, lastBackup: '2小时前' })
const handleBackup = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
