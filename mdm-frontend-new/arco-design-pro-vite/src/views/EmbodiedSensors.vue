<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-screenshot /> 环境感知配置</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-card title="感知参数配置">
            <a-form :model="formData" layout="vertical">
              <a-form-item label="摄像头开关">
                <a-switch v-model="formData.cameraEnabled" />
              </a-form-item>
              <a-form-item label="障碍物检测">
                <a-switch v-model="formData.obstacleEnabled" />
              </a-form-item>
              <a-form-item label="障碍物检测灵敏度">
                <a-slider v-model="formData.obstacleSensitivity" :min="0" :max="100" :show-input="true" :disabled="!formData.obstacleEnabled" />
              </a-form-item>
              <a-form-item label="感知范围半径">
                <a-input-number v-model="formData.sensingRange" :min="0.5" :max="10" :precision="1" suffix="米" />
              </a-form-item>
              <a-form-item label="感知刷新频率">
                <a-select v-model="formData.refreshRate">
                  <a-option value="1hz">1Hz (每秒1次)</a-option>
                  <a-option value="5hz">5Hz (每秒5次)</a-option>
                  <a-option value="10hz">10Hz (每秒10次)</a-option>
                </a-select>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleSave">保存配置</a-button>
                <a-button @click="handleCalibrate">标定向导</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="12">
          <a-card title="实时感知预览">
            <div class="preview-area">
              <svg width="400" height="300" viewBox="0 0 400 300">
                <rect x="0" y="0" width="400" height="300" fill="#f0f0f0" />
                <circle cx="200" cy="150" r="20" fill="#409EFF" />
                <text x="200" y="155" text-anchor="middle" fill="white" font-size="12">设备</text>
                <circle v-for="obs in obstacles" :key="obs.id" :cx="obs.x" :cy="obs.y" r="15" fill="#F56C6C" opacity="0.7" />
              </svg>
            </div>
            <a-space style="margin-top: 16px">
              <a-tag color="blue">设备</a-tag>
              <a-tag color="red">障碍物</a-tag>
            </a-space>
          </a-card>

          <a-card title="感知异常告警" style="margin-top: 16px">
            <a-form :model="alertForm" layout="vertical">
              <a-form-item label="启用异常告警">
                <a-switch v-model="alertForm.enabled" />
              </a-form-item>
              <a-form-item label="告警阈值">
                <a-input-number v-model="alertForm.threshold" :min="1" :max="20" suffix="个障碍物" />
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const formData = reactive({
  cameraEnabled: true, obstacleEnabled: true, obstacleSensitivity: 70,
  sensingRange: 3, refreshRate: '5hz'
})
const alertForm = reactive({ enabled: true, threshold: 5 })
const obstacles = ref([
  { id: 1, x: 150, y: 100 },
  { id: 2, x: 280, y: 200 }
])

const handleSave = () => { }
const handleCalibrate = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.preview-area { background: #f0f0f0; border-radius: 8px; padding: 16px; display: flex; justify-content: center; }
</style>
