<template>
  <div class="api-explorer-container">
    <a-card>
      <template #title>
        <span>API Explorer</span>
      </template>
      
      <a-form layout="vertical" style="max-width: 800px;">
        <a-form-item label="API路径">
          <a-input v-model="api.path" placeholder="/api/v1/devices" />
        </a-form-item>
        <a-form-item label="请求方法">
          <a-radio-group v-model="api.method">
            <a-radio value="GET">GET</a-radio>
            <a-radio value="POST">POST</a-radio>
            <a-radio value="PUT">PUT</a-radio>
            <a-radio value="DELETE">DELETE</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="Headers">
          <a-input v-model="api.headers" placeholder='{"Content-Type": "application/json"}' :rows="2" />
        </a-form-item>
        <a-form-item label="请求参数">
          <a-textarea v-model="api.body" placeholder='{"key": "value"}' :rows="4" />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSend">发送请求</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>

      <a-divider>响应</a-divider>
      <div v-if="response" class="response-area">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="状态码">
            <a-tag :color="response.status < 400 ? 'green' : 'red'">{{ response.status }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="耗时">{{ response.duration }}ms</a-descriptions-item>
          <a-descriptions-item label="响应内容">
            <pre class="response-json">{{ response.data }}</pre>
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </a-card>

    <a-card title="API文档" style="margin-top: 16px;">
      <a-tabs>
        <a-tab-pane v-for="group in apiDocs" :key="group.name" :tab="group.name" :title="group.name">
          <a-collapse>
            <a-collapse-item v-for="api in group.apis" :key="api.path" :header="api.method + ' ' + api.path">
              <div>
                <div><b>描述:</b> {{ api.description }}</div>
                <div><b>参数:</b> {{ api.params || '无' }}</div>
                <a-button type="text" size="small" @click="handleTryIt(api)">试用</a-button>
              </div>
            </a-collapse-item>
          </a-collapse>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const api = reactive({ path: '', method: 'GET', headers: '{"Content-Type": "application/json"}', body: '' });
const response = ref<any>(null);

const apiDocs = ref([
  {
    name: '设备API',
    apis: [
      { method: 'GET', path: '/api/v1/devices', description: '获取设备列表', params: 'page, page_size' },
      { method: 'GET', path: '/api/v1/devices/:id', description: '获取设备详情', params: 'device_id' },
      { method: 'POST', path: '/api/v1/devices', description: '创建设备', params: 'device_id, hardware_model' },
      { method: 'PUT', path: '/api/v1/devices/:id', description: '更新设备', params: 'device_id' },
      { method: 'DELETE', path: '/api/v1/devices/:id', description: '删除设备', params: 'device_id' },
    ],
  },
  {
    name: '会员API',
    apis: [
      { method: 'GET', path: '/api/v1/members', description: '获取会员列表', params: 'page, page_size' },
      { method: 'GET', path: '/api/v1/members/:id', description: '获取会员详情', params: 'id' },
      { method: 'POST', path: '/api/v1/members', description: '创建会员', params: 'member_name, phone' },
    ],
  },
]);

const handleSend = async () => {
  // 实际实现会调用后端API
  response.value = {
    status: 200,
    duration: 125,
    data: JSON.stringify({ code: 0, message: 'success', data: [] }, null, 2),
  };
};

const handleReset = () => {
  api.path = '';
  api.method = 'GET';
  api.headers = '{"Content-Type": "application/json"}';
  api.body = '';
  response.value = null;
};

const handleTryIt = (item: any) => {
  api.method = item.method;
  api.path = item.path;
};
</script>

<style scoped>
.api-explorer-container { padding: 20px; }
.response-area { background: #f7f8fa; padding: 16px; border-radius: 4px; }
.response-json { background: #1e1e1e; color: #d4d4d4; padding: 16px; border-radius: 4px; overflow: auto; max-height: 400px; font-family: monospace; font-size: 12px; }
</style>
