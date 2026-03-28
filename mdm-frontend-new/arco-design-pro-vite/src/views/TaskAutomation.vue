<template>
  <div class="task-automation-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="自动化任务" :value="128" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="执行次数" :value="56800" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="成功率" :value="98.5" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>任务自动化</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建任务
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="tasks" :loading="loading" :pagination="pagination">
        <template #trigger="{ record }">
          <a-tag>{{ record.triggerText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.enabled ? 'success' : 'default'" :text="record.enabled ? '启用' : '禁用'" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleRun(record)">执行</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑任务' : '创建任务'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="任务名称" required>
          <a-input v-model="form.name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-form-item label="触发条件">
          <a-select v-model="form.trigger" placeholder="选择触发条件">
            <a-option value="schedule">定时触发</a-option>
            <a-option value="device_online">设备上线</a-option>
            <a-option value="alert">告警触发</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="执行动作">
          <a-textarea v-model="form.action" :rows="3" placeholder="JSON格式动作" />
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

const form = reactive({ name: '', trigger: '', action: '', enabled: true });

const tasks = ref([
  { id: 1, name: '设备离线告警', trigger: 'device_online', triggerText: '设备上线', lastRun: '2026-03-28 18:00:00', enabled: true },
  { id: 2, name: '每日数据备份', trigger: 'schedule', triggerText: '定时触发', lastRun: '2026-03-28 02:00:00', enabled: true },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '任务名称', dataIndex: 'name', width: 200 },
  { title: '触发条件', slotName: 'trigger', width: 120 },
  { title: '最后执行', dataIndex: 'lastRun', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleRun = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.task-automation-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
