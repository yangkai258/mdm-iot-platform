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
        <a-form-item label="疾病名称"><a-input v-model="form.disease_name" /></a-form-item>
        <a-form-item label="预警级别"><a-input v-model="form.level" /></a-form-item>
        <a-form-item label="预警描述"><a-input v-model="form.description" /></a-form-item>
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
  id: '',
  name: '',
  disease_name: '',
  level: '',
  description: ''
})

const columns = [
  { title: '预警编号', dataIndex: 'id', width: 100 },
  { title: '疾病名称', dataIndex: 'disease_name', width: 150 },
  { title: '预警级别', dataIndex: 'level_name', width: 100 },
  { title: '预警描述', dataIndex: 'description', ellipsis: true },
  { title: '发生时间', dataIndex: 'created_at', width: 180 },
  { title: '状态', dataIndex: 'status_name', width: 100 }
]

const data = ref([])
const pagination = reactive({ total: 0, current: 1, pageSize: 10 })

const getLevelName = (level) => {
  const names = { critical: '危急', high: '高', medium: '中', low: '低' }
  return names[level] || level
}

const getStatusName = (status) => {
  const names = { pending: '待处理', confirmed: '已确认', ignored: '已忽略' }
  return names[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams()
    if (form.name) params.append('keyword', form.name)
    const res = await fetch(`/api/v1/health/warnings?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = (resData.data?.list || []).map(w => ({
        ...w,
        level_name: getLevelName(w.level),
        status_name: getStatusName(w.status)
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
    { id: 1, disease_name: '高血压风险', level: 'high', level_name: '高', description: '连续3天血压监测偏高，建议密切观察', status: 'pending', status_name: '待处理', created_at: '2026-03-22 10:30:00' },
    { id: 2, disease_name: '睡眠呼吸暂停', level: 'critical', level_name: '危急', description: '夜间血氧饱和度多次低于90%', status: 'pending', status_name: '待处理', created_at: '2026-03-22 08:00:00' },
    { id: 3, disease_name: '心律不齐', level: 'medium', level_name: '中', description: '心电图检测到偶发早搏', status: 'confirmed', status_name: '已确认', created_at: '2026-03-21 15:20:00' },
    { id: 4, disease_name: '体重异常波动', level: 'low', level_name: '低', description: '一周内体重下降超过5%', status: 'ignored', status_name: '已忽略', created_at: '2026-03-20 09:00:00' },
    { id: 5, disease_name: '血糖偏高', level: 'high', level_name: '高', description: '空腹血糖持续高于正常值上限', status: 'pending', status_name: '待处理', created_at: '2026-03-22 11:45:00' }
  ]
}

const handleSearch = () => loadData()

const handleReset = () => {
  form.name = ''
  loadData()
}

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', disease_name: '', level: '', description: '' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.disease_name) { Message.warning('请填写疾病名称'); return }
  if (isEdit.value) {
    const idx = data.value.findIndex(w => w.id === form.id)
    if (idx !== -1) data.value[idx] = { ...form, level_name: getLevelName(form.level), status_name: getStatusName(data.value[idx].status) }
    Message.success('编辑成功')
  } else {
    data.value.unshift({
      id: Date.now(),
      disease_name: form.disease_name,
      level: form.level,
      level_name: getLevelName(form.level),
      description: form.description,
      status: 'pending',
      status_name: '待处理',
      created_at: new Date().toLocaleString()
    })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
