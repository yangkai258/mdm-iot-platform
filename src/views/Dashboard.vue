<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>工作台</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
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

    <!-- 数据展示区 -->
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
            <template #empty>
              <a-empty description="暂无告警" />
            </template>
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
    const res = await fetch('/api/v1/dashboard/stats', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      const d = data.data || {}
      // 兼容后端可能返回的不同字段名
      stats.value = {
        total_devices: d.total_devices ?? d.total ?? d.device_total ?? 0,
        online_devices: d.online_devices ?? d.online ?? d.device_online ?? 0,
        offline_devices: d.offline_devices ?? d.offline ?? d.device_offline ?? 0,
        total_alerts: d.total_alerts ?? d.alerts_total ?? d.alert_total ?? 0,
        pending_alerts: d.pending_alerts ?? d.pending ?? d.alerts_pending ?? 0
      }
    }

    // Load recent alerts
    const alertRes = await fetch('/api/v1/alerts?status=1', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const alertData = await alertRes.json()
    if (alertData.code === 0) {
      recentAlerts.value = (alertData.data?.list || alertData.data || []).slice(0, 5)
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
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.stats-row {
  margin-bottom: 16px;
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
