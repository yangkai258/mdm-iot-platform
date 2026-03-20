<template>
  <div class="page-container">
<!-- 统计卡片 -->
        <a-row :gutter="16" class="stats-row">
          <a-col :span="6">
            <a-card>
              <a-statistic title="告警总数" :value="stats.total" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="待处理" :value="stats.pending" :value-style="{ color: '#ff4d4f' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="处理中" :value="stats.processing" :value-style="{ color: '#faad14' }" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card>
              <a-statistic title="已处理" :value="stats.resolved" :value-style="{ color: '#52c41a' }" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 标签页 -->
        <a-tabs v-model:active-key="activeTab" @change="handleTabChange" class="alert-tabs">
          <a-tab-pane key="pending">
            <template #title>待处理 ({{ stats.pending }})</template>
          </a-tab-pane>
          <a-tab-pane key="processing">
            <template #title>处理中 ({{ stats.processing }})</template>
          </a-tab-pane>
          <a-tab-pane key="resolved">
            <template #title>已处理 ({{ stats.resolved }})</template>
          </a-tab-pane>
          <a-tab-pane key="all">
            <template #title>全部</template>
          </a-tab-pane>
        </a-tabs>

        <!-- 操作栏 -->
        <a-card class="action-card">
          <a-space wrap>
            <a-select v-model="filters.severity" placeholder="严重级别" allow-clear style="width: 120px" @change="loadAlerts">
              <a-option value="critical">Critical</a-option>
              <a-option value="high">High</a-option>
              <a-option value="medium">Medium</a-option>
              <a-option value="low">Low</a-option>
            </a-select>
            <a-select v-model="filters.category" placeholder="告警类别" allow-clear style="width: 140px" @change="loadAlerts">
              <a-option value="device">设备告警</a-option>
              <a-option value="system">系统告警</a-option>
              <a-option value="security">安全告警</a-option>
              <a-option value="network">网络告警</a-option>
            </a-select>
            <a-input-search v-model="filters.keyword" placeholder="搜索告警内容" style="width: 200px" search-button @search="loadAlerts" />
            <a-button v-if="activeTab === 'pending'" type="primary" @click="batchAcknowledge">批量确认</a-button>
            <a-button @click="loadAlerts">刷新</a-button>
          </a-space>
        </a-card>

        <!-- 告警列表 -->
        <a-card class="alerts-card">
          <a-table
            :columns="columns"
            :data="alertList"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            :row-selection="activeTab !== 'resolved' ? rowSelection : null"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #severity="{ record }">
              <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #createdAt="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
                <a-button v-if="record.status === 'pending'" type="text" size="small" @click="acknowledgeAlert(record)">确认</a-button>
                <a-button v-if="record.status === 'acknowledged'" type="text" size="small" @click="handleAlert(record)">处理</a-button>
                <a-button type="text" size="small" status="danger" @click="deleteAlert(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
