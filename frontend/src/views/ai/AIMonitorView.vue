<template>
  <div class="pro-page-container">
    <!-- 闈㈠寘灞?-->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>棣栭〉</a-breadcrumb-item>
      <a-breadcrumb-item>AI 鍔熻兘</a-breadcrumb-item>
      <a-breadcrumb-item>妯″瀷鐩戞帶</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 鎸囨爣鍗＄墖 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="骞冲潎寤惰繜" :value="metrics.avg_latency" suffix="ms" :precision="0" :value-from="0" :animation="true">
            <template #prefix><icon-clock-circle style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">杈冩槰鏃?/span>
            <span class="trend-down">鈫?12%</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="妯″瀷鍑嗙‘鐜? :value="metrics.accuracy" suffix="%" :precision="2" :value-from="0" :animation="true">
            <template #prefix><icon-check-circle style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">杈冧笂鍛?/span>
            <span class="trend-up">鈫?2.1%</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="鍚炲悙閲? :value="metrics.throughput" suffix="req/s" :precision="0" :value-from="0" :animation="true">
            <template #prefix><icon-upload style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">宄板€?/span>
            <span>{{ metrics.peak_throughput }} req/s</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="娲昏穬妯″瀷" :value="metrics.active_models" :value-from="0" :animation="true">
            <template #prefix><icon-robot style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">鐗堟湰</span>
            <span>v{{ metrics.current_version }}</span>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 鍛婅鍗＄墖 -->
    <a-alert v-if="alerts.length > 0" class="alert-banner" type="warning">
      <template #title>
        <span>妫€娴嬪埌 {{ alerts.length }} 涓紓甯告寚鏍?/span>
      </template>
      <template #content>
        <a-space wrap>
          <a-tag v-for="a in alerts" :key="a.id" :color="a.level === 'critical' ? 'red' : 'orange'">
            {{ a.message }}
          </a-tag>
        </a-space>
      </template>
    </a-alert>

    <!-- 鍥捐〃鍖?-->
    <a-row :gutter="16" class="chart-row">
      <a-col :span="12">
        <a-card title="寤惰繜瓒嬪娍">
          <div ref="latencyChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="鍚炲悙閲忚秼鍔?>
          <div ref="throughputChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" class="chart-row">
      <a-col :span="12">
        <a-card title="鍑嗙‘鐜囪秼鍔?>
          <div ref="accuracyChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="妯″瀷鐗堟湰鍒嗗竷">
          <div ref="versionChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 鎼滅储琛ㄥ崟 -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="妯″瀷鍚嶇О">
          <a-input v-model="searchForm.model_name" placeholder="鎼滅储妯″瀷鍚嶇О" allow-clear style="width: 180px" />
        </a-form-item>
        <a-form-item label="鐗堟湰">
          <a-select v-model="searchForm.version" placeholder="閫夋嫨鐗堟湰" allow-clear style="width: 140px">
            <a-option v-for="v in versions" :key="v" :value="v">v{{ v }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="鐘舵€?>
          <a-select v-model="searchForm.status" placeholder="閫夋嫨鐘舵€? allow-clear style="width: 120px">
            <a-option value="online">鍦ㄧ嚎</a-option>
            <a-option value="offline">绂荤嚎</a-option>
            <a-option value="deprecated">宸插簾寮?/a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">鎼滅储</a-button>
            <a-button @click="handleReset">閲嶇疆</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 鎿嶄綔鏍?-->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="handleRefresh">鍒锋柊鐩戞帶</a-button>
        <a-button @click="exportReport">瀵煎嚭鎶ュ憡</a-button>
      </a-space>
    </div>

    <!-- 妯″瀷鐗堟湰鍒楄〃 -->
    <div class="pro-content-area">
      <a-table
        :columns="columns"
        :data="modelList"
        :loading="loading"
        :pagination="pagination"
        :scroll="{ x: 1200 }"
        @change="handleTableChange"
        row-key="id"
      >
        <template #status="{ record }">
          <a-badge :status="record.status === 'online' ? 'normal' : record.status === 'offline' ? 'error' : 'warning'" />
          <span>{{ getStatusText(record.status) }}</span>
        </template>
        <template #accuracy="{ record }">
          <a-progress :percent="record.accuracy" size="small" :show-text="true" />
        </template>
        <template #latency_ms="{ record }">
          <span :style="{ color: record.latency_ms > 300 ? '#f53f3f' : 'inherit' }">{{ record.latency_ms }}ms</span>
        </template>
        <template #actions="{ record }">
          <a-button type="text" size="small" @click="viewDetail(record)">璇︽儏</a-button>
        </template>
      </a-table>
    </div>

    <!-- 璇︽儏寮圭獥 -->
    <a-modal v-model:visible="detailVisible" title="妯″瀷璇︽儏" :width="700" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="妯″瀷ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="妯″瀷鍚嶇О">{{ currentRecord?.name }}</a-descriptions-item>
        <a-descriptions-item label="鐗堟湰">v{{ currentRecord?.version }}</a-descriptions-item>
        <a-descriptions-item label="鐘舵€?>{{ getStatusText(currentRecord?.status) }}</a-descriptions-item>
        <a-descriptions-item label="鍑嗙‘鐜?>{{ currentRecord?.accuracy }}%</a-descriptions-item>
        <a-descriptions-item label="骞冲潎寤惰繜">{{ currentRecord?.latency_ms }}ms</a-descriptions-item>
        <a-descriptions-item label="璋冪敤娆℃暟">{{ currentRecord?.call_count }}</a-descriptions-item>
        <a-descriptions-item label="鍙戝竷鏃堕棿">{{ currentRecord?.published_at }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const modelList = ref([])
const detailVisible = ref(false)
const currentRecord = ref(null)
const latencyChartRef = ref(null)
const throughputChartRef = ref(null)
const accuracyChartRef = ref(null)
const versionChartRef = ref(null)
let charts = []

const searchForm = reactive({
  model_name: '',
  version: '',
  status: ''
})

const versions = ref(['2.1.0', '2.0.5', '2.0.0', '1.9.0'])

const metrics = reactive({
  avg_latency: 124,
  accuracy: 94.72,
  throughput: 1280,
  peak_throughput: 3200,
  active_models: 3,
  current_version: '2.1.0'
})

const alerts = ref([
  { id: 1, level: 'warning', message: '琛屼负妯″瀷 v1.9.0 寤惰繜鍋忛珮' },
  { id: 2, level: 'info', message: '鎯呮劅妯″瀷 v2.0.5 鍑嗙‘鐜囦笅闄?0.5%' }
])

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '妯″瀷鍚嶇О', dataIndex: 'name', width: 180 },
  { title: '鐗堟湰', dataIndex: 'version', width: 100 },
  { title: '鐘舵€?, width: 100, slotName: 'status' },
  { title: '鍑嗙‘鐜?, dataIndex: 'accuracy', width: 160, slotName: 'accuracy' },
  { title: '骞冲潎寤惰繜', dataIndex: 'latency_ms', width: 110, slotName: 'latency_ms' },
  { title: '璋冪敤娆℃暟', dataIndex: 'call_count', width: 110 },
  { title: '鍙戝竷鏃堕棿', dataIndex: 'published_at', width: 170 },
  { title: '鎿嶄綔', width: 80, slotName: 'actions', fixed: 'right' }
]

const getStatusText = (s) => ({ online: '鍦ㄧ嚎', offline: '绂荤嚎', deprecated: '宸插簾寮? }[s] || s)

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const handleReset = () => {
  searchForm.model_name = ''
  searchForm.version = ''
  searchForm.status = ''
  pagination.current = 1
  loadData()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadData()
}

const handleRefresh = () => {
  loadData()
  Message.success('鐩戞帶鏁版嵁宸插埛鏂?)
}

const exportReport = () => {
  Message.success('鎶ュ憡瀵煎嚭涓?..')
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const initCharts = () => {
  // 寤惰繜瓒嬪娍
  const latencyChart = echarts.init(latencyChartRef.value)
  latencyChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
    yAxis: { type: 'value' },
    series: [{ type: 'line', smooth: true, data: [180, 150, 120, 95, 110, 130, 124], areaStyle: {} }]
  })
  charts.push(latencyChart)

  // 鍚炲悙閲?  const throughputChart = echarts.init(throughputChartRef.value)
  throughputChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: [320, 280, 980, 1450, 1280, 1020, 850] }]
  })
  charts.push(throughputChart)

  // 鍑嗙‘鐜?  const accuracyChart = echarts.init(accuracyChartRef.value)
  accuracyChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
    yAxis: { type: 'value', min: 90, max: 100 },
    series: [{ type: 'line', smooth: true, data: [93.1, 93.5, 94.0, 93.8, 94.2, 94.5, 94.72] }]
  })
  charts.push(accuracyChart)

  // 鐗堟湰鍒嗗竷
  const versionChart = echarts.init(versionChartRef.value)
  versionChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: [
        { value: 4800, name: 'v2.1.0' },
        { value: 2100, name: 'v2.0.5' },
        { value: 800, name: 'v2.0.0' },
        { value: 200, name: 'v1.9.0' }
      ]
    }]
  })
  charts.push(versionChart)
}

