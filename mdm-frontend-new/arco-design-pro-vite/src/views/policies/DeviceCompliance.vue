<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>合规策略</a-breadcrumb-item>
      <a-breadcrumb-item>设备合规状态</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card>
          <a-statistic title="设备总数" :value="stats.total">
            <template #prefix>📱</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="合规设备" :value="stats.compliant" :value-style="{ color: '#52c41a' }">
            <template #prefix>✅</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="不合规设备" :value="stats.nonCompliant" :value-style="{ color: '#ff4d4f' }">
            <template #prefix>❌</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="待检查" :value="stats.pending" :value-style="{ color: '#faad14' }">
            <template #prefix>⏳</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作栏 -->
    <a-card class="action-card">
      <a-space wrap>
        <a-select v-model="filters.complianceStatus" placeholder="合规状态" allow-clear style="width: 140px" @change="loadDeviceCompliance">
          <a-option value="compliant">合规</a-option>
          <a-option value="non_compliant">不合规</a-option>
          <a-option value="pending">待检查</a-option>
        </a-select>
        <a-select v-model="filters.platform" placeholder="平台" allow-clear style="width: 120px" @change="loadDeviceCompliance">
          <a-option value="ios">iOS</a-option>
          <a-option value="android">Android</a-option>
          <a-option value="windows">Windows</a-option>
          <a-option value="mac">macOS</a-option>
        </a-select>
        <a-input-search v-model="filters.keyword" placeholder="搜索设备ID/名称" style="width: 200px" search-button @search="loadDeviceCompliance" />
        <a-button type="primary" @click="showBatchCheckModal = true">批量合规检查</a-button>
        <a-button @click="loadDeviceCompliance">刷新</a-button>
      </a-space>
    </a-card>

    <!-- 设备合规列表 -->
    <a-card class="device-card">
      <a-table
        :columns="columns"
        :data="deviceList"
        :loading="loading"
        :pagination="pagination"
        row-key="device_id"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        :scroll="{ x: 1200 }"
      >
        <template #complianceStatus="{ record }">
          <a-tag :color="getComplianceColor(record.compliance_status)">
            {{ getComplianceText(record.compliance_status) }}
          </a-tag>
        </template>
        <template #platform="{ record }">
          <a-tag>{{ record.platform || '-' }}</a-tag>
        </template>
        <template #lastCheck="{ record }">
          {{ formatTime(record.last_compliance_check) }}
        </template>
        <template #violationCount="{ record }">
          <a-badge v-if="record.violation_count > 0" :count="record.violation_count" :max-count="99" :number-style="{ backgroundColor: '#ff4d4f' }" />
          <span v-else>0</span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="checkCompliance(record)">检查</a-button>
            <a-button type="text" size="small" status="warning" @click="viewViolations(record)">违规</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 设备合规详情抽屉 -->
    <a-drawer
      v-model:visible="showDetailDrawer"
      title="设备合规详情"
      :width="600"
    >
      <template v-if="currentDevice">
        <a-descriptions :column="2" bordered size="small">
          <a-descriptions-item label="设备ID">{{ currentDevice.device_id }}</a-descriptions-item>
          <a-descriptions-item label="设备名称">{{ currentDevice.device_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="平台">{{ currentDevice.platform || '-' }}</a-descriptions-item>
          <a-descriptions-item label="固件版本">{{ currentDevice.firmware_version || '-' }}</a-descriptions-item>
          <a-descriptions-item label="合规状态">
            <a-tag :color="getComplianceColor(currentDevice.compliance_status)">
              {{ getComplianceText(currentDevice.compliance_status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="违规数量">{{ currentDevice.violation_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="最后检查时间" :span="2">{{ formatTime(currentDevice.last_compliance_check) }}</a-descriptions-item>
          <a-descriptions-item label="绑定策略数">{{ currentDevice.policy_count || 0 }}</a-descriptions-item>
          <a-descriptions-item label="在线状态">
            <a-badge :status="currentDevice.is_online ? 'processing' : 'default'" :text="currentDevice.is_online ? '在线' : '离线'" />
          </a-descriptions-item>
        </a-descriptions>

        <!-- 违规详情列表 -->
        <a-divider>违规详情</a-divider>
        <a-table
          v-if="currentDeviceViolations.length > 0"
          :columns="violationColumns"
          :data="currentDeviceViolations"
          size="small"
          :pagination="false"
        >
          <template #severity="{ record }">
            <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
          </template>
          <template #status="{ record }">
            <a-tag :color="getViolationStatusColor(record.status)">{{ getViolationStatusText(record.status) }}</a-tag>
          </template>
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="resolveViolation(record)">处理</a-button>
          </template>
        </a-table>
        <a-empty v-else description="暂无违规记录" />
      </template>
    </a-drawer>

    <!-- 批量合规检查弹窗 -->
    <a-modal
      v-model:visible="showBatchCheckModal"
      title="批量合规检查"
      @before-ok="handleBatchCheck"
    >
      <a-form :model="batchCheckForm" layout="vertical">
        <a-form-item label="选择设备">
          <a-select
            v-model="batchCheckForm.deviceIds"
            placeholder="请选择要检查的设备"
            multiple
            allow-search
            :max-tag-count="3"
            :options="availableDevices.map(d => ({ label: `${d.device_id} (${d.device_name || d.device_id})`, value: d.device_id }))"
          />
        </a-form-item>
        <a-form-item label="或选择设备分组">
          <a-select
            v-model="batchCheckForm.groupIds"
            placeholder="请选择分组"
            multiple
            allow-search
            :max-tag-count="2"
            :options="availableGroups.map(g => ({ label: g.name, value: g.id }))"
          />
        </a-form-item>
        <a-form-item label="检查类型">
          <a-checkbox-group v-model="batchCheckForm.checkTypes">
            <a-checkbox value="all">全量检查</a-checkbox>
            <a-checkbox value="firmware">固件版本</a-checkbox>
            <a-checkbox value="encryption">加密状态</a-checkbox>
            <a-checkbox value="jailbreak">越狱检测</a-checkbox>
          </a-checkbox-group>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="showBatchCheckModal = false">取消</a-button>
        <a-button type="primary" @click="handleBatchCheck">开始检查</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const deviceList = ref([])
const showDetailDrawer = ref(false)
const showBatchCheckModal = ref(false)
const currentDevice = ref(null)
const currentDeviceViolations = ref([])
const availableDevices = ref([])
const availableGroups = ref([])

const filters = reactive({
  complianceStatus: undefined,
  platform: undefined,
  keyword: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const stats = reactive({
  total: 0,
  compliant: 0,
  nonCompliant: 0,
  pending: 0
})

const batchCheckForm = reactive({
  deviceIds: [],
  groupIds: [],
  checkTypes: ['all']
})

const columns = [
  { title: '设备ID', dataIndex: 'device_id', width: 180 },
  { title: '设备名称', dataIndex: 'device_name', width: 150 },
  { title: '平台', slotName: 'platform', width: 100 },
  { title: '固件版本', dataIndex: 'firmware_version', width: 120 },
  { title: '合规状态', slotName: 'complianceStatus', width: 100 },
  { title: '违规数', slotName: 'violationCount', width: 80 },
  { title: '最后检查', slotName: 'lastCheck', width: 160 },
  { title: '操作', slotName: 'actions', width: 150, fixed: 'right' }
]

const violationColumns = [
  { title: '策略名称', dataIndex: 'policy_name', width: 150 },
  { title: '策略类型', dataIndex: 'policy_type', width: 120 },
  { title: '期望值', dataIndex: 'expected_value', width: 100 },
  { title: '实际值', dataIndex: 'actual_value', width: 100 },
  { title: '严重级别', slotName: 'severity', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '操作', slotName: 'actions', width: 80 }
]

const loadDeviceCompliance = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, page_size: pagination.pageSize }
    if (filters.complianceStatus) params.compliance_status = filters.complianceStatus
    if (filters.platform) params.platform = filters.platform
    if (filters.keyword) params.keyword = filters.keyword

    const res = await axios.get('/api/v1/devices/compliance', { params })
    const data = res.data
    if (data.code === 0) {
      deviceList.value = data.data?.list || []
      pagination.total = data.data?.pagination?.total || 0
      updateStats()
    } else {
      // Fallback: load from violations endpoint
      await loadFromViolations()
    }
  } catch (err) {
    console.error('加载设备合规状态失败', err)
    // Fallback: load devices and check compliance manually
    await loadDevicesAndCheck()
  } finally {
    loading.value = false
  }
}

const loadFromViolations = async () => {
  try {
    const res = await axios.get('/api/v1/compliance-violations', {
      params: { page: 1, page_size: 1000 }
    })
    if (res.data.code === 0) {
      const violations = res.data.data?.list || []
      const deviceMap = new Map()
      violations.forEach(v => {
        if (!deviceMap.has(v.device_id)) {
          deviceMap.set(v.device_id, {
            device_id: v.device_id,
            device_name: v.device_name || v.device_id,
            compliance_status: 'non_compliant',
            violation_count: 0,
            violations: []
          })
        }
        const device = deviceMap.get(v.device_id)
        device.violation_count++
        device.violations.push(v)
      })
      deviceList.value = Array.from(deviceMap.values())
      pagination.total = deviceList.value.length
      updateStats()
    }
  } catch (err) {
    console.error('从违规记录加载失败', err)
  }
}

const loadDevicesAndCheck = async () => {
  try {
    const res = await axios.get('/api/v1/devices/list', {
      params: { page: 1, page_size: 1000 }
    })
    if (res.data.code === 0) {
      deviceList.value = (res.data.data?.list || []).map(d => ({
        device_id: d.device_id,
        device_name: d.device_name || d.sn_code || d.device_id,
        platform: d.platform || d.hardware_model || '-',
        firmware_version: d.firmware_version || '-',
        compliance_status: 'pending',
        violation_count: 0,
        is_online: d.is_online,
        last_compliance_check: null
      }))
      pagination.total = deviceList.value.length
      updateStats()
    }
  } catch (err) {
    console.error('加载设备列表失败', err)
    Message.error('加载设备列表失败')
  }
}

const updateStats = () => {
  stats.total = deviceList.value.length
  stats.compliant = deviceList.value.filter(d => d.compliance_status === 'compliant').length
  stats.nonCompliant = deviceList.value.filter(d => d.compliance_status === 'non_compliant').length
  stats.pending = deviceList.value.filter(d => d.compliance_status === 'pending').length
}

const handlePageChange = (page) => {
  pagination.current = page
  loadDeviceCompliance()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  loadDeviceCompliance()
}

const openDetail = async (record) => {
  currentDevice.value = record
  // Load violations for this device
  try {
    const res = await axios.get('/api/v1/compliance-violations', {
      params: { device_id: record.device_id }
    })
    if (res.data.code === 0) {
      currentDeviceViolations.value = res.data.data?.list || []
    }
  } catch (err) {
    console.error('加载违规记录失败', err)
    currentDeviceViolations.value = []
  }
  showDetailDrawer.value = true
}

const checkCompliance = async (record) => {
  try {
    const res = await axios.post(`/api/v1/devices/${record.device_id}/compliance-check`)
    if (res.data.code === 0) {
      Message.success('合规检查已触发')
      loadDeviceCompliance()
    } else {
      Message.error(res.data.message || '检查失败')
    }
  } catch (err) {
    Message.error('检查失败')
  }
}

const viewViolations = async (record) => {
  try {
    const res = await axios.get('/api/v1/compliance-violations', {
      params: { device_id: record.device_id }
    })
    if (res.data.code === 0) {
      currentDevice.value = record
      currentDeviceViolations.value = res.data.data?.list || []
      showDetailDrawer.value = true
    }
  } catch (err) {
    Message.error('加载违规记录失败')
  }
}

const resolveViolation = async (violation) => {
  try {
    const res = await axios.post(`/api/v1/compliance-violations/${violation.id}/resolve`, {
      status: '3' // 3: 已解决
    })
    if (res.data.code === 0) {
      Message.success('已标记为已解决')
      // Refresh violations
      if (currentDevice.value) {
        viewViolations(currentDevice.value)
      }
      loadDeviceCompliance()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('操作失败')
  }
}

const handleBatchCheck = async (done) => {
  if (batchCheckForm.deviceIds.length === 0 && batchCheckForm.groupIds.length === 0) {
    Message.error('请选择至少一个设备或分组')
    done(false)
    return
  }

  try {
    const res = await axios.post('/api/v1/devices/batch-compliance-check', {
      device_ids: batchCheckForm.deviceIds,
      group_ids: batchCheckForm.groupIds,
      check_types: batchCheckForm.checkTypes
    })
    if (res.data.code === 0) {
      Message.success('批量合规检查已启动')
      showBatchCheckModal.value = false
      // Reset form
      batchCheckForm.deviceIds = []
      batchCheckForm.groupIds = []
      batchCheckForm.checkTypes = ['all']
      // Refresh list
      setTimeout(() => loadDeviceCompliance(), 2000)
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (err) {
    Message.error('批量检查启动失败')
  }
  done(true)
}

const loadAvailableDevices = async () => {
  try {
    const res = await axios.get('/api/v1/devices/list', {
      params: { page: 1, page_size: 500 }
    })
    if (res.data.code === 0) {
      availableDevices.value = res.data.data?.list || []
    }
  } catch (err) {
    console.error('加载设备列表失败', err)
  }
}

const loadAvailableGroups = async () => {
  try {
    const res = await axios.get('/api/v1/groups', {
      params: { page: 1, page_size: 100 }
    })
    if (res.data.code === 0) {
      availableGroups.value = res.data.data?.list || []
    }
  } catch (err) {
    console.error('加载分组列表失败', err)
    // Fallback empty groups
    availableGroups.value = []
  }
}

const getComplianceColor = (status) => {
  const map = { compliant: 'green', non_compliant: 'red', pending: 'yellow' }
  return map[status] || 'gray'
}

const getComplianceText = (status) => {
  const map = { compliant: '合规', non_compliant: '不合规', pending: '待检查' }
  return map[status] || status
}

const getSeverityColor = (severity) => {
  const map = { 1: 'green', 2: 'blue', 3: 'orange', 4: 'red' }
  return map[severity] || 'gray'
}

const getSeverityText = (severity) => {
  const map = { 1: '低', 2: '中', 3: '高', 4: '严重' }
  return map[severity] || severity
}

const getViolationStatusColor = (status) => {
  const map = { 1: 'yellow', 2: 'blue', 3: 'green', 4: 'gray' }
  return map[status] || 'gray'
}

const getViolationStatusText = (status) => {
  const map = { 1: '待处理', 2: '处理中', 3: '已解决', 4: '已忽略' }
  return map[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadDeviceCompliance()
  loadAvailableDevices()
  loadAvailableGroups()
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.action-card {
  margin-bottom: 16px;
}

.device-card {
  background: #fff;
}
</style>
