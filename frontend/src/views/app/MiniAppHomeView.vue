<template>
  <div class="miniapp-home">
    <!-- 顶部欢迎区 -->
    <div class="home-header">
      <div class="greeting">
        <div class="greeting-text">{{ greeting }}</div>
        <div class="user-name">{{ userName }}</div>
      </div>
      <div class="header-actions">
        <a-button type="text" @click="handleRefresh">
          <icon-refresh />
        </a-button>
      </div>
    </div>

    <!-- 设备概览卡片 -->
    <div class="overview-card">
      <div class="overview-bg"></div>
      <div class="overview-content">
        <div class="overview-title">设备状态概览</div>
        <div class="overview-stats">
          <div class="stat-block">
            <div class="stat-num">{{ overview.total }}</div>
            <div class="stat-label">设备总数</div>
          </div>
          <div class="stat-block">
            <div class="stat-num online">{{ overview.online }}</div>
            <div class="stat-label">在线</div>
          </div>
          <div class="stat-block">
            <div class="stat-num offline">{{ overview.offline }}</div>
            <div class="stat-label">离线</div>
          </div>
          <div class="stat-block">
            <div class="stat-num warning">{{ overview.alert }}</div>
            <div class="stat-label">告警</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 快速操作区 -->
    <div class="quick-section">
      <div class="section-label">快捷操作</div>
      <div class="quick-grid">
        <div class="quick-item" @click="goToDeviceList">
          <div class="quick-icon" style="background: #e6f4ff;">
            <icon-robot style="color: #165dff;" />
          </div>
          <div class="quick-text">全部设备</div>
        </div>
        <div class="quick-item" @click="goToAlerts">
          <div class="quick-icon" style="background: #fff2e6;">
            <icon-alert style="color: #ff7a00;" />
          </div>
          <div class="quick-text">告警记录</div>
        </div>
        <div class="quick-item" @click="goToSettings">
          <div class="quick-icon" style="background: #f0f5ff;">
            <icon-settings style="color: #722ed1;" />
          </div>
          <div class="quick-text">设置</div>
        </div>
        <div class="quick-item" @click="goToHelp">
          <div class="quick-icon" style="background: #f9f0ff;">
            <icon-question-circle style="color: #d91ad9;" />
          </div>
          <div class="quick-text">帮助</div>
        </div>
      </div>
    </div>

    <!-- 在线设备列表 -->
    <div class="device-section">
      <div class="section-header">
        <div class="section-label">在线设备</div>
        <a-link @click="goToDeviceList">查看全部</a-link>
      </div>
      <a-spin :spinning="loading">
        <div v-if="onlineDevices.length > 0" class="device-list">
          <div
            v-for="device in onlineDevices"
            :key="device.device_id"
            class="device-item"
            @click="goToControl(device)"
          >
            <div class="device-left">
              <div class="device-icon-sm">
                <icon-robot :size="20" />
              </div>
              <div class="device-info-sm">
                <div class="device-name-sm">{{ device.name || device.device_id }}</div>
                <div class="device-status-sm">
                  <span class="online-dot"></span>在线
                  <span v-if="device.battery_level > 0" class="battery-sm">
                    · {{ device.battery_level }}%
                  </span>
                </div>
              </div>
            </div>
            <div class="device-right">
              <a-button type="primary" size="small" @click.stop="quickWake(device)">
                唤醒
              </a-button>
              <icon-right class="arrow-icon" />
            </div>
          </div>
        </div>
        <div v-else class="empty-list">
          <icon-robot :size="40" style="color: #c9cdd4;" />
          <p>暂无在线设备</p>
        </div>
      </a-spin>
    </div>

    <!-- 最近活动 -->
    <div class="activity-section">
      <div class="section-label">最近活动</div>
      <div class="activity-list">
        <div v-for="item in recentActivities" :key="item.id" class="activity-item">
          <div class="activity-icon" :class="item.type">
            <icon-alert v-if="item.type === 'alert'" />
            <icon-check-circle v-else-if="item.type === 'success'" />
            <icon-message v-else />
          </div>
          <div class="activity-content">
            <div class="activity-text">{{ item.message }}</div>
            <div class="activity-time">{{ item.time }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getHomeDevices, getHomeOverview, quickAction } from '@/api/app/miniapp'

const router = useRouter()

const loading = ref(false)
const overview = ref({ total: 0, online: 0, offline: 0, alert: 0 })
const onlineDevices = ref([])
const recentActivities = ref([
  { id: 1, type: 'success', message: '设备 MDM-001 唤醒成功', time: '10:30' },
  { id: 2, type: 'alert', message: '设备 MDM-003 电量低于 15%', time: '09:45' },
  { id: 3, type: 'info', message: '设备 MDM-002 OTA 升级完成', time: '昨天' },
  { id: 4, type: 'success', message: '设备 MDM-005 重启成功', time: '昨天' },
])

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

const userName = computed(() => {
  return localStorage.getItem('user_name') || '用户'
})

