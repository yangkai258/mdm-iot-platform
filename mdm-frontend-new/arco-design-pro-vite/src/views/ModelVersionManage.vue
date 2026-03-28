<template>
  <div class="model-version-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>模型版本管理</span>
          <a-space>
            <a-button @click="handleRollback">回滚</a-button>
            <a-button type="primary" @click="handleUpload">
              <template #icon><icon-upload /></template>
              上传新版本
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="tts" title="TTS模型">
          <a-table :columns="columns" :data="ttsVersions" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link v-if="record.status !== 'active'" @click="handleActivate(record)">激活</a-link>
                <a-link @click="handleViewLogs(record)">日志</a-link>
                <a-link v-if="record.status === 'active'" status="danger" @click="handleDelete(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="asr" title="ASR模型">
          <a-table :columns="columns" :data="asrVersions" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="nlu" title="NLU模型">
          <a-table :columns="columns" :data="nluVersions" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="llm" title="LLM模型">
          <a-table :columns="columns" :data="llmVersions" :pagination="pagination" />
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 回滚确认 -->
    <a-modal v-model:visible="rollbackVisible" title="确认回滚" @before-ok="handleConfirmRollback">
      <a-alert type="warning">回滚操作可能导致服务短暂中断，当前在线用户会话可能需要重新建立。</a-alert>
      <a-divider />
      <a-form :model="rollbackForm" layout="vertical">
        <a-form-item label="选择目标版本">
          <a-select v-model="rollbackForm.targetVersion" placeholder="选择要回滚到的版本">
            <a-option value="v2.0.0">v2.0.0 (推荐)</a-option>
            <a-option value="v1.5.0">v1.5.0</a-option>
            <a-option value="v1.0.0">v1.0.0</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="回滚原因">
          <a-textarea v-model="rollbackForm.reason" placeholder="请输入回滚原因" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 日志查看 -->
    <a-modal v-model:visible="logsVisible" title="版本日志" :width="800">
      <a-timeline>
        <a-timeline-item v-for="log in versionLogs" :key="log.time" :color="log.type === 'success' ? 'green' : log.type === 'error' ? 'red' : 'gray'">
          <b>{{ log.time }}</b> - {{ log.message }}
        </a-timeline-item>
      </a-timeline>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 3 });
const rollbackVisible = ref(false);
const logsVisible = ref(false);

const rollbackForm = reactive({ targetVersion: '', reason: '' });

const ttsVersions = ref([
  { id: 1, version: 'v2.1.0', size: '256MB', status: 'active', statusText: '当前使用', deployedAt: '2026-03-20', accuracy: 98.5, latency: 120 },
  { id: 2, version: 'v2.0.0', size: '240MB', status: 'inactive', statusText: '历史版本', deployedAt: '2026-03-01', accuracy: 97.8, latency: 135 },
  { id: 3, version: 'v1.0.0', size: '200MB', status: 'inactive', statusText: '历史版本', deployedAt: '2026-01-15', accuracy: 95.0, latency: 200 },
]);

const asrVersions = ref([
  { id: 1, version: 'v1.5.0', size: '180MB', status: 'active', statusText: '当前使用', deployedAt: '2026-03-15', accuracy: 97.2, latency: 80 },
]);

const nluVersions = ref([
  { id: 1, version: 'v3.0.0', size: '520MB', status: 'testing', statusText: '测试中', deployedAt: '2026-03-25', accuracy: 94.5, latency: 150 },
  { id: 2, version: 'v2.5.0', size: '480MB', status: 'active', statusText: '当前使用', deployedAt: '2026-02-20', accuracy: 92.0, latency: 120 },
]);

const llmVersions = ref([
  { id: 1, version: 'v1.0.0', size: '2GB', status: 'active', statusText: '当前使用', deployedAt: '2026-03-10', accuracy: 95.0, latency: 500 },
]);

const versionLogs = ref([
  { time: '2026-03-20 10:00:00', type: 'info', message: '开始部署 v2.1.0' },
  { time: '2026-03-20 10:00:05', type: 'success', message: '模型文件下载完成' },
  { time: '2026-03-20 10:00:30', type: 'success', message: '灰度发布10%流量' },
  { time: '2026-03-20 10:30:00', type: 'success', message: '全量发布完成' },
]);

const columns = [
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '大小', dataIndex: 'size', width: 100 },
  { title: '状态', slotName: 'status', width: 120 },
  { title: '准确率', dataIndex: 'accuracy', width: 100 },
  { title: '延迟(ms)', dataIndex: 'latency', width: 100 },
  { title: '发布时间', dataIndex: 'deployedAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const getStatusColor = (s: string) => ({ active: 'green', inactive: 'gray', testing: 'blue' }[s] || 'default');

const handleUpload = () => {};
const handleRollback = () => { rollbackVisible.value = true; };
const handleActivate = (record: any) => {};
const handleViewLogs = (record: any) => { logsVisible.value = true; };
const handleDelete = (record: any) => {};
const handleConfirmRollback = (done: boolean) => { done(true); rollbackVisible.value = false; };
</script>

<style scoped>
.model-version-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
