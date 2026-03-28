<template>
  <div class="multi-pet-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="协作组" :value="5" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="在线宠物" :value="12" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="协作次数" :value="128" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>多宠物协作</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新建协作组
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="groups" :pagination="pagination">
        <template #pets="{ record }">
          <a-space>
            <a-tag v-for="pet in record.petNames" :key="pet">{{ pet }}</a-tag>
          </a-space>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" title="协作组配置" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="协作组名称" required>
          <a-input v-model="form.name" placeholder="请输入协作组名称" />
        </a-form-item>
        <a-form-item label="参与宠物">
          <a-select v-model="form.petIds" multiple placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
            <a-option value="P002">小红</a-option>
            <a-option value="P003">咪咪</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="协作类型">
          <a-checkbox-group v-model="form.types">
            <a-checkbox value="play">互动玩耍</a-checkbox>
            <a-checkbox value="training">协作训练</a-checkbox>
            <a-checkbox value="social">社交活动</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', petIds: [], types: [] });

const groups = ref([
  { id: 1, name: '金毛俱乐部', petNames: ['小黄', '旺财', '球球'], types: ['play', 'training'], createdAt: '2026-03-20' },
  { id: 2, name: '猫咪联盟', petNames: ['咪咪', '小白'], types: ['play', 'social'], createdAt: '2026-03-15' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '协作组名称', dataIndex: 'name', width: 150 },
  { title: '参与宠物', slotName: 'pets', width: 250 },
  { title: '协作类型', dataIndex: 'types', width: 200 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleView = (record: any) => {};
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.multi-pet-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
