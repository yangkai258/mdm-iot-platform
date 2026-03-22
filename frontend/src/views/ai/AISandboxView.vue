<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>AI 沙箱测试</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 创建测试 -->
    <a-card class="create-test-card" title="创建测试">
      <a-form :model="testForm" layout="vertical">
        <a-form-item label="测试类型" required>
          <a-radio-group v-model="testForm.test_type">
            <a-radio value="unit">单元测试</a-radio>
            <a-radio value="integration">集成测试</a-radio>
            <a-radio value="performance">压力测试</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="选择模型" required>
          <a-select v-model="testForm.model_version_id" placeholder="请选择模型版本" style="width: 300px">
            <a-option
              v-for="m in modelOptions"
              :key="m.id"
              :value="m.id"
              :label="`${m.model_name} ${m.current_version || m.version_id}`"
            />
          </a-select>
        </a-form-item>
        <a-form-item label="测试数据">
          <a-space direction="vertical">
            <a-select v-model="testForm.test_dataset" placeholder="选择测试数据集" style="width: 300px" allow-create>
              <a-option value="dataset_intent">意图识别测试集</a-option>
              <a-option value="dataset_response">响应生成测试集</a-option>
              <a-option value="dataset_action">动作选择测试集</a-option>
            </a-select>
            <a-button type="outline" @click="showUploadDialog">
              <template #icon><icon-upload /></template>
              上传测试数据
            </a-button>
          </a-space>
        </a-form-item>
        <a-form-item label="测试描述">
          <a-textarea v-model="testForm.description" placeholder="输入测试说明..." :max-length="500" show-word-limit style="width: 500px" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleStartTest" :loading="starting">
            <template #icon><icon-play /></template>
            开始测试
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 测试历史 -->
    <a-card class="test-history-card" title="测试历史">
      <template #extra>
        <a-button type="text" @click="loadTests">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </template>
      <a-table
        :columns="columns"
        :data="testList"
        :loading="loading"
        row-key="id"
        :pagination="pagination"
        @page-change="handlePageChange"
      >
        <template #test_type="{ record }">
          <a-tag :color="getTestTypeColor(record.test_type)">
            {{ getTestTypeText(record.test_type) }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            <template #icon v-if="record.status === 'running'"><icon-loading /></template>
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewResult(record)">查看</a-button>
            <a-button
              v-if="record.status === 'pending' || record.status === 'running'"
              type="text"
              size="small"
              status="danger"
              @click="cancelTest(record)"
            >取消</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 测试结果弹窗 -->
    <a-modal
      v-model:visible="resultModalVisible"
      title="测试结果"
      :width="720"
      :footer="null"
    >
      <div v-if="currentResult" class="result-detail">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="任务ID">{{ currentResult.id }}</a-descriptions-item>
          <a-descriptions-item label="测试类型">{{ getTestTypeText(currentResult.test_type) }}</a-descriptions-item>
          <a-descriptions-item label="状态">{{ getStatusText(currentResult.status) }}</a-descriptions-item>
          <a-descriptions-item label="完成时间">{{ formatTime(currentResult.completed_at) }}</a-descriptions-item>
        </a-descriptions>

        <a-divider>测试结果</a-divider>
        <a-space direction="vertical" style="width: 100%">
          <a-result v-if="currentResult.status === 'completed'" status="success" title="测试通过" />
          <a-result v-else-if="currentResult.status === 'failed'" status="error" title="测试失败" />
          <a-result v-else status="info" title="测试进行中" />

          <div v-if="currentResult.results" class="result-metrics">
            <a-row :gutter="12">
              <a-col :span="6">
                <a-statistic title="通过用例" :value="currentResult.results.passed || 0" />
              </a-col>
              <a-col :span="6">
                <a-statistic title="失败用例" :value="currentResult.results.failed || 0" />
              </a-col>
              <a-col :span="6">
                <a-statistic title="平均延迟" :value="currentResult.results.avg_latency || 0" suffix="ms" />
              </a-col>
              <a-col :span="6">
                <a-statistic title="成功率" :value="currentResult.results.success_rate || 0" suffix="%" :precision="1" />
              </a-col>
            </a-row>
          </div>
        </a-space>
      </div>
    </a-modal>

    <!-- 上传测试数据弹窗 -->
    <a-modal
      v-model:visible="uploadModalVisible"
      title="上传测试数据"
      :width="480"
      @ok="handleUpload"
    >
      <a-upload
        draggable
        action="#"
        :custom-request="dummyUpload"
        accept=".json,.csv"
        @success="handleUploadSuccess"
      >
        <template #upload-button>
          <div>
            <icon-upload size="40" />
            <div class="upload-text">点击或拖拽上传 JSON/CSV 测试数据</div>
          </div>
        </template>
      </a-upload>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getAiSandboxTests,
  getAiSandboxTestResult,
  postAiSandboxTest,
  getAiModels
} from '@/api/ai'

