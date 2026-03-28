<template>
  <div class="billing-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="本月收入" :value="stats.monthlyRevenue" prefix="¥" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="本月支出" :value="stats.monthlyExpense" prefix="¥" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="本月订单" :value="stats.monthlyOrders" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="待结算" :value="stats.pendingSettlement" prefix="¥" status="warning" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>账单管理</span>
          <a-space>
            <a-button @click="handleExport">导出</a-button>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              创建账单
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs v-model="activeTab">
        <a-tab-pane key="records" title="账单记录">
          <a-table :columns="billColumns" :data="bills" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="record.type === 'income' ? 'green' : 'red'">
                {{ record.type === 'income' ? '收入' : '支出' }}
              </a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleView(record)">详情</a-link>
                <a-link v-if="record.status === 'pending'" @click="handlePay(record)">结算</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="invoices" title="发票管理">
          <a-table :columns="invoiceColumns" :data="invoices" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getInvoiceStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleViewInvoice(record)">查看</a-link>
                <a-link v-if="record.status === 'pending'" @click="handleIssueInvoice(record)">开票</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 创建账单弹窗 -->
    <a-modal v-model:visible="createVisible" title="创建账单" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="账单类型">
          <a-radio-group v-model="form.type">
            <a-radio value="income">收入</a-radio>
            <a-radio value="expense">支出</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="关联订单">
          <a-select v-model="form.orderId" placeholder="选择关联订单">
            <a-option value="O001">O001 - ¥299.00</a-option>
            <a-option value="O002">O002 - ¥99.00</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="金额">
          <a-input-number v-model="form.amount" :min="0" :precision="2" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const activeTab = ref('records');

const stats = reactive({
  monthlyRevenue: 158320,
  monthlyExpense: 45320,
  monthlyOrders: 1256,
  pendingSettlement: 28500,
});

const bills = ref([
  { id: 'B001', orderId: 'O001', type: 'income', amount: 299, status: 'paid', statusText: '已结算', createdAt: '2026-03-28 10:00:00', paidAt: '2026-03-28 10:30:00' },
  { id: 'B002', orderId: 'O002', type: 'income', amount: 99, status: 'pending', statusText: '待结算', createdAt: '2026-03-27 15:00:00', paidAt: null },
  { id: 'B003', orderId: null, type: 'expense', amount: 5000, status: 'paid', statusText: '已结算', createdAt: '2026-03-26 09:00:00', paidAt: '2026-03-26 18:00:00' },
]);

const invoices = ref([
  { id: 'INV001', title: 'XXX科技有限公司', taxNo: '91110000XXXXXXXX', amount: 2990, status: 'issued', statusText: '已开票', createdAt: '2026-03-28 10:00:00' },
  { id: 'INV002', title: 'YYY商贸公司', taxNo: '91110000XXXXXXXX', amount: 990, status: 'pending', statusText: '待开票', createdAt: '2026-03-27 15:00:00' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const createVisible = ref(false);

const form = reactive({ type: 'income', orderId: '', amount: 0, note: '' });

const billColumns = [
  { title: '账单ID', dataIndex: 'id', width: 120 },
  { title: '关联订单', dataIndex: 'orderId', width: 120 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '金额', dataIndex: 'amount', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
  { title: '结算时间', dataIndex: 'paidAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const invoiceColumns = [
  { title: '发票号', dataIndex: 'id', width: 120 },
  { title: '抬头', dataIndex: 'title', width: 200 },
  { title: '税号', dataIndex: 'taxNo', width: 180 },
  { title: '金额', dataIndex: 'amount', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '申请时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getStatusColor = (s: string) => ({ paid: 'green', pending: 'orange', cancelled: 'gray' }[s] || 'default');
const getInvoiceStatusColor = (s: string) => ({ issued: 'green', pending: 'orange', voided: 'gray' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleExport = () => {};
const handleView = (record: any) => {};
const handlePay = (record: any) => {};
const handleViewInvoice = (record: any) => {};
const handleIssueInvoice = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.billing-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
