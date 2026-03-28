<template>
  <div class="live-streaming-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="直播场次" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="观看人数" :value="12580" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="进行中" :value="2" status="processing" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="弹幕数" :value="8560" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>直播管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            创建直播
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="streams" :loading="loading" :pagination="pagination">
        <template #status="{ record }">
          <a-badge v-if="record.status === 'live'" status="processing" :text="record.statusText" />
          <a-badge v-else :status="record.status === 'ended' ? 'default' : 'warning'" :text="record.statusText" />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link v-if="record.status === 'live'" @click="handleView(record)">观看</a-link>
            <a-link @click="handleData(record)">数据</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="createVisible" title="创建直播" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="直播标题" required>
          <a-input v-model="form.title" placeholder="请输入直播标题" />
        </a-form-item>
        <a-form-item label="关联宠物">
          <a-select v-model="form.petId" placeholder="选择宠物">
            <a-option value="P001">小黄</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="直播类型">
          <a-select v-model="form.type" placeholder="选择类型">
            <a-option value="training">训练直播</a-option>
            <a-option value="interaction">互动直播</a-option>
            <a-option value="show">展示直播</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预计时长">
          <a-input-number v-model="form.duration" :min="30" /> 分钟
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

const form = reactive({ title: '', petId: '', type: '', duration: 60 });

const streams = ref([
  { id: 1, title: '小黄日常训练直播', petName: '小黄', type: 'training', typeText: '训练', status: 'live', statusText: '直播中', viewers: 856, duration: 45 },
  { id: 2, title: '宠物互动时间', petName: '小红', type: 'interaction', typeText: '互动', status: 'ended', statusText: '已结束', viewers: 2580, duration: 90 },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '直播标题', dataIndex: 'title', width: 200 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '类型', dataIndex: 'typeText', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '观看人数', dataIndex: 'viewers', width: 100 },
  { title: '时长', dataIndex: 'duration', width: 80 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const handleCreate = () => { createVisible.value = true; };
const handleView = (record: any) => {};
const handleData = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); createVisible.value = false; };
</script>

<style scoped>
.live-streaming-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
