<template>
  <div class="page-container">
    <a-card class="general-card" title="告警设置">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新增设置</a-button>
      </template>
      <div class="search-form">
        <a-form :model="searchForm" layout="inline">
          <a-row :gutter="16" style="width: 100%">
            <a-col :span="8">
              <a-form-item label="设置名称"><a-input v-model="searchForm.keyword" placeholder="搜索设置项" /></a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="类型">
                <a-select v-model="searchForm.type" placeholder="请选择" allow-clear>
                  <a-option value="email">邮件</a-option>
                  <a-option value="sms">短信</a-option>
                  <a-option value="webhook">Webhook</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="8" style="display: flex; justify-content: flex-end">
              <a-form-item><a-button type="primary" @click="handleSearch">查询</a-button><a-button style="margin-left: 8px" @click="handleReset">重置</a-button></a-form-item>
            </a-col>
          </a-row>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" />
      </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="设置名称"><a-input v-model="form.name" /></a-form-item>
        <a-form-item label="类型"><a-select v-model="form.type" style="width: 200px">
          <a-option value="email">邮件</a-option>
          <a-option value="sms">短信</a-option>
          <a-option value="webhook">Webhook</a-option>
        </a-select></a-form-item>
        <a-form-item label="配置值"><a-input v-model="form.config_value" /></a-form-item>
        <a-form-item label="描述"><a-input v-model="form.description" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
    </a-card>`n</div></template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const modalTitle = ref('新增设置')

const searchForm = reactive({
  keyword: '',
  type: ''
})

const form = reactive({
  name: '',
  type: 'email',
  config_value: '',
  description: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '设置名称', dataIndex: 'name', width: 200 },
  { title: '类型', dataIndex: 'type', width: 120 },
  { title: '配置值', dataIndex: 'config_value', ellipsis: true },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '启用状态', dataIndex: 'enabled', width: 100, render: ({ record }) => record.enabled ? '是' : '否' }
]

const handleSearch = () => {
  loadData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  loadData()
}

const handleCreate = () => {
  modalTitle.value = '新增设置'
  modalVisible.value = true
}

const handleSubmit = async () => {
  modalVisible.value = false
  try {
    const token = localStorage.getItem('token')
    await fetch('/api/v1/alerts/settings', {
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
    const res = await fetch('/api/v1/alerts/settings', { headers: { 'Authorization': `Bearer ${token}` } })
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
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

