<template>
  <div class="model-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="AI模型" :value="15" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="活跃模型" :value="8" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="调用次数" :value="256000" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均延迟" :value="45" suffix="ms" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>AI模型管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加模型
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="models" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'error'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleConfig(record)">配置</a-link>
            <a-link @click="handleMetrics(record)">指标</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="configVisible" :title="currentModel.name + ' 配置'" :width="600">
      <a-form layout="vertical">
        <a-form-item label="API Endpoint">
          <a-input v-model="currentModel.endpoint" />
        </a-form-item>
        <a-form-item label="API Key">
          <a-input-password v-model="currentModel.apiKey" />
        </a-form-item>
        <a-form-item label="最大Token">
          <a-input-number v-model="currentModel.maxTokens" :min="100" />
        </a-form-item>
        <a-form-item label="温度参数">
          <a-slider v-model="currentModel.temperature" :min="0" :max="1" :step="0.1" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSaveConfig">保存配置</a-button>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const configVisible = ref(false);
const currentModel = ref<any>({});

const models = ref([
  { id: 1, name: 'GPT-4', provider: 'OpenAI', status: 'active', statusText: '活跃', calls: 125600, avgLatency: 450 },
  { id: 2, name: 'Claude-3', provider: 'Anthropic', status: 'active', statusText: '活跃', calls: 98600, avgLatency: 380 },
  { id: 3, name: 'Llama-2', provider: 'Meta', status: 'inactive', statusText: '未启用', calls: 0, avgLatency: 0 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '模型名称', dataIndex: 'name', width: 150 },
  { title: '提供商', dataIndex: 'provider', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '调用次数', dataIndex: 'calls', width: 120 },
  { title: '平均延迟', dataIndex: 'avgLatency', width: 120 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const handleCreate = () => {};
const handleConfig = (record: any) => { currentModel.value = { ...record, endpoint: 'https://api.openai.com', apiKey: '', maxTokens: 4096, temperature: 0.7 }; configVisible.value = true; };
const handleMetrics = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSaveConfig = () => { configVisible.value = false; };
</script>

<style scoped>
.model-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
