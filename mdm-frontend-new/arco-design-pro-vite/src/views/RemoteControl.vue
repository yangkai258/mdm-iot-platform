<template>
  <div class="remote-control-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="远程控制记录" :value="856" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="执行成功率" :value="98" suffix="%" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="平均响应" :value="120" suffix="ms" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备远程控制</span>
          <a-space>
            <a-input-search placeholder="搜索设备ID" style="width: 200px;" />
            <a-button type="primary" @click="handleSendCommand">
              <template #icon><icon-send /></template>
              发送指令
            </a-button>
          </a-space>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="devices" title="在线设备">
          <a-table :columns="deviceColumns" :data="devices" :pagination="pagination">
            <template #status="{ record }">
              <a-badge :status="record.status === 'online' ? 'success' : 'error'" :text="record.statusText" />
            </template>
            <template #actions="{ record }">
              <a-button size="small" @click="handleControl(record)">控制</a-button>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="history" title="控制历史">
          <a-table :columns="historyColumns" :data="history" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="schedules" title="定时任务">
          <a-table :columns="scheduleColumns" :data="schedules" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <a-modal v-model:visible="controlVisible" title="设备控制" @before-ok="handleSubmitCommand">
      <a-form :model="controlForm" layout="vertical">
        <a-form-item label="设备">
          <a-input :model-value="controlForm.deviceId" disabled />
        </a-form-item>
        <a-form-item label="控制指令">
          <a-select v-model="controlForm.command" placeholder="选择指令">
            <a-option value="reboot">重启设备</a-option>
            <a-option value="shutdown">关机</a-option>
            <a-option value="reset">恢复出厂设置</a-option>
            <a-option value="sync">同步配置</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="执行参数">
          <a-textarea v-model="controlForm.params" :rows="2" placeholder="JSON格式参数" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 8 });
const controlVisible = ref(false);

const controlForm = reactive({ deviceId: '', deviceName: '', command: '', params: '' });

const devices = ref([
  { deviceId: 'DEV001', deviceName: '小黄-客厅', firmware: 'v2.0.0', status: 'online', statusText: '在线', lastSeen: '2026-03-28 18:00:00' },
  { deviceId: 'DEV002', deviceName: '小红-卧室', firmware: 'v2.0.0', status: 'online', statusText: '在线', lastSeen: '2026-03-28 18:00:00' },
]);

const history = ref([
  { id: 1, deviceId: 'DEV001', deviceName: '小黄-客厅', command: 'reboot', status: 'success', statusText: '成功', time: '2026-03-28 18:00:00' },
  { id: 2, deviceId: 'DEV002', deviceName: '小红-卧室', command: 'sync', status: 'success', statusText: '成功', time: '2026-03-28 17:00:00' },
]);

const schedules = ref([
  { id: 1, deviceName: '小黄-客厅', command: 'reboot', schedule: '每天 02:00', enabled: true },
]);

const deviceColumns = [
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'deviceName', width: 150 },
  { title: '固件版本', dataIndex: 'firmware', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '最后在线', dataIndex: 'lastSeen', width: 160 },
  { title: '操作', slotName: 'actions', width: 100, fixed: 'right' },
];

const historyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '设备名称', dataIndex: 'deviceName', width: 150 },
  { title: '指令', dataIndex: 'command', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
];

const scheduleColumns = [
  { title: '设备', dataIndex: 'deviceName', width: 150 },
  { title: '指令', dataIndex: 'command', width: 100 },
  { title: '执行时间', dataIndex: 'schedule', width: 120 },
  { title: '状态', slotName: 'status', width: 80 },
];

const getStatusColor = (s: string) => ({ success: 'green', failed: 'red', pending: 'orange' }[s] || 'default');

const handleSendCommand = () => {};
const handleControl = (record: any) => { Object.assign(controlForm, record); controlVisible.value = true; };
const handleSubmitCommand = (done: boolean) => { done(true); controlVisible.value = false; };
</script>

<style scoped>
.remote-control-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
