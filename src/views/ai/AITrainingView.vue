<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>AI 训练任务</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-statistic title="训练任务总数" :value="stats.total" :value-style="{ color: '#1650ff' }">
          <template #prefix><icon-history size="20" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="进行中" :value="stats.running" :value-style="{ color: '#1650ff' }">
          <template #prefix><icon-sync size="20" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="已完成" :value="stats.completed" :value-style="{ color: '#00b42a' }">
          <template #prefix><icon-check-circle size="20" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="失败" :value="stats.failed" :value-style="{ color: '#f53f3f' }">
          <template #prefix><icon-close-circle size="20" /></template>
        </a-statistic>
      </a-col>
    </a-row>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索任务名称"
          style="width: 200px"
          search-button
          @search="loadTasks"
        />
        <a-select v-model="filters.task_type" placeholder="任务类型" allow-clear style="width: 160px" @change="loadTasks">
          <a-option value="fine_tune">微调训练</a-option>
          <a-option value="pre_train">预训练</a-option>
          <a-option value="rlhf">RLHF 训练</a-option>
          <a-option value="distill">知识蒸馏</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 140px" @change="loadTasks">
          <a-option value="pending">等待中</a-option>
          <a-option value="running">进行中</a-option>
          <a-option value="completed">已完成</a-option>
          <a-option value="failed">失败</a-option>
          <a-option value="cancelled">已取消</a-option>
        </a-select>
        <a-range-picker
          v-model="filters.time_range"
          show-time
          format="YYYY-MM-DD HH:mm"
          @change="loadTasks"
          style="width: 280px"
        />
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><icon-plus /></template>
          创建训练任务
        </a-button>
        <a-button @click="loadTasks">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="taskList"
        :loading="loading"
        row-key="id"
        :pagination="pagination"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #task_name="{ record }">
          <div class="task-info">
            <span class="task-name">{{ record.task_name }}</span>
            <span class="task-id">{{ record.task_id }}</span>
          </div>
        </template>
        <template #task_type="{ record }">
          <a-tag :color="getTaskTypeColor(record.task_type)">
            {{ getTaskTypeText(record.task_type) }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            <template #icon v-if="record.status === 'running'"><icon-sync :style="{ animation: 'spin 1s linear infinite' }" /></template>
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #progress="{ record }">
          <div class="progress-cell" v-if="record.status === 'running' || record.status === 'completed'">
            <a-progress
              :percent="record.progress || 0"
              :status="record.status === 'completed' ? 'success' : 'normal'"
              :show-text="true"
              size="small"
            />
          </div>
          <span v-else>--</span>
        </template>
        <template #duration="{ record }">
          {{ formatDuration(record.start_time, record.end_time, record.status) }}
        </template>
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showLogModal(record)">日志</a-button>
            <a-button
              v-if="record.status === 'running'"
              type="text"
              size="small"
              status="warning"
              @click="handleCancel(record)"
            >取消</a-button>
            <a-button
              v-if="record.status === 'completed'"
              type="text"
              size="small"
              @click="handleDeploy(record)"
            >部署</a-button>
            <a-button
              v-if="['failed', 'completed', 'cancelled'].includes(record.status)"
              type="text"
              size="small"
              status="danger"
              @click="handleDelete(record)"
            >删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 创建训练任务弹窗 -->
    <a-modal
      v-model:visible="createModalVisible"
      title="创建训练任务"
      :width="640"
      @ok="handleCreate"
      @cancel="createModalVisible = false"
      :mask-closable="false"
    >
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="任务名称" required>
          <a-input v-model="createForm.task_name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-form-item label="任务类型" required>
          <a-select v-model="createForm.task_type" placeholder="请选择任务类型">
            <a-option value="fine_tune">微调训练</a-option>
            <a-option value="pre_train">预训练</a-option>
            <a-option value="rlhf">RLHF 训练</a-option>
            <a-option value="distill">知识蒸馏</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="基础模型" required>
          <a-select v-model="createForm.base_model_id" placeholder="请选择基础模型" allow-search>
            <a-option v-for="m in availableModels" :key="m.id" :value="m.id">{{ m.model_name }} ({{ m.model_id }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="训练数据集">
          <a-upload
            :custom-request="handleDataFileRequest"
            :limit="1"
            accept=".json,.jsonl,.csv"
          >
            <template #upload-button>
              <a-button>
                <template #icon><icon-upload /></template>
                上传数据集
              </a-button>
            </template>
          </a-upload>
          <span class="form-hint">支持 .json, .jsonl, .csv 格式</span>
        </a-form-item>
        <a-form-item label="验证数据集">
          <a-upload
            :custom-request="handleValidFileRequest"
            :limit="1"
            accept=".json,.jsonl,.csv"
          >
            <template #upload-button>
              <a-button>
                <template #icon><icon-upload /></template>
                上传验证集
              </a-button>
            </template>
          </a-upload>
        </a-form-item>
        <a-form-item label="训练参数">
          <a-space direction="vertical" :size="8" style="width: 100%">
            <a-space>
              <span style="width: 80px">学习率:</span>
              <a-input-number v-model="createForm.learning_rate" :min="0" :max="1" :step="0.0001" placeholder="默认 0.001" style="width: 160px" />
            </a-space>
            <a-space>
              <span style="width: 80px">批次大小:</span>
              <a-input-number v-model="createForm.batch_size" :min="1" :max="128" placeholder="默认 32" style="width: 160px" />
            </a-space>
            <a-space>
              <span style="width: 80px">训练轮次:</span>
              <a-input-number v-model="createForm.epochs" :min="1" :max="100" placeholder="默认 3" style="width: 160px" />
            </a-space>
            <a-space>
              <span style="width: 80px">最大长度:</span>
              <a-input-number v-model="createForm.max_length" :min="64" :max="4096" placeholder="默认 512" style="width: 160px" />
            </a-space>
          </a-space>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="createForm.description" placeholder="请输入备注信息" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 训练日志弹窗 -->
    <a-drawer
      v-model:visible="logDrawerVisible"
      title="训练日志"
      :width="720"
      :footer="false"
    >
      <div class="log-container" v-if="currentTask">
        <a-descriptions :column="2" bordered size="small" style="margin-bottom: 16px">
          <a-descriptions-item label="任务名称">{{ currentTask.task_name }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentTask.status)">{{ getStatusText(currentTask.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="开始时间">{{ formatDate(currentTask.start_time) }}</a-descriptions-item>
          <a-descriptions-item label="结束时间">{{ formatDate(currentTask.end_time) }}</a-descriptions-item>
          <a-descriptions-item label="当前轮次" :span="2">{{ currentTask.current_epoch || '--' }} / {{ currentTask.total_epochs || '--' }}</a-descriptions-item>
          <a-descriptions-item label="损失值" :span="2">{{ currentTask.loss || '--' }}</a-descriptions-item>
        </a-descriptions>
        <a-tabs default-active-key="logs">
          <a-tab-pane key="logs" title="训练日志">
            <div class="log-content" ref="logContentRef">
              <pre>{{ trainingLogs }}</pre>
            </div>
          </a-tab-pane>
          <a-tab-pane key="metrics" title="指标曲线">
            <a-empty v-if="!trainingMetrics.length" />
            <div v-else class="metrics-chart">
              <div v-for="(m, i) in trainingMetrics" :key="i" class="metric-item">
                <span class="metric-name">{{ m.name }}</span>
                <a-progress :percent="m.value" :status="m.status" size="small" />
              </div>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  getAiTrainingTasks,
  postAiTrainingTask,
  deleteAiTrainingTask,
  cancelAiTrainingTask,
  getAiTrainingTaskLogs
} from '@/api/ai'

const loading = ref(false)
const taskList = ref([])
const availableModels = ref([])
const stats = reactive({
  total: 0,
  running: 0,
  completed: 0,
  failed: 0
})

const filters = reactive({
  keyword: '',
  task_type: undefined,
  status: undefined,
  time_range: []
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const createModalVisible = ref(false)
const logDrawerVisible = ref(false)
const currentTask = ref(null)

const createForm = reactive({
  task_name: '',
  task_type: '',
  base_model_id: '',
  train_data_path: '',
  valid_data_path: '',
  learning_rate: 0.001,
  batch_size: 32,
  epochs: 3,
  max_length: 512,
  description: ''
})

const trainingLogs = ref('')
const trainingMetrics = ref([])
const logPollingTimer = ref(null)

const columns = [
  { title: '任务名称', slotName: 'task_name', width: 200, ellipsis: true },
  { title: '任务类型', slotName: 'task_type', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 160 },
  { title: '耗时', slotName: 'duration', width: 100 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const getTaskTypeColor = (type) => ({
  fine_tune: 'arcoblue',
  pre_train: 'purple',
  rlhf: 'orange',
  distill: 'green'
}[type] || 'gray')

const getTaskTypeText = (type) => ({
  fine_tune: '微调训练',
  pre_train: '预训练',
  rlhf: 'RLHF 训练',
  distill: '知识蒸馏'
}[type] || type)

const getStatusColor = (status) => ({
  pending: 'gray',
  running: 'arcoblue',
  completed: 'green',
  failed: 'red',
  cancelled: 'orange'
}[status] || 'gray')

const getStatusText = (status) => ({
  pending: '等待中',
  running: '进行中',
  completed: '已完成',
  failed: '失败',
  cancelled: '已取消'
}[status] || status)

const formatDate = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const formatDuration = (start, end, status) => {
  if (!start) return '--'
  const startTime = new Date(start).getTime()
  const endTime = end ? new Date(end).getTime() : Date.now()
  const diff = Math.floor((endTime - startTime) / 1000)
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ${diff % 60}s`
  return `${Math.floor(diff / 3600)}h ${Math.floor((diff % 3600) / 60)}m`
}

const loadTasks = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.task_type) params.task_type = filters.task_type
    if (filters.status) params.status = filters.status
    if (filters.time_range && filters.time_range.length === 2) {
      params.start_time = filters.time_range[0].toISOString()
      params.end_time = filters.time_range[1].toISOString()
    }

    const res = await getAiTrainingTasks(params)
    if (res.code === 0) {
      taskList.value = res.data.list || []
      pagination.total = res.data.total || 0

      // 计算统计数据
      stats.total = taskList.value.length
      stats.running = taskList.value.filter(t => t.status === 'running').length
      stats.completed = taskList.value.filter(t => t.status === 'completed').length
      stats.failed = taskList.value.filter(t => t.status === 'failed').length
    }
  } catch (e) {
    Message.error('加载训练任务失败')
  } finally {
    loading.value = false
  }
}

const loadAvailableModels = async () => {
  try {
    const res = await getAiTrainingTasks({ page: 1, page_size: 100 })
    if (res.code === 0) {
      // 从现有模型中提取可用的基础模型
      availableModels.value = (res.data.list || []).map(t => ({
        id: t.id,
        model_name: t.task_name,
        model_id: t.task_id
      })).slice(0, 5)
    }
  } catch (e) {
    // ignore
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadTasks()
}

const handlePageSizeChange = (size) => {
  pagination.pageSize = size
  pagination.current = 1
  loadTasks()
}

const showCreateModal = () => {
  Object.assign(createForm, {
    task_name: '',
    task_type: '',
    base_model_id: '',
    train_data_path: '',
    valid_data_path: '',
    learning_rate: 0.001,
    batch_size: 32,
    epochs: 3,
    max_length: 512,
    description: ''
  })
  createModalVisible.value = true
}

const showLogModal = async (record) => {
  currentTask.value = record
  logDrawerVisible.value = true
  trainingLogs.value = '加载中...'
  trainingMetrics.value = []

  // 加载日志
  await loadTaskLogs(record.id)

  // 如果任务还在运行，开始轮询
  if (record.status === 'running') {
    startLogPolling(record.id)
  }
}

const loadTaskLogs = async (taskId) => {
  try {
    const res = await getAiTrainingTaskLogs(taskId)
    if (res.code === 0) {
      trainingLogs.value = res.data.logs || '暂无日志'
      trainingMetrics.value = res.data.metrics || []

      // 更新任务状态
      if (currentTask.value && currentTask.value.id === taskId) {
        Object.assign(currentTask.value, res.data)
      }
    }
  } catch (e) {
    trainingLogs.value = '加载日志失败'
  }
}

const startLogPolling = (taskId) => {
  stopLogPolling()
  logPollingTimer.value = setInterval(() => {
    loadTaskLogs(taskId)
    // 如果任务完成或失败，停止轮询
    if (currentTask.value && !['running', 'pending'].includes(currentTask.value.status)) {
      stopLogPolling()
      loadTasks()
    }
  }, 5000)
}

const stopLogPolling = () => {
  if (logPollingTimer.value) {
    clearInterval(logPollingTimer.value)
    logPollingTimer.value = null
  }
}

const handleDataFileRequest = (options) => {
  const { fileItem, onSuccess } = options
  setTimeout(() => {
    createForm.train_data_path = fileItem.file.name
    onSuccess(fileItem)
    Message.success('数据集上传成功')
  }, 500)
}

const handleValidFileRequest = (options) => {
  const { fileItem, onSuccess } = options
  setTimeout(() => {
    createForm.valid_data_path = fileItem.file.name
    onSuccess(fileItem)
    Message.success('验证集上传成功')
  }, 500)
}

const handleCreate = async () => {
  if (!createForm.task_name || !createForm.task_type || !createForm.base_model_id) {
    Message.warning('请填写必填项')
    return
  }
  try {
    const res = await postAiTrainingTask(createForm)
    if (res.code === 0) {
      Message.success('训练任务创建成功')
      createModalVisible.value = false
      loadTasks()
    } else {
      Message.error(res.message || '创建失败')
    }
  } catch (e) {
    Message.error('创建失败')
  }
}

const handleCancel = (record) => {
  Modal.warning({
    title: '确认取消',
    content: `确定要取消训练任务 "${record.task_name}" 吗？`,
    okText: '确认',
    onOk: async () => {
      try {
        const res = await cancelAiTrainingTask(record.id)
        if (res.code === 0) {
          Message.success('任务已取消')
          loadTasks()
        } else {
          Message.error(res.message || '取消失败')
        }
      } catch (e) {
        Message.error('取消失败')
      }
    }
  })
}

const handleDeploy = (record) => {
  Message.info('部署功能开发中')
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除训练任务 "${record.task_name}" 吗？`,
    okText: '确认删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        const res = await deleteAiTrainingTask(record.id)
        if (res.code === 0) {
          Message.success('删除成功')
          loadTasks()
        } else {
          Message.error(res.message || '删除失败')
        }
      } catch (e) {
        Message.error('删除失败')
      }
    }
  })
}

onMounted(() => {
  loadTasks()
  loadAvailableModels()
})

onUnmounted(() => {
  stopLogPolling()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stats-row .arco-statistic { background: #fff; padding: 16px; border-radius: 8px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; display: flex; justify-content: flex-start; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
.task-info { display: flex; flex-direction: column; }
.task-name { font-weight: 500; }
.task-id { font-size: 12px; color: #999; }
.progress-cell { width: 140px; }
.form-hint { font-size: 12px; color: #999; margin-top: 4px; display: block; }
.log-container { height: 100%; }
.log-content {
  background: #1e1e1e; border-radius: 8px; padding: 16px;
  max-height: 500px; overflow-y: auto; font-family: 'Courier New', monospace;
  font-size: 12px; color: #d4d4d4;
}
.log-content pre { margin: 0; white-space: pre-wrap; word-break: break-all; }
.metrics-chart { display: flex; flex-direction: column; gap: 12px; }
.metric-item { display: flex; align-items: center; gap: 12px; }
.metric-name { width: 100px; font-size: 14px; color: #666; }

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
