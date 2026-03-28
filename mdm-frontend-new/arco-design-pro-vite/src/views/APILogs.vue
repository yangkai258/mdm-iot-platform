<template>
  <div class="api-logs-container">
    <a-card>
      <template #title>
        <span>API调用日志</span>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="5">
            <a-input v-model="searchForm.appId" placeholder="应用ID" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.method" placeholder="请求方法" allow-clear>
              <a-option value="GET">GET</a-option>
              <a-option value="POST">POST</a-option>
              <a-option value="PUT">PUT</a-option>
              <a-option value="DELETE">DELETE</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-input v-model="searchForm.path" placeholder="API路径" allow-clear />
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #method="{ record }">
          <a-tag>{{ record.method }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status < 400 ? 'green' : 'red'">{{ record.status }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const searchForm = reactive({ appId: '', method: '', path: '' });

const data = ref([
  { id: 1, appId: 'APP001', appName: '我的IoT应用', method: 'POST', path: '/api/v1/devices', status: 200, latency: 45, time: '2026-03-28 18:50:00' },
  { id: 2, appId: 'APP001', appName: '我的IoT应用', method: 'GET', path: '/api/v1/devices/1', status: 200, latency: 32, time: '2026-03-28 18:49:30' },
  { id: 3, appId: 'APP002', appName: '宠物健康App', method: 'POST', path: '/api/v1/commands', status: 401, latency: 15, time: '2026-03-28 18:49:00' },
  { id: 4, appId: 'APP002', appName: '宠物健康App', method: 'GET', path: '/api/v1/pets', status: 200, latency: 56, time: '2026-03-28 18:48:30' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '应用ID', dataIndex: 'appId', width: 100 },
  { title: '应用名称', dataIndex: 'appName', width: 150 },
  { title: '方法', slotName: 'method', width: 80 },
  { title: '路径', dataIndex: 'path', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '延迟(ms)', dataIndex: 'latency', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const handleSearch = () => {};
</script>

<style scoped>
.api-logs-container { padding: 20px; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
