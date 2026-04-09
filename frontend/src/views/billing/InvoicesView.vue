<template>
  <div class="page-container">
    <a-card class="general-card" title="发票管理">
      <template #extra>
        <a-button type="primary" @click="openApplyModal"><icon-plus />开票申请</a-button>
      </template>
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="list" title="发票列表">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="发票状态">
                <a-select v-model="form.status" placeholder="选择状态" allow-clear style="width: 140px">
                  <a-option value="pending">待开票</a-option>
                  <a-option value="issued">已开票</a-option>
                  <a-option value="delivered">已送达</a-option>
                  <a-option value="cancelled">已作废</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="时间范围"><a-range-picker v-model="form.time_range" style="width: 240px" /></a-form-item>
              <a-form-item><a-button type="primary" @click="loadData">查询</a-button><a-button @click="Object.assign(form, { status: '', time_range: [] }); loadData()">重置</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #amount="{ record }">¥{{ record.amount?.toFixed(2) }}</template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="viewInvoice(record)">详情</a-button>
              <a-button v-if="record.status === 'issued'" type="text" size="small" @click="downloadInvoice(record)">下载</a-button>
              <a-button v-if="record.status === 'pending'" type="text" size="small" status="danger" @click="cancelInvoice(record)">取消</a-button>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="history" title="开票历史">
          <a-table :columns="historyColumns" :data="historyData" :loading="historyLoading" :pagination="historyPagination" @page-change="onHistoryPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #amount="{ record }">¥{{ record.amount?.toFixed(2) }}</template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="viewInvoice(record)">详情</a-button>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    <!-- 开票申请弹窗 -->
    <a-modal v-model:visible="applyVisible" title="开票申请" @before-ok="handleApply" :loading="submitting" :width="600">
      <a-form :model="applyForm" layout="vertical">
        <a-form-item label="发票类型" required>
          <a-radio-group v-model="applyForm.invoice_type">
            <a-radio value="normal">普通发票</a-radio>
            <a-radio value="vat">增值税专用发票</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="发票抬头" required><a-input v-model="applyForm.title" placeholder="请输入发票抬头" /></a-form-item>
        <a-form-item label="税号" required><a-input v-model="applyForm.tax_no" placeholder="请输入税号" /></a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="开户银行"><a-input v-model="applyForm.bank" placeholder="请输入开户银行" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="银行账号"><a-input v-model="applyForm.bank_account" placeholder="请输入银行账号" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="注册地址"><a-input v-model="applyForm.address" placeholder="请输入注册地址" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系电话"><a-input v-model="applyForm.phone" placeholder="请输入联系电话" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="开票金额" required>
          <a-input-number v-model="applyForm.amount" :min="0" :precision="2" placeholder="请输入开票金额" style="width: 200px" />
          <span style="margin-left: 8px">元</span>
        </a-form-item>
        <a-form-item label="接收邮箱"><a-input v-model="applyForm.email" placeholder="请输入接收邮箱" /></a-form-item>
        <a-form-item label="备注"><a-textarea v-model="applyForm.notes" :rows="2" placeholder="备注信息" /></a-form-item>
      </a-form>
    </a-modal>
    <!-- 发票详情 -->
    <a-modal v-model:visible="detailVisible" title="发票详情" :footer="null" :width="600">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="发票号" :span="2">{{ currentInvoice?.invoice_no }}</a-descriptions-item>
        <a-descriptions-item label="发票类型">{{ currentInvoice?.invoice_type === 'normal' ? '普通发票' : '增值税专用发票' }}</a-descriptions-item>
        <a-descriptions-item label="发票金额">¥{{ currentInvoice?.amount?.toFixed(2) }}</a-descriptions-item>
        <a-descriptions-item label="发票抬头" :span="2">{{ currentInvoice?.title }}</a-descriptions-item>
        <a-descriptions-item label="税号" :span="2">{{ currentInvoice?.tax_no }}</a-descriptions-item>
        <a-descriptions-item label="状态"><a-tag :color="getStatusColor(currentInvoice?.status)">{{ getStatusText(currentInvoice?.status) }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="申请时间">{{ currentInvoice?.created_at }}</a-descriptions-item>
        <a-descriptions-item label="开票时间" :span="2">{{ currentInvoice?.issued_at || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const activeTab = ref('list')
const loading = ref(false)
const historyLoading = ref(false)
const submitting = ref(false)
const data = ref([])
const historyData = ref([])
const applyVisible = ref(false)
const detailVisible = ref(false)
const currentInvoice = ref(null)
const form = reactive({ status: '', time_range: [] })
const applyForm = reactive({ invoice_type: 'normal', title: '', tax_no: '', bank: '', bank_account: '', address: '', phone: '', amount: null, email: '', notes: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const historyPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '发票号', dataIndex: 'invoice_no', width: 140 },
  { title: '发票抬头', dataIndex: 'title', ellipsis: true },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '发票类型', dataIndex: 'invoice_type', width: 120 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '申请时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '发票号', dataIndex: 'invoice_no', width: 140 },
  { title: '发票抬头', dataIndex: 'title', ellipsis: true },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '开票时间', dataIndex: 'issued_at', width: 170 }
]

const getStatusColor = (s) => ({ pending: 'orange', issued: 'arcoblue', delivered: 'green', cancelled: 'red' }[s] || 'gray')
const getStatusText = (s) => ({ pending: '待开票', issued: '已开票', delivered: '已送达', cancelled: '已作废' }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/billing/invoices?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('加载失败') } finally { loading.value = false }
}

const loadHistory = async () => {
  historyLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: historyPagination.current, page_size: historyPagination.pageSize, status: 'issued,delivered' })
    const res = await fetch(`/api/v1/billing/invoices?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { historyData.value = res.data?.list || []; historyPagination.total = res.data?.total || 0 }
    else { historyData.value = [] }
  } catch (e) { Message.error('加载历史失败') } finally { historyLoading.value = false }
}

const openApplyModal = () => { Object.assign(applyForm, { invoice_type: 'normal', title: '', tax_no: '', bank: '', bank_account: '', address: '', phone: '', amount: null, email: '', notes: '' }); applyVisible.value = true }

const handleApply = async (done) => {
  if (!applyForm.title || !applyForm.tax_no || !applyForm.amount) { Message.warning('请填写必填项'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/billing/invoices', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(applyForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('申请已提交'); applyVisible.value = false; loadData() }
    else { Message.error(res.message || '申请失败') }
    done(true)
  } catch (e) { Message.error('申请失败'); done(false) } finally { submitting.value = false }
}

const viewInvoice = (record) => { currentInvoice.value = record; detailVisible.value = true }

const downloadInvoice = (record) => {
  window.open(`/api/v1/billing/invoices/${record.id}/download`, '_blank')
}

const cancelInvoice = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/billing/invoices/${record.id}/cancel`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('发票已取消'); loadData() }
    else Message.error('取消失败')
  } catch (e) { Message.error('取消失败') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onHistoryPageChange = (page) => { historyPagination.current = page; loadHistory() }

onMounted(() => { loadData(); loadHistory() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
