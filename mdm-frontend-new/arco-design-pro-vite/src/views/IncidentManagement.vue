<template>
  <div class="incident-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="事件总数" :value="568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="处理中" :value="12" status="processing" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="严重事件" :value="5" status="error" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="解决率" :value="98" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>故障事件管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建事件
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="incidents" :loading="loading" :pagination="pagination">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">{{ record.levelText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="getStatusBadge(record.status)" :text="record.statusText" />
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建事件" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="事件标题" required>
          <a-input v-model="form.title" placeholder="请输入事件标题" />
        </a-form-item>
        <a-form-item label="严重级别">
          <a-select v-model="form.level">
            <a-option value="low">低</a-option>
            <a-option value="medium">中</a-option>
            <a-option value="high">高</a-option>
            <a-option value="critical">严重</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="事件描述">
          <a-textarea v-model="form.description" :rows="4" />
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

const form = reactive({ title: '', level: '', description: '' });

const incidents = ref([
  { id: 1, title: '数据库连接超时', level: 'high', levelText: '高', status: 'investigating', statusText: '调查中', assignee: '张三', createdAt: '2026-03-28 18:00:00' },
  { id: 2, title: 'API响应缓慢', level: 'medium', levelText: '中', status: 'resolved', statusText: '已解决', assignee: '李四', createdAt: '2026-03-28 10:00:00' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '级别', slotName: 'level', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '负责人', dataIndex: 'assignee', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
];

const getLevelColor = (l: string) => ({ low: 'green', medium: 'blue', high: 'orange', critical: 'red' }[l] || 'default');
const getStatusBadge = (s: string) => ({ new: 'warning', investigating: 'processing', resolved: 'success', closed: 'default' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.incident-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
