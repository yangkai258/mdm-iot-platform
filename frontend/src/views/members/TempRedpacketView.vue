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
const tempMembers = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const grantModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const grantRecord = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, active: 0, issuedAmount: 0 })
const form = reactive({ redpacket_name: '', amount: 0, total_count: 0, min_amount: 0, dateRange: [], remark: '' })
const grantForm = reactive({ scope: 'single', temp_member_id: undefined, count: 1 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '红包名称', dataIndex: 'redpacket_name' },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '剩余数量', slotName: 'remain', width: 120 },
  { title: '有效期', dataIndex: 'expire_time', width: 170 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const getStatusColor = (s) => ({ active: 'green', paused: 'orange', finished: 'gray' }[s] || 'gray')
const getStatusText = (s) => ({ active: '进行中', paused: '已暂停', finished: '已结束' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/temp-redpackets?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.active = data.value.filter(d => d.status === 'active').length
      stats.issuedAmount = data.value.reduce((sum, d) => sum + (d.total_count - d.remain_count) * (d.amount || 0), 0)
    }
  } catch (e) { Message.error('加载红包失败') } finally { loading.value = false }
}

const loadTempMembers = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-members?page_size=100&status=pending`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) tempMembers.value = d.data?.list || []
  } catch (e) {}
}

const openCreate = () => { isEdit.value = false; currentId.value = null; Object.assign(form, { redpacket_name: '', amount: 0, total_count: 0, min_amount: 0, dateRange: [], remark: '' }); modalVisible.value = true }
const openEdit = (r) => { isEdit.value = true; currentId.value = r.id; Object.assign(form, r); modalVisible.value = true }
const openGrantModal = (r) => { grantRecord.value = r; grantForm.scope = 'single'; grantForm.temp_member_id = undefined; grantForm.count = 1; grantModal.value = true }

const handleSubmit = async () => {
  if (!form.redpacket_name || !form.amount || !form.total_count) { Message.warning('请填写完整信息'); return }
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/temp-redpackets/${currentId.value}` : `${API_BASE}/member/temp-redpackets`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const d = await res.json()
    if (d.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadData() }
    else Message.error(d.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleGrant = async () => {
  if (grantForm.scope === 'single' && !grantForm.temp_member_id) { Message.warning('请选择临时会员'); return }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-redpackets/${grantRecord.value.id}/grant`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(grantForm)
    })
    const d = await res.json()
    if (d.code === 0) { Message.success('发放成功'); grantModal.value = false; loadData() }
    else Message.error(d.message || '发放失败')
  } catch (e) { Message.error('发放失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除红包「${r.redpacket_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-redpackets/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) { Message.success('删除成功'); loadData() } else Message.error(d.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadTempMembers() })

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
