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
const pools = ref([])
const loading = ref(false)
const rechargeModal = ref(false)
const detailVisible = ref(false)
const current = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, sufficient: 0, issued: 0 })
const rechargeForm = reactive({ pool_id: undefined, points: 0, remark: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '积分池名称', dataIndex: 'pool_name' },
  { title: '剩余数量', slotName: 'remain', width: 140 },
  { title: '总库存', dataIndex: 'total', width: 120 },
  { title: '预警阈值', dataIndex: 'threshold', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const getStatusColor = (s) => ({ active: 'green', low: 'orange', depleted: 'red' }[s] || 'gray')
const getStatusText = (s) => ({ active: '正常', low: '库存不足', depleted: '已耗尽' }[s] || s)
const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/temp-points?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.sufficient = data.value.filter(d => d.status === 'active').length
    }
  } catch (e) { Message.error('加载积分池失败') } finally { loading.value = false }
}

const loadPools = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-points/pools?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) pools.value = d.data?.list || []
  } catch (e) {}
}

const openRechargeModal = () => { Object.assign(rechargeForm, { pool_id: undefined, points: 0, remark: '' }); rechargeModal.value = true }
const openRechargeModalFor = (r) => { Object.assign(rechargeForm, { pool_id: r.id, points: 0, remark: '' }); rechargeModal.value = true }
const viewDetail = (r) => { current.value = r; detailVisible.value = true }

const handleRecharge = async () => {
  if (!rechargeForm.pool_id || !rechargeForm.points) { Message.warning('请选择积分池并填写充值数量'); return }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-points/recharge`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(rechargeForm)
    })
    const d = await res.json()
    if (d.code === 0) { Message.success('充值成功'); rechargeModal.value = false; loadData() }
    else Message.error(d.message || '充值失败')
  } catch (e) { Message.error('充值失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadPools() })

</script>

<style scoped>
.page-container { background: #fff; border-radius: 4px; padding: 20px; }
.search-form { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
.toolbar { margin-bottom: 16px; }
</style>
