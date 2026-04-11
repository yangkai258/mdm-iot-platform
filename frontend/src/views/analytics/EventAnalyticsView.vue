<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>数据分析</a-breadcrumb-item>
      <a-breadcrumb-item>事件分析</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="openCreateModal">创建事件</a-button>
        <a-button @click="exportEvents">导出</a-button>
      </a-space>
    </div>

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card class="filter-card">
        <a-space wrap>
          <a-input-search v-model="searchKeyword" placeholder="搜索事件名称" style="width: 240px" search-button @search="loadEvents" />
          <a-select v-model="filterCategory" placeholder="事件类别" allow-clear style="width: 140px" @change="loadEvents">
            <a-option value="system">系统事件</a-option>
            <a-option value="user">用户行为</a-option>
            <a-option value="device">设备事件</a-option>
            <a-option value="business">业务事件</a-option>
          </a-select>
          <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="loadEvents">
            <a-option value="active">启用</a-option>
            <a-option value="inactive">停用</a-option>
          </a-select>
          <a-range-picker v-model="dateRange" style="width: 260px" @change="loadEvents" />
        </a-space>
      </a-card>
    </div>

    <!-- 事件列表 -->
    <div class="pro-content-area">
      <a-table :columns="columns" :data="events" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }" @page-change="onPageChange">
        <template #name="{ record }">
          <a-link @click="openEventDetail(record)">{{ record.name }}</a-link>
        </template>
      </a-table>
        <template #category="{ record }">
          <a-tag :color="getCategoryColor(record.category)">{{ getCategoryText(record.category) }}</a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 'active' ? 'green' : 'gray'">{{ record.status === 'active' ? '启用' : '停用' }}</a-tag>
        </template>
        <template #trend="{ record }">
          <span :class="record.trend >= 0 ? 'trend-up' : 'trend-down'">
            {{ record.trend >= 0 ? '+' : '' }}{{ record.trend }}%
          </span>
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEventDetail(record)">详情</a-button>
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">删除</a-button>
          </a-space>
        </template>
      
    </div>

    <!-- 事件详情抽屉 -->
    <a-drawer v-model:visible="detailVisible" :title="currentEvent?.name || '事件详情'" :width="900">
      <div v-if="currentEvent">
        <!-- 事件信息 -->
        <a-card title="事件信息" class="detail-card">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="事件ID">{{ currentEvent.id }}</a-descriptions-item>
            <a-descriptions-item label="事件名称">{{ currentEvent.name }}</a-descriptions-item>
            <a-descriptions-item label="类别">{{ getCategoryText(currentEvent.category) }}</a-descriptions-item>
            <a-descriptions-item label="状态">
              <a-tag :color="currentEvent.status === 'active' ? 'green' : 'gray'">
                {{ currentEvent.status === 'active' ? '启用' : '停用' }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="累计触发">{{ currentEvent.total_count || 0 }}</a-descriptions-item>
            <a-descriptions-item label="日均触发">{{ currentEvent.daily_avg || 0 }}</a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- 事件趋势 -->
        <a-card title="事件趋势" class="detail-card">
          <template #extra>
            <a-space>
              <a-select v-model="trendGranularity" style="width: 100px" @change="loadEventTrend">
                <a-option value="hour">按小时</a-option>
                <a-option value="day">按天</a-option>
                <a-option value="week">按周</a-option>
              </a-select>
              <a-select v-model="trendPeriod" style="width: 120px" @change="loadEventTrend">
                <a-option value="7">近7天</a-option>
                <a-option value="14">近14天</a-option>
                <a-option value="30">近30天</a-option>
              </a-select>
            </a-space>
          </template>
          <div ref="eventTrendChartRef" class="chart-container"></div>
        </a-card>

        <!-- 事件分布 -->
        <a-card title="事件分布" class="detail-card">
          <a-row :gutter="[16, 16]">
            <a-col :span="12">
              <div ref="eventDistByTypeRef" class="chart-container-sm"></div>
            </a-col>
            <a-col :span="12">
              <div ref="eventDistByDeviceRef" class="chart-container-sm"></div>
            </a-col>
          </a-row>
        </a-card>

        <!-- 事件漏斗 -->
        <a-card title="事件漏斗" class="detail-card">
          <div ref="eventFunnelChartRef" class="chart-container"></div>
        </a-card>
      </div>
    </a-drawer>

    <!-- 创建/编辑事件弹窗 -->
    <a-modal v-model:visible="formVisible" :title="isEditing ? '编辑事件' : '创建事件'" :width="560" @before-ok="handleSaveEvent" @cancel="formVisible = false">
      <a-form :model="eventForm" layout="vertical" ref="formRef">
        <a-form-item label="事件名称" field="name" required>
          <a-input v-model="eventForm.name" placeholder="请输入事件名称" />
        </a-form-item>
        <a-form-item label="事件标识" field="event_key" required>
          <a-input v-model="eventForm.event_key" placeholder="英文标识符，如 device_online" />
        </a-form-item>
        <a-form-item label="事件类别" field="category" required>
          <a-select v-model="eventForm.category" placeholder="请选择类别">
            <a-option value="system">系统事件</a-option>
            <a-option value="user">用户行为</a-option>
            <a-option value="device">设备事件</a-option>
            <a-option value="business">业务事件</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述" field="description">
          <a-textarea v-model="eventForm.description" placeholder="请输入描述" :max-length="200" />
        </a-form-item>
        <a-form-item label="状态" field="status">
          <a-switch v-model="eventForm.status" checked-value="active" unchecked-value="inactive" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import * as analytics from '@/api/analytics'
import * as echarts from 'echarts'

// 状态
const events = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const filterCategory = ref('')
const filterStatus = ref('')
const dateRange = ref([])
const detailVisible = ref(false)
const formVisible = ref(false)
const isEditing = ref(false)
const currentEvent = ref(null)
const trendGranularity = ref('day')
const trendPeriod = ref('7')

const eventForm = reactive({
  name: '',
  event_key: '',
  category: 'user',
  description: '',
  status: 'active'
})

// 图表 ref
const eventTrendChartRef = ref(null)
const eventDistByTypeRef = ref(null)
const eventDistByDeviceRef = ref(null)
const eventFunnelChartRef = ref(null)

const columns = [
  { title: '事件名称', slotName: 'name' },
  { title: '标识', dataIndex: 'event_key', width: 180 },
  { title: '类别', slotName: 'category', width: 100 },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '累计触发', dataIndex: 'total_count', width: 100 },
  { title: '日均触发', dataIndex: 'daily_avg', width: 100 },
  { title: '趋势', slotName: 'trend', width: 80 },
  { title: '操作', slotName: 'actions', width: 180 }
]

