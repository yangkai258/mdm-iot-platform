<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>应用生态</a-breadcrumb-item>
      <a-breadcrumb-item>安装管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="已安装" :value="stats.installed">
            <template #prefix>📲</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="运行中" :value="stats.running" :value-style="{ color: '#0fc6c2' }">
            <template #prefix>▶️</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="异常应用" :value="stats.error" :value-style="{ color: '#f53f3f' }">
            <template #prefix>⚠️</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <a-statistic title="存储占用" :value="formatBytes(stats.storage_used)">
            <template #prefix>💾</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选区 -->
    <a-card class="section-card" style="margin-top: 16px">
      <a-space wrap>
        <a-select v-model="filterForm.status" placeholder="运行状态" style="width: 150px" allow-clear @change="loadInstallations">
          <a-option value="running">运行中</a-option>
          <a-option value="stopped">已停止</a-option>
          <a-option value="error">异常</a-option>
          <a-option value="updating">更新中</a-option>
        </a-select>
        <a-input-search v-model="filterForm.keyword" placeholder="搜索应用名称/设备ID" style="width: 240px" search-button @search="loadInstallations" @change="e => !e.target.value && loadInstallations()" />
        <a-button @click="loadInstallations">
          <template #icon><icon-refresh :spin="loading" /></template>
          刷新
        </a-button>
        <a-button type="primary" @click="showBatchUpdate">
          <template #icon><icon-upload /></template>
          批量更新
        </a-button>
      </a-space>
    </a-card>

    <!-- 安装列表 -->
    <a-card class="section-card" style="margin-top: 16px">
      <template #title>
        <a-checkbox v-model="selectAll" @change="handleSelectAll">全选</a-checkbox>
        <span style="margin-left: 12px; font-size: 13px; color: #86909c">已选择 {{ selected.length }} 项</span>
      </template>
      <a-spin :loading="loading" tip="加载中...">
        <div class="install-grid">
          <div v-for="inst in installations" :key="inst.install_id" class="install-card" :class="{ 'card-selected': selected.includes(inst.install_id) }">
            <div class="install-header">
              <a-checkbox :model-value="selected.includes(inst.install_id)" @change="() => toggleSelect(inst.install_id)" />
              <span class="install-icon">{{ inst.icon }}</span>
              <div class="install-info">
                <div class="install-name">{{ inst.app_name }}</div>
                <div class="install-version">{{ inst.version }}</div>
              </div>
              <a-tag :color="getStatusColor(inst.status)" size="small">{{ getStatusText(inst.status) }}</a-tag>
            </div>
            <div class="install-meta">
              <div class="meta-item"><icon-desktop :size="12" /> {{ inst.device_id }}</div>
              <div class="meta-item"><icon-calendar :size="12" /> {{ formatDate(inst.installed_at) }}</div>
              <div class="meta-item"><icon-schedule :size="12" /> {{ formatDate(inst.last_updated) }}</div>
            </div>
            <div class="install-stats">
              <div class="stat-item">
                <span class="stat-label">内存</span>
                <a-progress :percent="inst.memory_usage" size="small" :stroke-color="inst.memory_usage > 70 ? '#ff7d00' : '#165dff'" style="width: 80px" />
                <span class="stat-val">{{ inst.memory_usage }}%</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">CPU</span>
                <a-progress :percent="inst.cpu_usage" size="small" :stroke-color="inst.cpu_usage > 60 ? '#f53f3f' : '#0fc6c2'" style="width: 80px" />
                <span class="stat-val">{{ inst.cpu_usage }}%</span>
              </div>
            </div>
            <div class="install-actions">
              <a-button type="text" size="small" @click="handleStart(inst)" v-if="inst.status === 'stopped'">启动</a-button>
              <a-button type="text" size="small" @click="handleStop(inst)" v-if="inst.status === 'running'">停止</a-button>
              <a-button type="text" size="small" @click="handleUpdate(inst)" v-if="inst.update_available">更新</a-button>
              <a-button type="text" size="small" status="warning" @click="showLogs(inst)">日志</a-button>
              <a-button type="text" size="small" status="danger" @click="handleUninstall(inst)">卸载</a-button>
            </div>
            <div v-if="inst.update_available" class="update-tip">
              <icon-up-circle :size="12" /> 发现新版本 {{ inst.latest_version }}
            </div>
          </div>
        </div>
      </a-spin>
    </a-card>

    <!-- 日志弹窗 -->
    <a-drawer v-model:visible="logVisible" :title="`${currentInst?.app_name} 运行日志`" :width="720" unmountOnHide>
      <a-spin :loading="logLoading" tip="加载日志中...">
        <div class="log-container">
          <div v-for="(line, i) in logs" :key="i" class="log-line" :class="`log-${line.level}`">
            <span class="log-time">{{ line.time }}</span>
            <span class="log-level">{{ line.level.toUpperCase() }}</span>
            <span class="log-msg">{{ line.message }}</span>
          </div>
        </div>
      </a-spin>
      <div style="margin-top: 12px; text-align: right">
        <a-button @click="loadLogs" :loading="logLoading">
          <template #icon><icon-refresh /></template>
          刷新日志
        </a-button>
      </div>
    </a-drawer>

    <!-- 批量更新弹窗 -->
    <a-modal v-model="visible="batchVisible" title="批量更新" @before-ok="handleBatchUpdate" @cancel="batchVisible = false">
      <a-alert style="margin-bottom: 16px">检测到 {{ selected.length }} 个应用有可用更新，是否全部更新？</a-alert>
      <a-list size="small">
        <a-list-item v-for="id in selected" :key="id">
          {{ installations.find(i => i.install_id === id)?.app_name }} → {{ installations.find(i => i.install_id === id)?.latest_version }}
        </a-list-item>
      </a-list>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const logLoading = ref(false)
