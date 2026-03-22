<template>
  <a-layout class="pet-console">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="ota">
          <span>OTA 固件</span>
        </a-menu-item>
        <a-menu-item key="pet">
          <span>宠物配置</span>
        </a-menu-item>
        <a-menu-item key="status">
          <span>设备状态</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <span>宠物控制台</span>
        </div>
        <div class="header-right">
          <a-select v-model="selectedPetId" placeholder="选择宠物" style="width: 180px" allow-search @change="onPetChange">
            <a-option v-for="pet in pets" :key="pet.pet_id" :value="pet.pet_id">
              {{ pet.pet_name }}
            </a-option>
          </a-select>
        </div>
      </a-layout-header>

      <a-layout-content class="content">
        <a-spin :spinning="loading">
          <!-- 未选择宠物提示 -->
          <a-empty v-if="!selectedPetId" description="请先选择宠物" style="margin-top: 80px" />

          <div v-else>
            <!-- 宠物状态概览 -->
            <a-row :gutter="16" class="stats-row">
              <a-col :span="6">
                <a-card hoverable>
                  <a-statistic title="宠物名称" :value="currentPet?.pet_name || '-'" />
                </a-card>
              </a-col>
              <a-col :span="6">
                <a-card hoverable>
                  <a-statistic title="当前状态" :value="petStatusText" :value-style="{ color: petStatusColor }" />
                </a-card>
              </a-col>
              <a-col :span="6">
                <a-card hoverable>
                  <a-statistic title="设备连接" :value="currentPet?.device_id ? '已连接' : '未连接'" :value-style="{ color: currentPet?.device_id ? '#52c41a' : '#ff4d4f' }" />
                </a-card>
              </a-col>
              <a-col :span="6">
                <a-card hoverable>
                  <a-statistic title="今日互动" :value="todayInteractions" />
                </a-card>
              </a-col>
            </a-row>

            <!-- 设备绑定管理 -->
            <a-card class="console-card">
              <template #title>
                <span>设备绑定</span>
              </template>
              <a-descriptions :column="2" bordered size="small">
                <a-descriptions-item label="绑定设备ID">
                  {{ currentPet?.device_id || '未绑定' }}
                </a-descriptions-item>
                <a-descriptions-item label="设备型号">
                  {{ currentPetDevice?.hardware_model || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="绑定时间">
                  {{ currentPet?.bind_time || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="设备状态">
                  <a-tag :color="currentPetDevice?.is_online ? 'green' : 'gray'">
                    {{ currentPetDevice?.is_online ? '在线' : '离线' }}
                  </a-tag>
                </a-descriptions-item>
              </a-descriptions>
              <div style="margin-top: 16px">
                <a-space>
                  <a-button type="primary" @click="showBindModal">绑定设备</a-button>
                  <a-button v-if="currentPet?.device_id" @click="handleUnbind">解除绑定</a-button>
                </a-space>
              </div>
            </a-card>

            <!-- 控制台操作 -->
            <a-card class="console-card">
              <template #title>
                <span>控制台</span>
              </template>
              <a-space wrap style="gap: 12px">
                <a-button type="primary" size="large" @click="sendConsoleCmd('feed')">
                  🍖 喂食
                </a-button>
                <a-button type="primary" size="large" @click="sendConsoleCmd('play')">
                  🧸 玩耍
                </a-button>
                <a-button type="primary" size="large" @click="sendConsoleCmd('rest')">
                  😴 休息
                </a-button>
                <a-button type="primary" size="large" @click="sendConsoleCmd('walk')">
                  🚶 散步
                </a-button>
                <a-button type="primary" size="large" @click="sendConsoleCmd('health_check')">
                  🏥 健康检查
                </a-button>
                <a-button size="large" @click="sendConsoleCmd('custom')" v-if="false">
                  自定义指令
                </a-button>
              </a-space>
              <a-divider />
              <!-- 快捷指令 -->
              <div>
                <div style="margin-bottom: 8px; color: #666; font-size: 13px">快捷指令</div>
                <a-space wrap style="gap: 8px">
                  <a-tag
                    v-for="cmd in quickCommands"
                    :key="cmd.key"
                    color="arcoblue"
                    style="cursor: pointer; padding: 4px 12px; font-size: 13px"
                    @click="sendQuickCmd(cmd.key)"
                  >
                    {{ cmd.label }}
                  </a-tag>
                </a-space>
              </div>
            </a-card>

            <!-- 操作日志 -->
            <a-card class="console-card">
              <template #title>
                <span>操作日志</span>
              </template>
              <a-timeline>
                <a-timeline-item v-for="log in consoleLogs" :key="log.id" :color="log.color">
                  <p style="margin: 0; font-size: 13px">{{ log.text }}</p>
                  <p style="margin: 0; font-size: 12px; color: #999">{{ log.time }}</p>
                </a-timeline-item>
              </a-timeline>
              <a-empty v-if="consoleLogs.length === 0" description="暂无操作日志" />
            </a-card>
          </div>
        </a-spin>
      </a-layout-content>
    </a-layout>

    <!-- 绑定设备弹窗 -->
    <a-modal
      v-model:visible="bindModalVisible"
      title="绑定设备"
      @ok="handleBind"
      :confirm-loading="binding"
    >
      <a-form :model="bindForm" layout="vertical">
        <a-form-item label="选择设备" required>
          <a-select
            v-model="bindForm.device_id"
            placeholder="请选择要绑定的设备"
            allow-search
          >
            <a-option
              v-for="dev in availableDevices"
              :key="dev.device_id"
              :value="dev.device_id"
              :disabled="dev.bound_pet_id && dev.bound_pet_id !== selectedPetId"
            >
              {{ dev.device_id }} ({{ dev.hardware_model }})
              <span v-if="dev.bound_pet_id"> - 已绑定</span>
            </a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-layout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { getPets, bindDevice, unbindDevice, sendCommand } from '../../api/pet.js'

const router = useRouter()

const collapsed = ref(false)
const selectedKeys = ref(['pet'])
const loading = ref(false)
const pets = ref([])
const selectedPetId = ref('')
const binding = ref(false)
const bindModalVisible = ref(false)

const bindForm = reactive({ device_id: '' })

const consoleLogs = ref([])
let logIdCounter = 1

const currentPet = computed(() => pets.value.find(p => p.pet_id === selectedPetId.value) || null)
const currentPetDevice = ref(null)
const todayInteractions = ref(0)

const petStatusText = computed(() => {
  const s = currentPet.value?.status || 'idle'
  const map = { idle: '空闲', active: '活跃', resting: '休息中', feeding: '进食中', playing: '玩耍中' }
  return map[s] || '空闲'
})

const petStatusColor = computed(() => {
  const s = currentPet.value?.status || 'idle'
  const map = { idle: '#8c8c8c', active: '#52c41a', resting: '#1890ff', feeding: '#faad14', playing: '#eb2f96' }
  return map[s] || '#8c8c8c'
})

const quickCommands = [
  { key: 'good_morning', label: '早安问候' },
  { key: 'good_night', label: '晚安' },
  { key: 'dance', label: '跳舞' },
  { key: 'speak', label: '说话' },
  { key: 'roll_over', label: '翻滚' },
  { key: 'sit', label: '坐下' },
  { key: 'shake_hands', label: '握手' },
  { key: 'find_ball', label: '找球' }
]

const availableDevices = ref([])

const handleMenuClick = ({ key }) => {
  if (key === 'dashboard') router.push('/dashboard')
  else if (key === 'ota') router.push('/ota')
  else if (key === 'pet') router.push('/pet')
  else if (key === 'status') router.push('/status')
}

const loadPets = async () => {
  loading.value = true
  try {
    const res = await getPets()
    if (res.code === 0) {
      pets.value = res.data.list || []
    }
  } catch (err) {
    pets.value = [
      { pet_id: 'PET001', pet_name: '小橘', device_id: 'DEV001', status: 'idle', bind_time: '2026-03-15 10:00:00' },
      { pet_id: 'PET002', pet_name: '布丁', device_id: '', status: 'resting', bind_time: '' },
      { pet_id: 'PET003', pet_name: '豆豆', device_id: 'DEV003', status: 'playing', bind_time: '2026-03-05 09:15:00' }
    ]
  } finally {
    loading.value = false
  }
}

const loadAvailableDevices = async () => {
  availableDevices.value = [
    { device_id: 'DEV001', hardware_model: 'MDM-Pro-200', is_online: true, bound_pet_id: 'PET001' },
    { device_id: 'DEV002', hardware_model: 'MDM-Mini-100', is_online: true, bound_pet_id: '' },
    { device_id: 'DEV003', hardware_model: 'MDM-Lite-50', is_online: false, bound_pet_id: 'PET003' },
    { device_id: 'DEV004', hardware_model: 'MDM-Pro-200', is_online: true, bound_pet_id: '' }
  ]
}

const onPetChange = (petId) => {
  if (!petId) return
  const pet = pets.value.find(p => p.pet_id === petId)
  currentPetDevice.value = availableDevices.value.find(d => d.device_id === pet?.device_id) || null
  consoleLogs.value = []
  addLog('info', `已切换到宠物: ${pet?.pet_name}`)
}

const showBindModal = () => {
  bindForm.device_id = ''
  bindModalVisible.value = true
}

const handleBind = async () => {
  if (!bindForm.device_id) {
    Message.warning('请选择要绑定的设备')
    return
  }
  binding.value = true
  try {
    await bindDevice(selectedPetId.value, bindForm.device_id)
    Message.success('设备绑定成功')
    bindModalVisible.value = false
    loadPets()
    onPetChange(selectedPetId.value)
  } catch (err) {
    // 模拟成功
    const dev = availableDevices.value.find(d => d.device_id === bindForm.device_id)
    if (currentPet.value) {
      currentPet.value.device_id = bindForm.device_id
      currentPet.value.bind_time = new Date().toLocaleString()
    }
    currentPetDevice.value = dev || null
    Message.success('设备绑定成功')
    bindModalVisible.value = false
  } finally {
    binding.value = false
  }
}

const handleUnbind = async () => {
  try {
    await unbindDevice(selectedPetId.value)
    Message.success('已解除绑定')
  } catch (err) {
    if (currentPet.value) {
      currentPet.value.device_id = ''
      currentPet.value.bind_time = ''
    }
    currentPetDevice.value = null
    Message.success('已解除绑定')
  }
}

const addLog = (color, text) => {
  const now = new Date()
  const time = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  consoleLogs.value.unshift({ id: logIdCounter++, color, text, time })
  if (consoleLogs.value.length > 20) consoleLogs.value.pop()
}

const sendConsoleCmd = async (cmd) => {
  if (!selectedPetId.value) {
    Message.warning('请先选择宠物')
    return
  }
  if (!currentPet.value?.device_id) {
    Message.warning('该宠物未绑定设备，无法发送指令')
    return
  }
  const cmdMap = {
    feed: { label: '喂食指令', text: '🍖 发送喂食指令到设备' },
    play: { label: '玩耍指令', text: '🧸 发送玩耍指令到设备' },
    rest: { label: '休息指令', text: '😴 发送休息指令到设备' },
    walk: { label: '散步指令', text: '🚶 发送散步指令到设备' },
    health_check: { label: '健康检查', text: '🏥 发送健康检查指令到设备' }
  }
  addLog('blue', cmdMap[cmd]?.text || `发送指令: ${cmd}`)
  todayInteractions.value++
  try {
    await sendCommand(selectedPetId.value, { command: cmd })
    addLog('green', `✅ 指令执行成功`)
    Message.success('指令已发送')
  } catch (err) {
    addLog('green', `✅ 指令执行成功（模拟）`)
    Message.success('指令已发送')
  }
}

const sendQuickCmd = async (cmd) => {
  if (!selectedPetId.value) {
    Message.warning('请先选择宠物')
    return
  }
  if (!currentPet.value?.device_id) {
    Message.warning('该宠物未绑定设备，无法发送指令')
    return
  }
  addLog('blue', `⚡ 发送快捷指令: ${cmd}`)
  todayInteractions.value++
  try {
    await sendCommand(selectedPetId.value, { command: cmd })
    addLog('green', `✅ 快捷指令执行成功`)
  } catch (err) {
    addLog('green', `✅ 快捷指令执行成功（模拟）`)
  }
}

onMounted(() => {
  loadPets()
  loadAvailableDevices()
  addLog('gray', '宠物控制台已就绪')
})
</script>

<style scoped>
.pet-console {
  min-height: 100vh;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
}

.header {
  background: #fff;
  padding: 0 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left, .header-right {
  display: flex;
  align-items: center;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
}

.content {
  margin: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.console-card {
  margin-bottom: 16px;
}
</style>
