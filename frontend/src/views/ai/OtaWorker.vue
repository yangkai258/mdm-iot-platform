<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>OTA升级</a-breadcrumb-item>
      <a-breadcrumb-item>Worker监控</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- Worker状态卡片 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card>
          <a-statistic title="Worker总数" :value="workerStats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="运行中" :value="workerStats.running" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="空闲" :value="workerStats.idle" :value-style="{ color: '#1890ff' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="错误" :value="workerStats.error" :value-style="{ color: '#ff4d4f' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- Worker列表 -->
    <a-card title="Worker节点" style="margin-bottom: 16px">
      <a-table :columns="workerColumns" :data="workers" :loading="loading" :pagination="false" row-key="id" :scroll="{ x: 900 }">
        <template #status="{ record }">
          <a-badge :status="getWorkerStatus(record.status)" :text="getWorkerStatusText(record.status)" />
        </template>
        <template #cpu="{ record }">
          <a-progress :percent="record.cpu_usage" :color="record.cpu_usage > 80 ? '#ff4d4f' : '#52c41a'" size="small" />
        </template>
        <template #memory="{ record }">
          <a-progress :percent="record.memory_usage" :color="record.memory_usage > 80 ? '#ff4d4f' : '#1890ff'" size="small" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" status="danger" @click="pauseWorker(record)" v-if="record.status === 'running'">暂停</a-button>
            <a-button type="text" size="small" @click="resumeWorker(record)" v-if="record.status === 'paused'">恢复</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 任务队列 -->
    <a-card title="任务队列">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showTriggerModal">手动触发任务</a-button>
          <a-button @click="loadTasks">刷新</a-button>
        </a-space>
      </template>
      <a-table :columns="taskColumns" :data="tasks" :loading="taskLoading" :pagination="pagination" row-key="id" @page-change="handlePageChange" :scroll="{ x: 1100 }">
        <template #status="{ record }">
          <a-tag :color="getTaskStatusColor(record.status)">{{ getTaskStatusText(record.status) }}</a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress :percent="record.progress" size="small" :status="record.status === 'failed' ? 'exception' : undefined" />
        </template>
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="retryTask(record)" v-if="record.status === 'failed'">重试</a-button>
            <a-button type="text" size="small" status="danger" @click="cancelTask(record)" v-if="record.status === 'pending' || record.status === 'running'">取消</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 手动触发任务对话框 -->
    <a-modal v-model:visible="triggerVisible" title="手动触发OTA任务" @ok="submitTrigger" :width="520" :loading="submitting">
      <a-form :model="triggerForm" layout="vertical">
        <a-form-item label="设备" required>
          <a-select v-model="triggerForm.device_id" placeholder="选择设备" show-search>
            <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }} ({{ d.id }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="固件版本" required>
          <a-select v-model="triggerForm.firmware_id" placeholder="选择固件">
            <a-option v-for="f in firmwares" :key="f.id" :value="f.id">{{ f.version }} ({{ f.device_model }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级">
          <a-select v-model="triggerForm.priority" placeholder="选择优先级">
            <a-option value="low">低</a-option>
            <a-option value="normal">普通</a-option>
            <a-option value="high">高</a-option>
            <a-option value="urgent">紧急</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const loading = ref(false)
const taskLoading = ref(false)
const triggerVisible = ref(false)
const submitting = ref(false)

const workers = ref<any[]>([])
const tasks = ref<any[]>([])
const devices = ref<any[]>([])
const firmwares = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const workerStats = reactive({ total: 0, running: 0, idle: 0, error: 0 })

const workerColumns = [
  { title: 'Worker ID', dataIndex: 'id', width: 120 },
  { title: '节点名称', dataIndex: 'name', width: 150 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '当前任务', dataIndex: 'current_task_id', width: 120 },
  { title: 'CPU使用率', slotName: 'cpu', width: 120 },
  { title: '内存使用率', slotName: 'memory', width: 120 },
  { title: '最后心跳', dataIndex: 'last_heartbeat', width: 180 },
  { title: '操作', slotName: 'actions', fixed: 'right', width: 120 },
]

const taskColumns = [
  { title: '任务ID', dataIndex: 'id', width: 80 },
  { title: '设备ID', dataIndex: 'device_id', width: 140, ellipsis: true },
  { title: '固件版本', dataIndex: 'firmware_version', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 140 },
  { title: '错误信息', dataIndex: 'error_msg', ellipsis: true },
  { title: '创建时间', dataIndex: 'created_at', slotName: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', fixed: 'right', width: 120 },
]

const triggerForm = reactive({ device_id: '', firmware_id: '', priority: 'normal' })

const loadWorkers = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/ota/workers')
    workers.value = res.data.items || []
    workerStats.total = workers.value.length
    workerStats.running = workers.value.filter((w: any) => w.status === 'running').length
    workerStats.idle = workers.value.filter((w: any) => w.status === 'idle').length
    workerStats.error = workers.value.filter((w: any) => w.status === 'error').length
  } catch {
    workers.value = [
      { id: 'worker-1', name: 'OTA-Worker-01', status: 'running', current_task_id: 'task-101', cpu_usage: 45, memory_usage: 62, last_heartbeat: new Date().toISOString() },
      { id: 'worker-2', name: 'OTA-Worker-02', status: 'idle', current_task_id: '-', cpu_usage: 5, memory_usage: 30, last_heartbeat: new Date().toISOString() },
      { id: 'worker-3', name: 'OTA-Worker-03', status: 'error', current_task_id: '-', cpu_usage: 0, memory_usage: 0, last_heartbeat: new Date(Date.now() - 300000).toISOString() },
    ]
    workerStats.total = 3; workerStats.running = 1; workerStats.idle = 1; workerStats.error = 1
  } finally { loading.value = false }
}

const loadTasks = async () => {
  taskLoading.value = true
  try {
    const res = await axios.get('/api/v1/ota/worker-tasks', { params: { page: pagination.current, page_size: pagination.pageSize } })
    tasks.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch {
    tasks.value = [
      { id: 1, device_id: 'dev-001', firmware_version: 'v2.1.0', status: 'running', progress: 65, error_msg: '', created_at: new Date().toISOString() },
      { id: 2, device_id: 'dev-002', firmware_version: 'v2.1.0', status: 'pending', progress: 0, error_msg: '', created_at: new Date(Date.now() - 60000).toISOString() },
      { id: 3, device_id: 'dev-003', firmware_version: 'v2.0.5', status: 'failed', progress: 30, error_msg: '下载超时', created_at: new Date(Date.now() - 120000).toISOString() },
    ]
    pagination.total = 3
  } finally { taskLoading.value = false }
}

const loadDevicesAndFirmwares = async () => {
  try {
    const [devRes, fwRes] = await Promise.all([axios.get('/api/v1/devices?page_size=100'), axios.get('/api/v1/ota/firmwares?page_size=100')])
    devices.value = devRes.data.items || []
    firmwares.value = fwRes.data.items || []
  } catch {
    devices.value = [{ id: 'dev-001', name: '设备1号' }, { id: 'dev-002', name: '设备2号' }]
    firmwares.value = [{ id: 1, version: 'v2.1.0', device_model: 'M5Stack-Core2' }, { id: 2, version: 'v2.0.5', device_model: 'M5Stack-Core' }]
  }
}

const showTriggerModal = () => { triggerVisible.value = true }
const submitTrigger = async () => {
  submitting.value = true
  try {
    await axios.post('/api/v1/ota/worker-tasks', triggerForm)
    Message.success('任务已创建')
    triggerVisible.value = false
    loadTasks()
  } catch (e) { Message.error('创建失败') } finally { submitting.value = false }
}
const pauseWorker = async (record: any) => {
  try { await axios.post(`/api/v1/ota/workers/${record.id}/pause`); Message.success('已暂停'); loadWorkers() } catch { Message.error('操作失败') }
}
const resumeWorker = async (record: any) => {
  try { await axios.post(`/api/v1/ota/workers/${record.id}/resume`); Message.success('已恢复'); loadWorkers() } catch { Message.error('操作失败') }
}
const retryTask = async (record: any) => {
  try { await axios.post(`/api/v1/ota/worker-tasks/${record.id}/retry`); Message.success('已重试'); loadTasks() } catch { Message.error('重试失败') }
}
const cancelTask = async (record: any) => {
  try { await axios.delete(`/api/v1/ota/worker-tasks/${record.id}`); Message.success('已取消'); loadTasks() } catch { Message.error('取消失败') }
}
const handlePageChange = (page: number) => { pagination.current = page; loadTasks() }

const getWorkerStatus = (s: string) => ({ running: 'success', idle: 'normal', error: 'error', paused: 'warning' }[s] || 'default')
const getWorkerStatusText = (s: string) => ({ running: '运行中', idle: '空闲', error: '错误', paused: '已暂停' }[s] || s)
const getTaskStatusColor = (s: string) => ({ pending: 'arcoblue', running: 'green', completed: 'gray', failed: 'red' }[s] || 'gray')
const getTaskStatusText = (s: string) => ({ pending: '等待中', running: '进行中', completed: '已完成', failed: '失败' }[s] || s)
const formatDate = (d: string) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(() => { loadWorkers(); loadTasks(); loadDevicesAndFirmwares() })
</script>