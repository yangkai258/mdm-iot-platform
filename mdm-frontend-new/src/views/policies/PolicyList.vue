<template>
  <div class="page-container">
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
            <a-select v-model="filters.category" placeholder="策略分类" allow-clear style="width: 140px" @change="loadPolicies">
              <a-option value="security">安全策略</a-option>
              <a-option value="network">网络策略</a-option>
              <a-option value="compliance">合规策略</a-option>
              <a-option value="content">内容策略</a-option>
            </a-select>
            <a-select v-model="filters.enabled" placeholder="状态" allow-clear style="width: 120px" @change="loadPolicies">
              <a-option value="true">已启用</a-option>
              <a-option value="false">已禁用</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索策略名称/编码" style="width: 200px" search-button @search="loadPolicies" />
            <a-button type="primary" @click="showCreateDrawer = true">创建策略</a-button>
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
          >
            <template #category="{ record }">
              <a-tag :color="getCategoryColor(record.category)">{{ getCategoryText(record.category) }}</a-tag>
            </template>
            <template #enabled="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'red'">{{ record.enabled ? '已启用' : '已禁用' }}</a-tag>
            </template>
            <template #priority="{ record }">
              <a-tag :color="getPriorityColor(record.priority)">{{ getPriorityText(record.priority) }}</a-tag>
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
                <a-button type="text" size="small" @click="editPolicy(record)">编辑</a-button>
                <a-button v-if="!record.enabled" type="text" size="small" @click="togglePolicy(record, true)">启用</a-button>
                <a-button v-else type="text" size="small" status="warning" @click="togglePolicy(record, false)">禁用</a-button>
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
        <a-form-item label="策略编码" required>
          <a-input v-model="form.policy_code" placeholder="如 SEC001" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="策略名称" required>
          <a-input v-model="form.policy_name" placeholder="如 密码复杂度策略" />
        </a-form-item>
        <a-form-item label="策略分类" required>
          <a-select v-model="form.category" placeholder="选择策略分类">
            <a-option value="security">安全策略</a-option>
            <a-option value="network">网络策略</a-option>
            <a-option value="compliance">合规策略</a-option>
            <a-option value="content">内容策略</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级">
          <a-select v-model="form.priority" placeholder="选择优先级">
            <a-option value="critical">Critical</a-option>
            <a-option value="high">High</a-option>
            <a-option value="medium">Medium</a-option>
            <a-option value="low">Low</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="策略描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="策略功能描述" />
        </a-form-item>
        <a-form-item label="策略规则 (JSON)">
          <a-textarea v-model="form.rules" :rows="8" placeholder='{"min_length": 8, "require_special": true}' />
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
          <a-descriptions-item label="策略编码">{{ currentPolicy.policy_code }}</a-descriptions-item>
          <a-descriptions-item label="策略名称">{{ currentPolicy.policy_name }}</a-descriptions-item>
          <a-descriptions-item label="策略分类">
            <a-tag :color="getCategoryColor(currentPolicy.category)">{{ getCategoryText(currentPolicy.category) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="优先级">
            <a-tag :color="getPriorityColor(currentPolicy.priority)">{{ getPriorityText(currentPolicy.priority) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="currentPolicy.enabled ? 'green' : 'red'">{{ currentPolicy.enabled ? '已启用' : '已禁用' }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="覆盖设备数">{{ currentPolicy.device_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="策略描述">{{ currentPolicy.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="策略规则">
            <pre style="max-height: 200px; overflow: auto;">{{ currentPolicy.rules || '{}' }}</pre>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentPolicy.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="更新时间">{{ formatTime(currentPolicy.updated_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const policyList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const currentPolicy = ref(null)

const filters = reactive({
  category: undefined,
  enabled: undefined,
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
  policy_code: '',
  policy_name: '',
  category: '',
  priority: 'medium',
  description: '',
  rules: '{}',
  enabled: true,
  created_by: 'admin'
})

const columns = [
  { title: '策略名称', dataIndex: 'policy_name', width: 200 },
  { title: '编码', dataIndex: 'policy_code', width: 100 },
  { title: '分类', slotName: 'category', width: 100 },
  { title: '优先级', slotName: 'priority', width: 100 },
  { title: '状态', slotName: 'enabled', width: 100 },
  { title: '覆盖设备', slotName: 'deviceCount', width: 100 },
  { title: '更新时间', slotName: 'updatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' }
]

const loadPolicies = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.category) params.category = filters.category
    if (filters.enabled !== undefined) params.enabled = filters.enabled
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/policies', { params })
    const data = res.data
    if (data.code === 0) {
      policyList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
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
  stats.enabled = policyList.value.filter(p => p.enabled).length
  stats.disabled = policyList.value.filter(p => !p.enabled).length
  stats.coveredDevices = policyList.value.reduce((sum, p) => sum + (p.device_count || 0), 0)
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
  form.policy_code = ''
  form.policy_name = ''
  form.category = ''
  form.priority = 'medium'
  form.description = ''
  form.rules = '{}'
  form.enabled = true
  form.id = null
}

const openDetail = (record) => {
  currentPolicy.value = record
  showDetailDrawer.value = true
}

const editPolicy = (record) => {
  isEdit.value = true
  Object.assign(form, record)
  showCreateDrawer.value = true
}

const togglePolicy = async (record, enabled) => {
  try {
    const res = await axios.post(`/api/v1/policies/${record.id}/toggle`, { enabled })
    if (res.data.code === 0) {
      Message.success(enabled ? '已启用' : '已禁用')
      loadPolicies()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const deletePolicy = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除策略「${record.policy_name}」吗？此操作不可恢复。`,
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

const getCategoryColor = (category) => {
  const map = { security: 'red', network: 'blue', compliance: 'orange', content: 'green' }
  return map[category] || 'gray'
}

const getCategoryText = (category) => {
  const map = { security: '安全策略', network: '网络策略', compliance: '合规策略', content: '内容策略' }
  return map[category] || category
}

const getPriorityColor = (priority) => {
  const map = { critical: 'red', high: 'orange', medium: 'blue', low: 'gray' }
  return map[priority] || 'gray'
}

const getPriorityText = (priority) => {
  const map = { critical: 'Critical', high: 'High', medium: 'Medium', low: 'Low' }
  return map[priority] || priority
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadPolicies()
})
</script>

<style scoped>
.policy-list { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
