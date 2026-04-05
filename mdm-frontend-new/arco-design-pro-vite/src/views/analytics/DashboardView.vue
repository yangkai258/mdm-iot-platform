<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>数据分析</a-breadcrumb-item>
      <a-breadcrumb-item>分析仪表板</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- Tab 切换区 -->
    <div class="pro-tabs-bar">
      <a-tabs v-model:active-tab="activeTab" @change="onTabChange">
        <a-tab-pane key="dashboard" title="Dashboard" />
        <a-tab-pane key="devices" title="设备统计" />
        <a-tab-pane key="ota" title="OTA统计" />
        <a-tab-pane key="members" title="会员分析" />
        <a-tab-pane key="alerts" title="告警统计" />
      </a-tabs>
    </div>

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card class="filter-card">
        <a-space wrap>
          <a-select v-model="timeRange" placeholder="时间范围" style="width: 120px" @change="loadData">
            <a-option value="today">今日</a-option>
            <a-option value="week">近7天</a-option>
            <a-option value="month">近30天</a-option>
          </a-select>
          <a-range-picker v-model="customRange" style="width: 260px" @change="onCustomRangeChange" />
          <a-button @click="loadData">刷新</a-button>
        </a-space>
      </a-card>
    </div>

    <!-- Dashboard Tab -->
    <div v-show="activeTab === 'dashboard'">
      <!-- 核心指标卡片 -->
      <a-row :gutter="[16, 16]" class="stat-cards-row">
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="设备总量" :value="stats.devices?.total || 0" :value-from="0" :animation-duration="800">
              <template #extra>
                <a-tag color="arcoblue" size="small">设备</a-tag>
              </template>
            </a-statistic>
            <div class="stat-trend" v-if="deviceTrend.length >= 2">
              <span class="trend-label">较上期:</span>
              <span :class="deviceTrendDelta >= 0 ? 'trend-up' : 'trend-down'">
                {{ deviceTrendDelta >= 0 ? '+' : '' }}{{ deviceTrendDelta }}
              </span>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="在线设备" :value="stats.devices?.online || 0" :value-from="0" :animation-duration="800">
              <template #extra>
                <a-tag color="green" size="small">在线</a-tag>
              </template>
            </a-statistic>
            <div class="online-rate">
              <span class="rate-label">在线率</span>
              <a-progress :percent="stats.devices?.online_rate || 0" :show-text="true" :stroke-width="6" size="small" />
            </div>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="离线设备" :value="stats.devices?.offline || 0" :value-from="0" :animation-duration="800">
              <template #extra>
                <a-tag color="red" size="small">离线</a-tag>
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="今日告警" :value="stats.alerts?.total_today || 0" :value-from="0" :animation-duration="800">
              <template #extra>
                <a-tag color="orangered" size="small">告警</a-tag>
              </template>
            </a-statistic>
            <div class="stat-trend">
              <span class="trend-label">待处理:</span>
              <span class="trend-warning">{{ stats.alerts?.pending || 0 }}</span>
            </div>
          </a-card>
        </a-col>
      </a-row>

      <!-- 会员 & OTA 指标 -->
      <a-row :gutter="[16, 16]" class="stat-cards-row">
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="会员总量" :value="stats.members?.total || 0" :value-from="0" :animation-duration="800">
              <template #extra>
                <a-tag color="purple" size="small">会员</a-tag>
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="今日活跃" :value="stats.members?.active_today || 0" :value-from="0" :animation-duration="800">
              <template #extra>
                <a-tag color="cyan" size="small">活跃</a-tag>
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="OTA成功率" :value="stats.ota?.avg_success_rate || 0" suffix="%" :value-from="0" :animation-duration="800" :precision="1">
              <template #extra>
                <a-tag color="green" size="small">OTA</a-tag>
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6">
          <a-card class="stat-card">
            <a-statistic title="进行中任务" :value="stats.ota?.running_tasks || 0" :value-from="0" :animation-duration="800">
              <template #extra>
                <a-tag color="arcoblue" size="small">OTA</a-tag>
              </template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>

      <!-- 图表区 -->
      <a-row :gutter="[16, 16]" class="charts-row">
        <a-col :xs="24" :lg="12">
          <a-card title="设备在线趋势" class="chart-card">
            <template #extra>
              <a-select v-model="deviceGranularity" style="width: 100px" @change="loadDeviceTrend">
                <a-option value="day">按天</a-option>
                <a-option value="week">按周</a-option>
                <a-option value="month">按月</a-option>
              </a-select>
            </template>
            <div ref="deviceTrendChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="设备分布" class="chart-card">
            <div ref="deviceDistChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="[16, 16]" class="charts-row">
        <a-col :xs="24" :lg="12">
          <a-card title="会员活跃趋势" class="chart-card">
            <div ref="memberTrendChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="OTA升级成功率趋势" class="chart-card">
            <div ref="otaTrendChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 设备统计 Tab -->
    <div v-show="activeTab === 'devices'">
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :md="8">
          <a-card title="设备概览" class="overview-card">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="总设备数">{{ deviceOverview.summary?.total || 0 }}</a-descriptions-item>
              <a-descriptions-item label="在线设备">{{ deviceOverview.summary?.online || 0 }}</a-descriptions-item>
              <a-descriptions-item label="离线设备">{{ deviceOverview.summary?.offline || 0 }}</a-descriptions-item>
              <a-descriptions-item label="在线率">{{ (deviceOverview.summary?.online_rate || 0).toFixed(1) }}%</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
        <a-col :xs="24" :md="16">
          <a-card title="设备生命周期分布" class="chart-card">
            <div ref="lifecycleChartRef" class="chart-container-sm"></div>
          </a-card>
        </a-col>
      </a-row>
      <a-row :gutter="[16, 16]" class="charts-row">
        <a-col :xs="24" :lg="12">
          <a-card title="设备趋势" class="chart-card">
            <div ref="deviceTrendTabChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="硬件型号分布" class="chart-card">
            <div ref="hardwareModelChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- OTA统计 Tab -->
    <div v-show="activeTab === 'ota'">
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="总任务数" :value="otaOverview.total_tasks || 0" />
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="进行中" :value="otaOverview.running_tasks || 0" />
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="成功率" :value="otaOverview.avg_success_rate || 0" suffix="%" :precision="1" />
          </a-card>
        </a-col>
      </a-row>
      <a-row :gutter="[16, 16]" class="charts-row">
        <a-col :xs="24" :lg="12">
          <a-card title="OTA版本分布" class="chart-card">
            <div ref="otaVersionChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="OTA任务列表" class="chart-card">
            <a-table :columns="otaTaskColumns" :data="otaTasks" :loading="otaLoading" :pagination="{ pageSize: 5 }" row-key="id" size="small">
              <template #status="{ record }">
                <a-tag :color="getOtaStatusColor(record.status)">{{ record.status }}</a-tag>
              </template>
      </a-table>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 会员分析 Tab -->
    <div v-show="activeTab === 'members'">
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="会员总量" :value="memberOverview.total || 0" />
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="今日新增" :value="memberOverview.new_today || 0" />
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="今日活跃" :value="memberOverview.active_today || 0" />
          </a-card>
        </a-col>
      </a-row>
      <a-row :gutter="[16, 16]" class="charts-row">
        <a-col :xs="24" :lg="12">
          <a-card title="会员等级分布" class="chart-card">
            <div ref="memberLevelChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="会员消费趋势" class="chart-card">
            <div ref="memberConsumptionChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 告警统计 Tab -->
    <div v-show="activeTab === 'alerts'">
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="今日告警" :value="alertOverview.total_today || 0" />
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="待处理" :value="alertOverview.pending || 0" />
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="8">
          <a-card class="stat-card">
            <a-statistic title="解决率" :value="alertOverview.resolution_rate || 0" suffix="%" :precision="1" />
          </a-card>
        </a-col>
      </a-row>
      <a-row :gutter="[16, 16]" class="charts-row">
        <a-col :xs="24" :lg="12">
          <a-card title="告警趋势" class="chart-card">
            <div ref="alertTrendChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="告警分布" class="chart-card">
            <div ref="alertDistChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import * as analytics from '@/api/analytics'
