<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>AI 行为日志</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索设备/用户"
          style="width: 200px"
          search-button
          @search="loadLogs"
        />
        <a-select v-model="filters.model_version" placeholder="模型版本" allow-clear style="width: 160px" @change="loadLogs">
          <a-option value="MiniClaw-v2.1.3">MiniClaw v2.1.3</a-option>
          <a-option value="MiniClaw-v2.1.2">MiniClaw v2.1.2</a-option>
          <a-option value="Behavior-v1.5.0">Behavior v1.5.0</a-option>
        </a-select>
        <a-select v-model="filters.behavior_type" placeholder="事件类型" allow-clear style="width: 160px" @change="loadLogs">
          <a-option value="intent_recognition">意图识别</a-option>
          <a-option value="response_generation">响应生成</a-option>
          <a-option value="action_selection">动作选择</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadLogs">
          <a-option value="success">成功</a-option>
          <a-option value="error">错误</a-option>
          <a-option value="anomaly">异常</a-option>
        </a-select>
        <a-range-picker
          v-model="filters.time_range"
          show-time
          format="YYYY-MM-DD HH:mm"
          @change="loadLogs"
          style="width: 280px"
        />
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="loadLogs">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button @click="handleExport">
          <template #icon><icon-download /></template>
          导出
        </a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="logList"
        :loading="logLoading"
        :pagination="pagination"
        row-key="id"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @row-click="handleRowClick"
        :row-class="getRowClass"
        :scroll="{ x: 1200 }"
      >
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #behavior_type="{ record }">
          <a-tag :color="getBehaviorTypeColor(record.behavior_type)">
            {{ getBehaviorTypeText(record.behavior_type) }}
          </a-tag>
        </template>
        <template #latency_ms="{ record }">
          {{ record.latency_ms }}ms
        </template>
        <template #confidence="{ record }">
          <span>{{ record.confidence ? (record.confidence * 100).toFixed(1) : '--' }}%</span>
        </template>
        <template #status="{ record }">
          <a-tag v-if="record.is_anomaly" color="red">异常</a-tag>
          <a-tag v-else-if="record.error_code" color="orange">错误</a-tag>
          <a-tag v-else color="green">正常</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click.stop="goToDetail(record.id)">详情</a-button>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { getAiMonitorEvents } from '@/api/ai'

const router = useRouter()

const logLoading = ref(false)
const logList = ref([])
const logTotal = ref(0)

const filters = reactive({
  keyword: '',
  model_version: undefined,
  behavior_type: undefined,
  status: undefined,
  time_range: []
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '时间', dataIndex: 'created_at', slotName: 'created_at', width: 170, fixed: 'left' },
  { title: '设备', dataIndex: 'device_id', ellipsis: true, width: 140 },
  { title: '用户', dataIndex: 'user_id', width: 100 },
  { title: '模型', dataIndex: 'model_version', ellipsis: true, width: 160 },
  { title: '事件类型', slotName: 'behavior_type', width: 120 },
  { title: '延迟', slotName: 'latency_ms', width: 90 },
  { title: '置信度', slotName: 'confidence', width: 90 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 80, fixed: 'right' }
]

const formatTime = (ts) => {
  if (!ts) return '--'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const getBehaviorTypeColor = (type) => ({
  intent_recognition: 'arcoblue',
  response_generation: 'green',
  action_selection: 'purple'
}[type] || 'gray')

const getBehaviorTypeText = (type) => ({
  intent_recognition: '意图识别',
  response_generation: '响应生成',
  action_selection: '动作选择'
}[type] || type)

const getRowClass = (record) => {
  if (record.is_anomaly) return 'row-anomaly'
  if (record.error_code) return 'row-error'
  return ''
}

const loadLogs = async () => {
  logLoading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.model_version) params.model_version = filters.model_version
    if (filters.behavior_type) params.behavior_type = filters.behavior_type
    if (filters.status) params.status = filters.status
    if (filters.time_range && filters.time_range.length === 2) {
      params.start_time = filters.time_range[0].toISOString()
      params.end_time = filters.time_range[1].toISOString()
    }

    const res = await getAiMonitorEvents(params)
    if (res.code === 0) {
      logList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    Message.error('加载失败')
  } finally {
    logLoading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.current = page
  loadLogs()
}

const handlePageSizeChange = (size) => {
  pagination.pageSize = size
  pagination.current = 1
  loadLogs()
}

const handleRowClick = (record) => {
  goToDetail(record.id)
}

const goToDetail = (id) => {
  router.push(`/ai/behavior-detail/${id}`)
}

const handleExport = () => {
  const data = logList.value
  if (!data.length) {
    Message.warning('暂无数据可导出')
    return
  }
  const headers = ['时间', '设备', '用户', '模型', '事件类型', '延迟(ms)', '置信度', '状态']
  const rows = data.map(r => [
    formatTime(r.created_at),
    r.device_id || '',
    r.user_id || '',
    r.model_version || '',
    getBehaviorTypeText(r.behavior_type),
    r.latency_ms || '',
    r.confidence ? (r.confidence * 100).toFixed(1) : '',
    r.is_anomaly ? '异常' : r.error_code ? '错误' : '正常'
  ])
  const csvContent = [headers, ...rows].map(r => r.map(c => `"${c}"`).join(',')).join('\n')
  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ai-behavior-log-${Date.now()}.csv`
  a.click()
  URL.revokeObjectURL(url)
  Message.success('导出成功')
}

onMounted(() => {
  loadLogs()
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

:deep(.row-anomaly) { background: #fff1f0; cursor: pointer; }
:deep(.row-error) { background: #fff7e6; cursor: pointer; }
:deep(.arco-table-tr) { cursor: pointer; }
</style>
