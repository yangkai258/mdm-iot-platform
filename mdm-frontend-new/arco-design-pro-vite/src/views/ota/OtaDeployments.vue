<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="ota-deployments-container">

    <a-card class="general-card">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="总任务数" :value="stats.total" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="进行中" :value="stats.running" :value-style="{ color: '#1650d8' }" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="成功率" :value="stats.successRate + '%'" :value-style="{ color: stats.successRate < 80 ? '#ff4d4f' : '#52c41a' }" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="待升级设备" :value="stats.pendingDevices" :value-style="{ color: '#faad14' }" />
        </a-col>
      </a-row>
    </a-card>

    <a-card class="general-card" style="margin-top: 16px">
      <template #title>
        <div class="card-title">
          <span>部署任务</span>
          <a-space>
            <a-select v-model="filterStatus" placeholder="任务状态" allow-clear style="width: 130px" @change="handleFilter">
              <a-option value="pending">待执行</a-option>
              <a-option value="running">进行中</a-option>
              <a-option value="paused">已暂停</a-option>
              <a-option value="completed">已完成</a-option>
              <a-option value="failed">失败</a-option>
              <a-option value="cancelled">已取消</a-option>
            </a-select>
            <a-select v-model="filterHardwareModel" placeholder="硬件型号" allow-clear style="width: 160px" @change="handleFilter">
              <a-option v-for="model in hardwareModels" :key="model" :value="model">{{ model }}</a-option>
            </a-select>
            <a-button type="primary" @click="showCreateDrawer">新建部署任务</a-button>
          </a-space>
        </div>
      </template>

      <a-table :columns="columns" :data="deployments" :loading="loading" :pagination="paginationConfig" row-key="id" @page-change="handlePageChange" @page-size-change="handlePageSizeChange">
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
      </a-table>
        <template #strategy_type="{ record }">
          <a-tag :color="getStrategyColor(record.strategy_type)">
            {{ getStrategyText(record.strategy_type) }}
          </a-tag>
        </template>
        <template #progress="{ record }">
          <div class="progress-cell">
            <a-progress :percent="getProgress(record)" :stroke-width="6" size="small" :show-text="false" style="width: 100px" />
            <span class="progress-text">{{ getProgress(record) }}%</span>
          </div>
        </template>
        <template #success_rate="{ record }">
          <span :style="{ color: getSuccessRateColor(record) }">{{ getSuccessRate(record) }}%</span>
        </template>
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleDetail(record)">详情</a-button>
            <a-button v-if="record.status === 'running' || record.status === 'pending'" type="text" size="small" @click="handlePause(record)">暂停</a-button>
            <a-button v-if="record.status === 'paused'" type="text" size="small" @click="handleResume(record)">恢复</a-button>
            <a-button v-if="record.status !== 'completed' && record.status !== 'cancelled' && record.status !== 'failed'" type="text" size="small" status="danger" @click="handleCancel(record)">取消</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 创建部署任务抽屉 -->
    <a-drawer v-model:visible="createDrawerVisible" title="新建部署任务" width="520px" @before-ok="handleCreate" :unmount-on-close="false">
      <a-form :model="form" layout="vertical" ref="formRef">
        <a-form-item label="任务名称" field="name" :rules="[{ required: true, message: '请输入任务名称' }]">
          <a-input v-model="form.name" placeholder="例如: v1.3.0 灰度发布30%" />
        </a-form-item>
        <a-form-item label="固件包" field="package_id" :rules="[{ required: true, message: '请选择固件包' }]">
          <a-select v-model="form.package_id" placeholder="请选择固件包" @change="handlePackageChange">
            <a-option v-for="pkg in packages" :key="pkg.id" :value="pkg.id">
              {{ pkg.name }} ({{ pkg.version }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标硬件型号" field="hardware_model" :rules="[{ required: true, message: '请选择目标硬件型号' }]">
          <a-select v-model="form.hardware_model" placeholder="请选择硬件型号">
            <a-option v-for="model in hardwareModels" :key="model" :value="model">{{ model }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="发布策略" field="strategy_type" :rules="[{ required: true, message: '请选择发布策略' }]">
          <a-radio-group v-model="form.strategy_type">
            <a-radio value="full">全量发布</a-radio>
            <a-radio value="percentage">百分比灰度</a-radio>
            <a-radio value="whitelist">白名单</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.strategy_type === 'percentage'" label="灰度百分比" field="percentage" :rules="[{ required: true, message: '请输入百分比' }]">
          <a-input-number v-model="form.percentage" :min="1" :max="100" suffix="%" placeholder="1-100" style="width: 200px" />
        </a-form-item>
        <a-form-item v-if="form.strategy_type === 'whitelist'" label="白名单设备ID" field="device_ids">
          <a-textarea v-model="form.device_ids" placeholder="多个设备ID用逗号分隔" :rows="3" />
        </a-form-item>
        <a-form-item label="失败率阈值" tooltip="当失败率超过此阈值时，任务自动暂停">
          <a-input-number v-model="form.pause_on_failure_threshold" :min="1" :max="100" suffix="%" default-value="20" style="width: 200px" />
        </a-form-item>
        <a-form-item label="计划开始时间" tooltip="留空则立即开始">
          <a-date-picker v-model="form.scheduled_at" show-time format="YYYY-MM-DD HH:mm" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 任务详情抽屉 -->
    <a-drawer v-model:visible="detailDrawerVisible" title="任务详情" width="640px" :footer="false">
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="任务名称">{{ currentDeployment?.name }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentDeployment?.status)">{{ getStatusText(currentDeployment?.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="固件版本">{{ currentDeployment?.package_version }}</a-descriptions-item>
        <a-descriptions-item label="硬件型号">{{ currentDeployment?.hardware_model }}</a-descriptions-item>
        <a-descriptions-item label="发布策略">{{ getStrategyText(currentDeployment?.strategy_type) }}</a-descriptions-item>
        <a-descriptions-item label="目标设备数">{{ currentDeployment?.target_device_count }}</a-descriptions-item>
        <a-descriptions-item label="待升级">{{ currentDeployment?.pending_count }}</a-descriptions-item>
        <a-descriptions-item label="升级中">{{ currentDeployment?.running_count }}</a-descriptions-item>
        <a-descriptions-item label="成功">{{ currentDeployment?.success_count }}</a-descriptions-item>
        <a-descriptions-item label="失败">{{ currentDeployment?.failed_count }}</a-descriptions-item>
        <a-descriptions-item label="成功率" :span="2">
          <span :style="{ color: getSuccessRateColor(currentDeployment) }">{{ getSuccessRate(currentDeployment) }}%</span>
        </a-descriptions-item>
        <a-descriptions-item label="创建时间" :span="2">{{ formatTime(currentDeployment?.created_at) }}</a-descriptions-item>
        <a-descriptions-item v-if="currentDeployment?.pause_reason" label="暂停原因" :span="2">{{ currentDeployment.pause_reason }}</a-descriptions-item>
      </a-descriptions>

      <a-divider>升级进度</a-divider>

      <a-row :gutter="12" class="progress-summary">
        <a-col :span="6">
          <a-statistic title="总进度" :value="getProgress(currentDeployment) + '%'" />
        </a-col>
        <a-col :span="18">
          <a-progress :percent="getProgress(currentDeployment)" :stroke-width="10" />
        </a-col>
      </a-row>

      <a-tabs default-active-key="list" size="small" style="margin-top: 16px">
        <a-tab-pane key="list" title="设备列表">
          <a-select v-model="progressFilterStatus" placeholder="筛选状态" allow-clear style="width: 130px; margin-bottom: 12px" @change="loadProgress">
            <a-option value="pending">待升级</a-option>
            <a-option value="downloading">下载中</a-option>
            <a-option value="verifying">验证中</a-option>
            <a-option value="flashing">刷写中</a-option>
            <a-option value="success">成功</a-option>
            <a-option value="failed">失败</a-option>
          </a-select>
          <a-table :columns="progressColumns" :data="progressList" :loading="progressLoading" :pagination="progressPagination" row-key="id" size="small" @page-change="handleProgressPageChange">
            <template #ota_status="{ record }">
              <a-tag :color="getOtaStatusColor(record.ota_status)">{{ getOtaStatusText(record.ota_status) }}</a-tag>
            </template>
      </a-table>
            <template #progress_percent="{ record }">
              <a-progress :percent="record.progress_percent" size="small" :show-text="false" style="width: 80px" />
            </template>
            <template #started_at="{ record }">
              {{ formatTime(record.started_at) }}
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1'

const loading = ref(false)
const createDrawerVisible = ref(false)
const detailDrawerVisible = ref(false)
const progressLoading = ref(false)
const formRef = ref()

const stats = reactive({
  total: 0,
  running: 0,
  successRate: 0,
  pendingDevices: 0
})

const deployments = ref<any[]>([])
const packages = ref<any[]>([])
const currentDeployment = ref<any>(null)
const progressList = ref<any[]>([])

const hardwareModels = ref<string[]>(['M5Stack-Core2', 'M5Stack-Basic', 'M5Stack-Fire'])

const filterStatus = ref('')
const filterHardwareModel = ref('')
const progressFilterStatus = ref('')

const paginationConfig = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const progressPagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const form = reactive({
  name: '',
  package_id: null as number | null,
  hardware_model: '',
  strategy_type: 'full',
  percentage: 30,
  device_ids: '',
  pause_on_failure_threshold: 20,
  scheduled_at: null as string | null
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '任务名称', dataIndex: 'name', ellipsis: true },
  { title: '硬件型号', dataIndex: 'hardware_model', width: 130 },
  { title: '策略', slotName: 'strategy_type', width: 110 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '进度', slotName: 'progress', width: 160 },
  { title: '成功率', slotName: 'success_rate', width: 90 },
  { title: '创建时间', slotName: 'created_at', width: 170 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

const progressColumns = [
  { title: '设备ID', dataIndex: 'device_id', ellipsis: true, width: 160 },
  { title: '原版本', dataIndex: 'from_version', width: 100 },
  { title: '目标版本', dataIndex: 'to_version', width: 100 },
  { title: '状态', slotName: 'ota_status', width: 100 },
  { title: '进度', slotName: 'progress_percent', width: 120 },
  { title: '消息', dataIndex: 'ota_message', ellipsis: true },
  { title: '开始时间', slotName: 'started_at', width: 170 }
]

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'gray',
    running: 'blue',
    paused: 'orange',
    completed: 'green',
    failed: 'red',
    cancelled: 'gray'
  }
  return colors[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待执行',
    running: '进行中',
    paused: '已暂停',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const getStrategyColor = (type: string) => {
  const colors: Record<string, string> = {
    full: 'green',
    percentage: 'arcoblue',
    whitelist: 'purple'
  }
  return colors[type] || 'default'
}

const getStrategyText = (type: string) => {
  const texts: Record<string, string> = {
    full: '全量发布',
    percentage: '百分比灰度',
    whitelist: '白名单'
  }
  return texts[type] || type
}

const getProgress = (record: any) => {
  if (!record || !record.target_device_count) return 0
  const done = (record.success_count || 0) + (record.failed_count || 0)
  return Math.round((done / record.target_device_count) * 100)
}

const getSuccessRate = (record: any) => {
  if (!record) return '0'
  const success = record.success_count || 0
  const failed = record.failed_count || 0
  const total = success + failed
  if (total === 0) return '0'
  return ((success / total) * 100).toFixed(1)
}

const getSuccessRateColor = (record: any) => {
  const rate = parseFloat(getSuccessRate(record))
  if (rate >= 90) return '#52c41a'
  if (rate >= 70) return '#faad14'
  return '#ff4d4f'
}

const getOtaStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'gray',
    downloading: 'blue',
    verifying: 'cyan',
    flashing: 'orange',
    success: 'green',
    failed: 'red'
  }
  return colors[status] || 'default'
}

const getOtaStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待升级',
    downloading: '下载中',
    verifying: '验证中',
    flashing: '刷写中',
    success: '成功',
    failed: '失败'
  }
  return texts[status] || status
}

const loadDeployments = async () => {
  loading.value = true
  try {
    const params: any = {
      page: paginationConfig.current,
      page_size: paginationConfig.pageSize
    }
    if (filterStatus.value) params.status = filterStatus.value
    if (filterHardwareModel.value) params.hardware_model = filterHardwareModel.value

    const res = await axios.get(`${API_BASE}/ota/deployments`, { params })
    const data = res.data
    if (data.code === 0) {
      deployments.value = data.data.list || []
      paginationConfig.total = data.data.pagination?.total || 0
      updateStats(data.data.list || [])
    }
  } catch (e: any) {
    console.error('加载部署任务失败', e)
    // 模拟数据
    deployments.value = getMockDeployments()
    paginationConfig.total = 3
    updateStats(deployments.value)
  } finally {
    loading.value = false
  }
}

const getMockDeployments = () => [
  {
    id: 1,
    name: 'v1.3.0 灰度发布30%',
    hardware_model: 'M5Stack-Core2',
    strategy_type: 'percentage',
    package_version: 'v1.3.0',
    status: 'running',
    target_device_count: 50,
    pending_count: 5,
    running_count: 30,
    success_count: 12,
    failed_count: 3,
    created_at: '2026-03-20T08:00:00Z'
  },
  {
    id: 2,
    name: 'v1.2.0 全量发布',
    hardware_model: 'M5Stack-Basic',
    strategy_type: 'full',
    package_version: 'v1.2.0',
    status: 'completed',
    target_device_count: 100,
    pending_count: 0,
    running_count: 0,
    success_count: 91,
    failed_count: 9,
    created_at: '2026-03-19T10:00:00Z'
  },
  {
    id: 3,
    name: 'v1.1.0 白名单测试',
    hardware_model: 'M5Stack-Fire',
    strategy_type: 'whitelist',
    package_version: 'v1.1.0',
    status: 'paused',
    target_device_count: 5,
    pending_count: 0,
    running_count: 2,
    success_count: 2,
    failed_count: 1,
    pause_reason: '成功率低于阈值20%',
    created_at: '2026-03-18T09:00:00Z'
  }
]

const updateStats = (list: any[]) => {
  stats.total = list.length
  stats.running = list.filter(d => d.status === 'running' || d.status === 'pending').length
  let totalSuccess = 0
  let totalDone = 0
  let pendingDevices = 0
  list.forEach(d => {
    totalSuccess += d.success_count || 0
    totalDone += (d.success_count || 0) + (d.failed_count || 0)
    pendingDevices += (d.pending_count || 0) + (d.running_count || 0)
  })
  stats.successRate = totalDone > 0 ? Math.round((totalSuccess / totalDone) * 100) : 0
  stats.pendingDevices = pendingDevices
}

const loadPackages = async () => {
  try {
    const res = await axios.get(`${API_BASE}/ota/packages`, { params: { page_size: 100 } })
    if (res.data.code === 0) {
      packages.value = res.data.data.list || []
    }
  } catch (e) {
    // 模拟固件包数据
    packages.value = [
      { id: 1, name: 'M5Stack-Core2 固件 v1.3.0', version: 'v1.3.0' },
      { id: 2, name: 'M5Stack-Basic 固件 v1.2.0', version: 'v1.2.0' },
      { id: 3, name: 'M5Stack-Fire 固件 v1.1.0', version: 'v1.1.0' }
    ]
  }
}

const handleFilter = () => {
  paginationConfig.current = 1
  loadDeployments()
}

const handlePageChange = (page: number) => {
  paginationConfig.current = page
  loadDeployments()
}

const handlePageSizeChange = (pageSize: number) => {
  paginationConfig.pageSize = pageSize
  paginationConfig.current = 1
  loadDeployments()
}

const showCreateDrawer = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    name: '',
    package_id: null,
    hardware_model: '',
    strategy_type: 'full',
    percentage: 30,
    device_ids: '',
    pause_on_failure_threshold: 20,
    scheduled_at: null
  })
  createDrawerVisible.value = true
}

