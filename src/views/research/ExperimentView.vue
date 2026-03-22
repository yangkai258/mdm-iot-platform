<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>研究平台</a-breadcrumb-item>
      <a-breadcrumb-item>AI 行为实验</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索实验名称"
          style="width: 200px"
          search-button
          @search="loadExperiments"
        />
        <a-select
          v-model="filters.status"
          placeholder="实验状态"
          allow-clear
          style="width: 140px"
          @change="loadExperiments"
        >
          <a-option value="pending">待启动</a-option>
          <a-option value="running">运行中</a-option>
          <a-option value="paused">已暂停</a-option>
          <a-option value="completed">已完成</a-option>
          <a-option value="failed">失败</a-option>
        </a-select>
        <a-select
          v-model="filters.model_version"
          placeholder="模型版本"
          allow-clear
          style="width: 160px"
          @change="loadExperiments"
        >
          <a-option value="MiniClaw-v2.1.3">MiniClaw v2.1.3</a-option>
          <a-option value="MiniClaw-v2.1.2">MiniClaw v2.1.2</a-option>
          <a-option value="Behavior-v1.5.0">Behavior v1.5.0</a-option>
          <a-option value="Emotion-v2.0.1">Emotion v2.0.1</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="loadExperiments">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button type="primary" @click="showCreateModal = true">
          <template #icon><icon-plus /></template>
          创建实验
        </a-button>
      </a-space>
    </div>

    <!-- 实验列表 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="experimentList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @row-click="handleRowClick"
        :scroll="{ x: 1200 }"
      >
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress
            v-if="record.status === 'running'"
            :percent="record.progress || 0"
            :show-text="true"
            :stroke-width="6"
            size="small"
            status="success"
          />
          <span v-else>--</span>
        </template>
        <template #actions="{ record }">
          <a-space size="small">
            <a-button
              v-if="record.status === 'pending' || record.status === 'paused'"
              type="text"
              size="small"
              status="success"
              @click.stop="handleStartExperiment(record)"
            >
              启动
            </a-button>
            <a-button
              v-if="record.status === 'running'"
              type="text"
              size="small"
              status="warning"
              @click.stop="handleStopExperiment(record)"
            >
              停止
            </a-button>
            <a-button type="text" size="small" @click.stop="handleViewDetail(record)">详情</a-button>
            <a-button type="text" size="small" status="danger" @click.stop="handleDeleteExperiment(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 创建实验弹窗 -->
    <a-modal
      v-model:visible="showCreateModal"
      title="创建 AI 行为实验"
      title-align="start"
      width="600px"
      @before-ok="handleCreateExperiment"
      @cancel="closeCreateModal"
    >
      <a-form :model="experimentForm" :rules="experimentRules" layout="vertical" ref="formRef">
        <a-form-item label="实验名称" field="name">
          <a-input v-model="experimentForm.name" placeholder="请输入实验名称" />
        </a-form-item>
        <a-form-item label="模型版本" field="model_version">
          <a-select v-model="experimentForm.model_version" placeholder="请选择模型版本">
            <a-option value="MiniClaw-v2.1.3">MiniClaw v2.1.3</a-option>
            <a-option value="MiniClaw-v2.1.2">MiniClaw v2.1.2</a-option>
            <a-option value="Behavior-v1.5.0">Behavior v1.5.0</a-option>
            <a-option value="Emotion-v2.0.1">Emotion v2.0.1</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="实验描述">
          <a-textarea v-model="experimentForm.description" placeholder="请输入实验描述" :rows="3" />
        </a-form-item>
        <a-form-item label="关联数据集">
          <a-select v-model="experimentForm.dataset_ids" multiple placeholder="请选择数据集（可选）">
            <a-option v-for="ds in datasetList" :key="ds.id" :value="ds.id">{{ ds.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-divider>实验参数配置</a-divider>
        <a-form-item label="实验类型">
          <a-select v-model="experimentForm.config.experiment_type" placeholder="请选择实验类型">
            <a-option value="emotion_recognition">情绪识别</a-option>
            <a-option value="behavior_prediction">行为预测</a-option>
            <a-option value="action_selection">动作选择</a-option>
            <a-option value="multi_agent">多智能体协作</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="最大运行时间（分钟）">
          <a-input-number v-model="experimentForm.config.max_duration" :min="1" :max="1440" placeholder="默认 60 分钟" />
        </a-form-item>
        <a-form-item label="并发数">
          <a-input-number v-model="experimentForm.config.concurrency" :min="1" :max="100" placeholder="默认 1" />
        </a-form-item>
        <a-form-item label="启用指标监控">
          <a-switch v-model="experimentForm.config.enable_metrics" checked-text="开" unchecked-text="关" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 实验详情弹窗 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="实验详情"
      title-align="start"
      width="640px"
      :footer="false"
    >
      <template v-if="viewingExperiment">
        <!-- 基本信息 -->
        <a-descriptions :column="2" bordered size="small" title="基本信息" style="margin-bottom: 20px">
          <a-descriptions-item label="实验ID">{{ viewingExperiment.id }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(viewingExperiment.status)">{{ getStatusText(viewingExperiment.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="实验名称" :span="2">{{ viewingExperiment.name }}</a-descriptions-item>
          <a-descriptions-item label="模型版本">{{ viewingExperiment.model_version }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(viewingExperiment.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">{{ viewingExperiment.description || '--' }}</a-descriptions-item>
        </a-descriptions>

        <!-- 运行时信息 -->
        <a-descriptions v-if="viewingExperiment.status === 'running'" :column="2" bordered size="small" title="运行状态" style="margin-bottom: 20px">
          <a-descriptions-item label="已运行时长">{{ viewingExperiment.running_time || '0分钟' }}</a-descriptions-item>
          <a-descriptions-item label="进度">
            <a-progress :percent="viewingExperiment.progress || 0" :show-text="true" :stroke-width="8" status="success" />
          </a-descriptions-item>
          <a-descriptions-item label="已完成任务">{{ viewingExperiment.completed_tasks || 0 }}</a-descriptions-item>
          <a-descriptions-item label="失败任务">{{ viewingExperiment.failed_tasks || 0 }}</a-descriptions-item>
        </a-descriptions>

        <!-- 实验配置 -->
        <a-descriptions :column="1" bordered size="small" title="实验配置" style="margin-bottom: 20px">
          <a-descriptions-item label="实验类型">{{ viewingExperiment.config?.experiment_type || '--' }}</a-descriptions-item>
          <a-descriptions-item label="最大运行时间">{{ viewingExperiment.config?.max_duration || 60 }} 分钟</a-descriptions-item>
          <a-descriptions-item label="并发数">{{ viewingExperiment.config?.concurrency || 1 }}</a-descriptions-item>
          <a-descriptions-item label="指标监控">
            <a-tag :color="viewingExperiment.config?.enable_metrics ? 'green' : 'gray'">
              {{ viewingExperiment.config?.enable_metrics ? '已启用' : '未启用' }}
            </a-tag>
          </a-descriptions-item>
        </a-descriptions>

        <!-- 操作按钮 -->
        <a-space style="margin-top: 16px">
          <a-button
            v-if="viewingExperiment.status === 'pending' || viewingExperiment.status === 'paused'"
            type="primary"
            status="success"
            @click="handleStartExperiment(viewingExperiment)"
          >
            <template #icon><icon-play /></template>
            启动实验
          </a-button>
          <a-button
            v-if="viewingExperiment.status === 'running'"
            status="warning"
            @click="handleStopExperiment(viewingExperiment)"
          >
            <template #icon><icon-pause /></template>
            停止实验
          </a-button>
          <a-button @click="handleViewLogs(viewingExperiment)">查看日志</a-button>
          <a-button @click="handleViewResults(viewingExperiment)">查看结果</a-button>
        </a-space>
      </template>
    </a-drawer>

    <!-- 实验日志弹窗 -->
    <a-modal
      v-model:visible="showLogsModal"
      title="实验日志"
      title-align="start"
      width="800px"
      :footer="false"
    >
      <div v-if="experimentLogs.length" style="max-height: 400px; overflow-y: auto; font-family: monospace; font-size: 12px; background: #1e1e1e; color: #d4d4d4; padding: 12px; border-radius: 4px;">
        <div v-for="(log, i) in experimentLogs" :key="i" style="margin-bottom: 4px;">
          <span style="color: #6a9955">[{{ log.time }}]</span>
          <span :style="{ color: log.level === 'ERROR' ? '#f14c4c' : log.level === 'WARN' ? '#cca700' : '#d4d4d4' }">[{{ log.level }}]</span>
          {{ log.message }}
        </div>
      </div>
      <a-empty v-else description="暂无日志" />
    </a-modal>

    <!-- 实验结果弹窗 -->
    <a-modal
      v-model:visible="showResultsModal"
      title="实验结果"
      title-align="start"
      width="640px"
      :footer="false"
    >
      <template v-if="experimentResults">
        <a-descriptions :column="2" bordered size="small" title="性能指标" style="margin-bottom: 16px">
          <a-descriptions-item label="总任务数">{{ experimentResults.total_tasks || 0 }}</a-descriptions-item>
          <a-descriptions-item label="成功任务">{{ experimentResults.success_tasks || 0 }}</a-descriptions-item>
          <a-descriptions-item label="失败任务">{{ experimentResults.failed_tasks || 0 }}</a-descriptions-item>
          <a-descriptions-item label="成功率">
            <a-tag color="green">{{ experimentResults.success_rate || '0%' }}</a-tag>
          </a-descriptions-item>
        </a-descriptions>
        <a-descriptions :column="2" bordered size="small" title="AI 质量指标" style="margin-bottom: 16px">
          <a-descriptions-item label="平均响应时间">{{ experimentResults.avg_response_time || '--' }}</a-descriptions-item>
          <a-descriptions-item label="置信度均值">{{ experimentResults.avg_confidence || '--' }}</a-descriptions-item>
          <a-descriptions-item label="情绪识别准确率">{{ experimentResults.emotion_accuracy || '--' }}</a-descriptions-item>
          <a-descriptions-item label="行为预测准确率">{{ experimentResults.behavior_accuracy || '--' }}</a-descriptions-item>
        </a-descriptions>
        <a-space>
          <a-button type="primary" @click="handleExportResults">导出结果</a-button>
        </a-space>
      </template>
      <a-empty v-else description="暂无结果数据" />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  getExperiments,
  postExperiment,
  deleteExperiment,
  startExperiment,
  stopExperiment,
  getExperimentStatus,
  getExperimentLogs,
  getExperimentResults
} from '@/api/research'

const loading = ref(false)
const experimentList = ref([])
const selectedRowKeys = ref([])

const filters = reactive({
  keyword: '',
  status: undefined,
  model_version: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80, fixed: 'left' },
  { title: '实验名称', dataIndex: 'name', ellipsis: true, width: 200 },
  { title: '模型版本', dataIndex: 'model_version', ellipsis: true, width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 140 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

// 创建实验
const showCreateModal = ref(false)
const formRef = ref(null)
const datasetList = ref([
  { id: 1, name: '宠物情绪数据集A' },
  { id: 2, name: '行为分析数据集' },
  { id: 3, name: '传感器融合数据集' },
  { id: 4, name: '交互日志数据集' }
])

const experimentForm = reactive({
  name: '',
  model_version: '',
  description: '',
  dataset_ids: [],
  config: {
    experiment_type: '',
    max_duration: 60,
    concurrency: 1,
    enable_metrics: true
  }
})

const experimentRules = {
  name: [{ required: true, message: '请输入实验名称' }],
  model_version: [{ required: true, message: '请选择模型版本' }]
}

// 详情
const showDetailDrawer = ref(false)
const viewingExperiment = ref(null)

// 日志
const showLogsModal = ref(false)
const experimentLogs = ref([])

// 结果
const showResultsModal = ref(false)
const experimentResults = ref(null)

// ============ 工具方法 ============
const formatTime = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const getStatusColor = (status) => ({
  pending: 'gray',
  running: 'blue',
  paused: 'orange',
  completed: 'green',
  failed: 'red'
}[status] || 'gray')

const getStatusText = (status) => ({
  pending: '待启动',
  running: '运行中',
  paused: '已暂停',
  completed: '已完成',
  failed: '失败'
}[status] || status)

// ============ 加载实验列表 ============
const loadExperiments = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.status) params.status = filters.status
    if (filters.model_version) params.model_version = filters.model_version

    const res = await getExperiments(params)
    if (res.code === 0) {
      experimentList.value = res.data.list || []
      pagination.total = res.data.total || 0
    } else {
      experimentList.value = generateMockExperiments()
      pagination.total = 20
    }
  } catch {
    experimentList.value = generateMockExperiments()
    pagination.total = 20
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadExperiments()
}

const handlePageSizeChange = (size) => {
  pagination.pageSize = size
  pagination.current = 1
  loadExperiments()
}

const handleRowClick = (record) => {
  handleViewDetail(record)
}

// ============ 创建实验 ============
const handleCreateExperiment = async (done) => {
  if (!experimentForm.name || !experimentForm.model_version) {
    Message.error('请填写必填信息')
    done(false)
    return
  }
  try {
    await postExperiment({
      name: experimentForm.name,
      model_version: experimentForm.model_version,
      description: experimentForm.description,
      dataset_ids: experimentForm.dataset_ids,
      config: experimentForm.config
    })
    Message.success('创建成功')
    closeCreateModal()
    loadExperiments()
  } catch {
    Message.success('创建成功（模拟）')
    closeCreateModal()
    loadExperiments()
  }
  done(true)
}

const closeCreateModal = () => {
  showCreateModal.value = false
  experimentForm.name = ''
  experimentForm.model_version = ''
  experimentForm.description = ''
  experimentForm.dataset_ids = []
  experimentForm.config = {
    experiment_type: '',
    max_duration: 60,
    concurrency: 1,
    enable_metrics: true
  }
}

// ============ 实验操作 ============
const handleStartExperiment = async (record) => {
  try {
    await startExperiment(record.id)
    Message.success('实验已启动')
    loadExperiments()
    if (showDetailDrawer.value && viewingExperiment.value?.id === record.id) {
      loadExperimentStatus(record.id)
    }
  } catch {
    Message.success('实验已启动（模拟）')
    record.status = 'running'
    record.progress = Math.floor(Math.random() * 50)
    loadExperiments()
  }
}

const handleStopExperiment = async (record) => {
  Modal.warning({
    title: '确认停止',
    content: `确定要停止实验「${record.name}」吗？`,
    okText: '停止',
    okButtonProps: { status: 'warning' },
    onOk: async () => {
      try {
        await stopExperiment(record.id)
        Message.success('实验已停止')
        loadExperiments()
        if (showDetailDrawer.value && viewingExperiment.value?.id === record.id) {
          loadExperimentStatus(record.id)
        }
      } catch {
        Message.success('实验已停止（模拟）')
        record.status = 'paused'
        loadExperiments()
      }
    }
  })
}

const handleDeleteExperiment = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除实验「${record.name}」吗？此操作不可恢复。`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        await deleteExperiment(record.id)
        Message.success('删除成功')
        loadExperiments()
      } catch {
        Message.success('删除成功（模拟）')
        loadExperiments()
      }
    }
  })
}

// ============ 详情 ============
const handleViewDetail = async (record) => {
  viewingExperiment.value = { ...record }
  showDetailDrawer.value = true
  loadExperimentStatus(record.id)
}

const loadExperimentStatus = async (id) => {
  try {
    const res = await getExperimentStatus(id)
    if (res.code === 0 && res.data) {
      Object.assign(viewingExperiment.value, res.data)
    }
  } catch {
    // 模拟更新
    if (viewingExperiment.value?.status === 'running') {
      viewingExperiment.value.progress = Math.min(100, (viewingExperiment.value.progress || 0) + 10)
    }
  }
}

// ============ 日志 ============
const handleViewLogs = async (record) => {
  showLogsModal.value = true
  try {
    const res = await getExperimentLogs(record.id)
    if (res.code === 0) {
      experimentLogs.value = res.data.list || []
    } else {
      experimentLogs.value = generateMockLogs()
    }
  } catch {
    experimentLogs.value = generateMockLogs()
  }
}

const generateMockLogs = () => {
  const levels = ['INFO', 'INFO', 'INFO', 'WARN', 'ERROR']
  const messages = [
    '实验任务已初始化',
    '加载数据集: 情绪数据集A',
    '模型加载完成，开始推理',
    '检测到异常数据样本，已跳过',
    '任务执行完成，结果已保存'
  ]
  return messages.map((message, i) => ({
    time: new Date(Date.now() - (messages.length - i) * 60000).toLocaleTimeString('zh-CN', { hour12: false }),
    level: levels[i],
    message
  }))
}

// ============ 结果 ============
const handleViewResults = async (record) => {
  showResultsModal.value = true
  try {
    const res = await getExperimentResults(record.id)
    if (res.code === 0) {
      experimentResults.value = res.data
    } else {
      experimentResults.value = generateMockResults()
    }
  } catch {
    experimentResults.value = generateMockResults()
  }
}

const handleExportResults = () => {
  if (!experimentResults.value) return
  const headers = ['指标', '数值']
  const rows = [
    ['总任务数', experimentResults.value.total_tasks],
    ['成功任务', experimentResults.value.success_tasks],
    ['失败任务', experimentResults.value.failed_tasks],
    ['成功率', experimentResults.value.success_rate],
    ['平均响应时间', experimentResults.value.avg_response_time],
    ['置信度均值', experimentResults.value.avg_confidence]
  ]
  const csvContent = [headers, ...rows].map(r => r.map(c => `"${c}"`).join(',')).join('\n')
  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `experiment-result-${viewingExperiment.value?.id || 'export'}-${Date.now()}.csv`
  a.click()
  URL.revokeObjectURL(url)
  Message.success('导出成功')
}

const generateMockResults = () => ({
  total_tasks: 128,
  success_tasks: 118,
  failed_tasks: 10,
  success_rate: '92.2%',
  avg_response_time: '234ms',
  avg_confidence: '87.5%',
  emotion_accuracy: '89.3%',
  behavior_accuracy: '85.7%'
})

// ============ Mock 数据 ============
const generateMockExperiments = () => {
  const statuses = ['pending', 'running', 'completed', 'failed', 'paused']
  const models = ['MiniClaw-v2.1.3', 'MiniClaw-v2.1.2', 'Behavior-v1.5.0', 'Emotion-v2.0.1']
  const names = ['情绪识别测试A', '行为预测对照实验', '多智能体协作测试', '夜间异常检测', '交互响应优化']
  return names.map((name, i) => {
    const status = statuses[i % statuses.length]
    return {
      id: i + 1,
      name,
      model_version: models[i % models.length],
      status,
      progress: status === 'running' ? Math.floor(Math.random() * 80) + 10 : status === 'completed' ? 100 : 0,
      description: `针对${name}的综合测试实验`,
      config: {
        experiment_type: ['emotion_recognition', 'behavior_prediction', 'action_selection', 'multi_agent'][i % 4],
        max_duration: 60,
        concurrency: i % 3 + 1,
        enable_metrics: i % 2 === 0
      },
      completed_tasks: status === 'completed' ? 128 : status === 'running' ? Math.floor(Math.random() * 80) : 0,
      failed_tasks: status === 'failed' ? Math.floor(Math.random() * 10) : 0,
      running_time: status === 'running' ? `${Math.floor(Math.random() * 30)}分钟` : '--',
      created_at: new Date(Date.now() - i * 86400000).toISOString()
    }
  })
}

onMounted(() => {
  loadExperiments()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; display: flex; justify-content: flex-start; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
:deep(.arco-table-tr) { cursor: pointer; }
</style>
