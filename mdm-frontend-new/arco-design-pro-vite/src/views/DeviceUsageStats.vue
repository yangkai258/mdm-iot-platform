<template>
  <div class="container">
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card>
          <a-statistic title="今日使用时长" :value="stats.todayMinutes" suffix="分钟">
            <template #trend>
              <icon-trending-up style="color: #67C23A" />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="本周使用时长" :value="stats.weekMinutes" suffix="分钟" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="本月使用时长" :value="stats.monthMinutes" suffix="分钟" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="总会话次数" :value="stats.totalSessions" suffix="次" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>使用趋势</template>
      <template #extra>
        <a-space>
          <a-radio-group v-model="timeRange" type="button" @change="handleTimeChange">
            <a-radio value="today">今日</a-radio>
            <a-radio value="week">本周</a-radio>
            <a-radio value="month">本月</a-radio>
            <a-radio value="custom">自定义</a-radio>
          </a-radio-group>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出数据
          </a-button>
        </a-space>
      </template>
      <a-chart :option="usageChart" style="height: 300px" />
    </a-card>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="12">
        <a-card title="会话次数统计">
          <a-chart :option="sessionChart" style="height: 250px" />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="平均响应时长">
          <a-chart :option="responseChart" style="height: 250px" />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="使用明细" style="margin-top: 16px">
      <template #extra>
        <a-space>
          <a-input-search v-model="keyword" placeholder="搜索设备ID" style="width: 200px" />
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出
          </a-button>
        </a-space>
      </template>
      <a-table :columns="columns" :data="tableData" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({
  todayMinutes: 125,
  weekMinutes: 856,
  monthMinutes: 3240,
  totalSessions: 156
})
const timeRange = ref('week')
const keyword = ref('')
const pagination = reactive({ current: 1, pageSize: 10, total: 50 })

const columns = [
  { title: '日期', dataIndex: 'date' },
  { title: '设备ID', dataIndex: 'deviceId' },
  { title: '使用时长(分钟)', dataIndex: 'minutes' },
  { title: '会话次数', dataIndex: 'sessions' },
  { title: '平均响应(ms)', dataIndex: 'avgResponse' }
]
const tableData = ref([
  { date: '2026-03-28', deviceId: 'D001', minutes: 45, sessions: 8, avgResponse: 320 },
  { date: '2026-03-28', deviceId: 'D002', minutes: 80, sessions: 12, avgResponse: 285 }
])

const usageChart = reactive({
  tooltip: { trigger: 'axis' },
  legend: { data: ['使用时长', '会话次数'] },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: [{ type: 'value', name: '时长(分钟)' }, { type: 'value', name: '次数' }],
  series: [
    { name: '使用时长', type: 'bar', data: [120, 135, 110, 145, 130, 160, 125] },
    { name: '会话次数', type: 'line', yAxisIndex: 1, data: [15, 18, 14, 20, 16, 22, 18] }
  ]
})

const sessionChart = reactive({
  tooltip: { trigger: 'item' },
  series: [{ type: 'pie', radius: ['40%', '70%'], data: [{ name: 'AI对话', value: 85 }, { name: '语音交互', value: 10 }, { name: '其他', value: 5 }] }]
})

const responseChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['09:00', '10:00', '11:00', '12:00', '13:00', '14:00'] },
  yAxis: { type: 'value', name: 'ms' },
  series: [{ type: 'line', smooth: true, data: [320, 280, 310, 290, 340, 300] }]
})

const handleTimeChange = () => { }
const handleExport = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
