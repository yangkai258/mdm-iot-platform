<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>AI 模型管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-statistic title="模型总数" :value="stats.total" :value-style="{ color: '#1650ff' }">
          <template #prefix><icon-robot size="20" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="生产环境" :value="stats.production" :value-style="{ color: '#00b42a' }">
          <template #prefix><icon-check-circle size="20" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="测试环境" :value="stats.testing" :value-style="{ color: '#ff7d00' }">
          <template #prefix><icon-clock size="20" /></template>
        </a-statistic>
      </a-col>
      <a-col :span="6">
        <a-statistic title="已部署" :value="stats.deployed" :value-style="{ color: '#14c9c0' }">
          <template #prefix><icon-upload size="20" /></template>
        </a-statistic>
      </a-col>
    </a-row>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索模型名称/标识"
          style="width: 200px"
          search-button
          @search="loadModels"
        />
        <a-select v-model="filters.model_type" placeholder="模型类型" allow-clear style="width: 160px" @change="loadModels">
          <a-option value="llm">大语言模型</a-option>
          <a-option value="embedding">Embedding</a-option>
          <a-option value="vision">视觉模型</a-option>
          <a-option value="speech">语音模型</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 140px" @change="loadModels">
          <a-option value="production">生产环境</a-option>
          <a-option value="testing">测试环境</a-option>
          <a-option value="deprecated">已废弃</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showUploadModal">
          <template #icon><icon-upload /></template>
          上传模型
        </a-button>
        <a-button @click="loadModels">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="modelList"
        :loading="loading"
        row-key="id"
        :pagination="pagination"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #model_name="{ record }">
          <div class="model-info">
            <span class="model-name">{{ record.model_name }}</span>
            <span class="model-id">{{ record.model_id }}</span>
          </div>
        </template>
        <template #model_type="{ record }">
          <a-tag :color="getModelTypeColor(record.model_type)">
            {{ getModelTypeText(record.model_type) }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #is_deployed="{ record }">
          <a-switch
            v-model="record.is_deployed"
            :loading="record.deploying"
            size="small"
            @change="(val) => handleDeployToggle(record, val)"
          />
        </template>
        <template #file_size="{ record }">
          {{ formatFileSize(record.file_size) }}
        </template>
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showDetailModal(record)">详情</a-button>
            <a-button type="text" size="small" @click="showDeployModal(record)">部署</a-button>
            <a-button
              v-if="record.status !== 'deprecated'"
              type="text"
              size="small"
              status="danger"
              @click="handleDelete(record)"
            >删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 上传模型弹窗 -->
    <a-modal
      v-model:visible="uploadModalVisible"
      title="上传模型"
      :width="600"
      @ok="handleUpload"
      @cancel="uploadModalVisible = false"
      :mask-closable="false"
    >
      <a-form :model="uploadForm" layout="vertical">
        <a-form-item label="模型名称" required>
          <a-input v-model="uploadForm.model_name" placeholder="请输入模型名称" />
        </a-form-item>
        <a-form-item label="模型标识" required>
          <a-input v-model="uploadForm.model_id" placeholder="请输入模型标识（如 mini-claw-v1）" />
        </a-form-item>
        <a-form-item label="模型类型" required>
          <a-select v-model="uploadForm.model_type" placeholder="请选择模型类型">
            <a-option value="llm">大语言模型</a-option>
            <a-option value="embedding">Embedding</a-option>
            <a-option value="vision">视觉模型</a-option>
            <a-option value="speech">语音模型</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模型描述">
          <a-textarea v-model="uploadForm.description" placeholder="请输入模型描述" :rows="3" />
        </a-form-item>
        <a-form-item label="模型文件" required>
          <a-upload
            ref="uploadRef"
            :custom-request="handleFileRequest"
            :limit="1"
            accept=".bin,.pt,.pth,.onnx,.safetensors,.gguf"
            @success="handleUploadSuccess"
            @error="handleUploadError"
          >
            <template #upload-button>
              <div class="upload-trigger">
                <icon-upload size="32" />
                <span>点击或拖拽文件到此处上传</span>
                <span class="upload-hint">支持 .bin, .pt, .onnx, .safetensors, .gguf 格式</span>
              </div>
            </template>
          </a-upload>
        </a-form-item>
        <a-form-item label="版本号">
          <a-input v-model="uploadForm.version" placeholder="如 1.0.0" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 模型详情弹窗 -->
    <a-modal
      v-model:visible="detailModalVisible"
      title="模型详情"
      :width="560"
      @ok="detailModalVisible = false"
      @cancel="detailModalVisible = false"
    >
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="模型名称">{{ currentModel?.model_name }}</a-descriptions-item>
        <a-descriptions-item label="模型标识">{{ currentModel?.model_id }}</a-descriptions-item>
        <a-descriptions-item label="模型类型">{{ getModelTypeText(currentModel?.model_type) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentModel?.status)">
            {{ getStatusText(currentModel?.status) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="文件大小">{{ formatFileSize(currentModel?.file_size) }}</a-descriptions-item>
        <a-descriptions-item label="版本号">{{ currentModel?.version || '--' }}</a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">{{ currentModel?.description || '--' }}</a-descriptions-item>
        <a-descriptions-item label="创建时间" :span="2">{{ formatDate(currentModel?.created_at) }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 部署配置弹窗 -->
    <a-modal
      v-model:visible="deployModalVisible"
      title="部署配置"
      :width="500"
      @ok="handleDeploy"
      @cancel="deployModalVisible = false"
    >
      <a-form :model="deployForm" layout="vertical">
        <a-form-item label="部署环境" required>
          <a-select v-model="deployForm.environment" placeholder="请选择部署环境">
            <a-option value="production">生产环境</a-option>
            <a-option value="staging">预发环境</a-option>
            <a-option value="testing">测试环境</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="实例数量">
          <a-input-number v-model="deployForm.instances" :min="1" :max="10" placeholder="默认 1" />
        </a-form-item>
        <a-form-item label="GPU 配置">
          <a-select v-model="deployForm.gpu_type" placeholder="请选择 GPU 类型" allow-clear>
            <a-option value="T4">T4 (16GB)</a-option>
            <a-option value="V100">V100 (32GB)</a-option>
            <a-option value="A100">A100 (80GB)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="内存限制">
          <a-input-number v-model="deployForm.memory_limit" :min="1" :max="128" placeholder="GB" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  getAiModels,
  postAiModel,
  deleteAiModel,
  postAiModelDeploy,
  postAiModelUndeploy
} from '@/api/ai'

const loading = ref(false)
const modelList = ref([])
const stats = reactive({
  total: 0,
  production: 0,
  testing: 0,
  deployed: 0
})

const filters = reactive({
  keyword: '',
  model_type: undefined,
  status: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const uploadModalVisible = ref(false)
const detailModalVisible = ref(false)
const deployModalVisible = ref(false)
const currentModel = ref(null)

const uploadForm = reactive({
  model_name: '',
  model_id: '',
  model_type: '',
  description: '',
  version: '',
  file_path: ''
})

const deployForm = reactive({
  environment: 'production',
  instances: 1,
  gpu_type: '',
  memory_limit: 16
})

const columns = [
  { title: '模型名称', slotName: 'model_name', width: 200, ellipsis: true },
  { title: '类型', slotName: 'model_type', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '已部署', slotName: 'is_deployed', width: 80 },
  { title: '文件大小', slotName: 'file_size', width: 100 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

const getModelTypeColor = (type) => ({
  llm: 'arcoblue',
  embedding: 'green',
  vision: 'purple',
  speech: 'orange'
}[type] || 'gray')

const getModelTypeText = (type) => ({
  llm: '大语言模型',
  embedding: 'Embedding',
  vision: '视觉模型',
  speech: '语音模型'
}[type] || type)

const getStatusColor = (status) => ({
  production: 'green',
  testing: 'orange',
  deprecated: 'red'
}[status] || 'gray')

const getStatusText = (status) => ({
  production: '生产环境',
  testing: '测试环境',
  deprecated: '已废弃'
}[status] || status)

const formatDate = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const formatFileSize = (bytes) => {
  if (!bytes) return '--'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (bytes >= 1024 && i < units.length - 1) {
    bytes /= 1024
    i++
  }
  return `${bytes.toFixed(1)} ${units[i]}`
}

const loadModels = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.model_type) params.model_type = filters.model_type
    if (filters.status) params.status = filters.status

    const res = await getAiModels(params)
    if (res.code === 0) {
      modelList.value = (res.data.list || []).map(m => ({ ...m, deploying: false }))
      pagination.total = res.data.total || 0

      // 计算统计数据
      stats.total = modelList.value.length
      stats.production = modelList.value.filter(m => m.status === 'production').length
      stats.testing = modelList.value.filter(m => m.status === 'testing').length
      stats.deployed = modelList.value.filter(m => m.is_deployed).length
    }
  } catch (e) {
    Message.error('加载模型列表失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadModels()
}

const handlePageSizeChange = (size) => {
  pagination.pageSize = size
  pagination.current = 1
  loadModels()
}

const showUploadModal = () => {
  Object.assign(uploadForm, {
    model_name: '',
    model_id: '',
    model_type: '',
    description: '',
    version: '',
    file_path: ''
  })
  uploadModalVisible.value = true
}

const showDetailModal = (record) => {
  currentModel.value = record
  detailModalVisible.value = true
}

const showDeployModal = (record) => {
  currentModel.value = record
  Object.assign(deployForm, {
    environment: 'production',
    instances: 1,
    gpu_type: '',
    memory_limit: 16
  })
  deployModalVisible.value = true
}

const handleFileRequest = (options) => {
  const { fileItem, onSuccess, onError } = options
  // 模拟文件上传
  setTimeout(() => {
    uploadForm.file_path = fileItem.file.name
    onSuccess(fileItem)
  }, 500)
}

const handleUploadSuccess = (fileItem) => {
  Message.success('文件上传成功')
}

const handleUploadError = () => {
  Message.error('文件上传失败')
}

const handleUpload = async () => {
  if (!uploadForm.model_name || !uploadForm.model_id || !uploadForm.model_type) {
    Message.warning('请填写必填项')
    return
  }
  try {
    const res = await postAiModel(uploadForm)
    if (res.code === 0) {
      Message.success('模型上传成功')
      uploadModalVisible.value = false
      loadModels()
    } else {
      Message.error(res.message || '上传失败')
    }
  } catch (e) {
    Message.error('上传失败')
  }
}

const handleDeployToggle = async (record, val) => {
  record.deploying = true
  try {
    if (val) {
      await postAiModelDeploy(record.id, { environment: 'production' })
      Message.success('部署成功')
    } else {
      await postAiModelUndeploy(record.id)
      Message.success('取消部署成功')
    }
    record.is_deployed = val
  } catch (e) {
    Message.error('操作失败')
    record.is_deployed = !val
  } finally {
    record.deploying = false
  }
}

const handleDeploy = async () => {
  if (!deployForm.environment) {
    Message.warning('请选择部署环境')
    return
  }
  try {
    const res = await postAiModelDeploy(currentModel.value.id, deployForm)
    if (res.code === 0) {
      Message.success('部署配置成功')
      deployModalVisible.value = false
      loadModels()
    } else {
      Message.error(res.message || '部署失败')
    }
  } catch (e) {
    Message.error('部署失败')
  }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除模型 "${record.model_name}" 吗？此操作不可恢复。`,
    okText: '确认删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        const res = await deleteAiModel(record.id)
        if (res.code === 0) {
          Message.success('删除成功')
          loadModels()
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
  loadModels()
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
.model-info { display: flex; flex-direction: column; }
.model-name { font-weight: 500; }
.model-id { font-size: 12px; color: #999; }
.upload-trigger {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  width: 100%; height: 120px; border: 1px dashed #ccc; border-radius: 8px;
  cursor: pointer; color: #666; transition: all 0.3s;
}
.upload-trigger:hover { border-color: #1650ff; color: #1650ff; }
.upload-hint { font-size: 12px; color: #999; margin-top: 8px; }
</style>
