<template>
  <div class="api-webhooks-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="Webhook总数" :value="25" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="触发次数" :value="5680" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="成功率" :value="98.5" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>Webhook管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建Webhook
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="webhooks" :loading="loading" :pagination="pagination">
        <template #event="{ record }">
          <a-tag>{{ record.eventText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.enabled ? 'success' : 'default'" :text="record.enabled ? '启用' : '禁用'" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleTest(record)">测试</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑Webhook' : '创建Webhook'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="Webhook名称" required>
          <a-input v-model="form.name" placeholder="请输入名称" />
        </a-form-item>
        <a-form-item label="触发事件">
          <a-select v-model="form.event" placeholder="选择事件">
            <a-option value="device.online">设备上线</a-option>
            <a-option value="device.offline">设备离线</a-option>
            <a-option value="alert.created">告警创建</a-option>
            <a-option value="member.registered">会员注册</a-option>
            <a-option value="order.created">订单创建</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="回调URL" required>
          <a-input v-model="form.url" placeholder="https://example.com/webhook" />
        </a-form-item>
        <a-form-item label="加密密钥">
          <a-input-password v-model="form.secret" placeholder="可选，用于签名验证" />
        </a-form-item>
        <a-form-item label="重试次数">
          <a-input-number v-model="form.retryCount" :min="0" :max="5" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model="form.enabled" />
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

const form = reactive({ name: '', event: '', url: '', secret: '', retryCount: 3, enabled: true });

const webhooks = ref([
  { id: 1, name: '设备状态同步', event: 'device.status', eventText: '设备状态', url: 'https://api.example.com/device', enabled: true, lastTrigger: '2026-03-28 18:00:00' },
  { id: 2, name: '会员注册通知', event: 'member.registered', eventText: '会员注册', url: 'https://api.example.com/member', enabled: true, lastTrigger: '2026-03-28 10:00:00' },
  { id: 3, name: '订单创建回调', event: 'order.created', eventText: '订单创建', url: 'https://api.example.com/order', enabled: false, lastTrigger: null },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: '事件', slotName: 'event', width: 120 },
  { title: '回调URL', dataIndex: 'url', ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '最后触发', dataIndex: 'lastTrigger', width: 160 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleTest = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.api-webhooks-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
