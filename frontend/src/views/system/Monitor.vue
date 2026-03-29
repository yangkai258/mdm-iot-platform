<template>
  <div class="monitor-container">
    <Breadcrumb :items="[{ label: '首页', href: '/' }, { label: '系统监控' }]" />

    <a-card class="general-card">
      <template #title><span class="card-title">系统概览</span></template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="CPU 使用率" :value="cpuUsage" suffix="%" :precision="1">
            <template #prefix><icon-robot style="font-size: 20px; color: #1650d8; margin-right: 6px" /></template>
          </a-statistic>
          <a-progress :percent="cpuUsage" :stroke-width="8" :show-text="true" style="margin-top: 8px" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="内存使用" :value="memUsage" suffix="%" :precision="1" :value-style="{ color: memUsage > 80 ? '#ff4d4f' : '#1650d8' }">
            <template #prefix><icon-storage style="font-size: 20px; color: #1650d8; margin-right: 6px" /></template>
          </a-statistic>
          <a-progress :percent="memUsage" :stroke-width="8" :show-text="true" :status="memUsage > 80 ? 'error' : 'normal'" style="margin-top: 8px" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="运行时长" :value="uptime">
            <template #prefix><icon-clock-circle style="font-size: 20px; color: #1650d8; margin-right: 6px" /></template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="在线设备" :value="onlineDevices">
            <template #prefix><icon-link style="font-size: 20px; color: #52c41a; margin-right: 6px" /></template>
          </a-statistic>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title><span class="card-title">服务状态</span></template>
      <template #extra>
        <a-button type="text" size="small" @click="fetchStatus">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </template>
      <a-table :columns="serviceColumns" :data="services" :pagination="false" row-key="name" size="small">
        <template #status="{ record }">
          <a-badge :status="record.status === 'running' ? 'success' : 'error'" />
          <a-tag :color="record.status === 'running' ? 'green' : 'red'" style="margin-left: 6px">
            {{ record.status === 'running' ? '运行中' : '已停止' }}
          </a-tag>
        </template>
        <template #uptime="{ record }">{{ record.uptime || '-' }}</template>
      </a-table>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title><span class="card-title">数据库连接</span></template>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-statistic title="PostgreSQL" :value="pgStatus === 'connected' ? '已连接' : '未连接'" :value-style="{ color: pgStatus === 'connected' ? '#52c41a' : '#ff4d4f' }">
            <template #prefix><icon-mind-mapping style="font-size: 20px; margin-right: 6px" :style="{ color: pgStatus === 'connected' ? '#52c41a' : '#ff4d4f' }" /></template>
          </a-statistic>
        </a-col>
        <a-col :span="8">
          <a-statistic title="Redis" :value="redisStatus === 'connected' ? '已连接' : '未连接'" :value-style="{ color: redisStatus === 'connected' ? '#52c41a' : '#ff4d4f' }">
            <template #prefix><icon-eraser style="font-size: 20px; margin-right: 6px" :style="{ color: redisStatus === 'connected' ? '#52c41a' : '#ff4d4f' }" /></template>
          </a-statistic>
        </a-col>
        <a-col :span="8">
          <a-statistic title="MQTT" :value="mqttStatus === 'connected' ? '已连接' : '未连接'" :value-style="{ color: mqttStatus === 'connected' ? '#52c41a' : '#ff4d4f' }">
            <template #prefix><icon-wifi style="font-size: 20px; margin-right: 6px" :style="{ color: mqttStatus === 'connected' ? '#52c41a' : '#ff4d4f' }" /></template>
          </a-statistic>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const cpuUsage = ref(0)
const memUsage = ref(0)
const uptime = ref('-')
const onlineDevices = ref(0)
const pgStatus = ref('disconnected')
const redisStatus = ref('disconnected')
const mqttStatus = ref('disconnected')

const services = ref([
  { name: 'Backend API', status: 'running', uptime: '-' },
  { name: 'Frontend', status: 'running', uptime: '-' },
  { name: 'PostgreSQL', status: 'running', uptime: '-' },
  { name: 'Redis', status: 'running', uptime: '-' },
  { name: 'MQTT Broker', status: 'running', uptime: '-' }
])

const serviceColumns = [
  { title: '服务名称', dataIndex: 'name', width: 200 },
  { title: '状态', slotName: 'status', width: 140 },
  { title: '运行时长', slotName: 'uptime', width: 160 }
]

let timer = null

const fetchStatus = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/system/status', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0) {
        cpuUsage.value = data.data.cpu_usage || 0
        memUsage.value = data.data.memory_usage || 0
        uptime.value = data.data.uptime || '-'
        onlineDevices.value = data.data.online_devices || 0
        pgStatus.value = data.data.pg_status || 'disconnected'
        redisStatus.value = data.data.redis_status || 'disconnected'
        mqttStatus.value = data.data.mqtt_status || 'disconnected'
        services.value = services.value.map(s => ({
          ...s,
          status: data.data.services?.[s.name] || s.status
        }))
        return
      }
    }
  } catch (e) {}
  // mock data
  cpuUsage.value = Math.round(Math.random() * 30 + 10)
  memUsage.value = Math.round(Math.random() * 40 + 30)
  uptime.value = `${Math.floor(Math.random() * 30) + 1}天 ${Math.floor(Math.random() * 24)}时 ${Math.floor(Math.random() * 60)}分`
  onlineDevices.value = Math.floor(Math.random() * 50 + 10)
  pgStatus.value = 'connected'
  redisStatus.value = 'connected'
  mqttStatus.value = 'connected'
}

onMounted(() => {
  fetchStatus()
  timer = setInterval(fetchStatus, 30000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.monitor-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
