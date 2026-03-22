<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>健康中心</a-breadcrumb-item>
      <a-breadcrumb-item>健康报告</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-radio-group v-model="reportType" type="button" @change="loadReport">
          <a-radio value="week">周报</a-radio>
          <a-radio value="month">月报</a-radio>
        </a-radio-group>
        <a-select v-model="selectedMember" placeholder="选择成员" style="width: 150px" @change="loadReport" allow-search>
          <a-option value="all">全部成员</a-option>
          <a-option value="1">张三</a-option>
          <a-option value="2">李四</a-option>
        </a-select>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="loadReport">刷新</a-button>
        <a-button @click="exportReport">导出报告</a-button>
      </a-space>
    </div>

    <!-- 健康指标汇总 -->
    <a-row :gutter="16" class="stats-card-row">
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="reportData.heart_rate" :precision="0" suffix="次/分">
            <template #prefix>
              <icon-heart :size="24" style="color: #f53f3f" />
            </template>
            <template #title>平均心率</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="reportData.blood_pressure" suffix="">
            <template #prefix>
              <icon-activity :size="24" style="color: #f53f3f" />
            </template>
            <template #title>平均血压</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="reportData.sleep_hours" :precision="1" suffix="小时">
            <template #prefix>
              <icon-moon :size="24" style="color: #722ed1" />
            </template>
            <template #title>平均睡眠</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card class="stat-card" hoverable>
          <a-statistic :value="reportData.exercise_minutes" :precision="0" suffix="分钟">
            <template #prefix>
              <icon-fire :size="24" style="color: #ff6700" />
            </template>
            <template #title>运动时长</template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- 健康评分卡片 -->
    <a-card class="score-card" title="综合健康评分">
      <div class="score-content">
        <div class="score-circle">
          <a-progress
            type="circle"
            :percent="reportData.overall_score"
            :color="getScoreColor(reportData.overall_score)"
            :stroke-width="12"
            :size="160"
          />
        </div>
        <div class="score-details">
          <div class="score-item">
            <span class="score-label">心肺功能</span>
            <a-progress :percent="reportData.cardio_score" :color="getScoreColor(reportData.cardio_score)" size="small" />
          </div>
          <div class="score-item">
            <span class="score-label">睡眠质量</span>
            <a-progress :percent="reportData.sleep_score" :color="getScoreColor(reportData.sleep_score)" size="small" />
          </div>
          <div class="score-item">
            <span class="score-label">运动活跃度</span>
            <a-progress :percent="reportData.exercise_score" :color="getScoreColor(reportData.exercise_score)" size="small" />
          </div>
          <div class="score-item">
            <span class="score-label">压力指数</span>
            <a-progress :percent="reportData.stress_score" :color="getScoreColor(reportData.stress_score)" size="small" />
          </div>
        </div>
      </div>
    </a-card>

    <!-- Tab 切换：周报/月报内容 -->
    <a-card class="report-tab-card">
      <a-tabs v-model:active-key="activeTab" @change="changeTab">
        <a-tab-pane key="summary" title="健康摘要">
          <a-descriptions :column="2" bordered>
            <a-descriptions-item label="报告周期" :span="2">{{ reportData.period }}</a-descriptions-item>
            <a-descriptions-item label="总体评价">{{ reportData.overall_comment }}</a-descriptions-item>
            <a-descriptions-item label="健康建议">{{ reportData.suggestion }}</a-descriptions-item>
          </a-descriptions>
        </a-tab-pane>
        <a-tab-pane key="heart" title="心率数据">
          <a-table :columns="heartColumns" :data="heartData" :loading="loading" row-key="id" :pagination="false" />
        </a-tab-pane>
        <a-tab-pane key="exercise" title="运动数据">
          <a-table :columns="exerciseColumns" :data="exerciseData" :loading="loading" row-key="id" :pagination="false" />
        </a-tab-pane>
        <a-tab-pane key="warnings" title="预警历史">
          <a-table :columns="warningColumns" :data="warningHistory" :loading="loading" row-key="id" :pagination="{ pageSize: 5 }">
            <template #level="{ record }">
              <a-tag :color="getWarningLevelColor(record.level)">
                {{ getWarningLevelName(record.level) }}
              </a-tag>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const reportType = ref('week')
const selectedMember = ref('all')
const activeTab = ref('summary')

const reportData = reactive({
  heart_rate: 72,
  blood_pressure: '120/80',
  sleep_hours: 7.5,
  exercise_minutes: 45,
  overall_score: 85,
  cardio_score: 88,
  sleep_score: 82,
  exercise_score: 75,
  stress_score: 70,
  period: '2026-03-16 至 2026-03-22',
  overall_comment: '本周健康状况良好，各项指标均在正常范围内。建议继续保持规律作息和适度运动。',
  suggestion: '建议每天保持7-8小时睡眠，每周进行3-5次有氧运动。'
})

