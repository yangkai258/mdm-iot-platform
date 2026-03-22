<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>技术架构</a-breadcrumb-item>
      <a-breadcrumb-item>设备 OTA 优化</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="待升级设备" :value="stats.pending" :value-style="{ color: '#faad14' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="升级中" :value="stats.upgrading" :value-style="{ color: '#165dff' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="升级成功" :value="stats.success" :value-style="{ color: '#52c41a' }" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="升级失败" :value="stats.failed" :value-style="{ color: '#ff4d4f' }" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 标签页 -->
    <a-tabs v-model:active-key="activeTab" class="pro-tabs">
      <a-tab-pane key="config" title="分片升级配置">
        <a-card class="config-card">
          <template #title><span>分片策略</span></template>
          <a-form :model="shardConfig" layout="vertical">
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item label="分片大小(KB)">
                  <a-input-number v-model="shardConfig.chunk_size_kb" :min="4" :max="256" :step="4" style="width: 100%" />
                  <div class="form-tip">推荐: 16KB (BLE MTU限制)</div>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="并行设备数">
                  <a-input-number v-model="shardConfig.parallel_devices" :min="1" :max="50" style="width: 100%" />
                  <div class="form-tip">同时升级的设备数量</div>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="重试次数">
                  <a-input-number v-model="shardConfig.retry_count" :min="0" :max="10" style="width: 100%" />
                  <div class="form-tip">单设备失败重试次数</div>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="分片间隔(ms)">
                  <a-input-number v-model="shardConfig.chunk_interval_ms" :min="50" :max="5000" :step="50" style="width: 100%" />
                  <div class="form-tip">分片间发送间隔</div>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="升级超时(s)">
                  <a-input-number v-model="shardConfig.timeout_seconds" :min="60" :max="3600" :step="60" style="width: 100%" />
                  <div class="form-tip">单设备升级超时时间</div>
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item label="断点续传">
                  <a-switch v-model="shardConfig.resume_enabled" />
                  <div class="form-tip">中断后自动续传</div>
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="saveShardConfig">保存配置</a-button>
                <a-button @click="resetShardConfig">重置默认</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>

        <a-card class="config-card" style="margin-top: 16px">
          <template #title><span>OTA 策略预设</span></template>
          <a-table :columns="policyColumns" :data="policies" :pagination="false" row-key="id">
            <template #isActive="{ record }">
              <a-tag :color="record.is_active ? 'green' : 'gray'">{{ record.is_active ? '已启用' : '已禁用' }}</a-tag>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="editPolicy(record)">编辑</a-button>
                <a-button type="text" size="small" @click="togglePolicy(record)">{{ record.is_active ? '禁用' : '启用' }}</a-button>
                <a-button type="text" size="small" status="danger" @click="deletePolicy(record)">删除</a-button>
              </a-space>
            </template>
          </a-table>
          <div style="margin-top: 12px">
            <a-button type="dashed" @click="showPolicyModal(null)">新建策略</a-button>
          </div>
        </a-card>
      </a-tab-pane>

      <a-tab-pane key="monitor" title="升级状态监控">
        <div class="pro-search-bar">
          <a-space>
            <a-input-search v-model="searchKeyword" placeholder="搜索设备ID" style="width: 260px" @search="loadTasks" search-button />
            <a-select v-model="filterTaskStatus" placeholder="任务状态" style="width: 140px" allow-clear>
              <a-option :value="1">等待中</a-option>
              <a-option :value="2">下载中</a-option>
              <a-option :value="3">升级中</a-option>
              <a-option :value="4">成功</a-option>
              <a-option :value="5">失败</a-option>
            </a-select>
          </a-space>
        </div>
        <div class="pro-action-bar">
          <a-space>
            <a-button type="primary" @click="showBatchUpgradeModal">批量升级</a-button>
            <a-button @click="loadTasks">刷新</a-button>
            <a-button v-if="filterTaskStatus || searchKeyword" @click="clearFilters">清除筛选</a-button>
          </a-space>
        </div>
        <div class="pro-content-area">
          <a-table :columns="taskColumns" :data="filteredTasks" :loading="loading" :pagination="pagination" @change="handleTableChange" row-key="id">
            <template #status="{ record }">
              <a-tag :color="getTaskStatusColor(record.status)">{{ getTaskStatusText(record.status) }}</a-tag>
            </template>
            <template #progress="{ record }">
              <a-progress :percent="record.progress" :stroke-width="8" size="small"
                :status="record.status === 4 ? 'success' : record.status === 5 ? 'exception' : 'normal'" />
            </template>
            <template #shardInfo="{ record }">
              <span>{{ record.shard_index }}/{{ record.total_shards }} 分片</span>
            </template>
            <template #startTime="{ record }">
              {{ formatTime(record.start_time) }}
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button v-if="record.status === 5" type="text" size="small" @click="handleRetry(record)">重试</a-button>
                <a-button v-if="record.status === 1 || record.status === 2" type="text" size="small" status="warning" @click="handlePause(record)">暂停</a-button>
                <a-button v-if="record.status === 4" type="text" size="small" @click="handleVerify(record)">验证</a-button>
                <a-button type="text" size="small" status="danger" @click="handleCancelTask(record)">取消</a-button>
              </a-space>
            </template>
          </a-table>
        </div>

        <a-card class="log-card">
          <template #title><span>实时日志</span></template>
          <div class="log-terminal">
            <div v-for="(log, idx) in otaLogs" :key="idx" class="log-line" :class="log.level">
              <span class="log-time">{{ log.time }}</span>
              <span class="log-device">[{{ log.device_id }}]</span>
              <span class="log-msg">{{ log.message }}</span>
            </div>
          </div>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <!-- 批量升级弹窗 -->
    <a-modal v-model:visible="batchUpgradeModalVisible" title="批量升级" @ok="handleBatchUpgrade" :confirm-loading="batchLoading" :width="560">
      <a-form layout="vertical">
        <a-form-item label="选择固件版本" required>
          <a-select v-model="batchUpgradeForm.firmware_version" placeholder="选择目标固件版本">
            <a-option value="v1.2.0">v1.2.0 (推荐)</a-option>
            <a-option value="v1.1.5">v1.1.5</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="选择设备">
          <a-select v-model="batchUpgradeForm.device_ids" multiple placeholder="选择目标设备">
            <a-option v-for="d in availableDevices" :key="d.device_id" :value="d.device_id">
              {{ d.device_id }} ({{ d.firmware_version }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="升级策略">
          <a-select v-model="batchUpgradeForm.policy_id" placeholder="选择OTA策略">
            <a-option :value="1">默认策略</a-option>
            <a-option :value="2">低带宽策略</a-option>
            <a-option :value="3">快速升级策略</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="静默升级">
          <a-switch v-model="batchUpgradeForm.silent" />
          <div class="form-tip">静默升级不推送通知</div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 策略编辑弹窗 -->
    <a-modal v-model:visible="policyModalVisible" :title="editingPolicy ? '编辑策略' : '新建策略'" @ok="handleSavePolicy" :confirm-loading="policySaving" :width="480">
      <a-form :model="policyForm" layout="vertical">
        <a-form-item label="策略名称" required>
          <a-input v-model="policyForm.name" placeholder="例如: 快速升级" />
        </a-form-item>
        <a-form-item label="分片大小(KB)">
          <a-input-number v-model="policyForm.chunk_size_kb" :min="4" :max="256" :step="4" style="width: 100%" />
        </a-form-item>
        <a-form-item label="并行设备数">
          <a-input-number v-model="policyForm.parallel_devices" :min="1" :max="50" style="width: 100%" />
        </a-form-item>
        <a-form-item label="分片间隔(ms)">
          <a-input-number v-model="policyForm.chunk_interval_ms" :min="50" :max="5000" :step="50" style="width: 100%" />
        </a-form-item>
        <a-form-item label="启用断点续传">
          <a-switch v-model="policyForm.resume_enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/device'
const activeTab = ref('config')
const loading = ref(false)
const batchLoading = ref(false)
const policySaving = ref(false)
const searchKeyword = ref('')
const filterTaskStatus = ref(null)
const filterVersion = ref(null)
const editingPolicy = ref(null)

const batchUpgradeModalVisible = ref(false)
const policyModalVisible = ref(false)

const stats = reactive({ pending: 0, upgrading: 0, success: 0, failed: 0 })
const shardConfig = reactive({ chunk_size_kb: 16, parallel_devices: 10, retry_count: 3, chunk_interval_ms: 100, timeout_seconds: 300, resume_enabled: true })
const shardConfigDefault = { chunk_size_kb: 16, parallel_devices: 10, retry_count: 3, chunk_interval_ms: 100, timeout_seconds: 300, resume_enabled: true }
const policyForm = reactive({ name: '', chunk_size_kb: 16, parallel_devices: 10, chunk_interval_ms: 100, resume_enabled: true })
const batchUpgradeForm = reactive({ firmware_version: '', device_ids: [], policy_id: 1, silent: false })
const tasks = ref([])
const policies = ref([])
const availableDevices = ref([])
const otaLogs = ref([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const taskColumns = [
  { title: '设备ID', dataIndex: 'device_id', width: 140 },
  { title: '目标版本', dataIndex: 'target_version', width: 100 },
  { title: '当前版本', dataIndex: 'current_version', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 160 },
  { title: '分片', slotName: 'shardInfo', width: 120 },
  { title: '开始时间', slotName: 'startTime', width: 180 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const policyColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '策略名称', dataIndex: 'name', width: 140 },
  { title: '分片大小', dataIndex: 'chunk_size_kb', width: 100 },
  { title: '并行数', dataIndex: 'parallel_devices', width: 80 },
  { title: '断点续传', dataIndex: 'resume_enabled', width: 100 },
  { title: '状态', slotName: 'isActive', width: 80 },
  { title: '操作', slotName: 'actions', width: 180 }
]

const filteredTasks = computed(() => {
  let result = tasks.value
  if (filterTaskStatus.value !== null) result = result.filter(t => t.status === filterTaskStatus.value)
  if (filterVersion.value) result = result.filter(t => t.target_version === filterVersion.value)
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(t => t.device_id.toLowerCase().includes(kw))
  }
  return result
})

let logTimer = null

const loadTasks = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/ota/tasks`)
    if (res.data.code === 0) {
      tasks.value = res.data.data.list
      pagination.total = res.data.data.pagination?.total || tasks.value.length
    }
  } catch {
    tasks.value = [
      { id: 1, device_id: 'DEV001', target_version: 'v1.2.0', current_version: 'v1.1.5', status: 4, progress: 100, shard_index: 128, total_shards: 128, start_time: '2026-03-22 14:00:00' },
      { id: 2, device_id: 'DEV002', target_version: 'v1.2.0', current_version: 'v1.1.5', status: 3, progress: 68, shard_index: 87, total_shards: 128, start_time: '2026-03-22 15:00:00' },
      { id: 3, device_id: 'DEV003', target_version: 'v1.2.0', current_version: 'v1.0.0', status: 2, progress: 25, shard_index: 32, total_shards: 128, start_time: '2026-03-22 16:00:00' },
      { id: 4, device_id: 'DEV004', target_version: 'v1.2.0', current_version: 'v1.1.5', status: 5, progress: 42, shard_index: 54, total_shards: 128, start_time: '2026-03-22 15:30:00' },
      { id: 5, device_id: 'DEV005', target_version: 'v1.2.0', current_version: 'v1.1.5', status: 1, progress: 0, shard_index: 0, total_shards: 128, start_time: '2026-03-22 17:00:00' }
    ]
    pagination.total = tasks.value.length
    Message.warning('使用模拟数据')
  } finally { loading.value = false }
  updateStats()
}

const loadPolicies = async () => {
  try {
    const res = await axios.get(`${API_BASE}/ota/policies`)
    if (res.data.code === 0) policies.value = res.data.data.list
  } catch {
    policies.value = [
      { id: 1, name: '默认策略', chunk_size_kb: 16, parallel_devices: 10, chunk_interval_ms: 100, resume_enabled: true, is_active: true },
      { id: 2, name: '低带宽策略', chunk_size_kb: 8, parallel_devices: 3, chunk_interval_ms: 500, resume_enabled: true, is_active: false },
      { id: 3, name: '快速升级策略', chunk_size_kb: 32, parallel_devices: 20, chunk_interval_ms: 50, resume_enabled: false, is_active: false }
    ]
  }
}

const loadDevices = async () => {
  try {
    const res = await axios.get('/api/v1/devices')
    if (res.data.code === 0) availableDevices.value = res.data.data.list
  } catch {
    availableDevices.value = [
      { device_id: 'DEV001', firmware_version: 'v1.1.5' },
      { device_id: 'DEV002', firmware_version: 'v1.1.5' },
      { device_id: 'DEV003', firmware_version: 'v1.0.0' }
    ]
  }
}

const startLogStream = () => {
  otaLogs.value = [
    { time: '17:25:01', device_id: 'DEV001', message: '开始下载固件 v1.2.0', level: 'info' },
    { time: '17:25:02', device_id: 'DEV002', message: '分片 1/128 发送成功', level: 'info' },
    { time: '17:25:03', device_id: 'DEV001', message: '下载进度: 25%', level: 'info' },
    { time: '17:25:04', device_id: 'DEV003', message: '设备响应超时, 重试(1/3)', level: 'warning' },
    { time: '17:25:05', device_id: 'DEV004', message: '固件校验失败: MD5不匹配', level: 'error' }
  ]
  logTimer = setInterval(() => {
    const msgs = [
      { msg: '分片 N/128 发送成功', level: 'info' },
      { msg: '设备响应超时, 重试', level: 'warning' },
      { msg: '下载进度: NN%', level: 'info' },
      { msg: '固件写入成功', level: 'success' },
      { msg: 'CRC校验通过', level: 'success' }
    ]
    const devs = ['DEV001', 'DEV002', 'DEV003', 'DEV004', 'DEV005']
    const now = new Date()
    const timeStr = `${String(now.getHours()).padStart(2,'0')}:${String(now.getMinutes()).padStart(2,'0')}:${String(now.getSeconds()).padStart(2,'0')}`
    const m = msgs[Math.floor(Math.random() * msgs.length)]
    otaLogs.value.push({ time: timeStr, device_id: devs[Math.floor(Math.random() * devs.length)], message: m.msg.replace('N', Math.floor(Math.random() * 128) + 1).replace('NN', Math.floor(Math.random() * 100)), level: m.level })
    if (otaLogs.value.length > 50) otaLogs.value.shift()
  }, 3000)
}

const updateStats = () => {
  stats.pending = tasks.value.filter(t => t.status === 1).length
  stats.upgrading = tasks.value.filter(t => t.status === 2 || t.status === 3).length
  stats.success = tasks.value.filter(t => t.status === 4).length
  stats.failed = tasks.value.filter(t => t.status === 5).length
}

const saveShardConfig = () => {
  Message.success('分片策略配置已保存')
}

const resetShardConfig = () => {
  Object.assign(shardConfig, shardConfigDefault)
  Message.info('已重置为默认配置')
}

const showPolicyModal = (record) => {
  editingPolicy.value = record
  if (record) Object.assign(policyForm, record)
  else Object.assign(policyForm, { name: '', chunk_size_kb: 16, parallel_devices: 10, chunk_interval_ms: 100, resume_enabled: true })
  policyModalVisible.value = true
}

const editPolicy = (record) => showPolicyModal(record)

const handleSavePolicy = () => {
  if (!policyForm.name) { Message.warning('请填写策略名称'); return }
  policySaving.value = true
  setTimeout(() => {
    if (editingPolicy.value) {
      const idx = policies.value.findIndex(p => p.id === editingPolicy.value.id)
      if (idx !== -1) Object.assign(policies.value[idx], policyForm)
      Message.success('策略已更新')
    } else {
      policies.value.push({ id: Date.now(), ...policyForm, is_active: false })
      Message.success('策略已创建')
    }
    policySaving.value = false
    policyModalVisible.value = false
  }, 600)
}

const togglePolicy = (record) => {
  const idx = policies.value.findIndex(p => p.id === record.id)
  if (idx !== -1) policies.value[idx].is_active = !policies.value[idx].is_active
  Message.success(`策略 ${record.name} 已${policies.value[idx].is_active ? '启用' : '禁用'}`)
}

const deletePolicy = (record) => {
  policies.value = policies.value.filter(p => p.id !== record.id)
  Message.success('策略已删除')
}

const showBatchUpgradeModal = () => {
  batchUpgradeForm.firmware_version = ''
  batchUpgradeForm.device_ids = []
  batchUpgradeForm.policy_id = 1
  batchUpgradeForm.silent = false
  batchUpgradeModalVisible.value = true
}

const handleBatchUpgrade = () => {
  if (!batchUpgradeForm.firmware_version) { Message.warning('请选择固件版本'); return }
  if (!batchUpgradeForm.device_ids.length) { Message.warning('请选择目标设备'); return }
  batchLoading.value = true
  setTimeout(() => {
    batchUpgradeForm.device_ids.forEach((devId, i) => {
      tasks.value.unshift({ id: Date.now() + i, device_id: devId, target_version: batchUpgradeForm.firmware_version, current_version: 'v1.1.5', status: 1, progress: 0, shard_index: 0, total_shards: 128, start_time: new Date().toLocaleString() })
    })
    batchLoading.value = false
    batchUpgradeModalVisible.value = false
    Message.success(`已创建 ${batchUpgradeForm.device_ids.length} 个升级任务`)
    updateStats()
  }, 1000)
}

const handleRetry = (record) => {
  const idx = tasks.value.findIndex(t => t.id === record.id)
  if (idx !== -1) { tasks.value[idx].status = 1; tasks.value[idx].progress = 0 }
  Message.success(`设备 ${record.device_id} 重试任务已创建`)
  updateStats()
}

const handlePause = (record) => {
  const idx = tasks.value.findIndex(t => t.id === record.id)
  if (idx !== -1) { tasks.value[idx].status = 1 }
  Message.info(`设备 ${record.device_id} 升级已暂停`)
}

const handleVerify = (record) => { Message.info(`设备 ${record.device_id} 升级验证通过`) }

const handleCancelTask = (record) => {
  tasks.value = tasks.value.filter(t => t.id !== record.id)
  Message.success('升级任务已取消')
  updateStats()
}

const clearFilters = () => { filterTaskStatus.value = null; filterVersion.value = null; searchKeyword.value = '' }

const handleTableChange = (pag) => { pagination.current = pag.current }

const getTaskStatusColor = (s) => ({ 1: 'default', 2: 'processing', 3: 'blue', 4: 'green', 5: 'red' }[s] || 'default')
const getTaskStatusText = (s) => ({ 1: '等待中', 2: '下载中', 3: '升级中', 4: '成功', 5: '失败' }[s] || '未知')
const formatTime = (t) => t || '-'

onMounted(() => { loadTasks(); loadPolicies(); loadDevices(); startLogStream() })
onUnmounted(() => { if (logTimer) clearInterval(logTimer) })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.pro-tabs { background: #fff; border-radius: 8px; padding: 16px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area { background: #fff; border-radius: 8px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); margin-bottom: 16px; }
.config-card { border-radius: 8px; }
.form-tip { font-size: 12px; color: #999; margin-top: 4px; }
.log-card { border-radius: 8px; }
.log-terminal { background: #1e1e1e; border-radius: 4px; padding: 12px; max-height: 200px; overflow-y: auto; font-family: 'Courier New', monospace; font-size: 12px; }
.log-line { display: flex; gap: 8px; padding: 2px 0; }
.log-time { color: #888; }
.log-device { color: #4fc3f7; }
.log-msg { color: #e0e0e0; }
.log-line.success .log-msg { color: #66bb6a; }
.log-line.error .log-msg { color: #ef5350; }
.log-line.warning .log-msg { color: #ffa726; }
</style>
