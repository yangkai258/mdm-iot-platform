<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-run /> 运动追踪</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="今日步数" :value="stats.todaySteps" suffix="步" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="今日活动时长" :value="stats.todayDuration" suffix="分钟" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="消耗卡路里" :value="stats.calories" suffix="kcal" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="目标完成率" :value="stats.goalRate" suffix="%" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="运动趋势">
            <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="运动类型分布">
            <div style="height:250px;background:#f5f5f5;border-radius:4px;display:flex;align-items:center;justify-content:center;color:#999;font-size:14px">Chart placeholder</div>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="运动记录" style="margin-top: 16px">
        <a-table :columns="columns" :data="records" />
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ todaySteps: 2500, todayDuration: 45, calories: 180, goalRate: 83 })

const trendChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value', name: '步数' },
  series: [{ type: 'bar', data: [2000, 2500, 2200, 2800, 3000, 3500, 2500] }]
})

const pieChart = reactive({
  tooltip: { trigger: 'item' },
  series: [{ type: 'pie', radius: '60%', data: [
    { name: '散步', value: 45 }, { name: '跑步', value: 30 }, { name: '玩耍', value: 25 }
  ]}]
})

const columns = [
  { title: '时间', dataIndex: 'time' },
  { title: '运动类型', dataIndex: 'type' },
  { title: '时长', dataIndex: 'duration' },
  { title: '消耗', dataIndex: 'calories' }
]
const records = ref([
  { time: '2026-03-28 10:00', type: '散步', duration: '30分钟', calories: '120kcal' }
])
</script>

<style scoped>
.container { padding: 16px; }
</style>
