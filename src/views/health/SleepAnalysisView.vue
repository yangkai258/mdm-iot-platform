<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>健康中心</a-breadcrumb-item>
      <a-breadcrumb-item>睡眠分析</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="search-form">
      <a-form :model="searchForm" layout="inline">
        <a-form-item label="时间范围">
          <a-select v-model="searchForm.timeRange" placeholder="选择范围" style="width: 120px" @change="loadStats">
            <a-option value="day">今日</a-option>
            <a-option value="week">本周</a-option>
            <a-option value="month">本月</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="自定义日期">
          <a-range-picker v-model="searchForm.dateRange" style="width: 260px" @change="loadStats" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="loadStats">刷新</a-button>
        </a-form-item>
      </a-form>
    </div>

    <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic :value="sleepData.total_hours" :precision="1" suffix="小时">
            <template #prefix>
              <span>⏰</span>
            </template>
            <template #title>睡眠时长</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic :value="sleepData.quality_score" :precision="0" suffix="分">
            <template #prefix>
              <span>⭐</span>
            </template>
            <template #title>睡眠质量评分</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic :value="sleepData.deep_sleep_ratio" :precision="0" suffix="%">
            <template #prefix>
              <span>🌙</span>
            </template>
            <template #title>深睡比例</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic :value="sleepData.light_sleep_ratio" :precision="0" suffix="%">
            <template #prefix>
              <span>☀️</span>
            </template>
            <template #title>浅睡比例</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 图表区 -->
    <a-row :gutter="16" style="margin-bottom: 16px;">
      <a-col :span="12">
        <a-card title="睡眠时长趋势">
          <div class="chart-placeholder">
            <a-empty v-if="loading" description="加载中..." />
            <div v-else class="chart-content">
              <div class="chart-line">
                <svg viewBox="0 0 400 150" class="line-chart">
                  <polyline
                    fill="none"
                    stroke="#722ed1"
                    stroke-width="2"
                    :points="durationChartPoints"
                  />
                  <circle v-for="(point, i) in durationChartPointsArray" :key="i" :cx="point.x" :cy="point.y" r="3" fill="#722ed1" />
                </svg>
              </div>
              <div class="chart-labels">
                <span v-for="(label, i) in chartLabels" :key="i">{{ label }}</span>
              </div>
            </div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="睡眠构成">
          <div class="chart-placeholder">
            <a-empty v-if="loading" description="加载中..." />
            <div v-else class="sleep-composition">
              <div class="composition-item deep-sleep">
                <span class="composition-label">深睡</span>
                <span class="composition-value">{{ sleepData.deep_sleep_ratio }}%</span>
                <div class="composition-bar">
                  <div class="composition-fill deep" :style="{ width: sleepData.deep_sleep_ratio + '%' }"></div>
                </div>
              </div>
              <div class="composition-item light-sleep">
                <span class="composition-label">浅睡</span>
                <span class="composition-value">{{ sleepData.light_sleep_ratio }}%</span>
                <div class="composition-bar">
                  <div class="composition-fill light" :style="{ width: sleepData.light_sleep_ratio + '%' }"></div>
                </div>
              </div>
              <div class="composition-item rem-sleep">
                <span class="composition-label">快速眼动</span>
                <span class="composition-value">{{ sleepData.rem_ratio || 15 }}%</span>
                <div class="composition-bar">
                  <div class="composition-fill rem" :style="{ width: (sleepData.rem_ratio || 15) + '%' }"></div>
                </div>
              </div>
              <div class="composition-item awake">
                <span class="composition-label">清醒</span>
                <span class="composition-value">{{ sleepData.awake_ratio || 5 }}%</span>
                <div class="composition-bar">
                  <div class="composition-fill awake" :style="{ width: (sleepData.awake_ratio || 5) + '%' }"></div>
                </div>
              </div>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 表格 -->
    <a-table :columns="columns" :data="sleepRecords" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
      <template #quality="{ record }">
        <a-progress :percent="record.quality" :color="getQualityColor(record.quality)" size="small" />
      </template>
      <template #deep_sleep="{ record }">
        {{ record.deep_sleep }} 小时
      </template>
      <template #light_sleep="{ record }">
        {{ record.light_sleep }} 小时
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'

const loading = ref(false)

const searchForm = reactive({
  timeRange: 'week',
  dateRange: []
})

