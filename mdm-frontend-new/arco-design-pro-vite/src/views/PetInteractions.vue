<template>
  <div class="pet-interactions-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="今日互动" :value="stats.todayInteractions" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="互动时长" :value="stats.totalDuration" suffix="分钟" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均情绪" :value="stats.avgMood" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="用户评分" :value="stats.userRating" suffix="/5" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>宠物互动记录</span>
          <a-space>
            <a-button @click="handleExport">导出</a-button>
            <a-button type="primary" @click="handleStartSession">
              <template #icon><icon-plus /></template>
              开始互动
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
        </template>
        <template #mood="{ record }">
          <a-tag :color="getMoodColor(record.moodEnd)">{{ record.moodEndText }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleReplay(record)">回放</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="互动详情" :width="700">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="互动ID">{{ currentSession.id }}</a-descriptions-item>
        <a-descriptions-item label="宠物">{{ currentSession.petName }}</a-descriptions-item>
        <a-descriptions-item label="互动类型">
          <a-tag>{{ currentSession.typeText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="开始情绪">
          <a-tag>{{ currentSession.moodStartText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="结束情绪">
          <a-tag :color="getMoodColor(currentSession.moodEnd)">{{ currentSession.moodEndText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="持续时间">{{ currentSession.duration }}</a-descriptions-item>
        <a-descriptions-item label="时间">{{ currentSession.time }}</a-descriptions-item>
        <a-descriptions-item label="用户评分">{{ currentSession.userRating }}/5</a-descriptions-item>
        <a-descriptions-item label="互动内容" :span="2">{{ currentSession.content }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const stats = reactive({ todayInteractions: 28, totalDuration: 156, avgMood: '开心', userRating: 4.5 });

const data = ref([
  { id: 'I001', petName: '小黄', type: 'play', typeText: '玩耍', moodStart: 'calm', moodStartText: '平静', moodEnd: 'happy', moodEndText: '开心', duration: '30分钟', userRating: 5, time: '2026-03-28 18:50:00', content: '和主人玩飞盘游戏' },
  { id: 'I002', petName: '小黄', type: 'training', typeText: '训练', moodStart: 'excited', moodStartText: '兴奋', moodEnd: 'happy', moodEndText: '开心', duration: '20分钟', userRating: 4, time: '2026-03-28 18:00:00', content: '练习坐下和握手' },
  { id: 'I003', petName: '小红', type: 'chat', typeText: '对话', moodStart: 'sad', moodStartText: '低落', moodEnd: 'calm', moodEndText: '平静', duration: '15分钟', userRating: 5, time: '2026-03-28 17:00:00', content: '陪伴聊天' },
  { id: 'I004', petName: '咪咪', type: 'play', typeText: '玩耍', moodStart: 'happy', moodStartText: '开心', moodEnd: 'happy', moodEndText: '开心', duration: '40分钟', userRating: 5, time: '2026-03-28 16:00:00', content: '逗猫棒游戏' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const detailVisible = ref(false);
const currentSession = ref<any>({});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 100 },
  { title: '宠物', dataIndex: 'petName', width: 80 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '开始情绪', dataIndex: 'moodStartText', width: 100 },
  { title: '结束情绪', slotName: 'mood', width: 100 },
  { title: '时长', dataIndex: 'duration', width: 100 },
  { title: '评分', dataIndex: 'userRating', width: 80 },
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '操作', slotName: 'actions', width: 120 },
];

const getTypeColor = (t: string) => ({ play: 'blue', training: 'orange', chat: 'green', feeding: 'purple' }[t] || 'default');
const getMoodColor = (m: string) => ({ happy: 'green', calm: 'blue', sad: 'orange', excited: 'purple', anxious: 'red' }[m] || 'default');

const handleExport = () => {};
const handleStartSession = () => {};
const handleView = (record: any) => { currentSession.value = record; detailVisible.value = true; };
const handleReplay = (record: any) => {};
</script>

<style scoped>
.pet-interactions-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
