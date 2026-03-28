<template>
  <div class="container">
    <a-card>
      <template #title><icon-monitor /> AI监控</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="模型调用" :value="stats.calls" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="成功率" :value="stats.successRate" suffix="%" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="平均延迟" :value="stats.latency" suffix="ms" /></a-card></a-col>
      </a-row>
      <a-chart :option="chartOption" style="height: 300px; margin-top: 16px" />
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ calls: 125680, successRate: 99.5, latency: 125 })
const chartOption = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五'] },
  yAxis: { type: 'value' },
  series: [{ type: 'line', smooth: true, data: [120, 132, 128, 145, 138] }]
})
</script>
<style scoped>.container { padding: 16px; }</style>
