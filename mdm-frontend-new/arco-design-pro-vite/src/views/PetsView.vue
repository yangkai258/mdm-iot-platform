<template>
  <div class="pets-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增宠物
          </a-button>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-input v-model="searchForm.keyword" placeholder="搜索宠物名称/ID" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="状态" allow-clear>
              <a-option value="active">活跃</a-option>
              <a-option value="idle">空闲</a-option>
              <a-option value="sleeping">休眠</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
        <template #avatar="{ record }">
          <a-avatar :size="40" :style="{ backgroundColor: record.color }">
            {{ record.name?.charAt(0) }}
          </a-avatar>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ record.status }}
          </a-tag>
        </template>
        <template #mood="{ record }">
          <a-tag :color="getMoodColor(record.mood)">{{ record.mood }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleMemory(record)">记忆</a-link>
            <a-link @click="handleAction(record)">动作</a-link>
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
  { id: 'PET001', name: '小白', deviceId: 'DEV001', breed: '柯基', age: 2, personality: '活泼', mood: '开心', status: 'active', color: '#165DFF', createdAt: '2026-01-15' },
  { id: 'PET002', name: '小黄', deviceId: 'DEV002', breed: '金毛', age: 3, personality: '温顺', mood: '平静', status: 'idle', color: '#FF7D00', createdAt: '2026-02-01' },
  { id: 'PET003', name: '小红', deviceId: 'DEV003', breed: '柴犬', age: 1, personality: '好奇', mood: '兴奋', status: 'active', color: '#F53F3F', createdAt: '2026-02-20' },
  { id: 'PET004', name: '小绿', deviceId: 'DEV004', breed: '边牧', age: 4, personality: '聪明', mood: '困倦', status: 'sleeping', color: '#00B42A', createdAt: '2026-03-01' },
]);

const searchForm = reactive({ keyword: '', status: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '', slotName: 'avatar', width: 60 },
  { title: '宠物ID', dataIndex: 'id', width: 100 },
  { title: '名称', dataIndex: 'name', width: 100 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '品种', dataIndex: 'breed', width: 80 },
  { title: '年龄', dataIndex: 'age', width: 60 },
  { title: '性格', dataIndex: 'personality', width: 80 },
  { title: '情绪', slotName: 'mood', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const getStatusColor = (status: string) => {
  const map: Record<string, string> = { active: 'green', idle: 'gray', sleeping: 'blue' };
  return map[status] || 'default';
};

const getMoodColor = (mood: string) => {
  const map: Record<string, string> = { '开心': 'green', '平静': 'blue', '兴奋': 'orange', '困倦': 'gray' };
  return map[mood] || 'default';
};

const handleSearch = () => {};
const handlePageChange = (page: number) => { pagination.current = page; };
const handleCreate = () => {};
const handleView = (record: any) => {};
const handleMemory = (record: any) => {};
const handleAction = (record: any) => {};
</script>

<style scoped>
.pets-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
