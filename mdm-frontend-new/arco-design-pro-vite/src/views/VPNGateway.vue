<template>
  <div class="vpn-gateway-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="VPN网关" :value="8" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="在线用户" :value="156" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="流量使用" :value="25" suffix="GB" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均延迟" :value="15" suffix="ms" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>VPN网关管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建网关
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="gateways" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'default'" :text="record.statusText" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang:ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });

const gateways = ref([
  { id: 1, name: 'VPN-北京', ip: '203.0.113.1', type: 'IPSec', users: 86, status: 'active', statusText: '在线', uptime: '99.9%' },
  { id: 2, name: 'VPN-上海', ip: '203.0.113.2', type: 'SSL', users: 70, status: 'active', statusText: '在线', uptime: '99.8%' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '网关名称', dataIndex: 'name', width: 150 },
  { title: 'IP地址', dataIndex: 'ip', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '在线用户', dataIndex: 'users', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '运行时间', dataIndex: 'uptime', width: 100 },
];

const handleCreate = () => {};
</script>

<style scoped>
.vpn-gateway-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
