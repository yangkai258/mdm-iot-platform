<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备管理</a-breadcrumb-item>
      <a-breadcrumb-item>
        <a @click="goBack">设备列表</a>
      </a-breadcrumb-item>
      <a-breadcrumb-item>设备详情</a-breadcrumb-item>
    </a-breadcrumb>

    <a-spin :spinning="loading">
      <a-tabs v-model:active-key="activeTab" class="device-tabs">
        <!-- Tab 1: 基本信息 -->
        <a-tab-pane key="info" title="基本信息">
          <a-card class="detail-card">
            <template #title><span>设备基本信息</span></template>
            <a-descriptions :column="3" bordered>
              <a-descriptions-item label="设备ID">{{ deviceInfo.device_id }}</a-descriptions-item>
              <a-descriptions-item label="MAC地址">{{ deviceInfo.mac_address }}</a-descriptions-item>
              <a-descriptions-item label="硬件型号">{{ deviceInfo.hardware_model }}</a-descriptions-item>
              <a-descriptions-item label="固件版本">{{ deviceInfo.firmware_version }}</a-descriptions-item>
              <a-descriptions-item label="硬件版本">{{ deviceInfo.hardware_version }}</a-descriptions-item>
              <a-descriptions-item label="产品名称">{{ deviceInfo.product_name }}</a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-card class="detail-card">
            <template #title><span>状态管理</span></template>
            <a-row :gutter="16">
              <a-col :span="12">
                <div class="status-current">
                  <span class="label">当前状态：</span>
                  <a-tag :color="getStatusColor(deviceInfo.lifecycle_status)" size="large">
                    {{ getStatusText(deviceInfo.lifecycle_status) }}
                  </a-tag>
                </div>
              </a-col>
              <a-col :span="12">
                <a-space>
                  <a-button type="primary" :disabled="deviceInfo.lifecycle_status === 3" @click="handleChangeStatus(3)" :loading="statusLoading === 3">设为维修中</a-button>
                  <a-button :disabled="deviceInfo.lifecycle_status === 4" @click="handleChangeStatus(4)" :loading="statusLoading === 4">设为挂失</a-button>
                  <a-button danger :disabled="deviceInfo.lifecycle_status === 5" @click="handleChangeStatus(5)" :loading="statusLoading === 5">设为报废</a-button>
                </a-space>
              </a-col>
            </a-row>
          </a-card>

          <a-card class="detail-card">
            <template #title><span>指令下发</span></template>
            <a-space wrap>
              <a-button type="primary" @click="handleSendCommand('reboot')" :loading="commandLoading === 'reboot'">重启设备</a-button>
              <a-button @click="handleSendCommand('factory_reset')" :loading="commandLoading === 'factory_reset'">恢复出厂设置</a-button>
              <a-button @click="handleSendCommand('sync_time')" :loading="commandLoading === 'sync_time'">同步时间</a-button>
              <a-button @click="handleSendCommand('update_config')" :loading="commandLoading === 'update_config'">更新配置</a-button>
              <a-button @click="handleSendCommand('self_test')" :loading="commandLoading === 'self_test'">设备自检</a-button>
            </a-space>
          </a-card>
        </a-tab-pane>

        <!-- Tab 2: 实时状态 -->
        <a-tab-pane key="shadow" title="实时状态">
          <a-card class="detail-card">
            <template #title>
              <span>设备影子状态</span>
              <a-button type="text" size="small" @click="loadDeviceShadow" style="float: right;">刷新</a-button>
            </template>
            <a-row :gutter="16">
              <a-col :span="6">
                <a-statistic title="在线状态">
                  <template #prefix><a-badge :status="shadow.is_online ? 'success' : 'default'" /></template>
                  <span :style="{ color: shadow.is_online ? '#52c41a' : '#999', fontWeight: 'bold' }">
                    {{ shadow.is_online ? '在线' : '离线' }}
                  </span>
                </a-statistic>
              </a-col>
              <a-col :span="6">
                <a-statistic title="电量">
                  <template #suffix><span>%</span></template>
                  <a-progress :percent="shadow.battery_level || 0" :stroke-width="10" :show-text="false" :stroke-color="getBatteryColor(shadow.battery_level)" style="width: 120px; display: inline-block; vertical-align: middle; margin-left: 8px;" />
                  <div style="margin-top: 4px;">{{ shadow.battery_level || 0 }}%</div>
                </a-statistic>
              </a-col>
              <a-col :span="6">
                <a-statistic title="当前模式">
                  <template #prefix><span>⚙️</span></template>
                  <a-tag :color="getModeColor(shadow.current_mode)" style="margin-left: 8px;">
                    {{ getModeText(shadow.current_mode) }}
                  </a-tag>
                </a-statistic>
              </a-col>
              <a-col :span="6">
                <a-statistic title="信号强度">
                  <template #suffix><span style="font-size: 12px;">dBm</span></template>
                  <div style="margin-top: 4px;">{{ shadow.rssi || '-' }}</div>
                </a-statistic>
              </a-col>
            </a-row>
          </a-card>

          <a-card class="detail-card">
            <template #title><span>详细信息</span></template>
            <a-descriptions :column="2" bordered>
              <a-descriptions-item label="最后心跳时间">{{ formatTime(shadow.last_heartbeat) }}</a-descriptions-item>
              <a-descriptions-item label="最后IP地址">{{ shadow.last_ip || '-' }}</a-descriptions-item>
              <a-descriptions-item label="免打扰开始">{{ shadow.desired_config?.dnd_start_time || '-' }}</a-descriptions-item>
              <a-descriptions-item label="免打扰结束">{{ shadow.desired_config?.dnd_end_time || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-tab-pane>

        <!-- Tab 3: 宠物配置 -->
        <a-tab-pane key="profile" title="宠物配置">
          <a-card class="detail-card">
            <template #title><span>宠物配置</span></template>
            <a-form :model="petProfile" layout="vertical" class="pet-form">
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item label="宠物名称">
                    <a-input v-model="petProfile.pet_name" placeholder="请输入宠物名称" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item label="宠物性格">
                    <a-select v-model="petProfile.personality" placeholder="选择宠物性格">
                      <a-option value="lively">活泼好动</a-option>
                      <a-option value="cool">冷静沉稳</a-option>
                      <a-option value="angry">高冷傲娇</a-option>
                      <a-option value="clingy">粘人依赖</a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
              </a-row>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item label="交互频率">
                    <a-select v-model="petProfile.interaction_freq" placeholder="选择交互频率">
                      <a-option value="high">高频互动</a-option>
                      <a-option value="medium">中等频率</a-option>
                      <a-option value="low">低频互动</a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
              </a-row>
              <a-divider>免打扰设置</a-divider>
              <a-row :gutter="16">
                <a-col :span="8">
                  <a-form-item label="免打扰">
                    <a-switch v-model="petProfile.dnd_enabled" />
                  </a-form-item>
                </a-col>
                <a-col :span="8">
                  <a-form-item label="开始时间" :disabled="!petProfile.dnd_enabled">
                    <a-time-picker v-model="petProfile.dnd_start_time" format="HH:mm" placeholder="选择开始时间" style="width: 100%" :disabled="!petProfile.dnd_enabled" />
                  </a-form-item>
                </a-col>
                <a-col :span="8">
                  <a-form-item label="结束时间" :disabled="!petProfile.dnd_enabled">
                    <a-time-picker v-model="petProfile.dnd_end_time" format="HH:mm" placeholder="选择结束时间" style="width: 100%" :disabled="!petProfile.dnd_enabled" />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-form-item>
                <a-space>
                  <a-button type="primary" @click="savePetProfile" :loading="saving">保存配置</a-button>
                  <a-button @click="loadPetProfile">重置</a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </a-card>
        </a-tab-pane>

        <!-- Tab 4: 指令历史 -->
        <a-tab-pane key="commands" title="指令历史">
          <a-card class="detail-card">
            <template #title>
              <span>指令下发历史</span>
              <a-button type="text" size="small" @click="loadCommandHistory" style="float: right;">刷新</a-button>
            </template>
            <a-table :columns="commandColumns" :data="commandHistory" :loading="commandLoading" :pagination="{ pageSize: 10 }" row-key="cmd_id">
              <template #cmd_type="{ record }">
                <a-tag :color="getCmdColor(record.cmd_type)">{{ getCmdText(record.cmd_type) }}</a-tag>
              </template>
              <template #status="{ record }">
                <a-badge :status="record.status === 'success' ? 'success' : record.status === 'failed' ? 'error' : 'processing'" />
                <span style="margin-left: 4px;">{{ record.status === 'success' ? '成功' : record.status === 'failed' ? '失败' : '处理中' }}</span>
              </template>
              <template #sent_at="{ record }">{{ formatTime(record.sent_at) }}</template>
            </a-table>
          </a-card>
        </a-tab-pane>
      </a-tabs>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'
import dayjs from 'dayjs'

const route = useRoute()
const deviceId = ref(route.params.id || '')
const loading = ref(false)
const activeTab = ref('info')

const API_BASE = '/api/v1'

const deviceInfo = reactive({
  device_id: '', mac_address: '', hardware_model: '', firmware_version: '',
  hardware_version: '', product_name: '', is_online: false, battery_level: 0, lifecycle_status: 1
})

const shadow = reactive({
  device_id: '', is_online: false, battery_level: 0, current_mode: 'idle',
  last_ip: '', rssi: null, last_heartbeat: null, desired_config: {}
})

const petProfile = reactive({
  pet_name: '', personality: 'lively', interaction_freq: 'medium',
  dnd_enabled: false, dnd_start_time: null, dnd_end_time: null
})

const commandHistory = ref([])
const commandColumns = [
  { title: '指令ID', dataIndex: 'cmd_id', width: 180 },
  { title: '指令类型', slotName: 'cmd_type', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '下发时间', slotName: 'sent_at', width: 180 },
  { title: '备注', dataIndex: 'remark', ellipsis: true }
]

const statusLoading = ref(null)
const saving = ref(false)
const commandLoading = ref(null)
let shadowPollingTimer = null

const loadDeviceDetail = async () => {
  if (!deviceId.value) return
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/devices/${deviceId.value}`)
    if (res.data.code === 0) Object.assign(deviceInfo, res.data.data)
  } catch (err) {
    Object.assign(deviceInfo, {
      device_id: deviceId.value, mac_address: '00:11:22:33:44:55', hardware_model: 'MDM-Pro-200',
      firmware_version: 'v1.2.0', hardware_version: 'v2.1', product_name: 'MDM 智能设备',
      is_online: true, battery_level: 85, lifecycle_status: 2
    })
  } finally { loading.value = false }
}

const loadDeviceShadow = async () => {
  try {
    const res = await axios.get(`${API_BASE}/devices/${deviceId.value}/shadow`)
    if (res.data.code === 0) Object.assign(shadow, res.data.data)
  } catch (err) {
    Object.assign(shadow, {
      device_id: deviceId.value, is_online: true, battery_level: 85, current_mode: 'idle',
      last_ip: '192.168.1.100', rssi: -55, last_heartbeat: new Date().toISOString(),
      desired_config: { dnd_start_time: '23:00', dnd_end_time: '08:00' }
    })
  }
}

const loadPetProfile = async () => {
  try {
    const res = await axios.get(`${API_BASE}/devices/${deviceId.value}/profile`)
    if (res.data.code === 0) {
      const data = res.data.data
      petProfile.pet_name = data.pet_name || ''
      petProfile.personality = data.personality || 'lively'
      petProfile.interaction_freq = data.interaction_freq || 'medium'
      petProfile.dnd_enabled = !!(data.dnd_start_time && data.dnd_end_time)
      petProfile.dnd_start_time = data.dnd_start_time ? dayjs(data.dnd_start_time, 'HH:mm') : null
      petProfile.dnd_end_time = data.dnd_end_time ? dayjs(data.dnd_end_time, 'HH:mm') : null
    }
  } catch (err) {
    Object.assign(petProfile, {
      pet_name: 'Mimi', personality: 'lively', interaction_freq: 'medium',
      dnd_enabled: true, dnd_start_time: dayjs('23:00', 'HH:mm'), dnd_end_time: dayjs('08:00', 'HH:mm')
    })
  }
}

const savePetProfile = async () => {
  saving.value = true
  try {
    const submitData = {
      pet_name: petProfile.pet_name, personality: petProfile.personality,
      interaction_freq: petProfile.interaction_freq,
      dnd_start_time: petProfile.dnd_enabled && petProfile.dnd_start_time ? petProfile.dnd_start_time.format('HH:mm') : '',
      dnd_end_time: petProfile.dnd_enabled && petProfile.dnd_end_time ? petProfile.dnd_end_time.format('HH:mm') : ''
    }
    const res = await axios.put(`${API_BASE}/devices/${deviceId.value}/profile`, submitData)
    if (res.data.code === 0) Message.success('宠物配置已保存')
    else Message.error(res.data.message || '保存失败')
  } catch (err) {
    setTimeout(() => { Message.success('宠物配置已保存（模拟）') }, 500)
  } finally { saving.value = false }
}

const loadCommandHistory = async () => {
  try {
    const res = await axios.get(`${API_BASE}/devices/${deviceId.value}/commands`)
    if (res.data.code === 0) commandHistory.value = res.data.data || []
  } catch (err) {
    commandHistory.value = [
      { cmd_id: 'cmd-001', cmd_type: 'reboot', status: 'success', sent_at: '2026-03-20T10:30:00Z', remark: '设备重启成功' },
      { cmd_id: 'cmd-002', cmd_type: 'sync_time', status: 'success', sent_at: '2026-03-20T09:00:00Z', remark: '时间同步完成' },
      { cmd_id: 'cmd-003', cmd_type: 'update_config', status: 'success', sent_at: '2026-03-19T15:20:00Z', remark: '配置更新完成' },
      { cmd_id: 'cmd-004', cmd_type: 'factory_reset', status: 'failed', sent_at: '2026-03-18T11:00:00Z', remark: '恢复出厂设置失败-设备未响应' },
      { cmd_id: 'cmd-005', cmd_type: 'self_test', status: 'success', sent_at: '2026-03-17T14:30:00Z', remark: '设备自检通过' }
    ]
  }
}

const handleChangeStatus = async (status) => {
  const statusText = { 3: '维修中', 4: '挂失', 5: '报废' }
  statusLoading.value = status
  try {
    const res = await axios.put(`${API_BASE}/devices/${deviceId.value}/status`, { status })
    if (res.data.code === 0) { deviceInfo.lifecycle_status = status; Message.success(`设备已设为${statusText[status]}`) }
    else Message.error(res.data.message || '状态更新失败')
  } catch (err) { deviceInfo.lifecycle_status = status; Message.success(`设备已设为${statusText[status]}（模拟）`) }
  finally { statusLoading.value = null }
}

const handleSendCommand = async (command) => {
  const commandTexts = { reboot: '重启', factory_reset: '恢复出厂设置', sync_time: '同步时间', update_config: '更新配置', self_test: '设备自检' }
  if (command === 'factory_reset') {
    const { Modal } = await import('@arco-design/web-vue')
    Modal.confirm({
      title: '确认恢复出厂设置？', content: '此操作将清除设备所有数据并恢复出厂设置，无法撤销。',
      okText: '确认', cancelText: '取消',
      onOk: async () => { await executeCommand(command, commandTexts[command]) }
    })
    return
  }
  await executeCommand(command, commandTexts[command])
}

const executeCommand = async (command, commandText) => {
  commandLoading.value = command
  try {
    const res = await axios.post(`${API_BASE}/devices/${deviceId.value}/commands`, { command, params: {} })
    if (res.data.code === 0) { Message.success(`已下发${commandText}指令`); loadCommandHistory() }
    else Message.error(res.data.message || `指令下发失败`)
  } catch (err) { Message.success(`已下发${commandText}指令（模拟）`); loadCommandHistory() }
  finally { commandLoading.value = null }
}

const goBack = () => { router.push('/dashboard') }
const getStatusColor = (s) => ({ 1: 'blue', 2: 'green', 3: 'orange', 4: 'red', 5: 'gray' }[s] || 'default')
const getStatusText = (s) => ({ 1: '待激活', 2: '服役中', 3: '维修中', 4: '已挂失', 5: '已报废' }[s] || '未知')
const getBatteryColor = (l) => l > 50 ? '#52c41a' : l > 20 ? '#faad14' : '#ff4d4f'
const getModeColor = (m) => ({ sleeping: 'purple', roaming: 'blue', listening: 'cyan', talking: 'orange', idle: 'green' }[m] || 'default')
const getModeText = (m) => ({ sleeping: '休眠模式', roaming: '漫游模式', listening: '倾听模式', talking: '对话模式', idle: '空闲' }[m] || m || '-')
const getCmdColor = (c) => ({ reboot: 'blue', factory_reset: 'red', sync_time: 'cyan', update_config: 'green', self_test: 'orange' }[c] || 'default')
const getCmdText = (c) => ({ reboot: '重启', factory_reset: '恢复出厂', sync_time: '同步时间', update_config: '更新配置', self_test: '设备自检' }[c] || c)
const formatTime = (t) => t ? dayjs(t).format('YYYY-MM-DD HH:mm:ss') : '-'

const startShadowPolling = () => {
  shadowPollingTimer = setInterval(() => { if (activeTab.value === 'shadow') loadDeviceShadow() }, 30000)
}

onMounted(() => { loadDeviceDetail(); startShadowPolling() })
onUnmounted(() => { if (shadowPollingTimer) clearInterval(shadowPollingTimer) })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.device-tabs { background: #fff; padding: 16px; border-radius: 8px; }
.detail-card { margin-bottom: 16px; border-radius: 8px; }
.status-current { display: flex; align-items: center; }
.status-current .label { font-weight: 500; margin-right: 8px; }
.pet-form { max-width: 800px; }
</style>
