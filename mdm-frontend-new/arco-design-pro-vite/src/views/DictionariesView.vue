<template>
  <div class="dicts-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>字典管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增字典
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">
            {{ record.status === 'active' ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleItems(record)">字典项</a-link>
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
  { id: 'D001', code: 'pet_type', name: '宠物类型', items: 5, description: '宠物类型字典', status: 'active', createdAt: '2026-03-01' },
  { id: 'D002', code: 'device_status', name: '设备状态', items: 4, description: '设备在线状态', status: 'active', createdAt: '2026-03-01' },
  { id: 'D003', code: 'alert_level', name: '告警级别', items: 3, description: '告警严重程度', status: 'active', createdAt: '2026-03-01' },
  { id: 'D004', code: 'member_level', name: '会员等级', items: 4, description: '会员等级定义', status: 'active', createdAt: '2026-03-01' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '字典ID', dataIndex: 'id', width: 80 },
  { title: '字典编码', dataIndex: 'code', width: 150 },
  { title: '字典名称', dataIndex: 'name', width: 150 },
  { title: '字典项数量', dataIndex: 'items', width: 100 },
  { title: '描述', dataIndex: 'description' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleCreate = () => {};
const handleItems = (record: any) => {};
const handleEdit = (record: any) => {};
const handleDelete = (record: any) => {};
</script>

<style scoped>
.dicts-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
