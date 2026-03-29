<template>
  <div class="page-container">
    <!-- 统计概览 -->
    <a-row :gutter="12" class="stat-row">
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic title="总导出任务" :value="stats.total" :value-from="0" :animation-duration="600">
            <template #icon><icon-file /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card green" hoverable>
          <a-statistic title="已完成" :value="stats.completed" :value-from="0" :animation-duration="600" color="green">
            <template #icon><icon-check-circle /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card orange" hoverable>
          <a-statistic title="进行中" :value="stats.processing" :value-from="0" :animation-duration="600" :colored="stats.processing > 0 ? 'orange' : undefined">
            <template #icon><icon-sync /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card red" hoverable>
          <a-statistic title="失败" :value="stats.failed" :value-from="0" :animation-duration="600" :colored="stats.failed > 0 ? 'red' : undefined">
            <template #icon><icon-close-circle /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Tab 切换 -->
    <a-card class="main-card">
      <a-tabs v-model:active-tab="activeTab">
        <!-- 导出任务列表 -->
        <a-tab key="tasks" title="导出任务">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-select v-model="filter.status" placeholder="任务状态" style="width: 130px" allow-clear @change="loadTasks">
                <a-option value="pending">等待中</a-option>
                <a-option value="processing">处理中</a-option>
                <a-option value="completed">已完成</a-option>
                <a-option value="failed">失败</a-option>
                <a-option value="cancelled">已取消</a-option>
              </a-select>
              <a-select v-model="filter.format" placeholder="导出格式" style="width: 130px" allow-clear @change="loadTasks">
                <a-option value="xlsx">Excel</a-option>
                <a-option value="csv">CSV</a-option>
                <a-option value="json">JSON</a-option>
                <a-option value="pdf">PDF</a-option>
              </a-select>
              <a-input-search v-model="filter.keyword" placeholder="搜索任务名称..." style="width: 200px" search-button @search="loadTasks" />
            </div>
            <div class="toolbar-right">
              <a-button type="primary" @click="showCreateModal = true">
                <template #icon><icon-plus /></template>
                创建导出
              </a-button>
              <a-button @click="loadTasks"><template #icon><icon-refresh /></template>刷新</a-button>
            </div>
          </div>
          <a-table :columns="taskColumns" :data="tasks" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange" :stripe="true" class="mt-3">
            <template #created_at="{ record }"><span>{{ formatDate(record.created_at) }}</span></template>
            <template #status="{ record }"><a-tag :color="taskStatusColor(record.status)">{{ taskStatusLabel(record.status) }}</a-tag></template>
            <template #format="{ record }"><span class="format-badge">{{ record.format?.toUpperCase() }}</span></template>
            <template #progress="{ record }">
              <a-progress v-if="record.status === 'processing'" :percent="record.progress || 0" :color="'arcoblue'" size="small" />
              <a-progress v-else-if="record.status === 'completed'" :percent="100" :color="'green'" size="small" />
              <a-tag v-else-if="record.status === 'failed'" color="red">失败</a-tag>
              <a-tag v-else color="gray">{{ taskStatusLabel(record.status) }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button v-if="record.status === 'completed'" type="text" size="small" @click="downloadFile(record)">
                  <template #icon><icon-download /></template>下载
                </a-button>
                <a-button v-if="record.status === 'pending' || record.status === 'processing'" type="text" size="small" status="warning" @click="cancelTask(record)">取消</a-button>
                <a-button type="text" size="small" @click="openTaskDetail(record)">详情</a-button>
                <a-button type="text" size="small" status="danger" @click="deleteTask(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab>

        <!-- 下载管理 -->
        <a-tab key="downloads" title="下载管理">
          <div class="toolbar-row">
            <div class="toolbar-left">
              <a-input-search v-model="downloadKeyword" placeholder="搜索文件名..." style="width: 200px" search-button @search="loadCompletedTasks" />
            </div>
            <div class="toolbar-right">
              <a-button @click="loadCompletedTasks"><template #icon><icon-refresh /></template>刷新</a-button>
            </div>
          </div>
          <a-table :columns="downloadColumns" :data="completedTasks" :loading="downloadLoading" :pagination="downloadPagination" row-key="id" @change="handleDownloadTableChange" :stripe="true" class="mt-3">
            <template #created_at="{ record }"><span>{{ formatDate(record.created_at) }}</span></template>
            <template #file_size="{ record }"><span>{{ formatSize(record.file_size) }}</span></template>
            <template #expires_at="{ record }"><span :class="isExpiringSoon(record.expires_at) ? 'expiring' : ''">{{ formatDate(record.expires_at) }}</span></template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="primary" size="small" @click="downloadFile(record)"><template #icon><icon-download /></template>下载</a-button>
                <a-button type="text" size="small" status="danger" @click="deleteTask(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab>
      </a-tabs>
    </a-card>

    <!-- 创建导出任务弹窗 -->
    <a-modal v-model:visible="showCreateModal" title="创建导出任务" width="600px" @before-ok="handleCreateTask" :loading="createLoading">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="任务名称" required>
          <a-input v-model="createForm.name" placeholder="请输入导出任务名称" />
        </a-form-item>
        <a-form-item label="数据范围" required>
          <a-select v-model="createForm.data_scope" placeholder="请选择数据范围">
            <a-option value="all">全部数据</a-option>
            <a-option value="device">设备数据</a-option>
            <a-option value="user">用户数据</a-option>
            <a-option value="operation">操作日志</a-option>
            <a-option value="custom">自定义</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker v-model="createForm.date_range" style="width: 100%" />
        </a-form-item>
        <a-form-item label="导出格式" required>
          <a-select v-model="createForm.format" placeholder="请选择导出格式">
            <a-option value="xlsx">Excel (.xlsx)</a-option>
            <a-option value="csv">CSV (.csv)</a-option>
            <a-option value="json">JSON (.json)</a-option>
            <a-option value="pdf">PDF (.pdf)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="压缩">
          <a-switch v-model="createForm.compress" />
          <span class="ml-3">{{ createForm.compress ? '启用压缩（ZIP）' : '不压缩' }}</span>
        </a-form-item>
        <a-form-item label="包含字段">
          <a-checkbox-group v-model="createForm.fields">
            <a-checkbox value="device_id">设备ID</a-checkbox>
            <a-checkbox value="device_name">设备名称</a-checkbox>
            <a-checkbox value="owner">所有者</a-checkbox>
            <a-checkbox value="status">状态</a-checkbox>
            <a-checkbox value="location">位置</a-checkbox>
            <a-checkbox value="created_at">创建时间</a-checkbox>
            <a-checkbox value="last_seen">最后在线</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 任务详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="导出任务详情" width="640px" :footer="null">
      <div v-if="currentTask" class="task-detail">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="任务ID" :span="2"><span class="mono">{{ currentTask.id }}</span></a-descriptions-item>
          <a-descriptions-item label="任务名称" :span="2">{{ currentTask.name }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="taskStatusColor(currentTask.status)">{{ taskStatusLabel(currentTask.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="格式"><span class="format-badge">{{ currentTask.format?.toUpperCase() }}</span></a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatDate(currentTask.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="完成时间">{{ currentTask.completed_at ? formatDate(currentTask.completed_at) : '-' }}</a-descriptions-item>
          <a-descriptions-item label="文件大小">{{ currentTask.file_size ? formatSize(currentTask.file_size) : '-' }}</a-descriptions-item>
          <a-descriptions-item label="过期时间">{{ currentTask.expires_at ? formatDate(currentTask.expires_at) : '-' }}</a-descriptions-item>
          <a-descriptions-item label="数据范围" :span="2">{{ currentTask.data_scope }}</a-descriptions-item>
        </a-descriptions>
        <div v-if="currentTask.error" class="mt-3">
          <a-divider>错误信息</a-divider>
          <div class="code-block error"><pre>{{ currentTask.error }}</pre></div>
        </div>
        <div class="mt-3" style="text-align: right">
          <a-button v-if="currentTask.status === 'completed'" type="primary" @click="downloadFile(currentTask)"><template #icon><icon-download /></template>下载</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getExportTasks,
  createExportTask,
  cancelExportTask,
  deleteExportTask,
  downloadExportFile
} from '@/api/security-evo'
import dayjs from 'dayjs'

const activeTab = ref('tasks')
const loading = ref(false)
const createLoading = ref(false)
const downloadLoading = ref(false)

const tasks = ref([])
const completedTasks = ref([])
const stats = ref({ total: 0, completed: 0, processing: 0, failed: 0 })

const showCreateModal = ref(false)
const detailVisible = ref(false)
const currentTask = ref(null)

const filter = reactive({ status: '', format: '', keyword: '' })
const downloadKeyword = ref('')

const createForm = reactive({
  name: '',
  data_scope: 'all',
  date_range: [],
  format: 'xlsx',
  compress: false,
  fields: ['device_id', 'device_name', 'owner', 'status']
})

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: true })
const downloadPagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: true })

const taskColumns = [
  { title: '任务名称', dataIndex: 'name', width: 180 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '格式', slotName: 'format', width: 80 },
  { title: '进度', slotName: 'progress', width: 160 },
  { title: '数据范围', dataIndex: 'data_scope', width: 120, ellipsis: true },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const downloadColumns = [
  { title: '文件名', dataIndex: 'name', width: 200 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '格式', slotName: 'format', width: 80 },
  { title: '文件大小', slotName: 'file_size', width: 100 },
  { title: '过期时间', slotName: 'expires_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' }
]

function taskStatusColor(s) { const m = { pending: 'gray', processing: 'arcoblue', completed: 'green', failed: 'red', cancelled: 'orange' }; return m[s] || 'default' }
function taskStatusLabel(s) { const m = { pending: '等待中', processing: '处理中', completed: '已完成', failed: '失败', cancelled: '已取消' }; return m[s] || s }
function formatDate(d) { return d ? dayjs(d).format('YYYY-MM-DD HH:mm') : '-' }
function formatSize(bytes) { if (!bytes) return '-'; const k = 1024; return `${(bytes / k / k).toFixed(1)} MB` }
function isExpiringSoon(d) { if (!d) return false; return dayjs(d).diff(dayjs()) < 86400000 * 2 }

function buildParams() { return { page: pagination.current, page_size: pagination.pageSize, status: filter.status || undefined, format: filter.format || undefined, keyword: filter.keyword || undefined } }

onMounted(() => { loadTasks(); loadCompletedTasks() })

async function loadTasks() {
  loading.value = true
  try {
    const r = await getExportTasks(buildParams())
    const d = r.data || r
    tasks.value = d.list || d.records || []
    pagination.total = d.total || tasks.value.length
    // stats
    stats.value.total = tasks.value.length
    stats.value.completed = tasks.value.filter(t => t.status === 'completed').length
    stats.value.processing = tasks.value.filter(t => t.status === 'processing' || t.status === 'pending').length
    stats.value.failed = tasks.value.filter(t => t.status === 'failed').length
  } catch (e) { Message.error('加载导出任务失败') } finally { loading.value = false }
}

async function loadCompletedTasks() {
  downloadLoading.value = true
  try {
    const r = await getExportTasks({ page: downloadPagination.current, page_size: downloadPagination.pageSize, status: 'completed', keyword: downloadKeyword.value || undefined })
    const d = r.data || r
    completedTasks.value = d.list || d.records || []
    downloadPagination.total = d.total || completedTasks.value.length
  } catch (e) { Message.error('加载下载列表失败') } finally { downloadLoading.value = false }
}

function handleTableChange(p) { pagination.current = p.current; pagination.pageSize = p.pageSize; loadTasks() }
function handleDownloadTableChange(p) { downloadPagination.current = p.current; downloadPagination.pageSize = p.pageSize; loadCompletedTasks() }

async function handleCreateTask(done) {
  createLoading.value = true
  try {
    await createExportTask({
      name: createForm.name,
      data_scope: createForm.data_scope,
      start_date: createForm.date_range[0] ? dayjs(createForm.date_range[0]).format('YYYY-MM-DD') : undefined,
      end_date: createForm.date_range[1] ? dayjs(createForm.date_range[1]).format('YYYY-MM-DD') : undefined,
      format: createForm.format,
      compress: createForm.compress,
      fields: createForm.fields
    })
    Message.success('导出任务已创建')
    showCreateModal.value = false
    loadTasks()
    done(true)
  } catch (e) { Message.error('创建失败'); done(false) } finally { createLoading.value = false }
}

async function cancelTask(record) {
  try { await cancelExportTask(record.id); Message.success('任务已取消'); loadTasks() } catch (e) { Message.error('取消失败') }
}

async function deleteTask(record) {
  try { await deleteExportTask(record.id); Message.success('删除成功'); loadTasks(); loadCompletedTasks() } catch (e) { Message.error('删除失败') }
}

async function downloadFile(record) {
  try {
    const blob = await downloadExportFile(record.id)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a'); a.href = url; a.download = `${record.name}.${record.format}`; a.click()
    URL.revokeObjectURL(url)
    Message.success('下载成功')
  } catch (e) { Message.error('下载失败') }
}

function openTaskDetail(record) { currentTask.value = record; detailVisible.value = true }
</script>

<style scoped>
.page-container { padding: 16px; display: flex; flex-direction: column; gap: 12px; height: 100%; box-sizing: border-box; }
.stat-row .stat-card { text-align: center; }
.main-card { flex: 1; overflow: auto; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.toolbar-left { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.toolbar-right { display: flex; gap: 8px; flex-shrink: 0; }
.mt-3 { margin-top: 12px; }
.ml-3 { margin-left: 12px; }
.mono { font-family: 'Courier New', monospace; font-size: 12px; }
.format-badge { background: var(--color-fill-2); border-radius: 4px; padding: 2px 8px; font-size: 12px; font-weight: 600; }
.expiring { color: var(--color-red); }
.task-detail { display: flex; flex-direction: column; gap: 12px; }
.code-block { background: var(--color-fill-1); border-radius: 4px; padding: 12px; }
.code-block pre { margin: 0; font-family: 'Courier New', monospace; font-size: 12px; white-space: pre-wrap; }
.code-block.error { background: rgba(255, 77, 79, 0.1); }
</style>
