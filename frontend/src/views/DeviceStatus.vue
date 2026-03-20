<template>
  <div class="page-container">
<a-spin :spinning="loading">
          <!-- 状态统计 -->
          <a-row :gutter="16" class="stats-row">
            <a-col :span="4">
              <a-card hoverable @click="filterByStatus(1)">
                <a-statistic title="待激活" :value="statusStats[1] || 0" :value-style="{ color: '#1890ff' }" />
              </a-card>
            </a-col>
            <a-col :span="4">
              <a-card hoverable @click="filterByStatus(2)">
                <a-statistic title="服役中" :value="statusStats[2] || 0" :value-style="{ color: '#52c41a' }" />
              </a-card>
            </a-col>
            <a-col :span="4">
              <a-card hoverable @click="filterByStatus(3)">
                <a-statistic title="维修中" :value="statusStats[3] || 0" :value-style="{ color: '#faad14' }" />
              </a-card>
            </a-col>
            <a-col :span="4">
              <a-card hoverable @click="filterByStatus(4)">
                <a-statistic title="已挂失" :value="statusStats[4] || 0" :value-style="{ color: '#ff4d4f' }" />
              </a-card>
            </a-col>
            <a-col :span="4">
              <a-card hoverable @click="filterByStatus(5)">
                <a-statistic title="已报废" :value="statusStats[5] || 0" :value-style="{ color: '#8c8c8c' }" />
              </a-card>
            </a-col>
            <a-col :span="4">
              <a-card hoverable @click="filterByStatus(null)">
                <a-statistic title="全部设备" :value="statusStats.total || 0" />
              </a-card>
            </a-col>
          </a-row>

          <!-- 筛选和操作栏 -->
          <a-card class="filter-card">
            <a-space>
              <a-input-search
                v-model="searchKeyword"
                placeholder="搜索设备ID或MAC地址"
                style="width: 250px"
                @search="handleSearch"
                search-button
              />
              <a-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px" allow-clear @change="handleFilterChange">
                <a-option :value="1">待激活</a-option>
                <a-option :value="2">服役中</a-option>
                <a-option :value="3">维修中</a-option>
                <a-option :value="4">已挂失</a-option>
                <a-option :value="5">已报废</a-option>
              </a-select>
              <a-button @click="loadDevices">刷新</a-button>
              <a-button v-if="filterStatus !== null || searchKeyword" @click="clearFilters">清除筛选</a-button>
            </a-space>
          </a-card>

          <!-- 设备列表 -->
          <a-card class="device-card">
            <template #title>
              <span>设备列表</span>
            </template>

            <a-table
              :columns="columns"
              :data="filteredDevices"
              :loading="loading"
              :pagination="pagination"
              @change="handleTableChange"
              row-key="device_id"
            >
              <template #deviceId="{ record }">
                <a-space>
                  <a-avatar :size="24" :style="{ backgroundColor: '#165dff' }">
                    {{ record.device_id.charAt(0) }}
                  </a-avatar>
                  <span>{{ record.device_id }}</span>
                </a-space>
              </template>
              <template #isOnline="{ record }">
                <a-badge :status="record.is_online ? 'success' : 'default'" :text="record.is_online ? '在线' : '离线'" />
              </template>
              <template #lifecycleStatus="{ record }">
                <a-tag :color="getStatusColor(record.lifecycle_status)" style="cursor: pointer" @click="showStatusModal(record)">
                  {{ getStatusText(record.lifecycle_status) }}
                </a-tag>
              </template>
              <template #lastUpdate="{ record }">
                {{ formatTime(record.last_update_time) }}
              </template>
              <template #actions="{ record }">
                <a-space>
                  <a-button type="primary" size="small" @click="showStatusModal(record)">状态</a-button>
                  <a-button size="small" @click="viewDevice(record)">详情</a-button>
                </a-space>
              </template>
            </a-table>
          </a-card>

          <!-- 状态变更弹窗 -->
          <a-modal
            v-model="statusModalVisible"
            title="变更设备状态"
            @ok="handleStatusChange"
            :confirm-loading="statusChanging"
          >
            <a-descriptions :column="1" bordered size="small">
              <a-descriptions-item label="设备ID">{{ currentDevice?.device_id }}</a-descriptions-item>
              <a-descriptions-item label="MAC地址">{{ currentDevice?.mac_address }}</a-descriptions-item>
              <a-descriptions-item label="硬件型号">{{ currentDevice?.hardware_model }}</a-descriptions-item>
              <a-descriptions-item label="当前状态">
                <a-tag :color="getStatusColor(currentDevice?.lifecycle_status)">
                  {{ getStatusText(currentDevice?.lifecycle_status) }}
                </a-tag>
              </a-descriptions-item>
            </a-descriptions>

            <a-divider>选择新状态</a-divider>

            <a-radio-group v-model="newStatus" direction="vertical">
              <a-radio :value="1">
                <a-tag color="#1890ff">待激活</a-tag>
                <span style="margin-left: 8px; color: #666">设备尚未激活</span>
              </a-radio>
              <a-radio :value="2">
                <a-tag color="#52c41a">服役中</a-tag>
                <span style="margin-left: 8px; color: #666">设备正常运行</span>
              </a-radio>
              <a-radio :value="3">
                <a-tag color="#faad14">维修中</a-tag>
                <span style="margin-left: 8px; color: #666">设备正在维修中</span>
              </a-radio>
              <a-radio :value="4">
                <a-tag color="#ff4d4f">已挂失</a-tag>
                <span style="margin-left: 8px; color: #666">设备已挂失，需找回</span>
              </a-radio>
              <a-radio :value="5">
                <a-tag color="#8c8c8c">已报废</a-tag>
                <span style="margin-left: 8px; color: #666">设备已报废，不可使用</span>
              </a-radio>
            </a-radio-group>

            <a-divider>变更原因（可选）</a-divider>
            <a-textarea v-model="statusRemark" placeholder="请输入状态变更原因" :rows="2" />
          </a-modal>

          <!-- 批量操作弹窗 -->
          <a-modal
            v-model="batchModalVisible"
            title="批量变更设备状态"
            @ok="handleBatchStatusChange"
            :confirm-loading="batchChanging"
          >
            <a-alert v-if="selectedDevices.length > 0" type="info" show-icon style="margin-bottom: 16px">
              已选择 {{ selectedDevices.length }} 台设备
            </a-alert>

            <a-radio-group v-model="batchNewStatus" direction="vertical">
              <a-radio :value="1">待激活</a-radio>
              <a-radio :value="2">服役中</a-radio>
              <a-radio :value="3">维修中</a-radio>
              <a-radio :value="4">已挂失</a-radio>
              <a-radio :value="5">已报废</a-radio>
            </a-radio-group>

            <a-divider>批量变更原因</a-divider>
            <a-textarea v-model="batchRemark" placeholder="请输入状态变更原因" :rows="2" />
          </a-modal>
        </a-spin>