</div>

    <!-- 告警详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="告警详情"
      :width="560"
    >
      <template v-if="currentAlert">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="告警ID">{{ currentAlert.id }}</a-descriptions-item>
          <a-descriptions-item label="告警规则">{{ currentAlert.rule_name }}</a-descriptions-item>
          <a-descriptions-item label="严重级别">
            <a-tag :color="getSeverityColor(currentAlert.severity)">{{ getSeverityText(currentAlert.severity) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="类别">{{ getCategoryText(currentAlert.category) }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentAlert.status)">{{ getStatusText(currentAlert.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="设备">{{ currentAlert.device_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="告警内容">{{ currentAlert.message }}</a-descriptions-item>
          <a-descriptions-item label="触发时间">{{ formatTime(currentAlert.created_at) }}</a-descriptions-item>
          <a-descriptions-item v-if="currentAlert.acknowledged_at" label="确认时间">{{ formatTime(currentAlert.acknowledged_at) }}</a-descriptions-item>
          <a-descriptions-item v-if="currentAlert.resolved_at" label="处理时间">{{ formatTime(currentAlert.resolved_at) }}</a-descriptions-item>
          <a-descriptions-item label="处理人">{{ currentAlert.handler || '-' }}</a-descriptions-item>
          <a-descriptions-item label="处理备注">{{ currentAlert.resolve_note || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>

    <!-- 处理告警抽屉 -->
    <a-drawer
      v-model:visible="showHandleDrawer"
      title="处理告警"
      :width="400"
    >
      <a-form :model="handleForm" layout="vertical">
        <a-form-item label="处理结果" required>
          <a-radio-group v-model="handleForm.result">
            <a-radio value="resolved">已解决</a-radio>
            <a-radio value="false_positive">误报</a-radio>
            <a-radio value="ignored">忽略</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="处理备注">
          <a-textarea v-model="handleForm.note" :rows="4" placeholder="请输入处理备注" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="submitHandle">确认处理</a-button>
            <a-button @click="showHandleDrawer = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { Message, Modal } from '@arco-design/web-vue'

const loading = ref(false)
const alertList = ref([])
const showDetailDrawer = ref(false)
const showHandleDrawer = ref(false)
const currentAlert = ref(null)
const activeTab = ref('pending')
const selectedRowKeys = ref([])

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
  pending: 0,
  processing: 0,
  resolved: 0
})

const handleForm = reactive({
  result: 'resolved',
  note: ''
})

const rowSelection = reactive({
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false
})

const columns = [
  { title: '告警内容', dataIndex: 'message', width: 280, ellipsis: true },
  { title: '规则名称', dataIndex: 'rule_name', width: 150 },
  { title: '设备', dataIndex: 'device_name', width: 120 },
  { title: '严重级别', slotName: 'severity', width: 100 },
  { title: '类别', dataIndex: 'category', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '触发时间', slotName: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

if (routes[key]) router.push(routes[key])
  selectedKeys.value = [key]
}

const loadAlerts = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.severity) params.severity = filters.severity
    if (filters.category) params.category = filters.category
    if (filters.keyword) params.keyword = filters.keyword
    if (activeTab.value !== 'all') params.status = activeTab.value

    const res = await axios.get('/api/v1/alerts', { params })
    const data = res.data
    if (data.code === 0) {
      alertList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    }
  } catch (err) {
    Message.error('加载告警列表失败')
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.total = alertList.value.length
  stats.pending = alertList.value.filter(a => a.status === 'pending').length
  stats.processing = alertList.value.filter(a => a.status === 'acknowledged').length
  stats.resolved = alertList.value.filter(a => a.status === 'resolved').length
}

const handleTabChange = () => {
  selectedRowKeys.value = []
  pagination.current = 1
  loadAlerts()
}

const handlePageChange = (page) => {
  pagination.current = page
  loadAlerts()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadAlerts()
}

const openDetail = (record) => {
  currentAlert.value = record
  showDetailDrawer.value = true
}

const acknowledgeAlert = async (record) => {
  try {
    const res = await axios.post(`/api/v1/alerts/${record.id}/acknowledge`)
    if (res.data.code === 0) {
      Message.success('已确认')
      loadAlerts()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const handleAlert = (record) => {
  currentAlert.value = record
  handleForm.result = 'resolved'
  handleForm.note = ''
  showHandleDrawer.value = true
}

const submitHandle = async () => {
  try {
    const res = await axios.post(`/api/v1/alerts/${currentAlert.value.id}/resolve`, handleForm)
    if (res.data.code === 0) {
      Message.success('处理成功')
      showHandleDrawer.value = false
      loadAlerts()
    } else {
      Message.error(res.data.message || '处理失败')
    }
  } catch (err) {
    Message.error('处理失败')
  }
}

const batchAcknowledge = () => {
  if (selectedRowKeys.value.length === 0) {
    Message.warning('请选择要确认的告警')
    return
  }
  Modal.warning({
    title: '批量确认',
    content: `确定要确认选中的 ${selectedRowKeys.value.length} 条告警吗？`,
    okText: '确认',
    onOk: async () => {
      try {
        const res = await axios.post('/api/v1/alerts/batch/acknowledge', { ids: selectedRowKeys.value })
        if (res.data.code === 0) {
          Message.success('批量确认成功')
          selectedRowKeys.value = []
          loadAlerts()
        } else {
          Message.error(res.data.message || '操作失败')
        }
      } catch (err) {
        Message.error('操作失败')
      }
    }
  })
}

const deleteAlert = (record) => {
  Modal.warning({
    title: '确认删除',
    content: '确定要删除该告警记录吗？',
    okText: '删除',
    onOk: async () => {
      try {
        const res = await axios.delete(`/api/v1/alerts/${record.id}`)
        if (res.data.code === 0) {
          Message.success('删除成功')
          loadAlerts()
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

const getStatusColor = (status) => {
  const map = { pending: 'red', acknowledged: 'yellow', resolved: 'green', false_positive: 'gray', ignored: 'gray' }
  return map[status] || 'gray'
}

const getStatusText = (status) => {
  const map = { pending: '待处理', acknowledged: '处理中', resolved: '已解决', false_positive: '误报', ignored: '忽略' }
  return map[status] || status
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
  loadAlerts()
})
</script>

<style scoped>
.alert-list { min-height: 100vh; }
.header { background: #fff; padding: 0 16px; display: flex; align-items: center; gap: 16px; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
.header-left { display: flex; align-items: center; }
.header-title { font-size: 16px; font-weight: 500; }
.content { padding: 16px; background: #f0f2f5; }
.stats-row { margin-bottom: 16px; }
.alert-tabs { margin-bottom: 16px; }
.action-card { margin-bottom: 16px; }
</style>
