<template>
  <div class="voice-emotion-container">
    <a-card>
      <template #title>
        <span>语音情绪识别</span>
      </template>

      <!-- 录音分析 -->
      <a-card class="record-card">
        <div class="record-area">
          <div class="waveform" :class="{ recording: isRecording }">
            <div v-for="i in 20" :key="i" class="wave-bar" :style="{ height: waveHeights[i-1] + 'px' }"></div>
          </div>
          <a-button v-if="!isRecording" type="primary" size="large" @click="startRecording">
            <template #icon><icon-mic /></template>
            开始录音
          </a-button>
          <a-button v-else type="danger" size="large" status="error" @click="stopRecording">
            <template #icon><icon-stop /></template>
            停止录音
          </a-button>
        </div>

        <!-- 分析结果 -->
        <div v-if="analysisResult" class="analysis-result">
          <a-divider>分析结果</a-divider>
          <a-descriptions :column="2">
            <a-descriptions-item label="情绪类型">
              <a-tag :color="getEmotionColor(analysisResult.emotion)">
                {{ getEmotionText(analysisResult.emotion) }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="强度">
              <a-progress :percent="analysisResult.intensity * 10" :show-text="false" />
              {{ analysisResult.intensity }}/10
            </a-descriptions-item>
            <a-descriptions-item label="置信度">
              {{ (analysisResult.confidence * 100).toFixed(1) }}%
            </a-descriptions-item>
            <a-descriptions-item label="转录文本">
              {{ analysisResult.transcript || '无' }}
            </a-descriptions-item>
          </a-descriptions>
        </div>
      </a-card>

      <!-- 历史记录 -->
      <a-tabs default-active-key="records">
        <a-tab key="records" title="情绪记录">
          <a-table :data="records" :loading="loading" :pagination="pagination">
            <a-table-column title="宠物" dataIndex="pet_id"></a-table-column>
            <a-table-column title="情绪类型" dataIndex="emotion_type">
              <template #cell="{ record }">
                <a-tag :color="getEmotionColor(record.emotion_type)">
                  {{ getEmotionText(record.emotion_type) }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column title="强度" dataIndex="intensity">
              <template #cell="{ record }">
                <a-progress :percent="record.intensity * 10" :show-text="false" size="small" />
              </template>
            </a-table-column>
            <a-table-column title="置信度" dataIndex="confidence">
              <template #cell="{ record }">
                {{ (record.confidence * 100).toFixed(1) }}%
              </a-template>
            </a-table-column>
            <a-table-column title="转录" dataIndex="transcript" :ellipsis="true"></a-table-column>
            <a-table-column title="时间" dataIndex="created_at"></a-table-column>
          </a-table>
        </a-tab>

        <a-tab key="trend" title="情绪趋势">
          <div class="trend-chart">
            <div v-for="(item, index) in trendData" :key="index" class="trend-bar">
              <div class="trend-label">{{ item.date }}</div>
              <div class="trend-bars">
                <div v-for="emotion in ['calm', 'happy', 'anxious', 'angry']" :key="emotion" 
                     class="trend-segment" 
                     :style="{ height: (item[emotion] || 0) * 3 + 'px', backgroundColor: getEmotionColor(emotion) }">
                </div>
              </div>
            </div>
          </div>
        </a-tab>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const isRecording = ref(false)
const waveHeights = ref(Array(20).fill(10))
const analysisResult = ref(null)
const records = ref([])

const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0
})

const trendData = ref([
  { date: '周一', calm: 5, happy: 3, anxious: 1, angry: 1 },
  { date: '周二', calm: 4, happy: 4, anxious: 2, angry: 0 },
  { date: '周三', calm: 6, happy: 2, anxious: 1, angry: 1 },
  { date: '周四', calm: 3, happy: 5, anxious: 1, angry: 1 },
  { date: '周五', calm: 4, happy: 4, anxious: 2, angry: 0 },
])

let waveInterval = null

const getEmotionColor = (emotion) => {
  const colors = { calm: 'green', happy: 'arcoblue', anxious: 'orange', angry: 'red', sad: 'purple' }
  return colors[emotion] || 'gray'
}

const getEmotionText = (emotion) => {
  const texts = { calm: '平静', happy: '开心', anxious: '焦虑', angry: '愤怒', sad: '悲伤' }
  return texts[emotion] || emotion
}

const startRecording = () => {
  isRecording.value = true
  waveInterval = setInterval(() => {
    waveHeights.value = waveHeights.value.map(() => Math.random() * 40 + 10)
  }, 100)
  Message.success('开始录音...')
}

const stopRecording = async () => {
  isRecording.value = false
  if (waveInterval) {
    clearInterval(waveInterval)
    waveInterval = null
  }
  
  try {
    const res = await fetch('/api/v1/voice-emotion/analyze', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ audio_url: '/temp/recording.mp3' })
    })
    const data = await res.json()
    if (data.code === 0) {
      analysisResult.value = data.data
      Message.success('分析完成')
      loadRecords()
    }
  } catch (e) {
    Message.error('分析失败')
  }
}

const loadRecords = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/v1/voice-emotion/records?user_id=1')
    const data = await res.json()
    if (data.code === 0) {
      records.value = data.data.list || []
      pagination.value.total = data.data.total || 0
    }
  } catch (e) {
    Message.error('加载记录失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRecords()
})
</script>

<style scoped>
.voice-emotion-container {
  padding: 16px;
}

.record-card {
  margin-bottom: 16px;
}

.record-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 20px;
}

.waveform {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  height: 60px;
}

.wave-bar {
  width: 8px;
  background: #165dff;
  border-radius: 4px;
  transition: height 0.1s;
}

.waveform.recording .wave-bar {
  background: #00b42a;
  animation: pulse 0.5s infinite alternate;
}

@keyframes pulse {
  from { opacity: 0.7; }
  to { opacity: 1; }
}

.analysis-result {
  margin-top: 20px;
}

.trend-chart {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  height: 200px;
  padding: 20px;
  background: #f7f8fa;
  border-radius: 8px;
}

.trend-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.trend-label {
  font-size: 12px;
  color: #666;
}

.trend-bars {
  display: flex;
  gap: 2px;
  align-items: flex-end;
}

.trend-segment {
  width: 16px;
  border-radius: 2px;
}
</style>
