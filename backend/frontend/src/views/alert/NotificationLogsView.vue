<template>
  <div class="notification-logs-view">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <a-form :model="filters" layout="inline" ref="filterFormRef">
        <a-form-item label="渠道类型">
          <a-select v-model="filters.channel_type" placeholder="选择渠道" style="width: 120px" allow-clear>
            <a-option value="email">邮件</a-option>
            <a-option value="sms">短信</a-option>
            <a-option value="webhook">Webhook</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发送状态">
          <a-select v-model="filters.status" placeholder="选择状态" style="width: 120px" allow-clear>
            <a-option value="pending">等待中</a-option>
            <a-option value="sent">已发送</a-option>
            <a-option value="failed">失败</a-option>
            <a-option value="retrying">重试中</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker
            v-model="dateRange"
            style="width: 280px"
            @change="onDateChange"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              「查询」
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              「重置」
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作区域 -->
    <div class="action-section">
      <a-button @click="handleRefresh">
        <template #icon><icon-refresh /></template>
        「刷新」
      </a-button>
      <span class="total-count">共 {{ pagination.total }} 条记录</span>
    </div>

    <!-- 表格区域 -->
    <div class="table-section">
      <a-table
        :data="logs"
        :loading="loading"
        :pagination="paginationConfig"
        :columns="columns"
        :scroll="{ x: 1200 }"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        row-key="id"
        stripe
      >
        <template #columns>
          <a-table-column title="时间" data-index="created_at" :width="160">
            <template #cell="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
          </a-table-column>
          <a-table-column title="渠道" data-index="channel_type" :width="90">
            <template #cell="{ record }">
              <a-tag :color="channelTypeColor(record.channel_type)" size="small">
                {{ channelTypeLabel(record.channel_type) }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="渠道名称" data-index="channel_name" :width="140">
            <template #cell="{ record }">
              {{ record.channel_name || '-' }}
            </template>
          </a-table-column>
          <a-table-column title="接收人" data-index="recipient" :width="160">
            <template #cell="{ record }">
              <a-tooltip :content="record.recipient">
                <span class="recipient-text">{{ record.recipient }}</span>
              </a-tooltip>
            </template>
          </a-table-column>
          <a-table-column title="主题/内容" data-index="subject">
            <template #cell="{ record }">
              <a-tooltip :content="record.body || record.subject || '-'">
                <span class="content-text">{{ record.subject || record.body || '-' }}</span>
              </a-tooltip>
            </template>
          </a-table-column>
          <a-table-column title="状态" data-index="status" :width="90">
            <template #cell="{ record }">
              <a-tag :color="LOG_STATUS_MAP[record.status]?.color || 'gray'" size="small">
                {{ LOG_STATUS_MAP[record.status]?.label || record.status }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="重试次数" data-index="attempt_count" :width="90">
            <template #cell="{ record }">
              {{ record.attempt_count }}次
            </template>
          </a-table-column>
          <a-table-column title="操作" :width="80" fixed="right">
            <template #cell="{ record }">
              <a-button type="text" size="small" @click="handleViewDetail(record)">
                「详情」
              </a-button>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>

    <!-- 详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      title="通知详情"
      :width="640"
      :footer="null"
    >
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="日志 ID">{{ currentLog?.id }}</a-descriptions-item>
        <a-descriptions-item label="渠道类型">
          <a-tag :color="channelTypeColor(currentLog?.channel_type)" size="small">
            {{ channelTypeLabel(currentLog?.channel_type) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="渠道名称">{{ currentLog?.channel_name || '-' }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="LOG_STATUS_MAP[currentLog?.status]?.color" size="small">
            {{ LOG_STATUS_MAP[currentLog?.status]?.label }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="接收人">{{ currentLog?.recipient }}</a-descriptions-item>
        <a-descriptions-item label="重试次数">{{ currentLog?.attempt_count }}次</a-descriptions-item>
        <a-descriptions-item label="主题" :span="2">{{ currentLog?.subject || '-' }}</a-descriptions-item>
        <a-descriptions-item label="消息内容" :span="2">
          <pre class="body-content">{{ currentLog?.body || '-' }}</pre>
        </a-descriptions-item>
        <a-descriptions-item label="错误码" :span="2">{{ currentLog?.error_code || '-' }}</a-descriptions-item>
        <a-descriptions-item label="错误信息" :span="2">
          <span class="error-text">{{ currentLog?.error_message || '-' }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ formatTime(currentLog?.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="发送时间">{{ formatTime(currentLog?.sent_at) || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useNotificationLogs, LOG_STATUS_MAP } from '@/composables/useNotification'
import { CHANNEL_TYPE_MAP } from '@/composables/useNotification'

const {
  loading,
  logs,
  pagination,
  filters,
  loadLogs,
  getLogDetail,
  resetFilters
} = useNotificationLogs()

const filterFormRef = ref(null)
const detailVisible = ref(false)
const currentLog = ref(null)
const dateRange = ref([])

const columns = [
  { title: '时间', dataIndex: 'created_at' },
  { title: '渠道', dataIndex: 'channel_type' },
  { title: '渠道名称', dataIndex: 'channel_name' },
  { title: '接收人', dataIndex: 'recipient' },
  { title: '主题/内容', dataIndex: 'content' },
  { title: '状态', dataIndex: 'status' },
  { title: '重试次数', dataIndex: 'attempt_count' },
  { title: '操作', dataIndex: 'action' }
]

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
}))

function channelTypeLabel(type) {
  return CHANNEL_TYPE_MAP[type]?.label || type
}

function channelTypeColor(type) {
  const colorMap = { email: 'blue', sms: 'green', webhook: 'orange' }
  return colorMap[type] || 'gray'
}

function formatTime(time) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

function onDateChange(value) {
  if (value && value.length === 2) {
    filters.start_time = value[0]?.toISOString()
    filters.end_time = value[1]?.toISOString()
  } else {
    filters.start_time = undefined
    filters.end_time = undefined
  }
}

function handleSearch() {
  pagination.current = 1
  loadLogs()
}

function handleReset() {
  dateRange.value = []
  resetFilters()
  loadLogs()
}

function handleRefresh() {
  loadLogs()
  Message.success('已刷新')
}

function handlePageChange(page) {
  pagination.current = page
  loadLogs()
}

function handlePageSizeChange(size) {
  pagination.pageSize = size
  pagination.current = 1
  loadLogs()
}

async function handleViewDetail(record) {
  currentLog.value = await getLogDetail(record.id)
  detailVisible.value = true
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.notification-logs-view {
  padding: 20px;
}
.filter-section {
  background: #fff;
  padding: 16px;
  border-radius: 8px 8px 0 0;
}
.action-section {
  background: #fff;
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e5e6eb;
}
.total-count {
  color: #86909c;
  font-size: 14px;
}
.table-section {
  background: #fff;
  border-radius: 0 0 8px 8px;
}
.recipient-text {
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}
.content-text {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}
.body-content {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  max-height: 200px;
  overflow-y: auto;
}
.error-text {
  color: #f53f3f;
}
</style>
