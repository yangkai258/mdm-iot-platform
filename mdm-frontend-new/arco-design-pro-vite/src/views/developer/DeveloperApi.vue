<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <a-card class="general-card" title="开发者应用">
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />创建应用</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="应用名称">
            <a-input v-model="form.app_name" placeholder="请输入应用名称" />
          </a-form-item>
          <a-form-item label="状态">
            <a-select v-model="form.status" placeholder="请选择" allow-clear style="width: 120px">
              <a-option value="active">启用</a-option>
              <a-option value="disabled">禁用</a-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="loadData">查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #status="{ record }">
        <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
          {{ record.status === 'active' ? '启用' : '禁用' }}
        </a-tag>
      </template>
      </a-table>
      <template #scopes="{ record }">
        <a-tag v-for="s in (record.scopes || [])" :key="s" size="small">{{ s }}</a-tag>
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleViewKey(record)" v-if="record.api_key">查看Key</a-button>
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleToggle(record)">
          {{ record.status === 'active' ? '禁用' : '启用' }}
        </a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="520" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="formData" layout="vertical" label-col-flex="100px">
        <a-form-item label="应用名称" required>
          <a-input v-model="formData.app_name" placeholder="请输入应用名称" />
        </a-form-item>
        <a-form-item label="权限范围" required>
          <a-select v-model="formData.scopes" multiple placeholder="请选择权限">
            <a-option value="device:read">设备读取</a-option>
            <a-option value="device:write">设备写入</a-option>
            <a-option value="device:control">设备控制</a-option>
            <a-option value="ota:read">OTA读取</a-option>
            <a-option value="ota:write">OTA写入</a-option>
            <a-option value="alert:read">告警读取</a-option>
            <a-option value="alert:write">告警写入</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="formData.description" placeholder="请输入描述" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 查看Key弹窗 -->
    <a-modal v-model:visible="keyVisible" title="API Key" :width="520">
      <a-alert type="warning" style="margin-bottom: 16px">
        请妥善保管您的 API Key，关闭弹窗后将无法再次查看完整 Key。
      </a-alert>
      <a-form :model="keyRecord" layout="vertical" label-col-flex="100px">
        <a-form-item label="应用名称">{{ keyRecord?.app_name }}</a-form-item>
        <a-form-item label="API Key">
          <a-input-group>
            <a-input v-model="keyRecord.api_key" readonly :style="{ fontFamily: 'monospace' }" />
            <template #append>
              <a-button @click="handleCopy">复制</a-button>
            </template>
          </a-input-group>
        </a-form-item>
      </a-form>
    </a-modal>
    </a-card>`n</div></template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const keyVisible = ref(false)
const modalTitle = ref('创建应用')
const editingId = ref(null)
const keyRecord = ref(null)
const formData = reactive({ app_name: '', scopes: [], description: '' })
const form = reactive({ app_name: '', status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '应用名称', dataIndex: 'app_name', width: 160 },
  { title: 'API Key', dataIndex: 'api_key_preview', width: 200, ellipsis: true },
  { title: '权限范围', slotName: 'scopes', width: 280 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 240, fixed: 'right' },
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/developer/apps?page=${pagination.current}&page_size=${pagination.pageSize}`)
    const json = await res.json()
    data.value = json.data?.list || json.data || []
    pagination.total = json.data?.total || 0
  } catch {
    data.value = []
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  modalTitle.value = '创建应用'
  editingId.value = null
  Object.assign(formData, { app_name: '', scopes: [], description: '' })
  modalVisible.value = true
}

const handleEdit = (record) => {
  modalTitle.value = '编辑应用'
  editingId.value = record.id
  Object.assign(formData, { app_name: record.app_name, scopes: record.scopes || [], description: record.description || '' })
  modalVisible.value = true
}

const handleSubmit = async () => {
  const method = editingId.value ? 'PUT' : 'POST'
  const url = editingId.value ? `/api/v1/developer/apps/${editingId.value}` : '/api/v1/developer/apps'
  const res = await fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData)
  })
  const json = await res.json()
  Message.success(editingId.value ? '更新成功' : '创建成功')
  if (!editingId.value && json.data?.api_key) {
    keyRecord.value = json.data
    keyVisible.value = true
  }
  modalVisible.value = false
  loadData()
}

const handleViewKey = (record) => {
  keyRecord.value = record
  keyVisible.value = true
}

const handleCopy = () => {
  navigator.clipboard.writeText(keyRecord.value?.api_key || '')
  Message.success('已复制到剪贴板')
}

const handleToggle = async (record) => {
  const newStatus = record.status === 'active' ? 'disabled' : 'active'
  await fetch(`/api/v1/developer/apps/${record.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status: newStatus })
  })
  Message.success(newStatus === 'active' ? '启用成功' : '禁用成功')
  loadData()
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: '删除后无法恢复，确定要删除吗？',
    onOk: async () => {
      await fetch(`/api/v1/developer/apps/${record.id}`, { method: 'DELETE' })
      Message.success('删除成功')
      loadData()
    }
  })
}

const handleReset = () => { form.app_name = ''; form.status = undefined; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>

