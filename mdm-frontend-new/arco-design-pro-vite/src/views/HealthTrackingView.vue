<template>
  <div class="container">
    <a-card>
      <template #title><icon-heart-pulse /> 健康追踪</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="今日步数" :value="stats.steps" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="睡眠时长" :value="stats.sleep" suffix="小时" /></a-card></a-col>
      </a-row>
      <a-chart :option="chartOption" style="height: 250px; margin-top: 16px" />
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ steps: 2500, sleep: 8.5 })
const chartOption = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五'] },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: [2200, 2800, 2500, 3000, 2500] }]
})
</script>
<style scoped>.container { padding: 16px; }</style>
