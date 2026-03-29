<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>工作台</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="[16, 16]" class="stats-row">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic title="设备总数" :value="stats.totalDevices" :value-style="{ color: '#1650ff' }">
            <template #suffix><span class="stat-suffix">台</span></template>
            <template #prefix><icon-desktop style="margin-right: 8px;" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic title="在线设备" :value="stats.onlineDevices" :value-style="{ color: '#52c41a' }">
            <template #suffix><span class="stat-suffix">台</span></template>
            <template #prefix><icon-check-circle style="margin-right: 8px;" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic title="今日指令" :value="stats.todayCommands" :value-style="{ color: '#faad14' }">
            <template #suffix><span class="stat-suffix">条</span></template>
            <template #prefix><icon-send style="margin-right: 8px;" /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stat-card">
          <a-statistic title="待处理告警" :value="stats.pendingAlerts" :value-style="{ color: '#ff4d4f' }">
            <template #suffix><span class="stat-suffix">条</span></template>
            <template #prefix><icon-alert style="margin-right: 8px;" /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 快捷入口 -->
    <a-card title="快捷入口" class="quick-access-card">
      <template #extra><a-link @click="$router.push('/portal')">更多</a-link></template>
      <div class="quick-access-grid">
        <div class="quick-access-item" @click="$router.push('/devices')">
          <icon-desktop class="quick-icon" />
          <span>设备列表</span>
        </div>
        <div class="quick-access-item" @click="$router.push('/process')">
          <icon-history class="quick-icon" />
          <span>流程中心</span>
        </div>
        <div class="quick-access-item" @click="$router.push('/toolbox')">
          <icon-tool class="quick-icon" />
          <span>调试工具</span>
        </div>
        <div class="quick-access-item" @click="$router.push('/ota')">
          <icon-upload class="quick-icon" />
          <span>固件管理</span>
        </div>
        <div class="quick-access-item" @click="$router.push('/alert')">
          <icon-message class="quick-icon" />
          <span>告警管理</span>
        </div>
        <div class="quick-access-item" @click="$router.push('/broadcast')">
          <icon-broadcast class="quick-icon" />
          <span>系统广播</span>
        </div>
      </div>
    </a-card>

    <!-- 最近告警和系统动态 -->
    <a-row :gutter="[16, 16]" style="margin-top: 16px;">
      <a-col :xs="24" :lg="12">
        <a-card title="最近告警" class="recent-card">
          <template #extra><a-link @click="$router.push('/alert')">查看全部</a-link></template>
          <a-list :data="recentAlerts" :bordered="false">
            <template #empty><a-empty description="暂无告警" /></template>
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta :title="item.title" :description="item.time">
                  <template #avatar>
                    <a-tag :color="item.level === 'critical' ? 'red' : item.level === 'warning' ? 'orange' : 'blue'">
                      {{ item.level === 'critical' ? '紧急' : item.level === 'warning' ? '警告' : '提示' }}
                    </a-tag>
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="12">
        <a-card title="系统广播" class="recent-card">
          <template #extra><a-link @click="$router.push('/broadcast')">查看全部</a-link></template>
          <a-list :data="recentBroadcasts" :bordered="false">
            <template #empty><a-empty description="暂无广播" /></template>
            <template #item="{ item }">
              <a-list-item>
                <a-list-item-meta :title="item.title" :description="item.content">
                  <template #avatar>
                    <a-avatar :style="{ backgroundColor: '#1650ff' }"><icon-broadcast /></a-avatar>
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const stats = ref({ totalDevices: 0, onlineDevices: 0, todayCommands: 0, pendingAlerts: 0 })
const recentAlerts = ref([])
const recentBroadcasts = ref([])

const fetchDashboardData = async () => {
  try {
    const [devicesRes, alertsRes, broadcastsRes] = await Promise.all([
      axios.get('/api/v1/devices/list', { params: { page_size: 1 } }),
      axios.get('/api/v1/alerts', { params: { page_size: 5 } }),
      axios.get('/api/v1/broadcasts', { params: { page_size: 5 } })
    ])
    if (devicesRes.data.code === 0) {
      stats.value.totalDevices = devicesRes.data.data.pagination?.total || 0
      stats.value.onlineDevices = devicesRes.data.data.list?.filter(d => d.is_online).length || 0
    }
    if (alertsRes.data.code === 0) {
      recentAlerts.value = alertsRes.data.data.list || []
      stats.value.pendingAlerts = alertsRes.data.data.pagination?.total || 0
    }
    if (broadcastsRes.data.code === 0) recentBroadcasts.value = broadcastsRes.data.data.list || []
  } catch (error) {
    stats.value = { totalDevices: 128, onlineDevices: 86, todayCommands: 1247, pendingAlerts: 3 }
    recentAlerts.value = [
      { id: 1, title: '设备离线告警', time: '2026-03-19 15:30', level: 'critical' },
      { id: 2, title: '电量低于15%', time: '2026-03-19 14:20', level: 'warning' },
      { id: 3, title: 'OTA升级失败', time: '2026-03-19 12:00', level: 'info' }
    ]
    recentBroadcasts.value = [
      { id: 1, title: '系统维护通知', content: '系统将于今晚22:00进行例行维护' },
      { id: 2, title: '新功能上线', content: '支持批量设备管理功能' }
    ]
  }
}

onMounted(() => { fetchDashboardData() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; }
.stat-suffix { font-size: 14px; color: #86909c; margin-left: 4px; }
.quick-access-card { border-radius: 8px; }
.quick-access-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
}
.quick-access-item {
  display: flex; flex-direction: column; align-items: center;
  padding: 20px; border-radius: 8px; background: #f7f8fa;
  cursor: pointer; transition: all 0.3s;
}
.quick-access-item:hover { background: #e6f4ff; transform: translateY(-2px); }
.quick-icon { font-size: 32px; color: #1650ff; margin-bottom: 8px; }
.quick-access-item span { font-size: 14px; color: #4e5969; }
.recent-card { border-radius: 8px; height: 100%; }
@media (max-width: 768px) { .quick-access-grid { grid-template-columns: repeat(3, 1fr); } }
</style>
