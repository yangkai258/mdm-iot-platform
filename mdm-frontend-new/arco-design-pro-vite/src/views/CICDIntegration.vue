<template>
  <div class="container">
    <a-card>
      <template #title><icon-sync /> CI/CD集成</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="平台"><a-select v-model="form.platform"><a-option value="github">GitHub</a-option></a-select></a-form-item>
        <a-form-item label="仓库"><a-input v-model="form.repo" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleConnect">连接</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ platform: '', repo: '' })
const handleConnect = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
