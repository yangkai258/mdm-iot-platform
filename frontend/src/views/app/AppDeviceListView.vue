<template>
  <div class="container" style="padding: 0 0 80px; min-height: 100vh; background: #f5f6f7;">
    <Breadcrumb :items="['menu.app', 'menu.app.devices']" />
    <div class="app-device-list">
    <!-- 顶部导航栏 -->
    <div class="app-header">
      <div class="header-title">我的设备</div>
      <div class="header-stats">
        <span class="stat-item online">
          <span class="dot"></span>{{ stats.online }} 在线
        </span>
        <span class="stat-item offline">
          <span class="dot"></span>{{ stats.offline }} 离线
        </span>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-section">
      <a-input-search
        v-model="searchKeyword"
        placeholder="搜索设备名称或ID"
        @search="handleSearch"
        search-button
        class="app-search"
      />
    </div>

    <!-- 设备分类标签 -->
    <div class="filter-tabs">
      <a-radio-group v-model="activeFilter" type="button" @change="handleFilterChange">
        <a-radio value="all">全部</a-radio>
        <a-radio value="online">在线</a-radio>
        <a-radio value="offline">离线</a-radio>
      </a-radio-group>
    </div>

    <!-- 加载状态 -->
    <a-spin :spinning="loading" tip="加载设备中...">
      <!-- 设备网格 -->
      <div v-if="filteredDevices.length > 0" class="device-grid">
        <div
          v-for="device in filteredDevices"
          :key="device.device_id"
          class="device-card"
          :class="{ online: device.is_online, offline: !device.is_online }"
          @click="goToControl(device)"
        >
          <!-- 设备图标 -->
          <div class="device-icon-wrap">
            <div class="device-icon">
              <icon-robot />
            </div>
            <span class="status-badge" :class="device.is_online ? 'badge-online' : 'badge-offline'">
              {{ device.is_online ? '在线' : '离线' }}
            </span>
          </div>

          <!-- 设备信息 -->
          <div class="device-info">
            <div class="device-name">{{ device.name || device.device_id }}</div>
            <div class="device-id">ID: {{ device.device_id }}</div>
          </div>

          <!-- 电量 -->
          <div v-if="device.battery_level > 0" class="device-battery">
            <icon-battery :level="getBatteryLevel(device.battery_level)" />
            <span>{{ device.battery_level }}%</span>
          </div>

          <!-- 快速操作按钮 -->
          <div class="quick-actions" @click.stop>
            <a-button
              v-if="device.is_online"
              type="primary"
              size="mini"
              @click="handleQuickAction(device, 'wake')"
            >
              唤醒
            </a-button>
            <a-button
              v-if="device.is_online"
              size="mini"
              @click="handleQuickAction(device, 'sleep')"
            >
              休眠
            </a-button>
            <a-button
              v-if="device.is_online"
              size="mini"
              @click="handleQuickAction(device, 'reboot')"
            >
              重启
            </a-button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading" class="empty-state">
        <icon-robot :size="64" style="color: #c9cdd4;" />
        <p>暂无设备</p>
        <a-button type="primary" size="small">添加设备</a-button>
      </div>
    </a-spin>

    <!-- Toast 提示 -->
    <a-modal
      v-model:visible="toastVisible"
      :footer="null"
      :header="false"
      :mask-closable="true"
      class="app-toast-modal"
    >
      <div class="toast-content" :class="toastType">
        <icon-check-circle v-if="toastType === 'success'" />
        <icon-close-circle v-else-if="toastType === 'error'" />
        <icon-alert v-else />
        <span>{{ toastMessage }}</span>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDeviceList, quickAction } from '@/api/app'
import Breadcrumb from '@/components/Breadcrumb.vue'

const router = useRouter()

const loading = ref(false)
const devices = ref([])
const searchKeyword = ref('')
const activeFilter = ref('all')

const toastVisible = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

const stats = computed(() => ({
  online: devices.value.filter(d => d.is_online).length,
  offline: devices.value.filter(d => !d.is_online).length
}))

