<template>
  <div class="container">
    <a-card>
      <template #title><icon-dashboard /> 数据仪表盘</template>
      <a-row :gutter="16">
        <a-col :span="6"><a-card><a-statistic title="总用户" :value="stats.users" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="总设备" :value="stats.devices" /></a-card></a-col>
        <a-col :span="6"><a-card><a-statistic title="今日活跃" :value="stats.active" /></a-card></a-col>
      </a-row>
      <a-chart :option="chartOption" style="height: 300px; margin-top: 16px" />
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const stats = reactive({ users: 12580, devices: 3200, active: 856 })
const chartOption = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: [820, 932, 901, 1234, 1290, 1330, 1320] }]
})
</script>
<style scoped>.container { padding: 16px; }</style>
