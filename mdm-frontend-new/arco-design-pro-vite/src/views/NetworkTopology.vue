<template>
  <div class="network-topology-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="节点数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="在线" :value="120" status="success" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="离线" :value="8" status="error" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="带宽使用" :value="45" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>网络拓扑管理</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="topology" title="拓扑图">
          <a-card class="topology-card">
            <div class="topology-placeholder">
              <icon-api />
              <span>网络拓扑图</span>
            </div>
          </a-card>
        </a-tab-pane>
        
        <a-tab-pane key="nodes" title="节点列表">
          <a-table :columns="nodeColumns" :data="nodes" :pagination="pagination">
            <template #status="{ record }">
              <a-badge :status="record.status === 'online' ? 'success' : 'error'" :text="record.statusText" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="connections" title="连接管理">
          <a-table :columns="connColumns" :data="connections" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const nodes = ref([
  { id: 1, name: '中心节点-1', type: 'gateway', ip: '192.168.1.1', status: 'online', statusText: '在线', load: 45 },
  { id: 2, name: '边缘节点-1', type: 'edge', ip: '192.168.1.10', status: 'online', statusText: '在线', load: 30 },
]);

const connections = ref([
  { id: 1, source: '中心节点-1', target: '边缘节点-1', type: 'wired', bandwidth: '100Mbps', latency: '5ms' },
]);

const nodeColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '节点名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: 'IP地址', dataIndex: 'ip', width: 150 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '负载', dataIndex: 'load', width: 100 },
];

const connColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '源节点', dataIndex: 'source', width: 150 },
  { title: '目标节点', dataIndex: 'target', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '带宽', dataIndex: 'bandwidth', width: 100 },
  { title: '延迟', dataIndex: 'latency', width: 100 },
];
</script>

<style scoped>
.network-topology-container { padding: 20px; }
.topology-card { min-height: 400px; display: flex; align-items: center; justify-content: center; }
.topology-placeholder { text-align: center; color: #86909c; }
.topology-placeholder .arco-icon { font-size: 64px; margin-bottom: 16px; }
</style>
