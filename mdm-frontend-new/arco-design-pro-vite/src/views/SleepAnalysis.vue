<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-moon /> 睡眠分析</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="睡眠时长" :value="stats.duration" suffix="小时" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="睡眠质量" :value="stats.quality" suffix="分" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="深睡比例" :value="stats.deepSleep" suffix="%" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="REM比例" :value="stats.remSleep" suffix="%" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="睡眠阶段分布">
            <a-chart :option="sleepChart" style="height: 250px" />
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="睡眠趋势">
            <a-chart :option="trendChart" style="height: 250px" />
          </a-card>
        </a-col>
      </a-row>

      <a-card title="睡眠报告" style="margin-top: 16px">
        <a-descriptions :column="3" bordered>
          <a-descriptions-item label="入睡时间">22:30</a-descriptions-item>
          <a-descriptions-item label="醒来时间">07:00</a-descriptions-item>
          <a-descriptions-item label="睡眠效率">85%</a-descriptions-item>
        </a-descriptions>
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ duration: 8.5, quality: 85, deepSleep: 25, remSleep: 20 })

const sleepChart = reactive({
  tooltip: { trigger: 'item' },
  series: [{ type: 'pie', radius: ['40%', '70%'], data: [
    { name: '清醒', value: 5 }, { name: 'REM', value: 20 }, { name: '浅睡', value: 50 }, { name: '深睡', value: 25 }
  ]}]
})

const trendChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value', name: '小时' },
  series: [{ type: 'line', smooth: true, data: [8, 7.5, 8.5, 8, 7, 9, 8.5] }]
})
</script>

<style scoped>
.container { padding: 16px; }
</style>
