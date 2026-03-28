<template>
  <div class="container">
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card>
          <a-statistic title="系统状态">
            <template #value>
              <a-tag color="green" size="large">正常</a-tag>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="在线设备" :value="stats.onlineDevices" suffix="/ 100" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="API响应时间" :value="stats.apiResponse" suffix="ms" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="数据库连接" :value="stats.dbConnections" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="服务组件状态">
          <a-table :columns="columns" :data="services" size="small">
            <template #status="{ record }">
              <a-badge :status="getStatus(record.status)" :text="getStatusText(record.status)" />
            </template>
            <template #responseTime="{ record }">
              <a-progress :percent="record.responseTime / 1000" :max="5" :show-text="true" size="small" />
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="响应时间趋势 (5分钟)">
          <a-chart :option="responseChart" style="height: 250px" />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="操作日志" style="margin-top: 16px">
      <template #extra>
        <a-space>
          <a-input-search v-model="keyword" placeholder="搜索" style="width: 200px" />
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出
          </a-button>
        </a-space>
      </template>
      <a-table :columns="logColumns" :data="logs" size="small" :pagination="pagination">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">{{ record.level }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ onlineDevices: 95, apiResponse: 125, dbConnections: 25 })
const keyword = ref('')
const pagination = reactive({ current: 1, pageSize: 10, total: 100 })

const columns = [
  { title: '服务', dataIndex: 'name' },
  { title: '状态', slotName: 'status' },
  { title: '响应时间', slotName: 'responseTime' },
  { title: '最后检查', dataIndex: 'lastCheck' }
]
const services = ref([
  { name: 'API Gateway', status: 'healthy', responseTime: 45, lastCheck: '2026-03-28 10:00:00' },
  { name: 'MQTT Broker', status: 'healthy', responseTime: 12, lastCheck: '2026-03-28 10:00:00' },
  { name: 'PostgreSQL', status: 'healthy', responseTime: 28, lastCheck: '2026-03-28 10:00:00' },
  { name: 'Redis', status: 'healthy', responseTime: 5, lastCheck: '2026-03-28 10:00:00' }
])

const logColumns = [
  { title: '时间', dataIndex: 'time', width: 180 },
  { title: '级别', slotName: 'level', width: 80 },
  { title: '服务', dataIndex: 'service' },
  { title: '消息', dataIndex: 'message' }
]
const logs = ref([
  { time: '2026-03-28 10:00:00', level: 'info', service: 'API', message: 'Health check passed' }
])

const responseChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['09:56', '09:57', '09:58', '09:59', '10:00'] },
  yAxis: { type: 'value', name: 'ms' },
  series: [{ type: 'line', smooth: true, data: [120, 135, 128, 140, 125] }]
})

const getStatus = (s) => ({ healthy: 'success', degraded: 'warning', down: 'error' }[s] || 'default')
const getStatusText = (s) => ({ healthy: '正常', degraded: '降级', down: '故障' }[s] || s)
const getLevelColor = (l) => ({ info: 'blue', warn: 'orange', error: 'red' }[l] || 'gray')
const handleExport = () => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
