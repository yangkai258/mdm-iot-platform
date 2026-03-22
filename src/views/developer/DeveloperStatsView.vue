<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>开发者平台</a-breadcrumb-item>
      <a-breadcrumb-item>API 使用统计</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计概览卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="今日 API 调用" :value="stats.todayCalls" :value-style="{ color: '#165dff' }">
            <template #suffix>次</template>
          </a-statistic>
          <div class="stat-trend">
            <span :style="{ color: stats.todayTrend >= 0 ? '#52c41a' : '#ff4d4f' }">
              {{ stats.todayTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(stats.todayTrend) }}%
            </span>
            <span class="trend-label">较昨日</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="本周 API 调用" :value="stats.weekCalls" :value-style="{ color: '#00d4ff' }">
            <template #suffix>次</template>
          </a-statistic>
          <div class="stat-trend">
            <span :style="{ color: stats.weekTrend >= 0 ? '#52c41a' : '#ff4d4f' }">
              {{ stats.weekTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(stats.weekTrend) }}%
            </span>
            <span class="trend-label">较上周</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="配额使用率" :value="stats.quotaUsage" :value-style="{ color: getQuotaColor(stats.quotaUsage) }" suffix="%">
            <template #extra>
              <a-progress :percent="stats.quotaUsage" :stroke-width="6" :color="getQuotaColor(stats.quotaUsage)" :show-text="false" size="small" style="width: 120px" />
            </template>
          </a-statistic>
          <div class="stat-trend">
            <span class="quota-label">{{ stats.quotaUsed }} / {{ stats.quotaTotal }} 次</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="活跃应用数" :value="stats.activeApps" :value-style="{ color: '#722ed1' }">
            <template #suffix>个</template>
          </a-statistic>
          <div class="stat-trend">
            <span class="trend-label">共 {{ stats.totalApps }} 个应用</span>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-radio-group v-model="timeRange" type="button" @change="onTimeRangeChange">
          <a-radio value="today">今日</a-radio>
          <a-radio value="7d">近7天</a-radio>
          <a-radio value="30d">近30天</a-radio>
          <a-radio value="90d">近90天</a-radio>
        </a-radio-group>
        <a-range-picker v-model="customDateRange" style="width: 260px" @change="onCustomDateChange" />
        <a-select v-model="selectedApp" placeholder="选择应用" style="width: 160px" allow-clear>
          <a-option :value="null">全部应用</a-option>
          <a-option v-for="app in appList" :key="app.id" :value="app.id">{{ app.name }}</a-option>
        </a-select>
      </a-space>
      <a-space>
        <a-button @click="exportReport">导出报表</a-button>
        <a-button @click="loadStats">刷新</a-button>
      </a-space>
    </div>

    <!-- 主内容区 -->
    <a-row :gutter="16" style="margin-top: 16px; margin-bottom: 16px">
      <!-- API 调用趋势图 -->
      <a-col :span="16">
        <a-card class="chart-card" title="API 调用趋势">
          <template #extra>
            <a-space>
              <a-tag v-for="t in chartTypes" :key="t.value" :color="chartType === t.value ? 'arcoblue' : 'gray'" style="cursor: pointer" @click="chartType = t.value; loadChartData()">{{ t.label }}</a-tag>
            </a-space>
          </template>
          <div ref="trendChartRef" style="height: 280px"></div>
        </a-card>
      </a-col>
      <!-- API 状态分布 -->
      <a-col :span="8">
        <a-card class="chart-card" title="API 响应状态分布">
          <div ref="pieChartRef" style="height: 280px"></div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-bottom: 16px">
      <!-- 应用排行 -->
      <a-col :span="12">
        <a-card class="chart-card" title="应用调用排行">
          <a-table :columns="rankColumns" :data="appRanking" :loading="loading" :pagination="false" row-key="rank" size="small">
            <template #rank="{ record }">
              <a-tag :color="record.rank <= 3 ? 'orangered' : 'gray'" style="border-radius: 50%; width: 28px; text-align: center">{{ record.rank }}</a-tag>
            </template>
            <template #appName="{ record }">
              <a-avatar :size="24" :style="{ backgroundColor: record.color || '#165dff', marginRight: '8px' }">{{ record.app_name?.charAt(0) }}</a-avatar>
              {{ record.app_name }}
            </template>
            <template #calls="{ record }">
              <span style="font-weight: 600">{{ formatNumber(record.calls) }}</span>
            </template>
            <template #percent="{ record }">
              <a-progress :percent="record.percent" :stroke-width="6" size="small" :color="getQuotaColor(record.percent)" :show-text="false" />
              <span style="font-size: 12px; color: #86909c">{{ record.percent.toFixed(1) }}%</span>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <!-- API 端点排行 -->
      <a-col :span="12">
        <a-card class="chart-card" title="API 端点调用排行">
          <a-table :columns="endpointColumns" :data="endpointRanking" :loading="loading" :pagination="false" row-key="rank" size="small">
            <template #rank="{ record }">
              <a-badge :count="record.rank" :number-style="{ backgroundColor: record.rank <= 3 ? '#ff4d4f' : '#86909c', borderRadius: '50%', width: '20px', height: '20px', lineHeight: '20px !important', fontSize: '11px' }" />
            </template>
            <template #method="{ record }">
              <a-tag :color="getMethodColor(record.method)">{{ record.method }}</a-tag>
            </template>
            <template #path="{ record }">
              <span style="font-family: monospace; font-size: 12px; color: #4e5969">{{ record.path }}</span>
            </template>
            <template #avgLatency="{ record }">
              <a-tag :color="record.avg_latency < 200 ? 'green' : record.avg_latency < 500 ? 'orange' : 'red'">
                {{ record.avg_latency }}ms
              </a-tag>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 配额详情 -->
    <a-card class="chart-card" title="配额使用详情">
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6" v-for="quota in quotaList" :key="quota.name">
          <a-card class="quota-item-card" :bordered="false">
            <a-statistic :title="quota.label" :value="quota.used" :value-style="{ color: getQuotaColor((quota.used / quota.total) * 100) }">
              <template #suffix>/ {{ quota.total }}</template>
            </a-statistic>
            <a-progress :percent="(quota.used / quota.total) * 100" :stroke-width="8" :color="getQuotaColor((quota.used / quota.total) * 100)" size="small" style="margin-top: 8px" />
            <div style="font-size: 12px; color: #86909c; margin-top: 4px">
              剩余 {{ quota.total - quota.used }} 次
            </div>
          </a-card>
        </a-col>
      </a-row>
      <a-divider />
      <div class="quota-upgrade">
        <a-space>
          <a-typography-title :heading="6" style="margin: 0">配额升级</a-typography-title>
          <a-typography-text type="secondary">当前套餐: <a-tag color="arcoblue">免费版 (1000次/天)</a-tag></a-typography-text>
        </a-space>
        <a-space>
          <a-button type="primary" size="small">升级套餐</a-button>
          <a-button size="small">查看定价</a-button>
        </a-space>
      </div>
    </a-card>

    <!-- 调用明细表 -->
    <a-card class="chart-card" style="margin-top: 16px" title="API 调用明细">
      <template #extra>
        <a-space>
          <a-select v-model="filterStatus" placeholder="响应状态" style="width: 120px" allow-clear>
            <a-option :value="200">200 成功</a-option>
            <a-option :value="400">400 错误</a-option>
            <a-option :value="401">401 未授权</a-option>
            <a-option :value="429">429 限流</a-option>
            <a-option :value="500">500 服务错误</a-option>
          </a-select>
          <a-input-search v-model="searchKeyword" placeholder="搜索接口路径" style="width: 200px" @search="loadDetailedLogs" search-button />
        </a-space>
      </template>
      <a-table :columns="logColumns" :data="detailedLogs" :loading="logLoading" :pagination="logPagination" @change="handleLogTableChange" row-key="id" size="small">
        <template #method="{ record }">
          <a-tag :color="getMethodColor(record.method)">{{ record.method }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.status }}</a-tag>
        </template>
        <template #latency="{ record }">
          <span :style="{ color: record.latency < 200 ? '#52c41a' : record.latency < 500 ? '#ff7d00' : '#ff4d4f' }">
            {{ record.latency }}ms
          </span>
        </template>
        <template #timestamp="{ record }">
          {{ formatDateTime(record.timestamp) }}
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/developer'
const loading = ref(false)
const logLoading = ref(false)
const timeRange = ref('7d')
const customDateRange = ref([])
const selectedApp = ref(null)
const chartType = ref('calls')
const filterStatus = ref(null)
const searchKeyword = ref('')

