<template>
  <div class="webhook-debug-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="调试次数" :value="568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="失败次数" :value="25" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均耗时" :value="120" suffix="ms" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="成功率" :value="96" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>Webhook调试工具</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="debug" title="在线调试">
          <a-form layout="vertical" style="max-width: 800px;">
            <a-form-item label="Webhook URL">
              <a-input v-model="debugForm.url" placeholder="https://example.com/webhook" />
            </a-form-item>
            <a-form-item label="请求方法">
              <a-select v-model="debugForm.method">
                <a-option value="POST">POST</a-option>
                <a-option value="GET">GET</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="请求头">
              <a-textarea v-model="debugForm.headers" :rows="3" placeholder='{"Content-Type": "application/json"}' />
            </a-form-item>
            <a-form-item label="请求体">
              <a-textarea v-model="debugForm.body" :rows="6" placeholder="JSON格式请求体" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="handleSend">发送请求</a-button>
                <a-button @click="handleClear">清空</a-button>
              </a-space>
            </a-form-item>
          </a-form>
          
          <a-divider>响应结果</a-divider>
          <a-card v-if="response" size="small">
            <a-descriptions :column="1">
              <a-descriptions-item label="状态码">{{ response.status }}</a-descriptions-item>
              <a-descriptions-item label="耗时">{{ response.duration }}ms</a-descriptions-item>
            </a-descriptions>
            <pre class="response-body">{{ response.body }}</pre>
          </a-card>
        </a-tab-pane>
        
        <a-tab-pane key="history" title="调试历史">
          <a-table :columns="historyColumns" :data="history" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const debugForm = reactive({ url: '', method: 'POST', headers: '', body: '' });
const response = ref<any>(null);

const history = ref([
  { id: 1, url: 'https://api.example.com/webhook', method: 'POST', status: 200, duration: 120, time: '2026-03-28 18:00:00' },
]);

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: 'URL', dataIndex: 'url', ellipsis: true },
  { title: '方法', dataIndex: 'method', width: 80 },
  { title: '状态码', dataIndex: 'status', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 80 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const handleSend = () => {
  response.value = { status: 200, duration: 120, body: '{"code": 0, "message": "success"}' };
};
const handleClear = () => { debugForm.url = ''; debugForm.headers = ''; debugForm.body = ''; response.value = null; };
</script>

<style scoped>
.webhook-debug-container { padding: 20px; }
.response-body { background: #f5f5f5; padding: 12px; border-radius: 4px; overflow-x: auto; }
</style>
