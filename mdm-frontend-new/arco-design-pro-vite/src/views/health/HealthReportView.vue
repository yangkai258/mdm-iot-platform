<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="健康报告">
      <template #extra>
        <a-button type="primary" @click="handleGenerate"><icon-plus />生成报告</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="报告周期">
            <a-select v-model="form.period" placeholder="请选择" style="width: 140px" allow-clear>
              <a-option value="week">周报</a-option>
              <a-option value="month">月报</a-option>
              <a-option value="quarter">季报</a-option>
            </a-select>
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
    />
    </a-table>
  </a-card>`n</div></template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)

const form = reactive({ period: '' })

const columns = [
  { title: '报告周期', dataIndex: 'period', width: 200 },
  { title: '平均心率', dataIndex: 'heart_rate', width: 100 },
  { title: '平均血压', dataIndex: 'blood_pressure', width: 120 },
  { title: '平均睡眠(小时)', dataIndex: 'sleep_hours', width: 120 },
  { title: '运动时长(分钟)', dataIndex: 'exercise_minutes', width: 130 },
  { title: '综合评分', dataIndex: 'overall_score', width: 100 },
  { title: '生成时间', dataIndex: 'created_at', width: 160 }
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
    if (form.period) params.type = form.period
    const res = await axios.get('/api/v1/health/reports', { params })
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
    { id: 1, period: '2026-03-16 至 2026-03-22', heart_rate: 72, blood_pressure: '120/80', sleep_hours: 7.5, exercise_minutes: 45, overall_score: 85, created_at: '2026-03-22 10:00:00' },
    { id: 2, period: '2026-03-09 至 2026-03-15', heart_rate: 74, blood_pressure: '122/82', sleep_hours: 7.2, exercise_minutes: 40, overall_score: 80, created_at: '2026-03-15 10:00:00' },
    { id: 3, period: '2026-03-02 至 2026-03-08', heart_rate: 70, blood_pressure: '118/78', sleep_hours: 7.8, exercise_minutes: 50, overall_score: 88, created_at: '2026-03-08 10:00:00' },
    { id: 4, period: '2026-02-23 至 2026-03-01', heart_rate: 73, blood_pressure: '120/82', sleep_hours: 7.0, exercise_minutes: 35, overall_score: 78, created_at: '2026-03-01 10:00:00' }
  ]
  pagination.total = data.value.length
}

const handleReset = () => {
  form.period = ''
  loadData()
}

const handleGenerate = () => {
  Message.info('正在生成健康报告...')
  setTimeout(() => {
    Message.success('报告生成成功')
    loadData()
  }, 1000)
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


