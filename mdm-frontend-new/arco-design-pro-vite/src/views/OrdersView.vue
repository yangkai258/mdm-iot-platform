<template>
  <div class="orders-container">
    <a-card>
      <template #title>
        <span>订单管理</span>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="4">
            <a-input v-model="searchForm.orderId" placeholder="订单号" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="订单状态" allow-clear>
              <a-option value="pending">待支付</a-option>
              <a-option value="paid">已支付</a-option>
              <a-option value="completed">已完成</a-option>
              <a-option value="cancelled">已取消</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #amount="{ record }">
          <span style="color: #FF7D00; font-weight: 500;">¥{{ record.amount }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-link @click="handleView(record)">详情</a-link>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'ORD001', orderId: '2026032800001', memberName: '张三', amount: 299, items: '专业版月卡 x1', status: 'paid', createdAt: '2026-03-28 10:30:00' },
  { id: 'ORD002', orderId: '2026032700001', memberName: '李四', amount: 99, items: '基础版月卡 x1', status: 'completed', createdAt: '2026-03-27 14:20:00' },
  { id: 'ORD003', orderId: '2026032700002', memberName: '王五', amount: 0, items: '积分兑换', status: 'completed', createdAt: '2026-03-27 16:45:00' },
  { id: 'ORD004', orderId: '2026032600001', memberName: '赵六', amount: 999, items: '企业版年卡 x1', status: 'cancelled', createdAt: '2026-03-26 09:15:00' },
]);

const searchForm = reactive({ orderId: '', status: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '订单号', dataIndex: 'orderId', width: 150 },
  { title: '会员', dataIndex: 'memberName', width: 100 },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '商品', dataIndex: 'items', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 80, fixed: 'right' },
];

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { pending: 'orange', paid: 'blue', completed: 'green', cancelled: 'gray' };
  return map[status] || 'default';
};

const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待支付', paid: '已支付', completed: '已完成', cancelled: '已取消' };
  return map[status] || status;
};

const handleSearch = () => {};
const handleView = (record: any) => {};
</script>

<style scoped>
.orders-container { padding: 20px; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