const chartTypes = [
  { label: '调用量', value: 'calls' },
  { label: '延迟', value: 'latency' },
  { label: '错误率', value: 'errors' }
]

const trendChartRef = ref(null)
const pieChartRef = ref(null)

const stats = reactive({
  todayCalls: 0,
  todayTrend: 0,
  weekCalls: 0,
  weekTrend: 0,
  quotaUsage: 0,
  quotaUsed: 0,
  quotaTotal: 0,
  activeApps: 0,
  totalApps: 0
})

const appList = ref([])
const appRanking = ref([])
const endpointRanking = ref([])
const detailedLogs = ref([])
const quotaList = ref([])

const logPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const rankColumns = [
  { title: '排名', slotName: 'rank', width: 60 },
  { title: '应用名称', slotName: 'appName', width: 180 },
  { title: '调用次数', slotName: 'calls', width: 120 },
  { title: '占比', slotName: 'percent', width: 180 }
]

const endpointColumns = [
  { title: '排名', slotName: 'rank', width: 60 },
  { title: '方法', slotName: 'method', width: 80 },
  { title: '接口路径', slotName: 'path', width: 260 },
  { title: '调用次数', dataIndex: 'calls', width: 100 },
  { title: '平均延迟', slotName: 'avgLatency', width: 100 }
]

