<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="规则名称"><a-input v-model="form.keyword" placeholder="搜索规则名称" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建规则</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="规则名称"><a-input v-model="form.name" /></a-form-item>
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
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新建规则')

const form = reactive({
  name: '',
  alert_type: 'battery_low',
  condition: '<',
  threshold: 20,
  severity: 2
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '规则名称', dataIndex: 'name', width: 200 },
  { title: '告警类型', dataIndex: 'alert_type', width: 150 },
  { title: '条件', dataIndex: 'condition', width: 80 },
  { title: '阈值', dataIndex: 'threshold', width: 80 },
  { title: '严重程度', dataIndex: 'severity', width: 100 },
  { title: '启用', dataIndex: 'enabled', width: 100 }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  form.keyword = ''
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '新建规则'
  modalVisible.value = true
}

const handleSubmit = async () => {
  modalVisible.value = false
  try {
    const token = localStorage.getItem('token')
    await fetch('/api/v1/alerts/rules', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    Message.success('保存成功')
    loadData()
  } catch (e) {
    Message.error('操作失败')
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/alerts/rules', { headers: { 'Authorization': `Bearer ${token}` } })
    const resData = await res.json()
    if (resData.code === 0) {
      data.value = resData.data?.list || []
      pagination.total = data.value.length
    }
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
