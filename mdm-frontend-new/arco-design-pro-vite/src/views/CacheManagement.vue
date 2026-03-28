<template>
  <div class="cache-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="缓存键" :value="2568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="命中率" :value="95" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="内存使用" :value="2.5" suffix="GB" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="过期键" :value="128" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>缓存管理</span>
          <a-space>
            <a-button @click="handleFlush">刷新缓存</a-button>
            <a-button type="primary" @click="handleAddRule">添加规则</a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="keys" title="缓存键">
          <a-table :columns="keyColumns" :data="keys" :pagination="pagination">
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleViewValue(record)">查看</a-link>
                <a-link @click="handleDeleteKey(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="rules" title="缓存规则">
          <a-table :columns="ruleColumns" :data="rules" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 10 });

const keys = ref([
  { id: 1, key: 'device:status:DEV001', type: 'String', ttl: 90, hits: 1256, size: '256B' },
  { id: 2, key: 'user:profile:1001', type: 'Hash', ttl: 3600, hits: 856, size: '1KB' },
]);

const rules = ref([
  { id: 1, pattern: 'device:*', ttl: 90, maxSize: '10MB', enabled: true },
  { id: 2, pattern: 'user:*', ttl: 3600, maxSize: '100MB', enabled: true },
]);

const keyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '键', dataIndex: 'key', width: 250, ellipsis: true },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: 'TTL(秒)', dataIndex: 'ttl', width: 100 },
  { title: '命中次数', dataIndex: 'hits', width: 100 },
  { title: '大小', dataIndex: 'size', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const ruleColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '匹配模式', dataIndex: 'pattern', width: 150 },
  { title: 'TTL', dataIndex: 'ttl', width: 100 },
  { title: '最大大小', dataIndex: 'maxSize', width: 100 },
  { title: '启用', dataIndex: 'enabled', width: 80 },
];

const handleFlush = () => {};
const handleAddRule = () => {};
const handleViewValue = (record: any) => {};
const handleDeleteKey = (record: any) => {};
</script>

<style scoped>
.cache-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
