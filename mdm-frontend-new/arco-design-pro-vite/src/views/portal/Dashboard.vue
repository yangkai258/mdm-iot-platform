<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="dashboard-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>仪表盘</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 设备统计卡片 -->
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card device-card">
          <a-statistic title="设备总数" :value="deviceStats.total" :value-style="{ color: '#165dff' }">
            <template #prefix><icon-desktop style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">台</span></template>
          </a-statistic>
          <div class="stat-trend up">
            <icon-trending-up /> 较上月 +12%
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card online-card">
          <a-statistic title="在线设备" :value="deviceStats.online" :value-style="{ color: '#52c41a' }">
            <template #prefix><icon-check-circle style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">台</span></template>
          </a-statistic>
          <div class="stat-trend up">
            <icon-trending-up /> 较上月 +8%
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card offline-card">
          <a-statistic title="离线设备" :value="deviceStats.offline" :value-style="{ color: '#ff4d4f' }">
            <template #prefix><icon-close-circle style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">台</span></template>
          </a-statistic>
          <div class="stat-trend down">
            <icon-trending-down /> 较上月 -3%
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card alert-card">
          <a-statistic title="待处理告警" :value="alertStats.pending" :value-style="{ color: '#faad14' }">
            <template #prefix><icon-alert style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">条</span></template>
          </a-statistic>
          <div class="stat-trend warn">
            <icon-alert /> 待处理
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 会员统计卡片 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card member-card">
          <a-statistic title="会员总数" :value="memberStats.total" :value-style="{ color: '#722ed1' }">
            <template #prefix><icon-user-group style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">人</span></template>
          </a-statistic>
          <div class="stat-trend up">
            <icon-trending-up /> 较上月 +15%
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card vip-card">
          <a-statistic title="VIP 会员" :value="memberStats.vip" :value-style="{ color: '#d91ad9' }">
            <template #prefix><icon-star style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">人</span></template>
          </a-statistic>
          <div class="stat-trend up">
            <icon-trending-up /> 较上月 +5%
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card active-card">
          <a-statistic title="活跃会员" :value="memberStats.active" :value-style="{ color: '#0fc6c2' }">
            <template #prefix><icon-activity style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">人</span></template>
          </a-statistic>
          <div class="stat-trend up">
            <icon-trending-up /> 本月
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card new-card">
          <a-statistic title="本月新增" :value="memberStats.newThisMonth" :value-style="{ color: '#165dff' }">
            <template #prefix><icon-user-add style="margin-right: 6px;" /></template>
            <template #suffix><span class="stat-unit">人</span></template>
          </a-statistic>
          <div class="stat-trend up">
            <icon-trending-up /> 日均 +3
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 图表区域 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :xs="24" :lg="12">
        <a-card title="设备状态分布" class="chart-card">
          <div class="chart-placeholder">
            <a-progress type="circle" :percent="deviceOnlineRate" :width="160" :color="deviceOnlineRate > 80 ? '#52c41a' : deviceOnlineRate > 50 ? '#faad14' : '#ff4d4f'">
              <template #formatter>
                <div style="text-align: center;">
                  <div style="font-size: 28px; font-weight: 600;">{{ deviceOnlineRate }}%</div>
                  <div style="font-size: 12px; color: #86909c;">在线率</div>
                </div>
              </template>
            </a-progress>
            <div class="device-dist">
              <div class="dist-item"><span class="dot green"></span>在线 {{ deviceStats.online }} 台</div>
              <div class="dist-item"><span class="dot red"></span>离线 {{ deviceStats.offline }} 台</div>
            </div>
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="12">
        <a-card title="告警级别分布" class="chart-card">
          <div class="chart-placeholder">
            <a-progress type="circle" :percent="alertResolveRate" :width="160" :color="alertResolveRate > 70 ? '#52c41a' : '#faad14'">
              <template #formatter>
                <div style="text-align: center;">
                  <div style="font-size: 28px; font-weight: 600;">{{ alertResolveRate }}%</div>
                  <div style="font-size: 12px; color: #86909c;">解决率</div>
                </div>
              </template>
            </a-progress>
            <div class="device-dist">
              <div class="dist-item"><span class="dot red"></span>紧急 {{ alertStats.critical }}</div>
              <div class="dist-item"><span class="dot orange"></span>警告 {{ alertStats.warning }}</div>
              <div class="dist-item"><span class="dot blue"></span>提示 {{ alertStats.info }}</div>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 详细数据表格 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :span="24">
        <a-card title="会员等级分布" class="table-card">
          <a-table :columns="memberLevelColumns" :data="memberLevelData" :pagination="false" size="small">
            <template #level="{ record }">
              <a-tag :color="record.color">{{ record.level }}</a-tag>
            </template>
      </a-table>
            <template #percent="{ record }">
              <a-progress :percent="record.percent" :color="record.color" size="small" />
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 告警列表 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :span="24">
        <a-card title="最新告警" class="table-card">
          <template #extra>
            <a-link @click="$router.push('/alert')">查看全部</a-link>
          </template>
          <a-table :columns="alertColumns" :data="alertList" :pagination="{ pageSize: 5 }" size="small">
            <template #severity="{ record }">
              <a-tag :color="getSeverityColor(record.severity)">{{ record.severity_text }}</a-tag>
            </template>
      </a-table>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const deviceStats = ref({ total: 128, online: 96, offline: 32 })
