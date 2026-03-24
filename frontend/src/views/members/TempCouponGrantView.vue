<template>
  <div class="pro-page-container">
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>会员管理</a-breadcrumb-item>
      <a-breadcrumb-item>临时优惠券发放</a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16" class="stats-row">
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="优惠券总数" :value="stats.total" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="已发放" :value="stats.issued" :value-style="{ color: '#1890ff' }" /></a-card></a-col>
      <a-col :span="8"><a-card class="stat-card"><a-statistic title="已使用" :value="stats.used" :value-style="{ color: '#52c41a' }" /></a-card></a-col>
    </a-row>

    <div class="pro-search-bar">
      <a-space wrap>
        <a-input-search v-model="filters.keyword" placeholder="搜索优惠券名称" style="width: 240px" search-button @search="loadData" />
        <a-select v-model="filters.status" placeholder="状态" allow-clear style="width: 120px" @change="loadData">
          <a-option value="active">可发放</a-option>
          <a-option value="issued">已发放</a-option>
          <a-option value="used">已使用</a-option>
          <a-option value="expired">已过期</a-option>
        </a-select>
      </a-space>
    </div>

    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openIssueModal = true">发放优惠券</a-button>
        <a-button @click="handleExport">导出</a-button>
        <a-button @click="loadData">刷新</a-button>
      </a-space>
    </div>

    <div class="pro-content-area">
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
        <template #type="{ record }"><a-tag>{{ getTypeText(record.coupon_type) }}</a-tag></template>
        <template #status="{ record }"><a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag></template>
        <template #value="{ record }">{{ record.discount ? record.discount + '折' : record.amount ? '￥' + record.amount : '-' }}</template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="viewDetail(record)">查看详情</a-button>
            <a-button type="text" size="small" @click="grantCoupon(record)" v-if="record.status === 'active'">发放</a-button>
          </a-space>
        </template>
      </a-table>
    </div>

    <a-modal v-model:visible="openIssueModal" title="发放优惠券" :width="600" :mask-closable="false">
      <a-form :model="issueForm" layout="vertical">
        <a-form-item label="选择优惠券" required>
          <a-select v-model="issueForm.coupon_id" placeholder="选择要发放的优惠券" filterable>
            <a-option v-for="c in coupons" :key="c.id" :value="c.id">
              {{ c.name }} ({{ c.coupon_type === 'discount' ? c.discount + '折' : '￥' + c.amount }})
            </a-option>
          </a-select>
        </a-form-item>

        <a-divider>发放方式</a-divider>

        <a-form-item label="发放范围">
          <a-radio-group v-model="issueForm.scope">
            <a-radio value="single">指定临时会员</a-radio>
            <a-radio value="batch">批量发放</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item v-if="issueForm.scope === 'single'" label="选择临时会员" required>
          <a-select v-model="issueForm.temp_member_id" placeholder="选择临时会员" filterable>
            <a-option v-for="m in tempMembers" :key="m.id" :value="m.id">{{ m.name }} ({{ m.phone }})</a-option>
          </a-select>
        </a-form-item>

        <a-form-item v-if="issueForm.scope === 'batch'" label="选择门店">
          <a-select v-model="issueForm.store_id" placeholder="选择门店(不选则发往所有)" allow-clear filterable>
            <a-option v-for="s in stores" :key="s.id" :value="s.id">{{ s.store_name }}</a-option>
          </a-select>
        </a-form-item>

        <a-form-item v-if="issueForm.scope === 'batch'" label="发放数量">
          <a-input-number v-model="issueForm.quantity" :min="1" :max="1000" style="width: 200px;" />
          <span style="margin-left: 8px; color: #999;">将为每个临时会员发放1张</span>
        </a-form-item>

        <a-form-item label="备注">
          <a-textarea v-model="issueForm.remark" :rows="2" placeholder="请输入备注" />
        </a-form-item>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="openIssueModal = false">取消</a-button>
          <a-button type="primary" @click="handleIssue">确认发放</a-button>
        </a-space>
      </template>
    </a-modal>

    <a-drawer v-model:visible="detailVisible" title="优惠券详情" :width="480">
      <template v-if="current">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="优惠券名称">{{ current.name }}</a-descriptions-item>
          <a-descriptions-item label="优惠券类型">{{ getTypeText(current.coupon_type) }}</a-descriptions-item>
          <a-descriptions-item label="面值/折扣">{{ current.discount ? current.discount + '折' : '￥' + current.amount }}</a-descriptions-item>
          <a-descriptions-item label="使用门槛">{{ current.min_amount ? '满￥' + current.min_amount : '无门槛' }}</a-descriptions-item>
          <a-descriptions-item label="有效期">{{ current.start_date }} 至 {{ current.end_date }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="getStatusColor(current.status)">{{ getStatusText(current.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="发放时间">{{ formatTime(current.issued_at) }}</a-descriptions-item>
          <a-descriptions-item label="使用时间">{{ formatTime(current.used_at) || '-' }}</a-descriptions-item>
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
const coupons = ref([])
const tempMembers = ref([])
const stores = ref([])
const loading = ref(false)
const openIssueModal = ref(false)
const detailVisible = ref(false)
const current = ref(null)

const filters = reactive({ keyword: '', status: '' })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const stats = reactive({ total: 0, issued: 0, used: 0 })
const issueForm = reactive({ coupon_id: undefined, scope: 'single', temp_member_id: undefined, store_id: undefined, quantity: 1, remark: '' })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '优惠券名称', dataIndex: 'name' },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '面值/折扣', slotName: 'value', width: 100 },
  { title: '发放方式', dataIndex: 'grant_type', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '发放时间', dataIndex: 'issued_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 150 }
]

const getTypeText = (t) => ({ discount: '折扣券', cash: '现金券', gift: '礼品券' }[t] || t)
const getStatusColor = (s) => ({ active: 'green', issued: 'blue', used: 'purple', expired: 'gray' }[s] || 'gray')
const getStatusText = (s) => ({ active: '可发放', issued: '已发放', used: '已使用', expired: '已过期' }[s] || s)
const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (filters.keyword) params.append('keyword', filters.keyword)
    if (filters.status) params.append('status', filters.status)
    const res = await fetch(`${API_BASE}/member/temp-coupons?${params}`, { headers: { 'Authorization': `Bearer ${token}` } })
    const resp = await res.json()
    if (resp.code === 0) {
      data.value = resp.data?.list || resp.data || []
      pagination.total = resp.data?.total || 0
      stats.total = pagination.total
      stats.issued = data.value.filter(d => d.status === 'issued').length
      stats.used = data.value.filter(d => d.status === 'used').length
    }
  } catch (e) { Message.error('加载优惠券失败') } finally { loading.value = false }
}