const loadData = async () => {
  loading.value = true
  try {
    modelList.value = [
      { id: 'mdl_001', name: '琛屼负璇嗗埆妯″瀷', version: '2.1.0', status: 'online', accuracy: 96.5, latency_ms: 98, call_count: 12480, published_at: '2026-03-20 10:00:00' },
      { id: 'mdl_002', name: '鎯呮劅鍒嗘瀽妯″瀷', version: '2.0.5', status: 'online', accuracy: 94.2, latency_ms: 145, call_count: 9800, published_at: '2026-03-15 14:30:00' },
      { id: 'mdl_003', name: '璇煶鍚堟垚妯″瀷', version: '2.0.0', status: 'online', accuracy: 93.8, latency_ms: 110, call_count: 7600, published_at: '2026-03-10 09:00:00' },
      { id: 'mdl_004', name: '濮挎€佷及璁℃ā鍨?, version: '1.9.0', status: 'deprecated', accuracy: 89.5, latency_ms: 320, call_count: 3200, published_at: '2026-02-20 11:00:00' }
    ]
    pagination.total = 4
  } catch (e) {
    Message.error('鍔犺浇澶辫触')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
  nextTick(() => initCharts())
  window.addEventListener('resize', () => charts.forEach(c => c.resize()))
})

onUnmounted(() => {
  charts.forEach(c => c.dispose())
})
</script>

<style scoped lang="less">
.stat-row { margin-bottom: 16px; }
.chart-row { margin-bottom: 16px; }
.metric-trend {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-3);
  .trend-label { margin-right: 4px; }
  .trend-up { color: #00b42a; }
  .trend-down { color: #f53f3f; }
}
.alert-banner { margin-bottom: 16px; }
</style>

