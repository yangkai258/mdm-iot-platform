<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>告警列表</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="告警总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="待处理" :value="stats.pending" :value-style="{ color: '#ff4d4f' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="处理中" :value="stats.processing" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="已处理" :value="stats.resolved" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索框 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input-search v-model="filters.keyword" placeholder="搜索告警内容" style="width: 280px" search-button @search="loadAlerts" />
        <a-select v-model="filters.severity" placeholder="严重级别" allow-clear style="width: 120px" @change="loadAlerts">
          <a-option value="critical">严重</a-option>
          <a-option value="high">高</a-option>
          <a-option value="medium">中</a-option>
          <a-option value="low">低</a-option>
        </a-select>
        <a-select v-model="filters.category" placeholder="告警类别" allow-clear style="width: 140px" @change="loadAlerts">
          <a-option value="device">设备告警</a-option>
          <a-option value="system">系统告警</a-option>
          <a-option value="security">安全告警</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮组 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="batchAcknowledge">批量确认</a-button>
        <a-button @click="loadAlerts">刷新</a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="alertList" :loading="loading" :pagination="pagination" row-key="id" @page-change="handlePageChange">
        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button v-if="record.status === 'pending'" type="text" size="small" @click="acknowledgeAlert(record)">确认</a-button>
            <a-button type="text" size="small" status="danger" @click="deleteAlert(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 告警详情抽屉 -->
    <a-drawer v-model:visible="showDetailDrawer" title="告警详情" :width="560">
      <template v-if="currentAlert">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="告警ID">{{ currentAlert.id }}</a-descriptions-item>
          <a-descriptions-item label="告警内容">{{ currentAlert.message }}</a-descriptions-item>
          <a-descriptions-item label="严重级别">
            <a-tag :color="getSeverityColor(currentAlert.severity)">{{ getSeverityText(currentAlert.severity) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="getStatusColor(currentAlert.status)">{{ getStatusText(currentAlert.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="触发时间">{{ formatTime(currentAlert.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const alertList = ref([])
const showDetailDrawer = ref(false)
const currentAlert = ref(null)

const filters = reactive({ severity: undefined, category: undefined, keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const stats = reactive({ total: 0, pending: 0, processing: 0, resolved: 0 })

const columns = [
  { title: '告警内容', dataIndex: 'message', width: 280, ellipsis: true },
  { title: '规则名称', dataIndex: 'rule_name', width: 150 },
  { title: '设备', dataIndex: 'device_name', width: 120 },
  { title: '严重级别', slotName: 'severity', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '触发时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getSeverityColor = (s) => ({ critical: 'red', high: 'orange', medium: 'blue', low: 'gray' }[s] || 'gray')
const getSeverityText = (s) => ({ critical: '严重', high: '高', medium: '中', low: '低' }[s] || s)
const getStatusColor = (s) => ({ pending: 'red', processing: 'yellow', resolved: 'green' }[s] || 'gray')
const getStatusText = (s) => ({ pending: '待处理', processing: '处理中', resolved: '已处理' }[s] || s)
const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadAlerts = () => { loading.value = true; setTimeout(() => { loading.value = false }, 300) }
const handlePageChange = (p) => { pagination.current = p; loadAlerts() }
const openDetail = (r) => { currentAlert.value = r; showDetailDrawer.value = true }
const acknowledgeAlert = (r) => { Message.success('已确认'); loadAlerts() }
const deleteAlert = (r) => { alertList.value = alertList.value.filter(x => x.id !== r.id); Message.success('删除成功') }
const batchAcknowledge = () => { Message.info('批量确认') }

onMounted(() => loadAlerts())
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
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
