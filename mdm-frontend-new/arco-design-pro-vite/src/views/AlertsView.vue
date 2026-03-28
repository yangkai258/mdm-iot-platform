<template>
  <div class="alerts-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>告警管理</span>
          <a-space>
            <a-button @click="handleExport">
              <template #icon><icon-download /></template>
              导出
            </a-button>
          </a-space>
        </div>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="4">
            <a-select v-model="searchForm.level" placeholder="告警级别" allow-clear>
              <a-option value="critical">Critical</a-option>
              <a-option value="warning">Warning</a-option>
              <a-option value="info">Info</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.status" placeholder="处理状态" allow-clear>
              <a-option value="pending">待处理</a-option>
              <a-option value="resolved">已处理</a-option>
            </a-select>
          </a-col>
          <a-col :span="6">
            <a-date-picker v-model="searchForm.startDate" placeholder="开始日期" style="width: 100%" />
          </a-col>
          <a-col :span="6">
            <a-date-picker v-model="searchForm.endDate" placeholder="结束日期" style="width: 100%" />
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination" @page-change="handlePageChange">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">
            {{ getLevelText(record.level) }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'resolved' ? 'green' : 'red'">
            {{ record.status === 'resolved' ? '已处理' : '待处理' }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link v-if="record.status !== 'resolved'" @click="handleResolve(record)">处理</a-link>
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
  { id: 'A001', deviceId: 'DEV001', level: 'critical', type: '设备离线', message: '设备超过30分钟无响应', status: 'pending', createdAt: '2026-03-28 09:00:00' },
  { id: 'A002', deviceId: 'DEV002', level: 'warning', type: '电量低', message: '设备电量低于20%', status: 'resolved', createdAt: '2026-03-28 08:30:00' },
  { id: 'A003', deviceId: 'DEV003', level: 'critical', type: '异常行为', message: '检测到可疑行为模式', status: 'pending', createdAt: '2026-03-27 22:00:00' },
  { id: 'A004', deviceId: 'DEV001', level: 'info', type: '固件更新', message: '新固件版本可用', status: 'resolved', createdAt: '2026-03-27 10:00:00' },
]);

const searchForm = reactive({ level: '', status: '', startDate: '', endDate: '' });
const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const columns = [
  { title: '告警ID', dataIndex: 'id', width: 100 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '级别', slotName: 'level', width: 100 },
  { title: '类型', dataIndex: 'type', width: 120 },
  { title: '消息', dataIndex: 'message' },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getLevelColor = (level: string) => {
  const map: Record<string, string> = { critical: 'red', warning: 'orange', info: 'blue' };
  return map[level] || 'default';
};

const getLevelText = (level: string) => {
  const map: Record<string, string> = { critical: 'Critical', warning: 'Warning', info: 'Info' };
  return map[level] || level;
};

const handleSearch = () => {};
const handleExport = () => {};
const handleView = (record: any) => {};
const handleResolve = (record: any) => {};
const handlePageChange = (page: number) => { pagination.current = page; };
</script>

<style scoped>
.alerts-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
