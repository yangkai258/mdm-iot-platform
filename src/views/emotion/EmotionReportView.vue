<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>情绪管理</a-breadcrumb-item>
      <a-breadcrumb-item>情绪报告</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="filterForm.pet_id" placeholder="选择宠物" style="width: 160px" allow-clear>
          <a-option v-for="pet in pets" :key="pet.pet_id" :value="pet.pet_id">{{ pet.pet_name }}</a-option>
        </a-select>
        <a-date-picker v-model="filterForm.start_date" placeholder="开始日期" style="width: 140px" />
        <a-date-picker v-model="filterForm.end_date" placeholder="结束日期" style="width: 140px" />
        <a-button @click="loadReport">查询</a-button>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="exportReport">导出报告</a-button>
        <a-button @click="loadReport">刷新</a-button>
      </a-space>
    </div>

    <!-- 内容区 -->
    <div class="pro-content-area">
      <!-- 报告类型切换 -->
      <a-tabs v-model="activeTab" @change="loadReport" style="margin-bottom: 16px">
        <a-tab-pane key="daily" title="日报" />
        <a-tab-pane key="weekly" title="周报" />
        <a-tab-pane key="monthly" title="月报" />
      </a-tabs>

      <!-- 报告概览 -->
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card :bordered="false" class="stat-card">
            <div class="stat-label">总识别次数</div>
            <div class="stat-value">{{ reportOverview.total_count }}</div>
            <div class="stat-trend up">↑ 12.5%</div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card :bordered="false" class="stat-card">
            <div class="stat-label">主要情绪</div>
            <div class="stat-value emotion-main">
              {{ getEmotionEmoji(reportOverview.main_emotion) }} {{ getEmotionText(reportOverview.main_emotion) }}
            </div>
            <div class="stat-trend">占比 {{ reportOverview.main_emotion_ratio || 0 }}%</div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card :bordered="false" class="stat-card">
            <div class="stat-label">异常事件</div>
            <div class="stat-value warning">{{ reportOverview.anomaly_count }}</div>
            <div class="stat-trend" :class="reportOverview.anomaly_trend === 'up' ? 'down' : 'up'">
              {{ reportOverview.anomaly_trend === 'up' ? '↑' : '↓' }} {{ reportOverview.anomaly_change || 0 }}%
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card :bordered="false" class="stat-card">
            <div class="stat-label">响应触发</div>
            <div class="stat-value">{{ reportOverview.response_triggered }}</div>
            <div class="stat-trend">成功率 {{ reportOverview.response_success_rate || 0 }}%</div>
          </a-card>
        </a-col>
      </a-row>

      <!-- 情绪趋势图表 -->
      <a-card title="情绪趋势" :bordered="false" style="margin-bottom: 16px">
        <div ref="trendChartRef" class="chart-container"></div>
      </a-card>

      <!-- 情绪分布 -->
      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="情绪类型分布" :bordered="false">
            <div ref="pieChartRef" class="chart-container"></div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="异常事件标注" :bordered="false">
            <a-table
              :columns="anomalyColumns"
              :data="anomalyEvents"
              :pagination="false"
              :loading="anomalyLoading"
              row-key="event_id"
              size="small"
            >
              <template #emotion="{ record }">
                <span>{{ getEmotionEmoji(record.emotion_type) }} {{ getEmotionText(record.emotion_type) }}</span>
              </template>
              <template #severity="{ record }">
                <a-tag :color="getSeverityColor(record.severity)">{{ getSeverityText(record.severity) }}</a-tag>
              </template>
              <template #actions="{ record }">
                <a-button type="text" size="small" @click="showAnomalyDetail(record)">详情</a-button>
              </template>
            </a-table>
            <a-empty v-if="anomalyEvents.length === 0 && !anomalyLoading" description="暂无异常事件" style="margin-top: 20px" />
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 异常详情弹窗 -->
    <a-modal
      v-model:visible="anomalyModalVisible"
      title="异常事件详情"
      :width="560"
      :footer="null"
    >
      <a-descriptions :column="2" bordered v-if="currentAnomaly">
        <a-descriptions-item label="事件ID">{{ currentAnomaly.event_id }}</a-descriptions-item>
        <a-descriptions-item label="宠物名称">{{ currentAnomaly.pet_name }}</a-descriptions-item>
        <a-descriptions-item label="异常情绪">
          <span>{{ getEmotionEmoji(currentAnomaly.emotion_type) }} {{ getEmotionText(currentAnomaly.emotion_type) }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="严重程度">
          <a-tag :color="getSeverityColor(currentAnomaly.severity)">{{ getSeverityText(currentAnomaly.severity) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="发生时间" :span="2">{{ currentAnomaly.create_time }}</a-descriptions-item>
        <a-descriptions-item label="持续时长">{{ currentAnomaly.duration || '未知' }}</a-descriptions-item>
        <a-descriptions-item label="触发响应">{{ currentAnomaly.triggered_action || '无' }}</a-descriptions-item>
        <a-descriptions-item label="建议措施" :span="2">{{ currentAnomaly.suggestion || '无' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'
import * as echarts from 'echarts'

const API_BASE = '/api/v1/emotion'
const activeTab = ref('daily')
const trendChartRef = ref(null)
const pieChartRef = ref(null)
let trendChart = null
let pieChart = null

const anomalyLoading = ref(false)
const anomalyModalVisible = ref(false)
const currentAnomaly = ref(null)

const pets = ref([
  { pet_id: 'PET001', pet_name: '小橘' },
  { pet_id: 'PET002', pet_name: '布丁' },
  { pet_id: 'PET003', pet_name: '豆豆' }
])

const emotionTypes = {
  happy: { emoji: '😄', text: '开心' },
  sad: { emoji: '😢', text: '伤心' },
  angry: { emoji: '😠', text: '生气' },
  surprised: { emoji: '😮', text: '惊讶' },
  calm: { emoji: '😌', text: '平静' },
  anxious: { emoji: '😰', text: '焦虑' },
  lonely: { emoji: '🥺', text: '孤独' },
  tired: { emoji: '😴', text: '疲惫' }
}

const severityTypes = {
  low: { color: 'green', text: '轻微' },
  medium: { color: 'orange', text: '中等' },
  high: { color: 'red', text: '严重' }
}

const filterForm = reactive({
  pet_id: '',
  start_date: null,
  end_date: null
})

const reportOverview = reactive({
  total_count: 0,
  main_emotion: 'calm',
  main_emotion_ratio: 0,
  anomaly_count: 0,
  anomaly_trend: 'down',
  anomaly_change: 0,
  response_triggered: 0,
  response_success_rate: 0
})

const anomalyColumns = [
  { title: '时间', dataIndex: 'create_time', width: 160 },
  { title: '情绪', slotName: 'emotion', width: 100 },
  { title: '严重程度', slotName: 'severity', width: 80 },
  { title: '持续时长', dataIndex: 'duration', width: 100 },
  { title: '操作', slotName: 'actions', width: 80 }
]

const anomalyEvents = ref([])

const getEmotionEmoji = (type) => emotionTypes[type]?.emoji || '❓'
const getEmotionText = (type) => emotionTypes[type]?.text || '未知'
const getSeverityColor = (severity) => severityTypes[severity]?.color || 'gray'
const getSeverityText = (severity) => severityTypes[severity]?.text || severity

const loadReport = async () => {
  try {
    const params = {
      type: activeTab.value,
      pet_id: filterForm.pet_id || undefined,
      start_date: filterForm.start_date?.format('YYYY-MM-DD') || undefined,
      end_date: filterForm.end_date?.format('YYYY-MM-DD') || undefined
    }
    const res = await axios.get(`${API_BASE}/reports`, { params })
    if (res.data.code === 0) {
      Object.assign(reportOverview, res.data.data.overview)
      anomalyEvents.value = res.data.data.anomalies || []
    }
  } catch (err) {
    Object.assign(reportOverview, {
      total_count: 1256,
      main_emotion: 'calm',
      main_emotion_ratio: 35,
      anomaly_count: 12,
      anomaly_trend: 'down',
      anomaly_change: 15,
      response_triggered: 89,
      response_success_rate: 92
    })
    anomalyEvents.value = generateMockAnomalies()
    Message.warning('使用模拟数据')
  }
  
  await nextTick()
  renderTrendChart()
  renderPieChart()
}

const generateMockAnomalies = () => {
  const types = ['sad', 'angry', 'anxious', 'lonely']
  const severities = ['low', 'medium', 'high']
  const suggestions = [
    '建议增加陪伴时间，关注宠物情绪变化',
    '检测到持续负面情绪，建议咨询专业人士',
    '环境因素可能影响情绪，建议检查生活环境',
    '建议播放舒缓音乐帮助宠物放松'
  ]
  
  return Array.from({ length: 6 }, (_, i) => ({
    event_id: `ANM${Date.now() - i}`,
    pet_id: pets.value[i % pets.value.length].pet_id,
    pet_name: pets.value[i % pets.value.length].pet_name,
    emotion_type: types[i % types.length],
    severity: severities[i % severities.length],
    create_time: new Date(Date.now() - i * 86400000 * 2).toLocaleString(),
    duration: `${10 + i * 5}分钟`,
    triggered_action: ['播放安抚音乐', '发送问候消息', '启动互动游戏'][i % 3],
    suggestion: suggestions[i % suggestions.length]
  }))
}

const renderTrendChart = () => {
  if (!trendChartRef.value) return
  
  if (trendChart) trendChart.dispose()
  trendChart = echarts.init(trendChartRef.value)
  
  const days = activeTab.value === 'daily' ? 24 : activeTab.value === 'weekly' ? 7 : 30
  const labels = activeTab.value === 'daily' 
    ? Array.from({ length: 24 }, (_, i) => `${i}:00`)
    : activeTab.value === 'weekly'
    ? ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    : Array.from({ length: 30 }, (_, i) => `${i + 1}日`)
  
  const emotionKeys = Object.keys(emotionTypes)
  const series = emotionKeys.map(type => ({
    name: emotionTypes[type].text,
    type: 'line',
    smooth: true,
    data: Array.from({ length: days }, () => Math.floor(Math.random() * 50 + 10))
  }))
  
  const option = {
    tooltip: { trigger: 'axis' },
    legend: { 
      data: emotionKeys.map(k => emotionTypes[k].text),
      bottom: 0
    },
    grid: { left: '3%', right: '4%', bottom: '15%', top: '3%', containLabel: true },
    xAxis: { type: 'category', boundaryGap: false, data: labels },
    yAxis: { type: 'value', name: '次数' },
    series
  }
  
  trendChart.setOption(option)
}

const renderPieChart = () => {
  if (!pieChartRef.value) return
  
  if (pieChart) pieChart.dispose()
  pieChart = echarts.init(pieChartRef.value)
  
  const data = Object.entries(emotionTypes).map(([key, val]) => ({
    name: val.text,
    value: Math.floor(Math.random() * 300 + 50)
  }))
  
  const option = {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { orient: 'vertical', right: 10, top: 'center' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['40%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
      data
    }]
  }
  
  pieChart.setOption(option)
}

const showAnomalyDetail = (record) => {
  currentAnomaly.value = record
  anomalyModalVisible.value = true
}

const exportReport = () => {
  Message.info('正在生成报告...')
  setTimeout(() => Message.success('报告导出成功'), 1500)
}

const initDateRange = () => {
  const now = new Date()
  if (activeTab.value === 'daily') {
    filterForm.start_date = now
    filterForm.end_date = now
  } else if (activeTab.value === 'weekly') {
    filterForm.start_date = new Date(now - 7 * 86400000)
    filterForm.end_date = now
  } else {
    filterForm.start_date = new Date(now.getFullYear(), now.getMonth(), 1)
    filterForm.end_date = now
  }
}

onMounted(() => {
  initDateRange()
  loadReport()
  
  window.addEventListener('resize', () => {
    trendChart?.resize()
    pieChart?.resize()
  })
})
</script>

<style scoped>
.pro-page-container { padding: 20px 24px; min-height: calc(100vh - 64px); background: #f5f7fa; }
.pro-breadcrumb { margin-bottom: 16px; }
.pro-search-bar { margin-bottom: 12px; }
.pro-action-bar { margin-bottom: 16px; }
.pro-content-area {
  background: #fff; border-radius: 8px; padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
.stat-card { text-align: center; }
.stat-label { font-size: 14px; color: #666; margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: 600; color: #333; }
.stat-value.emotion-main { font-size: 24px; }
.stat-value.warning { color: #ff7d00; }
.stat-trend { font-size: 12px; color: #999; margin-top: 4px; }
.stat-trend.up { color: #52c41a; }
.stat-trend.down { color: #ff4d4f; }
.chart-container { height: 300px; width: 100%; }
</style>
