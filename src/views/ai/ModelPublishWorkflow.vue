<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>模型版本管理</a-breadcrumb-item>
      <a-breadcrumb-item>发布工作流</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 选择模型 -->
    <a-card class="workflow-card" title="选择模型">
      <a-form :model="workflowForm" layout="vertical">
        <a-form-item label="选择模型" required>
          <a-select
            v-model="workflowForm.model_id"
            placeholder="请选择要发布的模型"
            style="width: 320px"
            @change="handleModelChange"
          >
            <a-option
              v-for="m in modelList"
              :key="m.id"
              :value="m.id"
              :label="m.model_name"
            />
          </a-select>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 发布信息 -->
    <a-card v-if="workflowForm.model_id" class="workflow-card" title="发布信息">
      <a-form :model="workflowForm" layout="vertical">
        <a-form-item label="目标版本" required>
          <a-input v-model="workflowForm.version_id" placeholder="如 v2.2.0" style="width: 300px" />
        </a-form-item>
        <a-form-item label="版本说明" required>
          <a-textarea
            v-model="workflowForm.change_log"
            placeholder="描述本次发布的主要变更..."
            :max-length="1000"
            show-word-limit
            style="width: 500px"
            :auto-size="{ minRows: 3, maxRows: 6 }"
          />
        </a-form-item>
        <a-form-item label="灰度发布比例">
          <a-slider
            v-model="workflowForm.rollout_percent"
            :min="0"
            :max="100"
            :step="10"
            show-input
            style="width: 300px"
          />
          <div class="slider-hint">
            {{ workflowForm.rollout_percent === 0 ? '不发布（仅保存）' :
               workflowForm.rollout_percent === 100 ? '全量发布' :
               `灰度发布至 ${workflowForm.rollout_percent}% 设备` }}
          </div>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSubmit" :loading="submitting">
              提交审核
            </a-button>
            <a-button @click="handleSaveDraft">
              保存草稿
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 审核流程 -->
    <a-card v-if="currentWorkflow" class="workflow-card" title="发布进度">
      <a-steps :current="getCurrentStep" linelayout="horizontal">
        <a-step title="填写信息" :status="getStepStatus('info')" />
        <a-step title="提交审核" :status="getStepStatus('review')" />
        <a-step title="审核通过" :status="getStepStatus('approved')" />
        <a-step title="发布上线" :status="getStepStatus('published')" />
      </a-steps>

      <!-- 审核详情 -->
      <div v-if="currentWorkflow.review_comments" class="review-comments">
        <a-divider>审核意见</a-divider>
        <a-alert :type="currentWorkflow.review_status === 'approved' ? 'success' : 'warning'">
          {{ currentWorkflow.review_comments }}
        </a-alert>
      </div>
    </a-card>

    <!-- 发布历史 -->
    <a-card class="workflow-card" title="发布历史">
      <a-table
        :columns="historyColumns"
        :data="publishHistory"
        :loading="historyLoading"
        row-key="id"
        :pagination="{ pageSize: 5 }"
      >
        <template #status="{ record }">
          <a-tag :color="getPublishStatusColor(record.status)">
            {{ getPublishStatusText(record.status) }}
          </a-tag>
        </template>
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { getAiModels, postAiModelVersion } from '@/api/ai'

const loading = ref(false)
const modelList = ref([])
const submitting = ref(false)
const historyLoading = ref(false)
const publishHistory = ref([])
const currentWorkflow = ref(null)

const workflowForm = reactive({
  model_id: undefined,
  version_id: '',
  change_log: '',
  rollout_percent: 100
})

const historyColumns = [
  { title: '版本', dataIndex: 'version_id', width: 120 },
  { title: '变更说明', dataIndex: 'change_log', ellipsis: true },
  { title: '状态', slotName: 'status', width: 120 },
  { title: '发布时间', slotName: 'created_at', width: 170 }
]

const getPublishStatusColor = (status) => ({
  draft: 'gray',
  pending: 'orange',
  approved: 'blue',
  rejected: 'red',
  published: 'green'
}[status] || 'gray')

const getPublishStatusText = (status) => ({
  draft: '草稿',
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  published: '已发布'
}[status] || status)

const formatTime = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const getCurrentStep = () => {
  if (!currentWorkflow.value) return 0
  const map = { draft: 0, pending: 1, approved: 2, published: 3 }
  return map[currentWorkflow.value.status] ?? 0
}

const getStepStatus = (step) => {
  if (!currentWorkflow.value) return 'wait'
  const stepMap = { info: 0, review: 1, approved: 2, published: 3 }
  const current = getCurrentStep()
  const target = stepMap[step]
  if (current > target) return 'finish'
  if (current === target) return 'process'
  return 'wait'
}

const loadModels = async () => {
  loading.value = true
  try {
    const res = await getAiModels({})
    if (res.code === 0) {
      modelList.value = res.data.list || []
    }
  } finally {
    loading.value = false
  }
}

const handleModelChange = (modelId) => {
  workflowForm.version_id = ''
  workflowForm.change_log = ''
  workflowForm.rollout_percent = 100
  currentWorkflow.value = null
  // 查找该模型的最新草稿
  loadPublishHistory(modelId)
}

const loadPublishHistory = async (modelId) => {
  historyLoading.value = true
  try {
    // TODO: 调用发布历史API
    publishHistory.value = []
  } finally {
    historyLoading.value = false
  }
}

const handleSubmit = async () => {
  if (!workflowForm.version_id) {
    Message.warning('请填写版本号')
    return
  }
  if (!workflowForm.change_log) {
    Message.warning('请填写变更说明')
    return
  }

  Modal.confirm({
    title: '确认提交审核',
    content: `确定要提交 ${workflowForm.version_id} 的发布审核吗？`,
    okText: '确认提交',
    onOk: async () => {
      submitting.value = true
      try {
        // TODO: 调用提交审核API
        Message.success('提交成功，等待审核')
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleSaveDraft = () => {
  if (!workflowForm.version_id) {
    Message.warning('请填写版本号')
    return
  }
  // TODO: 保存草稿API
  Message.success('草稿已保存')
}

onMounted(() => {
  loadModels()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.workflow-card { margin-bottom: 16px; border-radius: 8px; }
.slider-hint { color: var(--color-text-3); font-size: 12px; margin-top: 4px; }
.review-comments { margin-top: 16px; }
</style>
