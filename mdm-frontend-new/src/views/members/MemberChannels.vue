<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="渠道名称">
          <a-input v-model="form.keyword" placeholder="搜索渠道名称" style="width: 200px" />
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
      <template #type="{ record }"><a-tag>{{ getTypeText(record.channel_type) }}</a-tag></template>
      <template #status="{ record }"><a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag></template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="480" :mask-closable="false">
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

const form = reactive({ keyword: '', channel_name: '', channel_type: 'wechat', app_id: '', app_secret: '', remark: '' })
const formStatus = ref('1')
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const modalTitle = computed(() => isEdit.value ? '编辑渠道' : '新建渠道')

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '渠道名称', dataIndex: 'channel_name' },
  { title: '渠道类型', slotName: 'type', width: 120 },
  { title: 'AppID', dataIndex: 'app_id', width: 180 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getTypeText = (t) => ({ wechat: '微信公众号', miniprogram: '小程序', app: 'APP', web: '官网', offline: '线下', other: '其他' }[t] || t)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.keyword) params.append('keyword', form.keyword)
    const res = await fetch(`${API_BASE}/member/channels?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) { data.value = resp.data?.list || resp.data || []; pagination.total = resp.data?.total || 0 }
  } catch (e) { Message.error('加载渠道失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.keyword = ''; handleSearch() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { channel_name: '', channel_type: 'wechat', app_id: '', app_secret: '', remark: '' })
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
  if (!form.channel_name) { Message.warning('请填写渠道名称'); return }
  try {
    const token = localStorage.getItem('token')
    form.status = parseInt(formStatus.value)
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/channels/${currentId.value}` : `${API_BASE}/member/channels`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const resp = await res.json()
    if (resp.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadData() }
    else Message.error(resp.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (record) => {
  if (!confirm(`确定删除渠道「${record.channel_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/channels/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
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
