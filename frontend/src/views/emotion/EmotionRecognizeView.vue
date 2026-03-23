<template>
  <div class="emotion-recognize-view">
    <a-card title="情绪识别配置">
      <a-form :model="form" layout="vertical" @submit-success="onSave">
        <a-form-item label="识别模式">
          <a-select v-model="form.mode">
            <a-option value="audio">语音情绪识别</a-option>
            <a-option value="visual">视觉情绪识别</a-option>
            <a-option value="both">综合识别</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="识别灵敏度">
          <a-slider v-model="form.sensitivity" :min="1" :max="10" :marks="{1:'低', 5:'中', 10:'高'}" />
        </a-form-item>
        <a-form-item label="最小置信度阈值">
          <a-input-number v-model="form.threshold" :min="0" :max="100" suffix="%" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="saving">保存配置</a-button>
        </a-form-item>
      </a-form>

      <a-divider>实时测试</a-divider>
      <a-space style="margin-bottom: 16px">
        <a-button type="outline" @click="startTest">开始测试</a-button>
        <a-button type="outline" status="warning" @click="stopTest">停止</a-button>
        <a-tag v-if="currentEmotion">{{ currentEmotion }}</a-tag>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const saving = ref(false)
const currentEmotion = ref('')

const form = reactive({
  mode: 'both',
  sensitivity: 5,
  threshold: 70
})

async function onSave() {
  saving.value = true
  try {
    await new Promise(r => setTimeout(r, 500))
    Message.success('配置保存成功')
  } finally {
    saving.value = false
  }
}

let testTimer = null
function startTest() {
  const emotions = ['开心', '平静', '好奇', '兴奋', '困惑']
  testTimer = setInterval(() => {
    currentEmotion.value = emotions[Math.floor(Math.random() * emotions.length)]
  }, 2000)
}

function stopTest() {
  if (testTimer) clearInterval(testTimer)
  currentEmotion.value = ''
}
</script>
