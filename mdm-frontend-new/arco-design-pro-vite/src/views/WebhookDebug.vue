<template>
  <div class="container">
    <a-card>
      <template #title><icon-bug /> Webhook调试</template>
      <a-form layout="vertical" style="max-width: 600px">
        <a-form-item label="Webhook URL"><a-input v-model="form.url" placeholder="https://..." /></a-form-item>
        <a-form-item label="Method"><a-select v-model="form.method"><a-option value="POST">POST</a-option><a-option value="GET">GET</a-option></a-select></a-form-item>
        <a-form-item label="Body"><a-textarea v-model="form.body" :rows="6" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSend">发送</a-button></a-form-item>
      </a-form>
      <a-divider>响应</a-divider>
      <a-card v-if="response">{{ response }}</a-card>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ url: '', method: 'POST', body: '' })
const response = ref('')
const handleSend = () => { response.value = '{"status":"ok"}' }
</script>
<style scoped>.container { padding: 16px; }</style>
