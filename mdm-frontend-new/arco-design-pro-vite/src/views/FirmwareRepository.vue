<template>
  <div class="firmware-repo-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="固件版本" :value="24" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="总下载量" :value="12850" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="设备覆盖率" :value="98" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>固件仓库</span>
          <a-button type="primary" @click="handleUpload">
            <template #icon><icon-upload /></template>
            上传固件
          </a-button>
        </div>
      </template>
      
      <a-table :columns="columns" :data="firmwares" :loading="loading" :pagination="pagination">
        <template #version="{ record }">
          <a-tag color="arcoblue">{{ record.version }}</a-tag>
        </template>
        <template #type="{ record }">
          <a-tag>{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.isLatest ? 'green' : 'gray'">{{ record.isLatest ? '最新' : '历史' }}</a-tag>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleView(record)">详情</a-link>
            <a-link @click="handleDownload(record)">下载</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:visible="uploadVisible" title="上传固件" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="固件版本" required>
          <a-input v-model="form.version" placeholder="如: v2.0.0" />
        </a-form-item>
        <a-form-item label="硬件型号">
          <a-select v-model="form.hardwareModel" placeholder="选择硬件型号">
            <a-option value="M5Stack-PetBot-v2">M5Stack-PetBot-v2</a-option>
            <a-option value="M5Stack-PetBot-v1">M5Stack-PetBot-v1</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="固件类型">
          <a-radio-group v-model="form.type">
            <a-radio value="release">正式版</a-radio>
            <a-radio value="beta">测试版</a-radio>
            <a-radio value="dev">开发版</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="固件文件">
          <a-upload action="#" :limit="1" />
        </a-form-item>
        <a-form-item label="更新日志">
          <a-textarea v-model="form.changelog" :rows="4" placeholder="请输入更新内容" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const pagination = reactive({ current: 1, pageSize: 10, total: 6 });
const uploadVisible = ref(false);

const form = reactive({ version: '', hardwareModel: '', type: 'release', changelog: '' });

const firmwares = ref([
  { id: 1, version: 'v2.0.0', hardwareModel: 'M5Stack-PetBot-v2', type: 'release', typeText: '正式版', size: '2.5MB', downloads: 5200, isLatest: true, createdAt: '2026-03-20' },
  { id: 2, version: 'v1.9.0', hardwareModel: 'M5Stack-PetBot-v2', type: 'release', typeText: '正式版', size: '2.4MB', downloads: 3800, isLatest: false, createdAt: '2026-03-10' },
  { id: 3, version: 'v2.1.0-beta', hardwareModel: 'M5Stack-PetBot-v2', type: 'beta', typeText: '测试版', size: '2.6MB', downloads: 1200, isLatest: false, createdAt: '2026-03-28' },
]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '版本', slotName: 'version', width: 120 },
  { title: '硬件型号', dataIndex: 'hardwareModel', width: 150 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '文件大小', dataIndex: 'size', width: 100 },
  { title: '下载量', dataIndex: 'downloads', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '上传时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 180, fixed: 'right' },
];

const handleUpload = () => { uploadVisible.value = true; };
const handleView = (record: any) => {};
const handleDownload = (record: any) => {};
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); uploadVisible.value = false; };
</script>

<style scoped>
.firmware-repo-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