import * as echarts from 'echarts'

// 状态
const activeTab = ref('dashboard')
const timeRange = ref('week')
const customRange = ref([])
const deviceGranularity = ref('day')

const stats = ref({
  devices: {},
  alerts: {},
  members: {},
  ota: {}
})
const deviceTrend = ref([])
const deviceTrendDelta = ref(0)
const deviceOverview = ref({ summary: {} })
const otaOverview = ref({})
const otaTasks = ref([])
const otaLoading = ref(false)
const memberOverview = ref({})
const alertOverview = ref({})

// 图表 ref
const deviceTrendChartRef = ref(null)
const deviceDistChartRef = ref(null)
const memberTrendChartRef = ref(null)
const otaTrendChartRef = ref(null)
const lifecycleChartRef = ref(null)
const hardwareModelChartRef = ref(null)
const deviceTrendTabChartRef = ref(null)
const otaVersionChartRef = ref(null)
const memberLevelChartRef = ref(null)
const memberConsumptionChartRef = ref(null)
const alertTrendChartRef = ref(null)
const alertDistChartRef = ref(null)

const otaTaskColumns = [
  { title: '任务ID', dataIndex: 'id', width: 80 },
  { title: '固件版本', dataIndex: 'firmware_version', width: 120 },
  { title: '状态', slotName: 'status' },
  { title: '成功率', dataIndex: 'success_rate', width: 80 }
]

