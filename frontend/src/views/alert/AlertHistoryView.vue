<template>
  <div class="alert-history-view">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <a-form :model="filters" layout="inline" ref="filterFormRef">
        <a-form-item label="设备">
          <a-input v-model="filters.device_id" placeholder="设备 ID" style="width: 160px" allow-clear />
        </a-form-item>
        <a-form-item label="告警类型">
          <a-select v-model="filters.alert_type" placeholder="选择类型" style="width: 140px" allow-clear>
            <a-option value="temperature">温度告警</a-option>
            <a-option value="humidity">湿度告警</a-option>
            <a-option value="battery">电量告警</a-option>
            <a-option value="offline">离线告警</a-option>
            <a-option value="threshold">阈值告警</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="告警级别">
          <a-select v-model="filters.severity" placeholder="选择级别" style="width: 120px" allow-clear>
            <a-option :value="1">提示</a-option>
            <a-option :value="2">警告</a-option>
            <a-option :value="3">严重</a-option>
            <a-option :value="4">紧急</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="处理状态">
          <a-select v-model="filters.status" placeholder="选择状态" style="width: 120px" allow-clear>
            <a-option :value="0">未处理</a-option>
            <a-option :value="1">已确认</a-option>
            <a-option :value="2">已解决</a-option>
            <a-option :value="3">已忽略</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker
            v-model="dateRange"
            style="width: 280px"
            @change="onDateChange"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              「查询」
            </a-button>
            <a-button @click="handleReset">
              <template #icon><icon-refresh /></template>
              「重置」
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作区域 -->
    <div class="action-section">
      <a-space>
        <a-button type="primary" @click="handleExport">
          <template #icon><icon-download /></template>
          「导出」
        </a-button>
        <a-button @click="handleRefresh">
          <template #icon><icon-refresh /></template>
          「刷新」
        </a-button>
      </a-space>
      <span class="total-count">共 {{ pagination.total }} 条记录</span>
    </div>

    <!-- 表格区域 -->
    <div class="table-section">
      <a-table
        :data="alerts"
        :loading="loading"
        :pagination="paginationConfig"
        :columns="columns"
        :scroll="{ x: 1200 }"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        row-key="id"
        stripe
      >
        <template #columns>
          <a-table-column title="时间" data-index="created_at" :width="160">
            <template #cell="{ record }">
              {{ formatTime(record.created_at) }}
            </template>
          </a-table-column>
          <a-table-column title="设备" data-index="device_id" :width="140">
            <template #cell="{ record }">
              <a-tooltip :content="record.device_id">
                <span class="device-id">{{ record.device_name || record.device_id }}</span>
              </a-tooltip>
            </template>
          </a-table-column>
          <a-table-column title="告警类型" data-index="alert_type" :width="120">
            <template #cell="{ record }">
              {{ alertTypeLabel(record.alert_type) }}
            </template>
          </a-table-column>
          <a-table-column title="级别" data-index="severity" :width="80">
            <template #cell="{ record }">
              <a-tag :color="SEVERITY_MAP[record.severity]?.color || 'gray'" size="small">
                {{ SEVERITY_MAP[record.severity]?.label || record.severity }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="告警消息" data-index="message">
            <template #cell="{ record }">
              <a-tooltip :content="record.message">
                <span class="message-text">{{ record.message }}</span>
              </a-tooltip>
            </template>
          </a-table-column>
          <a-table-column title="触发值" data-index="trigger_value" :width="100">
            <template #cell="{ record }">
              {{ record.trigger_value || '-' }}
            </template>
          </a-table-column>
          <a-table-column title="处理状态" data-index="status" :width="90">
            <template #cell="{ record }">
              <a-tag :color="STATUS_MAP[record.status]?.color || 'gray'" size="small">
                {{ STATUS_MAP[record.status]?.label || record.status }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="操作" :width="160" fixed="right">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="handleViewDetail(record)">
                  「详情」
                </a-button>
                <a-button
                  v-if="record.status === 0"
                  type="text"
                  size="small"
                  @click="handleConfirm(record)"
                >
                  「确认」
                </a-button>
                <a-button
                  v-if="record.status !== 2"
                  type="text"
                  size="small"
                  @click="handleResolve(record)"
                >
                  「解决」
                </a-button>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>

    <!-- 详情弹窗 -->
    <a-modal
      v-model:visible="detailVisible"
      title="告警详情"
      :width="640"
      :footer="null"
    >
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="告警 ID">{{ currentAlert?.id }}</a-descriptions-item>
        <a-descriptions-item label="原始 ID">{{ currentAlert?.original_id }}</a-descriptions-item>
        <a-descriptions-item label="设备 ID">{{ currentAlert?.device_id }}</a-descriptions-item>
        <a-descriptions-item label="设备名称">{{ currentAlert?.device_name || '-' }}</a-descriptions-item>
        <a-descriptions-item label="告警类型">{{ alertTypeLabel(currentAlert?.alert_type) }}</a-descriptions-item>
        <a-descriptions-item label="告警级别">
          <a-tag :color="SEVERITY_MAP[currentAlert?.severity]?.color" size="small">
            {{ SEVERITY_MAP[currentAlert?.severity]?.label }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="触发值">{{ currentAlert?.trigger_value || '-' }}</a-descriptions-item>
        <a-descriptions-item label="阈值">{{ currentAlert?.threshold || '-' }}</a-descriptions-item>
        <a-descriptions-item label="处理状态">
          <a-tag :color="STATUS_MAP[currentAlert?.status]?.color" size="small">
            {{ STATUS_MAP[currentAlert?.status]?.label }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="通知渠道" :span="2">
          {{ currentAlert?.notified_channels?.join(', ') || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="告警消息" :span="2">{{ currentAlert?.message }}</a-descriptions-item>
        <a-descriptions-item label="发生时间">{{ formatTime(currentAlert?.created_at) }}</a-descriptions-item>
        <a-descriptions-item label="确认时间">{{ formatTime(currentAlert?.confirmed_at) || '-' }}</a-descriptions-item>
        <a-descriptions-item label="解决时间">{{ formatTime(currentAlert?.resolved_at) || '-' }}</a-descriptions-item>
        <a-descriptions-item label="解决备注" :span="2">{{ currentAlert?.resolve_remark || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 解决弹窗 -->
    <a-modal
      v-model:visible="resolveVisible"
      title="解决告警"
      :width="400"
      @ok="handleResolveSubmit"
    >
      <a-form :model="resolveForm" layout="vertical">
        <a-form-item label="解决备注">
          <a-textarea v-model="resolveForm.remark" :rows="3" placeholder="可选备注信息" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconSearch, IconRefresh, IconDownload } from '@arco-design/web-icons/vue'
import { useAlertHistory, SEVERITY_MAP, STATUS_MAP } from '@/composables/useNotification'

const {
  loading,
  alerts,
  pagination,
  filters,
  loadAlerts,
  confirmAlert,
  resolveAlert,
  exportAlerts,
  resetFilters
} = useAlertHistory()

const filterFormRef = ref(null)
const detailVisible = ref(false)
const resolveVisible = ref(false)
const currentAlert = ref(null)
const dateRange = ref([])
const resolveForm = reactive({
  remark: ''
})

const columns = [
  { title: '时间', dataIndex: 'created_at' },
  { title: '设备', dataIndex: 'device_id' },
  { title: '告警类型', dataIndex: 'alert_type' },
  { title: '级别', dataIndex: 'severity' },
  { title: '告警消息', dataIndex: 'message' },
  { title: '触发值', dataIndex: 'trigger_value' },
  { title: '处理状态', dataIndex: 'status' },
  { title: '操作', dataIndex: 'action' }
]

const paginationConfig = computed(() => ({
  current: pagination.current,
  pageSize: pagination.pageSize,
  total: pagination.total,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
}))

function alertTypeLabel(type) {
  const map = {
    temperature: '温度告警',
    humidity: '湿度告警',
    battery: '电量告警',
    offline: '离线告警',
    threshold: '阈值告警'
  }
  return map[type] || type
}

function formatTime(time) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

function onDateChange(value) {
  if (value && value.length === 2) {
    filters.start_time = value[0]?.toISOString()
    filters.end_time = value[1]?.toISOString()
  } else {
    filters.start_time = undefined
    filters.end_time = undefined
  }
}

function handleSearch() {
  pagination.current = 1
  loadAlerts()
}

function handleReset() {
  dateRange.value = []
  resetFilters()
  loadAlerts()
}

function handleRefresh() {
  loadAlerts()
  Message.success('已刷新')
}

function handlePageChange(page) {
  pagination.current = page
  loadAlerts()
}

function handlePageSizeChange(size) {
  pagination.pageSize = size
  pagination.current = 1
  loadAlerts()
}

function handleViewDetail(record) {
  currentAlert.value = record
  detailVisible.value = true
}

async function handleConfirm(record) {
  await confirmAlert(record.id)
}

function handleResolve(record) {
  currentAlert.value = record
  resolveForm.remark = ''
  resolveVisible.value = true
}

async function handleResolveSubmit() {
  await resolveAlert(currentAlert.value.id, resolveForm.remark)
  resolveVisible.value = false
}

async function handleExport() {
  await exportAlerts()
}

onMounted(() => {
  loadAlerts()
})
</script>

<style scoped>
.alert-history-view {
  padding: 20px;
}
.filter-section {
  background: #fff;
  padding: 16px;
  border-radius: 8px 8px 0 0;
  margin-bottom: 0;
}
.action-section {
  background: #fff;
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e5e6eb;
}
.total-count {
  color: #86909c;
  font-size: 14px;
}
.table-section {
  background: #fff;
  border-radius: 0 0 8px 8px;
}
.device-id {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}
.message-text {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}
</style>
