<template>
    <Breadcrumb :items="['Home','Console','']" />



  <div class="pro-page-container">

    <!-- 面包屑 -->

    <a-breadcrumb class="pro-breadcrumb">

      <a-breadcrumb-item>首页</a-breadcrumb-item>

      <a-breadcrumb-item>设备管理</a-breadcrumb-item>

      <a-breadcrumb-item>监控面板</a-breadcrumb-item>

    </a-breadcrumb>



    <!-- 筛选条件 -->

    <div class="pro-search-bar">

      <a-space>

        <a-select v-model="filterModel" placeholder="设备型号" allow-clear style="width: 160px" @change="loadDashboard">

          <a-option value="M5Stack">M5Stack</a-option>

          <a-option value="ESP32">ESP32</a-option>

          <a-option value="Raspberry Pi">Raspberry Pi</a-option>

        </a-select>

        <a-select v-model="timeRange" placeholder="时间范围" style="width: 130px" @change="loadMetrics">

          <a-option value="24h">近24小时</a-option>

          <a-option value="7d">近7天</a-option>

          <a-option value="30d">近30天</a-option>

        </a-select>

        <a-button @click="loadDashboard">刷新</a-button>

      </a-space>

    </div>



    <!-- 统计卡片 -->

    <a-row :gutter="16" class="stat-row">

      <a-col :span="6">

        <a-card class="stat-card" hoverable>

          <a-statistic :value="summary.total_devices" title="设备总数">

            <template #icon><icon-mobile style="font-size: 28px; color: #1650d8" /></template>

          </a-statistic>

        </a-card>

      </a-col>

      <a-col :span="6">

        <a-card class="stat-card" hoverable>

          <a-statistic :value="summary.online_devices" title="在线设备" :value-from="0" :animation="true">

            <template #icon><icon-check-circle style="font-size: 28px; color: #00b42a" /></template>

          </a-statistic>

        </a-card>

      </a-col>

      <a-col :span="6">

        <a-card class="stat-card" hoverable>

          <a-statistic

            :value="summary.online_rate"

            title="在线率"

            :precision="1"

            suffix="%"

            :value-from="0"

          >

            <template #icon><icon-line-chart style="font-size: 28px; color: #1650d8" /></template>

          </a-statistic>

        </a-card>

      </a-col>

      <a-col :span="6">

        <a-card class="stat-card" hoverable>

          <a-statistic :value="summary.active_alerts" title="活跃告警">

            <template #icon><icon-exclamation-circle style="font-size: 28px; color: #ff4d4f" /></template>

          </a-statistic>

        </a-card>

      </a-col>

    </a-row>



    <!-- 图表区 -->

    <a-row :gutter="16" class="chart-row">

      <a-col :span="12">

        <a-card title="在线设备列表" class="chart-card">

          <template #extra>

            <a-button size="mini" @click="loadDashboard">刷新</a-button>

          </template>

          <a-table

            :columns="onlineColumns"

            :data="onlineDevices"

            :loading="loading"

            :pagination="{ pageSize: 8 }"

            row-key="device_id"

            size="small"

          >

            <template #is_online="{ record }">

              <span class="online-dot" :class="record.is_online ? 'online' : 'offline'"></span>

              {{ record.is_online ? '在线' : '离线' }}

            </template>

            <template #battery_level="{ record }">

              <a-progress :percent="record.battery_level" :stroke-width="6" :show-text="true" size="small"

                v-if="record.battery_level > 0" :color="record.battery_level < 20 ? '#f53f3f' : '#00b42a'" />

              <span v-else>-</span>

            </template>

          </a-table>

        </a-card>

      </a-col>

      <a-col :span="12">

        <a-card title="告警趋势（7天）" class="chart-card">

          <div class="chart-container">

            <div v-if="alertTrend.length === 0" class="no-data">暂无数据</div>

            <div v-else class="bar-chart">

              <div v-for="item in alertTrend" :key="item.date" class="bar-item">

                <div class="bar" :style="{ height: getBarHeight(item.count) + 'px' }" :title="item.count + '条告警'"></div>

                <div class="bar-label">{{ formatDateLabel(item.date) }}</div>

              </div>

            </div>

          </div>

        </a-card>

      </a-col>

    </a-row>



    <!-- 分布图区 -->

    <a-row :gutter="16" class="dist-row">

      <a-col :span="8">

        <a-card title="设备型号分布" class="dist-card">

          <div class="pie-container">

            <div v-if="modelDistribution.length === 0" class="no-data">暂无数据</div>

            <div v-else class="pie-chart">

              <div v-for="(item, idx) in modelDistribution" :key="item.model" class="pie-item">

                <span class="pie-dot" :style="{ background: pieColors[idx % pieColors.length] }"></span>

                <span class="pie-label">{{ item.model }}</span>

                <span class="pie-count">{{ item.count }}</span>

                <span class="pie-percent">({{ getPercent(item.count) }}%)</span>

              </div>

            </div>

          </div>

        </a-card>

      </a-col>

      <a-col :span="8">

        <a-card title="设备状态分布" class="dist-card">

          <div class="pie-container">

            <div v-if="statusDistribution.length === 0" class="no-data">暂无数据</div>

            <div v-else class="pie-chart">

              <div v-for="(item, idx) in statusDistribution" :key="item.status" class="pie-item">

                <span class="pie-dot" :style="{ background: pieColors[idx % pieColors.length] }"></span>

                <span class="pie-label">{{ getStatusText(item.status) }}</span>

                <span class="pie-count">{{ item.count }}</span>

              </div>

            </div>

          </div>

        </a-card>

      </a-col>

      <a-col :span="8">

        <a-card title="告警类型分布" class="dist-card">

          <div class="pie-container">

            <div v-if="alertTypeDistribution.length === 0" class="no-data">暂无数据</div>

            <div v-else class="pie-chart">

              <div v-for="(item, idx) in alertTypeDistribution" :key="item.alert_type" class="pie-item">

                <span class="pie-dot" :style="{ background: pieColors[idx % pieColors.length] }"></span>

                <span class="pie-label">{{ item.alert_type }}</span>

                <span class="pie-count">{{ item.count }}</span>

              </div>

            </div>

          </div>

        </a-card>

      </a-col>

    </a-row>



    <!-- 指标详情 -->

    <a-row :gutter="16" class="metrics-row">

      <a-col :span="12">

        <a-card title="设备固件版本分布" class="dist-card">

          <a-table :columns="firmwareColumns" :data="firmwareDistribution" size="small" :pagination="false" row-key="version">

            <template #version="{ record }">

              <a-tag color="blue">{{ record.version }}</a-tag>

            </template>

      </a-table>

        </a-card>

      </a-col>

      <a-col :span="12">

        <a-card title="运行时长分布" class="dist-card">

          <a-table :columns="uptimeColumns" :data="uptimeDistribution" size="small" :pagination="false" row-key="bucket">

            <template #count="{ record }">

              <a-progress :percent="getUptimePercent(record.count)" :stroke-width="8" :show-text="true" size="small" />

            </template>

      </a-table>

        </a-card>

      </a-col>

    </a-row>

  </div>