async function loadEvents() {
  loading.value = true
  try {
    const res = await analytics.getEventList({
      keyword: searchKeyword.value,
      category: filterCategory.value,
      status: filterStatus.value
    })
    events.value = res.data?.list || []
  } catch (e) {
    console.error('loadEvents error:', e)
  } finally {
    loading.value = false
  }
}

async function loadEventTrend() {
  if (!currentEvent.value) return
  try {
    const res = await analytics.getEventTrend(currentEvent.value.id, {
      granularity: trendGranularity.value,
      period: trendPeriod.value
    })
    const trendData = res.data?.list || []
    renderEventTrendChart(trendData)

    const distRes = await analytics.getEventDistribution(currentEvent.value.id)
    const distData = distRes.data || {}
    renderEventDistCharts(distData)

    const funnelRes = await analytics.getEventFunnel(currentEvent.value.id)
    renderEventFunnel(funnelRes.data || [])
  } catch (e) {
    console.error('loadEventTrend error:', e)
  }
}

function renderEventTrendChart(data) {
  if (!eventTrendChartRef.value) return
  const chart = echarts.init(eventTrendChartRef.value)
  const times = data.map(d => d.time)
  const counts = data.map(d => d.count)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: times },
    yAxis: { type: 'value' },
    series: [{ type: 'line', data: counts, smooth: true, areaStyle: { opacity: 0.2 } }]
  })
}

function renderEventDistCharts(data) {
  if (eventDistByTypeRef.value) {
    const chart = echarts.init(eventDistByTypeRef.value)
    chart.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: 0 },
      series: [{ type: 'pie', radius: '60%', data: data.by_type || [] }]
    })
  }
  if (eventDistByDeviceRef.value) {
    const chart = echarts.init(eventDistByDeviceRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: (data.by_device || []).map(d => d.model) },
      yAxis: { type: 'value' },
      series: [{ type: 'bar', data: (data.by_device || []).map(d => d.count) }]
    })
  }
}

function renderEventFunnel(data) {
  if (!eventFunnelChartRef.value) return
  const chart = echarts.init(eventFunnelChartRef.value)
  chart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'funnel',
      left: '10%',
      top: 20,
      bottom: 20,
      width: '80%',
      min: 0,
      max: 100,
      minSize: '0%',
      maxSize: '100%',
      sort: 'descending',
      gap: 2,
      label: { show: true, position: 'inside', formatter: '{b}: {c}' },
      data: data.map((d, i) => ({ name: d.step, value: d.count }))
    }]
  })
}

function openEventDetail(record) {
  currentEvent.value = record
  detailVisible.value = true
  nextTick(() => loadEventTrend())
}

function openCreateModal() {
  isEditing.value = false
  Object.assign(eventForm, { name: '', event_key: '', category: 'user', description: '', status: 'active' })
  formVisible.value = true
}

function openEditModal(record) {
  isEditing.value = true
  currentEvent.value = record
  Object.assign(eventForm, {
    name: record.name,
    event_key: record.event_key,
    category: record.category,
    description: record.description,
    status: record.status
  })
  formVisible.value = true
}

async function handleSaveEvent() {
  try {
    if (isEditing.value) {
      await analytics.updateEvent(currentEvent.value.id, eventForm)
    } else {
      await analytics.createEvent(eventForm)
    }
    formVisible.value = false
    detailVisible.value = false
    loadEvents()
  } catch (e) {
    console.error('handleSaveEvent error:', e)
  }
}

async function handleDelete(record) {
  try {
    await analytics.deleteEvent(record.id)
    loadEvents()
  } catch (e) {
    console.error('handleDelete error:', e)
  }
}

async function exportEvents() {
  // 导出功能（预留）
}

function onPageChange(page) {
  loadEvents()
}

function getCategoryColor(category) {
  return { system: 'arcoblue', user: 'green', device: 'orange', business: 'purple' }[category] || 'gray'
}

function getCategoryText(category) {
  return { system: '系统', user: '用户', device: '设备', business: '业务' }[category] || category
}

onMounted(() => {
  loadEvents()
})
</script>

<style scoped>
.detail-card {
  margin-bottom: 16px;
}
.chart-container {
  height: 280px;
  width: 100%;
}
.chart-container-sm {
  height: 200px;
  width: 100%;
}
.trend-up {
  color: #f53f3f;
}
.trend-down {
  color: #0fbf60;
}
</style>
