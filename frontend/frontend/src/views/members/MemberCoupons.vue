<template>
  <div class="page-container">
    <div class="search-form">
      <a-form :model="form" layout="inline">
        <a-form-item label="优惠券名称">
          <a-input v-model="form.keyword" placeholder="搜索优惠券名称" style="width: 200px" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option value="active">有效</a-option>
            <a-option value="expired">已过期</a-option>
            <a-option value="used">已使用</a-option>
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
      <template #type="{ record }"><a-tag>{{ getTypeText(record.type) }}</a-tag></template>
      <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag></template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="text" size="small" @click="handleDelete(record)">删除</a-button>
        </a-space>
      </template>
    </a-table>
    <a-modal v-model:visible="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="480" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="优惠券名称" required><a-input v-model="form.name" placeholder="请输入优惠券名称" /></a-form-item>
        <a-form-item label="优惠券类型" required>
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="discount">折扣券</a-option>
            <a-option value="cash">现金券</a-option>
            <a-option value="gift">礼品券</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="面值/折扣"><a-input-number v-model="form.value" :min="0" style="width: 100%" /></a-form-item>
        <a-form-item label="使用门槛"><a-input-number v-model="form.min_amount" :min="0" style="width: 100%" /></a-form-item>
        <a-form-item label="有效期"><a-input v-model="form.expire_time" placeholder="如: 2024-12-31" /></a-form-item>
        <a-form-item label="备注"><a-textarea v-model="form.remark" :rows="2" placeholder="请输入备注" /></a-form-item>
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

const form = reactive({ keyword: '', status: '', name: '', type: 'discount', value: 0, min_amount: 0, expire_time: '', remark: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const modalTitle = computed(() => isEdit.value ? '编辑优惠券' : '新建优惠券')

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '优惠券名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值', dataIndex: 'value', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '有效期', dataIndex: 'expire_time', width: 160 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getTypeText = (type) => ({ discount: '折扣券', cash: '现金券', gift: '礼品券' }[type] || type)
const getStatusText = (s) => ({ active: '有效', expired: '已过期', used: '已使用' }[s] || s)
const getStatusColor = (s) => ({ active: 'green', expired: 'gray', used: 'blue' }[s] || 'gray')

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.keyword) params.append('keyword', form.keyword)
    if (form.status) params.append('status', form.status)
    const res = await fetch(`${API_BASE}/member/coupons?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) { data.value = resp.data?.list || resp.data || []; pagination.total = resp.data?.total || 0 }
  } catch (e) { Message.error('加载优惠券失败') } finally { loading.value = false }
}

const handleSearch = () => { pagination.current = 1; loadData() }
const handleReset = () => { form.keyword = ''; form.status = ''; handleSearch() }
const onPageChange = (page) => { pagination.current = page; loadData() }

const handleCreate = () => {
  isEdit.value = false
  currentId.value = null
  Object.assign(form, { name: '', type: 'discount', value: 0, min_amount: 0, expire_time: '', remark: '' })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(form, record)
  modalVisible.value = true
}

const handleSubmit = async () => {
  if (!form.name) { Message.warning('请填写优惠券名称'); return }
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/coupons/${currentId.value}` : `${API_BASE}/member/coupons`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const resp = await res.json()
    if (resp.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadData() }
    else Message.error(resp.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (record) => {
  if (!confirm(`确定删除优惠券「${record.name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/coupons/${record.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
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
