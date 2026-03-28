<template>
  <div class="ai-training-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="训练任务" :value="stats.totalTasks" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="运行中" :value="stats.running" status="processing" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="完成" :value="stats.completed" status="success" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="失败" :value="stats.failed" status="error" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>AI训练任务</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建任务
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress v-if="record.status === 'running'" :percent="record.progress" :status="record.progress === 100 ? 'success' : 'normal'" />
          <span v-else>{{ record.progress }}%</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link v-if="record.status === 'completed'" @click="handleDeploy(record)">部署</a-link>
            <a-link v-if="record.status === 'failed'" @click="handleRetry(record)">重试</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建任务弹窗 -->
    <a-modal v-model:visible="createVisible" title="新建训练任务" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="任务名称" required>
          <a-input v-model="form.name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-form-item label="模型类型">
          <a-select v-model="form.modelType" placeholder="选择模型类型">
            <a-option value="nlu">NLU意图识别</a-option>
            <a-option value="tts">TTS语音合成</a-option>
            <a-option value="cv">视觉识别</a-option>
            <a-option value="emotion">情绪识别</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="训练数据">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="训练配置">
          <a-space vertical>
            <a-space>
              <span>Epochs:</span>
              <a-input-number v-model="form.config.epochs" :min="1" :max="1000" />
            </a-space>
            <a-space>
              <span>Learning Rate:</span>
              <a-input-number v-model="form.config.learningRate" :min="0.0001" :max="1" :precision="4" />
            </a-space>
            <a-space>
              <span>Batch Size:</span>
              <a-input-number v-model="form.config.batchSize" :min="1" :max="256" />
            </a-space>
          </a-space>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const stats = reactive({ totalTasks: 28, running: 2, completed: 24, failed: 2 });

const data = ref([
  { id: 1, name: 'Pet-NLU-v3.1', type: 'nlu', typeText: 'NLU训练', modelVersion: 'v3.1', status: 'running', statusText: '运行中', progress: 68, config: { epochs: 100, learningRate: 0.001, batchSize: 32 }, accuracy: 0, startTime: '2026-03-28 18:00:00', endTime: null },
  { id: 2, name: 'Pet-NLU-v3.0', type: 'nlu', typeText: 'NLU训练', modelVersion: 'v3.0', status: 'completed', statusText: '已完成', progress: 100, config: { epochs: 100, learningRate: 0.001, batchSize: 32 }, accuracy: 94.5, startTime: '2026-03-27 10:00:00', endTime: '2026-03-27 18:00:00' },
  { id: 3, name: 'Emotion-v2.0', type: 'emotion', typeText: '情绪识别', modelVersion: 'v2.0', status: 'failed', statusText: '失败', progress: 45, config: { epochs: 50, learningRate: 0.001, batchSize: 16 }, accuracy: 0, startTime: '2026-03-26 10:00:00', endTime: '2026-03-26 12:00:00', error: '数据不足' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const createVisible = ref(false);

const form = reactive({
  name: '', modelType: 'nlu', note: '',
  config: { epochs: 100, learningRate: 0.001, batchSize: 32 }
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '任务名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '目标版本', dataIndex: 'modelVersion', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 150 },
  { title: '准确率', dataIndex: 'accuracy', width: 100 },
  { title: '开始时间', dataIndex: 'startTime', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const getStatusColor = (s: string) => ({ running: 'blue', completed: 'green', failed: 'red', pending: 'gray' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleView = (record: any) => {};
const handleDeploy = (record: any) => {};
const handleRetry = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.ai-training-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
