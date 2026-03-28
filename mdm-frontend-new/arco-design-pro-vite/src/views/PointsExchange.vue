<template>
  <div class="points-exchange-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="兑换商品" :value="128" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="兑换次数" :value="5680" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="消耗积分" :value="2560000" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>积分商城</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加商品
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="products" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-tag :color="record.stock > 0 ? 'green' : 'red'">{{ record.stock > 0 ? '有货' : '缺货' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑商品' : '添加商品'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="商品名称" required>
          <a-input v-model="form.name" placeholder="请输入商品名称" />
        </a-form-item>
        <a-form-item label="积分价格">
          <a-input-number v-model="form.points" :min="0" />
        </a-form-item>
        <a-form-item label="库存">
          <a-input-number v-model="form.stock" :min="0" />
        </a-form-item>
        <a-form-item label="商品描述">
          <a-textarea v-model="form.description" :rows="3" />
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

const form = reactive({ name: '', points: 0, stock: 0, description: '' });

const products = ref([
  { id: 1, name: '宠物玩具套装', points: 500, stock: 50, exchangedCount: 128 },
  { id: 2, name: '一个月高级会员', points: 1000, stock: 100, exchangedCount: 256 },
  { id: 3, name: '宠物食品大礼包', points: 2000, stock: 20, exchangedCount: 45 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '商品名称', dataIndex: 'name', width: 200 },
  { title: '积分价格', dataIndex: 'points', width: 120 },
  { title: '库存', dataIndex: 'stock', width: 80 },
  { title: '已兑换', dataIndex: 'exchangedCount', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.points-exchange-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
