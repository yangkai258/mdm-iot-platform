<template>
  <div class="data-export-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="导出任务" :value="56" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="成功导出" :value="52" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="总数据量" :value="2.5" suffix="GB" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>数据导出管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建导出任务
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="exports" :pagination="pagination">
        <template #format="{ record }">
          <a-tag>{{ record.formatText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link v-if="record.status === 'completed'" @click="handleDownload(record)">下载</a-link>
            <a-link v-if="record.status === 'failed'" @click="handleRetry(record)">重试</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="新建导出任务" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="数据模块" required>
          <a-select v-model="form.module" placeholder="选择数据模块">
            <a-option value="devices">设备数据</a-option>
            <a-option value="members">会员数据</a-option>
            <a-option value="pets">宠物数据</a-option>
            <a-option value="orders">订单数据</a-option>
            <a-option value="alerts">告警数据</a-option>
            <a-option value="activities">操作日志</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="导出格式">
          <a-radio-group v-model="form.format">
            <a-radio value="csv">CSV</a-radio>
            <a-radio value="excel">Excel</a-radio>
            <a-radio value="json">JSON</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker v-model="form.dateRange" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="字段选择">
          <a-checkbox-group v-model="form.fields">
            <a-checkbox value="id">ID</a-checkbox>
            <a-checkbox value="name">名称</a-checkbox>
            <a-checkbox value="created_at">创建时间</a-checkbox>
            <a-checkbox value="status">状态</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="压缩导出">
          <a-switch v-model="form.compressed" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ module: '', format: 'csv', dateRange: [], fields: [], compressed: true });

const exports = ref([
  { id: 'EXP001', name: '设备数据导出', module: 'devices', moduleText: '设备数据', format: 'csv', formatText: 'CSV', records: 1250, size: '2.5MB', status: 'completed', statusText: '已完成', createdAt: '2026-03-28 10:00:00' },
  { id: 'EXP002', name: '会员数据导出', module: 'members', moduleText: '会员数据', format: 'excel', formatText: 'Excel', records: 580, size: '1.2MB', status: 'completed', statusText: '已完成', createdAt: '2026-03-27 15:00:00' },
  { id: 'EXP003', name: '订单月报', module: 'orders', moduleText: '订单数据', format: 'excel', formatText: 'Excel', records: 0, size: '-', status: 'running', statusText: '处理中', createdAt: '2026-03-28 18:00:00' },
]);

const columns = [
  { title: '任务ID', dataIndex: 'id', width: 100 },
  { title: '任务名称', dataIndex: 'name', width: 150 },
  { title: '数据模块', dataIndex: 'moduleText', width: 100 },
  { title: '格式', slotName: 'format', width: 80 },
  { title: '记录数', dataIndex: 'records', width: 80 },
  { title: '文件大小', dataIndex: 'size', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 120 },
];

const getStatusColor = (s: string) => ({ completed: 'green', running: 'blue', failed: 'red', pending: 'gray' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleDownload = (record: any) => {};
const handleRetry = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.data-export-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
