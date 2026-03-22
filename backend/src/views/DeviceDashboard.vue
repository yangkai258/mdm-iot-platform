<template>
  <a-layout class="device-dashboard">
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
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-right">
          <a-badge :count="warningCount">
            <span>🔔</span>
          </a-badge>
        </div>
      </a-layout-header>

      <a-layout-content class="content">
        <a-row :gutter="16" class="stats-row">
          <a-col :span="8">
            <a-card>
              <a-statistic title="总设备数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card>
              <a-statistic title="当前在线" :value="stats.online" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card>
              <a-statistic title="离线告警" :value="stats.offline" :value-style="{ color: '#ff4d4f' }" />
            </a-card>
          </a-col>
        </a-row>

        <a-card class="device-table-card">
          <template #title>
            <div class="table-title">
              <span>设备列表</span>
              <a-button type="primary" @click="loadDevices">刷新</a-button>
            </div>
          </template>

          <a-table
            :columns="columns"
            :data="devices"
            :loading="loading"
            :pagination="pagination"
            @change="handleTableChange"
            row-key="device_id"
          >
            <template #isOnline="{ record }">
              <a-badge :status="record.is_online ? 'success' : 'default'" :text="record.is_online ? '在线' : '离线'" />
            </template>
            <template #batteryLevel="{ record }">
              <a-progress :percent="record.battery_level" :stroke-width="6" :show-text="true" v-if="record.battery_level > 0" />
              <span v-else>-</span>
            </template>
            <template #lifecycleStatus="{ record }">
              <a-tag :color="getStatusColor(record.lifecycle_status)">
                {{ getStatusText(record.lifecycle_status) }}
              </a-tag>
            </template>
            <template #actions="{ record }">
              <a-button type="primary" size="small" @click="viewDevice(record)">详情</a-button>
            </template>
          </a-table>
        </a-card>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const router = useRouter()

const collapsed = ref(false)
const selectedKeys = ref(['dashboard'])
const devices = ref([])
const loading = ref(false)
const warningCount = ref(0)

const stats = reactive({
  total: 0,
  online: 0,
  offline: 0
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showSizeChanger: true
})

const columns = [
  { title: '设备ID', dataIndex: 'device_id', ellipsis: true },
  { title: 'MAC地址', dataIndex: 'mac_address' },
  { title: '硬件型号', dataIndex: 'hardware_model' },
  { title: '固件版本', dataIndex: 'firmware_version' },
  { title: '在线状态', slotName: 'isOnline' },
  { title: '电量', slotName: 'batteryLevel', width: 150 },
  { title: '状态', slotName: 'lifecycleStatus', width: 100 },
  { title: '操作', slotName: 'actions', width: 100 }
]

const API_BASE = '/api/v1'

const loadDevices = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/devices`, {
      params: {
        page: pagination.current,
        page_size: pagination.pageSize
      }
    })
    if (res.data.code === 0) {
      devices.value = res.data.data.list
      pagination.total = res.data.data.pagination.total
      
      // 更新统计数据
      stats.total = res.data.data.pagination.total
      stats.online = devices.value.filter(d => d.is_online).length
      stats.offline = stats.total - stats.online
      warningCount.value = stats.offline
    }
  } catch (err) {
    // 使用模拟数据
    const mockDevices = [
      { device_id: 'DEV001', mac_address: '00:11:22:33:44:55', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 85, lifecycle_status: 2 },
      { device_id: 'DEV002', mac_address: '00:11:22:33:44:56', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: true, battery_level: 72, lifecycle_status: 2 },
      { device_id: 'DEV003', mac_address: '00:11:22:33:44:57', hardware_model: 'MDM-Lite-50', firmware_version: 'v1.0.0', is_online: false, battery_level: 0, lifecycle_status: 1 },
      { device_id: 'DEV004', mac_address: '00:11:22:33:44:58', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 95, lifecycle_status: 2 },
      { device_id: 'DEV005', mac_address: '00:11:22:33:44:59', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: false, battery_level: 0, lifecycle_status: 3 }
    ]
    devices.value = mockDevices
    pagination.total = mockDevices.length
    stats.total = mockDevices.length
    stats.online = mockDevices.filter(d => d.is_online).length
    stats.offline = stats.total - stats.online
    warningCount.value = stats.offline
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadDevices()
}

const viewDevice = (record) => {
  router.push(`/device/${record.device_id}`)
}

const handleMenuClick = ({ key }) => {
  if (key === 'ota') {
    router.push('/ota')
  } else if (key === 'pet') {
    router.push('/pet')
  } else if (key === 'status') {
    router.push('/status')
  }
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
  loadDevices()
})
</script>

<style scoped>
.device-dashboard {
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

.content {
  margin: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.table-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
