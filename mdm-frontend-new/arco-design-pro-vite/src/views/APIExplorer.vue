<template>
  <div class="container">
    <a-card>
      <template #title><icon-code /> APIExplorer</template>
      <a-form layout="vertical" style="max-width: 600px">
        <a-form-item label="Method"><a-select v-model="form.method"><a-option value="GET">GET</a-option><a-option value="POST">POST</a-option></a-select></a-form-item>
        <a-form-item label="URL"><a-input v-model="form.url" placeholder="/api/..." /></a-form-item>
        <a-form-item label="Body"><a-textarea v-model="form.body" :rows="4" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSend">发送</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ method: 'GET', url: '', body: '' })
const handleSend = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
