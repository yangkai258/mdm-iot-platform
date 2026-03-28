<template>
  <div class="health-container">
    <a-card>
      <template #title>
        <span>健康追踪</span>
      </template>
      
      <a-tabs default-active-key="exercise">
        <a-tab-pane key="exercise" title="运动记录">
          <a-table :columns="exerciseColumns" :data="exerciseData" :pagination="pagination">
            <template #duration="{ record }">
              {{ record.duration }}分钟
            </template>
            <template #calories="{ record }">
              {{ record.calories }} kcal
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="sleep" title="睡眠分析">
          <a-table :columns="sleepColumns" :data="sleepData" :pagination="pagination">
            <template #duration="{ record }">
              {{ record.hours }}小时{{ record.minutes }}分钟
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="vitals" title="生命体征">
          <a-descriptions bordered>
            <a-descriptions-item label="平均心率">72 bpm</a-descriptions-item>
            <a-descriptions-item label="平均体温">37.2°C</a-descriptions-item>
            <a-descriptions-item label="活动时长">2.5小时/天</a-descriptions-item>
            <a-descriptions-item label="睡眠质量">良好</a-descriptions-item>
          </a-descriptions>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });

const exerciseData = ref([
  { id: 'E001', petName: '小白', type: '散步', duration: 30, calories: 150, distance: 2.5, createdAt: '2026-03-28 08:00:00' },
  { id: 'E002', petName: '小黄', type: '跑步', duration: 20, calories: 200, distance: 3.2, createdAt: '2026-03-28 07:30:00' },
]);

const sleepData = ref([
  { id: 'S001', petName: '小白', hours: 8, minutes: 30, quality: '良好', score: 85, createdAt: '2026-03-28 07:00:00' },
  { id: 'S002', petName: '小黄', hours: 7, minutes: 45, quality: '一般', score: 70, createdAt: '2026-03-28 07:15:00' },
]);

const exerciseColumns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '运动类型', dataIndex: 'type', width: 100 },
  { title: '时长', slotName: 'duration', width: 100 },
  { title: '消耗热量', slotName: 'calories', width: 100 },
  { title: '距离(km)', dataIndex: 'distance', width: 100 },
  { title: '时间', dataIndex: 'createdAt', width: 160 },
];

const sleepColumns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '睡眠时长', slotName: 'duration', width: 120 },
  { title: '睡眠质量', dataIndex: 'quality', width: 100 },
  { title: '评分', dataIndex: 'score', width: 80 },
  { title: '时间', dataIndex: 'createdAt', width: 160 },
];
</script>

<style scoped>
.health-container { padding: 20px; }
</style>
