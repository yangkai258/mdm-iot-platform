<template>
  <div class="page-container">
    <a-card class="general-card" title="鍙戠エ绠＄悊">
      <template #extra>
        <a-button type="primary" @click="openApplyModal"><icon-plus />寮€绁ㄧ敵璇?/a-button>
      </template>
      <a-tabs v-model:active-key="activeTab">
        <a-tab-pane key="list" title="鍙戠エ鍒楄〃">
          <div class="search-form">
            <a-form :model="form" layout="inline">
              <a-form-item label="鍙戠エ鐘舵€?>
                <a-select v-model="form.status" placeholder="閫夋嫨鐘舵€" allow-clear style="width: 140px">
                  <a-option value="pending">寰呭紑绁?/a-option>
                  <a-option value="issued">宸插紑绁?/a-option>
                  <a-option value="delivered">宸查€佽揪</a-option>
                  <a-option value="cancelled">宸蹭綔搴?/a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="鏃堕棿鑼冨洿"><a-range-picker v-model="form.time_range" style="width: 240px" /></a-form-item>
              <a-form-item><a-button type="primary" @click="loadData">鏌ヨ</a-button><a-button @click="Object.assign(form, { status: '', time_range: [] }); loadData()">閲嶇疆</a-button></a-form-item>
            </a-form>
          </div>
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="onPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #amount="{ record }">楼{{ record.amount?.toFixed(2) }}</template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="viewInvoice(record)">璇︽儏</a-button>
              <a-button v-if="record.status === 'issued'" type="text" size="small" @click="downloadInvoice(record)">涓嬭浇</a-button>
              <a-button v-if="record.status === 'pending'" type="text" size="small" status="danger" @click="cancelInvoice(record)">鍙栨秷</a-button>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="history" title="寮€绁ㄥ巻鍙?>
          <a-table :columns="historyColumns" :data="historyData" :loading="historyLoading" :pagination="historyPagination" @page-change="onHistoryPageChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
            </template>
            <template #amount="{ record }">楼{{ record.amount?.toFixed(2) }}</template>
            <template #actions="{ record }">
              <a-button type="text" size="small" @click="viewInvoice(record)">璇︽儏</a-button>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
    <!-- 寮€绁ㄧ敵璇峰脊绐?-->
    <a-modal v-model:visible="applyVisible" title="寮€绁ㄧ敵璇? @before-ok="handleApply" :loading="submitting" :width="600">
      <a-form :model="applyForm" layout="vertical">
        <a-form-item label="鍙戠エ绫诲瀷" required>
          <a-radio-group v-model="applyForm.invoice_type">
            <a-radio value="normal">鏅€氬彂绁?/a-radio>
            <a-radio value="vat">澧炲€肩◣涓撶敤鍙戠エ</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="鍙戠エ鎶ご" required><a-input v-model="applyForm.title" placeholder="璇疯緭鍏ュ彂绁ㄦ姮澶" /></a-form-item>
        <a-form-item label="绋庡彿" required><a-input v-model="applyForm.tax_no" placeholder="璇疯緭鍏ョ◣鍙" /></a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="寮€鎴烽摱琛?><a-input v-model="applyForm.bank" placeholder="璇疯緭鍏ュ紑鎴烽摱琛" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="閾惰璐﹀彿"><a-input v-model="applyForm.bank_account" placeholder="璇疯緭鍏ラ摱琛岃处鍙" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="娉ㄥ唽鍦板潃"><a-input v-model="applyForm.address" placeholder="璇疯緭鍏ユ敞鍐屽湴鍧€" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="鑱旂郴鐢佃瘽"><a-input v-model="applyForm.phone" placeholder="璇疯緭鍏ヨ仈绯荤數璇" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="寮€绁ㄩ噾棰? required>
          <a-input-number v-model="applyForm.amount" :min="0" :precision="2" placeholder="璇疯緭鍏ュ紑绁ㄩ噾棰" style="width: 200px" />
          <span style="margin-left: 8px">鍏?/span>
        </a-form-item>
        <a-form-item label="鎺ユ敹閭"><a-input v-model="applyForm.email" placeholder="璇疯緭鍏ユ帴鏀堕偖绠" /></a-form-item>
        <a-form-item label="澶囨敞"><a-textarea v-model="applyForm.notes" :rows="2" placeholder="澶囨敞淇℃伅" /></a-form-item>
      </a-form>
    </a-modal>
    <!-- 鍙戠エ璇︽儏 -->
    <a-modal v-model:visible="detailVisible" title="鍙戠エ璇︽儏" :footer="null" :width="600">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="鍙戠エ鍙? :span="2">{{ currentInvoice?.invoice_no }}</a-descriptions-item>
        <a-descriptions-item label="鍙戠エ绫诲瀷">{{ currentInvoice?.invoice_type === 'normal' ? '鏅€氬彂绁? : '澧炲€肩◣涓撶敤鍙戠エ' }}</a-descriptions-item>
        <a-descriptions-item label="鍙戠エ閲戦">楼{{ currentInvoice?.amount?.toFixed(2) }}</a-descriptions-item>
        <a-descriptions-item label="鍙戠エ鎶ご" :span="2">{{ currentInvoice?.title }}</a-descriptions-item>
        <a-descriptions-item label="绋庡彿" :span="2">{{ currentInvoice?.tax_no }}</a-descriptions-item>
        <a-descriptions-item label="鐘舵€?><a-tag :color="getStatusColor(currentInvoice?.status)">{{ getStatusText(currentInvoice?.status) }}</a-tag></a-descriptions-item>
        <a-descriptions-item label="鐢宠鏃堕棿">{{ currentInvoice?.created_at }}</a-descriptions-item>
        <a-descriptions-item label="寮€绁ㄦ椂闂? :span="2">{{ currentInvoice?.issued_at || '-' }}</a-descriptions-item>
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
  { title: '鍙戠エ鍙?, dataIndex: 'invoice_no', width: 140 },
  { title: '鍙戠エ鎶ご', dataIndex: 'title', ellipsis: true },
  { title: '閲戦', slotName: 'amount', width: 100 },
  { title: '鍙戠エ绫诲瀷', dataIndex: 'invoice_type', width: 120 },
  { title: '鐘舵€?, slotName: 'status', width: 90 },
  { title: '鐢宠鏃堕棿', dataIndex: 'created_at', width: 170 },
  { title: '鎿嶄綔', slotName: 'actions', width: 140 }
]

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '鍙戠エ鍙?, dataIndex: 'invoice_no', width: 140 },
  { title: '鍙戠エ鎶ご', dataIndex: 'title', ellipsis: true },
  { title: '閲戦', slotName: 'amount', width: 100 },
  { title: '鐘舵€?, slotName: 'status', width: 90 },
  { title: '寮€绁ㄦ椂闂?, dataIndex: 'issued_at', width: 170 }
]

const getStatusColor = (s) => ({ pending: 'orange', issued: 'arcoblue', delivered: 'green', cancelled: 'red' }[s] || 'gray')
const getStatusText = (s) => ({ pending: '寰呭紑绁?, issued: '宸插紑绁?, delivered: '宸查€佽揪', cancelled: '宸蹭綔搴? }[s] || s)

const loadData = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: pagination.current, page_size: pagination.pageSize })
    if (form.status) params.append('status', form.status)
    const res = await fetch(`/api/v1/billing/invoices?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { data.value = res.data?.list || []; pagination.total = res.data?.total || 0 }
    else { data.value = [] }
  } catch (e) { Message.error('鍔犺浇澶辫触') } finally { loading.value = false }
}

const loadHistory = async () => {
  historyLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({ page: historyPagination.current, page_size: historyPagination.pageSize, status: 'issued,delivered' })
    const res = await fetch(`/api/v1/billing/invoices?${params}`, { headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { historyData.value = res.data?.list || []; historyPagination.total = res.data?.total || 0 }
    else { historyData.value = [] }
  } catch (e) { Message.error('鍔犺浇鍘嗗彶澶辫触') } finally { historyLoading.value = false }
}

const openApplyModal = () => { Object.assign(applyForm, { invoice_type: 'normal', title: '', tax_no: '', bank: '', bank_account: '', address: '', phone: '', amount: null, email: '', notes: '' }); applyVisible.value = true }

const handleApply = async (done) => {
  if (!applyForm.title || !applyForm.tax_no || !applyForm.amount) { Message.warning('璇峰～鍐欏繀濉」'); done(false); return }
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/v1/billing/invoices', { method: 'POST', headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' }, body: JSON.stringify(applyForm) }).then(r => r.json())
    if (res.code === 0) { Message.success('鐢宠宸叉彁浜?); applyVisible.value = false; loadData() }
    else { Message.error(res.message || '鐢宠澶辫触') }
    done(true)
  } catch (e) { Message.error('鐢宠澶辫触'); done(false) } finally { submitting.value = false }
}

const viewInvoice = (record) => { currentInvoice.value = record; detailVisible.value = true }

const downloadInvoice = (record) => {
  window.open(`/api/v1/billing/invoices/${record.id}/download`, '_blank')
}

const cancelInvoice = async (record) => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/billing/invoices/${record.id}/cancel`, { method: 'POST', headers: { 'Authorization': `Bearer ${token}` } }).then(r => r.json())
    if (res.code === 0) { Message.success('鍙戠エ宸插彇娑?); loadData() }
    else Message.error('鍙栨秷澶辫触')
  } catch (e) { Message.error('鍙栨秷澶辫触') }
}

const onPageChange = (page) => { pagination.current = page; loadData() }
const onHistoryPageChange = (page) => { historyPagination.current = page; loadHistory() }

onMounted(() => { loadData(); loadHistory() })
</script>

<style scoped>
.page-container { padding: 16px; }
.search-form { margin-bottom: 16px; padding: 16px; background: var(--color-fill-lightest); border-radius: 4px; }
</style>
