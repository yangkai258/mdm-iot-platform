<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>第三方集成</a-breadcrumb-item>
      <a-breadcrumb-item>智能家居</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 + 操作按钮 -->
    <div class="page-header">
      <div class="header-left">
        <h2>智能家居控制</h2>
      </div>
      <div class="header-right">
        <a-button type="primary" @click="loadDevices">
          <template #icon><icon-refresh /></template>
          刷新设备
        </a-button>
      </div>
    </div>

    <!-- 搜索筛选栏 -->
    <div class="filter-bar">
      <a-space wrap>
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索设备名称..."
          style="width: 240px"
          search-button
          @search="loadDevices"
        />
        <a-select
          v-model="filterRoom"
          placeholder="房间"
          style="width: 140px"
          allow-clear
          @change="loadDevices"
        >
          <a-option value="living">客厅</a-option>
          <a-option value="bedroom">卧室</a-option>
          <a-option value="kitchen">厨房</a-option>
          <a-option value="bathroom">卫生间</a-option>
          <a-option value="study">书房</a-option>
        </a-select>
        <a-select
          v-model="filterStatus"
          placeholder="设备状态"
          style="width: 120px"
          allow-clear
          @change="loadDevices"
        >
          <a-option value="online">在线</a-option>
          <a-option value="offline">离线</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 设备列表 -->
    <div class="device-list">
      <div v-if="loading" class="loading-state">
        <a-spin size="large" />
      </div>
      <div v-else-if="filteredDevices.length === 0" class="empty-state">
        <icon-drive-info class="empty-icon" />
        <p>暂无设备，请先在米家 App 中添加设备</p>
      </div>
      <div v-else class="device-grid">
        <div v-for="device in filteredDevices" :key="device.id" class="device-card">
          <div class="device-card-header">
            <div class="device-icon-wrapper" :class="getDeviceIconClass(device.device_type)">
              <component :is="getDeviceIcon(device.device_type)" :size="28" />
            </div>
            <div class="device-info">
              <div class="device-name-row">
                <span class="device-name">{{ device.name }}</span>
                <span class="device-status-badge" :class="device.online ? 'status-online' : 'status-offline'">
                  {{ device.online ? '在线' : '离线' }}
                </span>
              </div>
              <div class="device-meta">
                <span>{{ getDeviceTypeName(device.device_type) }}</span>
                <span>{{ getRoomName(device.room) }}</span>
              </div>
            </div>
          </div>

          <!-- 设备状态显示 -->
          <div class="device-status-area">
            <template v-if="device.device_type === 'light'">
              <div class="status-row">
                <span class="status-label">亮度</span>
                <a-slider
                  v-model="device.brightness"
                  :min="0"
                  :max="100"
                  :disabled="!device.online"
                  class="device-slider"
                  @change="(val) => updateDeviceState(device, 'brightness', val)"
                />
                <span class="status-value">{{ device.brightness }}%</span>
              </div>
              <div class="status-row">
                <span class="status-label">开关</span>
                <a-switch
                  v-model="device.power"
                  :disabled="!device.online"
                  size="small"
                  @change="(val) => updateDeviceState(device, 'power', val)"
                />
              </div>
            </template>
            <template v-else-if="device.device_type === 'ac'">
              <div class="status-row">
                <span class="status-label">温度</span>
                <a-input-number
                  v-model="device.target_temp"
                  :min="16"
                  :max="30"
                  :step="1"
                  :disabled="!device.online"
                  size="small"
                  @change="(val) => updateDeviceState(device, 'target_temp', val)"
                />
                <span class="status-value">°C</span>
              </div>
              <div class="status-row">
                <span class="status-label">模式</span>
                <a-select
                  v-model="device.mode"
                  size="small"
                  :disabled="!device.online"
                  style="width: 100px"
                  @change="(val) => updateDeviceState(device, 'mode', val)"
                >
                  <a-option value="cool">制冷</a-option>
                  <a-option value="heat">制热</a-option>
                  <a-option value="fan">送风</a-option>
                  <a-option value="auto">自动</a-option>
                </a-select>
              </div>
            </template>
            <template v-else-if="device.device_type === 'switch'">
              <div class="status-row">
                <span class="status-label">开关</span>
                <a-switch
                  v-model="device.power"
                  :disabled="!device.online"
                  @change="(val) => updateDeviceState(device, 'power', val)"
                />
              </div>
            </template>
            <template v-else-if="device.device_type === 'camera'">
              <div class="status-row">
                <span class="status-label">录像</span>
                <a-switch
                  v-model="device.recording"
                  :disabled="!device.online"
                  size="small"
                  @change="(val) => updateDeviceState(device, 'recording', val)"
                />
              </div>
              <div class="camera-preview">
                <icon-video-camera :size="32" style="color: #999" />
                <span>画面预览</span>
              </div>
            </template>
            <template v-else>
              <div class="status-text">
                <span>{{ device.online ? '设备正常' : '设备离线' }}</span>
              </div>
            </template>
          </div>

          <!-- 联动触发按钮 -->
          <div class="device-actions">
            <a-button size="small" type="primary" :disabled="!device.online" @click="showLinkageDialog(device)">
              <template #icon><icon-link /></template>
              联动设置
            </a-button>
            <a-button size="small" :disabled="!device.online" @click="quickAction(device)">
              {{ device.power ? '关闭' : '打开' }}
            </a-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 联动设置弹窗 -->
    <a-modal
      v-model:visible="linkageModalVisible"
      title="设备联动设置"
      :width="560"
      @ok="handleSaveLinkage"
      @cancel="linkageModalVisible = false"
    >
      <a-form :model="linkageForm" layout="vertical">
        <a-form-item label="联动名称" required>
          <a-input v-model="linkageForm.name" placeholder="如: 回家自动开灯" />
        </a-form-item>
        <a-form-item label="触发条件" required>
          <a-select v-model="linkageForm.trigger_type" placeholder="选择触发类型" style="width: 100%">
            <a-option value="time">定时触发</a-option>
            <a-option value="device_state">设备状态变化</a-option>
            <a-option value="location">位置变化</a-option>
            <a-option value="pet_status">宠物状态变化</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="linkageForm.trigger_type === 'time'" label="触发时间">
          <a-time-picker v-model="linkageForm.trigger_time" style="width: 100%" format="HH:mm" />
        </a-form-item>
        <a-form-item label="执行动作">
          <a-textarea v-model="linkageForm.action" placeholder="如: 打开客厅灯，亮度调至80%" :rows="3" />
        </a-form-item>
        <a-form-item label="关联宠物">
          <a-select v-model="linkageForm.pet_id" placeholder="选择关联宠物（可选）" style="width: 100%" allow-clear>
            <a-option v-for="pet in petList" :key="pet.pet_id" :value="pet.pet_id">{{ pet.pet_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="启用状态">
          <a-switch v-model="linkageForm.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  IconRefresh,
  IconDriveFile,
  IconSunFill,
  IconHome,
  IconVideoCamera,
  IconPoweroff,
  IconLink,
  IconClose
} from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const devices = ref([])
const searchKeyword = ref('')
const filterRoom = ref('')
const filterStatus = ref('')
const linkageModalVisible = ref(false)
const currentDevice = ref(null)
const petList = ref([])

const linkageForm = reactive({
  name: '',
  trigger_type: '',
  trigger_time: null,
  action: '',
  pet_id: null,
  enabled: true
})

const filteredDevices = computed(() => {
  let result = devices.value
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(d => d.name.toLowerCase().includes(kw))
  }
  if (filterRoom.value) {
    result = result.filter(d => d.room === filterRoom.value)
  }
  if (filterStatus.value) {
    const isOnline = filterStatus.value === 'online'
    result = result.filter(d => d.online === isOnline)
  }
  return result
})

