<template>
  <a-layout class="compliance-rules">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="policies">
          <span>策略管理</span>
        </a-menu-item>
        <a-menu-item key="alert">
          <span>告警中心</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <a-button type="text" @click="collapsed = !collapsed">
            <span v-if="collapsed">☰</span>
            <span v-else>✕</span>
          </a-button>
        </div>
        <div class="header-title">
          <a-breadcrumb>
            <a-breadcrumb-item>策略管理</a-breadcrumb-item>
            <a-breadcrumb-item>合规规则</a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <div class="header-right"></div>
      </a-layout-header>

      <a-layout-content class="content">
        <!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="规则总数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="合规中" :value="stats.compliant" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="违规中" :value="stats.violated" :value-style="{ color: '#ff4d4f' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="待修复" :value="stats.pending" :value-style="{ color: '#faad14' }" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.severity" placeholder="严重级别" allow-clear style="width: 120px" @change="loadRules">
              <a-option value="critical">Critical</a-option>
              <a-option value="high">High</a-option>
              <a-option value="medium">Medium</a-option>
              <a-option value="low">Low</a-option>
            </a-select>
            <a-select v-model="filters.status" placeholder="合规状态" allow-clear style="width: 120px" @change="loadRules">
              <a-option value="compliant">合规中</a-option>
              <a-option value="violated">违规中</a-option>
              <a-option value="pending">待修复</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索规则名称/编码" style="width: 200px" search-button @search="loadRules" />
            <a-button type="primary" @click="showCreateDrawer = true">创建规则</a-button>
            <a-button @click="loadRules">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 规则列表 -->
        <a-card class="rules-card">
          <a-table
            :columns="columns"
            :data="ruleList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #severity="{ record }">
              <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #violationCount="{ record }">
              <a-tag v-if="record.violation_count > 0" color="red">{{ record.violation_count }}</a-tag>
              <span v-else>0</span>
            </template>
            <template #updatedAt="{ record }">
              {{ formatTime(record.updated_at) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="editRule(record)">编辑</a-button>
                <a-button type="text" size="small" @click="viewViolations(record)">违规记录</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
      </a-layout-content>
    </a-layout>

    <!-- 创建/编辑规则抽屉 -->
    <a-drawer
      v-model:visible="showCreateDrawer"
      :title="isEdit ? '编辑规则' : '创建规则'"
      :width="560"
      @before-ok="handleSave"
    >
      <a-form :model="form" layout="vertical" @submit-success="handleSaveSubmit">
        <a-form-item label="规则编码" required>
          <a-input v-model="form.rule_code" placeholder="如 COMP001" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="规则名称" required>
          <a-input v-model="form.rule_name" placeholder="如 密码复杂度检测" />
        </a-form-item>
        <a-form-item label="严重级别" required>
          <a-select v-model="form.severity" placeholder="选择严重级别">
            <a-option value="critical">Critical</a-option>
            <a-option value="high">High</a-option>
            <a-option value="medium">Medium</a-option>
            <a-option value="low">Low</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="检查类型">
          <a-select v-model="form.check_type" placeholder="选择检查类型">
            <a-option value="password">密码策略</a-option>
            <a-option value="encryption">加密检测</a-option>
            <a-option value="jailbreak">越狱检测</a-option>
            <a-option value="os_version">系统版本</a-option>
            <a-option value="app_installed">应用安装</a-option>
            <a-option value="network">网络配置</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="规则描述">
          <a-textarea v-model="form.description" :rows="3" placeholder="规则功能描述" />
        </a-form-item>
        <a-form-item label="检查条件 (JSON)">
          <a-textarea v-model="form.conditions" :rows="6" placeholder='{"min_password_length": 8, "require_encryption": true}' />
        </a-form-item>
        <a-form-item label="自动修复">
          <a-switch v-model="form.auto_fix" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit">{{ isEdit ? '保存' : '创建' }}</a-button>
            <a-button @click="showCreateDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 规则详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="规则详情"
      :width="560"
    >
      <template v-if="currentRule">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="规则编码">{{ currentRule.rule_code }}</a-descriptions-item>
          <a-descriptions-item label="规则名称">{{ currentRule.rule_name }}</a-descriptions-item>
          <a-descriptions-item label="严重级别">
            <a-tag :color="getSeverityColor(currentRule.severity)">{{ getSeverityText(currentRule.severity) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="检查类型">{{ currentRule.check_type }}</a-descriptions-item>
          <a-descriptions-item label="合规状态">
            <a-tag :color="getStatusColor(currentRule.status)">{{ getStatusText(currentRule.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="违规次数">{{ currentRule.violation_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="自动修复">{{ currentRule.auto_fix ? '是' : '否' }}</a-descriptions-item>
          <a-descriptions-item label="规则描述">{{ currentRule.description || '-' }}</a-descriptions-item>
          <a-descriptions-item label="检查条件">
            <pre style="max-height: 200px; overflow: auto;">{{ currentRule.conditions || '{}' }}</pre>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentRule.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="更新时间">{{ formatTime(currentRule.updated_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const collapsed = ref(false)
const selectedKeys = ref(['policies'])
const loading = ref(false)
const ruleList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const currentRule = ref(null)

const filters = reactive({
  severity: undefined,
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
  compliant: 0,
  violated: 0,
  pending: 0
})

const form = reactive({
  rule_code: '',
  rule_name: '',
  severity: 'medium',
  check_type: '',
  description: '',
  conditions: '{}',
  auto_fix: false,
  created_by: 'admin'
})

const columns = [
  { title: '规则名称', dataIndex: 'rule_name', width: 200 },
  { title: '编码', dataIndex: 'rule_code', width: 100 },
  { title: '严重级别', slotName: 'severity', width: 100 },
  { title: '检查类型', dataIndex: 'check_type', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '违规次数', slotName: 'violationCount', width: 100 },
  { title: '更新时间', slotName: 'updatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const handleMenuClick = ({ key }) => {
  const routes = { dashboard: '/dashboard', policies: '/policies/list', alert: '/alerts/list' }
  if (routes[key]) router.push(routes[key])
  selectedKeys.value = [key]
}

const loadRules = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.severity) params.severity = filters.severity
    if (filters.status) params.status = filters.status
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/compliance-rules', { params })
    const data = res.data
    if (data.code === 0) {
      ruleList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载合规规则列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = ruleList.value.length
  stats.compliant = ruleList.value.filter(r => r.status === 'compliant').length
  stats.violated = ruleList.value.filter(r => r.status === 'violated').length
  stats.pending = ruleList.value.filter(r => r.status === 'pending').length
}

const handlePageChange = (page) => {
  pagination.current = page
  loadRules()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadRules()
}

const handleSave = (done) => { done(true) }

const handleSaveSubmit = async () => {
  try {
    const url = isEdit.value ? `/api/v1/compliance-rules/${form.id}` : '/api/v1/compliance-rules'
    const method = isEdit.value ? 'put' : 'post'
    const res = await axios[method](url, form)
    if (res.data.code === 0) {
      Message.success(isEdit.value ? '保存成功' : '创建成功')
      showCreateDrawer.value = false
      resetForm()
      loadRules()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const resetForm = () => {
  form.rule_code = ''
  form.rule_name = ''
  form.severity = 'medium'
  form.check_type = ''
  form.description = ''
  form.conditions = '{}'
  form.auto_fix = false
  form.id = null
}

const openDetail = (record) => {
  currentRule.value = record
  showDetailDrawer.value = true
}

const editRule = (record) => {
  isEdit.value = true
  Object.assign(form, record)
  showCreateDrawer.value = true
}

const viewViolations = (record) => {
  router.push({ path: '/compliance/violations', query: { ruleId: record.id } })
}

const getSeverityColor = (severity) => {
  const map = { critical: 'red', high: 'orange', medium: 'blue', low: 'gray' }
  return map[severity] || 'gray'
}

const getSeverityText = (severity) => {
  const map = { critical: 'Critical', high: 'High', medium: 'Medium', low: 'Low' }
  return map[severity] || severity
}

const getStatusColor = (status) => {
  const map = { compliant: 'green', violated: 'red', pending: 'yellow' }
  return map[status] || 'gray'
}

const getStatusText = (status) => {
  const map = { compliant: '合规中', violated: '违规中', pending: '待修复' }
  return map[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadRules()
})
</script>

<style scoped>
.compliance-rules { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
