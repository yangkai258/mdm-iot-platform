<template>
  <div class="container">
    <a-card>
      <template #title><icon-plus /> 设备注册</template>
      <a-form :model="formData" layout="vertical" style="max-width: 500px">
        <a-form-item label="设备ID"><a-input v-model="formData.deviceId" /></a-form-item>
        <a-form-item label="设备名称"><a-input v-model="formData.deviceName" /></a-form-item>
        <a-form-item label="设备型号"><a-select v-model="formData.model"><a-option value="M5Stack-1">M5Stack-1</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleRegister">注册</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const formData = reactive({ deviceId: '', deviceName: '', model: '' })
const handleRegister = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
