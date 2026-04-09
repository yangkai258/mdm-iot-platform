<template>
  <div class="pro-page-container">
    <!-- Breadcrumb -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>Home</a-breadcrumb-item>
      <a-breadcrumb-item>AI Functions</a-breadcrumb-item>
      <a-breadcrumb-item>Model Monitor</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- Metric Cards -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="Avg Latency" :value="metrics.avg_latency" suffix="ms" :precision="0" :value-from="0" :animation="true" />
          <div class="metric-trend">
            <span class="trend-down">-12% vs yesterday</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="Model Accuracy" :value="metrics.accuracy" suffix="%" :precision="2" :value-from="0" :animation="true" />
          <div class="metric-trend">
            <span class="trend-up">+2.1% vs last week</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="Throughput" :value="metrics.throughput" suffix="req/s" :precision="0" :value-from="0" :animation="true" />
          <div class="metric-trend">Peak: {{ metrics.peak_throughput }} req/s</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="Active Models" :value="metrics.active_models" :value-from="0" :animation="true" />
          <div class="metric-trend">Version: v{{ metrics.current_version }}</div>
        </a-card>
      </a-col>
    </a-row>

    <!-- Alert Banner -->
    <a-alert v-if="alerts.length > 0" class="alert-banner" type="warning">
      <template #title>
        <span>{{ alerts.length }} anomaly detected</span>
      </template>
      <template #content>
        <a-space wrap>
          <a-tag v-for="a in alerts" :key="a.id" :color="a.level === 'critical' ? 'red' : 'orange'">
            {{ a.message }}
          </a-tag>
        </a-space>
      </template>
    </a-alert>

    <!-- Charts Row 1 -->
    <a-row :gutter="16" class="chart-row">
      <a-col :span="12">
        <a-card title="Latency Trend">
          <div ref="latencyChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="Throughput Trend">
          <div ref="throughputChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" class="chart-row">
      <a-col :span="12">
        <a-card title="Accuracy Trend">
          <div ref="accuracyChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="Model Version Distribution">
          <div ref="versionChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
    </a-row>

    <!-- Search Form -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="Model Name">
          <a-input v-model="searchForm.model_name" placeholder="Search model" allow-clear style="width: 180px" />
        </a-form-item>
        <a-form-item label="Version">
          <a-select v-model="searchForm.version" placeholder="Select version" allow-clear style="width: 140px">
            <a-option v-for="v in versions" :key="v" :value="v">v{{ v }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="Status">
          <a-select v-model="searchForm.status" placeholder="Select status" allow-clear style="width: 120px">
            <a-option value="online">Online</a-option>
            <a-option value="offline">Offline</a-option>
            <a-option value="deprecated">Deprecated</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">Search</a-button>
            <a-button @click="handleReset">Reset</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- Action Bar -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="handleRefresh">Refresh</a-button>
        <a-button @click="exportReport">Export Report</a-button>
      </a-space>
    </div>

    <!-- Model Version Table -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="modelList"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1200 }"
        @change="handleTableChange"
        row-key="id"
      >
        <template #status="{ record }">
          <a-badge :status="record.status === 'online' ? 'normal' : record.status === 'offline' ? 'error' : 'warning'" />
          <span>{{ getStatusText(record.status) }}</span>
        </template>
        <template #accuracy="{ record }">
          <a-progress :percent="record.accuracy" size="small" :show-text="true" />
        </template>
        <template #latency_ms="{ record }">
          <span :style="{ color: record.latency_ms > 300 ? '#f53f3f' : 'inherit' }">{{ record.latency_ms }}ms</span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">Detail</a-button>
        </template>
      </a-table>
    </div>

    <!-- Detail Modal -->
    <a-modal v-model:visible="detailVisible" title="Model Detail" :width="700" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="Model ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="Model Name">{{ currentRecord?.name }}</a-descriptions-item>
        <a-descriptions-item label="Version">v{{ currentRecord?.version }}</a-descriptions-item>
        <a-descriptions-item label="Status">{{ getStatusText(currentRecord?.status) }}</a-descriptions-item>
        <a-descriptions-item label="Accuracy">{{ currentRecord?.accuracy }}%</a-descriptions-item>
        <a-descriptions-item label="Avg Latency">{{ currentRecord?.latency_ms }}ms</a-descriptions-item>
        <a-descriptions-item label="Call Count">{{ currentRecord?.call_count }}</a-descriptions-item>
        <a-descriptions-item label="Published At">{{ currentRecord?.published_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const modelList = ref([])
const detailVisible = ref(false)
const currentRecord = ref(null)
const latencyChartRef = ref(null)
const throughputChartRef = ref(null)
const accuracyChartRef = ref(null)
const versionChartRef = ref(null)
let charts = []

const searchForm = reactive({
  model_name: '',
  version: '',
  status: ''
})

const versions = ref(['2.1.0', '2.0.5', '2.0.0', '1.9.0'])

const metrics = reactive({
  avg_latency: 124,
  accuracy: 94.72,
  throughput: 1280,
  peak_throughput: 3200,
  active_models: 3,
  current_version: '2.1.0'
})

const alerts = ref([
  { id: 1, level: 'warning', message: 'Behavior model v1.9.0 high latency' },
  { id: 2, level: 'info', message: 'Emotion model v2.0.5 accuracy dropped 0.5%' }
])

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: 'Model Name', dataIndex: 'name', width: 180 },
  { title: 'Version', dataIndex: 'version', width: 100 },
  { title: 'Status', width: 100, slotName: 'status' },
  { title: 'Accuracy', dataIndex: 'accuracy', width: 160, slotName: 'accuracy' },
  { title: 'Avg Latency', dataIndex: 'latency_ms', width: 110, slotName: 'latency_ms' },
  { title: 'Call Count', dataIndex: 'call_count', width: 110 },
  { title: 'Published At', dataIndex: 'published_at', width: 170 },
  { title: 'Action', width: 80, slotName: 'actions', fixed: 'right' }
]

