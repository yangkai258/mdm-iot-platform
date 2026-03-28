<template>
  <div class="container">
    <a-card>
      <template #title><icon-growth /> 会员成长值</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="当前成长值" :value="stats.currentGrowth" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="本月增长" :value="stats.monthlyGrowth" /></a-card></a-col>
      </a-row>
      <a-chart :option="chartOption" style="height: 300px; margin-top: 16px" />
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ currentGrowth: 5800, monthlyGrowth: 520 })
const chartOption = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: [80, 95, 75, 85, 92, 88, 78] }]
})
</script>
<style scoped>.container { padding: 16px; }</style>