const heartColumns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '静息心率', dataIndex: 'resting_rate', width: 100 },
  { title: '最高心率', dataIndex: 'max_rate', width: 100 },
  { title: '最低心率', dataIndex: 'min_rate', width: 100 },
  { title: '平均心率', dataIndex: 'avg_rate', width: 100 }
]

const exerciseColumns = [
  { title: '日期', dataIndex: 'date', width: 120 },
  { title: '运动类型', dataIndex: 'type', width: 100 },
  { title: '时长(分钟)', dataIndex: 'duration', width: 120 },
  { title: '距离(公里)', dataIndex: 'distance', width: 120 },
  { title: '消耗(千卡)', dataIndex: 'calories', width: 120 }
]

const warningColumns = [
  { title: '预警时间', dataIndex: 'time', width: 160 },
  { title: '预警类型', dataIndex: 'type', width: 120 },
  { title: '预警级别', slotName: 'level', width: 100 },
  { title: '状态', dataIndex: 'status', width: 100 }
]

const heartData = ref([
  { id: 1, date: '2026-03-22', resting_rate: 68, max_rate: 142, min_rate: 55, avg_rate: 72 },
  { id: 2, date: '2026-03-21', resting_rate: 70, max_rate: 138, min_rate: 58, avg_rate: 74 },
  { id: 3, date: '2026-03-20', resting_rate: 72, max_rate: 145, min_rate: 56, avg_rate: 75 },
  { id: 4, date: '2026-03-19', resting_rate: 71, max_rate: 140, min_rate: 57, avg_rate: 73 },
  { id: 5, date: '2026-03-18', resting_rate: 69, max_rate: 148, min_rate: 54, avg_rate: 71 }
])

const exerciseData = ref([
  { id: 1, date: '2026-03-22', type: '跑步', duration: 45, distance: 5.2, calories: 320 },
  { id: 2, date: '2026-03-21', type: '步行', duration: 60, distance: 4.5, calories: 180 },
  { id: 3, date: '2026-03-20', type: '骑行', duration: 90, distance: 18.5, calories: 650 },
  { id: 4, date: '2026-03-19', type: '游泳', duration: 40, distance: 1.0, calories: 280 },
  { id: 5, date: '2026-03-18', type: '健身', duration: 60, distance: 0, calories: 400 }
])

const warningHistory = ref([
  { id: 1, time: '2026-03-22 10:30:00', type: '高血压风险', level: 'high', status: '已确认' },
  { id: 2, time: '2026-03-22 08:00:00', type: '睡眠呼吸暂停', level: 'critical', status: '待处理' },
  { id: 3, time: '2026-03-21 15:20:00', type: '心律不齐', level: 'medium', status: '已确认' },
  { id: 4, time: '2026-03-20 09:00:00', type: '体重异常波动', level: 'low', status: '已忽略' },
  { id: 5, time: '2026-03-19 14:00:00', type: '血糖偏高', level: 'high', status: '已确认' }
])

const getScoreColor = (score) => {
  if (score >= 80) return '#00b42a'
  if (score >= 60) return '#f7ba1e'
  return '#f53f3f'
}

const getWarningLevelColor = (level) => {
  const colors = { critical: 'red', high: 'orange', medium: 'blue', low: 'green' }
  return colors[level] || 'gray'
}

const getWarningLevelName = (level) => {
  const names = { critical: '危急', high: '高', medium: '中', low: '低' }
  return names[level] || '未知'
}

const loadReport = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch(`/api/v1/health/report?type=${reportType.value}&member=${selectedMember.value}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.code === 0) {
      Object.assign(reportData, data.data || {})
    } else {
      loadMockData()
    }
  } catch (e) {
    console.error('加载健康报告失败:', e)
    loadMockData()
  } finally {
    loading.value = false
  }
}

const loadMockData = () => {
  if (reportType.value === 'week') {
    reportData.period = '2026-03-16 至 2026-03-22'
  } else {
    reportData.period = '2026-03-01 至 2026-03-22'
  }
}

const changeTab = (key) => {
  activeTab.value = key
}

const exportReport = () => {
  Message.success('报告导出功能开发中')
}

onMounted(() => {
  loadReport()
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

.score-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.score-content {
  display: flex;
  align-items: center;
  gap: 40px;
  padding: 20px;
}

.score-circle {
  flex-shrink: 0;
}

.score-details {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.score-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.score-label {
  font-size: 14px;
  color: #666;
}

.report-tab-card {
  border-radius: 8px;
}

.report-tab-card :deep(.arco-tabs-nav) {
  padding-left: 16px;
}
</style>
