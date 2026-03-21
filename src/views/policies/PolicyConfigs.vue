<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>合规策略</a-breadcrumb-item>
      <a-breadcrumb-item>配置文件库</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 配置类型切换 -->
    <a-tabs v-model:active-key="activeConfigType" class="config-tabs" @change="handleTypeChange">
      <a-tab-pane key="wifi" title="Wi-Fi 配置">
        <template #icon><icon-wifi /></template>
      </a-tab-pane>
      <a-tab-pane key="vpn" title="VPN 配置" />
      <a-tab-pane key="certificate" title="证书配置" />
      <a-tab-pane key="restrictions" title="限制策略" />
      <a-tab-pane key="email" title="邮件配置" />
    </a-tabs>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card>
          <a-statistic title="配置总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已启用" :value="stats.enabled" :value-style="{ color: '#52c41a' }" />
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
        <a-select v-model="filters.status" placeholder="配置状态" allow-clear style="width: 120px" @change="loadConfigs">
          <a-option value="active">已启用</a-option>
          <a-option value="draft">草稿</a-option>
          <a-option value="deprecated">已废弃</a-option>
        </a-select>
        <a-input-search v-model="filters.keyword" placeholder="搜索配置名称" style="width: 200px" search-button @search="loadConfigs" />
        <a-button type="primary" @click="openCreateDrawer">新建配置</a-button>
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
        :scroll="{ x: 1000 }"
      >
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.enabled)">{{ record.enabled ? '已启用' : '已禁用' }}</a-tag>
        </template>
        <template #subType="{ record }">
          <a-tag>{{ record.sub_type || '-' }}</a-tag>
        </template>
        <template #version="{ record }">
          v{{ record.version || 1 }}
        </template>
        <template #updatedAt="{ record }">
          {{ formatTime(record.updated_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="editConfig(record)">编辑</a-button>
            <a-button v-if="!record.enabled" type="text" size="small" @click="toggleConfig(record, true)">启用</a-button>
            <a-button v-else type="text" size="small" status="warning" @click="toggleConfig(record, false)">禁用</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteConfig(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑配置抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      :title="isEdit ? '编辑配置' : '新建配置'"
      :width="560"
      @before-ok="handleSave"
    >
      <a-form :model="form" layout="vertical" @submit-success="handleSaveSubmit">
        <a-form-item label="配置名称" required>
          <a-input v-model="form.name" :placeholder="getNamePlaceholder()" />
        </a-form-item>
        <a-form-item v-if="activeConfigType === 'wifi'" label="Wi-Fi 类型">
          <a-select v-model="form.sub_type" placeholder="选择 Wi-Fi 类型">
            <a-option value="WPA2">WPA2</a-option>
            <a-option value="WPA3">WPA3</a-option>
            <a-option value="WPA2-Enterprise">WPA2-Enterprise</a-option>
            <a-option value="Open">Open</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="activeConfigType === 'vpn'" label="VPN 类型">
          <a-select v-model="form.sub_type" placeholder="选择 VPN 类型">
            <a-option value="OpenVPN">OpenVPN</a-option>
            <a-option value="IPSec">IPSec</a-option>
            <a-option value="L2TP">L2TP</a-option>
            <a-option value="WireGuard">WireGuard</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="activeConfigType === 'certificate'" label="证书类型">
          <a-select v-model="form.sub_type" placeholder="选择证书类型">
            <a-option value="CA">CA 证书</a-option>
            <a-option value="Client">客户端证书</a-option>
            <a-option value="Server">服务器证书</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="activeConfigType === 'restrictions'" label="限制类型">
          <a-select v-model="form.sub_type" placeholder="选择限制类型">
            <a-option value="AppInstall">应用安装限制</a-option>
            <a-option value="DeviceFeature">设备功能限制</a-option>
            <a-option value="ContentFilter">内容过滤</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="配置描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="配置功能描述" />
        </a-form-item>
        <a-form-item :label="getConfigDataLabel()">
          <a-textarea v-model="form.config_data" :rows="8" :placeholder="getConfigDataPlaceholder()" />
        </a-form-item>
        <a-form-item label="版本">
          <a-input v-model="form.version" placeholder="如 1.0.0" />
        </a-form-item>
        <a-form-item label="启用配置">
          <a-switch v-model="form.enabled" />
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
          <a-descriptions-item label="配置名称">{{ currentConfig.name }}</a-descriptions-item>
          <a-descriptions-item label="配置类型">{{ getConfigTypeText(currentConfig.config_type) }}</a-descriptions-item>
          <a-descriptions-item label="子类型">{{ currentConfig.sub_type || '-' }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentConfig.enabled)">{{ currentConfig.enabled ? '已启用' : '已禁用' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="版本">v{{ currentConfig.version || 1 }}</a-descriptions-item>
          <a-descriptions-item label="配置描述">{{ currentConfig.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="配置内容">
            <pre style="max-height: 300px; overflow: auto; background: #f5f5f5; padding: 8px; border-radius: 4px;">{{ formatJson(currentConfig.config_data) }}</pre>
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
import { Message, Modal } from '@arco-design/web-vue'

const activeConfigType = ref('wifi')
const loading = ref(false)
const configList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const currentConfig = ref(null)

const filters = reactive({
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
  enabled: 0,
  draft: 0,
  deprecated: 0
})

const form = reactive({
  name: '',
  config_type: 'wifi',
  sub_type: '',
  description: '',
  config_data: '{}',
  version: '1.0.0',
  enabled: true
})

const columns = [
  { title: '配置名称', dataIndex: 'name', width: 180 },
  { title: '子类型', slotName: 'subType', width: 140 },
  { title: '版本', slotName: 'version', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '描述', dataIndex: 'description', width: 200, ellipsis: true },
  { title: '更新时间', slotName: 'updatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' }
]

const handleTypeChange = (key) => {
  activeConfigType.value = key
  form.config_type = key
  form.sub_type = ''
  pagination.current = 1
  loadConfigs()
}

const loadConfigs = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    params.config_type = activeConfigType.value
    if (filters.status) params.enabled = filters.status === 'active'
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/policy-configs', { params })
    const data = res.data
    if (data.code === 0) {
      configList.value = data.data?.list || []
      pagination.total = data.data?.total || 0
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
  stats.enabled = configList.value.filter(c => c.enabled).length
  stats.draft = configList.value.filter(c => !c.enabled).length
  stats.deprecated = 0
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

const openCreateDrawer = () => {
  isEdit.value = false
  resetForm()
  showCreateDrawer.value = true
}

const handleSave = (done) => { done(true) }

const handleSaveSubmit = async () => {
  try {
    // Validate JSON
    if (form.config_data) {
      try {
        JSON.parse(form.config_data)
      } catch (e) {
        Message.error('配置内容必须是有效的 JSON 格式')
        return
      }
    }

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
  form.name = ''
  form.config_type = activeConfigType.value
  form.sub_type = ''
  form.description = ''
  form.config_data = '{}'
  form.version = '1.0.0'
  form.enabled = true
  form.id = null
}

const openDetail = (record) => {
  currentConfig.value = record
  showDetailDrawer.value = true
}

const editConfig = (record) => {
  isEdit.value = true
  Object.assign(form, {
    id: record.id,
    name: record.name,
    config_type: record.config_type,
    sub_type: record.sub_type || '',
    description: record.description || '',
    config_data: record.config_data || '{}',
    version: record.version ? String(record.version) : '1.0.0',
    enabled: record.enabled
  })
  activeConfigType.value = record.config_type
  showCreateDrawer.value = true
}

const toggleConfig = async (record, enabled) => {
  try {
    const res = await axios.put(`/api/v1/policy-configs/${record.id}`, { enabled })
    if (res.data.code === 0) {
      Message.success(enabled ? '已启用' : '已禁用')
      loadConfigs()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const deleteConfig = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除配置「${record.name}」吗？此操作不可恢复。`,
    okText: '删除',
    onOk: async () => {
      try {
        const res = await axios.delete(`/api/v1/policy-configs/${record.id}`)
        if (res.data.code === 0) {
          Message.success('删除成功')
          loadConfigs()
        } else {
          Message.error(res.data.message || '删除失败')
        }
      } catch (err) {
        Message.error('删除失败')
      }
    }
  })
}

const getNamePlaceholder = () => {
  const map = {
    wifi: '如 Corporate Wi-Fi',
    vpn: '如 企业VPN',
    certificate: '如 客户端证书',
    restrictions: '如 应用安装限制',
    email: '如 企业邮箱配置'
  }
  return map[activeConfigType.value] || '请输入配置名称'
}

const getConfigDataLabel = () => {
  const map = {
    wifi: 'Wi-Fi 配置 (JSON)',
    vpn: 'VPN 配置 (JSON)',
    certificate: '证书配置 (JSON)',
    restrictions: '限制策略 (JSON)',
    email: '邮件配置 (JSON)'
  }
  return map[activeConfigType.value] || '配置内容 (JSON)'
}

const getConfigDataPlaceholder = () => {
  const map = {
    wifi: '{"ssid": "Corporate", "security": "WPA2", "password": "***"}',
    vpn: '{"server": "vpn.example.com", "protocol": "OpenVPN", "port": 1194}',
    certificate: '{"type": "Client", "common_name": "user@example.com", "validity_days": 365}',
    restrictions: '{"forbidden_apps": ["game"], "max_screen_time": 3600}',
    email: '{"server": "smtp.example.com", "port": 587, "use_ssl": true}'
  }
  return map[activeConfigType.value] || '{"key": "value"}'
}

const getConfigTypeText = (type) => {
  const map = {
    wifi: 'Wi-Fi 配置',
    vpn: 'VPN 配置',
    certificate: '证书配置',
    restrictions: '限制策略',
    email: '邮件配置'
  }
  return map[type] || type
}

const getStatusColor = (enabled) => {
  return enabled ? 'green' : 'gray'
}

const formatJson = (jsonStr) => {
  if (!jsonStr) return '{}'
  try {
    return JSON.stringify(JSON.parse(jsonStr), null, 2)
  } catch (e) {
    return jsonStr
  }
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
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 16px;
}

.config-tabs {
  margin-bottom: 16px;
  background: #fff;
  border-radius: 4px;
  padding: 0 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.action-card {
  margin-bottom: 16px;
}

.config-card {
  background: #fff;
}
</style>
