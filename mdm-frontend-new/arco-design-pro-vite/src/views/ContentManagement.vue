<template>
  <div class="content-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="内容总数" :value="1280" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日新增" :value="28" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="浏览量" :value="125600" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="待审核" :value="12" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>内容管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建内容
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="list" title="内容列表">
          <a-table :columns="columns" :data="contents" :loading="loading" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleEdit(record)">编辑</a-link>
                <a-link @click="handleReview(record)">审核</a-link>
                <a-link status="danger" @click="handleDelete(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="categories" title="分类管理">
          <a-table :columns="catColumns" :data="categories" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑内容' : '创建内容'" :width="800" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题" required>
          <a-input v-model="form.title" placeholder="请输入标题" />
        </a-form-item>
        <a-form-item label="分类">
          <a-select v-model="form.categoryId" placeholder="选择分类">
            <a-option value="C001">新闻资讯</a-option>
            <a-option value="C002">养宠知识</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="内容类型">
          <a-radio-group v-model="form.type">
            <a-radio value="article">文章</a-radio>
            <a-radio value="video">视频</a-radio>
            <a-radio value="image">图片</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea v-model="form.content" :rows="10" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 10 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ title: '', categoryId: '', type: 'article', content: '' });

const contents = ref([
  { id: 1, title: '如何正确训练宠物', type: 'article', typeText: '文章', status: 'published', statusText: '已发布', views: 2560, createdAt: '2026-03-28' },
  { id: 2, title: '宠物喂食指南', type: 'video', typeText: '视频', status: 'pending', statusText: '待审核', views: 0, createdAt: '2026-03-28' },
]);

const categories = ref([
  { id: 1, name: '新闻资讯', count: 256, order: 1 },
  { id: 2, name: '养宠知识', count: 512, order: 2 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '浏览量', dataIndex: 'views', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const catColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '分类名称', dataIndex: 'name', width: 150 },
  { title: '内容量', dataIndex: 'count', width: 100 },
  { title: '排序', dataIndex: 'order', width: 80 },
];

const getStatusColor = (s: string) => ({ published: 'green', pending: 'orange', draft: 'gray', rejected: 'red' }[s] || 'default');

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleReview = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.content-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
