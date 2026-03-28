<template>
  <div class="container">
    <a-card>
      <template #title><icon-send /> 设备指令</template>
      <a-form :model="formData" layout="vertical" style="max-width: 500px">
        <a-form-item label="设备"><a-select v-model="formData.deviceId"><a-option value="D001">设备-001</a-option></a-select></a-form-item>
        <a-form-item label="指令类型"><a-select v-model="formData.commandType"><a-option value="restart">重启</a-option><a-option value="upgrade">升级</a-option></a-select></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSend">发送指令</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const formData = reactive({ deviceId: '', commandType: '' })
const handleSend = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
