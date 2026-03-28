<template>
  <div class="permissions-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>权限管理</span>
          <a-space>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              新增权限
            </a-button>
            <a-button @click="handleExport">
              <template #icon><icon-download /></template>
              导出
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
        <template #type="{ record }">
          <a-tag :color="record.type === 'menu' ? 'blue' : record.type === 'button' ? 'green' : 'orange'">
            {{ record.type === 'menu' ? '菜单' : record.type === 'button' ? '按钮' : 'API' }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-switch :checked="record.status === 'active'" disabled />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'P001', code: 'member:view', name: '查看会员', type: 'button', path: '/api/v1/members', status: 'active', description: '查看会员列表和详情' },
  { id: 'P002', code: 'member:edit', name: '编辑会员', type: 'button', path: '/api/v1/members/*', status: 'active', description: '创建和编辑会员信息' },
  { id: 'P003', code: 'device:view', name: '查看设备', type: 'button', path: '/api/v1/devices', status: 'active', description: '查看设备列表和状态' },
  { id: 'P004', code: 'device:command', name: '设备指令', type: 'button', path: '/api/v1/devices/*/command', status: 'active', description: '发送设备控制指令' },
  { id: 'P005', code: 'role:manage', name: '角色管理', type: 'menu', path: '/dashboard/roles', status: 'active', description: '管理系统角色' },
  { id: 'P006', code: 'ota:upgrade', name: 'OTA升级', type: 'button', path: '/api/v1/ota/*', status: 'active', description: '执行OTA固件升级' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 6 });

const columns = [
  { title: '权限ID', dataIndex: 'id', width: 80 },
  { title: '权限编码', dataIndex: 'code', width: 150 },
  { title: '权限名称', dataIndex: 'name', width: 120 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '路径', dataIndex: 'path' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const handleCreate = () => {};
const handleExport = () => {};
const handleEdit = (record: any) => {};
const handleDelete = (record: any) => {};
const handlePageChange = (page: number) => { pagination.current = page; };
</script>

<style scoped>
.permissions-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
