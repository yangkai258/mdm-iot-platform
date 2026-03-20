<template>
  <a-layout class="alert-rules">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
        <span v-if="!collapsed">MDM 控制台</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <a-menu-item key="dashboard">
          <span>设备大盘</span>
        </a-menu-item>
        <a-menu-item key="alerts">
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
            <a-breadcrumb-item>告警中心</a-breadcrumb-item>
            <a-breadcrumb-item>告警规则</a-breadcrumb-item>
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
              <a-statistic title="已启用" :value="stats.enabled" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="今日触发" :value="stats.todayTriggered" :value-style="{ color: '#ff4d4f' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="告警总数" :value="stats.totalAlerts" />
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
            <a-select v-model="filters.category" placeholder="告警类别" allow-clear style="width: 140px" @change="loadRules">
              <a-option value="device">设备告警</a-option>
              <a-option value="system">系统告警</a-option>
              <a-option value="security">安全告警</a-option>
              <a-option value="network">网络告警</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索规则名称" style="width: 200px" search-button @search="loadRules" />
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
            <template #category="{ record }">
              <a-tag>{{ getCategoryText(record.category) }}</a-tag>
            </template>
            <template #enabled="{ record }">
              <a-switch v-model="record.enabled" @change="toggleRule(record)" />
            </template>
            <template #triggerCount="{ record }">
              <span style="color: #ff4d4f;">{{ record.trigger_count || 0 }}</span>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="editRule(record)">编辑</a-button>
                <a-button type="text" size="small" status="danger" @click="deleteRule(record)">删除</a-button>
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
          <a-input v-model="form.rule_code" placeholder="如 ALR001" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="规则名称" required>
          <a-input v-model="form.rule_name" placeholder="如 设备离线告警" />
        </a-form-item>
        <a-form-item label="告警类别" required>
          <a-select v-model="form.category" placeholder="选择告警类别">
            <a-option value="device">设备告警</a-option>
            <a-option value="system">系统告警</a-option>
            <a-option value="security">安全告警</a-option>
            <a-option value="network">网络告警</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="严重级别" required>
          <a-select v-model="form.severity" placeholder="选择严重级别">
            <a-option value="critical">Critical</a-option>
            <a-option value="high">High</a-option>
            <a-option value="medium">Medium</a-option>
            <a-option value="low">Low</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="触发条件 (JSON)" required>
          <a-textarea v-model="form.conditions" :rows="5" placeholder='{"metric": "device_status", "operator": "==", "value": "offline", "duration": 300}' />
        </a-form-item>
        <a-form-item label="告警消息模板">
          <a-textarea v-model="form.message_template" :rows="2" placeholder="设备 ${device_name} 已离线超过 ${duration} 秒" />
        </a-form-item>
        <a-form-item label="通知方式">
          <a-checkbox-group v-model="form.notify_channels">
            <a-checkbox value="email">邮件</a-checkbox>
            <a-checkbox value="sms">短信</a-checkbox>
            <a-checkbox value="webhook">Webhook</a-checkbox>
            <a-checkbox value="in_app">应用内通知</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="告警抑制 (分钟)">
          <a-input-number v-model="form.suppress_minutes" :min="0" :max="1440" placeholder="相同告警的抑制时间" style="width: 100%" />
        </a-form-item>
        <a-form-item label="启用规则">
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
          <a-descriptions-item label="告警类别">{{ getCategoryText(currentRule.category) }}</a-descriptions-item>
          <a-descriptions-item label="严重级别">
            <a-tag :color="getSeverityColor(currentRule.severity)">{{ getSeverityText(currentRule.severity) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="启用状态">{{ currentRule.enabled ? '是' : '否' }}</a-descriptions-item>
          <a-descriptions-item label="触发次数">{{ currentRule.trigger_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="通知方式">{{ currentRule.notify_channels?.join(', ') || '-' }}</a-descriptions-item>
          <a-descriptions-item label="告警抑制">{{ currentRule.suppress_minutes || 0 }} 分钟</a-descriptions-item>
          <a-descriptions-item label="触发条件">
            <pre style="max-height: 150px; overflow: auto;">{{ currentRule.conditions || '{}' }}</pre>
          </a-descriptions-item>
          <a-descriptions-item label="消息模板">{{ currentRule.message_template || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(currentRule.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </a-layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Message, Modal } from '@arco-design/web-vue'

const router = useRouter()
const collapsed = ref(false)
const selectedKeys = ref(['alerts'])
const loading = ref(false)
const ruleList = ref([])
const showCreateDrawer = ref(false)
const showDetailDrawer = ref(false)
const isEdit = ref(false)
const currentRule = ref(null)

const filters = reactive({
  severity: undefined,
  category: undefined,
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
  todayTriggered: 0,
  totalAlerts: 0
})

const form = reactive({
  rule_code: '',
  rule_name: '',
  category: 'device',
  severity: 'medium',
  conditions: '{"metric": "", "operator": "", "value": ""}',
  message_template: '',
  notify_channels: ['in_app'],
  suppress_minutes: 5,
  enabled: true,
  created_by: 'admin'
})

const columns = [
  { title: '规则名称', dataIndex: 'rule_name', width: 200 },
  { title: '编码', dataIndex: 'rule_code', width: 100 },
  { title: '类别', slotName: 'category', width: 100 },
  { title: '严重级别', slotName: 'severity', width: 100 },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '触发次数', slotName: 'triggerCount', width: 100 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

const handleMenuClick = ({ key }) => {
  const routes = { dashboard: '/dashboard', alerts: '/alerts/rules' }
  if (routes[key]) router.push(routes[key])
  selectedKeys.value = [key]
}

const loadRules = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.severity) params.severity = filters.severity
    if (filters.category) params.category = filters.category
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/alerts/rules', { params })
    const data = res.data
    if (data.code === 0) {
      ruleList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载告警规则列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = ruleList.value.length
  stats.enabled = ruleList.value.filter(r => r.enabled).length
  stats.todayTriggered = ruleList.value.reduce((sum, r) => sum + (r.today_trigger_count || 0), 0)
  stats.totalAlerts = ruleList.value.reduce((sum, r) => sum + (r.trigger_count || 0), 0)
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
    const url = isEdit.value ? `/api/v1/alerts/rules/${form.id}` : '/api/v1/alerts/rules'
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
  form.category = 'device'
  form.severity = 'medium'
  form.conditions = '{"metric": "", "operator": "", "value": ""}'
  form.message_template = ''
  form.notify_channels = ['in_app']
  form.suppress_minutes = 5
  form.enabled = true
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

const toggleRule = async (record) => {
  try {
    const res = await axios.post(`/api/v1/alerts/rules/${record.id}/toggle`, { enabled: record.enabled })
    if (res.data.code === 0) {
      Message.success(record.enabled ? '已启用' : '已禁用')
    } else {
      Message.error(res.data.message || '操作失败')
      loadRules()
    }
  } catch (err) {
    Message.error('操作失败')
    loadRules()
  }
}

const deleteRule = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除告警规则「${record.rule_name}」吗？`,
    okText: '删除',
    onOk: async () => {
      try {
        const res = await axios.delete(`/api/v1/alerts/rules/${record.id}`)
        if (res.data.code === 0) {
          Message.success('删除成功')
          loadRules()
        } else {
          Message.error(res.data.message || '删除失败')
        }
      } catch (err) {
        Message.error('删除失败')
      }
    }
  })
}

const getSeverityColor = (severity) => {
  const map = { critical: 'red', high: 'orange', medium: 'blue', low: 'gray' }
  return map[severity] || 'gray'
}

const getSeverityText = (severity) => {
  const map = { critical: 'Critical', high: 'High', medium: 'Medium', low: 'Low' }
  return map[severity] || severity
}

const getCategoryText = (category) => {
  const map = { device: '设备告警', system: '系统告警', security: '安全告警', network: '网络告警' }
  return map[category] || category
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
.alert-rules { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
