<template>
  <div class="ai-monitor-container">
    <!-- 统计卡片 -->
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card>
          <a-statistic title="今日对话数" :value="stats.todayConversations" :value-from="0" animation />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="平均响应时间" :value="stats.avgResponseTime" suffix="ms" :value-from="0" animation />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="意图识别准确率" :value="stats.intentAccuracy" suffix="%" :value-from="0" animation />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="TTS成功率" :value="stats.ttsSuccessRate" suffix="%" :value-from="0" animation />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <!-- 实时对话 -->
      <a-col :span="16">
        <a-card>
          <template #title>
            <span>实时对话监控</span>
            <a-space style="float: right;">
              <a-switch v-model="autoRefresh" />
              <span>自动刷新</span>
            </a-space>
          </template>
          <a-table :columns="conversationColumns" :data="conversations" :pagination="pagination" :loading="loading">
            <template #type="{ record }">
              <a-tag :color="record.type === 'user' ? 'blue' : 'green'">
                {{ record.type === 'user' ? '用户' : 'AI' }}
              </a-tag>
            </template>
            <template #sentiment="{ record }">
              <a-tag :color="getSentimentColor(record.sentiment)">{{ record.sentimentText }}</a-tag>
            </template>
            <template #intent="{ record }">
              <a-tooltip :content="record.intent">{{ record.intent }}</a-tooltip>
            </template>
          </a-table>
        </a-card>
      </a-col>

      <!-- 质量趋势 -->
      <a-col :span="8">
        <a-card title="AI质量指标">
          <a-chart :option="qualityChart" style="height: 200px;" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px;">
      <!-- 异常告警 -->
      <a-col :span="12">
        <a-card title="AI异常事件">
          <a-table :columns="anomalyColumns" :data="anomalies" :pagination="paginationSmall">
            <template #level="{ record }">
              <a-tag :color="getAnomalyLevelColor(record.level)">{{ record.levelText }}</a-tag>
            </template>
          </a-table>
        </a-card>
      </a-col>

      <!-- 模型调用统计 -->
      <a-col :span="12">
        <a-card title="模型调用分布">
          <a-table :columns="modelColumns" :data="modelStats" :pagination="false">
            <template #rate="{ record }">
              <a-progress :percent="record.rate" :color="record.color" />
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- AI沙箱测试入口 -->
    <a-card title="AI沙箱测试" style="margin-top: 16px;">
      <a-space>
        <a-button type="primary" @click="handleSandboxTest">新建测试</a-button>
        <a-button @click="handleBatchTest">批量测试</a-button>
        <a-button @click="handleRegressionTest">回归测试</a-button>
      </a-space>
      <a-tabs style="margin-top: 16px;">
        <a-tab-pane key="results" title="测试结果">
          <a-table :columns="testResultColumns" :data="testResults" :pagination="paginationSmall">
            <template #status="{ record }">
              <a-tag :color="record.passed ? 'green' : 'red'">
                {{ record.passed ? '通过' : '失败' }}
              </a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="logs" title="测试日志">
          <a-log-list :data="testLogs" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const autoRefresh = ref(true);

const stats = reactive({
  todayConversations: 1583,
  avgResponseTime: 245,
  intentAccuracy: 94.5,
  ttsSuccessRate: 99.2,
});

const conversations = ref([
  { id: 'C001', deviceId: 'DEV001', type: 'user', content: '今天天气怎么样', intent: '查询天气', sentiment: 'neutral', sentimentText: '中性', responseTime: 230, time: '18:50:12' },
  { id: 'C002', deviceId: 'DEV001', type: 'ai', content: '今天天气晴朗，气温20-28度', intent: '查询天气', sentiment: 'positive', sentimentText: '正面', responseTime: 180, time: '18:50:13' },
  { id: 'C003', deviceId: 'DEV002', type: 'user', content: '陪我玩', intent: '玩耍', sentiment: 'positive', sentimentText: '正面', responseTime: 210, time: '18:49:55' },
  { id: 'C004', deviceId: 'DEV002', type: 'ai', content: '好的，我们来玩飞盘吧！', intent: '玩耍', sentiment: 'happy', sentimentText: '开心', responseTime: 195, time: '18:49:56' },
  { id: 'C005', deviceId: 'DEV003', type: 'user', content: '我不开心', intent: '情感陪伴', sentiment: 'negative', sentimentText: '负面', responseTime: 280, time: '18:48:30' },
]);

