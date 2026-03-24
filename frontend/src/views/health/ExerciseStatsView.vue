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
        <a-form-item label="运动类型"><a-input v-model="form.type" /></a-form-item>
        <a-form-item label="时长(分钟)"><a-input-number v-model="form.duration" :min="0" style="width:100%" /></a-form-item>
        <a-form-item label="距离(公里)"><a-input-number v-model="form.distance" :min="0" :precision="1" style="width:100%" /></a-form-item>
        <a-form-item label="消耗(千卡)"><a-input-number v-model="form.calories" :min="0" style="width:100%" /></a-form-item>
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
  type: '',
  duration: 0,
  distance: 0,
  calories: 0
})

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '运动类型', dataIndex: 'type_name', width: 120 },
  { title: '时长(分钟)', dataIndex: 'duration', width: 120 },
  { title: '距离(公里)', dataIndex: 'distance', width: 120 },
  { title: '消耗(千卡)', dataIndex: 'calories', width: 120 },
  { title: '完成度', dataIndex: 'completion', width: 100 }
]

const data = ref([])
const pagination = reactive({ total: 0, current: 1, pageSize: 10 })

const getExerciseTypeName = (type) => {
  const names = { run: '跑步', walk: '步行', cycling: '骑行', swim: '游泳', gym: '健身' }
  return names[type] || '其他'
}

const loadStats = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/exercise/stats?range=${form.name || 'week'}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = (resData.data?.records || []).map(r => ({
        ...r,
        type_name: getExerciseTypeName(r.type)
      }))
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
    { id: 1, date: '2026-03-22', type: 'run', type_name: '跑步', duration: 45, distance: 5.2, calories: 320, completion: '85%' },
    { id: 2, date: '2026-03-21', type: 'walk', type_name: '步行', duration: 60, distance: 4.5, calories: 180, completion: '100%' },
    { id: 3, date: '2026-03-20', type: 'cycling', type_name: '骑行', duration: 90, distance: 18.5, calories: 650, completion: '92%' },
    { id: 4, date: '2026-03-19', type: 'swim', type_name: '游泳', duration: 40, distance: 1.0, calories: 280, completion: '75%' },
    { id: 5, date: '2026-03-18', type: 'gym', type_name: '健身', duration: 60, distance: 0, calories: 400, completion: '80%' }
  ]
}

const handleSearch = () => loadStats()
const handleReset = () => { form.name = ''; loadStats() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', date: '', type: '', duration: 0, distance: 0, calories: 0 })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (isEdit.value) {
    const idx = data.value.findIndex(r => r.id === form.id)
    if (idx !== -1) data.value[idx] = { ...form, type_name: getExerciseTypeName(form.type) }
    Message.success('编辑成功')
  } else {
    data.value.unshift({
      id: Date.now(),
      date: form.date,
      type: form.type,
      type_name: getExerciseTypeName(form.type),
      duration: form.duration,
      distance: form.distance,
      calories: form.calories,
      completion: '0%'
    })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
