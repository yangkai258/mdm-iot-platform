<template>
  <div class="miniapp-device">
    <!-- 顶部设备信息 -->
    <div class="device-header">
      <a-button type="text" class="back-btn" @click="goBack">
        <icon-left /> 返回
      </a-button>
      <div class="header-title">{{ device.name || deviceId }}</div>
      <div class="header-right">
        <icon-more @click="showMore = true" />
      </div>
    </div>

    <!-- 设备状态卡片 -->
    <div class="status-card">
      <div class="status-main">
        <div class="device-avatar" :class="device.is_online ? 'online' : 'offline'">
          <icon-robot :size="36" />
        </div>
        <div class="status-info">
          <div class="status-text">
            <span class="status-indicator" :class="device.is_online ? 'online' : 'offline'"></span>
            {{ device.is_online ? '在线' : '离线' }}
          </div>
          <div class="status-mode">{{ modeText }}</div>
        </div>
      </div>
      <div v-if="device.battery_level > 0" class="battery-bar">
        <div class="battery-label">
          <icon-battery /> 电量
        </div>
        <a-progress
          :percent="device.battery_level"
          :color="device.battery_level > 20 ? '#52c41a' : '#ff4d4f'"
          size="small"
          :show-text="false"
        />
        <span class="battery-num">{{ device.battery_level }}%</span>
      </div>
    </div>

    <a-spin :spinning="loading">
      <!-- 语音指令入口 -->
      <div class="voice-section">
        <div class="voice-btn" @click="handleVoiceInput" :class="{ recording: isRecording }">
          <div class="voice-icon">
            <icon-mic v-if="!isRecording" />
            <icon-record v-else />
          </div>
          <div class="voice-text">
            {{ isRecording ? '录音中...' : '按住说话' }}
          </div>
        </div>
        <div class="voice-hint">
          <icon-info-circle /> 支持语音指令控制设备
        </div>
      </div>

      <!-- 常用控制按钮 -->
      <div class="control-section">
        <div class="section-title">常用操作</div>
        <div class="action-grid">
          <div
            v-for="action in commonActions"
            :key="action.cmd"
            class="action-item"
            :class="{ disabled: !device.is_online || commandLoading === action.cmd }"
            @click="handleAction(action.cmd)"
          >
            <div class="action-icon-wrap" :style="{ background: action.bg }">
              <span class="action-icon">{{ action.icon }}</span>
            </div>
            <div class="action-label">{{ action.label }}</div>
          </div>
        </div>
      </div>

      <!-- 设备状态详情 -->
      <div class="info-section">
        <div class="section-title">设备信息</div>
        <div class="info-card">
          <div class="info-row">
            <span class="info-label">设备ID</span>
            <span class="info-value">{{ device.device_id || deviceId }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">固件版本</span>
            <span class="info-value">{{ device.firmware_version || 'v1.0.0' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">硬件型号</span>
            <span class="info-value">{{ device.hardware_model || 'MDM-Pro-200' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">运行时间</span>
            <span class="info-value">{{ uptime }}</span>
          </div>
        </div>
      </div>

      <!-- 最近指令记录 -->
      <div class="log-section">
        <div class="section-title">最近指令</div>
        <div class="log-list">
          <div v-for="log in commandLogs" :key="log.id" class="log-item">
            <div class="log-icon" :class="log.result">
              <icon-check-circle v-if="log.result === 'success'" />
              <icon-close-circle v-else-if="log.result === 'failed'" />
              <icon-clock v-else />
            </div>
            <div class="log-content">
              <div class="log-cmd">{{ log.command }}</div>
              <div class="log-time">{{ log.time }}</div>
            </div>
            <div class="log-status" :class="log.result">
              {{ log.result === 'success' ? '成功' : log.result === 'failed' ? '失败' : '发送中' }}
            </div>
          </div>
        </div>
      </div>
    </a-spin>

    <!-- 语音输入弹窗 -->
    <a-modal
      v-model:visible="voiceModalVisible"
      :footer="null"
      title="语音指令"
      @cancel="cancelVoice"
    >
      <div class="voice-modal-content">
        <div class="voice-wave" :class="{ active: isRecording }">
          <div v-for="i in 5" :key="i" class="wave-bar"></div>
        </div>
        <div class="voice-hint-text">
          {{ isRecording ? '正在识别...' : '点击开始录音' }}
        </div>
        <div class="voice-actions">
          <a-button @click="cancelVoice">取消</a-button>
          <a-button type="primary" @click="toggleRecording">
            {{ isRecording ? '识别' : '开始' }}
          </a-button>
        </div>
        <div v-if="recognizedText" class="recognized-text">
          <icon-check-circle style="color: #52c41a;" /> 识别结果：{{ recognizedText }}
        </div>
      </div>
    </a-modal>

    <!-- Toast -->
    <a-modal
      v-model:visible="toastVisible"
      :footer="null"
      :header="false"
      :mask-closable="true"
      class="mini-toast-modal"
    >
      <div class="toast-content" :class="toastType">
        <a-spin v-if="commandLoading" />
        <icon-check-circle v-else-if="toastType === 'success'" />
        <icon-close-circle v-else-if="toastType === 'error'" />
        <span>{{ toastMessage }}</span>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getDeviceStatus, controlDevice, sendVoiceCommand } from '@/api/app/miniapp'

const router = useRouter()
const route = useRoute()

const deviceId = computed(() => route.params.id)
const loading = ref(false)
const commandLoading = ref(null)
const device = ref({})
const uptime = ref('0天 0小时')
const isRecording = ref(false)
const voiceModalVisible = ref(false)
const recognizedText = ref('')
const toastVisible = ref(false)
const toastMessage = ref('')
const toastType = ref('success')
const showMore = ref(false)
const commandLogs = ref([
  { id: 1, command: '唤醒设备', time: '今天 10:30', result: 'success' },
  { id: 2, command: '同步时间', time: '今天 09:15', result: 'success' },
  { id: 3, command: '重启设备', time: '昨天 18:00', result: 'success' },
  { id: 4, command: '恢复出厂', time: '昨天 14:20', result: 'failed' },
])

const modeText = computed(() => {
  if (!device.value.is_online) return '离线'
  return '智能模式'
})

const commonActions = [
  { cmd: 'wake', label: '唤醒', icon: '🌞', bg: '#fff2e6' },
  { cmd: 'sleep', label: '休眠', icon: '🌙', bg: '#f0f5ff' },
  { cmd: 'reboot', label: '重启', icon: '🔄', bg: '#f9f0ff' },
  { cmd: 'self_test', label: '自检', icon: '🔍', bg: '#f0f5ff' },
  { cmd: 'sync_time', label: '对时', icon: '⏰', bg: '#e6fffb' },
  { cmd: 'config', label: '配置', icon: '⚙️', bg: '#fff2e6' },
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
    const res = await getDeviceStatus(deviceId.value)
    if (res.code === 0 || res.code === 200) {
      device.value = res.data || {}
    } else {
      device.value = getMockDevice()
    }
  } catch (e) {
    device.value = getMockDevice()
  }
  loading.value = false
}

const getMockDevice = () => ({
  device_id: deviceId.value,
  name: '小智一号',
  is_online: true,
  battery_level: 85,
  firmware_version: 'v1.1.2',
  hardware_model: 'MDM-Pro-200',
})

const goBack = () => {
  router.back()
}

const handleAction = async (cmd) => {
  if (!device.value.is_online) {
    showToast('设备离线，无法发送指令', 'error')
    return
  }
  if (commandLoading.value) return

  commandLoading.value = cmd
  toastMessage.value = '发送中...'
  toastType.value = 'loading'
  toastVisible.value = true

  try {
    await controlDevice(deviceId.value, { command: cmd })
    showToast(getActionName(cmd) + ' 指令已发送')

    // 添加到日志
    commandLogs.value.unshift({
      id: Date.now(),
      command: getActionName(cmd),
      time: '刚刚',
      result: 'success'
    })
  } catch (e) {
    showToast('指令发送失败', 'error')
  } finally {
    commandLoading.value = null
  }
}

const getActionName = (cmd) => ({
  wake: '唤醒',
  sleep: '休眠',
  reboot: '重启',
  self_test: '设备自检',
  sync_time: '同步时间',
  config: '更新配置',
}[cmd] || cmd)

const handleVoiceInput = () => {
  voiceModalVisible.value = true
  isRecording.value = false
  recognizedText.value = ''
}

const toggleRecording = () => {
  isRecording.value = !isRecording.value
  if (!isRecording.value && recognizedText.value === '') {
    // 模拟识别结果
    recognizedText.value = '唤醒设备'
  }
}

const cancelVoice = () => {
  isRecording.value = false
  voiceModalVisible.value = false
  recognizedText.value = ''
}

onMounted(() => {
  loadDevice()
})
</script>

<style scoped>
.miniapp-device {
  min-height: 100vh;
  background: #f5f6f7;
  padding-bottom: 32px;
}

.device-header {
  background: #fff;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid #f0f0f0;
  position: sticky;
  top: 0;
  z-index: 10;
}

.back-btn {
  color: #165dff;
  padding: 4px 0;
}

.header-title {
  flex: 1;
  text-align: center;
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.header-right {
  color: #4e5969;
  font-size: 18px;
  cursor: pointer;
}

.status-card {
  background: #fff;
  margin: 12px;
  border-radius: 12px;
  padding: 16px;
}

.status-main {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 14px;
}

.device-avatar {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.device-avatar.online {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
}

.device-avatar.offline {
  background: #c9cdd4;
}

.status-info {
  flex: 1;
}

.status-text {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 4px;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.online { background: #52c41a; }
.status-indicator.offline { background: #c9cdd4; }

.status-mode {
  font-size: 13px;
  color: #86909c;
}

.battery-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
}

.battery-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #86909c;
  white-space: nowrap;
}

.battery-num {
  font-size: 12px;
  color: #86909c;
  white-space: nowrap;
}

.voice-section {
  padding: 16px 12px;
  text-align: center;
}

.voice-btn {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 20px 40px;
  background: linear-gradient(135deg, #165dff 0%, #4080ff 100%);
  border-radius: 16px;
  cursor: pointer;
  transition: transform 0.15s;
}

.voice-btn:active {
  transform: scale(0.96);
}

.voice-btn.recording {
  background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.8; }
}

.voice-icon {
  font-size: 32px;
  color: #fff;
}

.voice-text {
  font-size: 14px;
  color: #fff;
  font-weight: 500;
}

.voice-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 8px;
  font-size: 12px;
  color: #86909c;
}

.control-section {
  padding: 0 12px;
  margin-bottom: 12px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 12px;
  padding-left: 4px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.action-item {
  background: #fff;
  border-radius: 12px;
  padding: 14px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.15s;
}

.action-item:active {
  transform: scale(0.95);
}

.action-item.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.action-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-icon {
  font-size: 22px;
}

.action-label {
  font-size: 12px;
  color: #1d2129;
}

.info-section {
  padding: 0 12px;
  margin-bottom: 12px;
}

.info-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-bottom: 1px solid #f0f0f0;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 13px;
  color: #86909c;
}

.info-value {
  font-size: 13px;
  color: #1d2129;
  font-weight: 500;
}

.log-section {
  padding: 0 12px;
}

.log-list {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
}

.log-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-bottom: 1px solid #f0f0f0;
}

.log-item:last-child {
  border-bottom: none;
}

.log-icon {
  font-size: 16px;
}

.log-icon.success { color: #52c41a; }
.log-icon.failed { color: #ff4d4f; }
.log-icon.pending { color: #86909c; }

.log-content {
  flex: 1;
}

.log-cmd {
  font-size: 13px;
  color: #1d2129;
  margin-bottom: 2px;
}

.log-time {
  font-size: 11px;
  color: #86909c;
}

.log-status {
  font-size: 12px;
}

.log-status.success { color: #52c41a; }
.log-status.failed { color: #ff4d4f; }
.log-status.pending { color: #86909c; }

.voice-modal-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
}

.voice-wave {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 4px;
  height: 48px;
}

.wave-bar {
  width: 4px;
  height: 8px;
  background: #c9cdd4;
  border-radius: 2px;
  transition: height 0.2s;
}

.wave-bar:nth-child(1) { animation: none; }
.wave-bar:nth-child(2) { animation: none; }
.wave-bar:nth-child(3) { animation: none; }
.wave-bar:nth-child(4) { animation: none; }
.wave-bar:nth-child(5) { animation: none; }

.voice-wave.active .wave-bar {
  animation: wave 0.8s ease-in-out infinite;
  background: #165dff;
}

.voice-wave.active .wave-bar:nth-child(1) { animation-delay: 0s; }
.voice-wave.active .wave-bar:nth-child(2) { animation-delay: 0.1s; }
.voice-wave.active .wave-bar:nth-child(3) { animation-delay: 0.2s; }
.voice-wave.active .wave-bar:nth-child(4) { animation-delay: 0.3s; }
.voice-wave.active .wave-bar:nth-child(5) { animation-delay: 0.4s; }

@keyframes wave {
  0%, 100% { height: 8px; }
  50% { height: 40px; }
}

.voice-hint-text {
  font-size: 14px;
  color: #86909c;
}

.voice-actions {
  display: flex;
  gap: 12px;
}

.recognized-text {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: #f3ffed;
  border-radius: 8px;
  font-size: 13px;
  color: #1d2129;
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

.mini-toast-modal :deep(.arco-modal) {
  width: auto !important;
  background: transparent;
  box-shadow: none;
}
</style>
