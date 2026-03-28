<template>
  <div class="compliance-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="合规策略" :value="8" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="合规评分" :value="95" suffix="/100" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="待处理" :value="2" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="最近审计" value="2026-03-28" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>合规策略配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增策略
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑策略' : '新增策略'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="策略名称" required>
          <a-input v-model="form.name" placeholder="请输入策略名称" />
        </a-form-item>
        <a-form-item label="合规类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="gdpr">GDPR合规</a-option>
            <a-option value="ccpa">CCPA合规</a-option>
            <a-option value="security">安全合规</a-option>
            <a-option value="privacy">隐私合规</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="策略规则">
          <a-textarea v-model="form.rule" placeholder="请输入策略规则" :rows="4" />
        </a-form-item>
        <a-form-item label="执行动作">
          <a-select v-model="form.action" placeholder="选择动作">
            <a-option value="block">阻止操作</a-option>
            <a-option value="alert">发送告警</a-option>
            <a-option value="mask">数据脱敏</a-option>
            <a-option value="audit">记录审计</a-option>
          </a-select>
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
const data = ref([
  { id: 1, name: 'GDPR数据删除权', type: 'gdpr', typeText: 'GDPR合规', rule: '用户可请求删除个人数据', action: 'block', enabled: true },
  { id: 2, name: '数据访问审计', type: 'security', typeText: '安全合规', rule: '所有数据访问必须记录审计日志', action: 'audit', enabled: true },
  { id: 3, name: '敏感数据脱敏', type: 'privacy', typeText: '隐私合规', rule: '身份证、银行卡等敏感信息必须脱敏', action: 'mask', enabled: true },
  { id: 4, name: 'CCPA拒绝销售', type: 'ccpa', typeText: 'CCPA合规', rule: '用户可拒绝数据销售', action: 'block', enabled: false },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: '', rule: '', action: '', enabled: true });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '策略名称', dataIndex: 'name', width: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '策略规则', dataIndex: 'rule' },
  { title: '执行动作', dataIndex: 'action', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.compliance-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