const logVisible = ref(false)
const batchVisible = ref(false)
const selectAll = ref(false)
const selected = ref([])
const currentInst = ref(null)

const stats = reactive({ installed: 0, running: 0, error: 0, storage_used: 0 })

const filterForm = reactive({ status: '', keyword: '' })

const installations = ref([])
const logs = ref([])

const getStatusColor = (s) => ({ running: 'green', stopped: 'gray', error: 'red', updating: 'blue' }[s] || 'gray')
const getStatusText = (s) => ({ running: '运行中', stopped: '已停止', error: '异常', updating: '更新中' }[s] || s)

const formatDate = (d) => d ? new Date(d).toLocaleDateString('zh-CN') : '-'

const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const gb = bytes / (1024 * 1024 * 1024)
  if (gb >= 1) return `${gb.toFixed(1)} GB`
  const mb = bytes / (1024 * 1024)
  return `${mb.toFixed(1)} MB`
}

const loadInstallations = async () => {
  loading.value = true
  await new Promise(r => setTimeout(r, 600))
  installations.value = [
    { install_id: 'INST-001', app_id: 'APP-001', app_name: '智能家居控制', icon: '🏠', device_id: 'DEV-1001', version: 'v2.3.0', latest_version: 'v2.3.1', status: 'running', memory_usage: 45, cpu_usage: 12, installed_at: '2026-02-01', last_updated: '2026-03-15', update_available: true },
    { install_id: 'INST-002', app_id: 'APP-002', app_name: '设备监控大师', icon: '📊', device_id: 'DEV-1002', version: 'v1.8.4', latest_version: 'v1.8.5', status: 'running', memory_usage: 68, cpu_usage: 25, installed_at: '2026-01-20', last_updated: '2026-03-10', update_available: true },
    { install_id: 'INST-003', app_id: 'APP-003', app_name: 'AI 助手 Pro', icon: '🤖', device_id: 'DEV-1001', version: 'v3.1.0', latest_version: 'v3.1.0', status: 'running', memory_usage: 82, cpu_usage: 38, installed_at: '2026-02-15', last_updated: '2026-03-18', update_available: false },
    { install_id: 'INST-004', app_id: 'APP-004', app_name: '网络测速工具', icon: '⚡', device_id: 'DEV-1003', version: 'v1.2.0', latest_version: 'v1.2.0', status: 'stopped', memory_usage: 0, cpu_usage: 0, installed_at: '2026-01-25', last_updated: '2026-01-25', update_available: false },
    { install_id: 'INST-005', app_id: 'APP-005', app_name: '家庭保险柜', icon: '🛡️', device_id: 'DEV-1002', version: 'v1.0.3', latest_version: 'v1.0.5', status: 'error', memory_usage: 95, cpu_usage: 0, installed_at: '2026-02-10', last_updated: '2026-02-28', update_available: true },
    { install_id: 'INST-006', app_id: 'APP-006', app_name: '儿童学习乐园', icon: '📚', device_id: 'DEV-1004', version: 'v2.0.1', latest_version: 'v2.0.1', status: 'running', memory_usage: 55, cpu_usage: 18, installed_at: '2026-03-01', last_updated: '2026-03-01', update_available: false },
    { install_id: 'INST-007', app_id: 'APP-001', app_name: '智能家居控制', icon: '🏠', device_id: 'DEV-1005', version: 'v2.2.0', latest_version: 'v2.3.1', status: 'updating', memory_usage: 30, cpu_usage: 8, installed_at: '2026-01-15', last_updated: '2026-03-20', update_available: true }
  ]
  stats.installed = installations.value.length
  stats.running = installations.value.filter(i => i.status === 'running').length
  stats.error = installations.value.filter(i => i.status === 'error').length
  stats.storage_used = installations.value.reduce((s, i) => s + (i.memory_usage * 1024 * 1024 * 10), 0)
  loading.value = false
}

