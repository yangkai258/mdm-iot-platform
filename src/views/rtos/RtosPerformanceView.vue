<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备监控</a-breadcrumb-item>
      <a-breadcrumb-item>RTOS 性能监控</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="在线设备" :value="stats.online">
            <template #prefix>🟢</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="离线设备" :value="stats.offline" :value-style="{ color: '#86909c' }">
            <template #prefix>⚫</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="平均 CPU 使用率" :value="stats.avg_cpu" suffix="%" :precision="1">
            <template #prefix>📊</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="平均内存使用率" :value="stats.avg_mem" suffix="%" :precision="1" :value-style="{ color: getMemColor(stats.avg_mem) }">
            <template #prefix>💾</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选区 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-space wrap>
        <a-select v-model="filterForm.status" placeholder="设备状态" style="width: 140px" allow-clear @change="loadDevices">
          <a-option value="online">在线</a-option>
          <a-option value="offline">离线</a-option>
        </a-select>
        <a-input-search v-model="filterForm.keyword" placeholder="搜索设备ID/名称" style="width: 220px" search-button @search="loadDevices" @change="e => !e.target.value && loadDevices()" />
        <a-button @click="loadDevices">
          <template #icon><icon-refresh :spin="loading" /></template>
          刷新
        </a-button>
        <a-button type="primary" @click="showMonitorPanel">
          <template #icon><icon-line-chart /></template>
          监控面板
        </a-button>
      </a-space>
    </a-card>

    <!-- 设备列表 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-spin :loading="loading" tip="加载中...">
        <a-table :data="devices" :pagination="{ pageSize: 10, total: devices.length, showTotal: true }" :columns="columns" row-key="device_id" stripe>
          <template #status="{ record }">
            <a-tag :color="record.status === 'online' ? 'green' : 'gray'">
              {{ record.status === 'online' ? '在线' : '离线' }}
            </a-tag>
          </template>
          <template #cpu="{ record }">
            <a-progress :percent="record.cpu" size="small" :stroke-color="getCpuColor(record.cpu)" />
            <span style="font-size: 12px; color: #86909c">{{ record.cpu }}%</span>
          </template>
          <template #memory="{ record }">
            <a-progress :percent="record.memory" size="small" :stroke-color="getMemColor(record.memory)" />
            <span style="font-size: 12px; color: #86909c">{{ record.memory }}%</span>
          </template>
          <template #uptime="{ record }">
            {{ formatUptime(record.uptime) }}
          </template>
          <template #temperature="{ record }">
            <span :style="{ color: record.temperature > 70 ? '#f53f3f' : record.temperature > 50 ? '#ff7d00' : '#0fc6c2' }">
              {{ record.temperature }}°C
            </span>
          </template>
          <template #operations="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="viewMetrics(record)">详情</a-button>
              <a-button type="text" size="small" @click="sendCommand(record)">下发指令</a-button>
            </a-space>
          </template>
        </a-table>
      </a-spin>
    </a-card>

    <!-- 实时监控面板 -->
    <a-drawer v-model:visible="panelVisible" title="实时监控面板" :width="720" unmountOnHide>
      <a-space direction="vertical" style="width: 100%" size="large">
        <a-card title="设备基本信息">
          <a-descriptions :column="2">
            <a-descriptions-item label="设备ID">{{ currentDevice?.device_id }}</a-descriptions-item>
            <a-descriptions-item label="设备名称">{{ currentDevice?.device_name }}</a-descriptions-item>
            <a-descriptions-item label="固件版本">{{ currentDevice?.firmware_version }}</a-descriptions-item>
            <a-descriptions-item label="在线状态">
              <a-tag :color="currentDevice?.status === 'online' ? 'green' : 'gray'">
                {{ currentDevice?.status === 'online' ? '在线' : '离线' }}
              </a-tag>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-card title="CPU 使用率">
              <a-progress type="circle" :percent="currentDevice?.cpu || 0" :stroke-color="getCpuColor(currentDevice?.cpu)" />
              <div style="text-align: center; margin-top: 8px; font-size: 12px; color: #86909c">当前: {{ currentDevice?.cpu }}%</div>
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card title="内存使用率">
              <a-progress type="circle" :percent="currentDevice?.memory || 0" :stroke-color="getMemColor(currentDevice?.memory)" />
              <div style="text-align: center; margin-top: 8px; font-size: 12px; color: #86909c">当前: {{ currentDevice?.memory }}%</div>
            </a-card>
          </a-col>
        </a-row>

        <a-card title="任务列表">
          <a-table :data="currentDevice?.tasks || []" :pagination="false" :columns="taskColumns" row-key="task_id" size="small">
            <template #stack="{ record }">
              <a-progress :percent="record.stack_usage" size="small" :stroke-color="record.stack_usage > 80 ? 'red' : 'blue'" />
            </template>
          </a-table>
        </a-card>

        <a-card title="历史趋势（模拟）">
          <a-descriptions :column="3">
            <a-descriptions-item label="峰值 CPU">{{ currentDevice?.peak_cpu || 0 }}%</a-descriptions-item>
            <a-descriptions-item label="峰值内存">{{ currentDevice?.peak_memory || 0 }}%</a-descriptions-item>
            <a-descriptions-item label="平均负载">{{ currentDevice?.avg_load || 0 }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-space>
    </a-drawer>

    <!-- 指令下发弹窗 -->
    <a-modal v-model:visible="cmdVisible" title="下发设备指令" @before-ok="handleSendCommand" @cancel="cmdVisible = false">
      <a-form :model="cmdForm" layout="vertical">
        <a-form-item label="目标设备">
          <a-input :model-value="currentDevice?.device_id" disabled />
        </a-form-item>
        <a-form-item label="指令类型" field="cmd_type">
          <a-select v-model="cmdForm.cmd_type" placeholder="请选择指令">
            <a-option value="restart">重启设备</a-option>
            <a-option value="reboot">系统重启</a-option>
            <a-option value="collect">数据采集</a-option>
            <a-option value="update">固件检查</a-option>
            <a-option value="custom">自定义指令</a-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="cmdForm.cmd_type === 'custom'" label="自定义Payload" field="payload">
          <a-textarea v-model="cmdForm.payload" placeholder='{"cmd": "xxx"}' :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const panelVisible = ref(false)
const cmdVisible = ref(false)
const currentDevice = ref(null)

const stats = reactive({ online: 0, offline: 0, avg_cpu: 0, avg_mem: 0 })

const filterForm = reactive({ status: '', keyword: '' })
const cmdForm = reactive({ cmd_type: '', payload: '' })

const devices = ref([])

const columns = [
  { title: '设备ID', dataIndex: 'device_id', width: 140 },
  { title: '设备名称', dataIndex: 'device_name', width: 140 },
  { title: '固件版本', dataIndex: 'firmware_version', width: 100 },
  { title: '状态', dataIndex: 'status', width: 80, slotName: 'status' },
  { title: 'CPU', dataIndex: 'cpu', width: 140, slotName: 'cpu' },
  { title: '内存', dataIndex: 'memory', width: 140, slotName: 'memory' },
  { title: '温度', dataIndex: 'temperature', width: 80, slotName: 'temperature' },
  { title: '运行时长', dataIndex: 'uptime', width: 100, slotName: 'uptime' },
  { title: '操作', slotName: 'operations', width: 140 }
]

const taskColumns = [
  { title: '任务名', dataIndex: 'task_name', width: 120 },
  { title: '优先级', dataIndex: 'priority', width: 80 },
  { title: '状态', dataIndex: 'state', width: 80 },
  { title: '栈使用', dataIndex: 'stack_usage', slotName: 'stack', width: 160 }
]

const getCpuColor = (val) => {
  if (val > 80) return '#f53f3f'
  if (val > 50) return '#ff7d00'
  return '#0fc6c2'
}

const getMemColor = (val) => {
  if (val > 85) return '#f53f3f'
  if (val > 60) return '#ff7d00'
  return '#165dff'
}

const formatUptime = (seconds) => {
  if (!seconds) return '-'
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  return d > 0 ? `${d}天${h}时` : `${h}时${m}分`
}

const loadDevices = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 600))
  devices.value = [
    { device_id: 'DEV-RTOS-001', device_name: '边缘网关-A', firmware_version: 'v2.1.0', status: 'online', cpu: 42, memory: 65, temperature: 48, uptime: 259200, peak_cpu: 78, peak_memory: 82, avg_load: 0.85, tasks: [{ task_id: 1, task_name: 'net_manager', priority: 10, state: 'Running', stack_usage: 45 }, { task_id: 2, task_name: 'data_collector', priority: 15, state: 'Ready', stack_usage: 30 }, { task_id: 3, task_name: 'ota_service', priority: 5, state: 'Blocked', stack_usage: 60 }] },
    { device_id: 'DEV-RTOS-002', device_name: '传感器节点-B', firmware_version: 'v1.9.3', status: 'online', cpu: 18, memory: 32, temperature: 38, uptime: 864000, peak_cpu: 45, peak_memory: 55, avg_load: 0.32, tasks: [{ task_id: 1, task_name: 'sensor_read', priority: 12, state: 'Running', stack_usage: 20 }, { task_id: 2, task_name: 'data_report', priority: 8, state: 'Ready', stack_usage: 35 }] },
    { device_id: 'DEV-RTOS-003', device_name: '工业控制器-C', firmware_version: 'v3.0.1', status: 'online', cpu: 76, memory: 88, temperature: 72, uptime: 172800, peak_cpu: 95, peak_memory: 92, avg_load: 1.82, tasks: [{ task_id: 1, task_name: 'control_loop', priority: 1, state: 'Running', stack_usage: 85 }, { task_id: 2, task_name: 'safety_monitor', priority: 3, state: 'Running', stack_usage: 70 }, { task_id: 3, task_name: 'log_service', priority: 18, state: 'Ready', stack_usage: 25 }] },
    { device_id: 'DEV-RTOS-004', device_name: '智能终端-D', firmware_version: 'v2.5.0', status: 'offline', cpu: 0, memory: 0, temperature: 25, uptime: 0, peak_cpu: 55, peak_memory: 68, avg_load: 0.45, tasks: [] },
    { device_id: 'DEV-RTOS-005', device_name: '数据采集器-E', firmware_version: 'v1.7.2', status: 'online', cpu: 28, memory: 45, temperature: 52, uptime: 432000, peak_cpu: 62, peak_memory: 71, avg_load: 0.58, tasks: [{ task_id: 1, task_name: 'mqtt_client', priority: 10, state: 'Running', stack_usage: 40 }, { task_id: 2, task_name: 'store_forward', priority: 12, state: 'Ready', stack_usage: 55 }] }
  ]
  stats.online = devices.value.filter(d => d.status === 'online').length
  stats.offline = devices.value.filter(d => d.status === 'offline').length
  const onlineDevices = devices.value.filter(d => d.status === 'online')
  stats.avg_cpu = onlineDevices.length ? Math.round(onlineDevices.reduce((s, d) => s + d.cpu, 0) / onlineDevices.length) : 0
  stats.avg_mem = onlineDevices.length ? Math.round(onlineDevices.reduce((s, d) => s + d.memory, 0) / onlineDevices.length) : 0
  loading.value = false
}

const viewMetrics = (record) => {
  currentDevice.value = record
  panelVisible.value = true
}

const showMonitorPanel = () => {
  if (devices.value.length) {
    currentDevice.value = devices.value[0]
    panelVisible.value = true
  }
}

const sendCommand = (record) => {
  currentDevice.value = record
  Object.assign(cmdForm, { cmd_type: '', payload: '' })
  cmdVisible.value = true
}

const handleSendCommand = async (done) => {
  if (!cmdForm.cmd_type) {
    Message.error('请选择指令类型')
    done(false)
    return
  }
  await new Promise(r => setTimeout(r, 500))
  Message.success(`指令已下发: ${cmdForm.cmd_type}`)
  cmdVisible.value = false
  done(true)
}

loadDevices()
</script>

<style scoped>
.page-container { padding: 20px; }
.breadcrumb { margin-bottom: 12px; }
.stat-card { text-align: center; }
.section-card { margin-bottom: 16px; }
</style>
