<template>
  <div class="page-container">
    <!-- 概览统计 -->
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">总会话数</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-value">{{ stats.active }}</div>
          <div class="stat-label">活跃会话</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-value">{{ stats.current }}</div>
          <div class="stat-label">当前会话</div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <div class="stat-value">{{ stats.expired24h }}</div>
          <div class="stat-label">24h内过期</div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作工具栏 -->
    <a-card class="toolbar-card">
      <div class="toolbar-row">
        <div class="toolbar-left">
          <a-select
            v-model="filter.status"
            placeholder="会话状态"
            style="width: 120px"
            allow-clear
            @change="loadSessions"
          >
            <a-option value="active">活跃</a-option>
            <a-option value="expired">已过期</a-option>
            <a-option value="terminated">已终止</a-option>
          </a-select>
          <a-select
            v-model="filter.deviceType"
            placeholder="设备类型"
            style="width: 140px"
            allow-clear
            @change="loadSessions"
          >
            <a-option value="web">Web</a-option>
            <a-option value="mobile">移动端</a-option>
            <a-option value="tablet">平板</a-option>
            <a-option value="desktop">桌面端</a-option>
            <a-option value="api">API</a-option>
          </a-select>
          <a-range-picker
            v-model="filter.dateRange"
            style="width: 260px"
            @change="loadSessions"
          />
          <a-input-search
            v-model="filter.keyword"
            placeholder="搜索 IP / 设备名称..."
            style="width: 200px"
            @search="loadSessions"
            @press-enter="loadSessions"
          />
        </div>
        <div class="toolbar-right">
          <a-button @click="loadSessions">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </div>
      </div>
    </a-card>

    <!-- 会话列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="sessions"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
        :scroll="{ x: 1100 }"
      >
        <template #device="{ record }">
          <div class="device-cell">
            <icon-computer v-if="record.device_type === 'web' || record.device_type === 'desktop'" class="device-icon" />
            <icon-mobile v-else class="device-icon" />
            <div class="device-info">
              <div>{{ record.device_name || '未知设备' }}</div>
              <div class="device-meta">
                <a-tag size="small">{{ deviceTypeLabel(record.device_type) }}</a-tag>
                <span class="ip-text">{{ record.ip || '-' }}</span>
              </div>
            </div>
          </div>
        </template>
        <template #location="{ record }">
          <span>{{ record.location || record.country || '-' }}</span>
        </template>
        <template #browser="{ record }">
          <span>{{ record.browser || '-' }}</span>
        </template>
        <template #status="{ record }">
          <a-badge :status="sessionStatusBadge(record.status)" />
          <span>{{ statusLabel(record.status) }}</span>
        </template>
        <template #created_at="{ record }">
          <span>{{ formatTime(record.created_at) }}</span>
        </template>
        <template #last_active="{ record }">
          <span>{{ formatTime(record.last_active) }}</span>
        </template>
        <template #actions="{ record }">
          <a-button
            v-if="record.id !== currentSessionId && record.status === 'active'"
            type="text"
            size="small"
            status="danger"
            @click="terminateSession(record)"
          >
            终止
          </a-button>
          <span v-else-if="record.id === currentSessionId" class="current-tag">当前</span>
        </template>
      </a-table>
    </a-card>

    <!-- 批量操作 -->
    <a-card class="batch-card">
      <div class="batch-row">
        <div class="batch-info">
          <span>已选择 {{ selectedSessions.length }} 个会话</span>
        </div>
        <div class="batch-actions">
          <a-popconfirm content="确定终止选中的所有会话？" @ok="batchTerminate">
            <a-button
              type="primary"
              status="danger"
              :disabled="selectedSessions.length === 0"
            >
              批量终止会话
            </a-button>
          </a-popconfirm>
          <a-popconfirm content="确定终止所有其他会话（保留当前）？" @ok="terminateAllOthers">
            <a-button status="warning">
              终止所有其他会话
            </a-button>
          </a-popconfirm>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getSessions,
  terminateSession as terminateSessionApi,
  terminateAllSessions,
  terminateOtherSessions,
  getSessionStats
} from '@/api/security'

