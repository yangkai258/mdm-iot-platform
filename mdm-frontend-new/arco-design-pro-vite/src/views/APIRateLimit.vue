<template>
  <div class="container">
    <a-card>
      <template #title><icon-limit /> API限流配置</template>
      <a-form :model="formData" layout="vertical" style="max-width: 500px">
        <a-form-item label="每分钟请求数"><a-input-number v-model="formData.rpm" :min="1" /></a-form-item>
        <a-form-item label="每天请求数"><a-input-number v-model="formData.rpd" :min="1" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const formData = reactive({ rpm: 100, rpd: 10000 })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
