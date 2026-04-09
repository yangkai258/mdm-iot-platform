<template>
  <div class="page-container">
    <a-card class="general-card" title="动作�?>
      <template #extra>
        <a-button type="primary" @click="handleCreate"><icon-plus />新建动作</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="动作名称"><a-input v-model="form.action_name" placeholder="请输�? /></a-form-item>
          <a-form-item label="类型">
            <a-select v-model="form.action_type" placeholder="选择类型" allow-clear style="width: 120px">
              <a-option value="built-in">内置</a-option>
              <a-option value="learned">学习</a-option>
              <a-option value="custom">自定�?/a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="handleSearch">查询</a-button><a-button @click="handleReset">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="handleTest(record)">测试</a-button>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" :width="560">
      <a-form :model="form" layout="vertical">
        <a-form-item label="动作名称" required><a-input v-model="form.action_name" placeholder="请输入动作名�? /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model="form.description" :rows="3" placeholder="动作描述" /></a-form-item>
        <a-form-item label="触发条件"><a-input v-model="form.trigger_condition" placeholder="�? temperature > 30" /></a-form-item>
        <a-form-item label="类型">
          <a-select v-model="form.action_type" placeholder="选择类型">
            <a-option value="built-in">内置</a-option>
            <a-option value="learned">学习</a-option>
            <a-option value="custom">自定�?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="难度">
          <a-select v-model="form.difficulty" placeholder="选择难度">
            <a-option value="easy">简�?/a-option>
            <a-option value="medium">中等</a-option>
            <a-option value="hard">困难</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签">
          <a-select v-model="form.tags" multiple placeholder="选择标签" allow-create>
            <a-option value="dance">舞蹈</a-option>
            <a-option value="greeting">问�?/a-option>
            <a-option value="exercise">运动</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="testModalVisible" title="动作测试" :width="480">
      <a-form layout="vertical">
        <a-form-item label="选择设备">
          <a-select v-model="testDeviceId" placeholder="选择设备" style="width: 100%">
            <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }} ({{ d.id }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="测试动作">{{ testAction?.action_name }}</a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="testModalVisible = false">取消</a-button>
        <a-button type="primary" :loading="testing" @click="executeTest">执行测试</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const testing = ref(false)
const data = ref([])
const devices = ref([])
const modalVisible = ref(false)
const testModalVisible = ref(false)
const modalTitle = ref('新建动作')
const isEdit = ref(false)
const testAction = ref(null)
const testDeviceId = ref('')
const form = reactive({ id: null, action_name: '', description: '', action_type: 'custom', difficulty: 'medium', trigger_condition: '', tags: [] })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '动作名称', dataIndex: 'action_name', width: 160 },
  { title: '类型', dataIndex: 'action_type', width: 100 },
  { title: '难度', dataIndex: 'difficulty', width: 90 },
  { title: '触发条件', dataIndex: 'trigger_condition', ellipsis: true },
  { title: '评分', dataIndex: 'score', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.action_name) params.append('keyword', form.action_name)
    if (form.action_type) params.append('action_type', form.action_type)
    const res = await fetch(`/api/v1/ai/action-library?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const loadDevices = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/devices?page_size=200', { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    devices.value = res.data?.list || []
  } catch (e) { console.error('加载设备失败', e) }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { Object.assign(form, { action_name: '', action_type: '', difficulty: 'medium', trigger_condition: '', tags: [] }); loadData() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const handleCreate = () => { isEdit.value = false; Object.assign(form, { id: null, action_name: '', description: '', action_type: 'custom', difficulty: 'medium', trigger_condition: '', tags: [] }); modalVisible.value = true }
const handleEdit = (record) => { isEdit.value = true; Object.assign(form, record); modalVisible.value = true }

const handleSubmit = async (done) => {
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/v1/ai/action-library/${form.id}` : '/api/v1/ai/action-library'
    const res = await fetch(url, { method: isEdit.value ? 'PUT' : 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) }).then(r => r.json())
    if (res.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadData() }
    else { Message.error(res.message || '操作失败') }
    done(true)
  } catch (e) { Message.error('操作失败'); done(false) }
}

const handleDelete = async (record) => {
  try {
    const token = localStorage.getItem('token')
    await fetch(`/api/v1/ai/action-library/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    Message.success('删除成功'); loadData()
  } catch (e) { Message.error('删除失败') }
}

const handleTest = (record) => { testAction.value = record; testDeviceId.value = ''; testModalVisible.value = true }
const executeTest = async () => {
  if (!testDeviceId.value) { Message.warning('请选择设备'); return }
  testing.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/ai/action-library/test', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify({ device_id: testDeviceId.value, action_id: testAction.value.id }) }).then(r => r.json())
    if (res.code === 0) Message.success('测试执行成功')
    else Message.error(res.message || '测试执行失败')
    testModalVisible.value = false
  } catch (e) { Message.error('测试执行失败') } finally { testing.value = false }
}

onMounted(() => { loadData(); loadDevices() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>