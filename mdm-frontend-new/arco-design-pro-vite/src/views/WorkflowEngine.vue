<template>
  <div class="workflow-engine-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="工作流" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="运行中" :value="12" status="processing" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已完成" :value="12580" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="成功率" :value="98.5" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>工作流引擎</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建工作流
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="workflows" title="工作流列表">
          <a-table :columns="workflowColumns" :data="workflows" :pagination="pagination">
            <template #status="{ record }">
              <a-badge :status="record.enabled ? 'success' : 'default'" :text="record.enabled ? '启用' : '禁用'" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="runs" title="运行记录">
          <a-table :columns="runColumns" :data="runs" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建工作流" :width="900" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="工作流名称" required>
          <a-input v-model="form.name" placeholder="请输入工作流名称" />
        </a-form-item>
        <a-form-item label="工作流定义">
          <a-textarea v-model="form.definition" :rows="10" placeholder="JSON格式工作流定义" />
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

const form = reactive({ name: '', definition: '' });

const workflows = ref([
  { id: 1, name: '设备注册流程', steps: 5, enabled: true, totalRuns: 2560 },
  { id: 2, name: '会员注册流程', steps: 8, enabled: true, totalRuns: 5680 },
]);

const runs = ref([
  { id: 1, workflowName: '设备注册流程', status: 'completed', statusText: '完成', duration: '5s', startTime: '2026-03-28 18:00:00' },
]);

const workflowColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '工作流名称', dataIndex: 'name', width: 200 },
  { title: '步骤数', dataIndex: 'steps', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '总运行次数', dataIndex: 'totalRuns', width: 120 },
];

const runColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '工作流', dataIndex: 'workflowName', width: 150 },
  { title: '状态', dataIndex: 'statusText', width: 100 },
  { title: '耗时', dataIndex: 'duration', width: 80 },
  { title: '开始时间', dataIndex: 'startTime', width: 160 },
];

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.workflow-engine-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
