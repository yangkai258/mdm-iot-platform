<template>
  <div class="coupons-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>优惠券管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建优惠券
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="record.type === 'discount' ? 'blue' : 'orange'">
            {{ record.type === 'discount' ? '折扣' : '代金券' }}
          </a-tag>
        </template>
        <template #value="{ record }">
          <span style="color: #FF7D00;">
            {{ record.type === 'discount' ? record.value + '%' : '¥' + record.value }}
          </span>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
            {{ record.status === 'active' ? '进行中' : '已结束' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleGrant(record)">发放</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'C001', name: '新人专享券', type: 'discount', value: 20, minAmount: 100, validDays: 30, granted: 1250, used: 380, status: 'active', createdAt: '2026-03-01' },
  { id: 'C002', name: '满100减20', type: 'voucher', value: 20, minAmount: 100, validDays: 7, granted: 560, used: 120, status: 'active', createdAt: '2026-03-15' },
  { id: 'C003', name: '限时5折券', type: 'discount', value: 50, minAmount: 200, validDays: 3, granted: 200, used: 45, status: 'expired', createdAt: '2026-02-28' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });

const columns = [
  { title: '优惠券ID', dataIndex: 'id', width: 100 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '面值', slotName: 'value', width: 80 },
  { title: '满减条件', dataIndex: 'minAmount', width: 100 },
  { title: '有效期', dataIndex: 'validDays', width: 80 },
  { title: '已发放', dataIndex: 'granted', width: 80 },
  { title: '已使用', dataIndex: 'used', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const handleCreate = () => {};
const handleView = (record: any) => {};
const handleGrant = (record: any) => {};
</script>

<style scoped>
.coupons-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
