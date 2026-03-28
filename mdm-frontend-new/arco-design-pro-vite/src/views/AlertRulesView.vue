<template>
  <div class="alert-rules-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>告警规则</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增规则
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">{{ record.level }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-switch :checked="record.enabled" disabled />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
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
  { id: 'R001', name: '设备离线告警', type: 'device_offline', level: 'critical', enabled: true, threshold: '30分钟', notification: 'SMS,Email', createdAt: '2026-03-01' },
  { id: 'R002', name: '电量低告警', type: 'low_battery', level: 'warning', enabled: true, threshold: '20%', notification: 'SMS', createdAt: '2026-03-01' },
  { id: 'R003', name: '异常行为告警', type: 'abnormal_behavior', level: 'critical', enabled: true, threshold: '检测到', notification: 'Email,Webhook', createdAt: '2026-03-15' },
  { id: 'R004', name: '固件更新提醒', type: 'firmware_update', level: 'info', enabled: false, threshold: '新版本', notification: 'Email', createdAt: '2026-03-20' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '规则ID', dataIndex: 'id', width: 80 },
  { title: '规则名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 120 },
  { title: '级别', slotName: 'level', width: 100 },
  { title: '触发条件', dataIndex: 'threshold', width: 120 },
  { title: '通知方式', dataIndex: 'notification', width: 150 },
  { title: '启用', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getLevelColor = (level: string) => {
  const map: Record<string, string> = { critical: 'red', warning: 'orange', info: 'blue' };
  return map[level] || 'default';
};

const handleCreate = () => {};
const handleEdit = (record: any) => {};
const handleDelete = (record: any) => {};
</script>

<style scoped>
.alert-rules-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
