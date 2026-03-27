<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备管理</a-breadcrumb-item>
      <a-breadcrumb-item>设备列表</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索框 -->
    <div class="search-bar">
      <a-input-search
        v-model="searchKeyword"
        placeholder="搜索设备ID或MAC地址"
        style="width: 280px"
        @search="handleSearch"
        search-button
      />
    </div>

    <!-- 操作按钮组 -->
    <div class="action-bar">
      <a-space>
        <a-button type="primary" @click="showAddModal">新增设备</a-button>
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="loadDevices">刷新</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="content-area">
      <a-table
        :columns="columns"
        :data="filteredDevices"
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
          <a-button type="text" size="small" @click="viewDevice(record)">详情</a-button>
          <a-button type="text" size="small" @click="editDevice(record)">编辑</a-button>
        </template>
      </a-table>
    </div>

    <!-- 新增/编辑设备弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑设备' : '新增设备'" @ok="handleSubmit" :width="600">
      <a-form :model="form" layout="vertical">
        <a-form-item label="设备ID" required>
          <a-input v-model="form.device_id" placeholder="请输入设备ID" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="MAC地址" required>
          <a-input v-model="form.mac_address" placeholder="请输入MAC地址" />
        </a-form-item>
        <a-form-item label="硬件型号" required>
          <a-select v-model="form.hardware_model" placeholder="选择硬件型号">
            <a-option value="MDM-Pro-200">MDM-Pro-200</a-option>
            <a-option value="MDM-Mini-100">MDM-Mini-100</a-option>
            <a-option value="MDM-Lite-50">MDM-Lite-50</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="固件版本">
          <a-input v-model="form.firmware_version" placeholder="例如: v1.2.0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const devices = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const modalVisible = ref(false)
const isEdit = ref(false)
const currentDeviceId = ref('')

const form = reactive({
  device_id: '',
  mac_address: '',
  hardware_model: '',
  firmware_version: ''
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
  { title: '操作', slotName: 'actions', width: 150 }
]

const API_BASE = '/api/v1'

const filteredDevices = computed(() => {
  if (!searchKeyword.value) return devices.value
  const kw = searchKeyword.value.toLowerCase()
  return devices.value.filter(d =>
    d.device_id.toLowerCase().includes(kw) ||
    (d.mac_address && d.mac_address.toLowerCase().includes(kw))
  )
})

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
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadDevices()
}

const showAddModal = () => {
  isEdit.value = false
  Object.assign(form, { device_id: '', mac_address: '', hardware_model: '', firmware_version: '' })
  modalVisible.value = true
}

const editDevice = (record) => {
  isEdit.value = true
  currentDeviceId.value = record.device_id
  Object.assign(form, record)
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!form.device_id || !form.mac_address || !form.hardware_model) {
    Message.warning('请填写必填项')
    return
  }
  try {
    if (isEdit.value) {
      await axios.put(`${API_BASE}/devices/${currentDeviceId.value}`, form)
      Message.success('设备更新成功')
    } else {
      await axios.post(`${API_BASE}/devices`, form)
      Message.success('设备创建成功')
    }
    modalVisible.value = false
    loadDevices()
  } catch (err) {
    setTimeout(() => {
      if (isEdit.value) {
        const idx = devices.value.findIndex(d => d.device_id === currentDeviceId.value)
        if (idx !== -1) devices.value[idx] = { ...devices.value[idx], ...form }
        Message.success('设备更新成功（模拟）')
      } else {
        devices.value.unshift({ ...form, is_online: false, battery_level: 0, lifecycle_status: 1 })
        Message.success('设备创建成功（模拟）')
      }
      modalVisible.value = false
    }, 500)
  }
}

const handleExport = () => {
  Message.info('导出功能开发中')
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
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-bar {
  margin-bottom: 12px;
}

.action-bar {
  margin-bottom: 16px;
}

.content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
