<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>临时会员管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="临时会员总数" :value="stats.total" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="今日新增" :value="stats.todayNew" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="已转化" :value="stats.converted" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索姓名/手机号" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="pending">待转化</a-option>
          <a-option value="converted">已转化</a-option>
          <a-option value="expired">已过期</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">新建临时会员</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="convertToFormal(record)" v-if="record.status === 'pending'">转为正式会员</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" title="新建临时会员" @ok="handleSubmit" :width="480" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="姓名" required><a-input v-model="form.name" placeholder="请输入姓名" /></a-form-item>
        <a-form-item label="手机号" required><a-input v-model="form.phone" placeholder="请输入手机号" /></a-form-item>
        <a-form-item label="性别">
          <a-radio-group v-model="form.gender">
            <a-radio value="1">男</a-radio>
            <a-radio value="2">女</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="关联门店">
          <a-select v-model="form.store_id" placeholder="选择门店" allow-clear>
            <a-option v-for="s in stores" :key="s.id" :value="s.id">{{ s.store_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="有效期（天）">
          <a-input-number v-model="form.valid_days" :min="1" :max="365" style="width: 200px;" />
          <span style="margin-left: 8px; color: #999;">不填则使用系统默认值</span>
        </a-form-item>
        <a-form-item label="备注"><a-textarea v-model="form.remark" :rows="2" placeholder="请输入备注" /></a-form-item>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="detailVisible" title="临时会员详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="临时码">{{ current.temp_code || '-' }}</a-descriptions-item>
          <a-descriptions-item label="姓名">{{ current.name }}</a-descriptions-item>
          <a-descriptions-item label="手机号">{{ current.phone }}</a-descriptions-item>
          <a-descriptions-item label="性别">{{ current.gender === '1' ? '男' : current.gender === '2' ? '女' : '-' }}</a-descriptions-item>
          <a-descriptions-item label="关联门店">{{ current.store_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="getStatusColor(current.status)">{{ getStatusText(current.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="备注">{{ current.remark || '-' }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(current.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="过期时间">{{ formatTime(current.expire_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const data = ref([])
const stores = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const current = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, todayNew: 0, converted: 0 })
const form = reactive({ name: '', phone: '', gender: '1', store_id: undefined, valid_days: undefined, remark: '' })

const columns = [
  { title: '临时码', dataIndex: 'temp_code', width: 120 },
  { title: '姓名', dataIndex: 'name' },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '性别', dataIndex: 'gender', width: 70 },
  { title: '关联门店', dataIndex: 'store_name', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '有效期', dataIndex: 'expire_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 240 }
]

const getStatusColor = (s) => ({ pending: 'orange', converted: 'green', expired: 'gray' }[s] || 'gray')
const getStatusText = (s) => ({ pending: '待转化', converted: '已转化', expired: '已过期' }[s] || s)
const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/temp-members?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.converted = data.value.filter(d => d.status === 'converted').length
    }
  } catch (e) { Message.error('加载临时会员失败') } finally { loading.value = false }
}

const loadStores = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) stores.value = d.data?.list || []
  } catch (e) {}
}

const openCreate = () => { Object.assign(form, { name: '', phone: '', gender: '1', store_id: undefined, valid_days: undefined, remark: '' }); modalVisible.value = true }
const openDetail = (r) => { current.value = r; detailVisible.value = true }

const handleSubmit = async () => {
  if (!form.name || !form.phone) { Message.warning('请填写姓名和手机号'); return }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-members`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    const d = await res.json()
    if (d.code === 0) { Message.success('创建成功'); modalVisible.value = false; loadData() }
    else Message.error(d.message || '创建失败')
  } catch (e) { Message.error('创建失败') }
}

const convertToFormal = async (r) => {
  if (!confirm(`确定将「${r.name}」转为正式会员吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-members/${r.id}/convert`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const d = await res.json()
    if (d.code === 0) { Message.success('转化成功'); loadData() } else Message.error(d.message || '转化失败')
  } catch (e) { Message.error('转化失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除临时会员「${r.name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-members/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) { Message.success('删除成功'); loadData() } else Message.error(d.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadStores() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
