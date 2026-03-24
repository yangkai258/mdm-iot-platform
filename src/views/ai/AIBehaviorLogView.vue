<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>AI 行为日志</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="关键词">
          <a-input-search
            v-model="searchForm.keyword"
            placeholder="搜索设备/用户"
            style="width: 200px"
            search-button
            @search="loadLogs"
          />
        </a-form-item>
        <a-form-item label="模型版本">
          <a-select v-model="searchForm.model_version" placeholder="选择版本" allow-clear style="width: 160px">
            <a-option value="MiniClaw-v2.1.3">MiniClaw v2.1.3</a-option>
            <a-option value="MiniClaw-v2.1.2">MiniClaw v2.1.2</a-option>
            <a-option value="Behavior-v1.5.0">Behavior v1.5.0</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="事件类型">
          <a-select v-model="searchForm.behavior_type" placeholder="选择类型" allow-clear style="width: 160px">
            <a-option value="intent_recognition">意图识别</a-option>
            <a-option value="response_generation">响应生成</a-option>
            <a-option value="action_selection">动作选择</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option value="success">成功</a-option>
            <a-option value="error">错误</a-option>
            <a-option value="anomaly">异常</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker
            v-model="searchForm.time_range"
            show-time
            format="YYYY-MM-DD HH:mm"
            style="width: 280px"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="loadLogs">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-space>
        <a-button type="primary" @click="loadLogs">刷新</a-button>
        <a-button @click="handleExport">导出</a-button>
      </a-space>
    </div>

    <!-- 表格 -->
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
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'

const router = useRouter()

const logLoading = ref(false)
const logList = ref([])

const searchForm = reactive({
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
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.model_version) params.model_version = searchForm.model_version
    if (searchForm.behavior_type) params.behavior_type = searchForm.behavior_type
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.time_range && searchForm.time_range.length === 2) {
      params.start_time = searchForm.time_range[0].toISOString()
      params.end_time = searchForm.time_range[1].toISOString()
    }

    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/ai/events', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      logList.value = data.data?.list || []
      pagination.total = data.data?.total || 0
    } else {
      logList.value = []
    }
  } catch (e) {
    Message.error('加载失败')
  } finally {
    logLoading.value = false
  }
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.model_version = undefined
  searchForm.behavior_type = undefined
  searchForm.status = undefined
  searchForm.time_range = []
  pagination.current = 1
  loadLogs()
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
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.toolbar {
  margin-bottom: 16px;
}

:deep(.row-anomaly) {
  background: #fff1f0;
  cursor: pointer;
}

:deep(.row-error) {
  background: #fff7e6;
  cursor: pointer;
}

:deep(.arco-table-tr) {
  cursor: pointer;
}
</style>
