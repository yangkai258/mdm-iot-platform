<template>
  <div class="ota-gray-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>OTA灰度发布配置</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            新增灰度
          </a-button>
        </div>
      </template>
      
      <a-alert type="info" style="margin-bottom: 16px;">
        灰度发布允许您先将固件推送给小部分设备进行验证，确认无问题后再全量推送。
      </a-alert>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="record.type === 'percentage' ? 'blue' : 'green'">{{ record.type === 'percentage' ? '百分比' : '指定设备' }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
        <template #progress="{ record }">
          <a-progress v-if="record.status === 'deploying'" :percent="record.progress" status="success" />
          <span v-else>{{ record.progress }}%</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link v-if="record.status === 'draft'" @click="handleDeploy(record)">部署</a-link>
            <a-link v-if="record.status === 'deploying'" @click="handlePause(record)">暂停</a-link>
            <a-link v-if="record.status === 'paused'" @click="handleResume(record)">继续</a-link>
            <a-link v-if="record.status === 'draft'" status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑灰度' : '新增灰度'" :width="600" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="目标固件版本" required>
          <a-select v-model="form.targetVersion" placeholder="选择目标版本">
            <a-option value="v2.0.0">v2.0.0 (正式版)</a-option>
            <a-option value="v2.1.0">v2.1.0 (测试版)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="灰度类型">
          <a-radio-group v-model="form.type">
            <a-radio value="percentage">百分比</a-radio>
            <a-radio value="specific">指定设备</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.type === 'percentage'" label="百分比">
          <a-slider v-model="form.percentage" :marks="{0:'0%',25:'25%',50:'50%',75:'75%',100:'100%'}" :max="100" />
        </a-form-item>
        <a-form-item v-if="form.type === 'specific'" label="选择设备">
          <a-select v-model="form.deviceIds" multiple placeholder="选择设备">
            <a-option value="DEV001">DEV001</a-option>
            <a-option value="DEV002">DEV002</a-option>
            <a-option value="DEV003">DEV003</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="自动暂停条件">
          <a-space vertical>
            <a-checkbox v-model="form.pauseOnFailure">失败率超过</a-checkbox>
            <a-input-number v-model="form.failureThreshold" :min="1" :max="100" suffix="%" style="width: 120px;" />
          </a-space>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);

const data = ref([
  { id: 1, name: 'v2.0.0灰度测试', targetVersion: 'v2.0.0', type: 'percentage', percentage: 10, deviceIds: [], status: 'completed', statusText: '已完成', progress: 100, totalDevices: 50, successCount: 48, failCount: 2, createdAt: '2026-03-28' },
  { id: 2, name: 'v2.1.0内测', targetVersion: 'v2.1.0', type: 'specific', percentage: 0, deviceIds: ['DEV001', 'DEV002', 'DEV003'], status: 'deploying', statusText: '部署中', progress: 66, totalDevices: 3, successCount: 2, failCount: 0, createdAt: '2026-03-28' },
  { id: 3, name: '50%灰度发布', targetVersion: 'v2.0.0', type: 'percentage', percentage: 50, deviceIds: [], status: 'draft', statusText: '草稿', progress: 0, totalDevices: 0, successCount: 0, failCount: 0, createdAt: '2026-03-27' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 3 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({ name: '', targetVersion: '', type: 'percentage', percentage: 10, deviceIds: [], pauseOnFailure: true, failureThreshold: 5, note: '' });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: '目标版本', dataIndex: 'targetVersion', width: 120 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '进度', slotName: 'progress', width: 150 },
  { title: '成功/总数', dataIndex: 'successCount', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' },
];

const getStatusColor = (s: string) => ({ draft: 'gray', deploying: 'blue', paused: 'orange', completed: 'green', cancelled: 'red' }[s] || 'default');

const handleCreate = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleView = (record: any) => {};
const handleDeploy = (record: any) => {};
const handlePause = (record: any) => {};
const handleResume = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.ota-gray-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
