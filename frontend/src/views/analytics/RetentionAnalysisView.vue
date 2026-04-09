<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>数据分析</a-breadcrumb-item>
      <a-breadcrumb-item>留存分析</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 筛选区 -->
    <div class="pro-filter-bar">
      <a-card class="filter-card">
        <a-space wrap>
          <a-select v-model="retentionType" placeholder="留存类型" style="width: 140px" @change="loadRetentionData">
            <a-option value="daily">日留存</a-option>
            <a-option value="weekly">周留存</a-option>
            <a-option value="monthly">月留存</a-option>
          </a-select>
          <a-select v-model="segmentType" placeholder="用户群" style="width: 140px" @change="loadRetentionData">
            <a-option value="all">全部用户</a-option>
            <a-option value="new">新用户</a-option>
            <a-option value="active">活跃用户</a-option>
          </a-select>
          <a-select v-model="timeRange" placeholder="时间范围" style="width: 120px" @change="loadRetentionData">
            <a-option value="7">近7天</a-option>
            <a-option value="14">近14天</a-option>
            <a-option value="30">近30天</a-option>
            <a-option value="60">近60天</a-option>
          </a-select>
          <a-button @click="loadRetentionData">刷新</a-button>
        </a-space>
      </a-card>
    </div>

    <!-- 核心指标 -->
    <a-row :gutter="[16, 16]" class="stat-cards-row">
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="次日留存率" :value="overviewData.day1_retention || 0" suffix="%" :precision="2" :value-from="0" :animation-duration="800">
            <template #extra>
              <a-tag color="arcoblue" size="small">D+1</a-tag>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="7日留存率" :value="overviewData.day7_retention || 0" suffix="%" :precision="2" :value-from="0" :animation-duration="800">
            <template #extra>
              <a-tag color="green" size="small">D+7</a-tag>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="8">
        <a-card class="stat-card">
          <a-statistic title="30日留存率" :value="overviewData.day30_retention || 0" suffix="%" :precision="2" :value-from="0" :animation-duration="800">
            <template #extra>
              <a-tag color="purple" size="small">D+30</a-tag>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 留存曲线图 -->
    <a-row :gutter="[16, 16]" class="charts-row">
      <a-col :span="24">
        <a-card title="留存曲线" class="chart-card">
          <template #extra>
            <a-space>
              <a-checkbox v-model="showBenchmark">显示基准线</a-checkbox>
            </a-space>
          </template>
          <div ref="retentionCurveRef" class="chart-container-lg"></div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 留存报表 -->
    <a-row :gutter="[16, 16]" class="charts-row">
      <a-col :xs="24" :lg="12">
        <a-card title="留存报表" class="chart-card">
          <a-table :columns="reportColumns" :data="reportData" :loading="loading" :pagination="{ pageSize: 10 }" row-key="cohort_date" size="small">
            <template #retention="{ record, column }">
              <span :style="{ color: getRetentionColor(record[column.dataIndex]) }">
                {{ record[column.dataIndex] !== null ? record[column.dataIndex] + '%' : '-' }}
              </span>
            </template>
      </a-table>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="12">
        <a-card title="留存分布" class="chart-card">
          <div ref="retentionDistRef" class="chart-container"></div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 自定义留存 -->
    <a-row :gutter="[16, 16]" class="charts-row">
      <a-col :span="24">
        <a-card title="自定义留存分析">
          <template #extra>
            <a-space>
              <a-button size="small" @click="loadCustomRetention">刷新</a-button>
            </a-space>
          </template>
          <div class="custom-filter">
            <a-space wrap>
              <a-select v-model="customEventType" placeholder="事件类型" style="width: 140px" @change="loadCustomRetention">
                <a-option value="login">登录</a-option>
                <a-option value="purchase">付费</a-option>
                <a-option value="share">分享</a-option>
              </a-select>
              <a-select v-model="customPeriod" placeholder="周期" style="width: 120px" @change="loadCustomRetention">
                <a-option value="7">7天</a-option>
                <a-option value="14">14天</a-option>
                <a-option value="30">30天</a-option>
              </a-select>
            </a-space>
          </div>
          <div ref="customRetentionChartRef" class="chart-container"></div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import * as analytics from '@/api/analytics'
import * as echarts from 'echarts'

// 状态
const retentionType = ref('daily')
const segmentType = ref('all')
const timeRange = ref('30')
const showBenchmark = ref(false)
const loading = ref(false)
const customEventType = ref('login')
const customPeriod = ref('7')

