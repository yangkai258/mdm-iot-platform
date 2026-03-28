<template>
  <div class="questionnaire-survey-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="问卷总数" :value="36" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="收集答卷" :value="5680" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="进行中" :value="5" status="processing" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>问卷调查管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建问卷
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="questionnaires" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'default'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleViewResults(record)">结果</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建问卷" :width="800" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="问卷标题" required>
          <a-input v-model="form.title" placeholder="请输入问卷标题" />
        </a-form-item>
        <a-form-item label="问卷描述">
          <a-textarea v-model="form.description" :rows="2" />
        </a-form-item>
        <a-form-item label="问题">
          <a-textarea v-model="form.questions" :rows="6" placeholder="每行一个问题，格式: 问题|类型(单选/多选/填空)" />
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

const form = reactive({ title: '', description: '', questions: '', dateRange: [] });

const questionnaires = ref([
  { id: 1, title: '用户满意度调查', description: '了解用户对产品的满意度', responses: 1250, status: 'active', statusText: '进行中', endDate: '2026-03-31' },
  { id: 2, title: '功能需求调研', description: '收集用户对新功能的需求', responses: 580, status: 'ended', statusText: '已结束', endDate: '2026-03-15' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '答卷数', dataIndex: 'responses', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '结束日期', dataIndex: 'endDate', width: 120 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const handleCreate = () => { createVisible.value = true; };
const handleViewResults = (record: any) => {};
const handleEdit = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.questionnaire-survey-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
