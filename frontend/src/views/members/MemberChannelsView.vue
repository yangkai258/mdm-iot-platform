<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>会员渠道设置</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="渠道总数" :value="stats.total" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="已启用" :value="stats.enabled" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="本月新增会员" :value="stats.monthNew" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索渠道名称" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.type" placeholder="渠道类型" allow-clear style="width: 140px" @change="loadData">
          <a-option value="online">线上</a-option>
          <a-option value="offline">线下</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">新建渠道</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #type="{ record }"><a-tag :color="getTypeColor(record.channel_type)">{{ getTypeText(record.channel_type) }}</a-tag></template>
        <template #status="{ record }"><a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag></template>
        <template #memberCount="{ record }">{{ record.member_count || 0 }}</template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑渠道' : '新建渠道'" @ok="handleSubmit" :width="480" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="渠道名称" required><a-input v-model="form.channel_name" placeholder="请输入渠道名称" /></a-form-item>
        <a-form-item label="渠道类型">
          <a-select v-model="form.channel_type" placeholder="选择渠道类型">
            <a-option value="wechat">微信公众号</a-option>
            <a-option value="miniprogram">微信小程序</a-option>
            <a-option value="app">APP</a-option>
            <a-option value="web">官网</a-option>
            <a-option value="offline">线下</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="AppID"><a-input v-model="form.app_id" placeholder="微信公众号/小程序AppID" /></a-form-item>
        <a-form-item label="AppSecret"><a-input v-model="form.app_secret" placeholder="微信公众号/小程序AppSecret" type="password" /></a-form-item>
        <a-form-item label="备注"><a-textarea v-model="form.remark" :rows="2" placeholder="请输入备注" /></a-form-item>
        <a-form-item label="状态"><a-switch v-model="formStatus" checked-value="1" unchecked-value="0" /></a-form-item>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="detailVisible" title="渠道详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="渠道名称">{{ current.channel_name }}</a-descriptions-item>
          <a-descriptions-item label="渠道类型"><a-tag>{{ getTypeText(current.channel_type) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="AppID">{{ current.app_id || '-' }}</a-descriptions-item>
          <a-descriptions-item label="会员数">{{ current.member_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="current.status === 1 ? 'green' : 'gray'">{{ current.status === 1 ? '启用' : '禁用' }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="备注">{{ current.remark || '-' }}</a-descriptions-item>
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
const loading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const current = ref(null)
const currentId = ref(null)

const filters = reactive({ keyword: '', type: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, enabled: 0, monthNew: 0 })
const form = reactive({ channel_name: '', channel_type: 'wechat', app_id: '', app_secret: '', remark: '' })
const formStatus = ref('1')

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '渠道名称', dataIndex: 'channel_name' },
  { title: '渠道类型', slotName: 'type', width: 120 },
  { title: 'AppID', dataIndex: 'app_id', width: 180 },
  { title: '会员数', slotName: 'memberCount', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const getTypeColor = (t) => ({ wechat: 'green', miniprogram: 'cyan', app: 'blue', web: 'purple', offline: 'orange', other: 'gray' }[t] || 'gray')
const getTypeText = (t) => ({ wechat: '微信公众号', miniprogram: '小程序', app: 'APP', web: '官网', offline: '线下', other: '其他' }[t] || t)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.type) params.append('type', filters.type)
    const res = await fetch(`${API_BASE}/member/channels?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.enabled = data.value.filter(d => d.status === 1).length
    }
  } catch (e) { Message.error('加载渠道失败') } finally { loading.value = false }
}

const openCreate = () => { isEdit.value = false; currentId.value = null; Object.assign(form, { channel_name: '', channel_type: 'wechat', app_id: '', app_secret: '', remark: '' }); formStatus.value = '1'; modalVisible.value = true }
const openEdit = (r) => { isEdit.value = true; currentId.value = r.id; Object.assign(form, r); formStatus.value = String(r.status || 1); modalVisible.value = true }
const openDetail = (r) => { current.value = r; detailVisible.value = true }

const handleSubmit = async () => {
  if (!form.channel_name) { Message.warning('请填写渠道名称'); return }
  try {
    const token = localStorage.getItem('token')
    form.status = parseInt(formStatus.value)
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/channels/${currentId.value}` : `${API_BASE}/member/channels`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const d = await res.json()
    if (d.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadData() }
    else Message.error(d.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除渠道「${r.channel_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/channels/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) { Message.success('删除成功'); loadData() } else Message.error(d.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => loadData())
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
