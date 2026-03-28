<template>
  <div class="resource-quota-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="租户数" :value="25" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="配额使用" :value="68" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="超配告警" :value="3" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="总配额" :value="10000" suffix="API/日" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>资源配额管理</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="quotas" title="配额设置">
          <a-table :columns="quotaColumns" :data="quotas" :pagination="false">
            <template #usage="{ record }">
              <a-progress :percent="record.usagePercent" :color="record.usagePercent > 80 ? 'red' : 'blue'" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="tenants" title="租户配额">
          <a-table :columns="tenantColumns" :data="tenants" :pagination="pagination">
            <template #actions="{ record }">
              <a-link @click="handleEditQuota(record)">调整配额</a-link>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="plans" title="配额套餐">
          <a-row :gutter="16">
            <a-col :span="6" v-for="plan in plans" :key="plan.id">
              <a-card class="plan-card">
                <div class="plan-name">{{ plan.name }}</div>
                <div class="plan-price">{{ plan.price }}/月</div>
                <a-divider />
                <div v-for="(value, key) in plan.limits" :key="key">{{ key }}: {{ value }}</div>
                <a-button type="primary" long style="margin-top: 16px;">选择套餐</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const quotas = ref([
  { resource: 'API调用', used: 6800, total: 10000, usagePercent: 68 },
  { resource: '设备数量', used: 850, total: 1000, usagePercent: 85 },
  { resource: '存储空间', used: 256, total: 500, usagePercent: 51 },
  { resource: '带宽', used: 450, total: 1000, usagePercent: 45 },
]);

const tenants = ref([
  { id: 1, name: '租户A', plan: '专业版', apiQuota: 10000, usedQuota: 6800, deviceQuota: 1000, usedDevices: 850 },
  { id: 2, name: '租户B', plan: '基础版', apiQuota: 1000, usedQuota: 450, deviceQuota: 100, usedDevices: 45 },
]);

const plans = ref([
  { id: 1, name: '基础版', price: 99, limits: { 'API调用': '1,000/日', '设备数': '10台', '存储': '10GB' } },
  { id: 2, name: '专业版', price: 299, limits: { 'API调用': '10,000/日', '设备数': '100台', '存储': '100GB' } },
  { id: 3, name: '企业版', price: 999, limits: { 'API调用': '无限制', '设备数': '无限制', '存储': '1TB' } },
]);

const quotaColumns = [
  { title: '资源', dataIndex: 'resource', width: 150 },
  { title: '已使用', dataIndex: 'used', width: 120 },
  { title: '总额度', dataIndex: 'total', width: 120 },
  { title: '使用率', slotName: 'usage' },
];

const tenantColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '租户名称', dataIndex: 'name', width: 150 },
  { title: '当前套餐', dataIndex: 'plan', width: 100 },
  { title: 'API配额/使用', dataIndex: 'usedQuota', width: 150 },
  { title: '设备配额/使用', dataIndex: 'usedDevices', width: 150 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const handleEditQuota = (record: any) => {};
</script>

<style scoped>
.resource-quota-container { padding: 20px; }
.plan-card { text-align: center; }
.plan-name { font-weight: bold; font-size: 18px; }
.plan-price { color: #F53F3F; font-size: 24px; }
</style>
