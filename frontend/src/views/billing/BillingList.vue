<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>账单管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- Tab 切换 -->
    <a-tabs v-model:active-tab="activeTab" class="billing-tabs">
      <!-- Tab1: 账单记录 -->
      <a-tab-pane key="records" title="账单记录">
        <!-- 费用汇总卡片 -->
        <a-row :gutter="16" class="summary-row">
          <a-col :span="6">
            <a-card class="summary-card" hoverable>
              <a-statistic title="本月应结" :value="summary.monthly_total" :precision="2" prefix="¥" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card class="summary-card" hoverable>
              <a-statistic title="待支付" :value="summary.pending_amount" :precision="2" prefix="¥" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card class="summary-card" hoverable>
              <a-statistic title="已支付（本年）" :value="summary.paid_this_year" :precision="2" prefix="¥" />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card class="summary-card" hoverable>
              <a-statistic title="开票中" :value="summary.invoice_pending" :precision="2" prefix="¥" />
            </a-card>
          </a-col>
        </a-row>

        <!-- 搜索表单 -->
        <div class="pro-search-bar">
          <a-space wrap>
            <a-input-search v-model="form.user_id" placeholder="用户ID" style="width: 140px" @search="loadRecords" search-button />
            <a-select v-model="form.billing_type" placeholder="计费类型" allow-clear style="width: 160px" @change="loadRecords">
              <a-option value="subscription">订阅费用</a-option>
              <a-option value="usage">用量费用</a-option>
              <a-option value="API_quota">API配额</a-option>
            </a-select>
            <a-select v-model="form.status" placeholder="账单状态" allow-clear style="width: 140px" @change="loadRecords">
              <a-option value="pending">待支付</a-option>
              <a-option value="paid">已支付</a-option>
              <a-option value="overdue">逾期</a-option>
            </a-select>
            <a-button @click="handleResetRecords">重置</a-button>
          </a-space>
        </div>

        <!-- 操作按钮 -->
        <div class="pro-action-bar">
          <a-button @click="loadRecords">刷新</a-button>
          <a-button @click="handleExport">导出</a-button>
        </div>

        <!-- 账单列表 -->
        <div class="pro-content-area">
          <a-table
            :columns="recordColumns"
            :data="records"
            :loading="loading"
            :pagination="pagination"
            row-key="id"
            @page-change="handlePageChange"
          >
            <template #billing_type="{ record }">
              <a-tag :color="getBillingTypeColor(record.billing_type)">{{ getBillingTypeText(record.billing_type) }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getRecordStatusColor(record.status)">{{ getRecordStatusText(record.status) }}</a-tag>
            </template>
            <template #amount="{ record }">
              <span class="amount">¥{{ record.amount }}</span>
            </template>
            <template #period="{ record }">
              {{ formatDate(record.period_start) }} ~ {{ formatDate(record.period_end) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="showRecordDetail(record)">详情</a-button>
                <a-button type="text" size="small" @click="handlePay(record)" v-if="record.status === 'pending'">支付</a-button>
              </a-space>
            </template>
          </a-table>
        </div>
      </a-tab-pane>

      <!-- Tab2: 发票管理 -->
      <a-tab-pane key="invoices" title="发票管理">
        <template #extra>
          <a-button type="primary" @click="showInvoiceModal">新建发票</a-button>
        </template>

        <!-- 搜索 -->
        <div class="pro-search-bar">
          <a-space>
            <a-select v-model="invoiceForm.status" placeholder="发票状态" allow-clear style="width: 140px" @change="loadInvoices">
              <a-option value="pending">待审核</a-option>
              <a-option value="issued">已开票</a-option>
              <a-option value="void">已作废</a-option>
            </a-select>
            <a-button @click="loadInvoices">刷新</a-button>
          </a-space>
        </div>

        <!-- 发票列表 -->
        <div class="pro-content-area">
          <a-table
            :columns="invoiceColumns"
            :data="invoices"
            :loading="invoiceLoading"
            :pagination="invoicePagination"
            row-key="id"
            @page-change="handleInvoicePageChange"
          >
            <template #invoice_type="{ record }">
              <a-tag :color="record.invoice_type === 'VAT' ? 'purple' : 'blue'">
                {{ record.invoice_type === 'VAT' ? '增值税专用发票' : '普通发票' }}
              </a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getInvoiceStatusColor(record.status)">{{ getInvoiceStatusText(record.status) }}</a-tag>
            </template>
            <template #total_amount="{ record }">
              <span class="amount">¥{{ record.total_amount }}</span>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="showInvoiceDetail(record)">查看</a-button>
                <a-button type="text" size="small" v-if="record.status === 'issued'" @click="handleVoidInvoice(record)">作废</a-button>
                <a-button type="text" size="small" v-if="record.status === 'issued'" @click="handleShipInvoice(record)">寄送</a-button>
              </a-space>
            </template>
          </a-table>
        </div>
      </a-tab-pane>
    </a-tabs>

    <!-- 账单详情弹窗 -->
    <a-modal v-model:visible="recordDetailVisible" title="账单详情" :width="560" :footer="null">
      <a-descriptions :column="2" bordered v-if="currentRecord">
        <a-descriptions-item label="账单ID">{{ currentRecord.id }}</a-descriptions-item>
        <a-descriptions-item label="用户ID">{{ currentRecord.user_id }}</a-descriptions-item>
        <a-descriptions-item label="账单号">{{ currentRecord.invoice_no || '-' }}</a-descriptions-item>
        <a-descriptions-item label="计费类型">
          <a-tag :color="getBillingTypeColor(currentRecord.billing_type)">{{ getBillingTypeText(currentRecord.billing_type) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="金额">
          <span class="amount">¥{{ currentRecord.amount }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="货币">{{ currentRecord.currency || 'CNY' }}</a-descriptions-item>
        <a-descriptions-item label="计费周期" :span="2">
          {{ formatDate(currentRecord.period_start) }} ~ {{ formatDate(currentRecord.period_end) }}
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getRecordStatusColor(currentRecord.status)">{{ getRecordStatusText(currentRecord.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="支付时间">{{ formatDate(currentRecord.paid_at) }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 新建发票弹窗 -->
    <a-modal v-model:visible="invoiceModalVisible" title="创建发票申请" :width="560" :loading="invoiceSubmitting" @before-ok="handleCreateInvoice" @cancel="invoiceModalVisible = false">
      <a-form :model="invoiceCreateForm" layout="vertical">
        <a-form-item label="关联账单" required>
          <a-select v-model="invoiceCreateForm.record_id" placeholder="选择关联账单">
            <a-option v-for="r in pendingRecords" :key="r.id" :value="r.id">
              {{ r.invoice_no || '账单-' + r.id }} - ¥{{ r.amount }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发票类型" required>
          <a-select v-model="invoiceCreateForm.invoice_type" placeholder="选择发票类型">
            <a-option value="normal">普通发票</a-option>
            <a-option value="VAT">增值税专用发票</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发票抬头" required>
          <a-input v-model="invoiceCreateForm.title" placeholder="个人姓名或公司名称" />
        </a-form-item>
        <a-form-item label="接收邮箱">
          <a-input v-model="invoiceCreateForm.email" placeholder="接收电子发票的邮箱" />
        </a-form-item>
        <template v-if="invoiceCreateForm.invoice_type === 'VAT'">
          <a-form-item label="税号" required>
            <a-input v-model="invoiceCreateForm.tax_no" placeholder="纳税人识别号" />
          </a-form-item>
          <a-form-item label="开户行">
            <a-input v-model="invoiceCreateForm.bank_name" placeholder="开户银行名称" />
          </a-form-item>
          <a-form-item label="银行账号">
            <a-input v-model="invoiceCreateForm.bank_account" placeholder="银行账号" />
          </a-form-item>
        </template>
      </a-form>
    </a-modal>

    <!-- 发票详情弹窗 -->
    <a-modal v-model:visible="invoiceDetailVisible" title="发票详情" :width="600" :footer="null">
      <a-descriptions :column="2" bordered v-if="currentInvoice">
        <a-descriptions-item label="发票号">{{ currentInvoice.invoice_no }}</a-descriptions-item>
        <a-descriptions-item label="用户ID">{{ currentInvoice.user_id }}</a-descriptions-item>
        <a-descriptions-item label="发票类型">
          <a-tag :color="currentInvoice.invoice_type === 'VAT' ? 'purple' : 'blue'">
            {{ currentInvoice.invoice_type === 'VAT' ? '增值税专用发票' : '普通发票' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="金额">¥{{ currentInvoice.amount }}</a-descriptions-item>
        <a-descriptions-item label="税额">¥{{ currentInvoice.tax_amount }}</a-descriptions-item>
        <a-descriptions-item label="价税合计">
          <span class="amount">¥{{ currentInvoice.total_amount }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="发票抬头">{{ currentInvoice.title }}</a-descriptions-item>
        <a-descriptions-item label="税号">{{ currentInvoice.tax_no || '-' }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getInvoiceStatusColor(currentInvoice.status)">{{ getInvoiceStatusText(currentInvoice.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="开票日期">{{ formatDate(currentInvoice.issue_date) }}</a-descriptions-item>
        <a-descriptions-item label="接收邮箱" :span="2">{{ currentInvoice.email || '-' }}</a-descriptions-item>
      </a-descriptions>
      <div v-if="currentInvoice?.pdf_url" class="invoice-pdf-link">
        <a :href="currentInvoice.pdf_url" target="_blank">查看PDF</a>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'

const activeTab = ref('records')
const loading = ref(false)
const invoiceLoading = ref(false)
const invoiceSubmitting = ref(false)
const records = ref([])
const invoices = ref([])
const pendingRecords = ref([])
const currentRecord = ref(null)
const currentInvoice = ref(null)
const recordDetailVisible = ref(false)
const invoiceModalVisible = ref(false)
const invoiceDetailVisible = ref(false)

const summary = reactive({ monthly_total: 0, pending_amount: 0, paid_this_year: 0, invoice_pending: 0 })

const form = reactive({ user_id: '', billing_type: '', status: '' })
const invoiceForm = reactive({ status: '' })
const invoiceCreateForm = reactive({ record_id: null, invoice_type: 'normal', title: '', email: '', tax_no: '', bank_name: '', bank_account: '' })

const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const invoicePagination = reactive({ current: 1, pageSize: 20, total: 0 })

const recordColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '用户ID', dataIndex: 'user_id', width: 80 },
  { title: '计费类型', slotName: 'billing_type', width: 110 },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '周期', slotName: 'period', width: 200 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '支付时间', dataIndex: 'paid_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 140 }
]

const invoiceColumns = [
  { title: '发票号', dataIndex: 'invoice_no', width: 140, ellipsis: true },
  { title: '发票类型', slotName: 'invoice_type', width: 130 },
  { title: '抬头', dataIndex: 'title', ellipsis: true },
  { title: '金额', slotName: 'total_amount', width: 100 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '开票日期', dataIndex: 'issue_date', width: 170 },
  { title: '操作', slotName: 'actions', width: 160 }
]

const getToken = () => localStorage.getItem('token')

const getBillingTypeColor = (t) => ({ subscription: 'blue', usage: 'green', API_quota: 'orange' }[t] || 'gray')
const getBillingTypeText = (t) => ({ subscription: '订阅费用', usage: '用量费用', API_quota: 'API配额' }[t] || t)
const getRecordStatusColor = (s) => ({ pending: 'orange', paid: 'green', overdue: 'red' }[s] || 'gray')
const getRecordStatusText = (s) => ({ pending: '待支付', paid: '已支付', overdue: '逾期' }[s] || s)
const getInvoiceStatusColor = (s) => ({ pending: 'orange', issued: 'green', void: 'gray' }[s] || 'gray')
const getInvoiceStatusText = (s) => ({ pending: '待审核', issued: '已开票', void: '已作废' }[s] || s)
const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

// ========== 账单记录 ==========
const loadRecords = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (form.user_id) params.user_id = form.user_id
    if (form.billing_type) params.billing_type = form.billing_type
    if (form.status) params.status = form.status

    const res = await fetch(`/api/v1/billing/records?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      records.value = json.data?.list || json.data || []
      pagination.total = json.data?.total || 0
      // 加载待支付账单用于创建发票
      pendingRecords.value = records.value.filter(r => r.status === 'pending')
    }
    // 加载费用汇总
    loadSummary()
  } catch (e) {
    Message.error('加载账单列表失败')
  } finally {
    loading.value = false
  }
}

const loadSummary = async () => {
  try {
    const res = await fetch('/api/v1/billing/summary', {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) Object.assign(summary, json.data || {})
  } catch (e) { /* silent */ }
}

const handlePageChange = (page) => { pagination.current = page; loadRecords() }

const showRecordDetail = (record) => { currentRecord.value = record; recordDetailVisible.value = true }

const handlePay = async (record) => {
  Modal.confirm({ title: '确认支付', content: `确认支付 ¥${record.amount}？`, onOk: async () => {
    try {
      const res = await fetch('/api/v1/billing/pay', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${getToken()}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ record_id: record.id, payment_method: 'alipay' })
      })
      const json = await res.json()
      if (json.code === 0) {
        Message.success('支付成功')
        loadRecords()
      } else {
        Message.error(json.message || '支付失败')
      }
    } catch (e) {
      Message.error('支付失败')
    }
  }})
}

const handleExport = () => {
  window.open('/api/v1/billing/records/export', '_blank')
  Message.success('导出已开始')
}

const handleResetRecords = () => {
  form.user_id = ''
  form.billing_type = ''
  form.status = ''
  pagination.current = 1
  loadRecords()
}

// ========== 发票管理 ==========
const loadInvoices = async () => {
  invoiceLoading.value = true
  try {
    const params = { page: invoicePagination.current, page_size: invoicePagination.pageSize }
    if (invoiceForm.status) params.status = invoiceForm.status

    const res = await fetch(`/api/v1/billing/invoices?${new URLSearchParams(params)}`, {
      headers: { 'Authorization': `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code === 0) {
      invoices.value = json.data?.list || []
      invoicePagination.total = json.data?.total || 0
    }
  } catch (e) {
    Message.error('加载发票列表失败')
  } finally {
    invoiceLoading.value = false
  }
}

const handleInvoicePageChange = (page) => { invoicePagination.current = page; loadInvoices() }

const showInvoiceModal = () => {
  Object.assign(invoiceCreateForm, { record_id: null, invoice_type: 'normal', title: '', email: '', tax_no: '', bank_name: '', bank_account: '' })
  invoiceModalVisible.value = true
}

const handleCreateInvoice = async (done) => {
  if (!invoiceCreateForm.record_id || !invoiceCreateForm.title) {
    Message.warning('请填写必填字段')
    done(false)
    return
  }
  invoiceSubmitting.value = true
  try {
    const res = await fetch('/api/v1/billing/invoices', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${getToken()}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(invoiceCreateForm)
    })
    const json = await res.json()
    if (json.code === 0) {
      Message.success('发票申请已提交')
      invoiceModalVisible.value = false
      loadInvoices()
      done(true)
    } else {
      Message.error(json.message || '创建失败')
      done(false)
    }
  } catch (e) {
    Message.error('创建失败')
    done(false)
  } finally {
    invoiceSubmitting.value = false
  }
}

const showInvoiceDetail = (record) => { currentInvoice.value = record; invoiceDetailVisible.value = true }

const handleVoidInvoice = async (record) => {
  Modal.confirm({ title: '确认作废', content: '确定要作废该发票吗？', onOk: async () => {
    try {
      const res = await fetch(`/api/v1/billing/invoices/${record.id}/void`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${getToken()}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ reason: '用户申请作废' })
      })
      const json = await res.json()
      if (json.code === 0) {
        Message.success('发票已作废')
        loadInvoices()
      } else {
        Message.error(json.message || '作废失败')
      }
    } catch (e) {
      Message.error('作废失败')
    }
  }})
}

const handleShipInvoice = async (record) => {
  Message.info('请到发票寄送功能填写物流信息')
}

onMounted(() => { loadRecords(); loadInvoices() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.billing-tabs { background: #fff; border-radius: 8px; }
.summary-row { margin-bottom: 16px; }
.summary-card { border-radius: 8px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.amount { font-weight: 600; color: #f53f3f; }
.invoice-pdf-link { margin-top: 12px; text-align: center; }
</style>