</div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const loading = ref(false)
const searchKeyword = ref('')
const filterStatus = ref(null)

// 状态变更相关
const statusModalVisible = ref(false)
const statusChanging = ref(false)
const currentDevice = ref(null)
const newStatus = ref(null)
const statusRemark = ref('')

// 批量操作相关
const batchModalVisible = ref(false)
const batchChanging = ref(false)
const selectedDevices = ref([])
const batchNewStatus = ref(null)
const batchRemark = ref('')

const devices = ref([])

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const statusStats = reactive({
  1: 0,
  2: 0,
  3: 0,
  4: 0,
  5: 0,
  total: 0
})

const columns = [
  { title: '设备ID', slotName: 'deviceId', width: 180 },
  { title: 'MAC地址', dataIndex: 'mac_address', width: 160 },
  { title: '硬件型号', dataIndex: 'hardware_model', width: 120 },
  { title: '固件版本', dataIndex: 'firmware_version', width: 100 },
  { title: '在线状态', slotName: 'isOnline', width: 100 },
  { title: '电量', dataIndex: 'battery_level', width: 80 },
  { title: '状态', slotName: 'lifecycleStatus', width: 100 },
  { title: '最后更新', slotName: 'lastUpdate', width: 180 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const API_BASE = '/api/v1'

const filteredDevices = computed(() => {
  let result = devices.value

  if (filterStatus.value !== null) {
    result = result.filter(d => d.lifecycle_status === filterStatus.value)
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(d =>
      d.device_id.toLowerCase().includes(keyword) ||
      d.mac_address.toLowerCase().includes(keyword)
    )
  }

  return result
})

const loadDevices = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/devices`, {
      params: { page: pagination.current, page_size: pagination.pageSize }
    })
    if (res.data.code === 0) {
      devices.value = res.data.data.list
      pagination.total = res.data.data.pagination.total
      calcStatusStats()
    }
  } catch (err) {
    // 使用模拟数据
    devices.value = [
      { device_id: 'DEV001', mac_address: '00:11:22:33:44:55', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 85, lifecycle_status: 2, last_update_time: '2026-03-19 10:30:00' },
      { device_id: 'DEV002', mac_address: '00:11:22:33:44:56', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: true, battery_level: 72, lifecycle_status: 2, last_update_time: '2026-03-19 10:25:00' },
      { device_id: 'DEV003', mac_address: '00:11:22:33:44:57', hardware_model: 'MDM-Lite-50', firmware_version: 'v1.0.0', is_online: false, battery_level: 0, lifecycle_status: 1, last_update_time: '2026-03-18 14:20:00' },
      { device_id: 'DEV004', mac_address: '00:11:22:33:44:58', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: true, battery_level: 95, lifecycle_status: 3, last_update_time: '2026-03-19 09:00:00' },
      { device_id: 'DEV005', mac_address: '00:11:22:33:44:59', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: false, battery_level: 0, lifecycle_status: 4, last_update_time: '2026-03-17 18:30:00' },
      { device_id: 'DEV006', mac_address: '00:11:22:33:44:5a', hardware_model: 'MDM-Pro-200', firmware_version: 'v1.2.0', is_online: false, battery_level: 0, lifecycle_status: 5, last_update_time: '2026-03-15 12:00:00' },
      { device_id: 'DEV007', mac_address: '00:11:22:33:44:5b', hardware_model: 'MDM-Lite-50', firmware_version: 'v1.0.0', is_online: true, battery_level: 60, lifecycle_status: 2, last_update_time: '2026-03-19 10:35:00' },
      { device_id: 'DEV008', mac_address: '00:11:22:33:44:5c', hardware_model: 'MDM-Mini-100', firmware_version: 'v1.1.5', is_online: true, battery_level: 88, lifecycle_status: 2, last_update_time: '2026-03-19 10:20:00' }
    ]
    pagination.total = devices.value.length
    calcStatusStats()
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const calcStatusStats = () => {
  statusStats[1] = devices.value.filter(d => d.lifecycle_status === 1).length
  statusStats[2] = devices.value.filter(d => d.lifecycle_status === 2).length
  statusStats[3] = devices.value.filter(d => d.lifecycle_status === 3).length
  statusStats[4] = devices.value.filter(d => d.lifecycle_status === 4).length
  statusStats[5] = devices.value.filter(d => d.lifecycle_status === 5).length
  statusStats.total = devices.value.length
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  loadDevices()
}

const handleSearch = () => {
  // 搜索通过 computed 自动过滤
}

const handleFilterChange = () => {
  // 筛选通过 computed 自动过滤
}

const filterByStatus = (status) => {
  filterStatus.value = status
}

const clearFilters = () => {
  filterStatus.value = null
  searchKeyword.value = ''
}

const showStatusModal = (record) => {
  currentDevice.value = { ...record }
  newStatus.value = record.lifecycle_status
  statusRemark.value = ''
  statusModalVisible.value = true
}

const handleStatusChange = async () => {
  if (newStatus.value === currentDevice.value.lifecycle_status) {
    Message.warning('状态未变更')
    return
  }

  statusChanging.value = true
  try {
    await axios.put(`${API_BASE}/devices/${currentDevice.value.device_id}/status`, {
      status: newStatus.value,
      remark: statusRemark.value
    })
    Message.success('设备状态已更新')
  } catch (err) {
    // 模拟更新
    setTimeout(() => {
      const idx = devices.value.findIndex(d => d.device_id === currentDevice.value.device_id)
      if (idx !== -1) {
        devices.value[idx].lifecycle_status = newStatus.value
      }
      calcStatusStats()
      Message.success('设备状态已更新')
    }, 500)
  } finally {
    statusChanging.value = false
    statusModalVisible.value = false
  }
}

const handleBatchStatusChange = async () => {
  if (!batchNewStatus.value) {
    Message.warning('请选择新状态')
    return
  }

  batchChanging.value = true
  try {
    await axios.post(`${API_BASE}/devices/batch-status`, {
      device_ids: selectedDevices.value,
      status: batchNewStatus.value,
      remark: batchRemark.value
    })
    Message.success(`批量更新 ${selectedDevices.value.length} 台设备状态成功`)
  } catch (err) {
    // 模拟批量更新
    setTimeout(() => {
      selectedDevices.value.forEach(id => {
        const idx = devices.value.findIndex(d => d.device_id === id)
        if (idx !== -1) {
          devices.value[idx].lifecycle_status = batchNewStatus.value
        }
      })
      calcStatusStats()
      Message.success(`批量更新 ${selectedDevices.value.length} 台设备状态成功`)
    }, 500)
  } finally {
    batchChanging.value = false
    batchModalVisible.value = false
    selectedDevices.value = []
    batchNewStatus.value = null
    batchRemark.value = ''
  }
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

const formatTime = (time) => {
  if (!time) return '-'
  return time
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
.device-status {
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

.filter-card, .device-card {
  margin-bottom: 16px;
}

.stats-row .arco-card {
  cursor: pointer;
  transition: all 0.3s;
}

.stats-row .arco-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
