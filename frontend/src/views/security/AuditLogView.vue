<template>
  <div class="page-container">
    <Breadcrumb :items="['menu.security', 'menu.security.auditLog']" />
    <!-- 筛选工具栏 -->
    <a-card class="toolbar-card">
      <div class="toolbar-row">
        <div class="toolbar-left">
          <a-select
            v-model="filter.module"
            placeholder="模块"
            style="width: 120px"
            allow-clear
            @change="loadLogs"
          >
            <a-option value="auth">认证</a-option>
            <a-option value="device">设备</a-option>
            <a-option value="user">用户</a-option>
            <a-option value="config">配置</a-option>
            <a-option value="data">数据</a-option>
            <a-option value="security">安全</a-option>
          </a-select>
          <a-select
            v-model="filter.level"
            placeholder="级别"
            style="width: 100px"
            allow-clear
            @change="loadLogs"
          >
            <a-option value="info">信息</a-option>
            <a-option value="warning">警告</a-option>
            <a-option value="error">错误</a-option>
          </a-select>
          <a-select
            v-model="filter.action"
            placeholder="操作类型"
            style="width: 140px"
            allow-clear
            @change="loadLogs"
          >
            <a-option value="create">创建</a-option>
            <a-option value="update">更新</a-option>
            <a-option value="delete">删除</a-option>
            <a-option value="read">读取</a-option>
            <a-option value="login">登录</a-option>
            <a-option value="logout">登出</a-option>
          </a-select>
          <a-range-picker
            v-model="filter.dateRange"
            style="width: 260px"
            @change="loadLogs"
          />
          <a-input-search
            v-model="filter.keyword"
            placeholder="搜索操作人/内容..."
            style="width: 200px"
            @search="loadLogs"
            @press-enter="loadLogs"
          />
        </div>
        <div class="toolbar-right">
          <a-button @click="exportLogs">
            <template #icon><icon-download /></template>
            导出日志
          </a-button>
          <a-button @click="refreshLogs">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </div>
      </div>
    </a-card>

    <!-- 日志列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
        :scroll="{ x: 1200 }"
      >
        <template #timestamp="{ record }">
          <span class="timestamp">{{ formatTimestamp(record.timestamp) }}</span>
        </template>
        <template #level="{ record }">
          <a-tag :color="levelColor(record.level)">{{ levelLabel(record.level) }}</a-tag>
        </template>
        <template #module="{ record }">
          <span>{{ moduleLabel(record.module) }}</span>
        </template>
        <template #action="{ record }">
          <a-tag>{{ actionLabel(record.action) }}</a-tag>
        </template>
        <template #user="{ record }">
          <span>{{ record.user_name || record.user_id || '-' }}</span>
        </template>
        <template #ip="{ record }">
          <span class="ip-text">{{ record.ip || '-' }}</span>
        </template>
        <template #result="{ record }">
          <a-tag :color="resultColor(record.result)">{{ resultLabel(record.result) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
        </template>
      </a-table>
    </a-card>

    <!-- 日志详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      title="日志详情"
      width="640px"
      :footer="null"
    >
      <div class="log-detail" v-if="currentLog">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="日志 ID" :span="2">
            <span class="mono">{{ currentLog.id }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="时间">
            {{ formatTimestamp(currentLog.timestamp) }}
          </a-descriptions-item>
          <a-descriptions-item label="级别">
            <a-tag :color="levelColor(currentLog.level)">{{ levelLabel(currentLog.level) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="模块">
            {{ moduleLabel(currentLog.module) }}
          </a-descriptions-item>
          <a-descriptions-item label="操作">
            <a-tag>{{ actionLabel(currentLog.action) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="操作人">
            {{ currentLog.user_name || currentLog.user_id || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="IP 地址">
            {{ currentLog.ip || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="结果">
            <a-tag :color="resultColor(currentLog.result)">{{ resultLabel(currentLog.result) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="资源类型">
            {{ currentLog.resource_type || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="资源 ID" :span="2">
            <span class="mono">{{ currentLog.resource_id || '-' }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="操作描述" :span="2">
            {{ currentLog.description || '-' }}
          </a-descriptions-item>
        </a-descriptions>

        <a-divider>请求详情</a-divider>
        <div class="code-block" v-if="currentLog.request">
          <pre>{{ formatJson(currentLog.request) }}</pre>
        </div>
        <div class="code-block" v-else>
          <span class="no-data">无请求详情</span>
        </div>

        <a-divider>响应详情</a-divider>
        <div class="code-block" v-if="currentLog.response">
          <pre>{{ formatJson(currentLog.response) }}</pre>
        </div>
        <div class="code-block" v-else>
          <span class="no-data">无响应详情</span>
        </div>

        <a-divider>错误信息</a-divider>
        <div class="code-block error" v-if="currentLog.error">
          <pre>{{ currentLog.error }}</pre>
        </div>
        <div class="code-block" v-else>
          <span class="no-data">无错误</span>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getAuditLogs,
  getAuditLogDetail,
  exportAuditLogs
} from '@/api/audit'
import dayjs from 'dayjs'

const loading = ref(false)
const logs = ref([])
const detailVisible = ref(false)
const currentLog = ref(null)

const filter = reactive({
  module: '',
  level: '',
  action: '',
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
  { title: '时间', slotName: 'timestamp', width: 170 },
  { title: '级别', slotName: 'level', width: 80 },
  { title: '模块', slotName: 'module', width: 80 },
  { title: '操作', slotName: 'action', width: 80 },
  { title: '操作人', slotName: 'user', width: 120 },
  { title: 'IP', slotName: 'ip', width: 130 },
  { title: '结果', slotName: 'result', width: 80 },
  { title: '操作', slotName: 'actions', width: 80, fixed: 'right' }
]

function levelColor(level) {
  const map = { info: 'blue', warning: 'orange', error: 'red' }
  return map[level] || 'default'
}

function levelLabel(level) {
  const map = { info: '信息', warning: '警告', error: '错误' }
  return map[level] || level
}

function moduleLabel(module) {
  const map = {
    auth: '认证',
    device: '设备',
    user: '用户',
    config: '配置',
    data: '数据',
    security: '安全'
  }
  return map[module] || module
}

function actionLabel(action) {
  const map = {
    create: '创建',
    update: '更新',
    delete: '删除',
    read: '读取',
    login: '登录',
    logout: '登出'
  }
  return map[action] || action
}

function resultColor(result) {
  const map = { success: 'green', failure: 'red', partial: 'orange' }
  return map[result] || 'default'
}

function resultLabel(result) {
  const map = { success: '成功', failure: '失败', partial: '部分成功' }
  return map[result] || result
}

function formatTimestamp(timestamp) {
  if (!timestamp) return '-'
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

function formatJson(obj) {
  if (!obj) return ''
  try {
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(obj)
  }
}

onMounted(() => {
  loadLogs()
})

async function loadLogs() {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      module: filter.module || undefined,
      level: filter.level || undefined,
      action: filter.action || undefined,
      keyword: filter.keyword || undefined,
      start_date: filter.dateRange[0] ? dayjs(filter.dateRange[0]).format('YYYY-MM-DD') : undefined,
      end_date: filter.dateRange[1] ? dayjs(filter.dateRange[1]).format('YYYY-MM-DD') : undefined
    }
    const res = await getAuditLogs(params)
    const data = res.data || res
    logs.value = data.list || data.records || []
    pagination.total = data.total || logs.value.length
  } catch (e) {
    console.error('加载审计日志失败', e)
    Message.error('加载审计日志失败')
  } finally {
    loading.value = false
  }
}

function handleTableChange(pag) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadLogs()
}

async function refreshLogs() {
  loadLogs()
  Message.success('刷新成功')
}

function openDetail(log) {
  currentLog.value = log
  detailVisible.value = true
}

async function exportLogs() {
  try {
    const params = {
      module: filter.module || undefined,
      level: filter.level || undefined,
      action: filter.action || undefined,
      keyword: filter.keyword || undefined,
      start_date: filter.dateRange[0] ? dayjs(filter.dateRange[0]).format('YYYY-MM-DD') : undefined,
      end_date: filter.dateRange[1] ? dayjs(filter.dateRange[1]).format('YYYY-MM-DD') : undefined
    }
    const blob = await exportAuditLogs(params)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `audit_logs_${dayjs().format('YYYYMMDD_HHmmss')}.json`
    a.click()
    URL.revokeObjectURL(url)
    Message.success('导出成功')
  } catch (e) {
    Message.error('导出失败')
  }
}
</script>

<style scoped>
.page-container {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  box-sizing: border-box;
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

.timestamp {
  font-size: 12px;
  color: var(--color-text-3);
}

.ip-text {
  font-family: monospace;
  font-size: 12px;
}

.mono {
  font-family: monospace;
  font-size: 12px;
}

.log-detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.code-block {
  background: var(--color-fill-1);
  border-radius: 4px;
  padding: 12px;
  max-height: 200px;
  overflow: auto;
}

.code-block pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
}

.code-block.error {
  background: rgba(255, 77, 79, 0.1);
  border: 1px solid var(--color-border);
}

.no-data {
  color: var(--color-text-3);
  font-size: 13px;
}
</style>
