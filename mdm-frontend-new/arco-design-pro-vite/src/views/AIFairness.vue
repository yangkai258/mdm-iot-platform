<template>
  <div class="ai-fairness-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="测试用例" :value="128" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="公平性评分" :value="92" suffix="/100" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="发现问题" :value="3" status="warning" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>AI伦理与公平性测试</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="tests" title="测试用例">
          <a-table :columns="testColumns" :data="tests" :pagination="pagination">
            <template #result="{ record }">
              <a-tag :color="record.passed ? 'green' : 'red'">{{ record.passed ? '通过' : '失败' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="bias" title="偏见检测">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-card title="性别偏见检测">
                <a-chart :option="biasChart" style="height: 250px;" />
              </a-card>
            </a-col>
            <a-col :span="12">
              <a-card title="年龄偏见检测">
                <a-chart :option="ageBiasChart" style="height: 250px;" />
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="explain" title="决策可解释性">
          <a-table :columns="explainColumns" :data="decisions" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const tests = ref([
  { id: 1, name: '性别偏见测试', category: '公平性', passed: true, score: 95, time: '2026-03-28' },
  { id: 2, name: '年龄歧视测试', category: '公平性', passed: true, score: 88, time: '2026-03-28' },
  { id: 3, name: '种族偏见测试', category: '伦理', passed: false, score: 72, time: '2026-03-27' },
]);

const decisions = ref([
  { id: 1, input: '用户:我不开心', decision: '安慰+陪伴', confidence: 92, factors: ['情绪词:不开心', '历史情绪:低落'], time: '2026-03-28 18:00:00' },
  { id: 2, input: '用户:陪我玩', decision: '开始游戏', confidence: 98, factors: ['意图词:玩', '时间:下午'], time: '2026-03-28 17:00:00' },
]);

const biasChart = {
  xAxis: { type: 'category', data: ['男性', '女性', '其他'] },
  yAxis: { type: 'value', max: 100 },
  series: [{ type: 'bar', data: [95, 92, 88] }],
};

const ageBiasChart = {
  xAxis: { type: 'category', data: ['儿童', '青年', '中年', '老年'] },
  yAxis: { type: 'value', max: 100 },
  series: [{ type: 'bar', data: [85, 95, 90, 80] }],
};

const testColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '测试名称', dataIndex: 'name', width: 200 },
  { title: '类别', dataIndex: 'category', width: 100 },
  { title: '结果', slotName: 'result', width: 100 },
  { title: '评分', dataIndex: 'score', width: 100 },
  { title: '时间', dataIndex: 'time', width: 120 },
];

const explainColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '输入', dataIndex: 'input', width: 200 },
  { title: 'AI决策', dataIndex: 'decision', width: 150 },
  { title: '置信度', dataIndex: 'confidence', width: 100 },
  { title: '决策因素', dataIndex: 'factors', width: 250 },
  { title: '时间', dataIndex: 'time', width: 160 },
];
</script>

<style scoped>
.ai-fairness-container { padding: 20px; }
</style>
