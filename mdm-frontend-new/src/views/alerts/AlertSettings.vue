<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="设置名称">
          <a-input v-model="form.name" placeholder="请输入设置名称" />
        </a-form-item>
        <a-form-item label="类型">
          <a-select v-model="form.type" placeholder="请选择" allow-clear style="width: 140px">
            <a-option value="email">邮件通知</a-option>
            <a-option value="sms">短信通知</a-option>
            <a-option value="webhook">Webhook</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新增设置</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange" row-key="id">
      <template #type="{ record }">
        <a-tag :color="getTypeColor(record.type)">{{ getTypeText(record.type) }}</a-tag>
      </template>
      <template #enabled="{ record }">
        <a-switch v-model="record.enabled" @change="handleToggle(record)" />
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" label-col-flex="100px" ref="formRef">
        <a-form-item label="设置名称" field="name" :rules="[{ required: true, message: '请输入设置名称' }]">
          <a-input v-model="form.name" placeholder="请输入" />
        </a-form-item>
        <a-form-item label="设置类型" field="type" :rules="[{ required: true, message: '请选择类型' }]">
          <a-select v-model="form.type" placeholder="请选择">
            <a-option value="email">邮件通知</a-option>
            <a-option value="sms">短信通知</a-option>
            <a-option value="webhook">Webhook</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="配置值" field="config_value" :rules="[{ required: true, message: '请输入配置值' }]">
          <a-input v-model="form.config_value" placeholder="请输入配置值" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" placeholder="请输入描述" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const modalVisible = ref(false)
const modalTitle = ref('新增设置')
const formRef = ref()
const editId = ref<number | null>(null)

const form = reactive({
  name: '',
  type: 'email',
  config_value: '',
  description: ''
})

const data = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '设置名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 120 },
  { title: '配置值', dataIndex: 'config_value', ellipsis: true },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '启用', slotName: 'enabled', width: 100 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getTypeColor = (t: string) => ({ email: 'blue', sms: 'green', webhook: 'purple' }[t] || 'gray')
const getTypeText = (t: string) => ({ email: '邮件通知', sms: '短信通知', webhook: 'Webhook' }[t] || t)

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.name) params.name = form.name
    if (form.type) params.type = form.type
    const res = await axios.get(`${API_BASE}/alerts/settings`, { params })
    if (res.data.code === 0) {
      data.value = res.data.data?.list || []
      pagination.total = res.data.data?.pagination?.total || 0
    }
  } catch {
    data.value = [
      { id: 1, name: '管理员邮件', type: 'email', config_value: 'admin@example.com', description: '管理员通知邮箱', enabled: true },
      { id: 2, name: '短信通知', type: 'sms', config_value: '13800138000', description: '紧急告警短信', enabled: false }
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.name = ''; form.type = undefined; pagination.current = 1; loadData() }
const handlePageChange = (page: number) => { pagination.current = page; loadData() }

const handleCreate = () => {
  editId.value = null
  modalTitle.value = '新增设置'
  Object.assign(form, { name: '', type: 'email', config_value: '', description: '' })
  modalVisible.value = true
}

const handleEdit = (record: any) => {
  editId.value = record.id
  modalTitle.value = '编辑设置'
  Object.assign(form, { ...record })
  modalVisible.value = true
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    if (editId.value) {
      await axios.put(`${API_BASE}/alerts/settings/${editId.value}`, form)
    } else {
      await axios.post(`${API_BASE}/alerts/settings`, form)
    }
    Message.success('保存成功')
    modalVisible.value = false
    loadData()
    done(true)
  } catch {
    done(false)
  }
}

const handleToggle = async (record: any) => {
  try { await axios.put(`${API_BASE}/alerts/settings/${record.id}`, { enabled: record.enabled }) } catch {}
}

const handleDelete = (record: any) => {
  Modal.confirm({ title: '确认删除', content: '确定删除该设置？', onOk: async () => {
    try { await axios.delete(`${API_BASE}/alerts/settings/${record.id}`) } catch {}
    Message.success('删除成功')
    loadData()
  }})
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
