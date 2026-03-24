<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="端点URL">
          <a-input v-model="form.url" placeholder="请输入URL" />
        </a-form-item>
        <a-form-item label="事件类型">
          <a-select v-model="form.event_type" placeholder="请选择" allow-clear style="width: 160px">
            <a-option value="device.online">设备上线</a-option>
            <a-option value="device.offline">设备离线</a-option>
            <a-option value="device.data">设备数据</a-option>
            <a-option value="alert.created">告警创建</a-option>
            <a-option value="ota.started">OTA开始</a-option>
            <a-option value="ota.completed">OTA完成</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="form.status" placeholder="请选择" allow-clear style="width: 120px">
            <a-option value="enabled">启用</a-option>
            <a-option value="disabled">禁用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="loadData">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建Webhook</a-button>
      <a-button @click="handleViewLogs">查看日志</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #status="{ record }">
        <a-tag :color="record.status === 'enabled' ? 'green' : 'gray'">
          {{ record.status === 'enabled' ? '启用' : '禁用' }}
        </a-tag>
      </template>
      <template #event_types="{ record }">
        <a-tag v-for="e in (record.event_types || [])" :key="e" size="small">{{ e }}</a-tag>
      </template>
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleToggle(record)">
          {{ record.status === 'enabled' ? '禁用' : '启用' }}
        </a-button>
        <a-button type="text" size="small" @click="handleTest(record)">测试</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>

    <!-- 新建/编辑弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="modalTitle" :width="520" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="formData" layout="vertical" label-col-flex="100px">
        <a-form-item label="端点URL" required>
          <a-input v-model="formData.url" placeholder="https://example.com/webhook" />
        </a-form-item>
        <a-form-item label="事件类型" required>
          <a-select v-model="formData.event_types" multiple placeholder="请选择事件类型">
            <a-option value="device.online">设备上线</a-option>
            <a-option value="device.offline">设备离线</a-option>
            <a-option value="device.data">设备数据</a-option>
            <a-option value="alert.created">告警创建</a-option>
            <a-option value="ota.started">OTA开始</a-option>
            <a-option value="ota.completed">OTA完成</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="签名密钥">
          <a-input v-model="formData.secret" placeholder="用于签名验证" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="formData.description" placeholder="请输入描述" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 日志弹窗 -->
    <a-drawer v-model:visible="logsVisible" title="Webhook 日志" :width="800">
      <WebhookLogs :webhook-id="currentWebhookId" />
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message, Modal } from '@arco-design/web-vue'
import WebhookLogs from './WebhookLogs.vue'

const router = useRouter()
const loading = ref(false)
const data = ref([])
const modalVisible = ref(false)
const logsVisible = ref(false)
const modalTitle = ref('新建Webhook')
const currentWebhookId = ref(null)
const formData = reactive({ url: '', event_types: [], secret: '', description: '' })
const editingId = ref(null)

const form = reactive({ url: '', event_type: undefined, status: undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '端点URL', dataIndex: 'url', ellipsis: true },
  { title: '事件类型', slotName: 'event_types', width: 240 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '最近触发', dataIndex: 'last_triggered_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
]

const loadData = async () => {
  loading.value = true
  try {
    const res = await fetch(`/api/v1/webhooks?page=${pagination.current}&page_size=${pagination.pageSize}`)
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
  modalTitle.value = '新建Webhook'
  editingId.value = null
  Object.assign(formData, { url: '', event_types: [], secret: '', description: '' })
  modalVisible.value = true
}

const handleEdit = (record) => {
  modalTitle.value = '编辑Webhook'
  editingId.value = record.id
  Object.assign(formData, { url: record.url, event_types: record.event_types || [], secret: record.secret || '', description: record.description || '' })
  modalVisible.value = true
}

const handleSubmit = async () => {
  const method = editingId.value ? 'PUT' : 'POST'
  const url = editingId.value ? `/api/v1/webhooks/${editingId.value}` : '/api/v1/webhooks'
  await fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData)
  })
  Message.success('保存成功')
  modalVisible.value = false
  loadData()
}

const handleToggle = async (record) => {
  const newStatus = record.status === 'enabled' ? 'disabled' : 'enabled'
  await fetch(`/api/v1/webhooks/${record.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status: newStatus })
  })
  Message.success(newStatus === 'enabled' ? '启用成功' : '禁用成功')
  loadData()
}

const handleTest = async (record) => {
  await fetch(`/api/v1/webhooks/${record.id}/test`, { method: 'POST' })
  Message.success('测试请求已发送')
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: '删除后无法恢复，确定要删除吗？',
    onOk: async () => {
      await fetch(`/api/v1/webhooks/${record.id}`, { method: 'DELETE' })
      Message.success('删除成功')
      loadData()
    }
  })
}

const handleViewLogs = () => {
  router.push('/webhooks/logs')
}

const handleReset = () => { form.url = ''; form.event_type = undefined; form.status = undefined; loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData() })
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
