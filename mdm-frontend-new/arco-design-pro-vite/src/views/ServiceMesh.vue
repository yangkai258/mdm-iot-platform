<template>
  <div class="service-mesh-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="服务数" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="边车代理" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="请求数" :value="125600" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="成功率" :value="99.5" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>服务网格管理</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="services" title="服务列表">
          <a-table :columns="serviceColumns" :data="services" :pagination="pagination">
            <template #status="{ record }">
              <a-badge :status="record.status === 'healthy' ? 'success' : 'error'" :text="record.statusText" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="traffic" title="流量管理">
          <a-table :columns="trafficColumns" :data="traffic" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const services = ref([
  { id: 1, name: 'api-gateway', version: 'v2.0.0', status: 'healthy', statusText: '健康', requests: 56800, latency: '5ms' },
  { id: 2, name: 'user-service', version: 'v1.5.0', status: 'healthy', statusText: '健康', requests: 35600, latency: '3ms' },
]);

const traffic = ref([
  { id: 1, source: 'api-gateway', destination: 'user-service', requests: 125600, successRate: 99.5 },
]);

const serviceColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '服务名称', dataIndex: 'name', width: 200 },
  { title: '版本', dataIndex: 'version', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '请求数', dataIndex: 'requests', width: 120 },
  { title: '延迟', dataIndex: 'latency', width: 100 },
];

const trafficColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '源服务', dataIndex: 'source', width: 150 },
  { title: '目标服务', dataIndex: 'destination', width: 150 },
  { title: '请求数', dataIndex: 'requests', width: 120 },
  { title: '成功率', dataIndex: 'successRate', width: 100 },
];
</script>

<style scoped>
.service-mesh-container { padding: 20px; }
</style>