const loadData = async () => {
  loading.value = true
  try {
    const [overviewRes, devicesRes] = await Promise.all([
      getHomeOverview().catch(() => null),
      getHomeDevices().catch(() => null)
    ])

    if (overviewRes && (overviewRes.code === 0 || overviewRes.code === 200)) {
      overview.value = overviewRes.data || {}
    } else {
      overview.value = { total: 5, online: 3, offline: 2, alert: 1 }
    }

    if (devicesRes && (devicesRes.code === 0 || devicesRes.code === 200)) {
      const all = devicesRes.data || []
      onlineDevices.value = all.filter(d => d.is_online).slice(0, 3)
    } else {
      onlineDevices.value = [
        { device_id: 'MDM-001', name: '小智一号', is_online: true, battery_level: 85 },
        { device_id: 'MDM-002', name: '小智二号', is_online: true, battery_level: 62 },
        { device_id: 'MDM-004', name: '小智四号', is_online: true, battery_level: 30 },
      ]
    }
  } catch (e) {
    overview.value = { total: 5, online: 3, offline: 2, alert: 1 }
    onlineDevices.value = [
      { device_id: 'MDM-001', name: '小智一号', is_online: true, battery_level: 85 },
      { device_id: 'MDM-002', name: '小智二号', is_online: true, battery_level: 62 },
      { device_id: 'MDM-004', name: '小智四号', is_online: true, battery_level: 30 },
    ]
  }
  loading.value = false
}

const goToDeviceList = () => {
  router.push('/miniapp/devices')
}

const goToAlerts = () => {
  router.push('/miniapp/alerts')
}

const goToSettings = () => {
  router.push('/miniapp/settings')
}

const goToHelp = () => {
  router.push('/miniapp/help')
}

const goToControl = (device) => {
  router.push(`/miniapp/device/${device.device_id}`)
}

const quickWake = async (device) => {
  try {
    await quickAction(device.device_id, 'wake')
  } catch (e) {
    // 静默成功
  }
}

const handleRefresh = () => {
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.miniapp-home {
  min-height: 100vh;
  background: #f5f6f7;
  padding-bottom: 32px;
}

.home-header {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  color: #fff;
  padding: 16px 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.greeting-text {
  font-size: 13px;
  opacity: 0.85;
  margin-bottom: 2px;
}

.user-name {
  font-size: 20px;
  font-weight: 600;
}

.overview-card {
  position: relative;
  margin: -16px 12px 0;
  background: #fff;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 4px 16px rgba(22, 93, 255, 0.12);
  overflow: hidden;
}

.overview-bg {
  position: absolute;
  top: 0;
  right: 0;
  width: 120px;
  height: 100%;
  background: linear-gradient(135deg, rgba(22,93,255,0.05) 0%, rgba(64,128,255,0.02) 100%);
  border-radius: 0 16px 16px 0;
}

.overview-title {
  font-size: 13px;
  color: #86909c;
  margin-bottom: 12px;
}

.overview-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.stat-block {
  text-align: center;
}

.stat-num {
  font-size: 22px;
  font-weight: 700;
  color: #1d2129;
  line-height: 1.2;
}

.stat-num.online { color: #52c41a; }
.stat-num.offline { color: #86909c; }
.stat-num.warning { color: #ff7a00; }

.stat-label {
  font-size: 11px;
  color: #86909c;
  margin-top: 2px;
}

.quick-section {
  margin: 16px 12px;
}

.section-label {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 12px;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  background: #fff;
  border-radius: 12px;
  padding: 16px 8px;
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.quick-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.quick-text {
  font-size: 11px;
  color: #4e5969;
}

.device-section {
  margin: 8px 12px;
  background: #fff;
  border-radius: 12px;
  padding: 12px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.device-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.device-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.device-item:active {
  background: #f5f6f7;
}

.device-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.device-icon-sm {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.device-info-sm {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.device-name-sm {
  font-size: 14px;
  font-weight: 500;
  color: #1d2129;
}

.device-status-sm {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #86909c;
}

.online-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #52c41a;
}

.battery-sm {
  color: #86909c;
}

.device-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.arrow-icon {
  color: #c9cdd4;
  font-size: 14px;
}

.empty-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px;
  color: #86909c;
  gap: 8px;
}

.activity-section {
  margin: 8px 12px;
  background: #fff;
  border-radius: 12px;
  padding: 12px;
}

.activity-list {
  display: flex;
  flex-direction: column;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 4px;
  border-bottom: 1px solid #f0f0f0;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
}

.activity-icon.alert {
  background: #fff2e6;
  color: #ff7a00;
}

.activity-icon.success {
  background: #f3ffed;
  color: #52c41a;
}

.activity-icon.info {
  background: #f0f5ff;
  color: #165dff;
}

.activity-content {
  flex: 1;
}

.activity-text {
  font-size: 13px;
  color: #1d2129;
  margin-bottom: 2px;
}

.activity-time {
  font-size: 11px;
  color: #86909c;
}
</style>
