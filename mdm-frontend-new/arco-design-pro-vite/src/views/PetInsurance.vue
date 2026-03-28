<template>
  <div class="pet-insurance-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="投保宠物" :value="856" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="保费收入" :value="128500" prefix="¥" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="理赔次数" :value="45" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="理赔金额" :value="25800" prefix="¥" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物保险管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建保单
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="policies" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleClaim(record)">理赔</a-link>
            <a-link @click="handleView(record)">详情</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建保单" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="保险方案">
          <a-select v-model="form.plan" placeholder="选择方案">
            <a-option value="basic">基础版</a-option>
            <a-option value="standard">标准版</a-option>
            <a-option value="premium">高级版</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="投保期限">
          <a-select v-model="form.term">
            <a-option value="1year">1年</a-option>
            <a-option value="2year">2年</a-option>
          </a-select>
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

const form = reactive({ petId: '', plan: '', term: '' });

const policies = ref([
  { id: 1, policyNo: 'INS202603280001', petName: '小黄', plan: 'premium', planText: '高级版', premium: 2999, startDate: '2026-03-01', endDate: '2027-02-28', status: 'active', statusText: '生效中' },
  { id: 2, policyNo: 'INS202603150001', petName: '小红', plan: 'standard', planText: '标准版', premium: 1999, startDate: '2026-03-15', endDate: '2027-03-14', status: 'active', statusText: '生效中' },
]);

const columns = [
  { title: '保单号', dataIndex: 'policyNo', width: 160 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '方案', dataIndex: 'planText', width: 100 },
  { title: '保费', dataIndex: 'premium', width: 100 },
  { title: '生效日期', dataIndex: 'startDate', width: 120 },
  { title: '到期日期', dataIndex: 'endDate', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const getStatusColor = (s: string) => ({ active: 'green', expired: 'gray', cancelled: 'red' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleClaim = (record: any) => {};
const handleView = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.pet-insurance-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
