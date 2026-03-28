<template>
  <div class="action-library-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="动作总数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已学习" :value="45" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="热门动作" :value="10" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="社区贡献" :value="73" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>动作资源库</span>
          <a-space>
            <a-button @click="handleUpload">
              <template #icon><icon-upload /></template>
              上传动作
            </a-button>
            <a-button type="primary" @click="handleCreate">
              <template #icon><icon-plus /></template>
              创建动作
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="all" title="全部动作">
          <a-row :gutter="16">
            <a-col :span="4" v-for="action in actions" :key="action.id">
              <a-card size="small" class="action-card">
                <div class="action-preview">{{ action.icon }}</div>
                <div class="action-name">{{ action.name }}</div>
                <div class="action-author">by {{ action.author }}</div>
                <div class="action-stats">
                  <span>⬇️ {{ action.downloads }}</span>
                  <span>❤️ {{ action.likes }}</span>
                </div>
                <a-button type="primary" size="small" @click="handleDownload(action)">下载</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="learning" title="学习中">
          <a-table :columns="learningColumns" :data="learning" :pagination="pagination">
            <template #progress="{ record }">
              <a-progress :percent="record.progress" :status="record.progress === 100 ? 'success' : 'normal'" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="mastered" title="已掌握">
          <a-table :columns="masteredColumns" :data="mastered" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 20, total: 10 });

const actions = ref([
  { id: 1, name: '后空翻', icon: '🤸', author: 'User001', downloads: 1234, likes: 456 },
  { id: 2, name: '撒娇', icon: '🥺', author: 'User002', downloads: 987, likes: 321 },
  { id: 3, name: '跳舞', icon: '💃', author: 'User001', downloads: 876, likes: 234 },
  { id: 4, name: '敬礼', icon: '敬礼', author: 'User003', downloads: 654, likes: 123 },
]);

const learning = ref([
  { id: 1, name: '后空翻', icon: '🤸', petName: '小黄', progress: 45, lastPracticed: '2026-03-28' },
  { id: 2, name: '撒娇', icon: '🥺', petName: '小黄', progress: 78, lastPracticed: '2026-03-28' },
]);

const mastered = ref([
  { id: 1, name: '坐下', petName: '小黄', masteredAt: '2026-03-20' },
  { id: 2, name: '握手', petName: '小黄', masteredAt: '2026-03-15' },
]);

const learningColumns = [
  { title: '动作', dataIndex: 'name', width: 150 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '学习进度', slotName: 'progress', width: 200 },
  { title: '最近练习', dataIndex: 'lastPracticed', width: 120 },
];

const masteredColumns = [
  { title: '动作', dataIndex: 'name', width: 150 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '掌握时间', dataIndex: 'masteredAt', width: 120 },
];

const handleUpload = () => {};
const handleCreate = () => {};
const handleDownload = (action: any) => {};
</script>

<style scoped>
.action-library-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.action-card { text-align: center; margin-bottom: 12px; }
.action-preview { font-size: 48px; margin-bottom: 8px; }
.action-name { font-weight: bold; }
.action-author { color: #86909c; font-size: 12px; }
.action-stats { margin: 8px 0; color: #86909c; font-size: 12px; }
</style>