function getDeviceIconClass(type) {
  const map = {
    light: 'icon-light',
    ac: 'icon-ac',
    switch: 'icon-switch',
    camera: 'icon-camera'
  }
  return map[type] || 'icon-default'
}

function getDeviceIcon(type) {
  const map = {
    light: IconSunFill,
    ac: IconSunFill,
    switch: IconPoweroff,
    camera: IconVideoCamera
  }
  return map[type] || IconHome
}

function getDeviceTypeName(type) {
  const map = {
    light: '智能灯',
    ac: '空调',
    switch: '智能开关',
    camera: '摄像头',
    sensor: '传感器',
    curtain: '窗帘电机',
    lock: '智能门锁'
  }
  return map[type] || type
}

function getRoomName(room) {
  const map = {
    living: '客厅',
    bedroom: '卧室',
    kitchen: '厨房',
    bathroom: '卫生间',
    study: '书房'
  }
  return map[room] || room || '未分配房间'
}

async function loadDevices() {
  loading.value = true
  try {
    const res = await fetch('/api/v1/integrations/smarthome/devices', {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` }
    })
    const data = await res.json()
    if (data.data) {
      devices.value = data.data
    } else {
      loadMockDevices()
    }
  } catch {
    loadMockDevices()
  } finally {
    loading.value = false
  }
}

function loadMockDevices() {
  devices.value = [
    { id: '1', name: '客厅主灯', device_type: 'light', room: 'living', online: true, power: true, brightness: 80 },
    { id: '2', name: '卧室空调', device_type: 'ac', room: 'bedroom', online: true, power: true, target_temp: 24, mode: 'cool' },
    { id: '3', name: '玄关开关', device_type: 'switch', room: 'living', online: true, power: false },
    { id: '4', name: '客厅摄像头', device_type: 'camera', room: 'living', online: true, recording: false },
    { id: '5', name: '厨房灯', device_type: 'light', room: 'kitchen', online: false, power: false, brightness: 0 },
    { id: '6', name: '书房空调', device_type: 'ac', room: 'study', online: true, power: false, target_temp: 26, mode: 'heat' }
  ]
}

async function updateDeviceState(device, key, value) {
  try {
    await fetch(`/api/v1/integrations/smarthome/devices/${device.id}/state`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token') || ''}`
      },
      body: JSON.stringify({ [key]: value })
    })
    device[key] = value
    Message.success('设备状态已更新')
  } catch {
    Message.error('更新失败，请重试')
  }
}

