<template>
  <div class="white-label-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="定制品牌" :value="12" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="活跃品牌" :value="8" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="品牌用户" :value="5680" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>品牌定制管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建品牌
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="brands" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'default'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handlePreview(record)">预览</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:visible="configVisible" :title="currentBrand.name + ' 配置'" :width="700">
      <a-form layout="vertical">
        <a-form-item label="品牌名称">
          <a-input v-model="currentBrand.name" />
        </a-form-item>
        <a-form-item label="Logo">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="主题色">
          <a-color-picker v-model="currentBrand.primaryColor" />
        </a-form-item>
        <a-form-item label="域名">
          <a-input v-model="currentBrand.domain" />
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
const currentBrand = ref<any>({});

const brands = ref([
  { id: 1, name: '品牌A', logo: '', primaryColor: '#1659d5', domain: 'brand-a.example.com', status: 'active', statusText: '活跃', userCount: 2560 },
  { id: 2, name: '品牌B', logo: '', primaryColor: '#00b42a', domain: 'brand-b.example.com', status: 'active', statusText: '活跃', userCount: 1280 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '品牌名称', dataIndex: 'name', width: 150 },
  { title: '域名', dataIndex: 'domain', ellipsis: true },
  { title: '用户数', dataIndex: 'userCount', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => {};
const handlePreview = (record: any) => {};
const handleEdit = (record: any) => { currentBrand.value = { ...record }; configVisible.value = true; };
const handleSaveConfig = () => { configVisible.value = false; };
</script>

<style scoped>
.white-label-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
