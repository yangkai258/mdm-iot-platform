<template>
  <div class="cert-man-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="证书总数" :value="36" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="有效证书" :value="32" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="即将过期" :value="3" status="warning" /></a-card>
      </a-col>
    </a-row>
    <a-card>
      <template #title>证书管理</template>
      <a-table :columns="columns" :data="data" :pagination="pagination" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const data = ref([
  { id: 1, name: 'api.example.com', type: 'SSL', issuer: 'Let\'s Encrypt', expireDate: '2026-06-28', status: 'valid' },
]);
const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '域名', dataIndex: 'name', width: 200 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '颁发者', dataIndex: 'issuer', width: 150 },
  { title: '过期日期', dataIndex: 'expireDate', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100 },
];
</script>

<style scoped>
.cert-man-container { padding: 20px; }
</style>
