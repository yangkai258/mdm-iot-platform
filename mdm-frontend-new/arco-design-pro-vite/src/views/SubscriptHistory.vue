<template>
  <div class="subscription-history-container">
    <a-card>
      <template #title>
        <span>订阅历史</span>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #amount="{ record }">
          <span style="color: #F53F3F; font-weight: bold;">¥{{ record.amount }}</span>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const data = ref([
  { id: 'SUB001', memberName: '张三', type: 'renewal', typeText: '续费', amount: 299, status: 'paid', statusText: '已支付', paidAt: '2026-03-28 10:00:00', subscription: '年度高级版' },
  { id: 'SUB002', memberName: '李四', type: 'upgrade', typeText: '升级', amount: 100, status: 'paid', statusText: '已支付', paidAt: '2026-03-27 15:00:00', subscription: '月度基础版→年度基础版' },
  { id: 'SUB003', memberName: '王五', type: 'new', typeText: '新订阅', amount: 29, status: 'paid', statusText: '已支付', paidAt: '2026-03-26 10:00:00', subscription: '月度基础版' },
  { id: 'SUB004', memberName: '赵六', type: 'renewal', typeText: '续费', amount: 99, status: 'pending', statusText: '待支付', paidAt: null, subscription: '月度高级版' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '订阅ID', dataIndex: 'id', width: 120 },
  { title: '会员', dataIndex: 'memberName', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '订阅计划', dataIndex: 'subscription', width: 200 },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '支付时间', dataIndex: 'paidAt', width: 160 },
];

const getStatusColor = (s: string) => ({ paid: 'green', pending: 'orange', cancelled: 'gray', refunded: 'red' }[s] || 'default');
</script>

<style scoped>
.subscription-history-container { padding: 20px; }
</style>
