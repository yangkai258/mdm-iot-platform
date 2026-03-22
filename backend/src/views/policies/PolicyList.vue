<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>合规策略</a-breadcrumb-item>
      <a-breadcrumb-item>策略管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card>
          <a-statistic title="策略总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已启用" :value="stats.enabled" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已禁用" :value="stats.disabled" :value-style="{ color: '#ff4d4f' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="覆盖设备" :value="stats.coveredDevices" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作栏 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.policyType" placeholder="策略类型" allow-clear style="width: 140px" @change="loadPolicies">
          <a-option value="compliance">合规策略</a-option>
          <a-option value="security">安全策略</a-option>
          <a-option value="network">网络策略</a-option>
          <a-option value="app">应用策略</a-option>
          <a-option value="device">设备策略</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadPolicies">
          <a-option value="active">已启用</a-option>
          <a-option value="draft">草稿</a-option>
          <a-option value="archived">已归档</a-option>
        </a-select>
        <a-input-search v-model="filters.keyword" placeholder="搜索策略名称/编码" style="width: 200px" search-button @search="loadPolicies" />
        <a-button type="primary" @click="openCreateDrawer">创建策略</a-button>
        <a-button @click="loadPolicies">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 策略列表 -->
    <a-card class="policy-card">
      <a-table
        :columns="columns"
        :data="policyList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        :scroll="{ x: 1400 }"
      >
        <template #policyType="{ record }">
          <a-tag :color="getPolicyTypeColor(record.policy_type)">{{ getPolicyTypeText(record.policy_type) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : record.status === 'draft' ? 'yellow' : 'gray'">
            {{ record.status === 'active' ? '已启用' : record.status === 'draft' ? '草稿' : '已归档' }}
          </a-tag>
        </template>
        <template #priority="{ record }">
          <a-tag :color="getPriorityColor(record.priority)">{{ record.priority }}</a-tag>
        </template>
        <template #deviceCount="{ record }">
          <a-button type="text" size="small" @click="openBindingDrawer(record)">
            {{ record.binding_count || 0 }}
          </a-button>
        </template>
        <template #updatedAt="{ record }">
          {{ formatTime(record.updated_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="openBindingDrawer(record)">绑定</a-button>
            <a-button type="text" size="small" @click="openVersionHistory(record)">版本</a-button>
            <a-button type="text" size="small" @click="editPolicy(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="deletePolicy(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑策略抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      :title="isEdit ? '编辑策略' : '创建策略'"
      :width="560"
      @before-ok="handleSave"
    >
      <a-form :model="form" layout="vertical" @submit-success="handleSaveSubmit">
        <a-form-item label="策略名称" required>
          <a-input v-model="form.name" placeholder="如 基础安全策略" />
        </a-form-item>
        <a-form-item label="策略类型" required>
          <a-select v-model="form.policy_type" placeholder="选择策略类型">
            <a-option value="compliance">合规策略</a-option>
            <a-option value="security">安全策略</a-option>
            <a-option value="network">网络策略</a-option>
            <a-option value="app">应用策略</a-option>
            <a-option value="device">设备策略</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级">
          <a-input-number v-model="form.priority" :min="0" :max="100" placeholder="数字越大优先级越高" />
        </a-form-item>
        <a-form-item label="适用范围">
          <a-select v-model="form.scope" placeholder="选择适用范围">
            <a-option value="all">全部</a-option>
            <a-option value="group">分组</a-option>
            <a-option value="individual">单个设备</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="适用平台">
          <a-select v-model="form.platform" placeholder="选择适用平台">
            <a-option value="all">全部平台</a-option>
            <a-option value="ios">iOS</a-option>
            <a-option value="android">Android</a-option>
            <a-option value="windows">Windows</a-option>
            <a-option value="mac">macOS</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="策略描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="策略功能描述" />
        </a-form-item>
        <a-form-item label="关联配置文件">
          <a-select
            v-model="form.configIds"
            placeholder="选择关联的配置文件"
            multiple
            allow-search
            :options="availableConfigs.map(c => ({ label: c.name, value: c.id }))"
          />
        </a-form-item>
        <a-form-item label="关联合规规则">
          <a-select
            v-model="form.ruleIds"
            placeholder="选择关联的合规规则"
            multiple
            allow-search
            :options="availableRules.map(r => ({ label: r.name, value: r.id }))"
          />
        </a-form-item>
        <a-form-item label="启用策略">
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

    <!-- 策略详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="策略详情"
      :width="560"
    >
      <template v-if="currentPolicy">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="策略名称">{{ currentPolicy.name }}</a-descriptions-item>
          <a-descriptions-item label="策略类型">
            <a-tag :color="getPolicyTypeColor(currentPolicy.policy_type)">{{ getPolicyTypeText(currentPolicy.policy_type) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="currentPolicy.status === 'active' ? 'green' : 'gray'">
              {{ currentPolicy.status === 'active' ? '已启用' : '已禁用' }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="优先级">{{ currentPolicy.priority }}</a-descriptions-item>
          <a-descriptions-item label="适用范围">{{ getScopeText(currentPolicy.scope) }}</a-descriptions-item>
          <a-descriptions-item label="适用平台">{{ getPlatformText(currentPolicy.platform) }}</a-descriptions-item>
          <a-descriptions-item label="关联配置数">{{ currentPolicy.config_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="关联规则数">{{ currentPolicy.rule_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="绑定设备数">{{ currentPolicy.binding_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="策略描述">{{ currentPolicy.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentPolicy.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="更新时间">{{ formatTime(currentPolicy.updated_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>

    <!-- 策略绑定抽屉 -->
    <a-drawer
      v-model:visible="showBindingDrawer"
      title="策略绑定"
      :width="600"
    >
      <template v-if="currentPolicy">
        <a-tabs>
          <a-tab-pane key="devices" title="绑定设备">
            <a-space direction="vertical" style="width: 100%;">
              <a-select
                v-model="bindingDevices"
                placeholder="选择要绑定的设备"
                multiple
                allow-search
                :max-tag-count="5"
                :options="availableDevices.map(d => ({ label: `${d.device_id} (${d.device_name || d.sn_code || d.device_id})`, value: d.device_id }))"
                style="width: 100%;"
              />
              <a-button type="primary" @click="handleBindDevices" :loading="bindingLoading">绑定选中设备</a-button>
            </a-space>
            <a-divider>已绑定设备</a-divider>
            <a-table
              :columns="bindingColumns"
              :data="currentBindings.filter(b => b.target_type === 'device')"
              size="small"
              :pagination="{ pageSize: 5 }"
            >
              <template #actions="{ record }">
                <a-button type="text" size="small" status="danger" @click="handleUnbind(record)">解绑</a-button>
              </template>
            </a-table>
          </a-tab-pane>
          <a-tab-pane key="groups" title="绑定分组">
            <a-space direction="vertical" style="width: 100%;">
              <a-select
                v-model="bindingGroups"
                placeholder="选择要绑定的分组"
                multiple
                allow-search
                :max-tag-count="5"
                :options="availableGroups.map(g => ({ label: g.name, value: String(g.id) }))"
                style="width: 100%;"
              />
              <a-button type="primary" @click="handleBindGroups" :loading="bindingLoading">绑定选中分组</a-button>
            </a-space>
            <a-divider>已绑定分组</a-divider>
            <a-table
              :columns="bindingColumns"
              :data="currentBindings.filter(b => b.target_type === 'group')"
              size="small"
              :pagination="{ pageSize: 5 }"
            >
              <template #actions="{ record }">
                <a-button type="text" size="small" status="danger" @click="handleUnbind(record)">解绑</a-button>
              </template>
            </a-table>
          </a-tab-pane>
        </a-tabs>
      </template>
    </a-drawer>

    <!-- 策略版本历史抽屉 -->
    <a-drawer
      v-model:visible="showVersionDrawer"
      title="版本历史"
      :width="600"
    >
      <template v-if="currentPolicy">
        <a-timeline>
          <a-timeline-item v-for="version in versionHistory" :key="version.id" :color="version.is_current ? 'green' : 'gray'">
            <template #label>
              <span style="color: #999;">v{{ version.version }}</span>
            </template>
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="版本">v{{ version.version }}</a-descriptions-item>
              <a-descriptions-item label="变更说明">{{ version.change_note || '-' }}</a-descriptions-item>
              <a-descriptions-item label="操作人">{{ version.updated_by || '-' }}</a-descriptions-item>
              <a-descriptions-item label="更新时间">{{ formatTime(version.updated_at) }}</a-descriptions-item>
            </a-descriptions>
            <a-tag v-if="version.is_current" color="green">当前版本</a-tag>
            <a-space v-else>
              <a-button type="text" size="small" @click="rollbackVersion(version)">回滚至此版本</a-button>
            </a-space>
          </a-timeline-item>
        </a-timeline>
        <a-empty v-if="versionHistory.length === 0" description="暂无版本历史" />
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const bindingLoading = ref(false)
const policyList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const showBindingDrawer = ref(false)
const showVersionDrawer = ref(false)
const isEdit = ref(false)
const currentPolicy = ref(null)
const currentBindings = ref([])
const versionHistory = ref([])
const availableConfigs = ref([])
const availableRules = ref([])
const availableDevices = ref([])
const availableGroups = ref([])

const bindingDevices = ref([])
const bindingGroups = ref([])

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
  enabled: 0,
  disabled: 0,
  coveredDevices: 0
})

const form = reactive({
  name: '',
  policy_type: '',
  description: '',
  priority: 0,
  configIds: [],
  ruleIds: [],
  enabled: true,
  status: 'draft',
  platform: 'all',
  scope: 'all'
})

const columns = [
  { title: '策略名称', dataIndex: 'name', width: 200 },
  { title: '类型', slotName: 'policyType', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '优先级', slotName: 'priority', width: 80 },
  { title: '绑定设备', slotName: 'deviceCount', width: 100 },
  { title: '更新时间', slotName: 'updatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' }
]

const bindingColumns = [
  { title: '目标名称', dataIndex: 'target_name', width: 200 },
  { title: '绑定时间', dataIndex: 'bound_at', width: 160 },
  { title: '操作人', dataIndex: 'bound_by', width: 120 },
  { title: '操作', slotName: 'actions', width: 80 }
]

const loadPolicies = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.policyType) params.policy_type = filters.policyType
    if (filters.status) params.status = filters.status
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/policies', { params })
    const data = res.data
    if (data.code === 0) {
      policyList.value = data.data?.list || []
      pagination.total = data.data?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载策略列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = policyList.value.length
  stats.enabled = policyList.value.filter(p => p.status === 'active').length
  stats.disabled = policyList.value.filter(p => p.status !== 'active').length
  stats.coveredDevices = policyList.value.reduce((sum, p) => sum + (p.binding_count || 0), 0)
}

const handlePageChange = (page) => {
  pagination.current = page
  loadPolicies()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadPolicies()
}

const openCreateDrawer = () => {
  isEdit.value = false
  resetForm()
  showCreateDrawer.value = true
}

const handleSave = (done) => { done(true) }

const handleSaveSubmit = async () => {
  try {
    const url = isEdit.value ? `/api/v1/policies/${form.id}` : '/api/v1/policies'
    const method = isEdit.value ? 'put' : 'post'
    const res = await axios[method](url, form)
    if (res.data.code === 0) {
      Message.success(isEdit.value ? '保存成功' : '创建成功')
      showCreateDrawer.value = false
      resetForm()
      loadPolicies()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const resetForm = () => {
  form.name = ''
  form.policy_type = ''
  form.description = ''
  form.priority = 0
  form.configIds = []
  form.ruleIds = []
  form.enabled = true
  form.status = 'draft'
  form.platform = 'all'
  form.scope = 'all'
  form.id = null
}

const openDetail = (record) => {
  currentPolicy.value = record
  showDetailDrawer.value = true
}

const editPolicy = (record) => {
  isEdit.value = true
  Object.assign(form, {
    id: record.id,
    name: record.name,
    policy_type: record.policy_type,
    description: record.description,
    priority: record.priority || 0,
    enabled: record.status === 'active',
    status: record.status,
    platform: record.platform || 'all',
    scope: record.scope || 'all',
    configIds: record.config_ids ? JSON.parse(record.config_ids) : [],
    ruleIds: record.rule_ids ? JSON.parse(record.rule_ids) : []
  })
  showCreateDrawer.value = true
}

const openBindingDrawer = async (record) => {
  currentPolicy.value = record
  bindingDevices.value = []
  bindingGroups.value = []
  // Load current bindings
  try {
    const res = await axios.get(`/api/v1/policies/${record.id}/bindings`)
    if (res.data.code === 0) {
      currentBindings.value = res.data.data || []
    }
  } catch (err) {
    console.error('加载绑定列表失败', err)
    currentBindings.value = []
  }
  showBindingDrawer.value = true
}

const handleBindDevices = async () => {
  if (bindingDevices.value.length === 0) {
    Message.warning('请选择要绑定的设备')
    return
  }
  bindingLoading.value = true
  try {
    const res = await axios.post(`/api/v1/policies/${currentPolicy.value.id}/bind`, {
      target_type: 'device',
      target_ids: bindingDevices.value
    })
    if (res.data.code === 0) {
      Message.success(`成功绑定 ${res.data.data?.created || bindingDevices.value.length} 个设备`)
      bindingDevices.value = []
      // Refresh bindings
      const bindRes = await axios.get(`/api/v1/policies/${currentPolicy.value.id}/bindings`)
      if (bindRes.data.code === 0) {
        currentBindings.value = bindRes.data.data || []
      }
      loadPolicies()
    } else {
      Message.error(res.data.message || '绑定失败')
    }
  } catch (err) {
    Message.error('绑定失败')
  } finally {
    bindingLoading.value = false
  }
}

const handleBindGroups = async () => {
  if (bindingGroups.value.length === 0) {
    Message.warning('请选择要绑定的分组')
    return
  }
  bindingLoading.value = true
  try {
    const res = await axios.post(`/api/v1/policies/${currentPolicy.value.id}/bind`, {
      target_type: 'group',
      target_ids: bindingGroups.value
    })
    if (res.data.code === 0) {
      Message.success(`成功绑定 ${res.data.data?.created || bindingGroups.value.length} 个分组`)
      bindingGroups.value = []
      // Refresh bindings
      const bindRes = await axios.get(`/api/v1/policies/${currentPolicy.value.id}/bindings`)
      if (bindRes.data.code === 0) {
        currentBindings.value = bindRes.data.data || []
      }
      loadPolicies()
    } else {
      Message.error(res.data.message || '绑定失败')
    }
  } catch (err) {
    Message.error('绑定失败')
  } finally {
    bindingLoading.value = false
  }
}

const handleUnbind = async (binding) => {
  try {
    const res = await axios.post(`/api/v1/policies/${currentPolicy.value.id}/unbind`, {
      target_type: binding.target_type,
      target_ids: [binding.target_id]
    })
    if (res.data.code === 0) {
      Message.success('解绑成功')
      // Refresh bindings
      const bindRes = await axios.get(`/api/v1/policies/${currentPolicy.value.id}/bindings`)
      if (bindRes.data.code === 0) {
        currentBindings.value = bindRes.data.data || []
      }
      loadPolicies()
    } else {
      Message.error(res.data.message || '解绑失败')
    }
  } catch (err) {
    Message.error('解绑失败')
  }
}

const openVersionHistory = async (record) => {
  currentPolicy.value = record
  // Load version history
  try {
    const res = await axios.get(`/api/v1/policies/${record.id}/versions`)
    if (res.data.code === 0) {
      versionHistory.value = res.data.data || []
    } else {
      // Fallback: show current version as history
      versionHistory.value = [{
        id: record.id,
        version: record.version || 1,
        is_current: true,
        updated_by: record.updated_by,
        updated_at: record.updated_at
      }]
    }
  } catch (err) {
    console.error('加载版本历史失败', err)
    versionHistory.value = []
  }
  showVersionDrawer.value = true
}

const rollbackVersion = async (version) => {
  try {
    const res = await axios.post(`/api/v1/policies/${currentPolicy.value.id}/rollback`, {
      version_id: version.id
    })
    if (res.data.code === 0) {
      Message.success('回滚成功')
      loadPolicies()
      openVersionHistory(currentPolicy.value)
    } else {
      Message.error(res.data.message || '回滚失败')
    }
  } catch (err) {
    Message.error('回滚失败')
  }
}

const deletePolicy = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除策略「${record.name}」吗？此操作不可恢复。`,
    okText: '删除',
    onOk: async () => {
      try {
        const res = await axios.delete(`/api/v1/policies/${record.id}`)
        if (res.data.code === 0) {
          Message.success('删除成功')
          loadPolicies()
        } else {
          Message.error(res.data.message || '删除失败')
        }
      } catch (err) {
        Message.error('删除失败')
      }
    }
  })
}

const loadAvailableConfigs = async () => {
  try {
    const res = await axios.get('/api/v1/policy-configs', { params: { page: 1, page_size: 500 } })
    if (res.data.code === 0) {
      availableConfigs.value = res.data.data?.list || []
    }
  } catch (err) {
    console.error('加载配置文件列表失败', err)
  }
}

const loadAvailableRules = async () => {
  try {
    const res = await axios.get('/api/v1/compliance-rules', { params: { page: 1, page_size: 500 } })
    if (res.data.code === 0) {
      availableRules.value = res.data.data?.list || []
    }
  } catch (err) {
    console.error('加载合规规则列表失败', err)
  }
}

const loadAvailableDevices = async () => {
  try {
    const res = await axios.get('/api/v1/devices/list', { params: { page: 1, page_size: 500 } })
    if (res.data.code === 0) {
      availableDevices.value = res.data.data?.list || []
    }
  } catch (err) {
    console.error('加载设备列表失败', err)
  }
}

const loadAvailableGroups = async () => {
  try {
    const res = await axios.get('/api/v1/groups', { params: { page: 1, page_size: 100 } })
    if (res.data.code === 0) {
      availableGroups.value = res.data.data?.list || []
    }
  } catch (err) {
    console.error('加载分组列表失败', err)
    availableGroups.value = []
  }
}

const getPolicyTypeColor = (type) => {
  const map = { compliance: 'orange', security: 'red', network: 'blue', app: 'green', device: 'purple' }
  return map[type] || 'gray'
}

const getPolicyTypeText = (type) => {
  const map = { compliance: '合规策略', security: '安全策略', network: '网络策略', app: '应用策略', device: '设备策略' }
  return map[type] || type
}

const getPriorityColor = (priority) => {
  if (priority >= 80) return 'red'
  if (priority >= 60) return 'orange'
  if (priority >= 40) return 'blue'
  return 'gray'
}

const getScopeText = (scope) => {
  const map = { all: '全部', group: '分组', individual: '单个设备' }
  return map[scope] || scope
}

const getPlatformText = (platform) => {
  const map = { all: '全部平台', ios: 'iOS', android: 'Android', windows: 'Windows', mac: 'macOS' }
  return map[platform] || platform
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadPolicies()
  loadAvailableConfigs()
  loadAvailableRules()
  loadAvailableDevices()
  loadAvailableGroups()
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

.stats-row {
  margin-bottom: 16px;
}

.action-card {
  margin-bottom: 16px;
}

.policy-card {
  background: #fff;
}
</style>
