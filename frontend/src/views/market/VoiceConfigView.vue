<template>
  <div class="container">
    <Breadcrumb :items="['menu.market', 'menu.market.voiceConfig']" />
    <a-card class="general-card" title="语音配置">
      <template #extra>
        <a-button type="primary" @click="handleSave"><icon-save />保存</a-button>
      </template>
      <a-form :model="form" layout="vertical" style="max-width: 600px">
        <a-form-item label="默认语音"><a-input v-model="form.default_voice" /></a-form-item>
        <a-form-item label="语速"><a-input-number v-model="form.speed" :min="0.5" :max="2" :precision="1" style="width:100%" /></a-form-item>
        <a-form-item label="音调"><a-input-number v-model="form.pitch" :min="0.5" :max="2" :precision="1" style="width:100%" /></a-form-item>
        <a-form-item label="音量"><a-input-number v-model="form.volume" :min="0" :max="100" style="width:100%" /></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import Breadcrumb from '@/components/Breadcrumb.vue'

const form = reactive({ default_voice: '', speed: 1, pitch: 1, volume: 80 })
const handleSave = () => Message.success('保存成功')
</script>
