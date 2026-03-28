<template>
  <div class="device-diagnostics-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="诊断任务" :value="28" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="发现问题" :value="5" status="warning" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="修复率" :value="82" suffix="%" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="平均诊断时间" :value="30" suffix="秒" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备诊断中心</span>
          <a-space>
            <a-input-search placeholder="搜索设备ID" style="width: 200px;" />
            <a-button type="primary" @click="handleNewDiagnosis">
              <template #icon><icon-plus /></template>
              新建诊断
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="tasks" title="诊断任务">
          <a-table :columns="taskColumns" :data="tasks" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #result="{ record }">
              <a-tag v-if="record.issues > 0" color="red">{{ record.issues }}个问题</a-tag>
              <a-tag v-else color="green">正常</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="history" title="诊断历史">
          <a-table :columns="historyColumns" :data="history" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="templates" title="诊断模板">
          <a-row :gutter="16">
            <a-col :span="6" v-for="tpl in templates" :key="tpl.id">
              <a-card size="small" class="tpl-card">
                <div class="tpl-name">{{ tpl.name }}</div>
                <div class="tpl-desc">{{ tpl.description }}</div>
                <a-button type="primary" size="small" @click="handleUseTemplate(tpl)">使用</a-button>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-drawer v-model:visible="diagVisible" title="设备诊断" :width="600">
      <a-result v-if="diagResult" :status="diagResult.success ? 'success' : 'error'" :title="diagResult.success ? '诊断完成' : '发现异常'">
        <template #content>
          <a-descriptions :column="1">
            <a-descriptions-item label="设备ID">{{ diagResult.deviceId }}</a-descriptions-item>
            <a-descriptions-item label="网络状态">{{ diagResult.network }}</a-descriptions-item>
            <a-descriptions-item label="固件版本">{{ diagResult.firmware }}</a-descriptions-item>
            <a-descriptions-item label="电池电量">{{ diagResult.battery }}%</a-descriptions-item>
          </a-descriptions>
          <div v-if="diagResult.issues && diagResult.issues.length > 0" style="margin-top: 16px;">
            <div v-for="issue in diagResult.issues" :key="issue.id" style="color: #F53F3F;">
              {{ issue.id }}. {{ issue.description }}
            </div>
          </div>
        </template>
      </a-result>
      <a-spin v-else />
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const diagVisible = ref(false);
const diagResult = ref<any>(null);

const tasks = ref([
  { id: 1, deviceId: 'DEV001', deviceName: '小黄-客厅', status: 'completed', statusText: '完成', issues: 2, completedAt: '2026-03-28 18:00:00' },
  { id: 2, deviceId: 'DEV002', deviceName: '小红-卧室', status: 'running', statusText: '进行中', issues: 0, completedAt: null },
]);

const history = ref([
  { id: 1, deviceId: 'DEV001', type: 'full', status: 'warning', issues: 2, duration: '25秒', time: '2026-03-28 18:00:00' },
  { id: 2, deviceId: 'DEV003', type: 'quick', status: 'success', issues: 0, duration: '10秒', time: '2026-03-28 17:00:00' },
]);

const templates = ref([
  { id: 1, name: '全面诊断', description: '检查网络/固件/电池/传感器' },
  { id: 2, name: '快速扫描', description: '检查网络连接状态' },
  { id: 3, name: '固件检查', description: '检查固件版本和兼容性' },
]);

const taskColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'deviceName', width: 150 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '诊断结果', slotName: 'result', width: 120 },
  { title: '完成时间', dataIndex: 'completedAt', width: 160 },
];

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '诊断类型', dataIndex: 'type', width: 100 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '问题数', dataIndex: 'issues', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const getStatusColor = (s: string) => ({ completed: 'green', running: 'blue', failed: 'red' }[s] || 'default');

const handleNewDiagnosis = () => { diagVisible.value = true; diagResult.value = null; setTimeout(() => { diagResult.value = { deviceId: 'DEV001', network: '良好', firmware: 'v2.0.0', battery: 85, success: true, issues: [] }; }, 2000); };
const handleUseTemplate = (tpl: any) => {};
</script>

<style scoped>
.device-diagnostics-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.tpl-card { text-align: center; margin-bottom: 12px; }
.tpl-name { font-weight: bold; margin-bottom: 8px; }
.tpl-desc { color: #86909c; font-size: 12px; margin-bottom: 8px; }
</style>
