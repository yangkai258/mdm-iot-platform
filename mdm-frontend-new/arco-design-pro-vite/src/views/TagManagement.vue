<template>
  <div class="tag-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="标签总数" :value="156" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="用户数" :value="8560" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="使用率" :value="92" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>用户标签管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建标签
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="tags" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="record.typeColor">{{ record.typeText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑标签' : '创建标签'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标签名称" required>
          <a-input v-model="form.name" placeholder="请输入标签名称" />
        </a-form-item>
        <a-form-item label="标签类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="system">系统标签</a-option>
            <a-option value="custom">自定义标签</a-option>
            <a-option value="behavior">行为标签</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="标签颜色">
          <a-color-picker v-model="form.color" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" />
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

const form = reactive({ name: '', type: '', color: '', description: '' });

const tags = ref([
  { id: 1, name: '高价值用户', type: 'system', typeText: '系统', typeColor: 'red', userCount: 856, description: '消费金额前10%用户' },
  { id: 2, name: '活跃用户', type: 'behavior', typeText: '行为', typeColor: 'blue', userCount: 2560, description: '30天内有登录' },
  { id: 3, name: '潜在流失', type: 'behavior', typeText: '行为', typeColor: 'orange', userCount: 428, description: '60天未活跃' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标签名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '用户数', dataIndex: 'userCount', width: 100 },
  { title: '描述', dataIndex: 'description' },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.tag-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
