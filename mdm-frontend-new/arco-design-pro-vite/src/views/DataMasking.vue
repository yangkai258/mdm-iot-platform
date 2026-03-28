<template>
  <div class="data-masking-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>数据脱敏规则</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增规则
          </a-button>
        </div>
      </template>
      
      <a-alert type="info" style="margin-bottom: 16px;">
        配置数据脱敏规则，对敏感数据进行自动脱敏处理，符合GDPR等合规要求
      </a-alert>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #fieldType="{ record }">
          <a-tag>{{ record.fieldTypeText }}</a-tag>
        </template>
        <template #maskingType="{ record }">
          <a-tag :color="getMaskingColor(record.maskingType)">{{ record.maskingTypeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-switch :checked="record.enabled" disabled />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑规则' : '新增规则'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="字段名称" required>
          <a-input v-model="form.fieldName" placeholder="如: phone" />
        </a-form-item>
        <a-form-item label="字段类型" required>
          <a-select v-model="form.fieldType" placeholder="选择字段类型">
            <a-option value="phone">手机号</a-option>
            <a-option value="email">邮箱</a-option>
            <a-option value="idcard">身份证</a-option>
            <a-option value="bankcard">银行卡</a-option>
            <a-option value="name">姓名</a-option>
            <a-option value="address">地址</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="脱敏方式" required>
          <a-select v-model="form.maskingType" placeholder="选择脱敏方式">
            <a-option value="partial">部分隐藏</a-option>
            <a-option value="hash">哈希</a-option>
            <a-option value="encrypt">加密</a-option>
            <a-option value="replace">替换</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="脱敏规则">
          <a-input v-model="form.rule" placeholder="如: 3-7位显示，其余*代替" />
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
  { id: 1, fieldName: 'phone', fieldType: 'phone', fieldTypeText: '手机号', maskingType: 'partial', maskingTypeText: '部分隐藏', rule: '138****1234', enabled: true },
  { id: 2, fieldName: 'email', fieldType: 'email', fieldTypeText: '邮箱', maskingType: 'partial', maskingTypeText: '部分隐藏', rule: 'a***@example.com', enabled: true },
  { id: 3, fieldName: 'idcard', fieldType: 'idcard', fieldTypeText: '身份证', maskingType: 'partial', maskingTypeText: '部分隐藏', rule: '110***********1234', enabled: true },
  { id: 4, fieldName: 'bankcard', fieldType: 'bankcard', fieldTypeText: '银行卡', maskingType: 'encrypt', maskingTypeText: '加密', rule: 'AES-256加密', enabled: false },
  { id: 5, fieldName: 'realname', fieldType: 'name', fieldTypeText: '姓名', maskingType: 'replace', maskingTypeText: '替换', rule: '***', enabled: true },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 5 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({
  fieldName: '', fieldType: '', maskingType: '', rule: '', enabled: true,
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '字段名称', dataIndex: 'fieldName', width: 120 },
  { title: '字段类型', slotName: 'fieldType', width: 100 },
  { title: '脱敏方式', slotName: 'maskingType', width: 100 },
  { title: '脱敏规则', dataIndex: 'rule', width: 200 },
  { title: '启用', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getMaskingColor = (t: string) => ({ partial: 'blue', hash: 'orange', encrypt: 'green', replace: 'purple' }[t] || 'default');

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.data-masking-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
