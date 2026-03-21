<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>告警管理</a-breadcrumb-item>
      <a-breadcrumb-item>告警列表</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic title="告警总数" :value="stats.total" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card stat-pending">
          <a-statistic title="待处理" :value="stats.pending" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card stat-confirmed">
          <a-statistic title="已确认" :value="stats.confirmed" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card stat-resolved">
          <a-statistic title="已解决" :value="stats.resolved" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <a-space wrap>
        <a-input-search
          v-model="filters.device_id"
          placeholder="设备ID"
          style="width: 160px"
          search-button
          @search="loadAlerts"
          allow-clear
        />
        <a-select
          v-model="filters.status"
          placeholder="处理状态"
          style="width: 130px"
          allow-clear
          @change="loadAlerts"
        >
          <a-option :value="1">未处理</a-option>
          <a-option :value="2">已确认</a-option>
          <a-option :value="3">已解决</a-option>
          <a-option :value="4">已忽略</a-option>
        </a-select>
        <a-select
          v-model="filters.severity"
          placeholder="严重程度"
          style="width: 130px"
          allow-clear
          @change="loadAlerts"
        >
          <a-option :value="1">低</a-option>
          <a-option :value="2">中</a-option>
          <a-option :value="3">高</a-option>
          <a-option :value="4">严重</a-option>
        </a-select>
        <a-select
          v-model="filters.alert_type"
          placeholder="告警类型"
          style="width: 150px"
          allow-clear
          @change="loadAlerts"
        >
          <a-option value="battery_low">电量过低</a-option>
          <a-option value="offline">设备离线</a-option>
          <a-option value="temperature_high">温度过高</a-option>
          <a-option value="geofence_violation">地理围栏违规</a-option>
          <a-option value="jailbreak_detected">越狱/Root检测</a-option>
          <a-option value="storage_low">存储空间不足</a-option>
          <a-option value="network_unreachable">网络不可达</a-option>
        </a-select>
        <a-button @click="clearFilters">重置</a-button>
      </a-space>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button
          type="primary"
          :disabled="selectedKeys.length === 0"
          @click="batchConfirm"
        >
          <template #icon><icon-check /></template>
          批量确认 ({{ selectedKeys.length }})
        </a-button>
        <a-button
          :disabled="selectedKeys.length === 0"
          @click="batchResolve"
        >
          <template #icon><icon-check-circle /></template>
          批量解决
        </a-button>
        <a-button @click="loadAlerts">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 数据表格 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="alerts"
        :loading="loading"
        :pagination="{ pageSize: 20, total: pagination.total }"
        :row-selection="{ type: 'checkbox', showCheckedAll: true, onlyCurrent: false }"
        v-model:selected-keys="selectedKeys"
        row-key="id"
        @page-change="handlePageChange"
      >
        <template #id="{ record }">
          <span class="alert-id">#{{ record.id }}</span>
        </template>

        <template #message="{ record }">
          <a-tooltip :content="record.message">
            <span class="alert-message">{{ record.message }}</span>
          </a-tooltip>
        </template>

        <template #device_id="{ record }">
          <a-button type="text" size="small" @click="openDeviceDetail(record.device_id)">
            {{ record.device_id }}
          </a-button>
        </template>

        <template #alert_type="{ record }">
          <a-tag :color="getAlertTypeColor(record.alert_type)">
            {{ getAlertTypeText(record.alert_type) }}
          </a-tag>
        </template>

        <template #severity="{ record }">
          <a-tag :color="getSeverityColor(record.severity)">
            {{ getSeverityText(record.severity) }}
          </a-tag>
        </template>

        <template #status="{ record }">
          <a-badge
            :status="getStatusBadge(record.status)"
            :text="getStatusText(record.status)"
          />
        </template>

        <template #trigger_info="{ record }">
          <span class="trigger-info">
            {{ record.trigger_val ?? '-' }} / {{ record.threshold ?? '-' }}
          </span>
        </template>

        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>

        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button
              v-if="record.status === 1"
              type="text" size="small"
              @click="handleConfirm(record)"
            >确认</a-button>
            <a-button
              v-if="record.status !== 3 && record.status !== 4"
              type="text" size="small"
              @click="handleResolve(record)"
            >解决</a-button>
            <a-popconfirm content="确定删除该告警记录？" @ok="handleDelete(record.id)">
              <a-button type="text" size="small" status="danger">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </div>

    <!-- 告警详情抽屉 -->
    <a-drawer
      v-model:visible="detailVisible"
      :title="`告警详情 #${currentAlert?.id || ''}`"
      :width="560"
      ok-text="确认"
      cancel-text="关闭"
      @ok="currentAlert && currentAlert.status === 1 ? handleConfirm(currentAlert) : (detailVisible = false)"
    >
      <template v-if="currentAlert">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="告警ID">#{{ currentAlert.id }}</a-descriptions-item>
          <a-descriptions-item label="告警内容">{{ currentAlert.message }}</a-descriptions-item>
          <a-descriptions-item label="设备ID">
            <a-button type="text" size="small" @click="openDeviceDetail(currentAlert.device_id)">
              {{ currentAlert.device_id }}
            </a-button>
          </a-descriptions-item>
          <a-descriptions-item label="告警类型">
            <a-tag :color="getAlertTypeColor(currentAlert.alert_type)">
              {{ getAlertTypeText(currentAlert.alert_type) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="严重程度">
            <a-tag :color="getSeverityColor(currentAlert.severity)">
              {{ getSeverityText(currentAlert.severity) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="处理状态">
            <a-badge
              :status="getStatusBadge(currentAlert.status)"
              :text="getStatusText(currentAlert.status)"
            />
          </a-descriptions-item>
          <a-descriptions-item label="触发值 / 阈值">
            {{ currentAlert.trigger_val ?? '-' }} / {{ currentAlert.threshold ?? '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="触发时间">
            {{ formatTime(currentAlert.created_at) }}
          </a-descriptions-item>
          <a-descriptions-item v-if="currentAlert.confirmed_at" label="确认时间">
            {{ formatTime(currentAlert.confirmed_at) }} ({{ currentAlert.confirmed_by || '系统' }})
          </a-descriptions-item>
          <a-descriptions-item v-if="currentAlert.resolved_at" label="解决时间">
            {{ formatTime(currentAlert.resolved_at) }} ({{ currentAlert.resolved_by || '系统' }})
          </a-descriptions-item>
          <a-descriptions-item v-if="currentAlert.extra_data" label="附加数据">
            <pre style="margin:0;font-size:12px">{{ formatExtraData(currentAlert.extra_data) }}</pre>
          </a-descriptions-item>
        </a-descriptions>

        <a-divider>操作</a-divider>
        <a-space vertical>
          <a-button
            v-if="currentAlert.status === 1"
            type="primary"
            long
            @click="handleConfirm(currentAlert)"
          >确认告警</a-button>
          <a-button
            v-if="currentAlert.status !== 3 && currentAlert.status !== 4"
            type="normal"
            long
            @click="handleResolve(currentAlert)"
          >标记为已解决</a-button>
          <a-popconfirm content="确定忽略该告警？" @ok="handleIgnore(currentAlert.id)">
            <a-button type="outline" long status="warning">忽略告警</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as api from '@/api/alerts'

const loading = ref(false)
const alerts = ref([])
const selectedKeys = ref([])
const detailVisible = ref(false)
const currentAlert = ref(null)
const pagination = reactive({ total: 0 })

const filters = reactive({
  device_id: '',
  status: undefined,
  severity: undefined,
  alert_type: undefined
})

const stats = reactive({ total: 0, pending: 0, confirmed: 0, resolved: 0 })

const columns = [
  { title: 'ID', slotName: 'id', width: 80 },
  { title: '告警内容', slotName: 'message', minWidth: 200 },
  { title: '设备', slotName: 'device_id', width: 160 },
  { title: '类型', slotName: 'alert_type', width: 130 },
  { title: '严重程度', slotName: 'severity', width: 100 },
  { title: '触发值/阈值', slotName: 'trigger_info', width: 130 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '触发时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const getSeverityColor = (s) => ({ 1: 'green', 2: 'arcoblue', 3: 'orange', 4: 'red' })[s] || 'gray'
const getSeverityText = (s) => ({ 1: '低', 2: '中', 3: '高', 4: '严重' })[s] || '未知'
const getAlertTypeColor = (t) => ({
  battery_low: 'green', offline: 'orange', temperature_high: 'red',
  geofence_violation: 'purple', jailbreak_detected: 'red',
  storage_low: 'orange', network_unreachable: 'gray'
})[t] || 'arcoblue'
const getAlertTypeText = (t) => ({
  battery_low: '电量过低', offline: '设备离线', temperature_high: '温度过高',
  geofence_violation: '地理围栏', jailbreak_detected: '越狱/Root',
  storage_low: '存储不足', network_unreachable: '网络不可达'
})[t] || t
const getStatusBadge = (s) => ({ 1: 'error', 2: 'warning', 3: 'success', 4: 'default' })[s] || 'default'
const getStatusText = (s) => ({ 1: '未处理', 2: '已确认', 3: '已解决', 4: '已忽略' })[s] || '未知'

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function formatExtraData(data) {
  try {
    return JSON.stringify(JSON.parse(data), null, 2)
  } catch {
    return data
  }
}

function computeStats(list) {
  stats.total = list.length
  stats.pending = list.filter(a => a.status === 1).length
  stats.confirmed = list.filter(a => a.status === 2).length
  stats.resolved = list.filter(a => a.status === 3 || a.status === 4).length
}

async function loadAlerts() {
  loading.value = true
  selectedKeys.value = []
  try {
    const params = {}
    if (filters.device_id) params.device_id = filters.device_id
    if (filters.status !== undefined) params.status = filters.status
    if (filters.severity !== undefined) params.severity = filters.severity
    if (filters.alert_type) params.alert_type = filters.alert_type

    const data = await api.getAlerts(params)
    alerts.value = data.data?.list || []
    computeStats(alerts.value)
  } catch (e) {
    Message.error('加载告警列表失败')
  } finally {
    loading.value = false
  }
}

function clearFilters() {
  Object.assign(filters, { device_id: '', status: undefined, severity: undefined, alert_type: undefined })
  loadAlerts()
}

function handlePageChange(page) {
  // backend doesn't support pagination offset for now, just reload
  loadAlerts()
}

function openDetail(record) {
  currentAlert.value = record
  detailVisible.value = true
}

function openDeviceDetail(deviceId) {
  window.location.hash = `#/device/${deviceId}`
}

async function handleConfirm(record) {
  try {
    await api.confirmAlert(record.id)
    Message.success('告警已确认')
    record.status = 2
    computeStats(alerts.value)
    if (currentAlert.value?.id === record.id) currentAlert.value.status = 2
  } catch (e) {
    Message.error('确认失败')
  }
}

async function handleResolve(record) {
  try {
    await api.resolveAlert(record.id)
    Message.success('告警已解决')
    record.status = 3
    computeStats(alerts.value)
    if (currentAlert.value?.id === record.id) currentAlert.value.status = 3
  } catch (e) {
    Message.error('解决失败')
  }
}

async function handleIgnore(id) {
  try {
    await api.ignoreAlert(id)
    Message.success('告警已忽略')
    loadAlerts()
  } catch (e) {
    Message.error('忽略失败')
  }
}

async function handleDelete(id) {
  try {
    alerts.value = alerts.value.filter(a => a.id !== id)
    Message.success('删除成功')
    computeStats(alerts.value)
  } catch (e) {
    Message.error('删除失败')
  }
}

async function batchConfirm() {
  if (!selectedKeys.value.length) return
  try {
    await api.batchConfirmAlerts(selectedKeys.value)
    Message.success(`已确认 ${selectedKeys.value.length} 条告警`)
    selectedKeys.value = []
    loadAlerts()
  } catch (e) {
    Message.error('批量确认失败')
  }
}

async function batchResolve() {
  if (!selectedKeys.value.length) return
  try {
    await api.batchResolveAlerts(selectedKeys.value)
    Message.success(`已解决 ${selectedKeys.value.length} 条告警`)
    selectedKeys.value = []
    loadAlerts()
  } catch (e) {
    Message.error('批量解决失败')
  }
}

onMounted(() => { loadAlerts() })
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { border-radius: 8px; text-align: center; }
.alert-id { font-family: monospace; color: #888; font-size: 12px; }
.alert-message { font-size: 13px; }
.trigger-info { font-family: monospace; font-size: 12px; }
.filter-bar { margin-bottom: 12px; background: #fff; padding: 12px 16px; border-radius: 8px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