const loadCoupons = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/coupons?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) coupons.value = d.data?.list || []
  } catch (e) {}
}

const loadTempMembers = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-members?page_size=100&status=pending`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) tempMembers.value = d.data?.list || []
  } catch (e) {}
}

const loadStores = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/stores?page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
    const d = await res.json()
    if (d.code === 0) stores.value = d.data?.list || []
  } catch (e) {}
}

const viewDetail = (r) => { current.value = r; detailVisible.value = true }
const grantCoupon = (r) => { issueForm.coupon_id = r.id; issueForm.scope = 'single'; openIssueModal.value = true }

const handleIssue = async () => {
  if (!issueForm.coupon_id || (issueForm.scope === 'single' && !issueForm.temp_member_id)) {
    Message.warning('请选择优惠券和发放对象'); return
  }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`${API_BASE}/member/temp-coupons/issue`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(issueForm)
    })
    const d = await res.json()
    if (d.code === 0) { Message.success('发放成功'); openIssueModal.value = false; loadData() }
    else Message.error(d.message || '发放失败')
  } catch (e) { Message.error('发放失败') }
}

const handleExport = () => { Message.info('正在导出...') }

const onPageChange = (page) => { pagination.current = page; loadData() }
onMounted(() => { loadData(); loadCoupons(); loadTempMembers(); loadStores() })
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
