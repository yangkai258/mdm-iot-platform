<template>
  <div class="container-registry-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="镜像数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="镜像版本" :value="568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="拉取次数" :value="25600" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="存储大小" :value="256" suffix="GB" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>容器镜像仓库</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            上传镜像
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="images" :loading="loading" :pagination="pagination">
        <template #size="{ record }">
          <span>{{ record.size }}MB</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handlePull(record)">拉取</a-link>
            <a-link @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const images = ref([
  { id: 1, name: 'mdm-backend', version: 'v2.0.0', size: 256, pulls: 12560, lastPull: '2026-03-28 18:00:00' },
  { id: 2, name: 'mdm-frontend', version: 'v2.0.0', size: 128, pulls: 8560, lastPull: '2026-03-28 17:00:00' },
]);

const columns = [
  { title: '镜像名称', dataIndex: 'name', width: 200 },
  { title: '版本', dataIndex: 'version', width: 120 },
  { title: '大小', slotName: 'size', width: 100 },
  { title: '拉取次数', dataIndex: 'pulls', width: 120 },
  { title: '最后拉取', dataIndex: 'lastPull', width: 160 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => {};
const handlePull = (record: any) => {};
const handleDelete = (record: any) => {};
</script>

<style scoped>
.container-registry-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
