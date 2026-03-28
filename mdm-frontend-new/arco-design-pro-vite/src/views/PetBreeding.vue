<template>
  <div class="pet-breeding-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="配对记录" :value="36" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="待确认" :value="5" status="warning" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="配对成功率" :value="78" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物配对繁育</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建配对
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="records" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="宠物配对" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="父本" required>
          <a-select v-model="form.fatherId" placeholder="选择父本宠物">
            <a-option value="P001">小黄</a-option>
            <a-option value="P002">旺财</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="母本" required>
          <a-select v-model="form.motherId" placeholder="选择母本宠物">
            <a-option value="P003">小红</a-option>
            <a-option value="P004">咪咪</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预计分娩日期">
          <a-date-picker v-model="form.expectedDate" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ fatherId: '', motherId: '', expectedDate: '', note: '' });

const records = ref([
  { id: 1, fatherName: '小黄', motherName: '小红', status: 'pregnant', statusText: '待产', expectedDate: '2026-04-15', createdAt: '2026-03-20' },
  { id: 2, fatherName: '旺财', motherName: '咪咪', status: 'completed', statusText: '已完成', expectedDate: '2026-02-10', createdAt: '2026-01-15', births: 3 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '父本', dataIndex: 'fatherName', width: 120 },
  { title: '母本', dataIndex: 'motherName', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '预计日期', dataIndex: 'expectedDate', width: 120 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
];

const getStatusColor = (s: string) => ({ pending: 'orange', pregnant: 'blue', completed: 'green', cancelled: 'gray' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.pet-breeding-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
