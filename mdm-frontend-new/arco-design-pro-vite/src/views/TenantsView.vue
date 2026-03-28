<template>
  <div class="tenants-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>租户管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增租户
          </a-button>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索租户名称/编码" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="状态" allow-clear>
              <a-option value="active">正常</a-option>
              <a-option value="suspended">已暂停</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'orange'">
            {{ record.status === 'active' ? '正常' : '已暂停' }}
          </a-tag>
        </template>
        <template #expires="{ record }">
          <span :style="{ color: isExpiringSoon(record.expiresAt) ? '#FF7D00' : 'inherit' }">
            {{ record.expiresAt || '永久' }}
          </span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleConfig(record)">配置</a-link>
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
  { id: 'T001', code: 'corp_a', name: 'XX科技有限公司', contact: '张经理', phone: '138****1234', deviceLimit: 100, deviceCount: 45, status: 'active', expiresAt: '2027-01-01', createdAt: '2026-01-01' },
  { id: 'T002', code: 'corp_b', name: 'YY宠物医院', contact: '李医生', phone: '139****5678', deviceLimit: 50, deviceCount: 32, status: 'active', expiresAt: '2026-06-01', createdAt: '2026-02-15' },
  { id: 'T003', code: 'corp_c', name: 'ZZ智能家居', contact: '王总', phone: '137****9012', deviceLimit: 200, deviceCount: 0, status: 'suspended', expiresAt: '', createdAt: '2026-03-01' },
]);

const searchForm = reactive({ keyword: '', status: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 3 });

const columns = [
  { title: '租户ID', dataIndex: 'id', width: 80 },
  { title: '租户编码', dataIndex: 'code', width: 100 },
  { title: '租户名称', dataIndex: 'name', width: 180 },
  { title: '联系人', dataIndex: 'contact', width: 100 },
  { title: '联系电话', dataIndex: 'phone', width: 120 },
  { title: '设备数/限额', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '到期时间', slotName: 'expires', width: 120 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const isExpiringSoon = (date: string) => {
  if (!date) return false;
  const diff = new Date(date).getTime() - Date.now();
  return diff > 0 && diff < 30 * 24 * 60 * 60 * 1000;
};

const handleSearch = () => {};
const handleCreate = () => {};
const handleView = (record: any) => {};
const handleConfig = (record: any) => {};
const handleDelete = (record: any) => {};
const handlePageChange = (page: number) => { pagination.current = page; };
</script>

<style scoped>
.tenants-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
