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
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
