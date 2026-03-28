<template>
  <div class="price-strategy-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="进行中策略" :value="8" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="累计优惠" :value="12580" prefix="¥" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="触达用户" :value="2560" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="转化率" :value="12.5" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>价格策略</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建策略
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="strategies" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #discount="{ record }">
          <span style="color: #F53F3F; font-weight: bold;">{{ record.discount }}</span>
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

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑策略' : '新建策略'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="策略名称" required>
          <a-input v-model="form.name" placeholder="请输入策略名称" />
        </a-form-item>
        <a-form-item label="优惠类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="discount">折扣</a-option>
            <a-option value="coupon">优惠券</a-option>
            <a-option value="package">套餐</a-option>
            <a-option value="flash">限时特价</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="折扣">
          <a-input-number v-model="form.discountValue" :min="0" :max="100" />
          <span style="margin-left: 8px;">{{ form.type === 'discount' ? '%' : '元' }}</span>
        </a-form-item>
        <a-form-item label="适用商品">
          <a-select v-model="form.targetType" placeholder="选择适用范围">
            <a-option value="all">全部商品</a-option>
            <a-option value="category">指定分类</a-option>
            <a-option value="product">指定商品</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="开始时间">
          <a-date-picker v-model="form.startTime" style="width: 100%;" />
        </a-form-item>
        <a-form-item label="结束时间">
          <a-date-picker v-model="form.endTime" style="width: 100%;" />
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

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', type: 'discount', discountValue: 10, targetType: 'all', startTime: '', endTime: '', enabled: true });

const strategies = ref([
  { id: 1, name: '新用户首单7折', type: 'discount', typeText: '折扣', discount: '30%', targetType: 'all', targetText: '全部', startTime: '2026-03-01', endTime: '2026-03-31', enabled: true },
  { id: 2, name: '会员日8折', type: 'discount', typeText: '折扣', discount: '20%', targetType: 'category', targetText: '宠物食品', startTime: '2026-03-15', endTime: '2026-03-15', enabled: true },
  { id: 3, name: '满减优惠券', type: 'coupon', typeText: '优惠券', discount: '¥20', targetType: 'all', targetText: '全部', startTime: '2026-03-20', endTime: '2026-04-20', enabled: true },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '策略名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '优惠', slotName: 'discount', width: 100 },
  { title: '适用范围', dataIndex: 'targetText', width: 100 },
  { title: '开始时间', dataIndex: 'startTime', width: 120 },
  { title: '结束时间', dataIndex: 'endTime', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.price-strategy-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
