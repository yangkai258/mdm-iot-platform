<template>
  <div class="sleep-analysis-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="今日睡眠时长" :value="stats.todaySleep" suffix="小时" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="睡眠质量评分" :value="stats.qualityScore" suffix="/100" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="异常次数" :value="stats.abnormalCount" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>睡眠分析</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="today" title="今日数据">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="睡眠趋势">
                <a-chart :option="sleepTrendChart" style="height: 250px;" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="睡眠阶段">
                <a-chart :option="sleepPhaseChart" style="height: 250px;" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="history" title="历史记录">
          <a-table :columns="columns" :data="historyData" :pagination="pagination">
            < template #quality="{ record }">
              <a-rate :model-value="record.quality" disabled />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="report" title="睡眠报告">
          <a-card>
            <a-descriptions :column="2" bordered>
              <a-descriptions-item label="平均睡眠时长">{{ report.avgDuration }}</a-descriptions-item>
              <a-descriptions-item label="平均入睡时间">{{ report.avgSleepTime }}</a-descriptions-item>
              <a-descriptions-item label="平均起床时间">{{ report.avgWakeTime }}</a-descriptions-item>
              <a-descriptions-item label="睡眠效率">{{ report.efficiency }}%</a-descriptions-item>
              <a-descriptions-item label="深睡比例">{{ report.deepSleepRatio }}%</a-descriptions-item>
              <a-descriptions-item label="REM比例">{{ report.remRatio }}%</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const stats = reactive({ todaySleep: 8.5, qualityScore: 85, abnormalCount: 0 });
const pagination = reactive({ current: 1, pageSize: 10, total: 7 });

const historyData = ref([
  { date: '2026-03-28', sleepTime: '22:30', wakeTime: '07:00', duration: 8.5, quality: 5, deepSleep: 2.5, remSleep: 1.8, awakeCount: 1 },
  { date: '2026-03-27', sleepTime: '23:00', wakeTime: '07:30', duration: 8.5, quality: 4, deepSleep: 2.2, remSleep: 1.5, awakeCount: 2 },
  { date: '2026-03-26', sleepTime: '22:45', wakeTime: '06:45', duration: 8.0, quality: 4, deepSleep: 2.0, remSleep: 1.6, awakeCount: 1 },
]);

const report = reactive({ avgDuration: '8.2小时', avgSleepTime: '23:00', avgWakeTime: '07:15', efficiency: 92, deepSleepRatio: 28, remRatio: 22 });

const sleepTrendChart = {
  xAxis: { type: 'category', data: ['22:00', '23:00', '00:00', '01:00', '02:00', '03:00', '04:00', '05:00', '06:00', '07:00'] },
  yAxis: { type: 'value' },
  series: [
    { name: '清醒', type: 'bar', stack: 'total', data: [0, 0, 5, 0, 0, 0, 10, 15, 20, 50] },
    { name: 'REM', type: 'bar', stack: 'total', data: [0, 20, 30, 25, 20, 25, 15, 10, 0, 0] },
    { name: '浅睡', type: 'bar', stack: 'total', data: [30, 40, 35, 40, 45, 35, 30, 25, 20, 0] },
    { name: '深睡', type: 'bar', stack: 'total', data: [70, 40, 30, 35, 35, 40, 45, 50, 60, 50] },
  ],
};

const sleepPhaseChart = { series: [{ name: '睡眠阶段', type: 'pie', radius: '60%', data: [
    { value: 35, name: '深睡' },
    { value: 25, name: 'REM' },
    { value: 30, name: '浅睡' },
    { value: 10, name: '清醒' },
  ]}] };

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '入睡', dataIndex: 'sleepTime', width: 80 },
  { title: '起床', dataIndex: 'wakeTime', width: 80 },
  { title: '时长(h)', dataIndex: 'duration', width: 80 },
  { title: '质量', slotName: 'quality', width: 150 },
  { title: '深睡(h)', dataIndex: 'deepSleep', width: 80 },
  { title: 'REM(h)', dataIndex: 'remSleep', width: 80 },
  { title: '清醒次数', dataIndex: 'awakeCount', width: 100 },
];
</script>

<style scoped>
.sleep-analysis-container { padding: 20px; }
</style>
