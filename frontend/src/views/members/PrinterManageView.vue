<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>云打印设备管理</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="设备总数" :value="stats.total" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="在线" :value="stats.online" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="本月打印" :value="stats.printCount" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索设备名称/编号" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="online">在线</a-option>
          <a-option value="offline">离线</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreate">创建设备</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }"><a-tag :color="record.status === 'online' ? 'green' : 'gray'">{{ record.status === 'online' ? '在线' : '离线' }}</a-tag></template>
        <template #onlineStatus="{ record }"><a-tag :color="record.is_online ? 'green' : 'red'">{{ record.is_online ? '已连接' : '未连接' }}</a-tag></template>
        <template #lastUsed="{ record }">{{ formatTime(record.last_used_at) }}</template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="text" size="small" @click="testPrint(record)">测试打印</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="modalVisible" :title="isEdit ? '编辑设备' : '创建设备'" @ok="handleSubmit" :width="520" :mask-closable="false">
      <a-form :model="form" layout="vertical">
        <a-form-item label="设备名称" required><a-input v-model="form.printer_name" placeholder="请输入设备名称" /></a-form-item>
        <a-form-item label="设备编号"><a-input v-model="form.printer_code" placeholder="请输入设备编号" /></a-form-item>
        <a-form-item label="打印机品牌">
          <a-select v-model="form.brand" placeholder="选择打印机品牌">
            <a-option value="feie">飞鹅</a-option>
            <a-option value="yilian">易联</a-option>
            <a-option value="365">365云打印</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="设备ID"><a-input v-model="form.device_id" placeholder="云打印机设备ID" /></a-form-item>
        <a-form-item label="API密钥"><a-input v-model="form.api_key" placeholder="云打印机API密钥" type="password" /></a-form-item>
        <a-form-item label="关联门店">
          <a-select v-model="form.store_id" placeholder="选择关联门店" allow-clear>
            <a-option v-for="s in stores" :key="s.id" :value="s.id">{{ s.store_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注"><a-textarea v-model="form.remark" :rows="2" placeholder="请输入备注" /></a-form-item>
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
const isEdit = ref(false)
const currentId = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, online: 0, printCount: 0 })
const form = reactive({ printer_name: '', printer_code: '', brand: 'feie', device_id: '', api_key: '', store_id: undefined, remark: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '设备名称', dataIndex: 'printer_name' },
  { title: '设备编号', dataIndex: 'printer_code', width: 120 },
  { title: '品牌', dataIndex: 'brand', width: 100 },
  { title: '关联门店', dataIndex: 'store_name', width: 120 },
  { title: '最近使用', slotName: 'lastUsed', width: 170 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 200 }
]

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/printers?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.online = data.value.filter(d => d.status === 'online').length
    }
  } catch (e) { Message.error('加载设备失败') } finally { loading.value = false }
}

const loadStores = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) stores.value = d.data?.list || []
  } catch (e) {}
}

const openCreate = () => { isEdit.value = false; currentId.value = null; Object.assign(form, { printer_name: '', printer_code: '', brand: 'feie', device_id: '', api_key: '', store_id: undefined, remark: '' }); modalVisible.value = true }
const openEdit = (r) => { isEdit.value = true; currentId.value = r.id; Object.assign(form, r); modalVisible.value = true }

const handleSubmit = async () => {
  if (!form.printer_name) { Message.warning('请填写设备名称'); return }
  try {
    const token = localStorage.getItem('token')
    const method = isEdit.value ? 'PUT' : 'POST'
    const url = isEdit.value ? `${API_BASE}/member/printers/${currentId.value}` : `${API_BASE}/member/printers`
    const res = await fetch(url, { method, headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(form) })
    const d = await res.json()
    if (d.code === 0) { Message.success(isEdit.value ? '更新成功' : '创建设备成功'); modalVisible.value = false; loadData() }
    else Message.error(d.message || '操作失败')
  } catch (e) { Message.error('操作失败') }
}

const handleDelete = async (r) => {
  if (!confirm(`确定删除设备「${r.printer_name}」吗？`)) return
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/printers/${r.id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) { Message.success('删除成功'); loadData() } else Message.error(d.message || '删除失败')
  } catch (e) { Message.error('删除失败') }
}

const testPrint = (r) => { Message.info(`正在向 ${r.printer_name} 发送测试打印...`) }

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
