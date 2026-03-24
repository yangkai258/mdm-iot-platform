<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="检测项"><a-input v-model="form.check_item" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">运行测试</a-button>
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

const API_BASE = '/api/v1'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建')

const form = reactive({
  check_item: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '检测项', dataIndex: 'check_item', width: 200 },
  { title: '受影响群体', dataIndex: 'affected_group', width: 120 },
  { title: '偏差程度', dataIndex: 'bias_score', width: 100 },
  { title: '风险等级', dataIndex: 'level', width: 100 },
  { title: '建议', dataIndex: 'suggestion', ellipsis: true }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.check_item = ''
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '运行测试'
  modalVisible.value = true
}

const handleSubmit = async () => {
  modalVisible.value = false
  try {
    await fetch(`${API_BASE}/ai-fairness/run`, { method: 'POST' })
    Message.success('测试已启动')
    setTimeout(loadData, 2000)
  } catch (e) {
    Message.error('启动失败')
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/ai-fairness/bias-detection`)
    const resData = await res.json()
    data.value = resData.bias_results || []
    pagination.total = data.value.length
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
