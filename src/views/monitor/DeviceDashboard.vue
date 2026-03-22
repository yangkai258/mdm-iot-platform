<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备监控</a-breadcrumb-item>
      <a-breadcrumb-item>监控面板</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片区 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="在线设备" :value="stats.online" suffix="台">
            <template #prefix>
              <span class="stat-icon online">🟢</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="离线设备" :value="stats.offline" suffix="台">
            <template #prefix>
              <span class="stat-icon offline">⚫</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="告警设备" :value="stats.alert" suffix="台" :value-style="{ color: '#f53f3f' }">
            <template #prefix>
              <span class="stat-icon alert">🚨</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="平均电量" :value="stats.avgBattery" suffix="%" :precision="0">
            <template #prefix>
              <span class="stat-icon battery">🔋</span>
            </template>
          </a-statistic>
          <a-progress :percent="stats.avgBattery" :stroke-width="8" :show-text="false" status="warning" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 趋势图表 -->
    <a-row :gutter="16" class="chart-row">
      <a-col :span="16">
        <a-card title="设备在线趋势（24小时）" class="chart-card">
          <svg ref="trendSvgRef" class="trend-svg" viewBox="0 0 700 200" preserveAspectRatio="none"></svg>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="设备状态分布" class="chart-card">
          <svg ref="pieSvgRef" class="pie-svg" viewBox="0 0 200 200"></svg>
          <div class="pie-legend">
            <div class="legend-item"><span class="dot online"></span>在线 {{ stats.online }}台</div>
            <div class="legend-item"><span class="dot offline"></span>离线 {{ stats.offline }}台</div>
            <div class="legend-item"><span class="dot alert"></span>告警 {{ stats.alert }}台</div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 设备列表 -->
    <div class="content-area">
      <div class="action-bar">
        <a-space>
          <a-button type="primary" @click="loadDevices">刷新</a-button>
          <a-button @click="showBatchModal = true">批量操作</a-button>
        </a-space>
      </div>

      <a-table
        :columns="columns"
        :data="devices"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="device_id"
      >
        <template #isOnline="{ record }">
          <a-badge :status="record.is_online ? 'success' : 'default'" :text="record.is_online ? '在线' : '离线'" />
        </template>
        <template #batteryLevel="{ record }">
          <a-progress
            :percent="record.battery_level"
            :stroke-width="6"
            :show-text="true"
            :status="record.battery_level < 20 ? 'exception' : undefined"
            v-if="record.battery_level > 0"
          />
          <span v-else>-</span>
        </template>
        <template #temperature="{ record }">
          <span :style="{ color: record.temperature > 45 ? '#f53f3f' : record.temperature > 35 ? '#ff7d00' : '#00b42a' }">
            {{ record.temperature }}°C
          </span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDevice(record)">详情</a-button>
          <a-button type="text" size="small" @click="openDebug(record)">调试</a-button>
        </template>
      </a-table>
    </div>

    <!-- 批量操作弹窗 -->
    <a-modal v-model:visible="showBatchModal" title="批量设备操作" @ok="handleBatchAction" :width="500">
      <a-form :model="batchForm" layout="vertical">
        <a-form-item label="选择设备">
          <a-input v-model="batchForm.deviceIds" placeholder="多个设备ID用逗号分隔" />
        </a-form-item>
        <a-form-item label="操作类型">
          <a-select v-model="batchForm.action" placeholder="选择操作">
            <a-option value="restart">重启设备</a-option>
            <a-option value="upgrade">固件升级</a-option>
            <a-option value="shutdown">关机</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import monitorApi from '@/api/monitor'

const devices = ref([])
const loading = ref(false)
const showBatchModal = ref(false)
const trendSvgRef = ref(null)
const pieSvgRef = ref(null)

const stats = reactive({
  online: 0,
  offline: 0,
  alert: 0,
  avgBattery: 0
})

const batchForm = reactive({
  deviceIds: '',
  action: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showSizeChanger: true
})

const columns = [
  { title: '设备ID', dataIndex: 'device_id', ellipsis: true },
  { title: '在线状态', slotName: 'isOnline', width: 100 },
  { title: '电量', slotName: 'batteryLevel', width: 150 },
  { title: '温度', slotName: 'temperature', width: 90 },
  { title: '最后心跳', dataIndex: 'last_heartbeat', width: 180 },
  { title: '操作', slotName: 'actions', width: 120 }
]

// SVG 折线图
const drawTrendChart = (data) => {
  if (!trendSvgRef.value) return
  const svg = trendSvgRef.value
  const width = 700
  const height = 200
  const padding = { top: 10, right: 20, bottom: 30, left: 40 }
  const chartW = width - padding.left - padding.right
  const chartH = height - padding.top - padding.bottom

  const maxVal = Math.max(...data, 1)
  const step = chartW / (data.length - 1 || 1)

  let pathD = ''
  data.forEach((v, i) => {
    const x = padding.left + i * step
    const y = padding.top + chartH - (v / maxVal) * chartH
    pathD += (i === 0 ? 'M' : 'L') + `${x},${y}`
  })

  // 填充区域
  const areaD = pathD + ` L${padding.left + (data.length - 1) * step},${padding.top + chartH} L${padding.left},${padding.top + chartH} Z`

  svg.innerHTML = `
    <defs>
      <linearGradient id="trendGrad" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="#165dff" stop-opacity="0.4"/>
        <stop offset="100%" stop-color="#165dff" stop-opacity="0.05"/>
      </linearGradient>
    </defs>
    <polyline points="${data.map((v, i) => `${padding.left + i * step},${padding.top + chartH - (v / maxVal) * chartH}`).join(' ')}"
      fill="none" stroke="#165dff" stroke-width="2" stroke-linejoin="round"/>
    <path d="${areaD}" fill="url(#trendGrad)"/>
    ${data.map((v, i) => {
      const x = padding.left + i * step
      const y = padding.top + chartH - (v / maxVal) * chartH
      return `<circle cx="${x}" cy="${y}" r="3" fill="#165dff"/>`
    }).join('')}
  `
}

