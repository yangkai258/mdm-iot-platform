<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>订阅管理</a-breadcrumb-item>
      <a-breadcrumb-item>发票管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="pro-page-header">
      <h2 class="pro-page-title">发票管理</h2>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="handleApplyInvoice">
          <template #icon><icon-plus /></template>
          申请发票
        </a-button>
        <a-button @click="loadData">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <a-tabs v-model="activeTab" v-loading="loading">
      <!-- 账单记录 -->
      <a-tab-pane key="records" title="账单记录">
        <a-table :columns="recordColumns" :data="records" :loading="loading" :pagination="{ pageSize: 10 }">
          <template #status="{ record }">
            <a-tag :color="record.status === 'paid' ? 'green' : 'orange'">
              {{ record.status === 'paid' ? '已支付' : '待支付' }}
            </a-tag>
          </template>
          <template #operations="{ record }">
            <a-space>
              <a-button size="small" @click="handleCreateInvoice(record)">开发票</a-button>
              <a-button size="small" @click="handleRecordDetail(record)">详情</a-button>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>

      <!-- 发票列表 -->
      <a-tab-pane key="invoices" title="发票列表">
        <a-table :columns="invoiceColumns" :data="invoices" :loading="loading" :pagination="{ pageSize: 10 }">
          <template #status="{ record }">
            <a-tag :color="record.status === 'issued' ? 'green' : 'orange'">
              {{ record.status === 'issued' ? '已开' : '待开' }}
            </a-tag>
          </template>
          <template #operations="{ record }">
            <a-space>
              <a-button size="small" v-if="record.status === 'issued'" @click="handleDownload(record)">下载</a-button>
              <a-button size="small" @click="handleInvoiceDetail(record)">详情</a-button>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>
    </a-tabs>

    <!-- 申请发票弹窗 -->
    <a-modal
      v-model:visible="applyModalVisible"
      title="申请发票"
      @before-ok="confirmApplyInvoice"
    >
      <a-form :model="invoiceForm" layout="vertical">
        <a-form-item label="发票抬头" required>
          <a-input v-model="invoiceForm.title" placeholder="请输入发票抬头" />
        </a-form-item>
        <a-form-item label="税号">
          <a-input v-model="invoiceForm.tax_number" placeholder="请输入税号（可选）" />
        </a-form-item>
        <a-form-item label="发票金额" required>
          <a-input-number v-model="invoiceForm.amount" :min="0" :step="0.01" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="发票类型">
          <a-select v-model="invoiceForm.invoice_type">
            <a-option value="normal">普通发票</a-option>
            <a-option value="special">增值税专用发票</a-option>
            <a-option value="electronic">电子发票</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { billingApi } from '@/api/billing'

const loading = ref(false)
const activeTab = ref('records')
const applyModalVisible = ref(false)
const invoiceForm = ref({ title: '', tax_number: '', amount: 0, invoice_type: 'normal' })

const recordColumns = [
  { title: '日期', dataIndex: 'date' },
  { title: '类型', dataIndex: 'type' },
  { title: '金额', dataIndex: 'amount' },
  { title: '状态', dataIndex: 'status', slotName: 'status' },
  { title: '操作', slotName: 'operations', width: 180 }
]

const invoiceColumns = [
  { title: '发票抬头', dataIndex: 'title' },
  { title: '税号', dataIndex: 'tax_number' },
  { title: '金额', dataIndex: 'amount' },
  { title: '状态', dataIndex: 'status', slotName: 'status' },
  { title: '操作', slotName: 'operations', width: 180 }
]

const records = ref([
  { id: 1, date: '2026-03-01', type: '订阅', amount: '¥99.00', status: 'paid' },
  { id: 2, date: '2026-02-01', type: '订阅', amount: '¥99.00', status: 'paid' },
  { id: 3, date: '2026-01-01', type: '订阅', amount: '¥99.00', status: 'paid' }
])

const invoices = ref([
  { id: 1, title: '公司名称A', tax_number: '911xxxxx', amount: '¥297.00', status: 'issued' }
])

const loadData = async () => {
  loading.value = true
  try {
    const [recordsRes, invoicesRes] = await Promise.all([
      billingApi.getBillingRecords(),
      billingApi.getInvoices()
    ])
    if (recordsRes.code === 0 || recordsRes.code === 200) {
      records.value = recordsRes.data || records.value
    }
    if (invoicesRes.code === 0 || invoicesRes.code === 200) {
      invoices.value = invoicesRes.data || invoices.value
    }
  } catch (e) {
    // use mock data
  } finally {
    loading.value = false
  }
}

const handleApplyInvoice = () => {
  applyModalVisible.value = true
}

const confirmApplyInvoice = async (done: (val: boolean) => void) => {
  if (!invoiceForm.value.title || invoiceForm.value.amount <= 0) {
    Message.error('请填写完整的发票信息')
    done(false)
    return
  }
  try {
    const res = await billingApi.createInvoice(invoiceForm.value)
    if (res.code === 0 || res.code === 200) {
      Message.success('申请成功')
      applyModalVisible.value = false
      loadData()
      done(true)
    } else {
      Message.error('申请失败')
      done(false)
    }
  } catch (e) {
    Message.error('申请失败')
    done(false)
  }
}

const handleCreateInvoice = (record: any) => {
  invoiceForm.value.amount = parseFloat(record.amount.replace('¥', ''))
  applyModalVisible.value = true
}

const handleRecordDetail = (record: any) => {
  Message.info(`查看账单详情: ${record.id}`)
}

const handleDownload = (invoice: any) => {
  Message.success(`下载发票: ${invoice.title}`)
}

const handleInvoiceDetail = (invoice: any) => {
  Message.info(`查看发票详情: ${invoice.id}`)
}

onMounted(() => {
  loadData()
})
</script>
