<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({ name: '' })

const columns = [
  { title: '报告周期', dataIndex: 'period', width: 200 },
  { title: '平均心率', dataIndex: 'heart_rate', width: 100 },
  { title: '平均血压', dataIndex: 'blood_pressure', width: 100 },
  { title: '平均睡眠', dataIndex: 'sleep_hours', width: 100 },
  { title: '运动时长', dataIndex: 'exercise_minutes', width: 100 },
  { title: '综合评分', dataIndex: 'overall_score', width: 100 }
]

const data = ref([])
const pagination = reactive({ total: 0, current: 1, pageSize: 10 })

const loadReport = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/report?type=${form.name || 'week'}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = resData.data?.reports || []
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    loading.value = false
  }
  pagination.total = data.value.length
}

const loadMockData = () => {
  data.value = [
    { period: '2026-03-16 至 2026-03-22', heart_rate: 72, blood_pressure: '120/80', sleep_hours: 7.5, exercise_minutes: 45, overall_score: 85 },
    { period: '2026-03-09 至 2026-03-15', heart_rate: 74, blood_pressure: '122/82', sleep_hours: 7.2, exercise_minutes: 40, overall_score: 80 },
    { period: '2026-03-02 至 2026-03-08', heart_rate: 70, blood_pressure: '118/78', sleep_hours: 7.8, exercise_minutes: 50, overall_score: 88 }
  ]
}

const handleSearch = () => loadReport()
const handleReset = () => { form.name = ''; loadReport() }

const handleCreate = () => {
  modalTitle.value = '新建'
  modalVisible.value = true
}

const handleSubmit = () => {
  modalVisible.value = false
}

onMounted(() => { loadReport() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
