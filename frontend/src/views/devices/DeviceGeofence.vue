<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备管理</a-breadcrumb-item>
      <a-breadcrumb-item>地理围栏</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input-search v-model="searchKeyword" placeholder="搜索围栏名称" style="width: 260px" @search="loadGeofences" search-button />
        <a-select v-model="filterStatus" placeholder="围栏状态" allow-clear style="width: 120px" @change="loadGeofences">
          <a-option value="active">启用</a-option>
          <a-option value="inactive">停用</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">新建围栏</a-button>
        <a-button @click="loadGeofences">刷新</a-button>
      </a-space>
    </div>

    <!-- 围栏列表 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="geofences" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '启用' : '停用' }}</a-tag>
        </template>
      </a-table>
        <template #device_count="{ record }">
          <a-badge :count="record.device_count || 0" :max-count="999" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showBindModal(record)">绑定设备</a-button>
            <a-button type="text" size="small" @click="viewAlerts(record)">告警记录</a-button>
            <a-button type="text" size="small" @click="editGeofence(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteGeofence(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 新建/编辑围栏弹窗 -->
    <a-modal v-model:visible="createModalVisible" :title="isEdit ? '编辑围栏' : '新建围栏'" @ok="handleSubmit" :width="580" :loading="submitting">
      <a-form :model="form" layout="vertical">
        <a-form-item label="围栏名称" required>
          <a-input v-model="form.name" placeholder="请输入围栏名称" />
        </a-form-item>
        <a-form-item label="围栏类型" required>
          <a-select v-model="form.geofence_type" placeholder="选择围栏类型">
            <a-option value="circle">圆形区域</a-option>
            <a-option value="polygon">多边形区域</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="中心坐标">
          <a-space>
            <a-input-number v-model="form.latitude" placeholder="纬度" style="width: 140px" />
            <a-input-number v-model="form.longitude" placeholder="经度" style="width: 140px" />
          </a-space>
        </a-form-item>
        <a-form-item label="半径(米)">
          <a-input-number v-model="form.radius" placeholder="圆形围栏半径" :min="10" :max="10000" style="width: 200px" />
        </a-form-item>
        <a-form-item label="告警触发">
          <a-checkbox-group v-model="form.alarm_types">
            <a-checkbox value="enter">进入告警</a-checkbox>
            <a-checkbox value="exit">离开告警</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model="form.is_active" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 绑定设备弹窗 -->
    <a-modal v-model:visible="bindModalVisible" title="绑定设备到围栏" @ok="handleBind" :width="620" :loading="submitting">
      <a-form layout="vertical">
        <a-form-item label="当前围栏">
          <a-tag color="arcoblue">{{ selectedGeofence?.name }}</a-tag>
        </a-form-item>
        <a-form-item label="选择设备">
          <a-select v-model="selectedDeviceIds" multiple placeholder="选择要绑定的设备" style="width: 100%">
            <a-option v-for="d in allDevices" :key="d.id" :value="d.id">{{ d.name }} ({{ d.id }})</a-option>
          </a-select>
        </a-form-item>
        <div style="color: #8a8a8a; font-size: 12px">已选择 <strong>{{ selectedDeviceIds.length }}</strong> 台设备</div>
      </a-form>
    </a-modal>

    <!-- 围栏告警记录 -->
    <a-modal v-model:visible="alertsModalVisible" title="围栏告警记录" :width="800" :footer="null">
      <a-table :columns="alertColumns" :data="alerts" :loading="alertsLoading" :pagination="alertPagination" row-key="id" @page-change="handleAlertPageChange">
        <template #alert_type="{ record }">
          <a-tag :color="record.alert_type === 'enter' ? 'green' : 'orange'">{{ record.alert_type === 'enter' ? '进入' : '离开' }}</a-tag>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const alertsLoading = ref(false)
const geofences = ref([])
const allDevices = ref([])
const searchKeyword = ref('')
const filterStatus = ref('')
const createModalVisible = ref(false)
const bindModalVisible = ref(false)
const alertsModalVisible = ref(false)
const isEdit = ref(false)
const selectedGeofence = ref(null)
const selectedDeviceIds = ref([])
const alerts = ref([])

