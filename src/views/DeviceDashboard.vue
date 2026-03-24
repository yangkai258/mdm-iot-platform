<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备管理</a-breadcrumb-item>
      <a-breadcrumb-item>设备大盘</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="设备ID">
          <a-input v-model="searchForm.device_id" placeholder="请输入设备ID" allow-clear />
        </a-form-item>
        <a-form-item label="硬件型号">
          <a-select v-model="searchForm.hardware_model" placeholder="选择型号" allow-clear style="width: 160px">
            <a-option value="MDM-Pro-200">MDM-Pro-200</a-option>
            <a-option value="MDM-Mini-100">MDM-Mini-100</a-option>
            <a-option value="MDM-Lite-50">MDM-Lite-50</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option :value="1">待激活</a-option>
            <a-option :value="2">服役中</a-option>
            <a-option :value="3">维修中</a-option>
            <a-option :value="4">已挂失</a-option>
            <a-option :value="5">已报废</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 统计卡片 -->
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

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-button type="primary" @click="loadDevices">刷新</a-button>
    </div>

    <!-- 表格 -->
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
        <a-space>
          <a-button type="text" size="small" @click="viewDevice(record)">详情</a-button>
        </a-space>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const router = useRouter()

const devices = ref([])
const loading = ref(false)

const searchForm = reactive({
  device_id: '',
  hardware_model: '',
  status: ''
})

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
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (searchForm.device_id) params.device_id = searchForm.device_id
    if (searchForm.hardware_model) params.hardware_model = searchForm.hardware_model
    if (searchForm.status !== '') params.status = searchForm.status

    const res = await axios.get(`${API_BASE}/devices`, { params })
    if (res.data.code === 0) {
      devices.value = res.data.data.list
      pagination.total = res.data.data.pagination.total
      
      // 更新统计数据
      stats.total = res.data.data.pagination.total
      stats.online = devices.value.filter(d => d.is_online).length
      stats.offline = stats.total - stats.online
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
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadDevices()
}

const handleReset = () => {
  searchForm.device_id = ''
  searchForm.hardware_model = ''
  searchForm.status = ''
  pagination.current = 1
  loadDevices()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadDevices()
}

const viewDevice = (record) => {
  router.push(`/device/${record.device_id}`)
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
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.stats-row {
  margin-bottom: 16px;
}

.toolbar {
  margin-bottom: 16px;
}
</style>
