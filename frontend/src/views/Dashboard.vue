<template>
  <div class="dashboard-container">
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="总设备数" :value="stats.total_devices">
            <template #prefix>📱</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="在线设备" :value="stats.online_devices" :value-style="{ color: '#52c41a' }">
            <template #prefix>🟢</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="离线设备" :value="stats.offline_devices" :value-style="{ color: '#ff4d4f' }">
            <template #prefix>🔴</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="待处理告警" :value="stats.pending_alerts" :value-style="{ color: '#faad14' }">
            <template #prefix>⚠️</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px;">
      <a-col :span="12">
        <a-card title="设备状态分布">
          <div class="chart-placeholder">
            <a-progress type="circle" :percent="onlineRate" :width="150">
              <template #formatter>
                <span style="font-size: 24px;">{{ onlineRate }}%</span>
              </template>
            </a-progress>
            <p>设备在线率</p>
          </div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="最近告警">
          <a-list :data-source="recentAlerts">
            <a-list-item v-for="alert in recentAlerts" :key="alert.id">
              <a-list-item-meta>
                <template #avatar>
                  <a-avatar :style="{ backgroundColor: getSeverityColor(alert.severity) }">
                    {{ alert.severity }}
                  </a-avatar>
                </template>
                <template #title>{{ alert.message }}</template>
                <template #description>{{ alert.device_id }} - {{ alert.created_at }}</template>
              </a-list-item-meta>
            </a-list-item>
            <a-empty v-if="recentAlerts.length === 0" description="暂无告警" />
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const stats = ref({
  total_devices: 0,
  online_devices: 0,
  offline_devices: 0,
  total_alerts: 0,
  pending_alerts: 0
})

const recentAlerts = ref([])

const onlineRate = computed(() => {
  if (stats.value.total_devices === 0) return 0
  return Math.round((stats.value.online_devices / stats.value.total_devices) * 100)
})

const getSeverityColor = (severity) => {
  const colors = { 1: '#52c41a', 2: '#1890ff', 3: '#faad14', 4: '#ff4d4f' }
  return colors[severity] || '#999'
}

const loadStats = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('http://localhost:8080/api/v1/dashboard/stats', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      stats.value = data.data
    }
    
    // Load recent alerts
    const alertRes = await fetch('http://localhost:8080/api/v1/alerts?status=1', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const alertData = await alertRes.json()
    if (alertData.code === 0) {
      recentAlerts.value = (alertData.data.list || []).slice(0, 5)
    }
  } catch (e) {
    console.error('加载统计失败:', e)
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.dashboard-container {
  padding: 16px;
}

.stat-card {
  text-align: center;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
}
</style>