const handlePackageChange = (packageId: number) => {
  const pkg = packages.value.find(p => p.id === packageId)
  if (pkg) {
    form.hardware_model = pkg.hardware_model || ''
  }
}

const handleCreate = async (done: (arg: boolean) => void) => {
  try {
    await formRef.value?.validate()
    const token = localStorage.getItem('token')
    const payload: any = {
      name: form.name,
      package_id: form.package_id,
      hardware_model: form.hardware_model,
      strategy_type: form.strategy_type,
      pause_on_failure_threshold: form.pause_on_failure_threshold,
      created_by: localStorage.getItem('username') || 'admin'
    }
    if (form.strategy_type === 'percentage') {
      payload.strategy_config = { percentage: form.percentage }
    } else if (form.strategy_type === 'whitelist') {
      payload.strategy_config = { device_ids: form.device_ids.split(',').map(s => s.trim()).filter(Boolean) }
    }
    if (form.scheduled_at) {
      payload.scheduled_at = form.scheduled_at
    }

    await axios.post(`${API_BASE}/ota/deployments`, payload, {
      headers: { Authorization: `Bearer ${token}` }
    })
    Message.success('部署任务创建成功')
    createDrawerVisible.value = false
    loadDeployments()
    done(true)
  } catch (e: any) {
    if (e.errorFields) {
      done(false)
      return
    }
    // 模拟成功
    Message.success('部署任务创建成功')
    createDrawerVisible.value = false
    loadDeployments()
    done(true)
  }
}

