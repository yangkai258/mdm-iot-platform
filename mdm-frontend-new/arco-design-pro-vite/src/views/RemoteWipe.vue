<template>
  <div class="remote-wipe-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>远程锁定/擦除</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建任务
          </a-button>
        </div>
      </template>
      
      <a-alert type="warning" style="margin-bottom: 16px;">
        <template #title>危险操作</template>
        远程锁定和擦除是危险操作，可能导致设备数据永久丢失。请务必确认操作后果。
      </a-alert>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link v-if="record.status === 'pending'" @click="handleCancel(record)">取消</a-link>
            <a-link v-if="record.status === 'failed'" @click="handleRetry(record)">重试</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建任务弹窗 -->
    <a-modal v-model:visible="createVisible" title="创建远程操作任务" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="设备" required>
          <a-select v-model="form.deviceIds" multiple placeholder="选择目标设备" :max-tag-count="3">
            <a-option value="DEV001">DEV001 - 小黄(在线)</a-option>
            <a-option value="DEV002">DEV002 - 小红(在线)</a-option>
            <a-option value="DEV003">DEV003 - 小绿(离线)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="操作类型" required>
          <a-radio-group v-model="form.type">
            <a-radio value="lock">远程锁定</a-radio>
            <a-radio value="wipe_data">数据擦除</a-radio>
            <a-radio value="wipe_full">完全恢复出厂</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.type === 'lock'" label="锁定原因">
          <a-select v-model="form.reason" placeholder="选择锁定原因">
            <a-option value="lost">设备丢失</a-option>
            <a-option value="stolen">设备被盗</a-option>
            <a-option value="security">安全事件</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="form.type !== 'lock'" label="二次确认">
          <a-input v-model="form.confirmText" placeholder='请输入"CONFIRM"以确认操作' />
        </a-form-item>
        <a-form-item label="执行时间">
          <a-radio-group v-model="form.scheduleType">
            <a-radio value="now">立即执行</a-radio>
            <a-radio value="scheduled">定时执行</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.scheduleType === 'scheduled'" label="定时">
          <a-date-picker v-model="form.scheduledAt" show-time />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" placeholder="可选备注" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 任务详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="任务详情" :width="700">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="任务ID">{{ currentTask.id }}</a-descriptions-item>
        <a-descriptions-item label="操作类型">
          <a-tag :color="getTypeColor(currentTask.type)">{{ currentTask.typeText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentTask.status)">{{ currentTask.statusText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ currentTask.createdAt }}</a-descriptions-item>
        <a-descriptions-item label="执行时间">{{ currentTask.executedAt || '-' }}</a-descriptions-item>
        <a-descriptions-item label="完成时间">{{ currentTask.completedAt || '-' }}</a-descriptions-item>
        <a-descriptions-item label="目标设备" :span="2">
          <a-space wrap>
            <a-tag v-for="d in currentTask.devices" :key="d">{{ d }}</a-tag>
          </a-space>
        </a-descriptions-item>
        <a-descriptions-item label="执行结果" :span="2">{{ currentTask.result || '-' }}</a-descriptions-item>
      </a-descriptions>
      
      <a-divider>执行日志</a-divider>
      <a-timeline>
        <a-timeline-item v-for="log in currentTask.logs" :key="log.time" :color="log.type === 'success' ? 'green' : log.type === 'error' ? 'red' : 'gray'">
          <b>{{ log.time }}</b> - {{ log.message }}
        </a-timeline-item>
      </a-timeline>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'W001', type: 'lock', typeText: '远程锁定', deviceCount: 1, status: 'success', statusText: '成功', createdAt: '2026-03-28 10:00:00', executedAt: '2026-03-28 10:00:05', result: '设备已锁定' },
  { id: 'W002', type: 'wipe_data', typeText: '数据擦除', deviceCount: 2, status: 'running', statusText: '执行中', createdAt: '2026-03-28 09:30:00', executedAt: '2026-03-28 09:30:05', result: '' },
  { id: 'W003', type: 'wipe_full', typeText: '恢复出厂', deviceCount: 1, status: 'failed', statusText: '失败', createdAt: '2026-03-27 15:00:00', executedAt: '2026-03-27 15:00:03', result: '设备离线，操作失败' },
  { id: 'W004', type: 'lock', typeText: '远程锁定', deviceCount: 3, status: 'pending', statusText: '待执行', createdAt: '2026-03-28 11:00:00', executedAt: '-', result: '' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '任务ID', dataIndex: 'id', width: 100 },
  { title: '操作类型', slotName: 'type', width: 120 },
  { title: '设备数', dataIndex: 'deviceCount', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
  { title: '执行时间', dataIndex: 'executedAt', width: 160 },
  { title: '结果', dataIndex: 'result', width: 200 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const createVisible = ref(false);
const detailVisible = ref(false);
const currentTask = ref<any>({});

const form = reactive({
  deviceIds: [],
  type: 'lock',
  reason: '',
  confirmText: '',
  scheduleType: 'now',
  scheduledAt: '',
  note: '',
});

const getTypeColor = (t: string) => ({ lock: 'orange', wipe_data: 'red', wipe_full: '#F53F3F' }[t] || 'default');
const getStatusColor = (s: string) => ({ pending: 'blue', running: 'blue', success: 'green', failed: 'red' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleView = (record: any) => {
  currentTask.value = {
    ...record,
    devices: ['DEV001'],
    logs: [
      { time: '10:00:05', type: 'info', message: '任务开始执行' },
      { time: '10:00:05', type: 'success', message: '发送锁定指令到设备' },
      { time: '10:00:06', type: 'success', message: '设备确认锁定成功' },
    ],
  };
  detailVisible.value = true;
};
const handleCancel = (record: any) => {};
const handleRetry = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.remote-wipe-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
