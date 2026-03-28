<template>
  <div class="cluster-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="集群数" :value="8" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="节点总数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="运行中" :value="125" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="CPU使用" :value="45" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>集群管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建集群
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="clusters" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-badge :status="record.status === 'healthy' ? 'success' : 'error'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleScale(record)">扩容</a-link>
            <a-link @click="handleConfig(record)">配置</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="configVisible" :title="currentCluster.name + ' 配置'" :width="600">
      <a-form layout="vertical">
        <a-form-item label="集群名称">
          <a-input v-model="currentCluster.name" />
        </a-form-item>
        <a-form-item label="节点数">
          <a-input-number v-model="currentCluster.nodeCount" :min="1" />
        </a-form-item>
        <a-form-item label="资源配置">
          <a-space direction="vertical" fill>
            <a-input-number v-model="currentCluster.cpu" :min="1" suffix="CPU" />
            <a-input-number v-model="currentCluster.memory" :min="1" suffix="GB内存" />
            <a-input-number v-model="currentCluster.storage" :min="1" suffix="GB存储" />
          </a-space>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSaveConfig">保存配置</a-button>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const configVisible = ref(false);
const currentCluster = ref<any>({});

const clusters = ref([
  { id: 1, name: '生产集群', nodes: 50, status: 'healthy', statusText: '健康', cpu: 45, memory: 60, storage: 40 },
  { id: 2, name: '测试集群', nodes: 10, status: 'healthy', statusText: '健康', cpu: 20, memory: 30, storage: 25 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '集群名称', dataIndex: 'name', width: 150 },
  { title: '节点数', dataIndex: 'nodes', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: 'CPU', dataIndex: 'cpu', width: 80 },
  { title: '内存', dataIndex: 'memory', width: 80 },
  { title: '存储', dataIndex: 'storage', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => {};
const handleScale = (record: any) => {};
const handleConfig = (record: any) => { currentCluster.value = { ...record }; configVisible.value = true; };
const handleSaveConfig = () => { configVisible.value = false; };
</script>

<style scoped>
.cluster-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