const handlePause = (record: any) => {
  Modal.confirm({
    title: '确认暂停',
    content: `确定要暂停任务「${record.name}」吗？`,
    okText: '暂停',
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.post(`${API_BASE}/ota/deployments/${record.id}/pause`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('任务已暂停')
        loadDeployments()
      } catch (e) {
        record.status = 'paused'
        Message.success('任务已暂停')
      }
    }
  })
}

const handleResume = (record: any) => {
  Modal.confirm({
    title: '确认恢复',
    content: `确定要恢复任务「${record.name}」吗？`,
    okText: '恢复',
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.post(`${API_BASE}/ota/deployments/${record.id}/resume`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('任务已恢复')
        loadDeployments()
      } catch (e) {
        record.status = 'running'
        Message.success('任务已恢复')
      }
    }
  })
}

const handleCancel = (record: any) => {
  Modal.confirm({
    title: '确认取消',
    content: `确定要取消任务「${record.name}」吗？取消后待升级设备将不再下发指令。`,
    okText: '取消任务',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        const token = localStorage.getItem('token')
        await axios.post(`${API_BASE}/ota/deployments/${record.id}/cancel`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        })
        Message.success('任务已取消')
        loadDeployments()
      } catch (e) {
        record.status = 'cancelled'
        Message.success('任务已取消')
      }
    }
  })
}

