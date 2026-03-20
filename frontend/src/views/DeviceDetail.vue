<template>
  <a-layout class="device-detail">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="status">
          <span>设备状态</span>
        </a-menu-item>
        <a-menu-item key="pet">
          <span>宠物配置</span>
        </a-menu-item>
        <a-menu-item key="ota">
          <span>OTA 固件</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="goBack">
            <span>← 返回</span>
          </a-button>
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <span>设备详情 - {{ deviceId }}</span>
        </div>
        <div class="header-right">
        </div>
      </a-layout-header>

      <a-layout-content class="content">
        <a-spin :spinning="loading">
          <!-- 基本信息 -->
          <a-card class="detail-card">
            <template #title>
              <span>基本信息</span>
            </template>
            <a-descriptions :column="3" bordered>
              <a-descriptions-item label="设备ID">{{ deviceInfo.device_id }}</a-descriptions-item>
              <a-descriptions-item label="MAC地址">{{ deviceInfo.mac_address }}</a-descriptions-item>
              <a-descriptions-item label="硬件型号">{{ deviceInfo.hardware_model }}</a-descriptions-item>
              <a-descriptions-item label="固件版本">{{ deviceInfo.firmware_version }}</a-descriptions-item>
              <a-descriptions-item label="硬件版本">{{ deviceInfo.hardware_version }}</a-descriptions-item>
              <a-descriptions-item label="产品名称">{{ deviceInfo.product_name }}</a-descriptions-item>
            </a-descriptions>
          </a-card>

          <!-- 状态信息 -->
          <a-card class="detail-card">
            <template #title>
              <span>状态信息</span>
            </template>
            <a-row :gutter="16">
              <a-col :span="8">
                <a-statistic title="在线状态">
                  <template #prefix>
                    <a-badge :status="deviceInfo.is_online ? 'success' : 'default'" />
                  </template>
                  <span :style="{ color: deviceInfo.is_online ? '#52c41a' : '#999' }">
                    {{ deviceInfo.is_online ? '在线' : '离线' }}
                  </span>
                </a-statistic>
              </a-col>
              <a-col :span="8">
                <a-statistic title="电量">
                  <template #prefix>
                    <span>🔋</span>
                  </template>
                  <a-progress :percent="deviceInfo.battery_level" :stroke-width="8" :show-text="true" v-if="deviceInfo.battery_level > 0" />
                  <span v-else>-</span>
                </a-statistic>
              </a-col>
              <a-col :span="8">
                <a-statistic title="生命周期状态">
                  <template #prefix>
                    <a-tag :color="getStatusColor(deviceInfo.lifecycle_status)">
                      {{ getStatusText(deviceInfo.lifecycle_status) }}
                    </a-tag>
                  </template>
                </a-statistic>
              </a-col>
            </a-row>
          </a-card>

          <!-- 状态管理 -->
          <a-card class="detail-card">
            <template #title>
              <span>状态管理</span>
            </template>
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
                  <a-button 
                    type="primary" 
                    :disabled="deviceInfo.lifecycle_status === 3"
                    @click="handleChangeStatus(3)"
                    :loading="statusLoading === 3"
                  >
                    设为维修中
                  </a-button>
                  <a-button 
                    :disabled="deviceInfo.lifecycle_status === 4"
                    @click="handleChangeStatus(4)"
                    :loading="statusLoading === 4"
                  >
                    设为挂失
                  </a-button>
                  <a-button 
                    danger 
                    :disabled="deviceInfo.lifecycle_status === 5"
                    @click="handleChangeStatus(5)"
                    :loading="statusLoading === 5"
                  >
                    设为报废
                  </a-button>
                </a-space>
              </a-col>
            </a-row>
          </a-card>

          <!-- 指令下发 -->
          <a-card class="detail-card">
            <template #title>
              <span>指令下发</span>
            </template>
            <a-space wrap>
              <a-button type="primary" @click="handleSendCommand('reboot')" :loading="commandLoading === 'reboot'">
                🔄 重启设备
              </a-button>
              <a-button @click="handleSendCommand('factory_reset')" :loading="commandLoading === 'factory_reset'">
                恢复出厂设置
              </a-button>
              <a-button @click="handleSendCommand('sync_time')" :loading="commandLoading === 'sync_time'">
                同步时间
              </a-button>
              <a-button @click="handleSendCommand('update_config')" :loading="commandLoading === 'update_config'">
                更新配置
              </a-button>
              <a-button @click="handleSendCommand('self_test')" :loading="commandLoading === 'self_test'">
                设备自检
              </a-button>
            </a-space>
          </a-card>

          <!-- 设备操作 -->
          <a-card class="detail-card">
            <template #title>
              <span>设备操作</span>
            </template>
            <a-space>
              <a-button type="primary" @click="handleRefresh">刷新数据</a-button>
              <a-button @click="handleReboot">重启设备</a-button>
              <a-button @click="handleReset">恢复出厂</a-button>
              <a-button type="outline" @click="handleUpgrade">检查更新</a-button>
            </a-space>
          </a-card>

          <!-- 运行日志 -->
          <a-card class="detail-card">
            <template #title>
              <span>运行日志</span>
            </template>
            <a-timeline>
              <a-timeline-item v-for="log in logs" :key="log.id" :color="log.color">
                <p>{{ log.time }} - {{ log.message }}</p>
              </a-timeline-item>
            </a-timeline>
          </a-card>
        </a-spin>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const route = useRoute()
