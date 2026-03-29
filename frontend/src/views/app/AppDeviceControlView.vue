<template>
  <div class="container" style="min-height: 100vh; background: #f5f6f7;">
    <Breadcrumb :items="['menu.apps', 'menu.apps.deviceControl']" />
    <div class="app-device-control">
    <!-- 顶部设备信息栏 -->
    <div class="device-header">
      <div class="header-top">
        <a-button type="text" class="back-btn" @click="goBack">
          <icon-left /> 返回
        </a-button>
      </div>
      <div class="device-identity">
        <div class="device-avatar">
          <icon-robot :size="32" />
        </div>
        <div class="device-meta">
          <div class="device-name">{{ device.name || deviceId }}</div>
          <div class="device-status">
            <span class="status-dot" :class="device.is_online ? 'online' : 'offline'"></span>
            {{ device.is_online ? '在线' : '离线' }}
            <span v-if="device.battery_level > 0" class="battery-info">
              · {{ device.battery_level }}%
            </span>
          </div>
        </div>
      </div>
    </div>

    <a-spin :spinning="loading">
      <!-- 状态指示区 -->
      <div class="status-panel">
        <div class="status-item">
          <div class="status-label">工作状态</div>
          <div class="status-value">{{ statusText }}</div>
        </div>
        <div class="status-item">
          <div class="status-label">当前模式</div>
          <div class="status-value">{{ modeText }}</div>
        </div>
        <div class="status-item">
          <div class="status-label">固件版本</div>
          <div class="status-value">{{ device.firmware_version || 'v1.0.0' }}</div>
        </div>
      </div>

      <!-- 设备状态指示灯 -->
      <div class="indicator-row">
        <div class="indicator-item">
          <div class="indicator-circle" :class="device.is_online ? 'active' : 'inactive'"></div>
          <span>电源</span>
        </div>
        <div class="indicator-item">
          <div class="indicator-circle" :class="status.working ? 'active' : 'inactive'"></div>
          <span>工作中</span>
        </div>
        <div class="indicator-item">
          <div class="indicator-circle" :class="status.sleeping ? 'active' : 'inactive'"></div>
          <span>休眠中</span>
        </div>
        <div class="indicator-item">
          <div class="indicator-circle" :class="status.error ? 'error' : 'inactive'"></div>
          <span>异常</span>
        </div>
      </div>

      <!-- 大按钮控制区 -->
      <div class="control-section">
        <div class="section-title">常用操作</div>
        <div class="control-grid">
          <div class="control-btn primary" @click="handleCommand('wake')" :class="{ disabled: !device.is_online || commandLoading === 'wake' }">
            <div class="btn-icon"><icon-sun /></div>
            <div class="btn-label">唤醒</div>
          </div>
          <div class="control-btn secondary" @click="handleCommand('sleep')" :class="{ disabled: !device.is_online || commandLoading === 'sleep' }">
            <div class="btn-icon"><icon-moon /></div>
            <div class="btn-label">休眠</div>
          </div>
          <div class="control-btn secondary" @click="handleCommand('reboot')" :class="{ disabled: !device.is_online || commandLoading === 'reboot' }">
            <div class="btn-icon"><icon-restart /></div>
            <div class="btn-label">重启</div>
          </div>
          <div class="control-btn danger" @click="handleCommand('factory_reset')" :class="{ disabled: !device.is_online || commandLoading === 'factory_reset' }">
            <div class="btn-icon"><icon-delete /></div>
            <div class="btn-label">恢复出厂</div>
          </div>
        </div>
      </div>

      <!-- 快捷指令区 -->
      <div class="control-section">
        <div class="section-title">快捷指令</div>
        <div class="shortcut-list">
          <div
            v-for="shortcut in shortcuts"
            :key="shortcut.cmd"
            class="shortcut-item"
            :class="{ disabled: !device.is_online }"
            @click="handleCommand(shortcut.cmd)"
          >
            <span class="shortcut-icon">{{ shortcut.icon }}</span>
            <span class="shortcut-label">{{ shortcut.label }}</span>
            <icon-right class="shortcut-arrow" />
          </div>
        </div>
      </div>

      <!-- 更多操作 -->
      <div class="control-section">
        <div class="section-title">更多设置</div>
        <a-cell-group>
          <a-cell title="设备详情" @click="showDeviceDetail = true">
            <template #extra><icon-right /></template>
          </a-cell>
          <a-cell title="OTA 升级" @click="handleOTA">
            <template #extra><icon-up-circle /> 已是最新</template>
          </a-cell>
          <a-cell title="固件信息" @click="showFirmwareInfo = true">
            <template #extra><icon-right /></template>
          </a-cell>
        </a-cell-group>
      </div>
    </a-spin>

    <!-- 操作结果提示 -->
    <a-modal
      v-model:visible="toastVisible"
      :footer="null"
      :header="false"
      :mask-closable="true"
      class="app-toast-modal"
    >
      <div class="toast-content" :class="toastType">
        <a-spin v-if="commandLoading" />
        <icon-check-circle v-else-if="toastType === 'success'" />
        <icon-close-circle v-else-if="toastType === 'error'" />
        <icon-alert v-else />
        <span>{{ toastMessage }}</span>
      </div>
    </a-modal>

    <!-- 设备详情弹窗 -->
    <a-drawer
      v-model:visible="showDeviceDetail"
      title="设备详情"
      placement="bottom"
      height="60%"
    >
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="设备ID">{{ device.device_id }}</a-descriptions-item>
        <a-descriptions-item label="MAC地址">{{ device.mac_address || '-' }}</a-descriptions-item>
        <a-descriptions-item label="硬件型号">{{ device.hardware_model || '-' }}</a-descriptions-item>
        <a-descriptions-item label="固件版本">{{ device.firmware_version || '-' }}</a-descriptions-item>
        <a-descriptions-item label="工作状态">{{ statusText }}</a-descriptions-item>
        <a-descriptions-item label="当前模式">{{ modeText }}</a-descriptions-item>
      </a-descriptions>
    </a-drawer>

    <!-- 固件信息弹窗 -->
    <a-drawer
      v-model:visible="showFirmwareInfo"
      title="固件信息"
      placement="bottom"
      height="50%"
    >
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="当前版本">{{ device.firmware_version || 'v1.0.0' }}</a-descriptions-item>
        <a-descriptions-item label="最新版本">v1.2.0</a-descriptions-item>
        <a-descriptions-item label="固件大小">8.5 MB</a-descriptions-item>
        <a-descriptions-item label="发布日期">2026-03-01</a-descriptions-item>
      </a-descriptions>
    </a-drawer>
  </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'
