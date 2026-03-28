<template>
  <div class="ble-debug-container">
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="8">
        <a-card><a-statistic title="BLE设备" :value="12" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="Mesh网络" :value="3" /></a-card>
      </a-col>
      <a-col :span="8">
        <a-card><a-statistic title="在线设备" :value="10" status="success" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <template #title>
        <span>BLE Mesh配置</span>
      </template>
      
      <a-tabs>
        <a-tab-pane key="devices" title="BLE设备">
          <a-table :columns="deviceColumns" :data="devices" :pagination="pagination">
            <template #status="{ record }">
              <a-tag :color="record.online ? 'green' : 'gray'">{{ record.online ? '在线' : '离线' }}</a-tag>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="mesh" title="Mesh网络">
          <a-row :gutter="16">
            <a-col :span="8" v-for="net in meshNetworks" :key="net.id">
              <a-card size="small" class="mesh-card">
                <div class="mesh-name">{{ net.name }}</div>
                <div class="mesh-devices">设备数: {{ net.deviceCount }}</div>
                <a-progress :percent="net.healthScore" :color="net.healthScore > 80 ? '#00B42A' : '#FF7D00'" />
                <div class="mesh-health">健康度: {{ net.healthScore }}%</div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
        
        <a-tab-pane key="logs" title="调试日志">
          <a-table :columns="logColumns" :data="logs" :pagination="paginationSmall" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const pagination = reactive({ current: 1, pageSize: 10, total: 12 });
const paginationSmall = reactive({ current: 1, pageSize: 5, total: 20 });

const devices = ref([
  { id: 1, name: 'BLE-001', mac: 'AA:BB:CC:DD:EE:01', type: 'M5Stack', rssi: -65, online: true, lastSeen: '2026-03-28 18:00:00' },
  { id: 2, name: 'BLE-002', mac: 'AA:BB:CC:DD:EE:02', type: 'ESP32', rssi: -72, online: true, lastSeen: '2026-03-28 18:00:00' },
  { id: 3, name: 'BLE-003', mac: 'AA:BB:CC:DD:EE:03', type: 'M5Stack', rssi: -85, online: false, lastSeen: '2026-03-28 10:00:00' },
]);

const meshNetworks = ref([
  { id: 1, name: 'Mesh-Net-1', deviceCount: 8, healthScore: 92 },
  { id: 2, name: 'Mesh-Net-2', deviceCount: 5, healthScore: 78 },
  { id: 3, name: 'Mesh-Net-3', deviceCount: 3, healthScore: 65 },
]);

const logs = ref([
  { time: '18:00:00', device: 'BLE-001', level: 'info', message: '连接成功' },
  { time: '17:59:55', device: 'BLE-002', level: 'warn', message: '信号强度低' },
  { time: '17:59:30', device: 'BLE-003', level: 'error', message: '连接超时' },
]);

const deviceColumns = [
  { title: '设备ID', dataIndex: 'name', width: 100 },
  { title: 'MAC地址', dataIndex: 'mac', width: 180 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: 'RSSI', dataIndex: 'rssi', width: 80 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '最后活跃', dataIndex: 'lastSeen', width: 160 },
];

const logColumns = [
  { title: '时间', dataIndex: 'time', width: 100 },
  { title: '设备', dataIndex: 'device', width: 100 },
  { title: '级别', dataIndex: 'level', width: 80 },
  { title: '消息', dataIndex: 'message' },
];
</script>

<style scoped>
.ble-debug-container { padding: 20px; }
.mesh-card { text-align: center; margin-bottom: 12px; }
.mesh-name { font-weight: bold; }
.mesh-devices { color: #86909c; font-size: 12px; margin: 8px 0; }
.mesh-health { color: #86909c; font-size: 12px; margin-top: 8px; }
</style>
