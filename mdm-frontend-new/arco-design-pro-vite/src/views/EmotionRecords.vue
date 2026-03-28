<template>
  <div class="emotion-records-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card>
          <a-statistic title="今日情绪记录" :value="stats.todayRecords" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic title="正面情绪占比" :value="stats.positiveRate" suffix="%" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic title="需要关注" :value="stats.needAttention" status="warning" />
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>情绪记录</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="records" title="情绪记录">
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
            <template #emotion="{ record }">
              <a-tag :color="getEmotionColor(record.emotion)">{{ record.emotionText }}</a-tag>
            </template>
            <template #intensity="{ record }">
              <a-progress :percent="record.intensity" :color="getEmotionColor(record.emotion)" size="small" />
            </template>
            <template #actions="{ record }">
              <a-link @click="handleView(record)">详情</a-link>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="analysis" title="情绪分析">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="情绪分布">
                <a-chart :option="emotionChart" style="height: 300px;" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="情绪趋势">
                <a-chart :option="trendChart" style="height: 300px;" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="情绪详情" :width="600">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="记录ID">{{ currentRecord.id }}</a-descriptions-item>
        <a-descriptions-item label="设备ID">{{ currentRecord.deviceId }}</a-descriptions-item>
        <a-descriptions-item label="情绪">
          <a-tag :color="getEmotionColor(currentRecord.emotion)">{{ currentRecord.emotionText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="强度">{{ currentRecord.intensity }}%</a-descriptions-item>
        <a-descriptions-item label="时间">{{ currentRecord.time }}</a-descriptions-item>
        <a-descriptions-item label="触发原因">{{ currentRecord.trigger }}</a-descriptions-item>
        <a-descriptions-item label="AI响应" :span="2">{{ currentRecord.aiResponse }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const stats = reactive({ todayRecords: 156, positiveRate: 72, needAttention: 5 });

const data = ref([
  { id: 'E001', deviceId: 'DEV001', petName: '小黄', emotion: 'happy', emotionText: '开心', intensity: 85, trigger: '和主人玩耍', aiResponse: '播放欢快音乐', time: '2026-03-28 18:50:00' },
  { id: 'E002', deviceId: 'DEV001', petName: '小黄', emotion: 'calm', emotionText: '平静', intensity: 60, trigger: '休息中', aiResponse: '降低音量', time: '2026-03-28 18:40:00' },
  { id: 'E003', deviceId: 'DEV002', petName: '小红', emotion: 'sad', emotionText: '低落', intensity: 70, trigger: '主人不在', aiResponse: '播放安抚音乐', time: '2026-03-28 18:30:00' },
  { id: 'E004', deviceId: 'DEV003', petName: '咪咪', emotion: 'happy', emotionText: '开心', intensity: 90, trigger: '被抚摸', aiResponse: '发出满足的声音', time: '2026-03-28 18:20:00' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const detailVisible = ref(false);
const currentRecord = ref<any>({});

const emotionChart = {
  tooltip: { trigger: 'item' },
  series: [{ name: '情绪分布', type: 'pie', radius: '60%', data: [
    { value: 45, name: '开心' },
    { value: 25, name: '平静' },
    { value: 15, name: '低落' },
    { value: 10, name: '兴奋' },
    { value: 5, name: '焦虑' },
  ]}],
};

const trendChart = {
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value', max: 100 },
  series: [
    { name: '开心', type: 'line', data: [65, 70, 75, 80, 72, 85, 88], smooth: true },
    { name: '平静', type: 'line', data: [30, 28, 25, 30, 32, 28, 25], smooth: true },
  ],
};

const columns = [
  { title: 'ID', dataIndex: 'id', width: 100 },
  { title: '设备', dataIndex: 'deviceId', width: 100 },
  { title: '宠物', dataIndex: 'petName', width: 80 },
  { title: '情绪', slotName: 'emotion', width: 100 },
  { title: '强度', slotName: 'intensity', width: 120 },
  { title: '触发原因', dataIndex: 'trigger', width: 150 },
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '操作', slotName: 'actions', width: 80 },
];

const getEmotionColor = (e: string) => ({ happy: 'green', calm: 'blue', sad: 'orange', excited: 'purple', anxious: 'red' }[e] || 'default');

const handleView = (record: any) => { currentRecord.value = record; detailVisible.value = true; };
</script>

<style scoped>
.emotion-records-container { padding: 20px; }
</style>