const sleepData = reactive({
  total_hours: 0,
  quality_score: 0,
  deep_sleep_ratio: 0,
  light_sleep_ratio: 0,
  rem_ratio: 15,
  awake_ratio: 5
})

const chartLabels = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

const durationChartPoints = computed(() => {
  const points = [
    { x: 0, y: 100 },
    { x: 60, y: 70 },
    { x: 120, y: 80 },
    { x: 180, y: 50 },
    { x: 240, y: 60 },
    { x: 300, y: 40 },
    { x: 360, y: 30 }
  ]
  return points.map(p => `${p.x},${p.y}`).join(' ')
})

const durationChartPointsArray = computed(() => [
  { x: 0, y: 100 },
  { x: 60, y: 70 },
  { x: 120, y: 80 },
  { x: 180, y: 50 },
  { x: 240, y: 60 },
  { x: 300, y: 40 },
  { x: 360, y: 30 }
])

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '入睡时间', dataIndex: 'sleep_time', width: 120 },
  { title: '起床时间', dataIndex: 'wake_time', width: 120 },
  { title: '睡眠时长', dataIndex: 'total_hours', width: 120 },
  { title: '深睡', slotName: 'deep_sleep', width: 100 },
  { title: '浅睡', slotName: 'light_sleep', width: 100 },
  { title: '睡眠质量', slotName: 'quality', width: 150 }
]

const getQualityColor = (quality) => {
  if (quality >= 80) return '#00b42a'
  if (quality >= 60) return '#f7ba1e'
  if (quality >= 40) return '#ff7d00'
  return '#f53f3f'
}

const loadStats = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/sleep/stats?range=${searchForm.timeRange}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      Object.assign(sleepData, data.data || {})
    } else {
      loadMockData()
    }
  } catch (e) {
    console.error('加载睡眠统计失败:', e)
    loadMockData()
  } finally {
    loading.value = false
  }
}

const loadMockData = () => {
  sleepData.total_hours = 7.5
  sleepData.quality_score = 82
  sleepData.deep_sleep_ratio = 25
  sleepData.light_sleep_ratio = 55
  sleepData.rem_ratio = 15
  sleepData.awake_ratio = 5
}

const sleepRecords = ref([
  { id: 1, date: '2026-03-22', sleep_time: '22:30', wake_time: '06:00', total_hours: '7.5h', deep_sleep: 1.8, light_sleep: 4.2, quality: 85 },
  { id: 2, date: '2026-03-21', sleep_time: '23:00', wake_time: '06:30', total_hours: '7.5h', deep_sleep: 1.9, light_sleep: 4.1, quality: 88 },
  { id: 3, date: '2026-03-20', sleep_time: '22:45', wake_time: '06:15', total_hours: '7.5h', deep_sleep: 1.7, light_sleep: 4.3, quality: 78 },
  { id: 4, date: '2026-03-19', sleep_time: '23:30', wake_time: '07:00', total_hours: '7.5h', deep_sleep: 1.5, light_sleep: 4.5, quality: 65 },
  { id: 5, date: '2026-03-18', sleep_time: '22:00', wake_time: '05:30', total_hours: '7.5h', deep_sleep: 2.0, light_sleep: 4.0, quality: 92 }
])

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 4px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
}

.chart-placeholder {
  min-height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-content {
  width: 100%;
  padding: 10px 0;
}

.chart-line {
  width: 100%;
}

.line-chart {
  width: 100%;
  height: 150px;
}

.chart-labels {
  display: flex;
  justify-content: space-around;
  margin-top: 10px;
  font-size: 12px;
  color: #666;
}

.sleep-composition {
  width: 100%;
  padding: 10px 20px;
}

.composition-item {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.composition-label {
  width: 80px;
  font-size: 14px;
  color: #333;
}

.composition-value {
  width: 50px;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.composition-bar {
  flex: 1;
  height: 20px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.composition-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.composition-fill.deep {
  background: linear-gradient(90deg, #1650d8 0%, #3478f6 100%);
}

.composition-fill.light {
  background: linear-gradient(90deg, #0fc6c8 0%, #5cd6d6 100%);
}

.composition-fill.rem {
  background: linear-gradient(90deg, #722ed1 0%, #925dff 100%);
}

.composition-fill.awake {
  background: linear-gradient(90deg, #f7ba1e 0%, #ffc107 100%);
}
</style>
