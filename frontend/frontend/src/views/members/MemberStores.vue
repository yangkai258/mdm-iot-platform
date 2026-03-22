<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>店铺信息</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="门店总数" :value="stats.total" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="营业中" :value="stats.open" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="会员总数" :value="stats.memberCount" /></a-card></a-col>
      <a-col :span="6"><a-card class="stat-card"><a-statistic title="本月订单" :value="stats.orderCount" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索门店名称/编号" style="width: 240px" search-button @search="loadStores" />
        <a-select v-model="filters.status" placeholder="门店状态" allow-clear style="width: 120px" @change="loadStores">
          <a-option value="1">营业中</a-option>
          <a-option value="0">已关闭</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">新建门店</a-button>
        <a-button @click="loadStores">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="stores" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }"><a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '营业' : '关闭' }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑门店' : '新建门店'" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="门店名称" required><a-input v-model="form.store_name" placeholder="请输入门店名称" /></a-form-item>
        <a-form-item label="门店编号"><a-input v-model="form.store_code" placeholder="请输入门店编号" /></a-form-item>
        <a-form-item label="门店地址"><a-input v-model="form.address" placeholder="请输入门店地址" /></a-form-item>
        <a-form-item label="联系电话"><a-input v-model="form.phone" placeholder="请输入联系电话" /></a-form-item>
        <a-form-item label="营业时间"><a-input v-model="form.business_hours" placeholder="如 09:00-21:00" /></a-form-item>
        <a-form-item label="状态"><a-switch v-model="formStatus" checked-value="1" unchecked-value="0" /></a-form-item>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="detailVisible" title="门店详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="门店名称">{{ current.store_name }}</a-descriptions-item>
          <a-descriptions-item label="门店编号">{{ current.store_code || '-' }}</a-descriptions-item>
          <a-descriptions-item label="门店地址">{{ current.address || '-' }}</a-descriptions-item>
          <a-descriptions-item label="联系电话">{{ current.phone || '-' }}</a-descriptions-item>
          <a-descriptions-item label="营业时间">{{ current.business_hours || '-' }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="current.status === 1 ? 'green' : 'gray'">{{ current.status === 1 ? '营业中' : '已关闭' }}</a-tag></a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const stores = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const current = ref(null)
const currentId = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, open: 0, memberCount: 0, orderCount: 0 })
const form = reactive({ store_name: '', store_code: '', address: '', phone: '', business_hours: '' })
const formStatus = ref('1')

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '门店名称', dataIndex: 'store_name' },
  { title: '门店编号', dataIndex: 'store_code', width: 120 },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const loadStores = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/stores?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) { stores.value = data.data?.list || data.data || []; pagination.total = data.data?.total || 0 }
  } catch (e) { Message.error('加载门店失败') } finally { loading.value = false }
}

const openCreate = () => { isEdit.value = false; currentId.value = null; Object.assign(form, { store_name: '', store_code: '', address: '', phone: '', business_hours: '' }); formStatus.value = '1'; modalVisible.value = true }
const openEdit = (r) => { isEdit.value = true; currentId.value = r.id; Object.assign(form, r); formStatus.value = String(r.status || 1); modalVisible.value = true }
const openDetail = (r) => { current.value = r; detailVisible.value = true }

const handleSubmit = async () => {
  if (!form.store_name) { Message.warning('请填写门店名称'); return }
  try {
    const token = localStorage.getItem('token')
    form.status = parseInt(formStatus.value)
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/stores/${currentId.value}` : `${API_BASE}/member/stores`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const data = await res.json()
    if (data.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadStores() }
    else Message.error(data.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除门店「${r.store_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) { Message.success('删除成功'); loadStores() } else Message.error(data.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadStores() }
onMounted(() => loadStores())
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
