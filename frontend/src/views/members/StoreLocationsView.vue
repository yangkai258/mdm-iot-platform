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
        <a-input-search v-model="filters.keyword" placeholder="搜索门店名称" style="width: 240px" search-button @search="loadData" />
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">配置位置</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #hasLocation="{ record }"><a-tag :color="record.has_location ? 'green' : 'orange'">{{ record.has_location ? '已配置' : '未配置' }}</a-tag></template>
        <template #coordinate="{ record }"><span>{{ record.longitude && record.latitude ? record.longitude + ', ' + record.latitude : '-' }}</span></template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openMap(record)" v-if="record.has_location">查看地图</a-button>
            <a-button type="text" size="small" @click="openEdit(record)">{{ record.has_location ? '编辑' : '配置' }}位置</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)" v-if="record.has_location">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" title="配置门店位置" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="选择门店" required>
          <a-select v-model="form.store_id" placeholder="选择门店" filterable :disabled="!!currentId">
            <a-option v-for="s in stores" :key="s.id" :value="s.id">{{ s.store_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="经度"><a-input-number v-model="form.longitude" :min="-180" :max="180" :step="0.0001" style="width: 100%" placeholder="如: 116.397428" /></a-form-item>
        <a-form-item label="纬度"><a-input-number v-model="form.latitude" :min="-90" :max="90" :step="0.0001" style="width: 100%" placeholder="如: 39.90923" /></a-form-item>
        <a-form-item label="详细地址"><a-textarea v-model="form.address" :rows="2" placeholder="请输入详细地址" /></a-form-item>
        <a-form-item label="定位说明"><a-input v-model="form.location_remark" placeholder="如: 门店门口左手边" /></a-form-item>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="mapVisible" title="门店位置" :width="600">
      <template v-if="mapRecord">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="门店名称">{{ mapRecord.store_name }}</a-descriptions-item>
          <a-descriptions-item label="详细地址">{{ mapRecord.address || '-' }}</a-descriptions-item>
          <a-descriptions-item label="经度">{{ mapRecord.longitude }}</a-descriptions-item>
          <a-descriptions-item label="纬度">{{ mapRecord.latitude }}</a-descriptions-item>
        </a-descriptions>
        <div style="margin-top: 16px; text-align: center; color: #999; border: 1px dashed #ddd; padding: 40px; border-radius: 8px;">
          <div>地图展示区域</div>
          <div style="font-size: 12px; margin-top: 8px;">经度: {{ mapRecord.longitude }}, 纬度: {{ mapRecord.latitude }}</div>
          <a-button type="primary" style="margin-top: 12px;" @click="openGaodeMap(mapRecord)">在高德地图中打开</a-button>
        </div>
      </template>
    </a-drawer>
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
const mapVisible = ref(false)
const mapRecord = ref(null)
const currentId = ref(null)

const filters = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ configured: 0, unconfigured: 0 })
const form = reactive({ store_id: undefined, longitude: 0, latitude: 0, address: '', location_remark: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '门店名称', dataIndex: 'store_name' },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '经纬度', slotName: 'coordinate', width: 180 },
  { title: '状态', slotName: 'hasLocation', width: 100 },
  { title: '操作', slotName: 'actions', width: 220 }
]

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    const res = await fetch(`${API_BASE}/member/store-locations?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.configured = data.value.filter(d => d.has_location).length
      stats.unconfigured = data.value.filter(d => !d.has_location).length
    }
  } catch (e) { Message.error('加载位置失败') } finally { loading.value = false }
}

const loadStores = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) stores.value = d.data?.list || []
  } catch (e) {}
}

const openCreate = () => { currentId.value = null; Object.assign(form, { store_id: undefined, longitude: 0, latitude: 0, address: '', location_remark: '' }); modalVisible.value = true }
const openEdit = (r) => { currentId.value = r.store_id; Object.assign(form, r); modalVisible.value = true }
const openMap = (r) => { mapRecord.value = r; mapVisible.value = true }
const openGaodeMap = (r) => { window.open(`https://lbs.amap.com/regeo/?lng=${r.longitude}&lat=${r.latitude}`, '_blank') }

const handleSubmit = async () => {
  if (!form.store_id) { Message.warning('请选择门店'); return }
  try {
    const token = localStorage.getItem('token')
    const method = currentId.value ? 'PUT' : 'POST'
    const url = `${API_BASE}/member/store-locations`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const d = await res.json()
    if (d.code === 0) { Message.success('保存成功'); modalVisible.value = false; loadData() }
    else Message.error(d.message || '保存失败')
  } catch (e) { Message.error('保存失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除「${r.store_name}」的位置信息吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/store-locations/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) { Message.success('删除成功'); loadData() } else Message.error(d.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
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
