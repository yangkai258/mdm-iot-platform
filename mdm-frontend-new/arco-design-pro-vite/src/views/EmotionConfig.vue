<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-heart /> 情绪识别配置</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="识别开关">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="语音情绪识别">
                <a-switch v-model="formData.voiceEnabled" />
              </a-form-item>
              <a-form-item label="文字情绪识别">
                <a-switch v-model="formData.textEnabled" />
              </a-form-item>
              <a-form-item label="表情情绪识别">
                <a-switch v-model="formData.expressionEnabled" />
              </a-form-item>
            </a-form>
          </a-card>

          <a-card title="响应策略" style="margin-top: 16px">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="响应灵敏度">
                <a-slider v-model="formData.sensitivity" :min="0" :max="100" :show-input="true" />
              </a-form-item>
              <a-form-item label="情绪响应模式">
                <a-select v-model="formData.responseMode">
                  <a-option value="gentle">温和模式</a-option>
                  <a-option value="active">积极模式</a-option>
                  <a-option value="balanced">平衡模式</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="AI模型版本">
                <a-select v-model="formData.modelVersion">
                  <a-option value="v2.1">EmotionNet v2.1</a-option>
                  <a-option value="v2.0">EmotionNet v2.0</a-option>
                </a-select>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="实时预览">
            <a-space direction="vertical" fill>
              <a-card size="small">
                <template #title>当前情绪状态</template>
                <a-statistic title="情绪类型" value="happy" />
                <a-progress :percent="currentEmotion.confidence" :stroke-color="getEmotionColor(currentEmotion.type)" />
              </a-card>
              <a-card size="small">
                <template #title>最近识别记录</template>
                <a-timeline>
                  <a-timeline-item v-for="r in recentRecords" :key="r.time" :color="getEmotionColor(r.type)">
                    <p>{{ r.type }} ({{ r.confidence }}%)</p>
                    <span class="time">{{ r.time }}</span>
                  </a-timeline-item>
                </a-timeline>
              </a-card>
            </a-space>
          </a-card>
        </a-col>
      </a-row>

      <a-space style="margin-top: 16px">
        <a-button type="primary" @click="handleSave">保存配置</a-button>
        <a-button @click="handleReset">恢复默认</a-button>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({
  voiceEnabled: true, textEnabled: true, expressionEnabled: true,
  sensitivity: 70, responseMode: 'balanced', modelVersion: 'v2.1'
})

const currentEmotion = reactive({ type: 'happy', confidence: 85 })
const recentRecords = ref([
  { type: 'happy', confidence: 85, time: '10:30:00' },
  { type: 'neutral', confidence: 72, time: '10:25:00' },
  { type: 'sad', confidence: 45, time: '10:20:00' }
])

const getEmotionColor = (type) => ({
  happy: '#67C23A', sad: '#409EFF', angry: '#F56C6C', neutral: '#909399'
}[type] || 'gray')

const handleSave = () => { }
const handleReset = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.time { font-size: 12px; color: #909399; }
</style>
