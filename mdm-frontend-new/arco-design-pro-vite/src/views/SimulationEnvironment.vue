<template>
  <div class="simulation-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="仿真场景" :value="5" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="运行中仿真" :value="2" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="总运行次数" :value="128" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="通过率" :value="94" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>虚拟宠物仿真</span>
          <a-button type="primary" @click="handleCreateScene">
            <template #icon><icon-plus /></template>
            创建场景
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="scenes" title="仿真场景">
          <a-table :columns="sceneColumns" :data="scenes" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleRun(record)">运行</a-link>
                <a-link @click="handleView(record)">详情</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="results" title="仿真结果">
          <a-table :columns="resultColumns" :data="results" :pagination="pagination">
            <template #result="{ record }">
              <a-tag :color="record.passed ? 'green' : 'red'">{{ record.passed ? '通过' : '失败' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="reports" title="测试报告">
          <a-table :columns="reportColumns" :data="reports" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const scenes = ref([
  { id: 1, name: '日常互动仿真', type: 'behavior', typeText: '行为仿真', description: '模拟主人与宠物日常互动', status: 'draft', statusText: '草稿', runCount: 45, passRate: 96 },
  { id: 2, name: '情绪响应仿真', type: 'emotion', typeText: '情绪仿真', description: '测试宠物情绪识别与响应', status: 'published', statusText: '已发布', runCount: 32, passRate: 89 },
  { id: 3, name: 'OTA升级仿真', type: 'ota', typeText: 'OTA仿真', description: '模拟固件升级流程', status: 'published', statusText: '已发布', runCount: 28, passRate: 100 },
]);

const results = ref([
  { id: 1, sceneName: '日常互动仿真', result: true, duration: '5m30s', passRate: 96, time: '2026-03-28 18:00:00' },
  { id: 2, sceneName: '情绪响应仿真', result: false, duration: '3m00s', passRate: 78, time: '2026-03-28 17:00:00' },
]);

const reports = ref([
  { id: 1, name: '日常互动仿真报告', scene: '日常互动仿真', time: '2026-03-28', passRate: 96, issues: 2 },
  { id: 2, name: '情绪响应仿真报告', scene: '情绪响应仿真', time: '2026-03-27', passRate: 89, issues: 5 },
]);

const sceneColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '场景名称', dataIndex: 'name', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '描述', dataIndex: 'description' },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '运行次数', dataIndex: 'runCount', width: 100 },
  { title: '通过率', dataIndex: 'passRate', width: 100 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const resultColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '场景', dataIndex: 'sceneName', width: 150 },
  { title: '结果', slotName: 'result', width: 100 },
  { title: '耗时', dataIndex: 'duration', width: 100 },
  { title: '通过率', dataIndex: 'passRate', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const reportColumns = [
  { title: '报告名称', dataIndex: 'name', width: 200 },
  { title: '场景', dataIndex: 'scene', width: 150 },
  { title: '时间', dataIndex: 'time', width: 120 },
  { title: '通过率', dataIndex: 'passRate', width: 100 },
  { title: '问题数', dataIndex: 'issues', width: 100 },
];

const getStatusColor = (s: string) => ({ draft: 'gray', published: 'green', running: 'blue' }[s] || 'default');

const handleCreateScene = () => {};
const handleRun = (record: any) => {};
const handleView = (record: any) => {};
</script>

<style scoped>
.simulation-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
