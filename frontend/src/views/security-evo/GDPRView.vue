<template>
  <div class="page-container">
    <Breadcrumb :items="['menu.security', 'menu.security.gdpr']" />
    <!-- 统计概览 -->
    <a-row :gutter="12" class="stat-row">
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic title="总请求数" :value="stats.total" :value-from="0" :animation-duration="600">
            <template #icon><icon-file /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card orange" hoverable>
          <a-statistic title="待处理" :value="stats.pending" :value-from="0" :animation-duration="600" :colored="stats.pending > 0 ? 'orange' : undefined">
            <template #icon><icon-clock /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card blue" hoverable>
          <a-statistic title="处理中" :value="stats.processing" :value-from="0" :animation-duration="600" color="arcoblue">
            <template #icon><icon-sync /></template>
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
    </a-row>

    <!-- 工具栏 -->
    <a-card class="toolbar-card">
      <div class="toolbar-row">
        <div class="toolbar-left">
          <a-select v-model="filter.status" placeholder="请求状态" style="width: 130px" allow-clear @change="loadRequests">
            <a-option value="pending">待处理</a-option>
            <a-option value="approved">已批准</a-option>
            <a-option value="processing">处理中</a-option>
            <a-option value="completed">已完成</a-option>
            <a-option value="rejected">已拒绝</a-option>
          </a-select>
          <a-select v-model="filter.type" placeholder="请求类型" style="width: 160px" allow-clear @change="loadRequests">
            <a-option value="access">数据访问</a-option>
            <a-option value="rectification">数据更正</a-option>
            <a-option value="erasure">数据删除</a-option>
            <a-option value="portability">数据可携带</a-option>
            <a-option value="objection">拒绝处理</a-option>
            <a-option value="restriction">处理限制</a-option>
          </a-select>
          <a-range-picker v-model="filter.dateRange" style="width: 260px" @change="loadRequests" />
          <a-input-search v-model="filter.keyword" placeholder="搜索请求者..." style="width: 200px" search-button @search="loadRequests" />
        </div>
        <div class="toolbar-right">
          <a-button type="primary" @click="showCreateModal = true">
            <template #icon><icon-plus /></template>
            创建请求
          </a-button>
          <a-button @click="loadRequests"><template #icon><icon-refresh /></template>刷新</a-button>
        </div>
      </div>
    </a-card>

    <!-- 请求列表 -->
    <a-card class="table-card">
      <a-table
        :columns="columns"
        :data="requests"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
        :stripe="true"
        :scroll="{ x: 1400 }"
      >
        <template #created_at="{ record }"><span>{{ formatDate(record.created_at) }}</span></template>
        <template #type="{ record }"><a-tag :color="typeColor(record.type)">{{ typeLabel(record.type) }}</a-tag></template>
        <template #status="{ record }"><a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag></template>
        <template #priority="{ record }"><a-tag :color="priorityColor(record.priority)">{{ priorityLabel(record.priority) }}</a-tag></template>
        <template #due_date="{ record }">
          <span :class="isOverdue(record.due_date) ? 'overdue' : ''">
            {{ record.due_date ? formatDate(record.due_date) : '-' }}
          </span>
        </template>
        <template #progress="{ record }">
          <a-progress v-if="record.status === 'processing'" :percent="record.progress || 0" :color="'arcoblue'" size="small" />
          <span v-else>-</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-dropdown trigger="click">
              <a-button type="text" size="small">更多</a-button>
              <template #content>
                <a-doption v-if="record.status === 'pending'" @click="approveRequest(record)">批准</a-doption>
                <a-doption v-if="record.status === 'pending'" @click="rejectRequest(record)">拒绝</a-doption>
                <a-doption v-if="record.status === 'approved' || record.status === 'processing'" @click="markCompleted(record)">标记完成</a-doption>
                <a-doption v-if="record.status !== 'completed' && record.status !== 'rejected'" @click="addProgressNote(record)">添加进度</a-doption>
              </template>
            </a-dropdown>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 详情弹窗 -->
    <a-drawer v-model:visible="detailVisible" :title="`GDPR 请求详情 - ${currentRequest?.id?.slice(0, 8) || ''}`" width="640px" :footer="detailFooter">
      <div v-if="currentRequest" class="gdpr-detail">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="请求ID" :span="2"><span class="mono">{{ currentRequest.id }}</span></a-descriptions-item>
          <a-descriptions-item label="类型"><a-tag :color="typeColor(currentRequest.type)">{{ typeLabel(currentRequest.type) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="statusColor(currentRequest.status)">{{ statusLabel(currentRequest.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="请求者">{{ currentRequest.requester_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="请求者邮箱">{{ currentRequest.requester_email || '-' }}</a-descriptions-item>
          <a-descriptions-item label="优先级"><a-tag :color="priorityColor(currentRequest.priority)">{{ priorityLabel(currentRequest.priority) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="截止日期">{{ currentRequest.due_date ? formatDate(currentRequest.due_date) : '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatDate(currentRequest.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="处理人">{{ currentRequest.handler_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="完成时间">{{ currentRequest.completed_at ? formatDate(currentRequest.completed_at) : '-' }}</a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">{{ currentRequest.description || '-' }}</a-descriptions-item>
        </a-descriptions>

        <a-divider>处理进度</a-divider>
        <div class="progress-section">
          <a-steps :current="progressSteps(currentRequest)" status="process" size="small">
            <a-step title="待处理" description="" />
            <a-step title="已批准" description="" />
            <a-step title="处理中" description="" />
            <a-step title="已完成" description="" />
          </a-steps>
          <div v-if="currentRequest.progress_notes?.length" class="notes-list mt-3">
            <div v-for="note in currentRequest.progress_notes" :key="note.id" class="note-item">
              <div class="note-header">
                <span class="note-time">{{ formatDate(note.created_at) }}</span>
                <span class="note-author">{{ note.author_name }}</span>
              </div>
              <div class="note-content">{{ note.content }}</div>
            </div>
          </div>
          <div v-else class="no-data">暂无进度记录</div>
        </div>

        <a-divider>相关资源</a-divider>
        <div v-if="currentRequest.related_resources?.length" class="resource-list">
          <div v-for="r in currentRequest.related_resources" :key="r.id" class="resource-item">
            <span class="resource-type">{{ r.type }}</span>
            <span class="resource-id mono">{{ r.id }}</span>
          </div>
        </div>
        <div v-else class="no-data">暂无关联资源</div>
      </div>
      <template #footer>
        <a-space>
          <a-button v-if="currentRequest?.status === 'pending'" @click="rejectRequest(currentRequest)">拒绝</a-button>
          <a-button v-if="currentRequest?.status === 'pending'" type="primary" @click="approveRequest(currentRequest)">批准</a-button>
          <a-button v-if="currentRequest?.status === 'approved' || currentRequest?.status === 'processing'" type="primary" @click="markCompleted(currentRequest)">标记完成</a-button>
          <a-button @click="detailVisible = false">关闭</a-button>
        </a-space>
      </template>
    </a-drawer>

    <!-- 创建请求弹窗 -->
    <a-modal v-model:visible="showCreateModal" title="创建 GDPR 请求" width="560px" @before-ok="handleCreateRequest" :loading="createLoading">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="请求类型" required>
          <a-select v-model="createForm.type" placeholder="请选择请求类型">
            <a-option value="access">数据访问请求</a-option>
            <a-option value="rectification">数据更正请求</a-option>
            <a-option value="erasure">数据删除请求（被遗忘权）</a-option>
            <a-option value="portability">数据可携带请求</a-option>
            <a-option value="objection">拒绝处理申请</a-option>
            <a-option value="restriction">处理限制申请</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="请求者姓名" required>
          <a-input v-model="createForm.requester_name" placeholder="请输入请求者姓名" />
        </a-form-item>
        <a-form-item label="请求者邮箱" required>
          <a-input v-model="createForm.requester_email" placeholder="请输入请求者邮箱" />
        </a-form-item>
        <a-form-item label="优先级">
          <a-select v-model="createForm.priority">
            <a-option value="low">低</a-option>
            <a-option value="medium">中</a-option>
            <a-option value="high">高</a-option>
            <a-option value="urgent">紧急</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="截止日期">
          <a-date-picker v-model="createForm.due_date" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="createForm.description" :rows="4" placeholder="请详细描述请求内容及背景" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 拒绝弹窗 -->
    <a-modal v-model:visible="rejectModalVisible" title="拒绝请求" width="480px" @before-ok="handleReject" :loading="actionLoading">
      <a-form :model="rejectForm" layout="vertical">
        <a-form-item label="拒绝原因" required>
          <a-textarea v-model="rejectForm.reason" :rows="4" placeholder="请说明拒绝原因" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 添加进度弹窗 -->
    <a-modal v-model:visible="progressModalVisible" title="添加进度记录" width="480px" @before-ok="handleAddProgress" :loading="actionLoading">
      <a-form :model="progressForm" layout="vertical">
        <a-form-item label="进度说明" required>
          <a-textarea v-model="progressForm.content" :rows="4" placeholder="请说明当前处理进度" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  getGdprRequests,
  getGdprRequestDetail,
  createGdprRequest,
  approveGdprRequest,
  rejectGdprRequest,
  completeGdprRequest,
  updateGdprRequest,
  getGdprStatistics
} from '@/api/security-evo'
import dayjs from 'dayjs'

const loading = ref(false)
const createLoading = ref(false)
const actionLoading = ref(false)

const requests = ref([])
const stats = ref({ total: 0, pending: 0, processing: 0, completed: 0 })
const detailVisible = ref(false)
const showCreateModal = ref(false)
const rejectModalVisible = ref(false)
const progressModalVisible = ref(false)
const currentRequest = ref(null)
const currentActionRequest = ref(null)

const filter = reactive({ status: '', type: '', dateRange: [], keyword: '' })

const createForm = reactive({
  type: 'access',
  requester_name: '',
  requester_email: '',
  priority: 'medium',
  due_date: null,
  description: ''
})

const rejectForm = reactive({ reason: '' })
const progressForm = reactive({ content: '' })

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: true })

const detailFooter = ref(null)

const columns = [
  { title: '请求ID', dataIndex: 'id', width: 120, ellipsis: true },
  { title: '类型', slotName: 'type', width: 130 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '优先级', slotName: 'priority', width: 90 },
  { title: '请求者', dataIndex: 'requester_name', width: 130 },
  { title: '邮箱', dataIndex: 'requester_email', width: 180, ellipsis: true },
  { title: '截止日期', slotName: 'due_date', width: 160 },
  { title: '进度', slotName: 'progress', width: 120 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' }
]

function typeColor(t) { const m = { access: 'arcoblue', rectification: 'orange', erasure: 'red', portability: 'green', objection: 'purple', restriction: 'gold' }; return m[t] || 'default' }
function typeLabel(t) { const m = { access: '数据访问', rectification: '数据更正', erasure: '数据删除', portability: '数据可携带', objection: '拒绝处理', restriction: '处理限制' }; return m[t] || t }
function statusColor(s) { const m = { pending: 'gray', approved: 'arcoblue', processing: 'orange', completed: 'green', rejected: 'red' }; return m[s] || 'default' }
function statusLabel(s) { const m = { pending: '待处理', approved: '已批准', processing: '处理中', completed: '已完成', rejected: '已拒绝' }; return m[s] || s }
function priorityColor(p) { const m = { low: 'green', medium: 'arcoblue', high: 'orange', urgent: 'red' }; return m[p] || 'default' }
function priorityLabel(p) { const m = { low: '低', medium: '中', high: '高', urgent: '紧急' }; return m[p] || p }
function formatDate(d) { return d ? dayjs(d).format('YYYY-MM-DD HH:mm') : '-' }
function isOverdue(d) { return d && dayjs(d).isBefore(dayjs()) }
function progressSteps(r) { const m = { pending: 0, approved: 1, processing: 2, completed: 3 }; return m[r.status] ?? 0 }

function buildParams() {
  return {
    page: pagination.current,
    page_size: pagination.pageSize,
    status: filter.status || undefined,
    type: filter.type || undefined,
    keyword: filter.keyword || undefined,
    start_date: filter.dateRange[0] ? dayjs(filter.dateRange[0]).format('YYYY-MM-DD') : undefined,
    end_date: filter.dateRange[1] ? dayjs(filter.dateRange[1]).format('YYYY-MM-DD') : undefined
  }
}

onMounted(() => { loadRequests(); loadStats() })

async function loadStats() {
  try { const r = await getGdprStatistics(); stats.value = r.data || r } catch (e) { console.error(e) }
}

async function loadRequests() {
  loading.value = true
  try {
    const r = await getGdprRequests(buildParams())
    const d = r.data || r
    requests.value = d.list || d.records || []
    pagination.total = d.total || requests.value.length
  } catch (e) { Message.error('加载请求失败') } finally { loading.value = false }
}

function handleTableChange(p) { pagination.current = p.current; pagination.pageSize = p.pageSize; loadRequests() }

async function openDetail(record) {
  currentRequest.value = null
  detailVisible.value = true
  try {
    const r = await getGdprRequestDetail(record.id)
    currentRequest.value = r.data || r
  } catch (e) { Message.error('加载详情失败') }
}

async function handleCreateRequest(done) {
  createLoading.value = true
  try {
    await createGdprRequest({
      type: createForm.type,
      requester_name: createForm.requester_name,
      requester_email: createForm.requester_email,
      priority: createForm.priority,
      due_date: createForm.due_date ? dayjs(createForm.due_date).format('YYYY-MM-DD') : undefined,
      description: createForm.description
    })
    Message.success('请求创建成功')
    showCreateModal.value = false
    loadRequests()
    done(true)
  } catch (e) { Message.error('创建失败'); done(false) } finally { createLoading.value = false }
}

async function approveRequest(record) {
  actionLoading.value = true
  try {
    await approveGdprRequest(record.id)
    Message.success('请求已批准')
    loadRequests()
    if (detailVisible.value) openDetail({ id: record.id })
  } catch (e) { Message.error('批准失败') } finally { actionLoading.value = false }
}

function rejectRequest(record) {
  currentActionRequest.value = record
  rejectForm.reason = ''
  rejectModalVisible.value = true
}

async function handleReject(done) {
  actionLoading.value = true
  try {
    await rejectGdprRequest(currentActionRequest.value.id, { reason: rejectForm.reason })
    Message.success('请求已拒绝')
    rejectModalVisible.value = false
    loadRequests()
    if (detailVisible.value) openDetail({ id: currentActionRequest.value.id })
    done(true)
  } catch (e) { Message.error('拒绝失败'); done(false) } finally { actionLoading.value = false }
}

function addProgressNote(record) {
  currentActionRequest.value = record
  progressForm.content = ''
  progressModalVisible.value = true
}

async function handleAddProgress(done) {
  actionLoading.value = true
  try {
    await updateGdprRequest(currentActionRequest.value.id, { progress_note: progressForm.content })
    Message.success('进度已添加')
    progressModalVisible.value = false
    if (detailVisible.value) openDetail({ id: currentActionRequest.value.id })
    done(true)
  } catch (e) { Message.error('添加失败'); done(false) } finally { actionLoading.value = false }
}

async function markCompleted(record) {
  actionLoading.value = true
  try {
    await completeGdprRequest(record.id)
    Message.success('请求已标记完成')
    loadRequests()
    if (detailVisible.value) openDetail({ id: record.id })
  } catch (e) { Message.error('操作失败') } finally { actionLoading.value = false }
}
</script>

<style scoped>
.page-container { padding: 16px; display: flex; flex-direction: column; gap: 12px; height: 100%; box-sizing: border-box; }
.stat-row .stat-card { text-align: center; }
.toolbar-card { flex-shrink: 0; }
.toolbar-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.toolbar-left { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.toolbar-right { display: flex; gap: 8px; flex-shrink: 0; }
.table-card { flex: 1; overflow: auto; }
.mt-3 { margin-top: 12px; }
.mono { font-family: 'Courier New', monospace; font-size: 12px; }
.overdue { color: var(--color-red); }
.gdpr-detail { display: flex; flex-direction: column; gap: 12px; }
.notes-list { display: flex; flex-direction: column; gap: 10px; }
.note-item { background: var(--color-fill-1); border-radius: 4px; padding: 10px; }
.note-header { display: flex; justify-content: space-between; margin-bottom: 4px; font-size: 12px; color: var(--color-text-3); }
.note-content { font-size: 13px; }
.resource-list { display: flex; flex-direction: column; gap: 8px; }
.resource-item { display: flex; align-items: center; gap: 10px; }
.resource-type { background: var(--color-fill-2); border-radius: 4px; padding: 2px 8px; font-size: 12px; }
.no-data { color: var(--color-text-3); font-size: 13px; }
</style>
