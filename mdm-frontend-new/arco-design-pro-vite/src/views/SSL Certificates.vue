<template>
  <div class="ssl-certificates-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="证书总数" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="有效证书" :value="32" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="即将过期" :value="3" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="覆盖率" :value="100" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>SSL证书管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加证书
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="certificates" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="添加证书" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="证书名称" required>
          <a-input v-model="form.name" placeholder="请输入证书名称" />
        </a-form-item>
        <a-form-item label="域名">
          <a-input v-model="form.domain" placeholder="example.com" />
        </a-form-item>
        <a-form-item label="证书文件">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="私钥文件">
          <a-upload action="#" :limit="1" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const createVisible = ref(false);

const form = reactive({ name: '', domain: '', certFile: '', keyFile: '' });

const certificates = ref([
  { id: 1, name: 'api.example.com', domain: 'api.example.com', issuer: 'Let\'s Encrypt', expireDate: '2026-06-28', status: 'valid', statusText: '有效' },
  { id: 2, name: 'web.example.com', domain: 'web.example.com', issuer: 'Let\'s Encrypt', expireDate: '2026-04-15', status: 'expiring', statusText: '即将过期' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '证书名称', dataIndex: 'name', width: 200 },
  { title: '域名', dataIndex: 'domain', width: 180 },
  { title: '颁发者', dataIndex: 'issuer', width: 150 },
  { title: '过期日期', dataIndex: 'expireDate', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
];

const getStatusColor = (s: string) => ({ valid: 'green', expiring: 'orange', expired: 'red' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.ssl-certificates-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
