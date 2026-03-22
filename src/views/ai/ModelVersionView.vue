<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>模型版本管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="filters.model_id" placeholder="模型名称" allow-clear style="width: 160px" @change="loadModels">
          <a-option value="MiniClaw">MiniClaw</a-option>
          <a-option value="Behavior">Behavior</a-option>
          <a-option value="PetLLM">PetLLM</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 140px" @change="loadModels">
          <a-option value="stable">生产环境</a-option>
          <a-option value="testing">测试环境</a-option>
          <a-option value="deprecated">已废弃</a-option>
          <a-option value="rolling_back">回滚中</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showRegisterModal">
          <template #icon><icon-plus /></template>
          注册新模型
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
          <span class="model-name">{{ record.model_name }}</span>
        </template>
        <template #current_version="{ record }">
          <a-space>
            <span>{{ record.current_version || '--' }}</span>
            <a-tag v-if="record.status === 'stable'" color="green" size="small">🟢</a-tag>
          </a-space>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template #published_at="{ record }">
          {{ formatDate(record.published_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="showVersionModal(record)">版本</a-button>
            <a-button
              v-if="record.status !== 'deprecated'"
              type="text"
              size="small"
              status="warning"
              @click="handleRollback(record)"
            >回滚</a-button>
            <a-button
              v-if="record.status !== 'deprecated'"
              type="text"
              size="small"
              status="danger"
              @click="handleDeprecate(record)"
            >废弃</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 注册新模型弹窗 -->
    <a-modal
      v-model:visible="registerModalVisible"
      title="注册新模型"
      :width="560"
      @ok="handleRegister"
      @cancel="registerModalVisible = false"
    >
      <a-form :model="registerForm" layout="vertical">
        <a-form-item label="模型标识" required>
          <a-input v-model="registerForm.model_id" placeholder="如 MiniClaw" />
        </a-form-item>
        <a-form-item label="模型名称" required>
          <a-input v-model="registerForm.model_name" placeholder="如 MiniClaw v2.1" />
        </a-form-item>
        <a-form-item label="模型类型" required>
          <a-select v-model="registerForm.model_type" placeholder="选择模型类型">
            <a-option value="openai">OpenAI</a-option>
            <a-option value="anthropic">Anthropic</a-option>
            <a-option value="local">本地模型</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="版本号" required>
          <a-input v-model="registerForm.version_id" placeholder="如 v2.1.3" />
        </a-form-item>
        <a-form-item label="版本名称">
          <a-input v-model="registerForm.version_name" placeholder="如 正式版" />
        </a-form-item>
        <a-form-item label="端点地址">
          <a-input v-model="registerForm.endpoint_url" placeholder="https://api.openai.com/..." />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 版本历史弹窗 -->
    <a-modal
      v-model:visible="versionModalVisible"
      title="版本历史"
      :width="640"
      @cancel="versionModalVisible = false"
      :footer="null"
    >
      <a-table
        :columns="versionColumns"
        :data="versionHistory"
        :loading="versionLoading"
        row-key="id"
        :pagination="false"
      >
        <template #version_id="{ record }">
          <a-space>
            <span>{{ record.version_id }}</span>
            <a-tag v-if="record.is_current" color="green" size="small">当前</a-tag>
          </a-space>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #created_at="{ record }">
          {{ formatDate(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-button
            v-if="!record.is_current && record.status !== 'deprecated'"
            type="text"
            size="small"
            status="warning"
            @click="handleVersionRollback(record)"
          >回滚至此</a-button>
        </template>
      </a-table>

      <!-- 发布新版本 -->
      <div class="publish-section">
        <a-divider>发布新版本</a-divider>
        <a-form :model="publishForm" layout="inline">
          <a-form-item label="版本号" required>
            <a-input v-model="publishForm.version_id" placeholder="如 v2.2.0" style="width: 140px" />
          </a-form-item>
          <a-form-item label="版本说明">
            <a-input v-model="publishForm.version_name" placeholder="版本说明" style="width: 200px" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handlePublish" :loading="publishing">发布</a-button>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  getAiModels,
  postAiModel,
  postAiModelRollback,
  postAiModelDeprecate,
  getAiModelVersions,
  postAiModelVersion
} from '@/api/ai'

const loading = ref(false)
const modelList = ref([])

const filters = reactive({
  model_id: undefined,
  status: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '模型名称', slotName: 'model_name', width: 200 },
  { title: '当前版本', slotName: 'current_version', width: 160 },
  { title: '状态', slotName: 'status', width: 120 },
  { title: '发布事件', slotName: 'published_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const versionColumns = [
  { title: '版本号', slotName: 'version_id', width: 160 },
  { title: '版本名称', dataIndex: 'version_name', ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '发布时间', slotName: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const registerModalVisible = ref(false)
const registerForm = reactive({
  model_id: '',
  model_name: '',
  model_type: 'openai',
  version_id: '',
  version_name: '',
  endpoint_url: ''
})

const versionModalVisible = ref(false)
const versionHistory = ref([])
const versionLoading = ref(false)
const currentModel = ref(null)
const publishing = ref(false)
const publishForm = reactive({ version_id: '', version_name: '' })

const getStatusColor = (status) => ({
  stable: 'green',
  testing: 'blue',
  dev: 'gray',
  deprecated: 'red',
  rolling_back: 'orange'
}[status] || 'gray')

const getStatusText = (status) => ({
  stable: '生产环境',
  testing: '测试环境',
  dev: '开发中',
  deprecated: '已废弃',
  rolling_back: '回滚中'
}[status] || status)

const formatDate = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleDateString('zh-CN')
}

const loadModels = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.model_id) params.model_id = filters.model_id
    if (filters.status) params.status = filters.status

    const res = await getAiModels(params)
    if (res.code === 0) {
      modelList.value = res.data.list || []
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
  loadModels()
}

const handlePageSizeChange = (size) => {
  pagination.pageSize = size
  pagination.current = 1
  loadModels()
}

const showRegisterModal = () => {
  Object.assign(registerForm, { model_id: '', model_name: '', model_type: 'openai', version_id: '', version_name: '', endpoint_url: '' })
  registerModalVisible.value = true
}

const handleRegister = async () => {
  if (!registerForm.model_id || !registerForm.model_name || !registerForm.version_id) {
    Message.warning('请填写必填项')
    return
  }
  try {
    const res = await postAiModel(registerForm)
    if (res.code === 0) {
      Message.success('注册成功')
      registerModalVisible.value = false
      loadModels()
    } else {
      Message.error(res.message || '注册失败')
    }
  } catch (e) {
    Message.error('注册失败')
  }
}

const showVersionModal = async (record) => {
  currentModel.value = record
  versionModalVisible.value = true
  versionLoading.value = true
  try {
    const res = await getAiModelVersions(record.id)
    if (res.code === 0) {
      versionHistory.value = res.data || []
    }
  } finally {
    versionLoading.value = false
  }
}

const handleRollback = (record) => {
  Modal.warning({
    title: '确认回滚',
    content: `确定要回滚 ${record.model_name} 吗？这将切换到上一稳定版本。`,
    okText: '确认回滚',
    onOk: async () => {
      try {
        const res = await postAiModelRollback(record.id, { reason: '手动回滚' })
        if (res.code === 0) {
          Message.success('回滚成功')
          loadModels()
        } else {
          Message.error(res.message || '回滚失败')
        }
      } catch (e) {
        Message.error('回滚失败')
      }
    }
  })
}

const handleDeprecate = (record) => {
  Modal.danger({
    title: '确认废弃',
    content: `确定要废弃 ${record.model_name} ${record.current_version} 吗？`,
    okText: '确认废弃',
    onOk: async () => {
      try {
        const res = await postAiModelDeprecate(record.id, {})
        if (res.code === 0) {
          Message.success('已废弃')
          loadModels()
        } else {
          Message.error(res.message || '废弃失败')
        }
      } catch (e) {
        Message.error('废弃失败')
      }
    }
  })
}

const handleVersionRollback = async (record) => {
  try {
    const res = await postAiModelRollback(currentModel.value.id, {
      target_version: record.version_id,
      reason: '版本历史回滚'
    })
    if (res.code === 0) {
      Message.success('回滚成功')
      showVersionModal(currentModel.value)
      loadModels()
    }
  } catch (e) {
    Message.error('回滚失败')
  }
}

const handlePublish = async () => {
  if (!publishForm.version_id) {
    Message.warning('请填写版本号')
    return
  }
  publishing.value = true
  try {
    const res = await postAiModelVersion(currentModel.value.id, publishForm)
    if (res.code === 0) {
      Message.success('发布成功')
      publishForm.version_id = ''
      publishForm.version_name = ''
      showVersionModal(currentModel.value)
      loadModels()
    }
  } finally {
    publishing.value = false
  }
}

onMounted(() => {
  loadModels()
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
.model-name { font-weight: 500; }
.publish-section { margin-top: 16px; }
</style>
