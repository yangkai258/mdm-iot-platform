<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">搜索</a-button>
          <a-button @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </div>
    <div class="toolbar">
      <a-button type="primary" @click="handleCreate">新建</a-button>
    </div>
    <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
      <template #actions="{ record }">
        <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
        <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @before-ok="handleSubmit" @cancel="modalVisible = false">
      <a-form :model="form" label-col-flex="100px">
        <a-form-item label="名称"><a-input v-model="form.name" placeholder="请输入" /></a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>

import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const data = ref([])
const coupons = ref([])
const tempMembers = ref([])
const stores = ref([])
const loading = ref(false)
const openIssueModal = ref(false)
const detailVisible = ref(false)
const current = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, issued: 0, used: 0 })
const issueForm = reactive({ coupon_id: undefined, scope: 'single', temp_member_id: undefined, store_id: undefined, quantity: 1, remark: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '优惠券名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值/折扣', slotName: 'value', width: 100 },
  { title: '发放方式', dataIndex: 'grant_type', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '发放时间', dataIndex: 'issued_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getTypeText = (t) => ({ discount: '折扣券', cash: '现金券', gift: '礼品券' }[t] || t)
const getStatusColor = (s) => ({ active: 'green', issued: 'blue', used: 'purple', expired: 'gray' }[s] || 'gray')
const getStatusText = (s) => ({ active: '可发放', issued: '已发放', used: '已使用', expired: '已过期' }[s] || s)
const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/temp-coupons?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.issued = data.value.filter(d => d.status === 'issued').length
      stats.used = data.value.filter(d => d.status === 'used').length
    }
  } catch (e) { Message.error('加载优惠券失败') } finally { loading.value = false }
}

const loadCoupons = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/coupons?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) coupons.value = d.data?.list || []
  } catch (e) {}
}

const loadTempMembers = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-members?page_size=100&status=pending`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) tempMembers.value = d.data?.list || []
  } catch (e) {}
}

const loadStores = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) stores.value = d.data?.list || []
  } catch (e) {}
}

const viewDetail = (r) => { current.value = r; detailVisible.value = true }
const grantCoupon = (r) => { issueForm.coupon_id = r.id; issueForm.scope = 'single'; openIssueModal.value = true }

const handleIssue = async () => {
  if (!issueForm.coupon_id || (issueForm.scope === 'single' && !issueForm.temp_member_id)) {
    Message.warning('请选择优惠券和发放对象'); return
  }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-coupons/issue`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(issueForm)
    })
    const d = await res.json()
    if (d.code === 0) { Message.success('发放成功'); openIssueModal.value = false; loadData() }
    else Message.error(d.message || '发放失败')
  } catch (e) { Message.error('发放失败') }
}

const handleExport = () => { Message.info('正在导出...') }

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadCoupons(); loadTempMembers(); loadStores() })

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
