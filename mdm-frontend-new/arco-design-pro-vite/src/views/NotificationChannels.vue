<template>
  <div class="container">
    <a-card>
      <template #title><icon-broadcast /> 通知渠道配置</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="邮件通知"><a-switch v-model="form.emailEnabled" /></a-form-item>
        <a-form-item label="短信通知"><a-switch v-model="form.smsEnabled" /></a-form-item>
        <a-form-item label="推送通知"><a-switch v-model="form.pushEnabled" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ emailEnabled: true, smsEnabled: true, pushEnabled: true })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
