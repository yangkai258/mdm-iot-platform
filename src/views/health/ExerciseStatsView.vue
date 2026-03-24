<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>健康中心</a-breadcrumb-item>
      <a-breadcrumb-item>运动统计</a-breadcrumb-item>
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
          <a-statistic :value="statsData.duration" :precision="0" suffix="分钟">
            <template #prefix>
              <span>⏱️</span>
            </template>
            <template #title>运动时长</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic :value="statsData.distance" :precision="1" suffix="公里">
            <template #prefix>
              <span>📍</span>
            </template>
            <template #title>运动距离</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic :value="statsData.calories" :precision="0" suffix="千卡">
            <template #prefix>
              <span>🔥</span>
            </template>
            <template #title>消耗卡路里</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card">
          <a-statistic :value="statsData.completion_rate" :precision="0" suffix="%">
            <template #prefix>
              <span>✅</span>
            </template>
            <template #title>目标完成率</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 操作栏 -->
    <div class="toolbar">
      <a-radio-group v-model="chartType" type="button" size="small">
        <a-radio value="day">日</a-radio>
        <a-radio value="week">周</a-radio>
        <a-radio value="month">月</a-radio>
      </a-radio-group>
    </div>

    <!-- 图表卡片 -->
    <a-card class="chart-card">
      <template #title>运动趋势</template>
      <div class="chart-placeholder">
        <a-empty v-if="loading" description="加载中..." />
        <div v-else class="chart-content">
          <div class="chart-bar">
            <div v-for="(item, index) in trendData" :key="index" class="chart-bar-item">
              <div class="chart-bar-value" :style="{ height: (item.value / maxValue * 100) + '%' }">
                <span class="value-label">{{ item.value }}</span>
              </div>
              <span class="chart-bar-label">{{ item.label }}</span>
            </div>
          </div>
        </div>
      </div>
    </a-card>

    <!-- 表格 -->
    <a-table :columns="columns" :data="exerciseRecords" :loading="loading" row-key="id" :pagination="{ pageSize: 10 }">
      <template #type="{ record }">
        <a-tag :color="getExerciseColor(record.type)">
          {{ getExerciseTypeName(record.type) }}
        </a-tag>
      </template>
      <template #duration="{ record }">
        {{ record.duration }} 分钟
      </template>
      <template #distance="{ record }">
        {{ record.distance }} 公里
      </template>
      <template #calories="{ record }">
        {{ record.calories }} 千卡
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'

const loading = ref(false)
const chartType = ref('week')

const searchForm = reactive({
  timeRange: 'week',
  dateRange: []
})

const statsData = reactive({
  duration: 0,
  distance: 0,
  calories: 0,
  completion_rate: 0
})

const trendData = ref([])
const maxValue = computed(() => {
  if (trendData.value.length === 0) return 100
  return Math.max(...trendData.value.map(item => item.value))
})

const columns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '运动类型', slotName: 'type', width: 120 },
  { title: '时长', slotName: 'duration', width: 120 },
  { title: '距离', slotName: 'distance', width: 120 },
  { title: '消耗卡路里', slotName: 'calories', width: 120 },
  { title: '完成度', dataIndex: 'completion', width: 100 }
]

const getExerciseColor = (type) => {
  const colors = { run: 'blue', walk: 'green', cycling: 'orange', swim: 'cyan', gym: 'purple' }
  return colors[type] || 'gray'
}

const getExerciseTypeName = (type) => {
  const names = { run: '跑步', walk: '步行', cycling: '骑行', swim: '游泳', gym: '健身' }
  return names[type] || '其他'
}

const loadStats = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/exercise/stats?range=${searchForm.timeRange}&chart=${chartType.value}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      Object.assign(statsData, data.data.stats || {})
      trendData.value = data.data.trend || []
    } else {
      loadMockData()
    }
  } catch (e) {
    console.error('加载运动统计失败:', e)
    loadMockData()
  } finally {
    loading.value = false
  }
}

const loadMockData = () => {
  statsData.duration = 1250
  statsData.distance = 85.5
  statsData.calories = 4850
  statsData.completion_rate = 78

  const weekLabels = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  trendData.value = weekLabels.map((label, i) => ({
    label,
    value: Math.floor(Math.random() * 500) + 100
  }))
}

const exerciseRecords = ref([
  { id: 1, date: '2026-03-22', type: 'run', duration: 45, distance: 5.2, calories: 320, completion: '85%' },
  { id: 2, date: '2026-03-21', type: 'walk', duration: 60, distance: 4.5, calories: 180, completion: '100%' },
  { id: 3, date: '2026-03-20', type: 'cycling', duration: 90, distance: 18.5, calories: 650, completion: '92%' },
  { id: 4, date: '2026-03-19', type: 'swim', duration: 40, distance: 1.0, calories: 280, completion: '75%' },
  { id: 5, date: '2026-03-18', type: 'gym', duration: 60, distance: 0, calories: 400, completion: '80%' }
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

.chart-card {
  margin-bottom: 16px;
}

.chart-placeholder {
  min-height: 250px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-content {
  width: 100%;
  padding: 20px 0;
}

.chart-bar {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  height: 200px;
  padding: 0 20px;
}

.chart-bar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  max-width: 80px;
}

.chart-bar-value {
  width: 40px;
  background: linear-gradient(180deg, #1650d8 0%, #0fc6c8 100%);
  border-radius: 4px 4px 0 0;
  min-height: 20px;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  transition: height 0.3s ease;
}

.value-label {
  font-size: 12px;
  color: #fff;
  padding-top: 4px;
  font-weight: 500;
}

.chart-bar-label {
  margin-top: 8px;
  font-size: 12px;
  color: #666;
}

.toolbar {
  margin-bottom: 16px;
}
</style>
