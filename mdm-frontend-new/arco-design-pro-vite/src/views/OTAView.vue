<template>
  <div class="ota-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>OTA管理</span>
          <a-button type="primary" @click="handleUpload">
            <template #icon><icon-upload /></template>
            上传固件
          </a-button>
        </div>
      </template>
      
      <a-tabs default-active-key="packages">
        <a-tab-pane key="packages" title="固件包">
          <a-table :columns="pkgColumns" :data="packages" :loading="loading" :pagination="pagination">
            <template #version="{ record }">
              <span style="font-weight: 500;">{{ record.version }}</span>
            </template>
            <template #size="{ record }">
              {{ (record.size / 1024 / 1024).toFixed(2) }} MB
            </template>
            <template #status="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'default'">
                {{ record.status === 'active' ? '已启用' : '未启用' }}
              </a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleDeploy(record)">部署</a-link>
                <a-link @click="handleEdit(record)">编辑</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="deployments" title="部署记录">
          <a-table :columns="deployColumns" :data="deployments" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.status === 'success' ? 'green' : record.status === 'failed' ? 'red' : 'orange'">
                {{ record.statusText }}
              </a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const packages = ref([
  { id: 'F001', version: 'v2.0.0', deviceType: 'M5Stack', size: 2097152, checksum: 'abc123...', status: 'active', uploadedAt: '2026-03-20 10:00:00' },
  { id: 'F002', version: 'v1.1.0', deviceType: 'M5Stack', size: 1887436, checksum: 'def456...', status: 'inactive', uploadedAt: '2026-03-15 14:00:00' },
  { id: 'F003', version: 'v1.0.0', deviceType: 'ESP32', size: 1572864, checksum: 'ghi789...', status: 'inactive', uploadedAt: '2026-02-28 09:00:00' },
]);

const deployments = ref([
  { id: 'D001', version: 'v2.0.0', targetCount: 10, successCount: 8, failedCount: 2, status: 'success', statusText: '已完成', createdAt: '2026-03-28 08:00:00' },
  { id: 'D002', version: 'v1.1.0', targetCount: 5, successCount: 3, failedCount: 0, status: 'running', statusText: '进行中', createdAt: '2026-03-27 15:00:00' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });

const pkgColumns = [
  { title: '固件ID', dataIndex: 'id', width: 80 },
  { title: '版本', slotName: 'version', width: 100 },
  { title: '设备类型', dataIndex: 'deviceType', width: 100 },
  { title: '大小', slotName: 'size', width: 100 },
  { title: '校验码', dataIndex: 'checksum', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '上传时间', dataIndex: 'uploadedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const deployColumns = [
  { title: '部署ID', dataIndex: 'id', width: 80 },
  { title: '目标版本', dataIndex: 'version', width: 100 },
  { title: '目标设备', dataIndex: 'targetCount', width: 80 },
  { title: '成功', dataIndex: 'successCount', width: 80 },
  { title: '失败', dataIndex: 'failedCount', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'createdAt', width: 160 },
];

const handleUpload = () => {};
const handleDeploy = (record: any) => {};
const handleEdit = (record: any) => {};
</script>

<style scoped>
.ota-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
