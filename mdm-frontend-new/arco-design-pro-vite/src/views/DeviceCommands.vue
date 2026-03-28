<template>
  <div class="device-commands-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备指令下发</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon><icon-plus /></template>
            发送指令
          </a-button>
        </div>
      </template>
      
      <a-tabs>
        <a-tab-pane key="history" title="指令历史">
          <a-table :columns="columns" :data="data" :loading="loading" :pagination="pagination">
            <template #type="{ record }">
              <a-tag :color="getTypeColor(record.type)">{{ record.typeText }}</a-tag>
            </template>
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">{{ record.statusText }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="templates" title="指令模板">
          <a-row :gutter="16">
            <a-col :span="6" v-for="tpl in templates" :key="tpl.id">
              <a-card size="small" class="tpl-card" @click="handleUseTemplate(tpl)">
                <div class="tpl-icon">{{ tpl.icon }}</div>
                <div class="tpl-name">{{ tpl.name }}</div>
                <div class="tpl-desc">{{ tpl.description }}</div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 发送指令弹窗 -->
    <a-modal v-model:visible="sendVisible" title="发送设备指令" :width="600" @before-ok="handleSend">
      <a-form :model="form" layout="vertical">
        <a-form-item label="目标设备" required>
          <a-select v-model="form.deviceId" placeholder="选择设备" show-search>
            <a-option value="DEV001">DEV001 - 小黄(在线)</a-option>
            <a-option value="DEV002">DEV002 - 小红(在线)</a-option>
            <a-option value="DEV003">DEV003 - 小绿(离线)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="指令类型" required>
          <a-radio-group v-model="form.type">
            <a-radio value="reboot">重启设备</a-radio>
            <a-radio value="reset">恢复出厂</a-radio>
            <a-radio value="sync">同步状态</a-radio>
            <a-radio value="config">更新配置</a-radio>
            <a-radio value="custom">自定义</a-radio>
          </a-radio-group>
        </a-form-item>
        
        <a-form-item v-if="form.type === 'config'" label="配置参数">
          <a-form-item label="设备模式">
            <a-select v-model="form.params.mode">
              <a-option value="active">活跃模式</a-option>
              <a-option value="quiet">安静模式</a-option>
              <a-option value="dnd">勿扰模式</a-option>
              <a-option value="sleep">休眠模式</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="音量">
            <a-slider v-model="form.params.volume" :marks="{0:'0%',50:'50%',100:'100%'}" />
          </a-form-item>
          <a-form-item label="亮度">
            <a-slider v-model="form.params.brightness" :marks="{0:'0%',50:'50%',100:'100%'}" />
          </a-form-item>
          <a-form-item label="DND启用">
            <a-switch v-model="form.params.dnd" />
          </a-form-item>
        </a-form-item>
        
        <a-form-item v-if="form.type === 'custom'" label="自定义指令">
          <a-textarea v-model="form.customCommand" placeholder="JSON格式自定义指令" :rows="4" />
        </a-form-item>
        
        <a-form-item label="执行方式">
          <a-radio-group v-model="form.executeType">
            <a-radio value="immediate">立即执行</a-radio>
            <a-radio value="scheduled">定时执行</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.executeType === 'scheduled'" label="执行时间">
          <a-date-picker v-model="form.scheduledAt" show-time />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const loading = ref(false);
const data = ref([
  { id: 'C001', deviceId: 'DEV001', type: 'reboot', typeText: '重启', params: '-', status: 'success', statusText: '成功', sentAt: '2026-03-28 18:00:00', completedAt: '2026-03-28 18:00:05' },
  { id: 'C002', deviceId: 'DEV001', type: 'config', typeText: '更新配置', params: '{"mode":"active","volume":80}', status: 'success', statusText: '成功', sentAt: '2026-03-28 17:00:00', completedAt: '2026-03-28 17:00:03' },
  { id: 'C003', deviceId: 'DEV002', type: 'sync', typeText: '同步状态', params: '-', status: 'pending', statusText: '待发送', sentAt: '2026-03-28 16:00:00', completedAt: null },
  { id: 'C004', deviceId: 'DEV003', type: 'reboot', typeText: '重启', params: '-', status: 'failed', statusText: '失败', sentAt: '2026-03-28 15:00:00', completedAt: '2026-03-28 15:00:30', error: '设备离线' },
]);

const templates = ref([
  { id: 1, icon: '🔄', name: '重启设备', description: '重启M5Stack设备', type: 'reboot' },
  { id: 2, icon: '⚡', name: '切换活跃模式', description: '设备进入活跃响应模式', type: 'config', params: { mode: 'active' } },
  { id: 3, icon: '😴', name: '切换休眠模式', description: '设备进入低功耗休眠', type: 'config', params: { mode: 'sleep' } },
  { id: 4, icon: '🔇', name: '启用勿扰', description: '关闭声音和灯光提醒', type: 'config', params: { dnd: true } },
  { id: 5, icon: '🔊', name: '同步状态', description: '强制同步设备当前状态', type: 'sync' },
  { id: 6, icon: '🏭', name: '恢复出厂', description: '清除所有数据恢复出厂设置', type: 'reset' },
]);

const pagination = reactive({ current: 1, pageSize: 20, total: 4 });
const sendVisible = ref(false);

const form = reactive({
  deviceId: '',
  type: 'reboot',
  params: { mode: 'active', volume: 50, brightness: 50, dnd: false },
  customCommand: '',
  executeType: 'immediate',
  scheduledAt: null,
});

const columns = [
  { title: '指令ID', dataIndex: 'id', width: 100 },
  { title: '设备ID', dataIndex: 'deviceId', width: 120 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '参数', dataIndex: 'params', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '发送时间', dataIndex: 'sentAt', width: 160 },
  { title: '完成时间', dataIndex: 'completedAt', width: 160 },
];

const getTypeColor = (t: string) => ({ reboot: 'blue', config: 'green', sync: 'purple', reset: 'orange', custom: 'cyan' }[t] || 'default');
const getStatusColor = (s: string) => ({ success: 'green', failed: 'red', pending: 'orange' }[s] || 'default');

const handleCreate = () => { sendVisible.value = true; };
const handleUseTemplate = (tpl: any) => {
  form.type = tpl.type;
  if (tpl.params) Object.assign(form.params, tpl.params);
  sendVisible.value = true;
};
const handleSend = (done: boolean) => { done(true); sendVisible.value = false; };
</script>

<style scoped>
.device-commands-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
.tpl-card { text-align: center; cursor: pointer; transition: all 0.3s; margin-bottom: 12px; }
.tpl-card:hover { border-color: #165DFF; }
.tpl-icon { font-size: 32px; margin-bottom: 8px; }
.tpl-name { font-weight: bold; margin-bottom: 4px; }
.tpl-desc { font-size: 12px; color: #86909c; }
</style>
