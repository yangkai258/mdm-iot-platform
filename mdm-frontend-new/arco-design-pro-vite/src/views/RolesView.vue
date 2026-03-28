<template>
  <div class="roles-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>角色管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增角色
          </a-button>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索角色名称/编码" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="状态" allow-clear>
              <a-option value="active">启用</a-option>
              <a-option value="inactive">禁用</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'red'">
            {{ record.status === 'active' ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template #userCount="{ record }">
          <span style="color: #165DFF;">{{ record.userCount }} 人</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleAuth(record)">权限</a-link>
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
  { id: 'R001', code: 'admin', name: '超级管理员', description: '拥有系统所有权限', status: 'active', userCount: 3, createdAt: '2026-01-01 00:00:00' },
  { id: 'R002', code: 'operator', name: '运营管理员', description: '负责日常运营管理', status: 'active', userCount: 12, createdAt: '2026-01-15 10:00:00' },
  { id: 'R003', code: 'viewer', name: '只读用户', description: '仅可查看数据', status: 'active', userCount: 45, createdAt: '2026-02-01 09:00:00' },
  { id: 'R004', code: 'device_mgr', name: '设备管理员', description: '管理设备和固件', status: 'inactive', userCount: 5, createdAt: '2026-02-20 14:00:00' },
]);

const searchForm = reactive({ keyword: '', status: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '角色ID', dataIndex: 'id', width: 100 },
  { title: '角色编码', dataIndex: 'code', width: 120 },
  { title: '角色名称', dataIndex: 'name', width: 150 },
  { title: '描述', dataIndex: 'description' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '用户数', slotName: 'userCount', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleSearch = () => { loading.value = true; setTimeout(() => { loading.value = false; }, 500); };
const handlePageChange = (page: number) => { pagination.current = page; };
const handleCreate = () => {};
const handleView = (record: any) => {};
const handleEdit = (record: any) => {};
const handleAuth = (record: any) => {};
const handleDelete = (record: any) => {};
</script>

<style scoped>
.roles-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
