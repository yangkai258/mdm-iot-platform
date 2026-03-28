<template>
  <div class="pressure-test-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="测试任务" :value="15" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="进行中" :value="2" status="processing" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="通过率" :value="94" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="QPS峰值" :value="10000" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>压力测试</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建任务
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link v-if="record.status === 'draft'" @click="handleRun(record)">运行</a-link>
            <a-link @click="handleView(record)">详情</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建压力测试" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="任务名称" required>
          <a-input v-model="form.name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-form-item label="目标API">
          <a-input v-model="form.targetUrl" placeholder="https://api.example.com/..." />
        </a-form-item>
        <a-form-item label="并发数">
          <a-input-number v-model="form.concurrentUsers" :min="1" :max="10000" />
        </a-form-item>
        <a-form-item label="持续时间">
          <a-input-number v-model="form.duration" :min="1" /> 秒
        </a-form-item>
        <a-form-item label="QPS目标">
          <a-input-number v-model="form.targetQps" :min="1" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, name: '设备列表API压测', targetUrl: '/api/v1/devices', concurrentUsers: 100, duration: 60, targetQps: 1000, status: 'completed', statusText: '已完成', passRate: 98, maxQps: 1200, time: '2026-03-28 10:00:00' },
  { id: 2, name: '会员创建API压测', targetUrl: '/api/v1/members', concurrentUsers: 50, duration: 30, targetQps: 500, status: 'running', statusText: '运行中', passRate: 94, maxQps: 580, time: '2026-03-28 18:00:00' },
  { id: 3, name: 'OTA下发压测', targetUrl: '/api/v1/ota/deploy', concurrentUsers: 20, duration: 60, targetQps: 200, status: 'draft', statusText: '草稿', passRate: 0, maxQps: 0, time: null },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const createVisible = ref(false);

const form = reactive({ name: '', targetUrl: '', concurrentUsers: 100, duration: 60, targetQps: 1000 });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '任务名称', dataIndex: 'name', width: 200 },
  { title: '目标API', dataIndex: 'targetUrl', width: 200 },
  { title: '并发', dataIndex: 'concurrentUsers', width: 80 },
  { title: 'QPS目标', dataIndex: 'targetQps', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '通过率', dataIndex: 'passRate', width: 100 },
  { title: 'QPS峰值', dataIndex: 'maxQps', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '操作', slotName: 'actions', width: 120 },
];

const getStatusColor = (s: string) => ({ draft: 'gray', running: 'blue', completed: 'green', failed: 'red' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleRun = (record: any) => {};
const handleView = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.pressure-test-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
