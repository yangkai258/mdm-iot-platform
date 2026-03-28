<template>
  <div class="logs-container">
    <a-card>
      <template #title>
        <span>操作日志</span>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="4">
            <a-select v-model="searchForm.type" placeholder="操作类型" allow-clear>
              <a-option value="create">创建</a-option>
              <a-option value="update">更新</a-option>
              <a-option value="delete">删除</a-option>
              <a-option value="login">登录</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-input v-model="searchForm.operator" placeholder="操作人" allow-clear />
          </a-col>
          <a-col :span="6">
            <a-range-picker v-model="searchForm.dateRange" />
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.type }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'success' ? 'green' : 'red'">{{ record.status }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'L001', type: 'login', operator: 'admin', target: '-', ip: '192.168.1.100', status: 'success', detail: '用户登录成功', createdAt: '2026-03-28 10:30:00' },
  { id: 'L002', type: 'create', operator: 'admin', target: '设备 DEV001', ip: '192.168.1.100', status: 'success', detail: '创建设备成功', createdAt: '2026-03-28 10:25:00' },
  { id: 'L003', type: 'update', operator: 'admin', target: '会员 M001', ip: '192.168.1.100', status: 'success', detail: '更新会员信息', createdAt: '2026-03-28 10:20:00' },
  { id: 'L004', type: 'delete', operator: 'admin', target: '优惠券 C003', ip: '192.168.1.100', status: 'success', detail: '删除优惠券', createdAt: '2026-03-28 10:15:00' },
]);

const searchForm = reactive({ type: '', operator: '', dateRange: [] });
const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '日志ID', dataIndex: 'id', width: 80 },
  { title: '操作类型', slotName: 'type', width: 100 },
  { title: '操作人', dataIndex: 'operator', width: 100 },
  { title: '操作对象', dataIndex: 'target', width: 150 },
  { title: 'IP地址', dataIndex: 'ip', width: 130 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '详情', dataIndex: 'detail' },
  { title: '时间', dataIndex: 'createdAt', width: 160 },
];

const getTypeColor = (type: string) => {
  const map: Record<string, string> = { login: 'blue', create: 'green', update: 'orange', delete: 'red' };
  return map[type] || 'default';
};

const handleSearch = () => {};
</script>

<style scoped>
.logs-container { padding: 20px; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
