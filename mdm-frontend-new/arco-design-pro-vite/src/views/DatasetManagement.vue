<template>
  <div class="dataset-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="数据集" :value="15" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="总数据量" :value="1024" suffix="MB" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="已标注" :value="85" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>数据集管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建数据集
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-progress :percent="record.annotated" size="small" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">查看</a-link>
            <a-link @click="handleAnnotate(record)">标注</a-link>
            <a-link @click="handleExport(record)">导出</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建数据集" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="数据集名称" required>
          <a-input v-model="form.name" placeholder="请输入数据集名称" />
        </a-form-item>
        <a-form-item label="数据集类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="image">图像数据</a-option>
            <a-option value="audio">语音数据</a-option>
            <a-option value="text">文本数据</a-option>
            <a-option value="mixed">混合数据</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :rows="2" />
        </a-form-item>
        <a-form-item label="标注任务">
          <a-checkbox-group v-model="form.tasks">
            <a-checkbox value="classification">分类标注</a-checkbox>
            <a-checkbox value="detection">目标检测</a-checkbox>
            <a-checkbox value="segmentation">分割标注</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, name: '宠物图像数据集', type: 'image', typeText: '图像数据', size: 512, count: 10000, annotated: 95, description: '用于宠物识别训练' },
  { id: 2, name: '语音指令数据集', type: 'audio', typeText: '语音数据', size: 256, count: 5000, annotated: 80, description: '用于语音指令识别' },
  { id: 3, name: '情绪对话数据集', type: 'text', typeText: '文本数据', size: 128, count: 8000, annotated: 90, description: '用于情绪识别训练' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const createVisible = ref(false);

const form = reactive({ name: '', type: 'image', description: '', tasks: [] });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '数据集名称', dataIndex: 'name', width: 200 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '大小(MB)', dataIndex: 'size', width: 100 },
  { title: '数据量', dataIndex: 'count', width: 100 },
  { title: '已标注', slotName: 'status', width: 120 },
  { title: '描述', dataIndex: 'description' },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
];

const handleCreate = () => { createVisible.value = true; };
const handleView = (record: any) => {};
const handleAnnotate = (record: any) => {};
const handleExport = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.dataset-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
