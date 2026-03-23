<template>
  <div class="emotion-response-config-view">
    <a-card title="情绪响应配置">
      <a-alert type="info" style="margin-bottom: 16px">
        配置不同情绪状态下的响应策略，包括语音回复、动作表情、灯光颜色等。
      </a-alert>

      <a-form :model="configs" layout="vertical">
        <a-form-item label="开心响应">
          <a-select v-model="configs.happy.actions" multiple placeholder="选择响应动作">
            <a-option value="dance">跳舞</a-option>
            <a-option value="sing">唱歌</a-option>
            <a-option value="wave">挥手</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="难过响应">
          <a-select v-model="configs.sad.actions" multiple placeholder="选择响应动作">
            <a-option value="comfort">安慰语音</a-option>
            <a-option value="hug">拥抱</a-option>
            <a-option value="song">播放舒缓音乐</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="生气响应">
          <a-select v-model="configs.angry.actions" multiple placeholder="选择响应动作">
            <a-option value="calm">播放平复音乐</a-option>
            <a-option value="distract">转移注意力</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="害怕响应">
          <a-select v-model="configs.fear.actions" multiple placeholder="选择响应动作">
            <a-option value="safe">安全确认语音</a-option>
            <a-option value="hide">躲藏动作</a-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="onSave">保存配置</a-button>
            <a-button @click="onReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { Message } from '@arco-design/web-vue'

const configs = reactive({
  happy: { actions: ['dance', 'sing'] },
  sad: { actions: ['comfort', 'hug'] },
  angry: { actions: ['calm'] },
  fear: { actions: ['safe'] }
})

async function onSave() {
  try {
    await fetch('/api/v1/emotion/response-config', {
      method: 'PUT',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(configs)
    })
    Message.success('保存成功')
  } catch (e) {
    Message.error('保存失败')
  }
}

function onReset() {
  configs.happy.actions = []
  configs.sad.actions = []
  configs.angry.actions = []
  configs.fear.actions = []
}
</script>
