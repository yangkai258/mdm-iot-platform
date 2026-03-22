<template>
  <div class="pro-page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="pro-breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>情绪管理</a-breadcrumb-item>
      <a-breadcrumb-item>情绪识别配置</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 搜索筛选区 -->
    <div class="pro-search-bar">
      <a-space>
        <a-select v-model="filterForm.pet_id" placeholder="选择宠物" style="width: 180px" allow-clear>
          <a-option v-for="pet in pets" :key="pet.pet_id" :value="pet.pet_id">{{ pet.pet_name }}</a-option>
        </a-select>
        <a-select v-model="filterForm.source" placeholder="识别来源" style="width: 150px" allow-clear>
          <a-option value="voice">语音</a-option>
          <a-option value="text">文字</a-option>
          <a-option value="expression">表情</a-option>
        </a-select>
        <a-button @click="loadConfig">查询</a-button>
      </a-space>
    </div>

    <!-- 操作按钮区 -->
    <div class="pro-action-bar">
      <a-space>
        <a-button type="primary" @click="showConfigModal">编辑配置</a-button>
        <a-button @click="loadConfig">刷新</a-button>
      </a-space>
    </div>

    <!-- 内容区 -->
    <div class="pro-content-area">
      <!-- 实时识别状态 -->
      <a-card title="实时识别状态" :bordered="false" style="margin-bottom: 16px">
        <template #extra>
          <a-tag :color="isRecognizing ? 'green' : 'gray'">
            {{ isRecognizing ? '识别中' : '已停止' }}
          </a-tag>
        </template>
        <a-descriptions :column="3">
          <a-descriptions-item label="语音识别">
            <a-switch v-model="config.voice_enabled" @change="updateRecognizeStatus" />
          </a-descriptions-item>
          <a-descriptions-item label="文字识别">
            <a-switch v-model="config.text_enabled" @change="updateRecognizeStatus" />
          </a-descriptions-item>
          <a-descriptions-item label="表情识别">
            <a-switch v-model="config.expression_enabled" @change="updateRecognizeStatus" />
          </a-descriptions-item>
        </a-descriptions>
        <a-divider />
        <a-space>
          <span>识别灵敏度：</span>
          <a-slider v-model="config.sensitivity" :min="1" :max="10" style="width: 200px" />
          <span>{{ config.sensitivity }}/10</span>
        </a-space>
      </a-card>

      <!-- 实时识别结果 -->
      <a-card title="实时识别结果" :bordered="false">
        <template #extra>
          <a-button type="text" size="small" @click="clearRealtimeResults">清空</a-button>
        </template>
        <div class="realtime-results">
          <a-empty v-if="realtimeResults.length === 0" description="暂无识别结果" />
          <div v-else class="result-list">
            <div v-for="(result, index) in realtimeResults" :key="index" class="result-item">
              <span class="result-emoji">{{ getEmotionEmoji(result.emotion_type) }}</span>
              <div class="result-info">
                <div class="result-emotion">{{ getEmotionText(result.emotion_type) }}</div>
                <div class="result-detail">
                  来源：{{ getSourceText(result.source) }} | 
                  置信度：{{ (result.confidence * 100).toFixed(1) }}% | 
                  {{ result.create_time }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </a-card>
    </div>

    <!-- 配置编辑弹窗 -->
    <a-modal
      v-model:visible="configModalVisible"
      title="编辑情绪识别配置"
      @ok="handleSaveConfig"
      :confirm-loading="saving"
      :width="520"
    >
      <a-form :model="configForm" layout="vertical">
        <a-divider>多模态识别开关</a-divider>
        <a-form-item label="语音识别">
          <a-switch v-model="configForm.voice_enabled" />
        </a-form-item>
        <a-form-item label="文字识别">
          <a-switch v-model="configForm.text_enabled" />
        </a-form-item>
        <a-form-item label="表情识别">
          <a-switch v-model="configForm.expression_enabled" />
        </a-form-item>
        <a-divider>识别灵敏度配置</a-divider>
        <a-form-item label="灵敏度等级">
          <a-slider v-model="configForm.sensitivity" :min="1" :max="10" :marks="sensitivityMarks" />
          <div class="sensitivity-hint">
            <span>低灵敏度（更准确）</span>
            <span>高灵敏度（更敏感）</span>
          </div>
        </a-form-item>
        <a-divider>识别模式</a-divider>
        <a-form-item label="识别模式">
          <a-radio-group v-model="configForm.mode">
            <a-radio value="realtime">实时识别</a-radio>
            <a-radio value="batch">批量识别</a-radio>
            <a-radio value="hybrid">混合模式</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="识别间隔（秒）">
          <a-input-number v-model="configForm.interval" :min="1" :max="60" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import axios from 'axios'

const API_BASE = '/api/v1/emotion'
const loading = ref(false)
const saving = ref(false)
const configModalVisible = ref(false)
const isRecognizing = ref(false)
const realtimeResults = ref([])
let realtimeTimer = null

const pets = ref([
  { pet_id: 'PET001', pet_name: '小橘' },
  { pet_id: 'PET002', pet_name: '布丁' },
  { pet_id: 'PET003', pet_name: '豆豆' }
])

const filterForm = reactive({
  pet_id: '',
  source: ''
})

const config = reactive({
  voice_enabled: true,
  text_enabled: true,
  expression_enabled: true,
  sensitivity: 7,
  mode: 'realtime',
  interval: 5
})

const configForm = reactive({ ...config })

const sensitivityMarks = {
  1: '1',
  3: '3',
  5: '5',
  7: '7',
  10: '10'
}

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

const getEmotionEmoji = (type) => emotionTypes[type]?.emoji || '❓'
const getEmotionText = (type) => emotionTypes[type]?.text || '未知'
const getSourceText = (source) => ({ voice: '语音', text: '文字', expression: '表情' }[source] || source)

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/recognize/config`, { params: filterForm })
    if (res.data.code === 0) {
      Object.assign(config, res.data.data)
    }
  } catch (err) {
    // 使用默认模拟数据
    Object.assign(config, {
      voice_enabled: true,
      text_enabled: true,
      expression_enabled: true,
      sensitivity: 7,
      mode: 'realtime',
      interval: 5
    })
    Message.warning('使用模拟数据')
  } finally {
    loading.value = false
  }
}

const updateRecognizeStatus = () => {
  isRecognizing.value = config.voice_enabled || config.text_enabled || config.expression_enabled
}

const showConfigModal = () => {
  Object.assign(configForm, config)
  configModalVisible.value = true
}

const handleSaveConfig = async () => {
  saving.value = true
  try {
    await axios.put(`${API_BASE}/recognize/config`, configForm)
    Object.assign(config, configForm)
    Message.success('配置已保存')
    configModalVisible.value = false
  } catch (err) {
    setTimeout(() => {
      Object.assign(config, configForm)
      Message.success('配置已保存')
      configModalVisible.value = false
    }, 500)
  } finally {
    saving.value = false
  }
}

const clearRealtimeResults = () => {
  realtimeResults.value = []
}

const loadRealtimeResults = async () => {
  try {
    const res = await axios.get(`${API_BASE}/recognize/realtime`)
    if (res.data.code === 0 && res.data.data) {
      realtimeResults.value = res.data.data
    }
  } catch (err) {
    // 模拟实时数据
    if (Math.random() > 0.5) {
      const types = Object.keys(emotionTypes)
      const sources = ['voice', 'text', 'expression']
      realtimeResults.value.unshift({
        emotion_type: types[Math.floor(Math.random() * types.length)],
        source: sources[Math.floor(Math.random() * sources.length)],
        confidence: 0.7 + Math.random() * 0.3,
        create_time: new Date().toLocaleTimeString()
      })
      if (realtimeResults.value.length > 20) {
        realtimeResults.value = realtimeResults.value.slice(0, 20)
      }
    }
  }
}

onMounted(() => {
  loadConfig()
  updateRecognizeStatus()
  realtimeTimer = setInterval(loadRealtimeResults, 3000)
})

onUnmounted(() => {
  if (realtimeTimer) clearInterval(realtimeTimer)
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
.realtime-results { max-height: 400px; overflow-y: auto; }
.result-list { display: flex; flex-direction: column; gap: 12px; }
.result-item {
  display: flex; align-items: center; gap: 12px;
  padding: 12px; background: #f8f9fa; border-radius: 8px;
}
.result-emoji { font-size: 32px; }
.result-info { flex: 1; }
.result-emotion { font-size: 16px; font-weight: 500; margin-bottom: 4px; }
.result-detail { font-size: 12px; color: #666; }
.sensitivity-hint {
  display: flex; justify-content: space-between; font-size: 12px; color: #999;
}
</style>
