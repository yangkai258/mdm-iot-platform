<template>
  <div class="container">
    <a-card>
      <template #title><icon-launch /> 压力测试</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="并发数"><a-input-number v-model="form.concurrent" /></a-form-item>
        <a-form-item label="持续时间(秒)"><a-input-number v-model="form.duration" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleStart">开始测试</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ concurrent: 100, duration: 60 })
const handleStart = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
