<template>
  <div class="page-container">
<!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="配置文件总数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已激活" :value="stats.active" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="草稿" :value="stats.draft" :value-style="{ color: '#faad14' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已废弃" :value="stats.deprecated" :value-style="{ color: '#ff4d4f' }" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.policyType" placeholder="策略类型" allow-clear style="width: 140px" @change="loadConfigs">
              <a-option value="security">安全策略</a-option>
              <a-option value="network">网络策略</a-option>
              <a-option value="app">应用策略</a-option>
              <a-option value="device">设备策略</a-option>
            </a-select>
            <a-select v-model="filters.status" placeholder="配置状态" allow-clear style="width: 140px" @change="loadConfigs">
              <a-option value="active">已激活</a-option>
              <a-option value="draft">草稿</a-option>
              <a-option value="deprecated">已废弃</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索配置名称/编码" style="width: 200px" search-button @search="loadConfigs" />
            <a-button type="primary" @click="showCreateDrawer = true">新建配置</a-button>
            <a-button @click="loadConfigs">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 配置列表 -->
        <a-card class="config-card">
          <a-table
            :columns="columns"
            :data="configList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #policyType="{ record }">
              <a-tag :color="getPolicyTypeColor(record.policy_type)">{{ getPolicyTypeText(record.policy_type) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #deviceCount="{ record }">
              {{ record.device_count || 0 }}
            </template>
            <template #updatedAt="{ record }">
              {{ formatTime(record.updated_at) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="editConfig(record)">编辑</a-button>
                <a-button v-if="record.status === 'draft'" type="text" size="small" @click="activateConfig(record)">激活</a-button>
                <a-button v-if="record.status === 'active'" type="text" size="small" status="warning" @click="deprecateConfig(record)">废弃</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>

    <!-- 创建/编辑配置抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      :title="isEdit ? '编辑配置' : '新建配置'"
      :width="520"
      @before-ok="handleSave"
    >
      <a-form :model="form" layout="vertical" @submit-success="handleSaveSubmit">
        <a-form-item label="配置编码" required>
          <a-input v-model="form.config_code" placeholder="如 POL001" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="配置名称" required>
          <a-input v-model="form.config_name" placeholder="如 基础安全策略" />
        </a-form-item>
        <a-form-item label="策略类型" required>
          <a-select v-model="form.policy_type" placeholder="选择策略类型">
            <a-option value="security">安全策略</a-option>
            <a-option value="network">网络策略</a-option>
            <a-option value="app">应用策略</a-option>
            <a-option value="device">设备策略</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="配置描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="配置功能描述" />
        </a-form-item>
        <a-form-item label="配置内容 (JSON)">
          <a-textarea v-model="form.config_content" :rows="8" placeholder='{"key": "value"}' />
        </a-form-item>
        <a-form-item label="版本">
          <a-input v-model="form.version" placeholder="如 1.0.0" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">{{ isEdit ? '保存' : '创建' }}</a-button>
            <a-button @click="showCreateDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 配置详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="配置详情"
      :width="520"
    >
      <template v-if="currentConfig">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="配置编码">{{ currentConfig.config_code }}</a-descriptions-item>
          <a-descriptions-item label="配置名称">{{ currentConfig.config_name }}</a-descriptions-item>
          <a-descriptions-item label="策略类型">
            <a-tag :color="getPolicyTypeColor(currentConfig.policy_type)">{{ getPolicyTypeText(currentConfig.policy_type) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentConfig.status)">{{ getStatusText(currentConfig.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="版本">{{ currentConfig.version }}</a-descriptions-item>
          <a-descriptions-item label="关联设备数">{{ currentConfig.device_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="配置描述">{{ currentConfig.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="配置内容">
            <pre style="max-height: 200px; overflow: auto;">{{ currentConfig.config_content || '{}' }}</pre>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentConfig.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="更新时间">{{ formatTime(currentConfig.updated_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const configList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const currentConfig = ref(null)

const filters = reactive({
  policyType: undefined,
  status: undefined,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const stats = reactive({
  total: 0,
  active: 0,
  draft: 0,
  deprecated: 0
})

const form = reactive({
  config_code: '',
  config_name: '',
  policy_type: '',
  description: '',
  config_content: '{}',
  version: '1.0.0',
  status: 'draft',
  created_by: 'admin'
})

const columns = [
  { title: '配置名称', dataIndex: 'config_name', width: 180 },
  { title: '编码', dataIndex: 'config_code', width: 100 },
  { title: '策略类型', slotName: 'policyType', width: 100 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '关联设备', slotName: 'deviceCount', width: 100 },
  { title: '更新时间', slotName: 'updatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const loadConfigs = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.policyType) params.policy_type = filters.policyType
    if (filters.status) params.status = filters.status
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/policy-configs', { params })
    const data = res.data
    if (data.code === 0) {
      configList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载配置文件列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = configList.value.length
  stats.active = configList.value.filter(c => c.status === 'active').length
  stats.draft = configList.value.filter(c => c.status === 'draft').length
  stats.deprecated = configList.value.filter(c => c.status === 'deprecated').length
}

const handlePageChange = (page) => {
  pagination.current = page
  loadConfigs()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadConfigs()
}

const handleSave = (done) => { done(true) }

const handleSaveSubmit = async () => {
  try {
    const url = isEdit.value ? `/api/v1/policy-configs/${form.id}` : '/api/v1/policy-configs'
    const method = isEdit.value ? 'put' : 'post'
    const res = await axios[method](url, form)
    if (res.data.code === 0) {
      Message.success(isEdit.value ? '保存成功' : '创建成功')
      showCreateDrawer.value = false
      resetForm()
      loadConfigs()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const resetForm = () => {
  form.config_code = ''
  form.config_name = ''
  form.policy_type = ''
  form.description = ''
  form.config_content = '{}'
  form.version = '1.0.0'
  form.status = 'draft'
  form.id = null
}

const openDetail = (record) => {
  currentConfig.value = record
  showDetailDrawer.value = true
}

const editConfig = (record) => {
  isEdit.value = true
  Object.assign(form, record)
  showCreateDrawer.value = true
}

const activateConfig = async (record) => {
  try {
    const res = await axios.post(`/api/v1/policy-configs/${record.id}/activate`)
    if (res.data.code === 0) {
      Message.success('激活成功')
      loadConfigs()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const deprecateConfig = async (record) => {
  try {
    const res = await axios.post(`/api/v1/policy-configs/${record.id}/deprecate`)
    if (res.data.code === 0) {
      Message.success('已废弃')
      loadConfigs()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const getPolicyTypeColor = (type) => {
  const map = { security: 'red', network: 'blue', app: 'green', device: 'purple' }
  return map[type] || 'gray'
}

const getPolicyTypeText = (type) => {
  const map = { security: '安全策略', network: '网络策略', app: '应用策略', device: '设备策略' }
  return map[type] || type
}

const getStatusColor = (status) => {
  const map = { active: 'green', draft: 'yellow', deprecated: 'red' }
  return map[status] || 'gray'
}

const getStatusText = (status) => {
  const map = { active: '已激活', draft: '草稿', deprecated: '已废弃' }
  return map[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.policy-configs { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
