<template>
  <div class="device-detail-container">
    <a-card>
      <template #title>
        <div class="card-title">
          <span>设备详情</span>
          <a-space>
            <a-button type="primary" @click="handleCommand">
              <template #icon><icon-send /></template>
              发送指令
            </a-button>
            <a-button @click="handleRefresh">刷新</a-button>
          </a-space>
        </div>
      </template>
      
      <a-descriptions bordered :column="2">
        <a-descriptions-item label="设备ID">{{ device.deviceId }}</a-descriptions-item>
        <a-descriptions-item label="设备名称">{{ device.name }}</a-descriptions-item>
        <a-descriptions-item label="设备类型">{{ device.type }}</a-descriptions-item>
        <a-descriptions-item label="固件版本">{{ device.firmware }}</a-descriptions-item>
        <a-descriptions-item label="在线状态">
          <a-tag :color="device.online ? 'green' : 'red'">
            {{ device.online ? '在线' : '离线' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="最后活跃">{{ device.lastSeen }}</a-descriptions-item>
        <a-descriptions-item label="运行时长">{{ device.uptime }}</a-descriptions-item>
        <a-descriptions-item label="电池电量">
          <a-progress :percent="device.battery" :color="getBatteryColor(device.battery)" />
        </a-descriptions-item>
      </a-descriptions>
    </a-card>

    <!-- 设备影子 -->
    <a-card title="设备影子" style="margin-top: 16px;">
      <a-tabs>
        <a-tab-pane key="desired" title="期望状态">
          <a-descriptions bordered :column="2">
            <a-descriptions-item label="模式">{{ shadow.desired.mode }}</a-descriptions-item>
            <a-descriptions-item label="DND">
              <a-tag :color="shadow.desired.dnd ? 'orange' : 'green'">
                {{ shadow.desired.dnd ? '免打扰' : '正常' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="音量">{{ shadow.desired.volume }}%</a-descriptions-item>
            <a-descriptions-item label="亮度">{{ shadow.desired.brightness }}%</a-descriptions-item>
          </a-descriptions>
          <div style="margin-top: 16px;">
            <a-button type="primary" @click="handleEditShadow">编辑期望状态</a-button>
          </div>
        </a-tab-pane>
        <a-tab-pane key="reported" title="实际状态">
          <a-descriptions bordered :column="2">
            <a-descriptions-item label="当前模式">{{ shadow.reported.mode }}</a-descriptions-item>
            <a-descriptions-item label="DND">
              <a-tag :color="shadow.reported.dnd ? 'orange' : 'green'">
                {{ shadow.reported.dnd ? '免打扰' : '正常' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="音量">{{ shadow.reported.volume }}%</a-descriptions-item>
            <a-descriptions-item label="亮度">{{ shadow.reported.brightness }}%</a-descriptions-item>
          </a-descriptions>
        </a-tab-pane>
        <a-tab-pane key="delta" title="差异">
          <a-alert v-if="hasDelta" type="warning" :closable="false">
            检测到设备状态与期望不一致
          </a-alert>
          <a-descriptions v-else bordered :column="2">
            <a-descriptions-item label="状态">设备状态正常，无差异</a-descriptions-item>
          </a-descriptions>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 设备模式 -->
    <a-card title="设备模式" style="margin-top: 16px;">
      <a-radio-group v-model="deviceMode" type="button">
        <a-radio value="active">活跃模式</a-radio>
        <a-radio value="quiet">安静模式</a-radio>
        <a-radio value="dnd">勿扰模式</a-radio>
        <a-radio value="sleep">休眠模式</a-radio>
      </a-radio-group>
      <div style="margin-top: 16px;">
        <a-button type="primary" @click="handleSetMode">设置模式</a-button>
      </div>
    </a-card>

    <!-- 指令下发 -->
    <a-card title="指令记录" style="margin-top: 16px;">
      <a-table :columns="commandColumns" :data="commands" :pagination="pagination">
        <template #type="{ record }">
          <a-tag :color="getCommandTypeColor(record.type)">{{ record.typeText }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="getCommandStatusColor(record.status)">{{ record.statusText }}</a-tag>
        </template>
      </a-table>
    </a-card>

    <!-- OTA进度 -->
    <a-card title="OTA升级" style="margin-top: 16px;">
      <a-descriptions bordered :column="2">
        <a-descriptions-item label="当前版本">{{ ota.currentVersion }}</a-descriptions-item>
        <a-descriptions-item label="目标版本">{{ ota.targetVersion || '-' }}</a-descriptions-item>
        <a-descriptions-item label="升级状态">
          <a-tag :color="getOTAStatusColor(ota.status)">{{ ota.statusText }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="进度">
          <a-progress v-if="ota.status === 'upgrading'" :percent="ota.progress" status="success" />
          <span v-else>{{ ota.progress }}%</span>
        </a-descriptions-item>
      </a-descriptions>
      <div style="margin-top: 16px;">
        <a-button type="primary" @click="handleOTA">执行OTA升级</a-button>
      </div>
    </a-card>

    <!-- 版本历史 -->
    <a-card title="版本历史" style="margin-top: 16px;">
      <a-table :columns="historyColumns" :data="history" :pagination="false" />
    </a-card>

    <!-- 指令下发弹窗 -->
    <a-modal v-model:visible="commandVisible" title="发送指令" @before-ok="handleSendCommand">
      <a-form :model="commandForm">
        <a-form-item label="指令类型" required>
          <a-select v-model="commandForm.type" placeholder="选择指令类型">
            <a-option value="reboot">重启设备</a-option>
            <a-option value="reset">恢复出厂设置</a-option>
            <a-option value="sync">同步状态</a-option>
            <a-option value="config">更新配置</a-option>
            <a-option value="custom">自定义指令</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="commandForm.type === 'config'" label="配置参数">
          <a-textarea v-model="commandForm.params" placeholder="JSON格式配置参数" :rows="4" />
        </a-form-item>
        <a-form-item v-if="commandForm.type === 'custom'" label="自定义指令">
          <a-textarea v-model="commandForm.customCommand" placeholder="自定义指令内容" :rows="4" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue';

const device = reactive({
  deviceId: 'DEV001',
  name: '小黄-客厅',
  type: 'M5Stack',
  firmware: 'v2.0.0',
  online: true,
  lastSeen: '2026-03-28 18:30:00',
  uptime: '15天3小时',
  battery: 85,
});

const shadow = reactive({
  desired: { mode: 'active', dnd: false, volume: 80, brightness: 70 },
  reported: { mode: 'active', dnd: false, volume: 75, brightness: 68 },
});

const deviceMode = ref('active');

const commands = ref([
  { id: 'C001', type: 'reboot', typeText: '重启', params: '-', status: 'success', statusText: '成功', sentAt: '2026-03-28 10:00:00' },
  { id: 'C002', type: 'config', typeText: '更新配置', params: '{"volume":80}', status: 'success', statusText: '成功', sentAt: '2026-03-28 09:00:00' },
  { id: 'C003', type: 'sync', typeText: '同步状态', params: '-', status: 'success', statusText: '成功', sentAt: '2026-03-27 18:00:00' },
  { id: 'C004', type: 'reboot', typeText: '重启', params: '-', status: 'timeout', statusText: '超时', sentAt: '2026-03-27 10:00:00' },
]);

const history = ref([
  { version: 'v2.0.0', time: '2026-03-20 10:00:00', type: 'upgrade', description: '升级到v2.0.0' },
  { version: 'v1.1.0', time: '2026-03-15 14:00:00', type: 'upgrade', description: '升级到v1.1.0' },
  { version: 'v1.0.0', time: '2026-03-01 10:00:00', type: 'init', description: '初始版本' },
]);

const ota = reactive({
  currentVersion: 'v2.0.0',
  targetVersion: '',
  status: 'idle',
  statusText: '无升级',
  progress: 0,
});

const pagination = reactive({ current: 1, pageSize: 10, total: 4 });

const commandColumns = [
  { title: '指令ID', dataIndex: 'id', width: 100 },
  { title: '类型', slotName: 'type', width: 100 },
  { title: '参数', dataIndex: 'params', width: 200 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '发送时间', dataIndex: 'sentAt', width: 160 },
];

const historyColumns = [
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '时间', dataIndex: 'time', width: 160 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '描述', dataIndex: 'description' },
];

const commandVisible = ref(false);
const commandForm = reactive({
  type: '',
  params: '',
  customCommand: '',
});

const hasDelta = computed(() => {
  return shadow.desired.mode !== shadow.reported.mode ||
         shadow.desired.volume !== shadow.reported.volume;
});

const getBatteryColor = (battery: number) => {
  if (battery >= 60) return '#00B42A';
  if (battery >= 30) return '#FF7D00';
  return '#F53F3F';
};

const getCommandTypeColor = (type: string) => {
  const map: Record<string, string> = { reboot: 'blue', config: 'green', sync: 'purple', custom: 'orange' };
  return map[type] || 'default';
};

const getCommandStatusColor = (status: string) => {
  const map: Record<string, string> = { success: 'green', failed: 'red', timeout: 'orange', pending: 'blue' };
  return map[status] || 'default';
};

const getOTAStatusColor = (status: string) => {
  const map: Record<string, string> = { idle: 'gray', upgrading: 'blue', success: 'green', failed: 'red' };
  return map[status] || 'default';
};

const handleCommand = () => { commandVisible.value = true; };
const handleRefresh = () => { /* TODO */ };
const handleEditShadow = () => { /* TODO */ };
const handleSetMode = () => { /* TODO */ };
const handleSendCommand = (done: boolean) => { done(true); commandVisible.value = false; };
const handleOTA = () => { /* TODO */ };
</script>

<style scoped>
.device-detail-container { padding: 20px; }
.card-title { display: flex; justify-content: space-between; align-items: center; }
</style>
