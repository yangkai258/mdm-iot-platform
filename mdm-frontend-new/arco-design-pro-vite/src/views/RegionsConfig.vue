<template>
  <div class="regions-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>区域配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增区域
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleNodes(record)">边缘节点</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 边缘节点弹窗 -->
    <a-modal v-model:visible="nodesVisible" title="边缘节点" :width="700">
      <a-table :columns="nodeColumns" :data="nodes" :pagination="false">
        <template #status="{ record }">
          <a-tag :color="record.online ? 'green' : 'red'">{{ record.online ? '在线' : '离线' }}</a-tag>
        </template>
      </a-table>
    </a-modal>

    <!-- 编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑区域' : '新增区域'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="区域名称" required>
          <a-input v-model="form.name" placeholder="如: 中国大陆" />
        </a-form-item>
        <a-form-item label="区域代码" required>
          <a-input v-model="form.code" placeholder="如: CN" />
        </a-form-item>
        <a-form-item label="时区">
          <a-select v-model="form.timezone" placeholder="选择时区">
            <a-option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</a-option>
            <a-option value="America/New_York">America/New_York (UTC-5)</a-option>
            <a-option value="Europe/London">Europe/London (UTC+0)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="数据库">
          <a-select v-model="form.dbConfig" placeholder="选择数据库">
            <a-option value="primary">主数据库</a-option>
            <a-option value="cn-east">华东节点</a-option>
            <a-option value="cn-north">华北节点</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, name: '中国大陆', code: 'CN', timezone: 'Asia/Shanghai', dbConfig: '主数据库', enabled: true, nodeCount: 5 },
  { id: 2, name: '北美', code: 'US', timezone: 'America/New_York', dbConfig: 'us-east', enabled: true, nodeCount: 3 },
  { id: 3, name: '欧洲', code: 'EU', timezone: 'Europe/London', dbConfig: 'eu-west', enabled: false, nodeCount: 0 },
]);

const nodes = ref([
  { id: 'N001', name: 'CN-华东-1', ip: '10.0.1.1', online: true, load: 45, lastSeen: '2026-03-28 18:00:00' },
  { id: 'N002', name: 'CN-华东-2', ip: '10.0.1.2', online: true, load: 38, lastSeen: '2026-03-28 18:00:00' },
  { id: 'N003', name: 'CN-华北-1', ip: '10.0.2.1', online: false, load: 0, lastSeen: '2026-03-27 10:00:00' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const editVisible = ref(false);
const nodesVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', code: '', timezone: 'Asia/Shanghai', dbConfig: 'primary', enabled: true });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '区域名称', dataIndex: 'name', width: 120 },
  { title: '代码', dataIndex: 'code', width: 80 },
  { title: '时区', dataIndex: 'timezone', width: 180 },
  { title: '数据库', dataIndex: 'dbConfig', width: 120 },
  { title: '节点数', dataIndex: 'nodeCount', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' },
];

const nodeColumns = [
  { title: '节点ID', dataIndex: 'id', width: 100 },
  { title: '节点名称', dataIndex: 'name', width: 120 },
  { title: 'IP地址', dataIndex: 'ip', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '负载', dataIndex: 'load', width: 80 },
  { title: '最后活跃', dataIndex: 'lastSeen', width: 160 },
];

const handleCreate = () => { isEdit.value = false; Object.assign(form, { name: '', code: '', timezone: 'Asia/Shanghai', dbConfig: 'primary', enabled: true }); editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleNodes = (record: any) => { nodesVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.regions-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
