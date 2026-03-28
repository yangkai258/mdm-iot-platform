<template>
  <div class="cloud-backup-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="备份任务" :value="36" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="备份大小" :value="256" suffix="GB" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="恢复次数" :value="12" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="备份成功率" :value="99" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>云端备份管理</span>
          <a-space>
            <a-button @click="handleBackupNow">立即备份</a-button>
            <a-button type="primary" @click="handleRestore">恢复</a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="backups" title="备份记录">
          <a-table :columns="backupColumns" :data="backups" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.status === 'completed' ? 'green' : 'blue'">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="schedule" title="备份计划">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="自动备份">
              <a-switch v-model="settings.autoBackup" />
            </a-form-item>
            <a-form-item label="备份频率">
              <a-select v-model="settings.frequency">
                <a-option value="daily">每天</a-option>
                <a-option value="weekly">每周</a-option>
                <a-option value="monthly">每月</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="保留份数">
              <a-input-number v-model="settings.retention" :min="1" />
            </a-form-item>
            <a-form-item label="备份目标">
              <a-select v-model="settings.target">
                <a-option value="oss">阿里云OSS</a-option>
                <a-option value="s3">AWS S3</a-option>
                <a-option value="local">本地存储</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveSettings">保存设置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const settings = reactive({ autoBackup: true, frequency: 'daily', retention: 7, target: 'oss' });

const backups = ref([
  { id: 1, name: '全量备份-20260328', size: '5.2GB', status: 'completed', statusText: '完成', time: '2026-03-28 02:00:00' },
  { id: 2, name: '增量备份-20260327', size: '512MB', status: 'completed', statusText: '完成', time: '2026-03-27 02:00:00' },
]);

const backupColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '备份名称', dataIndex: 'name', width: 200 },
  { title: '大小', dataIndex: 'size', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '备份时间', dataIndex: 'time', width: 160 },
];

const handleBackupNow = () => {};
const handleRestore = () => {};
const handleSaveSettings = () => {};
</script>

<style scoped>
.cloud-backup-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
