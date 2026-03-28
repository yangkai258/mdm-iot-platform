<template>
  <div class="change-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="变更单" :value="256" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="待审批" :value="8" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已完成" :value="240" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="失败" :value="8" status="error" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>变更管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建变更单
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="changes" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建变更单" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="变更标题" required>
          <a-input v-model="form.title" placeholder="请输入变更标题" />
        </a-form-item>
        <a-form-item label="变更类型">
          <a-select v-model="form.type">
            <a-option value="normal">普通变更</a-option>
            <a-option value="emergency">紧急变更</a-option>
            <a-option value="rollback">回滚变更</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="变更内容">
          <a-textarea v-model="form.content" :rows="4" />
        </a-form-item>
        <a-form-item label="执行时间">
          <a-date-picker v-model="form.executeAt" style="width: 100%;" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const createVisible = ref(false);

const form = reactive({ title: '', type: '', content: '', executeAt: '' });

const changes = ref([
  { id: 1, title: '数据库索引优化', type: 'normal', typeText: '普通', status: 'approved', statusText: '已批准', applicant: '张三', executeAt: '2026-03-29 02:00:00' },
  { id: 2, title: '紧急修复登录bug', type: 'emergency', typeText: '紧急', status: 'pending', statusText: '待审批', applicant: '李四', executeAt: '2026-03-28 20:00:00' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '申请人', dataIndex: 'applicant', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '执行时间', dataIndex: 'executeAt', width: 160 },
];

const getStatusColor = (s: string) => ({ pending: 'orange', approved: 'blue', executing: 'processing', completed: 'green', failed: 'red' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.change-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
