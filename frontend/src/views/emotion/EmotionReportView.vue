<template>
  <div class="emotion-report-view">
    <a-card title="情绪报告">
      <a-row :gutter="16" style="margin-bottom: 24px">
        <a-col :span="6">
          <a-statistic title="今日平均情绪" :value="summary.today_avg" :suffix="summary.today_label" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="本周情绪趋势" :value="summary.weekly_trend" suffix="%">
            <template #prefix><icon-arrow-rise /></template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="情绪事件数" :value="summary.event_count" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="情绪稳定指数" :value="summary.stability" suffix="/100" />
        </a-col>
      </a-row>

      <a-tabs>
        <a-tab-pane key="daily" title="日报">
          <a-table :columns="dailyColumns" :data="dailyData" :loading="loading" :pagination="false" />
        </a-tab-pane>
        <a-tab-pane key="weekly" title="周报">
          <a-table :columns="weeklyColumns" :data="weeklyData" :loading="loading" :pagination="false" />
        </a-tab-pane>
        <a-tab-pane key="monthly" title="月报">
          <a-table :columns="monthlyColumns" :data="monthlyData" :loading="loading" :pagination="false" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

const loading = ref(false)
const summary = reactive({ today_avg: 75, today_label: '开心', weekly_trend: 12, event_count: 23, stability: 82 })
const dailyData = ref([])
const weeklyData = ref([])
const monthlyData = ref([])

const dailyColumns = [
  { title: '日期', dataIndex: 'date' },
  { title: '主要情绪', dataIndex: 'primary_emotion' },
  { title: '情绪占比', dataIndex: 'emotion_ratio' },
  { title: '异常事件', dataIndex: 'anomaly_count' }
]

const weeklyColumns = [
  { title: '周期', dataIndex: 'week' },
  { title: '平均情绪', dataIndex: 'avg_emotion' },
  { title: '情绪变化', dataIndex: 'change' },
  { title: '建议', dataIndex: 'suggestion' }
]

const monthlyColumns = [
  { title: '月份', dataIndex: 'month' },
  { title: '整体情绪', dataIndex: 'overall' },
  { title: '峰值情绪', dataIndex: 'peak' },
  { title: '情绪多样性', dataIndex: 'diversity' }
]

onMounted(async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/emotion/reports', {
      headers: { 'Authorization': 'Bearer ' + localStorage.getItem('token') }
    })
    const json = await res.json()
    dailyData.value = json.daily || []
    weeklyData.value = json.weekly || []
    monthlyData.value = json.monthly || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
