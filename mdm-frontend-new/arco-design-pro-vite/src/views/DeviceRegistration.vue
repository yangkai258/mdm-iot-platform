<template>
  <div class="device-registration-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="注册设备" :value="1256" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="今日注册" :value="28" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="激活率" :value="92" suffix="%" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备注册管理</span>
          <a-space>
            <a-input-search placeholder="搜索设备ID/名称" style="width: 200px;" />
            <a-button @click="handleExport">导出</a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="all" title="全部设备">
          <a-table :columns="columns" :data="devices" :pagination="pagination">
            <template #status="{ record }">
              <a-badge :status="getStatusBadge(record.status)" :text="record.statusText" />
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-link @click="handleView(record)">详情</a-link>
                <a-link @click="handleConfig(record)">配置</a-link>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="pending" title="待激活">
          <a-table :columns="columns" :data="pendingDevices" :pagination="pagination" />
        </a-tab-pane>
        
        <a-tab-pane key="templates" title="注册配置">
          <a-form layout="vertical" style="max-width: 600px;">
            <a-form-item label="自动激活">
              <a-switch v-model="config.autoActivate" />
              <span style="margin-left: 8px; color: #86909c;">开启后设备注册自动激活</span>
            </a-form-item>
            <a-form-item label="设备分组">
              <a-select v-model="config.defaultGroup" placeholder="选择默认分组">
                <a-option value="default">默认分组</a-option>
                <a-option value="home">家庭组</a-option>
                <a-option value="business">商用组</a-option>
              </a-select>
            </a-form-item>
            <a-form-item label="固件版本">
              <a-select v-model="config.firmwareVersion" placeholder="选择固件版本">
                <a-option value="v2.0.0">v2.0.0 最新</a-option>
                <a-option value="v1.8.0">v1.8.0 稳定</a-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button type="primary" @click="handleSaveConfig">保存配置</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-drawer v-model:visible="configVisible" title="设备配置" :width="500">
      <a-form :model="deviceConfig" layout="vertical">
        <a-form-item label="设备名称">
          <a-input v-model="deviceConfig.name" />
        </a-form-item>
        <a-form-item label="所属分组">
          <a-select v-model="deviceConfig.groupId">
            <a-option value="G001">家庭组</a-option>
            <a-option value="G002">商用组</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="固件版本">
          <a-select v-model="deviceConfig.firmware">
            <a-option value="v2.0.0">v2.0.0</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="deviceConfig.note" :rows="3" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSaveDevice">保存</a-button>
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 15 });
const configVisible = ref(false);
const config = reactive({ autoActivate: true, defaultGroup: 'default', firmwareVersion: 'v2.0.0' });
const deviceConfig = reactive({ name: '', groupId: '', firmware: '', note: '' });

const devices = ref([
  { id: 'DEV001', name: '小黄-客厅', hardwareId: 'HW12345', firmware: 'v2.0.0', status: 'online', statusText: '在线', registeredAt: '2026-03-20 10:00:00', activatedAt: '2026-03-20 10:05:00' },
  { id: 'DEV002', name: '小红-卧室', hardwareId: 'HW12346', firmware: 'v1.8.0', status: 'offline', statusText: '离线', registeredAt: '2026-03-19 15:00:00', activatedAt: '2026-03-19 15:10:00' },
  { id: 'DEV003', name: '咪咪-阳台', hardwareId: 'HW12347', firmware: 'v2.0.0', status: 'online', statusText: '在线', registeredAt: '2026-03-28 09:00:00', activatedAt: '2026-03-28 09:02:00' },
]);

const pendingDevices = ref([
  { id: 'DEV010', name: '待激活设备', hardwareId: 'HW99999', firmware: 'v2.0.0', status: 'pending', statusText: '待激活', registeredAt: '2026-03-28 18:00:00', activatedAt: null },
]);

const columns = [
  { title: '设备ID', dataIndex: 'id', width: 100 },
  { title: '设备名称', dataIndex: 'name', width: 150 },
  { title: '硬件ID', dataIndex: 'hardwareId', width: 120 },
  { title: '固件版本', dataIndex: 'firmware', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '注册时间', dataIndex: 'registeredAt', width: 160 },
  { title: '激活时间', dataIndex: 'activatedAt', width: 160 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' },
];

const getStatusBadge = (s: string) => ({ online: 'success', offline: 'error', pending: 'warning' }[s] || 'default');

const handleView = (record: any) => {};
const handleConfig = (record: any) => { Object.assign(deviceConfig, record); configVisible.value = true; };
const handleExport = () => {};
const handleSaveConfig = () => {};
const handleSaveDevice = () => { configVisible.value = false; };
</script>

<style scoped>
.device-registration-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