const loading = ref(false)
const sessions = ref([])
const selectedSessions = ref([])
const currentSessionId = ref('')

const stats = reactive({
  total: 0,
  active: 0,
  current: 1,
  expired24h: 0
})

const filter = reactive({
  status: '',
  deviceType: '',
  dateRange: [],
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showPageSize: true
})

const columns = [
  { title: '设备', slotName: 'device', width: 220 },
  { title: '位置', slotName: 'location', width: 130 },
  { title: '浏览器', slotName: 'browser', width: 130 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '最后活动', slotName: 'last_active', width: 170 },
  { title: '操作', slotName: 'actions', width: 100, fixed: 'right' }
]

function deviceTypeLabel(type) {
  const map = {
    web: 'Web',
    mobile: '移动端',
    tablet: '平板',
    desktop: '桌面端',
    api: 'API'
  }
  return map[type] || type || '未知'
}

function statusLabel(status) {
  const map = {
    active: '活跃',
    expired: '已过期',
    terminated: '已终止'
  }
  return map[status] || status
}

function sessionStatusBadge(status) {
  const map = {
    active: 'success',
    expired: 'default',
    terminated: 'danger'
  }
  return map[status] || 'default'
}

function formatTime(timestamp) {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString('zh-CN')
}

function handleTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadSessions()
}

async function loadStats() {
  try {
    const data = await getSessionStats()
    const res = data.data || data
    stats.total = res.total || 0
    stats.active = res.active || 0
    stats.expired24h = res.expired_24h || 0
  } catch (e) {
    console.error('加载会话统计失败', e)
  }
}

async function loadSessions() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      status: filter.status || undefined,
      device_type: filter.deviceType || undefined,
      keyword: filter.keyword || undefined,
      start_date: filter.dateRange[0] ? new Date(filter.dateRange[0]).toISOString() : undefined,
      end_date: filter.dateRange[1] ? new Date(filter.dateRange[1]).toISOString() : undefined
    }
    const data = await getSessions(params)
    const res = data.data || data
    sessions.value = res.list || res.records || []
    pagination.total = res.total || sessions.value.length
    currentSessionId.value = res.current_session_id || localStorage.getItem('session_id') || ''
  } catch (e) {
    console.error('加载会话列表失败', e)
    Message.error('加载会话列表失败')
  } finally {
    loading.value = false
  }
}

async function terminateSession(record) {
  try {
    await terminateSessionApi(record.id)
    Message.success('会话已终止')
    loadSessions()
    loadStats()
  } catch (e) {
    Message.error('终止失败')
  }
}

async function batchTerminate() {
  try {
    const ids = selectedSessions.value
    await Promise.all(ids.map(id => terminateSessionApi(id)))
    Message.success(`已终止 ${ids.length} 个会话`)
    selectedSessions.value = []
    loadSessions()
    loadStats()
  } catch (e) {
    Message.error('批量终止失败')
  }
}

async function terminateAllOthers() {
  try {
    await terminateOtherSessions()
    Message.success('已终止所有其他会话')
    loadSessions()
    loadStats()
  } catch (e) {
    Message.error('操作失败')
  }
}

onMounted(() => {
  currentSessionId.value = localStorage.getItem('session_id') || ''
  loadSessions()
  loadStats()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-card {
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: rgb(var(--primary-6));
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-3);
  margin-top: 4px;
}

.toolbar-card {
  flex-shrink: 0;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.table-card {
  flex: 1;
  overflow: auto;
}

.device-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.device-icon {
  font-size: 24px;
  color: var(--color-text-3);
}

.device-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.device-meta {
  display: flex;
  align-items: center;
  gap: 6px;
}

.ip-text {
  font-family: monospace;
  font-size: 12px;
  color: var(--color-text-3);
}

.current-tag {
  color: rgb(var(--primary-6));
  font-size: 12px;
  font-weight: 500;
}

.batch-card {
  flex-shrink: 0;
}

.batch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.batch-info {
  color: var(--color-text-3);
  font-size: 13px;
}

.batch-actions {
  display: flex;
  gap: 8px;
}
</style>
