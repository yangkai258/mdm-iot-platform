<template>
  <div class="container">
    <a-card>
      <template #title><icon-limit /> 仿真资源配额</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="CPU配额"><a-input-number v-model="form.cpuQuota" suffix="核" /></a-form-item>
        <a-form-item label="内存配额"><a-input-number v-model="form.memoryQuota" suffix="GB" /></a-form-item>
        <a-form-item label="GPU配额"><a-input-number v-model="form.gpuQuota" suffix="卡" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ cpuQuota: 8, memoryQuota: 32, gpuQuota: 1 })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
