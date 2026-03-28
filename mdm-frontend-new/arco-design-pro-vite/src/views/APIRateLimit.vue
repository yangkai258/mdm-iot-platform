<template>
  <div class="api-rate-limit-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="限流规则" :value="25" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="触发次数" :value="5680" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="拦截率" :value="2.5" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="白名单" :value="36" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>API限流管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建规则
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="rules" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #enabled="{ record }">
          <a-switch v-model="record.enabled" @change="handleToggle(record)" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑规则' : '创建规则'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model="form.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="限流类型">
          <a-select v-model="form.type">
            <a-option value="ip">IP限流</a-option>
            <a-option value="user">用户限流</a-option>
            <a-option value="api">API限流</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="限制阈值">
          <a-input-number v-model="form.threshold" :min="1" />
        </a-form-item>
        <a-form-item label="时间窗口">
          <a-input-number v-model="form.window" :min="1" /> 秒
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: '', threshold: 100, window: 60, enabled: true });

const rules = ref([
  { id: 1, name: '登录限流', type: 'ip', typeText: 'IP限流', threshold: 10, window: 60, enabled: true, triggeredCount: 256 },
  { id: 2, name: 'API调用限流', type: 'api', typeText: 'API限流', threshold: 1000, window: 60, enabled: true, triggeredCount: 1250 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '阈值', dataIndex: 'threshold', width: 80 },
  { title: '窗口(秒)', dataIndex: 'window', width: 100 },
  { title: '触发次数', dataIndex: 'triggeredCount', width: 100 },
  { title: '启用', slotName: 'enabled', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleToggle = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.api-rate-limit-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