const form = reactive({
  name: '', geofence_type: 'circle', latitude: 0, longitude: 0,
  radius: 500, alarm_types: ['enter', 'exit'], is_active: true
})

const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const alertPagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: '围栏ID', dataIndex: 'id', width: 80 },
  { title: '围栏名称', dataIndex: 'name' },
  { title: '围栏类型', dataIndex: 'geofence_type', width: 100 },
  { title: '绑定设备数', slotName: 'device_count', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' }
]

const alertColumns = [
  { title: '告警ID', dataIndex: 'id', width: 80 },
  { title: '设备ID', dataIndex: 'device_id', width: 100 },
  { title: '围栏名称', dataIndex: 'geofence_name' },
  { title: '告警类型', slotName: 'alert_type', width: 90 },
  { title: '触发时间', dataIndex: 'created_at', width: 170 }
]

const loadGeofences = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (searchKeyword.value) params.append('keyword', searchKeyword.value)
    if (filterStatus.value) params.append('status', filterStatus.value)
    const res = await fetch(`/api/v1/device/geofences?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) {
      geofences.value = data.data.list || []
      pagination.total = data.data.total || 0
    }
  } catch (e) { Message.error('加载围栏列表失败') }
  finally { loading.value = false }
}

const loadDevices = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/devices?page_size=200', { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) allDevices.value = data.data.list || []
  } catch (e) { console.error('加载设备失败', e) }
}

const loadAlerts = async (geofenceId) => {
  alertsLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ geofence_id: geofenceId, page: alertPagination.current, page_size: alertPagination.pageSize })
    const res = await fetch(`/api/v1/device/geofences/alerts?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) {
      alerts.value = data.data.list || []
      alertPagination.total = data.data.total || 0
    }
  } catch (e) { Message.error('加载告警记录失败') }
  finally { alertsLoading.value = false }
}

const showCreateModal = () => {
  isEdit.value = false
  Object.assign(form, { name: '', geofence_type: 'circle', latitude: 0, longitude: 0, radius: 500, alarm_types: ['enter', 'exit'], is_active: true })
  createModalVisible.value = true
}

const editGeofence = (record) => {
  isEdit.value = true
  Object.assign(form, record)
  createModalVisible.value = true
}

const handleSubmit = async () => {
  if (!form.name) { Message.warning('请输入围栏名称'); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/device/geofences/${form.id}` : '/api/v1/device/geofences'
    const res = await fetch(url, {
      method: isEdit.value ? 'PUT' : 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success(isEdit.value ? '围栏更新成功' : '围栏创建成功')
      createModalVisible.value = false
      loadGeofences()
    } else { Message.error(data.message || '操作失败') }
  } catch (e) { Message.error('操作失败') }
  finally { submitting.value = false }
}

const showBindModal = async (geofence) => {
  selectedGeofence.value = geofence
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/device/geofences/${geofence.id}/devices`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    selectedDeviceIds.value = data.data?.device_ids || []
  } catch (e) { selectedDeviceIds.value = [] }
  bindModalVisible.value = true
}

const handleBind = async () => {
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/device/geofences/${selectedGeofence.value.id}/devices`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ device_ids: selectedDeviceIds.value })
    })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('设备绑定成功')
      bindModalVisible.value = false
      loadGeofences()
    } else { Message.error(data.message || '绑定失败') }
  } catch (e) { Message.error('绑定失败') }
  finally { submitting.value = false }
}

const viewAlerts = (geofence) => {
  selectedGeofence.value = geofence
  alertPagination.current = 1
  loadAlerts(geofence.id)
  alertsModalVisible.value = true
}

const deleteGeofence = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/device/geofences/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) {
      Message.success('删除成功')
      loadGeofences()
    }
  } catch (e) { Message.error('删除失败') }
}

const handlePageChange = (page) => { pagination.current = page; loadGeofences() }
const handleAlertPageChange = (page) => { alertPagination.current = page; loadAlerts(selectedGeofence.value.id) }

onMounted(() => { loadGeofences(); loadDevices() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
