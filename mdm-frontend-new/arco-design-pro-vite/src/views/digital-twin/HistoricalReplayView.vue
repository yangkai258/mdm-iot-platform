<template>
    <Breadcrumb :items="['Home','Console','']" />


  <div class="page-container">
    <!-- 搜索筛选区 -->
    <div class="search-bar">
      <a-space wrap>
        <a-select
          v-model="selectedPetId"
          placeholder="选择宠物"
          style="width: 200px"
          allow-search
          @change="loadEvents"
        >
          <a-option v-for="pet in petList" :key="pet.device_id" :value="pet.device_id">
            {{ pet.pet_name }} ({{ pet.device_id }})
          </a-option>
        </a-select>
        <a-range-picker
          v-model="dateRange"
          style="width: 260px"
          @change="handleDateChange"
        />
        <a-select
          v-model="eventType"
          placeholder="事件类型"
          style="width: 140px"
          allow-clear
          @change="loadEvents"
        >
          <a-option value="eating">进食</a-option>
          <a-option value="drinking">饮水</a-option>
          <a-option value="sleeping">睡眠</a-option>
          <a-option value="playing">玩耍</a-option>
          <a-option value="walking">散步</a-option>
          <a-option value="abnormal">异常</a-option>
        </a-select>
        <a-button type="primary" @click="loadEvents">
          <template #icon><icon-search /></template>
          查询
        </a-button>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="action-bar">
      <a-space>
        <a-button @click="loadEvents">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <a-button @click="exportEvents">
          <template #icon><icon-download /></template>
          导出
        </a-button>
        <a-button @click="goBack">返回仪表盘</a-button>
      </a-space>
    </div>

    <!-- 时间轴展示 -->
    <div class="timeline-area">
      <a-card title="事件时间轴" class="timeline-card">
        <div v-if="loading" class="loading-state">
          <a-spin size="large" />
        </div>
        <div v-else-if="events.length === 0" class="empty-state">
          <a-empty description="暂无事件数据" />
        </div>
        <a-timeline v-else class="event-timeline">
          <a-timeline-item
            v-for="event in events"
            :key="event.id"
            :color="getEventColor(event.event_type)"
            :dot-type="event.severity === 'critical' ? 'solid' : 'hollow'"
          >
            <div class="timeline-item">
              <div class="timeline-header">
                <a-tag :color="getEventColor(event.event_type)">
                  {{ getEventName(event.event_type) }}
                </a-tag>
                <span class="timeline-time">{{ formatTime(event.occurred_at) }}</span>
                <a-badge
                  v-if="event.severity && event.severity !== 'normal'"
                  :status="event.severity === 'critical' ? 'error' : 'warning'"
                  :text="event.severity === 'critical' ? '严重' : '警告'"
                />
              </div>
              <div class="timeline-content">
                <div class="event-name">{{ event.event_name }}</div>
                <div class="event-desc">{{ event.description }}</div>
                <div v-if="event.duration" class="event-duration">
                  持续时间: {{ formatDuration(event.duration) }}
                </div>
              </div>
            </div>
          </a-timeline-item>
        </a-timeline>
      </a-card>
    </div>

    <!-- 事件统计 -->
    <div class="stats-area">
      <a-row :gutter="16">
        <a-col :span="4">
          <a-statistic title="进食次数" :value="stats.eating" suffix="次">
            <template #extra>
              <icon-heart-fill style="color: #ff7d00" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="4">
          <a-statistic title="饮水次数" :value="stats.drinking" suffix="次">
            <template #extra>
              <icon-water-fill style="color: #1650ff" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="4">
          <a-statistic title="睡眠次数" :value="stats.sleeping" suffix="次">
            <template #extra>
              <icon-moon-fill style="color: #722ed1" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="4">
          <a-statistic title="玩耍次数" :value="stats.playing" suffix="次">
            <template #extra>
              <icon-live-broadcast style="color: #00b42a" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="4">
          <a-statistic title="散步次数" :value="stats.walking" suffix="次">
            <template #extra>
              <icon-location style="color: #f53f3f" />
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="4">
          <a-statistic title="异常次数" :value="stats.abnormal" suffix="次">
            <template #extra>
              <icon-exclamation-circle style="color: #f53f3f" />
            </template>
          </a-statistic>
        </a-col>
      </a-row>
    </div>

    <!-- 事件列表 -->
    <div class="content-area">
      <a-card title="事件列表">
        <a-table
          :columns="columns"
          :data="events"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #event_type="{ record }">
            <a-tag :color="getEventColor(record.event_type)">
              {{ getEventName(record.event_type) }}
            </a-tag>
          </template>
          <template #severity="{ record }">
            <a-badge
              v-if="record.severity && record.severity !== 'normal'"
              :status="record.severity === 'critical' ? 'error' : 'warning'"
              :text="getSeverityText(record.severity)"
            />
            <span v-else>正常</span>
          </template>
          <template #occurred_at="{ record }">
            {{ formatTime(record.occurred_at) }}
          </template>
          <template #duration="{ record }">
            {{ record.duration ? formatDuration(record.duration) : '-' }}
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const router = useRouter()
const route = useRoute()

const API_BASE = '/api/digital-twin'

const petList = ref([])
const selectedPetId = ref('')
const dateRange = ref([])
const eventType = ref('')
const loading = ref(false)
const events = ref([])

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showSizeChanger: true
})

const stats = reactive({
  eating: 0,
  drinking: 0,
  sleeping: 0,
  playing: 0,
  walking: 0,
  abnormal: 0
})

