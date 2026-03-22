<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>不积分规则</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8">
        <a-card class="stat-card"><a-statistic title="规则总数" :value="stats.total" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card"><a-statistic title="已启用" :value="stats.enabled" :value-style="{ color: '#52c41a' }" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card class="stat-card"><a-statistic title="本月触发" :value="stats.triggerCount" /></a-card>
      </a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索规则名称" style="width: 240px" search-button @search="loadRules" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadRules">
          <a-option :value="1">启用</a-option>
          <a-option :value="0">禁用</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">新建规则</a-button>
        <a-button @click="loadRules">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="rules" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #scope="{ record }"><a-tag>{{ getScopeText(record.scope) }}</a-tag></template>
        <template #status="{ record }"><a-tag :color="record.status === 1 ? 'green' : 'gray'">{{ record.status === 1 ? '启用' : '禁用' }}</a-tag></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑规则' : '新建规则'" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required><a-input v-model="form.rule_name" placeholder="请输入规则名称" /></a-form-item>
        <a-form-item label="规则类型" required>
          <a-select v-model="form.scope" placeholder="选择规则类型">
            <a-option value="category">指定商品分类</a-option>
            <a-option value="product">指定商品</a-option>
            <a-option value="store">指定门店</a-option>
            <a-option value="payment">指定支付方式</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="关联内容"><a-input v-model="form.target_name" placeholder="请输入分类/商品/门店名称" /></a-form-item>
        <a-form-item label="规则说明"><a-textarea v-model="form.remark" :rows="2" placeholder="请输入规则说明" /></a-form-item>
        <a-form-item label="状态"><a-switch v-model="formStatus" checked-value="1" unchecked-value="0" /></a-form-item>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="detailVisible" title="规则详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="规则名称">{{ current.rule_name }}</a-descriptions-item>
          <a-descriptions-item label="规则类型">{{ getScopeText(current.scope) }}</a-descriptions-item>
          <a-descriptions-item label="关联内容">{{ current.target_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="current.status === 1 ? 'green' : 'gray'">{{ current.status === 1 ? '启用' : '禁用' }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="规则说明">{{ current.remark || '-' }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const rules = ref([])
const loading = ref(false)
const modalVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const current = ref(null)
const currentId = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, enabled: 0, triggerCount: 0 })
const form = reactive({ rule_name: '', scope: 'category', target_name: '', remark: '' })
const formStatus = ref('1')

const columns = [
  { title: '规则名称', dataIndex: 'rule_name' },
  { title: '规则类型', slotName: 'scope', width: 130 },
  { title: '关联内容', dataIndex: 'target_name', ellipsis: true },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const getScopeText = (s) => ({ category: '商品分类', product: '指定商品', store: '门店', payment: '支付方式' }[s] || s)

const loadRules = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status !== '') params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/points/exclude?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) { rules.value = data.data?.list || data.data || []; pagination.total = data.data?.total || 0 }
  } catch (e) { Message.error('加载规则失败') } finally { loading.value = false }
}

const openCreate = () => { isEdit.value = false; currentId.value = null; Object.assign(form, { rule_name: '', scope: 'category', target_name: '', remark: '' }); formStatus.value = '1'; modalVisible.value = true }
const openEdit = (r) => { isEdit.value = true; currentId.value = r.id; Object.assign(form, r); formStatus.value = String(r.status || 1); modalVisible.value = true }
const openDetail = (r) => { current.value = r; detailVisible.value = true }

const handleSubmit = async () => {
  if (!form.rule_name) { Message.warning('请填写规则名称'); return }
  try {
    const token = localStorage.getItem('token')
    form.status = parseInt(formStatus.value)
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/points/exclude/${currentId.value}` : `${API_BASE}/member/points/exclude`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const data = await res.json()
    if (data.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建成功'); modalVisible.value = false; loadRules() }
    else Message.error(data.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除规则「${r.rule_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/points/exclude/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) { Message.success('删除成功'); loadRules() } else Message.error(data.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const onPageChange = (page) => { pagination.current = page; loadRules() }
onMounted(() => loadRules())
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
