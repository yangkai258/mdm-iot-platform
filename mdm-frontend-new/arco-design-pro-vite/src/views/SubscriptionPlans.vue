<template>
  <div class="plans-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>订阅套餐</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增套餐
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="record.type === 'subscription' ? 'blue' : 'orange'">
            {{ record.type === 'subscription' ? '订阅' : '按量' }}
          </a-tag>
        </template>
        <template #price="{ record }">
          <span style="color: #F53F3F; font-weight: bold;">¥{{ record.price }}</span>
          <span v-if="record.period">/{{ record.period }}</span>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '上架' : '下架' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link @click="handleUsers(record)">用户</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 用户列表弹窗 -->
    <a-modal v-model:visible="usersVisible" title="套餐用户" :width="800">
      <a-table :columns="userColumns" :data="planUsers" :pagination="pagination">
        <template #subscription="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-modal>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑套餐' : '新增套餐'" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="套餐名称" required>
          <a-input v-model="form.name" placeholder="如: 年度高级版" />
        </a-form-item>
        <a-form-item label="套餐类型">
          <a-radio-group v-model="form.type">
            <a-radio value="subscription">订阅</a-radio>
            <a-radio value="usage">按量付费</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="价格">
          <a-input-number v-model="form.price" :min="0" :precision="2" />
        </a-form-item>
        <a-form-item label="计费周期">
          <a-select v-model="form.period" placeholder="选择周期">
            <a-option value="月">月</a-option>
            <a-option value="季">季</a-option>
            <a-option value="年">年</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="功能权限">
          <a-checkbox-group v-model="form.features">
            <a-checkbox value="devices">设备数量: {{ formlimits.devices }}台</a-checkbox>
            <a-checkbox value="ai_chat">AI对话</a-checkbox>
            <a-checkbox value="cloud_backup">云端备份</a-checkbox>
            <a-checkbox value="priority">优先客服</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
        <a-form-item label="设备限额">
          <a-input-number v-model="form.limits.devices" :min="1" />
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
  { id: 1, name: '免费版', type: 'subscription', price: 0, period: '-', features: ['devices:3', 'ai_chat'], limits: { devices: 3 }, enabled: true, userCount: 1250 },
  { id: 2, name: '月度基础版', type: 'subscription', price: 29, period: '月', features: ['devices:10', 'ai_chat', 'cloud_backup'], limits: { devices: 10 }, enabled: true, userCount: 580 },
  { id: 3, name: '年度基础版', type: 'subscription', price: 290, period: '年', features: ['devices:10', 'ai_chat', 'cloud_backup'], limits: { devices: 10 }, enabled: true, userCount: 320 },
  { id: 4, name: '月度高级版', type: 'subscription', price: 99, period: '月', features: ['devices:50', 'ai_chat', 'cloud_backup', 'priority'], limits: { devices: 50 }, enabled: true, userCount: 156 },
  { id: 5, name: '按量付费', type: 'usage', price: 0.01, period: '次', features: ['ai_chat'], limits: { devices: 999 }, enabled: false, userCount: 45 },
]);

const planUsers = ref([
  { id: 1, username: 'user1', phone: '138****1234', status: 'active', statusText: '正常', startDate: '2026-01-01', expireDate: '2027-01-01' },
  { id: 2, username: 'user2', phone: '139****5678', status: 'active', statusText: '正常', startDate: '2026-02-15', expireDate: '2027-02-15' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 5 });
const editVisible = ref(false);
const usersVisible = ref(false);
const isEdit = ref(false);

const form = reactive({
  name: '', type: 'subscription', price: 0, period: '月', features: [], enabled: true, limits: { devices: 10 },
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '套餐名称', dataIndex: 'name', width: 120 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '价格', slotName: 'price', width: 120 },
  { title: '设备限额', dataIndex: 'limits', width: 100 },
  { title: '用户数', dataIndex: 'userCount', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const userColumns = [
  { title: '用户', dataIndex: 'username', width: 120 },
  { title: '手机', dataIndex: 'phone', width: 130 },
  { title: '状态', slotName: 'subscription', width: 80 },
  { title: '开始日期', dataIndex: 'startDate', width: 120 },
  { title: '到期日期', dataIndex: 'expireDate', width: 120 },
];

const getStatusColor = (s: string) => ({ active: 'green', expired: 'red', cancelled: 'gray' }[s] || 'default');

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleUsers = (record: any) => { usersVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.plans-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
