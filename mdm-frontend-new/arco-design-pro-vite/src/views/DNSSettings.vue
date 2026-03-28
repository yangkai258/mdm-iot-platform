<template>
  <div class="dns-settings-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="DNS记录" :value="128" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="域名" :value="25" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="查询次数" :value="125600" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>DNS设置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加记录
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="records" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" title="编辑DNS记录" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="记录类型">
          <a-select v-model="form.type">
            <a-option value="A">A记录</a-option>
            <a-option value="CNAME">CNAME</a-option>
            <a-option value="MX">MX</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="主机记录">
          <a-input v-model="form.host" placeholder="@或子域名" />
        </a-form-item>
        <a-form-item label="记录值">
          <a-input v-model="form.value" placeholder="IP地址或域名" />
        </a-form-item>
        <a-form-item label="TTL">
          <a-input-number v-model="form.ttl" :min="60" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const editVisible = ref(false);

const form = reactive({ type: '', host: '', value: '', ttl: 600 });

const records = ref([
  { id: 1, domain: 'example.com', type: 'A', typeText: 'A', host: '@', value: '192.168.1.100', ttl: 600 },
  { id: 2, domain: 'example.com', type: 'CNAME', typeText: 'CNAME', host: 'www', value: 'example.com', ttl: 600 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '域名', dataIndex: 'domain', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '主机', dataIndex: 'host', width: 120 },
  { title: '记录值', dataIndex: 'value', ellipsis: true },
  { title: 'TTL', dataIndex: 'ttl', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { editVisible.value = true; };
const handleEdit = (record: any) => { Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.dns-settings-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