const columns = [
  { title: '时间', dataIndex: 'occurred_at', slotName: 'occurred_at', width: 180 },
  { title: '事件类型', dataIndex: 'event_type', slotName: 'event_type', width: 100 },
  { title: '事件名称', dataIndex: 'event_name', ellipsis: true },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '持续时间', slotName: 'duration', width: 120 },
  { title: '严重程度', slotName: 'severity', width: 100 }
]

// 获取宠物列表
const loadPets = async () => {
  try {
    const res = await axios.get(`${API_BASE}/pets`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      petList.value = res.data.data?.list || res.data.data || []
    }
  } catch {
    petList.value = [
      { device_id: 'PET001', pet_name: '小白' },
      { device_id: 'PET002', pet_name: '旺财' }
    ]
  }

  if (route.query.device) {
    selectedPetId.value = route.query.device
  } else if (petList.value.length > 0) {
    selectedPetId.value = petList.value[0].device_id
  }
  if (selectedPetId.value) {
    loadEvents()
  }
}

// 加载事件数据
const loadEvents = async () => {
  if (!selectedPetId.value) {
    Message.warning('请先选择宠物')
    return
  }
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }
    if (eventType.value) params.event_type = eventType.value
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0].format('YYYY-MM-DD')
      params.end_date = dateRange.value[1].format('YYYY-MM-DD')
    }

    const res = await axios.get(`${API_BASE}/history/events/${selectedPetId.value}`, {
      params,
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      events.value = res.data.data?.list || res.data.data || []
      pagination.total = res.data.data?.pagination?.total || events.value.length
      updateStats()
    }
  } catch {
    // 模拟数据
    events.value = generateMockEvents()
    pagination.total = events.value.length
    updateStats()
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

// 生成模拟事件
const generateMockEvents = () => {
  const types = ['eating', 'drinking', 'sleeping', 'playing', 'walking', 'abnormal']
  const events = []
  const now = Date.now()

  for (let i = 0; i < 20; i++) {
    const type = types[Math.floor(Math.random() * types.length)]
    const severity = type === 'abnormal' ? (Math.random() > 0.5 ? 'critical' : 'warning') : 'normal'
    events.push({
      id: i + 1,
      device_id: selectedPetId.value,
      event_type: type,
      event_name: getEventName(type),
      description: getEventDescription(type),
      occurred_at: new Date(now - i * 3600000 * 2).toISOString(),
      duration: type === 'sleeping' ? Math.floor(Math.random() * 7200) + 1800 : null,
      severity
    })
  }
  return events
}

// 更新统计
const updateStats = () => {
  const s = { eating: 0, drinking: 0, sleeping: 0, playing: 0, walking: 0, abnormal: 0 }
  events.value.forEach(e => {
    if (s[e.event_type] !== undefined) s[e.event_type]++
  })
  Object.assign(stats, s)
}

// 获取事件颜色
const getEventColor = (type) => {
  const colors = {
    eating: 'orange',
    drinking: 'blue',
    sleeping: 'purple',
    playing: 'green',
    walking: 'red',
    abnormal: 'red'
  }
  return colors[type] || 'gray'
}

// 获取事件名称
const getEventName = (type) => {
  const names = {
    eating: '进食',
    drinking: '饮水',
    sleeping: '睡眠',
    playing: '玩耍',
    walking: '散步',
    abnormal: '异常'
  }
  return names[type] || type
}

// 获取事件描述
const getEventDescription = (type) => {
  const descs = {
    eating: '宠物正在进食',
    drinking: '宠物正在饮水',
    sleeping: '宠物处于睡眠状态',
    playing: '宠物正在玩耍',
    walking: '宠物正在散步',
    abnormal: '检测到异常行为，请关注'
  }
  return descs[type] || ''
}

// 获取严重程度文本
const getSeverityText = (severity) => {
  const texts = { warning: '警告', critical: '严重' }
  return texts[severity] || severity
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return '--'
  const d = new Date(timeStr)
  return `${d.getFullYear()}-${(d.getMonth() + 1).toString().padStart(2, '0')}-${d.getDate().toString().padStart(2, '0')} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}

// 格式化持续时间
const formatDuration = (seconds) => {
  if (!seconds) return '-'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  if (h > 0) return `${h}小时${m}分钟`
  if (m > 0) return `${m}分钟${s}秒`
  return `${s}秒`
}

// 日期变化
const handleDateChange = () => {
  loadEvents()
}

// 表格变化
const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadEvents()
}

// 导出事件
const exportEvents = () => {
  Message.info('导出功能开发中')
}

// 返回仪表盘
const goBack = () => {
  router.push('/digital-twin/vitals')
}

onMounted(() => {
  loadPets()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
}

.search-bar {
  margin-bottom: 12px;
}

.action-bar {
  margin-bottom: 16px;
}

.timeline-area {
  margin-bottom: 16px;
}

.timeline-card {
  border-radius: 8px;
}

.loading-state,
.empty-state {
  padding: 60px 0;
  text-align: center;
}

.event-timeline {
  max-height: 400px;
  overflow-y: auto;
}

.timeline-item {
  padding-bottom: 8px;
}

.timeline-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.timeline-time {
  font-size: 12px;
  color: #86909c;
}

.timeline-content {
  padding-left: 4px;
}

.event-name {
  font-weight: 500;
  color: #1d2129;
  margin-bottom: 2px;
}

.event-desc {
  font-size: 13px;
  color: #86909c;
}

.event-duration {
  font-size: 12px;
  color: #1650ff;
  margin-top: 4px;
}

.stats-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  margin-bottom: 16px;
}

.content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>

