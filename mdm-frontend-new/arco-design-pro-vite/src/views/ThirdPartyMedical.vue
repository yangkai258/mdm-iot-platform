<template>
  <div class="container">
    <a-card>
      <template #title><icon-medicine-box /> 第三方医疗对接</template>
      <a-form layout="vertical" style="max-width: 500px">
        <a-form-item label="对接平台"><a-select v-model="form.platform"><a-option value="pet_medical">宠物医疗</a-option></a-select></a-form-item>
        <a-form-item label="API Key"><a-input-password v-model="form.apiKey" /></a-form-item>
        <a-form-item><a-button type="primary" @click="handleSave">保存</a-button></a-form-item>
      </a-form>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
const form = reactive({ platform: '', apiKey: '' })
const handleSave = () => {}
</script>
<style scoped>.container { padding: 16px; }</style>
