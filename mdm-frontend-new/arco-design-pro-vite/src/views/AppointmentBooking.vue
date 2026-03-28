<template>
  <div class="appointment-booking-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="预约总数" :value="568" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日预约" :value="28" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="待确认" :value="5" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="完成率" :value="92" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>预约服务管理</span>
          <a-space>
            <a-input-search placeholder="搜索宠物/主人" style="width: 200px;" />
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              新建预约
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-table :columns="columns" :data="appointments" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link v-if="record.status === 'pending'" @click="handleConfirm(record)">确认</a-link>
            <a-link @click="handleView(record)">详情</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="新建预约" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="宠物" required>
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="服务类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="checkup">体检</a-option>
            <a-option value="vaccine">疫苗</a-option>
            <a-option value="grooming">美容</a-option>
            <a-option value="training">训练</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预约时间">
          <a-date-picker v-model="form.appointmentTime" style="width: 100%;" />
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
const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const createVisible = ref(false);

const form = reactive({ petId: '', type: '', appointmentTime: '', note: '' });

const appointments = ref([
  { id: 1, petName: '小黄', ownerName: '张三', type: 'checkup', typeText: '体检', appointmentTime: '2026-03-29 10:00', status: 'pending', statusText: '待确认', note: '' },
  { id: 2, petName: '小红', ownerName: '李四', type: 'vaccine', typeText: '疫苗', appointmentTime: '2026-03-28 14:00', status: 'confirmed', statusText: '已确认', note: '' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '主人', dataIndex: 'ownerName', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '预约时间', dataIndex: 'appointmentTime', width: 160 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const getStatusColor = (s: string) => ({ pending: 'orange', confirmed: 'blue', completed: 'green', cancelled: 'gray' }[s] || 'default');

const handleCreate = () => { createVisible.value = true; };
const handleConfirm = (record: any) => {};
const handleView = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.appointment-booking-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
