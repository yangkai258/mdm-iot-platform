<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>宠物数字孪生</a-breadcrumb-item>
      <a-breadcrumb-item>实时体征曲线</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-bar">
      <a-space>
        <a-select
          v-model="selectedPetId"
          placeholder="选择宠物"
          style="width: 200px"
          allow-search
          @change="handlePetChange"
        >
          <a-option v-for="pet in petList" :key="pet.device_id" :value="pet.device_id">
            {{ pet.pet_name }} ({{ pet.device_id }})
          </a-option>
        </a-select>
        <a-select
          v-model="timeRange"
          placeholder="时间范围"
          style="width: 140px"
          @change="loadChartData"
        >
          <a-option value="1">最近1小时</a-option>
          <a-option value="6">最近6小时</a-option>
          <a-option value="12">最近12小时</a-option>
          <a-option value="24">最近24小时</a-option>
        </a-select>
        <a-button type="primary" @click="loadChartData">
          <template #icon><icon-search /></template>
          查询
        </a-button>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="action-bar">
      <a-space>
        <a-button @click="loadChartData">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button @click="goBack">返回仪表盘</a-button>
      </a-space>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <!-- 心率曲线 -->
      <a-card class="chart-card" title="心率曲线">
        <template #extra>
          <a-tag v-if="latestHeartRate" :color="isHeartRateAbnormal ? 'red' : 'green'">
            {{ latestHeartRate }} bpm
          </a-tag>
        </template>
        <div ref="heartRateChartRef" class="chart-container"></div>
      </a-card>

      <!-- 呼吸曲线 -->
      <a-card class="chart-card" title="呼吸频率曲线">
        <template #extra>
          <a-tag v-if="latestRespiratory" :color="isRespiratoryAbnormal ? 'red' : 'green'">
            {{ latestRespiratory }} 次/分
          </a-tag>
        </template>
        <div ref="respiratoryChartRef" class="chart-container"></div>
      </a-card>

      <!-- 体温曲线 -->
      <a-card class="chart-card" title="体温曲线">
        <template #extra>
          <a-tag v-if="latestTemp" :color="isTempAbnormal ? 'red' : 'green'">
            {{ latestTemp }} ℃
          </a-tag>
        </template>
        <div ref="tempChartRef" class="chart-container"></div>
      </a-card>

      <!-- 活动量曲线 -->
      <a-card class="chart-card" title="活动量曲线">
        <template #extra>
          <a-tag v-if="latestActivity != null">
            {{ latestActivity }} 分
          </a-tag>
        </template>
        <div ref="activityChartRef" class="chart-container"></div>
      </a-card>
    </div>

    <!-- 数据统计 -->
    <div class="stats-area">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="心率平均值" :value="stats.heartRateAvg" suffix="bpm" :precision="1" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="呼吸频率平均值" :value="stats.respiratoryAvg" suffix="次/分" :precision="1" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="体温平均值" :value="stats.tempAvg" suffix="℃" :precision="2" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="活动量平均值" :value="stats.activityAvg" suffix="分" :precision="1" />
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'
import * as echarts from 'echarts'

const router = useRouter()
const route = useRoute()

const API_BASE = '/api/digital-twin'

const petList = ref([])
const selectedPetId = ref('')
const timeRange = ref('24')
const loading = ref(false)

const chartData = ref([])

const heartRateChartRef = ref(null)
const respiratoryChartRef = ref(null)
const tempChartRef = ref(null)
const activityChartRef = ref(null)

let heartRateChart = null
let respiratoryChart = null
let tempChart = null
let activityChart = null

const stats = reactive({
  heartRateAvg: 0,
  respiratoryAvg: 0,
  tempAvg: 0,
  activityAvg: 0
})

const latestHeartRate = ref(null)
const latestRespiratory = ref(null)
const latestTemp = ref(null)
const latestActivity = ref(null)

const isHeartRateAbnormal = computed(() => {
  return latestHeartRate.value && (latestHeartRate.value < 60 || latestHeartRate.value > 140)
})

const isRespiratoryAbnormal = computed(() => {
  return latestRespiratory.value && (latestRespiratory.value < 10 || latestRespiratory.value > 30)
})

