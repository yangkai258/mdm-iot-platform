<template>
  <div class="announcements-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="公告总数" :value="36" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="进行中" :value="5" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="总浏览" :value="25600" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>系统公告管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            发布公告
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="announcements" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.status === 'published' ? 'success' : 'default'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">查看</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑公告' : '发布公告'" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="公告标题" required>
          <a-input v-model="form.title" placeholder="请输入公告标题" />
        </a-form-item>
        <a-form-item label="公告类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="system">系统通知</a-option>
            <a-option value="maintenance">维护公告</a-option>
            <a-option value="activity">活动公告</a-option>
            <a-option value="update">版本更新</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="公告内容">
          <a-textarea v-model="form.content" :rows="6" />
        </a-form-item>
        <a-form-item label="发布时间">
          <a-date-picker v-model="form.publishAt" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="置顶">
          <a-switch v-model="form.pinned" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ title: '', type: '', content: '', publishAt: '', pinned: false });

const announcements = ref([
  { id: 1, title: '系统维护通知', type: 'maintenance', typeText: '维护', status: 'published', statusText: '已发布', views: 2560, publishAt: '2026-03-28 10:00:00' },
  { id: 2, title: '新功能上线公告', type: 'update', typeText: '更新', status: 'draft', statusText: '草稿', views: 0, publishAt: null },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '浏览量', dataIndex: 'views', width: 100 },
  { title: '发布时间', dataIndex: 'publishAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleView = (record: any) => {};
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.announcements-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
