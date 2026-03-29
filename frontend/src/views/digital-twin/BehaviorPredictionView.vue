<template>
  <div class="page-container">
    <!-- 搜索筛选区 -->
    <div class="search-bar">
      <a-space>
        <a-select
          v-model="selectedPetId"
          placeholder="选择宠物"
          style="width: 200px"
          allow-search
          @change="loadPrediction"
        >
          <a-option v-for="pet in petList" :key="pet.device_id" :value="pet.device_id">
            {{ pet.pet_name }} ({{ pet.device_id }})
          </a-option>
        </a-select>
        <a-button type="primary" @click="loadPrediction">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="action-bar">
      <a-space>
        <a-button @click="goBack">返回仪表盘</a-button>
      </a-space>
    </div>

    <!-- 当前状态展示 -->
    <div class="current-state-area">
      <a-card title="当前状态" class="state-card">
        <div class="state-display">
          <div class="state-icon" :class="prediction?.current_state?.state">
            <component :is="getStateIcon(prediction?.current_state?.state)" :size="48" />
          </div>
          <div class="state-info">
            <div class="state-name">{{ getStateName(prediction?.current_state?.state) }}</div>
            <div class="state-confidence">
              置信度: <a-progress
                :percent="Math.round((prediction?.current_state?.confidence || 0) * 100)"
                :stroke-width="6"
                :show-text="false"
                :color="getConfidenceColor(prediction?.current_state?.confidence)"
                style="width: 100px; display: inline-block; vertical-align: middle;"
              /> {{ Math.round((prediction?.current_state?.confidence || 0) * 100) }}%
            </div>
            <div class="state-time">
              开始时间: {{ formatTime(prediction?.current_state?.started_at) }}
            </div>
          </div>
        </div>
      </a-card>
    </div>

    <!-- 短期预测结果 -->
    <div class="prediction-area">
      <a-row :gutter="16">
        <!-- 短期预测 -->
        <a-col :span="16">
          <a-card title="短期行为预测" class="prediction-card">
            <div v-if="loading" class="loading-state">
              <a-spin size="large" />
            </div>
            <a-table
              v-else
              :columns="predictionColumns"
              :data="prediction?.short_term_predictions || []"
              :pagination="false"
              row-key="behavior"
            >
              <template #behavior="{ record }">
                <a-tag :color="getPredictionColor(record.probability)">
                  {{ getBehaviorName(record.behavior) }}
                </a-tag>
              </template>
              <template #probability="{ record }">
                <div class="probability-cell">
                  <a-progress
                    :percent="Math.round(record.probability * 100)"
                    :stroke-width="8"
                    :show-text="false"
                    :color="getPredictionColor(record.probability)"
                    style="width: 100px"
                  />
                  <span class="probability-text">{{ Math.round(record.probability * 100) }}%</span>
                </div>
              </template>
              <template #expected_time="{ record }">
                {{ formatTime(record.expected_time) }}
              </template>
              <template #duration_estimate="{ record }">
                {{ formatDuration(record.duration_estimate) }}
              </template>
            </a-table>
          </a-card>
        </a-col>

        <!-- 意图识别 -->
        <a-col :span="8">
          <a-card title="意图识别" class="intent-card">
            <div v-if="prediction?.intent_recognition" class="intent-display">
              <div class="intent-icon" :class="prediction.intent_recognition.intent">
                <component :is="getIntentIcon(prediction.intent_recognition.intent)" :size="36" />
              </div>
              <div class="intent-name">{{ getIntentName(prediction.intent_recognition.intent) }}</div>
              <div class="intent-confidence">
                置信度: {{ Math.round(prediction.intent_recognition.confidence * 100) }}%
              </div>
              <a-divider />
              <div class="suggested-action">
                <div class="action-label">建议动作</div>
                <div class="action-content">{{ prediction.intent_recognition.suggested_action }}</div>
              </div>
            </div>
            <a-empty v-else description="暂无意图数据" />
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 预测详情图表 -->
    <div class="charts-area">
      <a-card title="行为概率分布" class="chart-card">
        <div ref="chartRef" class="chart-container"></div>
      </a-card>
    </div>

    <!-- 详细信息 -->
    <div class="content-area">
      <a-card title="详细信息">
        <a-descriptions :column="2" bordered>
          <a-descriptions-item label="设备ID">
            {{ prediction?.device_id || '--' }}
          </a-descriptions-item>
          <a-descriptions-item label="更新时间">
            {{ formatTime(prediction?.updated_at) }}
          </a-descriptions-item>
          <a-descriptions-item label="当前状态">
            {{ getStateName(prediction?.current_state?.state) }}
          </a-descriptions-item>
          <a-descriptions-item label="状态置信度">
            {{ Math.round((prediction?.current_state?.confidence || 0) * 100) }}%
          </a-descriptions-item>
          <a-descriptions-item label="识别意图">
            {{ getIntentName(prediction?.intent_recognition?.intent) }}
          </a-descriptions-item>
          <a-descriptions-item label="意图置信度">
            {{ Math.round((prediction?.intent_recognition?.confidence || 0) * 100) }}%
          </a-descriptions-item>
        </a-descriptions>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'
