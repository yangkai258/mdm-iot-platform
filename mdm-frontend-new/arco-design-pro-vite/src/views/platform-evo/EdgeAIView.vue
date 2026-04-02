<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>平台演进</a-breadcrumb-item>
      <a-breadcrumb-item>端侧 AI 推理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏 -->
    <div class="pro-search-bar">
      <a-space>
        <a-input v-model="searchKey" placeholder="搜索模型名称 / 版本" allow-clear style="width: 240px" @press-enter="loadModels" />
        <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 140px">
          <a-option value="deployed">已部署</a-option>
          <a-option value="idle">空闲</a-option>
          <a-option value="error">异常</a-option>
        </a-select>
        <a-select v-model="filterDevice" placeholder="所属设备" allow-clear style="width: 160px">
          <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showDeployModal = true">
          <template #icon><icon-upload /></template>
          部署模型
        </a-button>
        <a-button @click="loadModels">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 统计卡片 -->
    <div class="metric-cards">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="在线模型数" :value="stats.online_models" :loading="statsLoading" animation />
        </a-col>
        <a-col :span="6">
          <a-statistic title="今日推理次数" :value="stats.today_inferences" :loading="statsLoading" animation>
            <template #suffix>次</template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="平均延迟" :value="stats.avg_latency_ms" :loading="statsLoading" animation>
            <template #suffix>ms</template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="设备数" :value="stats.active_devices" :loading="statsLoading" animation />
        </a-col>
      </a-row>
    </div>

    <!-- 模型列表 -->
    <a-table :data="models" :loading="loading" :pagination="{ total: total, current: page, pageSize: pageSize, showTotal: true, showPageSize: true }" @page-change="onPageChange" @page-size-change="onPageSizeChange" row-key="id">
      <template #columns>
        <a-table-column title="模型名称" data-index="name" :width="180">
          <template #cell="{ record }">
            <a-link @click="viewDetail(record)">{{ record.name }}</a-link>
          </template>
      </a-table>
        </a-table-column>
        <a-table-column title="版本" data-index="version" :width="100" />
        <a-table-column title="算法" data-index="algorithm" :width="120" />
        <a-table-column title="精度" data-index="precision" :width="100">
          <template #cell="{ record }">{{ record.precision || '-' }}</template>
        </a-table-column>
        <a-table-column title="参数量" data-index="params" :width="100">
          <template #cell="{ record }">{{ record.params || '-' }}</template>
        </a-table-column>
        <a-table-column title="状态" data-index="status" :width="100">
          <template #cell="{ record }">
            <a-tag :color="statusColor(record.status)">{{ statusText(record.status) }}</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="部署设备" data-index="device_name" :width="140">
          <template #cell="{ record }">{{ record.device_name || '-' }}</template>
        </a-table-column>
        <a-table-column title="今日推理" data-index="today_inferences" :width="100">
          <template #cell="{ record }">{{ record.today_inferences || 0 }} 次</template>
        </a-table-column>
        <a-table-column title="最近推理" data-index="last_inference_at" :width="160">
          <template #cell="{ record }">{{ formatTime(record.last_inference_at) }}</template>
        </a-table-column>
        <a-table-column title="操作" :width="160" fixed="right">
          <template #cell="{ record }">
            <a-space>
              <a-button size="small" @click="viewDetail(record)">详情</a-button>
              <a-button v-if="record.status === 'idle'" size="small" type="primary" @click="deployModel(record)">部署</a-button>
              <a-button v-if="record.status === 'deployed'" size="small" status="danger" @click="undeployModel(record)">卸载</a-button>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 性能监控区域 -->
    <a-divider>推理性能监控</a-divider>
    <a-row :gutter="16">
      <a-col :span="24">
        <a-tabs v-model="activeTab">
          <a-tab-pane key="latency" title="推理延迟">
            <a-space style="margin-bottom: 12px">
              <a-radio-group v-model="metricsTimeRange" type="button" @change="loadMetrics">
                <a-radio value="1h">最近1小时</a-radio>
                <a-radio value="24h">最近24小时</a-radio>
                <a-radio value="7d">最近7天</a-radio>
              </a-radio-group>
            </a-space>
            <a-table :data="inferenceLogs" :loading="logsLoading" :pagination="{ total: logsTotal, current: logsPage, pageSize: logsPageSize, showTotal: true }" @page-change="onLogsPageChange" row-key="id" size="small">
              <template #columns>
                <a-table-column title="时间" data-index="created_at" :width="180">
                  <template #cell="{ record }">{{ formatTime(record.created_at) }}</template>
      </a-table>
                </a-table-column>
                <a-table-column title="设备" data-index="device_name" :width="140" />
                <a-table-column title="模型" data-index="model_name" :width="160" />
                <a-table-column title="输入大小" data-index="input_size" :width="100">
                  <template #cell="{ record }">{{ record.input_size || '-' }} KB</template>
                </a-table-column>
                <a-table-column title="输出大小" data-index="output_size" :width="100">
                  <template #cell="{ record }">{{ record.output_size || '-' }} KB</template>
                </a-table-column>
                <a-table-column title="延迟" data-index="latency_ms" :width="100">
                  <template #cell="{ record }">
                    <span :class="latencyClass(record.latency_ms)">{{ record.latency_ms }} ms</span>
                  </template>
                </a-table-column>
                <a-table-column title="状态" data-index="status" :width="80">
                  <template #cell="{ record }">
                    <a-tag :color="record.status === 'success' ? 'green' : 'red'" size="small">
                      {{ record.status === 'success' ? '成功' : '失败' }}
                    </a-tag>
                  </template>
                </a-table-column>
                <a-table-column title="错误信息" data-index="error_message">
                  <template #cell="{ record }">
                    <span class="text-danger">{{ record.error_message || '-' }}</span>
                  </template>
                </a-table-column>
              </template>
            </a-table>
          </a-tab-pane>
          <a-tab-pane key="devices" title="设备推理状态">
            <a-table :data="deviceStatuses" :loading="deviceStatusLoading" row-key="device_id">
              <template #columns>
                <a-table-column title="设备名称" data-index="device_name" :width="160" />
                <a-table-column title="在线状态" data-index="online" :width="100">
                  <template #cell="{ record }">
                    <a-badge :status="record.online ? 'success' : 'default'" :text="record.online ? '在线' : '离线'" />
                  </template>
      </a-table>
                </a-table-column>
                <a-table-column title="运行模型" data-index="running_model" :width="160" />
                <a-table-column title="推理次数" data-index="inference_count" :width="100" />
                <a-table-column title="平均延迟" data-index="avg_latency_ms" :width="100">
                  <template #cell="{ record }">{{ record.avg_latency_ms || '-' }} ms</template>
                </a-table-column>
                <a-table-column title="CPU 使用" data-index="cpu_percent" :width="140">
                  <template #cell="{ record }">
                    <a-progress :percent="record.cpu_percent || 0" :color="cpuColor(record.cpu_percent)" size="small" />
                  </template>
                </a-table-column>
                <a-table-column title="内存使用" data-index="mem_percent" :width="140">
                  <template #cell="{ record }">
                    <a-progress :percent="record.mem_percent || 0" color="arcoblue" size="small" />
                  </template>
                </a-table-column>
                <a-table-column title="温度" data-index="temperature_c" :width="100">
                  <template #cell="{ record }">{{ record.temperature_c ? record.temperature_c + '°C' : '-' }}</template>
                </a-table-column>
              </template>
            </a-table>
          </a-tab-pane>
        </a-tabs>
      </a-col>
    </a-row>

    <!-- 部署弹窗 -->
    <a-modal v-model:visible="showDeployModal" title="部署模型" :width="480" @before-ok="handleDeploy" @cancel="showDeployModal = false">
      <a-form :model="deployForm" layout="vertical">
        <a-form-item label="选择模型" required>
          <a-select v-model="deployForm.model_id" placeholder="请选择模型" @change="onModelSelect">
            <a-option v-for="m in models" :key="m.id" :value="m.id">{{ m.name }} ({{ m.version }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标设备" required>
          <a-select v-model="deployForm.device_id" placeholder="请选择设备">
            <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="推理优先级">
          <a-select v-model="deployForm.priority" placeholder="默认优先级">
            <a-option value="low">低</a-option>
            <a-option value="normal">普通</a-option>
            <a-option value="high">高</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情弹窗 -->
    <a-drawer v-model:visible="showDetailDrawer" :title="`模型详情: ${currentModel?.name}`" :width="560" @close="showDetailDrawer = false">
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="模型名称">{{ currentModel?.name }}</a-descriptions-item>
        <a-descriptions-item label="版本">{{ currentModel?.version }}</a-descriptions-item>
        <a-descriptions-item label="算法">{{ currentModel?.algorithm }}</a-descriptions-item>
        <a-descriptions-item label="精度">{{ currentModel?.precision || '-' }}</a-descriptions-item>
        <a-descriptions-item label="参数量">{{ currentModel?.params || '-' }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="statusColor(currentModel?.status)">{{ statusText(currentModel?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="部署设备">{{ currentModel?.device_name || '-' }}</a-descriptions-item>
        <a-descriptions-item label="今日推理">{{ currentModel?.today_inferences || 0 }} 次</a-descriptions-item>
        <a-descriptions-item label="总推理次数">{{ currentModel?.total_inferences || 0 }} 次</a-descriptions-item>
        <a-descriptions-item label="平均延迟">{{ currentModel?.avg_latency_ms || '-' }} ms</a-descriptions-item>
      </a-descriptions>
      <a-divider>推理统计</a-divider>
      <a-row :gutter="16">
        <a-col :span="12">
          <a-statistic title="今日推理" :value="currentModel?.today_inferences || 0" />
        </a-col>
        <a-col :span="12">
          <a-statistic title="平均延迟" :value="currentModel?.avg_latency_ms || 0" suffix="ms" />
        </a-col>
      </a-row>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

async function apiRequest(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${getToken()}`,
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  })
  const data = await res.json()
  if (data.code !== 0 && data.code !== 200) {
    throw new Error(data.message || '请求失败')
  }
  return data
}

// 列表数据
const models = ref([])
const devices = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const statsLoading = ref(false)
const logsLoading = ref(false)
const deviceStatusLoading = ref(false)
const searchKey = ref('')
const filterStatus = ref('')
const filterDevice = ref('')

// 统计数据
const stats = ref({ online_models: 0, today_inferences: 0, avg_latency_ms: 0, active_devices: 0 })

// 性能监控
const activeTab = ref('latency')
const metricsTimeRange = ref('24h')
const inferenceLogs = ref([])
const logsTotal = ref(0)
const logsPage = ref(1)
const logsPageSize = ref(10)
const deviceStatuses = ref([])

// 部署弹窗
const showDeployModal = ref(false)
const deployForm = reactive({ model_id: '', device_id: '', priority: 'normal' })

// 详情抽屉
const showDetailDrawer = ref(false)
const currentModel = ref(null)

function statusColor(s) {
  return { deployed: 'green', idle: 'gray', error: 'red' }[s] || 'gray'
}
function statusText(s) {
  return { deployed: '已部署', idle: '空闲', error: '异常' }[s] || s
}
function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}
function latencyClass(ms) {
  if (ms < 50) return 'text-success'
  if (ms < 200) return 'text-warning'
  return 'text-danger'
}
function cpuColor(p) {
  if (p < 50) return 'green'
  if (p < 80) return 'orange'
  return 'red'
}

async function loadModels() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (searchKey.value) params.keyword = searchKey.value
    if (filterStatus.value) params.status = filterStatus.value
    if (filterDevice.value) params.device_id = filterDevice.value
    const qs = new URLSearchParams(params).toString()
    const res = await apiRequest(`${API_BASE}/edge-ai/models?${qs}`)
    models.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    Message.error('加载模型列表失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

async function loadDevices() {
  try {
    const res = await apiRequest(`${API_BASE}/edge-ai/devices?page_size=100`)
    devices.value = res.data?.list || []
  } catch (e) {
    // ignore
  }
}

async function loadStats() {
  statsLoading.value = true
  try {
    const res = await apiRequest(`${API_BASE}/edge-ai/stats`)
    stats.value = res.data || {}
  } catch (e) {
    // ignore
  } finally {
    statsLoading.value = false
  }
}

async function loadMetrics() {
  if (activeTab.value === 'latency') {
    await loadInferenceLogs()
  } else {
    await loadDeviceStatuses()
  }
}

async function loadInferenceLogs() {
  logsLoading.value = true
  try {
    const params = { page: logsPage.value, page_size: logsPageSize.value, time_range: metricsTimeRange.value }
    const qs = new URLSearchParams(params).toString()
    const res = await apiRequest(`${API_BASE}/edge-ai/logs?${qs}`)
    inferenceLogs.value = res.data?.list || []
    logsTotal.value = res.data?.total || 0
  } catch (e) {
    Message.error('加载推理日志失败')
  } finally {
    logsLoading.value = false
  }
}

async function loadDeviceStatuses() {
  deviceStatusLoading.value = true
  try {
    const res = await apiRequest(`${API_BASE}/edge-ai/devices/status`)
    deviceStatuses.value = res.data?.list || []
  } catch (e) {
    Message.error('加载设备状态失败')
  } finally {
    deviceStatusLoading.value = false
  }
}

function onModelSelect(val) {
  deployForm.model_id = val
}

async function handleDeploy(done) {
  if (!deployForm.model_id || !deployForm.device_id) {
    Message.warning('请选择模型和设备')
    done(false)
    return
  }
  try {
    await apiRequest(`${API_BASE}/edge-ai/models/${deployForm.model_id}/deploy`, {
      method: 'POST',
      body: JSON.stringify({ device_id: deployForm.device_id, priority: deployForm.priority })
    })
    Message.success('部署成功')
    showDeployModal.value = false
    deployForm.model_id = ''
    deployForm.device_id = ''
    deployForm.priority = 'normal'
    loadModels()
    loadStats()
    done(true)
  } catch (e) {
    Message.error('部署失败: ' + e.message)
    done(false)
  }
}

async function undeployModel(record) {
  try {
    await apiRequest(`${API_BASE}/edge-ai/models/${record.id}/undeploy`, {
      method: 'POST',
      body: JSON.stringify({ device_id: record.device_id })
    })
    Message.success('卸载成功')
    loadModels()
    loadStats()
  } catch (e) {
    Message.error('卸载失败: ' + e.message)
  }
}

function viewDetail(record) {
  currentModel.value = record
  showDetailDrawer.value = true
}

function deployModel(record) {
  deployForm.model_id = record.id
  showDeployModal.value = true
}

function onPageChange(p) {
  page.value = p
  loadModels()
}
function onPageSizeChange(s) {
  pageSize.value = s
  loadModels()
}
function onLogsPageChange(p) {
  logsPage.value = p
  loadInferenceLogs()
}

// init
loadModels()
loadDevices()
loadStats()
</script>

<style scoped>
.text-success { color: var(--color-success); }
.text-warning { color: var(--color-warning); }
.text-danger { color: var(--color-danger); }
</style>
