<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="dashboard-container">

    <a-card class="general-card">
      <template #title>
        <span class="card-title">设备概览</span>
      </template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic
            title="总设备数"
            :value="stats.total_devices"
            :value-style="{ color: '#1650d8' }"
          >
            <template #prefix>
              <icon-apps style="font-size: 22px; margin-right: 6px; color: #1650d8" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic
            title="在线设备"
            :value="stats.online_devices"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <icon-check-circle style="font-size: 22px; margin-right: 6px; color: #52c41a" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic
            title="离线设备"
            :value="stats.offline_devices"
            :value-style="{ color: '#ff4d4f' }"
          >
            <template #prefix>
              <icon-close-circle style="font-size: 22px; margin-right: 6px; color: #ff4d4f" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic
            title="待处理告警"
            :value="stats.pending_alerts"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <icon-exclamation-circle style="font-size: 22px; margin-right: 6px; color: #faad14" />
            </template>
          </a-statistic>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title>
        <span class="card-title">设备状态分布</span>
      </template>
      <a-row :gutter="16">
        <a-col :span="8">
          <div class="online-rate-wrapper">
            <a-progress
              type="circle"
              :percent="onlineRate"
              :width="140"
              :color="onlineRate >= 70 ? '#52c41a' : onlineRate >= 40 ? '#faad14' : '#ff4d4f'"
            >
              <template #formatter>
                <div class="rate-formatter">
                  <span class="rate-value">{{ onlineRate }}%</span>
                  <span class="rate-label">在线率</span>
                </div>
              </template>
            </a-progress>
          </div>
        </a-col>
        <a-col :span="16">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="在线设备">{{ stats.online_devices }}</a-descriptions-item>
            <a-descriptions-item label="离线设备">{{ stats.offline_devices }}</a-descriptions-item>
            <a-descriptions-item label="设备总数">{{ stats.total_devices }}</a-descriptions-item>
            <a-descriptions-item label="告警设备">{{ stats.pending_alerts }}</a-descriptions-item>
          </a-descriptions>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title>
        <span class="card-title">最近告警</span>
      </template>
      <a-table
        :columns="alertColumns"
        :data="recentAlerts"
        :pagination="false"
        :loading="alertsLoading"
        row-key="id"
        size="small"
      >
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">
            {{ getSeverityText(record.severity) }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.status === 1 ? 'processing' : 'default'" />
          {{ record.status === 1 ? '处理中' : '已处理' }}
        </template>
      </a-table>
      <a-empty v-if="recentAlerts.length === 0 && !alertsLoading" description="暂无告警记录" style="margin-top: 16px" />
    </a-card>
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
const alertsLoading = ref(false)

const alertColumns = [
  { title: '告警级别', slotName: 'severity', width: 100 },
  { title: '告警消息', dataIndex: 'message', ellipsis: true },
  { title: '设备ID', dataIndex: 'device_id', width: 160 },
  { title: '发生时间', dataIndex: 'created_at', width: 180 },
  { title: '状态', slotName: 'status', width: 100 }
]

const onlineRate = computed(() => {
  if (stats.value.total_devices === 0) return 0
  return Math.round((stats.value.online_devices / stats.value.total_devices) * 100)
})

const getSeverityColor = (severity) => {
  const colors = { 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }
  return colors[severity] || 'gray'
}

const getSeverityText = (severity) => {
  const texts = { 1: '提示', 2: '一般', 3: '重要', 4: '紧急' }
  return texts[severity] || '未知'
}

const loadStats = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/dashboard/stats', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      stats.value = data.data
    }
  } catch (e) {
    // silent fail, keep defaults
  }
}

const loadAlerts = async () => {
  alertsLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts?page=1&page_size=10', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      recentAlerts.value = data.data.list || []
    }
  } catch (e) {
    recentAlerts.value = []
  } finally {
    alertsLoading.value = false
  }
}

onMounted(() => {
  loadStats()
  loadAlerts()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.general-card {
  border-radius: 8px;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}

.online-rate-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 140px;
}

.rate-formatter {
  display: flex;
  flex-direction: column;
  align-items: center;
  line-height: 1.3;
}

.rate-value {
  font-size: 26px;
  font-weight: 600;
  color: #1d2129;
}

.rate-label {
  font-size: 12px;
  color: #86909c;
}
</style>

