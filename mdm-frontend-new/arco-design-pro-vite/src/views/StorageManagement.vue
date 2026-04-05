<template>
  <div class="container">
    <a-card>
      <template #title><icon-database /> 存储管理</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="已用" :value="stats.used" suffix="GB" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="总计" :value="stats.total" suffix="GB" /></a-card></a-col>
      </a-row>
      <a-progress :percent="stats.used / stats.total * 100" :show-text="true" style="margin-top: 16px" />
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ used: 256, total: 500 })
</script>
<style scoped>.container { padding: 16px; }</style>
