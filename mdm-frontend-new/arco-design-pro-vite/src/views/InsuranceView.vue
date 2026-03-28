<template>
  <div class="insurance-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="在保宠物" :value="stats.insured" suffix="只" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="本月新增" :value="stats.thisMonth" suffix="只" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="理赔中" :value="stats.claims" suffix="件" status="warning" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="本月理赔" :value="stats.paid" prefix="¥" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物保险</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增保险
          </a-button>
        </div>
      </template>
      
      <a-tabs v-model="activeTab">
        <a-tab-pane key="list" title="保险列表">
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleView(record)">详情</a-link>
                <a-link @click="handleClaim(record)">理赔</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="claims" title="理赔记录">
          <a-table :columns="claimColumns" :data="claims" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getClaimStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 创建保险弹窗 -->
    <a-modal v-model:visible="createVisible" title="新增保险" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="选择宠物">
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄 - 金毛</a-option>
            <a-option value="P002">小红 - 柯基</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="保险方案">
          <a-select v-model="form.planId" placeholder="选择方案">
            <a-option value="basic">基础版 - 年费299元</a-option>
            <a-option value="standard">标准版 - 年费599元</a-option>
            <a-option value="premium">高级版 - 年费999元</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="被保险人">
          <a-input v-model="form.ownerName" placeholder="请输入被保险人姓名" />
        </a-form-item>
        <a-form-item label="身份证号">
          <a-input v-model="form.idCard" placeholder="请输入身份证号" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const activeTab = ref('list');

const stats = reactive({ insured: 156, thisMonth: 12, claims: 3, paid: 8500 });

const data = ref([
  { id: 'INS001', petName: '小黄', petType: 'dog', breed: '金毛', planName: '高级版', premium: 999, startDate: '2026-01-01', endDate: '2027-01-01', status: 'active', statusText: '生效中' },
  { id: 'INS002', petName: '小红', petType: 'dog', breed: '柯基', planName: '标准版', premium: 599, startDate: '2026-02-15', endDate: '2027-02-15', status: 'active', statusText: '生效中' },
  { id: 'INS003', petName: '咪咪', petType: 'cat', breed: '英短', planName: '基础版', premium: 299, startDate: '2025-06-01', endDate: '2026-06-01', status: 'expired', statusText: '已过期' },
]);

const claims = ref([
  { id: 'CL001', insuranceId: 'INS001', petName: '小黄', claimType: 'illness', claimTypeText: '疾病', amount: 2500, status: 'approved', statusText: '已赔付', claimDate: '2026-03-15' },
  { id: 'CL002', insuranceId: 'INS002', petName: '小红', claimType: 'accident', claimTypeText: '意外', amount: 1800, status: 'pending', statusText: '处理中', claimDate: '2026-03-20' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const createVisible = ref(false);

const form = reactive({ petId: '', planId: '', ownerName: '', idCard: '' });

const columns = [
  { title: '保险单号', dataIndex: 'id', width: 120 },
  { title: '宠物', dataIndex: 'petName', width: 80 },
  { title: '类型', dataIndex: 'petType', width: 60 },
  { title: '品种', dataIndex: 'breed', width: 80 },
  { title: '方案', dataIndex: 'planName', width: 100 },
  { title: '保费', dataIndex: 'premium', width: 80 },
  { title: '生效日期', dataIndex: 'startDate', width: 120 },
  { title: '到期日期', dataIndex: 'endDate', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const claimColumns = [
  { title: '理赔单号', dataIndex: 'id', width: 120 },
  { title: '保险单号', dataIndex: 'insuranceId', width: 120 },
  { title: '宠物', dataIndex: 'petName', width: 80 },
  { title: '理赔类型', dataIndex: 'claimTypeText', width: 100 },
  { title: '理赔金额', dataIndex: 'amount', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '申请日期', dataIndex: 'claimDate', width: 120 },
];

const getStatusColor = (s: string) => ({ active: 'green', expired: 'gray', cancelled: 'red' }[s] || 'default');
const getClaimStatusColor = (s: string) => ({ pending: 'blue', approved: 'green', rejected: 'red' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleView = (record: any) => {};
const handleClaim = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.insurance-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
