<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>AI 管理</a-breadcrumb-item>
      <a-breadcrumb-item>AI 质量监控</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索栏：时间范围 -->
    <div class="pro-search-bar">
      <a-space>
        <a-radio-group v-model="timeRange" type="button">
          <a-radio value="1h">最近1小时</a-radio>
          <a-radio value="24h">最近24小时</a-radio>
          <a-radio value="7d">最近7天</a-radio>
          <a-radio value="custom">自定义</a-radio>
        </a-radio-group>
        <a-date-picker.RangePicker
          v-if="timeRange === 'custom'"
          v-model="customRange"
          @change="loadAll"
        />
      </a-space>
    </div>

    <!-- 操作按钮 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button @click="loadAll">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 统计卡片区 -->
    <div class="metric-cards">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-statistic title="推理次数" :value="metrics.total_inferences" :loading="metricsLoading" :value-from="0" animation>
            <template #suffix>次</template>
            <template #extra>
              <span :class="changeClass(metrics.total_inferences_change)">
                {{ formatChange(metrics.total_inferences_change) }}
              </span>
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="平均延迟" :value="metrics.avg_latency_ms" :loading="metricsLoading" :value-from="0" animation>
            <template #suffix>ms</template>
            <template #extra>
              <span :class="changeClass(metrics.avg_latency_change, true)">
                {{ formatChange(metrics.avg_latency_change, true) }}
              </span>
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="错误率" :value="metrics.error_rate" :loading="metricsLoading" :value-from="0" animation :precision="2" suffix="%">
            <template #extra>
              <span :class="changeClass(metrics.error_rate_change, true)">
                {{ formatChange(metrics.error_rate_change, true) }}
              </span>
            </template>
          </a-statistic>
        </a-col>
        <a-col :span="6">
          <a-statistic title="置信度" :value="metrics.avg_confidence" :loading="metricsLoading" :value-from="0" animation :precision="1" suffix="%">
            <template #extra>
              <span :class="changeClass(metrics.avg_confidence_change)">
                {{ formatChange(metrics.avg_confidence_change) }}
              </span>
            </template>
          </a-statistic>
        </a-col>
      </a-row>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <!-- 推理延迟趋势 -->
      <a-card class="chart-card" title="推理延迟趋势">
        <div ref="latencyChartRef" class="chart-container"></div>
      </a-card>

      <!-- 错误率趋势 -->
      <a-card class="chart-card" title="错误率趋势">
        <div ref="errorRateChartRef" class="chart-container"></div>
      </a-card>

      <!-- 模型调用分布 -->
      <a-card class="chart-card" title="模型调用分布">
        <div ref="modelDistChartRef" class="chart-container"></div>
      </a-card>

      <!-- 异常事件告警 -->
      <a-card class="chart-card" title="异常事件告警">
        <a-table
          :columns="anomalyColumns"
          :data="anomalyAlerts"
          :loading="logLoading"
          :pagination="false"
          size="small"
          row-key="id"
        >
          <template #anomaly_score="{ record }">
            <a-tag :color="getAnomalyColor(record.anomaly_score)">
              {{ (record.anomaly_score * 100).toFixed(1) }}%
            </a-tag>
          </template>
          <template #actions="{ record }">
            <a-button type="text" size="small" @click="goToDetail(record.id)">详情</a-button>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAIQuality } from '@/composables/useAIQuality'
import * as echarts from 'echarts'

const router = useRouter()
const {
  metricsLoading,
  metrics,
  trendLoading,
  latencyTrend,
  errorRateTrend,
  logLoading,
  logList,
  anomalyAlerts,
  loadMetrics,
  loadTrend,
  loadLogs
} = useAIQuality()

// 时间范围
const timeRange = ref('24h')
const customRange = ref([])

// Chart refs
const latencyChartRef = ref(null)
const errorRateChartRef = ref(null)
const modelDistChartRef = ref(null)

let latencyChart = null
let errorRateChart = null
let modelDistChart = null

// 异常告警列
const anomalyColumns = [
  { title: '时间', dataIndex: 'created_at', width: 160 },
  { title: '设备', dataIndex: 'device_id', ellipsis: true },
  { title: '模型', dataIndex: 'model_version', ellipsis: true },
  { title: '异常分数', slotName: 'anomaly_score', width: 100 },
  { title: '操作', slotName: 'actions', width: 80 }
]

/**
 * 格式化环比变化
 */
