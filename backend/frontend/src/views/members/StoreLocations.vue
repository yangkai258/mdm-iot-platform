<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>门店位置管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="已配置位置" :value="stats.configured" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="未配置" :value="stats.unconfigured" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索门店名称" style="width: 240px" search-button @change="loadData" />
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-button type="primary" @click="openCreate">配置位置</a-button>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #hasLocation="{ record }"><a-tag :color="record.has_location ? 'green' : 'orange'">{{ record.has_location ? '已配置' : '未配置' }}</a-tag></template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="openEdit(record)">{{ record.has_location ? '编辑' : '配置' }}位置</a-button>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" title="配置门店位置" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="选择门店" required>
          <a-select v-model="form.store_id" placeholder="选择门店" filterable>
            <a-option v-for="s in stores" :key="s.id" :value="s.id">{{ s.store_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="经度"><a-input-number v-model="form.longitude" :min="-180" :max="180" :step="0.0001" style="width: 100%" placeholder="如: 116.397428" /></a-form-item>
        <a-form-item label="纬度"><a-input-number v-model="form.latitude" :min="-90" :max="90" :step="0.0001" style="width: 100%" placeholder="如: 39.90923" /></a-form-item>
        <a-form-item label="详细地址"><a-textarea v-model="form.address" :rows="2" placeholder="请输入详细地址" /></a-form-item>
        <a-form-item label="定位说明"><a-input v-model="form.location_remark" placeholder="如: 门店门口左手边" /></a-form-item>
      </a-form>
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
const currentId = ref(null)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ configured: 0, unconfigured: 0 })
const form = reactive({ store_id: undefined, longitude: 0, latitude: 0, address: '', location_remark: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '门店名称', dataIndex: 'store_name' },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '坐标', dataIndex: 'coordinate', width: 160 },
  { title: '状态', slotName: 'hasLocation', width: 100 },
  { title: '操作', slotName: 'actions', width: 120 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    const res = await fetch(`${API_BASE}/member/store-locations?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) { data.value = resp.data?.list || resp.data || []; pagination.total = resp.data?.total || 0 }
  } catch (e) { Message.error('加载位置失败') } finally { loading.value = false }
}

const loadStores = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const data = await res.json()
    if (data.code === 0) stores.value = data.data?.list || []
  } catch (e) {}
}

const openCreate = () => { currentId.value = null; Object.assign(form, { store_id: undefined, longitude: 0, latitude: 0, address: '', location_remark: '' }); modalVisible.value = true }
const openEdit = (r) => { currentId.value = r.store_id; Object.assign(form, r); modalVisible.value = true }

const handleSubmit = async () => {
  if (!form.store_id) { Message.warning('请选择门店'); return }
  try {
    const token = localStorage.getItem('token')
    const method = currentId.value ? 'PUT' : 'POST'
    const url = `${API_BASE}/member/store-locations`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const data = await res.json()
    if (data.code === 0) { Message.success('保存成功'); modalVisible.value = false; loadData() }
    else Message.error(data.message || '保存失败')
  } catch (e) { Message.error('保存失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadStores() })
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
