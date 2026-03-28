<template>
  <div class="device-grouping-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="设备分组" :value="25" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="分组设备" :value="580" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="平均分组" :value="23" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备分组管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建分组
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="groups" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">查看</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="detailVisible" :title="currentGroup.name" :width="600">
      <div v-if="currentGroup">
        <a-descriptions :column="1" bordered>
          <a-descriptions-item label="分组名称">{{ currentGroup.name }}</a-descriptions-item>
          <a-descriptions-item label="分组类型">{{ currentGroup.typeText }}</a-descriptions-item>
          <a-descriptions-item label="设备数量">{{ currentGroup.deviceCount }}</a-descriptions-item>
          <a-descriptions-item label="描述">{{ currentGroup.description }}</a-descriptions-item>
        </a-descriptions>
        
        <a-divider>分组设备列表</a-divider>
        <a-table :columns="deviceColumns" :data="currentDevices" :pagination="false" size="small" />
      </div>
    </a-drawer>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑分组' : '创建分组'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="分组名称" required>
          <a-input v-model="form.name" placeholder="请输入分组名称" />
        </a-form-item>
        <a-form-item label="分组类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="location">位置分组</a-option>
            <a-option value="function">功能分组</a-option>
            <a-option value="owner"> owner分组</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const detailVisible = ref(false);
const editVisible = ref(false);
const isEdit = ref(false);
const currentGroup = ref<any>(null);

const form = reactive({ name: '', type: '', description: '' });

const groups = ref([
  { id: 1, name: '客厅设备', type: 'location', typeText: '位置分组', deviceCount: 5, description: '客厅中的所有设备', createdAt: '2026-03-20' },
  { id: 2, name: '卧室设备', type: 'location', typeText: '位置分组', deviceCount: 3, description: '卧室中的所有设备', createdAt: '2026-03-20' },
  { id: 3, name: '健康监测', type: 'function', typeText: '功能分组', deviceCount: 8, description: '用于健康监测的设备', createdAt: '2026-03-15' },
]);

const currentDevices = ref([
  { deviceId: 'DEV001', deviceName: '小黄-客厅', status: 'online' },
  { deviceId: 'DEV002', deviceName: '小红-客厅', status: 'online' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '分组名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '设备数', dataIndex: 'deviceCount', width: 100 },
  { title: '描述', dataIndex: 'description' },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const deviceColumns = [
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'deviceName', width: 150 },
  { title: '状态', dataIndex: 'status', width: 100 },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleView = (record: any) => { currentGroup.value = record; detailVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.device-grouping-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