const toggleSelect = (id) => {
  const idx = selected.value.indexOf(id)
  if (idx >= 0) selected.value.splice(idx, 1)
  else selected.value.push(id)
  selectAll.value = selected.value.length === installations.value.length
}

const handleSelectAll = (checked) => {
  selected.value = checked ? installations.value.map(i => i.install_id) : []
}

const handleStart = async (inst) => {
  inst.status = 'running'
  Message.success(`${inst.app_name} 已启动`)
}

const handleStop = async (inst) => {
  inst.status = 'stopped'
  inst.memory_usage = 0
  inst.cpu_usage = 0
  Message.warning(`${inst.app_name} 已停止`)
}

const handleUpdate = async (inst) => {
  inst.status = 'updating'
  await new Promise(r => setTimeout(r, 1500))
  inst.status = 'running'
  inst.version = inst.latest_version
  Message.success(`${inst.app_name} 更新成功`)
}

const handleUninstall = async (inst) => {
  await new Promise(r => setTimeout(r, 500))
  installations.value = installations.value.filter(i => i.install_id !== inst.install_id)
  Message.warning(`${inst.app_name} 已卸载`)
  loadInstallations()
}

const showLogs = async (inst) => {
  currentInst.value = inst
  logVisible.value = true
  await loadLogs()
}

const loadLogs = async () => {
  logLoading.value = true
  await new Promise(r => setTimeout(r, 500))
  logs.value = [
    { time: '2026-03-22 21:00:01', level: 'info', message: `[${currentInst.value?.app_name}] 服务启动成功` },
    { time: '2026-03-22 21:00:02', level: 'info', message: '正在连接设备... 连接成功' },
    { time: '2026-03-22 21:00:05', level: 'info', message: '心跳上报: device_id=' + currentInst.value?.device_id },
    { time: '2026-03-22 21:00:10', level: 'warn', message: '内存使用率偏高: ' + currentInst.value?.memory_usage + '%' },
    { time: '2026-03-22 21:00:15', level: 'info', message: '数据同步完成，耗时 120ms' },
    { time: '2026-03-22 21:00:20', level: 'info', message: '定时任务执行成功' },
    { time: '2026-03-22 21:00:25', level: 'error', message: '网络请求超时，已自动重试（1/3）' },
    { time: '2026-03-22 21:00:28', level: 'info', message: '重试成功，数据上报完成' }
  ]
  logLoading.value = false
}

const showBatchUpdate = () => {
  if (!selected.value.length) {
    Message.warning('请先选择要更新的应用')
    return
  }
  batchVisible.value = true
}

const handleBatchUpdate = async (done) => {
  await new Promise(r => setTimeout(r, 1000))
  Message.success('批量更新任务已启动')
  batchVisible.value = false
  done(true)
}

loadInstallations()
</script>

<style scoped>
.page-container { padding: 20px; }
.breadcrumb { margin-bottom: 12px; }
.stat-card { text-align: center; }
.section-card { margin-bottom: 16px; }
.install-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(360px, 1fr)); gap: 16px; }
.install-card { border: 1px solid #e5e6e8; border-radius: 8px; padding: 16px; background: #fff; transition: box-shadow 0.2s; position: relative; }
.install-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); }
.install-card.card-selected { border-color: #165dff; background: #f2f3f5; }
.install-header { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.install-icon { font-size: 28px; }
.install-info { flex: 1; }
.install-name { font-weight: 600; font-size: 14px; color: #1d2129; }
.install-version { font-size: 12px; color: #86909c; }
.install-meta { display: flex; gap: 16px; margin-bottom: 10px; }
.meta-item { display: flex; align-items: center; gap: 4px; font-size: 12px; color: #86909c; }
.install-stats { display: flex; gap: 24px; margin-bottom: 10px; }
.stat-item { display: flex; align-items: center; gap: 6px; }
.stat-label { font-size: 12px; color: #86909c; width: 28px; }
.stat-val { font-size: 12px; color: #4e5969; }
.install-actions { display: flex; gap: 4px; flex-wrap: wrap; }
.update-tip { display: flex; align-items: center; gap: 4px; font-size: 12px; color: #165dff; margin-top: 8px; }
.log-container { background: #1d1d1d; border-radius: 4px; padding: 12px; max-height: 500px; overflow-y: auto; font-family: 'Courier New', monospace; font-size: 12px; }
.log-line { display: flex; gap: 8px; margin-bottom: 4px; line-height: 1.5; }
.log-time { color: #888; }
.log-level { width: 50px; font-weight: 600; }
.log-info .log-level { color: #4facfe; }
.log-warn .log-level { color: #ffc53d; }
.log-error .log-level { color: #f53f3f; }
.log-msg { color: #ddd; flex: 1; }
</style>
