<template>
  <div class="awards-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="成就总数" :value="48" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="已解锁" :value="1256" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="本周新增" :value="86" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>成就徽章管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建成就
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="awards" :loading="loading" :pagination="pagination">
        <template #rarity="{ record }">
          <a-tag :color="getRarityColor(record.rarity)">{{ record.rarityText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑成就' : '创建成就'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="成就名称" required>
          <a-input v-model="form.name" placeholder="请输入成就名称" />
        </a-form-item>
        <a-form-item label="成就图标">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="稀有度">
          <a-select v-model="form.rarity" placeholder="选择稀有度">
            <a-option value="common">普通</a-option>
            <a-option value="rare">稀有</a-option>
            <a-option value="epic">史诗</a-option>
            <a-option value="legendary">传说</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="解锁条件">
          <a-textarea v-model="form.condition" :rows="2" />
        </a-form-item>
        <a-form-item label="奖励积分">
          <a-input-number v-model="form.rewardPoints" :min="0" />
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

const form = reactive({ name: '', icon: '', rarity: '', condition: '', rewardPoints: 0 });

const awards = ref([
  { id: 1, name: '首次登录', icon: '🎯', rarity: 'common', rarityText: '普通', condition: '完成首次登录', unlockedCount: 8560, rewardPoints: 10 },
  { id: 2, name: '训练达人', icon: '🏆', rarity: 'rare', rarityText: '稀有', condition: '完成10次训练', unlockedCount: 1256, rewardPoints: 100 },
  { id: 3, name: '资深铲屎官', icon: '👑', rarity: 'epic', rarityText: '史诗', condition: '养宠超过1年', unlockedCount: 428, rewardPoints: 500 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '图标', dataIndex: 'icon', width: 80 },
  { title: '成就名称', dataIndex: 'name', width: 150 },
  { title: '稀有度', slotName: 'rarity', width: 100 },
  { title: '解锁条件', dataIndex: 'condition' },
  { title: '已解锁', dataIndex: 'unlockedCount', width: 100 },
  { title: '奖励积分', dataIndex: 'rewardPoints', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const getRarityColor = (r: string) => ({ common: 'gray', rare: 'blue', epic: 'purple', legendary: 'gold' }[r] || 'default');

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.awards-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
