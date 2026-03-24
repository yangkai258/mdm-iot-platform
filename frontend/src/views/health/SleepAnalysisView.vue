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
        <a-form-item label="日期"><a-input v-model="form.date" /></a-form-item>
        <a-form-item label="入睡时间"><a-input v-model="form.sleep_time" /></a-form-item>
        <a-form-item label="起床时间"><a-input v-model="form.wake_time" /></a-form-item>
        <a-form-item label="深睡(小时)"><a-input-number v-model="form.deep_sleep" :min="0" style="width:100%" /></a-form-item>
        <a-form-item label="浅睡(小时)"><a-input-number v-model="form.light_sleep" :min="0" style="width:100%" /></a-form-item>
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
const isEdit = ref(false)

const form = reactive({
  name: '',
  id: '',
  date: '',
  sleep_time: '',
  wake_time: '',
  deep_sleep: 0,
  light_sleep: 0
})

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '入睡时间', dataIndex: 'sleep_time', width: 120 },
  { title: '起床时间', dataIndex: 'wake_time', width: 120 },
  { title: '睡眠时长', dataIndex: 'total_hours', width: 120 },
  { title: '深睡(小时)', dataIndex: 'deep_sleep', width: 100 },
  { title: '浅睡(小时)', dataIndex: 'light_sleep', width: 100 },
  { title: '睡眠质量', dataIndex: 'quality', width: 100 }
]

const data = ref([])
const pagination = reactive({ total: 0, current: 1, pageSize: 10 })

const loadStats = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/sleep/stats?range=${form.name || 'week'}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = resData.data?.records || []
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
    { id: 1, date: '2026-03-22', sleep_time: '22:30', wake_time: '06:00', total_hours: '7.5h', deep_sleep: 1.8, light_sleep: 4.2, quality: 85 },
    { id: 2, date: '2026-03-21', sleep_time: '23:00', wake_time: '06:30', total_hours: '7.5h', deep_sleep: 1.9, light_sleep: 4.1, quality: 88 },
    { id: 3, date: '2026-03-20', sleep_time: '22:45', wake_time: '06:15', total_hours: '7.5h', deep_sleep: 1.7, light_sleep: 4.3, quality: 78 },
    { id: 4, date: '2026-03-19', sleep_time: '23:30', wake_time: '07:00', total_hours: '7.5h', deep_sleep: 1.5, light_sleep: 4.5, quality: 65 },
    { id: 5, date: '2026-03-18', sleep_time: '22:00', wake_time: '05:30', total_hours: '7.5h', deep_sleep: 2.0, light_sleep: 4.0, quality: 92 }
  ]
}

const handleSearch = () => loadStats()
const handleReset = () => { form.name = ''; loadStats() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', date: '', sleep_time: '', wake_time: '', deep_sleep: 0, light_sleep: 0 })
  modalVisible.value = true
}

const handleSubmit = () => {
  const totalHours = (form.deep_sleep + form.light_sleep).toFixed(1) + 'h'
  if (isEdit.value) {
    const idx = data.value.findIndex(r => r.id === form.id)
    if (idx !== -1) data.value[idx] = { ...form, total_hours: totalHours }
    Message.success('编辑成功')
  } else {
    data.value.unshift({
      id: Date.now(),
      date: form.date,
      sleep_time: form.sleep_time,
      wake_time: form.wake_time,
      total_hours: totalHours,
      deep_sleep: form.deep_sleep,
      light_sleep: form.light_sleep,
      quality: 0
    })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => { loadStats() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
