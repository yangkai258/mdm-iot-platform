<template>
  <div class="container">
    <a-card>
      <template #title><icon-file-text /> 日志分析</template>
      <a-form layout="inline" style="margin-bottom: 16px">
        <a-form-item label="时间范围"><a-range-picker v-model="dateRange" /></a-form-item>
        <a-form-item><a-button type="primary">查询</a-button></a-form-item>
      </a-form>
      <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const dateRange = ref([])
const chartOption = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五'] },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: [120, 200, 150, 80, 70] }]
})
</script>
<style scoped>.container { padding: 16px; }</style>
