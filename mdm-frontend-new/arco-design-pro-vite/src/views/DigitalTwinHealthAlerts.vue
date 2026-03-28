<template>
  <div class="container">
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card>
          <a-statistic title="预警总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="待处理" :value="stats.pending" :value-style="{ color: '#F56C6C' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已处理" :value="stats.resolved" :value-style="{ color: '#67C23A' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="处理率" :value="stats.resolved / stats.total * 100" suffix="%" :precision="1" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>健康预警列表</template>
      <template #extra>
        <a-space>
          <a-select v-model="filterType" placeholder="预警类型" allow-clear style="width: 150px">
            <a-option value="heart">心率异常</a-option>
            <a-option value="temp">体温异常</a-option>
            <a-option value="activity">活动异常</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px">
            <a-option value="pending">待处理</a-option>
            <a-option value="resolved">已处理</a-option>
          </a-select>
        </a-space>
      </template>

      <a-table :columns="columns" :data="alerts">
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">{{ record.severity }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'pending' ? 'orange' : 'green'">{{ record.status === 'pending' ? '待处理' : '已处理' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleView(record)">详情</a-link>
          <a-link v-if="record.status === 'pending'" @click="handleResolve(record)">处理</a-link>
          <a-link @click="handleIgnore(record)">忽略</a-link>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="detailVisible" :title="'预警详情 - ' + selectedAlert.type" :width="600">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="预警ID">{{ selectedAlert.id }}</a-descriptions-item>
        <a-descriptions-item label="严重程度">
          <a-tag :color="getSeverityColor(selectedAlert.severity)">{{ selectedAlert.severity }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="触发时间">{{ selectedAlert.time }}</a-descriptions-item>
        <a-descriptions-item label="当前状态">
          <a-tag :color="selectedAlert.status === 'pending' ? 'orange' : 'green'">{{ selectedAlert.status }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="预警描述" :span="2>{{ selectedAlert.description }}</a-descriptions-item>
        <a-descriptions-item label="建议措施" :span="2>{{ selectedAlert.suggestion }}</a-descriptions-item>
      </a-descriptions>

      <a-divider>趋势图</a-divider>
      <a-chart :option="trendChart" style="height: 200px" />
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ total: 45, pending: 8, resolved: 37 })
const filterType = ref('')
const filterStatus = ref('')
const detailVisible = ref(false)
const selectedAlert = reactive({})

const columns = [
  { title: '时间', dataIndex: 'time', width: 180 },
  { title: '类型', dataIndex: 'type' },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '数值', dataIndex: 'value' },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const alerts = ref([
  { id: 1, time: '2026-03-28 10:00', type: '心率异常', severity: 'warning', value: '95bpm', status: 'pending', description: '心率持续偏高', suggestion: '减少剧烈活动' },
  { id: 2, time: '2026-03-28 09:00', type: '体温异常', severity: 'major', value: '39.5°C', status: 'resolved', description: '体温偏高', suggestion: '观察休息' }
])

const trendChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['08:00', '09:00', '10:00', '11:00', '12:00'] },
  yAxis: { type: 'value' },
  series: [{ type: 'line', smooth: true, data: [85, 88, 95, 92, 88] }]
})

const getSeverityColor = (s) => ({ warning: 'orange', major: 'red', info: 'blue' }[s] || 'gray')
const handleView = (r) => { Object.assign(selectedAlert, r); detailVisible.value = true }
const handleResolve = (r) => { }
const handleIgnore = (r) => { }
</script>

<style scoped>
.container { padding: 16px; }
</style>