</template>



<script setup>

import { ref, reactive, onMounted } from 'vue'

import { Message } from '@arco-design/web-vue'



const loading = ref(false)

const filterModel = ref('')

const timeRange = ref('7d')



const summary = reactive({

  total_devices: 0,

  online_devices: 0,

  offline_devices: 0,

  online_rate: 0,

  active_alerts: 0,

  today_alerts: 0,

  today_resolved: 0

})



const onlineDevices = ref([])

const alertTrend = ref([])

const modelDistribution = ref([])

const statusDistribution = ref([])

const alertTypeDistribution = ref([])

const firmwareDistribution = ref([])

const uptimeDistribution = ref([])



const pieColors = ['#1650d8', '#00b42a', '#ff7d00', '#f53f3f', '#722ed1', '#3491fa']



const onlineColumns = [

  { title: '设备ID', dataIndex: 'device_id', width: 130, ellipsis: true },

  { title: '在线', slotName: 'is_online', width: 80 },

  { title: '最后心跳', dataIndex: 'last_seen', width: 160, ellipsis: true },

  { title: '电量', slotName: 'battery_level', width: 120 },

  { title: '当前模式', dataIndex: 'current_mode', width: 90 }

]



const firmwareColumns = [

  { title: '固件版本', slotName: 'version' },

  { title: '设备数量', dataIndex: 'count' }

]