const getStatusText = (s) => ({ online: 'Online', offline: 'Offline', deprecated: 'Deprecated' }[s] || s)

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  searchForm.model_name = ''
  searchForm.version = ''
  searchForm.status = ''
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleRefresh = () => {
  loadData()
  Message.success('Monitor data refreshed')
}

const exportReport = () => {
  Message.success('Exporting report...')
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const initCharts = () => {
  const latencyChart = echarts.init(latencyChartRef.value)
  latencyChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
    yAxis: { type: 'value' },
    series: [{ type: 'line', smooth: true, data: [180, 150, 120, 95, 110, 130, 124], areaStyle: {} }]
  })
  charts.push(latencyChart)

  const throughputChart = echarts.init(throughputChartRef.value)
  throughputChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: [320, 280, 980, 1450, 1280, 1020, 850] }]
  })
  charts.push(throughputChart)

  const accuracyChart = echarts.init(accuracyChartRef.value)
  accuracyChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
    yAxis: { type: 'value', min: 90, max: 100 },
    series: [{ type: 'line', smooth: true, data: [93.1, 93.5, 94.0, 93.8, 94.2, 94.5, 94.72] }]
  })
  charts.push(accuracyChart)

  const versionChart = echarts.init(versionChartRef.value)
  versionChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: [
        { value: 4800, name: 'v2.1.0' },
        { value: 2100, name: 'v2.0.5' },
        { value: 800, name: 'v2.0.0' },
        { value: 200, name: 'v1.9.0' }
      ]
    }]
  })
  charts.push(versionChart)
}

const loadData = async () => {
  loading.value = true
  try {
    modelList.value = [
      { id: 'mdl_001', name: 'Behavior Model', version: '2.1.0', status: 'online', accuracy: 96.5, latency_ms: 98, call_count: 12480, published_at: '2026-03-20 10:00:00' },
      { id: 'mdl_002', name: 'Emotion Model', version: '2.0.5', status: 'online', accuracy: 94.2, latency_ms: 145, call_count: 9800, published_at: '2026-03-15 14:30:00' },
      { id: 'mdl_003', name: 'Voice Model', version: '2.0.0', status: 'online', accuracy: 93.8, latency_ms: 110, call_count: 7600, published_at: '2026-03-10 09:00:00' },
      { id: 'mdl_004', name: 'Pose Model', version: '1.9.0', status: 'deprecated', accuracy: 89.5, latency_ms: 320, call_count: 3200, published_at: '2026-02-20 11:00:00' }
    ]
    pagination.total = 4
  } catch (e) {
    Message.error('Load failed')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
  nextTick(() => initCharts())
  window.addEventListener('resize', () => charts.forEach(c => c.resize()))
})

onUnmounted(() => {
  charts.forEach(c => c.dispose())
})
</script>

<style scoped lang="less">
.stat-row { margin-bottom: 16px; }
.chart-row { margin-bottom: 16px; }
.metric-trend {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-3);
}
.trend-up { color: #00b42a; }
.trend-down { color: #f53f3f; }
.alert-banner { margin-bottom: 16px; }
</style>