const anomalies = ref([
  { id: 'A001', type: 'high_latency', typeText: '高延迟', level: 'warning', levelText: '警告', deviceId: 'DEV005', message: '响应时间超过500ms', time: '18:45:00' },
  { id: 'A002', type: 'intent_fail', typeText: '意图识别失败', level: 'error', levelText: '错误', deviceId: 'DEV008', message: '无法识别用户意图', time: '18:40:15' },
  { id: 'A003', type: 'tts_fail', typeText: 'TTS失败', level: 'warning', levelText: '警告', deviceId: 'DEV012', message: 'TTS服务超时', time: '18:35:22' },
]);

const modelStats = ref([
  { model: 'MiniMax-TTS', calls: 15000, rate: 45, color: '#165DFF' },
  { model: 'Pet-NLU', calls: 12000, rate: 36, color: '#00B42A' },
  { model: 'Pet-LLM', calls: 5800, rate: 17, color: '#FF7D00' },
  { model: 'MiniMax-ASR', calls: 500, rate: 2, color: '#722ED1' },
]);

const testResults = ref([
  { id: 'T001', name: '情感陪伴测试', passed: true, duration: '5m30s', passRate: 98, time: '18:30:00' },
  { id: 'T002', name: '意图识别测试', passed: true, duration: '3m15s', passRate: 94, time: '18:25:00' },
  { id: 'T003', name: 'TTS质量测试', passed: false, duration: '2m00s', passRate: 78, time: '18:20:00' },
]);

const testLogs = ref([
  { time: '18:30:00', level: 'info', message: '开始情感陪伴测试' },
  { time: '18:30:05', level: 'info', message: '测试用例1: 开心情绪-通过' },
  { time: '18:30:10', level: 'info', message: '测试用例2: 悲伤情绪-通过' },
  { time: '18:35:00', level: 'warn', message: '发现2个非致命问题' },
  { time: '18:35:30', level: 'info', message: '测试完成，通过率98%' },
]);

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });
const paginationSmall = reactive({ current: 1, pageSize: 5, total: 3 });

const conversationColumns = [
  { title: '时间', dataIndex: 'time', width: 90 },
  { title: '设备', dataIndex: 'deviceId', width: 100 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '内容', dataIndex: 'content' },
  { title: '意图', slotName: 'intent', width: 100 },
  { title: '情绪', slotName: 'sentiment', width: 80 },
  { title: '响应(ms)', dataIndex: 'responseTime', width: 100 },
];

const anomalyColumns = [
  { title: '时间', dataIndex: 'time', width: 90 },
  { title: '类型', dataIndex: 'typeText', width: 120 },
  { title: '级别', slotName: 'level', width: 80 },
  { title: '设备', dataIndex: 'deviceId', width: 100 },
  { title: '消息', dataIndex: 'message' },
];

const modelColumns = [
  { title: '模型', dataIndex: 'model', width: 150 },
  { title: '调用次数', dataIndex: 'calls', width: 100 },
  { title: '占比', slotName: 'rate' },
];

const testResultColumns = [
  { title: '测试名称', dataIndex: 'name', width: 150 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 80 },
  { title: '通过率', dataIndex: 'passRate', width: 80 },
  { title: '时间', dataIndex: 'time', width: 100 },
];

const qualityChart = {
  xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
  yAxis: { type: 'value', min: 0, max: 100 },
  series: [
    { name: '意图识别准确率', type: 'line', data: [92, 94, 95, 93, 96, 94, 95], smooth: true },
    { name: 'TTS成功率', type: 'line', data: [99, 99, 98, 99, 99, 99, 99], smooth: true },
  ],
};

const getSentimentColor = (s: string) => ({ positive: 'green', negative: 'red', neutral: 'gray', happy: 'blue' }[s] || 'default');
const getAnomalyLevelColor = (l: string) => ({ warning: 'orange', error: 'red' }[l] || 'default');

const handleSandboxTest = () => {};
const handleBatchTest = () => {};
const handleRegressionTest = () => {};
</script>

<style scoped>
.ai-monitor-container { padding: 20px; }
</style>
