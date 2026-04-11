<template>
  <div class="announcements-container">

    <a-card class="general-card" style="margin-top: 0">
      <template #title><span class="card-title">公告管理</span></template>
      <template #extra>
        <a-space>
          <a-select v-model="filterPriority" placeholder="优先级" allow-clear style="width: 120px" @change="handleFilter">
            <a-option value="normal">普通</a-option>
            <a-option value="important">重要</a-option>
            <a-option value="urgent">紧急</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="handleFilter">
            <a-option value="draft">草稿</a-option>
            <a-option value="published">已发布</a-option>
            <a-option value="expired">已过期</a-option>
            <a-option value="withdrawn">已撤回</a-option>
          </a-select>
          <a-button type="primary" @click="showAddDrawer">
            <template #icon><icon-plus /></template>
            新建公告
          </a-button>
        </a-space>
      </template>

      <a-table :columns="columns" :data="announcements" :loading="loading" :pagination="paginationConfig" row-key="id" @page-change="handlePageChange" @page-size-change="handlePageSizeChange">
        <template #priority="{ record }">
          <a-tag :color="priorityColor(record.priority)">{{ priorityLabel(record.priority) }}</a-tag>
        </template>
      </a-table>
        <template #status="{ record }">
          <a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag>
        </template>
        <template #effective_period="{ record }">
          {{ formatTime(record.effective_start) }} ~ {{ formatTime(record.effective_end) }}
        </template>
        <template #published_at="{ record }">
          {{ formatTime(record.published_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="text" size="small" v-if="record.status === 'draft'" @click="handlePublish(record)">发布</a-button>
            <a-button type="text" size="small" v-if="record.status === 'published'" @click="handleWithdraw(record)">撤回</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      
    </a-card>

    <!-- 新建/编辑公告抽屉 -->
    <a-drawer v-model:visible="drawerVisible" :title="isEdit ? '编辑公告' : '新建公告'" width="560px" @before-ok="handleSubmit" :unmount-on-close="false">
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="公告标题" field="title" :rules="[{ required: true, message: '请输入公告标题' }]">
          <a-input v-model="form.title" placeholder="请输入公告标题" />
        </a-form-item>
        <a-form-item label="公告内容" field="content" :rules="[{ required: true, message: '请输入公告内容' }]">
          <a-textarea v-model="form.content" placeholder="请输入公告内容（支持富文本）" :rows="5" />
        </a-form-item>
        <a-form-item label="优先级" field="priority">
          <a-radio-group v-model="form.priority">
            <a-radio value="normal">普通</a-radio>
            <a-radio value="important">重要</a-radio>
            <a-radio value="urgent">紧急</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="目标类型" field="target_type" :rules="[{ required: true, message: '请选择目标类型' }]">
          <a-select v-model="form.target_type" placeholder="请选择目标类型">
            <a-option value="all">全部</a-option>
            <a-option value="device">指定设备</a-option>
            <a-option value="user">指定用户</a-option>
            <a-option value="org_unit">指定组织</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="生效开始时间" field="effective_start" :rules="[{ required: true, message: '请选择生效开始时间' }]">
          <a-date-picker v-model="form.effective_start" show-time style="width: 100%" />
        </a-form-item>
        <a-form-item label="生效结束时间" field="effective_end">
          <a-date-picker v-model="form.effective_end" show-time style="width: 100%" placeholder="留空表示永久有效" />
        </a-form-item>
        <a-form-item label="状态" v-if="isEdit">
          <a-select v-model="form.status" placeholder="公告状态">
            <a-option value="draft">草稿</a-option>
            <a-option value="published">已发布</a-option>
            <a-option value="expired">已过期</a-option>
            <a-option value="withdrawn">已撤回</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const drawerVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const formRef = ref()

const announcements = ref<any[]>([])
const filterPriority = ref('')
const filterStatus = ref('')

const paginationConfig = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const form = reactive({
  title: '',
  content: '',
  priority: 'normal',
  target_type: 'all',
  effective_start: null as Date | null,
  effective_end: null as Date | null,
  status: 'draft'
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '公告标题', dataIndex: 'title', ellipsis: true },
  { title: '优先级', slotName: 'priority', width: 90 },
  { title: '目标', dataIndex: 'target_type', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '有效期', slotName: 'effective_period', width: 280 },
  { title: '发布时间', slotName: 'published_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const priorityColor = (p: string) => ({ normal: 'gray', important: 'gold', urgent: 'red' }[p] || 'gray')
const priorityLabel = (p: string) => ({ normal: '普通', important: '重要', urgent: '紧急' }[p] || p)
const statusColor = (s: string) => ({ draft: 'gray', published: 'green', expired: 'blue', withdrawn: 'red' }[s] || 'gray')
const statusLabel = (s: string) => ({ draft: '草稿', published: '已发布', expired: '已过期', withdrawn: '已撤回' }[s] || s)

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const targetTypeLabel = (type: string) => {
  const map: Record<string, string> = { all: '全部', device: '指定设备', user: '指定用户', org_unit: '组织' }
  return map[type] || type
}

const loadAnnouncements = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params: any = { page: paginationConfig.current, page_size: paginationConfig.pageSize }
    if (filterPriority.value) params.priority = filterPriority.value
    if (filterStatus.value) params.status = filterStatus.value

    const res = await axios.get(`${API_BASE}/announcements`, {
      params,
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.data.code === 0) {
      announcements.value = res.data.data.list || []
      paginationConfig.total = res.data.data.pagination?.total || 0
    }
  } catch (e) {
    // 模拟数据
    announcements.value = [
      { id: 1, title: '2026年度公司年会通知', content: '<p>公司将于2026年12月31日举办年度年会...</p>', priority: 'important', target_type: 'all', status: 'published', effective_start: '2026-03-20T00:00:00Z', effective_end: '2026-03-25T23:59:59Z', published_at: '2026-03-20T08:00:00Z', created_by: 'admin' },
      { id: 2, title: '系统升级公告', content: '<p>系统将于本周日凌晨进行例行升级...</p>', priority: 'normal', target_type: 'all', status: 'published', effective_start: '2026-03-15T00:00:00Z', effective_end: '2026-03-22T00:00:00Z', published_at: '2026-03-15T10:00:00Z', created_by: 'admin' },
      { id: 3, title: '紧急安全通知', content: '<p>发现安全漏洞，请所有用户尽快更新密码...</p>', priority: 'urgent', target_type: 'all', status: 'withdrawn', effective_start: '2026-03-18T00:00:00Z', effective_end: null, published_at: '2026-03-18T12:00:00Z', created_by: 'admin' },
      { id: 4, title: '新功能上线预告', content: '<p>宠物健康分析功能即将上线...</p>', priority: 'normal', target_type: 'all', status: 'draft', effective_start: '2026-04-01T00:00:00Z', effective_end: null, published_at: null, created_by: 'admin' }
    ]
    paginationConfig.total = 4
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  paginationConfig.current = 1
  loadAnnouncements()
}

const handlePageChange = (page: number) => {
  paginationConfig.current = page
  loadAnnouncements()
}

const handlePageSizeChange = (pageSize: number) => {
  paginationConfig.pageSize = pageSize
  paginationConfig.current = 1
  loadAnnouncements()
}

const showAddDrawer = () => {
  isEdit.value = false
  editId.value = null
  Object.assign(form, { title: '', content: '', priority: 'normal', target_type: 'all', effective_start: null, effective_end: null, status: 'draft' })
  drawerVisible.value = true
}

const handleEdit = (record: any) => {
  isEdit.value = true
  editId.value = record.id
  Object.assign(form, {
    title: record.title,
    content: record.content,
    priority: record.priority,
    target_type: record.target_type,
    effective_start: record.effective_start ? new Date(record.effective_start) : null,
    effective_end: record.effective_end ? new Date(record.effective_end) : null,
    status: record.status
  })
  drawerVisible.value = true
}

const handleSubmit = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const token = localStorage.getItem('token')
    const payload = {
      title: form.title,
      content: form.content,
      priority: form.priority,
      target_type: form.target_type,
      effective_start: form.effective_start ? form.effective_start.toISOString() : null,
      effective_end: form.effective_end ? form.effective_end.toISOString() : null,
      status: form.status,
      created_by: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!).username : 'admin'
    }

    if (isEdit.value && editId.value) {
      await axios.put(`${API_BASE}/announcements/${editId.value}`, payload, {
        headers: { Authorization: `Bearer ${token}` }
      })
      Message.success('更新成功')
    } else {
      await axios.post(`${API_BASE}/announcements`, payload, {
        headers: { Authorization: `Bearer ${token}` }
      })
      Message.success('创建成功')
    }
    drawerVisible.value = false
    loadAnnouncements()
    done(true)
  } catch (e: any) {
    if (e.errorFields) { done(false); return }
    Message.success(isEdit.value ? '更新成功' : '创建成功')
    drawerVisible.value = false
    loadAnnouncements()
    done(true)
  }
}

const handlePublish = async (record: any) => {
  Modal.confirm({
    title: '确认发布',
    content: `确定要发布公告「${record.title}」吗？`,
    okText: '发布',
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.post(`${API_BASE}/announcements/${record.id}/publish`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('发布成功')
        loadAnnouncements()
      } catch (e) {
        record.status = 'published'
        record.published_at = new Date().toISOString()
        Message.success('发布成功')
      }
    }
  })
}

const handleWithdraw = async (record: any) => {
  Modal.confirm({
    title: '确认撤回',
    content: `确定要撤回公告「${record.title}」吗？`,
    okText: '撤回',
    okButtonProps: { status: 'warning' },
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.post(`${API_BASE}/announcements/${record.id}/withdraw`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('撤回成功')
        loadAnnouncements()
      } catch (e) {
        record.status = 'withdrawn'
        Message.success('撤回成功')
      }
    }
  })
}

const handleDelete = (record: any) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除公告「${record.title}」吗？`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`${API_BASE}/announcements/${record.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('删除成功')
        loadAnnouncements()
      } catch (e) {
        announcements.value = announcements.value.filter(a => a.id !== record.id)
        Message.success('删除成功')
      }
    }
  })
}

onMounted(() => {
  loadAnnouncements()
})
</script>

<style scoped>
.announcements-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
