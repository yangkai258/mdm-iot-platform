<template>
  <div class="pet-microsite-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="微站总数" :value="25" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="访问量" :value="125800" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="绑定宠物" :value="568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="转化率" :value="12.5" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物微站管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建微站
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="sites" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-badge :status="record.status === 'published' ? 'success' : 'default'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handlePreview(record)">预览</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleShare(record)">分享</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建微站" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="微站名称" required>
          <a-input v-model="form.name" placeholder="请输入微站名称" />
        </a-form-item>
        <a-form-item label="关联宠物">
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="主题模板">
          <a-select v-model="form.template">
            <a-option value="default">默认模板</a-option>
            <a-option value="cute">可爱风格</a-option>
            <a-option value="minimal">简约风格</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="自定义域名">
          <a-input v-model="form.domain" placeholder="pet.example.com" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const createVisible = ref(false);

const form = reactive({ name: '', petId: '', template: 'default', domain: '' });

const sites = ref([
  { id: 1, name: '小黄的个人主页', petName: '小黄', template: 'cute', visitors: 25600, status: 'published', statusText: '已发布', url: 'https://pet.example.com/xiaohuang' },
  { id: 2, name: '小红的故事站', petName: '小红', template: 'minimal', visitors: 12500, status: 'draft', statusText: '草稿', url: '' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '微站名称', dataIndex: 'name', width: 200 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '模板', dataIndex: 'template', width: 100 },
  { title: '访问量', dataIndex: 'visitors', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: 'URL', dataIndex: 'url', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const handleCreate = () => { createVisible.value = true; };
const handlePreview = (record: any) => {};
const handleEdit = (record: any) => {};
const handleShare = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.pet-microsite-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
