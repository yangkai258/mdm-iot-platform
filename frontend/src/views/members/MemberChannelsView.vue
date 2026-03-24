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
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