const isTempAbnormal = computed(() => {
  return latestTemp.value && (latestTemp.value < 38 || latestTemp.value > 39.5)
})

import { computed } from 'vue'

// 获取宠物列表
const loadPets = async () => {
  try {
    const res = await axios.get(`${API_BASE}/pets`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      petList.value = res.data.data?.list || res.data.data || []
    }
  } catch {
    petList.value = [
      { device_id: 'PET001', pet_name: '小白' },
      { device_id: 'PET002', pet_name: '旺财' }
    ]
  }

  // 从 URL 参数获取
  if (route.query.device) {
    selectedPetId.value = route.query.device
  } else if (petList.value.length > 0) {
    selectedPetId.value = petList.value[0].device_id
  }
  if (selectedPetId.value) {
    loadChartData()
  }
}

// 加载图表数据
const loadChartData = async () => {
  if (!selectedPetId.value) {
    Message.warning('请先选择宠物')
    return
  }
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/vitals/trend/${selectedPetId.value}`, {
      params: { hours: timeRange.value },
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      chartData.value = res.data.data || []
      updateStats()
      renderCharts()
    }
  } catch {
    // 模拟数据
    const mockData = generateMockData()
    chartData.value = mockData
    updateStats()
    renderCharts()
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

// 生成模拟数据
const generateMockData = () => {
  const data = []
  const hours = parseInt(timeRange.value)
  const now = Date.now()
  for (let i = hours * 12; i >= 0; i--) {
    data.push({
      time: new Date(now - i * 5 * 60 * 1000).toISOString(),
      heart_rate: Math.floor(Math.random() * 60) + 70,
      respiratory_rate: Math.floor(Math.random() * 15) + 15,
      body_temp: parseFloat((Math.random() * 1.5 + 38).toFixed(1)),
      activity_level: Math.floor(Math.random() * 100)
    })
  }
  return data
}

// 更新统计数据
const updateStats = () => {
  if (chartData.value.length === 0) return

  const hr = chartData.value.map(d => d.heart_rate)
  const rr = chartData.value.map(d => d.respiratory_rate)
  const tp = chartData.value.map(d => d.body_temp)
  const ac = chartData.value.map(d => d.activity_level)

  stats.heartRateAvg = Math.round(hr.reduce((a, b) => a + b, 0) / hr.length)
  stats.respiratoryAvg = Math.round(rr.reduce((a, b) => a + b, 0) / rr.length * 10) / 10
  stats.tempAvg = Math.round(tp.reduce((a, b) => a + b, 0) / tp.length * 100) / 100
  stats.activityAvg = Math.round(ac.reduce((a, b) => a + b, 0) / ac.length)

  latestHeartRate.value = chartData.value[chartData.value.length - 1]?.heart_rate
  latestRespiratory.value = chartData.value[chartData.value.length - 1]?.respiratory_rate
  latestTemp.value = chartData.value[chartData.value.length - 1]?.body_temp
  latestActivity.value = chartData.value[chartData.value.length - 1]?.activity_level
}

// 渲染图表
const renderCharts = async () => {
  await nextTick()
  renderHeartRateChart()
  renderRespiratoryChart()
  renderTempChart()
  renderActivityChart()
}

const formatTime = (timeStr) => {
  const d = new Date(timeStr)
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}

const renderHeartRateChart = () => {
  if (!heartRateChartRef.value) return
  if (!heartRateChart) heartRateChart = echarts.init(heartRateChartRef.value)

  const times = chartData.value.map(d => formatTime(d.time))
  const data = chartData.value.map(d => d.heart_rate)

  heartRateChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: times,
      axisLabel: { rotate: 45, interval: Math.floor(times.length / 8) }
    },
    yAxis: { type: 'value', name: 'bpm', min: 40, max: 160 },
    series: [{
      type: 'line',
      data,
      smooth: true,
      areaStyle: { opacity: 0.2 },
      lineStyle: { color: '#f53f3f' },
      itemStyle: { color: '#f53f3f' },
      markLine: {
        silent: true,
        data: [
          { yAxis: 60, lineStyle: { color: '#ff7d00', type: 'dashed' }, label: { formatter: '下限' } },
          { yAxis: 140, lineStyle: { color: '#ff7d00', type: 'dashed' }, label: { formatter: '上限' } }
        ]
      }
    }],
    grid: { left: 50, right: 16, top: 20, bottom: 50 }
  })
}

const renderRespiratoryChart = () => {
  if (!respiratoryChartRef.value) return
  if (!respiratoryChart) respiratoryChart = echarts.init(respiratoryChartRef.value)

  const times = chartData.value.map(d => formatTime(d.time))
  const data = chartData.value.map(d => d.respiratory_rate)

  respiratoryChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: times,
      axisLabel: { rotate: 45, interval: Math.floor(times.length / 8) }
    },
    yAxis: { type: 'value', name: '次/分', min: 5, max: 40 },
    series: [{
      type: 'line',
      data,
      smooth: true,
      areaStyle: { opacity: 0.2 },
      lineStyle: { color: '#00b42a' },
      itemStyle: { color: '#00b42a' },
      markLine: {
        silent: true,
        data: [
          { yAxis: 10, lineStyle: { color: '#ff7d00', type: 'dashed' }, label: { formatter: '下限' } },
          { yAxis: 30, lineStyle: { color: '#ff7d00', type: 'dashed' }, label: { formatter: '上限' } }
        ]
      }
    }],
    grid: { left: 50, right: 16, top: 20, bottom: 50 }
  })
}

const renderTempChart = () => {
  if (!tempChartRef.value) return
  if (!tempChart) tempChart = echarts.init(tempChartRef.value)

  const times = chartData.value.map(d => formatTime(d.time))
  const data = chartData.value.map(d => d.body_temp)

  tempChart.setOption({
    tooltip: { trigger: 'axis', formatter: (p) => `${p[0].axisValue}<br/>体温: ${p[0].value}℃` },
    xAxis: {
      type: 'category',
      data: times,
      axisLabel: { rotate: 45, interval: Math.floor(times.length / 8) }
    },
    yAxis: { type: 'value', name: '℃', min: 37, max: 41 },
    series: [{
      type: 'line',
      data,
      smooth: true,
      areaStyle: { opacity: 0.2 },
      lineStyle: { color: '#ff7d00' },
      itemStyle: { color: '#ff7d00' },
      markLine: {
        silent: true,
        data: [
          { yAxis: 38, lineStyle: { color: '#1650ff', type: 'dashed' }, label: { formatter: '下限' } },
          { yAxis: 39.5, lineStyle: { color: '#f53f3f', type: 'dashed' }, label: { formatter: '上限' } }
        ]
      }
    }],
    grid: { left: 50, right: 16, top: 20, bottom: 50 }
  })
}

const renderActivityChart = () => {
  if (!activityChartRef.value) return
  if (!activityChart) activityChart = echarts.init(activityChartRef.value)

  const times = chartData.value.map(d => formatTime(d.time))
  const data = chartData.value.map(d => d.activity_level)

  activityChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: times,
      axisLabel: { rotate: 45, interval: Math.floor(times.length / 8) }
    },
    yAxis: { type: 'value', name: '分', min: 0, max: 100 },
    series: [{
      type: 'bar',
      data,
      itemStyle: {
        color: (params) => {
          const val = params.value
          if (val < 30) return '#00b42a'
          if (val < 70) return '#1650ff'
          return '#ff7d00'
        }
      }
    }],
    grid: { left: 50, right: 16, top: 20, bottom: 50 }
  })
}

const handleResize = () => {
  heartRateChart?.resize()
  respiratoryChart?.resize()
  tempChart?.resize()
  activityChart?.resize()
}

const handlePetChange = () => {
  loadChartData()
}

const goBack = () => {
  router.push('/digital-twin/vitals')
}

onMounted(() => {
  loadPets()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  heartRateChart?.dispose()
  respiratoryChart?.dispose()
  tempChart?.dispose()
  activityChart?.dispose()
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-bar {
  margin-bottom: 12px;
}

.action-bar {
  margin-bottom: 16px;
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

.chart-card {
  border-radius: 8px;
}

.chart-container {
  width: 100%;
  height: 280px;
}

.stats-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
