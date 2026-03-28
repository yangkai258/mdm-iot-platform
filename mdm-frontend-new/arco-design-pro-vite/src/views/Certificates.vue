<template>
  <div class="certificates-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>证书管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            上传证书
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.valid ? 'green' : 'red'">{{ record.valid ? '有效' : '已过期' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleDownload(record)">下载</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="上传证书" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="证书名称" required>
          <a-input v-model="form.name" placeholder="请输入证书名称" />
        </a-form-item>
        <a-form-item label="证书类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="ssl">SSL证书</a-option>
            <a-option value="device">设备证书</a-option>
            <a-option value="code_signing">代码签名</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="证书文件">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="到期日期">
          <a-date-picker v-model="form.expireAt" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, name: 'api.example.com SSL', type: 'ssl', typeText: 'SSL证书', serial: '00:AB:CD:...', valid: true, expireAt: '2027-03-28', issuedAt: '2026-03-28' },
  { id: 2, name: '设备签名证书', type: 'device', typeText: '设备证书', serial: 'DE:AD:BE:...', valid: true, expireAt: '2028-03-28', issuedAt: '2026-03-28' },
  { id: 3, name: '代码签名证书', type: 'code_signing', typeText: '代码签名', serial: 'FA:CE:...', valid: false, expireAt: '2026-01-01', issuedAt: '2023-01-01' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const createVisible = ref(false);

const form = reactive({ name: '', type: 'ssl', expireAt: null, description: '' });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '证书名称', dataIndex: 'name', width: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '序列号', dataIndex: 'serial', width: 150 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '发放日期', dataIndex: 'issuedAt', width: 120 },
  { title: '到期日期', dataIndex: 'expireAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { createVisible.value = true; };
const handleDownload = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.certificates-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
