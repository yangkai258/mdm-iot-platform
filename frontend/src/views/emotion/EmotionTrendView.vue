<template>
  <div class="emotion-trend">
    <a-card title="情绪趋势分析">
      <template #extra>
        <a-select v-model="period" style="width: 120px" @change="load">
          <a-option value="week">近一周</a-option>
          <a-option value="month">近一月</a-option>
        </a-select>
      </template>
      
      <a-spin :loading="loading">
        <a-row :gutter="16" style="margin-bottom: 24px">
          <a-col :span="6">
            <a-statistic title="平均情绪强度" :value="stats.avg_intensity" suffix="/ 10" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="主要情绪" :value="stats.dominant_emotion" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="情绪趋势" :value="stats.trend">
              <template #suffix>
                <icon-arrow-up v-if="stats.trend === 'improving'" color="green" />
                <icon-arrow-down v-else-if="stats.trend === 'declining'" color="red" />
              </template>
            </a-statistic>
          </a-col>
          <a-col :span="6">
            <a-statistic title="情绪记录数" :value="stats.total_records" />
          </a-col>
        </a-row>
        
        <a-divider>情绪分布</a-divider>
        <div class="emotion-distribution">
          <div v-for="(pct, emotion) in stats.emotion_distribution" :key="emotion" class="emotion-bar">
            <span class="label">{{ emotion }}</span>
            <a-progress :percent="pct" :color="getEmotionColor(emotion)" />
          </div>
        </div>
        
        <a-divider>周趋势</a-divider>
        <div class="weekly-chart">
          <a-table :columns="chartColumns" :data="stats.weekly_data" size="small" />
        </div>
      </a-spin>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const loading = ref(false)
const period = ref('week')
const stats = ref({
  avg_intensity: 0,
  dominant_emotion: '-',
  trend: '-',
  total_records: 0,
  emotion_distribution: {},
  weekly_data: []
})

const chartColumns = [
  { title: '日期', dataIndex: 'date' },
  { title: '平均强度', dataIndex: 'avg_intensity' },
  { title: '主要情绪', dataIndex: 'dominant_emotion' }
]

const getEmotionColor = (e) => ({
  happy: '#52c41a', sad: '#1890ff', angry: '#ff4d4f', 
  calm: '#722ed1', excited: '#faad14', anxious: '#f5222d'
}[e] || '#1890ff')

const load = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/emotions/records/stats?period=${period.value}&pet_id=1`)
    stats.value = await res.json()
  } catch (e) {
    Message.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.emotion-distribution { max-width: 500px; }
.emotion-bar { display: flex; align-items: center; gap: 12px; margin-bottom: 8px; }
.emotion-bar .label { width: 60px; font-weight: 500; }
.emotion-bar :deep(.arco-progress-line) { flex: 1; }
.weekly-chart { margin-top: 16px; }
</style>
