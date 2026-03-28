<template>
  <div class="data-migration-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="迁移任务" :value="15" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="进行中" :value="2" status="processing" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="迁移记录" :value="568000" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>数据迁移工具</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建迁移
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="migrations" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="getStatusBadge(record.status)" :text="record.statusText" />
        </template>
        <template #progress="{ record }">
          <a-progress v-if="record.status === 'running'" :percent="record.progress" status="success" />
          <span v-else>{{ record.progress }}%</span>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="新建迁移任务" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="迁移类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="import">数据导入</a-option>
            <a-option value="export">数据导出</a-option>
            <a-option value="sync">数据同步</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="数据源">
          <a-select v-model="form.sourceType">
            <a-option value="csv">CSV文件</a-option>
            <a-option value="json">JSON文件</a-option>
            <a-option value="database">数据库</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标表">
          <a-select v-model="form.targetTable">
            <a-option value="devices">设备表</a-option>
            <a-option value="members">会员表</a-option>
            <a-option value="pets">宠物表</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="映射规则">
          <a-textarea v-model="form.mapping" :rows="4" placeholder="字段映射JSON" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ type: '', sourceType: '', targetTable: '', mapping: '' });

const migrations = ref([
  { id: 1, name: '设备数据导入', type: 'import', typeText: '导入', source: 'csv', targetTable: 'devices', status: 'completed', statusText: '完成', progress: 100 },
  { id: 2, name: '会员数据同步', type: 'sync', typeText: '同步', source: 'database', targetTable: 'members', status: 'running', statusText: '进行中', progress: 68 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '任务名称', dataIndex: 'name', width: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '目标表', dataIndex: 'targetTable', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 150 },
];

const getStatusBadge = (s: string) => ({ pending: 'default', running: 'processing', completed: 'success', failed: 'error' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.data-migration-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
