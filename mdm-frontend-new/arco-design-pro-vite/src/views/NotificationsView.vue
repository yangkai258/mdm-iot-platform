<template>
  <div class="notifications-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>通知渠道</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加渠道
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.type }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-switch :checked="record.enabled" disabled />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleTest(record)">测试</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'N001', name: '系统通知', type: 'Email', config: 'admin@example.com', enabled: true, lastSent: '2026-03-28 09:00:00' },
  { id: 'N002', name: '告警通知', type: 'SMS', config: '138****1234', enabled: true, lastSent: '2026-03-28 08:30:00' },
  { id: 'N003', name: '营销通知', type: 'Webhook', config: 'https://api.example.com/webhook', enabled: true, lastSent: '2026-03-27 10:00:00' },
  { id: 'N004', name: '企业微信', type: 'Webhook', config: 'wecom.example.com', enabled: false, lastSent: '-' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '渠道ID', dataIndex: 'id', width: 80 },
  { title: '渠道名称', dataIndex: 'name', width: 120 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '配置', dataIndex: 'config', width: 250 },
  { title: '启用', slotName: 'status', width: 80 },
  { title: '最后发送', dataIndex: 'lastSent', width: 160 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getTypeColor = (type: string) => {
  const map: Record<string, string> = { Email: 'blue', SMS: 'green', Webhook: 'purple' };
  return map[type] || 'default';
};

const handleCreate = () => {};
const handleEdit = (record: any) => {};
const handleTest = (record: any) => {};
</script>

<style scoped>
.notifications-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