const uptimeColumns = [

  { title: '运行时长', dataIndex: 'bucket' },

  { title: '设备数量', slotName: 'count' }

]



const getToken = () => localStorage.getItem('token')



const getBarHeight = (count) => {

  if (!alertTrend.value.length) return 0

  const max = Math.max(...alertTrend.value.map(t => t.count))

  return max > 0 ? Math.round((count / max) * 100) : 0

}



const formatDateLabel = (d) => {

  try { return d.substring(5) } catch { return d }

}



const getPercent = (count) => {

  const total = modelDistribution.value.reduce((s, i) => s + i.count, 0)

  return total > 0 ? Math.round((count / total) * 100) : 0

}



const getUptimePercent = (count) => {

  const total = uptimeDistribution.value.reduce((s, i) => s + i.count, 0)

  return total > 0 ? Math.round((count / total) * 100) : 0

}



const getStatusText = (s) => ({ 1: '离线', 2: '在线' }[s] || s)



const loadDashboard = async () => {

  loading.value = true

  try {

    const params = {}

    if (filterModel.value) params.model = filterModel.value



    const res = await fetch(`/api/devices/monitor/dashboard?${new URLSearchParams(params)}`, {

      headers: { 'Authorization': `Bearer ${getToken()}` }

    })

    const json = await res.json()

    if (json.code === 0) {

      const d = json.data

      Object.assign(summary, d.summary || d)

      onlineDevices.value = d.online_devices || []

      alertTrend.value = d.alert_trend || []

      modelDistribution.value = d.model_distribution || []

      statusDistribution.value = d.status_distribution || []

    }

  } catch (e) {

    Message.error('加载监控面板失败')

  } finally {

    loading.value = false

  }

}



const loadMetrics = async () => {

  try {

    const hoursMap = { '24h': 24, '7d': 168, '30d': 720 }

    const hours = hoursMap[timeRange.value] || 168



    const res = await fetch(`/api/devices/monitor/metrics?hours=${hours}`, {

      headers: { 'Authorization': `Bearer ${getToken()}` }

    })

    const json = await res.json()

    if (json.code === 0) {

      const d = json.data

      summary.total_devices = d.total_devices || summary.total_devices

      summary.active_alerts = d.high_severity_alerts || 0

      alertTypeDistribution.value = d.alert_type_distribution || []

      firmwareDistribution.value = d.firmware_distribution || []

      uptimeDistribution.value = d.uptime_distribution || []

    }

  } catch (e) { /* silent */ }

}



onMounted(() => { loadDashboard(); loadMetrics() })

</script>



<style scoped>

.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }

.pro-breadcrumb { margin-bottom: 16px; }

.pro-search-bar { margin-bottom: 12px; }

.stat-row { margin-bottom: 16px; }

.stat-card { border-radius: 8px; }

.chart-row { margin-bottom: 16px; }

.chart-card { border-radius: 8px; }

.dist-row { margin-bottom: 16px; }

.dist-card { border-radius: 8px; }

.metrics-row { margin-bottom: 16px; }

.chart-container { height: 200px; display: flex; align-items: flex-end; justify-content: center; }

.no-data { color: #999; text-align: center; width: 100%; padding: 60px 0; }

.bar-chart { display: flex; align-items: flex-end; gap: 8px; height: 200px; width: 100%; }

.bar-item { flex: 1; display: flex; flex-direction: column; align-items: center; height: 100%; justify-content: flex-end; }

.bar { width: 100%; max-width: 40px; background: linear-gradient(180deg, #3491fa, #1650d8); border-radius: 4px 4px 0 0; min-height: 4px; transition: height 0.3s; cursor: pointer; }

.bar-label { font-size: 10px; color: #999; margin-top: 4px; }

.pie-container { min-height: 120px; }

.pie-chart { display: flex; flex-direction: column; gap: 8px; }

.pie-item { display: flex; align-items: center; gap: 8px; font-size: 13px; }

.pie-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }

.pie-label { flex: 1; color: #333; }

.pie-count { font-weight: 600; color: #333; }

.pie-percent { color: #999; font-size: 12px; }

.online-dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 4px; }

.online-dot.online { background: #00b42a; }

.online-dot.offline { background: #8a8a8a; }

</style>



