<template>
  <div class="page-container">
    <a-card class="general-card" title="AI训练流水�?>
      <template #extra>
        <a-button type="primary" @click="openCreateModal"><icon-plus />新建训练任务</a-button>
      </template>
      <div class="search-form">
        <a-form :model="form" layout="inline">
          <a-form-item label="任务名称"><a-input v-model="form.task_name" placeholder="请输�? /></a-form-item>
          <a-form-item label="状�?>
            <a-select v-model="form.status" placeholder="选择状�? allow-clear style="width: 140px">
              <a-option value="pending">排队�?/a-option>
              <a-option value="running">训练�?/a-option>
              <a-option value="completed">已完�?/a-option>
              <a-option value="failed">失败</a-option>
            </a-select>
          </a-form-item>
          <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="Object.keys(form).forEach(k => form[k] = ''); loadData()">重置</a-button></a-form-item>
        </a-form>
      </div>
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress :percent="record.progress || 0" :status="record.status === 'failed' ? 'exception' : record.status === 'completed' ? 'success' : 'normal'" size="small" />
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewLogs(record)">日志</a-button>
          <a-button v-if="record.status === 'failed'" type="text" size="small" @click="retryTask(record)">重试</a-button>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:visible="createModalVisible" title="新建训练任务" @before-ok="handleCreate" :loading="submitting" :width="600">
      <a-form :model="trainForm" layout="vertical">
        <a-form-item label="任务名称" required><a-input v-model="trainForm.task_name" placeholder="请输入任务名�? /></a-form-item>
        <a-form-item label="选择数据�? required>
          <a-select v-model="trainForm.dataset_id" placeholder="选择数据�?>
            <a-option v-for="d in datasets" :key="d.id" :value="d.id">{{ d.name }} ({{ d.sample_count }}样本)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="选择模型" required>
          <a-select v-model="trainForm.model_id" placeholder="选择基础模型">
            <a-option v-for="m in models" :key="m.id" :value="m.id">{{ m.name }} ({{ m.version }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="训练轮次">
          <a-input-number v-model="trainForm.epochs" :min="1" :max="1000" placeholder="默认10" style="width: 200px" />
        </a-form-item>
        <a-form-item label="学习�?>
          <a-input-number v-model="trainForm.learning_rate" :min="0" :max="1" :step="0.001" placeholder="默认0.001" style="width: 200px" />
        </a-form-item>
        <a-form-item label="备注"><a-textarea v-model="trainForm.description" :rows="2" placeholder="训练任务描述" /></a-form-item>
      </a-form>
    </a-modal>
    <a-modal v-model:visible="logModalVisible" title="训练日志" :width="800" :footer="null">
      <a-spin :loading="logLoading">
        <div class="log-container">
          <pre>{{ logs }}</pre>
        </div>
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const submitting = ref(false)
const logLoading = ref(false)
const data = ref([])
const datasets = ref([])
const models = ref([])
const logs = ref('')
const createModalVisible = ref(false)
const logModalVisible = ref(false)
const selectedTask = ref(null)
const form = reactive({ task_name: '', status: '' })
const trainForm = reactive({ task_name: '', dataset_id: '', model_id: '', epochs: 10, learning_rate: 0.001, description: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '任务ID', dataIndex: 'id', width: 80 },
  { title: '任务名称', dataIndex: 'task_name', width: 180 },
  { title: '数据�?, dataIndex: 'dataset_name', width: 140 },
  { title: '基础模型', dataIndex: 'model_name', width: 120 },
  { title: '状�?, slotName: 'status', width: 90 },
  { title: '进度', slotName: 'progress', width: 160 },
  { title: '开始时�?, dataIndex: 'started_at', width: 170 },
  { title: '耗时', dataIndex: 'duration', width: 100 },
  { title: '操作', slotName: 'actions', width: 100 }
]

const getStatusColor = (s) => ({ pending: 'orange', running: 'arcoblue', completed: 'green', failed: 'red' }[s] || 'gray')
const getStatusText = (s) => ({ pending: '排队�?, running: '训练�?, completed: '已完�?, failed: '失败' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.task_name) params.append('task_name', form.task_name)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/ai/training/tasks?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const loadDatasets = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/ai/datasets?page_size=100', { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    datasets.value = res.data?.list || []
  } catch (e) { console.error('加载数据集失�?, e) }
}

const loadModels = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/ai/models?page_size=100', { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    models.value = res.data?.list || []
  } catch (e) { console.error('加载模型失败', e) }
}

const openCreateModal = () => { Object.assign(trainForm, { task_name: '', dataset_id: '', model_id: '', epochs: 10, learning_rate: 0.001, description: '' }); createModalVisible.value = true }

const handleCreate = async (done) => {
  if (!trainForm.task_name || !trainForm.dataset_id || !trainForm.model_id) { Message.warning('请填写必填项'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/ai/training/tasks', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(trainForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('训练任务已创�?); createModalVisible.value = false; loadData() }
    else { Message.error(res.message || '创建失败') }
    done(true)
  } catch (e) { Message.error('创建失败'); done(false) } finally { submitting.value = false }
}

const viewLogs = async (record) => {
  selectedTask.value = record
  logModalVisible.value = true
  logLoading.value = true
  logs.value = ''
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/ai/training/tasks/${record.id}/logs`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    logs.value = res.data?.logs || '暂无日志'
  } catch (e) { logs.value = '加载日志失败' } finally { logLoading.value = false }
}

const retryTask = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/ai/training/tasks/${record.id}/retry`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('任务已重新提�?); loadData() }
    else Message.error('重试失败')
  } catch (e) { Message.error('重试失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }

onMounted(() => { loadData(); loadDatasets(); loadModels() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
.log-container { max-height: 400px; overflow-y: auto; background: #1e1e1e; color: #d4d4d4; padding: 12px; border-radius: 4px; }
.log-container pre { margin: 0; white-space: pre-wrap; font-family: 'Courier New', monospace; font-size: 12px; }
</style>