<template>
  <div class="behavior-analysis-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="分析模型" :value="5" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="分析记录" :value="12580" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="异常检测" :value="28" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="预测准确率" :value="94" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>宠物行为分析</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="overview" title="分析概览">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="行为分布">
                <a-progress :percent="35" label="进食行为" />
                <a-progress :percent="25" label="睡眠行为" />
                <a-progress :percent="20" label="运动行为" />
                <a-progress :percent="15" label="社交行为" />
                <a-progress :percent="5" label="异常行为" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="活跃时段">
                <a-timeline>
                  <a-timeline-item>06:00-08:00 活跃高峰</a-timeline-item>
                  <a-timeline-item>12:00-14:00 午休时段</a-timeline-item>
                  <a-timeline-item>18:00-20:00 活跃高峰</a-timeline-item>
                  <a-timeline-item>22:00-06:00 睡眠时段</a-timeline-item>
                </a-timeline>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="records" title="分析记录">
          <a-table :columns="columns" :data="records" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="models" title="分析模型">
          <a-table :columns="modelColumns" :data="models" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 12 });

const records = ref([
  { id: 1, petName: '小黄', behaviorType: 'eating', behaviorText: '进食行为', duration: '15分钟', calories: 150, anomalyScore: 0.1, time: '2026-03-28 08:00:00' },
  { id: 2, petName: '小黄', behaviorType: 'sleeping', behaviorText: '睡眠行为', duration: '2小时', quality: '良好', anomalyScore: 0.05, time: '2026-03-28 12:00:00' },
  { id: 3, petName: '小红', behaviorType: 'exercise', behaviorText: '运动行为', duration: '30分钟', calories: 200, anomalyScore: 0.8, time: '2026-03-28 18:00:00' },
]);

const models = ref([
  { id: 1, name: '进食分析模型', type: 'eating', accuracy: 96, trainedAt: '2026-03-20', enabled: true },
  { id: 2, name: '睡眠质量模型', type: 'sleep', accuracy: 94, trainedAt: '2026-03-15', enabled: true },
  { id: 3, name: '异常行为检测', type: 'anomaly', accuracy: 89, trainedAt: '2026-03-10', enabled: true },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '宠物', dataIndex: 'petName', width: 100 },
  { title: '行为类型', dataIndex: 'behaviorText', width: 120 },
  { title: '持续时间', dataIndex: 'duration', width: 120 },
  { title: '异常分数', dataIndex: 'anomalyScore', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const modelColumns = [
  { title: '模型名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '准确率', dataIndex: 'accuracy', width: 100 },
  { title: '训练时间', dataIndex: 'trainedAt', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
];

const getAnomalyColor = (score: number) => score > 0.7 ? 'red' : score > 0.4 ? 'orange' : 'green';
</script>

<style scoped>
.behavior-analysis-container { padding: 20px; }
</style>
