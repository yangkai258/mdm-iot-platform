<template>
  <div class="container">
    <a-card>
      <template #title><icon-launch /> OTA灰度发布</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="灰度比例"><a-slider v-model="form.percentage" :min="1" :max="100" :show-input="true" /></a-form-item>
        <a-form-item label="目标固件"><a-select v-model="form.version"><a-option value="v2.1.0">v2.1.0</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleDeploy">开始发布</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ percentage: 10, version: '' })
const handleDeploy = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