const overviewData = ref({
  day1_retention: 0,
  day7_retention: 0,
  day30_retention: 0
})
const reportData = ref([])
const retentionCurve = ref([])

// 图表 ref
const retentionCurveRef = ref(null)
const retentionDistRef = ref(null)
const customRetentionChartRef = ref(null)

const reportColumns = [
  { title: 'Cohort', dataIndex: 'cohort_date', width: 120 },
  { title: 'D+1', dataIndex: 'day1', slotName: 'retention', width: 80 },
  { title: 'D+3', dataIndex: 'day3', slotName: 'retention', width: 80 },
  { title: 'D+7', dataIndex: 'day7', slotName: 'retention', width: 80 },
  { title: 'D+14', dataIndex: 'day14', slotName: 'retention', width: 80 },
  { title: 'D+30', dataIndex: 'day30', slotName: 'retention', width: 80 }
]

async function loadRetentionData() {
  loading.value = true
  try {
    const [overviewRes, curveRes, reportRes] = await Promise.allSettled([
      analytics.getRetentionOverview({
        type: retentionType.value,
        segment: segmentType.value
      }),
      analytics.getRetentionCurve({
        type: retentionType.value,
        segment: segmentType.value,
        period: timeRange.value
      }),
      analytics.getRetentionReport({
        type: retentionType.value,
        segment: segmentType.value
      })
    ])

    if (overviewRes.status === 'fulfilled') {
      overviewData.value = overviewRes.value.data || {}
    }
    if (curveRes.status === 'fulfilled') {
      retentionCurve.value = curveRes.value.data?.list || []
      renderRetentionCurve()
    }
    if (reportRes.status === 'fulfilled') {
      reportData.value = reportRes.value.data?.list || []
    }
  } catch (e) {
    console.error('loadRetentionData error:', e)
  } finally {
    loading.value = false
  }
}

async function loadCustomRetention() {
  try {
    const res = await analytics.getRetentionCustom({
      event_type: customEventType.value,
      period: customPeriod.value
    })
    const data = res.data || []
    renderCustomRetention(data)
  } catch (e) {
    console.error('loadCustomRetention error:', e)
  }
}

function renderRetentionCurve() {
  if (!retentionCurveRef.value) return
  const chart = echarts.init(retentionCurveRef.value)
  const days = retentionCurve.value.map(d => d.day)
  const rates = retentionCurve.value.map(d => d.rate)
  const options = {
    tooltip: { trigger: 'axis', formatter: (params) => `${params[0].name}: ${params[0].value}%` },
    xAxis: { type: 'category', data: days, name: '天数' },
    yAxis: { type: 'value', name: '留存率(%)', min: 0, max: 100 },
    series: [
      {
        name: '留存率',
        type: 'line',
        data: rates,
        smooth: true,
        areaStyle: { opacity: 0.2 },
        lineStyle: { width: 2 },
        itemStyle: { color: '#1650ff' }
      }
    ]
  }
  if (showBenchmark.value) {
    options.series.push({
      name: '基准',
      type: 'line',
      data: days.map(() => 30),
      smooth: true,
      lineStyle: { width: 1, type: 'dashed', color: '#999' }
    })
  }
  chart.setOption(options)
}

function renderCustomRetention(data) {
  if (!customRetentionChartRef.value) return
  const chart = echarts.init(customRetentionChartRef.value)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: data.map(d => d.date) },
    yAxis: { type: 'value' },
    series: [{ type: 'line', data: data.map(d => d.value), smooth: true }]
  })
}

function getRetentionColor(rate) {
  if (rate === null || rate === undefined) return '#999'
  if (rate >= 50) return '#0fbf60'
  if (rate >= 30) return '#1650ff'
  if (rate >= 10) return '#ff7a00'
  return '#f53f3f'
}

onMounted(() => {
  loadRetentionData()
  loadCustomRetention()
})
</script>

<style scoped>
.stat-cards-row {
  margin-bottom: 16px;
}
.charts-row {
  margin-bottom: 16px;
}
.stat-card {
  text-align: center;
}
.chart-card {
  height: 100%;
}
.chart-container {
  height: 280px;
  width: 100%;
}
.chart-container-lg {
  height: 360px;
  width: 100%;
}
.custom-filter {
  margin-bottom: 12px;
}
</style>