const logColumns = [
  { title: '时间', slotName: 'timestamp', width: 180 },
  { title: '方法', slotName: 'method', width: 80 },
  { title: '接口路径', dataIndex: 'path', width: 260 },
  { title: '应用', dataIndex: 'app_name', width: 140 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '延迟', slotName: 'latency', width: 90 },
  { title: 'IP', dataIndex: 'ip', width: 130 }
]

const formatNumber = (n) => {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n
}

const getQuotaColor = (percent) => {
  if (percent >= 90) return '#ff4d4f'
  if (percent >= 70) return '#ff7d00'
  return '#52c41a'
}

const getMethodColor = (method) => {
  const map = { GET: 'green', POST: 'blue', PUT: 'orange', DELETE: 'red', PATCH: 'purple' }
  return map[method] || 'gray'
}

const getStatusColor = (status) => {
  if (status >= 200 && status < 300) return 'green'
  if (status >= 400 && status < 500) return 'orange'
  if (status >= 500) return 'red'
  return 'gray'
}

const formatDateTime = (ts) => {
  if (!ts) return '-'
  const d = new Date(ts)
  return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')} ${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}:${String(d.getSeconds()).padStart(2,'0')}`
}

const loadStats = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/stats/overview`, { params: { time_range: timeRange.value, app_id: selectedApp.value } })
    if (res.data.code === 0) Object.assign(stats, res.data.data)
  } catch {
    Object.assign(stats, {
      todayCalls: 12453, todayTrend: 12.5, weekCalls: 87321, weekTrend: 8.3,
      quotaUsage: 62, quotaUsed: 620, quotaTotal: 1000, activeApps: 3, totalApps: 4
    })
    Message.warning('使用模拟数据')
  } finally { loading.value = false }
  loadChartData()
  loadRanking()
  loadQuota()
  loadDetailedLogs()
}

const loadChartData = async () => {
  try {
    await axios.get(`${API_BASE}/stats/trend`, { params: { time_range: timeRange.value, chart_type: chartType.value, app_id: selectedApp.value } })
  } catch {
    // 模拟图表数据
  }
  nextTick(() => {
    if (trendChartRef.value) trendChartRef.value.innerHTML = `<div style="display:flex;align-items:center;justify-content:center;height:100%;color:#86909c;font-size:14px">📊 ${chartType.value === 'calls' ? '调用量' : chartType.value === 'latency' ? '延迟' : '错误率'} 趋势图表 (集成 ECharts 可视化)</div>`
    if (pieChartRef.value) pieChartRef.value.innerHTML = `<div style="display:flex;align-items:center;justify-content:center;height:100%;color:#86909c;font-size:14px">🥧 API 响应状态分布 (集成 ECharts 可视化)</div>`
  })
}

const loadRanking = async () => {
  try {
    const res = await axios.get(`${API_BASE}/stats/ranking`, { params: { time_range: timeRange.value, app_id: selectedApp.value } })
    if (res.data.code === 0) {
      appRanking.value = res.data.data.app_ranking || []
      endpointRanking.value = res.data.data.endpoint_ranking || []
    }
  } catch {
    appRanking.value = [
      { rank: 1, app_name: 'PetCare App', calls: 45230, percent: 51.8, color: '#165dff' },
      { rank: 2, app_name: 'HomeDashboard', calls: 28450, percent: 32.6, color: '#00d4ff' },
      { rank: 3, app_name: 'EdgeGateway', calls: 9860, percent: 11.3, color: '#52c41a' },
      { rank: 4, app_name: 'DataAnalytics', calls: 3781, percent: 4.3, color: '#722ed1' }
    ]
    endpointRanking.value = [
      { rank: 1, method: 'GET', path: '/api/v1/devices/status', calls: 32100, avg_latency: 85 },
      { rank: 2, method: 'POST', path: '/api/v1/devices/command', calls: 21500, avg_latency: 142 },
      { rank: 3, method: 'GET', path: '/api/v1/data/query', calls: 18700, avg_latency: 320 },
      { rank: 4, method: 'PUT', path: '/api/v1/devices/config', calls: 9800, avg_latency: 110 },
      { rank: 5, method: 'POST', path: '/api/v1/data/report', calls: 6200, avg_latency: 95 }
    ]
  }
}

const loadQuota = async () => {
  try {
    const res = await axios.get(`${API_BASE}/quota`)
    if (res.data.code === 0) quotaList.value = res.data.data.list || []
  } catch {
    quotaList.value = [
      { name: 'daily', label: '日配额', used: 620, total: 1000 },
      { name: 'weekly', label: '周配额', used: 4120, total: 7000 },
      { name: 'monthly', label: '月配额', used: 18200, total: 30000 },
      { name: 'rate_limit', label: 'QPS 限制', used: 45, total: 100 }
    ]
  }
}

const loadDetailedLogs = async () => {
  logLoading.value = true
  try {
    const res = await axios.get(`${API_BASE}/stats/logs`, { params: { app_id: selectedApp.value, status: filterStatus.value, keyword: searchKeyword.value, page: logPagination.current, page_size: logPagination.pageSize } })
    if (res.data.code === 0) {
      detailedLogs.value = res.data.data.list || []
      logPagination.total = res.data.data.pagination?.total || detailedLogs.value.length
    }
  } catch {
    detailedLogs.value = [
      { id: 1, timestamp: '2026-03-22 18:20:00', method: 'GET', path: '/api/v1/devices/status', app_name: 'PetCare App', status: 200, latency: 45, ip: '10.0.0.1' },
      { id: 2, timestamp: '2026-03-22 18:19:55', method: 'POST', path: '/api/v1/devices/command', app_name: 'HomeDashboard', status: 200, latency: 128, ip: '10.0.0.2' },
      { id: 3, timestamp: '2026-03-22 18:19:50', method: 'GET', path: '/api/v1/data/query', app_name: 'EdgeGateway', status: 200, latency: 312, ip: '10.0.0.3' },
      { id: 4, timestamp: '2026-03-22 18:19:45', method: 'PUT', path: '/api/v1/devices/config', app_name: 'PetCare App', status: 401, latency: 30, ip: '10.0.0.4' },
      { id: 5, timestamp: '2026-03-22 18:19:40', method: 'GET', path: '/api/v1/devices/status', app_name: 'DataAnalytics', status: 429, latency: 15, ip: '10.0.0.5' }
    ]
    logPagination.total = detailedLogs.value.length
    Message.warning('使用模拟数据')
  } finally { logLoading.value = false }
}

const loadAppList = async () => {
  try {
    const res = await axios.get(`${API_BASE}/apps`)
    if (res.data.code === 0) appList.value = res.data.data.list || []
  } catch {
    appList.value = [
      { id: 1, name: 'PetCare App', color: '#165dff' },
      { id: 2, name: 'HomeDashboard', color: '#00d4ff' },
      { id: 3, name: 'EdgeGateway', color: '#52c41a' },
      { id: 4, name: 'DataAnalytics', color: '#722ed1' }
    ]
  }
}

const onTimeRangeChange = () => { loadStats() }
const onCustomDateChange = () => { if (customDateRange.value?.length === 2) loadStats() }
const exportReport = () => { Message.info('报表导出功能开发中') }
const handleLogTableChange = (pag) => { logPagination.current = pag.current; loadDetailedLogs() }

onMounted(() => {
  loadAppList()
  loadStats()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.stat-trend { margin-top: 8px; font-size: 12px; color: #86909c; text-align: center; }
.trend-label { margin-left: 4px; color: #86909c; }
.quota-label { font-size: 12px; }
.pro-action-bar { background: #fff; border-radius: 8px; padding: 12px 16px; display: flex; justify-content: space-between; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.chart-card { border-radius: 8px; }
.quota-item-card { background: #f5f7fa; border-radius: 8px; text-align: center; }
.quota-upgrade { display: flex; justify-content: space-between; align-items: center; }
</style>
