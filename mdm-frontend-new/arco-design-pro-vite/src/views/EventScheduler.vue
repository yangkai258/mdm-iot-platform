<template>
  <div class="event-scheduler-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="定时事件" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="执行次数" :value="56800" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="成功率" :value="99" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="下次执行" :value="5" suffix="分钟" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>定时事件调度器</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建事件
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="events" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.enabled ? 'success' : 'default'" :text="record.enabled ? '启用' : '禁用'" />
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建定时事件" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="事件名称" required>
          <a-input v-model="form.name" placeholder="请输入事件名称" />
        </a-form-item>
        <a-form-item label="执行周期">
          <a-select v-model="form.schedule">
            <a-option value="once">一次性</a-option>
            <a-option value="daily">每天</a-option>
            <a-option value="weekly">每周</a-option>
            <a-option value="cron">Cron表达式</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="执行时间">
          <a-time-picker v-model="form.executeTime" format="HH:mm" />
        </a-form-item>
        <a-form-item label="执行动作">
          <a-textarea v-model="form.action" :rows="3" placeholder="JSON格式动作" />
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

const form = reactive({ name: '', schedule: '', executeTime: '', action: '' });

const events = ref([
  { id: 1, name: '数据备份', type: 'backup', typeText: '备份', schedule: 'daily', executeTime: '02:00', nextRun: '2026-03-29 02:00:00', enabled: true },
  { id: 2, name: '报表生成', type: 'report', typeText: '报表', schedule: 'weekly', executeTime: '09:00', nextRun: '2026-04-01 09:00:00', enabled: true },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '事件名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '周期', dataIndex: 'schedule', width: 100 },
  { title: '执行时间', dataIndex: 'executeTime', width: 100 },
  { title: '下次执行', dataIndex: 'nextRun', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
];

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.event-scheduler-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