const filteredDevices = computed(() => {
  let result = devices.value

  if (activeFilter.value === 'online') {
    result = result.filter(d => d.is_online)
  } else if (activeFilter.value === 'offline') {
    result = result.filter(d => !d.is_online)
  }

  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(d =>
      (d.name || '').toLowerCase().includes(kw) ||
      (d.device_id || '').toLowerCase().includes(kw)
    )
  }

  return result
})

const getBatteryLevel = (level) => {
  if (level > 75) return 'high'
  if (level > 25) return 'medium'
  return 'low'
}

const showToast = (message, type = 'success') => {
  toastMessage.value = message
  toastType.value = type
  toastVisible.value = true
  setTimeout(() => { toastVisible.value = false }, 2000)
}

const loadDevices = async () => {
  loading.value = true
  try {
    const res = await getDeviceList()
    if (res.code === 0 || res.code === 200) {
      devices.value = res.data || []
    } else {
      // 使用模拟数据
      devices.value = getMockDevices()
    }
  } catch (e) {
    console.warn('加载设备失败，使用模拟数据:', e)
    devices.value = getMockDevices()
  } finally {
    loading.value = false
  }
}

const getMockDevices = () => [
  { device_id: 'MDM-001', name: '小智一号', is_online: true, battery_level: 85 },
  { device_id: 'MDM-002', name: '小智二号', is_online: true, battery_level: 62 },
  { device_id: 'MDM-003', name: '小智三号', is_online: false, battery_level: 0 },
  { device_id: 'MDM-004', name: '小智四号', is_online: true, battery_level: 30 },
  { device_id: 'MDM-005', name: '小智五号', is_online: true, battery_level: 100 },
  { device_id: 'MDM-006', name: '小智六号', is_online: false, battery_level: 0 },
]

const handleSearch = () => {
  // 搜索由 computed 过滤实现
}

const handleFilterChange = () => {
  // 过滤由 computed 实现
}

const goToControl = (device) => {
  router.push(`/app/device/${device.device_id}`)
}

const handleQuickAction = async (device, action) => {
  try {
    await quickAction(device.device_id, action)
    showToast(`${action === 'wake' ? '唤醒' : action === 'sleep' ? '休眠' : '重启'}指令已发送`)
  } catch (e) {
    showToast('指令发送失败', 'error')
  }
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
.app-device-list {
  min-height: 100vh;
  background: #f5f6f7;
  padding: 0 0 80px;
}

.app-header {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  color: #fff;
  padding: 20px 16px 16px;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 8px;
}

.header-stats {
  display: flex;
  gap: 16px;
  font-size: 13px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  opacity: 0.9;
}

.stat-item .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.stat-item.online .dot { background: #52c41a; }
.stat-item.offline .dot { background: #ff4d4f; }

.search-section {
  padding: 12px 16px;
  background: #fff;
}

.app-search {
  width: 100%;
}

.filter-tabs {
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}

.device-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  padding: 12px 16px;
}

.device-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s;
  position: relative;
}

.device-card:active {
  transform: scale(0.98);
}

.device-card.online {
  border-left: 3px solid #52c41a;
}

.device-card.offline {
  border-left: 3px solid #c9cdd4;
  opacity: 0.7;
}

.device-icon-wrap {
  position: relative;
  display: inline-block;
  margin-bottom: 12px;
}

.device-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 24px;
}

.status-badge {
  position: absolute;
  top: -6px;
  right: -6px;
  padding: 2px 6px;
  border-radius: 8px;
  font-size: 10px;
  font-weight: 600;
}

.badge-online {
  background: #52c41a;
  color: #fff;
}

.badge-offline {
  background: #c9cdd4;
  color: #fff;
}

.device-info {
  margin-bottom: 8px;
}

.device-name {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.device-id {
  font-size: 11px;
  color: #86909c;
}

.device-battery {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #4e5969;
  margin-bottom: 8px;
}

.quick-actions {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 0;
  gap: 12px;
  color: #86909c;
}

.toast-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  font-size: 14px;
}

.toast-content.success { color: #52c41a; }
.toast-content.error { color: #ff4d4f; }
.toast-content.warning { color: #faad14; }

.app-toast-modal :deep(.arco-modal) {
  width: auto !important;
  background: transparent;
  box-shadow: none;
}
</style>
