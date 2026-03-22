<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备监控</a-breadcrumb-item>
      <a-breadcrumb-item>设备日志</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <a-card class="filter-card">
      <a-space wrap>
        <a-input-search
          v-model="filters.keyword"
          placeholder="搜索关键词"
          style="width: 260px"
          search-button
          @search="handleSearch"
        />
        <a-select
          v-model="filters.level"
          placeholder="日志级别"
          style="width: 140px"
          allow-clear
          @change="handleSearch"
        >
          <a-option value="info">INFO</a-option>
          <a-option value="warn">WARN</a-option>
          <a-option value="error">ERROR</a-option>
        </a-select>
        <a-select
          v-model="filters.device_id"
          placeholder="设备ID"
          style="width: 180px"
          allow-clear
          filterable
          @change="handleSearch"
        >
          <a-option v-for="d in deviceList" :key="d" :value="d">{{ d }}</a-option>
        </a-select>
        <a-range-picker
          v-model="filters.dateRange"
          style="width: 280px"
          @change="handleSearch"
        />
        <a-space>
          <a-button type="primary" @click="handleSearch">查询</a-button>
          <a-button @click="resetFilters">重置</a-button>
          <a-button @click="loadLogs">刷新</a-button>
        </a-space>
      </a-space>
    </a-card>

    <!-- 日志列表 -->
    <a-card class="log-card">
      <a-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
        :scroll="{ x: 900 }"
      >
        <template #level="{ record }">
          <a-tag :color="getLevelColor(record.level)" class="level-tag">
            {{ record.level?.toUpperCase() }}
          </a-tag>
        </template>
        <template #message="{ record }">
          <span class="log-message" :class="{ 'error-msg': record.level === 'error' }">
            {{ record.message }}
          </span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
        </template>
      </a-table>
    </a-card>

    <!-- 日志详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="日志详情" :width="700" :footer="null">
      <a-descriptions :column="2" bordered size="small">
        <a-descriptions-item label="设备ID">{{ currentLog.device_id }}</a-descriptions-item>
        <a-descriptions-item label="日志级别">
          <a-tag :color="getLevelColor(currentLog.level)">{{ currentLog.level?.toUpperCase() }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="时间">{{ currentLog.timestamp || currentLog.created_at }}</a-descriptions-item>
        <a-descriptions-item label="日志ID">{{ currentLog.id }}</a-descriptions-item>
        <a-descriptions-item label="消息" :span="2">{{ currentLog.message }}</a-descriptions-item>
        <a-descriptions-item label="详细信息" :span="2">
          <pre class="log-detail-pre">{{ currentLog.detail || currentLog.stack || '无' }}</pre>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import monitorApi from '@/api/monitor'

const logs = ref([])
const loading = ref(false)
const detailVisible = ref(false)
const currentLog = ref({})

const deviceList = ref(['DEV001', 'DEV002', 'DEV003', 'DEV004', 'DEV005'])

const filters = reactive({
  keyword: '',
  level: '',
  device_id: '',
  dateRange: []
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showSizeChanger: true
})

const columns = [
  { title: '时间', dataIndex: 'timestamp', width: 180 },
  { title: '设备ID', dataIndex: 'device_id', width: 120, ellipsis: true },
  { title: '级别', slotName: 'level', width: 90 },
  { title: '消息', slotName: 'message', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 80, fixed: 'right' }
]

const getLevelColor = (level) => {
  const colors = { info: 'blue', warn: 'orange', error: 'red' }
  return colors[level] || 'gray'
}

const loadLogs = async () => {
  loading.value = true
  const params = {
    page: pagination.current,
    page_size: pagination.pageSize
  }
  if (filters.level) params.level = filters.level
  if (filters.device_id) params.device_id = filters.device_id
  if (filters.keyword) params.keyword = filters.keyword

  try {
    const res = await monitorApi.getDeviceLogs(params)
    if (res.code === 0 || res.data) {
      logs.value = res.data?.list || res.data || []
      pagination.total = res.data?.pagination?.total || logs.value.length
    }
  } catch {
    // 模拟数据
    const mockLevels = ['info', 'warn', 'error']
    const mockMsgs = [
      '设备连接成功，心跳正常',
      '检测到电量低于20%，请及时充电',
      '固件升级失败，连接超时',
      '传感器数据上报异常',
      '设备重启成功',
      'OTA升级包下载完成',
      '网络信号弱，切换到备用信道',
      '温度过高，自动降频保护'
    ]
    const mockLogs = Array.from({ length: 20 }, (_, i) => {
      const level = mockLevels[Math.floor(Math.random() * 3)]
      const devs = ['DEV001', 'DEV002', 'DEV003', 'DEV004', 'DEV005']
      return {
        id: `LOG${String(i + 1).padStart(5, '0')}`,
        device_id: devs[Math.floor(Math.random() * devs.length)],
        level,
        message: mockMsgs[Math.floor(Math.random() * mockMsgs.length)],
        timestamp: new Date(Date.now() - i * 60000).toLocaleString('zh-CN'),
        detail: level === 'error' ? 'Error: connection timeout at module/sensor.go:123' : undefined
      }
    })
    logs.value = mockLogs
    pagination.total = mockLogs.length
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadLogs()
}

const resetFilters = () => {
  Object.assign(filters, { keyword: '', level: '', device_id: '', dateRange: [] })
  handleSearch()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadLogs()
}

const viewDetail = (record) => {
  currentLog.value = record
  detailVisible.value = true
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb { margin-bottom: 16px; }

.filter-card {
  margin-bottom: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}

.log-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}

.level-tag {
  font-weight: 600;
  letter-spacing: 0.5px;
}

.log-message {
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.error-msg {
  color: #f53f3f;
}

.log-detail-pre {
  background: #f5f7fa;
  padding: 8px 12px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
  max-height: 200px;
  overflow-y: auto;
}
</style>
