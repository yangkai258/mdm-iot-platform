<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>健康中心</a-breadcrumb-item>
      <a-breadcrumb-item>睡眠分析</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="timeRange" placeholder="时间范围" style="width: 120px" @change="loadStats">
          <a-option value="day">今日</a-option>
          <a-option value="week">本周</a-option>
          <a-option value="month">本月</a-option>
        </a-select>
        <a-range-picker v-model="dateRange" style="width: 260px" @change="loadStats" />
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="loadStats">刷新</a-button>
      </a-space>
    </div>

    <!-- 睡眠统计卡片 -->
    <a-row :gutter="16" class="stats-card-row">
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="sleepData.total_hours" :precision="1" suffix="小时">
            <template #prefix>
              <icon-clock :size="24" style="color: #722ed1" />
            </template>
            <template #title>睡眠时长</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="sleepData.quality_score" :precision="0" suffix="分">
            <template #prefix>
              <icon-star :size="24" style="color: #f7ba1e" />
            </template>
            <template #title>睡眠质量评分</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="sleepData.deep_sleep_ratio" :precision="0" suffix="%">
            <template #prefix>
              <icon-moon :size="24" style="color: #1650d8" />
            </template>
            <template #title>深睡比例</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="sleepData.light_sleep_ratio" :precision="0" suffix="%">
            <template #prefix>
              <icon-sun :size="24" style="color: #f7ba1e" />
            </template>
            <template #title>浅睡比例</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 睡眠质量趋势 -->
    <a-row :gutter="16" class="chart-row">
      <a-col :span="12">
        <a-card class="chart-card" title="睡眠时长趋势">
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
        <a-card class="chart-card" title="睡眠构成">
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

    <!-- 睡眠记录列表 -->
    <div class="pro-content-area">
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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const timeRange = ref('week')
const dateRange = ref([])

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
    const res = await fetch(`/api/v1/health/sleep/stats?range=${timeRange.value}`, {
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
  { id: 5, date: '2026-03-18', sleep_time: '22:00', wake_time: '05:30', total_hours: '7.5h', deep_sleep: 2.0, light_sleep: 4.0, quality: 92 },
  { id: 6, date: '2026-03-17', sleep_time: '23:15', wake_time: '06:45', total_hours: '7.5h', deep_sleep: 1.6, light_sleep: 4.4, quality: 72 },
  { id: 7, date: '2026-03-16', sleep_time: '22:30', wake_time: '06:00', total_hours: '7.5h', deep_sleep: 1.8, light_sleep: 4.2, quality: 82 }
])

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.pro-page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.pro-breadcrumb {
  margin-bottom: 16px;
}

.pro-search-bar {
  margin-bottom: 12px;
}

.pro-action-bar {
  margin-bottom: 16px;
}

.stats-card-row {
  margin-bottom: 16px;
}

.stat-card {
  border-radius: 8px;
  text-align: center;
}

.stat-card :deep(.arco-statistic .arco-statistic-title) {
  margin-top: 8px;
  font-size: 14px;
  color: #666;
}

.stat-card :deep(.arco-statistic .arco-statistic-value) {
  font-size: 28px;
  font-weight: 600;
}

.chart-row {
  margin-bottom: 16px;
}

.chart-card {
  border-radius: 8px;
  height: 100%;
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

.pro-content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