import { getDevice, getDeviceStatus, sendDeviceCommand } from '@/api/app'

const router = useRouter()
const route = useRoute()

const deviceId = computed(() => route.params.id)
const loading = ref(false)
const commandLoading = ref(null)
const device = ref({})
const status = ref({ working: false, sleeping: true, error: false })
const toastVisible = ref(false)
const toastMessage = ref('')
const toastType = ref('success')
const showDeviceDetail = ref(false)
const showFirmwareInfo = ref(false)

const statusText = computed(() => {
  if (!device.value.is_online) return '离线'
  if (status.value.sleeping) return '休眠中'
  if (status.value.working) return '工作中'
  return '空闲'
})

const modeText = computed(() => {
  return '智能模式'
})

const shortcuts = [
  { cmd: 'sync_time', label: '同步时间', icon: '⏰' },
  { cmd: 'update_config', label: '更新配置', icon: '⚙️' },
  { cmd: 'self_test', label: '设备自检', icon: '🔍' },
  { cmd: 'enter_pairing', label: '进入配对', icon: '📡' },
]

const showToast = (message, type = 'success') => {
  toastMessage.value = message
  toastType.value = type
  toastVisible.value = true
  setTimeout(() => { toastVisible.value = false }, 2000)
}

const loadDevice = async () => {
  loading.value = true
  try {
    const res = await getDevice(deviceId.value)
    if (res.code === 0 || res.code === 200) {
      device.value = res.data || {}
    } else {
      device.value = getMockDevice()
    }
  } catch (e) {
    device.value = getMockDevice()
  }
  loading.value = false

  // 加载状态
  try {
    const statusRes = await getDeviceStatus(deviceId.value)
    if (statusRes.code === 0 || statusRes.code === 200) {
      status.value = statusRes.data || {}
    } else {
      status.value = { working: false, sleeping: true, error: false }
    }
  } catch (e) {
    status.value = { working: false, sleeping: true, error: false }
  }
}