const loading = ref(false)
const testList = ref([])
const starting = ref(false)
const modelOptions = ref([])

const testForm = reactive({
  test_type: 'unit',
  model_version_id: undefined,
  test_dataset: undefined,
  description: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: '任务ID', dataIndex: 'id', ellipsis: true, width: 120 },
  { title: '模型', dataIndex: 'model_name', ellipsis: true },
  { title: '测试类型', slotName: 'test_type', width: 120 },
  { title: '状态', slotName: 'status', width: 120 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const resultModalVisible = ref(false)
const currentResult = ref(null)

const uploadModalVisible = ref(false)

const getTestTypeColor = (type) => ({
  unit: 'arcoblue',
  integration: 'green',
  performance: 'purple'
}[type] || 'gray')

const getTestTypeText = (type) => ({
  unit: '单元测试',
  integration: '集成测试',
  performance: '压力测试'
}[type] || type)

const getStatusColor = (status) => ({
  pending: 'gray',
  running: 'blue',
  completed: 'green',
  failed: 'red'
}[status] || 'gray')

const getStatusText = (status) => ({
  pending: '等待中',
  running: '进行中',
  completed: '通过',
  failed: '失败'
}[status] || status)

const formatTime = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const loadModels = async () => {
  try {
    const res = await getAiModels({ status: 'stable' })
    if (res.code === 0) {
      modelOptions.value = res.data.list || []
    }
  } catch (e) {
    // ignore
  }
}

const loadTests = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    const res = await getAiSandboxTests(params)
    if (res.code === 0) {
      testList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadTests()
}

const handleStartTest = async () => {
  if (!testForm.model_version_id) {
    Message.warning('请选择模型')
    return
  }
  starting.value = true
  try {
    const res = await postAiSandboxTest({
      test_name: testForm.description || `${getTestTypeText(testForm.test_type)}-${Date.now()}`,
      model_version_id: testForm.model_version_id,
      test_type: testForm.test_type,
      test_cases: [],
      description: testForm.description
    })
    if (res.code === 0) {
      Message.success('测试任务已创建')
      loadTests()
      testForm.description = ''
    } else {
      Message.error(res.message || '创建失败')
    }
  } finally {
    starting.value = false
  }
}

const viewResult = async (record) => {
  currentResult.value = record
  resultModalVisible.value = true
  // 如果已完成，尝试刷新结果
  if (record.status === 'completed' || record.status === 'failed') {
    try {
      const res = await getAiSandboxTestResult(record.id)
      if (res.code === 0) {
        currentResult.value = { ...currentResult.value, ...res.data }
      }
    } catch (e) {
      // ignore
    }
  }
}

const cancelTest = (record) => {
  Message.info('取消测试功能开发中')
}

const showUploadDialog = () => {
  uploadModalVisible.value = true
}

const dummyUpload = ({ file, onSuccess }) => {
  setTimeout(() => { onSuccess(file) }, 500)
}

const handleUpload = () => {
  uploadModalVisible.value = false
  Message.success('测试数据上传成功')
}

const handleUploadSuccess = (file) => {
  testForm.test_dataset = file.name
}

onMounted(() => {
  loadModels()
  loadTests()
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }

.create-test-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.test-history-card {
  border-radius: 8px;
}

.result-detail { padding: 8px 0; }
.result-metrics { margin-top: 16px; }

.upload-text {
  margin-top: 8px;
  color: var(--color-text-3);
  font-size: 13px;
}
</style>
