<template>
  <div class="log-analysis-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="日志总量" :value="1258000" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日日志" :value="56800" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="错误日志" :value="128" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="存储大小" :value="25" suffix="GB" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>日志分析</span>
          <a-space>
            <a-input-search placeholder="搜索日志内容" style="width: 200px;" />
            <a-button @click="handleExport">导出</a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="logs" title="日志列表">
          <a-table :columns="columns" :data="logs" :loading="loading" :pagination="pagination">
            <template #level="{ record }">
              <a-tag :color="getLevelColor(record.level)">{{ record.levelText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="statistics" title="统计分析">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="日志级别分布">
                <a-progress :percent="60" label="Info" />
                <a-progress :percent="25" label="Warning" />
                <a-progress :percent="10" label="Error" />
                <a-progress :percent="5" label="Debug" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="模块分布">
                <a-progress :percent="40" label="API模块" />
                <a-progress :percent="30" label="设备模块" />
                <a-progress :percent="20" label="会员模块" />
                <a-progress :percent="10" label="其他" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 20 });

const logs = ref([
  { id: 1, timestamp: '2026-03-28 18:00:00', level: 'info', levelText: 'Info', module: 'API', message: '用户登录成功', ip: '192.168.1.100' },
  { id: 2, timestamp: '2026-03-28 17:00:00', level: 'warning', levelText: 'Warning', module: 'Device', message: '设备离线', ip: '-' },
  { id: 3, timestamp: '2026-03-28 16:00:00', level: 'error', levelText: 'Error', module: 'API', message: 'API调用超时', ip: '192.168.1.101' },
]);

const columns = [
  { title: '时间', dataIndex: 'timestamp', width: 180 },
  { title: '级别', slotName: 'level', width: 100 },
  { title: '模块', dataIndex: 'module', width: 120 },
  { title: '消息', dataIndex: 'message', ellipsis: true },
  { title: 'IP地址', dataIndex: 'ip', width: 150 },
];

const getLevelColor = (l: string) => ({ info: 'blue', warning: 'orange', error: 'red', debug: 'gray' }[l] || 'default');

const handleExport = () => {};
</script>

<style scoped>
.log-analysis-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