import * as echarts from 'echarts'

const router = useRouter()
const route = useRoute()

const API_BASE = '/api/v1/digital-twin'

const petList = ref([])
const selectedPetId = ref('')
const loading = ref(false)
const prediction = ref(null)
const chartRef = ref(null)
let predictionChart = null

const predictionColumns = [
  { title: '行为', slotName: 'behavior', width: 120 },
  { title: '概率', slotName: 'probability', width: 200 },
  { title: '预计时间', slotName: 'expected_time', width: 180 },
  { title: '预估时长', slotName: 'duration_estimate', width: 120 }
]

// 获取宠物列表
const loadPets = async () => {
  try {
    const res = await axios.get(`${API_BASE}/pets`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      petList.value = res.data.data?.list || res.data.data || []
    }
  } catch {
    petList.value = [
      { device_id: 'PET001', pet_name: '小白' },
      { device_id: 'PET002', pet_name: '旺财' }
    ]
  }

  if (route.query.device) {
    selectedPetId.value = route.query.device
  } else if (petList.value.length > 0) {
    selectedPetId.value = petList.value[0].device_id
  }
  if (selectedPetId.value) {
    loadPrediction()
  }
}

// 加载预测数据
const loadPrediction = async () => {
  if (!selectedPetId.value) {
    Message.warning('请先选择宠物')
    return
  }
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/behavior/prediction/${selectedPetId.value}`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.data.code === 0 || res.data.code === 200) {
      prediction.value = res.data.data
      renderChart()
    }
  } catch {
    // 模拟数据
    prediction.value = generateMockPrediction()
    renderChart()
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

// 生成模拟预测数据
const generateMockPrediction = () => {
  const states = ['awake', 'sleeping', 'eating', 'playing']
  const intents = ['hunger', 'thirst', 'play', 'rest', 'seek_attention']
  const behaviors = ['sleeping', 'eating', 'playing', 'walking', 'drinking']

  const currentState = states[Math.floor(Math.random() * states.length)]
  const intent = intents[Math.floor(Math.random() * intents.length)]

  const shortTerm = behaviors.map(b => ({
    behavior: b,
    probability: Math.random() * 0.5 + 0.1,
    expected_time: new Date(Date.now() + Math.random() * 3600000).toISOString(),
    duration_estimate: Math.floor(Math.random() * 3600) + 600
  })).sort((a, b) => b.probability - a.probability)

  return {
    device_id: selectedPetId.value,
    current_state: {
      state: currentState,
      confidence: Math.random() * 0.3 + 0.7,
      started_at: new Date(Date.now() - Math.random() * 7200000).toISOString()
    },
    short_term_predictions: shortTerm,
    intent_recognition: {
      intent,
      confidence: Math.random() * 0.4 + 0.5,
      suggested_action: getSuggestedAction(intent)
    },
    updated_at: new Date().toISOString()
  }
}

// 获取建议动作
const getSuggestedAction = (intent) => {
  const actions = {
    hunger: '建议投喂食物',
    thirst: '建议补充水分',
    play: '建议陪宠物玩耍',
    rest: '建议让宠物休息',
    seek_attention: '建议给予关注和抚摸'
  }
  return actions[intent] || '观察宠物状态'
}

// 获取状态名称
const getStateName = (state) => {
  const names = {
    awake: '清醒',
    sleeping: '睡眠中',
    eating: '进食中',
    playing: '玩耍中'
  }
  return names[state] || state || '--'
}

// 获取状态图标
const getStateIcon = (state) => {
  const icons = {
    awake: 'IconBulb',
    sleeping: 'IconMoonFill',
    eating: 'IconHeartFill',
    playing: 'IconLiveBroadcast'
  }
  return icons[state] || 'IconBulb'
}

// 获取意图名称
const getIntentName = (intent) => {
  const names = {
    hunger: '饥饿',
    thirst: '口渴',
    play: '想要玩耍',
    rest: '需要休息',
    seek_attention: '寻求关注'
  }
  return names[intent] || intent || '--'
}

// 获取意图图标
const getIntentIcon = (intent) => {
  const icons = {
    hunger: 'IconHeartFill',
    thirst: 'IconWaterFill',
    play: 'IconLiveBroadcast',
    rest: 'IconMoonFill',
    seek_attention: 'IconUser'
  }
  return icons[intent] || 'IconBulb'
}

// 获取行为名称
const getBehaviorName = (behavior) => {
  const names = {
    sleeping: '睡眠',
    eating: '进食',
    playing: '玩耍',
    walking: '散步',
    drinking: '饮水'
  }
  return names[behavior] || behavior
}

// 获取置信度颜色
const getConfidenceColor = (confidence) => {
  if (confidence == null) return '#1650ff'
  if (confidence >= 0.8) return '#00b42a'
  if (confidence >= 0.5) return '#1650ff'
  return '#ff7d00'
}

// 获取预测概率颜色
const getPredictionColor = (probability) => {
  if (probability >= 0.7) return '#00b42a'
  if (probability >= 0.4) return '#1650ff'
  return '#86909c'
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return '--'
  const d = new Date(timeStr)
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}

// 格式化时长
const formatDuration = (seconds) => {
  if (!seconds) return '-'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  if (h > 0) return `${h}小时${m}分钟`
  return `${m}分钟`
}

// 渲染图表
const renderChart = async () => {
  if (!prediction.value?.short_term_predictions) return
  await nextTick()

  if (!chartRef.value) return
  if (!predictionChart) predictionChart = echarts.init(chartRef.value)

  const data = prediction.value.short_term_predictions.map(p => ({
    name: getBehaviorName(p.behavior),
    value: Math.round(p.probability * 100)
  }))

  predictionChart.setOption({
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', right: 10, top: 'center' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data,
      label: { formatter: '{b}: {d}%' },
      itemStyle: {
        color: (params) => {
          const colors = ['#00b42a', '#1650ff', '#ff7d00', '#722ed1', '#f53f3f']
          return colors[params.dataIndex % colors.length]
        }
      }
    }],
    grid: { left: 10, right: 10, top: 10, bottom: 10 }
  })
}

const handleResize = () => {
  predictionChart?.resize()
}

const goBack = () => {
  router.push('/digital-twin/vitals')
}

onMounted(() => {
  loadPets()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  predictionChart?.dispose()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
}

.search-bar {
  margin-bottom: 12px;
}

.action-bar {
  margin-bottom: 16px;
}

.current-state-area {
  margin-bottom: 16px;
}

.state-card {
  border-radius: 8px;
}

.state-display {
  display: flex;
  align-items: center;
  gap: 24px;
}

.state-icon {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.state-icon.awake {
  background: linear-gradient(135deg, #1650ff 0%, #722ed1 100%);
}

.state-icon.sleeping {
  background: linear-gradient(135deg, #722ed1 0%, #0c0c22 100%);
}

.state-icon.eating {
  background: linear-gradient(135deg, #ff7d00 0%, #f53f3f 100%);
}

.state-icon.playing {
  background: linear-gradient(135deg, #00b42a 0%, #1650ff 100%);
}

.state-info {
  flex: 1;
}

.state-name {
  font-size: 24px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 8px;
}

.state-confidence {
  font-size: 14px;
  color: #4e5969;
  margin-bottom: 4px;
}

.state-time {
  font-size: 13px;
  color: #86909c;
}

.prediction-area {
  margin-bottom: 16px;
}

.prediction-card,
.intent-card {
  border-radius: 8px;
  height: 100%;
}

.loading-state {
  padding: 40px;
  text-align: center;
}

.probability-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.probability-text {
  font-size: 13px;
  color: #1d2129;
}

.intent-display {
  text-align: center;
}

.intent-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-bottom: 12px;
}

.intent-icon.hunger {
  background: linear-gradient(135deg, #ff7d00 0%, #f53f3f 100%);
}

.intent-icon.thirst {
  background: linear-gradient(135deg, #1650ff 0%, #0c8ae6 100%);
}

.intent-icon.play {
  background: linear-gradient(135deg, #00b42a 0%, #00dd67 100%);
}

.intent-icon.rest {
  background: linear-gradient(135deg, #722ed1 0%, #0c0c22 100%);
}

.intent-icon.seek_attention {
  background: linear-gradient(135deg, #eb4a5d 0%, #ff7d00 100%);
}

.intent-name {
  font-size: 18px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 4px;
}

.intent-confidence {
  font-size: 13px;
  color: #86909c;
}

.suggested-action {
  text-align: left;
}

.action-label {
  font-size: 12px;
  color: #86909c;
  margin-bottom: 4px;
}

.action-content {
  font-size: 14px;
  color: #1d2129;
  background: #f2f3f5;
  padding: 8px 12px;
  border-radius: 4px;
}

.charts-area {
  margin-bottom: 16px;
}

.chart-card {
  border-radius: 8px;
}

.chart-container {
  width: 100%;
  height: 300px;
}

.content-area {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
</style>
