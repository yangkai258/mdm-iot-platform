<template>
  <div class="subscriptions-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>订阅管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建订阅
          </a-button>
        </div>
      </template>
      
      <a-tabs default-active-key="subscriptions">
        <a-tab-pane key="subscriptions" title="订阅列表">
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
            <template #plan="{ record }">
              <span style="font-weight: 500;">{{ record.planName }}</span>
              <div style="font-size: 12px; color: #86909c;">{{ record.deviceLimit }}设备</div>
            </template>
            <template #price="{ record }">
              <span style="color: #FF7D00; font-weight: 500;">¥{{ record.price }}/{{ record.period }}</span>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            <template #autoRenew="{ record }">
              <a-switch :checked="record.autoRenew" disabled />
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleView(record)">详情</a-link>
                <a-link @click="handleRenew(record)">续费</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="plans" title="套餐管理">
          <a-table :columns="planColumns" :data="plans" :pagination="pagination">
            <template #price="{ record }">
              <span style="color: #FF7D00; font-weight: 500;">¥{{ record.price }}</span>
            </template>
            <template #status="{ record }">
              <a-switch :checked="record.status === 'active'" disabled />
            </template>
            <template #actions="{ record }">
              <a-link @click="handleEditPlan(record)">编辑</a-link>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'S001', userId: 'U001', tenantId: 'T001', planName: '专业版', deviceLimit: 50, deviceUsed: 32, price: 299, period: '月', status: 'active', autoRenew: true, expiresAt: '2026-04-28', createdAt: '2026-03-28' },
  { id: 'S002', userId: 'U002', tenantId: 'T002', planName: '基础版', deviceLimit: 20, deviceUsed: 15, price: 99, period: '月', status: 'active', autoRenew: false, expiresAt: '2026-04-15', createdAt: '2026-03-15' },
  { id: 'S003', userId: 'U003', tenantId: 'T003', planName: '企业版', deviceLimit: 200, deviceUsed: 0, price: 999, period: '年', status: 'expired', autoRenew: false, expiresAt: '2026-03-01', createdAt: '2025-03-01' },
]);

const plans = ref([
  { id: 'P001', name: '基础版', price: 99, period: '月', deviceLimit: 20, features: '基础功能', status: 'active' },
  { id: 'P002', name: '专业版', price: 299, period: '月', deviceLimit: 50, features: '高级功能+数据分析', status: 'active' },
  { id: 'P003', name: '企业版', price: 999, period: '月', deviceLimit: 200, features: '全功能+专属客服', status: 'active' },
  { id: 'P004', name: '旗舰版', price: 2999, period: '年', deviceLimit: -1, features: '不限设备+私有部署', status: 'inactive' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });

const columns = [
  { title: '订阅ID', dataIndex: 'id', width: 80 },
  { title: '用户ID', dataIndex: 'userId', width: 80 },
  { title: '套餐', slotName: 'plan', width: 150 },
  { title: '价格', slotName: 'price', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '自动续费', slotName: 'autoRenew', width: 80 },
  { title: '到期时间', dataIndex: 'expiresAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const planColumns = [
  { title: '套餐ID', dataIndex: 'id', width: 80 },
  { title: '套餐名称', dataIndex: 'name', width: 100 },
  { title: '价格', slotName: 'price', width: 100 },
  { title: '设备限额', dataIndex: 'deviceLimit', width: 80 },
  { title: '功能', dataIndex: 'features', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 80 },
];

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { active: 'green', expired: 'red', pending: 'orange' };
  return map[status] || 'default';
};

const getStatusText = (status: string) => {
  const map: Record<string, string> = { active: '生效中', expired: '已过期', pending: '待生效' };
  return map[status] || status;
};

const handleCreate = () => {};
const handleView = (record: any) => {};
const handleRenew = (record: any) => {};
const handleEditPlan = (record: any) => {};
const handlePageChange = (page: number) => { pagination.current = page; };
</script>

<style scoped>
.subscriptions-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