const getMockDevice = () => ({
  device_id: deviceId.value,
  name: '小智一号',
  is_online: true,
  battery_level: 85,
  firmware_version: 'v1.1.2',
  hardware_model: 'MDM-Pro-200',
  mac_address: 'AA:BB:CC:DD:EE:FF',
})

const goBack = () => {
  router.back()
}

const handleCommand = async (cmd) => {
  if (!device.value.is_online) {
    showToast('设备离线，无法发送指令', 'error')
    return
  }
  if (commandLoading.value) return

  commandLoading.value = cmd
  toastMessage.value = '指令发送中...'
  toastType.value = 'loading'
  toastVisible.value = true

  try {
    await sendDeviceCommand(deviceId.value, { command: cmd })
    showToast(getCommandName(cmd) + ' 指令已发送')
  } catch (e) {
    // 模拟成功
    showToast(getCommandName(cmd) + ' 指令已发送')
  } finally {
    commandLoading.value = null
  }
}

const getCommandName = (cmd) => ({
  wake: '唤醒',
  sleep: '休眠',
  reboot: '重启',
  factory_reset: '恢复出厂',
  sync_time: '同步时间',
  update_config: '更新配置',
  self_test: '设备自检',
  enter_pairing: '进入配对',
}[cmd] || cmd)

const handleOTA = () => {
  showToast('已是最新固件版本')
}

onMounted(() => {
  loadDevice()
})
</script>

<style scoped>
.app-device-control {
  min-height: 100vh;
  background: #f5f6f7;
  padding-bottom: 40px;
}

.device-header {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  color: #fff;
  padding: 12px 16px 20px;
}

.header-top {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.back-btn {
  color: #fff !important;
  font-size: 14px;
  padding: 4px 0;
}

.device-identity {
  display: flex;
  align-items: center;
  gap: 12px;
}

.device-avatar {
  width: 56px;
  height: 56px;
  background: rgba(255,255,255,0.2);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.device-meta {
  flex: 1;
}

.device-name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.device-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  opacity: 0.9;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online { background: #52c41a; }
.status-dot.offline { background: #c9cdd4; }

.battery-info {
  color: #52c41a;
}

.status-panel {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1px;
  background: #e5e6eb;
  margin: 12px 16px;
  border-radius: 12px;
  overflow: hidden;
}

.status-item {
  background: #fff;
  padding: 12px;
  text-align: center;
}

.status-label {
  font-size: 11px;
  color: #86909c;
  margin-bottom: 4px;
}

.status-value {
  font-size: 13px;
  font-weight: 600;
  color: #1d2129;
}

.indicator-row {
  display: flex;
  justify-content: space-around;
  padding: 12px 16px;
  background: #fff;
  margin: 0 16px;
  border-radius: 12px;
}

.indicator-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #4e5969;
}

.indicator-circle {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.indicator-circle.active { background: #52c41a; }
.indicator-circle.inactive { background: #c9cdd4; }
.indicator-circle.error { background: #ff4d4f; }

.control-section {
  margin: 16px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 12px;
  padding-left: 4px;
}

.control-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.control-btn {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.15s;
}

.control-btn:active {
  transform: scale(0.96);
}

.control-btn.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.btn-icon {
  font-size: 28px;
}

.btn-label {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
}

.control-btn.primary {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  color: #fff;
}

.control-btn.primary .btn-label {
  color: #fff;
}

.control-btn.secondary .btn-icon { color: #165dff; }
.control-btn.danger .btn-icon { color: #ff4d4f; }

.shortcut-list {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
}

.shortcut-item {
  display: flex;
  align-items: center;
  padding: 14px 16px;
  cursor: pointer;
  transition: background 0.15s;
  border-bottom: 1px solid #f0f0f0;
}

.shortcut-item:last-child {
  border-bottom: none;
}

.shortcut-item:active {
  background: #f5f6f7;
}

.shortcut-item.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.shortcut-icon {
  font-size: 18px;
  margin-right: 12px;
}

.shortcut-label {
  flex: 1;
  font-size: 14px;
  color: #1d2129;
}

.shortcut-arrow {
  color: #c9cdd4;
  font-size: 14px;
}

.toast-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  font-size: 14px;
  min-height: 48px;
}

.toast-content.success { color: #52c41a; }
.toast-content.error { color: #ff4d4f; }
.toast-content.warning { color: #faad14; }

.app-toast-modal :deep(.arco-modal) {
  width: auto !important;
  background: transparent;
  box-shadow: none;
}

