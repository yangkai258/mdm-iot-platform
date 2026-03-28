<template>
  <div class="product-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="商品总数" :value="156" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="在售" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="库存预警" :value="5" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="本月销量" :value="2580" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>商品管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            添加商品
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="products" :loading="loading" :pagination="pagination">
        <template #image="{ record }">
          <a-image width="50" :src="record.image" />
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.statusText }}</a-tag>
        </template>
        <template #price="{ record }">
          <span style="color: #F53F3F;">¥{{ record.price }}</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑商品' : '添加商品'" :width="700" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="商品名称" required>
          <a-input v-model="form.name" placeholder="请输入商品名称" />
        </a-form-item>
        <a-form-item label="商品分类">
          <a-select v-model="form.category" placeholder="选择分类">
            <a-option value="food">宠物食品</a-option>
            <a-option value="toy">宠物玩具</a-option>
            <a-option value="health">宠物保健</a-option>
            <a-option value="accessory">宠物配件</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="商品图片">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="商品描述">
          <a-textarea v-model="form.description" :rows="3" />
        </a-form-item>
        <a-form-item label="价格">
          <a-input-number v-model="form.price" :min="0" :precision="2" style="width: 200px;" />
        </a-form-item>
        <a-form-item label="库存">
          <a-input-number v-model="form.stock" :min="0" style="width: 200px;" />
        </a-form-item>
        <a-form-item label="上架">
          <a-switch v-model="form.status" />
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

const form = reactive({ name: '', category: '', description: '', price: 0, stock: 0, status: true, image: '' });

const products = ref([
  { id: 1, name: '智能喂食器', category: 'food', categoryText: '宠物食品', image: 'https://placeholder.com/feed.jpg', price: 299, stock: 50, status: 'active', statusText: '在售' },
  { id: 2, name: '宠物摄像头', category: 'accessory', categoryText: '宠物配件', image: 'https://placeholder.com/cam.jpg', price: 199, stock: 30, status: 'active', statusText: '在售' },
  { id: 3, name: '自动饮水机', category: 'toy', categoryText: '宠物玩具', image: 'https://placeholder.com/water.jpg', price: 149, stock: 5, status: 'warning', statusText: '库存不足' },
]);

const columns = [
  { title: '商品图片', slotName: 'image', width: 80 },
  { title: '商品名称', dataIndex: 'name', width: 150 },
  { title: '分类', dataIndex: 'categoryText', width: 100 },
  { title: '价格', slotName: 'price', width: 100 },
  { title: '库存', dataIndex: 'stock', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.product-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
