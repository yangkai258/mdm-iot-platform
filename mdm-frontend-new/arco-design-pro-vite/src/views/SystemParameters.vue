<template>
  <div class="container">
    <a-card>
      <template #title><icon-settings /> 系统参数</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="会话超时"><a-input-number v-model="form.sessionTimeout" /> 秒</a-form-item>
        <a-form-item label="最大登录数"><a-input-number v-model="form.maxLogin" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ sessionTimeout: 3600, maxLogin: 5 })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