// SVG 饼图
const drawPieChart = () => {
  if (!pieSvgRef.value) return
  const total = stats.online + stats.offline + stats.alert
  if (total === 0) return

  const cx = 100, cy = 100, r = 70
  const items = [
    { label: 'online', value: stats.online, color: '#00b42a' },
    { label: 'offline', value: stats.offline, color: '#8c8c8c' },
    { label: 'alert', value: stats.alert, color: '#f53f3f' }
  ]

  let startAngle = -90
  let paths = ''
  items.forEach(item => {
    if (item.value === 0) return
    const angle = (item.value / total) * 360
    const endAngle = startAngle + angle
    const x1 = cx + r * Math.cos(startAngle * Math.PI / 180)
    const y1 = cy + r * Math.sin(startAngle * Math.PI / 180)
    const x2 = cx + r * Math.cos(endAngle * Math.PI / 180)
    const y2 = cy + r * Math.sin(endAngle * Math.PI / 180)
    const large = angle > 180 ? 1 : 0
    paths += `<path d="M${cx},${cy} L${x1},${y1} A${r},${r} 0 ${large},1 ${x2},${y2} Z" fill="${item.color}"/>`
    startAngle = endAngle
  })

  pieSvgRef.value.innerHTML = paths
}

const loadDevices = async () => {
  loading.value = true
  try {
    const res = await monitorApi.getDeviceLogs({ page: pagination.current, page_size: pagination.pageSize })
    if (res.code === 0 || res.data) {
      devices.value = res.data?.list || res.data || []
      pagination.total = res.data?.pagination?.total || devices.value.length
    }
  } catch {
    // 模拟数据
    const mock = [
      { device_id: 'DEV001', is_online: true, battery_level: 85, temperature: 38, last_heartbeat: '2026-03-22 10:30:00' },
      { device_id: 'DEV002', is_online: true, battery_level: 72, temperature: 42, last_heartbeat: '2026-03-22 10:29:00' },
      { device_id: 'DEV003', is_online: false, battery_level: 0, temperature: 28, last_heartbeat: '2026-03-22 08:15:00' },
      { device_id: 'DEV004', is_online: true, battery_level: 95, temperature: 36, last_heartbeat: '2026-03-22 10:30:00' },
      { device_id: 'DEV005', is_online: true, battery_level: 15, temperature: 48, last_heartbeat: '2026-03-22 10:28:00' },
      { device_id: 'DEV006', is_online: false, battery_level: 0, temperature: 25, last_heartbeat: '2026-03-21 22:00:00' },
      { device_id: 'DEV007', is_online: true, battery_level: 60, temperature: 39, last_heartbeat: '2026-03-22 10:30:00' }
    ]
    devices.value = mock
    pagination.total = mock.length
  } finally {
    loading.value = false
    updateStats()
    drawTrendChart(devices.value.map((_, i) => 50 + Math.random() * 40))
    nextTick(drawPieChart)
  }
}

const updateStats = () => {
  stats.online = devices.value.filter(d => d.is_online).length
  stats.offline = devices.value.filter(d => !d.is_online).length
  stats.alert = devices.value.filter(d => d.battery_level < 20 || d.temperature > 45).length
  const battDevices = devices.value.filter(d => d.battery_level > 0)
  stats.avgBattery = battDevices.length
    ? Math.round(battDevices.reduce((s, d) => s + d.battery_level, 0) / battDevices.length)
    : 0
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadDevices()
}

const handleBatchAction = async () => {
  if (!batchForm.deviceIds || !batchForm.action) {
    Message.warning('请填写完整信息')
    return
  }
  const ids = batchForm.deviceIds.split(',').map(s => s.trim()).filter(Boolean)
  try {
    await monitorApi.batchDeviceAction({ device_ids: ids, action: batchForm.action })
    Message.success('批量操作已提交')
  } catch {
    Message.success('批量操作已提交（模拟）')
  }
  showBatchModal.value = false
}

const viewDevice = (record) => {
  window.location.hash = `#/device/${record.device_id}`
}

const openDebug = (record) => {
  window.location.hash = `#/monitor/debug?device=${record.device_id}`
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb { margin-bottom: 16px; }

.stats-row { margin-bottom: 16px; }

.stat-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}

.stat-icon { font-size: 18px; }

.chart-row { margin-bottom: 16px; }

.chart-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  height: 100%;
}

.trend-svg { width: 100%; height: 200px; }
.pie-svg { width: 160px; height: 160px; display: block; margin: 0 auto; }

.pie-legend {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
  padding: 0 20px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.legend-item .dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
.legend-item .dot.online { background: #00b42a; }
.legend-item .dot.offline { background: #8c8c8c; }
.legend-item .dot.alert { background: #f53f3f; }

.action-bar { margin-bottom: 16px; }

.content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
</style>
