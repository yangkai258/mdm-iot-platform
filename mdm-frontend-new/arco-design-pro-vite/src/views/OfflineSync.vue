<template>
  <div class="offline-sync-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="离线缓存" :value="128" suffix="条" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="待同步" :value="12" status="warning" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="同步成功率" :value="98" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>离线缓存与同步</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="cache" title="离线缓存">
          <a-table :columns="cacheColumns" :data="caches" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.synced ? 'green' : 'orange'">{{ record.synced ? '已同步' : '待同步' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="sync" title="同步队列">
          <a-table :columns="syncColumns" :data="queue" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="settings" title="同步设置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="自动同步">
              <a-switch v-model="settings.autoSync" />
            </a-form-item>
            <a-form-item label="同步间隔">
              <a-select v-model="settings.syncInterval">
                <a-option value="1min">1分钟</a-option>
                <a-option value="5min">5分钟</a-option>
                <a-option value="15min">15分钟</a-option>
                <a-option value="manual">手动</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="WiFi下仅同步">
              <a-switch v-model="settings.wifiOnly" />
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSave">保存</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 5 });

const caches = ref([
  { id: 1, deviceId: 'DEV001', type: 'status', content: '设备状态快照', size: '2KB', synced: true, lastSync: '2026-03-28 18:00:00' },
  { id: 2, deviceId: 'DEV001', type: 'health', content: '健康数据', size: '5KB', synced: false, lastSync: null },
]);

const queue = ref([
  { id: 1, deviceId: 'DEV001', type: 'health', createdAt: '2026-03-28 17:00:00', retryCount: 0 },
  { id: 2, deviceId: 'DEV002', type: 'status', createdAt: '2026-03-28 16:30:00', retryCount: 1 },
]);

const settings = reactive({ autoSync: true, syncInterval: '5min', wifiOnly: true });

const cacheColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '内容', dataIndex: 'content', width: 200 },
  { title: '大小', dataIndex: 'size', width: 80 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '最后同步', dataIndex: 'lastSync', width: 160 },
];

const syncColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
  { title: '重试次数', dataIndex: 'retryCount', width: 100 },
];

const handleSave = () => {};
</script>

<style scoped>
.offline-sync-container { padding: 20px; }
</style>
