<template>
  <div class="container">
    <a-card>
      <template #title><icon-chart-bar /> 资源使用统计</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="CPU" :value="stats.cpu" suffix="%" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="内存" :value="stats.memory" suffix="%" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="磁盘" :value="stats.disk" suffix="%" /></a-card></a-col>
      </a-row>
      <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999">Chart</div>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ cpu: 45, memory: 62, disk: 38 })
const chartOption = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五'] },
  yAxis: { type: 'value' },
  series: [{ type: 'line', smooth: true, data: [40, 45, 42, 48, 45] }]
})
</script>
<style scoped>.container { padding: 16px; }</style>
