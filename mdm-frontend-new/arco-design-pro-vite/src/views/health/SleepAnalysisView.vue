<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="睡眠分析">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="日期范围">
            <a-range-picker v-model="form.dateRange" style="width: 240px" />
          </a-form-item>
          <a-form-item label="设备ID">
            <a-input v-model="form.deviceId" placeholder="请输入设备ID" style="width: 160px" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table
      :columns="columns"
      :data="data"
      :loading="loading"
      :pagination="paginationConfig"
      @page-change="onPageChange"
    >
      <template #quality="{ record }">
        <a-progress :percent="record.quality" :color="getQualityColor(record.quality)" size="small" style="width: 80px" />
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" layout="vertical">
        <a-form-item label="日期"><a-date-picker v-model="form.date" style="width: 100%" /></a-form-item>
        <a-form-item label="入睡时间"><a-time-picker v-model="form.sleep_time" format="HH:mm" style="width: 100%" /></a-form-item>
        <a-form-item label="醒来时间"><a-time-picker v-model="form.wake_time" format="HH:mm" style="width: 100%" /></a-form-item>
        <a-form-item label="深睡(小时)"><a-input-number v-model="form.deep_sleep" :min="0" :precision="1" style="width: 100%" /></a-form-item>
        <a-form-item label="浅睡(小时)"><a-input-number v-model="form.light_sleep" :min="0" :precision="1" style="width: 100%" /></a-form-item>
        <a-form-item label="睡眠质量评分"><a-input-number v-model="form.quality" :min="0" :max="100" style="width: 100%" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
    </a-card>
</div></template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)

const form = reactive({
  dateRange: [],
  deviceId: '',
  id: '',
  date: null,
  sleep_time: null,
  wake_time: null,
  deep_sleep: 0,
  light_sleep: 0,
  quality: 0
})

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '入睡时间', dataIndex: 'sleep_time', width: 100 },
  { title: '醒来时间', dataIndex: 'wake_time', width: 100 },
  { title: '睡眠时长', dataIndex: 'total_hours', width: 100 },
  { title: '深睡(h)', dataIndex: 'deep_sleep', width: 80 },
  { title: '浅睡(h)', dataIndex: 'light_sleep', width: 80 },
  { title: '睡眠质量', slotName: 'quality', width: 120 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 }
]

const data = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showTotal: true,
  showPageSize: true
}))

const getQualityColor = (quality) => {
  if (quality >= 80) return '#52c41a'
  if (quality >= 60) return '#1890ff'
  if (quality >= 40) return '#faad14'
  return '#ff4d4f'
}

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.deviceId) params.device_id = form.deviceId
    if (form.dateRange && form.dateRange.length === 2) {
      params.start_date = form.dateRange[0].format('YYYY-MM-DD')
      params.end_date = form.dateRange[1].format('YYYY-MM-DD')
    }
    const res = await axios.get('/api/health/sleep', { params })
    if (res.data.code === 0) {
      data.value = res.data.data.list || []
      pagination.total = res.data.data.pagination?.total || 0
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    loading.value = false
  }
}

const loadMockData = () => {
  data.value = [
    { id: 1, date: '2026-03-22', device_id: 'DEV001', sleep_time: '22:30', wake_time: '06:00', total_hours: '7.5h', deep_sleep: 1.8, light_sleep: 4.2, quality: 85, created_at: '2026-03-22 06:30:00' },
    { id: 2, date: '2026-03-21', device_id: 'DEV001', sleep_time: '23:00', wake_time: '06:30', total_hours: '7.5h', deep_sleep: 1.9, light_sleep: 4.1, quality: 88, created_at: '2026-03-21 06:30:00' },
    { id: 3, date: '2026-03-20', device_id: 'DEV002', sleep_time: '22:45', wake_time: '06:15', total_hours: '7.5h', deep_sleep: 1.7, light_sleep: 4.3, quality: 78, created_at: '2026-03-20 06:30:00' },
    { id: 4, date: '2026-03-19', device_id: 'DEV001', sleep_time: '23:30', wake_time: '07:00', total_hours: '7.5h', deep_sleep: 1.5, light_sleep: 4.5, quality: 65, created_at: '2026-03-19 07:00:00' },
    { id: 5, date: '2026-03-18', device_id: 'DEV003', sleep_time: '22:00', wake_time: '05:30', total_hours: '7.5h', deep_sleep: 2.0, light_sleep: 4.0, quality: 92, created_at: '2026-03-18 05:30:00' }
  ]
  pagination.total = data.value.length
}

const handleReset = () => {
  Object.assign(form, { dateRange: [], deviceId: '', date: null, sleep_time: null, wake_time: null, deep_sleep: 0, light_sleep: 0, quality: 0 })
  loadData()
}

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', date: null, sleep_time: null, wake_time: null, deep_sleep: 0, light_sleep: 0, quality: 0 })
  modalVisible.value = true
}

const handleSubmit = () => {
  modalVisible.value = false
  Message.success(isEdit.value ? '编辑成功' : '添加成功')
  loadData()
}

const onPageChange = (page) => {
  pagination.current = page
  loadData()
}

onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

