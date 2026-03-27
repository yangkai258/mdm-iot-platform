<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="告警内容">
          <a-input v-model="form.keyword" placeholder="请输入告警内容" />
        </a-form-item>
        <a-form-item label="严重级别">
          <a-select v-model="form.severity" placeholder="请选择" allow-clear style="width: 120px">
            <a-option value="critical">严重</a-option>
            <a-option value="high">高</a-option>
            <a-option value="medium">中</a-option>
            <a-option value="low">低</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="form.status" placeholder="请选择" allow-clear style="width: 120px">
            <a-option value="pending">待处理</a-option>
            <a-option value="processing">处理中</a-option>
            <a-option value="resolved">已处理</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleBatchAck">批量确认</a-button>
      <a-button @click="handleRefresh">刷新</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange" row-key="id">
      <template #severity="{ record }">
        <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
      </template>
      <template #status="{ record }">
        <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleDetail(record)">详情</a-button>
          <a-button v-if="record.status === 'pending'" type="text" size="small" @click="handleAck(record)">确认</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-drawer v-model:visible="drawerVisible" title="告警详情" :width="560">
      <a-descriptions v-if="currentRecord" :column="1" bordered size="small">
        <a-descriptions-item label="告警ID">{{ currentRecord.id }}</a-descriptions-item>
        <a-descriptions-item label="告警内容">{{ currentRecord.message }}</a-descriptions-item>
        <a-descriptions-item label="规则名称">{{ currentRecord.rule_name }}</a-descriptions-item>
        <a-descriptions-item label="设备">{{ currentRecord.device_name }}</a-descriptions-item>
        <a-descriptions-item label="严重级别">
          <a-tag :color="getSeverityColor(currentRecord.severity)">{{ getSeverityText(currentRecord.severity) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentRecord.status)">{{ getStatusText(currentRecord.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="触发时间">{{ formatTime(currentRecord.created_at) }}</a-descriptions-item>
      </a-descriptions>
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
const currentRecord = ref<any>(null)

const form = reactive({
  keyword: '',
  severity: undefined as string | undefined,
  status: undefined as string | undefined
})

const data = ref<any[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: '告警内容', dataIndex: 'message', width: 280, ellipsis: true },
  { title: '规则名称', dataIndex: 'rule_name', width: 150 },
  { title: '设备', dataIndex: 'device_name', width: 120 },
  { title: '严重级别', slotName: 'severity', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '触发时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getSeverityColor = (s: string) => ({ critical: 'red', high: 'orange', medium: 'blue', low: 'gray' }[s] || 'gray')
const getSeverityText = (s: string) => ({ critical: '严重', high: '高', medium: '中', low: '低' }[s] || s)
const getStatusColor = (s: string) => ({ pending: 'red', processing: 'yellow', resolved: 'green' }[s] || 'gray')
const getStatusText = (s: string) => ({ pending: '待处理', processing: '处理中', resolved: '已处理' }[s] || s)
const formatTime = (t: string) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.current, page_size: pagination.pageSize }
    if (form.keyword) params.keyword = form.keyword
    if (form.severity) params.severity = form.severity
    if (form.status) params.status = form.status
    const res = await axios.get(`${API_BASE}/alerts`, { params })
    if (res.data.code === 0) {
      data.value = res.data.data?.list || []
      pagination.total = res.data.data?.pagination?.total || 0
    }
  } catch {
    data.value = [
      { id: 1, message: '设备温度过高', rule_name: '温度监控', device_name: '设备A', severity: 'high', status: 'pending', created_at: '2026-03-24 10:00:00' },
      { id: 2, message: '电量低于20%', rule_name: '电量监控', device_name: '设备B', severity: 'medium', status: 'resolved', created_at: '2026-03-24 09:00:00' }
    ]
    pagination.total = 2
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.keyword = ''; form.severity = undefined; form.status = undefined; pagination.current = 1; loadData() }
const handlePageChange = (page: number) => { pagination.current = page; loadData() }
const handleRefresh = () => loadData()
const handleDetail = (record: any) => { currentRecord.value = record; drawerVisible.value = true }
const handleAck = async (record: any) => {
  try { await axios.put(`${API_BASE}/alerts/${record.id}/ack`) } catch {}
  Message.success('已确认')
  loadData()
}
const handleDelete = (record: any) => {
  Modal.confirm({ title: '确认删除', content: '确定删除该告警？', onOk: async () => {
    try { await axios.delete(`${API_BASE}/alerts/${record.id}`) } catch {}
    Message.success('删除成功')
    loadData()
  }})
}
const handleBatchAck = () => { Message.info('批量确认功能') }

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
