<template>
  <div class="error-tracking-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="错误总数" :value="568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="未解决" :value="28" status="error" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日新增" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="解决率" :value="95" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>错误追踪</span>
          <a-space>
            <a-input-search placeholder="搜索错误" style="width: 200px;" />
            <a-button @click="handleRefresh">刷新</a-button>
          </a-space>
        </div>
      </template>
      
      <a-table :columns="columns" :data="errors" :loading="loading" :pagination="pagination">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">{{ record.levelText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'resolved' ? 'green' : 'orange'">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 10 });

const errors = ref([
  { id: 1, errorType: 'NullPointerException', level: 'error', levelText: '错误', message: '空指针异常', module: 'DeviceService', count: 125, status: 'resolved', statusText: '已解决', firstOccurrence: '2026-03-28 10:00:00' },
  { id: 2, errorType: 'TimeoutException', level: 'warning', levelText: '警告', message: '请求超时', module: 'APIController', count: 56, status: 'pending', statusText: '待处理', firstOccurrence: '2026-03-28 18:00:00' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '错误类型', dataIndex: 'errorType', width: 180 },
  { title: '级别', slotName: 'level', width: 80 },
  { title: '消息', dataIndex: 'message', ellipsis: true },
  { title: '模块', dataIndex: 'module', width: 120 },
  { title: '次数', dataIndex: 'count', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '首次出现', dataIndex: 'firstOccurrence', width: 160 },
];

const getLevelColor = (l: string) => ({ error: 'red', warning: 'orange', info: 'blue' }[l] || 'default');

const handleRefresh = () => {};
</script>

<style scoped>
.error-tracking-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
