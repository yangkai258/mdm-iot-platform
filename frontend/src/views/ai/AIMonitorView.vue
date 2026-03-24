<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 功能</a-breadcrumb-item>
      <a-breadcrumb-item>模型监控</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 指标卡片 -->
    <a-row :gutter="16" class="stat-row">
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="平均延迟" :value="metrics.avg_latency" suffix="ms" :precision="0" :value-from="0" :animation="true">
            <template #prefix><icon-clock style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">较昨日</span>
            <span class="trend-down">↓ 12%</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="模型准确率" :value="metrics.accuracy" suffix="%" :precision="2" :value-from="0" :animation="true">
            <template #prefix><icon-check-circle style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">较上周</span>
            <span class="trend-up">↑ 2.1%</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="吞吐量" :value="metrics.throughput" suffix="req/s" :precision="0" :value-from="0" :animation="true">
            <template #prefix><icon-upload style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">峰值</span>
            <span>{{ metrics.peak_throughput }} req/s</span>
          </div>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card hoverable>
          <a-statistic title="活跃模型" :value="metrics.active_models" :value-from="0" :animation="true">
            <template #prefix><icon-robot style="margin-right: 4px" /></template>
          </a-statistic>
          <div class="metric-trend">
            <span class="trend-label">版本</span>
            <span>v{{ metrics.current_version }}</span>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 告警卡片 -->
    <a-alert v-if="alerts.length > 0" class="alert-banner" type="warning">
      <template #title>
        <span>检测到 {{ alerts.length }} 个异常指标</span>
      </template>
      <template #content>
        <a-space wrap>
          <a-tag v-for="a in alerts" :key="a.id" :color="a.level === 'critical' ? 'red' : 'orange'">
            {{ a.message }}
          </a-tag>
        </a-space>
      </template>
    </a-alert>

    <!-- 图表区 -->
    <a-row :gutter="16" class="chart-row">
      <a-col :span="12">
        <a-card title="延迟趋势">
          <div ref="latencyChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="吞吐量趋势">
          <div ref="throughputChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" class="chart-row">
      <a-col :span="12">
        <a-card title="准确率趋势">
          <div ref="accuracyChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="模型版本分布">
          <div ref="versionChartRef" style="height: 220px"></div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 搜索表单 -->
    <div class="pro-search-bar">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="模型名称">
          <a-input v-model="searchForm.model_name" placeholder="搜索模型名称" allow-clear style="width: 180px" />
        </a-form-item>
        <a-form-item label="版本">
          <a-select v-model="searchForm.version" placeholder="选择版本" allow-clear style="width: 140px">
            <a-option v-for="v in versions" :key="v" :value="v">v{{ v }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model="searchForm.status" placeholder="选择状态" allow-clear style="width: 120px">
            <a-option value="online">在线</a-option>
            <a-option value="offline">离线</a-option>
            <a-option value="deprecated">已废弃</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>

    <!-- 操作栏 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="handleRefresh">刷新监控</a-button>
        <a-button @click="exportReport">导出报告</a-button>
      </a-space>
    </div>

    <!-- 模型版本列表 -->
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
          <a-button type="text" size="small" @click="viewDetail(record)">详情</a-button>
        </template>
      </a-table>
    </div>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="模型详情" :width="700" footer="null">
      <a-descriptions :column="2" bordered size="large">
        <a-descriptions-item label="模型ID">{{ currentRecord?.id }}</a-descriptions-item>
        <a-descriptions-item label="模型名称">{{ currentRecord?.name }}</a-descriptions-item>
        <a-descriptions-item label="版本">v{{ currentRecord?.version }}</a-descriptions-item>
        <a-descriptions-item label="状态">{{ getStatusText(currentRecord?.status) }}</a-descriptions-item>
        <a-descriptions-item label="准确率">{{ currentRecord?.accuracy }}%</a-descriptions-item>
        <a-descriptions-item label="平均延迟">{{ currentRecord?.latency_ms }}ms</a-descriptions-item>
        <a-descriptions-item label="调用次数">{{ currentRecord?.call_count }}</a-descriptions-item>
        <a-descriptions-item label="发布时间">{{ currentRecord?.published_at }}</a-descriptions-item>
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
  { id: 1, level: 'warning', message: '行为模型 v1.9.0 延迟偏高' },
  { id: 2, level: 'info', message: '情感模型 v2.0.5 准确率下降 0.5%' }
])

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

const columns = [
  { title: '模型名称', dataIndex: 'name', width: 180 },
  { title: '版本', dataIndex: 'version', width: 100 },
  { title: '状态', width: 100, slotName: 'status' },
  { title: '准确率', dataIndex: 'accuracy', width: 160, slotName: 'accuracy' },
  { title: '平均延迟', dataIndex: 'latency_ms', width: 110, slotName: 'latency_ms' },
  { title: '调用次数', dataIndex: 'call_count', width: 110 },
  { title: '发布时间', dataIndex: 'published_at', width: 170 },
  { title: '操作', width: 80, slotName: 'actions', fixed: 'right' }
]

const getStatusText = (s) => ({ online: '在线', offline: '离线', deprecated: '已废弃' }[s] || s)

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
  Message.success('监控数据已刷新')
}

const exportReport = () => {
  Message.success('报告导出中...')
}

const viewDetail = (record) => {
  currentRecord.value = record
  detailVisible.value = true
}

const initCharts = () => {
  // 延迟趋势
  const latencyChart = echarts.init(latencyChartRef.value)
  latencyChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
    yAxis: { type: 'value' },
    series: [{ type: 'line', smooth: true, data: [180, 150, 120, 95, 110, 130, 124], areaStyle: {} }]
  })
  charts.push(latencyChart)

  // 吞吐量
  const throughputChart = echarts.init(throughputChartRef.value)
  throughputChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: [320, 280, 980, 1450, 1280, 1020, 850] }]
  })
  charts.push(throughputChart)

  // 准确率
  const accuracyChart = echarts.init(accuracyChartRef.value)
  accuracyChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
    yAxis: { type: 'value', min: 90, max: 100 },
    series: [{ type: 'line', smooth: true, data: [93.1, 93.5, 94.0, 93.8, 94.2, 94.5, 94.72] }]
  })
  charts.push(accuracyChart)

  // 版本分布
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
      { id: 'mdl_001', name: '行为识别模型', version: '2.1.0', status: 'online', accuracy: 96.5, latency_ms: 98, call_count: 12480, published_at: '2026-03-20 10:00:00' },
      { id: 'mdl_002', name: '情感分析模型', version: '2.0.5', status: 'online', accuracy: 94.2, latency_ms: 145, call_count: 9800, published_at: '2026-03-15 14:30:00' },
      { id: 'mdl_003', name: '语音合成模型', version: '2.0.0', status: 'online', accuracy: 93.8, latency_ms: 110, call_count: 7600, published_at: '2026-03-10 09:00:00' },
      { id: 'mdl_004', name: '姿态估计模型', version: '1.9.0', status: 'deprecated', accuracy: 89.5, latency_ms: 320, call_count: 3200, published_at: '2026-02-20 11:00:00' }
    ]
    pagination.total = 4
  } catch (e) {
    Message.error('加载失败')
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
