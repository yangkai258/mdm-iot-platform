<template>
  <div class="webhooks-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="今日事件" :value="stats.todayEvents" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="成功率" :value="stats.successRate" suffix="%" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="平均延迟" :value="stats.avgLatency" suffix="ms" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="失败数" :value="stats.failedCount" status="error" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>Webhook管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增订阅
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #event="{ record }">
          <a-tag>{{ record.eventType }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleTest(record)">测试</a-link>
            <a-link @click="handleLogs(record)">日志</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑订阅' : '新增订阅'" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="订阅名称" required>
          <a-input v-model="form.name" placeholder="如: 订单支付通知" />
        </a-form-item>
        <a-form-item label="订阅事件" required>
          <a-select v-model="form.eventTypes" multiple placeholder="选择事件类型">
            <a-option value="device.online">设备上线</a-option>
            <a-option value="device.offline">设备离线</a-option>
            <a-option value="device.alert">设备告警</a-option>
            <a-option value="order.paid">订单支付</a-option>
            <a-option value="order.completed">订单完成</a-option>
            <a-option value="subscription.renewed">订阅续费</a-option>
            <a-option value="subscription.expired">订阅过期</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="Webhook地址" required>
          <a-input v-model="form.url" placeholder="https://your-server.com/webhook" />
        </a-form-item>
        <a-form-item label="Secret密钥">
          <a-input-password v-model="form.secret" placeholder="用于签名验证" />
        </a-form-item>
        <a-form-item label="重试次数">
          <a-input-number v-model="form.retryCount" :min="0" :max="10" />
        </a-form-item>
        <a-form-item label="超时时间">
          <a-input-number v-model="form.timeout" :min="1" :max="60" /> 秒
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 日志弹窗 -->
    <a-modal v-model:visible="logsVisible" title="Webhook日志" :width="900">
      <a-table :columns="logColumns" :data="webhookLogs" :pagination="paginationSmall">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const stats = reactive({
  todayEvents: 1583,
  successRate: 99.2,
  avgLatency: 245,
  failedCount: 12,
});

const data = ref([
  { id: 1, name: '订单支付通知', eventTypes: ['order.paid'], url: 'https://api.example.com/webhook/order', secret: '****', retryCount: 3, timeout: 30, enabled: true },
  { id: 2, name: '设备告警通知', eventTypes: ['device.alert'], url: 'https://api.example.com/webhook/alert', secret: '****', retryCount: 5, timeout: 10, enabled: true },
  { id: 3, name: '订阅变更通知', eventTypes: ['subscription.renewed', 'subscription.expired'], url: 'https://api.example.com/webhook/sub', secret: '****', retryCount: 3, timeout: 30, enabled: true },
]);

const webhookLogs = ref([
  { id: 1, eventType: 'order.paid', url: 'https://api...', status: 200, statusText: '成功', latency: 120, time: '2026-03-28 18:00:00' },
  { id: 2, eventType: 'device.alert', url: 'https://api...', status: 500, statusText: '失败', latency: 5000, time: '2026-03-28 17:55:00' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const paginationSmall = reactive({ current: 1, pageSize: 10, total: 2 });
const editVisible = ref(false);
const logsVisible = ref(false);
const isEdit = ref(false);

const form = reactive({
  name: '', eventTypes: [], url: '', secret: '', retryCount: 3, timeout: 30, enabled: true,
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '订阅名称', dataIndex: 'name', width: 150 },
  { title: '事件类型', slotName: 'event', width: 200 },
  { title: 'URL', dataIndex: 'url', width: 250 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const logColumns = [
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '事件', dataIndex: 'eventType', width: 120 },
  { title: '状态码', slotName: 'status', width: 80 },
  { title: '延迟', dataIndex: 'latency', width: 80 },
];

const getStatusColor = (s: number) => s < 300 ? 'green' : s < 500 ? 'orange' : 'red';

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleTest = (record: any) => {};
const handleLogs = (record: any) => { logsVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: (closed: boolean) => void) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.webhooks-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
