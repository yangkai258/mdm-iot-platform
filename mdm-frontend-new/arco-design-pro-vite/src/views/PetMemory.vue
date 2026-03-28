<template>
  <div class="pet-memory-container">
    <a-card>
      <template #title>
        <span>宠物记忆系统</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="short" title="短期记忆">
          <a-table :columns="shortColumns" :data="shortTerm" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="long" title="长期记忆">
          <a-table :columns="longColumns" :data="longTerm" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="important" title="重要记忆">
          <a-table :columns="importantColumns" :data="importantMemories" :pagination="pagination">
            <template #type="{ record }">
              <a-tag color="gold">{{ record.typeText }}</a-tag>
            </template>
            <template #importance="{ record }">
              <a-rate :model-value="record.importance" disabled />
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 20, total: 5 });

const shortTerm = ref([
  { id: 'S001', type: 'play', typeText: '玩耍', content: '今天和主人玩了飞盘', time: '2026-03-28 10:00:00', duration: '30分钟' },
  { id: 'S002', type: 'food', typeText: '进食', content: '午餐吃了狗粮', time: '2026-03-28 12:00:00', duration: '10分钟' },
  { id: 'S003', type: 'walk', typeText: '散步', content: '公园散步', time: '2026-03-28 15:00:00', duration: '45分钟' },
]);

const longTerm = ref([
  { id: 'L001', type: 'command', typeText: '指令', content: '学会了"坐下"指令', learnedAt: '2026-02-15', mastery: '95%' },
  { id: 'L002', type: 'command', typeText: '指令', content: '学会了"握手"指令', learnedAt: '2026-02-20', mastery: '88%' },
  { id: 'L003', type: 'experience', typeText: '经历', content: '第一次去海边', learnedAt: '2026-01-10', mastery: '-' },
]);

const importantMemories = ref([
  { id: 'I001', type: 'milestone', typeText: '里程碑', content: '第一次学会游泳', date: '2026-03-01', importance: 5 },
  { id: 'I002', type: 'milestone', typeText: '里程碑', content: '完成了第一次服从训练', date: '2026-02-15', importance: 4 },
  { id: 'I003', type: 'trauma', typeText: '创伤', content: '被烟花吓到', date: '2026-01-28', importance: 5 },
]);

const shortColumns = [
  { title: 'ID', dataIndex: 'id', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '内容', dataIndex: 'content' },
  { title: '持续时间', dataIndex: 'duration', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const longColumns = [
  { title: 'ID', dataIndex: 'id', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '内容', dataIndex: 'content' },
  { title: '学习时间', dataIndex: 'learnedAt', width: 120 },
  { title: '掌握度', dataIndex: 'mastery', width: 100 },
];

const importantColumns = [
  { title: 'ID', dataIndex: 'id', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '内容', dataIndex: 'content' },
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '重要度', slotName: 'importance', width: 150 },
];

const getTypeColor = (t: string) => ({ play: 'blue', food: 'green', walk: 'orange', command: 'purple', experience: 'cyan', milestone: 'gold', trauma: 'red' }[t] || 'default');
</script>

<style scoped>
.pet-memory-container { padding: 20px; }
</style>
