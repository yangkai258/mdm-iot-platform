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
        <a-form-item label="动作名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="分类"><a-select v-model="form.category_id" placeholder="选择分类" style="width: 100%">
          <a-option v-for="c in categories" :key="c.category_id" :value="c.category_id">{{ c.category_name }}</a-option>
        </a-select></a-form-item>
        <a-form-item label="时长(秒)"><a-input-number v-model="form.duration" :min="1" style="width:100%" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="2" /></a-form-item>
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
  category_id: null,
  preview_icon: '',
  duration: 5,
  description: ''
})

const columns = [
  { title: '动作名称', dataIndex: 'name' },
  { title: '分类', dataIndex: 'category_name', width: 120 },
  { title: '来源', dataIndex: 'source_name', width: 80 },
  { title: '时长(秒)', dataIndex: 'duration', width: 80 },
  { title: '使用次数', dataIndex: 'usage_count', width: 100 },
  { title: '状态', dataIndex: 'status_name', width: 80 }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 12 })
const data = ref([])
const categories = ref([
  { category_id: 'ac1', category_name: '舞蹈' },
  { category_id: 'ac2', category_name: '问候' },
  { category_id: 'ac3', category_name: '卖萌' },
  { category_id: 'ac4', category_name: '休息' },
  { category_id: 'ac5', category_name: '互动' }
])

const loadActions = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/market/actions?keyword=${form.name || ''}`)
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = (resData.data?.list || []).map(a => ({
        ...a,
        source_name: a.is_official ? '官方' : '用户',
        status_name: a.is_published ? '已发布' : '未发布'
      }))
    } else {
      loadMockData()
    }
  } catch {
    loadMockData()
  } finally {
    pagination.total = data.value.length
    loading.value = false
  }
}

const loadMockData = () => {
  data.value = [
    { action_id: 'act1', name: '开心舞', category_id: 'ac1', category_name: '舞蹈', is_official: true, source_name: '官方', duration: 8, usage_count: 2345, is_published: true, status_name: '已发布' },
    { action_id: 'act2', name: '挥手问候', category_id: 'ac2', category_name: '问候', is_official: true, source_name: '官方', duration: 3, usage_count: 5678, is_published: true, status_name: '已发布' },
    { action_id: 'act3', name: '撒娇', category_id: 'ac3', category_name: '卖萌', is_official: false, source_name: '用户', duration: 5, usage_count: 890, is_published: false, status_name: '未发布' },
    { action_id: 'act4', name: '打盹', category_id: 'ac4', category_name: '休息', is_official: true, source_name: '官方', duration: 10, usage_count: 3456, is_published: true, status_name: '已发布' }
  ]
}

const handleSearch = () => loadActions()
const handleReset = () => { form.name = ''; loadActions() }

const handleCreate = () => {
  isEdit.value = false
  modalTitle.value = '新建'
  Object.assign(form, { id: '', name: '', category_id: null, preview_icon: '', duration: 5, description: '' })
  modalVisible.value = true
}

const handleSubmit = () => {
  if (!form.name) { Message.warning('请填写动作名称'); return }
  if (isEdit.value) {
    const idx = data.value.findIndex(a => a.action_id === form.id)
    if (idx !== -1) data.value[idx] = { ...data.value[idx], name: form.name }
    Message.success('编辑成功')
  } else {
    data.value.unshift({
      action_id: `act${Date.now()}`,
      name: form.name,
      category_id: form.category_id,
      category_name: categories.value.find(c => c.category_id === form.category_id)?.category_name || '默认',
      source_name: '用户',
      duration: form.duration,
      usage_count: 0,
      is_published: false,
      status_name: '未发布'
    })
    pagination.total++
    Message.success('添加成功')
  }
  modalVisible.value = false
}

onMounted(() => { loadActions() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
