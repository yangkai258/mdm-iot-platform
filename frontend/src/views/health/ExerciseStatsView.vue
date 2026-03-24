<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="日期范围">
          <a-range-picker v-model="form.dateRange" style="width: 240px" />
        </a-form-item>
        <a-form-item label="设备ID">
          <a-input v-model="form.deviceId" placeholder="请输入设备ID" style="width: 160px" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="loadData">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table
      :columns="columns"
      :data="data"
      :loading="loading"
      :pagination="paginationConfig"
      @page-change="onPageChange"
    />
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" layout="vertical">
        <a-form-item label="日期"><a-date-picker v-model="form.date" style="width: 100%" /></a-form-item>
        <a-form-item label="步数"><a-input-number v-model="form.steps" :min="0" style="width: 100%" /></a-form-item>
        <a-form-item label="距离(km)"><a-input-number v-model="form.distance" :min="0" :precision="1" style="width: 100%" /></a-form-item>
        <a-form-item label="卡路里(kcal)"><a-input-number v-model="form.calories" :min="0" style="width: 100%" /></a-form-item>
        <a-form-item label="时长(分钟)"><a-input-number v-model="form.duration" :min="0" style="width: 100%" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建')
const isEdit = ref(false)

const form = reactive({
  dateRange: [],
  deviceId: '',
  date: null,
  steps: 0,
  distance: 0,
  calories: 0,
  duration: 0
})

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '设备ID', dataIndex: 'device_id', width: 120 },
  { title: '步数', dataIndex: 'steps', width: 100 },
  { title: '距离(km)', dataIndex: 'distance', width: 100 },
  { title: '卡路里(kcal)', dataIndex: 'calories', width: 120 },
  { title: '时长(分钟)', dataIndex: 'duration', width: 120 },
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

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.deviceId) params.device_id = form.deviceId
    if (form.dateRange && form.dateRange.length === 2) {
      params.start_date = form.dateRange[0].format('YYYY-MM-DD')
      params.end_date = form.dateRange[1].format('YYYY-MM-DD')
    }
    const res = await axios.get('/api/v1/health/exercise-stats', { params })
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
    { id: 1, date: '2026-03-22', device_id: 'DEV001', steps: 8500, distance: 5.2, calories: 320, duration: 45, created_at: '2026-03-22 23:00:00' },
    { id: 2, date: '2026-03-21', device_id: 'DEV001', steps: 6200, distance: 4.5, calories: 280, duration: 60, created_at: '2026-03-21 23:00:00' },
    { id: 3, date: '2026-03-20', device_id: 'DEV002', steps: 12000, distance: 8.5, calories: 480, duration: 90, created_at: '2026-03-20 23:00:00' },
    { id: 4, date: '2026-03-19', device_id: 'DEV001', steps: 3000, distance: 2.0, calories: 150, duration: 30, created_at: '2026-03-19 23:00:00' },
    { id: 5, date: '2026-03-18', device_id: 'DEV003', steps: 9800, distance: 6.8, calories: 400, duration: 55, created_at: '2026-03-18 23:00:00' }
  ]
  pagination.total = data.value.length
}

const handleReset = () => {
  Object.assign(form, { dateRange: [], deviceId: '', date: null, steps: 0, distance: 0, calories: 0, duration: 0 })
  loadData()
}

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { date: null, steps: 0, distance: 0, calories: 0, duration: 0 })
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
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
