<template>
  <div class="container">
    <a-card :bordered="false">
      <template #title>
        <a-space>
          <icon-desktop /> 监控大屏
          <a-button type="text" @click="handleExit">
            <template #icon><icon-close /></template>
            退出全屏
          </a-button>
        </a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="在线设备" :value="stats.online" suffix="/ 100" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="总消息数" :value="stats.messages" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="告警数" :value="stats.alerts" :value-style="{ color: stats.alerts > 0 ? '#F56C6C' : '#67C23A' }" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="系统状态" value="正常" />
        </a-col>
      </a-row>

      <a-row :gutter="16" style="margin-top: 24px">
        <a-col :span="16">
          <a-card :bordered="false" title="设备状态地图">
            <div class="map-placeholder">
              <a-empty description="设备分布地图" />
            </div>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card :bordered="false" title="实时告警">
            <a-timeline>
              <a-timeline-item v-for="alert in alerts" :key="alert.id" :color="getAlertColor(alert.level)">
                <p>{{ alert.message }}</p>
                <span class="time">{{ alert.time }}</span>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16" style="margin-top: 24px">
        <a-col :span="12">
          <a-card :bordered="false" title="消息趋势">
            <a-chart :option="msgChart" style="height: 200px" />
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card :bordered="false" title="设备分布">
            <a-chart :option="deviceChart" style="height: 200px" />
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const stats = reactive({ online: 95, messages: 125680, alerts: 3 })
const alerts = ref([
  { id: 1, message: '设备D001温度异常', level: 'warning', time: '10:00' }
])

const msgChart = reactive({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: ['10:00', '10:05', '10:10', '10:15', '10:20'] },
  yAxis: { type: 'value' },
  series: [{ type: 'line', smooth: true, data: [120, 135, 128, 142, 138] }]
})

const deviceChart = reactive({
  tooltip: { trigger: 'item' },
  series: [{ type: 'pie', radius: '60%', data: [
    { name: '在线', value: 95 }, { name: '离线', value: 5 }
  ]}]
})

const getAlertColor = (level) => ({ warning: 'orange', error: 'red', info: 'blue' }[level] || 'gray')
const handleExit = () => { }
</script>

<style scoped>
.container { padding: 16px; background: #0d1117; min-height: 100vh; }
.map-placeholder { height: 300px; display: flex; align-items: center; justify-content: center; background: #161b22; border-radius: 8px; }
.time { font-size: 12px; color: #909399; }
</style>