const router = useRouter()

const deviceId = ref(route.params.id || '')
const collapsed = ref(false)
const selectedKeys = ref(['dashboard'])
const loading = ref(false)

const deviceInfo = reactive({
  device_id: '',
  mac_address: '',
  hardware_model: '',
  firmware_version: '',
  hardware_version: '',
  product_name: '',
  is_online: false,
  battery_level: 0,
  lifecycle_status: 1
})

const statusLoading = ref(null)
const commandLoading = ref(null)

const logs = ref([
  { id: 1, time: '2026-03-19 10:00:00', message: '设备已连接', color: 'green' },
  { id: 2, time: '2026-03-19 09:30:00', message: '固件版本更新至 v1.2.0', color: 'blue' },
  { id: 3, time: '2026-03-18 14:20:00', message: '设备重启完成', color: 'gray' }
])

const API_BASE = '/api/v1'

const loadDeviceDetail = async () => {
  if (!deviceId.value) return
  
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/devices/${deviceId.value}`)
    if (res.data.code === 0) {
      Object.assign(deviceInfo, res.data.data)
    }
  } catch (err) {
    Message.error('加载设备详情失败: ' + err.message)
    // 使用模拟数据
    deviceInfo.device_id = deviceId.value
    deviceInfo.mac_address = '00:11:22:33:44:55'
    deviceInfo.hardware_model = 'MDM-Pro-200'
    deviceInfo.firmware_version = 'v1.2.0'
    deviceInfo.hardware_version = 'v2.1'
    deviceInfo.product_name = 'MDM 智能设备'
    deviceInfo.is_online = true
    deviceInfo.battery_level = 85
    deviceInfo.lifecycle_status = 2
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/dashboard')
}

const handleMenuClick = ({ key }) => {
  if (key === 'dashboard') {
    router.push('/dashboard')
  } else if (key === 'ota') {
    router.push('/ota')
  } else if (key === 'pet') {
    router.push('/pet')
  } else if (key === 'status') {
    router.push('/status')
  }
}

const handleRefresh = () => {
  loadDeviceDetail()
  Message.success('数据已刷新')
}

const handleChangeStatus = async (status) => {
  const statusText = { 3: '维修中', 4: '挂失', 5: '报废' }
  
  statusLoading.value = status
  try {
    const res = await axios.put(`${API_BASE}/devices/${deviceId.value}/status`, {
      status: status
    })
    if (res.data.code === 0) {
      deviceInfo.lifecycle_status = status
      Message.success(`设备已设为${statusText[status]}`)
    } else {
      Message.error(res.data.message || '状态更新失败')
    }
  } catch (err) {
    // 模拟成功
    deviceInfo.lifecycle_status = status
    Message.success(`设备已设为${statusText[status]}（模拟）`)
  } finally {
    statusLoading.value = null
  }
}

const handleSendCommand = async (command) => {
  const commandTexts = {
    reboot: '重启',
    factory_reset: '恢复出厂设置',
    sync_time: '同步时间',
    update_config: '更新配置',
    self_test: '设备自检'
  }

  // 确认危险操作
  if (command === 'factory_reset') {
    // 使用 Modal.confirm 来确认
    import('@arco-design/web-vue').then(({ Modal }) => {
      Modal.confirm({
        title: '确认恢复出厂设置？',
        content: '此操作将清除设备所有数据并恢复出厂设置，无法撤销，请谨慎操作。',
        okText: '确认',
        cancelText: '取消',
        onOk: async () => {
          await executeCommand(command, commandTexts[command])
        }
      })
    })
    return
  }

  await executeCommand(command, commandTexts[command])
}

const executeCommand = async (command, commandText) => {
  commandLoading.value = command
  try {
    const res = await axios.post(`${API_BASE}/devices/${deviceId.value}/commands`, {
      command: command,
      params: {}
    })
    if (res.data.code === 0) {
      Message.success(`已下发${commandText}指令`)
    } else {
      Message.error(res.data.message || `指令下发失败`)
    }
  } catch (err) {
    // 模拟成功
    Message.success(`已下发${commandText}指令（模拟）`)
  } finally {
    commandLoading.value = null
  }
}

const handleReboot = () => {
  Message.info('正在重启设备...')
}

const handleReset = () => {
  Message.warning('恢复出厂设置功能开发中')
}

const handleUpgrade = () => {
  router.push('/ota')
}

const getStatusColor = (status) => {
  const colors = { 1: 'blue', 2: 'green', 3: 'orange', 4: 'red', 5: 'gray' }
  return colors[status] || 'default'
}

const getStatusText = (status) => {
  const texts = { 1: '待激活', 2: '服役中', 3: '维修中', 4: '已挂失', 5: '已报废' }
  return texts[status] || '未知'
}

onMounted(() => {
  loadDeviceDetail()
})
</script>

<style scoped>
.device-detail {
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
  padding: 0;
}

.detail-card {
  margin-bottom: 16px;
}
</style>
