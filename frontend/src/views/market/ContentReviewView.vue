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

const form = reactive({
  id: '',
  name: '',
  content_type_name: '',
  developer_name: '',
  status_name: ''
})

const columns = [
  { title: '内容名称', dataIndex: 'name' },
  { title: '内容类型', dataIndex: 'content_type_name', width: 100 },
  { title: '开发者', dataIndex: 'developer_name', width: 120 },
  { title: '提交时间', dataIndex: 'submitted_at', width: 160 },
  { title: '状态', dataIndex: 'status_name', width: 100 }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 10 })
const data = ref([])

const getTypeName = (type) => {
  const names = { plugin: '插件', emoticon: '表情包', action: '动作', voice: '声音' }
  return names[type] || type
}

const getStatusName = (status) => {
  const names = { pending: '待审核', approved: '已通过', rejected: '已拒绝' }
  return names[status] || status
}

const loadReviews = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/market/reviews?keyword=${form.name || ''}`)
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = (resData.data?.list || []).map(r => ({
        ...r,
        content_type_name: getTypeName(r.content_type),
        status_name: getStatusName(r.status)
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
    { id: 'rev1', name: '可爱猫咪表情包', content_type: 'emoticon', content_type_name: '表情包', developer_name: '用户小明', submitted_at: '2026-03-22 10:30:00', status: 'pending', status_name: '待审核' },
    { id: 'rev2', name: '夜灯控制插件', content_type: 'plugin', content_type_name: '插件', developer_name: '开发者张三', submitted_at: '2026-03-22 09:15:00', status: 'pending', status_name: '待审核' },
    { id: 'rev3', name: '跳舞动作', content_type: 'action', content_type_name: '动作', developer_name: '用户小红', submitted_at: '2026-03-21 16:45:00', status: 'approved', status_name: '已通过' },
    { id: 'rev4', name: '童声朗读', content_type: 'voice', content_type_name: '声音', developer_name: '用户小刚', submitted_at: '2026-03-21 14:20:00', status: 'rejected', status_name: '已拒绝' },
    { id: 'rev5', name: '节日烟花动作', content_type: 'action', content_type_name: '动作', developer_name: '官方', submitted_at: '2026-03-20 11:00:00', status: 'approved', status_name: '已通过' }
  ]
}

const handleSearch = () => loadReviews()
const handleReset = () => { form.name = ''; loadReviews() }

const handleCreate = () => {
  modalTitle.value = '新建'
  modalVisible.value = true
}

const handleSubmit = () => {
  modalVisible.value = false
}

onMounted(() => { loadReviews() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
