<template>
  <div class="digital-twin-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="数字孪生体" :value="5" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="实时同步" :value="5" status="success" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="最后更新" value="刚刚" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>数字孪生体</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="realtime" title="实时状态">
          <a-row :gutter="16">
            <a-col :span="8" v-for="pet in twins" :key="pet.id">
              <a-card size="small" class="twin-card">
                <div class="pet-avatar">{{ pet.avatar }}</div>
                <div class="pet-name">{{ pet.name }}</div>
                <a-descriptions :column="1" size="small">
                  <a-descriptions-item label="心率">{{ pet.heartRate }} bpm</a-descriptions-item>
                  <a-descriptions-item label="呼吸">{{ pet.breathing }} /min</a-descriptions-item>
                  <a-descriptions-item label="体温">{{ pet.temperature }}°C</a-descriptions-item>
                  <a-descriptions-item label="活跃度">{{ pet.activity }}%</a-descriptions-item>
                </a-descriptions>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="behavior" title="行为预测">
          <a-table :columns="predictionColumns" :data="predictions" :pagination="pagination">
            <template #confidence="{ record }">
              <a-progress :percent="record.confidence" :color="record.confidence >= 80 ? '#00B42A' : '#FF7D00'" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="highlights" title="精彩瞬间">
          <a-row :gutter="16">
            <a-col :span="6" v-for="h in highlights" :key="h.id">
              <a-card size="small" :cover="h.image">
                <div>{{ h.title }}</div>
                <div style="color: #86909c; font-size: 12px;">{{ h.time }}</div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="replay" title="历史回放">
          <a-table :columns="replayColumns" :data="replays" :pagination="pagination">
            <template #actions="{ record }">
              <a-button type="primary" size="small" @click="handleReplay(record)">回放</a-button>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const twins = ref([
  { id: 1, name: '小黄', avatar: '🐕', heartRate: 85, breathing: 18, temperature: 38.5, activity: 75 },
  { id: 2, name: '小红', avatar: '🐕', heartRate: 78, breathing: 15, temperature: 38.2, activity: 60 },
  { id: 3, name: '咪咪', avatar: '🐱', heartRate: 120, breathing: 25, temperature: 38.8, activity: 80 },
]);

const predictions = ref([
  { petName: '小黄', behavior: '即将入睡', confidence: 85, time: '5分钟后' },
  { petName: '小红', behavior: '准备觅食', confidence: 72, time: '10分钟后' },
  { petName: '咪咪', behavior: '精力充沛', confidence: 90, time: '持续1小时' },
]);

const highlights = ref([
  { id: 1, title: '第一次学会握手', time: '2026-03-25', image: '' },
  { id: 2, title: '海边玩耍', time: '2026-03-20', image: '' },
  { id: 3, title: '生日派对', time: '2026-03-15', image: '' },
]);

const replays = ref([
  { id: 1, petName: '小黄', startTime: '2026-03-28 10:00:00', endTime: '2026-03-28 11:00:00', duration: '1小时' },
  { id: 2, petName: '小红', startTime: '2026-03-27 15:00:00', endTime: '2026-03-27 16:00:00', duration: '1小时' },
]);

const predictionColumns = [
  { title: '宠物', dataIndex: 'petName', width: 120 },
  { title: '预测行为', dataIndex: 'behavior', width: 150 },
  { title: '置信度', slotName: 'confidence', width: 150 },
  { title: '预计时间', dataIndex: 'time', width: 120 },
];

const replayColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 120 },
  { title: '开始时间', dataIndex: 'startTime', width: 160 },
  { title: '结束时间', dataIndex: 'endTime', width: 160 },
  { title: '时长', dataIndex: 'duration', width: 100 },
  { title: '操作', slotName: 'actions', width: 100 },
];

const handleReplay = (record: any) => {};
</script>

<style scoped>
.digital-twin-container { padding: 20px; }
.twin-card { text-align: center; margin-bottom: 12px; }
.pet-avatar { font-size: 48px; margin-bottom: 8px; }
.pet-name { font-weight: bold; margin-bottom: 8px; }
</style>