function quickAction(device) {
  const newPower = !device.power
  updateDeviceState(device, 'power', newPower)
}

function showLinkageDialog(device) {
  currentDevice.value = device
  linkageForm.name = `${device.name} 联动`
  linkageForm.trigger_type = ''
  linkageForm.trigger_time = null
  linkageForm.action = ''
  linkageForm.pet_id = null
  linkageForm.enabled = true
  linkageModalVisible.value = true
}

async function handleSaveLinkage() {
  if (!linkageForm.name || !linkageForm.trigger_type) {
    Message.warning('请填写完整的联动信息')
    return
  }
  try {
    await fetch('/api/v1/integrations/smarthome/linkages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token') || ''}`
      },
      body: JSON.stringify({
        device_id: currentDevice.value.id,
        ...linkageForm
      })
    })
    Message.success('联动设置已保存')
    linkageModalVisible.value = false
  } catch {
    Message.error('保存失败，请重试')
  }
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  min-height: 100vh;
  background: #f5f6f7;
}

.breadcrumb {
  margin-bottom: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
}

.filter-bar {
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.device-list {
  min-height: 200px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  background: #fff;
  border-radius: 8px;
}

.empty-icon {
  font-size: 48px;
  color: #ccc;
  margin-bottom: 16px;
}

.empty-state p {
  color: #999;
  font-size: 14px;
}

.device-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.device-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  transition: box-shadow 0.2s;
}

.device-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.device-card-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.device-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon-light { background: #fffbe6; color: #feac00; }
.icon-ac { background: #e6f7ff; color: #0fc6c8; }
.icon-switch { background: #f0f5ff; color: #1650d8; }
.icon-camera { background: #fff1f0; color: #ff4d4f; }
.icon-default { background: #f7f8fa; color: #86909c; }

.device-info {
  flex: 1;
  min-width: 0;
}

.device-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.device-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-1);
}

.device-status-badge {
  font-size: 12px;
  padding: 1px 6px;
  border-radius: 4px;
}

.status-online {
  background: #e6fff4;
  color: #00b42a;
}

.status-offline {
  background: #f5f5f5;
  color: #86909c;
}

.device-meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: #86909c;
}

.device-status-area {
  background: #f7f8fa;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 12px;
}

.status-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.status-row:last-child {
  margin-bottom: 0;
}

.status-label {
  font-size: 13px;
  color: #4e5969;
  width: 48px;
}

.status-value {
  font-size: 13px;
  color: #4e5969;
  min-width: 36px;
}

.device-slider {
  flex: 1;
}

.camera-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 80px;
  background: #e6e8eb;
  border-radius: 6px;
  gap: 4px;
}

.camera-preview span {
  font-size: 12px;
  color: #86909c;
}

.status-text {
  text-align: center;
  padding: 8px;
  color: #86909c;
  font-size: 13px;
}

.device-actions {
  display: flex;
  gap: 8px;
}
</style>