// 加载数据
async function loadData() {
  if (activeTab.value === 'dashboard') {
    await loadDashboard()
  }
}

async function loadDashboard() {
  try {
    const [dashRes, deviceRes, otaRes, memberRes, alertRes] = await Promise.allSettled([
      analytics.getDashboardStats(),
      analytics.getDeviceStatsOverview(),
      analytics.getOtaStatsOverview(),
      analytics.getMemberStatsOverview(),
      analytics.getAlertStatsOverview()
    ])

    if (dashRes.status === 'fulfilled') stats.value = dashRes.value.data || {}
    if (deviceRes.status === 'fulfilled') deviceOverview.value = deviceRes.value.data || {}
    if (otaRes.status === 'fulfilled') otaOverview.value = otaRes.value.data || {}
    if (memberRes.status === 'fulfilled') memberOverview.value = memberRes.value.data || {}
    if (alertRes.status === 'fulfilled') alertOverview.value = alertRes.value.data || {}

    await loadDeviceTrend()
    await nextTick()
    renderDashboardCharts()
  } catch (e) {
    console.error('loadDashboard error:', e)
  }
}

async function loadDeviceTrend() {
  try {
    const res = await analytics.getDeviceStatsTrend({
      granularity: deviceGranularity.value,
      start_date: getStartDate(),
      end_date: getEndDate()
    })
    deviceTrend.value = res.data?.list || []
    if (deviceTrend.value.length >= 2) {
      const last = deviceTrend.value[deviceTrend.value.length - 1]?.total || 0
      const prev = deviceTrend.value[deviceTrend.value.length - 2]?.total || 0
      deviceTrendDelta.value = last - prev
    }
  } catch (e) {
    console.error('loadDeviceTrend error:', e)
  }
}

function getStartDate() {
  if (customRange.value && customRange.value.length === 2) {
    return customRange.value[0].format('YYYY-MM-DD')
  }
  const days = timeRange.value === 'today' ? 1 : timeRange.value === 'week' ? 7 : 30
  const d = new Date()
  d.setDate(d.getDate() - days)
  return d.toISOString().split('T')[0]
}

