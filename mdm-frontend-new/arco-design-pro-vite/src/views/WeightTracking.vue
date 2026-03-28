<template>
  <div class="weight-tracking-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="当前体重" :value="stats.currentWeight" suffix="kg" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="目标体重" :value="stats.targetWeight" suffix="kg" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="本周变化" :value="stats.weeklyChange" :precision="1" suffix="kg" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>体重追踪</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="trend" title="体重趋势">
          <a-chart :option="weightChart" style="height: 300px;" />
        </a-tab-pane>
        
        <a-tab-pane key="records" title="记录">
          <a-table :columns="columns" :data="records" :pagination="pagination">
            <template #change="{ record }">
              <span :style="{ color: record.change > 0 ? '#F53F3F' : record.change < 0 ? '#00B42A' : '#86909c' }">
                {{ record.change > 0 ? '+' : '' }}{{ record.change }}
              </span>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="goals" title="目标设定">
          <a-form layout="vertical" style="max-width: 400px;">
            <a-form-item label="目标体重">
              <a-input-number v-model="goal.targetWeight" :min="1" :precision="1" suffix="kg" />
            </a-form-item>
            <a-form-item label="目标日期">
              <a-date-picker v-model="goal.targetDate" style="width: 100%;" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveGoal">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const stats = reactive({ currentWeight: 12.5, targetWeight: 12.0, weeklyChange: -0.2 });
const pagination = reactive({ current: 1, pageSize: 10, total: 7 });

const records = ref([
  { date: '2026-03-28', weight: 12.5, change: -0.1, note: '' },
  { date: '2026-03-27', weight: 12.6, change: -0.1, note: '' },
  { date: '2026-03-26', weight: 12.7, change: 0, note: '' },
  { date: '2026-03-25', weight: 12.7, change: -0.2, note: '' },
  { date: '2026-03-24', weight: 12.9, change: -0.1, note: '' },
]);

const goal = reactive({ targetWeight: 12.0, targetDate: null });

const weightChart = {
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [
    { name: '实际体重', type: 'line', data: [13.1, 13.0, 12.9, 12.7, 12.7, 12.6, 12.5], smooth: true },
    { name: '目标体重', type: 'line', data: [12.5, 12.4, 12.3, 12.2, 12.1, 12.0, 12.0], smooth: true, lineStyle: { type: 'dashed' } },
  ],
};

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '体重(kg)', dataIndex: 'weight', width: 100 },
  { title: '变化(kg)', slotName: 'change', width: 100 },
  { title: '备注', dataIndex: 'note' },
];

const handleSaveGoal = () => {};
</script>

<style scoped>
.weight-tracking-container { padding: 20px; }
</style>
