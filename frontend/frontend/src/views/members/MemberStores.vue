<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="门店名称">
          <a-input v-model="form.keyword" placeholder="搜索门店名称/编号" style="width: 200px" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option value="1">营业中</a-option>
            <a-option value="0">已关闭</a-option>
          </a-select>
        </a-form-item>
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
      <template #status="{ record }"><a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '营业' : '关闭' }}</a-tag></template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="门店名称" required><a-input v-model="form.store_name" placeholder="请输入门店名称" /></a-form-item>
        <a-form-item label="门店编号"><a-input v-model="form.store_code" placeholder="请输入门店编号" /></a-form-item>
        <a-form-item label="门店地址"><a-input v-model="form.address" placeholder="请输入门店地址" /></a-form-item>
        <a-form-item label="联系电话"><a-input v-model="form.phone" placeholder="请输入联系电话" /></a-form-item>
        <a-form-item label="营业时间"><a-input v-model="form.business_hours" placeholder="如 09:00-21:00" /></a-form-item>
        <a-form-item label="状态"><a-switch v-model="formStatus" checked-value="1" unchecked-value="0" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const data = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({ keyword: '', status: '', store_name: '', store_code: '', address: '', phone: '', business_hours: '' })
const formStatus = ref('1')
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const modalTitle = computed(() => isEdit.value ? '编辑门店' : '新建门店')

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '门店名称', dataIndex: 'store_name' },
  { title: '门店编号', dataIndex: 'store_code', width: 120 },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.keyword) params.append('keyword', form.keyword)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`${API_BASE}/member/stores?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) { data.value = resp.data?.list || resp.data || []; pagination.total = resp.data?.total || 0 }
  } catch (e) { Message.error('加载门店失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.keyword = ''; form.status = ''; handleSearch() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { store_name: '', store_code: '', address: '', phone: '', business_hours: '' })
  formStatus.value = '1'
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, record)
  formStatus.value = String(record.status || 1)
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!form.store_name) { Message.warning('请填写门店名称'); return }
  try {
    const token = localStorage.getItem('token')
    form.status = parseInt(formStatus.value)
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/stores/${currentId.value}` : `${API_BASE}/member/stores`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const resp = await res.json()
    if (resp.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadData() }
    else Message.error(resp.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (record) => {
  if (!confirm(`确定删除门店「${record.store_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) { Message.success('删除成功'); loadData() } else Message.error(resp.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

onMounted(() => loadData())
</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