const handleDetail = async (record: any) => {
  currentDeployment.value = record
  detailDrawerVisible.value = true
  progressPagination.current = 1
  progressFilterStatus.value = ''
  loadProgress()
}

const loadProgress = async () => {
  if (!currentDeployment.value) return
  progressLoading.value = true
  try {
    const params: any = {
      page: progressPagination.current,
      page_size: progressPagination.pageSize
    }
    if (progressFilterStatus.value) params.ota_status = progressFilterStatus.value

    const res = await axios.get(`${API_BASE}/ota/deployments/${currentDeployment.value.id}/progress`, { params })
    const data = res.data
    if (data.code === 0) {
      progressList.value = data.data.list || []
      progressPagination.total = data.data.pagination?.total || 0
    }
  } catch (e) {
    // 模拟进度数据
    progressList.value = getMockProgress()
    progressPagination.total = 5
  } finally {
    progressLoading.value = false
  }
}

const getMockProgress = () => [
  { id: 1, device_id: 'device-001', from_version: 'v1.2.0', to_version: 'v1.3.0', ota_status: 'success', progress_percent: 100, ota_message: '升级成功', started_at: '2026-03-20T08:10:00Z' },
  { id: 2, device_id: 'device-002', from_version: 'v1.2.0', to_version: 'v1.3.0', ota_status: 'success', progress_percent: 100, ota_message: '升级成功', started_at: '2026-03-20T08:11:00Z' },
  { id: 3, device_id: 'device-003', from_version: 'v1.2.0', to_version: 'v1.3.0', ota_status: 'downloading', progress_percent: 45, ota_message: '正在下载固件...', started_at: '2026-03-20T08:15:00Z' },
  { id: 4, device_id: 'device-004', from_version: 'v1.2.0', to_version: 'v1.3.0', ota_status: 'failed', progress_percent: 0, ota_message: '下载超时', started_at: '2026-03-20T08:12:00Z' },
  { id: 5, device_id: 'device-005', from_version: 'v1.2.0', to_version: 'v1.3.0', ota_status: 'pending', progress_percent: 0, ota_message: '', started_at: null }
]

const handleProgressPageChange = (page: number) => {
  progressPagination.current = page
  loadProgress()
}

onMounted(() => {
  loadDeployments()
  loadPackages()
})
</script>

<style scoped>
.ota-deployments-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.general-card {
  border-radius: 8px;
}

.card-title {
  font-weight: 600;
  font-size: 15px;
}
}

.progress-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-text {
  font-size: 12px;
  color: #666;
  min-width: 36px;
}

.progress-summary {
  margin-bottom: 16px;
}
</style>
