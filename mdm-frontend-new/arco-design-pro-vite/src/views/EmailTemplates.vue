<template>
  <div class="email-templates-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>邮件模板管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建模板
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="templates" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handlePreview(record)">预览</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑模板' : '创建模板'" :width="800" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="模板名称" required>
          <a-input v-model="form.name" placeholder="请输入模板名称" />
        </a-form-item>
        <a-form-item label="模板类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="welcome">欢迎邮件</a-option>
            <a-option value="notification">通知邮件</a-option>
            <a-option value="marketing">营销邮件</a-option>
            <a-option value="alert">告警邮件</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="邮件主题">
          <a-input v-model="form.subject" placeholder="请输入邮件主题" />
        </a-form-item>
        <a-form-item label="模板内容">
          <a-textarea v-model="form.content" :rows="10" placeholder="支持HTML和变量: {{name}}, {{time}}" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="previewVisible" title="预览" :width="700">
      <div v-html="previewContent" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const editVisible = ref(false);
const previewVisible = ref(false);
const isEdit = ref(false);
const previewContent = ref('');

const form = reactive({ name: '', type: '', subject: '', content: '' });

const templates = ref([
  { id: 1, name: '新用户欢迎', type: 'welcome', typeText: '欢迎邮件', subject: '欢迎加入MDM大家庭', usageCount: 1250, createdAt: '2026-03-20' },
  { id: 2, name: '设备绑定通知', type: 'notification', typeText: '通知邮件', subject: '您的设备已成功绑定', usageCount: 856, createdAt: '2026-03-15' },
  { id: 3, name: '会员日促销', type: 'marketing', typeText: '营销邮件', subject: '🎉会员日专属优惠', usageCount: 2500, createdAt: '2026-03-10' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '模板名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '邮件主题', dataIndex: 'subject', ellipsis: true },
  { title: '使用次数', dataIndex: 'usageCount', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handlePreview = (record: any) => { previewContent.value = record.content; previewVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.email-templates-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
