<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>技术架构</a-breadcrumb-item>
      <a-breadcrumb-item>端侧推理模型</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="模型总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="已部署" :value="stats.deployed" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="推理中" :value="stats.inferencing" :value-style="{ color: '#165dff' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="推理次数" :value="stats.totalInferences" suffix="次" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索框 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索模型名称或版本"
          style="width: 280px"
          @search="loadModels"
          search-button
        />
        <a-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px" allow-clear>
          <a-option :value="1">待部署</a-option>
          <a-option :value="2">部署中</a-option>
          <a-option :value="3">运行中</a-option>
          <a-option :value="4">已停止</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showUploadModal">上传模型</a-button>
        <a-button @click="loadModels">刷新</a-button>
      </a-space>
    </div>

    <!-- 模型列表 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="filteredModels"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
      >
        <template #modelName="{ record }">
          <a-space>
            <a-avatar :size="24" :style="{ backgroundColor: '#165dff' }">
              {{ record.name.charAt(0).toUpperCase() }}
            </a-avatar>
            <span>{{ record.name }}</span>
          </a-space>
        </template>
        <template #status="{ record }">
          <a-badge :status="getStatusBadge(record.status)" :text="getStatusText(record.status)" />
        </template>
        <template #fileSize="{ record }">
          {{ formatFileSize(record.file_size) }}
        </template>
        <template #deployTime="{ record }">
          {{ formatTime(record.deploy_time) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button v-if="record.status === 1 || record.status === 4" type="text" size="small" @click="handleDeploy(record)">部署</a-button>
            <a-button v-if="record.status === 3" type="text" size="small" status="warning" @click="handleStop(record)">停止</a-button>
            <a-button v-if="record.status === 2" type="text" size="small" loading>部署中</a-button>
            <a-button type="text" size="small" @click="showConfigModal(record)">配置</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 上传模型弹窗 -->
    <a-modal v-model:visible="uploadModalVisible" title="上传模型" @ok="handleUpload" :confirm-loading="uploading" :width="520">
      <a-form :model="uploadForm" layout="vertical">
        <a-form-item label="模型名称" required>
          <a-input v-model="uploadForm.name" placeholder="例如: yolo-nano-v3" />
        </a-form-item>
        <a-form-item label="模型版本" required>
          <a-input v-model="uploadForm.version" placeholder="例如: 1.0.0" />
        </a-form-item>
        <a-form-item label="模型类型" required>
          <a-select v-model="uploadForm.model_type" placeholder="选择模型类型">
            <a-option value="detection">目标检测</a-option>
            <a-option value="classification">图像分类</a-option>
            <a-option value="segmentation">图像分割</a-option>
            <a-option value="nlp">NLP处理</a-option>
            <a-option value="audio">音频识别</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模型描述">
          <a-textarea v-model="uploadForm.description" placeholder="模型功能描述" :rows="3" />
        </a-form-item>
        <a-form-item label="模型文件" required>
          <a-upload :auto-upload="false" :show-upload-list="true" @change="handleFileChange">
            <a-button>选择文件</a-button>
          </a-upload>
        </a-form-item>
        <a-form-item label="硬件平台">
          <a-select v-model="uploadForm.hardware_target" placeholder="选择目标硬件">
            <a-option value="esp32">ESP32</a-option>
            <a-option value="stm32">STM32</a-option>
            <a-option value="riscv">RISC-V</a-option>
            <a-option value="generic">通用ARM</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 模型配置弹窗 -->
    <a-modal v-model:visible="configModalVisible" title="模型配置" @ok="handleConfigSave" :confirm-loading="configSaving" :width="600">
      <a-form :model="configForm" layout="vertical">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="模型名称">{{ currentModel?.name }}</a-descriptions-item>
          <a-descriptions-item label="模型版本">{{ currentModel?.version }}</a-descriptions-item>
          <a-descriptions-item label="模型类型">{{ currentModel?.model_type }}</a-descriptions-item>
          <a-descriptions-item label="文件大小">{{ formatFileSize(currentModel?.file_size) }}</a-descriptions-item>
        </a-descriptions>
        <a-divider>推理参数</a-divider>
        <a-form-item label="最大并发数">
          <a-input-number v-model="configForm.max_instances" :min="1" :max="10" style="width: 100%" />
        </a-form-item>
        <a-form-item label="推理超时(ms)">
          <a-input-number v-model="configForm.timeout_ms" :min="100" :max="60000" :step="100" style="width: 100%" />
        </a-form-item>
        <a-form-item label="启用缓存">
          <a-switch v-model="configForm.enable_cache" />
        </a-form-item>
        <a-form-item label="最低内存阈值(MB)">
          <a-input-number v-model="configForm.min_memory_mb" :min="64" :max="512" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 部署目标选择弹窗 -->
    <a-modal v-model:visible="deployModalVisible" title="选择部署目标" @ok="handleDeployConfirm" :confirm-loading="deployLoading" :width="520">
      <a-form layout="vertical">
        <a-form-item label="目标设备">
          <a-select v-model="deployForm.device_ids" multiple placeholder="选择目标设备">
            <a-option v-for="d in onlineDevices" :key="d.device_id" :value="d.device_id">
              {{ d.device_id }} ({{ d.hardware_model }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级">
          <a-radio-group v-model="deployForm.priority">
            <a-radio value="high">高优先级</a-radio>
            <a-radio value="normal">普通</a-radio>
            <a-radio value="low">低优先级</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/edge'

const loading = ref(false)
const uploading = ref(false)
const configSaving = ref(false)
const deployLoading = ref(false)

const searchKeyword = ref('')
const filterStatus = ref(null)

const uploadModalVisible = ref(false)
const configModalVisible = ref(false)
const deployModalVisible = ref(false)

const currentModel = ref(null)

const uploadForm = reactive({
  name: '', version: '', model_type: '', description: '', file: null, hardware_target: ''
})

const configForm = reactive({
  max_instances: 1, timeout_ms: 5000, enable_cache: true, min_memory_mb: 128
})

const deployForm = reactive({
  device_ids: [], priority: 'normal'
})

const stats = reactive({ total: 0, deployed: 0, inferencing: 0, totalInferences: 0 })

const models = ref([])
const onlineDevices = ref([])

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '模型名称', slotName: 'modelName', width: 180 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '类型', dataIndex: 'model_type', width: 100 },
  { title: '文件大小', slotName: 'fileSize', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '部署时间', slotName: 'deployTime', width: 180 },
  { title: '操作', slotName: 'actions', width: 260 }
]

const filteredModels = computed(() => {
  let result = models.value
  if (filterStatus.value !== null) {
    result = result.filter(m => m.status === filterStatus.value)
  }
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(m => m.name.toLowerCase().includes(kw) || m.version.toLowerCase().includes(kw))
  }
  return result
})

const loadModels = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/models`)
    if (res.data.code === 0) {
      models.value = res.data.data.list
      pagination.total = res.data.data.pagination?.total || models.value.length
    }
  } catch {
    models.value = [
      { id: 1, name: 'yolo-nano-v3', version: '1.0.0', model_type: 'detection', file_size: 2097152, status: 3, deploy_time: '2026-03-19 10:00:00', description: '轻量目标检测模型', hardware_target: 'esp32' },
      { id: 2, name: 'mobilenet-v2', version: '2.1.0', model_type: 'classification', file_size: 4194304, status: 3, deploy_time: '2026-03-18 14:30:00', description: '图像分类模型', hardware_target: 'stm32' },
      { id: 3, name: 'bert-tiny', version: '1.0.0', model_type: 'nlp', file_size: 52428800, status: 1, deploy_time: null, description: '轻量NLP模型', hardware_target: 'riscv' },
      { id: 4, name: 'resnet-18', version: '1.2.0', model_type: 'classification', file_size: 44695552, status: 2, deploy_time: '2026-03-20 09:00:00', description: 'ResNet图像分类', hardware_target: 'generic' },
      { id: 5, name: 'yolov5s', version: '1.0.0', model_type: 'detection', file_size: 14680064, status: 4, deploy_time: '2026-03-15 11:00:00', description: 'YOLOv5目标检测', hardware_target: 'esp32' }
    ]
    pagination.total = models.value.length
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
  updateStats()
}

const loadDevices = async () => {
  try {
    const res = await axios.get('/api/v1/devices', { params: { is_online: true } })
    if (res.data.code === 0) onlineDevices.value = res.data.data.list
  } catch {
    onlineDevices.value = [
      { device_id: 'DEV001', hardware_model: 'MDM-Pro-200' },
      { device_id: 'DEV002', hardware_model: 'MDM-Mini-100' },
      { device_id: 'DEV003', hardware_model: 'MDM-Pro-200' }
    ]
  }
}

const updateStats = () => {
  stats.total = models.value.length
  stats.deployed = models.value.filter(m => m.status === 3).length
  stats.inferencing = models.value.filter(m => m.status === 3).length
  stats.totalInferences = 12847
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
}

const showUploadModal = () => { uploadModalVisible.value = true }

const handleFileChange = (file) => { uploadForm.file = file }

const handleUpload = () => {
  if (!uploadForm.name || !uploadForm.version || !uploadForm.model_type) {
    Message.warning('请填写必填项'); return
  }
  uploading.value = true
  setTimeout(() => {
    models.value.unshift({
      id: Date.now(), name: uploadForm.name, version: uploadForm.version,
      model_type: uploadForm.model_type, file_size: 2097152, status: 1,
      deploy_time: null, description: uploadForm.description, hardware_target: uploadForm.hardware_target
    })
    uploadModalVisible.value = false
    uploading.value = false
    Message.success('模型上传成功')
    Object.assign(uploadForm, { name: '', version: '', model_type: '', description: '', file: null, hardware_target: '' })
    updateStats()
  }, 1500)
}

const showConfigModal = (record) => {
  currentModel.value = { ...record }
  Object.assign(configForm, { max_instances: 1, timeout_ms: 5000, enable_cache: true, min_memory_mb: 128 })
  configModalVisible.value = true
}

const handleConfigSave = () => {
  configSaving.value = true
  setTimeout(() => {
    configSaving.value = false
    configModalVisible.value = false
    Message.success('配置已保存')
  }, 800)
}

const handleDeploy = (record) => {
  currentModel.value = record
  deployForm.device_ids = []
  deployForm.priority = 'normal'
  deployModalVisible.value = true
}

const handleDeployConfirm = () => {
  if (!deployForm.device_ids.length) {
    Message.warning('请选择目标设备'); return
  }
  deployLoading.value = true
  setTimeout(() => {
    const idx = models.value.findIndex(m => m.id === currentModel.value.id)
    if (idx !== -1) models.value[idx].status = 2
    deployLoading.value = false
    deployModalVisible.value = false
    Message.success(`模型 ${currentModel.value.name} 正在部署到 ${deployForm.device_ids.length} 个设备`)
    updateStats()
  }, 1200)
}

const handleStop = (record) => {
  const idx = models.value.findIndex(m => m.id === record.id)
  if (idx !== -1) models.value[idx].status = 4
  Message.success(`模型 ${record.name} 已停止`)
  updateStats()
}

const handleDelete = (record) => {
  models.value = models.value.filter(m => m.id !== record.id)
  Message.success('模型已删除')
  updateStats()
}

const getStatusBadge = (status) => ({ 1: 'default', 2: 'processing', 3: 'success', 4: 'warning' }[status] || 'default')
const getStatusText = (status) => ({ 1: '待部署', 2: '部署中', 3: '运行中', 4: '已停止' }[status] || '未知')
const formatFileSize = (bytes) => bytes ? (bytes / 1024 / 1024).toFixed(2) + ' MB' : '-'
const formatTime = (time) => time || '-'

onMounted(() => { loadModels(); loadDevices() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04); margin-bottom: 16px;
}
</style>
