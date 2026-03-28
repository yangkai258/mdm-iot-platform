<template>
  <div class="firmware-matrix-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>固件兼容性矩阵</span>
          <a-space>
            <a-button @click="handleExport">导出</a-button>
            <a-button type="primary" @click="handleAdd">添加规则</a-button>
          </a-space>
        </div>
      </template>
      
      <a-alert type="info" style="margin-bottom: 16px;">
        定义固件版本与设备型号的兼容性关系，控制OTA升级范围
      </a-alert>
      
      <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
        <template #compatible="{ record }">
          <a-tag :color="record.compatible ? 'green' : 'red'">
            {{ record.compatible ? '兼容' : '不兼容' }}
          </a-tag>
        </template>
        <template #autoUpgrade="{ record }">
          <a-switch :checked="record.autoUpgrade" disabled />
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-link @click="handleEdit(record)">编辑</a-link>
            <a-link status="danger" @click="handleDelete(record)">删除</a-link>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 兼容设备列表 -->
    <a-card title="设备兼容性详情" style="margin-top: 16px;">
      <a-tabs>
        <a-tab-pane key="compatible" title="兼容设备">
          <a-row :gutter="16">
            <a-col :span="6" v-for="device in compatibleDevices" :key="device.deviceId">
              <a-card size="small" class="device-card">
                <a-descriptions :column="1" size="small">
                  <a-descriptions-item label="设备ID">{{ device.deviceId }}</a-descriptions-item>
                  <a-descriptions-item label="当前固件">{{ device.currentFirmware }}</a-descriptions-item>
                  <a-descriptions-item label="状态">
                    <a-tag :color="device.online ? 'green' : 'red'">{{ device.online ? '在线' : '离线' }}</a-tag>
                  </a-descriptions-item>
                </a-descriptions>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        <a-tab-pane key="incompatible" title="不兼容设备">
          <a-row :gutter="16">
            <a-col :span="6" v-for="device in incompatibleDevices" :key="device.deviceId">
              <a-card size="small" class="device-card incompatible">
                <a-descriptions :column="1" size="small">
                  <a-descriptions-item label="设备ID">{{ device.deviceId }}</a-descriptions-item>
                  <a-descriptions-item label="当前固件">{{ device.currentFirmware }}</a-descriptions-item>
                  <a-descriptions-item label="原因">{{ device.reason }}</a-descriptions-item>
                </a-descriptions>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 添加/编辑弹窗 -->
    <a-modal v-model:visible="editVisible" :title="isEdit ? '编辑规则' : '添加规则'" @before-ok="handleSubmit">
      <a-form :model="form" layout="vertical">
        <a-form-item label="目标固件版本" required>
          <a-input v-model="form.targetVersion" placeholder="如: v2.0.0" />
        </a-form-item>
        <a-form-item label="设备型号" required>
          <a-select v-model="form.deviceModel" placeholder="选择设备型号" multiple>
            <a-option value="M5Stack-Core">M5Stack-Core</a-option>
            <a-option value="M5Stack-Fire">M5Stack-Fire</a-option>
            <a-option value="M5Stack-Atom">M5Stack-Atom</a-option>
            <a-option value="ESP32-DevKit">ESP32-DevKit</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="兼容的当前版本">
          <a-select v-model="form.fromVersions" placeholder="选择可升级的源版本" multiple allow-create>
            <a-option value="v1.0.0">v1.0.0</a-option>
            <a-option value="v1.1.0">v1.1.0</a-option>
            <a-option value="v1.5.0">v1.5.0</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="兼容性">
          <a-radio-group v-model="form.compatible">
            <a-radio :value="true">兼容</a-radio>
            <a-radio :value="false">不兼容</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="自动升级">
          <a-switch v-model="form.autoUpgrade" />
          <span style="margin-left: 8px; color: #86909c;">允许自动升级到此版本</span>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="form.note" placeholder="可选备注" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 1, targetVersion: 'v2.0.0', deviceModel: 'M5Stack-Core', fromVersions: ['v1.0.0', 'v1.1.0'], compatible: true, autoUpgrade: true, note: '正式版', createdAt: '2026-03-01' },
  { id: 2, targetVersion: 'v2.0.0', deviceModel: 'M5Stack-Fire', fromVersions: ['v1.5.0'], compatible: true, autoUpgrade: false, note: 'Fire版本需手动确认', createdAt: '2026-03-01' },
  { id: 3, targetVersion: 'v2.0.0', deviceModel: 'M5Stack-Atom', fromVersions: [], compatible: false, autoUpgrade: false, note: 'Atom不支持此版本', createdAt: '2026-03-05' },
  { id: 4, targetVersion: 'v1.5.0', deviceModel: 'ESP32-DevKit', fromVersions: ['v1.0.0'], compatible: true, autoUpgrade: true, note: 'DevKit专用版本', createdAt: '2026-02-15' },
]);

const compatibleDevices = ref([
  { deviceId: 'DEV001', deviceName: '小黄-客厅', currentFirmware: 'v1.0.0', online: true },
  { deviceId: 'DEV002', deviceName: '小红-卧室', currentFirmware: 'v1.1.0', online: true },
  { deviceId: 'DEV003', deviceName: '小绿-书房', currentFirmware: 'v1.0.0', online: false },
]);

const incompatibleDevices = ref([
  { deviceId: 'DEV010', deviceName: '小Atom-厨房', currentFirmware: 'v1.0.0', reason: '内存不足' },
  { deviceId: 'DEV011', deviceName: '小Atom-阳台', currentFirmware: 'v1.1.0', reason: '存储空间不足' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const editVisible = ref(false);
const isEdit = ref(false);

const form = reactive({
  targetVersion: '',
  deviceModel: [],
  fromVersions: [],
  compatible: true,
  autoUpgrade: false,
  note: '',
});

const columns = [
  { title: '规则ID', dataIndex: 'id', width: 80 },
  { title: '目标版本', dataIndex: 'targetVersion', width: 120 },
  { title: '设备型号', dataIndex: 'deviceModel', width: 150 },
  { title: '源版本', dataIndex: 'fromVersions', width: 200 },
  { title: '兼容', slotName: 'compatible', width: 100 },
  { title: '自动升级', slotName: 'autoUpgrade', width: 100 },
  { title: '备注', dataIndex: 'note' },
  { title: '创建时间', dataIndex: 'createdAt', width: 120 },
  { title: '操作', slotName: 'actions', width: 120, fixed: 'right' },
];

const handleExport = () => {};
const handleAdd = () => { isEdit.value = false; editVisible.value = true; };
const handleEdit = (record: any) => { isEdit.value = true; Object.assign(form, record); editVisible.value = true; };
const handleDelete = (record: any) => {};
const handleSubmit = (done: boolean) => { done(true); editVisible.value = false; };
</script>

<style scoped>
.firmware-matrix-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.device-card { margin-bottom: 8px; }
.device-card.incompatible { opacity: 0.7; border: 1px solid #F53F3F; }
</style>
