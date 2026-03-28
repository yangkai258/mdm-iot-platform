<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-book-open /> 具身AI学习进度</a-space>
      </template>

      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="总学习进度" :value="progress.total" suffix="%">
              <template #prefix>
                <a-progress :percent="progress.total" :show-text="false" :stroke-width="10" />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="已学习动作" :value="progress.learnedActions" suffix="个" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="动作成功率" :value="progress.successRate" suffix="%" />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="学习时长" :value="progress.totalHours" suffix="小时" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="动作模仿成功率趋势">
            <a-chart :option="trendChart" style="height: 250px" />
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="已学习动作库">
            <a-list>
              <a-list-item v-for="action in learnedActions" :key="action.id">
                <a-list-item-meta :title="action.name" :description="'成功率: ' + action.successRate + '% | 学习次数: ' + action.learnCount" />
                <template #actions>
                  <a-link @click="handleRelearn(action)">重新学习</a-link>
                  <a-link @click="handleViewDetail(action)">详情</a-link>
                </template>
              </a-list-item>
            </a-list>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="设备学习历史" style="margin-top: 16px">
        <template #extra>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出
          </a-button>
        </template>
        <a-table :columns="columns" :data="historyData" />
      </a-card>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const progress = reactive({ total: 68, learnedActions: 15, successRate: 82, totalHours: 24 })

const trendChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['第1周', '第2周', '第3周', '第4周'] },
  yAxis: { type: 'value', min: 0, max: 100, name: '成功率%' },
  series: [{ type: 'line', smooth: true, data: [60, 68, 75, 82], areaStyle: {} }]
})

const learnedActions = ref([
  { id: 1, name: '挥手', successRate: 95, learnCount: 25 },
  { id: 2, name: '转圈', successRate: 88, learnCount: 20 },
  { id: 3, name: '坐下', successRate: 92, learnCount: 30 }
])

const columns = [
  { title: '时间', dataIndex: 'time' },
  { title: '设备ID', dataIndex: 'deviceId' },
  { title: '动作', dataIndex: 'action' },
  { title: '成功率', dataIndex: 'successRate' },
  { title: '学习时长', dataIndex: 'duration' }
]
const historyData = ref([
  { time: '2026-03-28', deviceId: 'D001', action: '挥手', successRate: '95%', duration: '10分钟' }
])

const handleRelearn = (action) => { }
const handleViewDetail = (action) => { }
const handleExport = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
