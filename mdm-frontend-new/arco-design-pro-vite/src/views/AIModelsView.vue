<template>
  <div class="ai-models-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>AI模型管理</span>
          <a-button type="primary" @click="handleUpload">
            <template #icon><icon-upload /></template>
            上传模型
          </a-button>
        </div>
      </template>
      
      <a-tabs default-active-key="models">
        <a-tab-pane key="models" title="模型列表">
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.status === 'online' ? 'green' : record.status === 'training' ? 'blue' : 'gray'">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            <template #type="{ record }">
              <a-tag>{{ record.type }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link v-if="record.status === 'online'" @click="handleRollback(record)">回滚</a-link>
                <a-link v-if="record.status !== 'training'" @click="handleDeploy(record)">部署</a-link>
                <a-link @click="handleView(record)">详情</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="training" title="训练任务">
          <a-table :columns="trainColumns" :data="trainingTasks" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.status === 'running' ? 'blue' : record.status === 'completed' ? 'green' : 'orange'">
                {{ record.status }}
              </a-tag>
            </template>
            <template #progress="{ record }">
              <a-progress :percent="record.progress" :color="record.status === 'failed' ? 'red' : 'blue'" />
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="ab" title="A/B实验">
          <a-table :columns="abColumns" :data="experiments" :pagination="pagination">
            <template #status="{ record }">
              <a-switch :checked="record.status === 'running'" disabled />
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleToggle(record)">{{ record.status === 'running' ? '暂停' : '启动' }}</a-link>
                <a-link @click="handleView(record)">详情</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'M001', name: 'PetBrain-v2.0', version: '2.0.5', type: '行为决策', status: 'online', accuracy: 96.5, latency: 45, deployedAt: '2026-03-20' },
  { id: 'M002', name: 'PetBrain-v1.9', version: '1.9.2', type: '行为决策', status: 'offline', accuracy: 94.2, latency: 52, deployedAt: '2026-03-15' },
  { id: 'M003', name: 'EmotionNet-v1.0', version: '1.0.3', type: '情绪识别', status: 'training', accuracy: 0, latency: 0, deployedAt: '-' },
  { id: 'M004', name: 'VoiceNet-v2.1', version: '2.1.0', type: '语音合成', status: 'online', accuracy: 98.1, latency: 28, deployedAt: '2026-03-18' },
]);

const trainingTasks = ref([
  { id: 'T001', name: 'PetBrain-v3.0训练', model: 'PetBrain', status: 'running', progress: 67, eta: '2小时', createdAt: '2026-03-28 08:00:00' },
  { id: 'T002', name: 'EmotionNet-v2.0训练', model: 'EmotionNet', status: 'completed', progress: 100, eta: '-', createdAt: '2026-03-27 10:00:00' },
  { id: 'T003', name: 'VoiceNet-v3.0训练', model: 'VoiceNet', status: 'pending', progress: 0, eta: '8小时', createdAt: '2026-03-28 14:00:00' },
]);

const experiments = ref([
  { id: 'E001', name: '新推荐算法测试', modelA: 'PetBrain-v2.0', modelB: 'PetBrain-v2.1', status: 'running', traffic: '50%', conversion: 12.5 },
  { id: 'E002', name: '情绪响应策略B测试', modelA: 'EmotionNet-v1.0', modelB: 'EmotionNet-v1.1', status: 'paused', traffic: '30%', conversion: 8.3 },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '模型ID', dataIndex: 'id', width: 80 },
  { title: '模型名称', dataIndex: 'name', width: 150 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '准确率', dataIndex: 'accuracy', width: 100 },
  { title: '延迟ms', dataIndex: 'latency', width: 80 },
  { title: '部署时间', dataIndex: 'deployedAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const trainColumns = [
  { title: '任务ID', dataIndex: 'id', width: 80 },
  { title: '任务名称', dataIndex: 'name', width: 200 },
  { title: '模型', dataIndex: 'model', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 150 },
  { title: '预计剩余', dataIndex: 'eta', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
];

const abColumns = [
  { title: '实验ID', dataIndex: 'id', width: 80 },
  { title: '实验名称', dataIndex: 'name', width: 200 },
  { title: 'A模型', dataIndex: 'modelA', width: 120 },
  { title: 'B模型', dataIndex: 'modelB', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '流量分配', dataIndex: 'traffic', width: 80 },
  { title: '转化率', dataIndex: 'conversion', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getStatusText = (status: string) => {
  const map: Record<string, string> = { online: '在线', offline: '离线', training: '训练中' };
  return map[status] || status;
};

const handleUpload = () => {};
const handleDeploy = (record: any) => {};
const handleRollback = (record: any) => {};
const handleView = (record: any) => {};
const handleToggle = (record: any) => {};
</script>

<style scoped>
.ai-models-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
