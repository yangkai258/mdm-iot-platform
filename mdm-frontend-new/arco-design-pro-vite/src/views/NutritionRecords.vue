<template>
  <div class="nutrition-records-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="今日摄入" :value="stats.todayIntake" suffix="kcal" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="目标" :value="stats.target" suffix="kcal" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="剩余" :value="stats.remaining" suffix="kcal" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日记录" :value="stats.records" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>饮食记录</span>
          <a-button type="primary" @click="handleAdd">
            <template #icon><icon-plus /></template>
            添加记录
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="addVisible" title="添加饮食记录" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
            <a-option value="P002">小红</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="食物类型" required>
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="breakfast">早餐</a-option>
            <a-option value="lunch">午餐</a-option>
            <a-option value="dinner">晚餐</a-option>
            <a-option value="snack">零食</a-option>
            <a-option value="treat">宠物零食</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="食物名称">
          <a-input v-model="form.foodName" placeholder="如: 狗粮、鸡胸肉" />
        </a-form-item>
        <a-form-item label="重量">
          <a-input-number v-model="form.weight" :min="0" /> 克
        </a-form-item>
        <a-form-item label="热量">
          <a-input-number v-model="form.calories" :min="0" /> kcal
        </a-form-item>
        <a-form-item label="时间">
          <a-time-picker v-model="form.time" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const stats = reactive({ todayIntake: 580, target: 800, remaining: 220, records: 4 });

const data = ref([
  { id: 1, petName: '小黄', type: 'breakfast', typeText: '早餐', foodName: '狗粮', weight: 100, calories: 150, time: '08:00', note: '' },
  { id: 2, petName: '小黄', type: 'lunch', typeText: '午餐', foodName: '鸡胸肉', weight: 80, calories: 120, time: '12:00', note: '煮熟切块' },
  { id: 3, petName: '小黄', type: 'snack', typeText: '零食', foodName: '宠物饼干', weight: 30, calories: 80, time: '15:00', note: '' },
  { id: 4, petName: '小黄', type: 'dinner', typeText: '晚餐', foodName: '狗粮+蔬菜', weight: 120, calories: 230, time: '18:00', note: '' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const addVisible = ref(false);

const form = reactive({ petId: '', type: '', foodName: '', weight: 0, calories: 0, time: null, note: '' });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '食物', dataIndex: 'foodName', width: 150 },
  { title: '重量(g)', dataIndex: 'weight', width: 100 },
  { title: '热量(kcal)', dataIndex: 'calories', width: 120 },
  { title: '时间', dataIndex: 'time', width: 80 },
  { title: '备注', dataIndex: 'note' },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const getTypeColor = (t: string) => ({ breakfast: 'orange', lunch: 'green', dinner: 'blue', snack: 'purple', treat: 'cyan' }[t] || 'default');

const handleAdd = () => { addVisible.value = true; };
const handleEdit = (record: any) => { Object.assign(form, record); addVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); addVisible.value = false; };
</script>

<style scoped>
.nutrition-records-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
