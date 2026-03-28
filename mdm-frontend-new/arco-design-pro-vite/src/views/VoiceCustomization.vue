<template>
  <div class="container">
    <a-card>
      <template #title><icon-music /> 语音定制</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="语音类型"><a-select v-model="form.voiceType"><a-option value="female">女声</a-option><a-option value="male">男声</a-option></a-select></a-form-item>
        <a-form-item label="语速"><a-slider v-model="form.speed" :min="0.5" :max="2" /></a-form-item>
        <a-form-item label="音调"><a-slider v-model="form.pitch" :min="0.5" :max="2" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handlePreview">预览</a-button><a-button @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ voiceType: 'female', speed: 1, pitch: 1 })
const handlePreview = () => {}
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
