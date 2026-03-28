<template>
  <div class="coupon-center-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="优惠券总数" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已发放" :value="5680" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已使用" :value="2850" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="核销率" :value="50" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>优惠券中心</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建优惠券
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="coupons" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #discount="{ record }">
          <span style="color: #F53F3F; font-weight: bold;">{{ record.discount }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建优惠券" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="优惠券名称" required>
          <a-input v-model="form.name" placeholder="请输入优惠券名称" />
        </a-form-item>
        <a-form-item label="优惠类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="discount">折扣券</a-option>
            <a-option value="cash">代金券</a-option>
            <a-option value="gift">礼品券</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="折扣金额">
          <a-input-number v-model="form.discount" :min="0" :precision="2" />
        </a-form-item>
        <a-form-item label="使用门槛">
          <a-input-number v-model="form.minSpend" :min="0" placeholder="满多少可用" />
        </a-form-item>
        <a-form-item label="发放数量">
          <a-input-number v-model="form.totalCount" :min="1" />
        </a-form-item>
        <a-form-item label="有效期">
          <a-range-picker v-model="form.dateRange" style="width: 100%;" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ name: '', type: 'discount', discount: 10, minSpend: 100, totalCount: 1000, dateRange: [] });

const coupons = ref([
  { id: 1, name: '新人10元券', type: 'cash', typeText: '代金券', discount: '¥10', minSpend: 100, totalCount: 1000, issuedCount: 850, usedCount: 320, status: 'active', statusText: '进行中' },
  { id: 2, name: '满减20元券', type: 'cash', typeText: '代金券', discount: '¥20', minSpend: 200, totalCount: 500, issuedCount: 500, usedCount: 280, status: 'active', statusText: '进行中' },
  { id: 3, name: '8折折扣券', type: 'discount', typeText: '折扣券', discount: '20%', minSpend: 100, totalCount: 300, issuedCount: 300, usedCount: 150, status: 'expired', statusText: '已过期' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '优惠券名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '优惠', slotName: 'discount', width: 100 },
  { title: '门槛', dataIndex: 'minSpend', width: 100 },
  { title: '发放/总量', dataIndex: 'issuedCount', width: 120 },
  { title: '已用', dataIndex: 'usedCount', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
];

const getStatusColor = (s: string) => ({ active: 'green', paused: 'orange', expired: 'gray' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.coupon-center-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
