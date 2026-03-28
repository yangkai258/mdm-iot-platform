<template>
  <div class="depts-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>组织架构</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增部门
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
            {{ record.status === 'active' ? '正常' : '停用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleMembers(record)">成员</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
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
  { id: 'DEPT001', name: '技术研发部', parentId: '', manager: '张工', members: 25, devices: 120, status: 'active', createdAt: '2026-01-01' },
  { id: 'DEPT002', name: '产品运营部', parentId: '', manager: '李经理', members: 15, devices: 45, status: 'active', createdAt: '2026-01-01' },
  { id: 'DEPT003', name: '前端组', parentId: 'DEPT001', manager: '王工', members: 8, devices: 30, status: 'active', createdAt: '2026-01-15' },
  { id: 'DEPT004', name: '后端组', parentId: 'DEPT001', manager: '刘工', members: 10, devices: 40, status: 'active', createdAt: '2026-01-15' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '部门ID', dataIndex: 'id', width: 100 },
  { title: '部门名称', dataIndex: 'name', width: 150 },
  { title: '上级部门', dataIndex: 'parentId', width: 100 },
  { title: '负责人', dataIndex: 'manager', width: 100 },
  { title: '成员数', dataIndex: 'members', width: 80 },
  { title: '设备数', dataIndex: 'devices', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => {};
const handleMembers = (record: any) => {};
const handleEdit = (record: any) => {};
</script>

<style scoped>
.depts-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
