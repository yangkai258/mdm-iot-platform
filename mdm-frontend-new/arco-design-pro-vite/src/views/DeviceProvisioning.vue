<template>
  <div class="container">
    <a-card>
      <template #title><icon-setup /> 设备配置</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="WiFi SSID"><a-input v-model="form.wifiSsid" /></a-form-item>
        <a-form-item label="WiFi 密码"><a-input-password v-model="form.wifiPassword" /></a-form-item>
        <a-form-item label="MQTT服务器"><a-input v-model="form.mqttServer" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存配置</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ wifiSsid: '', wifiPassword: '', mqttServer: '' })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
