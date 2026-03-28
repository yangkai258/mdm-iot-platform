<template>
  <div class="cicd-integration-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="CI/CD流水线" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="运行中" :value="2" status="processing" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="成功率" :value="94" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="总构建次数" :value="1289" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>CI/CD集成</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建流水线
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="pipelines" title="流水线">
          <a-table :columns="pipelineColumns" :data="pipelines" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link v-if="record.status === 'idle'" @click="handleRun(record)">运行</a-link>
                <a-link @click="handleView(record)">详情</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="runs" title="构建记录">
          <a-table :columns="runColumns" :data="runs" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建流水线" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="流水线名称" required>
          <a-input v-model="form.name" placeholder="请输入流水线名称" />
        </a-form-item>
        <a-form-item label="触发方式">
          <a-checkbox-group v-model="form.triggers">
            <a-checkbox value="push">代码推送</a-checkbox>
            <a-checkbox value="schedule">定时触发</a-checkbox>
            <a-checkbox value="manual">手动触发</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="阶段">
          <a-select v-model="form.stages" multiple placeholder="选择阶段">
            <a-option value="build">构建</a-option>
            <a-option value="test">测试</a-option>
            <a-option value="deploy">部署</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 12 });
const createVisible = ref(false);

const form = reactive({ name: '', triggers: ['push'], stages: ['build', 'test'] });

const pipelines = ref([
  { id: 1, name: '后端构建', trigger: 'push', stages: '构建/测试/部署', status: 'idle', statusText: '空闲', lastRun: '2026-03-28 10:00:00' },
  { id: 2, name: '前端构建', trigger: 'push', stages: '构建/部署', status: 'running', statusText: '运行中', lastRun: '2026-03-28 18:00:00' },
]);

const runs = ref([
  { id: 1, pipeline: '后端构建', status: 'success', duration: '5m30s', commit: 'abc123', time: '2026-03-28 10:00:00' },
  { id: 2, pipeline: '前端构建', status: 'failed', duration: '3m00s', commit: 'def456', time: '2026-03-28 09:00:00' },
]);

const pipelineColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '流水线名称', dataIndex: 'name', width: 150 },
  { title: '触发方式', dataIndex: 'trigger', width: 100 },
  { title: '阶段', dataIndex: 'stages', width: 150 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '最后运行', dataIndex: 'lastRun', width: 160 },
  { title: '操作', slotName: 'actions', width: 150 },
];

const runColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '流水线', dataIndex: 'pipeline', width: 150 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '耗时', dataIndex: 'duration', width: 100 },
  { title: 'Commit', dataIndex: 'commit', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const getStatusColor = (s: string) => ({ idle: 'gray', running: 'blue', success: 'green', failed: 'red' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleRun = (record: any) => {};
const handleView = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.cicd-integration-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