function getEndDate() {
  if (customRange.value && customRange.value.length === 2) {
    return customRange.value[1].format('YYYY-MM-DD')
  }
  return new Date().toISOString().split('T')[0]
}

function onTabChange() {
  nextTick(() => {
    if (activeTab.value === 'devices') renderDeviceCharts()
    if (activeTab.value === 'ota') renderOtaCharts()
    if (activeTab.value === 'members') renderMemberCharts()
    if (activeTab.value === 'alerts') renderAlertCharts()
  })
}

function onCustomRangeChange() {
  loadData()
}

function getOtaStatusColor(status) {
  const map = { running: 'arcoblue', completed: 'green', failed: 'red', paused: 'gray' }
  return map[status] || 'gray'
}

// ========== 图表渲染 ==========

function renderDashboardCharts() {
  if (deviceTrendChartRef.value) {
    const chart = echarts.init(deviceTrendChartRef.value)
    const dates = deviceTrend.value.map(d => d.date)
    const totals = deviceTrend.value.map(d => d.total)
    const onlineAvgs = deviceTrend.value.map(d => d.online_avg)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['总设备', '在线平均'] },
      xAxis: { type: 'category', data: dates },
      yAxis: { type: 'value' },
      series: [
        { name: '总设备', type: 'line', data: totals, smooth: true },
        { name: '在线平均', type: 'line', data: onlineAvgs, smooth: true }
      ]
    })
  }
  if (deviceDistChartRef.value) {
    const chart = echarts.init(deviceDistChartRef.value)
    const data = deviceOverview.value.by_hardware_model || []
    chart.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: 0 },
      series: [{ type: 'pie', radius: ['40%', '70%'], data: data.map(d => ({ name: d.model, value: d.count })) }]
    })
  }
  if (memberTrendChartRef.value) {
    const chart = echarts.init(memberTrendChartRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
      yAxis: { type: 'value' },
      series: [{ type: 'line', data: [120, 200, 150, 80, 70, 110, 130], smooth: true }]
    })
  }
  if (otaTrendChartRef.value) {
    const chart = echarts.init(otaTrendChartRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
      yAxis: { type: 'value', min: 0, max: 100 },
      series: [{ type: 'line', data: [90, 95, 88, 92, 97, 85, 91], smooth: true }]
    })
  }
}

function renderDeviceCharts() {
  if (lifecycleChartRef.value) {
    const chart = echarts.init(lifecycleChartRef.value)
    const data = deviceOverview.value.by_lifecycle || []
    chart.setOption({
      tooltip: { trigger: 'item' },
      series: [{ type: 'pie', radius: '60%', data: data.map(d => ({ name: d.status, value: d.count })) }]
    })
  }
  if (hardwareModelChartRef.value) {
    const chart = echarts.init(hardwareModelChartRef.value)
    const data = deviceOverview.value.by_hardware_model || []
    chart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: data.map(d => d.model) },
      yAxis: { type: 'value' },
      series: [{ type: 'bar', data: data.map(d => d.count) }]
    })
  }
}

function renderOtaCharts() {
  // OTA version distribution
}

function renderMemberCharts() {
  // Member charts
}

function renderAlertCharts() {
  // Alert charts
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.stat-cards-row {
  margin-bottom: 16px;
}
.charts-row {
  margin-bottom: 16px;
}
.stat-card {
  text-align: center;
}
.stat-trend {
  margin-top: 8px;
  font-size: 12px;
  color: #666;
}
.trend-up {
  color: #f53f3f;
}
.trend-down {
  color: #0fbf60;
}
.trend-warning {
  color: #ff7a00;
}
.online-rate {
  margin-top: 8px;
}
.rate-label {
  font-size: 12px;
  color: #666;
  margin-right: 8px;
}
.chart-card {
  height: 100%;
}
.chart-container {
  height: 300px;
  width: 100%;
}
.chart-container-sm {
  height: 200px;
  width: 100%;
}
.overview-card {
  margin-bottom: 16px;
}
</style>
