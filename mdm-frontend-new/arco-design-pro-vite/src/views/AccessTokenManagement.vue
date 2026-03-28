<template>
  <div class="access-token-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="访问令牌" :value="36" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="活跃令牌" :value="28" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="API调用" :value="12580" suffix="次" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>开放平台令牌管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建令牌
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="tokens" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-badge :status="record.status === 'active' ? 'success' : 'default'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link v-if="record.status === 'active'" @click="handleRevoke(record)">撤销</a-link>
            <a-link @click="handleView(record)">详情</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建访问令牌" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="令牌名称" required>
          <a-input v-model="form.name" placeholder="请输入令牌名称" />
        </a-form-item>
        <a-form-item label="令牌类型">
          <a-radio-group v-model="form.type">
            <a-radio value="read">只读</a-radio>
            <a-radio value="write">读写</a-radio>
            <a-radio value="admin">管理员</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="权限范围">
          <a-checkbox-group v-model="form.scopes">
            <a-checkbox value="devices">设备管理</a-checkbox>
            <a-checkbox value="members">会员管理</a-checkbox>
            <a-checkbox value="pets">宠物管理</a-checkbox>
            <a-checkbox value="ai">AI服务</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="过期时间">
          <a-select v-model="form.expiresIn">
            <a-option value="7d">7天</a-option>
            <a-option value="30d">30天</a-option>
            <a-option value="90d">90天</a-option>
            <a-option value="never">永不过期</a-option>
          </a-select>
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

const form = reactive({ name: '', type: 'read', scopes: [], expiresIn: '30d' });

const tokens = ref([
  { id: 1, name: 'MyApp-Production', type: 'read', typeText: '只读', key: 'tok_live_xxxxx', scopes: 'devices,members', status: 'active', statusText: '活跃', expiresAt: '2026-04-28', createdAt: '2026-03-28' },
  { id: 2, name: 'MyApp-Dev', type: 'write', typeText: '读写', key: 'tok_test_xxxxx', scopes: 'devices,members,pets', status: 'active', statusText: '活跃', expiresAt: '2026-04-15', createdAt: '2026-03-15' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '令牌名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '令牌Key', dataIndex: 'key', width: 150, ellipsis: true },
  { title: '权限范围', dataIndex: 'scopes', width: 150 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '过期时间', dataIndex: 'expiresAt', width: 120 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const handleCreate = () => { createVisible.value = true; };
const handleRevoke = (record: any) => {};
const handleView = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.access-token-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
