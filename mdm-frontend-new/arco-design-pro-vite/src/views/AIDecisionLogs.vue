<template>
  <div class="ai-decision-container">
    <a-card>
      <template #title>
        <span>AI决策日志</span>
      </template>
      
      <div class="search-area">
        <a-row :gutter="16">
          <a-col :span="5">
            <a-input v-model="searchForm.keyword" placeholder="搜索设备/内容" allow-clear />
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.decisionType" placeholder="决策类型" allow-clear>
              <a-option value="intent">意图识别</a-option>
              <a-option value="emotion">情绪响应</a-option>
              <a-option value="action">动作推荐</a-option>
              <a-option value="response">对话生成</a-option>
            </a-select>
          </a-col>
          <a-col :span="4">
            <a-select v-model="searchForm.confidence" placeholder="置信度" allow-clear>
              <a-option value="high">高(>90%)</a-option>
              <a-option value="medium">中(70-90%)</a-option>
              <a-option value="low">低(<70%)</a-option>
            </a-select>
          </a-col>
          <a-col :span="2">
            <a-button type="primary" @click="handleSearch">筛选</a-button>
          </a-col>
        </a-row>
      </div>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
        </template>
        <template #confidence="{ record }">
          <a-progress :percent="record.confidence" :color="getConfidenceColor(record.confidence)" size="small" />
        </template>
        <template #actions="{ record }">
          <a-link @click="handleViewDetail(record)">详情</a-link>
        </template>
      </a-table>
    </a-card>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="决策详情" :width="700">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="决策ID">{{ currentLog.id }}</a-descriptions-item>
        <a-descriptions-item label="决策类型">
          <a-tag>{{ currentLog.typeText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="设备ID">{{ currentLog.deviceId }}</a-descriptions-item>
        <a-descriptions-item label="置信度">{{ currentLog.confidence }}%</a-descriptions-item>
        <a-descriptions-item label="时间">{{ currentLog.time }}</a-descriptions-item>
        <a-descriptions-item label="耗时">{{ currentLog.duration }}ms</a-descriptions-item>
        <a-descriptions-item label="输入" :span="2">{{ currentLog.input }}</a-descriptions-item>
        <a-descriptions-item label="AI决策" :span="2">{{ currentLog.decision }}</a-descriptions-item>
        <a-descriptions-item label="执行结果" :span="2">{{ currentLog.result }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const searchForm = reactive({ keyword: '', decisionType: '', confidence: '' });

const data = ref([
  { id: 'DL001', deviceId: 'DEV001', type: 'intent', typeText: '意图识别', input: '今天天气怎么样', decision: '查询天气', confidence: 95, duration: 120, time: '2026-03-28 18:50:00', result: '成功' },
  { id: 'DL002', deviceId: 'DEV001', type: 'emotion', typeText: '情绪响应', input: '我不开心', decision: '安慰+陪伴', confidence: 88, duration: 80, time: '2026-03-28 18:45:00', result: '成功' },
  { id: 'DL003', deviceId: 'DEV002', type: 'action', typeText: '动作推荐', input: '用户心情低落', decision: '播放舒缓音乐', confidence: 72, duration: 50, time: '2026-03-28 18:40:00', result: '成功' },
  { id: 'DL004', deviceId: 'DEV003', type: 'response', typeText: '对话生成', input: '陪我聊天', decision: '好的呀', confidence: 92, duration: 200, time: '2026-03-28 18:35:00', result: '成功' },
  { id: 'DL005', deviceId: 'DEV001', type: 'intent', typeText: '意图识别', input: '啊啊啊', decision: '无法识别', confidence: 35, duration: 100, time: '2026-03-28 18:30:00', result: '无法响应' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 5 });
const detailVisible = ref(false);
const currentLog = ref<any>({});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 100 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '输入', dataIndex: 'input', width: 200 },
  { title: 'AI决策', dataIndex: 'decision', width: 150 },
  { title: '置信度', slotName: 'confidence', width: 120 },
  { title: '耗时', dataIndex: 'duration', width: 80 },
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '操作', slotName: 'actions', width: 80 },
];

const getTypeColor = (t: string) => ({ intent: 'blue', emotion: 'purple', action: 'orange', response: 'green' }[t] || 'default');
const getConfidenceColor = (c: number) => c >= 90 ? '#00B42A' : c >= 70 ? '#FF7D00' : '#F53F3F';

const handleSearch = () => {};
const handleViewDetail = (record: any) => { currentLog.value = record; detailVisible.value = true; };
</script>

<style scoped>
.ai-decision-container { padding: 20px; }
.search-area { margin-bottom: 16px; padding: 16px; background: #f7f8fa; border-radius: 4px; }
</style>
