<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>健康中心</a-breadcrumb-item>
      <a-breadcrumb-item>运动统计</a-breadcrumb-item>
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

    <!-- 数据卡片区 -->
    <a-row :gutter="16" class="stats-card-row">
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="statsData.duration" :precision="0" suffix="分钟">
            <template #prefix>
              <icon-clock :size="24" style="color: #0fc6c8" />
            </template>
            <template #title>运动时长</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="statsData.distance" :precision="1" suffix="公里">
            <template #prefix>
              <icon-location :size="24" style="color: #1650d8" />
            </template>
            <template #title>运动距离</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="statsData.calories" :precision="0" suffix="千卡">
            <template #prefix>
              <icon-fire :size="24" style="color: #ff6700" />
            </template>
            <template #title>消耗卡路里</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="statsData.completion_rate" :precision="0" suffix="%">
            <template #prefix>
              <icon-check-circle :size="24" style="color: #00b42a" />
            </template>
            <template #title>目标完成率</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 趋势图表区 -->
    <a-card class="chart-card">
      <template #title>
        <a-space>
          <span>运动趋势</span>
          <a-radio-group v-model="chartType" type="button" size="small">
            <a-radio value="day">日</a-radio>
            <a-radio value="week">周</a-radio>
            <a-radio value="month">月</a-radio>
          </a-radio-group>
        </a-space>
      </template>
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

    <!-- 运动记录列表 -->
    <div class="pro-content-area">
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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const timeRange = ref('week')
const chartType = ref('week')
const dateRange = ref([])

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
    const res = await fetch(`/api/v1/health/exercise/stats?range=${timeRange.value}&chart=${chartType.value}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      Object.assign(statsData, data.data.stats || {})
      trendData.value = data.data.trend || []
    } else {
      // 使用模拟数据
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
  { id: 5, date: '2026-03-18', type: 'gym', duration: 60, distance: 0, calories: 400, completion: '80%' },
  { id: 6, date: '2026-03-17', type: 'run', duration: 30, distance: 3.0, calories: 200, completion: '60%' },
  { id: 7, date: '2026-03-16', type: 'walk', duration: 45, distance: 3.5, calories: 140, completion: '90%' }
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

.chart-card {
  margin-bottom: 16px;
  border-radius: 8px;
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

.pro-content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
