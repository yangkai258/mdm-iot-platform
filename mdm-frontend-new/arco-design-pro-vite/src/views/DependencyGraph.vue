<template>
  <div class="dependency-graph-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="服务数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="依赖关系" :value="256" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="循环依赖" :value="0" status="success" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="覆盖率" :value="95" suffix="%" /></a-card>
      </a-col>
    </a-row>
    <a-card>
      <template #title>服务依赖图</template>
      <a-tabs>
        <a-tab-pane key="graph" title="依赖图">
          <a-empty description="服务依赖图展示区域" />
        </a-tab-pane>
        <a-tab-pane key="list" title="依赖列表">
          <a-table :columns="columns" :data="dependencies" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const dependencies = ref([
  { id: 1, service: 'api-gateway', dependsOn: 'user-service,auth-service', type: 'internal' },
  { id: 2, service: 'user-service', dependsOn: 'database,cache', type: 'internal' },
]);
const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '服务', dataIndex: 'service', width: 150 },
  { title: '依赖服务', dataIndex: 'dependsOn', width: 250 },
  { title: '类型', dataIndex: 'type', width: 100 },
];
</script>

<style scoped>
.dependency-graph-container { padding: 20px; }
</style>
