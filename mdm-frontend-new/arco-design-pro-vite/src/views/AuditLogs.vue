<template>
  <div class="audit-logs-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="操作日志" :value="12580" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="今日操作" :value="256" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="异常操作" :value="12" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>操作审计日志</span>
          <a-space>
            <a-input-search placeholder="搜索用户/操作" style="width: 200px;" />
            <a-button @click="handleExport">导出</a-button>
          </a-space>
        </div>
      </template>
      
      <a-table :columns="columns" :data="logs" :loading="loading" :pagination="pagination">
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)">{{ record.levelText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.success ? 'green' : 'red'">{{ record.success ? '成功' : '失败' }}</a-tag>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 20 });

const logs = ref([
  { id: 1, userName: 'admin', action: '登录系统', module: '认证', level: 'info', levelText: '信息', success: true, ip: '192.168.1.100', time: '2026-03-28 18:00:00' },
  { id: 2, userName: 'admin', action: '删除设备', module: '设备管理', level: 'warning', levelText: '警告', success: true, ip: '192.168.1.100', time: '2026-03-28 17:00:00' },
  { id: 3, userName: 'user01', action: '修改密码', module: '用户管理', level: 'warning', levelText: '警告', success: true, ip: '192.168.1.101', time: '2026-03-28 16:00:00' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户', dataIndex: 'userName', width: 120 },
  { title: '操作', dataIndex: 'action', width: 150 },
  { title: '模块', dataIndex: 'module', width: 100 },
  { title: '级别', slotName: 'level', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: 'IP地址', dataIndex: 'ip', width: 150 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const getLevelColor = (l: string) => ({ info: 'blue', warning: 'orange', error: 'red' }[l] || 'default');

const handleExport = () => {};
</script>

<style scoped>
.audit-logs-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
