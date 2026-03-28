<template>
  <div class="storage-management-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="6">
        <a-card><a-statistic title="总存储" :value="500" suffix="GB" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已使用" :value="256" suffix="GB" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="可用" :value="244" suffix="GB" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-progress :percent="51" label="使用率" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>存储管理</span>
          <a-space>
            <a-button @click="handleCleanup">清理缓存</a-button>
            <a-button type="primary" @click="handleExpand">扩容存储</a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="usage" title="存储使用">
          <a-table :columns="usageColumns" :data="usageData" :pagination="false">
            <template #size="{ record }">
              <span>{{ record.size }}GB</span>
            </template>
            <template #percent="{ record }">
              <a-progress :percent="record.percent" :color="getColor(record.percent)" />
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="files" title="文件管理">
          <a-table :columns="fileColumns" :data="files" :pagination="pagination">
            <template #type="{ record }">
              <a-tag>{{ record.typeText }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleDownload(record)">下载</a-link>
                <a-link status="danger" @click="handleDelete(record)">删除</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="backup" title="备份管理">
          <a-table :columns="backupColumns" :data="backups" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.status === 'completed' ? 'green' : 'blue'">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });

const usageData = ref([
  { id: 1, category: 'media', categoryText: '媒体文件', size: 120, percent: 47, count: 2560 },
  { id: 2, category: 'logs', categoryText: '系统日志', size: 80, percent: 31, count: 12500 },
  { id: 3, category: 'data', categoryText: '业务数据', size: 45, percent: 18, count: 856 },
  { id: 4, category: 'cache', categoryText: '缓存文件', size: 11, percent: 4, count: 3200 },
]);

const files = ref([
  { id: 1, name: 'video_20260328.mp4', type: 'media', typeText: '视频', size: '256MB', createdAt: '2026-03-28 10:00:00' },
  { id: 2, name: 'photo_album.zip', type: 'media', typeText: '图片', size: '128MB', createdAt: '2026-03-27 15:00:00' },
  { id: 3, name: 'system_log_2026.log', type: 'logs', typeText: '日志', size: '50MB', createdAt: '2026-03-28 00:00:00' },
]);

const backups = ref([
  { id: 1, name: '全量备份-20260328', size: '2.5GB', status: 'completed', statusText: '已完成', time: '2026-03-28 02:00:00' },
  { id: 2, name: '增量备份-20260327', size: '512MB', status: 'completed', statusText: '已完成', time: '2026-03-27 02:00:00' },
]);

const usageColumns = [
  { title: '类别', dataIndex: 'categoryText', width: 150 },
  { title: '大小', slotName: 'size', width: 120 },
  { title: '占比', slotName: 'percent' },
  { title: '文件数', dataIndex: 'count', width: 100 },
];

const fileColumns = [
  { title: '文件名', dataIndex: 'name', width: 250 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '大小', dataIndex: 'size', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const backupColumns = [
  { title: '备份名称', dataIndex: 'name', width: 200 },
  { title: '大小', dataIndex: 'size', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const getColor = (p: number) => p > 80 ? 'red' : p > 50 ? 'orange' : 'green';

const handleCleanup = () => {};
const handleExpand = () => {};
const handleDownload = (record: any) => {};
const handleDelete = (record: any) => {};
</script>

<style scoped>
.storage-management-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
