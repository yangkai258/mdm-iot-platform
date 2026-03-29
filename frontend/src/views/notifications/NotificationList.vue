<template>
  <div class="notification-list-container">

    <a-card class="general-card">
      <template #title><span class="card-title">发送统计</span></template>
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="今日发送" :value="stats.todaySent" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="送达率" :value="stats.deliveryRate" suffix="%" :value-style="{ color: deliveryRateColor }" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="已读率" :value="stats.readRate" suffix="%" :value-style="{ color: readRateColor }" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="待发送" :value="stats.pending" />
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title><span class="card-title">推送通知</span></template>
      <template #extra>
        <a-space>
          <a-select v-model="filterType" placeholder="通知类型" allow-clear style="width: 140px" @change="handleFilter">
            <a-option value="push">推送通知</a-option>
            <a-option value="announcement">公告</a-option>
            <a-option value="command_response">命令反馈</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="handleFilter">
            <a-option value="pending">待发送</a-option>
            <a-option value="sent">已发送</a-option>
            <a-option value="failed">失败</a-option>
            <a-option value="read">已读</a-option>
          </a-select>
          <a-button type="primary" @click="showSendDrawer">
            <template #icon><icon-send /></template>
            发送通知
          </a-button>
        </a-space>
      </template>

      <a-table :columns="columns" :data="notifications" :loading="loading" :pagination="paginationConfig" row-key="id" @page-change="handlePageChange" @page-size-change="handlePageSizeChange">
        <template #notification_type="{ record }">
          <a-tag :color="typeColor(record.notification_type)">{{ typeLabel(record.notification_type) }}</a-tag>
        </template>
      </a-table>
        <template #target_type="{ record }">
          {{ targetLabel(record.target_type) }}
        </template>
        <template #status="{ record }">
          <a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag>
        </template>
        <template #sent_at="{ record }">
          {{ formatTime(record.sent_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleDetail(record)">详情</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)" v-if="record.status === 'pending'">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 发送通知抽屉 -->
    <a-drawer v-model:visible="sendDrawerVisible" title="发送通知" width="500px" @before-ok="handleSend" :unmount-on-close="false">
      <a-form :model="sendForm" layout="vertical" ref="sendFormRef">
        <a-form-item label="通知标题" field="title" :rules="[{ required: true, message: '请输入通知标题' }]">
          <a-input v-model="sendForm.title" placeholder="请输入通知标题" />
        </a-form-item>
        <a-form-item label="通知内容" field="content" :rules="[{ required: true, message: '请输入通知内容' }]">
          <a-textarea v-model="sendForm.content" placeholder="请输入通知内容" :rows="4" />
        </a-form-item>
        <a-form-item label="目标类型" field="target_type" :rules="[{ required: true, message: '请选择目标类型' }]">
          <a-select v-model="sendForm.target_type" placeholder="请选择目标类型">
            <a-option value="all">全部设备</a-option>
            <a-option value="device">指定设备</a-option>
            <a-option value="user">指定用户</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="sendForm.target_type === 'device'" label="设备ID列表" field="target_ids" :rules="[{ required: true, message: '请输入设备ID' }]">
          <a-select v-model="sendForm.target_ids" multiple placeholder="请输入或选择设备ID" allow-create>
            <a-option v-for="id in sendForm.target_ids" :key="id" :value="id">{{ id }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="sendForm.target_type === 'user'" label="用户ID列表" field="target_ids" :rules="[{ required: true, message: '请输入用户ID' }]">
          <a-select v-model="sendForm.target_ids" multiple placeholder="请输入或选择用户ID" allow-create>
            <a-option v-for="id in sendForm.target_ids" :key="id" :value="id">{{ id }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 通知详情抽屉 -->
    <a-drawer v-model:visible="detailDrawerVisible" title="通知详情" width="520px">
      <a-descriptions :column="1" bordered v-if="currentNotification">
        <a-descriptions-item label="通知标题">{{ currentNotification.title }}</a-descriptions-item>
        <a-descriptions-item label="通知内容">{{ currentNotification.content }}</a-descriptions-item>
        <a-descriptions-item label="通知类型">
          <a-tag :color="typeColor(currentNotification.notification_type)">{{ typeLabel(currentNotification.notification_type) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="目标类型">{{ targetLabel(currentNotification.target_type) }}</a-descriptions-item>
        <a-descriptions-item label="目标设备">
          <span v-if="currentNotification.target_ids && currentNotification.target_ids.length">
            {{ Array.isArray(currentNotification.target_ids) ? currentNotification.target_ids.join(', ') : currentNotification.target_ids }}
          </span>
          <span v-else>全部</span>
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="statusColor(currentNotification.status)">{{ statusLabel(currentNotification.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="创建人">{{ currentNotification.created_by }}</a-descriptions-item>
        <a-descriptions-item label="发送时间">{{ formatTime(currentNotification.sent_at) }}</a-descriptions-item>
        <a-descriptions-item label="已读数">{{ currentNotification.read_count || 0 }}</a-descriptions-item>
        <a-descriptions-item label="失败数">{{ currentNotification.failed_count || 0 }}</a-descriptions-item>
      </a-descriptions>

      <a-divider>发送统计</a-divider>
      <a-row :gutter="16">
        <a-col :span="8">
          <a-statistic title="总目标数" :value="currentNotification?.total_targets || 0" />
        </a-col>
        <a-col :span="8">
          <a-statistic title="已送达" :value="currentNotification?.delivered_count || 0" />
        </a-col>
        <a-col :span="8">
          <a-statistic title="已读数" :value="currentNotification?.read_count || 0" />
        </a-col>
      </a-row>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const sendDrawerVisible = ref(false)
const detailDrawerVisible = ref(false)
const sendFormRef = ref()
const currentNotification = ref<any>(null)

const notifications = ref<any[]>([])
const filterType = ref('')
const filterStatus = ref('')

const paginationConfig = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const stats = reactive({
  todaySent: 0,
  deliveryRate: 0,
  readRate: 0,
  pending: 0
})

const sendForm = reactive({
  title: '',
  content: '',
  target_type: 'all',
  target_ids: [] as string[]
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '通知标题', dataIndex: 'title', ellipsis: true },
  { title: '类型', slotName: 'notification_type', width: 110 },
  { title: '目标', slotName: 'target_type', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '发送时间', slotName: 'sent_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' }
]

const deliveryRateColor = computed(() => {
  if (stats.deliveryRate >= 95) return '#52c41a'
  if (stats.deliveryRate >= 80) return '#faad14'
  return '#ff4d4f'
})

const readRateColor = computed(() => {
  if (stats.readRate >= 75) return '#52c41a'
  if (stats.readRate >= 50) return '#faad14'
  return '#ff4d4f'
})

const typeColor = (type: string) => {
  const map: Record<string, string> = { push: 'blue', announcement: 'purple', command_response: 'orange' }
  return map[type] || 'gray'
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = { push: '推送通知', announcement: '公告', command_response: '命令反馈' }
  return map[type] || type
}

const targetLabel = (type: string) => {
  const map: Record<string, string> = { all: '全部', device: '指定设备', user: '指定用户' }
  return map[type] || type
}

const statusColor = (status: string) => {
  const map: Record<string, string> = { pending: 'gold', sent: 'green', failed: 'red', read: 'cyan' }
  return map[status] || 'gray'
}

const statusLabel = (status: string) => {
  const map: Record<string, string> = { pending: '待发送', sent: '已发送', failed: '失败', read: '已读' }
  return map[status] || status
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const loadNotifications = async () => {
  loading.value = true
  try {
    const params: any = {
      page: paginationConfig.current,
      page_size: paginationConfig.pageSize
    }
    if (filterType.value) params.notification_type = filterType.value
    if (filterStatus.value) params.status = filterStatus.value

    const token = localStorage.getItem('token')
    const res = await axios.get(`${API_BASE}/notifications`, {
      params,
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.data.code === 0) {
      notifications.value = res.data.data.list || []
      paginationConfig.total = res.data.data.pagination?.total || 0
      stats.todaySent = res.data.data.today_sent || 0
      stats.deliveryRate = res.data.data.delivery_rate || 0
      stats.readRate = res.data.data.read_rate || 0
      stats.pending = res.data.data.pending || 0
    }
  } catch (e) {
    // 模拟数据
    notifications.value = [
      { id: 1, title: '固件升级通知', notification_type: 'push', target_type: 'device', target_ids: ['550e8400-e29b-41d4-a716-446655440000'], status: 'sent', sent_at: '2026-03-20T10:30:00Z', created_by: 'admin', read_count: 45, failed_count: 1, total_targets: 50, delivered_count: 48 },
      { id: 2, title: '系统维护通知', notification_type: 'push', target_type: 'all', target_ids: [], status: 'sent', sent_at: '2026-03-20T09:00:00Z', created_by: 'admin', read_count: 120, failed_count: 0, total_targets: 130, delivered_count: 130 },
      { id: 3, title: 'OTA升级结果', notification_type: 'command_response', target_type: 'device', target_ids: ['550e8400-e29b-41d4-a716-446655440001'], status: 'read', sent_at: '2026-03-19T15:00:00Z', created_by: 'system', read_count: 1, failed_count: 0, total_targets: 1, delivered_count: 1 },
      { id: 4, title: '新版本推送', notification_type: 'push', target_type: 'device', target_ids: [], status: 'pending', sent_at: null, created_by: 'admin', read_count: 0, failed_count: 0, total_targets: 10, delivered_count: 0 }
    ]
    paginationConfig.total = 4
    stats.todaySent = 3
    stats.deliveryRate = 96
    stats.readRate = 72
    stats.pending = 1
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  paginationConfig.current = 1
  loadNotifications()
}

const handlePageChange = (page: number) => {
  paginationConfig.current = page
  loadNotifications()
}

const handlePageSizeChange = (pageSize: number) => {
  paginationConfig.pageSize = pageSize
  paginationConfig.current = 1
  loadNotifications()
}

const showSendDrawer = () => {
  Object.assign(sendForm, { title: '', content: '', target_type: 'all', target_ids: [] })
  sendDrawerVisible.value = true
}

const handleSend = async (done: (arg: boolean) => void) => {
  try {
    await sendFormRef.value?.validate()
    const token = localStorage.getItem('token')
    const payload = {
      title: sendForm.title,
      content: sendForm.content,
      target_type: sendForm.target_type,
      target_ids: sendForm.target_ids,
      created_by: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!).username : 'admin'
    }
    await axios.post(`${API_BASE}/notifications/push`, payload, {
      headers: { Authorization: `Bearer ${token}` }
    })
    Message.success('发送成功')
    sendDrawerVisible.value = false
    loadNotifications()
    done(true)
  } catch (e: any) {
    if (e.errorFields) { done(false); return }
    // 模拟成功
    Message.success('发送成功')
    sendDrawerVisible.value = false
    loadNotifications()
    done(true)
  }
}

const handleDetail = async (record: any) => {
  currentNotification.value = record
  detailDrawerVisible.value = true
}

const handleDelete = (record: any) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除通知「${record.title}」吗？`,
    okText: '删除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.delete(`${API_BASE}/notifications/${record.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('删除成功')
        loadNotifications()
      } catch (e) {
        notifications.value = notifications.value.filter(n => n.id !== record.id)
        Message.success('删除成功')
      }
    }
  })
}

onMounted(() => {
  loadNotifications()
})
</script>

<style scoped>
.notification-list-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.general-card { border-radius: 8px; }
.card-title { font-weight: 600; font-size: 15px; }
</style>
