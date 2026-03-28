<template>
  <div class="load-balancer-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="负载均衡器" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="后端服务器" :value="48" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="活跃连接" :value="8560" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="带宽" :value="10" suffix="Gbps" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>负载均衡管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建配置
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="balancers" :loading="loading" :pagination="pagination">
        <template #algorithm="{ record }">
          <a-tag>{{ record.algorithmText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'default'" :text="record.statusText" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });

const balancers = ref([
  { id: 1, name: 'API负载均衡', vip: '192.168.1.100', algorithm: 'roundrobin', algorithmText: '轮询', backends: 8, status: 'active', statusText: '活跃' },
  { id: 2, name: 'Web负载均衡', vip: '192.168.1.101', algorithm: 'leastconn', algorithmText: '最小连接', backends: 6, status: 'active', statusText: '活跃' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: 'VIP', dataIndex: 'vip', width: 150 },
  { title: '算法', slotName: 'algorithm', width: 120 },
  { title: '后端数', dataIndex: 'backends', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
];

const handleCreate = () => {};
</script>

<style scoped>
.load-balancer-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