const memberStats = ref({ total: 845, vip: 56, active: 312, newThisMonth: 43 })
const alertStats = ref({ pending: 7, critical: 2, warning: 3, info: 8, resolved: 6 })

const deviceOnlineRate = computed(() => {
  if (deviceStats.value.total === 0) return 0
  return Math.round((deviceStats.value.online / deviceStats.value.total) * 100)
})

const alertResolveRate = computed(() => {
  const total = alertStats.value.critical + alertStats.value.warning + alertStats.value.info
  if (total === 0) return 100
  return Math.round((alertStats.value.resolved / total) * 100)
})

const memberLevelColumns = [
  { title: '等级', dataIndex: 'level', slot: 'level' },
  { title: '人数', dataIndex: 'count' },
  { title: '占比', dataIndex: 'percent', slot: 'percent' }
]

const memberLevelData = ref([
  { level: '普通会员', count: 520, percent: 62, color: 'gray' },
  { level: '银卡会员', count: 186, percent: 22, color: '#86909c' },
  { level: '金卡会员', count: 83, percent: 10, color: '#faad14' },
  { level: 'VIP 会员', count: 56, percent: 6, color: '#d91ad9' }
])

const alertColumns = [
  { title: '设备ID', dataIndex: 'device_id' },
  { title: '告警内容', dataIndex: 'message' },
  { title: '级别', dataIndex: 'severity', slot: 'severity' },
  { title: '时间', dataIndex: 'created_at' }
]

const alertList = ref([
  { device_id: 'MDM-001', message: '设备离线超过30分钟', severity: 4, severity_text: '紧急', created_at: '10:30' },
  { device_id: 'MDM-015', message: '电量低于15%', severity: 3, severity_text: '警告', created_at: '09:45' },
  { device_id: 'MDM-008', message: 'OTA 升级失败', severity: 3, severity_text: '警告', created_at: '09:20' },
  { device_id: 'MDM-022', message: '网络信号弱', severity: 2, severity_text: '提示', created_at: '08:55' },
  { device_id: 'MDM-003', message: '设备重启', severity: 2, severity_text: '提示', created_at: '08:30' }
])

const getSeverityColor = (severity) => ({ 4: 'red', 3: 'orange', 2: 'blue', 1: 'green' }[severity] || 'gray')

const loadData = async () => {
  try {
    const token = localStorage.getItem('token')
    // 尝试从后端加载真实数据
    const res = await fetch('/api/v1/dashboard/stats', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (res.ok) {
      const data = await res.json()
      if (data.code === 0) {
        const d = data.data || {}
        deviceStats.value.total = d.device_total || d.total_devices || deviceStats.value.total
        deviceStats.value.online = d.device_online || d.online_devices || deviceStats.value.online
        deviceStats.value.offline = d.device_offline || d.offline_devices || deviceStats.value.offline
      }
    }
  } catch (e) {
    console.warn('加载仪表盘数据失败，使用模拟数据:', e)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.dashboard-container {
  padding: 0;
}

.breadcrumb {
  margin-bottom: 16px;
}

.stat-card {
  border-radius: 8px;
  transition: box-shadow 0.3s;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-unit {
  font-size: 14px;
  color: #86909c;
  margin-left: 4px;
}

.stat-trend {
  margin-top: 8px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.stat-trend.up { color: #52c41a; }
.stat-trend.down { color: #ff4d4f; }
.stat-trend.warn { color: #faad14; }

.chart-card {
  border-radius: 8px;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 0;
  gap: 20px;
}

.device-dist {
  display: flex;
  gap: 24px;
}

.dist-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #4e5969;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
}

.dot.green { background: #52c41a; }
.dot.red { background: #ff4d4f; }
.dot.orange { background: #faad14; }
.dot.blue { background: #1890ff; }

.table-card {
  border-radius: 8px;
}
</style>
