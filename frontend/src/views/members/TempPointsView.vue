<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>临时积分库存</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="积分池总数" :value="stats.total" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="库存充足" :value="stats.sufficient" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="本月已发放" :value="stats.issued" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索积分池名称" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="active">正常</a-option>
          <a-option value="low">库存不足</a-option>
          <a-option value="depleted">已耗尽</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openRechargeModal = true">充值积分</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag></template>
        <template #remain="{ record }">
          <span :style="{ color: record.remain < record.threshold ? '#ff4d4f' : '#52c41a', fontWeight: 600 }">
            {{ record.remain.toLocaleString() }}
          </span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openRechargeModalFor(record)">充值</a-button>
            <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="rechargeModal" title="积分池充值" @ok="handleRecharge" :width="480" :mask-closable="false">
      <a-form :model="rechargeForm" layout="vertical">
        <a-form-item label="积分池" required>
          <a-select v-model="rechargeForm.pool_id" placeholder="选择积分池" filterable :disabled="!!rechargeForm.pool_id">
            <a-option v-for="p in pools" :key="p.id" :value="p.id">{{ p.pool_name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="充值积分数量" required>
          <a-input-number v-model="rechargeForm.points" :min="1" style="width: 100%;" placeholder="请输入充值积分数量" />
        </a-form-item>
        <a-form-item label="充值说明">
          <a-textarea v-model="rechargeForm.remark" :rows="2" placeholder="请输入充值说明" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-space>
          <a-button @click="rechargeModal = false">取消</a-button>
          <a-button type="primary" @click="handleRecharge">确认充值</a-button>
        </a-space>
      </template>
    </a-modal>

    <a-drawer v-model:visible="detailVisible" title="积分池详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="积分池名称">{{ current.pool_name }}</a-descriptions-item>
          <a-descriptions-item label="剩余数量"><span style="font-weight: 600; color: #1890ff;">{{ current.remain.toLocaleString() }}</span></a-descriptions-item>
          <a-descriptions-item label="总库存">{{ current.total.toLocaleString() }}</a-descriptions-item>
          <a-descriptions-item label="预警阈值">{{ current.threshold }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="getStatusColor(current.status)">{{ getStatusText(current.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="本月发放">{{ current.month_issued || 0 }}</a-descriptions-item>
          <a-descriptions-item label="创建时间">{{ formatTime(current.created_at) }}</a-descriptions-item>
        </a-descriptions>
      </template>
    </a-drawer>
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
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
</style>
