<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="device-dashboard-container">

    <a-card class="general-card">
      <template #title>
        <span class="card-title">设备查询</span>
      </template>
      <!-- 搜索区 -->
      <a-row :gutter="16">
        <a-col :flex="1">
          <a-form :model="searchForm" layout="vertical" size="small">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="设备ID / MAC">
                  <a-input v-model="searchForm.keyword" placeholder="搜索设备ID或MAC地址" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="硬件型号">
                  <a-select v-model="searchForm.hardware_model" placeholder="全部" allow-clear>
                    <a-option value="MDM-Pro-200">MDM-Pro-200</a-option>
                    <a-option value="MDM-Mini-100">MDM-Mini-100</a-option>
                    <a-option value="MDM-Lite-50">MDM-Lite-50</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="在线状态">
                  <a-select v-model="searchForm.online" placeholder="全部" allow-clear>
                    <a-option :value="true">在线</a-option>
                    <a-option :value="false">离线</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'86px'" style="text-align: right">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title>
        <span class="card-title">设备列表</span>
      </template>
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showAddModal">
            <template #icon><icon-plus /></template>
            新增设备
          </a-button>
          <a-button @click="handleExport">
            <template #icon><icon-download /></template>
            导出
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data="filteredDevices"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
        row-key="device_id"
      >
        <template #isOnline="{ record }">
          <a-badge :status="record.is_online ? 'success' : 'default'" :text="record.is_online ? '在线' : '离线'" />
        </template>
        <template #batteryLevel="{ record }">
          <a-progress
            v-if="record.battery_level > 0"
            :percent="record.battery_level"
            :stroke-width="6"
            :show-text="true"
            size="small"
          />
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
            <a-button type="text" size="small" @click="editDevice(record)">编辑</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑设备弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑设备' : '新增设备'"
      @ok="handleSubmit"
      :width="560"
    >
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
import { useRouter } from 'vue-router'

const router = useRouter()
const API_BASE = '/api/v1'

const devices = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const currentDeviceId = ref('')

const searchForm = reactive({
  keyword: '',
  hardware_model: '',
  online: undefined
})

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
  { title: '设备ID', dataIndex: 'device_id', ellipsis: true, width: 160 },
  { title: 'MAC地址', dataIndex: 'mac_address', width: 160 },
  { title: '硬件型号', dataIndex: 'hardware_model', width: 130 },
  { title: '固件版本', dataIndex: 'firmware_version', width: 100 },
  { title: '在线状态', slotName: 'isOnline', width: 90 },
  { title: '电量', slotName: 'batteryLevel', width: 140 },
  { title: '状态', slotName: 'lifecycleStatus', width: 90 },
  { title: '操作', slotName: 'actions', width: 130, fixed: 'right' }
]

const filteredDevices = computed(() => {
  let result = devices.value
  if (searchForm.keyword) {
    const kw = searchForm.keyword.toLowerCase()
    result = result.filter(d =>
      d.device_id.toLowerCase().includes(kw) ||
      (d.mac_address && d.mac_address.toLowerCase().includes(kw))
    )
  }
  if (searchForm.hardware_model) {
    result = result.filter(d => d.hardware_model === searchForm.hardware_model)
  }
  if (searchForm.online !== undefined && searchForm.online !== null) {
    result = result.filter(d => d.is_online === searchForm.online)
  }
  return result
})

const loadDevices = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(
      `${API_BASE}/devices?page=${pagination.current}&page_size=${pagination.pageSize}`,
      { headers: { 'Authorization': `Bearer ${token}` } }
    )
    const data = await res.json()
    if (data.code === 0) {
      devices.value = data.data.list || []
      pagination.total = data.data.pagination?.total || 0
    }
  } catch (e) {
    // use mock data
    devices.value = getMockDevices()
    pagination.total = devices.value.length
  } finally {
    loading.value = false
  }
}

const getMockDevices = () => [
  { device_id: 'DEV001', mac_address: '00:11:22:33:44:55', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 85, lifecycle_status: 2 },
  { device_id: 'DEV002', mac_address: '00:11:22:33:44:56', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: true, battery_level: 72, lifecycle_status: 2 },
  { device_id: 'DEV003', mac_address: '00:11:22:33:44:57', hardware_model: 'MDM-Lite-50', firmware_version: 'v1.0.0', is_online: false, battery_level: 0, lifecycle_status: 1 },
  { device_id: 'DEV004', mac_address: '00:11:22:33:44:58', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 95, lifecycle_status: 2 },
  { device_id: 'DEV005', mac_address: '00:11:22:33:44:59', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: false, battery_level: 0, lifecycle_status: 3 }
]

const handleSearch = () => {
  pagination.current = 1
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.hardware_model = ''
  searchForm.online = undefined
  pagination.current = 1
}

const onPageChange = (page) => {
  pagination.current = page
  loadDevices()
}

const onPageSizeChange = (size) => {
  pagination.pageSize = size
  pagination.current = 1
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
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/devices/${currentDeviceId.value}` : `${API_BASE}/devices`
    const res = await fetch(url, {
      method,
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success(isEdit.value ? '更新成功' : '创建成功')
      modalVisible.value = false
      loadDevices()
      return
    }
  } catch (e) {}
  // mock success
  if (isEdit.value) {
    const idx = devices.value.findIndex(d => d.device_id === currentDeviceId.value)
    if (idx !== -1) devices.value[idx] = { ...devices.value[idx], ...form }
  } else {
    devices.value.unshift({ ...form, is_online: false, battery_level: 0, lifecycle_status: 1 })
  }
  Message.success(isEdit.value ? '设备更新成功（模拟）' : '设备创建成功（模拟）')
  modalVisible.value = false
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
.device-dashboard-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.general-card {
  border-radius: 8px;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}
</style>
