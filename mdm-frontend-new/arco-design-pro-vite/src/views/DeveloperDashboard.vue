<template>
  <div class="container">
    <a-card>
      <template #title><icon-dashboard /> 开发者仪表盘</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="API调用" :value="stats.apiCalls" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="错误率" :value="stats.errorRate" suffix="%" /></a-card></a-col>
      </a-row>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ apiCalls: 125680, errorRate: 0.5 })
</script>
<style scoped>.container { padding: 16px; }</style>
