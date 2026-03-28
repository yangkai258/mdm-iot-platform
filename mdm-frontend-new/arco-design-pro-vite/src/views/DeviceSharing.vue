<template>
  <div class="device-sharing-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="共享设备" :value="156" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="共享成员" :value="428" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="共享次数" :value="856" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备共享管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建共享
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="shares" :loading="loading" :pagination="pagination">
        <template #permission="{ record }">
          <a-tag>{{ record.permissionText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleRevoke(record)">取消共享</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="设备共享" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="设备" required>
          <a-select v-model="form.deviceId" placeholder="选择设备">
            <a-option value="DEV001">小黄-客厅</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="共享给" required>
          <a-select v-model="form.sharedTo" placeholder="选择用户">
            <a-option value="U001">用户A</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="权限">
          <a-select v-model="form.permission" placeholder="选择权限">
            <a-option value="view">查看</a-option>
            <a-option value="control">控制</a-option>
            <a-option value="admin">管理</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="有效期">
          <a-date-picker v-model="form.expireAt" style="width: 100%;" />
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

const form = reactive({ deviceId: '', sharedTo: '', permission: 'view', expireAt: '' });

const shares = ref([
  { id: 1, deviceName: '小黄-客厅', sharedTo: '用户A', permission: 'control', permissionText: '控制', createdAt: '2026-03-20', expireAt: '2026-04-20' },
  { id: 2, deviceName: '小红-卧室', sharedTo: '用户B', permission: 'view', permissionText: '查看', createdAt: '2026-03-15', expireAt: null },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '设备', dataIndex: 'deviceName', width: 150 },
  { title: '共享给', dataIndex: 'sharedTo', width: 120 },
  { title: '权限', slotName: 'permission', width: 100 },
  { title: '共享时间', dataIndex: 'createdAt', width: 120 },
  { title: '过期时间', dataIndex: 'expireAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { createVisible.value = true; };
const handleRevoke = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.device-sharing-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
