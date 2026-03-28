<template>
  <div class="container">
    <a-card>
      <template #title><icon-user /> 儿童模式</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="模式开关"><a-switch v-model="form.enabled" /></a-form-item>
        <a-form-item label="可用时长"><a-input-number v-model="form.hours" :min="1" /> 小时/天</a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ enabled: false, hours: 2 })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
