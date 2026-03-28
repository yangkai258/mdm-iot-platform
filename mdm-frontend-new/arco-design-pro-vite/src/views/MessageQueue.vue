<template>
  <div class="message-queue-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="队列数" :value="15" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="消息总数" :value="1256000" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="消费中" :value="856" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="堆积" :value="128" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>消息队列监控</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="queues" title="队列列表">
          <a-table :columns="queueColumns" :data="queues" :pagination="pagination">
            <template #status="{ record }">
              <a-badge :status="record.status === 'healthy' ? 'success' : 'error'" :text="record.statusText" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="messages" title="消息监控">
          <a-table :columns="messageColumns" :data="messages" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const queues = ref([
  { id: 1, name: 'device.command', type: 'Command', messages: 125600, consumers: 5, status: 'healthy', statusText: '健康' },
  { id: 2, name: 'alert.notification', type: 'Notification', messages: 25600, consumers: 3, status: 'healthy', statusText: '健康' },
  { id: 3, name: 'data.sync', type: 'Sync', messages: 56800, consumers: 2, status: 'warning', statusText: '堆积' },
]);

const messages = ref([
  { id: 1, queue: 'device.command', messageId: 'msg_001', payload: '{}', status: 'sent', time: '2026-03-28 18:00:00' },
]);

const queueColumns = [
  { title: '队列名称', dataIndex: 'name', width: 200 },
  { title: '类型', dataIndex: 'type', width: 120 },
  { title: '消息数', dataIndex: 'messages', width: 120 },
  { title: '消费者', dataIndex: 'consumers', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
];

const messageColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '队列', dataIndex: 'queue', width: 150 },
  { title: '消息ID', dataIndex: 'messageId', width: 150 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];
</script>

<style scoped>
.message-queue-container { padding: 20px; }
</style>