const formatChange = (val, isNegativeGood = false) => {
  if (val === null || val === undefined) return '--'
  const abs = Math.abs(val)
  const arrow = val > 0 ? '↑' : val < 0 ? '↓' : ''
  const suffix = isNegativeGood ? (val < 0 ? '↓' : '↑') : (val > 0 ? '↑' : '↓')
  return `${abs.toFixed(1)}%`
}

const changeClass = (val, isNegativeGood = false) => {
  if (val === null || val === undefined) return 'change-neutral'
  const good = isNegativeGood ? val <= 0 : val >= 0
  return good ? 'change-good' : 'change-bad'
}

const getAnomalyColor = (score) => {
  if (score >= 0.8) return 'red'
  if (score >= 0.5) return 'orange'
  return 'yellow'
}

const goToDetail = (id) => {
  router.push(`/ai/behavior-detail/${id}`)
}

// ============ 图表渲染 ============
const renderLatencyChart = () => {
  if (!latencyChartRef.value) return
  if (!latencyChart) latencyChart = echarts.init(latencyChartRef.value)

  const data = latencyTrend.value.map(d => [d.time, d.value])
  latencyChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'time', axisLabel: { formatter: '{HH:mm}' } },
    yAxis: { type: 'value', name: 'ms' },
    series: [{
      type: 'line',
      data,
      smooth: true,
      areaStyle: { opacity: 0.15 },
      lineStyle: { color: '#1650ff' },
      itemStyle: { color: '#1650ff' }
    }],
    grid: { left: 50, right: 16, top: 20, bottom: 30 }
  })
}

const renderErrorRateChart = () => {
  if (!errorRateChartRef.value) return
  if (!errorRateChart) errorRateChart = echarts.init(errorRateChartRef.value)

  const data = errorRateTrend.value.map(d => [d.time, d.value])
  errorRateChart.setOption({
    tooltip: { trigger: 'axis', formatter: (p) => `${p[0].axisValue}<br/>错误率: ${(p[0].value[1] * 100).toFixed(2)}%` },
    xAxis: { type: 'time', axisLabel: { formatter: '{HH:mm}' } },
    yAxis: { type: 'value', name: '%', axisLabel: { formatter: (v) => (v * 100).toFixed(2) } },
    series: [{
      type: 'line',
      data,
      smooth: true,
      areaStyle: { opacity: 0.15 },
      lineStyle: { color: '#f53f3f' },
      itemStyle: { color: '#f53f3f' }
    }],
    grid: { left: 50, right: 16, top: 20, bottom: 30 }
  })
}

const renderModelDistChart = () => {
  if (!modelDistChartRef.value) return
  if (!modelDistChart) modelDistChart = echarts.init(modelDistChartRef.value)

  // 聚合模型调用分布
  const dist = {}
  logList.value.forEach(log => {
    const model = log.model_version || '未知'
    dist[model] = (dist[model] || 0) + 1
  })
  const pieData = Object.entries(dist).map(([name, value]) => ({ name, value }))

  modelDistChart.setOption({
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', right: 10, top: 'center' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: pieData.length ? pieData : [{ name: '无数据', value: 0 }],
      label: { formatter: '{b}: {d}%' }
    }],
    grid: { left: 10, right: 10, top: 10, bottom: 10 }
  })
}

const loadAll = async () => {
  const params = timeRange.value !== 'custom' ? { time_range: timeRange.value } : {}

  await Promise.all([
    loadMetrics(params),
    loadTrend(params),
    loadLogs({ ...params, page_size: 500 })
  ])

  await nextTick()
  renderLatencyChart()
  renderErrorRateChart()
  renderModelDistChart()
}

// 监听时间范围变化
watch(timeRange, () => loadAll())

// 窗口 resize 时重绘图表
const handleResize = () => {
  latencyChart?.resize()
  errorRateChart?.resize()
  modelDistChart?.resize()
}

onMounted(() => {
  loadAll()
  window.addEventListener('resize', handleResize)
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; display: flex; justify-content: flex-start; }

.metric-cards {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  margin-bottom: 16px;
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.chart-card {
  border-radius: 8px;
}

.chart-container {
  width: 100%;
  height: 280px;
}

.change-good { color: #00b42a; font-size: 12px; }
.change-bad { color: #f53f3f; font-size: 12px; }
.change-neutral { color: #8a8a8a; font-size: 12px; }
</style>
